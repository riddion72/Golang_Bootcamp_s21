package main

import (
	"fmt"
	pres "main/presents"
	tree "main/trees"
)

func main() {
	key := make([]uint8, 0, 10)
	key = append(key, 1, 1, 0, 0, 1, 1, 1, 0, 1, 1, 1, 1)
	root := tree.CreateTree(key)
	tree.PrintTree(root, 0, "[root]")
	fmt.Println()
	fmt.Println("============================= Ex00 =============================")
	fmt.Println(tree.AreToysBalanced(root))
	fmt.Println("============================= Ex01 =============================")
	tree.PrintGarland(tree.UnrollGarland(root))

	fmt.Println()
	fmt.Println()

	boxes := []pres.Present{
		{Value: 5, Size: 2},
		{Value: 5, Size: 1},
		{Value: 3, Size: 1},
		{Value: 4, Size: 3},
	}

	fmt.Println("============================= Ex02 =============================")
	fmt.Println()

	result, err := pres.GetNCoolestPresents(boxes, 3)
	if err != nil {
		fmt.Println(err)
		return
	}

	fmt.Println("N Coolest Presents:")
	for _, present := range result {
		fmt.Printf("Value: %d, Size: %d\n", present.Value, present.Size)
	}

	fmt.Println("============================= Ex03_0 =============================")
	fmt.Println()
	answer := pres.GrabPresents(boxes, 3)

	fmt.Println("")
	for _, present := range answer {
		fmt.Printf("Value: %d, Size: %d\n", present.Value, present.Size)
	}

	box := []pres.Present{
		{Value: 11, Size: 1},
		{Value: 11, Size: 1},
		{Value: 11, Size: 1},
		{Value: 11, Size: 1},
		{Value: 20, Size: 2},
		{Value: 20, Size: 2},
	}

	fmt.Println("============================= Ex03_1 =============================")
	ans := pres.GrabPresents(box, 4)

	fmt.Println("")
	for _, present := range ans {
		fmt.Printf("Value: %d, Size: %d\n", present.Value, present.Size)
	}

}
