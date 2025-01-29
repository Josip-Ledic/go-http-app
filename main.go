package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"time"
)

func checkConnection(w http.ResponseWriter, r *http.Request) {
	// Default to 8.8.8.8, but use EXTERNAL_API environment variable if set
	externalAPI := os.Getenv("EXTERNAL_API")
	if externalAPI == "" {
		externalAPI = "8.8.8.8"
	}

	client := http.Client{
		Timeout: 10 * time.Second,
	}

	var responseMessage string

	log.Printf("Incoming request")

	resp, err := client.Get("https://" + externalAPI)
	if err != nil {
		// If we can't connect, return an error message
		log.Printf("Failed to connect to %s: %v", externalAPI, err)
		responseMessage = fmt.Sprintf("Error: Failed to reach external API %s\nDetails: %v\n", externalAPI, err)

		// Set appropriate status code for error
		http.Error(w, responseMessage, http.StatusServiceUnavailable)
		return
	}

	// If we get a response, check the status code
	if resp.StatusCode != 200 {
		log.Printf("Received non-OK status code from %s: %d", externalAPI, resp.StatusCode)
		responseMessage = fmt.Sprintf("Error: Received non-OK status code %d from external API %s\n\n", externalAPI, resp.StatusCode)

		http.Error(w, responseMessage, http.StatusServiceUnavailable)
		return
	}

	responseMessage = fmt.Sprintf("Successfully connected to external API %s\n Status: %d\n\n", externalAPI, resp.StatusCode)

	w.Header().Set("Content-Type", "text/plain")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(responseMessage))

	defer resp.Body.Close()
}

func main() {
	// Handle the root route
	http.HandleFunc("/", checkConnection)

	// Start the HTTP server
	log.Println("Server started on port 8080")
	err := http.ListenAndServe(":8080", nil)
	if err != nil {
		log.Fatalf("Server failed to start: %v", err)
	}
}
