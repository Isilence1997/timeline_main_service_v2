package avatarlist

import (
	"context"
	"fmt"
	"sort"

	"git.code.oa.com/grpc-go/grpc-go/log"
	videoTimelineRpc "git.code.oa.com/grpcprotocol/component_plat/video_timeline_timeline_rpc"
	"git.code.oa.com/v/main_logic/feeds/grpc_timeline_main_service_v2/common/constant"
	"git.code.oa.com/v/main_logic/feeds/grpc_timeline_main_service_v2/dao/xvkj"
	"git.code.oa.com/v/main_logic/feeds/grpc_timeline_main_service_v2/logic/common"
	"git.code.oa.com/v/main_logic/feeds/grpc_timeline_main_service_v2/logic/idlist"
	"git.code.oa.com/v/main_logic/feeds/grpc_timeline_main_service_v2/logic/sortcache"
	"git.code.oa.com/v/main_logic/feeds/grpc_timeline_main_service_v2/model"
)

// GetReadDiffusionAvatarList 读扩散获取头像横滑列表
func GetReadDiffusionAvatarList(ctx context.Context, entityID, pageContext string,
	sceneConfig *model.SceneConfig) (*model.ProfPicFeedRspDTO, error) {
	// 获取读扩散配置
	readConfigID := sceneConfig.ReadConfigID
	readConfig, err := xvkj.GetReadConfig(ctx, readConfigID)
	if err != nil {
		log.Errorf("GetReadDiffusionFeed GetReadConfig error, readConfigID:%s, err:%v", readConfigID, err)
		return nil, err
	}
	scene := sceneConfig.Scene
	scenePageSize := sceneConfig.PageSize

	// 获取关系链列表
	relationIDCacheInfo := readConfig.CacheConfig[constant.RelationIDListXvkjCacheKey]
	relationIDList, err := idlist.GetRelationIDList(ctx, entityID, scene, readConfig.RelationIDRoute,
		&relationIDCacheInfo)
	if err != nil {
		log.Errorf("GetReadDiffusionFeed GetRelationIDList error, entityID:%s, err:%v", entityID, err)
		return nil, err
	}
	relationIDStrings := common.ItemDTOs2Strings(relationIDList)

	// 当存在排序项的时候，需要进行获取缓存排序
	if len(readConfig.RelationIDOrder) != 0 {
		relationIDList = getAfterSortingRelationIDList(ctx, entityID, scene, relationIDStrings, readConfig,
			relationIDList)
	}

	// 根据pageContext截取
	var resRelationIDList []model.IDItemDTO
	if len(pageContext) == 0 {
		resRelationIDList = relationIDList[0:scenePageSize]
	} else {

	}
	return &model.ProfPicFeedRspDTO{
		ItemList:    transIDItemDTOs2FeedItem(resRelationIDList),
		PageContext: "",
		HasNextPage: false,
	}, err
}

// getAfterSortingRelationIDList 获取排序后的关系链id列表
// relationIDStrings 关系string id列表
func getAfterSortingRelationIDList(ctx context.Context, entityID, scene string, relationIDStrings []string,
	readConfig *model.ReadConfigDTO, sourceRelationIDList []model.IDItemDTO) []model.IDItemDTO {
	sortInfoCh := make(chan model.ChannelSortInfoDTO, len(readConfig.RelationIDOrder))
	// 拉取所有的排序信息
	for _, orderName := range readConfig.RelationIDOrder {
		go func(ctx context.Context, orderName string) {
			channelSortInfoDTO := model.ChannelSortInfoDTO{}
			defer func() {
				sortInfoCh <- channelSortInfoDTO
			}()

			orderCache := readConfig.CacheConfig[orderName]
			sortInfoMap, err := sortcache.GetSortInfoMapAndBackSource(ctx, relationIDStrings, scene, orderName, // 按顺序获取每个排序指标的map[id]score
				&orderCache)
			if err != nil {
				log.Errorf("GetReadDiffusionAvatarList GetSortInfo error, entityID:%s, "+
					"orderName:%s, relationIDList:+v%, err:%v", entityID, sourceRelationIDList, err)
				channelSortInfoDTO.Err = err
				return
			}
			channelSortInfoDTO.SortInfoMap = sortInfoMap
			channelSortInfoDTO.OrderName = orderName
		}(ctx, orderName)
	}

	// 将排序信息塞入到待排序的item中
	for i := 0; i < len(readConfig.RelationIDOrder); i++ {
		channelSortInfoDTO := <- sortInfoCh
		if channelSortInfoDTO.Err != nil {
			//TODO 做些上报之类的
			continue
		}
		sourceRelationIDList = common.SetSortInfo2ItemDTO(sourceRelationIDList, channelSortInfoDTO.SortInfoMap,
			channelSortInfoDTO.OrderName) //排序值写入IDItemDTO中
	}

	// 多值排序
	sort.Slice(sourceRelationIDList, func(i, j int) bool {
		for _, orderName := range readConfig.RelationIDOrder { // 各排序值的优先级由无极表配置
			if sourceRelationIDList[i].ScoreMap[orderName] > sourceRelationIDList[j].ScoreMap[orderName] {
				return true
			}
		}
		return false
	})
	return sourceRelationIDList
}

// transIDItemDTOs2FeedItem 将IDItemDTO列表转化为timeline的FeedItem
func transIDItemDTOs2FeedItem(items []model.IDItemDTO) []*videoTimelineRpc.FeedItem {
	dataList := make([]*videoTimelineRpc.FeedItem, 0, len(items))
	for _, item := range items {
		listItem := &videoTimelineRpc.FeedItem{
			UserId:    item.ID,
			ExtraData: transMapStringInt642MapStringString(item.ScoreMap),
		}
		dataList = append(dataList, listItem)
	}
	return dataList
}

// transMapStringInt642MapStringString map[string]int64 转换位map[string]string
func transMapStringInt642MapStringString(source map[string]int64) map[string]string {
	res := make(map[string]string, len(source))
	for s, i := range source {
		res[s] = fmt.Sprintf("%d", i)
	}
	return res
}
