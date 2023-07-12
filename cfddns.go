package main

import (
	"bufio"
	"bytes"
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"regexp"
	"strings"
	"time"

	"flag"

	"github.com/cloudflare/cloudflare-go"
)

func updateDNS(w http.ResponseWriter, ctx context.Context, keyFlag string, targetFlag string, commentFlag string) {
	// fetch public address
	resp, err := http.Get("https://cloudflare.com/cdn-cgi/trace")

	if err != nil {
		fmt.Printf("failed to fetch public address")

		if w != nil {
			w.Write([]byte("failed to fetch public address \n"))
		}

		return
	}

	defer resp.Body.Close()

	scan := bufio.NewScanner(resp.Body)

	for scan.Scan() {
		if bytes.Contains(scan.Bytes(), []byte("ip=")) {
			// parse ip address
			address := strings.Split(scan.Text(), "=")[1]

			// make sure address is valid ipv4
			match, _ := regexp.MatchString("(25[0-5]|2[0-4][0-9]|1[0-9][0-9]|[1-9]?[0-9])\\.(25[0-5]|2[0-4][0-9]|1[0-9][0-9]|[1-9]?[0-9])\\.(25[0-5]|2[0-4][0-9]|1[0-9][0-9]|[1-9]?[0-9])\\.(25[0-5]|2[0-4][0-9]|1[0-9][0-9]|[1-9]?[0-9])", address)

			if !match {
				log.Fatalf("ip does not match ipv4 format")
				os.Exit(1)
			}

			// create cloudflare api client
			api, err := cloudflare.NewWithAPIToken(keyFlag)

			if err != nil {
				log.Fatal(err)
				os.Exit(1)
			}

			// find zone id
			zoneID, err := api.ZoneIDByName(targetFlag)

			if err != nil {
				log.Fatalln(err)
				return
			}

			// list dns records
			records, _, err := api.ListDNSRecords(ctx, cloudflare.ZoneIdentifier(zoneID), cloudflare.ListDNSRecordsParams{})
			if err != nil {
				fmt.Println(err)
				return
			}

			// find dns record with comment indentifier
			for _, record := range records {
				if record.Comment == commentFlag {
					// check if the dns record matches current public address
					if record.Content == address {
						fmt.Printf("no update required, %s: %s -> %s\n", record.Name, record.Content, address)

						if w != nil {
							w.Write([]byte("no update required, " + record.Name + ": " + record.Content + " -> " + address + "\n"))
						}
					} else {
						// update DNS record
						_, err := api.UpdateDNSRecord(ctx, cloudflare.ZoneIdentifier(zoneID), cloudflare.UpdateDNSRecordParams{
							Type:    record.Type,
							ID:      record.ID,
							Name:    record.Name,
							Content: address,
							Proxied: record.Proxied,
							Comment: record.Comment,
						})

						if err != nil {
							log.Fatalln("failed to update dns record")
							return
						}

						fmt.Printf("successfully updated %s: %s -> %s\n", record.Name, record.Content, address)

						if w != nil {
							w.Write([]byte("successfully updated " + record.Name + ": " + record.Content + " -> " + address + "\n"))
						}

					}
				}
			}

			return
		}
	}
}

func main() {
	ctx := context.Background()

	// fetch user defined flags
	keyFlag := flag.String("key", "null", "cloudflare api key")
	targetFlag := flag.String("target", "null", "domain name to find zone id")
	commentFlag := flag.String("comment", "null", "dns record comment used for targeting")
	timer := flag.Int("timer", 1440, "time interval between updates")
	addr := flag.String("addr", "0.0.0.0", "http server port")
	port := flag.String("port", "3000", "http server port")
	endpoint := flag.String("endpoint", "/cfddns/update", "http server endpoint")

	flag.Parse()

	// check if required flags are provided
	switch {
	case *keyFlag == "null":
		log.Fatalf("api key not provided")
		os.Exit(1)
	case *targetFlag == "null":
		log.Fatalf("target domain not provided")
		os.Exit(1)
	case *commentFlag == "null":
		log.Fatalf("dns record comment not provided")
		os.Exit(1)
	}

	if *timer < 1 {
		http.HandleFunc(*endpoint, func(w http.ResponseWriter, r *http.Request) {
			updateDNS(w, ctx, *keyFlag, *targetFlag, *commentFlag)
		})

		fmt.Println("http server starting on 0.0.0.0:" + *port)

		err := http.ListenAndServe(*addr+":"+*port, nil)

		if err != nil {
			log.Fatal(err)
			os.Exit(1)
		}

	} else {
		ticker := time.NewTicker(time.Duration(*timer) * time.Minute)

		for {
			select {
			case <-ticker.C:
				updateDNS(nil, ctx, *keyFlag, *targetFlag, *commentFlag)
			}
		}

	}

}
