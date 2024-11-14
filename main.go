package main

import (
	"os"

	"github.com/zvxte/kera/server"
)

func main() {
	address := os.Getenv("ADDRESS")
	if address == "" {
		panic("ADDRESS is not set")
	}

	server, err := server.NewServer()
	if err != nil {
		panic(err)
	}
	server.Run(address)
}
