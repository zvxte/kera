package main

import (
	"github.com/zvxte/kera/server"
)

func main() {
	server, err := server.NewServer()
	if err != nil {
		panic(err)
	}
	server.Run(":5000")
}
