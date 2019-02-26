package main

import (
	"net"
	"sync"
	"time"
)

type natState uint8

const (
	new natState = iota
	assured
	confirmed
)

type natMap struct {
	m             sync.Map
	timeout       time.Duration
	streamTimeout time.Duration
}

type natConn struct {
	nm    *natMap
	conn  *net.UDPConn
	state natState
}

func (nm *natMap) Get(clientAddr *net.UDPAddr) *natConn {
	v, ok := nm.m.Load(clientAddr.String())
	if !ok {
		return nil
	}
	return v.(*natConn)
}

func (nm *natMap) Set(clientAddr *net.UDPAddr, nc *natConn) {
	nm.m.Store(clientAddr.String(), nc)
}

func (nm *natMap) Remove(clientAddr *net.UDPAddr) {
	nm.m.Delete(clientAddr.String())
}

//noinspection GoUnhandledErrorResult
func (nc *natConn) UpdateDeadline() {
	ddl := time.Now()
	if nc.state != confirmed {
		ddl = ddl.Add(nc.nm.timeout)
	} else {
		ddl = ddl.Add(nc.nm.streamTimeout)
	}
	nc.conn.SetDeadline(ddl)
}

func (nc *natConn) SwapState(old, new natState) {
	if nc.state == old {
		nc.state = new
	}
}

//noinspection GoUnhandledErrorResult
func (nc *natConn) Close(clientAddr *net.UDPAddr) {
	nc.conn.Close()
	nc.nm.Remove(clientAddr)
}
