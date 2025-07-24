package main

import (
	"fmt"
	"html/template"
	"log"
	"net/http"
)

// handleHome handles GET /
func handleHome(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	// Parse and execute the template for the homepage
	tmpl, err := template.ParseFiles("templates/index.html")
	if err != nil {
		http.Error(w, "Template not found", http.StatusNotFound)
		log.Println("Error parsing template:", err)
		return
	}

	err = tmpl.Execute(w, nil)
	if err != nil {
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		log.Println("Error executing template:", err)
	}
}

// handlePost handles POST /ascii-art
func handlePost(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	fmt.Fprintln(w, "ASCII art generation coming soon!")
}

// main sets up routes and starts the server
func main() {
	http.HandleFunc("/", handleHome)          // For GET requests to homepage
	http.HandleFunc("/ascii-art", handlePost) // For POST requests

	fmt.Println("Server running at http://localhost:8080/")
	err := http.ListenAndServe(":8080", nil)
	if err != nil {
		fmt.Println("Error starting server:", err)
	}
}
