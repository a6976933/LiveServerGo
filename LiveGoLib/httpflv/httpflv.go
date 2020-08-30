package httpflv

import (
	/*"github.com/LiveGoLib/av"
	  "github.com/LiveGoLib/encap"*/

	"fmt"
	"LiveGoLib/av"
	"github.com/gin-gonic/gin"
	"LiveGoLib/proxy"
)

type ClientInfo struct {
	Prox proxy.MiddleBuf
	TimeStampShift uint32
	TimeStampBase  uint32
}

func (cli *ClientInfo) Handleflv(w *gin.Context) {
	cli.Prox.SetClientCnt(true)
	w.Writer.Header().Set("Access-Control-Allow-Origin", "*")
	w.Data(200,"flv",[]byte{0x46, 0x4C, 0x56, 0x01, 0x05, 0x00, 0x00, 0x00, 0x09})
	w.Data(200,"flv",[]byte{0x00, 0x00, 0x00, 0x00})
	metaBuff := cli.Prox.SendMetaData()
	metalen := av.TransUINT32_2_3Byte(uint32(len(metaBuff)))
	//fmt.Println("******************************************************")
	//fmt.Println(len(metaBuff))
	//fmt.Println(metaBuff)
	w.Data(200,"flv",[]byte{0x12})
	w.Data(200,"flv",metalen)
	w.Data(200,"flv",[]byte{0x00,0x00,0x00,0x00,0x00,0x00,0x00})
	w.Data(200,"flv",metaBuff)
	w.Data(200,"flv",av.TransUINT32_2_4Byte(uint32(len(metaBuff)+11)))
	avcData := cli.Prox.GetAVCData()
	w.Data(200,"flv",[]byte{0x09})
	w.Data(200,"flv",av.TransUINT32_2_3Byte(uint32(len(avcData))))
	w.Data(200,"flv",[]byte{0x00,0x00,0x00,0x00,0x00,0x00,0x00})
	w.Data(200,"flv",avcData)
	w.Data(200,"flv",av.TransUINT32_2_4Byte(uint32(len(avcData)+11)))
	aacData := cli.Prox.GetAACData()
	w.Data(200,"flv",[]byte{0x08})
	w.Data(200,"flv",av.TransUINT32_2_3Byte(uint32(len(aacData))))
	w.Data(200,"flv",[]byte{0x00,0x00,0x00,0x00,0x00,0x00,0x00})
	w.Data(200,"flv",aacData)
	w.Data(200,"flv",av.TransUINT32_2_4Byte(uint32(len(aacData)+11)))
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
	w.Data(200,"flv",tagInfo)
	w.Data(200,"flv",av.TransUINT32_2_4Byte(uint32(msgLen+11)))
	chn := cli.Prox.GetSChan()
	for {
		/*var msgData []byte
		var tagType int
		var time uint32
		var streamID uint32*/
		//var b = make([]byte, 1)
		//b[0] = byte(uint8(tagType))
		tagInfo = tagInfo[:0]
		msgData = msgData[:0]
		<-chn
		msgData, tagType, msgLen, time, _ := cli.Prox.GetSendInfo();

		if uint32(len(msgData)+11) != msgLen+11 {
			fmt.Println("Invalid PreTagSize: msgData+11: ",len(msgData)+11," msg Len: ", msgLen)
		}
		cli.TimeStampShift = time - cli.TimeStampBase
		//fmt.Println("FLV Sending Data : ",tagType,msgLen,cli.TimeStampShift)
		tagInfo = append(tagInfo, byte(tagType))
		tagInfo = append(tagInfo, av.TransUINT32_2_3Byte(msgLen)...)
		tagInfo = append(tagInfo, av.TransUINT32_2_3Byte(time)...)
		tagInfo = append(tagInfo, byte(0))
		tagInfo = append(tagInfo, av.TransUINT32_2_3Byte(0)...)
		tagInfo = append(tagInfo, msgData...)
		w.Data(200,"flv",tagInfo)
		//fmt.Println(tagInfo)
		w.Data(200,"flv",av.TransUINT32_2_4Byte(uint32(msgLen+11)))
		//fmt.Println(av.TransUINT32_2_4Byte(uint32(msgLen+11)))
	}
}
