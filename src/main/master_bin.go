package main

import (
	"common"
	"mapreduce"
)

func main() {
	mr := mapreduce.NewMaster()

	err := mr.StartRPCServer()
	if err != nil {
		common.Debug("Master: Started Failed %s", err)
		return
	}

	<-mr.ExitChannel
}
