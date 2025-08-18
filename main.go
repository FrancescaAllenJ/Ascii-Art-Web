package main

import (
	"encoding/json"
	"html/template"
	"log"
	"net/http"
	"strings"

	"learn.01founders.co/git/ftafrial/ascii-art-web/asciiart"
)

type PageData struct {
	Art   string
	Error string
}

type APIRequest struct {
	InputText string `json:"inputText"`
	Banner    string `json:"banner"`
}

type APIResponse struct {
	Art   string `json:"art,omitempty"`
	Error string `json:"error,omitempty"`
}

// Serve homepage
func handleHome(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/" {
		handleNotFound(w, r)
		return
	}

	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	tmpl, err := template.ParseFiles("templates/index.html")
	if err != nil {
		log.Println("❌ Failed to load index.html:", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	if err := tmpl.Execute(w, nil); err != nil {
		log.Println("❌ Error rendering index.html:", err)
		http.Error(w, "Error rendering homepage", http.StatusInternalServerError)
	}
}

// Handle HTML form
func handlePost(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	input := r.FormValue("inputText")
	banner := r.FormValue("banner")

	if strings.TrimSpace(input) == "" || strings.TrimSpace(banner) == "" {
		http.Error(w, "Bad Request: Missing input or banner", http.StatusBadRequest)
		return
	}

	validBanners := map[string]bool{"standard": true, "shadow": true, "thinkertoy": true}
	if !validBanners[banner] {
		http.Error(w, "Bad Request: Invalid banner", http.StatusBadRequest)
		return
	}

	log.Println("📥 Input:", input)
	log.Println("🎨 Banner:", banner)

	ascii, err := asciiart.Convert(input, banner)
	data := PageData{}

	if err != nil || ascii == "" {
		if err != nil {
			log.Println("❌ ASCII error:", err)
			data.Error = err.Error()
		} else {
			log.Println("❌ Empty ASCII output")
			data.Error = "No ASCII output was generated."
		}
	} else {
		data.Art = ascii
	}

	tmpl, err := template.ParseFiles("templates/result.html")
	if err != nil {
		log.Println("❌ Failed to load result.html:", err)
		http.Error(w, "Result template not found", http.StatusNotFound)
		return
	}

	if err := tmpl.Execute(w, data); err != nil {
		log.Println("❌ Template execution error:", err)
		http.Error(w, "Error rendering result", http.StatusInternalServerError)
	}
}

// ✅ New: JSON API handler
func handleAPI(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var req APIRequest
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		http.Error(w, "Invalid JSON body", http.StatusBadRequest)
		return
	}

	input := strings.TrimSpace(req.InputText)
	banner := strings.TrimSpace(req.Banner)

	if input == "" || banner == "" {
		http.Error(w, "Missing inputText or banner", http.StatusBadRequest)
		return
	}

	validBanners := map[string]bool{"standard": true, "shadow": true, "thinkertoy": true}
	if !validBanners[banner] {
		http.Error(w, "Invalid banner selected", http.StatusBadRequest)
		return
	}

	ascii, err := asciiart.Convert(input, banner)
	resp := APIResponse{}

	if err != nil {
		resp.Error = err.Error()
	} else {
		resp.Art = ascii
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(resp)
}

// Catch-all 404
func handleNotFound(w http.ResponseWriter, r *http.Request) {
	log.Println("🔍 404 Not Found:", r.URL.Path)
	http.Error(w, "404 - Page Not Found", http.StatusNotFound)
}

func main() {
	mux := http.NewServeMux()

	mux.HandleFunc("/", handleHome)
	mux.HandleFunc("/ascii-art", handlePost)
	mux.HandleFunc("/api/ascii-art", handleAPI)
	mux.HandleFunc("/favicon.ico", http.NotFound)

	log.Println("🚀 Server running on http://localhost:8080")

	err := http.ListenAndServe(":8080", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		muxHandler, pattern := mux.Handler(r)
		if pattern == "" {
			handleNotFound(w, r)
			return
		}
		muxHandler.ServeHTTP(w, r)
	}))
	if err != nil {
		log.Fatal("Server failed:", err)
	}
}
