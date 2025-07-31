package constant

import (
	"time"
)

const (
	// DefaultKafkaProducerTimeout 默认kafka生产者超时时间
	DefaultKafkaProducerTimeout = time.Second
	// DefaultMemoryCacheCapacity 默认内存缓存大小
	DefaultMemoryCacheCapacity = 5000
	// DefaultMemoryCacheExpire 默认内存缓存过期时间
	DefaultMemoryCacheExpire = time.Minute
	// DefaultRouteTimeout 默认路由超时时间
	DefaultRouteTimeout = time.Millisecond * 100
	// RelationIDListXvkjCacheKey 关系链ID列表在，在无极缓存配置中的key
	RelationIDListXvkjCacheKey = "relation_id_cache"
	// WorksIDListXvkjCacheKey 作品ID列表在，在无极缓存配置中的key
	WorksIDListXvkjCacheKey = "works_id_cache"
)
