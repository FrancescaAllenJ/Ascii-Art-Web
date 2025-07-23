package main

import (
	"fmt"
	"net/http"
)

// main sets up the routes and starts the server
func main() {
	http.HandleFunc("/", handleHome)          // For GET requests
	http.HandleFunc("/ascii-art", handlePost) // For POST requests

	fmt.Println("Server running at http://localhost:8080/")
	err := http.ListenAndServe(":8080", nil)
	if err != nil {
		fmt.Println("Error starting server:", err)
	}
}

// handleHome handles GET /
func handleHome(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	fmt.Fprintln(w, "Welcome to the ASCII Art Web App!")
}

// handlePost handles POST /ascii-art
func handlePost(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	fmt.Fprintln(w, "ASCII art generation coming soon!")
}
