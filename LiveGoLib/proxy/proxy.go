package proxy

type MiddleBuf interface {
	SetSendBytes(buf []byte)
	GetSendBytes() []byte
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

type Queue struct {
	buf  [][]byte
	cap  uint
	head uint
	tail uint
}

func (q *Queue) Enqueue(buff []byte) {
	q.tail = (q.tail + 1) % q.cap
	if q.head == q.tail {
		q.buf[q.head] = buff
		q.head = (q.head + 1) % q.cap
	} else {
		q.buf[q.tail] = buff
	}
}

func (q *Queue) Peek() []byte {
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
		q.buf = make([][]byte, 10)
	}
}

func (fg *FLVGetter) GetSendBytes(buf []byte) []byte {
	return buf
}

func (rs *RTMPSender) SetSendBytes(buf []byte) {
	rs.ChunkQueue.Enqueue(buf)
}

func (prox *ProxyPush) SetSendBytes(buf []byte) {
	prox.rs.SetSendBytes(buf)
}

func (prox *ProxyPush) GetSendBytes() []byte {
	buf := prox.rs.ChunkQueue.Peek()
	return buf
}
