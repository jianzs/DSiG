package master

import (
	"net"
	"sync"
)

type Master struct {
	sync.Mutex

	doneChannel chan bool // 一个任务的结束

	// protected by the mutex
	newCond *sync.Cond // signals when Register() adds to workers[]
	workers []*Worker  // each worker's IP socket address -- its RPC address

	jobs map[string]*Job

	ExitChannel chan struct{}
	l           net.Listener
}

func NewMaster() (mr *Master) {
	mr = new(Master)
	mr.doneChannel = make(chan bool)
	mr.newCond = sync.NewCond(mr)
	mr.jobs = make(map[string]*Job)
	mr.ExitChannel = make(chan struct{})
	return
}
