package main

func main() {
	forwardTCP(forwardConf{
		LocalAddr:  "127.0.0.1:8080",
		RemoteAddr: "httpbin.org:80",
	})
	forwardTCP(forwardConf{
		LocalAddr:  "127.0.0.1:8081",
		RemoteAddr: "httpbin.org:80",
	})
	forwardUDP(forwardConf{
		LocalAddr:  "127.0.0.1:53",
		RemoteAddr: "119.29.29.29:53",
	})
	forwardUDP(forwardConf{
		LocalAddr:  "127.0.0.1:5300",
		RemoteAddr: "119.29.29.29:53",
	})

	select {}
}
