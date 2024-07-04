package main

import (
	"fmt"

	"github.com/lfhy/cachex"
)

// 内存缓存使用Fifo
// 持久缓存使用Level
func main() {
	// 初始化缓存配置
	cache1 := cachex.NewMemFifoCacheStroage(30000)
	cache2, err := cachex.NewLevelDBCacheStroage("./tmp")
	if err != nil {
		panic(err)
	}
	// 建立缓存
	cache := cachex.NewCache(cache1, cache2)
	// 设置缓存
	cache.Set("hello", "world")
	// 读取缓存
	value, ok := cache.Get("hello")
	if ok {
		fmt.Printf("第一次从第一层读取: %v\n", value)
	} else {
		fmt.Println("读取失败")
	}
	// 将缓存删除
	cache1.Delete("hello")
	// 重新读取缓存
	value, ok = cache.Get("hello")
	if ok {
		fmt.Printf("第二次从第二层读: %v\n", value)
	} else {
		fmt.Println("读取失败")
	}
}
