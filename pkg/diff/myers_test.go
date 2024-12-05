package diff

import (
	"fmt"
	"testing"
)

func TestMyersStringDiffer_Diff(t *testing.T) {
	differ := &MyersStringDiffer{}
	diffs, err := differ.Diff("AAA", "BBB")
	if err != nil {
		// 处理错误
	}
	for _, diff := range diffs {
		fmt.Println(diff.String())
	}
}
