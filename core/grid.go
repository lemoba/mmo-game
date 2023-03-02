package core

import (
	"fmt"
	"sync"
)

/*
一个AOI地图中的格子类型
*/

type Grid struct {
	GID       uint32          // 格子类型
	MinX      uint32          // 格子左边界坐标
	MaxX      uint32          // 格子右边界坐标
	MinY      uint32          // 格子上边界坐标
	MaxY      uint32          // 格子下边界坐标
	playerIDs map[uint32]bool // 当前格子内玩家或者物体成员的ID集合
	pIDLock   sync.RWMutex    // 保护集合的读写锁
}

// 初始化当前的格子的方法
func NewGrid(gID, minX, maxX, minY, maxY uint32) *Grid {
	return &Grid{
		GID:       gID,
		MinX:      minX,
		MaxX:      maxX,
		MinY:      minY,
		MaxY:      maxY,
		playerIDs: make(map[uint32]bool),
	}
}

// 给格子添加一个玩家
func (g *Grid) Add(playerID uint32) {
	g.pIDLock.Lock()
	defer g.pIDLock.Unlock()

	g.playerIDs[playerID] = true
}

func (g *Grid) Remove(playerID uint32) {
	g.pIDLock.Lock()
	defer g.pIDLock.Unlock()

	delete(g.playerIDs, playerID)
}

// 得到当前格子中所有的玩家ID
func (g *Grid) GetPlayerIDs() (playerIDs []uint32) {
	g.pIDLock.RLock()
	defer g.pIDLock.RUnlock()

	for k, _ := range g.playerIDs {
		playerIDs = append(playerIDs, k)
	}
	return
}

// 调式使用-打印出格子的基本信息
func (g *Grid) String() string {
	return fmt.Sprintf("Grid id: %d, minX: %d, maxX: %d, minY: %d, maxY: %d, playerIDs: %v",
		g.GID, g.MinX, g.MaxX, g.MinY, g.MaxY, g.playerIDs)
}
