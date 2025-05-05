package workflow

import (
	"context"
	"fmt"
	"github.com/kimxuanhong/go-utils/safe"
	"log"
)

type Workflow struct {
	Name   string
	Tasks  []Task
	Result *Data
	Error  error
}

func (wf *Workflow) AddTask(task Task) {
	wf.Tasks = append(wf.Tasks, task)
}

func (wf *Workflow) Run(ctx context.Context, taskData *Data, whenDone Handler) {
	go func(ctx context.Context, taskData *Data) {
		safe.SafeGo(func(ex error) {
			if ex != nil {
				whenDone(ctx, taskData, ex)
				return
			}

			log.Printf("---------------------- Workflow %s starting! ----------------------\n", wf.Name)
			wf.Result = taskData
			for _, taskStep := range wf.Tasks {
				log.Printf("Run %s", taskStep.GetName())
				taskChannel := make(chan struct {
					*Data
					error
				}, 1)

				go func(taskStep Task) {
					safe.SafeGo(func(ex error) {
						if ex != nil {
							taskChannel <- struct {
								*Data
								error
							}{taskData, ex}
							return
						}

						select {
						case <-ctx.Done():
							taskChannel <- struct {
								*Data
								error
							}{wf.Result, fmt.Errorf("context cancled")}
						default:
							taskStep.Execute(ctx, wf.Result, func(ctx context.Context, result *Data, err error) {
								taskChannel <- struct {
									*Data
									error
								}{result, err}
							})
						}
					})
				}(taskStep)

				result := <-taskChannel
				close(taskChannel)

				if result.error != nil {
					wf.Error = result.error
					wf.Result = result.Data
					whenDone(ctx, wf.Result, wf.Error)
					return
				}
				wf.Result = result.Data
			}

			log.Printf("---------------------- Workflow %s success! ----------------------\n", wf.Name)
			whenDone(ctx, wf.Result, wf.Error)
		})

	}(ctx, taskData)
}

func NewWorkflow(workflowName string) *Workflow {
	return &Workflow{
		Name:  workflowName,
		Tasks: make([]Task, 0),
	}
}
