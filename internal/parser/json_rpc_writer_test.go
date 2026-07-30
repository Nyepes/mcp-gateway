package parser

import (
	"bytes"
	"context"
	"encoding/json"
	"testing"
	"time"
	"io"
)

func TestWriter_WriteResponses(t *testing.T) {
	writer := NewJSONRPCWriter()
	ctx, cancel := context.WithTimeout(context.Background(), time.Second * 2)
	defer cancel()

	pr, pw := io.Pipe()
	inChan := make(chan Response, 10)
	go func() {
		writer.Write(ctx, pw, inChan)
	}()

	res1 := Response{
		JSONRPC: "2.0",
		ID:      1,
		Result:  json.RawMessage(`{"status":"ok"}`),
	}
	inChan <- res1
	close(inChan)
	var actual Response
	err := json.NewDecoder(pr).Decode(&actual)
	if err != nil {
		t.Fatalf("writer produced invalid JSON output: %v", err)
	}

	if actual.JSONRPC != "2.0" {
		t.Errorf("expected jsonrpc '2.0', got '%s'", actual.JSONRPC)
	}
}

func TestWriter_ContextCancellation(t *testing.T) {
	var buf bytes.Buffer
	writer := NewJSONRPCWriter()

	ctx, cancel := context.WithCancel(context.Background())
	inChan := make(chan Response)

	done := make(chan struct{})
	go func() {
		writer.Write(ctx, &buf, inChan)
		close(done)
	}()

	cancel()
	close(inChan)

	select {
	case <-done:
		// Success
	case <-time.After(500 * time.Millisecond):
		t.Fatal("writer failed to exit on context cancellation")
	}
}
