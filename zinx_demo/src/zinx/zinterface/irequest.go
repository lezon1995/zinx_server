package zinterface

/**
将客户端请求的链接 与 请求数据 包装到一个request中
*/
type IRequest interface {
	//得到当前链接
	GetIConnection() IConnection

	//得到请求的消息数据
	GetData() []byte

	GetMsgId() uint32
}
