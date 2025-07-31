package feed

import (
	"reflect"
	"testing"

	videoTimeline "git.code.oa.com/grpcprotocol/component_plat/video_timeline_timeline_data"
	"git.code.oa.com/v/main_logic/feeds/grpc_timeline_main_service_v2/model"
)

func Test_transIDItemDTOs2ListItem(t *testing.T) {
	type args struct {
		items     []model.IDItemDTO
		orderName string
	}
	tests := []struct {
		name string
		args args
		want []*videoTimeline.ListItem
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
							"upload_time": 123456,
						},
					},
					{
						ID:       "id2",
						BelongID: "entityID",
						IdType:   1,
						ScoreMap: map[string]int64{},
					},
				},
				orderName: "upload_time",
			},
			want: []*videoTimeline.ListItem{
				{
					Member: "id1",
					Value:  "123456",
					IdType: 1,
				},
				{
					Member: "id2",
					Value:  "0",
					IdType: 1,
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := transIDItemDTOs2ListItem(tt.args.items, tt.args.orderName); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("transIDItemDTOs2ListItem() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_pageContextStr2Map(t *testing.T) {
	type args struct {
		pageContext string
	}
	tests := []struct {
		name string
		args args
		want map[string]int64
	}{
		{
			name: "empty",
			args: args{
				pageContext: "",
			},
			want: map[string]int64{},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := pageContextStr2Map(tt.args.pageContext); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("pageContextStr2Map() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_getRspPageContextStr(t *testing.T) {
	type args struct {
		reqPageContextMap  map[string]int64
		resWorksIDList     []model.IDItemDTO
		surplusWorksIDList []model.IDItemDTO
		orderName          string
	}
	tests := []struct {
		name string
		args args
		want string
	}{}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := getRspPageContextStr(tt.args.reqPageContextMap, tt.args.resWorksIDList, tt.args.surplusWorksIDList, tt.args.orderName); got != tt.want {
				t.Errorf("getRspPageContextStr() = %v, want %v", got, tt.want)
			}
		})
	}
}
