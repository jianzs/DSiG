package main

import (
	"common"
	"constant"
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

	kp := file.NewKeeper(constant.MASTER_FILE_RPC)
	err = kp.StartRPCServer()
	if err != nil {
		common.Debug("File Keeper: Started Failed %s", err)
		return
	}

	<-mr.ExitChannel
}
