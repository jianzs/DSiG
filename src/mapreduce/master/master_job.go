package master

import "common"

func (mr *Master) SubmitJob(args *common.SubmitJobArgs, reply *common.SubmitJobReply) error {
	job := args.Job

	common.Debug("Master: Receive job: %s", job)

	err := mr.schedule(job)
	if err != nil {
		reply.Err = err
		return err
	}
	return nil
}
