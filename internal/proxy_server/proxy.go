// ConnectionHandler defines the threading logic connecting receiving a request and processing it
package proxyserver

import (
	"context"
	"io"
	"sync"
	"mcp-gateway/internal/parser"
)

// Connection Handler configures how handle a user that connects to a client
// Creates routines to read from sockets, forward sockets and write back to sockets

type Task struct {
	req             parser.Request
	responseChannel chan parser.Response
}


type TaskReader interface {
	Read(context.Context, io.Reader, chan parser.Request)
}

type ResponseWriter interface {
	Write(context.Context, io.Writer, chan parser.Response)
}

type ProxyServer struct {
	PendingTasks chan Task
	NumWorkers   int
	Engine       Engine
	Reader       TaskReader
	Writer       ResponseWriter

	ctx    context.Context
	cancel context.CancelFunc
	group  sync.WaitGroup
}

func NewProxyServer(numWorkers int, engine Engine, reader TaskReader, writer ResponseWriter) *ProxyServer {
	pendingTasks := make(chan Task, 1000)
	ctx, cancel := context.WithCancel(context.Background())

	proxy := &ProxyServer{
		PendingTasks: pendingTasks,
		NumWorkers:   numWorkers,
		Engine:       engine,
		Reader:       reader,
		Writer:       writer,
		ctx:          ctx,
		cancel:       cancel,
	}

	proxy.group.Add(numWorkers)
	for i := 0; i < proxy.NumWorkers; i++ {
		go startWorker(ctx, proxy.Engine, proxy.PendingTasks)
	}

	return proxy
}

func startWorker(ctx context.Context, engine Engine, taskChannel chan Task) {
	for {
		select {
		case task := <-taskChannel:
			resp := engine.ProcessTask(task)
			task.responseChannel <- resp
		case <-ctx.Done():
			return
		}

	}
}

func (proxy *ProxyServer) HandleConnection(rw io.ReadWriter) {

	requestChannel := make(chan parser.Request, 50)
	responseChannel := make(chan parser.Response, 50)

	go proxy.Reader.Read(proxy.ctx, rw, requestChannel)
	go proxy.Writer.Write(proxy.ctx, rw, responseChannel)

	for req := range requestChannel {
		proxy.PendingTasks <- Task{
			req:             req,
			responseChannel: responseChannel,
		}
	}
}

func (proxy *ProxyServer) CloseConnection() {
	proxy.cancel()
	close(proxy.PendingTasks)
	proxy.group.Wait()
}
