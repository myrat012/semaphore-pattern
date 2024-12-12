package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"runtime"

	"golang.org/x/sync/errgroup"
	"golang.org/x/sync/semaphore"
)

type Task struct {
	ID        int    `json:"id"`
	UserID    int    `json:"user_id"`
	Title     string `json:"title"`
	Complated bool   `json:"complated"`
}

func main() {
	var t Task
	sem := semaphore.NewWeighted(10)
	errGroup, ctx := errgroup.WithContext(context.Background())
	for i := 0; i < 100; i++ {
		if err := sem.Acquire(ctx, 1); err != nil {
			log.Fatal(err)
		}
		i := i
		errGroup.Go(func() error {
			defer sem.Release(1)
			url := fmt.Sprintf("https://jsonplaceholder.typicode.com/todos/%d", i)
			if i == 30 {
				url = ""
			}
			res, err := http.Get(url)
			if err != nil {
				log.Fatal(err)
				return err
			}
			defer res.Body.Close()
			if err := json.NewDecoder(res.Body).Decode(&t); err != nil {
				log.Fatal(err)
				return err
			}
			fmt.Printf("%d. %s - gorutines(%d)\n", i, t.Title, runtime.NumGoroutine())
			return nil
		})
	}
	if err := errGroup.Wait(); err != nil {
		log.Fatal(err)
	}
}
