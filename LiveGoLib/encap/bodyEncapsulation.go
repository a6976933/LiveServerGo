package encap

type TAG struct { //FLV TAGs
	previousTagSize   [4]byte
	dataSize          [3]byte
	timeStamp         [3]byte
	timeStampExtended byte
	streamID          [3]byte
}

func NewTag0() *TAG { //New PreviousTagSize0
	ret := &TAG{
		previousTagSize: [4]byte{0x00, 0x00, 0x00, 0x00},
	}
	return ret
}
func NewTag() *TAG {
	ret := &TAG{}
	return ret
}
