package cachex

type CacheStroage[T any] interface {
	// 获取缓存
	Get(key string) (T, bool)
	// 设置缓存
	Set(key string, value T)
	// 删除缓存
	Delete(key string)
	// 清空缓存
	Free()
}

type JsonByteData[T any] struct {
	Data T `json:"Data"`
}
