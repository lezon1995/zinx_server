package zimpl

import (
	"fmt"
	"strconv"
	"zinx/util"
	"zinx/zinterface"
)

type MsgHandler struct {
	//存放msgId和IRouter的对应关系
	Apis map[uint32]zinterface.IRouter

	//负责消息的任务队列
	TaskQueue []chan zinterface.IRequest

	//worker工作池的worker数量
	WorkPoolSize uint32
}

func (m *MsgHandler) SendRequestToTaskQueue(request zinterface.IRequest) {
	//使用负载均衡的策略(轮训)将消息分配给不同的worker
	workerId := request.GetIConnection().GetConnID() % m.WorkPoolSize
	fmt.Println("add ConnID = [", request.GetIConnection().GetConnID(), "] request MsgID = [", request.GetMsgId(), "] to workerID = [", workerId, "]")
	//向workerId对应的channel发送request
	m.TaskQueue[workerId] <- request
}

func NewMsgHandler() *MsgHandler {
	return &MsgHandler{
		Apis:         make(map[uint32]zinterface.IRouter),
		WorkPoolSize: util.GlobalConfig.WorkerPoolSize,
		TaskQueue:    make([]chan zinterface.IRequest, util.GlobalConfig.WorkerPoolSize),
	}
}

func (m *MsgHandler) ExecuteHandler(request zinterface.IRequest) {
	router, ok := m.Apis[request.GetMsgId()]
	if !ok {
		fmt.Println("router msgId = ", request.GetMsgId(), "do not exit !")
		return
	}
	router.PreHandle(request)
	router.Handle(request)
	router.PostHandle(request)
}

func (m *MsgHandler) AddRouter(msgId uint32, router zinterface.IRouter) {
	//判断当前msgId是否已经绑定router
	if _, ok := m.Apis[msgId]; ok {
		//msgId已经注册过
		panic("repeat api , msgId=" + strconv.Itoa(int(msgId)))
	}
	m.Apis[msgId] = router
	fmt.Println("add router msgId = ", msgId, "success !")
}

//启动一个worker工作池(开启工作池的动作只能发生一次 一个zinx框架只能有一个worker工作池)
func (m *MsgHandler) StartWorkerPool() {
	//根据workerPoolSize分别开启worker 每一个worker用一个goroutine来重载
	for i := 0; i < int(m.WorkPoolSize); i++ {
		//当前的worker对应的channel消息队列 开辟空间
		m.TaskQueue[i] = make(chan zinterface.IRequest, util.GlobalConfig.TaskQueueSize)
		//启动当前worker
		go m.startOneWorker(i, m.TaskQueue[i])
	}
}

//启动一个worker协程
func (m *MsgHandler) startOneWorker(workerId int, taskQueue chan zinterface.IRequest) {
	fmt.Println("worker ID = ", workerId, "is started . . .")
	//不断阻塞监听任务队列中的消息
	for true {
		select {
		//如果有消息过来 出列的就是一个来自客户端的request 使用当前request中msgID所绑定的router来处理request
		case requests := <-taskQueue:
			m.ExecuteHandler(requests)
		}
	}
}
