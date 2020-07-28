package core

import (
	"encoding/binary"
	"fmt"
	"os"
	"io"
	"LiveGoLib/amf"
	"LiveGoLib/av"
	"LiveGoLib/bitop"
	"bufio"
)

type CommandR interface {
	Handle(clientInfo *ClientInfo)
}

type CommandW interface {
	Write(clientInfo *ClientInfo)
}

type CommandContext struct {
	clientInfo *ClientInfo
	cmdType    CommandR
	wCmdType   CommandW
}

type SetChunkSize struct{}
type AbortMessage struct{}
type Acknowledgement struct{}
type StreamBegin struct{}
type WindowAcknowledgementSize struct{ windowSize uint }
type SetPeerBandwidth struct{}
type AudioMessage struct{}
type VideoMessage struct{}

type CmdMessageAMF0 struct{}

type ConnectAMF0 struct{}
type ConnectRespAMF0 struct{}

type ReleaseStream struct{}
type FCPublish struct{}
type CreateStream struct{}

type Publish struct{}
type OnStatus struct{}

type SetDataFrame struct{}

func CmdFactory(clientInfo *ClientInfo) {
	var cmdContext CommandContext
	switch clientInfo.MessageTypeID {
	case 1:
		var scs *SetChunkSize
		cmdContext.cmdType = scs
		cmdContext.cmdType.Handle(clientInfo)
	case 2:

	case 8:
		var am *AudioMessage
		cmdContext.cmdType = am
		cmdContext.cmdType.Handle(clientInfo)
	case 9:
		var vm *VideoMessage
		cmdContext.cmdType = vm
		cmdContext.cmdType.Handle(clientInfo)
	case 18:
		var cMsgAMF0 *CmdMessageAMF0
		cmdContext.cmdType = cMsgAMF0
		cmdContext.cmdType.Handle(clientInfo)
	case 20:
		var cMsgAMF0 *CmdMessageAMF0
		cmdContext.cmdType = cMsgAMF0
		cmdContext.cmdType.Handle(clientInfo)
	default:
		bufio.NewReader(os.Stdin).ReadBytes('\n')
	}
}

func WriteCmdHeader(ch ChunkHeader) []byte {
	var b = make([]byte, 1)
	b[0] = bitop.WRITE_8(ch.Fmt, 0, 2)
	if ch.Csid > 63 {
		if ch.Csid > 319 {
			b = append(b, make([]byte, 2)...)
			b[0] = bitop.WRITE_8(1, 2, 8)
			b[1] = byte(ch.Csid - 64)
		} else {
			b = append(b, make([]byte, 1)...)
			b[0] = bitop.WRITE_8(0, 2, 8)
			binary.BigEndian.PutUint16(b[1:3], uint16(ch.Csid-64))
		}
	} else {
		b[0] = bitop.WRITE_8(uint8(ch.Csid), 2, 8)
	}
	switch ch.Fmt {
	case 0:
		b = append(b, av.TransUINT32_2_3Byte(ch.Timestamp)...)
		b = append(b, av.TransUINT32_2_3Byte(ch.MessageLength)...)
		b = append(b, av.TransUINT8_2_1Byte(ch.MessageTypeID)...)
		b = append(b, av.TransUINT32_2_4Byte(ch.MessageStreamID)...)
	case 1:
		b = append(b, av.TransUINT32_2_3Byte(ch.Timestamp)...)
		b = append(b, av.TransUINT32_2_3Byte(ch.MessageLength)...)
		b = append(b, av.TransUINT8_2_1Byte(ch.MessageTypeID)...)
	case 2:
		b = append(b, av.TransUINT32_2_3Byte(ch.Timestamp)...)
	case 3:

	}
	//ExtendTimestamp
	return b
}

func SetCmdChunk(timestamp, msgLength uint32, typeID uint8, streamID uint32) ChunkHeader {
	var ch ChunkHeader
	ch.Fmt = 0
	ch.Csid = 2
	ch.Timestamp = timestamp
	ch.MessageLength = msgLength
	ch.MessageTypeID = typeID
	ch.MessageStreamID = streamID
	return ch
}

