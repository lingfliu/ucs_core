package gis

type GSegment struct {
	Start GPos
	End   GPos
}

// 计算线段是否相交
func (seg *GSegment) Intersect(seg2 *GSegment) bool {
	if seg.Start.Longitude == seg.End.Longitude && seg.Start.Latitude == seg.End.Latitude {
		return false
	}
	return false
}
