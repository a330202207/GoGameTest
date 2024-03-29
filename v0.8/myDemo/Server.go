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
	s := gnet.NewServer("gInx V0.8")

	//当前框架添加自定义router
	s.AddRouter(0, &PingRouter{})

	s.AddRouter(1, &HelloGinxRouter{})

	//启动server
	s.Server()
}
