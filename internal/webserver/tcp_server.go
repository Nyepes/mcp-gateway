// webserver creates webserver to accept incoming TCP connections
package webserver

import "net"

// ConnectionHandler defines how to handle a connection after server accepts it
// It is responsible for cleaning up connection after it is done
type ConnectionHandler interface {
	HandleConnection(connection net.Conn)
}

// TCPServer configures TCP server
type TCPServer struct {
	Address string
	Handler ConnectionHandler
	Ready	chan struct{}
}

// NewTCPServer creates a TCPServer Object
func NewTCPServer(address string, handler ConnectionHandler, ready_channel chan struct{}) *TCPServer {
	return &TCPServer{
		Address: address,
		Handler: handler,
		Ready:	ready_channel,
	}
}

// Start the TCP server, it listens and accepts connections
func (server *TCPServer) Start() {
	listener, err := net.Listen("tcp", server.Address)
	if err != nil {
		panic(err)
	}
	server.Address = listener.Addr().String()
	// Server notifies that it has started
	defer listener.Close()
	if server.Ready != nil {
		close(server.Ready)
	}
	for {
		connection, err := listener.Accept()
		if err != nil {
			// Ignore
			continue
		}
		go server.Handler.HandleConnection(connection)
	}

}