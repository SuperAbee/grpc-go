package transport

import (
	"context"
	"fmt"
	"net"

	"google.golang.org/grpc/resolver"
)

func NewTcpClient(addr resolver.Address, opts ConnectOptions) *tcpClient {
	fmt.Printf("hello tcp client %s\n", addr.Addr)
	var c net.Conn
	var err error
	if opts.Dialer != nil {
		c, err = opts.Dialer(context.Background(), addr.Addr)
	} else {
		c, err = net.Dial("tcp", addr.Addr)
	}

	if err != nil {
		panic(err)
	}
	return &tcpClient{
		streams: make(map[uint32]*Stream),
		conn:    c,
	}
}

type tcpClient struct {
	streams map[uint32]*Stream
	conn    net.Conn
}

func (t *tcpClient) Close(err error) {

}

// GracefulClose starts to tear down the transport: the transport will stop
// accepting new RPCs and NewStream will return error. Once all streams are
// finished, the transport will close.
//
// It does not block.
func (t *tcpClient) GracefulClose() {

}

// Write sends the data for the given stream. A nil stream indicates
// the write is to be performed on the transport as a whole.
func (t *tcpClient) Write(s *Stream, hdr []byte, data []byte, opts *Options) error {
	fmt.Printf("tcp client ready write %v %v\n", hdr, data)
	n, err := t.conn.Write(append(hdr, data...))
	if err != nil {
		panic(err)
	}
	fmt.Printf("tcp server write %d bytes\n", n)
	return nil
}

// NewStream creates a Stream for an RPC.
func (t *tcpClient) NewStream(ctx context.Context, callHdr *CallHdr) (*Stream, error) {
	id := uint32(1)
	s := Stream{
		id:     id,
		ct:     t,
		method: callHdr.Method,
		trReader: &transportReader{
			reader:        t.conn,
			windowHandler: func(int) {},
		},
		requestRead: func(i int) {},
	}
	t.streams[id] = &s
	return &s, nil
}

// CloseStream clears the footprint of a stream when the stream is
// not needed any more. The err indicates the error incurred when
// CloseStream is called. Must be called when a stream is finished
// unless the associated transport is closing.
func (t *tcpClient) CloseStream(stream *Stream, err error) {

}

// Error returns a channel that is closed when some I/O error
// happens. Typically the caller should have a goroutine to monitor
// this in order to take action (e.g., close the current transport
// and create a new one) in error case. It should not return nil
// once the transport is initiated.
func (t *tcpClient) Error() <-chan struct{} {
	return make(<-chan struct{}, 10)
}

// GoAway returns a channel that is closed when ClientTransport
// receives the draining signal from the server (e.g., GOAWAY frame in
// HTTP/2).
func (t *tcpClient) GoAway() <-chan struct{} {
	return make(<-chan struct{}, 10)
}

// GetGoAwayReason returns the reason why GoAway frame was received, along
// with a human readable string with debug info.
func (t *tcpClient) GetGoAwayReason() (GoAwayReason, string) {
	return 0, ""
}

// RemoteAddr returns the remote network address.
func (t *tcpClient) RemoteAddr() net.Addr {
	return t.conn.RemoteAddr()
}

// IncrMsgSent increments the number of message sent through this transport.
func (t *tcpClient) IncrMsgSent() {
}

// IncrMsgRecv increments the number of message received through this transport.
func (t *tcpClient) IncrMsgRecv() {

}
