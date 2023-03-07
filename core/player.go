package core

import (
	"fmt"
	"github.com/lemoba/mmo-game/proto/pb"
	"github.com/lemoba/zinx/ziface"
	"google.golang.org/protobuf/proto"
	"math/rand"
	"sync"
)

// 玩家对象
type Player struct {
	Pid  uint32             // 玩家ID
	Conn ziface.IConnection // 当前玩家的连接(与客户端的连接)
	X    float32            // 平面x坐标
	Y    float32            // 高度
	Z    float32            // 平面y坐标
	V    float32            // 旋转角度(0-360)
}

/*
PlayerID gen
*/
var PidGen uint32 = 1
var IdLock sync.Mutex

// 创建一个玩家的方法
func NewPlayer(conn ziface.IConnection) *Player {
	// 生成玩家ID
	IdLock.Lock()
	id := PidGen
	PidGen++
	defer IdLock.Unlock()

	// 创建一个玩家对象
	return &Player{
		Pid:  id,
		Conn: conn,
		X:    float32(160 + rand.Intn(10)),
		Y:    0,
		Z:    float32(140 + rand.Intn(20)),
		V:    0,
	}
}

/*
提供一个发送给客户端消息的方法
主要将pb的protobuf数据序列化后，再调用zinx的SendMsg方法
*/
func (p *Player) SendMsg(msgID uint32, data proto.Message) {
	// 将proto Message结构体序列化 转换成二进制
	msg, err := proto.Marshal(data)
	if err != nil {
		fmt.Println("marshal msg error: ", err)
	}
	// 将二进制文件 通过zinx框架的sendmsg将数据发送给客户端
	if p.Conn == nil {
		fmt.Println("connection in player is nil")
		return
	}
	if err := p.Conn.SendMsg(msgID, msg); err != nil {
		fmt.Println("Player send msg error: ", err)
	}
}

// 告知客户端玩家Pid, 同步已经生成的玩家ID给客户端
func (p *Player) SyncPid() {
	// 组建MsgID: 0的proto数据
	proto_msg := &pb.SyncPid{
		Pid: p.Pid,
	}
	// 将消息发送给客户端
	p.SendMsg(1, proto_msg)
}

// 广播玩家子的出生地点
func (p *Player) BroadCastStartPosition() {
	// 组建MsgID: 200的proto数据
	proto_msg := &pb.BroadCast{
		Pid: p.Pid,
		Tp:  2, // 广播位置坐标
		Data: &pb.BroadCast_P{
			P: &pb.Position{
				X: p.X,
				Y: p.Y,
				Z: p.Z,
				V: p.V,
			},
		},
	}
	// 将消息发送给客户端
	p.SendMsg(200, proto_msg)
}

// 玩家广播新消息
func (p *Player) Talk(content string) {
	// 1. 组建MsgID: 200 proto数据
	proto_msg := &pb.BroadCast{
		Pid: p.Pid,
		Tp:  1, // 聊天广播
		Data: &pb.BroadCast_Content{
			content,
		},
	}
	// 2. 得到当前世界的所有玩家
	players := WorldMgrObj.GetAllPalyers()

	for _, player := range players {
		// 给对应的客户端发送消息
		player.SendMsg(200, proto_msg)
	}
}

// 同步玩家上线的位置信息
func (p *Player) SyncSurrounding() {
	// 1. 获取当前玩家周围有哪些
	pids := WorldMgrObj.AOIManager.GetPidsByPos(p.X, p.Z)

	players := make([]*Player, 0, len(pids))
	for _, pid := range pids {
		players = append(players, WorldMgrObj.GetPlayerByPid(pid))
	}
	// 2. 将当前的位置信息通过MsgID: 200 发送给周围的玩家(让其他玩家看到自己)
	// 2.1 MsgID: 200
	proto_msg := &pb.BroadCast{
		Pid: p.Pid,
		Tp:  2,
		Data: &pb.BroadCast_P{
			P: &pb.Position{
				X: p.X,
				Y: p.Y,
				Z: p.Z,
				V: p.V,
			},
		},
	}

	// 2.2 广播消息给周围玩家
	for _, player := range players {
		player.SendMsg(200, proto_msg)
	}

	// 3. 将周围的全部玩家的位置信息发送给当前的玩家MsgID: 202 客户端(让自己看到其他玩家)
	// 3.1 组建MsgID: 202 proto数据
	players_proto_msg := make([]*pb.Player, 0, len(players))

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
		players_proto_msg = append(players_proto_msg, p)
	}

	// 3.2 封装SyncPlayer protobuf数据
	syncPlayers_proto_msg := &pb.SyncPlayers{
		Ps: players_proto_msg[:],
	}

	// 3.3 将组建好的数据发送给当前的客户端
	p.SendMsg(202, syncPlayers_proto_msg)
}
