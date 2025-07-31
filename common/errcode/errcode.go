package errcode

const (
	// ParamsInvalidError 参数错误
	ParamsInvalidError = -880000
	// GetWorksIDLisgrpcError 获取作品ID列表rpc接口错误
	GetWorksIDLisgrpcError = -880001
	// GetRelationIDLisgrpc 获取关系链ID列表rpc接口错误
	GetRelationIDLisgrpc = -880002
	// XvkjAccessConfigNotExist 无极配置信息不存在
	XvkjAccessConfigNotExist = -880003
	// AccessKeyError 接入appkey
	AccessKeyError = -880004
	// XvkjSceneConfigNotExist 无极场景配置信息不存在
	XvkjSceneConfigNotExist = -880005
	// XvkjSceneConfigError 无极场景信息配置有误
	XvkjSceneConfigError = -880006
	// XvkjGetReadConfigError 无极获取读扩散配置错误
	XvkjGetReadConfigError = -880007
	// XvkjRelationIDConfigError 无极关系链id配置错误
	XvkjRelationIDConfigError = -880008
	// RedisPipelineNewError redis 创建Pipeline错误
	RedisPipelineNewError = -880009
	// RedisPipelineFlushError redis Pipeline flush错误
	RedisPipelineFlushError = -880010
	// BatchGetWorksIDListAllFail 批量获取作品id列表全部失败
	BatchGetWorksIDListAllFail = -880011
	// BatchGetWorksIDListPartFail 批量获取作品id列表部分失败
	BatchGetWorksIDListPartFail = -880012
)
