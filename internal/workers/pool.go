package workers

import (
	"strconv"
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
	numWorkers int

	jobs    chan TaskRequest
	Results chan TaskResponse
	wg      sync.WaitGroup
}

func NewWorkers(workers int) *Workers {
	w := &Workers{
		numWorkers: workers,
		wg:         sync.WaitGroup{},
		jobs:       make(chan TaskRequest, workers),
		Results:    make(chan TaskResponse, workers),
	}

	w.wg.Add(workers)
	for i := 0; i < workers; i++ {
		go w.work(i)
	}

	return w
}

func (w *Workers) AddTask(task TaskRequest) {
	w.jobs <- task
}

func (w *Workers) Close() {
	close(w.jobs)
}

func (w *Workers) Wait() {
	w.wg.Wait()
	close(w.Results)
}

func (w *Workers) work(id int) {
	defer w.wg.Done()

	for task := range w.jobs {
		w.Results <- TaskResponse{Err: task.Func(), Response: "worker: " + strconv.Itoa(id) + " " + task.Request}
	}
}
