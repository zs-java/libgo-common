package packet

import (
	"io"
)

const (
	DefaultReaderBufLength = 16
)

type Reader struct {
	cache     []byte
	reader    io.Reader
	bufLength int
}

func NewReader(reader io.Reader) *Reader {
	return NewReaderSize(reader, DefaultReaderBufLength)
}

func NewReaderSize(reader io.Reader, bufLength int) *Reader {
	return &Reader{
		cache:     []byte{},
		reader:    reader,
		bufLength: bufLength,
	}
}

func (c *Reader) Read() (*Packet, error) {
	isHeader := true
	var bodyLength int32
	buf := make([]byte, 16)
	for {
		n, err := c.reader.Read(buf)
		if err != nil {
			return nil, err
		}
		// reset cache and append cache to data
		data := append(c.Reset(), buf[:n]...)
		length := int32(len(data))

		for {
			if isHeader {
				// 解析 bodyLength
				if length < HeaderLength {
					// 读取到的数据长度不足 HeaderLength 长度
					// 记录缓存，break 进入下一次读取
					c.cache = append(c.cache, data...)
					break
				} else if length == HeaderLength {
					// 数据长度正好等于 HeaderLength 长度
					// 解析 bodyLength 完成，break 进入下一次读取
					bodyLength = Byte2Int(data[:HeaderLength])
					isHeader = false
					break
				} else {
					// 数据长度大于 HeaderLength
					// 解析 bodyLength 并将剩余的数据放到 data 中后继续解析 body
					bodyLength = Byte2Int(data[:length])
					isHeader = false
					data = data[HeaderLength:]
					length = int32(len(data))
				}
			} else {
				// 解析 body
				if length < bodyLength {
					// 读取到的数据长度不足 body 长度
					// 记录缓存，break 进入下一次读取
					c.cache = append(c.cache, data...)
					break
				} else {
					// 数据长度大于等于 bodyLength 长度
					// 解析 body，返回数据 packet，将剩余的数据放到 cache 中
					packet := UnmarshalBody(data[:bodyLength])
					c.cache = append(c.cache, data[bodyLength:]...)
					return packet, nil
				}
			}
		}
	}
}

func (c *Reader) Reset() []byte {
	cache := c.cache
	c.cache = []byte{}
	return cache
}
