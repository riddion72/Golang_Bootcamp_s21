package trees

import (
	"fmt"
	"strings"
)

type TreeNode struct {
	HasToy bool
	Left   *TreeNode
	Right  *TreeNode
}

func getHeight(root *TreeNode) int {
	if root == nil {
		return 0
	}
	HasToy := 0
	if root.HasToy {
		HasToy = 1
	}
	HasToy += getHeight(root.Left) + getHeight(root.Right)
	return HasToy
}

func abs(x int) int {
	if x < 0 {
		return -x
	}
	return x
}

//========================== Ex00 ============================

func AreToysBalanced(root *TreeNode) bool {
	if root == nil {
		return true
	}

	return getHeight(root.Left) == getHeight(root.Right)
}

//========================== Ex01 ============================

type QuNode struct {
	node *TreeNode
	next *QuNode
}

type Queue struct {
	front, back *QuNode
	len         uint
}

func newQueue() *Queue {
	return &Queue{}
}

func (q *Queue) push(node *TreeNode) {
	newNode := &QuNode{node: node}
	if q.len == 0 {
		q.front = newNode
	} else {
		q.back.next = newNode
	}
	q.back = newNode
	q.len++
}

func (q *Queue) pop() *TreeNode {
	if q.len == 0 {
		return nil
	}
	node := q.front.node
	q.front = q.front.next
	if q.front == nil {
		q.back = nil
	}
	q.len--
	return node
}

func reverse(arr []bool) {
	start, end := 0, len(arr)-1
	for start < end {
		arr[start], arr[end] = arr[end], arr[start]
		start++
		end--
	}
}

func UnrollGarland(root *TreeNode) []bool {
	if root == nil {
		return nil
	}
	answer := make([]bool, 0, 100)
	queue := newQueue()

	queue.push(root)
	level := 0
	size := 1
	fromLevel := make([]bool, 0, size)
	nextLevelSize := 0
	for queue.len > 0 {
		for i := 0; i < size && queue.len != 0; i++ {
			curNode := queue.pop()
			if curNode.Left != nil {
				queue.push(curNode.Left)
				nextLevelSize++
			}
			if curNode.Right != nil {
				queue.push(curNode.Right)
				nextLevelSize++
			}
			fromLevel = append(fromLevel, curNode.HasToy)
		}

		if level%2 == 0 {
			reverse(fromLevel)
		}
		answer = append(answer, fromLevel...)
		fromLevel = fromLevel[:0]
		level++

		size = nextLevelSize
		nextLevelSize = 0
	}

	return answer
}

func PrintGarland(answer []bool) {
	for _, v := range answer {
		if v {
			fmt.Print("1 ")
		} else {
			fmt.Print("0 ")
		}

	}
	fmt.Println()
}

//========================== Create Tree ============================

func CreateTree(nodes []uint8) *TreeNode {
	if len(nodes) == 0 {
		return nil
	}
	root := make([]*TreeNode, len(nodes))
	for i, n := range nodes {
		if n != 2 {
			root[i] = &TreeNode{HasToy: n == 1}
		}
	}
	for i := 0; i < len(root)/2; i++ {
		if root[i] != nil {
			if 2*i+1 < len(root) {
				root[i].Left = root[2*i+1]
			}
			if 2*i+2 < len(root) {
				root[i].Right = root[2*i+2]
			}
		}
	}
	return root[0]
}

//========================== Print Tree ============================

func PrintTree(node *TreeNode, depth int, ch string) {
	if node != nil {
		PrintTree(node.Right, depth+1, "/[R]")
		printNode(node, depth, ch)
		PrintTree(node.Left, depth+1, "\\[L]")
	}
}

func printNode(node *TreeNode, depth int, ch string) {
	var val uint8
	if node.HasToy {
		val = 1
	}
	fmt.Printf("%s%s=%v\n", strings.Repeat(" ", depth*4), ch, val)
}
