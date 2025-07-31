package constant

// DiffusionType 扩散类型
type DiffusionType = string

const (
	// Write 写扩散
	Write DiffusionType = "write"
	// Read 读扩散
	Read DiffusionType = "read"
)

// SceneType 场景类型
type SceneType = string

const (
	// Feed feed流
	Feed SceneType = "feed"
	// AvatarList 头像列表
	AvatarList SceneType = "avatar_list"
	// Refresh 下拉刷新
	Refresh SceneType = "refresh"
)

// RelationIDMode 关系链id拉取模式
type RelationIDMode = string

const (
	// Offset 偏移量拉取
	Offset RelationIDMode = "offset"
	// PageContext pageContext分页拉取
	PageContext RelationIDMode = "pageContext"
)