func WriteCmd4Byte(n uint32) []byte {
	var b = make([]byte, 4)
	binary.BigEndian.PutUint32(b, n)
	return b
}

func ParseSplitChunk(clientInfo *ClientInfo) ([]byte, error) {
	var TempBuff []byte
	var err error
	TempBuff, err = clientInfo.rw.ReadNByte(int(clientInfo.chunksize), clientInfo.Buff)
	if err != nil {
		return TempBuff, err
	}
	clientInfo.MessageLength = clientInfo.MessageLength - clientInfo.chunksize
	for {
		var TempBuff2 []byte
		err = ParseBasicHeader(clientInfo)
		if err != nil {
			fmt.Println(err)
			return TempBuff, err
		}
		if clientInfo.Fmt != 3 {
			bufio.NewReader(os.Stdin).ReadBytes('\n')
		}
		err = ParseMessageHeader(clientInfo)
		if err != nil {
			fmt.Println(err)
			return TempBuff, err
		}
		if clientInfo.MessageLength < clientInfo.chunksize {
			fmt.Println("end end end")
			TempBuff2, err = clientInfo.rw.ReadNByte(int(clientInfo.MessageLength), clientInfo.Buff)
			if err != nil {
				return TempBuff, err
			}
			TempBuff = append(TempBuff, TempBuff2...)
			return TempBuff, nil
		} else {
			fmt.Println("continue")
			TempBuff2, err = clientInfo.rw.ReadNByte(int(clientInfo.chunksize), clientInfo.Buff)
			if err != nil {
				fmt.Println("has err ",err)
				return TempBuff, err
			}
			clientInfo.MessageLength = clientInfo.MessageLength - clientInfo.chunksize
			//fmt.Println(TempBuff2)
			TempBuff = append(TempBuff, TempBuff2...)
		}
	}
}

func ParseChunk(clientInfo *ClientInfo) ([]byte, error) {
	var Buff []byte
	var err error
	fmt.Println("Message length:", clientInfo.MessageLength)
	OrgMsgLen := clientInfo.MessageLength
	if clientInfo.MessageLength > clientInfo.chunksize {
		Buff, err = ParseSplitChunk(clientInfo)
		if err != nil {
			return Buff, err
		}
	} else {
		Buff, err = clientInfo.rw.ReadNByte(int(clientInfo.MessageLength), clientInfo.Buff)
		if err != nil {
			fmt.Println(err)
			return Buff, err
		}
	}
	clientInfo.MessageLength = OrgMsgLen
	return Buff, nil
}

func (*CmdMessageAMF0) Handle(clientInfo *ClientInfo) {
	var cmdContext CommandContext
	var amfCmdType interface{}
	var err error
	clientInfo.Buff, err = ParseChunk(clientInfo)
	if err != nil {
		fmt.Println(err)
		HandleDisConnErr(clientInfo, err)
	}
	amfCmdType, clientInfo.Buff = amf.Decode_AMF0(clientInfo.Buff)
	fmt.Println(amfCmdType)
	switch amfCmdType {
	case "connect":
		var cnt *ConnectAMF0
		cmdContext.cmdType = cnt
		cmdContext.cmdType.Handle(clientInfo)
	case "releaseStream":
		var rls *ReleaseStream
		cmdContext.cmdType = rls
		cmdContext.cmdType.Handle(clientInfo)
	case "FCPublish":
		var fcp *FCPublish
		cmdContext.cmdType = fcp
		cmdContext.cmdType.Handle(clientInfo)
	case "createStream":
		var cs *CreateStream
		var cr *ConnectRespAMF0
		cmdContext.cmdType = cs
		cmdContext.cmdType.Handle(clientInfo)
		cmdContext.wCmdType = cr
		cmdContext.wCmdType.Write(clientInfo)
	case "publish":
		var pb *Publish
		var sb *StreamBegin
		var cr *ConnectRespAMF0
		var ost *OnStatus
		cmdContext.cmdType = pb
		cmdContext.cmdType.Handle(clientInfo)
		cmdContext.wCmdType = sb
		cmdContext.wCmdType.Write(clientInfo)
		cmdContext.wCmdType = cr
		cmdContext.wCmdType.Write(clientInfo)
		cmdContext.wCmdType = ost
		cmdContext.wCmdType.Write(clientInfo)
		/*Buff, _, _ := clientInfo.rw.ReadBuff()
		fmt.Println(Buff)*/
	case "@setDataFrame":
		var sdf *SetDataFrame
		cmdContext.cmdType = sdf
		cmdContext.cmdType.Handle(clientInfo)
	}
}

