package msg

import (
	"encoding/binary"

	"github.com/lingfliu/ucs_core/model/meta"
)

/**
 * 寻址模式下的数据消息结构
 */
const (
	D_MSG_CLASS_AUTO = 0
	D_MSG_CLASS_TRIG = 1
	D_MSG_CLASS_POLL = 2
)

type DMsgData struct {
	Offset int
	Meta   *meta.DataMeta
	Data   []byte
}

/**
 * 监测消息，每一个消息携带一个DNode所有点位的数据
 */
type DMsg struct {
	DNodeId    int64
	DNodeAddr  string
	DNodeClass string
	Ts         int64       //timestamp of first sample
	Idx        int         //序号， 用于辅助判断是否丢包
	Session    string      //会话标识
	Mode       int         //模式，对应DNode的Mode: 0-定时采样，1-事件触发，2-轮询
	Sps        int64       //采样频率, 仅在Mode=0时有效
	SampleLen  int         //采样长度
	DataList   []*DMsgData //offset as the key
}

func (msg *DMsg) DataAsByte(offset int) ([][]byte, []int64) {
	var data *DMsgData
	for i := 0; i < len(msg.DataList); i++ {
		if msg.DataList[i].Offset == offset {
			data = msg.DataList[i]
			break
		}
	}
	if data == nil {
		return nil, nil
	}

	dimen := data.Meta.Dimen
	sampleLen := msg.SampleLen
	valueList := make([][]byte, sampleLen)
	tsList := make([]int64, sampleLen)
	for i := 0; i <= sampleLen; i++ {
		valueList[i] = data.Data[i*dimen : (i+1)*dimen]
		tsList[i] = msg.Ts + int64(i)*msg.Sps
	}
	return valueList, tsList
}

func (msg *DMsg) DataAsInt16(offset int) ([][]int16, []int64) {
	var data *DMsgData
	for i := 0; i < len(msg.DataList); i++ {
		if msg.DataList[i].Offset == offset {
			data = msg.DataList[i]
			break
		}
	}
	if data == nil {
		return nil, nil
	}

	dimen := data.Meta.Dimen
	sampleLen := msg.SampleLen
	valueList := make([][]int16, sampleLen)
	tsList := make([]int64, sampleLen)

	for i := 0; i <= sampleLen; i++ {

		sample := make([]int16, dimen)
		slice := data.Data[i*dimen : (i+1)*dimen]

		for j := 0; j < dimen; j++ {
			if data.Meta.Msb {
				sample[j] = int16(binary.BigEndian.Uint16(slice[j*2 : j*2+2]))
			} else {
				sample[j] = int16(binary.LittleEndian.Uint16(slice[j*2 : j*2+2]))
			}
		}
		valueList[i] = sample

		tsList[i] = msg.Ts + int64(i)*msg.Sps
	}
	return valueList, tsList

}

func (msg *DMsg) DatAsInt32(offset int) ([][]int32, []int64) {
	var data *DMsgData
	for i := 0; i < len(msg.DataList); i++ {
		if msg.DataList[i].Offset == offset {
			data = msg.DataList[i]
			break
		}
	}
	if data == nil {
		return nil, nil
	}

	dimen := data.Meta.Dimen
	sampleLen := msg.SampleLen
	valueList := make([][]int32, sampleLen)
	tsList := make([]int64, sampleLen)
	for i := 0; i < sampleLen; i++ {

		sample := make([]int32, dimen)
		slice := data.Data[i*dimen*4 : (i+1)*dimen*4]

		for j := 0; j < dimen; j++ {
			if data.Meta.Msb {
				sample[j] = int32(binary.BigEndian.Uint32(slice[j*4 : j*4+4]))
			} else {
				sample[j] = int32(binary.LittleEndian.Uint32(slice[j*4 : j*4+4]))
			}
		}
		valueList[i] = sample

		tsList[i] = msg.Ts + int64(i)*msg.Sps
	}
	return valueList, tsList
}

func (msg *DMsg) DataAsInt64(offset int) ([][]int64, []int64) {
	var data *DMsgData
	for i := 0; i < len(msg.DataList); i++ {
		if msg.DataList[i].Offset == offset {
			data = msg.DataList[i]
			break
		}
	}
	if data == nil {
		return nil, nil
	}

	sampleLen := msg.SampleLen
	valueList := make([][]int64, sampleLen)
	tsList := make([]int64, sampleLen)
	dimen := data.Meta.Dimen
	for i := 0; i <= sampleLen; i++ {

		sample := make([]int64, dimen)
		slice := data.Data[i*dimen : (i+1)*dimen]

		for j := 0; j < dimen; j++ {
			if data.Meta.Msb {
				sample[j] = int64(binary.BigEndian.Uint64(slice[j*8 : j*8+8]))
			} else {
				sample[j] = int64(binary.LittleEndian.Uint64(slice[j*8 : j*8+8]))
			}
		}
		valueList[i] = sample

		tsList[i] = msg.Ts + int64(i)*msg.Sps
	}
	return valueList, tsList
}
