package av

import (
	"net"

	"LiveGoLib/rwop"
)

type Object map[string]interface{}

type ClientConn struct {
	Conn net.Conn
	Buff []byte
	ChunkHeader
	Nrw *rwop.NetReadWriter
}

type ChunkMessageHeader struct {
	Timestamp       uint32
	MessageLength   uint32
	MessageTypeID   uint8
	MessageStreamID uint32
}

type ChunkBasicHeader struct {
	Fmt  uint8
	Csid uint32
}

type ChunkExtendedTimestamp struct {
	ExtendTimestamp []byte
}

type ChunkHeader struct {
	ChunkBasicHeader
	ChunkMessageHeader
	ChunkExtendedTimestamp
}
