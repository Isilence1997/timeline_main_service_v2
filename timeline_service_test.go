package main

import (
	"context"
	"testing"

	grpc "git.code.oa.com/grpc-go/grpc-go"
	_ "git.code.oa.com/grpc-go/grpc-go/http"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"

	pb "git.code.oa.com/grpcprotocol/component_plat/video_timeline_timeline_rpc"
)

var timelineServiceService = &timelineServiceServiceImpl{}

//go:generate go mod tidy
//go:generate mockgen -destination=stub/git.code.oa.com/grpcprotocol/component_plat/video_timeline_timeline_rpc/timeline_service_mock.go -package=video_timeline_timeline_rpc -self_package=git.code.oa.com/grpcprotocol/component_plat/video_timeline_timeline_rpc git.code.oa.com/grpcprotocol/component_plat/video_timeline_timeline_rpc TimelineServiceClientProxy

func Test_TimelineService_GetTimelineList(t *testing.T) {

	// 开始写mock逻辑
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	timelineServiceClientProxy := pb.NewMockTimelineServiceClientProxy(ctrl)

	// 预期行为
	m := timelineServiceClientProxy.EXPECT().GetTimelineList(gomock.Any(), gomock.Any(), gomock.Any()).AnyTimes()

	m.DoAndReturn(func(ctx context.Context, req interface{}, opts ...interface{}) (interface{}, error) {

		r, ok := req.(*pb.TimelineListReq)
		if !ok {
			panic("invalid request")
		}

		rsp := &pb.TimelineListRsp{}
		err := timelineServiceService.GetTimelineList(grpc.BackgroundContext(), r, rsp)
		return rsp, err
	})

	// 开始写单元测试逻辑
	req := &pb.TimelineListReq{}

	rsp, err := timelineServiceClientProxy.GetTimelineList(grpc.BackgroundContext(), req)

	// 输出入参和返回 (检查t.Logf输出，运行 `go test -v`)
	t.Logf("TimelineService_GetTimelineList req: %v", req)
	t.Logf("TimelineService_GetTimelineList rsp: %v, err: %v", rsp, err)

	assert.Nil(t, err)
}

func Test_TimelineService_GetUnReadMark(t *testing.T) {

	// 开始写mock逻辑
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	timelineServiceClientProxy := pb.NewMockTimelineServiceClientProxy(ctrl)

	// 预期行为
	m := timelineServiceClientProxy.EXPECT().GetUnReadMark(gomock.Any(), gomock.Any(), gomock.Any()).AnyTimes()

	m.DoAndReturn(func(ctx context.Context, req interface{}, opts ...interface{}) (interface{}, error) {

		r, ok := req.(*pb.UnReadMarkRequest)
		if !ok {
			panic("invalid request")
		}

		rsp := &pb.UnReadMarkResponse{}
		err := timelineServiceService.GetUnReadMark(grpc.BackgroundContext(), r, rsp)
		return rsp, err
	})

	// 开始写单元测试逻辑
	req := &pb.UnReadMarkRequest{}

	rsp, err := timelineServiceClientProxy.GetUnReadMark(grpc.BackgroundContext(), req)

	// 输出入参和返回 (检查t.Logf输出，运行 `go test -v`)
	t.Logf("TimelineService_GetUnReadMark req: %v", req)
	t.Logf("TimelineService_GetUnReadMark rsp: %v, err: %v", rsp, err)

	assert.Nil(t, err)
}

func Test_TimelineService_GetProfPicFeedData(t *testing.T) {

	// 开始写mock逻辑
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	timelineServiceClientProxy := pb.NewMockTimelineServiceClientProxy(ctrl)

	// 预期行为
	m := timelineServiceClientProxy.EXPECT().GetProfPicFeedData(gomock.Any(), gomock.Any(), gomock.Any()).AnyTimes()

	m.DoAndReturn(func(ctx context.Context, req interface{}, opts ...interface{}) (interface{}, error) {

		r, ok := req.(*pb.ProfPicFeedRequest)
		if !ok {
			panic("invalid request")
		}

		rsp := &pb.ProfPicFeedResponse{}
		err := timelineServiceService.GetProfPicFeedData(grpc.BackgroundContext(), r, rsp)
		return rsp, err
	})

	// 开始写单元测试逻辑
	req := &pb.ProfPicFeedRequest{}

	rsp, err := timelineServiceClientProxy.GetProfPicFeedData(grpc.BackgroundContext(), req)

	// 输出入参和返回 (检查t.Logf输出，运行 `go test -v`)
	t.Logf("TimelineService_GetProfPicFeedData req: %v", req)
	t.Logf("TimelineService_GetProfPicFeedData rsp: %v, err: %v", rsp, err)

	assert.Nil(t, err)
}
