package logger

import (
	"fmt"
	"time"
)

const LogTime = false

func Time(msg string, t time.Duration, s bool) {
	if LogTime {
		fmt.Printf("[TIME] %s \t\t%s\tsuccess=%v", msg, t.String(), s)
	}
}
