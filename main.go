package main

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net"
	"net/http"
	"os"
	"time"
)

type Todo struct {
	UserID    int    `json:"userId"`
	ID        int    `json:"id"`
	Title     string `json:"title"`
	Completed bool   `json:"completed"`
}

func createHTTPClient() *http.Client {
	// Custom transport to force new connections
	transport := &http.Transport{
		DisableKeepAlives: true, // Ensures every request uses a fresh connection
		DialContext: (&net.Dialer{
			Timeout:   10 * time.Second,
			KeepAlive: 0, // Disable keep-alive at TCP level
		}).DialContext,
	}

	return &http.Client{
		Transport: transport,
		Timeout:   10 * time.Second,
	}
}

func checkConnection(w http.ResponseWriter, r *http.Request) {
	externalAPI := os.Getenv("EXTERNAL_API")
	if externalAPI == "" {
		externalAPI = "https://jsonplaceholder.typicode.com/todos/1"
	}

	client := createHTTPClient() // New client per request

	log.Println("Incoming request, attempting connection to:", externalAPI)

	resp, err := client.Get(externalAPI)
	if err != nil {
		log.Printf("Failed to connect to %s: %v", externalAPI, err)
		http.Error(w, fmt.Sprintf("Error: Failed to reach external API %s\nDetails: %v\n", externalAPI, err), http.StatusServiceUnavailable)
		return
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		log.Printf("Received non-OK status code from %s: %d", externalAPI, resp.StatusCode)
		http.Error(w, fmt.Sprintf("Error: Received non-OK status code %d from external API %s\n\n", resp.StatusCode, externalAPI), http.StatusServiceUnavailable)
		return
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil || len(body) == 0 {
		log.Printf("Failed to read response from %s", externalAPI)
		http.Error(w, fmt.Sprintf("Error: Received empty response from %s\n\n", externalAPI), http.StatusServiceUnavailable)
		return
	}

	var todo Todo
	if err := json.Unmarshal(body, &todo); err != nil {
		log.Printf("Invalid JSON received from %s", externalAPI)
		http.Error(w, fmt.Sprintf("Error: Invalid JSON received from %s\n\n", externalAPI), http.StatusServiceUnavailable)
		return
	}

	responseMessage := fmt.Sprintf(
		"Successfully connected to external API %s\nStatus: %d\nPayload: %+v\n\n",
		externalAPI, resp.StatusCode, todo,
	)

	w.Header().Set("Content-Type", "text/plain")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(responseMessage))
}

func main() {
	http.HandleFunc("/", checkConnection)

	log.Println("Server started on port 8080")
	if err := http.ListenAndServe(":8080", nil); err != nil {
		log.Fatalf("Server failed to start: %v", err)
	}
}
