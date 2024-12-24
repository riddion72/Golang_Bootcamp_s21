package main

import (
	"reflect"
	"testing"
)

func TestMinCoins1(t *testing.T) {
	coins := []int{1, 5, 10}
	amount := 13
	expected := []int{10, 1, 1, 1}
	actual := minCoins2(amount, coins)
	if !reflect.DeepEqual(expected, actual) {
		t.Errorf("Expected %v, got %v", expected, actual)
	}
}

func TestMinCoins2(t *testing.T) {
	tests := []struct {
		name     string
		val      int
		coins    []int
		expected []int
	}{
		{
			name:     "Basic case",
			val:      10,
			coins:    []int{1, 2, 5},
			expected: []int{5, 5},
		},
		{
			name:     "No coins available",
			val:      3,
			coins:    []int{2, 4},
			expected: []int{},
		},
		{
			name:     "Single coin",
			val:      7,
			coins:    []int{7},
			expected: []int{7},
		},
		{
			name:     "Multiple coins with same value",
			val:      12,
			coins:    []int{1, 2, 3, 4, 5},
			expected: []int{5, 5, 2},
		},
		{
			name:     "Large value",
			val:      100,
			coins:    []int{1, 5, 10, 25},
			expected: []int{25, 25, 25, 25},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := minCoins2(tt.val, tt.coins)

			if !reflect.DeepEqual(result, tt.expected) {
				t.Errorf("Test %q: Got %v, Want %v", tt.name, result, tt.expected)
			}
		})
	}
}
