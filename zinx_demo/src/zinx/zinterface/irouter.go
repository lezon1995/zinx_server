package zinterface

/**
路由抽象借口
路由里的数据全是IRequest
*/
type IRouter interface {
	PreHandle(request IRequest)
	Handle(request IRequest)
	PostHandle(request IRequest)
}
