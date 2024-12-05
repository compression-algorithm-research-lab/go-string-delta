package diff

// LCSStringDiffer 实现了 StringDiffer 接口
type LCSStringDiffer struct{}

// NewLCSStringDiffer 创建 LCSStringDiffer 实例
func NewLCSStringDiffer() *LCSStringDiffer {
	return &LCSStringDiffer{}
}

// Diff 使用 LCS 算法比较两个字符串并返回差异
func (d *LCSStringDiffer) Diff(s1, s2 string) ([]*StringDiff, error) {
	if s1 == s2 {
		return nil, nil // 没有差异
	}

	// 计算 LCS 矩阵
	l1, l2 := len(s1), len(s2)
	lcsMatrix := make([][]int, l1+1)
	for i := range lcsMatrix {
		lcsMatrix[i] = make([]int, l2+1)
	}

	for i := 1; i <= l1; i++ {
		for j := 1; j <= l2; j++ {
			if s1[i-1] == s2[j-1] {
				lcsMatrix[i][j] = lcsMatrix[i-1][j-1] + 1
			} else {
				lcsMatrix[i][j] = max(lcsMatrix[i-1][j], lcsMatrix[i][j-1])
			}
		}
	}

	// 从 LCS 矩阵中回溯差异
	var diffs []*StringDiff
	i, j := l1, l2
	for i > 0 && j > 0 {
		if s1[i-1] == s2[j-1] {
			if i == l1 || j == l2 || s1[i] != s2[j] {
				// 当遇到第一个不同的字符时，标记之前的内容为无变化
				if len(diffs) == 0 || diffs[len(diffs)-1].ChangeType != ChangeTypeNone {
					diffs = append(diffs, &StringDiff{
						ChangeType:  ChangeTypeNone,
						BeginOffset: uint64(i),
						EndOffset:   uint64(i + 1),
						Content:     string(s1[i-1]),
					})
				} else {
					// 如果上一个差异是无变化，则合并内容
					diffs[len(diffs)-1].EndOffset++
					diffs[len(diffs)-1].Content += string(s1[i-1])
				}
			}
			i--
			j--
		} else if lcsMatrix[i][j-1] > lcsMatrix[i-1][j] {
			// 插入差异
			diffs = append(diffs, &StringDiff{
				ChangeType:  ChangeTypeInsert,
				BeginOffset: uint64(i),
				EndOffset:   uint64(i),
				Content:     string(s2[j-1]),
			})
			j--
		} else {
			// 删除差异
			diffs = append(diffs, &StringDiff{
				ChangeType:  ChangeTypeDelete,
				BeginOffset: uint64(i),
				EndOffset:   uint64(i + 1),
				Content:     string(s1[i-1]),
			})
			i--
		}
	}

	// 处理剩余的差异
	for i > 0 {
		diffs = append(diffs, &StringDiff{
			ChangeType:  ChangeTypeDelete,
			BeginOffset: uint64(i),
			EndOffset:   uint64(i + 1),
			Content:     string(s1[i-1]),
		})
		i--
	}
	for j > 0 {
		diffs = append(diffs, &StringDiff{
			ChangeType:  ChangeTypeInsert,
			BeginOffset: uint64(l1),
			EndOffset:   uint64(l1),
			Content:     string(s2[j-1]),
		})
		j--
	}

	// 反转差异数组以获得正确的顺序
	for i, j := 0, len(diffs)-1; i < j; i, j = i+1, j-1 {
		diffs[i], diffs[j] = diffs[j], diffs[i]
	}

	return diffs, nil
}

// 辅助函数，用于获取两个整数中的最大值
func max(a, b int) int {
	if a > b {
		return a
	}
	return b
}
