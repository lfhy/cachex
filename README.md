# Cachex
一个简单的缓存库。

# 读取
数据查询时会先命中内存缓存，命中失败会从持久化存储中读取至缓存中

# 写入
数据会先写入持久化存储再自动存储至内存缓存中

# 支持的缓存类型

## 内存存储
### Lru算法
Lru是一种内存数据淘汰策略，最近最少被使用的数据最先被淘汰。使用```cachex.NewMemLruCacheStroage```创建Lru内存缓存存储。

### Lfu算法
Lfu是一种内存数据淘汰策略，使用频率最低的数据最先被淘汰。使用```cachex.NewMemLfuCacheStroage```创建Lru内存缓存存储。

## 持久化存储
### Redis
使用Redis进行存储，使用```cachex.NewRedisCacheStroage```创建Redis缓存存储。

### LevelDB
使用LevelDB进行存储，使用```cachex.NewLevelDBCacheStroage```创建Redis缓存存储。


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