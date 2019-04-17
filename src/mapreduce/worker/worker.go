package worker

import (
	"net"
	"sync"
)

type Worker struct {
	sync.Mutex

	address    string
	concurrent int // number of parallel DoTasks in this worker; mutex
	l          net.Listener

	master string

	ExitChannel chan struct{}
}

func NewWorker(ip string) (wk *Worker) {
	wk = new(Worker)
	wk.address = ip
	wk.concurrent = 0
	wk.ExitChannel = make(chan struct{})
	return
}
