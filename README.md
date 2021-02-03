Just an app that pings regularly pings other IP adresses. Responds to pings on `/ping` with a `pong`.

### Usage
`pinger --port <port> --interval <interval in seconds> --address <IP address 1>:<some port> [--address <IP address 2>:<some other port>]`

### Consumption
* boshrelease: https://github.com/voelzmo/pinger-app-release
* deploy ping-app on Cloud Foundry with the golang buildpack
```
cf push ping-app
```
* deploy ping-app on Cloud Foundry with the binary buildpack
```
GOOS=linux GOARCH=amd64 go build -o ping-app
cf push ping-app -c './ping-app' -b https://github.com/cloudfoundry/binary-buildpack.git
```

