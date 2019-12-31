package zinterface

type IMsgHandler interface {
	//调度/执行 对应的router消息处理方法
	ExecuteHandler(request IRequest)
	//为消息类型添加对应路由
	AddRouter(msgId uint32, router IRouter)
	//启动工作池(全局唯一)
	StartWorkerPool()

	//将request发送给任务队列
	SendRequestToTaskQueue(request IRequest)
}
