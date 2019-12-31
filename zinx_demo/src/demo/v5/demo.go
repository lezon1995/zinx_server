package main

import (
	"fmt"
	"zinx/zimpl"
	"zinx/zinterface"
)

type PingRouter struct {
	zimpl.BaseRouter
}

func NewPingRouter() *PingRouter {
	return &PingRouter{}
}

var count uint32 = 1

func (b *PingRouter) Handle(request zinterface.IRequest) {
	fmt.Println("--->server receive MsgId=", request.GetMsgId(), "msg=", string(request.GetData()))
	request.GetIConnection().Send(count, []byte(fmt.Sprintf("hello from server [%d]", count)))
	count++
}

/**
基于Zinx框架来开发 服务端应用程序
*/
func main() {
	//创建一个server句柄
	server := zimpl.NewServer("[zinx V0.5]")
	//server.AddRouter(NewPingRouter())
	//启动server
	server.Serve()
}
