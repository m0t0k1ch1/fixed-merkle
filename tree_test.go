package merkle

import (
	"crypto/sha256"
	"encoding/hex"
	"hash"
	"testing"
)

func newTestTree(t *testing.T) *Tree {
	tree, err := NewTree(
		sha256.New(),
		3,
		[][]byte{
			[]byte{0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00},
			[]byte{0x01, 0x01, 0x01, 0x01, 0x01, 0x01, 0x01, 0x01},
			[]byte{0x02, 0x02, 0x02, 0x02, 0x02, 0x02, 0x02, 0x02},
		},
	)
	if err != nil {
		t.Fatal(err)
	}
	return tree
}

func TestTree(t *testing.T) {
	type input struct {
		hasher hash.Hash
		depth  uint64
		leaves [][]byte
	}
	type output struct {
		rootHex string
		err     error
	}
	testCases := []struct {
		name string
		in   input
		out  output
	}{
		{
			"failure: too large depth",
			input{
				sha256.New(),
				33,
				nil,
			},
			output{
				"",
				ErrTooLargeDepth,
			},
		},
		{
			"failure: too many leaves",
			input{
				sha256.New(),
				2,
				make([][]byte, 5),
			},
			output{
				"",
				ErrTooManyLeaves,
			},
		},
		{
			"success",
			input{
				sha256.New(),
				3,
				[][]byte{
					[]byte{0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00},
					[]byte{0x01, 0x01, 0x01, 0x01, 0x01, 0x01, 0x01, 0x01},
					[]byte{0x02, 0x02, 0x02, 0x02, 0x02, 0x02, 0x02, 0x02},
				},
			},
			output{
				"f7a08c267a5f438acae772d3bd3c5721188cf4eec29f544d2621d049ec24b4c5",
				nil,
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			in, out := tc.in, tc.out

			tree, err := NewTree(in.hasher, in.depth, in.leaves)
			if err != out.err {
				t.Errorf("expected: %v, actual: %v", out.err, err)
			}
			if err == nil {
				rootHex := hex.EncodeToString(tree.Root().Bytes())
				if rootHex != out.rootHex {
					t.Errorf("expected: %s, actual: %s", out.rootHex, rootHex)
				}
			}
		})
	}
}

func TestTree_CreateMembershipProof(t *testing.T) {
	type input struct {
		index uint64
	}
	type output struct {
		proofHex string
		err      error
	}
	testCases := []struct {
		name string
		tree *Tree
		in   input
		out  output
	}{
		{
			"failure: too large leaf index",
			newTestTree(t),
			input{
				8,
			},
			output{
				"",
				ErrTooLargeLeafIndex,
			},
		},
		{
			"success",
			newTestTree(t),
			input{
				2,
			},
			output{
				"66687aadf862bd776c8fc18b8e9f8e20089714856ee233b3902a591d0d5f2925" +
					"54117bad1f06fb064c25d24002bb68de5835a800d7ad60f679b222a6810c290f" +
					"1223349a40d2ee10bd1bebb5889ef8018c8bc13359ed94b387810af96c6e4268",
				nil,
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			tree, in, out := tc.tree, tc.in, tc.out

			proof, err := tree.CreateMembershipProof(in.index)
			if err != out.err {
				t.Errorf("expected: %v, actual: %v", out.err, err)
			}
			if err == nil {
				proofHex := hex.EncodeToString(proof)
				if proofHex != out.proofHex {
					t.Errorf("expected: %s, actual: %s", out.proofHex, proofHex)
				}
			}
		})
	}
}

func TestTree_VerifyMembershipProof(t *testing.T) {
	type input struct {
		index    uint64
		proofHex string
	}
	type output struct {
		ok  bool
		err error
	}
	testCases := []struct {
		name string
		tree *Tree
		in   input
		out  output
	}{
		{
			"failure: too large leaf index",
			newTestTree(t),
			input{
				8,
				"",
			},
			output{
				false,
				ErrTooLargeLeafIndex,
			},
		},
		{
			"failure: invalid proof size",
			newTestTree(t),
			input{
				0,
				"0000000000000000000000000000000000000000000000000000000000000000" +
					"0000000000000000000000000000000000000000000000000000000000000000" +
					"00000000000000000000000000000000000000000000000000000000000000",
			},
			output{
				false,
				ErrInvalidProofSize,
			},
		},
		{
			"success",
			newTestTree(t),
			input{
				2,
				"66687aadf862bd776c8fc18b8e9f8e20089714856ee233b3902a591d0d5f2925" +
					"54117bad1f06fb064c25d24002bb68de5835a800d7ad60f679b222a6810c290f" +
					"1223349a40d2ee10bd1bebb5889ef8018c8bc13359ed94b387810af96c6e4268",
			},
			output{
				true,
				nil,
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			tree, in, out := tc.tree, tc.in, tc.out

			proof, err := hex.DecodeString(in.proofHex)
			if err != nil {
				t.Fatal(err)
			}
			ok, err := tree.VerifyMembershipProof(in.index, proof)
			if err != out.err {
				t.Errorf("expected: %v, actual: %v", out.err, err)
			}
			if err == nil {
				if ok != out.ok {
					t.Errorf("expected: %t, actual: %t", out.ok, ok)
				}
			}
		})
	}
}
