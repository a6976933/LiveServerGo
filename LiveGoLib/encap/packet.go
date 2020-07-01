package encap

type Packet struct {
	Data []byte
	FLV_HEADER
	TAG
}
