FROM golang:1.22-alpine AS builder
RUN apk --no-cache add git ca-certificates
WORKDIR /app
COPY go.mod go.su* ./
RUN go mod download
COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -ldflags="-s -w" -o mcp-gateway ./src
FROM alpine:3.19 AS runner

WORKDIR /app
COPY --from=builder /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/
COPY --from=builder /app/mcp-gateway .
EXPOSE 8080

CMD ["./mcp-gateway"]