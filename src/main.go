package main

import (
	"log"
	"net/http"
)

func main() {
	server, err := NewServer()
	if err != nil {
		log.Fatalf("Error creating server: %v", err)
	}
	http.HandleFunc("/gettoken", server.GetToken)
	log.Println(http.ListenAndServe(":8081", nil))
}
