package main

import (
	"context"
	"fmt"
	"net/http"
	"net/url"
	"strconv"

	"github.com/ethereum/go-ethereum/log"
)

type Server struct {
	transactionHandler *TransactionHandler
}

func NewServer() (*Server, error) {
	ctx := context.Background()
	cfg, err := NewConfig()
	if err != nil {
		return nil, err
	}
	transactionHandler, err := NewTransactionHandler(ctx, cfg)
	if err != nil {
		return nil, err
	}
	return &Server{
		transactionHandler: transactionHandler,
	}, nil
}

func (s *Server) GetToken(w http.ResponseWriter, r *http.Request) {
	query := r.URL.Query()

	to := query.Get("to")
	id, err := getInt64(&query, "id")
	if err != nil {
		handleError(w, err)
		return
	}
	quantity, err := getInt64(&query, "quantity")
	if err != nil {
		handleError(w, err)
		return
	}

	txHash, err := s.transactionHandler.Erc1155Transfer(r.Context(), to, id, quantity)
	if err != nil {
		handleError(w, err)
		return
	}
	_, writeErr := w.Write([]byte(fmt.Sprintf("Transaction hash: %s", txHash)))
	if writeErr != nil {
		log.Error(fmt.Sprintf("Error writing response %v", writeErr))
	}
}

func getInt64(query *url.Values, field string) (int64, error) {
	val := query.Get(field)
	return strconv.ParseInt(val, 10, 64)
}

func handleError(w http.ResponseWriter, err error) {
	w.WriteHeader(http.StatusInternalServerError)
	_, writeErr := w.Write([]byte(fmt.Sprintf("500 - Internal Server Error %v", err)))
	if writeErr != nil {
		log.Error(fmt.Sprintf("Error writing error response %v", writeErr))
	}
}
