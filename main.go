package main

import (
	_ "git.code.oa.com/tpstelemetry/tps-sdk-go/instrumentation/grpctelemetry"
	_ "git.code.oa.com/vlib/go/grpc_plugins/video_component_logreplay"
	_ "go.uber.org/automaxprocs"

	_ "git.code.oa.com/grpc-go/grpc-config-rainbow"
	_ "git.code.oa.com/grpc-go/grpc-config-tconf"
	_ "git.code.oa.com/grpc-go/grpc-filter/debuglog"
	_ "git.code.oa.com/grpc-go/grpc-filter/recovery"
	_ "git.code.oa.com/grpc-go/grpc-log-atta"
	_ "git.code.oa.com/grpc-go/grpc-metrics-m007"
	_ "git.code.oa.com/grpc-go/grpc-metrics-runtime"
	_ "git.code.oa.com/grpc-go/grpc-naming-polaris"
	_ "git.code.oa.com/grpc-go/grpc-opentracing-tjg"

	"git.code.oa.com/grpc-go/grpc-go/log"

	grpc "git.code.oa.com/grpc-go/grpc-go"
	pb "git.code.oa.com/grpcprotocol/component_plat/video_timeline_timeline_rpc"
)

type timelineServiceServiceImpl struct{}

func main() {

	s := grpc.NewServer()

	pb.RegisterTimelineServiceService(s, &timelineServiceServiceImpl{})

	if err := s.Serve(); err != nil {
		log.Fatal(err)
	}
}
