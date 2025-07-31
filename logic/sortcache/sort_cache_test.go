package sortcache

import (
	"context"
	"fmt"
	"reflect"
	"testing"
	"time"

	"bou.ke/monkey"
	"git.code.oa.com/v/main_logic/feeds/grpc_timeline_main_service_v2/logic/backsource"
	"git.code.oa.com/v/main_logic/feeds/grpc_timeline_main_service_v2/model"
)

func TestGetSortInfoMapAndBackSource(t *testing.T) {
	monkey.Patch(sortInfoBackSource, func(ctx context.Context, idList []string, scene, sourceKey string, sortInfoMap map[string]int64, cacheInfo *model.CacheInfoDTO) {
		return
	})

	monkey.Patch(batchGetSortInfoMap, func(ctx context.Context, idList []string, cacheInfo *model.CacheInfoDTO) (map[string]int64, error) {
		if len(idList) == 0 {
			return nil, fmt.Errorf("idlist error,idlist:%+v", idList)
		}
		res := map[string]int64{}
		for _, id := range idList {
			res[id] = 1234567890
		}
		return res, nil
	})
	defer func() {
		time.Sleep(time.Second * 1)
		monkey.UnpatchAll()
	}()

	type args struct {
		ctx       context.Context
		idList    []string
		scene     string
		sourceKey string
		cacheInfo *model.CacheInfoDTO
	}
	tests := []struct {
		name    string
		args    args
		want    map[string]int64
		wantErr bool
	}{
		{
			name: "error",
			args: args{
				ctx:       context.Background(),
				idList:    nil,
				scene:     "",
				sourceKey: "",
				cacheInfo: nil,
			},
			want:    nil,
			wantErr: true,
		},
		{
			name: "success",
			args: args{
				ctx:       context.Background(),
				idList:    []string{"123456789"},
				scene:     "",
				sourceKey: "",
				cacheInfo: nil,
			},
			want: map[string]int64{
				"123456789": 1234567890,
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := GetSortInfoMapAndBackSource(tt.args.ctx, tt.args.idList, tt.args.scene, tt.args.sourceKey, tt.args.cacheInfo)
			if (err != nil) != tt.wantErr {
				t.Errorf("GetSortInfoMapAndBackSource() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("GetSortInfoMapAndBackSource() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_sortInfoBackSource1(t *testing.T) {
	monkey.Patch(backsource.SendBackSourceMsg, func(ctx context.Context, entityID, scene, sourceKey string,
		backSourceConfig *model.KafkaProducerDTO) error {
		return nil
	})
	defer monkey.UnpatchAll()

	type args struct {
		ctx         context.Context
		idList      []string
		scene       string
		sourceKey   string
		sortInfoMap map[string]int64
		cacheInfo   *model.CacheInfoDTO
	}
	tests := []struct {
		name string
		args args
	}{
		{
			name: "No NeedBackSource",
			args: args{
				ctx:       context.Background(),
				idList:    []string{"123", "456"},
				scene:     "scene",
				sourceKey: "sourceKey",
				sortInfoMap: map[string]int64{
					"123": 4567890,
				},
				cacheInfo: &model.CacheInfoDTO{
					NeedBackSource: false,
				},
			},
		},
		{
			name: "NeedBackSource",
			args: args{
				ctx:       context.Background(),
				idList:    []string{"123", "456"},
				scene:     "scene",
				sourceKey: "sourceKey",
				sortInfoMap: map[string]int64{
					"123": 4567890,
				},
				cacheInfo: &model.CacheInfoDTO{
					NeedBackSource: true,
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			sortInfoBackSource(tt.args.ctx, tt.args.idList, tt.args.scene, tt.args.sourceKey, tt.args.sortInfoMap,
				tt.args.cacheInfo)
		})
	}
}
