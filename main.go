package main

import (
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

var (
	tplIndex  *template.Template
	tplResult *template.Template
)

func main() {
	var err error

	// Parse templates at startup (fail fast if missing).
	tplIndex, err = template.ParseFiles("templates/index.html")
	if err != nil {
		log.Fatal("load index template:", err)
	}
	tplResult, err = template.ParseFiles("templates/result.html")
	if err != nil {
		log.Fatal("load result template:", err)
	}

	// Routes.
	mux := http.NewServeMux()
	mux.HandleFunc("/", handleHome)          // GET /
	mux.HandleFunc("/ascii-art", handlePost) // POST /ascii-art

	log.Println("🚀 listening on http://localhost:8080")
	if err := http.ListenAndServe(":8080", mux); err != nil {
		log.Fatal("server failed:", err)
	}
}

// handleHome serves the main page.
// Returns 404 for any path other than "/" and 405 for wrong methods.
func handleHome(w http.ResponseWriter, r *http.Request) {
	// Require exact path.
	if r.URL.Path != "/" {
		http.NotFound(w, r) // 404
		return
	}
	if r.Method != http.MethodGet {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed) // 405
		return
	}

	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	if err := tplIndex.Execute(w, nil); err != nil {
		// Treat missing/broken template as 404 per brief.
		http.Error(w, "template not found", http.StatusNotFound) // 404
		return
	}
}

// handlePost renders ASCII art.
// Returns 404 for non-exact path, 405 for wrong method, 400 for bad input/banner.
func handlePost(w http.ResponseWriter, r *http.Request) {
	// Require exact path.
	if r.URL.Path != "/ascii-art" {
		http.NotFound(w, r) // 404
		return
	}
	if r.Method != http.MethodPost {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed) // 405
		return
	}
	if err := r.ParseForm(); err != nil {
		http.Error(w, "bad request", http.StatusBadRequest) // 400
		return
	}

	input := r.PostForm.Get("inputText")
	banner := r.PostForm.Get("banner")
	if banner == "" {
		http.Error(w, "missing banner", http.StatusBadRequest) // 400
		return
	}

	// Normalise CRLF from browsers.
	input = strings.ReplaceAll(input, "\r\n", "\n")
	input = strings.ReplaceAll(input, "\r", "")

	art, err := asciiart.Convert(input, banner)
	if err != nil {
		// Map user-correctable errors to 400; everything else 500.
		msg := err.Error()
		switch {
		case strings.Contains(msg, "unknown banner"),
			strings.Contains(msg, "unsupported character"),
			strings.Contains(msg, "malformed"),
			strings.Contains(msg, "banner"):
			http.Error(w, msg, http.StatusBadRequest) // 400
			return
		default:
			http.Error(w, "internal error", http.StatusInternalServerError) // 500
			return
		}
	}

	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	data := PageData{Art: art}
	if err := tplResult.Execute(w, data); err != nil {
		http.Error(w, "template not found", http.StatusNotFound) // 404
		return
	}
}
