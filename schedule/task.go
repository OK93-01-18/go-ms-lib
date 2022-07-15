package schedule

import (
	"sync"
)

type TaskInterface interface {
	Do() error
}

type Task struct {
	wg       *sync.WaitGroup
	shutdown *bool
}

func (t *Task) run(f func() error) error {
	if *t.shutdown {
		return nil
	}
	t.wg.Add(1)
	defer t.wg.Done()
	return f()
}

func NewTask(wg *sync.WaitGroup, shutdown *bool) *Task {
	return &Task{
		wg:       wg,
		shutdown: shutdown,
	}
}
