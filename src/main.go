package main

func main() {
	forwardTCP(&forwardConf{
		LocalAddr:  "127.0.0.1:8080",
		RemoteAddr: "httpbin.org:80",
	})
}
