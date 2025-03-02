package main

import "fmt"

// Key — это тип для ключей в кэше
type Key string

// Cache — интерфейс для кэша
type Cache interface {
	Set(key Key, value interface{}) bool
	Get(key Key) (interface{}, bool)
	Clear()
}

// lruCache — структура для LRU-кэша
type lruCache struct {
	// Cache // Убираем, так как это не нужно после реализации

	capacity int
	queue    *List
	items    map[Key]*ListItem
}

// NewCache создаёт новый LRU-кэш с заданной ёмкостью
func NewCache(capacity int) Cache {
	return &lruCache{
		capacity: capacity,
		queue:    NewList(),
		items:    make(map[Key]*ListItem, capacity),
	}
}

// Set добавляет или обновляет элемент в кэше
func (c *lruCache) Set(key Key, value interface{}) bool {
	// Если элемент уже существует, обновляем его и перемещаем в начало
	if item, found := c.items[key]; found {
		item.Value = value
		c.queue.MoveToFront(item)
		return false // Элемент был обновлён
	}

	// Если элемента нет в кэше, добавляем его
	item := c.queue.PushFront(value)
	c.items[key] = item

	// Если кэш переполнен, удаляем наименее используемый элемент (в конце списка)
	if len(c.items) > c.capacity {
		c.evict()
		return true // Элемент был добавлен, а старый вытолкнут
	}

	return true // Новый элемент был добавлен
}

// Get возвращает значение по ключу, если оно существует
func (c *lruCache) Get(key Key) (interface{}, bool) {
	// Если элемент найден, перемещаем его в начало списка
	if item, found := c.items[key]; found {
		c.queue.MoveToFront(item)
		return item.Value, true
	}

	// Если элемента нет, возвращаем nil
	return nil, false
}

// Clear очищает кэш
func (c *lruCache) Clear() {
	c.items = make(map[Key]*ListItem, c.capacity)
	c.queue = NewList() // Очищаем список
}

// evict удаляет наименее используемый элемент (в конце списка)
func (c *lruCache) evict() {
	// Наименее используемый элемент — это последний элемент в списке
	if c.queue.tail != nil {
		delete(c.items, c.queue.tail.Key)
		c.queue.Remove(c.queue.tail)
	}
}

func main() {
	// Создаём кэш с ёмкостью 3
	cache := NewCache(3)

	// Добавляем элементы
	cache.Set("a", 1)
	cache.Set("b", 2)
	cache.Set("c", 3)

	// Получаем элемент
	if value, found := cache.Get("a"); found {
		fmt.Println("Get a:", value) // Output: 1
	}

	// Добавляем еще один элемент, что приведет к вытолкнутому элементу
	cache.Set("d", 4)

	// Проверяем содержимое кэша
	if value, found := cache.Get("b"); found {
		fmt.Println("Get b:", value) // Output: 2
	}

	// Добавляем еще один элемент, что приведет к вытолкнутому элементу
	cache.Set("e", 5)

	// Проверяем содержимое кэша
	if value, found := cache.Get("c"); found {
		fmt.Println("Get c:", value) // Output: nil (не найден)
	}

	// Очищаем кэш
	cache.Clear()
	fmt.Println("Cache cleared")

	// Проверяем содержимое после очистки
	if value, found := cache.Get("a"); found {
		fmt.Println("Get a after clear:", value)
	} else {
		fmt.Println("a not found after clear")
	}
}
