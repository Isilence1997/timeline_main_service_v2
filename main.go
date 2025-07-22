package main

import (
	_ "git.code.oa.com/tpstelemetry/tps-sdk-go/instrumentation/trpctelemetry"
	_ "git.code.oa.com/vlib/go/trpc_plugins/video_component_logreplay"
	_ "go.uber.org/automaxprocs"

	_ "git.code.oa.com/trpc-go/trpc-config-rainbow"
	_ "git.code.oa.com/trpc-go/trpc-config-tconf"
	_ "git.code.oa.com/trpc-go/trpc-filter/debuglog"
	_ "git.code.oa.com/trpc-go/trpc-filter/recovery"
	_ "git.code.oa.com/trpc-go/trpc-log-atta"
	_ "git.code.oa.com/trpc-go/trpc-metrics-m007"
	_ "git.code.oa.com/trpc-go/trpc-metrics-runtime"
	_ "git.code.oa.com/trpc-go/trpc-naming-polaris"
	_ "git.code.oa.com/trpc-go/trpc-opentracing-tjg"

	"git.code.oa.com/trpc-go/trpc-go/log"

	trpc "git.code.oa.com/trpc-go/trpc-go"
	pb "git.code.oa.com/trpcprotocol/component_plat/video_timeline_timeline_rpc"
)

type timelineServiceServiceImpl struct{}

func main() {

	s := trpc.NewServer()

	pb.RegisterTimelineServiceService(s, &timelineServiceServiceImpl{})

	if err := s.Serve(); err != nil {
		log.Fatal(err)
	}
}
