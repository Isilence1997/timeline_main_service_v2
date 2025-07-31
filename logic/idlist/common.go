package idlist

import (
	videoTimeline "git.code.oa.com/trpcprotocol/component_plat/video_timeline_timeline_id_list"
	"git.code.oa.com/v/main_logic/feeds/trpc_timeline_main_service_v2/model"
)

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
