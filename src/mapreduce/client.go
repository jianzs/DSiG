package mapreduce

import (
	"common"
	"constant"
)

type Client struct {
	job    common.Job
	master string
}

func NewClient(job common.Job, mr string) Client {
	return Client{job, mr}
}

func (ct *Client) Submit() error {
	ct.job.Timestamp = common.GetNowTimestamp()

	content, err := common.ReadFile(ct.job.FuncFile)
	if err != nil {
		return err
	}
	common.Debug("Client: Read function file successfully")

	err = common.WriteFileToMaster(ct.master, common.ExecutorFile(ct.job.Name, ct.job.Timestamp), content)
	if err != nil {
		return err
	}
	common.Debug("Client: Write File successfully")

	args := common.SubmitJobArgs{Job: ct.job}
	var reply common.SubmitJobReply

	common.Debug("Client: Start to submit job")
	err = common.Call(common.SrvAddr(ct.master, constant.MASTER_RPC), constant.SUBMIT_JOB, args, &reply)
	if err != nil {
		return err
	}
	if reply.Err != nil {
		return reply.Err
	}
	common.Debug("Client: Job executed successfully")

	return nil
}

func (ct *Client) GetResult() (content string, err error) {
	filename := common.FinalName(ct.job.Name, ct.job.Timestamp, ct.job.OutFile)
	content, err = common.ReadFileFrMaster(ct.master, filename)
	return
}
