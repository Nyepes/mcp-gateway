package parser

import (
	"context"
	"io"
	"strings"
	"testing"
	"time"
)

func TestReader_SingleAndMultiLineJSON(t *testing.T) {
	rawStream := `{
		"jsonrpc": "2.0",
		"id": 1,
		"method": "initialize",
		"params": {"version": "1.0"}
	}` + "\n" + `{"jsonrpc": "2.0", "id": 2, "method": "tools/list"}`

	input := strings.NewReader(rawStream)
	reader := NewJSONRPCReader()

	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()

	outChan := make(chan Request, 10)
	go reader.Read(ctx, input, outChan)

	select {
	case req := <-outChan:
		if req.Method != "initialize" {
			t.Errorf("expected method 'initialize', got '%s'", req.Method)
		}
		if req.JSONRPC != "2.0" {
			t.Errorf("expected jsonrpc '2.0', got '%s'", req.JSONRPC)
		}
	case <-time.After(500 * time.Millisecond):
		t.Fatal("timed out waiting for first request")
	}


	select {
	case req := <-outChan:
		if req.Method != "tools/list" {
			t.Errorf("expected method 'tools/list', got '%s'", req.Method)
		}
	case <-time.After(500 * time.Millisecond):
		t.Fatal("timed out waiting for second request")
	}
}

func TestReader_ContextCancellation(t *testing.T) {
	// Infinite reader mock using io.Pipe
	r, w := io.Pipe()
	defer r.Close()

	reader := NewJSONRPCReader()
	ctx, cancel := context.WithCancel(context.Background())
	outChan := make(chan Request)
	done := make(chan struct{})
	
	go func() {
		reader.Read(ctx, r, outChan)
		close(done)
	}()
	cancel()
	_ = w.Close()

	select {
	case <-done:
		// Success
	case <-time.After(500 * time.Millisecond):
		t.Fatal("reader failed to exit on context cancellation")
	}
}