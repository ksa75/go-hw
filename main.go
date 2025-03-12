package main

import (
	"fmt"

	hw04 "hw04lrucache/hw04_lru_cache"
)

func main() {
	list := hw04.NewList()

	// Добавляем элементы в список
	list.PushBack(1)
	fmt.Println("1st element:", list.Front().Value)
	fmt.Println("last element:", list.Back().Value)

	list.PushBack(2)
	list.PushBack(3)
	list.PushFront(0)

	list.Print()
	fmt.Println("1st element:", list.Front().Value)
	fmt.Println("last element:", list.Back().Value)

	// Проверим длину списка
	fmt.Println("Length of list:", list.Len()) // Output: 4

	// Перемещаем элемент в начало
	list.MoveToFront(list.Back()) // Перемещаем последний элемент в начало

	list.Print()

	// Удаляем элемент
	list.Remove(list.Front()) // Удаляем первый элемент

	// Выводим элементы после удаления
	fmt.Println("After remove:")
	list.Print()

	fmt.Println("/////////////////////////////////////////////////////////////////////////")

	// Создаём кэш с ёмкостью 3
	cache := hw04.NewCache(3)

	// Добавляем элементы в кэш
	fmt.Println("Добавляем элементы в кэш:")
	cache.Set("a", "Первый")
	cache.Set("b", "Второй")
	cache.Set("c", "Третий")
	cache.Print()

	// Добавляем ещё один элемент, что должно вытолкнуть "a"
	cache.Set("d", "Четвертый")
	cache.Print()

	// Получаем элемент "b" — он должен быть перемещён в начало очереди
	if value, found := cache.Get("b"); found {
		fmt.Println("Получено значение по ключу 'b':", value)
	}
	cache.Print()

	// Добавляем ещё один элемент, что должно вытолкнуть "c"
	cache.Set("e", "Пятый")
	if value, found := cache.Get("e"); found {
		fmt.Println("Получено значение по ключу 'e':", value)
	}
	cache.Print()
}
