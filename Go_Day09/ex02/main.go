package main

import (
	"fmt"
	"sync"
)

func main() {
	chInt1 := make(chan interface{})
	chInt2 := make(chan interface{})
	chInt3 := make(chan interface{})
	go func() {
		for i := 0; i < 10; i++ {
			chInt1 <- "i"
			chInt2 <- i * 10
			chInt3 <- float32(i) * 0.1
		}
		close(chInt1)
		close(chInt2)
		close(chInt3)
	}()
	multiplexed := multiplex(chInt1, chInt2, chInt3)
	for data := range multiplexed {
		fmt.Println(data)
	}
}

func multiplex(chanIn ...chan interface{}) chan interface{} {
	multiplexed := make(chan interface{})

	var wg sync.WaitGroup

	wg.Add(len(chanIn))
	go func() {
		for i, ch := range chanIn {
			go func(lCh <-chan interface{}, j int) {
				fmt.Println("start", j)
				for n := range lCh {
					fmt.Println(j, " take", n)
					multiplexed <- n
				}
				wg.Done()
			}(ch, i)
		}
		wg.Wait()
		close(multiplexed)
	}()
	return multiplexed
}
