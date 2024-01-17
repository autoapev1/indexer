package logger

import (
	"fmt"
	"time"
)

const LogTime = true

func Time(msg string, t time.Duration, s bool) {
	if LogTime {
		fmt.Printf("[TIME] %s \t\t%s\tsuccess=%v\n", msg, t.String(), s)
	}
}
