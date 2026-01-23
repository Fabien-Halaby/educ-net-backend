.PHONY: test test-unit test-db test-coverage setup-test-db clean

# Variables
TEST_DB_NAME=educnet_test
TEST_DB_USER=postgres
TEST_DB_PASSWORD=postgres

# Lancer tous les tests
test:
	go test -v ./...

# Tests unitaires seulement (sans DB)
test-unit:
	go test -short -v ./internal/domain ./internal/utils

# Tests avec DB seulement
test-db:
	@echo "Testing repositories..."
	go test -v ./internal/repository

# Coverage
test-coverage:
	go test -coverprofile=coverage.out -covermode=atomic ./...
	go tool cover -html=coverage.out -o coverage.html
	@echo "✅ Coverage report: coverage.html"

# Setup DB de test
setup-test-db:
	@echo "🗄️  Creating test database..."
	PGPASSWORD=$(TEST_DB_PASSWORD) psql -h localhost -U $(TEST_DB_USER) -c "DROP DATABASE IF EXISTS $(TEST_DB_NAME);" || true
	PGPASSWORD=$(TEST_DB_PASSWORD) psql -h localhost -U $(TEST_DB_USER) -c "CREATE DATABASE $(TEST_DB_NAME);"
	PGPASSWORD=$(TEST_DB_PASSWORD) psql -h localhost -U $(TEST_DB_USER) -d $(TEST_DB_NAME) -f migrations/001_init.sql
	@echo "✅ Test database created"

# Clean test DB
clean-test-db:
	PGPASSWORD=$(TEST_DB_PASSWORD) psql -h localhost -U $(TEST_DB_USER) -c "DROP DATABASE IF EXISTS $(TEST_DB_NAME);"
	@echo "🗑️  Test database dropped"

# Clean
clean:
	rm -f coverage.out coverage.html
	go clean -testcache

# Run server
run:
	go run cmd/api/main.go

# Build
build:
	go build -o bin/api cmd/api/main.go

# Help
help:
	@echo "Available targets:"
	@echo "  make test          - Run all tests"
	@echo "  make test-unit     - Run unit tests only (fast)"
	@echo "  make test-db       - Run database tests only"
	@echo "  make setup-test-db - Create test database"
	@echo "  make clean-test-db - Drop test database"
	@echo "  make test-coverage - Generate coverage report"
	@echo "  make run           - Run the server"
	@echo "  make build         - Build the binary"
