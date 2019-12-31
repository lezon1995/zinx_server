package zimpl

import (
	"fmt"
	"net"
	"zinx/util"
	"zinx/zinterface"
)

type Server struct {
	//服务名称
	Name string
	//服务器ip版本
	IPVersion string
	//服务监听的ip
	IP string
	//服务监听的端口
	Port int
	//处理该链接方法的Router
	//Router zinterface.IRouter
	MsgHandler zinterface.IMsgHandler

	ConnectionManager zinterface.IConnManager

	//该server创建链接之后 自动调用的hook函数
	OnConnStart func(connection zinterface.IConnection)
	//该server销毁链接之后 自动调用的hook函数
	OnConnStop func(connection zinterface.IConnection)
}

func (s *Server) SetOnConnStart(function func(connection zinterface.IConnection)) {
	s.OnConnStart = function
}

func (s *Server) SetOnConnStop(function func(connection zinterface.IConnection)) {
	s.OnConnStop = function
}

func (s *Server) CallOnConnStart(connection zinterface.IConnection) {
	if s.OnConnStart == nil {
		return
	}
	fmt.Println("<-----Calling OnConnStart----->")
	s.OnConnStart(connection)
	fmt.Println("<-----Ending OnConnStart----->")
}

func (s *Server) CallOnConnStop(connection zinterface.IConnection) {
	if s.OnConnStop == nil {
		return
	}
	fmt.Println("<-----Calling OnConnStop----->")
	s.OnConnStop(connection)
	fmt.Println("<-----Ending OnConnStop----->")
}

func (s *Server) GetMsgHandler() zinterface.IMsgHandler {
	return s.MsgHandler
}

func (s *Server) GetConnectionManager() zinterface.IConnManager {
	return s.ConnectionManager
}

func NewServer(name string) *Server {
	return &Server{
		Name:              util.GlobalConfig.ServerName,
		IPVersion:         "tcp4",
		IP:                util.GlobalConfig.Host,
		Port:              util.GlobalConfig.Port,
		MsgHandler:        NewMsgHandler(),
		ConnectionManager: NewConnManager(),
	}
}

func (s *Server) AddRouter(msgId uint32, router zinterface.IRouter) {
	s.MsgHandler.AddRouter(msgId, router)
}

func (s *Server) Start() {
	//开启任务队列和worker工作池
	go s.MsgHandler.StartWorkerPool()

	fmt.Printf("[STARTING] Server [%s] Listener at IP:[%s] Port:[%d] is starting . . . \n", s.Name, s.IP, s.Port)
	//获取一个TCP的addr
	addr, err := net.ResolveTCPAddr(s.IPVersion, fmt.Sprintf("%s:%d", s.IP, s.Port))
	if err != nil {
		fmt.Println("resolve tcp addr error : ", err)
		return
	}
	//监听服务器的地址
	listener, err := net.ListenTCP(s.IPVersion, addr)
	if err != nil {
		fmt.Println("listen ", s.IPVersion, " err ", err)
		return
	}
	fmt.Printf("[STARTED] Server [%s] Listener at IP:[%s] Port:[%d] is successful !\n", s.Name, s.IP, s.Port)
	//阻塞等待客户端连接 处理客户端业务
	var cid uint32
	cid = 0
	for true {
		//util.PrintWaiting()
		fmt.Println("listening . . . ")
		//如果有客户端链接过来 阻塞会返回
		conn, err := listener.AcceptTCP()
		if err != nil {
			fmt.Println("Accept connection error : ", err)
			continue
		}

		//链接之前判断当前连接数是否超过最大值
		if s.ConnectionManager.Len() >= util.GlobalConfig.MaxConn {
			//TODO 给客户端响应连接失败信息 超出最大连接数
			fmt.Println("----------------------------------------------------------------")
			_ = conn.Close()
			continue
		}

		//已经与客户端建立连接
		newConn := NewConnection(s, conn, cid)

		cid++
		newConn.Start()
	}
}

/*func handle(conn *net.TCPConn, msg []byte, count int) error {
	fmt.Println("server receive : ", string(msg[:count]))
	if _, err := conn.Write(msg[:count]); err != nil {
		fmt.Println("write error : ", err)
		return errors.New("handler error")
	}
	return nil
}*/

func (s *Server) Stop() {
	fmt.Println("[STOP] Zinx server name = ", s.Name)
	//将服务器的资源 状态 或者一些已经开启的链接 进行停止或者回收
	s.ConnectionManager.ClearConn()
}

func (s *Server) Serve() {
	//启动服务
	go s.Start()

	//TODO 做一些启动服务器之后的额外业务

	//阻塞
	select {}
}
