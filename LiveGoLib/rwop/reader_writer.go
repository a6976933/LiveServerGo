package rwop

import (
	"bufio"
	"io"
	"net"
)

const (
	BUFFER_SIZE = 8192
)

type NetReadWriter struct {
	rw       io.ReadWriter
	readErr  error
	writeErr error
	buf      []byte
}

func NewNetReadWriter(nrw net.Conn) *NetReadWriter {
	ret := &NetReadWriter{
		rw: bufio.NewReadWriter(bufio.NewReaderSize(nrw, BUFFER_SIZE), bufio.NewWriterSize(nrw, BUFFER_SIZE)),
	}
	return ret
}
func (nrw *NetReadWriter) ReadBufLen() ([]byte, error) {
	var b = make([]byte, BUFFER_SIZE)
	_, err := nrw.rw.Read(b)
	return b, err
}
