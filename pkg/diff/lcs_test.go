package diff

import (
	"fmt"
	"testing"
)

func TestLCSStringDiffer_Diff(t *testing.T) {
	differ := NewLCSStringDiffer()
	diffs, err := differ.Diff("ABCABBA", "CBABAC")
	if err != nil {
		// 处理错误
	}
	for _, diff := range diffs {
		fmt.Println(diff.String())
	}
}
