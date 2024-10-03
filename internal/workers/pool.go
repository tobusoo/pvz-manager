package workers

import (
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

	jobs    chan *TaskRequest
	Results chan *TaskResponse
	wg      sync.WaitGroup
}

func NewWorkers(workers uint) *Workers {
	w := &Workers{
		numWorkers:   workers,
		closedJobs:   false,
		closedResult: false,
		wg:           sync.WaitGroup{},
		jobs:         make(chan *TaskRequest, workers),
		Results:      make(chan *TaskResponse, workers),
	}

	w.wg.Add(int(workers))
	for i := uint(0); i < workers; i++ {
		go w.work()
	}

	return w
}

func (w *Workers) GetSize() uint {
	return w.numWorkers
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

func (w *Workers) work() {
	defer w.wg.Done()

	for task := range w.jobs {
		w.Results <- &TaskResponse{Err: task.Func(), Response: task.Request}
	}
}
