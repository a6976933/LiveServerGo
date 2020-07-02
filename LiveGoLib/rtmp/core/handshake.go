package core

import (
	"fmt"
	"io"
	"math/rand"
)

const (
	VERSION = 0x03
)

func makeRandBytes(buf []byte) []byte {
	rand.Read(buf)
	return buf
}
func HandShakeClientPkg(clientConn ClientConn) {
	var shareArr [(1 + (1536 * 2)) * 2]byte
	C0C1C2 := shareArr[:1+1536*2]
	C0 := C0C1C2[:1]
	C1 := C0C1C2[1:1537]
	C2 := C0C1C2[1537 : 1+1536*2]
	C0C1 := C0C1C2[:1537]
	n, err := io.ReadFull(clientConn.GetNetConn(), C0C1)
	if err != nil {
		return
	}
	if n != 1537 {
		return
	}
	if C0[0] != VERSION {
		return
	}
	S0S1S2 := shareArr[1+1536*2:]
	S0 := S0S1S2[:1]
	S1 := S0S1S2[1:1537]
	S2 := S0S1S2[1537 : 1+1536*2]
	S0[0] = 0x03
	clientConn.GetNetConn().Write(S0)
	copy(S1[0:4], C1[0:4])
	copy(S1[4:8], []byte{0x00, 0x00, 0x00, 0x00})
	makeRandBytes(S1)
	clientConn.GetNetConn().Write(S1)
	copy(S2, C1)
	/*copy(S2[0:4], C1[0:4])
	copy(S2[4:8], S1[0:4])
	copy(S2, C1[8:])*/
	clientConn.GetNetConn().Write(S2)
	io.ReadFull(clientConn.GetNetConn(), C2)
	fmt.Printf("Handshake Finished")
}
