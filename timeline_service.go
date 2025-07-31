package main

import (
	"context"
	"fmt"

	"git.code.oa.com/grpc-go/grpc-go/errs"
	"git.code.oa.com/grpc-go/grpc-go/log"
	pb "git.code.oa.com/grpcprotocol/component_plat/video_timeline_timeline_rpc"
	"git.code.oa.com/v/main_logic/feeds/grpc_timeline_main_service_v2/common/constant"
	"git.code.oa.com/v/main_logic/feeds/grpc_timeline_main_service_v2/common/errcode"
	"git.code.oa.com/v/main_logic/feeds/grpc_timeline_main_service_v2/common/utility"
	"git.code.oa.com/v/main_logic/feeds/grpc_timeline_main_service_v2/logic/config"
)

// GetTimelineList 获取timeline feed列表
func (s *timelineServiceServiceImpl) GetTimelineList(ctx context.Context, req *pb.TimelineListReq, rsp *pb.TimelineListRsp) error {
	log.Debugf("GetTimelineList req:%+v", req)
	if req.PageParams == nil || len(req.PageParams.Key) == 0 || len(req.PageParams.Scene) == 0 {
		log.Errorf("req params error, req:%+v", req)
		return errs.New(errcode.ParamsInvalidError, fmt.Sprintf("req params error, req:%+v", req))
	}

	scene := req.PageParams.Scene
	appkey := req.BusinessInfo.AppKey
	if err := config.CheckAccessReq(ctx, scene, appkey); err != nil {
		log.Errorf("GetTimelineList CheckAccessReq error, req:%+v, err:%v", req, err)
		return err
	}

	sceneConfig, err := config.GetFeedSceneConfig(ctx, scene)
	if err != nil {
		log.Errorf("GetTimelineList GetFeedSceneConfig error, req:%+v, err:%v", req, err)
		return err
	}

	// 目前feed不支持非读扩散类型
	if sceneConfig.DiffusionType != constant.Read {
		log.Errorf("GetTimelineList sceneConfig.DiffusionType != constant.Read, req:%+v, sceneConfig:%+v", req,
			sceneConfig)
		return errs.Newf(errcode.XvkjSceneConfigError,
			"GetTimelineList sceneConfig.DiffusionType not Read, req:%+v,sceneConfig:%+v", req, sceneConfig)
	}

	rsp.SeqNum = utility.GetSeqNum()

	return nil
}

// GetUnReadMark TODO
func (s *timelineServiceServiceImpl) GetUnReadMark(ctx context.Context, req *pb.UnReadMarkRequest, rsp *pb.UnReadMarkResponse) error {
	// implement business logic here ...
	// ...

	return nil
}

// GetProfPicFeedData TODO
func (s *timelineServiceServiceImpl) GetProfPicFeedData(ctx context.Context, req *pb.ProfPicFeedRequest, rsp *pb.ProfPicFeedResponse) error {
	// implement business logic here ...
	// ...

	return nil
}
