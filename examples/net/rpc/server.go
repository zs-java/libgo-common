package main

import (
	"bytes"
	"encoding/binary"
	"encoding/json"
	"github.com/zs-java/libgo-common/nets"
	"github.com/zs-java/libgo-common/nets/handshake"
	"github.com/zs-java/libgo-common/nets/packet"
	"github.com/zs-java/libgo-common/nets/transaction"
	"io"
	"log"
	"net"
	"time"
)

func NewCalcPacketHandler() *nets.PacketHandler {
	return &nets.PacketHandler{
		Name: "CalcPacketHandler",
		Cmd:  CmdCalc,
		Action: func(ctx *nets.Context, p *packet.Packet) error {
			// 解析请求
			var body CalcBody
			err := json.Unmarshal(p.Data, &body)
			if err != nil {
				return err
			}
			// 计算
			result := body.Num1 * body.Num2
			buf := bytes.NewBuffer([]byte{})
			_ = binary.Write(buf, binary.BigEndian, result)
			// 响应事务
			callbackPacket := ctx.TransactionManager.CreateCallbackPacket(p.TransactionId, buf.Bytes())
			time.Sleep(time.Second)
			return ctx.WritePacket(callbackPacket)
		},
	}
}

func main() {
	listen, err := nets.Listen("tcp", ServerAddr, func(conn net.Conn) *nets.Conn {
		return &nets.Conn{
			Conn:               conn,
			HandshakeHandler:   handshake.NewVersionServerHandler(ServerVersion),
			Reader:             packet.NewReader(conn),
			Writer:             packet.NewWriter(conn),
			TransactionManager: transaction.NewManager(),
			PacketHandlers: []*nets.PacketHandler{
				nets.DefaultTransactionPacketHandler,
				NewCalcPacketHandler(),
			},
			ExceptionHandler: func(ctx *nets.Context, err error) {
				if err != io.EOF {
					log.Fatalln(err)
				}
			},
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
