package datastorage

import (
	"database/sql"
	"fmt"
	"math/rand"
)

const codeLength = 10

var allowedChars = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789_")

// Storage определяет интерфейс для работы с хранилищем URL
type Storage interface {
	Get(code string) (string, error)
	Save(url, code string) error
	FindByURL(url string) (string, error)
}

// generateCode генерирует случайный код длиной codeLength
func GenerateCode() string {
	b := make([]rune, codeLength)
	for i := range b {
		b[i] = allowedChars[rand.Intn(len(allowedChars))]
	}
	return string(b)
}

// ---------- Реализация in-memory хранилища ----------

type InMemoryStorage struct {
	codeToURL map[string]string
	urlToCode map[string]string
}

func NewInMemoryStorage() *InMemoryStorage {
	return &InMemoryStorage{
		codeToURL: make(map[string]string),
		urlToCode: make(map[string]string),
	}
}

func (s *InMemoryStorage) Get(code string) (string, error) {
	url, ok := s.codeToURL[code]
	if !ok {
		return "", fmt.Errorf("not found")
	}
	return url, nil
}

func (s *InMemoryStorage) Save(url, code string) error {
	if _, exists := s.codeToURL[code]; exists {
		return fmt.Errorf("code already exists")
	}
	s.codeToURL[code] = url
	s.urlToCode[url] = code
	return nil
}

func (s *InMemoryStorage) FindByURL(url string) (string, error) {
	code, ok := s.urlToCode[url]
	if !ok {
		return "", fmt.Errorf("not found")
	}
	return code, nil
}

// ---------- Реализация PostgreSQL хранилища ----------

type PostgresStorage struct {
	db *sql.DB
}

func NewPostgresStorage(connStr string) (*PostgresStorage, error) {
	db, err := sql.Open("postgres", connStr)
	if err != nil {
		return nil, err
	}
	// Создаем таблицу, если она отсутствует
	query := `
	CREATE TABLE IF NOT EXISTS urls (
		code VARCHAR(10) PRIMARY KEY,
		url TEXT NOT NULL UNIQUE
	);
	`
	_, err = db.Exec(query)
	if err != nil {
		return nil, err
	}
	return &PostgresStorage{db: db}, nil
}

func (s *PostgresStorage) Get(code string) (string, error) {
	var url string
	err := s.db.QueryRow("SELECT url FROM urls WHERE code=$1", code).Scan(&url)
	if err != nil {
		return "", err
	}
	return url, nil
}

func (s *PostgresStorage) Save(url, code string) error {
	_, err := s.db.Exec("INSERT INTO urls (code, url) VALUES ($1, $2)", code, url)
	return err
}

func (s *PostgresStorage) FindByURL(url string) (string, error) {
	var code string
	err := s.db.QueryRow("SELECT code FROM urls WHERE url=$1", url).Scan(&code)
	if err != nil {
		return "", err
	}
	return code, nil
}
