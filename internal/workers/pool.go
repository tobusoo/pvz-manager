package workers

import (
	"fmt"
	"sync"
)

type TaskRequest struct {
	Request string
	Func    func() error
}

type TaskResponse struct {
	Response string
	Err      error
}

type Workers struct {
	numWorkers   uint
	closedJobs   bool
	closedResult bool

	jobs     chan *TaskRequest
	to_close chan struct{}
	Results  chan *TaskResponse
	wg       sync.WaitGroup
}

func NewWorkers(workers uint) *Workers {
	w := &Workers{
		numWorkers:   0,
		closedJobs:   false,
		closedResult: false,
		wg:           sync.WaitGroup{},
		jobs:         make(chan *TaskRequest, workers),
		to_close:     make(chan struct{}, workers),
		Results:      make(chan *TaskResponse, workers),
	}

	w.AddWorkers(workers)

	return w
}

func (w *Workers) GetSize() uint {
	return w.numWorkers
}

func (w *Workers) AddWorkers(count uint) {
	w.wg.Add(int(count))
	for i := uint(0); i < count; i++ {
		go w.work()
	}
	w.numWorkers += count
}

func (w *Workers) CloseNworkers(count uint) error {
	if count > w.numWorkers {
		return fmt.Errorf("can't close %d workers: current workers count = %d", count, w.numWorkers)
	}

	for i := uint(0); i < count; i++ {
		w.to_close <- struct{}{}
	}
	w.numWorkers -= count
	return nil
}

func (w *Workers) AddTask(task *TaskRequest) {
	w.jobs <- task
}

func (w *Workers) CloseJobs() {
	if !w.closedJobs {
		close(w.jobs)
	}
	w.closedJobs = true
}

func (w *Workers) Wait() {
	w.wg.Wait()
	if !w.closedResult {
		close(w.Results)
	}
	w.closedResult = true
}

//gocognit:ignore
func (w *Workers) work() {
	defer w.wg.Done()
	for {
		select {
		case task, ok := <-w.jobs:
			if !ok {
				return
			}
			w.Results <- &TaskResponse{Err: task.Func(), Response: task.Request}
		case <-w.to_close:
			return
		}
	}
}
