package hashtable

import "fmt"

// 开放寻址法，实现hashtab

// 对象复用池 TODO
const (
	MaxFreeList = 80
	SmallTableSize = 1 << 3
	PerturbShift = 5
	LargeTableSize = 500000
	SmallScale = 4
	LargeScale = 2
)
var (
	freeList [MaxFreeList]*HashTable
	numFree = 0  // free 元素的指针
)

type entry struct {
	hashcode int
	key HashKeyer
	value interface{}
}

func newEntry() entry {
	return entry{-1, nil, nil}
}

func (et entry) String() string {
	return fmt.Sprintf("Entry: (%v, %v)", et.key, et.value)
}

/*
   used: active
   fill: active + dummy
   mask: len(table) size - 1
*/
type HashTable struct {
	used, fill, mask int
	table []entry  // 默认指向 smallTable &table[0] == &smallTable[0]
	smallTable [SmallTableSize]entry
	lookup func(*HashTable, HashKeyer) *entry // 寻找逻辑
}

func New() *HashTable {
	// 加入复用机制 销毁的时候，将其放入freelist中
	ht := &HashTable{
		used: 0,
		fill: 0,
		mask: SmallTableSize - 1,
	}
	st := [SmallTableSize]entry{}
	//for i := 0; i < SmallTableSize; i++ {
	//	st[i] = newEntry()
	//}
	ht.smallTable = st
	ht.table = st[:]
	ht.lookup = lookup
	return ht
}

// value could not be nil
func (ht *HashTable) Get(key HashKeyer) (interface{}, bool) {
	et := ht.lookup(ht, key)
	if et != nil && et.value != nil { // Active
		return et.value, true
	}
	return nil, false
}

func (ht *HashTable) Len() int { return ht.used }

// TODO Put will do 2 mallocs, 1. int, string -> HashKeyer 2. value -> interface{}
func (ht *HashTable) Put(key HashKeyer, value interface{}) {
	et := ht.lookup(ht, key)
	if et.value != nil { // Active
		et.value = value
	} else if et.key != nil { // Dummy
		et.key = key
		et.value = value
		ht.used++
	} else { // Unused ---
		et.key = key
		et.value = value
		ht.used++
		ht.fill++
		ht.resize()  // only insert deal with this
	}
}

func (ht *HashTable) Delete(key HashKeyer) {
	//fmt.Println("In delete: ", key)
	et := ht.lookup(ht, key)
	if et.value != nil { // Active -> Dummy
		et.value = nil
		ht.used -= 1
	}
}

func (ht *HashTable) String() string {
	return fmt.Sprintf("HashTable: mask[%v] used[%v] fill[%v]", ht.mask, ht.used, ht.fill)
}

func (ht *HashTable) resize() {
	var (
		oldSize = ht.mask + 1
		newSize = SmallTableSize
		expectSize = 0
		newTable []entry
		oldTable []entry
	)
	if !(3 * ht.fill >= 2 * (ht.mask + 1)) {
		return
	}
	//fmt.Println("Resize: ", ht.String())
	if ht.used > LargeTableSize {
		expectSize = LargeScale * ht.used
	} else {
		expectSize = SmallScale * ht.used
	}
	// find newSize TODO find out why could not equal
	for ; newSize < expectSize; newSize <<= 1 { }
	// newSize == SmallTableSize: Do clean dummy
	ht.used = 0
	ht.fill = 0
	ht.mask = newSize - 1
	if newSize == SmallTableSize {
		oldTable = make([]entry, oldSize) // old
		copy(oldTable, ht.table)
		ht.table = ht.smallTable[:]
		// fill zero
		for i := 0; i < SmallTableSize; i++ {
			ht.table[i].key = nil
			ht.table[i].value = nil
		}
	} else {
		oldTable = ht.table
		newTable = make([]entry, newSize)
		ht.table = newTable
	}

	// 搬迁
	for _, et := range oldTable {
		if et.value != nil {  // Active
			ht.Put(et.key, et.value)
		}
	}
}

// find exact entry 2 allocs/op 为什么
func lookup(ht *HashTable, key HashKeyer) *entry {
	var (
		i uint
		freeslot *entry
		startEntry *entry
		hash = uint(key.HashCode())
	)
	i = hash & uint(ht.mask)
	startEntry = &ht.table[i]
	// find: Unused || shot key(Active)
	if startEntry.key == nil || startEntry.key.IsEqual(key) {
		return startEntry
	}
	// Dummy / not key(Active)
	if startEntry.value == nil {  // Dummy(Could use)
		freeslot = startEntry
	} else {                      // not key(Active, could not use)
		freeslot = nil
	}
	// loop
	for perturb := hash; ; perturb >>= PerturbShift {
		i = (i << 2) + i + perturb + 1
		startEntry = &ht.table[i & uint(ht.mask)]
		if startEntry.key == nil { // find Unused
			if freeslot != nil {
				return freeslot
			} else {
				return startEntry
			}
		}
		if startEntry.key.IsEqual(key) { // find exact(Active and Key Equal)
			return startEntry
		}
		if startEntry.value == nil && freeslot == nil {  // dummy
			freeslot = startEntry
		}
	}
}
