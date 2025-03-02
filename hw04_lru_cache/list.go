package hw04lrucache

// type List interface {
// 	Len() int
// 	Front() *ListItem
// 	Back() *ListItem
// 	PushFront(v interface{}) *ListItem
// 	PushBack(v interface{}) *ListItem
// 	Remove(i *ListItem)
// 	MoveToFront(i *ListItem)
// }

type ListItem struct {
	Value interface{}
	Next  *ListItem
	Prev  *ListItem
}

type List struct {
	len  int       // длина списка
	head *ListItem // первый элемент списка
	tail *ListItem // последний элемент списка
}

// Конструктор нового списка
func NewList() *List {
	return &List{}
}

// Метод Len возвращает длину списка
func (l *List) Len() int {
	return l.len
}

// Метод Front возвращает первый элемент списка
func (l *List) Front() *ListItem {
	return l.head
}

// Метод Back возвращает последний элемент списка
func (l *List) Back() *ListItem {
	return l.tail
}

func (l *List) PushFront(v interface{}) *ListItem {
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

// Метод PushBack добавляет элемент в конец списка
func (l *List) PushBack(v interface{}) *ListItem {
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

// Метод Remove удаляет элемент из списка
func (l *List) Remove(i *ListItem) {
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

// Метод MoveToFront перемещает элемент в начало списка
func (l *List) MoveToFront(i *ListItem) {
	if i == nil || i == l.head {
		return
	}

	// Удаляем элемент из текущего места
	l.Remove(i)

	// Добавляем элемент в начало
	l.PushFront(i.Value)
}
