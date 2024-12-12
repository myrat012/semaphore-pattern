package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"runtime"
	"time"

	"golang.org/x/sync/semaphore"
)

type Task struct {
	ID        int    `json:"id"`
	UserID    int    `json:"user_id"`
	Title     string `json:"title"`
	Complated bool   `json:"complated"`
}

type TaskResult struct {
	Task    *Task
	Success bool
	Err     error
}

func main() {
	var t *Task
	sem := semaphore.NewWeighted(10)
	lenth := 100

	// Channel for transmitting the result of task execution
	result := make(chan TaskResult, lenth)

	// Timeout to complete all tasks
	globalCtx, globalCancel := context.WithTimeout(context.Background(), 50*time.Second)
	defer globalCancel()

	for i := 0; i < lenth; i++ {
		go func(i int) {
			// Local context with timeout for each task
			ctx, cancel := context.WithTimeout(globalCtx, 10*time.Second)
			defer cancel()

			// Capturing the semaphore
			if err := sem.Acquire(ctx, 1); err != nil {
				result <- TaskResult{
					Task:    t,
					Success: false,
					Err:     fmt.Errorf("failed to acquire semaphore: %v", err),
				}
				return
			}

			// Release the permission
			defer sem.Release(1)

			// Perform the task
			task, err := processTask(i)
			result <- TaskResult{
				Task:    task,
				Success: err == nil,
				Err:     err,
			}
			fmt.Println(runtime.NumGoroutine())
		}(i)
	}

	// Wait for all tasks to complete and process the results
	for i := 0; i < lenth; i++ {
		result := <-result
		if result.Success {
			log.Printf("Task %v completed successfully\n", result.Task)
		} else {
			log.Printf("Task %v failed: %v\n", result.Task, result.Err)
		}

	}
	fmt.Println("All tasks are processed")
}

func processTask(id int) (*Task, error) {
	var t Task
	url := fmt.Sprintf("https://jsonplaceholder.typicode.com/todos/%d", id)
	res, err := http.Get(url)
	if err != nil {
		log.Fatal(err)
		return nil, err
	}
	defer res.Body.Close()
	if err := json.NewDecoder(res.Body).Decode(&t); err != nil {
		log.Fatal(err)
		return nil, err
	}

	return &t, nil
}
