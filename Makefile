URL ?=

build:
	@go build -0 bin cmd/main.go

test:
	@go test ./...

dev-run:
	@go run cmd/main.go

migrations:
	@if [ -z "$(URL)" ]; then \
		echo "Running migrations without URL"; \
		dbmate up; \
	else \
		echo "Running migrations with URL: $(URL)"; \
		dbmate -u "$(URL)" up; \
	fi

seed-db:
	go run db/seed.go
