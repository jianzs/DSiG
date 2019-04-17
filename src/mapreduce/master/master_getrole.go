package master

import "common"

func (mr *Master) GetClient(args *common.GetClientArgs, reply *common.GetClientReply) error {
	jobname := args.JobName
	reply.Client = mr.jobs[jobname].client
	return nil
}

func (mr *Master) GetMapWorkers(args *common.GetMapWkArgs, reply *common.GetMapWkReply) error {
	jobname := args.JobName
	for _, wk := range mr.jobs[jobname].mapWorkers {
		reply.Workers = append(reply.Workers, wk.address)
	}
	return nil
}
