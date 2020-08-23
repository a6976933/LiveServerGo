package proxy

import (
	"sync"
)
type MiddleBuf interface {
	SetSendInfo(buf []byte,tagtype int,msgLen uint32,time uint32,id uint32)
	GetSendInfo() ([]byte, int,uint32, uint32, uint32)
	SendMetaData() []byte
	SetMetaData([]byte)
	SetAVCData ([]byte)
	GetAVCData () []byte
	SetAACData ([]byte)
	GetAACData () []byte
	GetClientCnt () uint32
	SetClientCnt (bool)
	GetSChan  () chan struct{}
}

type FLVGetter struct {
	Buff []byte
}

type RTMPSender struct {

}

type ProxyPush struct {
	rs *RTMPSender
	fg *FLVGetter
	si SendInfo
	sendChan chan struct{}
	clientCnt  uint32
	cCntm    sync.Mutex
	metaData []byte
	avcData  []byte
	aacData  []byte
}

type SendInfo struct {
	buf []byte
	tagtype int
	msgLen uint32
	time uint32
	streamID uint32
}

func NewProxyPush() *ProxyPush {
	var ret = new(ProxyPush)
	ret.sendChan = make(chan struct{})
	return ret
}

func (prox *ProxyPush) SetSendInfo(buf []byte,tagtype int,msgLen uint32,time uint32,streamID uint32) {
	var ob SendInfo
	ob.buf = buf
	ob.tagtype = tagtype
	ob.msgLen = msgLen
	ob.time = time
	ob.streamID = streamID
	prox.si = ob
}

func (prox *ProxyPush) GetSendInfo() ([]byte,int,uint32,uint32,uint32) {
	buf := prox.si
	return buf.buf, buf.tagtype,buf.msgLen, buf.time,buf.streamID
}

func (prox *ProxyPush) SetMetaData(meta []byte){
	prox.metaData = meta
}

func (prox *ProxyPush) SendMetaData() []byte {
	return prox.metaData
}

func (prox *ProxyPush) SetAVCData(data []byte) {
	prox.avcData = data
}

func (prox *ProxyPush) GetAVCData () []byte {
	return prox.avcData
}

func (prox *ProxyPush) SetAACData(data []byte) {
	prox.aacData = data
}

func (prox *ProxyPush) GetAACData () []byte {
	return prox.aacData
}

func (prox *ProxyPush) GetClientCnt () uint32 {
	return prox.clientCnt
}

func (prox *ProxyPush) SetClientCnt(op bool) {
	if op {
		prox.cCntm.Lock()
		prox.clientCnt++
		prox.cCntm.Unlock()
	} else {
		prox.cCntm.Lock()
		prox.clientCnt--
		prox.cCntm.Unlock()
	}
}

func (prox *ProxyPush) GetSChan  () chan struct{} {
	return prox.sendChan
}
