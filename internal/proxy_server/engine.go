package proxyserver

import (
	"mcp-gateway/internal/parser"
)

type Engine struct {
}

func (engine *Engine) ProcessTask(Task) parser.Response {
	return parser.Response{}
}
