.PHONY: build test run docker-up docker-down run-autotests lint local-tests

build:
	@echo "Building gophermart..."
	(cd cmd/gophermart && go build -buildvcs=false -o gophermart)
	@echo "Build complete: cmd/gophermart/gophermart"

test:
	@echo "Running unit tests..."
	go test ./...

run:
	air

docker-up:
	docker-compose up --build

docker-down:
	docker-compose down

run-autotests: build
	echo "Running autotests with local binaries..."; \
	./.tools/gophermarttest \
	-test.v -test.run=^TestGophermart$$ \
		  -gophermart-binary-path=cmd/gophermart/gophermart \
		  -accrual-binary-path=cmd/accrual/accrual_darwin_arm64 \
		  -gophermart-host=localhost \
		  -gophermart-port=8080 \
		  -gophermart-database-uri="postgres://gophermart:password@localhost:5433/gophermart?sslmode=disable" \
		  -accrual-host=localhost \
		  -accrual-port=8084 \
		  -accrual-database-uri="postgres://gophermart:password@localhost:5433/gophermart?sslmode=disable";

local-tests:
	@echo "Running local unit tests..."
	go test ./...

lint: build
	@echo "Running linter..."
	go vet -vettool=$(PWD)/.tools/statictest ./...
