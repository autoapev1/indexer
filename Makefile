

build:
	@echo "Building..."
	@go build -o ./bin/test ./cmd/test/main.go

run: build
	@./bin/test


postgres-up:
	docker compose -f ./docker/postgres.yml up -d --remove-orphans

postgres-down:
	docker compose -f ./docker/postgres.yml down