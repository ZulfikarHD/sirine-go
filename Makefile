.PHONY: help dev-backend dev-frontend build-frontend build run clean install test db-create db-drop db-migrate db-rollback db-fresh db-seed db-reset hash-make hash-check

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
	cd backend && go run cmd/migrate/main.go create

db-drop: ## Hapus database
	cd backend && go run cmd/migrate/main.go drop

db-migrate: ## Jalankan migrations (up)
	cd backend && go run cmd/migrate/main.go up

db-rollback: ## Rollback migrations (down)
	cd backend && go run cmd/migrate/main.go down

db-fresh: ## Drop database, create ulang, dan migrate
	cd backend && go run cmd/migrate/main.go fresh

db-seed: ## Jalankan database seeder
	cd backend && go run cmd/seed/main.go

db-reset: ## Fresh migration + seed (untuk development)
	@echo "🔄 Resetting database..."
	@$(MAKE) db-fresh
	@$(MAKE) db-seed
	@echo "✅ Database reset complete!"

hash-make: ## Generate bcrypt hash untuk password (Usage: make hash-make PASSWORD="YourPassword")
	@if [ -z "$(PASSWORD)" ]; then \
		echo "❌ Error: PASSWORD tidak diberikan"; \
		echo ""; \
		echo "Usage: make hash-make PASSWORD=\"YourPassword\""; \
		echo "Example: make hash-make PASSWORD=\"Admin@123\""; \
		exit 1; \
	fi
	@cd backend && go run cmd/hash/main.go make "$(PASSWORD)"

hash-check: ## Verify password dengan hash (direct go run recommended)
	@echo "⚠️  Note: Karena $ escaping di Makefile, gunakan command berikut:"
	@echo ""
	@echo "cd backend && go run cmd/hash/main.go check \"YourPassword\" '\$$2a\$$12\$$...'"
	@echo ""
	@echo "Example:"
	@echo "cd backend && go run cmd/hash/main.go check \"Admin@123\" '\$$2a\$$12\$$nNag4FTJB0fiX/f22aINOuYcuP8cUOWOyCZub6tBCom0Evxv4ahTK'"
