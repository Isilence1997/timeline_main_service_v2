package redis

import (
	"context"

	"git.code.oa.com/trpc-go/trpc-database/redis"
	"git.code.oa.com/trpc-go/trpc-go/errs"
	"git.code.oa.com/trpc-go/trpc-go/log"
	"git.code.oa.com/v/main_logic/feeds/trpc_timeline_main_service_v2/common/errcode"
	"git.code.oa.com/v/main_logic/feeds/trpc_timeline_main_service_v2/model"
)

// GetSortInfoMap 获取排序信息
func GetSortInfoMap(ctx context.Context, idList []string, cacheInfo *model.CacheInfoDTO) (map[string]int64, error) {
	readCacheRoute := cacheInfo.ReadCacheRoute
	serviceName := readCacheRoute.ServiceName
	redisProxy := redis.NewClientProxy(serviceName, readCacheRoute.RouteInfo2Options()...)
	conn, err := redisProxy.Pipeline(ctx)
	if err != nil {
		return nil, errs.Newf(errcode.RedisPipelineNewError, "redis new pipeline error, ServiceName:%s, err:%v",
			serviceName, err)
	}

	succIDs := make([]string, 0)
	errIDs := make([]string, 0)
	for _, id := range idList {
		key := cacheInfo.KeyConfig.KeyPrefix + id
		err := conn.Send("GET", key)
		if err != nil {
			errIDs = append(errIDs, id)
			log.Errorf("redis pipline send error, serviceName:%s, key:%s, err:%v", serviceName, key, err)
			continue
		}
		succIDs = append(succIDs, id)
	}

	if err = conn.Flush(); err != nil {
		return nil, errs.Newf(errcode.RedisPipelineNewError, "redis pipeline flush error, ServiceName:%s, err:%v",
			serviceName, err)
	}
	res := make(map[string]int64, 0)
	for _, id := range succIDs {
		receive, err := redis.Int64(conn.Receive())
		key := cacheInfo.KeyConfig.KeyPrefix + id
		if err != nil {
			errIDs = append(errIDs, id)
			log.Errorf("redis pipline receive error, serviceName:%s, key:%s, err:%v", serviceName, key, err)
			continue
		}
		res[id] = receive
	}
	return res, nil
}
