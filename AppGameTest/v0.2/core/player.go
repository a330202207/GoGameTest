package core

import (
	"fmt"
	"gFrame/giface"
	"game_test/pb"
	"github.com/golang/protobuf/proto"
	"math/rand"
	"sync"
)

//玩家对象
type Player struct {
	Pid  int32              //玩家ID
	Conn giface.IConnection //当前玩家的链接（用于和客户端链接）
	X    float32            //平面的x坐标
	Y    float32            //高度
	Z    float32            //平面Y坐标
	V    float32            //旋转的0-360角度
}

var PidGen int32 = 1  //用来生成玩家ID的计数器
var IdLock sync.Mutex //保护PIdGen的Mutex

//创建一个玩家
func NewPlayer(conn giface.IConnection) *Player {

	//生成一个玩家ID
	IdLock.Lock()
	id := PidGen
	PidGen++
	IdLock.Unlock()

	p := &Player{
		Pid:  id,
		Conn: conn,
		X:    float32(160 + rand.Intn(20)), //随机在160坐标掉，基于X轴若干偏移
		Y:    0,
		Z:    float32(140 + rand.Intn(20)), //随机在140坐标点,基于Y轴若干偏移
		V:    0,                            //角度为0
	}

	return p
}

//提供一个发送给客户端的消息方法
//主要是讲pb的protobuf数据序列化之后，再调用ginx的SendMsg方法
func (p *Player) SendMsg(msgId uint32, data proto.Message) {
	//将proto Message结构体序列化转换成二进制
	msg, err := proto.Marshal(data)
	if err != nil {
		fmt.Println("Marshal Msg Err", err)
		return
	}

	//将二进制文件，通过框架饿Sendmsg将数据方式给客户端
	if p.Conn == nil {
		fmt.Println("Connnection In Player Is Nil")
		return
	}

	if err := p.Conn.SendMsg(msgId, msg); err != nil {
		fmt.Println("Player SendMsg Error!")
		return
	}

	return
}

//将PlayerID同步给客户端
func (p *Player) SyncPid() {
	//组建MsgID：0 的proto数据
	protoMsg := &pb.SyncPid{
		Pid: p.Pid,
	}
	p.SendMsg(1, protoMsg)
}

//将PlayerID上线初始位置同步给客户端
func (p *Player) BroadCastStartPosition() {
	//组建MsgID：200 的proto数据
	protoMsg := &pb.BroadCast{
		Pid: p.Pid,
		Tp:  2, //广播玩家的位置
		Data: &pb.BroadCast_P{
			P: &pb.Position{
				X: p.X,
				Y: p.Y,
				Z: p.Z,
				V: p.V,
			},
		},
	}
	p.SendMsg(200, protoMsg)
}

//玩家广播世界聊天信息
func (p *Player) Talk(content string) {
	//组建MsgID：200 proto数据
	protoMsg := &pb.BroadCast{
		Pid: p.Pid,
		Tp:  1, //代表聊天广播
		Data: &pb.BroadCast_Content{
			Content: content,
		},
	}

	//得到当前世界所有的在线玩家
	players := WorldMgrObj.GetAllPlayers()

	//向所有的玩家（包括自己）发送MsgID：200消息
	for _, player := range players {
		//player分别给对应的客户端发送消息
		player.SendMsg(200, protoMsg)
	}

}

//同步玩家上线的位置信息
func (p *Player) SyncSurrounding() {
	//获取当前周围玩家有哪些（九宫格）
	pIds := WorldMgrObj.AoiMgr.GetPIdsByPos(p.X, p.Z)
	players := make([]*Player, 0, len(pIds))
	for _, pid := range pIds {
		players = append(players, WorldMgrObj.GetPlayerByPid(int32(pid)))
	}

	//将当前玩家位置通过MsgID:200发给周围的玩家（让其他玩家看到自己）
	protoMsg := &pb.BroadCast{
		Pid: p.Pid,
		Tp:  2, //代表广播坐标
		Data: &pb.BroadCast_P{
			P: &pb.Position{
				X: p.X,
				Y: p.Y,
				Z: p.Z,
				V: p.V,
			},
		},
	}

	//当前玩家的客户端
	for _, player := range players {
		player.SendMsg(200, protoMsg)
	}

	//将周围的全部玩家位置信息发送给当前玩家MsgID:202，让自己看的其他玩家
	playersProtoMsg := make([]*pb.Player, 0, len(players))
	for _, player := range players {
		p := &pb.Player{
			Pid: player.Pid,
			P: &pb.Position{
				X: player.X,
				Y: player.Y,
				Z: player.Z,
				V: player.V,
			},
		}

		playersProtoMsg = append(playersProtoMsg, p)
	}

	SyncPlayerProtoMsg := &pb.SyncPlayers{
		Ps: playersProtoMsg[:],
	}

	//将组建好的数据给发送给当前玩家的客户端
	p.SendMsg(202, SyncPlayerProtoMsg)
}

//广播当前玩家的位置移动信息
func (p *Player) UnpdatePos(x float32, y float32, v float32, z float32) {
	//更新当前玩家player对象的坐标
	p.X = x
	p.Y = y
	p.V = v
	p.Z = z

	//组建广播proto协议 MsgID：200
	protoMsg := &pb.BroadCast{
		Pid: p.Pid,
		Tp:  4, //移动之后的坐标信息
		Data: &pb.BroadCast_P{
			P: &pb.Position{
				X: p.X,
				Y: p.Y,
				V: p.V,
				Z: p.Z,
			},
		},
	}

	//获取当前玩家的周边玩家AOI九宫格之内的玩家
	players := p.GetSurroundingPlayers()

	//一次给每个玩家对应的和客户端发送当前玩家位置更新的消息
	for _, player := range players {
		player.SendMsg(200, protoMsg)
	}
}

//获取当前玩家的周边玩家AOI九宫格内的玩家
func (p *Player) GetSurroundingPlayers() []*Player {
	pids := WorldMgrObj.AoiMgr.GetPIdsByPos(p.X, p.Z)

	players := make([]*Player, 0, len(pids))
	for _, pid := range pids {
		players = append(players, WorldMgrObj.GetPlayerByPid(int32(pid)))
	}

	return players
}

//玩家下线
func (p *Player) Offline() {
	//得到当前玩家周边的九宫格都有哪些玩家
	players := p.GetSurroundingPlayers()

	//给周围玩家广播MgsID:201消息
	protoMsg := &pb.SyncPid{
		Pid: p.Pid,
	}

	for _, player := range players {
		player.SendMsg(201, protoMsg)
	}

	//将当前玩家从世界管理器删除
	WorldMgrObj.AoiMgr.RemoveFromGridByPos(int(p.Pid), p.X, p.Z)

	//将当前玩家从AOI管理器删除
	WorldMgrObj.RemovePlayer(p.Pid)
}
