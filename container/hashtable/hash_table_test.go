package hashtable

import (
	"testing"
	"fmt"
	"water-go/container"
)

func errorString(tip string, expect interface{}, get interface{}) string {
	return fmt.Sprintf("%v expect [%v] get [%v].\n", tip, expect, get)
}

func TestHashTable(t *testing.T) {
	fmt.Println("new...")
	ht := New()
	// init
	fmt.Println("init len...")
	initLen := ht.Len()
	if initLen != 0 {
		t.Errorf(errorString("Init len", 0, initLen))
		t.Failed()
	}
	// put
	fmt.Println("put...")
	for key, value := range container.SearchMap {
		ht.Put(NewIntHashKey(key), value)
	}
	// len
	fmt.Println("full len...")
	fullLen := ht.Len()
	if fullLen != container.SearchMapLength {
		t.Errorf(errorString("Full len", container.SearchMapLength, fullLen))
		t.Failed()
	}
	// get
	fmt.Println("get...")
	for key, value := range container.SearchMap {
		realValue, ok := ht.Get(NewIntHashKey(key))
		if ok != true {
			t.Errorf(errorString(
				fmt.Sprintf("Get [%v]-ok", key),
				true, ok))
		}
		if sv, ok := realValue.(string); !ok || sv != value {
			t.Errorf(errorString(
				fmt.Sprintf("Get [%v]", key),
				value, sv))
		}
	}
	// not key
	fmt.Println("not keys get...")
	for _, key := range container.SearchMapNotKeys {
		_, ok := ht.Get(NewIntHashKey(key))
		if ok != false {
			t.Errorf(errorString(
				fmt.Sprintf("Get [%v]-ok", key),
				false, ok))
		}
	}

	// delete
	fmt.Println("delete...")
	for key, _ := range container.SearchMap {
		ht.Delete(NewIntHashKey(key))
	}
	// empty len
	fmt.Println("empty len...")
	emptyLen := ht.Len()
	if emptyLen != 0 {
		t.Errorf(errorString("Empty len", 0, emptyLen))
		t.Failed()
	}
	// TODO 模拟实际使用，增删改，freeslot
}

func BenchmarkHashTable(b *testing.B) {
	b.ReportAllocs()
	ht := New()
	for i := 0; i < b.N; i++ {
		k := NewIntHashKey(i) // k, i per allocs
		ht.Put(k, i)
		ht.Get(k)
	}
}

func BenchmarkMap(b *testing.B) {
	b.ReportAllocs()
	m := make(map[int]int, b.N)
	for i := 0; i < b.N; i++ {
		m[i] = i
	}
}
