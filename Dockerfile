FROM golang:latest as builder
WORKDIR /app
COPY go.mod .
COPY go.sum .
RUN go mod download
COPY . .
RUN go build -a -o ./app ./cmd/app/main.go
RUN go build -a -o ./crawler ./cmd/crawler/main.go

FROM ubuntu:latest
WORKDIR /app
COPY --from=builder /app/app .
COPY --from=builder /app/crawler .
COPY --from=builder /app/entrypoint.sh .
RUN chmod +x /app/entrypoint.sh
ENTRYPOINT ["/app/entrypoint.sh"]