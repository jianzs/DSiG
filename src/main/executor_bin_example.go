package main

import (
	"common"
	"constant"
	"fmt"
	"hash/fnv"
	"net"
	"net/rpc"
	"os"
	"sort"
	"strconv"
	"strings"
)

// 用于注册RPC
type Executor struct {
	l      net.Listener
	master string
}

func NewExecutor(mr string) Executor {
	return Executor{master: mr}
}

func main() {
	args := os.Args
	exe := NewExecutor(args[1]) // 命令行得知Master地址
	err := exe.StartRPCServer()
	if err != nil {
		fmt.Printf("Executor: Start RPC Server error %s\n", err)
		return
	}
	fmt.Println("Executor: Start RPC server successfully!")

	exe.StopRPCServer()
	fmt.Println("Executor: End RPC Server successfully")
}

func (exe *Executor) StartRPCServer() error {
	server := rpc.NewServer()
	err := server.Register(exe)
	if err != nil {
		return err
	}
	server.HandleHTTP(rpc.DefaultRPCPath, rpc.DefaultDebugPath)

	l, err := net.Listen("tcp", common.ListenPort(constant.EXECUTOR_RPC))
	if err != nil {
		return err
	}

	exe.l = l
	conn, err := exe.l.Accept()
	if err != nil {
		common.Debug("Executor: RPC Server Connect TaskError %s", err)
	} else {
		server.ServeConn(conn)
		conn.Close()
	}

	return nil
}

func (exe *Executor) StopRPCServer() {
	exe.l.Close()
}

func (exe *Executor) DoTask(args *common.DoTaskArgs, reply *common.ExecuteReply) error {
	id := args.TaskId
	jobName := args.JobName
	timestamp := args.Timestamp
	nOther := args.NOther
	common.Debug("Executor: Receive %s#%d ", args.Phase, id)

	switch args.Phase {
	case constant.MAP_PHASE:
		err := doMap(jobName, timestamp, exe.master, id, nOther, args.Filename, MapFunc, reply)
		if err != nil {
			return err
		}
	case constant.REDUCE_PHASE:
		err := doReduce(jobName, timestamp, exe.master, id, nOther, ReduceFunc, reply)
		if err != nil {
			return err
		}
	}
	common.Debug("Executor: %s#%d done successfully", args.Phase, id)

	return nil
}

func doMap(jobName, timestamp, mr string, id, nReduce int, filename string,
	mapFunc func(string, string) []common.KeyValue,
	reply *common.ExecuteReply) error {

	// Read file
	data, err := common.ReadFileFrMaster(mr, filename)
	if err != nil {
		common.Debug("Executor: Read file failed %s", err)
		reply.Code = constant.READ_ERROR
		reply.Err = err
		return err
	}
	common.Debug("Executor: Read file successfully")

	// Run Map
	kvs := mapFunc(filename, data)
	common.Debug("Executor: MapFunc done successfully")

	// Grouping
	groupRes := make([][]common.KeyValue, nReduce)
	for i := 0; i < nReduce; i++ {
		groupRes[i] = make([]common.KeyValue, 0)
	}
	for _, kv := range kvs {
		r := ihash(kv.Key) % nReduce
		groupRes[r] = append(groupRes[r], kv)
	}
	common.Debug("Executor: Grouping done successfully")

	// Write file
	for i := 0; i < nReduce; i++ {
		// encode
		json, err := common.Encode(groupRes[i])
		if err != nil {
			common.Debug("Executor: Encode failed %s", err)
			reply.Code = constant.ENCODER_ERROR
			reply.Err = err
			return err
		}

		itmdName := common.IntermediateName(jobName, timestamp, id, i)

		err = common.WriteFileToMaster(mr, itmdName, json)
		if err != nil {
			common.Debug("Executor: Write file %s, error %s", itmdName, err)
			reply.Code = constant.WRITE_ERROR
			reply.Err = err
			return err
		}
		common.Debug("Executor: Write file '%s' successfully", itmdName)
	}

	reply.Code = constant.SUCCESS
	return nil
}

func ihash(s string) int {
	h := fnv.New32a()
	h.Write([]byte(s))
	return int(h.Sum32() & 0x7fffffff)
}

func doReduce(jobName, timestamp, mr string, id, nMap int,
	redFunc func(string, []string) string,
	reply *common.ExecuteReply) error {

	// Read Intermediate Files
	kvs := make([]common.KeyValue, 0)
	for i := 0; i < nMap; i++ {
		itmdName := common.IntermediateName(jobName, timestamp, i, id)
		content, err := common.ReadFileFrMaster(mr, itmdName)
		if err != nil {
			common.Debug("Executor: Read file %s failed, error %s", itmdName, err)
			reply.Code = constant.READ_ERROR
			reply.Err = err
			return err
		}

		kv, err := common.Decode(content)
		if err != nil {
			common.Debug("Executor: Decode failed %s", err)
			reply.Code = constant.DECODER_ERROR
			reply.Err = err
			return err
		}

		kvs = append(kvs, kv...)
	}

	// Sort
	sort.Slice(kvs, func(i, j int) bool {
		return kvs[i].Key < kvs[j].Key
	})

	// Group Reduce
	resKvs := make([]common.KeyValue, 0)
	curKey := ""
	curVals := make([]string, 0)
	for _, kv := range kvs {
		if kv.Key != curKey {
			if curKey != "" {
				// 执行 Reduce 得到结果
				redRes := redFunc(curKey, curVals)
				resKvs = append(resKvs, common.KeyValue{curKey, redRes})
			}

			curKey = kv.Key
			curVals = curVals[0:0]
		}
		curVals = append(curVals, kv.Value)
	}
	redRes := redFunc(curKey, curVals)
	resKvs = append(resKvs, common.KeyValue{curKey, redRes})

	// encode
	jsonstr, err := common.Encode(resKvs)
	if err != nil {
		common.Debug("Executor: Encode failed %s", err)
		reply.Code = constant.ENCODER_ERROR
		reply.Err = err
		return err
	}

	// Write File
	err = common.WriteFileToMaster(mr, common.ReduceName(jobName, timestamp, id), jsonstr)
	if err != nil {
		common.Debug("Executor: Write file  error %s", err)
		reply.Code = constant.WRITE_ERROR
		reply.Err = err
		return err
	}
	common.Debug("Executor: Write file successfully")

	return nil
}

// 用户自定义的Map函数
func MapFunc(file string, value string) (res []common.KeyValue) {
	words := strings.FieldsFunc(value, func(r rune) bool {
		return !((r >= 'A' && r <= 'Z') || (r >= 'a' && r <= 'z'))
	})
	for _, w := range words {
		kv := common.KeyValue{w, "1"}
		res = append(res, kv)
	}
	return
}

// 用户自定义的Reduce函数
func ReduceFunc(key string, values []string) string {
	cnt := 0
	for _, val := range values {
		tmp, err := strconv.Atoi(val)
		if err != nil {
			fmt.Println(err)
			continue
		}
		cnt += tmp
	}

	return strconv.Itoa(cnt)
}
