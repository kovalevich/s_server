package main

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"strings"
	"sync"
)

var (
	requests []map[string]interface{}
	mu       sync.Mutex
)

func main() {
	port := os.Getenv("PORT")

	if port == "" {
		port = "8080"
	}

	http.HandleFunc("/api/add", handleAdd)
	http.HandleFunc("/api/get", handleGet)
	log.Println("Listening on port", port)
	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%s", port), nil))
}

func handleAdd(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	body, err := io.ReadAll(r.Body)
	if err != nil {
		http.Error(w, "Error reading request body", http.StatusInternalServerError)
		return
	}
	defer r.Body.Close()

	// JSON validation
	var data map[string]interface{}
	if err := json.Unmarshal(body, &data); err != nil {
		http.Error(w, "Invalid JSON format", http.StatusBadRequest)
		return
	}

	abck, ok := data["cookie"].(string)
	if !ok {
		http.Error(w, "Invalid JSON format: 'cookie' key is missing or not a string", http.StatusBadRequest)
		return
	}

	if !strings.Contains(abck, "_abck") {
		http.Error(w, "Invalid JSON format: 'cookie' value should contains '_abck'", http.StatusBadRequest)
		return
	}

	mu.Lock()
	defer mu.Unlock()
	requests = append(requests, data)
}

func handleGet(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	mu.Lock()
	defer mu.Unlock()

	if len(requests) == 0 {
		w.Header().Set("Content-Type", "application/json")
		w.Write([]byte("[]"))
		return
	}

	response, err := json.Marshal(requests)
	if err != nil {
		http.Error(w, "Error encoding response", http.StatusInternalServerError)
		return
	}

	requests = nil // Clear the array after getting its contents

	w.Header().Set("Content-Type", "application/json")
	w.Write(response)
}
