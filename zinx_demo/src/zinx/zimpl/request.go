package zimpl

import "zinx/zinterface"

type Request struct {
	//已经和客户端建立好的链接
	conn zinterface.IConnection
	//客户端请求数据
	msg zinterface.IMessage
}

func NewRequest(conn zinterface.IConnection, msg zinterface.IMessage) *Request {
	return &Request{
		conn: conn,
		msg:  msg,
	}
}

func (r *Request) GetIConnection() zinterface.IConnection {
	return r.conn
}

func (r *Request) GetData() []byte {
	return r.msg.GetData()
}

func (r *Request) GetMsgId() uint32 {
	return r.msg.GetMsgId()
}
