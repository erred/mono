package ring

import (
	"testing"

	"golang.org/x/exp/slices"
)

func TestRingPutDump(t *testing.T) {
	tcs := []struct {
		in  []int
		out []int
	}{
		{
			[]int{1, 2, 3},
			[]int{1, 2, 3},
		}, {
			[]int{1, 2, 3, 4, 5},
			[]int{1, 2, 3, 4, 5},
		}, {
			[]int{1, 2, 3, 4, 5, 6, 7},
			[]int{3, 4, 5, 6, 7},
		},
	}

	for i, tc := range tcs {
		ring := Of[int](5)
		for _, v := range tc.in {
			ring.Put(v)
		}
		out := ring.Dump()
		if !slices.Equal(out, tc.out) {
			t.Logf("ring = %v; head = %v tail = %v", ring.ring, ring.head, ring.tail)
			t.Errorf("Dump(%d) = %v; want %v", i, out, tc.out)
		}
	}
}

func TestRingGet(t *testing.T) {
	tcs := []struct {
		in  []int
		out []int
	}{
		{
			[]int{1, 2, 3},
			[]int{1, 2},
		}, {
			[]int{1, 2, 3},
			[]int{1, 2, 3, 0},
		}, {
			[]int{1, 2, 3, 4, 5, 6, 7},
			[]int{3, 4, 5, 6, 7, 0},
		},
	}
	for i, tc := range tcs {
		ring := Of[int](5)
		for _, v := range tc.in {
			ring.Put(v)
		}
		for j, v := range tc.out {
			out, ok := ring.Get()
			if out != v || (v == 0 && ok != false) {
				t.Logf("ring = %v; head = %v tail = %v", ring.ring, ring.head, ring.tail)
				t.Errorf("Get(%d,%d) = %v; want %v", i, j, out, v)
			}
		}
	}
}
