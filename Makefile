.PHONY: up down restart build run seed logs clean

# Docker Compose Commands
up:
	docker compose up -d

down:
	docker compose down -v

restart: down up

logs:
	docker compose logs -f

# Application Commands
run:
	go run cmd/main.go

build:
	go build -o bin/api cmd/main.go

# Database Seeding
seed:
	docker compose exec -T db psql -U user -d itsware < migrations/seeds/seed.sql

# All-in-one setup
init: down up
	@echo "Waiting for DB to start..."
	@sleep 5
	@echo "Applying Schema..."
	cat migrations/*.sql | docker compose exec -T db psql -U user -d itsware
	@echo "Seeding Data..."
	@make seed
	@echo "Application ready! Run 'make run' to start."
