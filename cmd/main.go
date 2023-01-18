package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"

	"github.cbhq.net/engineering/sff-workshop/internal/server"
	"github.com/apex/gateway"
	"github.com/rs/cors"
)

func main() {
	port := flag.Int("port", -1, "port for local http dev")
	flag.Parse()
	server, err := server.NewServer(port)
	if err != nil {
		log.Fatalf("Error creating server: %v", err)
	}
	mux := http.NewServeMux()
	corsOpts := cors.New(cors.Options{
		AllowedOrigins: []string{"*"},
		AllowedMethods: []string{
			http.MethodPost,
			http.MethodGet,
			http.MethodOptions,
		},
		AllowedHeaders:   []string{"*"},
		AllowCredentials: false,
	})

	mux.HandleFunc("/api/gettoken", server.GetToken)
	handler := corsOpts.Handler(mux)

	if *port != -1 {
		portStr := fmt.Sprintf(":%d", *port)
		log.Println(http.ListenAndServe(portStr, handler))
	} else {
		log.Println(gateway.ListenAndServe("n/a", handler))
	}
}
