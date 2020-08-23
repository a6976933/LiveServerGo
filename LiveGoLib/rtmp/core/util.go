package core

import (
	"bufio"
	"errors"
	"fmt"
	"io"
	"net"
	"os"
	"LiveGoLib/av"
	"LiveGoLib/proxy"
)

const (
	BUFFER_SIZE = 4096
)

type ClientConn interface {
	GetNetConn() net.Conn
}

type BReadWriter interface {
	ReadNByte(n int, buff []byte) ([]byte, error)
	ReadBuff() ([]byte, int, error) //Read Buff Length
	Write([]byte) (int, error)
}

type ClientInfo struct {
	ChunkHeader
	DataFrame
	AVCType
	onMetadata     []byte
	avcSpec        bool
	aacSpec        bool
	prox           proxy.MiddleBuf
	Conn           net.Conn
	Buff           []byte
	rw             BReadWriter
	chunksize      uint32
	windowAckSize  uint32
	obj            av.Object
	transactionID  float64
	publishName    string
	publishType    string
	dataFrame      map[string]interface{}
	videoFrameType uint32
	media         *os.File
	disConn       bool
}

type DataFrame struct {
	audioCodecID    uint32
	audioDataRate   float64
	audioSampleRate uint32
	audioSampleSize uint32
	compatibleBrand string
	duration        float64
	encoder         string
	fileSize        float64
	frameRate       float64
	height          uint32
	width           uint32
	majorBrand      string
	minorVersion    string
	stereo          bool
	videoCodecID    uint32
	videoDataRate   float64
	title           string
	soundFormat     uint32
	soundRate       uint32
	soundSize       uint32
	soundType       uint32
}

type AVCType struct {
	avcPacketType   uint8
	compositionTime uint32
}

type ByteReadWriter struct {
	brw      *bufio.ReadWriter
	readerr  error
	writeerr error
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

func NewClientInfo(Conn net.Conn, buff []byte) *ClientInfo {
	ret := &ClientInfo{
		Conn:      Conn,
		Buff:      buff,
		rw:        NewByteReadWriter(Conn),
		chunksize: 128,
	}
	return ret
}

func NewChunkBasicHeader(fMt uint8, csid uint32) *ChunkBasicHeader {
	ret := &ChunkBasicHeader{
		Fmt:  fMt,
		Csid: csid,
	}
	return ret
}

func NewChunkMessageHeader() *ChunkMessageHeader {
	ret := &ChunkMessageHeader{}
	return ret
}

func NewByteReadWriter(r io.ReadWriter) *ByteReadWriter {
	ret := &ByteReadWriter{
		brw: bufio.NewReadWriter(bufio.NewReaderSize(r, BUFFER_SIZE), bufio.NewWriterSize(r, BUFFER_SIZE)),
	}
	return ret
}

func (brw *ByteReadWriter) ReadNByte(n int, b []byte) ([]byte, error) {
	b = make([]byte, n)
	var haveRead int
	haveRead, err := brw.brw.Read(b)
	if err != nil {
		return b, err
	}
	//var count = 0
	fmt.Println("Now return buff size: ", len(b))
	for haveRead < n {
		others := make([]byte, n-haveRead)
		d, err := brw.brw.Read(others)
		if err != nil {
			return b, err
		}
		fmt.Println("Now d : ", d," Now HaveRead : ", haveRead)
		if d != 0 {
			for i := 0;i < d; i++ {
				b[haveRead] = others[i]
				haveRead++
			}
		}
		fmt.Println("Now return read again buff size : ", len(b))
	}
	fmt.Println(haveRead)
	if haveRead != n {
		fmt.Println("Read over Buffer", haveRead, "bytes")
		return b, errors.New("Read over Buffer" + string(haveRead) + "bytes")
	}
	return b, nil
}

func (brw *ByteReadWriter) ReadBuff() ([]byte, int, error) {
	b := make([]byte, BUFFER_SIZE)
	n, err := brw.brw.Read(b)
	if err != nil {
		return b, n, err
	}
	return b, n, nil
}

func (brw *ByteReadWriter) Write(b []byte) (int, error) {
	n, err := brw.brw.Write(b)
	brw.brw.Flush()
	if err != nil {
		return n, err
	}
	return n, nil
}




/*
func (brw *ByteReadWriter) ReadFull() ([]byte, int, error) {
	var b = make([]byte, BUFFER_SIZE)
	var n int
	var err error
	n, err = io.ReadFull(brw.brw, b)
	if err != nil {
		return b, n, err
	}
	return b, n, nil
}*/
