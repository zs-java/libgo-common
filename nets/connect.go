package nets

import (
	"github.com/zs-java/libgo-common/nets/handshake"
	"github.com/zs-java/libgo-common/nets/packet"
	"github.com/zs-java/libgo-common/nets/transaction"
	"net"
	"time"
)

type ExceptionHandler func(ctx *Context, err error)
type ClosedHandler func(ctx *Context)
type HandshakeDoneHandler func(ctx *Context)

type Conn struct {
	net.Conn
	HandshakeHandler     handshake.Handler
	Reader               *packet.Reader
	Writer               *packet.Writer
	PacketHandlers       []*PacketHandler
	TransactionManager   transaction.Manager
	Attributes           map[string]interface{}
	ExceptionHandler     ExceptionHandler
	ClosedHandler        ClosedHandler
	HandshakeDoneHandler HandshakeDoneHandler
	closing              bool
	closed               bool
}

func (c *Conn) ReadPacket() (*packet.Packet, error) {
	p, err := c.Reader.Read()
	if c.ExceptionHandler != nil {
		c.handleException(err)
		return p, err
	}
	return p, err
}

func (c *Conn) WritePacket(p *packet.Packet) error {
	return c.Writer.Write(p)
}

func (c *Conn) Handshake() error {
	err := c.HandshakeHandler.Handshake(c)
	if err != nil {
		if c.ExceptionHandler != nil {
			c.handleException(err)
		}
	} else {
		if c.HandshakeDoneHandler != nil {
			c.HandshakeDoneHandler(c.wrapperContext())
		}
	}
	return err
}

func (c *Conn) Close() error {
	c.closing = true
	return c.SetReadDeadline(time.Now())
}

func (c *Conn) HandlePacket() {
	defer c.handleRecoverException()
	for {
		p, err := c.ReadPacket()

		if err != nil {
			panic(err)
		}
		go func() {
			for _, handler := range c.PacketHandlers {
				if p.Cmd == handler.Cmd {
					if c.Attributes == nil {
						c.Attributes = make(map[string]interface{})
					}
					err = handler.Action(c.wrapperContext(), p)
					if err != nil {
						// c.handleException(err)
						panic(err)
					}
				}
			}
		}()
	}
}

func (c *Conn) wrapperContext() *Context {
	return &Context{c}
}

func (c *Conn) handleException(err any) {
	if err == nil || c.closed {
		return
	}
	defer func() {
		c.closed = true
		_ = c.Conn.Close()
		if c.ClosedHandler != nil {
			c.ClosedHandler(c.wrapperContext())
		}
	}()
	if netErr, ok := err.(net.Error); ok && netErr.Timeout() && c.closing {
		return
	}
	if e, ok := err.(error); ok && c.ExceptionHandler != nil {
		c.ExceptionHandler(&Context{c}, e)
	} else {
		panic(err)
	}
}

func (c *Conn) handleRecoverException() {
	c.handleException(recover())
}

type ConnWrapper func(conn net.Conn) *Conn
