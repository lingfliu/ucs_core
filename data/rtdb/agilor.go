package rtdb

const (
	DATA_TYPE_INT   = 1
	DATA_TYPE_FLOAT = 2

	DB_STATE_DISCONNECTED = 0
	DB_STATE_CONNECTED    = 1
	DB_STATE_CONNECTING   = 2
)

type AgilorRtdbDevice struct {
	Id         string
	Name       string
	DataPoints []*AgilorRtdbDataPoint
}

type AgilorRtdbDataPoint struct {
	Id       string
	Name     string
	DeviceId string
	Type     int
	Data     any
	Ts       int64
}

type AgilorRtdb struct {
	Host     string
	Port     int
	Username string
	Passwd   string
	State    int
}

func (db *AgilorRtdb) Connect() {

}

func (db *AgilorRtdb) Disconnect() {

}

func (db *AgilorRtdb) CreateDevice(d *AgilorRtdbDevice) {

}

func (db *AgilorRtdb) CreateDataPoint(p *AgilorRtdbDataPoint) {

}

/**
 * update (store) data point for point p
 */
func (db *AgilorRtdb) Update(p *AgilorRtdbDataPoint) {
}

func (db *AgilorRtdb) QueryCurrentByDevice(dId string) []*AgilorRtdbDataPoint {
	return make([]*AgilorRtdbDataPoint, 0)
}

func (db *AgilorRtdb) QueryCurrent(pId string) *AgilorRtdbDataPoint {
	return &AgilorRtdbDataPoint{}
}

func (db *AgilorRtdb) QueryByTimeAndDevice(dId string, pId string, tic int64, toc int64) [][]*AgilorRtdbDataPoint {
	return make([][]*AgilorRtdbDataPoint, 0)
}

func (db *AgilorRtdb) QueryByTime(dId string, pId string, tic int64, toc int64) []*AgilorRtdbDataPoint {
	return make([]*AgilorRtdbDataPoint, 0)
}

// func (db *AgilorRtdb) QueryAll(p *AgilorRtdbDataPoint) []*AgilorRtdbDataPoint {
// 	return make([]*AgilorRtdbDataPoint, 0)
// }

func (db *AgilorRtdb) DeleteByTime(pId string, tic int64, toc int64) {
}

/**
 * delete all data of data point pId
 */
func (db *AgilorRtdb) DropDataPoint(pId string) {
}

func (db *AgilorRtdb) DeleteDataPoint(pId string) {
}

func (db *AgilorRtdb) DropDevice(dId string) {
}

/**
 * delete all data of data point pId
 */
func (db *AgilorRtdb) DeleteDevice(dId string) {
	for _, p := range db.QueryCurrentByDevice(dId) {
		db.DeleteDataPoint(p.Id)
		db.DropDataPoint(p.Id)
	}
	db.DropDevice(dId)
}

/**
 * subscribe to a device
 */
func (db *AgilorRtdb) SubscribeDevice(dId string) int {
	return 0
}

/**
 * subscribe to a set of data points
 */
func (db *AgilorRtdb) SubscribeDataPoints(pIds []string) int {
	return 0
}

func (db *AgilorRtdb) UnsubscribeDevice(did string) int {
	return 0
}

func (db *AgilorRtdb) UnsubscribeDataPoints(pIds []string) int {
	return 0
}

//////////////////////////////////
//basic data aggregation
/////////////////////////////////

/*
* time alignment of a set of data points by the NN principle
@param pids: the list of data point ids
@param ts_baseline: the baseline timestamp in us
@param ts_step: the step of the alignment in us
@param tic: the start timestamp in us
@param toc: the end timestamp in us
*/
func (db *AgilorRtdb) Align(pids []string, ts_baseline int64, ts_step int64, tic int64, toc int64) []*AgilorRtdbDataPoint {
	return make([]*AgilorRtdbDataPoint, 0)
}

func (db *AgilorRtdb) Mean(points []*AgilorRtdbDataPoint) *AgilorRtdbDataPoint {
	return &AgilorRtdbDataPoint{}
}

func (db *AgilorRtdb) Std(points []*AgilorRtdbDataPoint) *AgilorRtdbDataPoint {
	return &AgilorRtdbDataPoint{}
}

func (db *AgilorRtdb) Max(points []*AgilorRtdbDataPoint) *AgilorRtdbDataPoint {
	return &AgilorRtdbDataPoint{}
}

func (db *AgilorRtdb) Min(points []*AgilorRtdbDataPoint) *AgilorRtdbDataPoint {
	return &AgilorRtdbDataPoint{}
}
