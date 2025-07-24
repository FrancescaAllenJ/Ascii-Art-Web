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

	// Parse form input
	err := r.ParseForm()
	if err != nil {
		http.Error(w, "Bad Request: Unable to parse form", http.StatusBadRequest)
		log.Println("Error parsing form:", err)
		return
	}

	inputText := r.FormValue("inputText")
	banner := r.FormValue("banner")

	// Validate form input
	if inputText == "" || banner == "" {
		http.Error(w, "Bad Request: Missing input or banner", http.StatusBadRequest)
		return
	}

	// Call ASCII generator (replace this with your own logic)
	asciiArt, err := GenerateASCIIArt(inputText, banner)
	if err != nil {
		http.Error(w, "Internal Server Error: Failed to generate ASCII art", http.StatusInternalServerError)
		log.Println("ASCII generation error:", err)
		return
	}

	// Load template again to render result
	tmpl, err := template.ParseFiles("templates/index.html")
	if err != nil {
		http.Error(w, "Template not found", http.StatusNotFound)
		log.Println("Error loading result template:", err)
		return
	}

	// Data to pass into template
	data := struct {
		Art    string
		Text   string
		Banner string
	}{
		Art:    asciiArt,
		Text:   inputText,
		Banner: banner,
	}

	// Render filled template with ASCII art
	err = tmpl.Execute(w, data)
	if err != nil {
		http.Error(w, "Internal Server Error: Template execution failed", http.StatusInternalServerError)
		log.Println("Template execution error:", err)
	}
}

// main sets up routes and starts the server
func main() {
	http.HandleFunc("/", handleHome)
	http.HandleFunc("/ascii-art", handlePost)

	fmt.Println("Server running at http://localhost:8080/")
	err := http.ListenAndServe(":8080", nil)
	if err != nil {
		fmt.Println("Error starting server:", err)
	}
}
