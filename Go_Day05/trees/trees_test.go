package trees

import (
	"reflect"
	"testing"
)

func TestCreateTree(t *testing.T) {
	nodes := []uint8{1, 1, 0, 0, 1, 1, 1, 0, 1, 1, 1, 1}
	root := CreateTree(nodes)

	PrintTree(root, 0, "[root]")

	if root.Left.Left.HasToy == true ||
		root.Left.Right.HasToy == false ||
		root.Right.Left.HasToy == false ||
		root.Right.Right.HasToy == false {
		t.Error("Tree creation failed")
	}
}

func TestUnrollGarland(t *testing.T) {
	nodes := []uint8{1, 1, 0, 0, 1, 1, 1, 0, 1, 1, 1, 1}
	root := CreateTree(nodes)

	answer := UnrollGarland(root)
	// fmt.Println(answer)

	expected := []bool{true, true, false, true, true, true, false, false, true, true, true, true}
	if !reflect.DeepEqual(answer, expected) {
		t.Error("Garland unrolling failed")
	}
}

func TestAreToysBalanced(t *testing.T) {
	nodes := []uint8{1, 1, 0, 0, 1, 1, 1, 0, 1, 1, 1, 1}
	root := CreateTree(nodes)

	if AreToysBalanced(root) {
		t.Error("Toys are not balanced")
	}
}

func areEqual(a, b []bool) bool {
	if len(a) != len(b) {
		return false
	}
	for i := range a {
		if a[i] != b[i] {
			return false
		}
	}
	return true
}
