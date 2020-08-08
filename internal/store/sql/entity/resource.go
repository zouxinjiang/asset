package entity

type (
	ResourceType    int64
	ResourceRelType string
	Resource        struct {
		ID   int64
		Type ResourceType
	}

	ResourceResourceRel struct {
		ID       int64
		ParentID int64
		ChildID  int64
		RelType  ResourceRelType
	}
)

const (
	// 顺序不能动
	// 不属于任何类型。 作为代码中的特殊值使用
	ResourceUnknown ResourceType = iota
	ResourceUser
)

var (
	resourceTypeName = map[ResourceType]string{
		ResourceUser: "User",
	}
)

func (t ResourceType) String() string {
	return resourceTypeName[t]
}
