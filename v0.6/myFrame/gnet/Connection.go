package gnet

import (
	"fmt"
	"gFrame/giface"
	"github.com/pkg/errors"
	"io"
	"net"
)

type Connection struct {
	//当前链接的socket TCP套接字
	Conn *net.TCPConn

	//链接的ID
	ConnID uint32

	//当前的连接状态
	isClosed bool

	//告知当前链接已经退出的/停止 Channel
	ExitChan chan bool

	//该链接处理的方法Router
	Router giface.IRouter
}

//链接的读业务
func (c *Connection) StartReader() {
	fmt.Println(" Reader Goroutine Is Running...")
	defer fmt.Println("ConnID = ", c.ConnID, " Reader Is Exit, Remote Addr Is ", c.RemoteAddr().String())
	defer c.Stop()

	for {

		//创建一个拆包解包对象
		dp := NewDataPack()

		//读取客户端的Msg Head 二进制流 8个字节
		headData := make([]byte, dp.GetHeadLen())
		if _, err := io.ReadFull(c.GetTCPConnection(), headData); err != nil {
			fmt.Println("Read msg Head Error", err)
			break
		}

		//拆包，得到MsgID和 MsgDataLen，放在Msg消息中
		msg, err := dp.Unpack(headData)
		if err != nil {
			fmt.Println("Unpack Error", err)
			break
		}

		//根据DataLen，再次读取Data，放在Msg.Data中
		var data []byte
		if msg.GetMsgLen() > 0 {
			data = make([]byte, msg.GetMsgLen())
			if _, err := io.ReadFull(c.GetTCPConnection(), data); err != nil {
				fmt.Println("Read Msg Data Error", err)
				break
			}
		}

		msg.SetData(data)

		//得到当前Conn数据的Request请求数据
		req := Request{
			conn: c,
			msg:  msg,
		}

		//执行注册的路由方法
		go func(request giface.IRequest) {
			//从路由中，找到注册绑定的Conn对应的Router调用
			c.Router.PreHandle(request)
			c.Router.Handle(request)
			c.Router.PostHandle(request)
		}(&req)

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

//提供一个SendMsg方法 将我们要发送给客户端的数据，先进行封包，再发送
func (c *Connection) SendMsg(MsgId uint32, data []byte) error {
	if c.isClosed == true {
		return errors.New("Connection closed when send msg")
	}

	//将Data进行封包 MsgDataLen|Data
	dp := NewDataPack()

	binaryMsg, err := dp.Pack(NewMsgPackage(MsgId, data))
	if err != nil {
		fmt.Println("Pack Error Msg Id=", MsgId)
		return errors.New("Pack Error msg")
	}

	//将数据发送给客户端
	if _, err := c.Conn.Write(binaryMsg); err != nil {
		fmt.Println("Write Msg Id", MsgId, ",Error:", err)
		return errors.New("Conn Write Error")
	}

	return nil
}

//初始化链接模块的方法
func NewConnection(conn *net.TCPConn, ConnID uint32, router giface.IRouter) *Connection {
	c := &Connection{
		Conn:     conn,
		ConnID:   ConnID,
		Router:   router,
		isClosed: false,
		ExitChan: make(chan bool, 1),
	}

	return c
}
