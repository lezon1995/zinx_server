package util

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"zinx/zinterface"
)

/**
存储一切有关zinx框架的参数
*/
type Config struct {
	//当前ZINX全局的server对象
	TcpServer zinterface.IServer
	//当前服务器主机监听的IP
	Host string
	//当前服务器监听的端口号
	Port int
	//当前服务器的名称
	ServerName string

	//ZINX版本号
	Version string
	//最大链接数
	MaxConn int
	//数据包最大长度
	BufferSize uint32
	//工作池worker数量
	WorkerPoolSize uint32
	//任务队列长度
	TaskQueueSize uint32
}

/**
定义一个全局的Config对象
*/
var GlobalConfig *Config

/**
导包的时候会执行init方法
*/
func init() {
	GlobalConfig = &Config{
		TcpServer:      nil,
		Host:           "0.0.0.0",
		Port:           8999,
		ServerName:     "ZinxServer",
		Version:        "V0.4",
		MaxConn:        2,
		BufferSize:     1024,
		WorkerPoolSize: 10,
		TaskQueueSize:  1024,
	}
	//GlobalConfig.Reload()
}

func (g *Config) Reload() {
	data, err := ioutil.ReadFile("src/demo/v4/conf/zinx.json")
	if err != nil {
		fmt.Println("read file error : ", err)
	}
	err = json.Unmarshal(data, &GlobalConfig)
	if err != nil {
		fmt.Println("unmarshal json error : ", err)
	}
}
