package model

import (
	"git.code.oa.com/grpcprotocol/component_plat/video_timeline_timeline_data"
)

// IDItemDTO id项dto
type IDItemDTO struct {
	ID       string //当前的item的ID
	BelongID string //归属ID，代表当前ID的归属者是哪个ID
	IdType   video_timeline_timeline_data.IdType
	ScoreMap map[string]int64
}

// ChannelIDITemDTO 当Err为nil时，IDItemDTO才可用
type ChannelIDITemDTO struct {
	IDItemDTO []IDItemDTO
	Err       error
}

// ChannelSortInfoDTO 当Err为nil时，SortInfoMap才可用
type ChannelSortInfoDTO struct {
	SortInfoMap map[string]int64
	Err         error
	OrderName   string
}
