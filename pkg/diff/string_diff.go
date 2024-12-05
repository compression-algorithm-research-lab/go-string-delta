package diff

import "fmt"

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

func (x ChangeType) String() string {
	switch x {
	case ChangeTypeNone:
		return "None"
	case ChangeTypeInsert:
		return "Insert"
	case ChangeTypeDelete:
		return "Delete"
	case ChangeTypeReplace:
		return "Replace"
	default:
		return ""
	}
}

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

func (x *StringDiff) String() string {
	return fmt.Sprintf("%s at [%d, %d], Content = %s", x.ChangeType.String(), x.BeginOffset, x.EndOffset, x.Content)
}

// ------------------------------------------------ ---------------------------------------------------------------------

// StringDiffer 用于比较字符串的变化
type StringDiffer interface {
	Diff(s1, s2 string) ([]*StringDiff, error)
}

// ------------------------------------------------ ---------------------------------------------------------------------
