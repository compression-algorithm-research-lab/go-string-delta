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
	return fmt.Sprintf("%s at [%d, %d), Content = %s", x.ChangeType.String(), x.BeginOffset, x.EndOffset, x.Content)
}

// ------------------------------------------------ ---------------------------------------------------------------------

// StringDiffer 用于比较字符串的变化
type StringDiffer interface {
	Diff(s1, s2 string) ([]*StringDiff, error)
}

// TODO 测试 
// TODO 2024-12-06 01:44:33 当有很多个diff操作的时候，尝试对多个操作合并为一个
func TryMerge(diffs []*StringDiff) []*StringDiff {
	if len(diffs) == 0 {
		return diffs
	}

	mergedDiffs := []*StringDiff{diffs[0]}
	for _, diff := range diffs[1:] {
		last := mergedDiffs[len(mergedDiffs)-1]
		if diff.ChangeType == last.ChangeType && diff.BeginOffset == last.EndOffset {
			// 如果是相同类型的操作并且是连续的，合并它们
			if diff.ChangeType == ChangeTypeInsert {
				last.Content += diff.Content
				last.EndOffset = diff.EndOffset
			} else if diff.ChangeType == ChangeTypeDelete {
				last.EndOffset = diff.EndOffset
			} else if diff.ChangeType == ChangeTypeReplace {
				last.Content = diff.Content
				last.EndOffset = diff.EndOffset
			}
		} else {
			// 如果不能合并，添加到结果中
			mergedDiffs = append(mergedDiffs, diff)
		}
	}
	return mergedDiffs
}

// TODO 测试
// From 在一个给定的字符串上应用diff，得到另一个字符串
func From(s string, diffs []*StringDiff) (string, error) {
	result := s
	for _, diff := range diffs {
		switch diff.ChangeType {
		case ChangeTypeInsert:
			result = result[:int(diff.BeginOffset)] + diff.Content + result[int(diff.BeginOffset):]
		case ChangeTypeDelete:
			result = result[:int(diff.BeginOffset)] + result[int(diff.EndOffset):]
		case ChangeTypeReplace:
			result = result[:int(diff.BeginOffset)] + diff.Content + result[int(diff.EndOffset):]
		case ChangeTypeNone:
			// 无操作，继续下一个
			continue
		default:
			return "", fmt.Errorf("unknown change type: %v", diff.ChangeType)
		}
	}
	return result, nil
}

// ------------------------------------------------ ---------------------------------------------------------------------
