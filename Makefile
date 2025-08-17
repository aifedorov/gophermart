.PHONY: docker-up docker-down test statictest autotests autotests-docker ci-test ci-test-docker build help

docker-up:
	docker-compose up --build

docker-down:
	docker-compose down

test:
	go test ./...

run:
	go run cmd/gophermart/main.go

build:
	@echo "Building gophermart binary..."
	(cd cmd/gophermart && go build -buildvcs=false -o gophermart)
	@echo "Build complete: cmd/gophermart/gophermart"

statictest: build
	@echo "Running statictest..."
	go vet -vettool=$(PWD)/.tools/statictest ./...

autotests: build
	@echo "Running autotests..."
	@echo "Note: This requires PostgreSQL running locally with database 'praktikum'"
	@echo "Example connection: postgresql://postgres:postgres@localhost/praktikum?sslmode=disable"
	@if [ -z "$(DATABASE_URI)" ]; then \
		echo "Please set DATABASE_URI environment variable"; \
		echo "Example: make autotests DATABASE_URI='postgresql://postgres:postgres@localhost/praktikum?sslmode=disable'"; \
		exit 1; \
	fi
	./.tools/gophermarttest \
		-test.v -test.run=^TestGophermart$$ \
		-gophermart-binary-path=cmd/gophermart/gophermart \
		-gophermart-host=localhost \
		-gophermart-port=8080 \
		-gophermart-database-uri="$(DATABASE_URI)" \
		-accrual-binary-path=cmd/accrual/accrual_darwin_arm64 \
		-accrual-host=localhost \
		-accrual-port=$$((./.tools/random unused-port) | tr -d '%\n') \
		-accrual-database-uri="$(DATABASE_URI)"

start-services:
	@echo "Starting Docker Compose services..."
	docker-compose up -d postgres migrate
	@echo "Waiting for database to be ready..."
	@for i in {1..30}; do \
		if docker-compose exec -T postgres pg_isready -U gophermart -d gophermart -p 5433 >/dev/null 2>&1; then \
			echo "Database is ready!"; \
			break; \
		fi; \
		echo "Waiting for database... ($$i/30)"; \
		sleep 1; \
	done

autotests-docker: build start-services
	@echo "Running autotests with Docker Compose database..."
	./.tools/gophermarttest \
		-test.v -test.run=^TestGophermart$$ \
		-gophermart-binary-path=cmd/gophermart/gophermart \
		-gophermart-host=localhost \
		-gophermart-port=8080 \
		-gophermart-database-uri="postgres://gophermart:password@localhost:5433/gophermart?sslmode=disable" \
		-accrual-binary-path=cmd/accrual/accrual_darwin_arm64 \
		-accrual-host=localhost \
		-accrual-port=$$((./.tools/random unused-port) | tr -d '%\n') \
		-accrual-database-uri="postgres://gophermart:password@localhost:5433/gophermart?sslmode=disable"

ci-test-docker: statictest autotests-docker
	@echo "All CI tests with Docker completed successfully!"

ci-test: statictest autotests
	@echo "All CI tests completed successfully!"
