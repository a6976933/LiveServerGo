package httpflv

import (
	/*"github.com/LiveGoLib/av"
	  "github.com/LiveGoLib/encap"*/

	"net/http"
	"fmt"
	"LiveGoLib/av"
	"LiveGoLib/proxy"
)

type ClientInfo struct {
	Prox proxy.MiddleBuf
	TimeStampShift uint32
	TimeStampBase  uint32
}

func (cli *ClientInfo) Handleflv(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Write([]byte{0x46, 0x4C, 0x56, 0x01, 0x05, 0x00, 0x00, 0x00, 0x09})
	w.Write([]byte{0x00, 0x00, 0x00, 0x00})
	msgData, tagType, msgLen, time, _ := cli.Prox.GetSendInfo()
	var tagInfo []byte
	cli.TimeStampBase = time
	cli.TimeStampShift = 0
	tagInfo = append(tagInfo, byte(tagType))
	tagInfo = append(tagInfo, av.TransUINT32_2_3Byte(msgLen)...)
	tagInfo = append(tagInfo, av.TransUINT32_2_3Byte(cli.TimeStampShift)...)
	tagInfo = append(tagInfo, byte(0))
	tagInfo = append(tagInfo, av.TransUINT32_2_3Byte(0)...)
	tagInfo = append(tagInfo, msgData...)
	w.Write(tagInfo)
	w.Write(av.TransUINT32_2_4Byte(uint32(msgLen+11)))
	for {
		/*var msgData []byte
		var tagType int
		var time uint32
		var streamID uint32*/
		//var b = make([]byte, 1)
		//b[0] = byte(uint8(tagType))
		var tempEmp []byte
		for true {
			if _, tagTypeTest, msgLenTest, timeTest, _ := cli.Prox.GetSendInfo(); (tagTypeTest == tagType)||(msgLenTest == msgLen)||(timeTest == time) {

			} else {
				break
			}
		}
		tagInfo = tempEmp
		msgData, tagType, msgLen, time, _ = cli.Prox.GetSendInfo()
		cli.TimeStampShift = time - cli.TimeStampBase
		fmt.Println("FLV Sending Data : ",tagType,msgLen,cli.TimeStampShift)
		tagInfo = append(tagInfo, byte(tagType))
		tagInfo = append(tagInfo, av.TransUINT32_2_3Byte(msgLen)...)
		tagInfo = append(tagInfo, av.TransUINT32_2_3Byte(cli.TimeStampShift)...)
		tagInfo = append(tagInfo, byte(0))
		tagInfo = append(tagInfo, av.TransUINT32_2_3Byte(0)...)
		tagInfo = append(tagInfo, msgData...)
		w.Write(tagInfo)
		fmt.Println(tagInfo)
		w.Write(av.TransUINT32_2_4Byte(uint32(msgLen+11)))
		fmt.Println(av.TransUINT32_2_4Byte(uint32(msgLen+11)))
	}
}
