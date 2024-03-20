package rollinghash

import (
	"testing"
)

func TestRollHash(t *testing.T) {
	rh := NewRollingHash(3).(*RollingHash)
	testData := []struct {
		input     byte
		expected  byte
		hashValue uint32
	}{
		{input: 'a', expected: 0, hashValue: 97},
		{input: 'b', expected: 0, hashValue: 783},
		{input: 'c', expected: 0, hashValue: 5634},
		{input: 'd', expected: 'a', hashValue: 5691},
		{input: 'e', expected: 'b', hashValue: 5748},
	}

	for i, tt := range testData {
		output := rh.RollHash(tt.input)
		if output != tt.expected {
			t.Errorf("Test %d: RollHash(%c) = %c; want %c", i+1, tt.input, output, tt.expected)
		}
		if rh.Signature() != tt.hashValue {
			t.Errorf("Test %d: Hash value after rolling is incorrect. Got: %d, Expected: %d", i+1, rh.Signature(), tt.hashValue)
		}
	}
}
