package core

import "fmt"

/*
AOI区域管理模块
*/

type AOIManager struct {
	// 区域左边界坐标
	MinX uint32
	// 区域右边界坐标
	MaxX uint32
	// X方向格子的数量
	CntsX uint32
	// 区域上边界坐标
	MinY uint32
	// 区域下边界坐标
	MaxY uint32
	// Y方向格子的数量
	CntsY uint32
	// 当前区域中有哪些格子map-key=格子的ID,value=格子对象
	grids map[uint32]*Grid
}

/*
初始化一个AOI区域管理模块
*/
func NewAOIManager(minX, maxX, cntsX, minY, maxY, cntsY uint32) *AOIManager {
	aoiMgr := &AOIManager{
		MinX:  minX,
		MaxX:  maxX,
		CntsX: cntsY,
		MinY:  minY,
		MaxY:  maxY,
		CntsY: cntsY,
		grids: make(map[uint32]*Grid),
	}

	// 给AOI初始化区域的格子所有格子进行编号和初始化
	var x, y uint32
	for y = 0; y < cntsY; y++ {
		for x = 0; x < cntsX; x++ {
			// 计算格子ID，根据x, y编号
			// 各自编号： id = idy * cntsX + idx
			gid := y*cntsX + x

			// 初始化gid格子
			aoiMgr.grids[gid] = NewGrid(gid,
				aoiMgr.MinX+x*aoiMgr.gridWidth(),
				aoiMgr.MinX+(x+1)*aoiMgr.gridWidth(),
				aoiMgr.MinY+y*aoiMgr.gridHeight(),
				aoiMgr.MinY+(y+1)*aoiMgr.gridHeight())
		}
	}
	return aoiMgr
}

// 得到每个格子在X轴方向的宽度
func (m *AOIManager) gridWidth() uint32 {
	return (m.MaxX - m.MinX) / m.CntsX
}

// 得到每个格子在Y轴方向的宽度
func (m *AOIManager) gridHeight() uint32 {
	return (m.MaxY - m.MinY) / m.CntsY
}

// 打印格子信息
func (m *AOIManager) String() string {
	s := fmt.Sprintf("AOIManager: \nMinX:%d\nMaxX:%d\ncntsX:%d\nMinY:%d\nMaxY:%d\ncntsY:%d\nGrids in AOIManager:\n",
		m.MinX, m.MinY, m.CntsX, m.MinY, m.MaxY, m.CntsY)

	for _, grid := range m.grids {
		s += fmt.Sprintln(grid)
	}
	return s
}

// 根据GID得到周边九宫格的ID集合
func (m *AOIManager) GetSurroundGridsByGid(gID uint32) (grids []*Grid) {
	// 判断gID是否在AOIManager中
	if _, ok := m.grids[gID]; !ok {
		return
	}

	// 将当前gid加入九宫格切片中
	grids = append(grids, m.grids[gID])
	// 需要判断gID左右是否有格子
	// 需要通过gID得到当前格子x的编号 idx = id % nx
	idx := gID % m.CntsX

	// 判断idx编号左边是否还有格子，如果有放在gridsX集合中
	if idx > 0 {
		grids = append(grids, m.grids[gID-1])
	}

	// 判断idx编号右边是否还有格子，如果有放在gridsX集合中
	if idx < m.CntsX-1 {
		grids = append(grids, m.grids[gID+1])
	}

	// 遍历gridsX集合中每个格子的gid
	gridsX := make([]uint32, 0, len(grids))

	for _, v := range grids {
		gridsX = append(gridsX, v.GID)
	}

	for _, v := range gridsX {
		idy := v / m.CntsY
		// 判断grid上边是否有格子，如果有放在gridsY集合中
		if idy > 0 {
			grids = append(grids, m.grids[v-m.CntsX])
		}
		// 判断grid下边是否有格子，如果有放在gridsY集合中
		if idy < m.CntsY-1 {
			grids = append(grids, m.grids[v+m.CntsX])
		}
	}
	return
}

func (m *AOIManager) GetGidByPos(x, y float32) uint32 {
	idx := (uint32(x) - m.MinX) / m.gridWidth()
	idy := (uint32(y) - m.MinY) / m.gridHeight()
	return idy*m.CntsX + idx
}

// 通过横纵坐标得到周边九宫格内全部的PlayerIDs
func (m *AOIManager) GetPidsByPos(x, y float32) (playerIDs []uint32) {
	// 得到当前玩家的GID
	gID := m.GetGidByPos(x, y)
	// 通过GID得到九宫格信息
	grids := m.GetSurroundGridsByGid(gID)

	// 将九宫格信息里的playerID
	for _, grid := range grids {
		playerIDs = append(playerIDs, grid.GetPlayerIDs()...)
		fmt.Printf("===> girid ID: %d, pids: %v", grid.GID, grid.GetPlayerIDs())
	}
	return
}

// 添加一个playerID到格子中
func (m *AOIManager) AddPidToGrid(pID, gID uint32) {
	m.grids[gID].Add(pID)
}

// 移除一个格子中的playerID
func (m *AOIManager) RemovePidFromGrid(pID, gID uint32) {
	m.grids[gID].Remove(pID)
}

// 通过GID获取全部的playerID
func (m *AOIManager) GetPidsByGid(gID uint32) (playerIDs []uint32) {
	playerIDs = m.grids[gID].GetPlayerIDs()
	return
}

// 通过坐标把player添加到格子中
func (m *AOIManager) AddToGridByPos(pID uint32, x, y float32) {
	gId := m.GetGidByPos(x, y)
	grid := m.grids[gId]
	grid.Add(pID)
}

// 通过坐标把player从格子中删除
func (m *AOIManager) RemoveFromGridByPos(pID uint32, x, y float32) {
	gId := m.GetGidByPos(x, y)
	grid := m.grids[gId]
	grid.Remove(pID)
}
