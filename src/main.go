package main

import (
	"log"
	"net/url"
	"os"
	"strings"
)

func main() {
	for _, arg := range os.Args[1:] {
		arg = strings.TrimSpace(arg)
		if len(arg) > 0 {
			forwardURL(arg)
		}
	}
	select {}
}

func forwardURL(u string) {
	p, err := url.Parse(u)
	if err != nil {
		log.Fatalf("Failed to parse URL: %v", err)
	}
	switch p.Scheme {
	case "tcp":
		forwardTCP(forwardConf{
			LocalAddr:  p.Host,
			RemoteAddr: p.User.String(),
		})
	case "udp":
		conf := forwardConf{
			LocalAddr:  p.Host,
			RemoteAddr: p.User.String(),
			Options:    map[string]string{},
		}
		for k, v := range p.Query() {
			conf.Options[k] = v[0]
		}
		forwardUDP(conf)
	case "pprof":
		listenPProf(p.Host)
	default:
		log.Fatalf("Unsupported protocol: %s", p.Scheme)
	}
}
