package master

import "common"

type Job struct {
	job          common.Job
	mapWorkers   []*Worker
	reduceWorker []*Worker
	client       string
}

func newJob(job common.Job, client string) *Job {
	j := &Job{job: job, client: client}
	j.reduceWorker = make([]*Worker, job.NReduce)
	j.mapWorkers = make([]*Worker, job.NMap)
	return j
}

func (mr *Master) SubmitJob(args *common.SubmitJobArgs, reply *common.SubmitJobReply) error {
	job := args.Job

	mr.jobs[job.Name] = newJob(job, args.Client)

	common.Debug("Master: Receive job: %s", job)

	err := mr.schedule(job)
	if err != nil {
		reply.Err = err
		return err
	}
	return nil
}
