package hw04lrucache

import (
	"fmt"
	"math/rand"
	"strconv"
	"sync"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestCache(t *testing.T) {
	t.Run("empty cache", func(t *testing.T) {
		c := NewCache(10)

		_, ok := c.Get("aaa")
		require.False(t, ok)

		_, ok = c.Get("bbb")
		require.False(t, ok)
	})

	t.Run("simple", func(t *testing.T) {
		c := NewCache(5)

		wasInCache := c.Set("aaa", 100)
		require.False(t, wasInCache)

		wasInCache = c.Set("bbb", 200)
		require.False(t, wasInCache)

		c.Print()

		val, ok := c.Get("aaa")
		require.True(t, ok)
		require.Equal(t, 100, val)

		val, ok = c.Get("bbb")
		require.True(t, ok)
		require.Equal(t, 200, val)
		c.Print()

		val, ok = c.Get("aaa")
		require.True(t, ok)
		require.NotEqual(t, 300, val)
		c.Print()

		val, ok = c.Get("ccc")
		require.False(t, ok)
		require.Nil(t, val)

		wasInCache = c.Set("aaa", 300)
		require.True(t, wasInCache)
		c.Print()
		fmt.Println("////////////////////////////")
	})
	/*
	   Ожидаются следующие тесты:
	   - на логику выталкивания элементов из-за размера очереди
	   (например: n = 3, добавили 4 элемента - 1й из кэша вытолкнулся);
	   - на логику выталкивания давно используемых элементов
	   (например: n = 3, добавили 3 элемента, обратились несколько раз к разным элементам:
	   изменили значение, получили значение и пр. - добавили 4й элемент,
	   из первой тройки вытолкнется тот элемент, что был затронут наиболее давно).
	*/
	t.Run("purge logic", func(t *testing.T) {
		c := NewCache(3)

		c.Set("aaa", 600)
		c.Set("bbb", 700)
		c.Set("ccc", 800)
		c.Print()

		c.Set("ddd", 900)
		// c.Print()

		_, ok := c.Get("aaa")
		require.False(t, ok)

		fmt.Println("////////////////////////////")

		c.Clear()

		c.Set("aaa", 600)
		c.Set("bbb", 700)
		c.Set("ccc", 800)
		// c.Print()

		fmt.Println(c.Get("aaa"))
		// c.Print()

		c.Set("ddd", 900)
		c.Print()
	})
}

func TestCacheMultithreading(t *testing.T) {
	t.Skip() // Remove me if task with asterisk completed.

	c := NewCache(10)
	wg := &sync.WaitGroup{}
	wg.Add(2)

	go func() {
		defer wg.Done()
		for i := 0; i < 1_000_000; i++ {
			c.Set(Key(strconv.Itoa(i)), i)
		}
	}()

	go func() {
		defer wg.Done()
		for i := 0; i < 1_000_000; i++ {
			c.Get(Key(strconv.Itoa(rand.Intn(1_000_000))))
		}
	}()

	wg.Wait()
}
