package worker

import (
	"common"
	"constant"
	"net"
	"net/rpc"
)

func (wk *Worker) StartRPCServer() error {
	server := rpc.NewServer()
	server.Register(wk)
	server.HandleHTTP(constant.WORKER_PRC_PATH, constant.WORKER_DEBUG_PATH)

	l, err := net.Listen("tcp", common.ListenPort(constant.WORKER_RPC))
	if err != nil {
		return err
	}

	wk.l = l
	go func() {
		for {
			conn, err := wk.l.Accept()
			if err != nil {
				common.Debug("Worker: RPC Server Connect TaskError %s", err)
				break
			} else {
				go func() {
					server.ServeConn(conn)
					conn.Close()
				}()
			}
		}
	}()

	return nil
}

func (wk *Worker) StopRPCServer() {
	close(wk.ExitChannel)
	wk.l.Close()
	common.Debug("Worker: RPC Server Stopped")
}
