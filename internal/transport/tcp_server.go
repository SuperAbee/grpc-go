package transport

import (
	"context"
	"fmt"
	"net"

	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
)

func NewTcpServer(conn net.Conn) *tcpServer {
	fmt.Println("hello server")
	return &tcpServer{
		conn: conn,
	}
}

type tcpServer struct {
	conn net.Conn
}

// HandleStreams receives incoming streams using the given handler.
func (t *tcpServer) HandleStreams(handle func(*Stream), traceCtx func(context.Context, string) context.Context) {
	handle(&Stream{
		id:     1,
		st:     t,
		method: "test.Tester/MyTest",
		trReader: &transportReader{
			reader:        t.conn,
			windowHandler: func(int) {},
		},
		requestRead: func(int) {},
		ctx:         context.Background(),
	})
}

// WriteHeader sends the header metadata for the given stream.
// WriteHeader may not be called on all streams.
func (t *tcpServer) WriteHeader(s *Stream, md metadata.MD) error {
	return nil
}

// Write sends the data for the given stream.
// Write may not be called on all streams.
func (t *tcpServer) Write(s *Stream, hdr []byte, data []byte, opts *Options) error {
	fmt.Printf("tcp server ready write %v %v\n", hdr, data)
	n, err := t.conn.Write(append(hdr, data...))
	if err != nil {
		panic(err)
	}
	fmt.Printf("tcp server write %d bytes\n", n)
	return nil
}

// WriteStatus sends the status of a stream to the client.  WriteStatus is
// the final call made on a stream and always occurs.
func (t *tcpServer) WriteStatus(s *Stream, st *status.Status) error { return nil }

// Close tears down the transport. Once it is called, the transport
// should not be accessed any more. All the pending streams and their
// handlers will be terminated asynchronously.
func (t *tcpServer) Close(err error) {}

// RemoteAddr returns the remote network address.
func (t *tcpServer) RemoteAddr() net.Addr { return nil }

// Drain notifies the client this ServerTransport stops accepting new RPCs.
func (t *tcpServer) Drain(debugData string) {}

// IncrMsgSent increments the number of message sent through this transport.
func (t *tcpServer) IncrMsgSent() {}

// IncrMsgRecv increments the number of message received through this transport.
func (t *tcpServer) IncrMsgRecv() {}
