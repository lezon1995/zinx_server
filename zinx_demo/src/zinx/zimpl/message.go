package zimpl

type Message struct {
	msgId  uint32
	msgLen uint32
	data   []byte
}

func NewMessage1(msgId uint32, data []byte) *Message {
	return &Message{
		msgId:  msgId,
		msgLen: uint32(len(data)),
		data:   data,
	}
}

func NewMessage() *Message {
	return &Message{}
}

func (m *Message) GetMsgId() uint32 {
	return m.msgId
}

func (m *Message) GetMsgLen() uint32 {
	return m.msgLen
}

func (m *Message) GetData() []byte {
	return m.data
}

func (m *Message) SetMsgId(msgId uint32) {
	m.msgId = msgId
}

func (m *Message) SetMsgLen(msgLen uint32) {
	m.msgLen = msgLen
}

func (m *Message) SetData(data []byte) {
	m.data = data
}
