package main

import (
	"flag"
	"fmt"
	"os"
)

var remoteAddr string
var serverPort string

func init() {
	flag.StringVar(&remoteAddr, "c", "", "Provide host and port to which to connect, or leave this flag to start locally")
	flag.StringVar(&serverPort, "p", "", "Provide server port for other player to connect to")
}

func printUsage() {
	fmt.Println("Usage: ttt -h 127.0.0.1:9001 to connect to a server")
	fmt.Println("ttt -p 9001 to start a server listening on 9001 port")
}

func main() {
	flag.Parse()
	if remoteAddr == "" {
		// try to run a server and connect to it locally as a client
		if serverPort == "" {
			fmt.Println("Server port missing")
			printUsage()
			os.Exit(1)
		} else {
			runAsServer(serverPort)
		}
	} else {
		// try to connect to a remote server
		if serverPort != "" {
			fmt.Println("Cannot use both server port and host to connect to: please choose your role")
			printUsage()
			os.Exit(1)
		} else {
			runAsClient(remoteAddr)
		}
	}
}

func runAsServer(port string) {
	exit := make(chan bool, 1)
	startServer(port)
	go startClient("127.0.0.1:" + port)
	<-exit
}

func runAsClient(addr string) {
	exit := make(chan bool, 1)
	go startClient(addr)
	<-exit
}
