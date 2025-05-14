package safe

import (
	"fmt"
	"log"
	"runtime/debug"
)

func PanicHandler(whenDone func(err error)) {
	if r := recover(); r != nil {
		log.Printf("Panic recovered: %v\nStack trace: %s", r, string(debug.Stack()))
		whenDone(fmt.Errorf("Recovered from panic: %v\n", r))
		return
	}
}

func SafeGo(fnc func(ex error)) {
	defer PanicHandler(func(err error) {
		fnc(err)
		return
	})
	fnc(nil)
}

func GetOrDefault[T any](fn func() T, defaultValue T) (result T) {
	defer func() {
		if r := recover(); r != nil {
			log.Printf("GetOrDefault panic recovered: %v\nStack trace: %s", r, string(debug.Stack()))
			result = defaultValue
		}
	}()
	return fn()
}
