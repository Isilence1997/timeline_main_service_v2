package idlist

import (
	"reflect"
	"testing"

	videoTimeline "git.code.oa.com/trpcprotocol/component_plat/video_timeline_timeline_id_list"
	"git.code.oa.com/v/main_logic/feeds/trpc_timeline_main_service_v2/model"
)

func Test_transRspItem2ItemDTO(t *testing.T) {
	type args struct {
		item      *videoTimeline.Item
		entityID  string
		scoreName string
	}
	tests := []struct {
		name string
		args args
		want *model.IDItemDTO
	}{
		{
			name: "nil",
			args: args{
				item:      nil,
				entityID:  "nil",
				scoreName: "nil",
			},
			want: nil,
		},
		{
			name: "succ",
			args: args{
				item: &videoTimeline.Item{
					Id:     "id",
					Score:  123,
					IdType: 1,
				},
				entityID:  "entityID",
				scoreName: "follow_time",
			},
			want: &model.IDItemDTO{
				ID:       "id",
				BelongID: "entityID",
				IdType:   1,
				ScoreMap: map[string]int64{
					"follow_time": 123,
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := transRspItem2ItemDTO(tt.args.item, tt.args.entityID, tt.args.scoreName); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("transRspItem2ItemDTO() = %v, want %v", got, tt.want)
			}
		})
	}
}
