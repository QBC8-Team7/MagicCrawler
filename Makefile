.PHONY: app-linux crawler-linux migrate_up test run

APP_LINUX_BIN=build/app-linux-amd64
CRAWLER_LINUX_BIN=build/crawler-linux-amd64

app:
	GOOS=linux GOARCH=amd64 go build -o $(APP_LINUX_BIN) cmd/app/main.go

crawler:
	GOOS=linux GOARCH=amd64 go build -o build/crawler-linux-amd64 cmd/crawler/main.go

migrate_up:
	migrate -database postgres://postgres:postgres@localhost:5432/crawl_db?sslmode=disable -path migrations force 1 up

test:
	go test ./...

