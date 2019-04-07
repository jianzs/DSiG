package mapreduce

import (
	"common"
	"constant"
	"strconv"
	"sync"
)

func (mr *Master) schedule(job common.Job) error {
	defer mr.CleanupFiles(job) // 结束后，删除所有中间文件

	ch := make(chan string, 0)
	go forwardRegistrations(mr, ch)

	errs := make([]error, 0)

	var wg sync.WaitGroup
	for i := 0; i < job.NMap; i++ {
		wg.Add(1)
		go func(i int) {
			for {
				wk := <-ch

				common.Debug("Master: Schedule Map#%d to %s", i, wk)

				args := common.DoTaskArgs{i, job.Name,
					job.NReduce, job.InFiles[i], job.Timestamp, constant.MAP_PHASE}
				var reply common.DoTaskReply
				err := common.Call(common.SrvAddr(wk, constant.WORKER_RPC), constant.DO_TASK, args, &reply)

				if err != nil {
					common.Debug("Master: Worker do map error %s", err)
					errs = append(errs, err)
					return
				}

				if reply.Code != constant.SUCCESS {
					common.Debug("Master: Worker Do Map Error %s",
						&TaskError{reply.Code, reply.Error})
				} else {
					wg.Done()
					ch <- wk
					common.Debug("Master: Schedule Map#%d to %s done successfully", i, wk)
					break
				}
			}
		}(i)
	}
	wg.Wait()
	if len(errs) > 0 {
		return errs[0]
	}

	for i := 0; i < job.NReduce; i++ {
		wg.Add(1)
		go func(i int) {
			for {
				wk := <-ch

				common.Debug("Master: Schedule Reduce#%d to %s", i, wk)

				args := common.DoTaskArgs{i, job.Name, job.NMap,
					"", job.Timestamp, constant.REDUCE_PHASE}
				var reply common.DoTaskReply

				err := common.Call(common.SrvAddr(wk, constant.WORKER_RPC), constant.DO_TASK, args, &reply)

				if err != nil {
					errs = append(errs, err)
					return
				}

				if reply.Code != constant.SUCCESS {
					common.Debug("Worker: Do Reduce Error %s",
						&TaskError{reply.Code, reply.Error})
				} else {
					wg.Done()
					ch <- wk
					common.Debug("Master: Schedule Reduce#%d to %s done successfully", i, wk)
					break
				}
			}
		}(i)
	}
	wg.Wait()

	if len(errs) > 0 {
		return errs[0]
	}

	// Merge Reduces' result
	reduceFiles := make([]string, 0)
	for i := 0; i < job.NReduce; i++ {
		reduceFiles = append(reduceFiles, common.ReduceName(job.Name, job.Timestamp, i))
	}
	outFile := common.FinalName(job.Name, job.Timestamp, job.OutFile)
	err := mr.mergeFiles(reduceFiles, outFile)
	if err != nil {
		return err
	}

	return nil
}

func forwardRegistrations(mr *Master, ch chan string) {
	i := 0
	for {
		mr.Lock()
		if i < len(mr.workers) {
			ch <- mr.workers[i]
			i++
		} else {
			mr.newCond.Wait()
		}
		mr.Unlock()
	}
}

type TaskError struct {
	Code int
	Err  error
}

func (e *TaskError) Error() string { return strconv.Itoa(e.Code) + " " + e.Err.Error() }
