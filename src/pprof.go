package main

import (
	"log"
	"net/http"
	_ "net/http/pprof"
)

func listenPProf(addr string) {
	log.Printf("[PProf] Listening on %s", addr)
	go func() {
		err := http.ListenAndServe(addr, nil)
		if err != nil {
			log.Fatalf("[PProf] Listen error: %v", err)
		}
	}()
}
