package zinterface

/**
链接管理模块抽象层
*/

type IConnManager interface {
	//添加链接
	AddConn(connection IConnection)
	//删除链接
	RemoveConn(connection IConnection)
	//根据connID获取链接
	GetConn(connId uint32) (IConnection, error)
	//得到当前链接总数
	Len() int
	//清除并终止所有链接
	ClearConn()
}
