package common

// 通用的函数

import (
	"bytes"
	"encoding/json"
	"io"
	"strconv"
	"strings"
)

// RPC Server 监听的端口
func ListenPort(port int) string {
	return ":" + strconv.Itoa(port)
}

// RPC 请求的地址
func SrvAddr(ip string, port int) string {
	return ip + ":" + strconv.Itoa(port)
}

func ExecutorFile(jobName string) string {
	return jobName + "/main/executor_bin.go"
}

func ExecutorLogFile(jobName string, phase string, id int) string {
	return jobName + "/executor-" + phase + "-" + strconv.Itoa(id) + ".out"
}

// the map's result
func IntermediateName(jobName string, mapId, redId int) string {
	return jobName + "/mrtmp." + jobName + "-" + strconv.Itoa(mapId) + "-" + strconv.Itoa(redId)
}

// the reduce's result
func ReduceName(jobName string, reduceTask int) string {
	return jobName + "/mrtmp." + jobName + "-res-" + strconv.Itoa(reduceTask)
}

func FinalName(jobName, outFile string) string {
	return jobName + "/" + outFile
}

func Encode(kvs []KeyValue) (string, error) {
	jsons := make([]string, 0)
	for _, kv := range kvs {
		json, err := json.Marshal(kv)
		if err != nil {
			return "", err
		}

		jsons = append(jsons, string(json))
	}

	return string(strings.Join(jsons, "\n")), nil
}

func Decode(str string) ([]KeyValue, error) {
	buffer := bytes.NewBufferString(str)
	dc := json.NewDecoder(buffer)

	kvs := make([]KeyValue, 0)
	for {
		var kv KeyValue
		err := dc.Decode(&kv)
		if err == io.EOF {
			break
		} else if err != nil {
			return nil, err
		}
		kvs = append(kvs, kv)
	}
	return kvs, nil
}
