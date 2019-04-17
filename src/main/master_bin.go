package main

import (
	"common"
	"mapreduce/master"
)

func main() {
	mr := master.NewMaster()

	err := mr.StartRPCServer()
	if err != nil {
		common.Debug("Master: Started Failed %s", err)
		return
	}

	<-mr.ExitChannel
}
