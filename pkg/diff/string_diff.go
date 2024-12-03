package diff

// ------------------------------------------------ ---------------------------------------------------------------------

// ChangeType 用于表示变化类型
type ChangeType int

const (

	// ChangeTypeNone 没有变化
	ChangeTypeNone ChangeType = iota

	// ChangeTypeInsert 插入了内容
	ChangeTypeInsert ChangeType = iota

	// ChangeTypeDelete 删除了内容
	ChangeTypeDelete ChangeType = iota

	// ChangeTypeReplace 修改了内容
	ChangeTypeReplace ChangeType = iota
)

// ------------------------------------------------ ---------------------------------------------------------------------

// StringDiff 标识一个变化
type StringDiff struct {

	// 变化的类型
	ChangeType ChangeType

	// 变化的开始位置
	BeginOffset uint64

	// 变化的结束位置
	EndOffset uint64

	// 变化的内容
	Content string
}

// ------------------------------------------------ ---------------------------------------------------------------------

// StringDiffer 用于比较字符串的变化
type StringDiffer interface {
	Diff(s1, s2 string) ([]*StringDiff, error)
}

// ------------------------------------------------ ---------------------------------------------------------------------