package main

import (
	"fmt"
	"testing"

	"github.com/rapl27/eDiff/delta"
)

const (
	t1_old_file = "./testData/test1_old.txt"
	t1_new_file = "./testData/test1_new.txt"
	t2_old_file = "./testData/test2_old.txt"
	t2_new_file = "./testData/test2_new.txt"
)

func TestParseArgs(t *testing.T) {
	tests := []struct {
		name          string
		old           string
		new           string
		expectedDelta []delta.Delta
	}{
		{
			name: "TestDeltaTC1",
			old:  t1_old_file,
			new:  t1_new_file,
			expectedDelta: []delta.Delta{
				{
					Offset:    1,
					Operation: "U",
				},
				{
					Offset:    2,
					Operation: "U",
				},
				{
					Offset:    3,
					Operation: "M",
					Data:      []byte("e  upd"),
				},
				{
					Offset:    4,
					Operation: "M",
					Data:      []byte("ated "),
				},
				{
					Offset:    5,
					Operation: "U",
				},
			},
		},
		{
			name: "TestDeltaTC2",
			old:  t2_old_file,
			new:  t2_new_file,
			expectedDelta: []delta.Delta{
				{
					Offset:    1,
					Operation: "U",
				},
				{
					Offset:    2,
					Operation: "U",
				},
				{
					Offset:    3,
					Operation: "R",
				},
				{
					Offset:    4,
					Operation: "U",
				},
				{
					Offset:    5,
					Operation: "I",
					Data:      []byte("word3"),
				},
				{
					Offset:    5,
					Operation: "I",
					Data:      []byte("word5"),
				},
				{
					Offset:    5,
					Operation: "U",
				},
				{
					Offset:    6,
					Operation: "U",
				},
			},
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			fmt.Println(tc.name)
			fmt.Println("-----------")
			delta := runDelta(tc.old, tc.new, 5)

			for i, d := range delta {
				if d.Offset != tc.expectedDelta[i].Offset || d.Operation != tc.expectedDelta[i].Operation {
					t.Errorf("Test %s failed: want %v", tc.name, tc.expectedDelta[i])
				}
			}

		})
	}
}
