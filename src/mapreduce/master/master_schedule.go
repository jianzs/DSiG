package master

import (
	"common"
	"constant"
	"strconv"
	"sync"
)

func (mr *Master) schedule(job common.Job) error {
	defer mr.CleanupFiles(job) // 结束后，删除所有中间文件

	// Worker Channel
	ch := make(chan *Worker, 0)
	go forwardRegistrations(mr, ch)

	// Map Task
	var wg sync.WaitGroup
	for i := 0; i < job.NMap; i++ {
		wg.Add(1)
		go mr.handOutTask(i, &job, constant.MAP_PHASE, ch, &wg)
	}
	wg.Wait()

	// Reduce Task
	for i := 0; i < job.NReduce; i++ {
		wg.Add(1)
		go mr.handOutTask(i, &job, constant.REDUCE_PHASE, ch, &wg)
	}
	wg.Wait()

	// Merge Reduces' result
	err := mr.mergeFiles(job.Name)
	if err != nil {
		return err
	}
	common.Debug("Master: Merge Output file successfully")
	common.Debug("Master: %s done successfully", job.Name)

	return nil
}

func (mr *Master) handOutTask(taskId int, job *common.Job,
	phase string, wkCh chan *Worker, wg *sync.WaitGroup) (*common.DoTaskReply, error) {
	var nOther int

	switch phase {
	case constant.MAP_PHASE:
		nOther = job.NReduce
	case constant.REDUCE_PHASE:
		nOther = job.NMap
	}

	args := common.DoTaskArgs{taskId, job.Name, nOther, job.InFiles[taskId], phase}
	var reply common.DoTaskReply
	var err error

	tryCnt := 0
	for {
		if tryCnt >= 3 {
			break
		}
		tryCnt++

		// Update worker task info
		wk := <-wkCh
		wk.Lock()
		wk.taskCount++
		wk.activeTaskNum++
		wk.Unlock()

		common.Debug("Master: Schedule %s#%d to %s", phase, taskId, wk)
		err = common.Call(common.SrvAddr(wk.address, constant.WORKER_RPC), constant.DO_TASK, args, &reply)

		if err != nil {
			common.Debug("Master: Worker do %s error %s", phase, err)
			continue
		}

		if reply.Code != constant.SUCCESS {
			common.Debug("Master: Worker Do %s Error %s",
				phase, &TaskError{reply.Code, reply.Error})
			err = reply.Error
		} else {
			mr.Lock()
			switch phase {
			case constant.MAP_PHASE:
				mr.jobs[job.Name].mapWorkers[taskId] = wk
			case constant.REDUCE_PHASE:
				mr.jobs[job.Name].reduceWorker[taskId] = wk
			}
			mr.Unlock()

			wg.Done()
			wkCh <- wk
			common.Debug("Master: Schedule %s#%d to %s done successfully", phase, taskId, wk)
			break
		}
	}

	if tryCnt <= 3 {
		return &reply, nil
	} else {
		return nil, err
	}
}

func forwardRegistrations(mr *Master, ch chan *Worker) {
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
