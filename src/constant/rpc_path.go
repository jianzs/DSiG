package constant

import "net/rpc"

const (
	MASTER_RPC_PATH   = rpc.DefaultRPCPath + "/master"
	MASTER_DEBUG_PATH = rpc.DefaultDebugPath + "/master"
	WORKER_PRC_PATH   = rpc.DefaultRPCPath + "/worker"
	WORKER_DEBUG_PATH = rpc.DefaultDebugPath + "/worker"
	KEEPER_PRC_PATH   = rpc.DefaultRPCPath + "/keeper"
	KEEPER_DEBUG_PATH = rpc.DefaultDebugPath + "/keeper"
)
