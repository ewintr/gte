package storage_test

import (
	"testing"

	"git.ewintr.nl/go-kit/test"
	"git.ewintr.nl/gte/internal/storage"
)

func TestNextLocalId(t *testing.T) {
	for _, tc := range []struct {
		name string
		used []int
		exp  int
	}{
		{
			name: "empty",
			used: []int{},
			exp:  1,
		},
		{
			name: "not empty",
			used: []int{5},
			exp:  6,
		},
		{
			name: "multiple",
			used: []int{2, 3, 4},
			exp:  5,
		},
		{
			name: "holes",
			used: []int{1, 5, 8},
			exp:  9,
		},
		{
			name: "expand limit",
			used: []int{1, 2, 3, 4, 5, 6, 7, 8, 9, 10},
			exp:  11,
		},
		{
			name: "wrap if possible",
			used: []int{8, 9},
			exp:  1,
		},
		{
			name: "find hole",
			used: []int{1, 2, 3, 4, 5, 7, 8, 9},
			exp:  6,
		},
		{
			name: "dont wrap if expanded before",
			used: []int{15, 16},
			exp:  17,
		},
		{
			name: "do wrap if expanded limit is reached",
			used: []int{99},
			exp:  1,
		},
	} {
		t.Run(tc.name, func(t *testing.T) {
			test.Equals(t, tc.exp, storage.NextLocalId(tc.used))
		})
	}
}
