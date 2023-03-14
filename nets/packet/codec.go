package packet

import (
	"bytes"
	"encoding/binary"
)

// HeaderLength int32 4byte
const HeaderLength int32 = 4

// Marshal Header + Body
// Header is Body length, int32, 4byte
// Body see MarshalBody
func Marshal(packet *Packet) []byte {
	// MarshalBody
	body := MarshalBody(packet)
	// header = int32 body length
	header := Int2Byte(int32(len(body)))
	return append(header, body...)
}

// MarshalBody encode Packet to []byte
// format: cmd(int32) + transactionId(int32) + data([]byte)
func MarshalBody(packet *Packet) []byte {
	var buf []byte
	buf = append(buf, Int2Byte(packet.Cmd)...)
	buf = append(buf, Int2Byte(packet.TransactionId)...)
	if packet.Data != nil {
		buf = append(buf, packet.Data...)
	}
	return buf[:]
}

// UnmarshalBody decode []byte to Packet
func UnmarshalBody(buf []byte) *Packet {
	reader := bytes.NewReader(buf)

	var cmd, transactionId int32
	// read int32(4byte) Cmd And TransactionId
	_ = binary.Read(reader, binary.BigEndian, &cmd)
	_ = binary.Read(reader, binary.BigEndian, &transactionId)

	// read other bytes
	data := make([]byte, reader.Len())
	n, _ := reader.Read(data)
	return &Packet{
		Cmd:           cmd,
		TransactionId: transactionId,
		Data:          data[:n],
	}
}
