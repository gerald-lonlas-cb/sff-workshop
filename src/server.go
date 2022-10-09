package main

import (
	"context"
	"fmt"
	"html"
	"net/http"
)

type Server struct {
	cloudClient *CloudClient
}

func NewServer(ctx context.Context) (*Server, error) {
	cloudClient := NewCloudClient(ctx)
	return &Server{
		cloudClient: cloudClient,
	}, nil
}

func (s *Server) GetBalance(w http.ResponseWriter, r *http.Request) {
	s.cloudClient.GetBalance(r.Context())
}

func (s *Server) GetToken(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Hello, %q", html.EscapeString(r.URL.Path))
}
