package main

import (
	"io"
	"log"
	"net"
	"time"
)

//noinspection GoUnhandledErrorResult
func forwardTCP(conf *forwardConf) {
	l, err := net.Listen("tcp", conf.LocalAddr)
	if err != nil {
		log.Fatalf("TCP: Failed to listen on %s: %v", conf.LocalAddr, err)
	}
	log.Printf("TCP: Listening on %s", conf.LocalAddr)

	for {
		c, err := l.Accept()
		if err != nil {
			log.Printf("TCP: Failed to accept on %s: %v", conf.LocalAddr, err)
			continue
		}
		log.Printf("TCP: New connection from %v to %v", c.RemoteAddr(), c.LocalAddr())

		go func() {
			defer c.Close()
			_ = c.(*net.TCPConn).SetKeepAlive(true)

			rc, err := net.Dial("tcp", conf.RemoteAddr)
			if err != nil {
				log.Printf("TCP: Failed to dial to %s: %v", conf.RemoteAddr, err)
				return
			}
			defer rc.Close()
			_ = rc.(*net.TCPConn).SetKeepAlive(true)

			go tcpRelay(c, rc)
			tcpRelay(rc, c)
		}()
	}
}

func tcpRelay(dst, src net.Conn) error {
	_, err := io.Copy(dst, src)
	_ = src.SetDeadline(time.Now())
	_ = dst.SetDeadline(time.Now())
	return err
}
