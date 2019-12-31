package zimpl

import (
	"zinx/zinterface"
)
/**
	golang 没有抽象类的概念 这里虽然实现IRouter接口
	但是没有具体重写方法 目的是为了让BaseRouter的子类去选择性实现 体现抽象类的概念
 */
type BaseRouter struct {
}

func (b *BaseRouter) PreHandle(request zinterface.IRequest) {}

func (b *BaseRouter) Handle(request zinterface.IRequest) {}

func (b *BaseRouter) PostHandle(request zinterface.IRequest) {}
