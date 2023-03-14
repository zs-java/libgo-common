package handshake

import (
	"encoding/binary"
	"io"
)

// ================ Client ===============

type VersionClientHandler struct {
	Version float64
}

func NewVersionClientHandler(version float64) *VersionClientHandler {
	return &VersionClientHandler{
		Version: version,
	}
}

func (t *VersionClientHandler) Handshake(conn io.ReadWriteCloser) error {
	// send version float64
	err := binary.Write(conn, binary.BigEndian, t.Version)
	if err != nil {
		return NewError(UnknownError, err.Error())
	}

	// read handshake resp (1byte)
	buf := make([]byte, 1)
	n, err := conn.Read(buf)
	if err != nil {
		return NewError(UnknownError, err.Error())
	}
	status := Status(buf[:n][0])
	if status != Success {
		return NewError(status, "")
	}
	return nil
}

// ================ Server ===============

type VersionServerHandler struct {
	Version float64
}

func NewVersionServerHandler(version float64) *VersionServerHandler {
	return &VersionServerHandler{Version: version}
}

func (t *VersionServerHandler) Handshake(conn io.ReadWriteCloser) error {
	var version float64
	err := binary.Read(conn, binary.BigEndian, &version)
	if err != nil {
		return respHandshakeStatus(UnknownError, conn)
	}
	// var status Status
	var status Status
	if version != t.Version {
		status = VersionNotMatch
	} else {
		status = Success
	}
	return respHandshakeStatus(status, conn)
}

func respHandshakeStatus(status Status, conn io.ReadWriteCloser) error {
	data := [1]byte{byte(status)}
	_, err := conn.Write(data[:])
	if status == Success && err == nil {
		return nil
	}
	// close conn
	_ = conn.Close()
	if err != nil {
		return NewError(status, err.Error())
	}
	return NewError(status, "")
}
