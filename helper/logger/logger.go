package logger

import (
	"fmt"
	"time"
)

func Info(log string) {
	fmt.Println(fmt.Sprint("[stdio_helper]", time.Now().Format("[2006-01-02 15:04:05]"), " [i] ", log))
}

func Warning(log string) {
	fmt.Println(fmt.Sprint("[stdio_helper]", time.Now().Format("[2006-01-02 15:04:05]"), " [w] ", log))
}
