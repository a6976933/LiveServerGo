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
		chunk := cli.Prox.GetSendBytes()
		chunklen := len(chunk)
		w.Write(chunk)
		w.Write(av.TransUINT32_2_4Byte(uint32(chunklen)))
	}
}
