package main

import (
	"context"
	"fmt"
	"sync"

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
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	wp := workerpool.NewWorkerPool(3, ctx)
	wp.DeleteWorker(1)
	wp.Run()
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
	for job := range jobs {
		wp.AddJob(ctx, job)
	}

	wp.Close()
	wg.Wait()
	fmt.Println("Worker pool closed.")
}
