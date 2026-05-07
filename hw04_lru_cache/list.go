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
	Head   *ListItem
	Tail   *ListItem
	Length int
}

func NewList() List {
	return new(list)
}

func (l *list) Len() int {
	return l.Length
}

func (l *list) Front() *ListItem {
	return l.Head
}

func (l *list) Back() *ListItem {
	return l.Tail
}

func (l *list) PushFront(v interface{}) *ListItem {
	newHead := &ListItem{Value: v}
	l.Length++

	if l.Head == nil {
		l.Head = newHead
		l.Tail = newHead
		return newHead
	}
	//добавляем указатель на следующий элемент, то есть на старую голову
	newHead.Next = l.Head
	//Нужно у старой головы засетить prev element так как он nil
	l.Head.Prev = newHead
	//Новая голова списка
	l.Head = newHead
	return newHead
}

func (l *list) PushBack(v interface{}) *ListItem {
	l.Length++
	if l.Tail == nil {
		newItem := &ListItem{Value: v}
		l.Tail = newItem

		if l.Head == nil {
			l.Head = newItem
		}
		return l.Tail
	}
	newTail := new(ListItem)
	newTail.Value = v
	//аналогично head
	newTail.Prev = l.Tail
	l.Tail.Next = newTail
	l.Tail = newTail
	return newTail
}

func (l *list) Remove(i *ListItem) {
	if l.Length == 1 {
		l.Head = nil
		l.Tail = nil
	} else if l.Head == i {
		l.Head = i.Next
		l.Head.Prev = nil
	} else if l.Tail == i {
		l.Tail = i.Prev
		l.Tail.Next = nil
	} else {
		i.Prev.Next = i.Next
		i.Next.Prev = i.Prev
	}

	l.Length--
}

func (l *list) MoveToFront(i *ListItem) {
	if l.Head == nil || i == l.Head {
		return
	}

	if i.Prev != nil {
		i.Prev.Next = i.Next
	}
	if i.Next != nil {
		i.Next.Prev = i.Prev
	}
	if i == l.Tail {
		l.Tail = i.Prev
	}

	i.Prev = nil
	i.Next = l.Head
	l.Head.Prev = i
	l.Head = i
}
