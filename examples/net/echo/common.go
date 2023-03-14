package main

import (
	"fmt"
	"github.com/zs-java/libgo-common/nets"
	"github.com/zs-java/libgo-common/nets/packet"
	"os"
	"os/signal"
	"syscall"
)

const (
	CmdMsg int32 = iota + 1
	CmdExit
)

const (
	ServerVersion = 1.0
	ClientVersion = 1.0
	ServerAddr    = ":51000"
)

func notifySignalExit(conn *nets.Conn) {
	sigs := make(chan os.Signal, 1)
	signal.Notify(sigs, syscall.SIGINT)
	go func() {
		for {
			select {
			case c := <-sigs:
				if c == syscall.SIGINT {
					// send exit msg
					fmt.Println("\nSend Exit!")
					err := conn.WritePacket(packet.NewPacket(CmdExit, nil))
					if err != nil {
						panic(err)
					}
				}
			}
		}
	}()
}
