package main

import (
	"common"
	"constant"
	"file"
	"mapreduce/worker"
	"os"
)

// go run worker.go master-ip worker-ip
func main() {
	args := os.Args
	wkIp := args[1]
	mr := args[2]

	wk := worker.NewWorker(wkIp)

	err := wk.StartRPCServer()
	if err != nil {
		common.Debug("Worker: Started Failed %s", err)
		return
	}

	kp := file.NewKeeper(constant.WORKER_FILE_RPC)
	err = kp.StartRPCServer()
	if err != nil {
		common.Debug("File Keeper: Started Failed %s", err)
		return
	}

	err = wk.Register(mr)
	if err != nil {
		common.Debug("Worker: Register %s", err)
		return
	}

	<-wk.ExitChannel
}
