package internal

// RowResult 用于返回修改成功的行数
type RowResult struct {
	// 行数
	Row int64 `json:"row"`
}

// IDResult 用于返回添加成功的数据库 ID
type IDResult[T any] struct {
	// 数据库 ID
	ID T `json:"id"`
}

// IDPath 用于 binding 路径中的 id
type IDPath[T any] struct {
	// 数据库 ID
	ID T `uri:"id" binding:"required,min=1"`
}

// BatchDelete 用于 绑定删除的条件
type BatchDelete[T any] struct {
	// 数据库 ID 数组
	ID []T `json:"id" binding:"required"`
}
