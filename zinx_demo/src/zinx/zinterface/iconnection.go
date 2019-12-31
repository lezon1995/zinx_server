package zinterface

import "net"

type IConnection interface {
	//启动链接 让当前的链接准备开始工作
	Start()
	//停止链接 结束当前链接的工作
	Stop()
	//获取当前链接绑定的socket
	GetTCPConn() *net.TCPConn
	//获取当前链接模块的ID
	GetConnID() uint32
	//获取客户端的TCP状态 IP Port
	GetRemoteAddr() net.Addr
	//发送数据
	Send(msgId uint32, data []byte) error

	//设置链接属性
	SetProperty(key string, value interface{})
	//获取链接属性
	GetProperty(key string) interface{}
	//移除链接属性
	RemoveProperty(key string)
}

//抽象业务处理方法
type AbstractHandler func(conn *net.TCPConn, data []byte, count int) error
