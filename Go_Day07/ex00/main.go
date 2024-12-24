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

func minCoins2(val int, coins []int) []int {
	res := make([]int, 0)
	normCoins, ok := PrepareBag(coins)
	if !ok {
		return []int{}
	}
	i := len(normCoins) - 1
	// fmt.Println(i)
	innerVal := val
	for i >= 0 && innerVal != 0 {
		res = []int{}
		j := i
		// fmt.Println("try:", innerVal)
		innerVal = val
		for j >= 0 {
			for innerVal >= normCoins[j] {
				innerVal -= normCoins[j]
				res = append(res, normCoins[j])
				// fmt.Println("take >", normCoins[j], "remainder: ", innerVal)
				// fmt.Println(res)
			}
			j--
		}
		i--
	}
	// fmt.Println(innerVal)
	if innerVal != 0 {
		return []int{}
	}
	return res
}

func PrepareBag(coins []int) ([]int, bool) {
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
	// fmt.Println(coins)
	// fmt.Println(normCoins)
	return normCoins, true
}

func main() {
	coins := []int{2, 4}
	amount := 3

	actual := minCoins(amount, coins)
	fmt.Println("old ", actual)
	actual = minCoins2(amount, coins)
	fmt.Println("new ", actual)
}
