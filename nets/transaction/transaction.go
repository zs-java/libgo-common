package transaction

import "github.com/zs-java/libgo-common/nets/packet"

type Transaction interface {

	// GetId transactionId
	GetId() int32

	// GetPacket Get Request Packet
	GetPacket() *packet.Packet

	// GetCallbackPacket Get Callback Packet
	GetCallbackPacket() *packet.Packet

	// Wait Transaction has Callback Done
	Wait() error

	// ThenCallback New Goroutine to Wait and consumer callback packet
	ThenCallback(func(*packet.Packet))
}

type Manager interface {

	// CreateTransaction Create And Start Transaction
	CreateTransaction(cmd int32, data []byte) Transaction

	// CreateCallbackPacket with callbackCmd……
	CreateCallbackPacket(transactionId int32, data []byte) *packet.Packet

	IsCallbackPacket(p *packet.Packet) bool

	// DoneTransaction Close Transaction Wait
	DoneTransaction(p *packet.Packet) error
}
