package master

import "time"

type Worker struct {
	address       string
	lastHeartbeat int64
	activeTaskNum int
	taskCount     int
}

func newWorker(address string) (wk *Worker) {
	wk = &Worker{address:address,lastHeartbeat:time.Now().Unix()}
	return
}

func (wk *Worker) less(o *Worker) int {
	return wk.activeTaskNum - o.activeTaskNum
}

