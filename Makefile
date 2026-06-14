run:
	go run ./cmd/main.go

compile:
	CGO_ENABLED=0 GOOS=darwin GOARCH=arm64 go build -ldflags="-s -w" -o mcp-gateway ./cmd/main.go