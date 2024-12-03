package main

import (
	"fmt"
)

type DiffOp int

const (
	Equal DiffOp = iota
	Add
	Remove
)

type Diff struct {
	Op   DiffOp
	Text string
}

func myersDiff(a, b string) []Diff {
	m, n := len(a), len(b)
	max := m + n
	v := make([]int, 2*max+1)
	v[1] = 0

	var diffs []Diff

	for d := 0; d <= max; d++ {
		for k := -d; k <= d; k += 2 {
			var x int
			if k == -d || (k != d && v[k-1+max] < v[k+1+max]) {
				x = v[k+1+max]
			} else {
				x = v[k-1+max] + 1
			}
			y := x - k
			for x < m && y < n && a[x] == b[y] {
				x++
				y++
			}
			v[k+max] = x
			if x >= m && y >= n {
				k = 0
				x, y = m, n
				for d := d; d > 0; d-- {
					if x > 0 && v[k-1+max] < x {
						k--
					} else {
						k++
					}
					oldX := v[k+max]
					for x > oldX {
						diffs = append([]Diff{{Remove, string(a[x-1])}}, diffs...)
						x--
					}
					for y > x-k {
						diffs = append([]Diff{{Add, string(b[y-1])}}, diffs...)
						y--
					}
					if x > 0 && y > 0 {
						diffs = append([]Diff{{Equal, string(a[x-1])}}, diffs...)
					}
					x--
					y--
				}
				return diffs
			}
		}
	}

	return diffs
}

func main() {
	s1 := "abcdefghijklmnop"
	s2 := "abcdefghiyklmnop"
	diffs := myersDiff(s1, s2)
	for _, diff := range diffs {
		switch diff.Op {
		case Equal:
			fmt.Printf("  %s\n", diff.Text)
		case Add:
			fmt.Printf("+ %s\n", diff.Text)
		case Remove:
			fmt.Printf("- %s\n", diff.Text)
		}
	}
}
