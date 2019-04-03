package skiplist

import (
	"math/rand"
	"fmt"
	"time"
)

// TODO 实现有问题，慢得一笔

const (
	CompareBigger = 1
	CompareEqual = 0
	CompareLess = -1
	PutUpdate = 1
	PutNew = 2
	DeleteYes = 1
	DeleteNo = 2
)

type Comparer interface {
	CompareTo(Comparer) int
}

type Skiplister interface {
	Get(Comparer) (interface{}, bool)  // value, ok
	Put(Comparer, interface{}) int  // update, new
	Delete(Comparer) int
	Len() int
	MaxLevel() int
	IsEmpty() bool
	fmt.Stringer
}

type IntKey int

func (i IntKey) CompareTo(j Comparer) int {
	if k, ok := j.(IntKey); ok {
		if i > k {
			return CompareBigger
		} else if i == k {
			return CompareEqual
		} else {
			return CompareLess
		}
	} else {
		panic("j is not IntKey")
	}
}

type StringKey string

func (s StringKey) CompareTo(s2 Comparer) int {
	if k, ok := s2.(StringKey); ok {
		if s > k {
			return CompareBigger
		} else if s == k {
			return CompareEqual
		} else {
			return CompareLess
		}
	} else {
		panic("j is not IntKey")
	}
}

type node struct {
	next, down *node
	isGuard bool  // every linked list's first element is guard
	level int  // start from 0
	size int   // how many nodes after this node
	key Comparer
	value interface{}
}

// 将 newNd 插入有序的 nd 中
func (nd *node) insert(newNd *node) *node {
	if nd == nil {
		return newNd
	}
	if nd.isGuard {
		nd.next = nd.next.insert(newNd)
		nd.size = nd.next.size + 1
		return nd
	}
	comp := nd.key.CompareTo(newNd.key)
	if comp == CompareBigger {
		newNd.next = nd
		newNd.size = nd.size + 1
		return newNd
	} else if comp == CompareEqual {
		nd.value = newNd.value
		return nd
	} else {
		nd.next = nd.next.insert(newNd)
		nd.size = nd.next.size + 1
		return nd
	}
}

//
func (nd *node) delete(dlNd *node) *node {
	if nd == nil {
		return nil
	}
	if nd.isGuard {
		nd.next = nd.next.delete(dlNd)
		if nd.next != nil {
			nd.size = nd.next.size + 1
		}
		return nd
	}
	comp := nd.key.CompareTo(dlNd.key)
	if comp == CompareBigger {
		return nd
	} else if comp == CompareEqual {
		return nd.next
	} else {
		nd.next = nd.next.delete(dlNd)
		if nd.next != nil {
			nd.size = nd.next.size + 1
		}
		return nd
	}
}

type Skiplist struct {
	head *node  // point to highest level guard
	size int    // total element amount
}

func New() *Skiplist {
	rand.Seed(int64(time.Now().Nanosecond()))
	return &Skiplist{
		head: &node{
			next: nil,
			down: nil,
			isGuard: true,
			level: 0,
			key: nil,
			value: nil,
		},
		size: 0,
	}
}

func (sl *Skiplist) get(nd *node, key Comparer) (*node, bool) {
	if nd == nil {
		return nil, false
	}
	if !nd.isGuard { // 比较自身
		comp := nd.key.CompareTo(key)
		if comp == CompareEqual {
			return nd, true
		} else if comp == CompareBigger {
			return nil, false
		}
	}

	// 比较 next
	var nextNode *node
	if nd.next == nil {
		nextNode = nd.down
	} else {
		comp2 := nd.next.key.CompareTo(key)
		if comp2 == CompareEqual {
			return nd.next, true
		} else if comp2 == CompareBigger {
			nextNode = nd.down
		} else {
			nextNode = nd.next
		}
	}
	return sl.get(nextNode, key)
}

func (sl *Skiplist) Get(key Comparer) (interface{}, bool) {
	if n, ok := sl.get(sl.head, key); ok {
		return n.value, true
	} else {
		return nil, false
	}
}

func (sl *Skiplist) randomLevel() int {
	level := 0
	for k := rand.Intn(2); k > 0; k = rand.Intn(2) {
		level += 1
	}
	return level
}

func min(x, y int) int {
	if x >= y {
		return y
	} else {
		return x
	}
}

func (sl *Skiplist) Put(key Comparer, value interface{}) int {
	// 首先查找吗？Yes
	if foundNode, ok := sl.get(sl.head, key); ok {  // update
		foundNode.value = value
		return PutUpdate
	} else {
		level := sl.randomLevel()
		newNdList := []*node{}
		// 插入 newNd 到 [0, level] 中
		for levelGuard := sl.head; levelGuard != nil; levelGuard = levelGuard.down {
			if levelGuard.level > level {
				continue
			}
			newNd := &node{
				next:    nil,
				down:    nil,
				isGuard: false,
				level:   levelGuard.level,
				key:     key,
				value:   value,
			}
			newNdList = append(newNdList, newNd)  // levelMax -> level0
			levelGuard = levelGuard.insert(newNd)
		}
		// 插入 level - MaxLevel 个新链表
		for l := sl.MaxLevel()+1; l <= level; l++ {
			newNd := &node{
				next: nil,
				down: nil,
				isGuard: false,
				level: l,
				key: key,
				value: value,
			}
			newGuard := &node{
				next: newNd,
				down: nil,
				isGuard: true,
				level: l,
				key: nil,
				value: nil,
			}
			newGuard.insert(newNd)
			 newGuard.down = sl.head
			 sl.head = newGuard
			 newNdList = append([]*node{newNd}, newNdList...)
		}
		// set newNd.down
		for i := 0; i < len(newNdList)-1 ; i++ {
			newNdList[i].down = newNdList[i+1]
		}
		sl.size++
		return PutNew
	}
}

// 删除 [0, nd.level] 的所有值
func (sl *Skiplist) Delete(key Comparer) int {
	if foundNd, ok := sl.get(sl.head, key); ok { // delete it
		for currentGuard := sl.head; currentGuard != nil; currentGuard = currentGuard.down {
			if foundNd.level < currentGuard.level {
				continue
			}
			currentGuard = currentGuard.delete(foundNd)
		}
		sl.size--
		return DeleteYes
	} else {
		return DeleteNo
	}
}

func (sl *Skiplist) MaxLevel() int {
	return sl.head.level
}

func (sl *Skiplist) Len() int {
	return sl.size
}

func (sl *Skiplist) IsEmpty() bool {
	return sl.Len() == 0
}

func (sl *Skiplist) String() string {
	s := ""
	// stat info
	s += fmt.Sprintf("Skiplist: size [%v] level [%v]\n---\n", sl.Len(), sl.MaxLevel())
	for currentGuard := sl.head; currentGuard != nil; currentGuard = currentGuard.down {
		s += fmt.Sprintf("L[%v] ", currentGuard.level)
		for nd := currentGuard; nd != nil; nd = nd.next {
			if nd.isGuard {
				s += fmt.Sprintf("[G:%v] ", nd.size)
			} else {
				s += fmt.Sprintf("[%v:%v]", nd.key, nd.value)
				if nd.down != nil {
					s += fmt.Sprintf("(%v)", nd.down.key)
				}
				s += " "
			}
		}
		s += "\n"
	}
	s += "---\n"
	return s
}
