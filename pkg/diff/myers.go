package diff

// MyersStringDiffer 实现了 StringDiffer 接口
type MyersStringDiffer struct{}

// Diff 使用 Myers 差分算法比较两个字符串
func (d *MyersStringDiffer) Diff(s1, s2 string) ([]*StringDiff, error) {
	editScript := myersDiff(s1, s2)
	if editScript == nil {
		return nil, nil // 没有变化
	}

	var diffs []*StringDiff
	offset1, offset2 := uint64(0), uint64(0)
	for _, op := range editScript {
		switch op {
		case OP_INSERT:
			diffs = append(diffs, &StringDiff{
				ChangeType:  ChangeTypeInsert,
				BeginOffset: offset1,
				EndOffset:   offset1,
				Content:     string(s2[offset2]),
			})
			offset2++
		case OP_DELETE:
			diffs = append(diffs, &StringDiff{
				ChangeType:  ChangeTypeDelete,
				BeginOffset: offset1,
				EndOffset:   offset1 + 1,
				Content:     string(s1[offset1]),
			})
			offset1++
		case OP_REPLACE:
			diffs = append(diffs, &StringDiff{
				ChangeType:  ChangeTypeReplace,
				BeginOffset: offset1,
				EndOffset:   offset1 + 1,
				Content:     string(s2[offset2]),
			})
			offset1++
			offset2++
		case OP_MATCH:
			offset1++
			offset2++
		}
	}

	return diffs, nil
}

// 定义编辑操作类型
const (
	OP_MATCH = iota
	OP_INSERT
	OP_DELETE
	OP_REPLACE
)

// myersDiff 执行 Myers 差分算法
func myersDiff(a, b string) []int {
	n, m := len(a), len(b)
	max := n + m
	offsets := make([]int, 2*max+1)
	for i := range offsets {
		offsets[i] = -1
	}
	offsets[max+1] = 0

	var snake, x, y int
	for d := 0; d <= max; d++ {
		for k := -d; k <= d; k += 2 {
			if k == -d || (k != d && offsets[max+k-1] < offsets[max+k+1]) {
				x = offsets[max+k+1]
			} else {
				x = offsets[max+k-1] + 1
			}
			y = x - k

			for x < n && y < m && a[x] == b[y] {
				x++
				y++
			}

			offsets[max+k] = x
			if x >= n && y >= m {
				snake = k
				goto end
			}
		}
	}

end:
	if offsets[max+snake] == -1 {
		return nil // 没有变化
	}

	editScript := make([]int, 0, max)
	x, y = n, m
	for k := snake; k != 0; {
		if offsets[max+k-1] < offsets[max+k+1] {
			editScript = append(editScript, OP_DELETE)
			k--
			x--
		} else {
			editScript = append(editScript, OP_INSERT)
			k++
			y--
		}
	}
	for i := x - 1; i >= 0; i-- {
		if a[i] == b[i] {
			editScript = append(editScript, OP_MATCH)
		} else {
			editScript = append(editScript, OP_REPLACE)
		}
	}

	// Reverse the edit script to get the correct order
	for i, j := 0, len(editScript)-1; i < j; i, j = i+1, j-1 {
		editScript[i], editScript[j] = editScript[j], editScript[i]
	}

	return editScript
}
