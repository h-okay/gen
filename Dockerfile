FROM golang:1.22-alpine AS builder

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN CGO_ENABLED=0 GOOS=linux go build -o /main

FROM alpine:latest

WORKDIR /root/

COPY --from=builder /main .
COPY --from=builder /app/internal/db/migrations ./migrations

EXPOSE 8080

CMD ["./main"]