package main

import (
	"fmt"
	"sort"
)

func minCoins(val int, coins []int) []int {
	res := make([]int, 0)
	i := len(coins) - 1
	for i >= 0 {
		for val >= coins[i] {
			val -= coins[i]
			res = append(res, coins[i])
		}
		i -= 1
	}
	return res
}

func minCoins2(amount int, coins []int) []int {
	doubleAmount := amount * 2
	normCoins, ok := prepareBag(coins)
	if !ok {
		return []int{}
	}
	dp := make([]int, doubleAmount+1)
	for i := range dp {
		dp[i] = amount + 1
	}
	dp[0] = 0

	coinUsed := make([][]int, doubleAmount+1)
	for i := range coinUsed {
		coinUsed[i] = []int{}
	}

	for _, coin := range normCoins {
		// fmt.Println("coin: ", coin)

		for i := coin; i <= doubleAmount; i++ {
			// fmt.Println("i: ", i, "dp ", dp, "coinUsed ", coinUsed)
			if dp[i-coin]+1 < dp[i] {
				dp[i] = dp[i-coin] + 1
				coinUsed[i] = make([]int, len(coinUsed[i-coin]))
				copy(coinUsed[i], coinUsed[i-coin])
				coinUsed[i] = append(coinUsed[i], coin)
			}
		}
	}
	i := amount
	for ; (i < doubleAmount) && (dp[i] > amount); i++ {
	}

	if i == doubleAmount {
		return []int{}
	}
	// fmt.Println("answer: ", coinUsed)

	return coinUsed[i]
}

func prepareBag(coins []int) ([]int, bool) {
	if len(coins) == 0 {
		return []int{}, false
	}
	sort.Slice(coins, func(i, j int) bool {
		return coins[i] < coins[j]
	})
	normCoins := []int{}
	previus := 0
	for i, c := range coins {
		if coins[i] < 1 {
			return []int{}, false
		}
		if c != previus {
			normCoins = append(normCoins, c)
		}
		previus = c
	}
	return normCoins, true
}

func main() {
	coins := []int{3, 5, 10}
	amount := 13

	fmt.Println(minCoins2(amount, coins))

	// actual := minCoins2(amount, coins)
	// fmt.Println(actual)
	// fmt.Println(minCoins(13, []int{1, 5, 10}))
}
