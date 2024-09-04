package worker

import (
	"github.com/xarest/gobs-collection/schema"
	gocronwork "github.com/xarest/gobs-collection/worker/gocron"
)

type IScheduler interface {
	AddTask(task *schema.Task) error
}

func NewIScheduler() IScheduler {
	return &gocronwork.Scheduler{}
}
