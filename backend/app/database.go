package application

import (
	"database/sql"
	"fmt"
	"log"
	"os"
	"sync"

	_ "github.com/mattn/go-sqlite3"
)

var (
	db   *sql.DB
	once sync.Once
)

func InitDB() error {
	var initErr error
	once.Do(func() {
		dbPath := getEnv("DB_PATH", "../database/betterforum.db")
		var err error
		db, err = sql.Open("sqlite3", dbPath)
		if err != nil {
			initErr = fmt.Errorf("failed to open database: %w", err)
			return
		}
		if err = db.Ping(); err != nil {
			initErr = fmt.Errorf("failed to connect to database: %w", err)
			return
		}
		if err = applySchema(getEnv("SCHEMA_PATH", "../database/schema.sql")); err != nil {
			initErr = fmt.Errorf("failed to apply schema: %w", err)
			return
		}
		log.Println("Database initialized successfully")
	})
	return initErr
}

func CloseDB() {
	if db != nil {
		if err := db.Close(); err != nil {
			log.Printf("Error closing database: %v", err)
		} else {
			log.Println("Database connection closed")
		}
	}
}

func GetDB() *sql.DB {
	return db
}

func applySchema(path string) error {
	content, err := os.ReadFile(path)
	if err != nil {
		return fmt.Errorf("failed to read schema file: %w", err)
	}
	if _, err := db.Exec(string(content)); err != nil {
		return fmt.Errorf("failed to execute schema: %w", err)
	}
	return nil
}

func getEnv(key, fallback string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}
	return fallback
}
