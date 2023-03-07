package apis

import (
	"fmt"
	"github.com/lemoba/mmo-game/core"
	"github.com/lemoba/mmo-game/proto/pb"
	"github.com/lemoba/zinx/ziface"
	"github.com/lemoba/zinx/znet"
	"google.golang.org/protobuf/proto"
)

// 世界聊天 路由业务
type WorldChatApi struct {
	znet.BaseRouter
}

func (wc *WorldChatApi) Handle(request ziface.IRequest) {
	// 1. 解析客户端场传递的proto协议
	proto_msg := &pb.Talk{}

	err := proto.Unmarshal(request.GetData(), proto_msg)
	if err != nil {
		fmt.Println("Talk Unmarshal error: ", err)
		return
	}
	// 2. 当前的聊天数据是哪个玩家发送的
	pid, err := request.GetConnection().GetProperty("pid")

	if err != nil {
		fmt.Println("get pid error: ", err)
	}

	// 3. 得到player玩家信息
	player := core.WorldMgrObj.GetPlayerByPid(pid.(uint32))

	// 4. 发送消息给全部玩家
	player.Talk(proto_msg.Content)
}
