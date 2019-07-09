package core

import (
	"fmt"
	"testing"
)

func TestNewAOIManager(t *testing.T) {
	//初始化AOIManager
	aoiMgr := NewAOIManager(0, 250, 5, 0, 250, 5)

	//打印AOIManner
	fmt.Println(aoiMgr)
}

func TestAOIManager_GetSurroundGridsByGid(t *testing.T) {
	//初始化AOIManager
	aoiMgr := NewAOIManager(0, 250, 5, 0, 250, 5)

	for gid, _ := range grids {

		//得到当前gid的周五九宫格信息
		grids := GetSurroundGridsByGid(gid)
		fmt.Println("gid : ", gid, "，grids len = ", len(grids))

		gIDs := make([]int, 0, len(grids))
		for _, grid := range grids {
			gIDs = append(gIDs, GID)
		}

		fmt.Println("Surrounding Grid IDs Are", gIDs)
	}

}

func TestAOIManager_GetPIdsByPos(t *testing.T) {
	//初始化AOIManager
	aoiMgr := NewAOIManager(0, 250, 5, 0, 250, 5)
}
