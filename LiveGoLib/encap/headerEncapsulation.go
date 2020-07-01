package encap

import (
	"LiveGoLib/av"
)

type FLV_HEADER struct {
	btSignature  [3]byte
	btVersion    byte
	btFlags      byte
	btDataOffset [4]byte
}

func NewFLV_HEADER() *FLV_HEADER {
	fv := &FLV_HEADER{
		btSignature:  [3]byte{av.F, av.L, av.V},
		btVersion:    av.FLV_VERSION,
		btFlags:      av.VIDEO_AND_AUDIO,
		btDataOffset: [4]byte{0x00, 0x00, 0x00, 0x09},
	}
	return fv
}
