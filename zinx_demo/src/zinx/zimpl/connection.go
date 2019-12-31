package zimpl

import (
	"errors"
	"fmt"
	"io"
	"net"
	"sync"
	. "zinx/util"
	"zinx/zinterface"
)

type Connection struct {
	//Connection隶属于的server对象
	TcpServer zinterface.IServer

	//当前链接的socket
	Conn *net.TCPConn

	//当前链接的ID
	ConnId uint32

	//当前链接的状态
	isClosed bool

	//当前链接绑定的业务处理方法
	//handler zinterface.AbstractHandler

	//告知当前链接已经退出/停止的channel
	MsgChan chan []byte

	//无缓冲管道 用于读/写goroutine直接的通信
	ExitChan chan bool

	//链接的额外属性(metadata)
	property map[string]interface{}

	//链接属性的读写锁
	propertyLock sync.RWMutex

	//处理该链接方法的Router
	//Router zinterface.IRouter
}

func NewConnection(tcpServer zinterface.IServer, conn *net.TCPConn, connId uint32) *Connection {
	c := &Connection{
		TcpServer: tcpServer,
		Conn:      conn,
		ConnId:    connId,
		isClosed:  false,
		ExitChan:  make(chan bool, 1),
		MsgChan:   make(chan []byte),
		property:  make(map[string]interface{}),
	}
	//在新建Connection的时候加入到ConnectionManager中
	c.TcpServer.GetConnectionManager().AddConn(c)
	return c
}

func (c *Connection) SetProperty(key string, value interface{}) {
	c.propertyLock.Lock()
	defer c.propertyLock.Unlock()
	c.property[key] = value
}

func (c *Connection) GetProperty(key string) interface{} {
	c.propertyLock.RLock()
	defer c.propertyLock.RUnlock()
	return c.property[key]
}

func (c *Connection) RemoveProperty(key string) {
	c.propertyLock.Lock()
	defer c.propertyLock.Unlock()
	delete(c.property, key)
}

func (c *Connection) Start() {
	fmt.Println("Starting Conn ID = ", c.ConnId)
	go c.StartReading()
	go c.StartWriting()

	//调用链接创建时的hook函数
	c.TcpServer.CallOnConnStart(c)
}

func (c *Connection) StartWriting() {
	fmt.Println("[writer goroutine is running]")
	defer fmt.Println(c.GetRemoteAddr().String(), "[writer goroutine stopped]")

	//不断阻塞的等待channel的消息 写给客户端
	for true {
		select {
		case data := <-c.MsgChan:
			//如果管道有数据可读 则需要写给客户端
			_, err := c.Conn.Write(data)
			if err != nil {
				fmt.Println("send data error :", err)
				return
			}
		case <-c.ExitChan:
			//如果ExitChan可读 代表reader已经退出 此时writer也需要退出
			return
		}
	}
}
func (c *Connection) StartReading() {
	fmt.Println("[read goroutine is running]")
	defer fmt.Println("Conn ID = ", c.ConnId, " is exit , remote addr is ", c.GetRemoteAddr().String())
	defer c.Stop()

	for true {
		//buf := NewBuffer()
		//count, err := c.Conn.Read(buf)
		//if err != nil {
		//	fmt.Println("read buf err = ", err)
		//	continue
		//}
		//创建dataPack对象
		pack := NewDataPack()
		buffer := NewLenBuffer(pack.GetHeadLen())
		_, err := io.ReadFull(c.Conn, buffer)
		if err != nil {
			fmt.Errorf("server read error %v", err)
			break
		}

		message, _ := pack.Unpack(buffer)

		var data []byte
		if message.GetMsgLen() > 0 {
			data = NewLenBuffer(message.GetMsgLen())
			_, err := io.ReadFull(c.Conn, data)
			if err != nil {
				fmt.Errorf("server read error %v", err)
				break
			}
		}
		message.SetData(data)

		//err = c.handler(c.Conn, buf, count)
		//if err != nil {
		//	fmt.Println("handle err = ", err)
		//}

		//每一个request都会开启一个goroutine 这样会导致goroutine无限增多
		//go c.MsgHandler.ExecuteHandler(NewRequest(c, message))

		if GlobalConfig.WorkerPoolSize > 0 {
			//如果开启了工作池机制 则将消息发送给worker工作池处理即可
			c.TcpServer.GetMsgHandler().SendRequestToTaskQueue(NewRequest(c, message))
		} else {
			//如果没开启 则单独开启goroutine发送request
			go c.TcpServer.GetMsgHandler().ExecuteHandler(NewRequest(c, message))
		}

	}
}

func (c *Connection) Stop() {
	fmt.Println("Stopping Conn ID = ", c.ConnId)
	//如果当前链接已经关闭
	if c.isClosed == true {
		return
	}

	c.isClosed = true

	//调用销毁链接之前的hook函数
	c.TcpServer.CallOnConnStop(c)

	//关闭链接
	_ = c.Conn.Close()

	//告知 writer goroutine 关闭writer
	c.ExitChan <- true

	//当链接关闭的时候从ConnectionManager中删除对应链接
	c.TcpServer.GetConnectionManager().RemoveConn(c)

	//回收资源
	close(c.ExitChan)
	close(c.MsgChan)
}

func (c *Connection) GetTCPConn() *net.TCPConn {
	return c.Conn
}

func (c *Connection) GetConnID() uint32 {
	return c.ConnId
}

func (c *Connection) GetRemoteAddr() net.Addr {
	return c.Conn.RemoteAddr()
}

func (c *Connection) Send(msgId uint32, data []byte) error {
	if c.isClosed {
		return errors.New("connection already closed")
	}
	pack := NewDataPack()
	binaryData, _ := pack.Pack(NewMessage1(msgId, data))
	//读/写 goroutine分离之后 不再在send方法里面写数据 而是将数据写入管道 让 写goroutine发送数据
	//c.Conn.Write(binaryData)
	c.MsgChan <- binaryData
	return nil
}
