package amf

import (
	"encoding/binary"
	"fmt"
	"math"

	"LiveGoLib/av"
)

type Object map[string]interface{}

func Decode_AMF0(Buff []byte) (interface{}, []byte) {
	var ret interface{}
	switch uint(Buff[0]) {
	case 0:
		ret, Buff = Decode_Number(Buff, true)
	case 1:
		ret, Buff = Decode_Boolean(Buff, true)
	case 2:
		ret, Buff = Decode_String(Buff, true)
	case 3:
		ret, Buff = Decode_Object(Buff, true)
	case 4:

	case 5:
		ret, Buff = Decode_Null(Buff)
	case 6:
		ret, Buff = Decode_Undefined(Buff)
	case 8:
		ret, Buff = Decode_ECMAArray(Buff)
	}
	//fmt.Println(ret)
	/*	if ret == "connect" {
		TransactionID, Buff := Decode_Number(Buff, true)
		fmt.Println("Trans ID: ", TransactionID)
		Decode_Object(Buff)
	}*/
	return ret, Buff
}

func Decode_Number(Buff []byte, clear1byte bool) (float64, []byte) {
	if clear1byte {
		Buff = av.ShiftBytesRight(Buff, 1)
	}
	bits := binary.BigEndian.Uint64(Buff[:8])
	Buff = av.ShiftBytesRight(Buff, 8)
	return math.Float64frombits(bits), Buff
}

func Decode_Boolean(Buff []byte, clear1byte bool) (bool, []byte) {
	if clear1byte {
		Buff = av.ShiftBytesRight(Buff, 1)
	}
	if Buff[0] == 0x01 {
		return true, av.ShiftBytesRight(Buff, 1)
	} else {
		return false, av.ShiftBytesRight(Buff, 1)
	}
}

func Decode_String(Buff []byte, clear1byte bool) (string, []byte) {
	if clear1byte {
		Buff = av.ShiftBytesRight(Buff, 1)
	}
	length := av.Trans2Byte2UINT16(Buff[:2])
	//fmt.Println(Buff)
	return string(Buff[2 : length+2]), av.ShiftBytesRight(Buff, uint(length+2))
}

func Decode_Object(Buff []byte, clear1byte bool) (av.Object, []byte) {
	if clear1byte {
		Buff = av.ShiftBytesRight(Buff, 1)
	}
	var obj = make(av.Object)
	for {
		Buf := Buff
		key, Buf := Decode_String(Buf, false)
		if key == "" {
			if Buf[0] == 9 {
				fmt.Printf("OBJ End")
				break
			}
		}
		//fmt.Println(key)
		val, Buf := Decode_AMF0(Buf)
		obj[key] = val
		Buff = Buf
	}
	return obj, Buff
}

func Decode_Null(Buff []byte) (interface{}, []byte) {
	Buff = av.ShiftBytesRight(Buff, 1)
	return nil, Buff
}

func Decode_Undefined(Buff []byte) (interface{}, []byte) {
	Buff = av.ShiftBytesRight(Buff, 1)
	return nil, Buff
}

func Decode_ECMAArray(Buff []byte) (av.Object, []byte) {
	Buff = av.ShiftBytesRight(Buff, 1)
	_ = av.Trans4Byte2UINT32(Buff)
	Buff = av.ShiftBytesRight(Buff, 4)
	var obj map[string]interface{}
	obj, Buff = Decode_Object(Buff, false)
	return obj, Buff
}
