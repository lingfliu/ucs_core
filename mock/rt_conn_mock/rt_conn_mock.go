package main

const (
	dbName = "ucs-rtdata"
	dbHost = "localhost:6030"
)

// Mock code here
func main() {

	flowCli := &rtdb.FlowCli{}
	flowCli.RegFlow("ucs-rtdata", filter func(data *rtdb.Data) string)

	tdbCli := &rtdb.TdbCli{}
	tdbCli.Open(dbHost, dbName)
	tdbCli.Use(dbName)

	ddsCli := &rtdb.DdsCli{}

	ddsCli.Subscribe("dds://localhost:6030/ucs/eval_demo", func(msg *rtdb.DdsMessage) {
		//handling message here
		tableName = msg.TableName //raw data
		mData = msg.ToData()
		tdbCli.Store(tableName, mData)

		flowCli.Submit(flowCli)
	})

}
