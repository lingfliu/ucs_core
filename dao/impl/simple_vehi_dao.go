package impl

import (
	"github.com/lingfliu/ucs_core/model"
	"github.com/lingfliu/ucs_core/model/meta"
)

const ()

func CreateSimpleVehiTemplate() *model.Mach {
	return nil
}

func CreateSimpleVehiFromTeamplate(mach *model.Mach) *model.Mach {
	return nil
}

// 车辆：1个转速监测，1个油门监测，1个刹车，3个摄像头，2个控制点（方向盘，油门，刹车）
func CreateSimpleVehi(name string, addr string, camAddr []string) *model.Mach {

	mach := &model.Mach{
		Id:         0,
		Name:       name,
		Addr:       addr,
		Class:      model.MACH_CLASS_SIMPLE_VEHI,
		PropSet:    map[string]string{"model": "simple_vehicle"},
		DNodeSet:   make(map[int64]*model.DNode),
		CtlNodeSet: make(map[int64]*model.CtlNode),
		CamSet:     make(map[int64]*model.Cam),
	}

	//添加DNode
	dnode := &model.DNode{
		Id:        0,
		Name:      "motor",
		Addr:      addr,
		Class:     "motor",
		Mode:      model.DNODE_MODE_AUTO,
		Sps:       1000,
		PropSet:   make(map[string]string),
		DPointSet: make(map[int64]*model.DPoint),
	}

	dnode.PropSet["type"] = "vehi"
	//TODO:  insert dnode

	dpMotor := &model.DPoint{
		Id:   0,
		Name: "motor",
		DataMeta: &meta.DataMeta{
			ByteLen:   4,
			Dimen:     1,
			DataClass: meta.DATA_CLASS_INT,
			Unit:      "rpm",
			Msb:       true,
		},
	}

	//TODO: insert dpMotor

	dpAcc := &model.DPoint{
		Id:   0,
		Name: "油门",
	}

	dpWheel := &model.DPoint{
		Id:   0,
		Name: "方向盘",
	}

	dnode.DPointSet[dpMotor.Id] = dpMotor
	dnode.DPointSet[dpAcc.Id] = dpAcc
	dnode.DPointSet[dpWheel.Id] = dpWheel //方向盘

	return mach
}
