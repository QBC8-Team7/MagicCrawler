.PHONY: app-linux crawler-linux migrate-up migrate-down test run

APP_LINUX_BIN=build/app-linux-amd64
CRAWLER_LINUX_BIN=build/crawler-linux-amd64

app:
	GOOS=linux GOARCH=amd64 go build -o $(APP_LINUX_BIN) cmd/app/main.go

crawler:
	GOOS=linux GOARCH=amd64 go build -o build/crawler-linux-amd64 cmd/crawler/main.go

migrate-up:
	migrate -database "postgres://postgres:postgres@localhost:5432/magic-crawler?sslmode=disable" -path ./pkg/db/migration up

migrate-down:
	migrate -database "postgres://postgres:postgres@localhost:5432/magic-crawler?sslmode=disable" -path ./pkg/db/migration down

sqlc-check:
	sqlc compile -f config/sqlc.yml

sqlc-gen:
	sqlc generate -f config/sqlc.yml

test:
	go test ./...

