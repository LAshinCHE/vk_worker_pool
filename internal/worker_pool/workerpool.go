package workerpool

import (
	"context"
	"sync"
)

type Job interface {
	execute(string) error
}

type WorkerPool struct {
	workersCount int
	jobs         chan Job
	result       chan string
	ctx          context.Context
	cancel       context.CancelFunc
	wg           sync.WaitGroup
}

func NewWorkerPool(workerNum int, ctx context.Context) WorkerPool {
	ctx, cancel := context.WithCancel(ctx)

	return WorkerPool{
		workersCount: workerNum,
		jobs:         make(chan Job),
		result:       make(chan string),
		ctx:          ctx,
		cancel:       cancel,
	}
}

// Todo: сделать функцию которая делает примитивную работу
// Читает сообщение из канала jobs и выводит его на экран и выводит номер нашего воркера
// так как wp рабоает запускается с одного канала нужно будет реализовать паттерн fanOut для записи задач по всем воркерам
// считывание будет использоваться при помощи fanIn будет читать данные с воркеров и записывать их в один общий канал result
func (wp WorkerPool) worker(id int) error {
	return nil
}

// Todo:сделать функцию которая будет добавлять нового воркера в наш
func (wp WorkerPool) AddWorker(numWorkersToAdd int) error {
	return nil
}

// Todo: Сделать функцию удалеия числа воркеров
func (wp WorkerPool) DeleteWorker(numWorkersToAdd int) error {
	return nil
}
