# fixed-merkle

[![GoDoc](https://godoc.org/github.com/m0t0k1ch1/fixed-merkle?status.svg)](https://godoc.org/github.com/m0t0k1ch1/fixed-merkle) [![wercker status](https://app.wercker.com/status/358afd292494eaeb78e289507d4af25c/s/master "wercker status")](https://app.wercker.com/project/byKey/358afd292494eaeb78e289507d4af25c) [![codecov](https://codecov.io/gh/m0t0k1ch1/fixed-merkle/branch/master/graph/badge.svg)](https://codecov.io/gh/m0t0k1ch1/fixed-merkle)

an implementation of fixed Merkle tree written in golang

``` sh
$ go get github.com/m0t0k1ch1/fixed-merkle
```

## Example

``` go
package main

import (
	"crypto/sha256"
	"fmt"

	merkle "github.com/m0t0k1ch1/fixed-merkle"
)

func main() {
	conf, err := merkle.NewConfig(
		sha256.New(), // hasher
		2,            // depth
		32,           // leaf size (bytes)
	)
	if err != nil {
		panic(err)
	}

	tree, err := merkle.NewTree(
		conf,
		[][]byte{ // leaves
			[]byte{0x01, 0x01, 0x01, 0x01, 0x01, 0x01, 0x01, 0x01},
			[]byte{0x02, 0x02, 0x02, 0x02, 0x02, 0x02, 0x02, 0x02},
			[]byte{0x03, 0x03, 0x03, 0x03, 0x03, 0x03, 0x03, 0x03},
		},
	)
	if err != nil {
		panic(err)
	}

	fmt.Println("[tree]")

	// check nodes in each levels
	fmt.Println("- root:", tree.Root())
	fmt.Println("- level1:", tree.Level(1))
	fmt.Println("- level2:", tree.Level(2))

	fmt.Println("")
	fmt.Println("[membership proof]")

	// create membership proof of index 0
	proofBytes, err := tree.CreateMembershipProof(0)
	if err != nil {
		panic(err)
	}
	fmt.Println(fmt.Sprintf("- proofOfIndex0: 0x%x", proofBytes))

	// verify membership proof of index 0
	fmt.Println("  - verification:")
	for i := 0; i < 4; i++ {
		ok, err := tree.VerifyMembershipProof(i, proofBytes)
		if err != nil {
			panic(err)
		}
		fmt.Println(fmt.Sprintf("    - proofOfIndex%d?: %t", i, ok))
	}
}
```

```
[tree]
- root: 0x5236e8d7c1384c10d767b92b8cd33e569ff33c17f26473c8f95df899f37c47fc
- level1: [0x1e2abc6e477b5ac3b15d7f153989f49db219c0244ac94b9a1b778c9dbdd5b7e4 0x904f88d8f64ff1aa68eda446f8ecf0eb0fc79252ef8d55ca3c17611c00b75f6f]
- level2: [0x04abc8821a06e5a30937967d11ad10221cb5ac3b5273e434f1284ee87129a061 0x10ae0fdbf8c4f1f2b5e708fd7478abd2bf03b190edc878dc62ada645aa7e0310 0xd155d4b4a5d82abdc42ce8dcc31a7339a003b872ec0332c856f69d6ccc59c967 0x66687aadf862bd776c8fc18b8e9f8e20089714856ee233b3902a591d0d5f2925]

[membership proof]
- proofOfIndex0: 0x10ae0fdbf8c4f1f2b5e708fd7478abd2bf03b190edc878dc62ada645aa7e0310904f88d8f64ff1aa68eda446f8ecf0eb0fc79252ef8d55ca3c17611c00b75f6f
  - verification:
    - proofOfIndex0?: true
    - proofOfIndex1?: false
    - proofOfIndex2?: false
    - proofOfIndex3?: false
```
