package merkle

import (
	"crypto/sha256"
	"hash"
	"testing"
)

func TestNewConfig(t *testing.T) {
	var newConfigTestCases = []struct {
		hasher   hash.Hash
		depth    int
		hashSize int
		config   *Config
		err      error
	}{
		{sha256.New(), 1, 1, &Config{allLeavesNum: 2, allNodesNum: 3}, nil},
		{sha256.New(), 16, 64, &Config{allLeavesNum: 65536, allNodesNum: 131071}, nil},
		{sha256.New(), 0, 32, nil, ErrTooSmallDepth},
		{sha256.New(), 17, 32, nil, ErrTooLargeDepth},
		{sha256.New(), 1, 0, nil, ErrTooSmallHashSize},
		{sha256.New(), 16, 65, nil, ErrTooLargeHashSize},
	}

	for _, tc := range newConfigTestCases {
		conf, err := NewConfig(tc.hasher, tc.depth, tc.hashSize)
		if err != tc.err {
			t.Errorf("expected: %v, actual: %v", tc.err, err)
		}
		if conf != nil {
			if conf.allLeavesNum != tc.config.allLeavesNum {
				t.Errorf("expected: %d, actual: %d", tc.config.allLeavesNum, conf.allLeavesNum)
			}
			if conf.allNodesNum != tc.config.allNodesNum {
				t.Errorf("expected: %d, actual: %d", tc.config.allNodesNum, conf.allNodesNum)
			}
		}
	}
}
