package main

import (
	"fmt"
	"sync"
	"time"
)

func main() {
	maxGorutines := 5
	semaphore := make(chan struct{}, maxGorutines)

	var wg sync.WaitGroup

	for i := 0; i < 10; i++ {
		wg.Add(1)
		go func(i int) {
			defer wg.Done()
			semaphore <- struct{}{}
			defer func() { <-semaphore }()

			// simulate work
			fmt.Println("Running task", i)
			time.Sleep(5 * time.Second)
		}(i)
	}
	wg.Wait()
}
