package main

import (
	"encoding/json"
	"net/http"
	"strings"

	"./datastorage"
)

type App struct {
	storage datastorage.Storage
}

// shortenHandler обрабатывает POST /shorten
func (a *App) shortenHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var req struct {
		URL string `json:"url"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil || req.URL == "" {
		http.Error(w, "Invalid request", http.StatusBadRequest)
		return
	}

	// Если для данного URL уже существует код, возвращаем его
	if code, err := a.storage.FindByURL(req.URL); err == nil {
		json.NewEncoder(w).Encode(map[string]string{"short_url": code})
		return
	}

	// Генерируем уникальный код
	var code string
	for {
		code = datastorage.GenerateCode()
		if _, err := a.storage.Get(code); err != nil {
			break
		}
	}

	if err := a.storage.Save(req.URL, code); err != nil {
		http.Error(w, "Error saving URL", http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(map[string]string{"short_url": code})
}

// redirectHandler обрабатывает GET /{code}
func (a *App) redirectHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}
	code := strings.TrimPrefix(r.URL.Path, "/")
	if code == "" {
		http.Error(w, "Code not provided", http.StatusBadRequest)
		return
	}
	url, err := a.storage.Get(code)
	if err != nil {
		http.Error(w, "Not found", http.StatusNotFound)
		return
	}
	json.NewEncoder(w).Encode(map[string]string{"url": url})
}
