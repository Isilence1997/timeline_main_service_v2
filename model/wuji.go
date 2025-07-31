package model

import (
	"fmt"
	"time"

	"git.code.oa.com/grpc-go/grpc-go/client"
	"git.code.oa.com/v/main_logic/feeds/grpc_timeline_main_service_v2/common/constant"
	"git.code.oa.com/v/main_logic/feeds/grpc_timeline_main_service_v2/common/retry"
)

// AccessConfig timeline 接入配置
// http://xvkj.oa.com/p/edit?appid=timeline_config&schemaid=access_config_test
type AccessConfig struct {
	Appkey       string `json:"appkey"`
	AppkeyEnable bool   `json:"appkey_enable"`
	Scene        string `json:"scene"`
}

// SceneConfig 场景配置表
// http://xvkj.oa.com/p/edit?appid=timeline_config&schemaid=sence_config_test
type SceneConfig struct {
	Scene         string `json:"scene"`
	DiffusionType string `json:"diffusion_type"`
	SceneType     string `json:"scene_type"`
	PageSize      int    `json:"page_size"`
	ReadConfigID  string `json:"read_config_id"`
	WriteConfigID string `json:"write_config_id"`
}

// ReadConfig 读扩散配置
// http://xvkj.oa.com/p/edit?appid=timeline_config&schemaid=read_config_test
type ReadConfig struct {
	WorksIDRoute    string `json:"works_id_route"`
	RelationIDRoute string `json:"relation_id_route"`
	CacheConfig     string `json:"cache_config"`
	RelationIDOrder string `json:"relation_id_order"`
	WorksIDOrder    string `json:"works_id_order"`
}

// ReadConfigDTO 读扩散配置dto
type ReadConfigDTO struct {
	WorksIDRoute    *WorksIDRouteDTO        `json:"works_id_route"`
	RelationIDRoute *RelationIDRouteDTO     `json:"relation_id_route"`
	CacheConfig     map[string]CacheInfoDTO `json:"cache_config"`
	RelationIDOrder []string                `json:"relation_id_order"`
	WorksIDOrder    []string                `json:"works_id_order"`
}

// WorksIDRouteDTO 作品id列表路由信息
type WorksIDRouteDTO struct {
	Order     int          `json:"order"`
	PageSize  int64        `json:"page_size"`
	OrderName string       `json:"order_name"`
	RouteInfo RouteInfoDTO `json:"route_info"`
}

// RelationIDRouteDTO 关系链id列表路由
type RelationIDRouteDTO struct {
	Order     int          `json:"order"`
	Mode      string       `json:"mode"`
	PageSize  int64        `json:"page_size"`
	OrderName string       `json:"order_name"`
	RouteInfo RouteInfoDTO `json:"route_info"`
}

// SortCacheConfig 排序配置信息
type SortCacheConfig struct {
	Order []string                `json:"order"`
	Cache map[string]CacheInfoDTO `json:"cache"`
}

// CacheInfoDTO 缓存信息
type CacheInfoDTO struct {
	KeyConfig        *CacheKeyConfigDTO `json:"key_config"`
	ReadCacheRoute   *RouteInfoDTO      `json:"read_cache"`
	WriteCache       *RouteInfoDTO      `json:"write_cache"`
	MemoryCache      *MemoryCacheDTO    `json:"memory_cache"`
	NeedBackSource   bool               `json:"need_back_source"`
	BackSourceConfig *KafkaProducerDTO  `json:"back_source_config"`
}

// KafkaProducerDTO kafka路由信息
type KafkaProducerDTO struct {
	ServiceName string `json:"service_name"`
	Address     string `json:"address"`
	Topic       string `json:"topic"`
	ClientID    string `json:"client_id"`
	Timeout     int    `json:"timeout"`
}

// KafkaProducerInfo2Options kafka生产者信息转换为options
func (p KafkaProducerDTO) KafkaProducerInfo2Options() []client.Option {
	target := fmt.Sprintf(`kafka://%s?topic=%s&clientid=%s&compression=none`, p.Address, p.Topic, p.ClientID)
	timeout := constant.DefaultKafkaProducerTimeout
	if p.Timeout > 0 {
		timeout = time.Duration(p.Timeout) * time.Millisecond
	}

	res := []client.Option{
		client.WithTarget(target),
		client.WithTimeout(timeout),
	}
	return res
}

// GetAddress 获取地址
func (p KafkaProducerDTO) GetAddress() string {
	address := fmt.Sprintf(`kafka://%s?topic=%s&clientid=%s&compression=none`, p.Address, p.Topic, p.ClientID)
	return address
}

// GetTimeout 获取超时时间
func (p KafkaProducerDTO) GetTimeout() time.Duration {
	if p.Timeout <= 0 {
		return constant.DefaultKafkaProducerTimeout
	}
	return time.Duration(p.Timeout) * time.Millisecond
}

// CacheKeyConfigDTO 缓存key配置信息
type CacheKeyConfigDTO struct {
	KeyPrefix string `json:"key_prefix"`
	//KeyMod    string `json:"key_mod"`
	//KeyHash   string `json:"key_hash"`
	//ValueType string `json:"value_type"`
}

// MemoryCacheDTO 内存缓存DTO信息
type MemoryCacheDTO struct {
	Enable   bool `json:"enable"`
	Expire   int  `json:"expire"`
	Capacity int  `json:"capacity"`
}

// GetExpire 获取过期时间
func (p MemoryCacheDTO) GetExpire() time.Duration {
	if p.Expire <= 0 {
		return constant.DefaultMemoryCacheExpire
	}
	return time.Duration(p.Expire) * time.Second
}

// GetCapacity 获取缓存大小
func (p MemoryCacheDTO) GetCapacity() int {
	if p.Capacity <= 0 {
		return constant.DefaultMemoryCacheCapacity
	}
	return p.Capacity
}

// RouteInfoDTO 路由信息
type RouteInfoDTO struct {
	ServiceName string `json:"service_name"`
	Target      string `json:"target"`
	Namespace   string `json:"namespace"`
	Set         string `json:"set"`
	Retry       int    `json:"retry"`
	Timeout     int    `json:"timeout"`
	Password    string `json:"password"`
}

// RouteInfo2Options 路由信息转换为option切片
func (routeInfo RouteInfoDTO) RouteInfo2Options() []client.Option {
	res := make([]client.Option, 0)
	// 路由调用
	if len(routeInfo.Set) == 0 {
		res = append(res,
			client.WithTarget(routeInfo.Target),
			client.WithDisableServiceRouter(),
		)
	} else {
		res = append(res,
			client.WithServiceName(routeInfo.ServiceName),
			client.WithCalleeSetName(routeInfo.Set),
		)
	}

	// 密码主要是redis用到
	if len(routeInfo.Password) != 0 {
		res = append(res, client.WithPassword(routeInfo.Password))
	}

	// 路由超时
	if routeInfo.Timeout <= 0 {
		res = append(res, client.WithTimeout(constant.DefaultRouteTimeout))
	} else {
		res = append(res, client.WithTimeout(time.Duration(routeInfo.Timeout)*time.Millisecond))
	}
	res = append(res,
		client.WithNamespace(routeInfo.Namespace),
		client.WithFilter(retry.GetNetworkErrRetry(routeInfo.Retry).Invoke),
	)
	return res
}
