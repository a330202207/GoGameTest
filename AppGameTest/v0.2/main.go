package main

import (
	"fmt"
	"gFrame/giface"
	"gFrame/gnet"
	"game_test/apis"
	"game_test/core"
)

//当前客户端简历链接之后的Hook函数
func OnConnectionAdd(conn giface.IConnection) {
	//创建一个Player对象
	player := core.NewPlayer(conn)

	//给客户端发送MsgID：1消息: 同步当前Player的ID给客户端
	player.SyncPid()

	//给客户端发送MsgID：200的消息：同步当前Player的初始位置给客户端
	player.BroadCastStartPosition()

	//将当前新上线的玩家添加到WorldManager中
	core.WorldMgrObj.AddPlayer(player)

	//将该链接绑定一个Pid 玩家ID属性
	conn.SetProperty("pid", player.Pid)

	//同步周边玩家，告知他们当前玩家已经上线，广播当前玩家的信息
	player.SyncSurrounding()

	fmt.Println("===============> Player Pid = ", player.Pid, " Is Arrived <===============")
}

//给当前链接断开之前除非的Hook钩子函数
func OnConnectionLost(conn giface.IConnection) {

	//通过链接属性得到当前链接所绑定pid
	pid, _ := conn.GetProperty("pid")

	player := core.WorldMgrObj.GetPlayerByPid(pid.(int32))

	//玩家下线
	player.Offline()

	fmt.Println("===============> Player Pid = ", pid, " Offline <===============")
}

func main() {

	//创建Ginx Server句柄
	s := gnet.NewServer("MMO Game Ginx")

	//连接创建和销毁的Hook钩子函数
	s.SetOnConnStart(OnConnectionAdd)
	s.SetOnConnStop(OnConnectionLost)

	//注册一些路由业务
	s.AddRouter(2, &apis.WorldChatApi{})

	s.AddRouter(3, &apis.MoveAPi{})

	//启动服务
	s.Server()
}
