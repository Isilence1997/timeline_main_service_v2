package xvkj

import (
	"context"
	"encoding/json"
	"strings"

	"git.code.oa.com/open-xvkj/go-sdk/xvkjclient"
	grpc "git.code.oa.com/grpc-go/grpc-go"
	"git.code.oa.com/grpc-go/grpc-go/errs"
	"git.code.oa.com/grpc-go/grpc-go/log"
	"git.code.oa.com/v/main_logic/feeds/grpc_timeline_main_service_v2/common/errcode"
	"git.code.oa.com/v/main_logic/feeds/grpc_timeline_main_service_v2/config"
	"git.code.oa.com/v/main_logic/feeds/grpc_timeline_main_service_v2/model"
)

var (
	accessConfigFilter xvkjclient.FilterInterface
	sceneConfigFilter  xvkjclient.FilterInterface
	readConfigXvkj     xvkjclient.ClientInterface
)

// InitXvkj 初始化xvkj配置
func InitXvkj() error {
	err := grpc.GoAndWait(func() error {
		return initAccessConfigFilter()
	}, func() error {
		return initSceneConfigFilter()
	}, func() error {
		return initReadConfigXvkj()
	})
	return err
}

// initReadConfigXvkj 初始化读扩散配置表
func initReadConfigXvkj() error {
	readConfig := config.GetConfig().ReadConfigXvkj
	readConfigXvkjTmp, err := xvkjclient.NewClient(
		xvkjclient.WithAppID(readConfig.AppID),
		xvkjclient.WithSchemaKey(readConfig.SchemaKey),
		xvkjclient.WithSchemaID(readConfig.SchemaID),
		xvkjclient.EnsureInitLatestData(),
	)
	if err != nil {
		log.Errorf("initReadConfigXvkj init accessXvkj error, err:%+v", err)
		return err
	}
	readConfigXvkj = readConfigXvkjTmp
	return nil
}

// initSceneConfigFilter 初始化场景配置表
func initSceneConfigFilter() error {
	sceneConfig := config.GetConfig().SceneConfigXvkj
	sceneConfigXvkj, err := xvkjclient.NewClient(
		xvkjclient.WithAppID(sceneConfig.AppID),
		xvkjclient.WithSchemaKey(sceneConfig.SchemaKey),
		xvkjclient.WithSchemaID(sceneConfig.SchemaID),
		xvkjclient.EnsureInitLatestData(),
		xvkjclient.EnableFilter(),
	)
	if err != nil {
		log.Errorf("initSceneConfigFilter init accessXvkj error, err:%+v", err)
		return err
	}
	sceneConfigFilter, err = sceneConfigXvkj.AddFilter([]string{"scene"}, "", model.SceneConfig{})
	if err != nil {
		log.Errorf("initSceneConfigFilter accessXvkj addfilter error, err:%+v", err)
		return err
	}
	return nil
}

// initAccessConfigFilter 初始化timeline接入配置filter
func initAccessConfigFilter() error {
	accessConfig := config.GetConfig().AccessConfigXvkj
	accessXvkj, err := xvkjclient.NewClient(
		xvkjclient.WithAppID(accessConfig.AppID),
		xvkjclient.WithSchemaKey(accessConfig.SchemaKey),
		xvkjclient.WithSchemaID(accessConfig.SchemaID),
		xvkjclient.EnsureInitLatestData(),
		xvkjclient.EnableFilter(),
	)
	if err != nil {
		log.Errorf("InitXvkj init accessXvkj error, err:%+v", err)
		return err
	}
	accessConfigFilter, err = accessXvkj.AddFilter([]string{"scene"}, "", model.AccessConfig{})
	if err != nil {
		log.Errorf("InitXvkj accessXvkj addfilter error, err:%+v", err)
		return err
	}
	return nil
}

// GetAccessConfig 获取接入配置
func GetAccessConfig(scene string) (*model.AccessConfig, error) {
	accessConfig, ok := accessConfigFilter.Get(scene).(*model.AccessConfig)
	if !ok {
		return nil, errs.Newf(errcode.XvkjAccessConfigNotExist, "GetAccessConfig not exist, scene:%s", scene)
	}
	return accessConfig, nil
}

// GetSceneConfigConfig 获取场景配置
func GetSceneConfigConfig(scene string) (*model.SceneConfig, error) {
	sceneConfig, ok := sceneConfigFilter.Get(scene).(*model.SceneConfig)
	if !ok {
		return nil, errs.Newf(errcode.XvkjSceneConfigNotExist, "GetSceneConfigConfig not exist, scene:%s", scene)
	}
	return sceneConfig, nil
}

// GetReadConfig 获取读扩散配置
func GetReadConfig(ctx context.Context, configID string) (*model.ReadConfigDTO, error) {
	readConfig := &model.ReadConfig{}
	if err := readConfigXvkj.Get(configID, readConfig); err != nil {
		return nil, errs.Newf(errcode.XvkjGetReadConfigError, "GetReadConfig error, configID:%s, err:%v", configID, err)
	}
	dto := &model.ReadConfigDTO{}
	if len(readConfig.WorksIDRoute) != 0 {
		if err := json.Unmarshal([]byte(readConfig.WorksIDRoute), &dto.WorksIDRoute); err != nil {
			log.Errorf("GetReadConfig json.Unmarshal error, WorksIDRoute:%s, err:%v", readConfig.WorksIDRoute, err)
			return nil, err
		}
	}

	if len(readConfig.RelationIDRoute) != 0 {
		if err := json.Unmarshal([]byte(readConfig.RelationIDRoute), &dto.RelationIDRoute); err != nil {
			log.Errorf("GetReadConfig json.Unmarshal error, RelationIDRoute:%s, err:%v", readConfig.RelationIDRoute, err)
			return nil, err
		}
	}

	if len(readConfig.CacheConfig) != 0 {
		if err := json.Unmarshal([]byte(readConfig.CacheConfig), &dto.CacheConfig); err != nil {
			log.Errorf("GetReadConfig json.Unmarshal error, CacheConfig:%s, err:%v", readConfig.CacheConfig, err)
			return nil, err
		}
	}

	if len(readConfig.RelationIDOrder) != 0 {
		dto.RelationIDOrder = strings.Split(readConfig.RelationIDOrder, ",")
	}
	if len(readConfig.WorksIDOrder) != 0 {
		dto.WorksIDOrder = strings.Split(readConfig.WorksIDOrder, ",")
	}
	return dto, nil
}
