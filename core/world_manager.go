package core

import (
	"fmt"
	"sync"
)

/*
当前游戏的世界总管理模块
*/
type WorldManager struct {
	// AOIManager当前世界地图AOI的管理模块
	AOIManager *AOIManager
	// 当前全部在线的Players集合
	Players map[uint32]*Player
	// 保护Players集合的锁
	pLock sync.RWMutex
}

// 提供全局的管理句柄
var WorldMgrObj *WorldManager

// 初始化方法
func init() {
	WorldMgrObj = &WorldManager{
		// 创建AOI
		AOIManager: NewAOIManager(AOI_MIN_X, AOI_MAX_X, AOI_CNTS_X, AOI_MIN_Y, AOI_MAX_Y, AOI_CNTS_Y),
		Players:    make(map[uint32]*Player),
	}
}

// 添加玩家
func (wm *WorldManager) AddPlayer(player *Player) {
	wm.pLock.Lock()
	wm.Players[player.Pid] = player
	wm.pLock.Unlock()

	// 将当前玩家添加到AOIManager中
	wm.AOIManager.AddToGridByPos(player.Pid, player.X, player.Z)

}

// 删除玩家
func (wm *WorldManager) RemovePlayerByPid(pid uint32) {
	player := wm.Players[pid]

	wm.pLock.Lock()
	delete(wm.Players, player.Pid)
	wm.pLock.Unlock()
	// 将当前玩家从AOIManager中移除
	wm.AOIManager.RemoveFromGridByPos(player.Pid, player.X, player.Z)
}

// 通过玩家ID查询Player对象
func (wm *WorldManager) GetPlayerByPid(pid uint32) *Player {
	wm.pLock.RLock()
	defer wm.pLock.RUnlock()

	player, ok := wm.Players[pid]
	if !ok {
		fmt.Println("Not exists player by pid = ", pid)
	}
	return player
}

// 获取全部的在线玩家
func (wm *WorldManager) GetAllPalyers() []*Player {
	wm.pLock.RLock()
	defer wm.pLock.RUnlock()

	players := make([]*Player, 0)

	for _, player := range wm.Players {
		players = append(players, player)
	}
	return players
}
