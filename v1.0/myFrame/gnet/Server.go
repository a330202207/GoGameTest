package gnet

import (
	"fmt"
	"gFrame/giface"
	"gFrame/utils"
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

	//当前Server的消息管理模块，用来绑定MsgID和对应的处理业务关系
	MsgHandler giface.IMsgHandler

	//该Server的链接管理器
	ConnMgr giface.IConnManager

	//该Server创建链接之后自动调用Hook函数
	OnConnStart func(conn giface.IConnection)

	//该Server销毁链接之前自动调用Hook函数
	OnConnStop func(conn giface.IConnection)
}

//启动服务器
func (s *Server) Start() {
	fmt.Printf("[Ginx] Server Name :%s, listenner at ip :%s, Port:%d is staring\n",
		utils.GlobalObject.Name, utils.GlobalObject.Host, utils.GlobalObject.TcpPort)
	fmt.Printf("[Ginx] Version :%s, MaxConn:%d, MaxPackageSize:%d\n",
		utils.GlobalObject.Version,
		utils.GlobalObject.MaxConn,
		utils.GlobalObject.MaxPackageSize)

	go func() {

		//开启消息队列及Worker工作池
		s.MsgHandler.StartWorkerPool()

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

			//判断当前的链接数量已经超出最大值MaxConn，如果超过，则关闭新链接
			if s.ConnMgr.Len() >= utils.GlobalObject.MaxConn {
				//TODO 给客户端相应一个错误超出最大连接的错误包
				fmt.Println("================Too Many Connections MaxConn = ", utils.GlobalObject.MaxConn, "================")
				conn.Close()
				continue
			}

			//将处理新链接的业务方法和 conn 进行绑定，得到链接模块
			dealConn := NewConnection(s, conn, cid, s.MsgHandler)
			cid++

			//启动当前链接业务处理
			go dealConn.Start()
		}
	}()
}

//停止服务器
func (s *Server) Stop() {
	fmt.Println("[Stop] Ginx Server name", s.Name)
	//将一些资源进行回收
	s.ConnMgr.ClearConn()
}

//运行服务器
func (s *Server) Server() {
	//启动Server服务器
	s.Start()

	//TODO 做一些启动服务器之前后的额外业务

	//阻塞状态
	select {}
}

//路由功能，给当前的服务注册一个路由方法，供客户端的链接处理使用
func (s *Server) AddRouter(MsgID uint32, router giface.IRouter) {
	s.MsgHandler.AddRouter(MsgID, router)
	fmt.Println("Add Router Success!")
}

//获取当前Server的链接管理器
func (s *Server) GetConnMgr() giface.IConnManager {
	return s.ConnMgr
}

//注册OnConnStart钩子方法
func (s *Server) SetOnConnStart(hookFunc func(connection giface.IConnection)) {
	s.OnConnStart = hookFunc
}

//注册OnConnStop钩子方法
func (s *Server) SetOnConnStop(hookFunc func(connection giface.IConnection)) {
	s.OnConnStop = hookFunc
}

//调用OnConnStart钩子方法
func (s *Server) CallOnConnStart(conn giface.IConnection) {
	if s.OnConnStop != nil {
		fmt.Println("---------->Call OnConnStart()...")
		s.OnConnStart(conn)
	}
}

//调用OnConnStop钩子方法
func (s *Server) CallOnConnStop(conn giface.IConnection) {
	if s.OnConnStop != nil {
		fmt.Println("---------->Call OnConnStop()...")
		s.OnConnStop(conn)
	}
}

//初始化
func NewServer(name string) giface.IServer {
	s := &Server{
		Name:       utils.GlobalObject.Name,
		IPVersion:  "tcp4",
		IP:         utils.GlobalObject.Host,
		Port:       utils.GlobalObject.TcpPort,
		MsgHandler: NewMsgHandle(),
		ConnMgr:    NewConnManage(),
	}
	return s
}
