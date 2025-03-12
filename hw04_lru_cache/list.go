package hw04lrucache

import "fmt"

type List interface {
	Len() int
	Front() *ListItem
	Back() *ListItem
	PushFront(v interface{}) *ListItem
	PushBack(v interface{}) *ListItem
	Remove(i *ListItem)
	MoveToFront(i *ListItem)
	// Print()
}

type ListItem struct {
	Value interface{}
	Next  *ListItem
	Prev  *ListItem
}

type list struct {
	len  int       // длина списка
	head *ListItem // первый элемент списка
	tail *ListItem // последний элемент списка
}

// Конструктор нового списка.
func NewList() List {
	return &list{}
}

// Метод Len возвращает длину списка.
func (l *list) Len() int {
	return l.len
}

// Метод Front возвращает первый элемент списка.
func (l *list) Front() *ListItem {
	return l.head
}

// Метод Back возвращает последний элемент списка.
func (l *list) Back() *ListItem {
	return l.tail
}

func (l *list) PushFront(v interface{}) *ListItem {
	item := &ListItem{Value: v}

	// Если список пуст, новый элемент будет и первым, и последним
	if l.len == 0 {
		l.head = item
		l.tail = item
	} else {
		// Добавляем элемент в начало
		item.Next = l.head
		l.head.Prev = item
		l.head = item
	}

	l.len++
	return item
}

// Метод PushBack добавляет элемент в конец списка.
func (l *list) PushBack(v interface{}) *ListItem {
	item := &ListItem{Value: v}

	// Если список пуст, новый элемент будет и первым, и последним
	if l.len == 0 {
		l.head = item
		l.tail = item
	} else {
		// Добавляем элемент в конец
		item.Prev = l.tail
		l.tail.Next = item
		l.tail = item
	}

	l.len++
	return item
}

// Метод Remove удаляет элемент из списка.
func (l *list) Remove(i *ListItem) {
	if i == nil {
		return
	}

	// Если элемент первый
	if i == l.head {
		l.head = i.Next
	}

	// Если элемент последний
	if i == l.tail {
		l.tail = i.Prev
	}

	// Переподключаем предыдущий и следующий элементы
	if i.Prev != nil {
		i.Prev.Next = i.Next
	}
	if i.Next != nil {
		i.Next.Prev = i.Prev
	}

	l.len--
}

// Метод MoveToFront перемещает элемент в начало списка.
func (l *list) MoveToFront(i *ListItem) {
	if i == nil || i == l.head {
		return
	}

	// Удаляем элемент из текущего места
	l.Remove(i)

	// Перемещаем элемент в начало
	i.Next = l.head
	l.head.Prev = i
	l.head = i
	l.len++
}

func (l *list) Print() {
	for item := l.Front(); item != nil; item = item.Next {
		fmt.Println(item.Value)
	}
}
