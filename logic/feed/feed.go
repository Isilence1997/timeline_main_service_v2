package feed

import (
	"context"
	"encoding/json"
	"fmt"
	"sort"

	"git.code.oa.com/grpc-go/grpc-go/log"
	videoTimeline "git.code.oa.com/grpcprotocol/component_plat/video_timeline_timeline_data"
	"git.code.oa.com/v/main_logic/feeds/grpc_timeline_main_service_v2/common/constant"
	"git.code.oa.com/v/main_logic/feeds/grpc_timeline_main_service_v2/common/utility"
	"git.code.oa.com/v/main_logic/feeds/grpc_timeline_main_service_v2/dao/xvkj"
	"git.code.oa.com/v/main_logic/feeds/grpc_timeline_main_service_v2/logic/common"
	"git.code.oa.com/v/main_logic/feeds/grpc_timeline_main_service_v2/logic/idlist"
	"git.code.oa.com/v/main_logic/feeds/grpc_timeline_main_service_v2/logic/sortcache"
	"git.code.oa.com/v/main_logic/feeds/grpc_timeline_main_service_v2/model"
)

// GetReadDiffusionFeed 获取读扩散feed流
func GetReadDiffusionFeed(ctx context.Context, entityID, pageContext string, sceneConfig *model.SceneConfig) (*model.
	TimelineListRspDTO, error) {
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

	//获取排序信息，feed目前只有一个排序值，就是用户更新时间，所以取第0个
	relationIDOrderName := readConfig.RelationIDOrder[0]
	relationOderCache := readConfig.CacheConfig[relationIDOrderName]
	sortInfoMap, err := sortcache.GetSortInfoMapAndBackSource(ctx, relationIDStrings, scene, relationIDOrderName,
		&relationOderCache)
	if err != nil {
		log.Errorf("GetReadDiffusionFeed GetSortInfo error, entityID:%s, relationIDList:+v%, err:%v", entityID,
			relationIDList,
			err)
		return nil, err
	}

	// 设置排序信息
	relationIDList = common.SetSortInfo2ItemDTO(relationIDList, sortInfoMap, relationIDOrderName)
	// 设置pageContext的排序值，覆盖原有排序值
	reqPageContextMap := pageContextStr2Map(pageContext)
	relationIDList = common.SetSortInfo2ItemDTO(relationIDList, reqPageContextMap, relationIDOrderName)

	// 全排序，取top k个用户
	sort.Slice(relationIDList, func(i, j int) bool {
		return relationIDList[i].ScoreMap[relationIDOrderName] > relationIDList[j].ScoreMap[relationIDOrderName]
	})
	relationIDList = relationIDList[0:scenePageSize]

	// 获取作品列表
	worksIDCacheInfo := readConfig.CacheConfig[constant.WorksIDListXvkjCacheKey]
	worksIDList, err := idlist.BatchGetWorksIDList(ctx, relationIDList, scene, relationIDOrderName, readConfig.WorksIDRoute, &worksIDCacheInfo)
	if err != nil {
		log.Errorf("GetReadDiffusionFeed BatchGetWorksIDList error, entityID:%s, err:%v", entityID, err)
		return nil, err
	}

	// top k个作品
	worksIDOrderName := readConfig.WorksIDOrder[0]
	sort.Slice(worksIDList, func(i, j int) bool {
		return relationIDList[i].ScoreMap[worksIDOrderName] > relationIDList[j].ScoreMap[worksIDOrderName]
	})
	resWorksIDList := worksIDList[0:scenePageSize]

	// 获取rsp的pageContext
	rspPageContextStr := getRspPageContextStr(reqPageContextMap, resWorksIDList, worksIDList[scenePageSize:], worksIDOrderName)
	return &model.TimelineListRspDTO{
		DataList:    transIDItemDTOs2ListItem(resWorksIDList, worksIDOrderName),
		PageContext: rspPageContextStr,
		HasNextPage: len(resWorksIDList) > 0,
	}, nil
}

// getRspPageContextStr 获取响应的pageContextStr
// reqPageContextMap 这个是请求过来的pageContext
// resWorksIDList 结果返回的作品列表
// surplusWorksIDList 剩余作品的作品列表
func getRspPageContextStr(reqPageContextMap map[string]int64, resWorksIDList, surplusWorksIDList []model.IDItemDTO,
	orderName string) string {
	// 用户标识map，这个map主要存储了返回的作品列表中，哪些用户被返回了
	userFlagMap := make(map[string]bool)
	for _, dto := range resWorksIDList {
		userFlagMap[dto.BelongID] = true
		reqPageContextMap[dto.BelongID] = dto.ScoreMap[orderName]
	}

	// 遍历剩余的作品
	// 如果这个作品在用户标识map中，则更新pageContextMap中的用户作品时间
	for _, dto := range surplusWorksIDList {
		belongID := dto.BelongID
		if len(userFlagMap) == 0 {
			break
		}
		// 判断该用户是否在用户标识map中，在的话证明用户被返回了
		//  需要需要在pageContext返回他的作品时间（剩余作品中的第一个作品时间）
		_, ok := userFlagMap[belongID]
		if !ok {
			continue
		}
		// 用户被返回了，所以把这用户的作品发布时间更新到pageContext中
		reqPageContextMap[belongID] = dto.ScoreMap[orderName]
		// 删除用户标识
		delete(userFlagMap, belongID)
	}

	marshal, _ := json.Marshal(reqPageContextMap)
	log.Debugf("getRspPageContextStr unzip res:%s", string(marshal))
	// 压缩pageContext
	rspPageContext := utility.ZipString(string(marshal))
	log.Debugf("getRspPageContextStr zip res:%s", string(marshal))
	return rspPageContext
}

// pageContextStr2Map 将pageContext转换为map
func pageContextStr2Map(pageContext string) map[string]int64 {
	res := make(map[string]int64)
	if len(pageContext) == 0 {
		return res
	}

	// 解压缩字符串
	unZipPageCtxStr := utility.UnZipString(pageContext)
	if err := json.Unmarshal([]byte(unZipPageCtxStr), res); err != nil {
		//should not happen，这里不返回错误了
		//TODO 加个上报
		log.Errorf("pageContextStr2Map Unmarshal error, json:%s, err:%v", unZipPageCtxStr, err)
	}
	return res
}

// transIDItemDTOs2ListItem 将IDItemDTO列表转化为timeline的listItem
func transIDItemDTOs2ListItem(items []model.IDItemDTO, orderName string) []*videoTimeline.ListItem {
	dataList := make([]*videoTimeline.ListItem, 0, len(items))
	for _, item := range items {
		listItem := &videoTimeline.ListItem{
			Member: item.ID,
			Value:  fmt.Sprintf("%d", item.ScoreMap[orderName]),
			IdType: item.IdType,
		}
		dataList = append(dataList, listItem)
	}
	return dataList
}

// itemDTOs2Strings 将itemDTO切片转化为string切片，主要是去itemDTO中的id
func itemDTOs2Strings(items []model.IDItemDTO) []string {
	res := make([]string, 0, len(items))
	for _, item := range items {
		res = append(res, item.ID)
	}
	return res
}
