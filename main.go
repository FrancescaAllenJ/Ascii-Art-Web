package main

import (
	"fmt"
	"html/template"
	"log"
	"net/http"

	"learn.01founders.co/git/ftafrial/ascii-art-web/asciiart"
)

// handleHome renders the homepage form
func handleHome(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	tmpl, err := template.ParseFiles("templates/index.html")
	if err != nil {
		http.Error(w, "Template not found", http.StatusNotFound)
		log.Println("Error loading index.html:", err)
		return
	}

	err = tmpl.Execute(w, nil)
	if err != nil {
		http.Error(w, "Error rendering page", http.StatusInternalServerError)
		log.Println("Template execution error:", err)
	}
}

// handlePost handles the submitted form and displays ASCII result
func handlePost(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	input := r.FormValue("inputText") // Matches the name="inputText" in HTML form
	banner := r.FormValue("banner")   // Matches name="banner" in form

	if banner == "" {
		banner = "standard"
	}

	asciiOutput, err := asciiart.GenerateASCII(input, banner)
	if err != nil {
		log.Println("ASCII generation error:", err)
		tmpl, _ := template.ParseFiles("templates/result.html")
		tmpl.Execute(w, map[string]string{
			"Error": "There was a problem generating the ASCII art.",
		})
		return
	}

	tmpl, err := template.ParseFiles("templates/result.html")
	if err != nil {
		http.Error(w, "Result template not found", http.StatusNotFound)
		log.Println("Error loading result.html:", err)
		return
	}

	data := map[string]string{
		"Art": asciiOutput,
	}
	tmpl.Execute(w, data)
}

// main starts the server
func main() {
	http.HandleFunc("/", handleHome)
	http.HandleFunc("/ascii-art", handlePost)

	fmt.Println("Server running at http://localhost:8080/")
	if err := http.ListenAndServe(":8080", nil); err != nil {
		log.Fatal("Server error:", err)
	}
}
