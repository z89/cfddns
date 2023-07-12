# cfddns

a ddns client for cloudflare written in golang

## Install
installation instructions for unix based systems

#### Download
```bash
wget https://github.com/z89/cfddns/releases/download/v0.1.0-alpha/cfddns
```
```bash
chmod +x cfddns
```
#### Compile 
```bash
git clone https://github.com/z89/cfddns.git
```
```bash
cd cfddns && go build
```
## Usage
### Description
cfddns works by using the cloudflare API to retrieve your public IPv4 address. It then compares this address to your targeted DNS records on cloudflare. If the public address differs from your DNS records, cfddns will update the DNS records with the new address. This allows you to setup a custom DDNS client for cloudflare domains. It is possible to update multiple DNS records, as long as they share the same comment. By default, cfddns will update every 24 hours. To disable the timer, and enable the http endpoint, use the `-timer` flag. This is often done so another service can trigger the updates manually, such as a router.

### Required Flags
* `-key` - the API key for your cloudflare zone. DO NOT use a global API key, this will not work. Scope the API key to the specific zone you are targeting. A guide can be found here: https://developers.cloudflare.com/fundamentals/api/get-started/create-token/
* `-target` - the domain name of your zone. For example, if your zone is _example.com_, you would enter _example.com_ as the target. Do not include any subdomains. This is used to find the zone id.
* `-comment` - the comment on the DNS records you want to update. For example, if you want to update the record _home.example.com_, create a comment on that record with any string you'd like. This will be used to identify the records so they can be updated.  

### Optional Flags
* `-timer` - The time interval between updates in minutes. Default is 24 hours. To disable, set the timer to 0. This will enable the http endpoint.

* `-addr` - the address to listen on. Default is 0.0.0.0, which listens on all interfaces.
* `-port` - the port to listen on. Default is port 3000.

### Example


```bash
./cfddns -key=<api-key> -target=<domain-name> -comment=<comment> -timer=<interval> -addr=<address> -port=<port>
```
## Limitations
* currently only supports ipv4
* updates DNS A/AAAA records only
* must be ran 24/7 on a dedicated machine for maximum uptime

## Contributing

pull requests are welcome to:

- revise the docs
- fix bugs or bad logic
- suggest/add features or improvements

## License

MIT