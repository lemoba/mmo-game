package core

import (
	"fmt"
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
