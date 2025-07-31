package idlist

import (
	"context"
	"errors"
	"reflect"
	"testing"

	"bou.ke/monkey"
	"git.code.oa.com/trpc-go/trpc-go/client"
	videoTimeline "git.code.oa.com/trpcprotocol/component_plat/video_timeline_timeline_id_list"
	"git.code.oa.com/v/main_logic/feeds/trpc_timeline_main_service_v2/model"
	"github.com/golang/mock/gomock"
)

func TestGetRelationIDListRpc(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	proxy := videoTimeline.NewMockIDListServiceClientProxy(ctrl)
	monkey.Patch(videoTimeline.NewIDListServiceClientProxy, func(opts ...client.Option) videoTimeline.
		IDListServiceClientProxy {
		return proxy
	})
	proxy.EXPECT().GetRelationIDList(gomock.Any(), gomock.Any(), gomock.Any()).Return(&videoTimeline.
		GetRelationIDListRsp{
		PageInfo:    nil,
		HasNextPage: false,
		Items: []*videoTimeline.Item{{
			Id:     "1234",
			Score:  123456,
			IdType: 1,
		}},
	}, nil)
	proxy.EXPECT().GetRelationIDList(gomock.Any(), gomock.Any(), gomock.Any()).Return(nil,
		errors.New("GetRelationIDListRsp error"))

	defer monkey.UnpatchAll()

	type args struct {
		ctx       context.Context
		req       *videoTimeline.GetRelationIDListReq
		routeInfo *model.RelationIDRouteDTO
	}
	tests := []struct {
		name string
		args args
		want *videoTimeline.
			GetRelationIDListRsp
		wantErr bool
	}{
		{
			"normal",
			args{
				ctx:       context.Background(),
				req:       nil,
				routeInfo: &model.RelationIDRouteDTO{},
			},
			&videoTimeline.GetRelationIDListRsp{
				PageInfo:    nil,
				HasNextPage: false,
				Items: []*videoTimeline.Item{{
					Id:     "1234",
					Score:  123456,
					IdType: 1,
				}},
			},
			false,
		},
		{
			"error",
			args{
				ctx:       context.Background(),
				req:       nil,
				routeInfo: &model.RelationIDRouteDTO{},
			},
			nil,
			true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := getRelationIDListRpc(tt.args.ctx, tt.args.req, tt.args.routeInfo)
			if (err != nil) != tt.wantErr {
				t.Errorf("GetRelationIDListRpc() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("GetRelationIDListRpc() got = %v, want %v", got, tt.want)
			}
		})
	}
}
