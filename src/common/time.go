package common

import (
	"fmt"
	"time"
)

func GetNowTimestamp() string {
	return fmt.Sprintf("%d", time.Now().Unix())
}
