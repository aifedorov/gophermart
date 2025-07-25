.PHONY: dev test clean

# Start development environment
dev:
	docker-compose up --build

# Run tests
test:
	go test ./...

# Clean up
clean:
	docker-compose down

# Run locally (need PostgreSQL running)
local:
	air