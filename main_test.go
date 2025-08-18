package main

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestHandlePost_MissingInput(t *testing.T) {
	req := httptest.NewRequest("POST", "/ascii-art", strings.NewReader("inputText=&banner=standard"))
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	w := httptest.NewRecorder()

	handlePost(w, req)

	if w.Code != http.StatusBadRequest {
		t.Errorf("expected 400, got %d", w.Code)
	}
}

func TestHandleAPI_InvalidJSON(t *testing.T) {
	req := httptest.NewRequest("POST", "/api/ascii", strings.NewReader("not-json"))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	handleAPI(w, req)

	if w.Code != http.StatusBadRequest {
		t.Errorf("expected 400, got %d", w.Code)
	}
}