func (sb *StreamBegin) Write(clientInfo *ClientInfo) {
	cmdData := WriteCmd4Byte(1)
	var ch ChunkHeader
	ch = SetCmdChunk(clientInfo.Timestamp, 4, 4, 0)
	cmdHd := WriteCmdHeader(ch)
	cmdHd = append(cmdHd, cmdData...)
	_, err := clientInfo.rw.Write(cmdHd)
	if err != nil {
		fmt.Println(err)
	} else {
		fmt.Println("Stream Begin Success")
	}
}

func (scs *SetChunkSize) Handle(clientInfo *ClientInfo) {
	var err error
	clientInfo.Buff, err = clientInfo.rw.ReadNByte(4, clientInfo.Buff)
	if err != nil {
		fmt.Println(err)
	}
	clientInfo.chunksize = binary.BigEndian.Uint32(clientInfo.Buff[0:4])
	fmt.Println("Set Chunk Size to ", clientInfo.chunksize)
}

func (scs *SetChunkSize) Write(clientInfo *ClientInfo) {
	cmdData := WriteCmd4Byte(clientInfo.chunksize)
	var ch ChunkHeader
	ch = SetCmdChunk(clientInfo.Timestamp, 4, 1, 0)
	cmdHd := WriteCmdHeader(ch)
	cmdHd = append(cmdHd, cmdData...)
	//_, err := clientInfo.Conn.Write(cmdHd)
	_, err := clientInfo.rw.Write(cmdHd)
	if err != nil {
		fmt.Println(err)
	} else {
		fmt.Println("Set Chunk Success")
	}
}

func (was *WindowAcknowledgementSize) Write(clientInfo *ClientInfo) {
	cmdData := WriteCmd4Byte(clientInfo.windowAckSize)
	var ch ChunkHeader
	ch = SetCmdChunk(clientInfo.Timestamp, 4, 5, 0)
	cmdHd := WriteCmdHeader(ch)
	cmdHd = append(cmdHd, cmdData...)
	//_, err := clientInfo.Conn.Write(cmdHd)
	_, err := clientInfo.rw.Write(cmdHd)
	if err != nil {
		fmt.Println(err)
	} else {
		fmt.Println("Window Ack Success")
	}
}

func (spb *SetPeerBandwidth) Write(clientInfo *ClientInfo) {
	cmdData := WriteCmd4Byte(2500000)
	var ch ChunkHeader
	ch = SetCmdChunk(clientInfo.Timestamp, 5, 6, 0)
	cmdHd := WriteCmdHeader(ch)
	cmdData = append(cmdData, byte(2))
	cmdHd = append(cmdHd, cmdData...)
	_, err := clientInfo.rw.Write(cmdHd)
	if err != nil {
		fmt.Println(err)
	} else {
		fmt.Println("Set Peer Success")
	}
}

