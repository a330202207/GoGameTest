package gnet

import (
	"fmt"
	"gFrame/giface"
	"github.com/go-errors/errors"
	"sync"
)

type ConnManager struct {
	//管理的链接信息
	connections map[uint32]giface.IConnection

	//保护链接集合的读写锁
	connLock sync.RWMutex
}

//创建当前链接方法
func NewConnManage() *ConnManager {
	return &ConnManager{
		connections: make(map[uint32]giface.IConnection),
	}
}

//添加链接
func (connMgr *ConnManager) Add(conn giface.IConnection) {
	//保护共享资源map，加写锁
	connMgr.connLock.Lock()
	defer connMgr.connLock.Unlock()

	//将Conn加入ConnManage
	connMgr.connections[conn.GetConnID()] = conn
	fmt.Println("ConnID = ", conn.GetConnID(), " Connection Add To ConnManager Successfully:Conn Num = ", connMgr.Len())
}

//删除链接
func (connMgr *ConnManager) Remove(conn giface.IConnection) {
	//保护共享资源map，加写锁
	connMgr.connLock.Lock()
	defer connMgr.connLock.Unlock()

	//删除链接信息
	delete(connMgr.connections, conn.GetConnID())
	fmt.Println("ConnID = ", conn.GetConnID(), " Connection Remove To ConnManager Successfully:Conn Num = ", connMgr.Len())
}

//根据ConnID获取链接
func (connMgr *ConnManager) Get(connID uint32) (giface.IConnection, error) {
	//保护共享资源map，加读锁
	connMgr.connLock.RLock()
	defer connMgr.connLock.RUnlock()

	if conn, ok := connMgr.connections[connID]; ok {
		return conn, nil
	} else {
		return nil, errors.New("Connection Not Fount")
	}
}

//总链接个数
func (connMgr *ConnManager) Len() int {
	return len(connMgr.connections)
}

//清除并终止所有的链接
func (connMgr *ConnManager) ClearConn() {
	//保护共享资源map，加写锁
	connMgr.connLock.Lock()
	defer connMgr.connLock.Unlock()

	//删除Conn并停止Conn的工作
	for connID, conn := range connMgr.connections {
		//停止
		conn.Stop()
		//删除
		delete(connMgr.connections, connID)
		fmt.Println("Clear All Connections Success! Conn Num = ", connMgr.Len(), "ConnID = ", conn.GetConnID())
	}
}
