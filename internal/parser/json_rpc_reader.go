package parser

import (
	"context"
	"encoding/json"
	"io"
)

type JSONRPCReader struct {

}

func (reader *JSONRPCReader) Read(ctx context.Context, input io.Reader, requestChannel chan Request) {
	decoder := json.NewDecoder(input)
	
	for {
		var req Request
		if err := decoder.Decode(&req); err != nil {
			break
		}
		select {
		case <-ctx.Done():
			return
		case requestChannel <- req:
		}
	}
}