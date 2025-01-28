package main

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"time"
)

func main() {
	// Default to 8.8.8.8 if EXTERNAL_API is not set
	externalAPI := os.Getenv("EXTERNAL_API")
	if externalAPI == "" {
		externalAPI = "http://8.8.8.8"
	}

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		// Set a timeout for the outgoing request
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()

		// Create a new HTTP request with the context
		req, err := http.NewRequestWithContext(ctx, http.MethodGet, externalAPI, nil)
		if err != nil {
			http.Error(w, fmt.Sprintf("Failed to create request: %v", err), http.StatusInternalServerError)
			return
		}

		// Perform the request
		client := &http.Client{}
		resp, err := client.Do(req)
		if err != nil {
			http.Error(w, fmt.Sprintf("Error reaching external API: %v", err), http.StatusServiceUnavailable)
			return
		}
		defer resp.Body.Close()

		// If the request succeeds, return a success message
		w.WriteHeader(http.StatusOK)
		fmt.Fprintf(w, "Successfully connected to external API: %s\n", externalAPI)
	})

	fmt.Printf("Server running on http://localhost:8080\nExternal API: %s\n", externalAPI)
	if err := http.ListenAndServe(":8080", nil); err != nil {
		fmt.Printf("Server failed: %v\n", err)
	}
}
