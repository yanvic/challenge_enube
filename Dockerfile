FROM golang:1.25.1-alpine

RUN apk add --no-cache git ca-certificates

WORKDIR /app

COPY back/go.mod back/go.sum ./
RUN go mod download

COPY back/ ./
RUN go build -o app ./cmd/api

EXPOSE 8080

CMD ["./app"]
