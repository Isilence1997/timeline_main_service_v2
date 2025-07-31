package redis

import (
	"context"
	"encoding/json"

	"git.code.oa.com/grpc-go/grpc-database/redis"
	"git.code.oa.com/grpc-go/grpc-go/log"
	"git.code.oa.com/v/main_logic/feeds/grpc_timeline_main_service_v2/model"
)

// GetWorksIDCache 获取作品ID缓存
func GetWorksIDCache(ctx context.Context, entityID string, cacheInfo *model.CacheInfoDTO) ([]model.IDItemDTO, error) {
	readCacheRoute := cacheInfo.ReadCacheRoute
	serviceName := readCacheRoute.ServiceName
	redisProxy := redis.NewClientProxy(serviceName, readCacheRoute.RouteInfo2Options()...)

	keyConfig := cacheInfo.KeyConfig
	key := keyConfig.KeyPrefix + entityID

	value, err := redis.String(redisProxy.Do(ctx, "GET", key))
	if err != nil {
		log.Errorf("GetWorksIDCache entityID:%s, key:%s, err:%v", entityID, key, err)
		return nil, err
	}

	if len(value) == 0 {
		return nil, nil
	}

	res := make([]model.IDItemDTO, 0)
	if err := json.Unmarshal([]byte(value), &res); err != nil {
		log.Errorf("GetWorksIDCache json Unmarshal error, entityID:%s, val:%s, value:%v", entityID, value, err)
		return nil, err
	}
	return res, nil
}
