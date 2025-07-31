package common

import (
	"git.code.oa.com/v/main_logic/feeds/grpc_timeline_main_service_v2/model"
)

// ItemDTOs2Strings 将itemDTO切片转化为string切片，主要是去itemDTO中的id
func ItemDTOs2Strings(items []model.IDItemDTO) []string {
	res := make([]string, 0, len(items))
	for _, item := range items {
		res = append(res, item.ID)
	}
	return res
}

// SetSortInfo2ItemDTO 设置排序值到IDItemDTO中
func SetSortInfo2ItemDTO(items []model.IDItemDTO, sortInfoMap map[string]int64, oderName string) []model.IDItemDTO {
	for _, item := range items {
		score, ok := sortInfoMap[item.ID]
		if ok {
			item.ScoreMap[oderName] = score
		}
	}
	return items
}
