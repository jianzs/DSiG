package main

import (
	"common"
	"constant"
	"file"
	"os"
)

func main() {
	args := os.Args

	var kp *file.Keeper
	switch args[1] {
	case "master":
		kp = file.NewKeeper(constant.MASTER_FILE_RPC)
	case "worker":
		kp = file.NewKeeper(constant.WORKER_FILE_RPC)
	case "client":
		kp = file.NewKeeper(constant.CLIENT_FILE_RPC)
	}

	err := kp.StartRPCServer()
	if err != nil {
		common.Debug("File Keeper: Started Failed %s", err)
		return
	}
	common.Debug("File Keeper: Started Successfully")

	<-kp.ExitChannel
}
