package diff

import (
	"fmt"
	"testing"
)

func TestMyersStringDiffer_Diff(t *testing.T) {
	differ := &MyersStringDiffer{}
	diffs, err := differ.Diff("ABCABBA", "CBABAC")
	if err != nil {
		// 处理错误
	}
	for _, diff := range diffs {
		fmt.Printf("%v: %d, %d, %s\n", diff.ChangeType.String(), diff.BeginOffset, diff.EndOffset, diff.Content)
	}
}
