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

	//告知当前链接已经退出的/停止 Channel(由Reader告知Writer退出)
	ExitChan chan bool

	//用于无缓冲的管道，用于读写Goroutine直接的消息通信
	msgChan chan []byte

	//消息的管理MsgID和对应的处理业务API关系
	MsgHandler giface.IMsgHandler
}

//链接的读业务
func (c *Connection) StartReader() {
	fmt.Println(" [Reader Goroutine Is Running!]")
	defer fmt.Println("[Reader Is Exit!], ConnID = ", c.ConnID, " Remote Addr Is ", c.RemoteAddr().String())
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

		//从路由中，找到注册绑定的Conn对应的Router调用
		//根据绑定好的MsgID找到对应处理API业务，执行
		go c.MsgHandler.DoMsgHandler(&req)
	}
}

//写消息Goroutine，专门发送给客户端消息的模块
func (c *Connection) StartWriter() {
	fmt.Println("[Writer Goroutine Is Running]")
	defer fmt.Println("[Conn Writer Exit!]", c.RemoteAddr().String())

	for {
		select {
		case data := <-c.msgChan:
			//有数据要写给客户端
			if _, err := c.Conn.Write(data); err != nil {
				fmt.Println("Send Data Error,", err)
				return
			}
		case <-c.ExitChan:
			//代表Reader已经退出，此时Writer也要退出
			return
		}

	}
}

//启动链接
func (c *Connection) Start() {
	fmt.Println("Conn Start... ConnID=", c.ConnID)

	//启动从当前链接的读数据的业务
	go c.StartReader()

	//启动从当前链接写数据的业务
	go c.StartWriter()

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

	//告知Writer关闭
	c.ExitChan <- true

	//回收资源
	close(c.ExitChan)
	close(c.msgChan)
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

	//将数据发送给Channel
	c.msgChan <- binaryMsg

	return nil
}

//初始化链接模块的方法
func NewConnection(conn *net.TCPConn, ConnID uint32, msgHandler giface.IMsgHandler) *Connection {
	c := &Connection{
		Conn:       conn,
		ConnID:     ConnID,
		MsgHandler: msgHandler,
		isClosed:   false,
		ExitChan:   make(chan bool, 1),
		msgChan:    make(chan []byte),
	}

	return c
}
