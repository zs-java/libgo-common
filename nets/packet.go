package nets

import (
	"errors"
	"github.com/zs-java/libgo-common/nets/packet"
	"github.com/zs-java/libgo-common/nets/transaction"
)

type PacketActionFunc func(ctx *Context, p *packet.Packet) error

type PacketHandler struct {
	Name   string
	Cmd    int32
	Action PacketActionFunc
}

var DefaultTransactionPacketAction PacketActionFunc = func(ctx *Context, p *packet.Packet) error {
	tm := ctx.TransactionManager
	if tm == nil {
		return errors.New("transactionManager is nil")
	}
	if tm.IsCallbackPacket(p) {
		_ = tm.DoneTransaction(p)
	}
	return nil
}

var DefaultTransactionPacketHandler = &PacketHandler{
	Name:   "TransactionCallbackPacketHandler",
	Cmd:    transaction.DefaultCallbackCmd,
	Action: DefaultTransactionPacketAction,
}
