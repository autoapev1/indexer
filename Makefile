build:
	@echo "Building..."
	@go build -o ./bin/api ./cmd/api/main.go
	@go build -o ./bin/ingest-eth ./cmd/ingest/eth/main.go
	@go build -o ./bin/ingest-bsc ./cmd/ingest/bsc/main.go

run: build
	@./bin/api --config config.toml

ingest-eth: build
	@./bin/ingest-eth --config config.toml

ingest-bsc: build
	@./bin/ingest-bsc --config config.toml

postgres-up:
	docker compose -f ./docker/postgres.yml up -d --remove-orphans

postgres-down:
	docker compose -f ./docker/postgres.yml down