package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"

	"github.com/apex/gateway"
	"github.com/glonlas/sff-aicp/internal/server"
	"github.com/rs/cors"
)

func main() {
	displayTerminalBanner()

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

	
	mux = router(mux, server)

	handler := corsOpts.Handler(mux)

	if *port != -1 {
		portStr := fmt.Sprintf(":%d", *port)
		log.Println(http.ListenAndServe(portStr, handler))
	} else {
		log.Println(gateway.ListenAndServe("n/a", handler))
	}
}

func displayTerminalBanner() {
	banner := `
   █████╗     ██╗     ██████╗    ██████╗ 
  ██╔══██╗    ██║    ██╔════╝    ██╔══██╗
  ███████║    ██║    ██║         ██████╔╝
  ██╔══██║    ██║    ██║         ██╔═══╝ 
  ██║  ██║    ██║    ╚██████╗    ██║     
  ╚═╝  ╚═╝    ╚═╝     ╚═════╝    ╚═╝     
`
	fmt.Println(banner)
}

// Here is where you define the public API of the service
func router(mux *http.ServeMux, server *server.Server) *http.ServeMux {
	// APIs for Customer actions
	mux.HandleFunc("/api/getCustomerAvailableBalance", server.GetCustomerAvailableBalance)
	mux.HandleFunc("/api/createOrder", server.CreateOrder)
	mux.HandleFunc("/api/getOrderStatus", server.GetOrderStatus)
	
	
	// APIs for Merchant actions
	// Note: This API should be in a separate micro-service to ensure a good separation of concerns.
	mux.HandleFunc("/api/cancelOrder", server.CancelOrder)

	// Utilities APIs
	// Internal API used to help you test and debug this POC
	mux.HandleFunc("/api/getCustomerPublicAddress", server.GetCustomerPublicAddress)
	mux.HandleFunc("/api/getTransactionStatus", server.GetTransactionStatus)
	mux.HandleFunc("/api/mintTokens", server.MintTokens)

	return mux
}
