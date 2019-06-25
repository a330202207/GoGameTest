package main

import (
	"fmt"
	"gFrame/giface"
	"gFrame/gnet"
)

//ping test 自定义路由
type PingRouter struct {
	gnet.BaseRouter
}

//Test Handle
func (this *PingRouter) Handle(request giface.IRequest) {
	fmt.Println("Call Router Handle")
	//先读取客户端的数据，再会写ping...ping...ping
	fmt.Println("Recv From Client:msgID=", request.GetMsgID(),
		",Data=", string(request.GetData()))

	err := request.GetConnection().SendMsg(200, []byte("ping...ping...ping"))
	if err != nil {
		fmt.Println(err)
	}
}

type HelloGinxRouter struct {
	gnet.BaseRouter
}

//创建链接之后执行的钩子函数
func DoConnectionBegin(conn giface.IConnection) {
	fmt.Println("======>DoConnectionBegin is Called...")
	if err := conn.SendMsg(202, []byte("DoConnection BEGIN")); err != nil {
		fmt.Println(err)
	}

	fmt.Println("Set Conn Name, Hoe......")
	//给当前链接设置属性
	conn.SetProperty("Name", "Ned")
	conn.SetProperty("Age", "28")
}

//链接断开之前的需要执行的函数
func DoConnectionLost(conn giface.IConnection) {
	fmt.Println("======>DoConnectionLost is Called...")
	fmt.Println("Conn ID = ", conn.GetConnID(), " Is Lost...")

	//获取链接属性
	if name, err := conn.GetProperty("Name"); err == nil {
		fmt.Println("Name = ", name)
	}

	if age, err := conn.GetProperty("Age"); err == nil {
		fmt.Println("Age = ", age)
	}
}

func (this *HelloGinxRouter) Handle(request giface.IRequest) {
	fmt.Println("Call HelloGinxRouter Handle")
	//先读取客户端的数据，再会写ping...ping...ping
	fmt.Println("Recv From Client:msgID=", request.GetMsgID(),
		",Data=", string(request.GetData()))

	err := request.GetConnection().SendMsg(201, []byte("Hello"))
	if err != nil {
		fmt.Println(err)
	}
}

func main() {
	//创建server句柄,使用
	s := gnet.NewServer("gInx V1.0")

	//注册链接Hook钩子函数
	s.SetOnConnStart(DoConnectionBegin)

	s.SetOnConnStop(DoConnectionLost)

	//当前框架添加自定义router
	s.AddRouter(0, &PingRouter{})

	s.AddRouter(1, &HelloGinxRouter{})

	//启动server
	s.Server()
}
