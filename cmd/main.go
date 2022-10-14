package main

import (
	"log"
	"net/http"

	"github.cbhq.net/engineering/sff-workshop/internal/server"
)

func main() {
	server, err := server.NewServer()
	if err != nil {
		log.Fatalf("Error creating server: %v", err)
	}
	http.HandleFunc("/gettoken", server.GetToken)
	log.Println(http.ListenAndServe(":8081", nil))
}
