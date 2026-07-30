package parser

import (
	"context"
	"encoding/json"
	"io"
)

type JSONRPCWriter struct{}

func NewJSONRPCWriter() *JSONRPCWriter {
	return &JSONRPCWriter{}
}

func (writer *JSONRPCWriter) Write(ctx context.Context, output io.Writer, responseChannel chan Response) {
	encoder := json.NewEncoder(output)
	for {
		select {
		case <-ctx.Done():
			return
		case response, ok := <-responseChannel:
			if !ok {
				return
			}
			err := encoder.Encode(response)
			if err != nil {
				return
			}
		}
	}
}
