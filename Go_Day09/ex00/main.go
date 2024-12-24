package main

import (
	"fmt"
	"sync"
	"time"
)

func main() {
	slices := []int{10, 2, 30, 4, 5}
	ch := sleepSort(slices)
	for num := range ch {
		fmt.Println(num)
	}
}

func sleepSort(inPut []int) chan int {
	output := make(chan int, len(inPut))

	wg := sync.WaitGroup{}
	for _, num := range inPut {
		wg.Add(1)
		go func(num int) {
			time.Sleep(time.Duration(num) * time.Millisecond)
			output <- num
			wg.Done()
		}(num)
	}
	wg.Wait()
	close(output)

	return output
}
