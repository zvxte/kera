package main

import (
	"log"
	"net/http"
)

func main() {
	router := NewRouter()
	server := &http.Server{
		Addr:    ":5000",
		Handler: router,
	}
	err := server.ListenAndServe()
	if err != nil {
		log.Fatal(err)
	}
}
