package mapreduce

import (
	"common"
	"constant"
)

func (wk *Worker) Register(mr string) (err error) {
	args := common.RegisterArgs{wk.address}
	err = common.Call(common.SrvAddr(mr, constant.MASTER_RPC), constant.REGISETR, args, nil)

	if err != nil {
		return err
	}

	wk.master = mr

	return nil
}
