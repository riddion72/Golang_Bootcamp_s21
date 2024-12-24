package presents

import (
	"container/heap"
	"fmt"
)

type Present struct {
	Value int
	Size  int
}

type PresentsHeap []Present

func (ph PresentsHeap) Len() int { return len(ph) }

func (ph PresentsHeap) Less(i, j int) bool {
	if ph[i].Value == ph[j].Value {
		return ph[i].Size < ph[j].Size
	}
	return ph[i].Value > ph[j].Value
}

func (ph PresentsHeap) Swap(i, j int) {
	ph[i], ph[j] = ph[j], ph[i]
}

func (ph *PresentsHeap) Push(x any) {
	item := x.(Present)
	*ph = append(*ph, item)
}

func (ph *PresentsHeap) Pop() any {
	old := *ph
	n := len(old)
	pres := old[n-1]
	*ph = old[0 : n-1]
	return pres
}

func CreateHeapFromSlise(pres []Present) *PresentsHeap {
	h := &PresentsHeap{}
	for _, el := range pres {
		heap.Push(h, el)
	}
	heap.Init(h)
	return h
}

//============================== Ex02 ==============================

func GetNCoolestPresents(boxes []Present, size int) ([]Present, error) {
	// fmt.Println(len(boxes))
	if len(boxes) < size || size < 1 {
		return nil, fmt.Errorf("invalid number of presents")
	}
	prHeap := CreateHeapFromSlise(boxes)
	result := make([]Present, size)
	for i := 0; i < size; i++ {
		result[i] = heap.Pop(prHeap).(Present)
	}
	return result, nil
}

//============================== Ex03 ==============================

func GrabPresents(boxes []Present, size int) []Present {
	pakValue := make([][]int, (len(boxes) + 1))
	pakValueCap := make([]int, (size+1)*(len(boxes)+1))
	boxesL := make([]Present, 1, (len(boxes) + 1))
	boxAns := make([]Present, 0, (len(boxes)))
	boxesL = append(boxesL, boxes...)
	for i := range pakValue {
		pakValue[i], pakValueCap = pakValueCap[:(size+1)], pakValueCap[(size+1):]
	}
	for i := 1; i < len(boxesL); i++ {
		for j := 1; j <= size; j++ {
			if boxesL[i].Size > j {
				pakValue[i][j] = pakValue[i-1][j]
			} else if pakValue[i-1][j] < pakValue[i-1][j-boxesL[i].Size]+boxesL[i].Value {
				pakValue[i][j] = pakValue[i-1][j-boxesL[i].Size] + boxesL[i].Value
			} else {
				pakValue[i][j] = pakValue[i-1][j]
			}
		}

	}

	findAns(pakValue, len(boxes), size, boxesL, &boxAns)

	// for _, line := range pakValue {
	// 	fmt.Println(line)
	// }
	// fmt.Println(boxAns)
	return boxAns
}

func findAns(A [][]int, k int, s int, box []Present, ans *[]Present) {
	if A[k][s] == 0 {
		return
	}
	if A[k-1][s] == A[k][s] {
		findAns(A, k-1, s, box, ans)
	} else {
		findAns(A, k-1, s-box[k].Size, box, ans)
		*ans = append(*ans, box[k])
	}
}

// func GrabPresents(boxes []Present, size int) []Present {
// 	h := CreateHeapFromSlise(boxes)
// 	result := make([]Present, 0, 1)
// 	var weight int = 0
// 	for weight < size && h.Len() > 0 {
// 		pr := heap.Pop(h).(Present)
// 		if weight+int(pr.Size) <= size {
// 			result = append(result, pr)
// 			weight += int(pr.Size)
// 		}
// 	}
// 	return result
// }
