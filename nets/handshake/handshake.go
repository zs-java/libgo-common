package handshake

import (
	"errors"
	"fmt"
	"io"
)

// Handler HandshakeHandler
type Handler interface {
	Handshake(conn io.ReadWriteCloser) error
}

// Status HandshakeStatus
type Status byte

const (
	UnknownError    Status = 0
	Success         Status = 1
	VersionNotMatch Status = 2
)

func NewError(status Status, message string) error {
	return errors.New(fmt.Sprintf("Handshake Handler Status is %d, Message: %v", int(status), message))
}
