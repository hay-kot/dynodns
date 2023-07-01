# dynodns

DynoDNS is an Dynamic DNS client build for the PorkBun API. It currently only supports porkbun, but we can support other providers on request or PR. 

The intended usecase for DynoDNS is to run on your infrastracture as a container or background process and on some specified interval

1. Check your public IP address
2. If different change your IP or the subdomain in PorkBun


## Installation

```bash
go install https://github.com/hay-kot/dynodns
```

## Usage


To use this application you first need to obtain your API credentials from porkbun. You can see how to do that [here](https://kb.porkbun.com/article/190-getting-started-with-the-porkbun-api)

### Docker

```yaml
---
version: "3.7"
services:
  dynodns:
    image: ghcr.io/hay-kot/dynodns:v0.1.0
    container_name: dynodns
    environment:
      - INTERVAL=300 # Time In Seconds
      - LOG_LEVEL=info
      - PORKBUN_DOMAIN=example.com
      - PORKBUN_SUBDOMAIN=dns 
      - PORKBUN_API_KEY=abc123_key
      - PORKBUN_API_SECRET=abc123_secret
      - PING_URL=https://up.example.com/ping/dynodns # health check URL (Up Time Kuma)
    restart: unless-stopped
```


### CLI 

DynoDNS can also be used directly as a CLI. 

```sh
NAME:
   dynodns - client for setting up dynamic DNS

USAGE:
   dynodns [global options] command [command options] [arguments...]

VERSION:
   dev (HEAD) now

COMMANDS:
   run      runs the client
   test-ip  test external IP finder
   help, h  Shows a list of commands or help for one command

GLOBAL OPTIONS:
   --help, -h     show help
   --version, -v  print the version

   Options

   --interval value   interval in seconds to check for ip changes (default: 300) [$INTERVAL]
   --log-level value  log level (debug, info, warn, error, fatal, panic) (default: "info") [$LOG_LEVEL]
   --ping-url value   Healthcheck Ping URL [$PING_URL]

   Porkbun

   --porkbun.domain value     porkbun domain to update [$PORKBUN_DOMAIN]
   --porkbun.endpoint value   porkbun api endpoint (default: "https://porkbun.com/api/json/v3") [$PORKBUN_API_ENDPOINT]
   --porkbun.key value        porkbun api key [$PORKBUN_API_KEY]
   --porkbun.secret value     porkbun api secret [$PORKBUN_API_SECRET]
   --porkbun.subdomain value  porkbun subdomain to update (default: "dns") [$PORKBUN_SUBDOMAIN]
```

