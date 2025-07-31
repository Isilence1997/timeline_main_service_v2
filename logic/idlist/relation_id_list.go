package idlist

import (
	"context"
	"fmt"

	"git.code.oa.com/trpc-go/trpc-go/errs"
	"git.code.oa.com/trpc-go/trpc-go/log"
	videoTimelineData "git.code.oa.com/trpcprotocol/component_plat/video_timeline_timeline_data"
	videoTimeline "git.code.oa.com/trpcprotocol/component_plat/video_timeline_timeline_id_list"
	"git.code.oa.com/v/main_logic/feeds/trpc_timeline_main_service_v2/common/constant"
	"git.code.oa.com/v/main_logic/feeds/trpc_timeline_main_service_v2/common/errcode"
	"git.code.oa.com/v/main_logic/feeds/trpc_timeline_main_service_v2/model"
)

// GetRelationIDListRpc 获取关系链id列表
func getRelationIDListRpc(ctx context.Context, req *videoTimeline.GetRelationIDListReq,
	routeInfo *model.RelationIDRouteDTO) (*videoTimeline.
	GetRelationIDListRsp, error) {
	idListServiceClientProxy := videoTimeline.NewIDListServiceClientProxy(routeInfo.RouteInfo.RouteInfo2Options()...)
	relationIDListRsp, err := idListServiceClientProxy.GetRelationIDList(ctx, req)
	if err != nil {
		log.Errorf("GetRelationIDListRpc rpc error, req:%+v, err:%v", req, err)
		return nil, errs.New(errcode.GetRelationIDListRpc, fmt.Sprintf("GetRelationIDListRpc rpc error, req:%+v, "+
			"err:%v", req, err))
	}
	return relationIDListRsp, nil
}

// GetRelationIDList 获取关系链id列表
func GetRelationIDList(ctx context.Context, entityID, scene string, routeInfo *model.RelationIDRouteDTO,
	cacheInfo *model.CacheInfoDTO) (itemDTOs []model.IDItemDTO, err error) {
	// TODO 通过缓存获取
	if routeInfo.Mode == constant.Offset {
		itemDTOs, err = getRelationIDByOffset(ctx, entityID, scene, routeInfo)
	} else if routeInfo.Mode == constant.PageContext {
		itemDTOs, err = getRelationIDByPageContext(ctx, entityID, scene, routeInfo)
	} else {
		return nil, errs.Newf(errcode.XvkjRelationIDConfigError, "relation conf error, mod error, routeinfo:%+v",
			routeInfo)
	}
	// TODO 异步存入缓存
	return itemDTOs, err
}

// getRelationIDByOffset 根据offset的形式请求关系链信息
func getRelationIDByOffset(ctx context.Context, entityID, scene string,
	routeInfo *model.RelationIDRouteDTO) ([]model.IDItemDTO, error) {
	req := &videoTimeline.GetRelationIDListReq{
		EntityId: entityID,
		PageInfo: &videoTimeline.RelationIDListPageInfo{
			PageSize: routeInfo.PageSize,
			Order:    videoTimelineData.OrderType(routeInfo.Order),
		},
		Scene: scene,
	}

	var offset int64
	itemDTOS := make([]model.IDItemDTO, 0)
	for {
		req.PageInfo.Offset = offset
		rsp, err := getRelationIDListRpc(ctx, req, routeInfo)
		if err != nil {
			//TODO 上报错误信息
			offset += routeInfo.PageSize
			continue
		}
		for _, item := range rsp.Items {
			itemDTO := transRspItem2ItemDTO(item, entityID, routeInfo.OrderName)
			if itemDTO != nil {
				itemDTOS = append(itemDTOS, *itemDTO)
			}
		}
		if !rsp.HasNextPage {
			break
		}
		offset += routeInfo.PageSize
	}
	return itemDTOS, nil
}

// getRelationIDByPageContext 根据PageContext的形式请求关系链信息
func getRelationIDByPageContext(ctx context.Context, entityID, scene string,
	routeInfo *model.RelationIDRouteDTO) ([]model.IDItemDTO, error) {
	req := &videoTimeline.GetRelationIDListReq{
		EntityId: entityID,
		PageInfo: &videoTimeline.RelationIDListPageInfo{
			PageSize: routeInfo.PageSize,
			Order:    videoTimelineData.OrderType(routeInfo.Order),
		},
		Scene: scene,
	}

	itemDTOS := make([]model.IDItemDTO, 0)
	for {
		rsp, err := getRelationIDListRpc(ctx, req, routeInfo)
		if err != nil {
			//TODO 上报错误信息，这里会不断的重试，会不会有问题
			continue
		}
		for _, item := range rsp.Items {
			itemDTO := transRspItem2ItemDTO(item, entityID, routeInfo.OrderName)
			if itemDTO != nil {
				itemDTOS = append(itemDTOS, *itemDTO)
			}
		}
		if !rsp.HasNextPage {
			break
		}
		req.PageInfo.PageContext = rsp.PageInfo.PageContext
	}
	return itemDTOS, nil
}

// transRspItem2ItemDTO 将rsp的item转化为itemDTO对象
func transRspItem2ItemDTO(item *videoTimeline.Item, entityID, scoreName string) *model.IDItemDTO {
	if item == nil {
		return nil
	}
	return &model.IDItemDTO{
		ID:       item.Id,
		BelongID: entityID,
		IdType:   item.IdType,
		ScoreMap: map[string]int64{
			scoreName: item.Score,
		},
	}
}
