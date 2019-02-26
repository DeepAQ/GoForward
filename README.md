## GoForward

Simple TCP/UDP flow forwarder implemented in Go.

### Usage

```
GoForward URL [URL...]
```

### URL Format

```
protocol://remoteAddress:port@[localAddress]:port[/?parameters]
```

- protocol: `tcp` or `udp`
- remoteAddress / localAddress: IPv4 or IPv6 or hostname
- parameters:
    - `timeout` (UDP only): I/O timeout for new connections, default `10s`
    - `streamTimeout` (UDP only): I/O timeout for assured connections, default `3m`

### Example

```
GoForward \
  "tcp://httpbin.org:80@:8080" \
  "udp://1.1.1.1:53@:53/?timeout=5s&streamTimeout=5m"
```
