package diff

import (
	"testing"
)

func TestLcsStringDiffer_Diff(t *testing.T) {
	s1 := "abcdefghijklmnop"
	s2 := "abcdefghiyklmnop"
	differ := LcsStringDiffer{}
	diff, err := differ.Diff(s1, s2)
	if err != nil {
		panic(err)
	}
	print(diff)
}
