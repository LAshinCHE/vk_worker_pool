package main

import (
	"fmt"
	"sync"
)

func producer(strs ...string) <-chan string {
	prod := make(chan string)

	go func() {
		for _, str := range strs {
			prod <- str
		}
		close(prod)
	}()

	return prod
}

func worker(id int, wg *sync.WaitGroup, chanData <-chan string, result chan<- struct{}) {
	defer wg.Done()
	for d := range chanData {
		fmt.Printf("Worker id:%d message: %s\n", id, d)
		result <- struct{}{}
	}
}

const (
	workerAmount = 3
)

func main() {
	wg := &sync.WaitGroup{}
	resultSig := make(chan struct{})

	wg.Add(workerAmount)

	prod := producer("some data 4", "some data 3", "some data 2", "some data 1", "some data 5")

	for i := 1; i <= workerAmount; i++ {
		go worker(i, wg, prod, resultSig)
	}
	go func() {
		wg.Wait()
		close(resultSig)
	}()
	for range resultSig {
	}
}
