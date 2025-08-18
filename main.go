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
	tplIndex    *template.Template
	tplResult   *template.Template
	tplNotFound *template.Template
)

func main() {
	var err error

	// Parse templates at startup (fail fast for index/result).
	tplIndex, err = template.ParseFiles("templates/index.html")
	if err != nil {
		log.Fatal("load index template:", err)
	}
	tplResult, err = template.ParseFiles("templates/result.html")
	if err != nil {
		log.Fatal("load result template:", err)
	}
	// Best-effort: pretty 404 page; fallback to plain 404 if missing
	tplNotFound, _ = template.ParseFiles("templates/404.html")

	mux := http.NewServeMux()
	mux.HandleFunc("/", handleHome)          // GET /
	mux.HandleFunc("/ascii-art", handlePost) // POST /ascii-art

	// Serve static files: /static/... -> ./static/...
	mux.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir("static"))))

	log.Println("listening on http://localhost:8080")
	if err := http.ListenAndServe(":8080", mux); err != nil {
		log.Fatal("server failed:", err)
	}
}

// renderNotFound tries the template and falls back to plain 404 text
func renderNotFound(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusNotFound)
	if tplNotFound != nil {
		w.Header().Set("Content-Type", "text/html; charset=utf-8")
		_ = tplNotFound.Execute(w, nil)
		return
	}
	http.NotFound(w, r)
}

// GET /
func handleHome(w http.ResponseWriter, r *http.Request) {
	// exact path only
	if r.URL.Path != "/" {
		renderNotFound(w, r)
		return
	}
	if r.Method != http.MethodGet {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	if err := tplIndex.Execute(w, nil); err != nil {
		renderNotFound(w, r)
		return
	}
}

// POST /ascii-art
func handlePost(w http.ResponseWriter, r *http.Request) {
	// exact path only
	if r.URL.Path != "/ascii-art" {
		renderNotFound(w, r)
		return
	}
	if r.Method != http.MethodPost {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}
	if err := r.ParseForm(); err != nil {
		http.Error(w, "bad request", http.StatusBadRequest)
		return
	}

	input := r.PostForm.Get("inputText")
	banner := r.PostForm.Get("banner")
	if banner == "" {
		http.Error(w, "missing banner", http.StatusBadRequest)
		return
	}

	// normalise CRLF from browsers
	input = strings.ReplaceAll(input, "\r\n", "\n")
	input = strings.ReplaceAll(input, "\r", "")

	art, err := asciiart.Convert(input, banner)
	if err != nil {
		msg := err.Error()
		switch {
		case strings.Contains(msg, "unknown banner"),
			strings.Contains(msg, "unsupported character"),
			strings.Contains(msg, "malformed"),
			strings.Contains(msg, "banner"):
			http.Error(w, msg, http.StatusBadRequest)
			return
		default:
			http.Error(w, "internal error", http.StatusInternalServerError)
			return
		}
	}

	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	data := PageData{Art: art}
	if err := tplResult.Execute(w, data); err != nil {
		renderNotFound(w, r)
		return
	}
}
