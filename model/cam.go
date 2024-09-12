package model

const (
	CAM_CLASS_BOLT      = 0 //球机
	CAM_CLASS_VD        = 1 //固定枪机
	CAM_CLASS_FISHEYE   = 2 //鱼眼
	CAM_CLASS_BINOCULAR = 3 //双目
)

type Cam struct {
	Id          int64
	ParentId    int64
	Name        string
	Addr        string //安装位置
	StreamAddr  string //流地址
	StreamProto string //流协议
	Resolution  string //分辨率
}
