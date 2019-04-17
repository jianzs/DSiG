package common

// 包括RPC请求和响应的参数类型，以及RPC客户端Stub

import (
	"net/rpc"
)

type RegisterArgs struct {
	Worker string
}

type DoTaskArgs struct {
	TaskId   int
	JobName  string
	NOther   int
	Filename string
	Phase    string
}

type DoTaskReply struct {
	Code  int
	Error error
}

type ExecuteReply struct {
	Code int
	Err  error
}

type FileReply struct {
	Content string
	Err     error
}

type FileArgs struct {
	Filename string
	Content  string
}

type SubmitJobArgs struct {
	Job    Job
	Client string
}

type SubmitJobReply struct {
	Err error
}

type GetClientArgs struct {
	JobName string
}

type GetClientReply struct {
	Client string
}

type GetMapWkArgs struct {
	JobName string
}

type GetMapWkReply struct {
	Workers []string
}

func Call(srv string, rpcname string,
	args interface{}, reply interface{}) error {
	c, errx := rpc.Dial("tcp", srv)
	if errx != nil {
		return errx
	}
	defer c.Close()

	err := c.Call(rpcname, args, reply)
	if err == nil {
		return nil
	}
	return err
}