func (cnt *ConnectAMF0) Handle(clientInfo *ClientInfo) {
	var TransactionID float64
	clientInfo.transactionID, clientInfo.Buff = amf.Decode_Number(clientInfo.Buff, true)
	fmt.Println("Trans ID: ", TransactionID)
	clientInfo.obj, _ = amf.Decode_Object(clientInfo.Buff, true)
	fmt.Println(clientInfo.obj)
	clientInfo.windowAckSize = 2500000
	var cmdContext CommandContext
	var windowSent *WindowAcknowledgementSize
	var setPeer *SetPeerBandwidth
	//var streamBegin *StreamBegin
	var conResp *ConnectRespAMF0
	var setChunkSize *SetChunkSize
	cmdContext.wCmdType = windowSent
	cmdContext.wCmdType.Write(clientInfo)
	cmdContext.wCmdType = setPeer
	cmdContext.wCmdType.Write(clientInfo)
	clientInfo.chunksize = 4096
	cmdContext.wCmdType = setChunkSize
	cmdContext.wCmdType.Write(clientInfo)
	//cmdContext.wCmdType = streamBegin
	//cmdContext.wCmdType.Write(clientInfo)
	//bufio.NewReader(os.Stdin).ReadBytes('\n')clientInfo.chunksize = 1024

	cmdContext.wCmdType = conResp
	cmdContext.wCmdType.Write(clientInfo)
	clientInfo.media = SaveFLVHeader()
}

func SaveFLVHeader() *os.File {
	fo ,err := os.Create("stv.flv")
	if err != nil {
		fmt.Println(err)
	}
	fo.Write([]byte{0x46, 0x4C, 0x56, 0x01, 0x05, 0x00, 0x00, 0x00, 0x09})
	fo.Write([]byte{0x00, 0x00, 0x00, 0x00})
	return fo
}

func (rs *ReleaseStream) Handle(clientInfo *ClientInfo) {
	var TransactionID float64
	clientInfo.transactionID, clientInfo.Buff = amf.Decode_Number(clientInfo.Buff, true)
	fmt.Println("Trans ID: ", TransactionID)
	var NullType interface{}
	NullType, clientInfo.Buff = amf.Decode_AMF0(clientInfo.Buff)
	fmt.Println(NullType)
	var RlsStreamStr interface{}
	RlsStreamStr, clientInfo.Buff = amf.Decode_AMF0(clientInfo.Buff)
	fmt.Println(RlsStreamStr)
}

func (fcp *FCPublish) Handle(clientInfo *ClientInfo) {
	clientInfo.transactionID, clientInfo.Buff = amf.Decode_Number(clientInfo.Buff, true)
	fmt.Println("Trans ID: ", clientInfo.transactionID)
	var NullType interface{}
	NullType, clientInfo.Buff = amf.Decode_AMF0(clientInfo.Buff)
	fmt.Println(NullType)
	var FcpStr interface{}
	FcpStr, clientInfo.Buff = amf.Decode_AMF0(clientInfo.Buff)
	fmt.Println(FcpStr)
}

func (cs *CreateStream) Handle(clientInfo *ClientInfo) {
	fmt.Println("prvent trans ID: ", clientInfo.transactionID)
	clientInfo.transactionID, clientInfo.Buff = amf.Decode_Number(clientInfo.Buff, true)
	fmt.Println("Trans ID: ", clientInfo.transactionID)
	var obj interface{}
	obj, clientInfo.Buff = amf.Decode_AMF0(clientInfo.Buff)
	fmt.Println(obj)
	if obj == nil {
		return
	}
}

func (pb *Publish) Handle(clientInfo *ClientInfo) {
	clientInfo.transactionID, clientInfo.Buff = amf.Decode_Number(clientInfo.Buff, true)
	fmt.Println("Trans ID: ", clientInfo.transactionID)
	var obj interface{}
	obj, clientInfo.Buff = amf.Decode_AMF0(clientInfo.Buff)
	fmt.Println(obj)
	clientInfo.publishName, clientInfo.Buff = amf.Decode_String(clientInfo.Buff, true)
	clientInfo.publishType, clientInfo.Buff = amf.Decode_String(clientInfo.Buff, true)
	fmt.Println("Publish Name: ", clientInfo.publishName)
	fmt.Println("Publish Type: ", clientInfo.publishType)
}

func (sdf *SetDataFrame) Handle(clientInfo *ClientInfo) {
	var onMeta string
	onMeta, clientInfo.Buff = amf.Decode_String(clientInfo.Buff, true)
	fmt.Println(onMeta)
	var obj map[string]interface{}
	obj, clientInfo.Buff = amf.Decode_ECMAArray(clientInfo.Buff)
	fmt.Println(obj)
	clientInfo.dataFrame = obj
}

