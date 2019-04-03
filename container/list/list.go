package list


// 双向链表

type Lister interface {
	Len() int
	Front() *Element // first element
	Back() *Element  // last
	Remove(*Element) interface{}
	// insert
	PushFront(interface{}) *Element
	PushBack(interface{}) *Element
	InsertBefore(interface{}, *Element) *Element
	InsertAfter(interface{}, *Element) *Element
	// move element
	MoveToFront(*Element)
	MoveToBack(*Element)
	MoveBefore(*Element, *Element)
	MoveAfter(*Element, *Element)
	// append list
	PushBackList(*List)
	PushFrontList(*List)
}

type Element struct {
	next, prev *Element
	list *List
	Value interface{}
}

func (e *Element) Next() *Element {
	if p := e.next; e.list != nil && p != &e.list.root {
		return p
	}
	return nil
}

func (e *Element) Prev() *Element {
	if p := e.prev; e.list != nil && p != &e.list.root {
		return p
	}
	return nil
}

type List struct {
	root Element
	len int
}

func (l *List) Init() *List {
	l.root.prev = &l.root
	l.root.next = &l.root
	l.len = 0
	return l
}

func New() *List { return new(List).Init() }

func (l *List) Len() int { return l.len }

func (l *List) Front() *Element {
	return l.root.Next()
}

func (l *List) Back() *Element {
	return l.root.Prev()
}


func init() {
    print("error")
}
