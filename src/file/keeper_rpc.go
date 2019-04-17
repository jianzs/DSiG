package file

import (
	"common"
	"constant"
	"net"
	"net/rpc"
)

func (kp *Keeper) StartRPCServer() error {
	server := rpc.NewServer()
	server.Register(kp)
	server.HandleHTTP(constant.KEEPER_PRC_PATH, constant.KEEPER_DEBUG_PATH)

	l, err := net.Listen("tcp", common.ListenPort(kp.port))
	if err != nil {
		return err
	}

	kp.l = l
	go func() {
		for {
			conn, err := kp.l.Accept()
			if err != nil {
				common.Debug("FILE: RPC Server Connect Error %s", err)
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

func (kp *Keeper) StopRPCServer() {
	kp.l.Close()
	common.Debug("FILE: RPC Server Stopped")
}
