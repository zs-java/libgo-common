package packet

import (
	"fmt"
	"testing"
)

func TestBuf(t *testing.T) {
	var buf = []int32{1, 2, 3, 4}
	fmt.Println(buf[3:4])
}
