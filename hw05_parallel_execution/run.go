package hw05parallelexecution

import (
	"errors"
	"sync"
)

var ErrErrorsLimitExceeded = errors.New("errors limit exceeded")

type Task func() error

type worker struct {
	jobs             <-chan Task
	countErrorToStop int
	wg               *sync.WaitGroup

	mu          *sync.Mutex
	countErrors *int
}

func newWorker(jobs <-chan Task, countErrors *int, countErrorToStop int, mu *sync.Mutex, wg *sync.WaitGroup) *worker {
	wg.Add(1)
	return &worker{
		jobs:             jobs,
		wg:               wg,
		mu:               mu,
		countErrors:      countErrors,
		countErrorToStop: countErrorToStop,
	}
}

func (w *worker) run() {
	defer w.wg.Done()
	for {
		if w.checkStop() {
			return
		}

		job, ok := <-w.jobs
		if !ok {
			return
		}
		if err := job(); err != nil {
			w.mu.Lock()
			*w.countErrors++
			w.mu.Unlock()
		}
	}
}

func (w *worker) checkStop() bool {
	w.mu.Lock()
	defer w.mu.Unlock()

	return *w.countErrors >= w.countErrorToStop
}

// Run starts tasks in n goroutines and stops its work when receiving m errors from tasks.
func Run(tasks []Task, n, m int) error {
	if n == 0 {
		return nil
	}

	if m == 0 {
		return ErrErrorsLimitExceeded
	}

	wg := sync.WaitGroup{}
	mu := sync.Mutex{}
	countError := 0

	jobs := make(chan Task, len(tasks))
	for _, task := range tasks {
		jobs <- task
	}
	close(jobs)

	for i := 0; i < n; i++ {
		w := newWorker(jobs, &countError, m, &mu, &wg)
		go w.run()
	}

	wg.Wait()
	if countError >= m {
		return ErrErrorsLimitExceeded
	}
	return nil
}
