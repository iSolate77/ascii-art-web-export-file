package main

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestGenerateASCIIArt(t *testing.T) {
	result := generateASCIIArt("A", "standard")
	if strings.TrimSpace(result) == "" {
		t.Errorf("Expected non-empty string, got empty")
	}
}

func TestMainPageHandler(t *testing.T) {
	req, err := http.NewRequest("GET", "/", nil)
	if err != nil {
		t.Fatal(err)
	}

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(mainPageHandler)

	handler.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusOK)
	}
}

func TestAsciiArtHandler(t *testing.T) {
	req, err := http.NewRequest("POST", "/ascii-art", strings.NewReader("text=A&banner=standard"))
	if err != nil {
		t.Fatal(err)
	}
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(asciiArtHandler)

	handler.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusOK)
	}
}
