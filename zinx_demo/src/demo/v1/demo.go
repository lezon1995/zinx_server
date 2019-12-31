package main

import "zinx/zimpl"

/**
基于Zinx框架来开发 服务端应用程序
*/
func main() {
	//创建一个server句柄
	server := zimpl.NewServer("[zinx V0.1]")
	//启动server
	server.Serve()

}
