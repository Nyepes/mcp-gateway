// ConnectionHandler defines the threading logic connecting receiving a request and processing it
package ConnectionHandler


import(
	"net"
)
// Connection Handler configures how handle a user that connects to a client
// Creates routines to read from sockets, forward sockets and write back to sockets
type Connection struct {
	ReadChannel chan []string
	WriteChannel chan []string
	

}



func (handler *Connection) HandleConnection(net.Conn) {
	// TODO: setup reader thread
	// TODO: setup writer thread


}