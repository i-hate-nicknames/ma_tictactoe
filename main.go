package main

import (
	"flag"
	"fmt"
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
	testGame()
}

func testGame() {
	game := MakeBoard(3)
	err := game.MakeMove(PLAYER_X, 0, 0)
	err = game.MakeMove(PLAYER_O, 1, 0)
	err = game.MakeMove(PLAYER_X, 0, 1)
	err = game.MakeMove(PLAYER_O, 2, 0)
	err = game.MakeMove(PLAYER_X, 0, 2)
	if err != nil {
		panic(err)
	}
	fmt.Println(game)
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
