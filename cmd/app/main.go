package main

import (
	"context"
	"fmt"
	"os/signal"
	"sync"
	"syscall"

	workerpool "github.com/LAshinCHE/vk_worker_pool/internal/worker_pool"
)

func producer(jobs ...string) <-chan string {
	jobChan := make(chan string)
	go func() {
		defer close(jobChan)
		for _, job := range jobs {
			jobChan <- job
		}
	}()
	return jobChan
}

func main() {
	var (
		ctx                   = context.Background()
		ctxWithCancel, cancel = signal.NotifyContext(ctx, syscall.SIGINT)
		syncCh                = make(chan struct{})
	)

	defer cancel()

	wp := workerpool.NewWorkerPool(3, ctxWithCancel)
	wp.Run()
	wp.AddWorker(3)
	wp.DeleteWorker(1)
	// Генерируем канал с заданиями
	jobs := producer("Task1", "Task2", "Task3", "Task4", "Task5", "Task6", "Task7", "Task8", "Task9", "Task0")

	var wg sync.WaitGroup
	wg.Add(1)

	// Собираем данные
	go func() {
		defer wg.Done()
		for result := range wp.Result {
			fmt.Println("Result:", result)
		}
	}()

	// отправляем работу воркеру
	go func(ctx context.Context, syncCh chan<- struct{}) {
		for {
			select {
			case <-ctx.Done():
				fmt.Println("gracefull shutdown")
				syncCh <- struct{}{}
				close(syncCh)
				wp.Close()
				wg.Wait()
				fmt.Println("Worker pool closed.")
				return
			default:
				for job := range jobs {
					wp.AddJob(ctx, job)
				}
			}
		}
	}(ctxWithCancel, syncCh)

	for range syncCh {
		fmt.Println("Все горутины завершили работу")
	}

}
