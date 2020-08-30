package core

import (
	"fmt"

	"LiveGoLib/av"
	"LiveGoLib/bitop"
	"LiveGoLib/proxy"
)

func StartHandleConn(clientConn ClientConn, prox proxy.MiddleBuf) {
	cntInfo := initClientInfo(clientConn, prox)
	for {
		ParseBasicHeader(cntInfo)
		ParseMessageHeader(cntInfo)
		if IsExtendedTimeStamp(cntInfo) {
			ParseExtendedTimeStamp(cntInfo)
		}
		CmdFactory(cntInfo)
		if cntInfo.disConn {
			return
		}
	}
	//CmdHandle(cntInfo)
}

func initClientInfo(clientConn ClientConn, prox proxy.MiddleBuf) *ClientInfo {
	var buff = make([]byte, BUFFER_SIZE)
	cntInfo := NewClientInfo(clientConn.GetNetConn(), buff)
	cntInfo.prox = prox
	return cntInfo
}

func ParseBasicHeader(clientInfo *ClientInfo) error {
	var err error
	clientInfo.Buff, err = clientInfo.rw.ReadNByte(1, clientInfo.Buff)
	if err != nil {
		fmt.Println(err)
		return err
	}
	clientInfo.Fmt = bitop.RESERVED_8(clientInfo.Buff[0], 0, 2) // Parse FMT
	//fmt.Println("Now fmt: ", clientInfo.Fmt)
	clientInfo.Csid = uint32(bitop.RESERVED_8(clientInfo.Buff[0], 2, 8))
	//fmt.Println("Now csid: ", clientInfo.Csid)
	switch clientInfo.Csid {
	case 0:
		clientInfo.Buff, err = clientInfo.rw.ReadNByte(1, clientInfo.Buff)
		if err != nil {
			fmt.Println(err)
			return err
		}
		clientInfo.Csid = uint32(clientInfo.Buff[0]) + 64
		//msgHeaderStartPoint = 2
	case 1:
		clientInfo.Buff, err = clientInfo.rw.ReadNByte(2, clientInfo.Buff)
		if err != nil {
			fmt.Println(err)
			return err
		}
		upper := 256 * uint32(clientInfo.Buff[0])
		clientInfo.Csid = upper + uint32(clientInfo.Buff[1]) + 64
		//msgHeaderStartPoint = 3
	default:
		//msgHeaderStartPoint = 1
	}
	//fmt.Println(clientInfo.Buff)
	//fmt.Println(clientInfo.ChunkBasicHeader)
	return nil
}

func ParseMessageHeader(clientInfo *ClientInfo) error {
	var err error
	switch clientInfo.Fmt {
	case 0:
		clientInfo.Buff, err = clientInfo.rw.ReadNByte(11, clientInfo.Buff)
		if err != nil {
			fmt.Println(err)
			return err
		}
		clientInfo.Timestamp = av.Trans3Byte2UINT32(clientInfo.Buff[0:3])
		clientInfo.MessageLength = av.Trans3Byte2UINT32(clientInfo.Buff[3:6])
		clientInfo.MessageTypeID = av.Trans1Byte2UINT8(clientInfo.Buff[6:7])
		clientInfo.MessageStreamID = av.TransInverse4Byte2UINT32(clientInfo.Buff[7:11])
		//timeExtendStartPoint = msgHeaderStartPoint + 11
	case 1:
		clientInfo.Buff, err = clientInfo.rw.ReadNByte(7, clientInfo.Buff)
		if err != nil {
			fmt.Println(err)
			return err
		}
		clientInfo.Timestamp = clientInfo.Timestamp + av.Trans3Byte2UINT32(clientInfo.Buff[0:3])
		clientInfo.MessageLength = av.Trans3Byte2UINT32(clientInfo.Buff[3:6])
		clientInfo.MessageTypeID = av.Trans1Byte2UINT8(clientInfo.Buff[6:7])
		//timeExtendStartPoint = msgHeaderStartPoint + 7
	case 2:
		clientInfo.Buff, err = clientInfo.rw.ReadNByte(3, clientInfo.Buff)
		if err != nil {
			fmt.Println(err)
			return err
		}
		clientInfo.Timestamp = clientInfo.Timestamp + av.Trans3Byte2UINT32(clientInfo.Buff[0:3])
		//timeExtendStartPoint = msgHeaderStartPoint + 3
	case 3:
		//timeExtendStartPoint = msgHeaderStartPoint
	}
	//fmt.Println(clientInfo.Buff)
	//fmt.Println(clientInfo.ChunkMessageHeader)
	return nil
}

func IsExtendedTimeStamp(clientInfo *ClientInfo) bool {
	if clientInfo.Timestamp > 16777215 {
		return true
	} else {
		return false
	}
}

func ParseExtendedTimeStamp(clientInfo *ClientInfo) error {
	var err error
	clientInfo.Buff, err = clientInfo.rw.ReadNByte(4, clientInfo.Buff)
	if err != nil {
		fmt.Println(err)
		return err
	}
	clientInfo.Timestamp = av.Trans4Byte2UINT32(clientInfo.Buff[0:4])
	return nil
}

func CmdHandle(clientInfo *ClientInfo) {
	CmdFactory(clientInfo)
}
