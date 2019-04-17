package master

import (
	"common"
	"constant"
	"net"
	"net/rpc"
)

func (mr *Master) StartRPCServer() error {
	server := rpc.NewServer()
	server.Register(mr)
	server.HandleHTTP(constant.MASTER_RPC_PATH, constant.MASTER_DEBUG_PATH)

	l, err := net.Listen("tcp", common.ListenPort(constant.MASTER_RPC))
	if err != nil {
		return err
	}

	mr.l = l
	go func() {
		for {
			conn, err := mr.l.Accept()
			if err != nil {
				common.Debug("Master: RPC Server Connect Error %s", err)
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

func (mr *Master) StopRPCServer() {
	mr.l.Close()
	common.Debug("Master: RPC Server Stopped")
}
