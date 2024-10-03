package main

import "log"

func main() {
	router := NewRouter()
	server := NewServer(":5000", router)
	err := server.Run()
	if err != nil {
		log.Fatal(err)
	}
}
