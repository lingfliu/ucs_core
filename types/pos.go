package types

/**
 * 位置信息
 */
type GPos struct {
	// GNSS position
	Longitude float64
	Latitude  float64
	Altitude  float64

	// local position
	X float64
	Y float64
	H float64
	//辅助地址编码信息，如楼层，区域等
	AddrMaj int
	AddrMin int
}

/**
 * Local空间坐标
 */
type LPos struct {
	//position (m)
	X float64
	Y float64
	Z float64

	//rotation (rad)
	Raw   float64
	Yaw   float64
	Pitch float64
}

/**
 * 速度信息
 */
type Velo struct {
	//position (m/s)
	X float64
	Y float64
	Z float64

	//角速度 (rad/s)
	Raw   float64
	Yaw   float64
	Pitch float64
}
