package config

import (
	"encoding/json"

	"git.code.oa.com/grpc-go/grpc-go/config"
	"git.code.oa.com/grpc-go/grpc-go/log"
)

var (
	// serviceConfig 配置信息对象
	serviceConfig ServiceConfig
)

// ServiceConfig 配置信息
type ServiceConfig struct {
	AccessConfigXvkj struct {
		AppID     string `json:"app_id" yaml:"app_id"`
		SchemaKey string `json:"schema_key" yaml:"schema_key"`
		SchemaID  string `json:"schema_id" yaml:"schema_id"`
	} `json:"access_config_xvkj" yaml:"access_config_xvkj"`
	SceneConfigXvkj struct {
		AppID     string `json:"app_id" yaml:"app_id"`
		SchemaKey string `json:"schema_key" yaml:"schema_key"`
		SchemaID  string `json:"schema_id" yaml:"schema_id"`
	} `json:"scene_config_xvkj" yaml:"scene_config_xvkj"`
	ReadConfigXvkj struct {
		AppID     string `json:"app_id" yaml:"app_id"`
		SchemaKey string `json:"schema_key" yaml:"schema_key"`
		SchemaID  string `json:"schema_id" yaml:"schema_id"`
	} `json:"read_config_xvkj" yaml:"read_config_xvkj"`
}

// InitServiceConfig 初始化服务配置
func InitServiceConfig() error {
	// 加载配置文件
	confName := "timeline_main_service_v2.yaml"
	serviceConfig = ServiceConfig{}
	err := config.GetYAML(confName, &serviceConfig)
	if err != nil {
		log.Errorf("get yaml conf error,err:%v", err)
		return err
	}
	marshal, _ := json.Marshal(serviceConfig)
	log.Infof("yaml conf, conf:%+v", string(marshal))
	return nil
}

// GetConfig 获取配置
func GetConfig() ServiceConfig {
	return serviceConfig
}
