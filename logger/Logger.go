package logger

import (
	"bytes"
	"log"
	"runtime"
)

func getGoroutineID() string {
	buf := make([]byte, 64)
	n := runtime.Stack(buf, false)
	buf = bytes.TrimSpace(buf[:n])
	return string(buf)
}

func Info(msg string) {
	log.Printf("[%s] - %s", getGoroutineID(), msg)
}
