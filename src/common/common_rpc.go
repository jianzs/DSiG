package common

// 包括RPC请求和响应的参数类型，以及RPC客户端Stub

import (
	"net/rpc"
)

type RegisterArgs struct {
	Worker string
}

type DoTaskArgs struct {
	TaskId    int
	JobName   string
	NOther    int
	Filename  string
	Timestamp string
	Phase     string
}

type DoTaskReply struct {
	Code  int
	Error error
}

type ExecuteReply struct {
	Code int
	Err  error
}

type ReadFileArgs struct {
	Filename  string
	JobName   string
	Timestamp string
}

type ReadFileReply struct {
	Content string
	Err     error
}

type WriteFileArgs struct {
	Filename string
	Content  string
}

type WriteFileReply struct {
	Err error
}

type SubmitJobArgs struct {
	Job Job
}

type SubmitJobReply struct {
	Err error
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
