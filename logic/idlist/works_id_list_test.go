package idlist

import (
	"context"
	"errors"
	"reflect"
	"testing"

	"bou.ke/monkey"
	"git.code.oa.com/trpc-go/trpc-go/client"
	videoTimeline "git.code.oa.com/trpcprotocol/component_plat/video_timeline_timeline_id_list"
	"github.com/golang/mock/gomock"
)

func TestGetWorksIDListRpc(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	proxy := videoTimeline.NewMockIDListServiceClientProxy(ctrl)
	monkey.Patch(videoTimeline.NewIDListServiceClientProxy, func(opts ...client.Option) videoTimeline.
		IDListServiceClientProxy {
		return proxy
	})
	proxy.EXPECT().GetWorksIDList(gomock.Any(), gomock.Any(), gomock.Any()).Return(&videoTimeline.
		GetWorksIDListRsp{
		PageInfo:    nil,
		HasNextPage: false,
		Items: []*videoTimeline.Item{{
			Id:     "1234",
			Score:  123456,
			IdType: 1,
		}},
	}, nil)
	proxy.EXPECT().GetWorksIDList(gomock.Any(), gomock.Any(), gomock.Any()).Return(nil,
		errors.New("GetRelationIDListRsp error"))

	type args struct {
		ctx context.Context
		req *videoTimeline.GetWorksIDListReq
	}
	tests := []struct {
		name string
		args args
		want *videoTimeline.
			GetWorksIDListRsp
		wantErr bool
	}{
		{
			"normal",
			args{
				ctx: nil,
				req: nil,
			},
			&videoTimeline.GetWorksIDListRsp{
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
				ctx: nil,
				req: nil,
			},
			nil,
			true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := GetWorksIDListRpc(tt.args.ctx, tt.args.req)
			if (err != nil) != tt.wantErr {
				t.Errorf("GetWorksIDListRpc() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("GetWorksIDListRpc() got = %v, want %v", got, tt.want)
			}
		})
	}
}
