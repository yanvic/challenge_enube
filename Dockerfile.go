FROM golang:1.25.1-alpine AS builder

RUN apk add --no-cache git ca-certificates

WORKDIR /app

COPY back/go.mod back/go.sum ./
RUN go mod download

COPY back/cmd ./cmd
COPY back/internal ./internal

RUN go build -o app ./cmd/api

FROM alpine:3.18

RUN apk add --no-cache ca-certificates

WORKDIR /app

COPY --from=builder /app/app .

EXPOSE 8080

CMD ["./app"]
