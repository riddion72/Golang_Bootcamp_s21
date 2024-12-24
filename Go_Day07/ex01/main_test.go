package main

import (
	"fmt"
	"os"
	"sort"
	"testing"
	"time"
)

type coinsTest struct {
	value             int
	coins, exp_output []int
}

type testCaseAttrs struct {
	number           int
	result1, result2 time.Duration
}

type testResult struct {
	result []testCaseAttrs
}

var testCases = []coinsTest{

	{
		value:      10,
		coins:      []int{1, 2, 5},
		exp_output: []int{5, 5},
	},
	{
		value:      3,
		coins:      []int{2, 4},
		exp_output: []int{4},
	},
	{
		value:      3,
		coins:      []int{1, 5, 10},
		exp_output: []int{1, 1, 1},
	},
	{
		value:      7,
		coins:      []int{7},
		exp_output: []int{7},
	},
	{
		value:      12,
		coins:      []int{1, 2, 3, 4, 5},
		exp_output: []int{5, 5, 2},
	},
	{
		value:      100,
		coins:      []int{1, 5, 10, 25},
		exp_output: []int{25, 25, 25, 25},
	},
	{
		value:      10,
		coins:      []int{1, 2, 5},
		exp_output: []int{5, 5},
	},
	{
		value:      3,
		coins:      []int{2, 4},
		exp_output: []int{4},
	},
	{
		value:      3,
		coins:      []int{1, 5, 10},
		exp_output: []int{1, 1, 1},
	},
	{
		value:      7,
		coins:      []int{7},
		exp_output: []int{7},
	},
	{
		value:      12,
		coins:      []int{1, 2, 3, 4, 5},
		exp_output: []int{5, 5, 2},
	},
	{
		value:      100,
		coins:      []int{1, 5, 10, 25},
		exp_output: []int{25, 25, 25, 25},
	},
}

func BenchmarkBestTime(b *testing.B) {
	timesResult := make([]testResult, len(testCases))
	for i, testCase := range testCases {
		timesResult[i].result = append(timesResult[i].result, testCaseAttrs{i, 0, 0})
		b.StartTimer()
		for j := 0; j < b.N; j++ {
			b.ResetTimer()
			start := time.Now()
			minCoins2(testCase.value, testCase.coins)
			timesResult[i].result = append(timesResult[i].result, testCaseAttrs{number: i, result1: time.Duration(time.Since(start).Nanoseconds())})
		}
		b.StopTimer()
	}
	for t := range timesResult {
		sort.Slice(timesResult[t].result, func(i, j int) bool {
			return timesResult[t].result[i].result1 > timesResult[t].result[j].result1
		})
	}
	file, err := os.Create("top10.txt")
	if err != nil {
		b.Fatalf("Error creating a file: %s", err.Error())
	}
	defer file.Close()
	for i := 0; i < 10; i++ {
		fmt.Fprintf(file, "Test case number %d finished with %d ㎲\n", timesResult[i].result[0].number, timesResult[i].result[0].result1)
	}
}

func TestSimplestNTimes(t *testing.T) {
	timesResult := make([]testResult, len(testCases))
	for i, testCase := range testCases {
		timesResult[i].result = append(timesResult[i].result, testCaseAttrs{i, 0, 0})
		// b.StartTimer()
		for j := 0; j < 10; j++ {
			// b.ResetTimer()
			start := time.Now()
			minCoins2(testCase.value, testCase.coins)
			res1Time := time.Duration(time.Since(start).Nanoseconds())
			start = time.Now()
			minCoins(testCase.value, testCase.coins)
			res2Time := time.Duration(time.Since(start).Nanoseconds())
			timesResult[i].result = append(timesResult[i].result, testCaseAttrs{number: i, result1: res1Time, result2: res2Time})
		}
		// b.StopTimer()
	}
	for t := range timesResult {
		sort.Slice(timesResult[t].result, func(i, j int) bool {
			return timesResult[t].result[i].result1-timesResult[t].result[i].result2 > timesResult[t].result[j].result1-timesResult[t].result[j].result1
		})
	}
	file, err := os.Create("compare_times.txt")
	if err != nil {
		t.Fatalf("Error creating a file: %s", err.Error())
	}
	defer file.Close()
	for i, testTime := range timesResult {
		fmt.Fprintf(file, "Test case number %d minCoins2 finished with %d ㎲, minCoins with %d ㎲. The difference is: %d\n", testTime.result[0].number, testTime.result[0].result1, testTime.result[0].result2, testTime.result[0].result1-testTime.result[0].result2)
		if i == 10 {
			break
		}
	}
}
