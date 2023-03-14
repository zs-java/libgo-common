package packet

// Packet Network Transfer Packet
type Packet struct {
	Cmd           int32
	TransactionId int32
	Data          []byte
}

func NewPacket(cmd int32, data []byte) *Packet {
	return &Packet{
		Cmd:  cmd,
		Data: data,
	}
}

func NewPacketTransaction(cmd, transactionId int32, data []byte) *Packet {
	return &Packet{
		Cmd:           cmd,
		TransactionId: transactionId,
		Data:          data,
	}
}
