package main

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"strconv"
	"sync"
	"syscall"
)

func worker(id int, jobs <-chan string, results chan<- *string, wg *sync.WaitGroup, ctx context.Context) {
	fmt.Println("worker", id, "starting")
	defer wg.Done()
	// time.Sleep(time.Second)
	for {
		select {
		case <-ctx.Done():
			// t := time.NewTimer(time.Second * time.Duration(rand.Intn(10)))
			// <-t.C
			// fmt.Println("worker", id, "cancelled")
			return
		case url, ok := <-jobs:
			if !ok {
				// fmt.Println("worker", id, "finished")
				return
			}
			// fmt.Println("worker", id, "take", url)
			// time.Sleep(time.Second)
			results <- &url
		}
	}
}

func crawlWeb(url chan string, ctx context.Context) chan *string {
	const numWorker = 8
	results := make(chan *string)
	go func() {
		wg := sync.WaitGroup{}
		for w := 1; w <= numWorker; w++ {
			wg.Add(1)
			go worker(w, url, results, &wg, ctx)
		}
		wg.Wait()
		// fmt.Println("wDone")
		close(results)
	}()

	return results
}

func main() {
	const numJobs = 5_000

	jobs := make(chan string)
	cntx, cancel := context.WithCancel(context.Background())

	go func() {
		exit := make(chan os.Signal, 1)
		signal.Notify(exit, os.Interrupt, syscall.SIGTERM)
		<-exit
		cancel()
	}()

	go func() {
		for j := 0; j < numJobs; j++ {
			str := strconv.Itoa(j)
			jobs <- ("examples " + str)
		}
		close(jobs)
	}()

	results := crawlWeb(jobs, cntx)

	i := 0
	for a := range results {
		fmt.Println(*a)
		fmt.Println(i)
		i++
	}
}
