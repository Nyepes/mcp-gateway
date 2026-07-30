package parser

import (
	"context"
	"io"
)

type JSONRPCWriter struct {

}

func (writer *JSONRPCWriter) Write(ctx context.Context, output io.Writer, responseChannel chan Response) {
	responseChannel <- Response{
	}
}