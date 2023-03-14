package utils

import (
	"bytes"
	"encoding/binary"
	"log"
	"os"
	"os/exec"
)

func Int2Byte(data int32) []byte {
	buf := bytes.NewBuffer([]byte{})
	err := binary.Write(buf, binary.BigEndian, data)
	if err != nil {
		log.Panic(err)
	}
	return buf.Bytes()

	// var len uintptr = unsafe.Sizeof(data)
	// ret = make([]byte, len)
	// var tmp int = 0xff
	// var index uint = 0
	// for index=0; index<uint(len); index++{
	// 	ret[index] = byte((tmp<<(index*8) & data)>>(index*8))
	// }
	// return ret
}

func Byte2Int(data []byte) (ret int32) {
	err := binary.Read(bytes.NewBuffer(data), binary.BigEndian, &ret)
	if err != nil {
		log.Panic(err)
	}

	return ret

	// var ret int = 0
	// var len int = len(data)
	// var i uint = 0
	// for i=0; i<uint(len); i++{
	// 	ret = ret | (int(data[i]) << (i*8))
	// }
	// return ret
}

func ExecCmd(c string, args ...string) {
	log.Printf("exec cmd: %v %v:", c, args)
	cmd := exec.Command(c, args...)
	cmd.Stderr = os.Stderr
	cmd.Stdout = os.Stdout
	cmd.Stdin = os.Stdin
	err := cmd.Run()
	if err != nil {
		log.Println("failed to exec cmd:", err)
	}
}
