package workerpool

import (
	"context"
	"fmt"
	"sync"
	"time"
)

type WorkerPool struct {
	workersCount int
	Jobs         chan string
	Result       chan string
	wg           *sync.WaitGroup
	ctx          context.Context
	cancel       context.CancelFunc
	mtx          sync.Mutex
}

func NewWorkerPool(workerNum int, ctx context.Context) *WorkerPool {
	ctx, cancel := context.WithCancel(ctx)
	return &WorkerPool{
		workersCount: workerNum,
		Jobs:         make(chan string),
		Result:       make(chan string),
		wg:           &sync.WaitGroup{},
		ctx:          ctx,
		cancel:       cancel,
	}
}

func (wp *WorkerPool) worker(id int) {
	defer wp.wg.Done()
	for {
		select {
		case job, ok := <-wp.Jobs:
			if !ok {
				return
			}
			result := fmt.Sprintf("Worker %d jobs: %s", id, job)
			time.Sleep(1 * time.Second)
			wp.Result <- result
		case <-wp.ctx.Done():
			return
		}
	}
}

func (wp *WorkerPool) Run() {
	wp.mtx.Lock()
	defer wp.mtx.Unlock()
	for i := 1; i <= wp.workersCount; i++ {
		wp.wg.Add(1)
		go wp.worker(i)
	}
}

func (wp *WorkerPool) AddJob(ctx context.Context, job string) {
	select {
	case wp.Jobs <- job:
	case <-ctx.Done():
		fmt.Println("Context cancel, cant add job")
	}
}

func (wp *WorkerPool) Close() {
	wp.cancel()
	wp.wg.Wait()
	close(wp.Result)
	close(wp.Jobs)
}

func (wp *WorkerPool) AddWorker(num int) {
	wp.mtx.Lock()
	defer wp.mtx.Unlock()
	for i := 1; i <= num; i++ {
		wp.wg.Add(1)
		go wp.worker(wp.workersCount + i)
	}
	wp.workersCount += num
}

func (wp *WorkerPool) DeleteWorker(num int) {
	wp.mtx.Lock()
	defer wp.mtx.Unlock()

	if wp.workersCount < 0 {
		wp.workersCount = 0
	}

	for i := 0; i < num; i++ {
		ctxJob, cancel := context.WithCancel(wp.ctx)
		cancel()
		wp.AddJob(ctxJob, "")
	}
	wp.workersCount -= num
	return
}
