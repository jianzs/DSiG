package file

import "net"

type Keeper struct {
	ExitChannel chan struct{}
	l           net.Listener
}

func NewKeeper() (kp *Keeper) {
	kp = new(Keeper)
	kp.ExitChannel = make(chan struct{}, 0)
	return
}
