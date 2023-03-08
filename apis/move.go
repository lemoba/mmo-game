package apis

import (
	"fmt"
	"github.com/lemoba/mmo-game/core"
	"github.com/lemoba/mmo-game/proto/pb"
	"github.com/lemoba/zinx/ziface"
	"github.com/lemoba/zinx/znet"
	"google.golang.org/protobuf/proto"
)

// 玩家移动
type MoveApi struct {
	znet.BaseRouter
}

func (m *MoveApi) Handle(request ziface.IRequest) {
	proto_msg := &pb.Position{}
	err := proto.Unmarshal(request.GetData(), proto_msg)
	if err != nil {
		fmt.Println("Move: Position Unmarshal error: ", err)
		return
	}

	// 得到当前发送位置是哪个玩家
	pid, err := request.GetConnection().GetProperty("pid")
	if err != nil {
		fmt.Println("GetProperty Pid error: ", err)
		return
	}
	fmt.Printf("Player pid = %d, move(%f, %f, %f, %f)", pid, proto_msg.X, proto_msg.Y, proto_msg.Z, proto_msg.V)

	// 给其他玩家进行位置信息广播
	player := core.WorldMgrObj.GetPlayerByPid(pid.(uint32))

	player.UpdatePosition(proto_msg.X, proto_msg.Y, proto_msg.Z, proto_msg.V)
}