func (vm *VideoMessage) Handle(clientInfo *ClientInfo) {
	msgByte, err := ParseChunk(clientInfo)
	//bufio.NewReader(os.Stdin).ReadBytes('\n')
	if err != nil {
		fmt.Println(err)
		HandleDisConnErr(clientInfo, err)
	}
	SaveFLVBody(clientInfo, msgByte, 9)
	clientInfo.prox.SetSendInfo(msgByte,9,clientInfo.MessageLength,clientInfo.Timestamp,clientInfo.MessageStreamID)
	//bufio.NewReader(os.Stdin).ReadBytes('\n')
}

func (am *AudioMessage) Handle(clientInfo *ClientInfo) {
	msgByte, err := ParseChunk(clientInfo)
	if err != nil {
		fmt.Println(err)
		HandleDisConnErr(clientInfo, err)
	}
	SaveFLVBody(clientInfo, msgByte, 8)
	clientInfo.prox.SetSendInfo(msgByte,8,clientInfo.MessageLength,clientInfo.Timestamp,clientInfo.MessageStreamID)
	//bufio.NewReader(os.Stdin).ReadBytes('\n')
}

func SaveFLVBody(clientInfo *ClientInfo,msgByte []byte, tagType uint){
	var b = make([]byte, 1)
	b[0] = byte(uint8(tagType))
	clientInfo.media.Write(b)
	clientInfo.media.Write(av.TransUINT32_2_3Byte(uint32(len(msgByte))))
	clientInfo.media.Write(av.TransUINT32_2_3Byte(clientInfo.Timestamp))
	b[0] = byte(uint8(0))
	clientInfo.media.Write(b)
	clientInfo.media.Write(av.TransUINT32_2_3Byte(0))
	clientInfo.media.Write(msgByte)
	preTagSize := av.TransUINT32_2_4Byte(uint32(len(msgByte)+11))
	clientInfo.media.Write(preTagSize)
}

/*
func (vm *VideoMessage) Handle(clientInfo *ClientInfo) {
	var err error
	var dataLen uint
	clientInfo.Buff, err = clientInfo.rw.ReadNByte(1, clientInfo.Buff)
	if err != nil {
		fmt.Println(err)
	}
	dataLen += 1
	frameType := bitop.RESERVED_8(clientInfo.Buff[0], 0, 4)
	codecID := bitop.RESERVED_8(clientInfo.Buff[0], 4, 8)
	clientInfo.videoCodecID = uint32(codecID)
	clientInfo.videoFrameType = uint32(frameType)
	fmt.Println("frame type: ", frameType, " CodecID: ", codecID)
	clientInfo.Buff, err = clientInfo.rw.ReadNByte(1, clientInfo.Buff)
	if err != nil {
		fmt.Println(err)
	}
	dataLen += 1
	clientInfo.avcPacketType = uint8(clientInfo.Buff[0])
	if clientInfo.avcPacketType == 1 {
		clientInfo.Buff, err = clientInfo.rw.ReadNByte(3, clientInfo.Buff)
		if err != nil {
			fmt.Println(err)
		}
		dataLen += 3
		clientInfo.compositionTime = av.Trans3Byte2UINT32(clientInfo.Buff[0:3])
	}
	dataLen = uint(clientInfo.MessageLength) - dataLen
	clientInfo.Buff, err = clientInfo.rw.ReadNByte(int(dataLen), clientInfo.Buff)
	fmt.Println(clientInfo.Buff)
}*/
/*
func (am *AudioMessage) Handle(clientInfo *ClientInfo) {
	var err error
	var dataLen uint
	clientInfo.Buff, err = clientInfo.rw.ReadNByte(1, clientInfo.Buff)
	if err != nil {
		fmt.Println(err)
	}
	dataLen += 1
	clientInfo.soundFormat = uint32(bitop.RESERVED_8(clientInfo.Buff[0], 0, 4))
	clientInfo.soundRate = uint32(bitop.RESERVED_8(clientInfo.Buff[0], 4, 6))
	clientInfo.soundSize = uint32(bitop.RESERVED_8(clientInfo.Buff[0], 6, 7))
	clientInfo.soundType = uint32(bitop.RESERVED_8(clientInfo.Buff[0], 7, 8))
	fmt.Println("sound format: ", clientInfo.soundFormat, " sound rate: ", clientInfo.soundRate)
	fmt.Println(" sound size: ", clientInfo.soundSize, " sound type: ", clientInfo.soundType)
	if clientInfo.soundFormat == 2 {
		clientInfo.Buff, err = clientInfo.rw.ReadNByte(1, clientInfo.Buff)
		if err != nil {
			fmt.Println(err)
		}
		dataLen += 1
	}

	clientInfo.avcPacketType = uint8(clientInfo.Buff[0])
	if clientInfo.avcPacketType == 1 {
		clientInfo.Buff, err = clientInfo.rw.ReadNByte(3, clientInfo.Buff)
		if err != nil {
			fmt.Println(err)
		}
		dataLen += 3
		clientInfo.compositionTime = av.Trans3Byte2UINT32(clientInfo.Buff[0:3])
	}
	dataLen = uint(clientInfo.MessageLength) - dataLen
	clientInfo.Buff, err = clientInfo.rw.ReadNByte(int(dataLen), clientInfo.Buff)
	fmt.Println(clientInfo.Buff)
}*/

