package proxy

type MiddleBuf interface {
	SetSendInfo(buf []byte,tagtype int,msgLen uint32,time uint32,id uint32)
	GetSendInfo() ([]byte, int,uint32, uint32, uint32)
}

type FLVGetter struct {
	Buff []byte
}

type RTMPSender struct {
	ChunkQueue *Queue
}

type ProxyPush struct {
	rs *RTMPSender
	fg *FLVGetter
}

type SendInfo struct {
	buf []byte
	tagtype int
	msgLen uint32
	time uint32
	streamID uint32
}

type Queue struct {
	buf  []SendInfo
	cap  uint
	head uint
	tail uint
}

func (q *Queue) Enqueue(buff SendInfo) {
	q.tail = (q.tail + 1) % q.cap
	if q.head == q.tail {
		q.buf[q.head] = buff
		q.head = (q.head + 1) % q.cap
	} else {
		q.buf[q.tail] = buff
	}
}

func (q *Queue) Peek() SendInfo {
	return q.buf[q.head]
}

func NewProxyPush(cap uint) *ProxyPush {
	var ret = new(ProxyPush)
	ret.rs = NewRTMPSender()
	ret.fg = NewFLVGetter()
	ret.rs.ChunkQueue.SetCap(cap)
	return ret
}

func NewRTMPSender() *RTMPSender {
	ret := &RTMPSender{
		ChunkQueue: NewQueue(),
	}
	return ret
}

func NewFLVGetter() *FLVGetter {
	return new(FLVGetter)
}

func NewQueue() *Queue {
	return new(Queue)
}

func (q *Queue) SetCap(cp uint) {
	q.cap = cp
	for i := 0; i < int(cp); i++ {
		q.buf = make([]SendInfo, 10)
	}
}

func (fg *FLVGetter) GetSendBytes(buf []byte) []byte {
	return buf
}

func (rs *RTMPSender) SetSendInfo(ob SendInfo) {
	rs.ChunkQueue.Enqueue(ob)
}

func (prox *ProxyPush) SetSendInfo(buf []byte,tagtype int,msgLen uint32,time uint32,streamID uint32) {
	var ob SendInfo
	ob.buf = buf
	ob.tagtype = tagtype
	ob.msgLen = msgLen
	ob.time = time
	ob.streamID = streamID
	prox.rs.SetSendInfo(ob)
}

func (prox *ProxyPush) GetSendInfo() ([]byte,int,uint32,uint32,uint32) {
	buf := prox.rs.ChunkQueue.Peek()
	return buf.buf, buf.tagtype,buf.msgLen, buf.time,buf.streamID
}
