package gnet

import (
	"fmt"
	"gFrame/giface"
	"github.com/pkg/errors"
	"net"
)

//IServer的接口实现，定义一个Server的服务器模块
type Server struct {
	//服务器名称
	Name string
	//服务器绑定IP版本
	IPVersion string
	//服务器监听IP
	IP string
	//服务器监听端口号
	Port int
}

//定义当前客户度链接的所绑定 handle api（目前这个handle是写死的，以后优化）
func CallBackClient(conn *net.TCPConn, data []byte, cnt int) error {
	fmt.Println("Conn Handle CallBackToClient...")

	if _, err := conn.Write(data[:cnt]); err != nil {
		fmt.Println("Write Back Buf Error", err)
		return errors.New("CallBackToClient Error")
	}

	return nil
}

//启动服务器
func (s *Server) Start() {
	fmt.Printf("[Start] Server Listenner At IP :%s, Port:%d, is starting\n", s.IP, s.Port)

	go func() {
		//获取一个TCP的Addr
		addr, err := net.ResolveTCPAddr(s.IPVersion, fmt.Sprintf("%s:%d", s.IP, s.Port))
		if err != nil {
			fmt.Println("Resolve tcp addr error:", err)
			return
		}

		//监听服务器地址
		Listenner, err := net.ListenTCP(s.IPVersion, addr)
		if err != nil {
			fmt.Println("Listen error:", s.IPVersion, "error", err)
			return
		}

		fmt.Println("Start Server Success", s.Name, "Success Listening...")

		var cid uint32
		cid = 0

		//阻塞的等待客户端链接，出来客户端链接业务（读写）
		for {

			//如果有客户端链接过来，阻塞会返回
			conn, err := Listenner.AcceptTCP()
			if err != nil {
				fmt.Println("Accept error", err)
				continue
			}

			//将处理新链接的业务方法和 conn 进行绑定，得到链接模块
			dealConn := NewConnection(conn, cid, CallBackClient)
			cid++

			//启动当前链接业务处理
			go dealConn.Start()
		}
	}()
}

//停止服务器
func (s *Server) Stop() {
	//TODO 将一些服务器的资源、状态或者一些已经开辟的链接信息 进行停止或者回收
}

//运行服务器
func (s *Server) Server() {
	//启动Server服务器
	s.Start()

	//TODO 做一些启动服务器之前后的额外业务

	//阻塞状态
	select {}
}

//初始化
func NewServer(name string) giface.IServer {
	s := &Server{
		Name:      name,
		IPVersion: "tcp4",
		IP:        "0.0.0.0",
		Port:      8999,
	}
	return s
}
