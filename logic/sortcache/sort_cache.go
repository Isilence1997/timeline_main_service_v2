package sortcache

import (
	"context"
	"fmt"
	"strings"
	"sync"

	grpc "git.code.oa.com/grpc-go/grpc-go"
	"git.code.oa.com/grpc-go/grpc-go/errs"
	"git.code.oa.com/v/main_logic/feeds/grpc_timeline_main_service_v2/dao/redis"
	"git.code.oa.com/v/main_logic/feeds/grpc_timeline_main_service_v2/logic/backsource"
	"git.code.oa.com/v/main_logic/feeds/grpc_timeline_main_service_v2/model"
)

// GetSortInfoMapAndBackSource 获取排序信并且回源
func GetSortInfoMapAndBackSource(ctx context.Context, idList []string, scene, sourceKey string,
	cacheInfo *model.CacheInfoDTO) (map[string]int64, error) {
	sortInfoMap, err := batchGetSortInfoMap(ctx, idList, cacheInfo)
	if err != nil {
		return nil, err
	}
	// 由于是异步操作，所以把ctx clone一份，可以防止流程结束，导致ctx cancel
	go sortInfoBackSource(grpc.CloneContext(ctx), idList, scene, sourceKey, sortInfoMap, cacheInfo)
	return sortInfoMap, nil
}

// sortInfoBackSource 排序信息回源
func sortInfoBackSource(ctx context.Context, idList []string, scene, sourceKey string, sortInfoMap map[string]int64, cacheInfo *model.CacheInfoDTO) {
	// 不需要回源
	if !cacheInfo.NeedBackSource {
		return
	}
	for _, id := range idList {
		_, ok := sortInfoMap[id]
		if ok {
			continue
		}
		_ = backsource.SendBackSourceMsg(ctx, id, scene, sourceKey, cacheInfo.BackSourceConfig)
	}
}

const (
	batchGetSize = 50
)

// batchGetSortInfoMap  批量获取排序信息，key->id, value->排序值
func batchGetSortInfoMap(ctx context.Context, idList []string, cacheInfo *model.CacheInfoDTO) (map[string]int64, error) {
	idListLen := len(idList)
	if idListLen <= batchGetSize {
		return redis.GetSortInfoMap(ctx, idList, cacheInfo)
	}

	// 计算批次
	num := idListLen / batchGetSize
	if idListLen%batchGetSize != 0 {
		num += 1
	}

	res := map[string]int64{}
	errMsgs := make([]string, 0, num)
	lock := sync.Mutex{}
	for i := 0; i < num; i++ {
		left := i * batchGetSize
		right := (i + 1) * batchGetSize
		if right > idListLen {
			right = idListLen
		}

		go func(ctx context.Context, idList []string, cacheInfo *model.CacheInfoDTO) {
			infoMap, err := redis.GetSortInfoMap(ctx, idList, cacheInfo)
			lock.Lock()
			defer lock.Unlock()

			if err != nil {
				errMsgs = append(errMsgs, err.Error())
				return
			}
			for k, v := range infoMap {
				res[k] = v
			}
		}(ctx, idList[left:right], cacheInfo)
	}

	if len(errMsgs) == 0 {
		// 没有失败
		return res, nil
	} else if len(errMsgs) == num {
		// 全失败
		return nil, errs.Newf(0, fmt.Sprintf("batchGetSortInfoMap error, err:%s", strings.Join(errMsgs, "||")))
	} else {
		// 部分失败
		return res, nil
	}
}
