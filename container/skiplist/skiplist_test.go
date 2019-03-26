package skiplist

import (
	"testing"
	"water-go/container"
	"fmt"
	"math/rand"
)

var (
	sl = NewSkiplist()
)

func init() {

}

func errorString(tip string, expect interface{}, get interface{}) string {
	return fmt.Sprintf("%v expect [%v] get [%v].\n", tip, expect, get)
}

func TestSkiplist(t *testing.T) {
	sl := NewSkiplist()
	// init length
	initLen := sl.Len()
	if initLen != 0 {
		t.Errorf(errorString("Init length", 0, initLen))
	}
	// empty
	initEmpty := sl.IsEmpty()
	if initEmpty != true {
		t.Errorf(errorString("Init list", true, initEmpty))
	}
	// insert
	for key, value := range container.SearchMap {
		sl.Put(IntKey(key), value)
	}
	// length
	insertLen := sl.Len()
	if insertLen != container.SearchMapLength {
		t.Errorf(errorString("Insert length", container.SearchMapLength, insertLen))
	}
	// empty
	insertEmpty := sl.IsEmpty()
	if insertEmpty != false {
		t.Errorf(errorString("Insert empty", false, insertEmpty))
	}
	// value
	for key, value := range container.SearchMap {
		slValue, ok := sl.Get(IntKey(key))
		if ok != true {
			t.Errorf(errorString(fmt.Sprintf("Get [%v]-ok", key), true, ok))
		}
		if slValue.(string) != value {
			t.Errorf(errorString(fmt.Sprintf("Get [%v]", key), value, slValue))
		}
	}
	// not in value
	for _, key := range container.SearchMapNotKeys {
		slValue, ok := sl.Get(IntKey(key))
		if ok != false {
			t.Errorf(errorString(fmt.Sprintf("Get [%v]-ok", key), false, ok))
		}
		if slValue != nil {
			t.Errorf(errorString(fmt.Sprintf("Get [%v]", key), nil, slValue))
		}
	}
	// show
	fmt.Println(sl.String())
	// delete
	for key, _ := range container.SearchMap {
		delRe := sl.Delete(IntKey(key))
		if delRe != DeleteYes {
			t.Errorf(errorString(fmt.Sprintf("Delete [%v]-ok", key), DeleteYes, delRe))
		}
	}
	// delete not keys
	for _, key := range container.SearchMapNotKeys {
		delRe := sl.Delete(IntKey(key))
		if delRe != DeleteNo {
			t.Errorf(errorString(fmt.Sprintf("Delete [%v]-ok", key), DeleteNo, delRe))
		}
	}
	// delete size
	deleteLen := sl.Len()
	if deleteLen != 0 {
		t.Errorf(errorString("Delete length", 0, initLen))
	}
	// delete empty
	deleteEmpty := sl.IsEmpty()
	if deleteEmpty != true {
		t.Errorf(errorString("Delete list", true, initEmpty))
	}
	// insert many
	randKey := rand.Perm(10000) // 10W

	for index, key := range randKey {
		sl.Put(IntKey(key), index)
	}
	fmt.Println(sl.String())
}

func BenchmarkSkiplist_Put(b *testing.B) {
	b.ReportAllocs()

	// data
	intSlice := rand.Perm(10000)  // 1M

	// put
	b.Log("Puting...")
	for _, key := range intSlice {
		sl.Put(IntKey(key), 1)
	}

	// get
	b.Logf("Geting... b.N = [%v]\n", b.N)
	for i := 0; i < b.N; i++ {
		sl.Get(IntKey(i))
	}
}
