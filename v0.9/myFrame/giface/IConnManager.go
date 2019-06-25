package giface

//定义一个链接管理接口
type IConnManager interface {
	//添加链接
	Add(conn IConnection)

	//删除链接
	Remove(conn IConnection)

	//根据ConnID获取链接
	Get(connID uint32) (IConnection, error)

	//总链接个数
	Len() int

	//清除并终止所有的链接
	ClearConn()
}
