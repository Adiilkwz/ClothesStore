FROM golang:1.25-alpine AS builder
WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -o /app/web_app ./cmd

FROM alpine:latest
WORKDIR /app

COPY --from=builder /app/web_app /app/web_app
COPY html /app/ui

EXPOSE 5050
CMD ["/app/web_app"]