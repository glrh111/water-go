package hashtable

import (
	"testing"
	"fmt"
)

func TestIntHashKey_HashCode(t *testing.T) {
	intList := []int{1, 2, 3, 4, 5, 0, -1}
	for _, value := range intList {
		key := NewIntHashKey(value)
		if value != key.HashCode() {
			t.Failed()
		}
	}
}

func TestStringHashKey_HashCode(t *testing.T) {
	stringList := []string{"namea", "nameb", "namec", "named"}
	for _, value := range stringList {
		key := NewStringHashKey(value)
		t.Logf("Key: [%v], hash: [%v]\n", value, key.HashCode())
	}
}

func TestMask(t *testing.T) {
	mask := 7
	for i := -20; i < 20; i++ {
		t.Logf("mask: %v, hash: %v, slot: %v\n", mask, i, mask & i)
	}
}

func TestSome(t *testing.T) {
	a := [3]int{1,2,3}
	b := a[:]
	fmt.Println(a, b)
	b[2] = 555
	fmt.Println(a, b)
	c := make([]entry, 5)
	for _, item := range c {
		fmt.Printf("%+v\n", item)
	}
}
