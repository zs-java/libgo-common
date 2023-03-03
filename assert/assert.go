package assert

import "log"

func NotNil(val interface{}) {
	if val != nil {
		log.Panic("not nil", val)
	}
}

func GetFirstByteArr(b []byte, a ...any) []byte {
	return b
}

func GetFirst(a ...any) any {
	if len(a) > 0 {
		return a[0]
	}
	return nil
}
