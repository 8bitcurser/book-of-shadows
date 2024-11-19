package main

import (
	"encoding/json"
	"fmt"
	"net/http"
)

func handleConnection(w http.ResponseWriter, r *http.Request) {
	// Log the path that was requested
	fmt.Printf("Method: %s\nPath: %s\n", r.Method, r.URL.Path)

	// Send response
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{"message": "Got it"})
}

func main() {
	fmt.Println("Server listening on :8080")
	http.HandleFunc("/", handleConnection)
	if err := http.ListenAndServe(":8080", nil); err != nil {
		panic(err)
	}
}
