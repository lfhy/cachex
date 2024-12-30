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
	// 关闭
	Close()
	// 遍历
	Range(fn func(key string, value T) (Continue bool))
}

type noCloseStroage struct{}

func (noCloseStroage) Close() {}

type JsonByteData[T any] struct {
	Data T `json:"Data"`
}
