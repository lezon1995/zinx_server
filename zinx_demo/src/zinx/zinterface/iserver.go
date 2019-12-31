package zinterface

type IServer interface {
	//启动服务
	Start()
	//停止服务
	Stop()
	//运行服务
	Serve()
	//添加Router
	AddRouter(msgId uint32, router IRouter)

	GetConnectionManager() IConnManager

	GetMsgHandler() IMsgHandler

	//设置OnConnStart的hook函数
	SetOnConnStart(func(connection IConnection))
	//设置OnConnStop的hook函数
	SetOnConnStop(func(connection IConnection))
	//调用OnConnStart的hook函数
	CallOnConnStart(connection IConnection)
	//调用OnConnStop的hook函数
	CallOnConnStop(connection IConnection)
}
