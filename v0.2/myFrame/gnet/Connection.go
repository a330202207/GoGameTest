package gnet

import (
	"fmt"
	"gFrame/giface"
	"net"
)

type Connection struct {
	//当前链接的socket TCP套接字
	Conn *net.TCPConn

	//链接的ID
	ConnID uint32

	//当前的连接状态
	isClosed bool

	//当前链接所绑定的处理业务方法API
	handleAPI giface.HandleFunc

	//告知当前链接已经退出的/停止 Channel
	ExitChan chan bool
}

//链接的读业务
func (c *Connection) StartReader() {
	fmt.Println("Reader Goroutine is running...")
	defer fmt.Println("connID=", c.ConnID, "Reader Is Exit,remote addr is", c.RemoteAddr().String())
	defer c.Stop()

	for {
		//读取客户端的数据到buf中，目前最大512字节
		buf := make([]byte, 512)
		cnt, err := c.Conn.Read(buf)
		if err != nil {
			fmt.Println("Recv buf err", err)
			continue
		}

		//调用当前链接所绑定的HandleAPI
		if err := c.handleAPI(c.Conn, buf, cnt); err != nil {
			fmt.Println("ConnID", c.ConnID, "Handle is error", err)
			break
		}
	}
}

//启动链接
func (c *Connection) Start() {
	fmt.Println("Conn Start... ConnID=", c.ConnID)

	//启动从当前链接的读数据的业务
	go c.StartReader()

	//TODO 启动从当前链接写数据的业务

}

//停止链接
func (c *Connection) Stop() {
	fmt.Println("Conn Stop... ConnID=", c.ConnID)

	//如果当前已关闭
	if c.isClosed == true {
		return
	}

	c.isClosed = true

	//关闭socket
	c.Conn.Close()

	//回收资源
	close(c.ExitChan)
}

//获取当前的绑定socket、conn
func (c *Connection) GetTCPConnection() *net.TCPConn {
	return c.Conn
}

//获取当前链接模块的链接ID
func (c *Connection) GetConnID() uint32 {
	return c.ConnID
}

//获取远程客户端的 TCP状态 IP Port
func (c *Connection) RemoteAddr() net.Addr {
	return c.Conn.RemoteAddr()
}

//发送数据，将数据发送给远程的客户端
func (c *Connection) Send(data []byte) error {
	return nil
}

//初始化链接模块的方法
func NewConnection(conn *net.TCPConn, ConnID uint32, callbackAPI giface.HandleFunc) *Connection {
	c := &Connection{
		Conn:      conn,
		ConnID:    ConnID,
		handleAPI: callbackAPI,
		isClosed:  false,
		ExitChan:  make(chan bool, 1),
	}

	return c
}
