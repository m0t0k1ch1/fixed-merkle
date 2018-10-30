package merkle

import (
	"crypto/sha256"
	"hash"
	"testing"
)

func TestNewConfig(t *testing.T) {
	type input struct {
		hasher   hash.Hash
		depth    uint64
		hashSize uint64
	}
	type output struct {
		config *Config
		err    error
	}
	testCases := []struct {
		name string
		in   input
		out  output
	}{
		{
			"success: min size",
			input{sha256.New(), 1, 1},
			output{&Config{allLeavesNum: 2, allNodesNum: 3}, nil},
		},
		{
			"success: max size",
			input{sha256.New(), 16, 64},
			output{&Config{allLeavesNum: 65536, allNodesNum: 131071}, nil},
		},
		{
			"failure: too small depth",
			input{sha256.New(), 0, 32},
			output{nil, ErrTooSmallDepth},
		},
		{
			"failure: too large depth",
			input{sha256.New(), 17, 32},
			output{nil, ErrTooLargeDepth},
		},
		{
			"failure: too small hash size",
			input{sha256.New(), 8, 0},
			output{nil, ErrTooSmallHashSize},
		},
		{
			"failure: too large hash size",
			input{sha256.New(), 8, 65},
			output{nil, ErrTooLargeHashSize},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			in, out := tc.in, tc.out

			conf, err := NewConfig(in.hasher, in.depth, in.hashSize)
			if err != out.err {
				t.Errorf("expected: %v, actual: %v", out.err, err)
			}

			if conf != nil {
				if conf.allLeavesNum != out.config.allLeavesNum {
					t.Errorf("expected: %d, actual: %d", out.config.allLeavesNum, conf.allLeavesNum)
				}
				if conf.allNodesNum != out.config.allNodesNum {
					t.Errorf("expected: %d, actual: %d", out.config.allNodesNum, conf.allNodesNum)
				}
			}
		})
	}
}
