package main

import (
	"fmt"
	"html/template"
	"log"
	"net/http"

	"ascii-art-web/asciiart" // Adjust this if needed depending on your module name
)

// handleHome renders the homepage with the form
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

// handlePost processes the submitted text and displays ASCII art
func handlePost(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	input := r.FormValue("text")
	banner := "standard" // We’ll support this first

	asciiOutput, err := asciiart.GenerateASCII(input, banner)
	if err != nil {
		http.Error(w, "Error generating ASCII art", http.StatusInternalServerError)
		log.Println("ASCII generation error:", err)
		return
	}

	// Display the result as plain text (or update to HTML later)
	w.Header().Set("Content-Type", "text/plain")
	fmt.Fprint(w, asciiOutput)
}

// main starts the server and registers routes
func main() {
	http.HandleFunc("/", handleHome)
	http.HandleFunc("/ascii-art", handlePost)

	fmt.Println("Server running at http://localhost:8080/")
	err := http.ListenAndServe(":8080", nil)
	if err != nil {
		log.Fatal("Error starting server:", err)
	}
}
