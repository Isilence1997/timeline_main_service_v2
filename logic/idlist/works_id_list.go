package idlist

import (
	"context"
	"fmt"
	"strings"
	"time"

	trpc "git.code.oa.com/trpc-go/trpc-go"
	"git.code.oa.com/trpc-go/trpc-go/errs"
	"git.code.oa.com/trpc-go/trpc-go/log"
	videoTimelineData "git.code.oa.com/trpcprotocol/component_plat/video_timeline_timeline_data"
	videoTimelineIDList "git.code.oa.com/trpcprotocol/component_plat/video_timeline_timeline_id_list"
	"git.code.oa.com/v/main_logic/feeds/trpc_timeline_main_service_v2/common/constant"
	"git.code.oa.com/v/main_logic/feeds/trpc_timeline_main_service_v2/common/errcode"
	"git.code.oa.com/v/main_logic/feeds/trpc_timeline_main_service_v2/logic/backsource"
	"git.code.oa.com/v/main_logic/feeds/trpc_timeline_main_service_v2/model"
)

// GetWorksIDListRpc 获取作品id列表
func getWorksIDListRpc(ctx context.Context, req *videoTimelineIDList.GetWorksIDListReq,
	worksIDRouteDTO *model.WorksIDRouteDTO) (*videoTimelineIDList.
	GetWorksIDListRsp, error) {
	idListServiceClientProxy := videoTimelineIDList.NewIDListServiceClientProxy(worksIDRouteDTO.RouteInfo.RouteInfo2Options()...)
	worksIDListRsp, err := idListServiceClientProxy.GetWorksIDList(ctx, req)
	if err != nil {
		log.Errorf("GetWorksIDListRpc rpc error, req:%+v, err:%v", req, err)
		return nil, errs.New(errcode.GetWorksIDListRpcError, fmt.Sprintf("GetWorksIDListRpc rpc error, req:%+v, "+
			"err:%v", req, err))
	}
	return worksIDListRsp, nil
}

// getWorksIDList 获取作品id列表
func getWorksIDList(ctx context.Context, entityID, scene string, timestamp int64,
	worksIDRouteDTO *model.WorksIDRouteDTO,
	cacheInfo *model.CacheInfoDTO) (itemDTOs []model.IDItemDTO, err error) {
	// TODO 通过缓存获取
	if timestamp <= 0 {
		timestamp = time.Now().Unix()
	}
	// 多取一个，这样可以做pageContext透传优化
	pageSize := worksIDRouteDTO.PageSize + 1

	// 缓存没有，发起回源消息
	if cacheInfo.NeedBackSource {
		go backsource.SendBackSourceMsg(trpc.CloneContext(ctx), entityID, scene, constant.WorksIDListXvkjCacheKey,
			cacheInfo.BackSourceConfig)
	}

	worksIDListReq := &videoTimelineIDList.GetWorksIDListReq{
		EntityId: entityID,
		PageInfo: &videoTimelineIDList.WorksIDListPageInfo{
			Timestamp: timestamp,
			PageSize:  pageSize,
			Order:     videoTimelineData.OrderType(worksIDRouteDTO.Order),
		},
		Scene: scene,
	}
	rsp, err := getWorksIDListRpc(ctx, worksIDListReq, worksIDRouteDTO)
	if err != nil {
		log.Errorf("GetWorksIDList getWorksIDListRpc error, err:%v", err)
		return nil, err
	}

	itemDTOS := make([]model.IDItemDTO, 0)
	for _, item := range rsp.Items {
		itemDTO := transRspItem2ItemDTO(item, entityID, worksIDRouteDTO.OrderName)
		if itemDTO != nil {
			itemDTOS = append(itemDTOS, *itemDTO)
		}
	}
	return itemDTOs, err
}

// BatchGetWorksIDList 批量获取作品id列表
func BatchGetWorksIDList(ctx context.Context, items []model.IDItemDTO, scene, orderName string,
	worksIDRouteDTO *model.WorksIDRouteDTO,
	cacheInfo *model.CacheInfoDTO) ([]model.IDItemDTO, error) {
	itemsDTOsCh := make(chan model.ChannelIDITemDTO, len(items))
	for _, item := range items {
		go func(ctx context.Context, entityID string, timestamp int64) {
			channelIDITemDTO := model.ChannelIDITemDTO{}
			defer func() {
				itemsDTOsCh <- channelIDITemDTO
			}()
			list, err := getWorksIDList(ctx, entityID, scene, timestamp, worksIDRouteDTO, cacheInfo)
			if err != nil {
				log.Errorf("BatchGetWorksIDList getWorksIDList error, err:%v", err)
				channelIDITemDTO.Err = err
				return
			}
			channelIDITemDTO.IDItemDTO = list
		}(ctx, item.ID, item.ScoreMap[orderName])
	}

	res := make([]model.IDItemDTO, 0, len(items))
	errMsg := make([]string, 0)
	for i := 0; i < len(items); i++ {
		channelIDITemDTO := <-itemsDTOsCh
		if channelIDITemDTO.Err != nil {
			errMsg = append(errMsg, channelIDITemDTO.Err.Error())
			//TODO 做一些上报之类的操作
			continue
		}
		if len(channelIDITemDTO.IDItemDTO) > 0 {
			res = append(res, channelIDITemDTO.IDItemDTO...)
		}
	}

	// 全部都失败了，就返回错误，部分失败不返回错误
	if len(errMsg) == len(items) {
		return nil, errs.New(errcode.BatchGetWorksIDListAllFail, strings.Join(errMsg, "||"))
	}
	return res, nil
}
