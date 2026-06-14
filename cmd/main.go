package main

import (
	"fmt"
	"net"
	"mcp-gateway/internal/webserver"
)

type handler struct{}

func NewHandler() *handler {
	return &handler{}
}

func (h *handler) HandleConnection(conn net.Conn) {
	defer conn.Close()
	fmt.Println("Hello World")
}

func main() {
	h := NewHandler()
	server := webserver.NewTCPServer("127.0.0.1", "8000", h)
	server.Start()
}