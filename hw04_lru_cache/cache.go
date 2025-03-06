package hw04lrucache

import "fmt"

type Key string

type Cache interface {
	Set(key Key, value interface{}) bool
	Get(key Key) (interface{}, bool)
	// Print()
	Clear()
}

type lruCache struct {
	capacity int
	queue    List
	items    map[Key]*ListItem
	revMap   map[*ListItem]Key
}

func NewCache(capacity int) Cache {
	return &lruCache{
		capacity: capacity,
		queue:    NewList(),
		items:    make(map[Key]*ListItem, capacity),
		revMap:   make(map[*ListItem]Key, capacity),
	}
}

// Set добавляет или обновляет элемент в кэше.
func (c *lruCache) Set(key Key, value interface{}) bool {
	if item, exists := c.items[key]; exists {
		// Если элемент существует, обновляем его значение
		item.Value = value
		// Перемещаем элемент в начало очереди
		c.queue.MoveToFront(item)
		return true
	}

	// Если элемента нет, создаём новый
	item := &ListItem{Value: value}
	// Добавляем в очередь и в словарь
	c.queue.PushFront(item)
	c.items[key] = c.queue.Front()
	c.revMap[c.queue.Front()] = key

	// Если кэш переполнен, удаляем последний элемент
	if c.queue.Len() > c.capacity {
		// Удаляем элемент из очереди и из словаря
		removed := c.queue.Back()
		c.queue.Remove(removed)
		delete(c.items, c.revMap[removed])
		delete(c.revMap, removed)
	}

	return false
}

// Get возвращает элемент из кэша по ключу.
func (c *lruCache) Get(key Key) (interface{}, bool) {
	if item, exists := c.items[key]; exists {
		c.queue.MoveToFront(item)
		i, _ := item.Value.(*ListItem)
		return i.Value, true
	}
	return nil, false
}

func (c *lruCache) Clear() {
	c.queue = NewList()                           // Перезапускаем очередь
	c.items = make(map[Key]*ListItem, c.capacity) // Очищаем словарь
	c.revMap = make(map[*ListItem]Key, c.capacity)
}

// Print печатает элементы кэша для отладки.
func (c *lruCache) Print() {
	// Печатаем элементы кэша (содержимое очереди)
	fmt.Println("Состояние кэша:")

	// Получаем все элементы из очереди
	current := c.queue.Front()
	for current != nil {
		if item, ok := current.Value.(*ListItem); ok {
			fmt.Printf("%v\n", item.Value)
		}
		current = current.Next
	}
	// Печатаем словарь
	fmt.Println("Словарь:")
	for key, item := range c.items {
		fmt.Printf("%v: %v\n", key, item)
		if i, ok := item.Value.(*ListItem); ok {
			fmt.Printf("%v\n", i.Value)
		}
	}
	// Печатаем rev словарь
	for key, item := range c.revMap {
		fmt.Printf("%v: %v\n", key, item)
	}
}
