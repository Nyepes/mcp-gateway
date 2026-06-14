package webserver

import (
	"net"
	"testing"
	"time"
)

type mockHandler struct {
	wasCalled chan bool
}

func (handler *mockHandler) HandleConnection(conn net.Conn) {
	defer conn.Close()
	handler.wasCalled <- true
}

func TestTCPServer(t *testing.T) {
	handlerSignal := make(chan bool, 1)
	serverReadySignal := make(chan struct{}, 1)
	handler := &mockHandler{wasCalled: handlerSignal}
	address := "127.0.0.1:0"
	server := NewTCPServer(address, handler, serverReadySignal)
	
	go server.Start()
	select {
	case <- serverReadySignal:
		address = server.Address
	case <- time.After(1 * time.Second):
		t.Fatal("Timeout starting server")
	}

	conn, err := net.Dial("tcp", address)
	if err != nil {
		t.Fatal(err)
	}
	conn.Close()
	
	select {
		case <- handlerSignal:
			t.Log("Succcess")
		case <- time.After(1 * time.Second):
			t.Error("Timeout")
	}
}