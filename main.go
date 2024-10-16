package main

import "github.com/zvxte/kera/server"

func main() {
	server := server.NewServer()
	server.Run(":5000")
}
