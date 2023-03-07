package main

import (
	"fmt"
	"github.com/lemoba/mmo-game/apis"
	"github.com/lemoba/mmo-game/core"
	"github.com/lemoba/zinx/ziface"
	"github.com/lemoba/zinx/znet"
)

// 当前客户端建立连接之后的hook函数
func OnConnectionAdd(conn ziface.IConnection) {
	// 创建一个Player对象
	player := core.NewPlayer(conn)
	// 给客户端发送MsgID: 1的消息
	player.SyncPid()
	// 给客户端发送MsgID: 200的消息
	player.BroadCastStartPosition()
	// 将当前新上线的玩家添加到WorldManager中
	core.WorldMgrObj.AddPlayer(player)

	// 将改连接绑定一个玩家ID的属性
	conn.SetProperty("pid", player.Pid)

	fmt.Println("====> Player pid = ", player.Pid, " is arrived <====")
}

func main() {
	// 创建zinx
	s := znet.NewServer("MMO Game")

	// 创建HOOK函数
	s.SetOnConnStart(OnConnectionAdd)
	// 注册路由业务
	s.AddRouter(2, &apis.WorldChatApi{})
	// 启动服务
	s.Serve()
}
