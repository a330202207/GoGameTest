package apis

import (
	"fmt"
	"gFrame/giface"
	"gFrame/gnet"
	"game_test/core"
	"game_test/pb"
	"github.com/golang/protobuf/proto"
)

type WorldChatApi struct {
	gnet.BaseRouter
}

func (wc *WorldChatApi) Handle(request giface.IRequest) {

	//解析客户端传递进来的proto协议
	protoMsg := &pb.Talk{}
	err := proto.Unmarshal(request.GetData(), protoMsg)
	if err != nil {
		fmt.Println("Talk Unmarshal Error", err)
		return
	}

	//将当前新上线的玩家添加到WorldManager中
	pid, err := request.GetConnection().GetProperty("pid")

	//根据Pid得到对应Player对象
	player := core.WorldMgrObj.GetPlayerByPid(pid.(int32))

	//将这个消息广播给其他全部在线玩家
	player.Talk(protoMsg.Content)
}
