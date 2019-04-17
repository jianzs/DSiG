package master

import "common"

func (mr *Master) Register(args *common.RegisterArgs, _ *struct{}) error {
	mr.Lock()
	defer mr.Unlock()

	for _, w := range mr.workers {
		if w.address == args.Worker {
			common.Debug("Master: %s already in workers", w)
			return nil
		}
	}

	mr.workers = append(mr.workers, *newWorker(args.Worker))
	common.Debug("Master: A new worker '%s' join workers successfully", args.Worker)
	mr.newCond.Broadcast() // 通知正在调度的任务

	return nil
}
