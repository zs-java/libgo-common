package nets

import (
	"github.com/zs-java/libgo-common/nets/handshake"
	"github.com/zs-java/libgo-common/nets/packet"
	"net"
)

type Listener struct {
	listener    net.Listener
	ConnWrapper ConnWrapper
}

func (t *Listener) Accept() (*Conn, error) {
	conn, err := t.listener.Accept()
	if err != nil {
		return nil, err
	}
	return t.ConnWrapper(conn), nil
}

func (t *Listener) Close() error {
	return t.listener.Close()
}

func (t *Listener) Addr() net.Addr {
	return t.listener.Addr()
}

func Listen(network string, addr string, wrapper ConnWrapper) (*Listener, error) {
	listen, err := net.Listen(network, addr)
	if err != nil {
		return nil, err
	}
	return &Listener{
		listener:    listen,
		ConnWrapper: wrapper,
	}, nil
}

func ListenTCP(addr string, wrapper ConnWrapper) (*Listener, error) {
	return Listen("tcp", addr, wrapper)
}

func ListenSimple(network string, addr string, version float64) (*Listener, error) {
	return Listen(network, addr, func(conn net.Conn) *Conn {
		return &Conn{
			Conn:             conn,
			HandshakeHandler: handshake.NewVersionServerHandler(version),
			Reader:           packet.NewReader(conn),
			Writer:           packet.NewWriter(conn),
		}
	})
}
