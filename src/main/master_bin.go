package main

import (
	"common"
	"file"
	"mapreduce/master"
)

func main() {
	mr := master.NewMaster()

	err := mr.StartRPCServer()
	if err != nil {
		common.Debug("Master: Started Failed %s", err)
		return
	}

	kp := file.NewKeeper()
	err = kp.StartRPCServer()
	if err != nil {
		common.Debug("File Keeper: Started Failed %s", err)
		return
	}

	<-mr.ExitChannel
}
