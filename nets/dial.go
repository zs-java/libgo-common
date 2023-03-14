package nets

import (
	"github.com/zs-java/libgo-common/nets/handshake"
	"github.com/zs-java/libgo-common/nets/packet"
	"net"
)

func Dial(network string, addr string, wrapper ConnWrapper) (*Conn, error) {
	c, err := net.Dial(network, addr)
	if err != nil {
		return nil, err
	}
	return wrapper(c), nil
}

func DialTCP(addr string, wrapper ConnWrapper) (*Conn, error) {
	return Dial("tcp", addr, wrapper)
}

func DialSimple(network string, addr string, version float64) (*Conn, error) {
	return Dial(network, addr, func(conn net.Conn) *Conn {
		return &Conn{
			Conn:             conn,
			HandshakeHandler: handshake.NewVersionClientHandler(version),
			Reader:           packet.NewReader(conn),
			Writer:           packet.NewWriter(conn),
		}
	})
}
