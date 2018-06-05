# fixed-merkle

[![GoDoc](https://godoc.org/github.com/m0t0k1ch1/fixed-merkle?status.svg)](https://godoc.org/github.com/m0t0k1ch1/fixed-merkle) [![wercker status](https://app.wercker.com/status/358afd292494eaeb78e289507d4af25c/s/master "wercker status")](https://app.wercker.com/project/byKey/358afd292494eaeb78e289507d4af25c)

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
		false, // leaves are already hashed?
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
- root: 0x8bb60088191ff98eeb1d3050dd58eeba0e064417bf1d4dd9f6b837112b8aace3
- level1: [0x1e2abc6e477b5ac3b15d7f153989f49db219c0244ac94b9a1b778c9dbdd5b7e4 0xbf60c7d68271686dce02687a486ee4e339b42a6fef939c071f4578547234dc08]
- level2: [0x04abc8821a06e5a30937967d11ad10221cb5ac3b5273e434f1284ee87129a061 0x10ae0fdbf8c4f1f2b5e708fd7478abd2bf03b190edc878dc62ada645aa7e0310 0xd155d4b4a5d82abdc42ce8dcc31a7339a003b872ec0332c856f69d6ccc59c967 0x0000000000000000000000000000000000000000000000000000000000000000]

[membership proof]
- proofOfIndex0: 0x10ae0fdbf8c4f1f2b5e708fd7478abd2bf03b190edc878dc62ada645aa7e0310bf60c7d68271686dce02687a486ee4e339b42a6fef939c071f4578547234dc08
  - verification:
    - proofOfIndex0?: true
    - proofOfIndex1?: false
    - proofOfIndex2?: false
    - proofOfIndex3?: false
```
