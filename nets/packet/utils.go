package packet

import (
	"bytes"
	"encoding/binary"
	"log"
)

func Int2Byte(data int32) []byte {
	buf := bytes.NewBuffer([]byte{})
	err := binary.Write(buf, binary.BigEndian, data)
	if err != nil {
		log.Panic(err)
	}
	return buf.Bytes()
}

func Byte2Int(data []byte) (ret int32) {
	err := binary.Read(bytes.NewBuffer(data), binary.BigEndian, &ret)
	if err != nil {
		log.Panic(err)
	}
	return ret
}
