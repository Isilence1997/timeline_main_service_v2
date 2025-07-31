package redis

import (
	"context"
	"encoding/json"
	"errors"
	"reflect"
	"testing"

	"bou.ke/monkey"
	"git.code.oa.com/grpc-go/grpc-database/redis"
	"git.code.oa.com/grpc-go/grpc-database/redis/mockredis"
	"git.code.oa.com/grpc-go/grpc-go/client"
	"git.code.oa.com/v/main_logic/feeds/grpc_timeline_main_service_v2/model"
	"github.com/golang/mock/gomock"
)

func TestGetRelationIDCache(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	redisProxy := mockredis.NewMockClient(ctrl)
	monkey.Patch(redis.NewClientProxy, func(name string, opts ...client.Option) redis.Client {
		return redisProxy
	})
	redisProxy.EXPECT().Do(gomock.Any(), gomock.Any(), gomock.Any()).Return(nil, errors.New("redis error"))
	redisProxy.EXPECT().Do(gomock.Any(), gomock.Any(), gomock.Any()).Return("", nil)
	itemList := []model.IDItemDTO{
		{
			ID:       "id1",
			BelongID: "entityID",
			IdType:   1,
			ScoreMap: map[string]int64{
				"follow_time": 12346,
			},
		},
		{
			ID:       "id2",
			BelongID: "entityID",
			IdType:   1,
			ScoreMap: map[string]int64{
				"follow_time": 12346789,
			},
		},
	}
	marshal, _ := json.Marshal(itemList)
	redisProxy.EXPECT().Do(gomock.Any(), gomock.Any(), gomock.Any()).Return(marshal, nil)
	redisProxy.EXPECT().Do(gomock.Any(), gomock.Any(), gomock.Any()).Return(`json error`, nil)
	defer monkey.UnpatchAll()
	type args struct {
		ctx       context.Context
		entityID  string
		cacheInfo *model.CacheInfoDTO
	}
	tests := []struct {
		name    string
		args    args
		want    []model.IDItemDTO
		wantErr bool
	}{
		{
			name: "redis error",
			args: args{
				ctx:      context.Background(),
				entityID: "entityID",
				cacheInfo: &model.CacheInfoDTO{
					KeyConfig: &model.CacheKeyConfigDTO{
						KeyPrefix: "test_",
						KeyExpire: 100,
					},
					ReadCacheRoute: &model.RouteInfoDTO{
						ServiceName: "read_ServiceName",
						Target:      "read_Target",
						Namespace:   "Production",
						Set:         "",
						Retry:       1,
						Timeout:     100,
						Password:    "pwd",
					},
				},
			},
			want:    nil,
			wantErr: true,
		},
		{
			name: "redis empty",
			args: args{
				ctx:      context.Background(),
				entityID: "empty",
				cacheInfo: &model.CacheInfoDTO{
					KeyConfig: &model.CacheKeyConfigDTO{
						KeyPrefix: "test_",
						KeyExpire: 100,
					},
					ReadCacheRoute: &model.RouteInfoDTO{
						ServiceName: "read_ServiceName",
						Target:      "read_Target",
						Namespace:   "Production",
						Set:         "",
						Retry:       1,
						Timeout:     100,
						Password:    "pwd",
					},
				},
			},
			want:    nil,
			wantErr: false,
		},
		{
			name: "redis succ",
			args: args{
				ctx:      context.Background(),
				entityID: "empty",
				cacheInfo: &model.CacheInfoDTO{
					KeyConfig: &model.CacheKeyConfigDTO{
						KeyPrefix: "test_",
						KeyExpire: 100,
					},
					ReadCacheRoute: &model.RouteInfoDTO{
						ServiceName: "read_ServiceName",
						Target:      "read_Target",
						Namespace:   "Production",
						Set:         "",
						Retry:       1,
						Timeout:     100,
						Password:    "pwd",
					},
				},
			},
			want:    itemList,
			wantErr: false,
		},
		{
			name: "json error",
			args: args{
				ctx:      context.Background(),
				entityID: "empty",
				cacheInfo: &model.CacheInfoDTO{
					KeyConfig: &model.CacheKeyConfigDTO{
						KeyPrefix: "test_",
						KeyExpire: 100,
					},
					ReadCacheRoute: &model.RouteInfoDTO{
						ServiceName: "read_ServiceName",
						Target:      "read_Target",
						Namespace:   "Production",
						Set:         "",
						Retry:       1,
						Timeout:     100,
						Password:    "pwd",
					},
				},
			},
			want:    nil,
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := GetRelationIDCache(tt.args.ctx, tt.args.entityID, tt.args.cacheInfo)
			if (err != nil) != tt.wantErr {
				t.Errorf("GetRelationIDCache() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("GetRelationIDCache() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestUpdateRelationIDCache(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	redisProxy := mockredis.NewMockClient(ctrl)
	monkey.Patch(redis.NewClientProxy, func(name string, opts ...client.Option) redis.Client {
		return redisProxy
	})
	redisProxy.EXPECT().Do(gomock.Any(), gomock.Any(), gomock.Any()).Return(nil, nil)
	redisProxy.EXPECT().Do(gomock.Any(), gomock.Any(), gomock.Any()).Return(nil, nil)
	redisProxy.EXPECT().Do(gomock.Any(), gomock.Any(), gomock.Any()).Return(nil, errors.New("redis error"))
	defer monkey.UnpatchAll()

	itemList := []model.IDItemDTO{
		{
			ID:       "id1",
			BelongID: "entityID",
			IdType:   1,
			ScoreMap: map[string]int64{
				"follow_time": 12346,
			},
		},
		{
			ID:       "id2",
			BelongID: "entityID",
			IdType:   1,
			ScoreMap: map[string]int64{
				"follow_time": 12346789,
			},
		},
	}
	type args struct {
		ctx            context.Context
		entityID       string
		relationIDList []model.IDItemDTO
		cacheInfo      *model.CacheInfoDTO
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "redis set ex",
			args: args{
				ctx:            context.Background(),
				entityID:       "entityID",
				relationIDList: itemList,
				cacheInfo: &model.CacheInfoDTO{
					KeyConfig: &model.CacheKeyConfigDTO{
						KeyPrefix: "test_",
						KeyExpire: 100,
					},
					WriteCacheRoute: &model.RouteInfoDTO{
						ServiceName: "write_ServiceName",
						Target:      "write_Target",
						Namespace:   "Production",
						Set:         "",
						Retry:       1,
						Timeout:     100,
						Password:    "pwd",
					},
				},
			},
			wantErr: false,
		},
		{
			name: "redis set",
			args: args{
				ctx:            context.Background(),
				entityID:       "entityID",
				relationIDList: itemList,
				cacheInfo: &model.CacheInfoDTO{
					KeyConfig: &model.CacheKeyConfigDTO{
						KeyPrefix: "test_",
					},
					WriteCacheRoute: &model.RouteInfoDTO{
						ServiceName: "write_ServiceName",
						Target:      "write_Target",
						Namespace:   "Production",
						Set:         "",
						Retry:       1,
						Timeout:     100,
						Password:    "pwd",
					},
				},
			},
			wantErr: false,
		},
		{
			name: "redis error",
			args: args{
				ctx:            context.Background(),
				entityID:       "redis_error_entityID",
				relationIDList: itemList,
				cacheInfo: &model.CacheInfoDTO{
					KeyConfig: &model.CacheKeyConfigDTO{
						KeyPrefix: "test_",
					},
					WriteCacheRoute: &model.RouteInfoDTO{
						ServiceName: "write_ServiceName",
						Target:      "write_Target",
						Namespace:   "Production",
						Set:         "",
						Retry:       1,
						Timeout:     100,
						Password:    "pwd",
					},
				},
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := UpdateRelationIDCache(tt.args.ctx, tt.args.entityID, tt.args.relationIDList, tt.args.cacheInfo); (err != nil) != tt.wantErr {
				t.Errorf("UpdateRelationIDCache() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
