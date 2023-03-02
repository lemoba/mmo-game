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
	s := fmt.Sprintf("AOIManager: \nMinX:%d\nMaxX:%d\ncntsX:%d\nMinY:%\nMaxY:%d\ncntsY:%d\nGrids in AOIManager:\n")

	for _, grid := range m.grids {
		s += fmt.Sprintln(grid)
	}
	return s
}
