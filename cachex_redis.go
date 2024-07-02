package cachex

import (
	"context"
	"time"

	"github.com/bytedance/sonic"
	"github.com/redis/go-redis/v9"
)

type RedisConfig = redis.Options

type RedisCache[T any] struct {
	rdb     *redis.Client
	timeout time.Duration
}

// Redis存储
// 传入连接参数和缓存超时时间
func NewRedisCacheStroage(config *RedisConfig, timeout ...time.Duration) *RedisCache[any] {
	return NewRedisCacheStroageWithType[any](config, timeout...)
}

// Redis存储
// 传入连接参数和缓存超时时间
func NewRedisCacheStroageWithType[T any](config *RedisConfig, timeout ...time.Duration) *RedisCache[T] {
	t := time.Duration(0)
	if len(timeout) > 0 {
		t = timeout[0]
	}
	return &RedisCache[T]{
		rdb:     redis.NewClient(config),
		timeout: t,
	}
}

// 获取缓存
func (c *RedisCache[T]) Get(key string) (T, bool) {
	var getData JsonByteData[T]
	data, err := c.rdb.Get(context.Background(), key).Result()
	if err != nil {
		return getData.Data, false
	}
	return getData.Data, sonic.UnmarshalString(data, &getData) == nil
}

// 设置缓存
func (c *RedisCache[T]) Set(key string, value T) {
	var setData JsonByteData[T]
	setData.Data = value
	data, _ := sonic.MarshalString(setData)
	c.rdb.Set(context.Background(), key, data, c.timeout).Err()
}

// 删除缓存
func (c *RedisCache[T]) Delete(key string) {
	c.rdb.Del(context.Background(), key).Err()
}

// 清空缓存
func (c *RedisCache[T]) Free() {
	c.rdb.FlushDB(context.Background()).Err()
}
