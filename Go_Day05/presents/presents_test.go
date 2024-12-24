package presents

import (
	"testing"
)

func TestGetNCoolestPresents(t *testing.T) {
	presents := PresentsHeap{
		{Value: 5, Size: 1},
		{Value: 4, Size: 5},
		{Value: 3, Size: 1},
		{Value: 5, Size: 2},
	}

	capacity := 0
	coolestPresents, err := GetNCoolestPresents(presents, capacity)
	if err == nil {
		t.Error(err)
	}
	capacity = 30
	coolestPresents, err = GetNCoolestPresents(presents, capacity)
	if err == nil {
		t.Error(err)
	}
	capacity = 2
	coolestPresents, err = GetNCoolestPresents(presents, capacity)
	if len(coolestPresents) != 2 {
		t.Error("Expected 2 coolest presents, got", len(coolestPresents))
	}
	if coolestPresents[0].Value != 5 || coolestPresents[0].Size != 1 {
		t.Error("Expected first")
	}
}

func TestGrabPresents1(t *testing.T) {
	presents := PresentsHeap{
		{Value: 5, Size: 1},
		{Value: 4, Size: 5},
		{Value: 3, Size: 1},
		{Value: 5, Size: 2},
	}

	capacity := 0
	grabbedPresents := GrabPresents(presents, capacity)
	if len(grabbedPresents) != 0 {
		t.Error("Expected no presents, got", len(grabbedPresents))
	}
	capacity = 30
	grabbedPresents = GrabPresents(presents, capacity)
	if len(grabbedPresents) != 4 {
		t.Error("Expected no presents, got", len(grabbedPresents))
	}
	capacity = 2
	grabbedPresents = GrabPresents(presents, capacity)
	if len(grabbedPresents) != 2 {
		t.Error("Expected 2 presents, got", len(grabbedPresents))
	}
	if grabbedPresents[1].Value != 3 || grabbedPresents[1].Size != 1 {
		t.Error("Expected first")
	}
}

func TestGrabPresents2(t *testing.T) {
	presents := PresentsHeap{
		{Value: 11, Size: 1},
		{Value: 11, Size: 1},
		{Value: 11, Size: 1},
		{Value: 11, Size: 1},
		{Value: 20, Size: 2},
		{Value: 20, Size: 2},
	}

	capacity := 0
	grabbedPresents := GrabPresents(presents, capacity)
	if len(grabbedPresents) != 0 {
		t.Error("Expected no presents, got", len(grabbedPresents))
	}
	capacity = 30
	grabbedPresents = GrabPresents(presents, capacity)
	if len(grabbedPresents) != 6 {
		t.Error("Expected no presents, got", len(grabbedPresents))
	}
	capacity = 4
	grabbedPresents = GrabPresents(presents, capacity)
	if len(grabbedPresents) != 4 {
		t.Error("Expected 2 presents, got", len(grabbedPresents))
	}
	if grabbedPresents[3].Value != 11 || grabbedPresents[3].Size != 1 {
		t.Error("Expected first")
	}
}
