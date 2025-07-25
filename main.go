package main

import (
	"html/template"
	"log"
	"net/http"

	"learn.01founders.co/git/ftafrial/ascii-art-web/asciiart"
)

type PageData struct {
	Art   string
	Error string
}

func handleHome(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	tmpl, err := template.ParseFiles("templates/index.html")
	if err != nil {
		log.Println("Failed to load index.html:", err)
		http.Error(w, "Template not found", http.StatusNotFound)
		return
	}

	err = tmpl.Execute(w, nil)
	if err != nil {
		log.Println("Error rendering index.html:", err)
		http.Error(w, "Error rendering homepage", http.StatusInternalServerError)
	}
}

func handlePost(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	input := r.FormValue("inputText")
	banner := r.FormValue("banner")

	log.Println("📥 Input:", input)
	log.Println("📦 Banner:", banner)

	art, err := asciiart.GenerateASCII(input, banner)
	log.Println("🔍 ASCII output:\n" + art)

	data := PageData{}

	if err != nil || art == "" {
		if err != nil {
			log.Println("❌ ASCII generation error:", err)
			data.Error = err.Error()
		} else {
			log.Println("❌ ASCII output was empty")
			data.Error = "No ASCII output was generated. Please try again."
		}
	} else {
		log.Println("✅ ASCII Output:\n" + art)
		data.Art = art
	}

	tmpl, err := template.ParseFiles("templates/result.html")
	if err != nil {
		http.Error(w, "Template not found", http.StatusNotFound)
		return
	}

	err = tmpl.Execute(w, data)
	if err != nil {
		log.Println("❌ Template execution error:", err)
		http.Error(w, "Failed to render result page", http.StatusInternalServerError)
	}
}

func main() {
	http.HandleFunc("/", handleHome)
	http.HandleFunc("/ascii-art", handlePost)

	log.Println("🚀 Server running at http://localhost:8080")
	err := http.ListenAndServe(":8080", nil)
	if err != nil {
		log.Fatal("ListenAndServe error:", err)
	}
}
