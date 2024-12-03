package diff

type LcsStringDiffer struct {
}

var _ StringDiffer = &LcsStringDiffer{}

func (x *LcsStringDiffer) Diff(s1, s2 string) ([]*StringDiff, error) {
	m := len(s1)
	s := len(s2)
	lcsMatrix := make([][]int, m+1)
	for i := range lcsMatrix {
		lcsMatrix[i] = make([]int, s+1)
	}

	for i := 1; i <= m; i++ {
		for j := 1; j <= s; j++ {
			if s1[i-1] == s2[j-1] {
				lcsMatrix[i][j] = lcsMatrix[i-1][j-1] + 1
			} else {
				if lcsMatrix[i-1][j] > lcsMatrix[i][j-1] {
					lcsMatrix[i][j] = lcsMatrix[i-1][j]
				} else {
					lcsMatrix[i][j] = lcsMatrix[i][j-1]
				}
			}
		}
	}

	return nil, nil
}
