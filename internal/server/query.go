package server

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"net/url"
	"strings"
)


func (s *Server) verifyQueryParams(query url.Values, requiredParams []string) error {
	missingParams := []string{}
	for _, param := range requiredParams {
		if query.Get(param) == "" {
			missingParams = append(missingParams, param)
		}
	}

	if len(missingParams) > 0 {
		// Handle the missing parameters
		err := "Missing GET parameters: " + strings.Join(missingParams, ", ")
		return fmt.Errorf(err)
	}
	
	return  nil
}

func (s *Server) responseWithError(w http.ResponseWriter, err error) {
	w.WriteHeader(http.StatusInternalServerError)
	_, writeErr := w.Write([]byte(fmt.Sprintf("%v", err)))
	if writeErr != nil {
		log.Printf("Error writing error response %v", writeErr)
	}
}

func (s *Server) respondWithJSON(w http.ResponseWriter, statusCode int, data interface{}) {
	response, err := json.Marshal(data)
	if err != nil {
		s.responseWithError(w, fmt.Errorf("Failed to marshal JSON response: %v", err))
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	_, err = w.Write(response)
	if err != nil {
		log.Printf("Failed to write JSON response: %v\n", err)
	}
}

