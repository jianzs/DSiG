package master

import (
	"common"
	"os"
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
		data, err := common.ReadFileRemote(wk.address, common.ReduceName(jobName, idx))
		if err != nil {
			return err
		}
		err = common.AppendFileRemote(job.client, common.FinalName(jobName, job.job.OutFile), data)
		if err != nil {
			return err
		}
	}
	return nil
}

func (mr *Master) CleanupFiles(job common.Job) {
	for i := 0; i < job.NMap; i++ {
		for j := 0; j < job.NReduce; j++ {
			removeFile(common.IntermediateName(job.Name, i, j))
		}
	}
	for i := 0; i < job.NReduce; i++ {
		removeFile(common.ReduceName(job.Name, i))
	}
}

// removeFile is a simple wrapper around os.Remove that logs errors.
func removeFile(n string) error {
	err := os.Remove(n)
	if err == os.ErrNotExist {
		return nil
	} else if err != nil {
		return err
	}
	return nil
}
