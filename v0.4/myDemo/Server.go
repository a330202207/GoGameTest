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

//Test PreHandle
func (this *PingRouter) PreHandle(request giface.IRequest) {
	fmt.Println("Call Router PreHandle")
	_, err := request.GetConnection().GetTCPConnection().Write([]byte("before ping...\n"))
	if err != nil {
		fmt.Println("Call Back Before Ping Error!")
	}
}

//Test Handle
func (this *PingRouter) Handle(request giface.IRequest) {
	fmt.Println("Call Router PreHandle")
	_, err := request.GetConnection().GetTCPConnection().Write([]byte("ping...ping...ping...\n"))
	if err != nil {
		fmt.Println("Call Back Ping Error!")
	}
}

//Test PostHandle
func (this *PingRouter) PostHandle(request giface.IRequest) {
	fmt.Println("Call Router PreHandle")
	_, err := request.GetConnection().GetTCPConnection().Write([]byte("after ping...\n"))
	if err != nil {
		fmt.Println("Call Back After Ping Error!")
	}
}

func main() {
	//创建server句柄,使用
	s := gnet.NewServer("gInx V0.4")

	//当前框架添加一个自定义router
	s.AddRouter(&PingRouter{})

	//启动server
	s.Server()
}
