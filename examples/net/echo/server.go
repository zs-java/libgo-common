package main

import (
	"fmt"
	"github.com/zs-java/libgo-common/nets"
	"github.com/zs-java/libgo-common/nets/handshake"
	"github.com/zs-java/libgo-common/nets/packet"
	"github.com/zs-java/libgo-common/nets/transaction"
	"net"
)

var serverExitHandler = &nets.PacketHandler{
	Name: "ServerExitHandler",
	Cmd:  CmdExit,
	Action: func(ctx *nets.Context, p *packet.Packet) error {
		fmt.Println("exit")
		return ctx.Close()
	},
}

func NewServerMsgHandler(storage *nets.SimpleConnStorage) *nets.PacketHandler {
	return &nets.PacketHandler{
		Name: "ServerMsgHandler",
		Cmd:  CmdMsg,
		Action: func(ctx *nets.Context, p *packet.Packet) error {
			fmt.Printf("receive msg: %s\n", string(p.Data))
			for conn, _ := range storage.ConnMap {
				_ = conn.WritePacket(p)
			}
			return nil
		},
	}
}

func main() {

	storage := nets.NewSimpleConnStorage(nil)

	listen, err := nets.ListenTCP(ServerAddr, func(conn net.Conn) *nets.Conn {
		return &nets.Conn{
			Conn:             conn,
			HandshakeHandler: handshake.NewVersionServerHandler(ServerVersion),
			Reader:           packet.NewReader(conn),
			Writer:           packet.NewWriter(conn),
			PacketHandlers: []*nets.PacketHandler{
				nets.DefaultTransactionPacketHandler, NewServerMsgHandler(storage), serverExitHandler,
			},
			TransactionManager: transaction.NewManager(),
			ExceptionHandler: func(ctx *nets.Context, err error) {
				panic(err)
			},
			HandshakeDoneHandler: storage.DefaultHandshakeDoneHandler(),
			ClosedHandler:        storage.DefaultClosedHandler(),
		}
	})
	if err != nil {
		panic(err)
	}
	for {
		conn, _ := listen.Accept()
		go func() {
			_ = conn.Handshake()
			conn.HandlePacket()
		}()
	}

}
