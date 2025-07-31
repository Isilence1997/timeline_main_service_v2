package model

import (
	videoTimeline "git.code.oa.com/grpcprotocol/component_plat/video_timeline_timeline_data"
	videoTimelineRpc "git.code.oa.com/grpcprotocol/component_plat/video_timeline_timeline_rpc"
)

// TimelineListRspDTO timeline列表返回体DTO
type TimelineListRspDTO struct {
	DataList    []*videoTimeline.ListItem
	PageContext string
	HasNextPage bool
}

// ProfPicFeedRspDTO 头像流返回体
type ProfPicFeedRspDTO struct {
	ItemList    []*videoTimelineRpc.FeedItem
	PageContext string
	HasNextPage bool
}
