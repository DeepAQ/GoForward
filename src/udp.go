package main

import (
	"log"
	"net"
	"time"
)

const udpBufSize = 64 * 1024

//noinspection GoUnhandledErrorResult
func forwardUDP(conf forwardConf) {
	localAddr, err := net.ResolveUDPAddr("udp", conf.LocalAddr)
	if err != nil {
		log.Fatalf("[UDP] Failed to resolve local addr: %v", err)
	}
	remoteAddr, err := net.ResolveUDPAddr("udp", conf.RemoteAddr)
	if err != nil {
		log.Fatalf("[UDP] Failed to resolve remote addr: %v", err)
	}

	l, err := net.ListenUDP("udp", localAddr)
	if err != nil {
		log.Fatalf("[UDP] Local error: %v", err)
	}
	log.Printf("[UDP] Listening on %v", l.LocalAddr())

	go func() {
		nm := &natMap{
			timeout:       10 * time.Second,
			streamTimeout: 3 * time.Minute,
		}
		if v, err := time.ParseDuration(conf.Options["timeout"]); err == nil {
			nm.timeout = v
		}
		if v, err := time.ParseDuration(conf.Options["streamTimeout"]); err == nil {
			nm.streamTimeout = v
		}
		buf := make([]byte, udpBufSize)

		for {
			n, clientAddr, err := l.ReadFromUDP(buf)
			if err != nil {
				if err, ok := err.(net.Error); ok && !err.Timeout() {
					log.Printf("[UDP] Local error: %v", err)
				}
				continue
			}

			nc := nm.Get(clientAddr)
			if nc == nil {
				rc, err := net.DialUDP("udp", nil, remoteAddr)
				if err != nil {
					log.Printf("[UDP] Remote error: %v", err)
					continue
				}

				log.Printf("[UDP] New connection from %v to %v", clientAddr, l.LocalAddr())
				nc = &natConn{
					nm:    nm,
					conn:  rc,
					state: new,
				}
				nm.Set(clientAddr, nc)
				go udpRelay(nc, l, clientAddr)
			} else {
				nc.SwapState(assured, confirmed)
			}

			nc.UpdateDeadline()
			_, err = nc.conn.Write(buf[:n])
			if err != nil {
				log.Printf("[UDP] Remote error: %v", err)
				nc.Close(clientAddr)
				continue
			}
		}
	}()
}

//noinspection GoUnhandledErrorResult
func udpRelay(nc *natConn, local *net.UDPConn, clientAddr *net.UDPAddr) error {
	buf := make([]byte, udpBufSize)

	for {
		nc.UpdateDeadline()
		n, err := nc.conn.Read(buf)
		if err != nil {
			if err, ok := err.(net.Error); ok && !err.Timeout() {
				log.Printf("[UDP] Remote error: %v", err)
			}
			nc.Close(clientAddr)
			return err
		}
		nc.SwapState(new, assured)

		nc.UpdateDeadline()
		_, err = local.WriteToUDP(buf[:n], clientAddr)
		if err != nil {
			log.Printf("[UDP] Local error: %v", err)
			return err
		}
	}
}
