Just an app that pings regularly pings other IP adresses. Responds to pings on `/ping` with a `pong`.

### Usage
`pinger --listenPort <port> --interval <interval in seconds> --address <IP address 1>:<some port> [--address <IP address 2>:<some other port>]`

### Consumption
find a boshrelease at https://github.com/voelzmo/pinger-app-release or just deploy the pinger on Cloud Foundry e.g. with the binary buildpack.
