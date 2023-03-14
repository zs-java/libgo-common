package main

import (
	"bufio"
	"fmt"
	"github.com/zs-java/libgo-common/nets"
	"github.com/zs-java/libgo-common/nets/handshake"
	"github.com/zs-java/libgo-common/nets/packet"
	"io"
	"net"
	"os"
)

var clientMsgHandler = &nets.PacketHandler{
	Name: "ClientMsgHandler",
	Cmd:  CmdMsg,
	Action: func(ctx *nets.Context, packet *packet.Packet) error {
		countInfc, ok := ctx.Attributes["count"]
		var count int64
		if ok {
			count = countInfc.(int64) + 1
		} else {
			count = 1
		}
		ctx.Attributes["count"] = count
		fmt.Printf("echo%d: %s\n", count, string(packet.Data))
		return nil
	},
}

func main() {
	conn, err := nets.Dial("tcp", ServerAddr, func(conn net.Conn) *nets.Conn {
		return &nets.Conn{
			Conn:             conn,
			HandshakeHandler: handshake.NewVersionClientHandler(ClientVersion),
			Reader:           packet.NewReader(conn),
			Writer:           packet.NewWriter(conn),
			PacketHandlers: []*nets.PacketHandler{
				clientMsgHandler,
			},
			ExceptionHandler: func(ctx *nets.Context, err error) {
				if err == io.EOF {
					fmt.Println("disconnect...")
					os.Exit(0)
				} else {
					panic(err)
				}
			},
		}
	})
	if err != nil {
		panic(err)
	}

	// do handshake
	_ = conn.Handshake()
	fmt.Println("Handshake Done.")

	go conn.HandlePacket()

	go notifySignalExit(conn)

	stdinReader := bufio.NewReader(os.Stdin)
	for {
		data, _, _ := stdinReader.ReadLine()
		_ = conn.WritePacket(packet.NewPacket(CmdMsg, data))
	}
}
