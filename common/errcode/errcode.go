package errcode

const (
	// ParamsInvalidError 参数错误
	ParamsInvalidError = -880000
	// GetWorksIDListRpcError 获取作品ID列表rpc接口错误
	GetWorksIDListRpcError = -880001
	// GetRelationIDListRpc 获取关系链ID列表rpc接口错误
	GetRelationIDListRpc = -880002
	// WujiAccessConfigNotExist 无极配置信息不存在
	WujiAccessConfigNotExist = -880003
	// AccessKeyError 接入appkey
	AccessKeyError = -880004
	// WujiSceneConfigNotExist 无极场景配置信息不存在
	WujiSceneConfigNotExist = -880005
	// WujiSceneConfigError 无极场景信息配置有误
	WujiSceneConfigError = -880006
	// WujiGetReadConfigError 无极获取读扩散配置错误
	WujiGetReadConfigError = -880007
	// WujiRelationIDConfigError 无极关系链id配置错误
	WujiRelationIDConfigError = -880008
	// RedisPipelineNewError redis 创建Pipeline错误
	RedisPipelineNewError = -880009
	// RedisPipelineFlushError redis Pipeline flush错误
	RedisPipelineFlushError = -880010
	// BatchGetWorksIDListAllFail 批量获取作品id列表全部失败
	BatchGetWorksIDListAllFail = -880011
	// BatchGetWorksIDListPartFail 批量获取作品id列表部分失败
	BatchGetWorksIDListPartFail = -880012
)
