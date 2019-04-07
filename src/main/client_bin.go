package main

import (
	"common"
	"fmt"
	"mapreduce"
	"os"
)

func main() {
	args := os.Args
	mapRedFile := args[1] + ".go"
	inFiles := args[2:]

	var job common.Job
	job.FuncFile = mapRedFile
	job.NReduce = 5
	job.InFiles = inFiles
	job.NMap = len(job.InFiles)
	job.Name = "test"
	job.OutFile = "testOut"

	client := mapreduce.NewClient(job, "127.0.0.1")
	err := client.Submit()
	if err != nil {
		fmt.Println(err)
	} else {
		str, err := client.GetResult()
		if err != nil {
			fmt.Println(err)
		} else {
			fmt.Println(str)
		}
	}
}
