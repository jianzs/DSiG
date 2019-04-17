package master

import (
	"common"
	"os"
)

func (mr *Master) ReadFile(args *common.ReadFileArgs, reply *common.ReadFileReply) error {
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

func (mr *Master) WriteFile(args *common.WriteFileArgs, reply *common.WriteFileReply) error {
	common.Debug("Master: Receive RPC write file '%s'", args.Filename)
	err := common.WriteFile(args.Filename, args.Content)
	if err != nil {
		reply.Err = err
		return err
	}
	common.Debug("Master: write file '%s' successfully", args.Filename)
	return nil
}

func (mr *Master) mergeFiles(inFiles []string, outFile string) error {
	for _, file := range inFiles {
		data, err := common.ReadFile(file)
		if err != nil {
			return err
		}

		err = common.AppendFile(outFile, data)
		if err != nil {
			return err
		}
	}
	return nil
}

func (mr *Master) CleanupFiles(job common.Job) {
	for i := 0; i < job.NMap; i++ {
		for j := 0; j < job.NReduce; j++ {
			removeFile(common.IntermediateName(job.Name, job.Timestamp, i, j))
		}
	}
	for i := 0; i < job.NReduce; i++ {
		removeFile(common.ReduceName(job.Name, job.Timestamp, i))
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
