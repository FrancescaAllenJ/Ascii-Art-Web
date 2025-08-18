package main

import (
	"html/template"
	"log"
	"net/http"
	"strings"

	"ascii-art-web/asciiart"
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

	// parse templates at startup
	tplIndex, err = template.ParseFiles("templates/index.html")
	if err != nil {
		log.Fatal("load index template:", err)
	}
	tplResult, err = template.ParseFiles("templates/result.html")
	if err != nil {
		log.Fatal("load result template:", err)
	}

	// routes
	mux := http.NewServeMux()
	mux.HandleFunc("/", handleHome)          // GET /
	mux.HandleFunc("/ascii-art", handlePost) // POST /ascii-art

	// 404 wrapper: only serve known routes
	notFound := http.NotFoundHandler()
	root := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		_, pattern := mux.Handler(r)
		if pattern == "" {
			// unknown route
			notFound.ServeHTTP(w, r)
			return
		}
		mux.ServeHTTP(w, r)
	})

	log.Println("🚀 listening on http://localhost:8080")
	if err := http.ListenAndServe(":8080", root); err != nil {
		log.Fatal("server failed:", err)
	}
}

func handleHome(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}
	if err := tplIndex.Execute(w, nil); err != nil {
		// template missing = 404 as per brief
		http.Error(w, "template not found", http.StatusNotFound)
		return
	}
}

func handlePost(w http.ResponseWriter, r *http.Request) {
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
		// map known user errors to 400; everything else 500
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

	data := PageData{Art: art}
	if err := tplResult.Execute(w, data); err != nil {
		http.Error(w, "template not found", http.StatusNotFound)
		return
	}
}
