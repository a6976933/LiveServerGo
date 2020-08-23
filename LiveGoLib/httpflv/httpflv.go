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
	cli.Prox.SetClientCnt(true)
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Write([]byte{0x46, 0x4C, 0x56, 0x01, 0x05, 0x00, 0x00, 0x00, 0x09})
	w.Write([]byte{0x00, 0x00, 0x00, 0x00})
	metaBuff := cli.Prox.SendMetaData()
	metalen := av.TransUINT32_2_3Byte(uint32(len(metaBuff)))
	fmt.Println("******************************************************")
	fmt.Println(len(metaBuff))
	fmt.Println(metaBuff)
	w.Write([]byte{0x12})
	w.Write(metalen)
	w.Write([]byte{0x00,0x00,0x00,0x00,0x00,0x00,0x00})
	w.Write(metaBuff)
	w.Write(av.TransUINT32_2_4Byte(uint32(len(metaBuff)+11)))
	avcData := cli.Prox.GetAVCData()
	w.Write([]byte{0x09})
	w.Write(av.TransUINT32_2_3Byte(uint32(len(avcData))))
	w.Write([]byte{0x00,0x00,0x00,0x00,0x00,0x00,0x00})
	w.Write(avcData)
	w.Write(av.TransUINT32_2_4Byte(uint32(len(avcData)+11)))
	aacData := cli.Prox.GetAACData()
	w.Write([]byte{0x08})
	w.Write(av.TransUINT32_2_3Byte(uint32(len(aacData))))
	w.Write([]byte{0x00,0x00,0x00,0x00,0x00,0x00,0x00})
	w.Write(aacData)
	w.Write(av.TransUINT32_2_4Byte(uint32(len(aacData)+11)))
	msgData, tagType, msgLen, time, _ := cli.Prox.GetSendInfo()
	var tagInfo []byte
	cli.TimeStampBase = time
	cli.TimeStampShift = 0
	tagInfo = append(tagInfo, byte(tagType))
	tagInfo = append(tagInfo, av.TransUINT32_2_3Byte(msgLen)...)
	tagInfo = append(tagInfo, av.TransUINT32_2_3Byte(cli.TimeStampBase)...)
	tagInfo = append(tagInfo, byte(0))
	tagInfo = append(tagInfo, av.TransUINT32_2_3Byte(0)...)
	tagInfo = append(tagInfo, msgData...)
	w.Write(tagInfo)
	w.Write(av.TransUINT32_2_4Byte(uint32(msgLen+11)))
	chn := cli.Prox.GetSChan()
	for {
		/*var msgData []byte
		var tagType int
		var time uint32
		var streamID uint32*/
		//var b = make([]byte, 1)
		//b[0] = byte(uint8(tagType))
		var tempEmp []byte
		<-chn
		msgData, tagType, msgLen, time, _ := cli.Prox.GetSendInfo();

		if uint32(len(msgData)+11) != msgLen+11 {
			fmt.Println("Invalid PreTagSize: msgData+11: ",len(msgData)+11," msg Len: ", msgLen)
		}
		tagInfo = tempEmp
		cli.TimeStampShift = time - cli.TimeStampBase
		fmt.Println("FLV Sending Data : ",tagType,msgLen,cli.TimeStampShift)
		tagInfo = append(tagInfo, byte(tagType))
		tagInfo = append(tagInfo, av.TransUINT32_2_3Byte(msgLen)...)
		tagInfo = append(tagInfo, av.TransUINT32_2_3Byte(time)...)
		tagInfo = append(tagInfo, byte(0))
		tagInfo = append(tagInfo, av.TransUINT32_2_3Byte(0)...)
		tagInfo = append(tagInfo, msgData...)
		w.Write(tagInfo)
		fmt.Println(tagInfo)
		w.Write(av.TransUINT32_2_4Byte(uint32(msgLen+11)))
		fmt.Println(av.TransUINT32_2_4Byte(uint32(msgLen+11)))
	}
}
