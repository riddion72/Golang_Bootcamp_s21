package bagTask

import (
	"sort"
)

// MinCoins возвращает минимальное набор монет для достижения заданного значения.
//
// Parameters:
//
//	val - целевое значение, для которого нужно найти минимальный набор монет
//	coins - слайс с доступными номиналами монет
//
// Returns:
//
//	[]int - минимальный набор монет (значения монет в порядке их использования)
func MinCoins(val int, coins []int) []int {
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

// MinCoins2 возвращает минимальное набор монет для достижения заданного значения,
// используя динамическое программирование.
//
// Parameters:
//
//	amount - целевое значение, для которого нужно найти минимальный набор монет
//	coins - слайс с доступными номиналами монет
//
// Returns:
//
//	[]int - минимальный набор монет (значения монет в порядке их использования)
//
// Note: Этот вариант использует дополнительное пространство для хранения промежуточных результатов.
func MinCoins2(amount int, coins []int) []int {
	doubleAmount := amount * 2
	normCoins, ok := PrepareBag(coins)
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
		for i := coin; i <= doubleAmount; i++ {
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
	return coinUsed[i]
}

// PrepareBag преобразует список монет в отсортированный массив уникальных значений.
//
// Parameters:
//
//	coins - слайс с номиналами монет
//
// Returns:
//
//	[]int - отсортированный массив уникальных значений монет
//	bool - флаг успешности обработки (false, если массив не валидный)
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
	return normCoins, true
}
