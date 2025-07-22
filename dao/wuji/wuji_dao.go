package wuji

import (
	"context"
	"encoding/json"
	"strings"

	"git.code.oa.com/open-wuji/go-sdk/wujiclient"
	trpc "git.code.oa.com/trpc-go/trpc-go"
	"git.code.oa.com/trpc-go/trpc-go/errs"
	"git.code.oa.com/trpc-go/trpc-go/log"
	"git.code.oa.com/v/main_logic/feeds/trpc_timeline_main_service_v2/common/errcode"
	"git.code.oa.com/v/main_logic/feeds/trpc_timeline_main_service_v2/config"
	"git.code.oa.com/v/main_logic/feeds/trpc_timeline_main_service_v2/model"
)

var (
	accessConfigFilter wujiclient.FilterInterface
	sceneConfigFilter  wujiclient.FilterInterface
	readConfigWuji     wujiclient.ClientInterface
)

// InitWuJi 初始化wuji配置
func InitWuJi() error {
	err := trpc.GoAndWait(func() error {
		return initAccessConfigFilter()
	}, func() error {
		return initSceneConfigFilter()
	}, func() error {
		return initReadConfigWuji()
	})
	return err
}

// initReadConfigWuji 初始化读扩散配置表
func initReadConfigWuji() error {
	readConfig := config.GetConfig().ReadConfigWuji
	readConfigWujiTmp, err := wujiclient.NewClient(
		wujiclient.WithAppID(readConfig.AppID),
		wujiclient.WithSchemaKey(readConfig.SchemaKey),
		wujiclient.WithSchemaID(readConfig.SchemaID),
		wujiclient.EnsureInitLatestData(),
	)
	if err != nil {
		log.Errorf("initReadConfigWuji init accessWuji error, err:%+v", err)
		return err
	}
	readConfigWuji = readConfigWujiTmp
	return nil
}

// initSceneConfigFilter 初始化场景配置表
func initSceneConfigFilter() error {
	sceneConfig := config.GetConfig().SceneConfigWuji
	sceneConfigWuji, err := wujiclient.NewClient(
		wujiclient.WithAppID(sceneConfig.AppID),
		wujiclient.WithSchemaKey(sceneConfig.SchemaKey),
		wujiclient.WithSchemaID(sceneConfig.SchemaID),
		wujiclient.EnsureInitLatestData(),
		wujiclient.EnableFilter(),
	)
	if err != nil {
		log.Errorf("initSceneConfigFilter init accessWuji error, err:%+v", err)
		return err
	}
	sceneConfigFilter, err = sceneConfigWuji.AddFilter([]string{"scene"}, "", model.SceneConfig{})
	if err != nil {
		log.Errorf("initSceneConfigFilter accessWuji addfilter error, err:%+v", err)
		return err
	}
	return nil
}

// initAccessConfigFilter 初始化timeline接入配置filter
func initAccessConfigFilter() error {
	accessConfig := config.GetConfig().AccessConfigWuji
	accessWuji, err := wujiclient.NewClient(
		wujiclient.WithAppID(accessConfig.AppID),
		wujiclient.WithSchemaKey(accessConfig.SchemaKey),
		wujiclient.WithSchemaID(accessConfig.SchemaID),
		wujiclient.EnsureInitLatestData(),
		wujiclient.EnableFilter(),
	)
	if err != nil {
		log.Errorf("InitWuJi init accessWuji error, err:%+v", err)
		return err
	}
	accessConfigFilter, err = accessWuji.AddFilter([]string{"scene"}, "", model.AccessConfig{})
	if err != nil {
		log.Errorf("InitWuJi accessWuji addfilter error, err:%+v", err)
		return err
	}
	return nil
}

// GetAccessConfig 获取接入配置
func GetAccessConfig(scene string) (*model.AccessConfig, error) {
	accessConfig, ok := accessConfigFilter.Get(scene).(*model.AccessConfig)
	if !ok {
		return nil, errs.Newf(errcode.WujiAccessConfigNotExist, "GetAccessConfig not exist, scene:%s", scene)
	}
	return accessConfig, nil
}

// GetSceneConfigConfig 获取场景配置
func GetSceneConfigConfig(scene string) (*model.SceneConfig, error) {
	sceneConfig, ok := sceneConfigFilter.Get(scene).(*model.SceneConfig)
	if !ok {
		return nil, errs.Newf(errcode.WujiSceneConfigNotExist, "GetSceneConfigConfig not exist, scene:%s", scene)
	}
	return sceneConfig, nil
}

// GetReadConfig 获取读扩散配置
func GetReadConfig(ctx context.Context, configID string) (*model.ReadConfigDTO, error) {
	readConfig := &model.ReadConfig{}
	if err := readConfigWuji.Get(configID, readConfig); err != nil {
		return nil, errs.Newf(errcode.WujiGetReadConfigError, "GetReadConfig error, configID:%s, err:%v", configID, err)
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
