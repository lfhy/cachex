# Cachex
一个简单的缓存库。

# 读取
数据查询时会先命中内存缓存，命中失败会从持久化存储中读取至缓存中

# 写入
数据会先写入持久化存储再自动存储至内存缓存中

# 使用方法
见`example`包
```golang
package main

import (
	"fmt"

	"github.com/lfhy/cachex"
)

// 内存缓存使用lru
// 远程缓存使用redis
func main() {
	// 初始化缓存配置
	lruCache := cachex.NewMemLruCacheStroage(30000)
	redisCache := cachex.NewRedisCacheStroage(&cachex.RedisConfig{
		Addr:     "192.168.188.230:6379",
		Password: "password",
		DB:       4,
	}, 0)
	// 建立缓存
	cache := cachex.NewCache(lruCache, redisCache)
	// 设置缓存
	cache.Set("hello", "world")
	// 读取缓存
	value, ok := cache.Get("hello")
	if ok {
		fmt.Printf("第一次从lru读取: %v\n", value)
	} else {
		fmt.Println("读取失败")
	}
	// Lru将缓存删除
	lruCache.Delete("hello")
	// 重新读取缓存
	value, ok = cache.Get("hello")
	if ok {
		fmt.Printf("第二次从redis读: %v\n", value)
	} else {
		fmt.Println("读取失败")
	}
}
```