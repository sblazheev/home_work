package hw04lrucache

type List interface {
	Len() int
	Front() *ListItem
	Back() *ListItem
	PushFront(v interface{}) *ListItem
	PushBack(v interface{}) *ListItem
	Remove(i *ListItem)
	MoveToFront(i *ListItem)
}

type ListItem struct {
	Value interface{}
	Next  *ListItem
	Prev  *ListItem
}

type list struct {
	first *ListItem
	last  *ListItem
	len   int
}

func (l *list) PushFront(v interface{}) *ListItem {
	item := &ListItem{Value: v, Next: l.first}
	if l.first != nil {
		l.first.Prev = item
	}
	l.first = item
	if l.last == nil {
		l.last = item
	}
	l.len++
	return item
}

func (l *list) PushBack(v interface{}) *ListItem {
	item := &ListItem{Value: v, Prev: l.last}
	if l.last != nil {
		l.last.Next = item
	}
	l.last = item
	if l.first == nil {
		l.first = item
	}
	l.len++
	return item
}

func (l *list) Front() *ListItem {
	return l.first
}

func (l *list) Back() *ListItem {
	return l.last
}

func (l *list) Len() int {
	return l.len
}

func (l *list) Remove(i *ListItem) {
	if i.Next != nil {
		i.Next.Prev = i.Prev
	} else {
		l.last = i.Prev
	}
	if i.Prev == nil {
		l.first = i.Next
	} else {
		i.Prev.Next = i.Next
	}
	l.len--
}

func (l *list) MoveToFront(i *ListItem) {
	if l.first != i {
		l.Remove(i)
		front := l.Front()
		l.first = i
		if front != nil {
			front.Prev = i
			i.Next = front
		}
		i.Prev = nil
		if l.last == nil {
			l.last = front
		}
		l.len++
	}
}

func NewList() List {
	return new(list)
}
