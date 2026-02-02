package testutil

import (
	"database/sql"
	"fmt"
	"log"
	"os"
	"testing"

	_ "github.com/lib/pq"
)

// ! SetupTestDB crée une DB de test et retourne la connexion
func SetupTestDB(t *testing.T) *sql.DB {
	t.Helper()

	host := getEnv("TEST_DB_HOST", "localhost")
	port := getEnv("TEST_DB_PORT", "5432")
	user := getEnv("TEST_DB_USER", "postgres")
	password := getEnv("TEST_DB_PASSWORD", "postgres")
	dbname := getEnv("TEST_DB_NAME", "educnet_test")

	dsn := fmt.Sprintf(
		"host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		host, port, user, password, dbname,
	)

	db, err := sql.Open("postgres", dsn)
	if err != nil {
		t.Fatalf("Failed to connect to test database: %v", err)
	}

	if err := db.Ping(); err != nil {
		t.Fatalf("Failed to ping test database: %v", err)
	}

	//! Cleanup à la fin du test
	t.Cleanup(func() {
		db.Close()
	})

	return db
}

// ! CleanupTestDB nettoie toutes les tables
func CleanupTestDB(t *testing.T, db *sql.DB) {
	t.Helper()

	queries := []string{
		"TRUNCATE TABLE users CASCADE",
		"TRUNCATE TABLE schools CASCADE",
		"ALTER SEQUENCE schools_id_seq RESTART WITH 1",
		"ALTER SEQUENCE users_id_seq RESTART WITH 1",
	}

	for _, query := range queries {
		if _, err := db.Exec(query); err != nil {
			log.Printf("Warning: cleanup query failed: %v", err)
		}
	}
}

// ! SeedTestSchool crée une école de test
func SeedTestSchool(t *testing.T, db *sql.DB, name, slug, email string) int {
	t.Helper()

	var id int
	query := `
		INSERT INTO schools (name, slug, email, address, phone, status)
		VALUES ($1, $2, $3, 'Test Address', '+261 34 00 000 00', 'active')
		RETURNING id
	`

	err := db.QueryRow(query, name, slug, email).Scan(&id)
	if err != nil {
		t.Fatalf("Failed to seed test school: %v", err)
	}

	return id
}

// ! SeedTestUser crée un utilisateur de test
func SeedTestUser(t *testing.T, db *sql.DB, schoolID int, email, role string) int {
	t.Helper()

	var id int
	query := `
		INSERT INTO users (school_id, email, password_hash, first_name, last_name, role, status)
		VALUES ($1, $2, 'hashed', 'Test', 'User', $3, 'approved')
		RETURNING id
	`

	err := db.QueryRow(query, schoolID, email, role).Scan(&id)
	if err != nil {
		t.Fatalf("Failed to seed test user: %v", err)
	}

	return id
}

func SeedTestClass(t *testing.T, db *sql.DB, schoolID int, name, level, section, year string) int {
	t.Helper()

	var id int
	query := `
        INSERT INTO classes (school_id, name, level, section, capacity, academic_year, status)
        VALUES ($1, $2, $3, $4, 40, $5, 'active')
        RETURNING id
    `

	err := db.QueryRow(query, schoolID, name, level, section, year).Scan(&id)
	if err != nil {
		t.Fatalf("Failed to seed test class: %v", err)
	}

	return id
}

func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}
