package impl

// /**
//  * Three-boom Drilling Jumbo 三臂凿岩台车
//  */
// //TODO: improve template generation 模板应当在前端页面编辑
// func CreateTbjVehiTemplate() *model.Mach {
// 	mach := &model.Mach{
// 		PropSet:    make(map[string]string),
// 		DNodeSet:   make(map[int64]*model.DNode),
// 		CtlNodeSet: make(map[int64]*model.CtlNode),
// 		CamSet:     make(map[int64]*model.Cam),
// 	}

// 	//所有状态点位均设置至一个监测点位内
// 	dn := &model.DNode{
// 		Class:     "三臂凿岩台车",
// 		Name:      "257710",
// 		PropSet:   make(map[string]string),
// 		DPointSet: make(map[int64]*model.DPoint),
// 	}
// 	metaAngle1 := &meta.DataMeta{
// 		ByteLen:   2,
// 		Dimen:     1,
// 		SampleLen: 1,
// 		DataClass: meta.DATA_CLASS_INT16,
// 	}
// 	pid := int64(0)
// 	p1 := &model.DPoint{
// 		Id:       pid,
// 		NodeId:   dn.Id,
// 		Class:    "臂1大臂旋转角",
// 		NodeName: dn.Name,
// 		DataMeta: metaAngle1}
// 	pid++
// 	p2 := &model.DPoint{
// 		Id:       pid,
// 		NodeId:   dn.Id,
// 		Class:    "臂1小臂旋转角",
// 		NodeName: dn.Name,
// 		DataMeta: metaAngle1,
// 	}
// 	pid++

// 	p3 := &model.DPoint{
// 		Id:       pid,
// 		NodeId:   dn.Id,
// 		Class:    "臂2大臂旋转角",
// 		NodeName: dn.Name,
// 		DataMeta: metaAngle1,
// 	}
// 	pid++

// 	p4 := &model.DPoint{
// 		Id:       pid,
// 		NodeId:   dn.Id,
// 		Class:    "臂2小臂旋转角",
// 		NodeName: dn.Name,
// 		DataMeta: metaAngle1,
// 	}

// 	dn.DPointSet[0] = p1
// 	dn.DPointSet[1] = p2
// 	dn.DPointSet[0] = p3
// 	dn.DPointSet[1] = p4

// 	mach.DNodeSet[0] = dn

// 	return mach
// }

// func CreateTbjVehi(name string, addr string, camAddr []string, template *model.Mach) *model.Mach {
// 	//基于模板生成
// 	mach := template
// 	//填充对应字段
// 	mach.Addr = addr
// 	mach.Name = name

// 	return mach
// }
