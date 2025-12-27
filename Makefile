.PHONY: help dev-backend dev-frontend build-frontend build run clean install test

help: ## Tampilkan bantuan
	@echo "Available commands:"
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | sort | awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-20s\033[0m %s\n", $$1, $$2}'

install: ## Install semua dependencies
	@echo "Installing Go dependencies..."
	cd backend && go mod download
	@echo "Installing frontend dependencies..."
	cd frontend && yarn install
	@echo "✅ Dependencies installed!"

dev-backend: ## Jalankan backend development server
	@echo "Starting backend server..."
	cd backend && go run cmd/server/main.go

dev-frontend: ## Jalankan frontend development server
	@echo "Starting frontend server..."
	cd frontend && yarn dev

build-frontend: ## Build frontend untuk production
	@echo "Building frontend..."
	cd frontend && yarn build
	@echo "✅ Frontend built!"

build: build-frontend ## Build aplikasi untuk production
	@echo "Building backend..."
	cd backend && go build -o ../sirine-go cmd/server/main.go
	@echo "✅ Build complete! Run with: ./sirine-go"

run: ## Jalankan aplikasi production
	./sirine-go

clean: ## Bersihkan build artifacts
	@echo "Cleaning..."
	rm -f sirine-go
	rm -rf frontend/dist
	rm -rf frontend/node_modules
	@echo "✅ Cleaned!"

test: ## Jalankan tests
	@echo "Running Go tests..."
	cd backend && go test -v ./...
	@echo "Running frontend tests..."
	cd frontend && yarn test

db-create: ## Buat database
	@echo "Creating database..."
	mysql -u root -p -e "CREATE DATABASE IF NOT EXISTS sirine_go CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci;"
	@echo "✅ Database created!"

db-migrate: ## Jalankan migrations
	@echo "Running migrations..."
	cd backend && go run cmd/server/main.go
	@echo "✅ Migrations complete!"
