package app

import "sync"

type WorkerPool struct {
	workerCount int
	taskQueue   chan func()
	wg          sync.WaitGroup
}

func NewWorkerPool(workerCount int) *WorkerPool {
	return &WorkerPool{
		workerCount: workerCount,
		taskQueue:   make(chan func(), workerCount),
	}
}

func (p *WorkerPool) Start() {
	for i := 0; i < p.workerCount; i++ {
		p.wg.Add(1)
		go p.worker()
	}
}

func (p *WorkerPool) worker() {
	defer p.wg.Done()
	for task := range p.taskQueue {
		task()
	}
}

func (p *WorkerPool) SubmitTask(task func()) {
	p.taskQueue <- task
}

func (p *WorkerPool) Stop() {
	close(p.taskQueue)
	p.wg.Wait()
}
