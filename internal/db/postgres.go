package db

import (
	"database/sql"
	"fmt"
	"log"

	_ "github.com/lib/pq"
)

//! Connect crée et retourne une connexion database
func Connect(dsn string) (*sql.DB, error) {
	db, err := sql.Open("postgres", dsn)
	if err != nil {
		return nil, fmt.Errorf("failed to open database: %w", err)
	}

	db.SetMaxOpenConns(25)
	db.SetMaxIdleConns(5)

	if err := db.Ping(); err != nil {
		return nil, fmt.Errorf("failed to ping database: %w", err)
	}

	log.Println("Database connected successfully")

	return db, nil
}

//! Close ferme proprement la connexion
func Close(db *sql.DB) error {
	if db != nil {
		log.Println("Closing database connection")
		return db.Close()
	}

	return nil
}
