package main

import (
	"bytes"
	"encoding/binary"
	"encoding/json"
	"fmt"
	"github.com/zs-java/libgo-common/nets"
	"github.com/zs-java/libgo-common/nets/handshake"
	"github.com/zs-java/libgo-common/nets/packet"
	"github.com/zs-java/libgo-common/nets/transaction"
	"net"
	"sync"
)

func main() {
	conn, err := nets.Dial("tcp", ServerAddr, func(conn net.Conn) *nets.Conn {
		return &nets.Conn{
			Conn:               conn,
			HandshakeHandler:   handshake.NewVersionClientHandler(ClientVersion),
			Reader:             packet.NewReader(conn),
			Writer:             packet.NewWriter(conn),
			PacketHandlers:     []*nets.PacketHandler{nets.DefaultTransactionPacketHandler},
			TransactionManager: transaction.NewManager(),
		}
	})
	if err != nil {
		panic(err)
	}

	_ = conn.Handshake()
	go conn.HandlePacket()

	var wg sync.WaitGroup

	calcFunc := func(num1, num2 int32) {
		data, err := json.Marshal(CalcBody{Num1: num1, Num2: num2})
		if err != nil {
			panic(err)
		}
		t := conn.TransactionManager.CreateTransaction(CmdCalc, data)
		_ = conn.WritePacket(t.GetPacket())
		t.ThenCallback(func(p *packet.Packet) {
			defer wg.Done()
			var result int32
			err = binary.Read(bytes.NewBuffer(p.Data), binary.BigEndian, &result)
			if err != nil {
				panic(err)
			}
			fmt.Printf("%d * %d = %d\n", num1, num2, result)
		})
	}

	count := 50
	for i := 0; i < count; i++ {
		wg.Add(1)
		go calcFunc(int32(i), int32(i*2))
	}

	wg.Wait()

	fmt.Println("Done!")

}
