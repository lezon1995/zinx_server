package zimpl

import (
	"errors"
	"fmt"
	"sync"
	"zinx/zinterface"
)

/**
链接管理模块
负责链接的增删查
*/
type ConnManager struct {
	//链接集合
	connections map[uint32]zinterface.IConnection
	//保护链接的读写锁
	connLock sync.RWMutex
}

func NewConnManager() *ConnManager {
	return &ConnManager{connections: make(map[uint32]zinterface.IConnection)}
}

func (cm *ConnManager) AddConn(connection zinterface.IConnection) {
	//由于同时间会有大量请求一起进来 所以需要对链接集合加锁 避免并发操作问题
	cm.connLock.Lock()
	defer cm.connLock.Unlock()

	//将connection加入到链接map集合中
	cm.connections[connection.GetConnID()] = connection

	fmt.Println("connID =", connection.GetConnID(), " add to ConnManager successfully , current connections size = ", cm.Len())
}

func (cm *ConnManager) RemoveConn(connection zinterface.IConnection) {
	//由于同时间会有大量请求一起进来 所以需要对链接集合加锁 避免并发操作问题
	cm.connLock.Lock()
	defer cm.connLock.Unlock()

	//删除链接信息
	delete(cm.connections, connection.GetConnID())

	fmt.Println("connID =", connection.GetConnID(), " remove successfully , current connections size = ", cm.Len())
}

func (cm *ConnManager) GetConn(connId uint32) (zinterface.IConnection, error) {
	cm.connLock.RLock()
	defer cm.connLock.RUnlock()
	connection, ok := cm.connections[connId]
	if ok {
		return connection, nil
	} else {
		return nil, errors.New(fmt.Sprintf("connID = %d do not exist !", connId))
	}
}

func (cm *ConnManager) Len() int {
	return len(cm.connections)
}

func (cm *ConnManager) ClearConn() {
	//由于同时间会有大量请求一起进来 所以需要对链接集合加锁 避免并发操作问题
	cm.connLock.Lock()
	defer cm.connLock.Unlock()

	for connID, connection := range cm.connections {
		//停止connection工作
		connection.Stop()
		//删除connID对应地connection
		delete(cm.connections, connID)
	}

	fmt.Println("Clear all connections successfully !  current connections size = ", cm.Len())
}
