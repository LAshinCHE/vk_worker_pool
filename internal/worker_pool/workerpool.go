package workerpool

import (
	"context"
	"fmt"
	"sync"
	"sync/atomic"
	"time"
)

type WorkerPool struct {
	workersCount int
	Jobs         chan Job
	Result       chan string
	wg           *sync.WaitGroup
	ctx          context.Context
	cancel       context.CancelFunc
	mtx          sync.Mutex
}

type Job struct {
	ctx  context.Context
	task string
}

var (
	numJobs = atomic.Int64{}
)

func NewWorkerPool(workerNum int, ctx context.Context) *WorkerPool {
	ctx, cancel := context.WithCancel(ctx)
	return &WorkerPool{
		workersCount: workerNum,
		Jobs:         make(chan Job, 0),
		Result:       make(chan string),
		wg:           &sync.WaitGroup{},
		ctx:          ctx,
		cancel:       cancel,
	}
}

// worker читает работу с канала и выполняет ее, если приходит фиктивная работа то завершает воркер
func (wp *WorkerPool) worker(id int) {
	defer wp.wg.Done()
	for {
		job, ok := <-wp.Jobs
		if !ok {
			if numJobs.Load() == 0 {
				return
			}
		}
		select {
		case <-job.ctx.Done():
			return
		default:
		}
		result := fmt.Sprintf("Worker %d jobs: %s", id, job.task)
		wp.Result <- result
		time.Sleep(time.Second * 2)
		numJobs.Add(-1)
	}
}

// Run() Запуское всех воркеров
func (wp *WorkerPool) Run() {
	wp.mtx.Lock()
	defer wp.mtx.Unlock()
	for i := 1; i <= wp.workersCount; i++ {
		wp.wg.Add(1)
		go wp.worker(i)
	}
	return
}

// Добавляет работу в канал работ
func (wp *WorkerPool) AddJob(ctx context.Context, task string) {
	numJobs.Add(1)
	job := Job{
		ctx:  ctx,
		task: task,
	}

	wp.Jobs <- job

}

func (wp *WorkerPool) Close() {
	wp.cancel()
	wp.wg.Wait()
	close(wp.Result)
}

// AddWorker - Запускает новых worker
func (wp *WorkerPool) AddWorker(num int) {
	wp.mtx.Lock()
	defer wp.mtx.Unlock()
	for i := 1; i <= num; i++ {
		wp.wg.Add(1)
		go wp.worker(wp.workersCount + i)
	}
	wp.workersCount += num
}

// DeleteWorker - Удаляет воркера
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
}
