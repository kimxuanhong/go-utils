package workflow

import (
	"context"
)

type Data struct {
	Request  interface{}
	Response interface{}
}

type Handler func(ctx context.Context, taskData *Data, err error)

type Task interface {
	Execute(ctx context.Context, taskData *Data, whenDone Handler)
	GetName() string
}
