package bitop

func RESERVED_8(buf byte, start, end uint) uint8 {
	buf = buf << start
	buf = buf >> (start + (8 - end))
	return buf
}
func WRITE_8(i uint8, start, end uint) (buf byte) {
	buf = buf + i<<(8-end)
	return buf
}

/*func RESERVED_16(buf []byte, start, end uint) byte {
	if start >= 8 {
		buf[0] = 0x00
		buf[1] = buf[1] << (start - 8)
		buf[1] = buf[1] >> start + end
	}
	else if end >= 8 {
		buf[0] = buf[0] << start
		buf[1] = buf[1] >> (8 - )
	}
	buf = buf << start
	buf = buf >> (start + (8 - end))
	return buf
}*/
