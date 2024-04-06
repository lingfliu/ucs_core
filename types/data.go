package types

type Data struct {
	Meta    *PropMeta
	Ts      int64 //timestamp
	Payload []any
}
