package common

import (
	"fmt"
	"strings"
)

// Job 一次任务的所有内容

type Job struct {
	Name string

	NMap    int
	NReduce int

	InFiles []string
	OutFile string
	Status  int

	FuncFile string

	Timestamp string
}

func (b Job) String() string {
	return fmt.Sprintf("{name:%s, nMap:%d, nReduce:%d, inFiles:%s, outFile:%s, funcFile:%s}",
		b.Name, b.NMap, b.NReduce,
		strings.Join(b.InFiles, ","), b.OutFile, b.FuncFile)
}
