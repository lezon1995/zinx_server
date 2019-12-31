package main

import (
	"fmt"
	"zinx/zimpl"
	"zinx/zinterface"
)

//------------------------------------------------------------------------------------------------------------------------------------------------------
type PingRouter struct {
	zimpl.BaseRouter
}

func NewPingRouter() *PingRouter {
	return &PingRouter{}
}

var count uint32 = 1

func (b *PingRouter) Handle(request zinterface.IRequest) {
	fmt.Println(request.GetIConnection().GetRemoteAddr().String(), "--->server receive MsgId=", request.GetMsgId(), "msg=", string(request.GetData()))
	request.GetIConnection().Send(request.GetMsgId(), []byte(fmt.Sprintf("--------ping from server [%d]", count)))
	count++
}

//------------------------------------------------------------------------------------------------------------------------------------------------------

type HelloRouter struct {
	zimpl.BaseRouter
}

func NewHelloRouter() *HelloRouter {
	return &HelloRouter{}
}

func (b *HelloRouter) Handle(request zinterface.IRequest) {
	fmt.Println(request.GetIConnection().GetRemoteAddr().String(), "--->server receive MsgId=", request.GetMsgId(), "msg=", string(request.GetData()))
	request.GetIConnection().Send(request.GetMsgId(), []byte(fmt.Sprintf("=====hello from server")))
}

//------------------------------------------------------------------------------------------------------------------------------------------------------
/**
基于Zinx框架来开发 服务端应用程序
*/
func main() {
	//创建一个server句柄
	server := zimpl.NewServer("[zinx V0.9]")

	server.SetOnConnStart(func(connection zinterface.IConnection) {
		fmt.Println("===========CONNECTION BEGIN=============")
	})

	server.SetOnConnStop(func(connection zinterface.IConnection) {
		fmt.Println("===========CONNECTION STOP=============")
	})

	server.AddRouter(1, NewPingRouter())
	server.AddRouter(2, NewHelloRouter())
	//启动server
	server.Serve()
}
