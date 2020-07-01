package amf

import (
	"encoding/binary"
	//"fmt"
	"math"
	"reflect"
)

func Encode_AMF0(enObj interface{}) []byte {
	if enObj == nil {
		return Encode_Null()
	}
	encodeObj := reflect.ValueOf(enObj)
	var ret []byte
	switch encodeObj.Kind() {
	case reflect.Float64:
		ret = append(ret, Encode_Number(encodeObj.Float(), false)...)
	case reflect.Bool:
		ret = append(ret, Encode_Boolean(encodeObj.Bool(), false)...)
	case reflect.String:
		ret = Encode_String(encodeObj.String(), false)
	case reflect.Map:
		obj, _ := enObj.(Object)
		ret = Encode_Object(obj)
	}
	/*	if ret == "connect" {
		TransactionID, Buff := Decode_Number(Buff, true)
		fmt.Println("Trans ID: ", TransactionID)
		Decode_Object(Buff)
	}*/
	return ret
}

func Encode_Number(num float64, skip bool) []byte {
	var b = make([]byte, 8)
	bits := math.Float64bits(num)
	binary.BigEndian.PutUint64(b, bits)
	if !skip {
		var ret = make([]byte, 1)
		ret[0] = byte(0)
		b = append(ret, b...)
	}
	return b
}

func Encode_Boolean(input bool, skip bool) []byte {

	var b = make([]byte, 1)
	if input {
		b[0] = byte(1)
	} else {
		b[0] = byte(0)
	}
	if !skip {
		var ret = make([]byte, 1)
		ret[0] = byte(1)
		b = append(ret, b...)
	}
	return b
}

func Encode_String(input string, skip bool) []byte {
	var leng = make([]byte, 2)
	binary.BigEndian.PutUint16(leng, uint16(len(input)))
	var b = []byte(input)
	b = append(leng, b...)
	if !skip {
		var ret = make([]byte, 1)
		ret[0] = byte(2)
		b = append(ret, b...)
	}
	return b
}

func Encode_Object(input Object) []byte {
	var b = make([]byte, 1)
	b[0] = byte(3)
	for k, v := range input {
		key := Encode_String(k, true)
		b = append(b, key...)
		val := Encode_AMF0(v)
		b = append(b, val...)
	}
	var end = make([]byte, 3)
	end[0] = byte(0)
	end[1] = byte(0)
	end = append(end, Encode_ObjectEnd()...)
	b = append(b, end...)
	return b
}

func Encode_Null() []byte {
	var b = make([]byte, 1)
	b[0] = byte(5)
	return b
}

func Encode_ObjectEnd() []byte {
	var b = make([]byte, 1)
	b[0] = byte(9)
	return b
}
