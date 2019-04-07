package common

import "fmt"

const debugEnable = true

func Debug(format string, a ...interface{}) {
	if debugEnable {
		fmt.Printf(format+"\n", a...)
	}
}
