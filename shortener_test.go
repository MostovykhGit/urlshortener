package main

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/MostovykhGit/urlshortener/datastorage"
)

func TestShortenAndRedirectInMemory(t *testing.T) {
	storage := datastorage.NewInMemoryStorage()
	app := &App{storage: storage}

	// тест POST /shorten
	reqBody, _ := json.Marshal(map[string]string{"url": "http://example.com"})
	req := httptest.NewRequest(http.MethodPost, "/shorten", bytes.NewBuffer(reqBody))
	w := httptest.NewRecorder()
	app.shortenHandler(w, req)
	res := w.Result()
	if res.StatusCode != http.StatusOK {
		t.Fatalf("Expected status OK, got %d", res.StatusCode)
	}
	var respBody map[string]string
	json.NewDecoder(res.Body).Decode(&respBody)
	code, exists := respBody["short_url"]
	if !exists || len(code) != datastorage.CodeLength {
		t.Fatalf("Invalid short_url in response: %v", respBody)
	}

	// тест GET /{code}
	req2 := httptest.NewRequest(http.MethodGet, "/"+code, nil)
	w2 := httptest.NewRecorder()
	app.redirectHandler(w2, req2)
	res2 := w2.Result()
	if res2.StatusCode != http.StatusOK {
		t.Fatalf("Expected status OK, got %d", res2.StatusCode)
	}
	var getResp map[string]string
	json.NewDecoder(res2.Body).Decode(&getResp)
	original, exists := getResp["url"]
	if !exists || original != "http://example.com" {
		t.Fatalf("Expected original url to be http://example.com, got %v", getResp)
	}
}
