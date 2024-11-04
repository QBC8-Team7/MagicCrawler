app-linux:
	GOOS=linux GOARCH=amd64 go build -o build/app-linux-amd64 cmd/app/main.go

crawler-linux:
	GOOS=linux GOARCH=amd64 go build -o build/crawler-linux-amd64 cmd/crawler/main.go

test:
	go test ./...

