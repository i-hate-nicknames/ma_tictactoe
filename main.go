package main

func main() {
	exit := make(chan bool, 1)
	startServer(9001)
	go startClient("127.0.0.1", "9001")
	go startClient("127.0.0.1", "9001")
	<-exit
}
