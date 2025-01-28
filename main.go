package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"time"
)

// Check the external API connection
func checkConnection(w http.ResponseWriter, r *http.Request) {
	// Default to 8.8.8.8, but use EXTERNAL_API environment variable if set
	externalAPI := os.Getenv("EXTERNAL_API")
	if externalAPI == "" {
		externalAPI = "8.8.8.8" // default
	}

	// Set a timeout for the connection attempt
	client := http.Client{
		Timeout: 5 * time.Second,
	}

	// Prepare the message to send back in the response
	var responseMessage string

	// Make the GET request to the external API
	resp, err := client.Get("http://" + externalAPI)
	if err != nil {
		// If we can't connect, return an error message
		log.Printf("Failed to connect to %s: %v", externalAPI, err)
		responseMessage = fmt.Sprintf("Error: Failed to reach external API %s\nDetails: %v\n", externalAPI, err)
		http.Error(w, responseMessage, http.StatusInternalServerError)
		return
	}
	defer resp.Body.Close()

	// If we get a response, check the status code
	if resp.StatusCode != 200 {
		log.Printf("Received non-OK status code from %s: %d", externalAPI, resp.StatusCode)
		responseMessage = fmt.Sprintf("Error: Received non-OK status code %d from external API %s\nDetails: Expected 200, but got %d\n", externalAPI, resp.StatusCode, resp.StatusCode)
		http.Error(w, responseMessage, http.StatusInternalServerError)
		return
	}

	// If successful, prepare the success message
	responseMessage = fmt.Sprintf("Success: Successfully connected to external API %s\nResponse Status: %d\nDetails: Received a 200 OK response.\n", externalAPI, resp.StatusCode)

	// Return the response message
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(responseMessage))
}

func main() {
	// Handle the root route
	http.HandleFunc("/", checkConnection)

	// Start the HTTP server
	log.Println("Server started on port 8080")
	http.ListenAndServe(":8080", nil)
}
