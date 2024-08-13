# Build And Development
include .env

dev:
	@go run cmd/api/main.go

doc:
	@swag init -g cmd/api/main.go

dev-reload:
	@air

install:
	@go mod tidy

migration:
	migrate create -ext sql -dir pkg/migrations $(command)

migrateUp:
	migrate -path pkg/migrations -database $(DATABASE_URL) -verbose up

migrateDown:
	migrate -path pkg/migrations -database $(DATABASE_URL) -verbose down

migrateForce:
	migrate -path pkg/migrations -database $(DATABASE_URL) -verbose force $(command)

migrateDrop:
	migrate -path pkg/migrations -database $(DATABASE_URL) -verbose drop

.PHONY: dev doc dev-reload install migration migrateUp migrateDown migrateForce migrateDrop
