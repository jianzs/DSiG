package worker

import (
	"common"
	"constant"
	"net"
	"os"
	"os/exec"
	"time"
)

func (wk *Worker) DoTask(args *common.DoTaskArgs, reply *common.DoTaskReply) error {
	common.Debug("Worker: Receive %s#%d", args.Phase, args.TaskId)

	timestamp := args.Timestamp
	jobName := args.JobName
	executorFile := common.ExecutorFile(jobName, timestamp)
	logFileName := common.ExecutorLogFile(jobName, timestamp, args.Phase, args.TaskId)

	// 保证Executor存在
	err := ensureExecutor(executorFile, wk.master, reply)
	if err != nil {
		return err
	}

	// 启动Executor
	common.Debug("Worker: Start run executor %s", executorFile)
	err = startExecutor(wk.master, executorFile, logFileName, reply)
	if err != nil {
		common.Debug("Worker: Executor started failed %s", err)
		return err
	}
	common.Debug("Worker: Executor is executed ")

	// 通过RPC调用Executor
	common.Debug("Worker: Send %s#%d Task to executor", args.Phase, args.TaskId)
	var exeReply common.ExecuteReply
	// 避免Executor刚启动，RPC Server还没完成启动
	failedCnt := 0
	for {
		err = common.Call(common.SrvAddr("127.0.0.1", constant.EXECUTOR_RPC), constant.EXECUTE_TASK, args, &exeReply)
		if _, ok := err.(*net.OpError); ok {
			time.Sleep(time.Millisecond * 50)
			failedCnt++
		} else {
			break
		}
		if failedCnt > 20 {
			break
		}
	}
	if err != nil || exeReply.Err != nil {
		if err == nil {
			reply.Code = exeReply.Code
			reply.Error = exeReply.Err
		} else {
			reply.Code = constant.OTHER_ERROR
			reply.Error = err
		}
		return err
	}
	common.Debug("Worker: %s#%d task is finished successfully", args.Phase, args.TaskId)

	reply.Code = constant.SUCCESS
	return nil
}

func startExecutor(mr, executorFile, logName string, reply *common.DoTaskReply) error {
	cmd := exec.Command("go", "run", executorFile, mr)

	logFile, err := os.OpenFile(logName, os.O_WRONLY|os.O_CREATE, 0666)
	if err != nil {
		reply.Code = constant.OTHER_ERROR
		reply.Error = err
		return err
	}
	cmd.Stdout = logFile

	err = cmd.Start()
	if err != nil {
		reply.Code = constant.OTHER_ERROR
		reply.Error = err
		return err
	}
	return nil
}

func ensureExecutor(executorFile string, mr string, reply *common.DoTaskReply) error {
	if !common.Exists(executorFile) {
		common.Debug("Worker: Get Func File")
		content, err := common.ReadFileFrMaster(mr, executorFile)
		if err != nil {
			reply.Code = constant.READ_ERROR
			reply.Error = err
			return err
		}

		common.Debug("Worker: Write function file %s", executorFile)
		err = common.WriteFile(executorFile, content)
		if err != nil {
			reply.Code = constant.WRITE_ERROR
			reply.Error = err
			return err
		}
		common.Debug("Worker: Get function file successfully")
	}

	return nil
}
