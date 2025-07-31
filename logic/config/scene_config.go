package config

import (
	"context"

	"git.code.oa.com/grpc-go/grpc-go/errs"
	"git.code.oa.com/v/main_logic/feeds/grpc_timeline_main_service_v2/common/constant"
	"git.code.oa.com/v/main_logic/feeds/grpc_timeline_main_service_v2/common/errcode"
	"git.code.oa.com/v/main_logic/feeds/grpc_timeline_main_service_v2/dao/xvkj"
	"git.code.oa.com/v/main_logic/feeds/grpc_timeline_main_service_v2/model"
)

// GetFeedSceneConfig 获取feed流场景配置
func GetFeedSceneConfig(ctx context.Context, scene string) (*model.SceneConfig, error) {
	sceneConfig, err := xvkj.GetSceneConfigConfig(scene)
	if err != nil {
		return nil, err
	}

	if sceneConfig.SceneType != constant.Feed {
		return nil, errs.Newf(errcode.XvkjSceneConfigError, "GetFeedSceneConfig SceneType error, scene:%s, "+
			"sceneConfig:%s", scene, sceneConfig)
	}

	// 目前feed流只支持读扩散模式，后续支持写扩散需要将此判断去掉
	if sceneConfig.DiffusionType != constant.Read {
		return nil, errs.Newf(errcode.XvkjSceneConfigError, "GetFeedSceneConfig DiffusionType error, scene:%s, "+
			"sceneConfig:%s", scene, sceneConfig)
	}
	return sceneConfig, err
}

// GetAvatarListSceneConfig 获取头像列表场景配置
func GetAvatarListSceneConfig(ctx context.Context, scene string) (*model.SceneConfig, error) {
	sceneConfig, err := xvkj.GetSceneConfigConfig(scene)
	if err != nil {
		return nil, err
	}

	if sceneConfig.SceneType != constant.AvatarList {
		return nil, errs.Newf(errcode.XvkjSceneConfigError, "GetAvatarListSceneConfig SceneType error, scene:%s, "+
			"sceneConfig:%s", scene, sceneConfig)
	}

	// 目前头像列表只支持读扩散模式，后续支持写扩散需要将此判断去掉
	if sceneConfig.DiffusionType != constant.Read {
		return nil, errs.Newf(errcode.XvkjSceneConfigError, "GetAvatarListSceneConfig DiffusionType error, scene:%s, "+
			"sceneConfig:%s", scene, sceneConfig)
	}
	return sceneConfig, err
}

// GetRefreshSceneConfig 获取下拉刷新（红点）场景配置
func GetRefreshSceneConfig(ctx context.Context, scene string) (*model.SceneConfig, error) {
	sceneConfig, err := xvkj.GetSceneConfigConfig(scene)
	if err != nil {
		return nil, err
	}

	if sceneConfig.SceneType != constant.Refresh {
		return nil, errs.Newf(errcode.XvkjSceneConfigError, "GetRefreshSceneConfig SceneType error, scene:%s, "+
			"sceneConfig:%s", scene, sceneConfig)
	}

	// 目前下拉刷新（红点）只支持写扩散模式，后续支持读扩散将此判断删除
	if sceneConfig.DiffusionType != constant.Write {
		return nil, errs.Newf(errcode.XvkjSceneConfigError, "GetRefreshSceneConfig DiffusionType error, scene:%s, "+
			"sceneConfig:%s", scene, sceneConfig)
	}
	return sceneConfig, err
}
