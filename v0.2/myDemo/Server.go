package main

import "gFrame/gnet"

func main() {
	//创建server句柄,使用
	s := gnet.NewServer("gIndex V0.2")

	//启动server
	s.Server()
}