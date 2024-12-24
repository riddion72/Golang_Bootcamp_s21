package main

import (
	"fmt"
	"unsafe"
)

func main() {
	arr := []int{5, 4, 3, 2, 1}
	idx := 3

	element, err := getElement(arr, idx)
	if err != nil {
		fmt.Println(err)
		return
	}

	fmt.Println("Element at index", idx, "is:", element)
}

func getElement(arr []int, idx int) (int, error) {
	if idx < 0 || idx >= len(arr) {
		return 0, fmt.Errorf("Index out of range")
	}
	if len(arr) == 0 || cap(arr) == 0 {
		return 0, fmt.Errorf("Array is empty")
	}
	const size = unsafe.Sizeof(int(0))
	ans := unsafe.Add(unsafe.Pointer(&arr[0]), idx*int(size))
	return *(*int)(ans), nil
}
