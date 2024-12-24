package main

import (
	"fmt"
	"math/rand"
	"testing"
)

const (
	NUM_OF_CHANS int = 1000
)

func TestMultiplex(t *testing.T) {
	randChans := make([]chan interface{}, NUM_OF_CHANS)
	testIndexes := make([]interface{}, 0, len(randChans))

	for i := 0; i < len(randChans); i++ {
		testIndexes = append(testIndexes, rand.Intn(100))
	}

	for i, val := range testIndexes {
		randChans[i] = make(chan interface{})
		elem := val
		go func(ch chan interface{}, i interface{}) {
			ch <- elem
			close(ch)
		}(randChans[i], elem)
	}
	res := multiplex(randChans...)

	results := make(map[interface{}]int)

	for val := range res {
		results[val]++
	}

	fmt.Println(results)

	for _, val := range testIndexes {
		fmt.Println(val, results[val])
		if results[val] == 0 {
			t.Fatalf("Test data %v not recieved\n", val)
		} else {
			results[val]--
		}
	}
	fmt.Println(results)
	for _, val := range testIndexes {
		if results[val] != 0 {
			t.Fatalf("Test data %v not recieved\n", val)
		}
	}
}
