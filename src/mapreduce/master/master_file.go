package master

import (
	"common"
	"constant"
)

func (mr *Master) ReadFile(args *common.FileArgs, reply *common.FileReply) error {
	common.Debug("Master: Receive RPC read file '%s'", args.Filename)
	content, err := common.ReadFile(args.Filename)
	if err != nil {
		common.Debug("Master Read file error %s", err)
		reply.Err = err
		return err
	}
	reply.Content = content
	common.Debug("Master: read file '%s' successfully", args.Filename)
	return nil
}

func (mr *Master) WriteFile(args *common.FileArgs, reply *common.FileReply) error {
	common.Debug("Master: Receive RPC write file '%s'", args.Filename)
	err := common.WriteFile(args.Filename, args.Content)
	if err != nil {
		reply.Err = err
		return err
	}
	common.Debug("Master: write file '%s' successfully", args.Filename)
	return nil
}

func (mr *Master) mergeFiles(jobName string) error {
	job := mr.jobs[jobName]

	for idx, wk := range job.reduceWorker {
		data, err := common.ReadFileRemote(common.SrvAddr(wk.address, constant.WORKER_FILE_RPC), common.ReduceName(jobName, idx))
		if err != nil {
			return err
		}
		err = common.AppendFileRemote(common.SrvAddr(job.client, constant.CLIENT_FILE_RPC), common.FinalName(jobName, job.job.OutFile), data)
		if err != nil {
			return err
		}
	}
	return nil
}

func (mr *Master) CleanupFiles(job common.Job) {
	mapWks := mr.jobs[job.Name].mapWorkers
	redWks := mr.jobs[job.Name].reduceWorker

	for i := 0; i < job.NMap; i++ {
		for j := 0; j < job.NReduce; j++ {
			common.RemoveFileRemote(common.SrvAddr(mapWks[i].address, constant.WORKER_FILE_RPC), common.IntermediateName(job.Name, i, j))
		}
	}
	for i := 0; i < job.NReduce; i++ {
		common.RemoveFileRemote(common.SrvAddr(redWks[i].address, constant.WORKER_FILE_RPC), common.ReduceName(job.Name, i))
	}
}
