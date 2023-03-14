package packet

import "io"

type Writer struct {
	writer io.Writer
}

func NewWriter(writer io.Writer) *Writer {
	return &Writer{writer: writer}
}

func (c *Writer) Write(packet *Packet) (err error) {
	_, err = c.writer.Write(Marshal(packet))
	return err
}
