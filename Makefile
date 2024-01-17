build:
	@echo "Building..."
	@go build -o ./bin/test ./cmd/test/main.go
	@go build -o ./bin/ingest-eth ./cmd/ingest/eth/main.go
	@go build -o ./bin/ingest-bsc ./cmd/ingest/bsc/main.go

run: build
	@./bin/test

ingest-eth: build
	@./bin/ingest-eth

ingest-bsc: build
	@./bin/ingest-bsc

postgres-up:
	docker compose -f ./docker/postgres.yml up -d --remove-orphans

postgres-down:
	docker compose -f ./docker/postgres.yml down