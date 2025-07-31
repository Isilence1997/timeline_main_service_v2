package avatarlist

import (
	"reflect"
	"testing"

	videoTimelineRpc "git.code.oa.com/grpcprotocol/component_plat/video_timeline_timeline_rpc"
	"git.code.oa.com/v/main_logic/feeds/grpc_timeline_main_service_v2/model"
)

func Test_transIDItemDTOs2FeedItem(t *testing.T) {
	type args struct {
		items []model.IDItemDTO
	}
	tests := []struct {
		name string
		args args
		want []*videoTimelineRpc.FeedItem
	}{
		{
			name: "succ",
			args: args{
				items: []model.IDItemDTO{
					{
						ID:       "id1",
						BelongID: "entityID",
						IdType:   1,
						ScoreMap: map[string]int64{
							"follow_time": 123,
							"live_time":   456,
						},
					},
					{
						ID:       "id2",
						BelongID: "entityID",
						IdType:   1,
						ScoreMap: map[string]int64{
							"follow_time": 456,
							"live_time":   789,
						},
					},
				},
			},
			want: []*videoTimelineRpc.FeedItem{
				{
					UserId: "id1",
					ExtraData: map[string]string{
						"follow_time": "123",
						"live_time":   "456",
					},
				},
				{
					UserId: "id2",
					ExtraData: map[string]string{
						"follow_time": "456",
						"live_time":   "789",
					},
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := transIDItemDTOs2FeedItem(tt.args.items); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("transIDItemDTOs2FeedItem() = %v, want %v", got, tt.want)
			}
		})
	}
}
