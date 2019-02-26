package main

import (
	"io"
	"log"
	"net"
	"time"
)

//noinspection GoUnhandledErrorResult
func forwardTCP(conf forwardConf) {
	localAddr, err := net.ResolveTCPAddr("tcp", conf.LocalAddr)
	if err != nil {
		log.Fatalf("[TCP] Failed to resolve local addr: %v", err)
	}
	remoteAddr, err := net.ResolveTCPAddr("tcp", conf.RemoteAddr)
	if err != nil {
		log.Fatalf("[TCP] Failed to resolve remote addr: %v", err)
	}

	l, err := net.ListenTCP("tcp", localAddr)
	if err != nil {
		log.Fatalf("[TCP] Local error: %v", err)
	}
	log.Printf("[TCP] Listening on %v", l.Addr())

	go func() {
		for {
			c, err := l.AcceptTCP()
			if err != nil {
				log.Printf("[TCP] Local error: %v", err)
				continue
			}
			log.Printf("[TCP] New connection from %v to %v", c.RemoteAddr(), c.LocalAddr())

			go func() {
				defer c.Close()
				c.SetKeepAlive(true)

				rc, err := net.DialTCP("tcp", nil, remoteAddr)
				if err != nil {
					log.Printf("[TCP] Remote error: %v", err)
					return
				}
				defer rc.Close()
				rc.SetKeepAlive(true)

				go tcpRelay(c, rc)
				tcpRelay(rc, c)
			}()
		}
	}()
}

//noinspection GoUnhandledErrorResult
func tcpRelay(dst, src *net.TCPConn) error {
	_, err := io.Copy(dst, src)
	src.SetDeadline(time.Now())
	dst.SetDeadline(time.Now())
	return err
}