func (cra *ConnectRespAMF0) Write(clientInfo *ClientInfo) {
	respObj := make(amf.Object)
	respObj["fmsVer"] = "FMS/3,0,1,123"
	respObj["capabilities"] = float64(31)
	respEvent := make(amf.Object)
	respEvent["level"] = "status"
	respEvent["code"] = "NetConnection.Connect.Success"
	respEvent["description"] = "Connection succeeded."
	respEvent["objectEncoding"] = float64(0)
	var conRespAMF0 []byte
	conRespAMF0 = amf.Encode_AMF0("_result")
	conRespAMF0 = append(conRespAMF0, amf.Encode_AMF0(float64(clientInfo.transactionID))...)
	conRespAMF0 = append(conRespAMF0, amf.Encode_AMF0(respObj)...)
	conRespAMF0 = append(conRespAMF0, amf.Encode_AMF0(respEvent)...)
	leng := len(conRespAMF0)
	var ch ChunkHeader
	ch = SetCmdChunk(clientInfo.Timestamp, uint32(leng), 20, clientInfo.MessageStreamID)
	cmdHd := WriteCmdHeader(ch)
	cmdHd = append(cmdHd, conRespAMF0...)
	_, err := clientInfo.rw.Write(cmdHd)
	if err != nil {
		fmt.Println(err)
	} else {
		fmt.Println("Conn Resp Success")
	}
}

func (ost *OnStatus) Write(clientInfo *ClientInfo) {
	respEvent := make(amf.Object)
	respEvent["level"] = "status"
	respEvent["code"] = "NetStream.Publish.Start"
	respEvent["description"] = "Start publishing."
	var ostAMF0 []byte
	ostAMF0 = amf.Encode_AMF0("onStatus")
	ostAMF0 = append(ostAMF0, amf.Encode_AMF0(float64(clientInfo.transactionID))...)
	ostAMF0 = append(ostAMF0, amf.Encode_AMF0(nil)...)
	ostAMF0 = append(ostAMF0, amf.Encode_AMF0(respEvent)...)
	leng := len(ostAMF0)
	var ch ChunkHeader
	ch = SetCmdChunk(clientInfo.Timestamp, uint32(leng), 20, 1)
	cmdHd := WriteCmdHeader(ch)
	cmdHd = append(cmdHd, ostAMF0...)
	_, err := clientInfo.rw.Write(cmdHd)
	if err != nil {
		fmt.Println(err)
	} else {
		fmt.Println("onStatus Send Success")
	}
}

func HandleDisConnErr(clientInfo *ClientInfo, err error) {
	if err == io.EOF {
		clientInfo.media.Close()
		clientInfo.disConn = true
	}
}
