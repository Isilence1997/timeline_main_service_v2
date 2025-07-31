package config

import (
	"context"

	"git.code.oa.com/grpc-go/grpc-go/errs"
	"git.code.oa.com/v/main_logic/feeds/grpc_timeline_main_service_v2/common/errcode"
	"git.code.oa.com/v/main_logic/feeds/grpc_timeline_main_service_v2/dao/xvkj"
)

// CheckAccessReq 检查接入请求
func CheckAccessReq(ctx context.Context, scene, appkey string) error {
	config, err := xvkj.GetAccessConfig(scene)
	if err != nil {
		return err

	}
	// 是否开启appkey校验
	if !config.AppkeyEnable {
		return nil
	}
	if appkey == config.Appkey {
		return nil
	}
	return errs.Newf(errcode.AccessKeyError, "CheckAccessReq appkey error, config:%+v, req.appkey:%s", config, appkey)
}
