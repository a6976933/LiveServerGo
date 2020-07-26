package httpflv

import (
	/*"github.com/LiveGoLib/av"
	  "github.com/LiveGoLib/encap"*/

	"net/http"

	"LiveGoLib/av"
	"LiveGoLib/proxy"
)

type ClientInfo struct {
	Prox proxy.MiddleBuf
}

func (cli *ClientInfo) Handleflv(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Write([]byte{0x46, 0x4C, 0x56, 0x01, 0x05, 0x00, 0x00, 0x00, 0x09})
	w.Write([]byte{0x00, 0x00, 0x00, 0x00})
	for {
		/*var msgData []byte
		var tagType int
		var time uint32
		var streamID uint32*/
		var tagInfo []byte
		msgData, tagType, msgLen, time, streamID := cli.Prox.GetSendInfo()
		tagInfo = append(tagInfo, byte(tagType))
		tagInfo = append(tagInfo, av.TransUINT32_2_3Byte(msgLen)...)
		tagInfo = append(tagInfo, av.TransUINT32_2_3Byte(time)...)
		tagInfo = append(tagInfo, byte(0))
		tagInfo = append(tagInfo, av.TransUINT32_2_3Byte(streamID)...)
		tagInfo = append(tagInfo, msgData...)
		w.Write(tagInfo)
		w.Write(av.TransUINT32_2_4Byte(uint32(msgLen+11)))
	}
}
