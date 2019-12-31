package main

import (
	"zinx/zimpl"
	"zinx/zinterface"
)

type PingRouter struct {
	zimpl.BaseRouter
}

func NewPingRouter() *PingRouter {
	return &PingRouter{}
}

func (b *PingRouter) PreHandle(request zinterface.IRequest) {
	request.GetIConnection().GetTCPConn().Write([]byte("before ping \n"))
}

func (b *PingRouter) Handle(request zinterface.IRequest) {
	request.GetIConnection().GetTCPConn().Write([]byte("pinging \n"))
}

func (b *PingRouter) PostHandle(request zinterface.IRequest) {
	request.GetIConnection().GetTCPConn().Write([]byte("after ping \n"))
}

/**
基于Zinx框架来开发 服务端应用程序
*/
func main() {
	//创建一个server句柄
	server := zimpl.NewServer("[zinx V0.3]")
	server.AddRouter(NewPingRouter())
	//启动server
	server.Serve()
}
