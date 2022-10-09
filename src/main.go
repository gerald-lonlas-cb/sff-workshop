package main

import (
	"context"
	"fmt"
	"log"
)

func main() {
	ctx := context.Background()

	handler := newTransactionHandler(nil, ctx)
	txHash, err := handler.erc1155Transfer("0x4d120d7d8019c7616d4e14249fb696c6a5fe0b6b", 2, 3)
	if err != nil {
		log.Fatalf("error transfering erc1155: %v\n", err)
	}

	fmt.Printf("successfully transfered erc1155: %s\n", txHash)

	//server, err := NewServer(ctx)
	//if err != nil {
	//	log.Fatal("Error creating server")
	//}
	//http.HandleFunc("/gettoken", server.GetToken)
	//http.HandleFunc("/getbalance", server.GetBalance)
	//log.Println(http.ListenAndServe(":8081", nil))
}
