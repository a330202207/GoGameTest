package apis

import (
	"fmt"
	"gFrame/giface"
	"gFrame/gnet"
	"game_test/core"
	"game_test/pb"
	"github.com/golang/protobuf/proto"
)

type MoveAPi struct {
	gnet.BaseRouter
}

func (m *MoveAPi) Handle(request giface.IRequest) {
	//解析客户端传递进来的proto协议
	protoMsg := &pb.Position{}

	err := proto.Unmarshal(request.GetData(), protoMsg)
	if err != nil {
		fmt.Println("Move : Position Unmarshal Error", err)
		return
	}

	//得到当前发送位置的是那个玩家
	pid, err := request.GetConnection().GetProperty("pid")
	if err != nil {
		fmt.Println("GetProperty Pid Error", err)
		return
	}
	//fmt.Printf("Player Pid = %d, Move(%f,%f,%f,%f)\n", pid, protoMsg.X, protoMsg.Y, protoMsg.Z, protoMsg.V)

	//给其他玩家进行当前玩家的位置信息广播
	player := core.WorldMgrObj.GetPlayerByPid(pid.(int32))

	//广播并更新当前玩家的坐标
	player.UnpdatePos(protoMsg.X, protoMsg.Y, protoMsg.Z, protoMsg.V)
}
