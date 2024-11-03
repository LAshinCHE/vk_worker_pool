package workerpool

import (
	"context"
	"fmt"
	"sync"
)

// type Job interface {
// 	execute() string
// }

type WorkerPool struct {
	workersCount int
	jobs         chan string
	result       chan string
	wg           *sync.WaitGroup
}

func NewWorkerPool(workerNum int) WorkerPool {

	return WorkerPool{
		workersCount: workerNum,
		jobs:         make(chan string),
		result:       make(chan string),
	}
}

// Todo: сделать функцию которая делает примитивную работу
// Читает сообщение из канала jobs и выводит его на экран и выводит номер нашего воркера
// так как wp рабоает запускается с одного канала нужно будет реализовать паттерн fanOut для записи задач по всем воркерам
// считывание будет использоваться при помощи fanIn будет читать данные с воркеров и записывать их в один общий канал result
func (wp WorkerPool) worker(id int, ctx context.Context) {
	defer wp.wg.Done() // завершаем работу нашего воркера

	for {
		select {
		case job, ok := <-wp.jobs:
			if !ok {
				return
			}
			fmt.Printf("Worker id: %d message: %s", id, job)
		case <-ctx.Done():
			fmt.Println("Context finish")
			return
		}
	}
	return
}

func (wp WorkerPool) Run(ctx context.Context) {
	// for i := 0; i < wp.workersCount; i++{

	// }

}

// Todo:сделать функцию которая будет добавлять нового воркера в наш
func (wp WorkerPool) AddWorker(numWorkersToAdd int) {
	return
}

// Todo: Сделать функцию удалеия числа воркеров
func (wp WorkerPool) DeleteWorker(numWorkersToAdd int) {
	return
}
