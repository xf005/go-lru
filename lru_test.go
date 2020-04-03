package lru

import (
	"fmt"
	"testing"
	"time"
)

func TestGet(t *testing.T) {
	lru := NewCache(10, 5)
	for i := 0; i < 10; i++ {
		lru.Set(i, fmt.Sprint(i)+"-"+fmt.Sprint(i))
		fmt.Println(lru.Keys())
		time.Sleep(1 * time.Second)
	}
	fmt.Println(lru.Get(1))
	fmt.Println(lru.Get(8))
	fmt.Println(lru.Keys())
	fmt.Println(lru.Values())
	lru.Flush()
	fmt.Println(lru.Len())
	fmt.Println(lru.Keys())
}
