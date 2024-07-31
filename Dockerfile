# build
FROM golang:1.22.5-alpine AS builder
WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY main.go .
RUN go build -o main .
RUN adduser --disabled-password appuser
USER appuser

# runtime
FROM alpine:latest
COPY --from=builder /app/main /app/
WORKDIR /app
EXPOSE 4444
CMD ["./main"]
