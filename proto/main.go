package main

import "github.com/lemoba/zinx/znet"

func main() {
	// 创建zinx
	s := znet.NewServer("MMO Game")

	// 创建HOOK函数
	// 注册路由业务
	// 启动服务
	s.Start()
}
