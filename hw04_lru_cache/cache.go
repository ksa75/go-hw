package hw04lrucache

import "fmt"

type Key string

type Cache interface {
	Set(key Key, value interface{}) bool
	Get(key Key) (interface{}, bool)
	// Print()
	Clear()
}

type cacheItem struct {
	value interface{}
	key   Key
}

type lruCache struct {
	capacity int
	queue    List
	items    map[Key]*ListItem
}

func NewCache(capacity int) Cache {
	return &lruCache{
		capacity: capacity,
		queue:    NewList(),
		items:    make(map[Key]*ListItem, capacity),
	}
}

// Set добавляет или обновляет элемент в кэше.
func (c *lruCache) Set(key Key, value interface{}) bool {
	if item, exists := c.items[key]; exists {
		// Если элемент существует, обновляем его значение
		item.Value = cacheItem{value, key}
		// Перемещаем элемент в начало очереди
		c.queue.MoveToFront(item)
		return true
	}

	// Если элемента нет, создаём новый
	item := &ListItem{Value: cacheItem{value, key}}
	// Добавляем в очередь и в словарь
	c.queue.PushFront(item)
	c.items[key] = c.queue.Front()

	// Если кэш переполнен, удаляем последний элемент
	if c.queue.Len() > c.capacity {
		// Удаляем элемент из очереди и из словаря
		removed := c.queue.Back()
		if i, ok := removed.Value.(*ListItem); ok {
			if j, ok := i.Value.(cacheItem); ok {
				delete(c.items, j.key)
				c.queue.Remove(removed)
			}
		}
	}
	return false
}

// Get возвращает элемент из кэша по ключу.
func (c *lruCache) Get(key Key) (interface{}, bool) {
	if item, exists := c.items[key]; exists {
		c.queue.MoveToFront(item)
		if i, ok := item.Value.(*ListItem); ok {
			if j, ok := i.Value.(cacheItem); ok {
				return j.value, true
			}
		}
	}
	return nil, false
}

func (c *lruCache) Clear() {
	c.queue = NewList()                           // Перезапускаем очередь
	c.items = make(map[Key]*ListItem, c.capacity) // Очищаем словарь
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
	}
}
