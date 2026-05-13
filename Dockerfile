FROM golang:1.22-alpine AS builder

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o /ela-api ./cmd/api

FROM alpine:3.20

WORKDIR /app

RUN addgroup -S app && adduser -S app -G app

COPY --from=builder /ela-api /usr/local/bin/ela-api
COPY --from=builder /app/configs /app/configs

USER app

EXPOSE 8080

ENTRYPOINT ["ela-api"]
