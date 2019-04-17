package file

import "net"

type Keeper struct {
	port int

	ExitChannel chan struct{}
	l           net.Listener
}

func NewKeeper(port int) (kp *Keeper) {
	kp = new(Keeper)
	kp.port = port
	kp.ExitChannel = make(chan struct{}, 0)
	return
}
