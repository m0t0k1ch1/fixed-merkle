# fixed-merkle

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

	// check nodes in each levels
	fmt.Println("Root :", tree.Root())
	fmt.Println("Level1 :", tree.Level(1))
	fmt.Println("Level2 :", tree.Level(2))

	// create membership proof of index 0
	proofBytes, err := tree.CreateMembershipProof(0)
	if err != nil {
		panic(err)
	}
	fmt.Println(fmt.Sprintf("Membership proof of index 0 : 0x%x", proofBytes))

	// verify membership proof of index 0
	for i := 0; i < 4; i++ {
		ok, err := tree.VerifyMembershipProof(i, proofBytes)
		if err != nil {
			panic(err)
		}
		fmt.Println(fmt.Sprintf("Member for index %d? : %t", i, ok))
	}
}
```

```
Root : 0x8bb60088191ff98eeb1d3050dd58eeba0e064417bf1d4dd9f6b837112b8aace3
Level1 : [0x1e2abc6e477b5ac3b15d7f153989f49db219c0244ac94b9a1b778c9dbdd5b7e4 0xbf60c7d68271686dce02687a486ee4e339b42a6fef939c071f4578547234dc08]
Level2 : [0x04abc8821a06e5a30937967d11ad10221cb5ac3b5273e434f1284ee87129a061 0x10ae0fdbf8c4f1f2b5e708fd7478abd2bf03b190edc878dc62ada645aa7e0310 0xd155d4b4a5d82abdc42ce8dcc31a7339a003b872ec0332c856f69d6ccc59c967 0x0000000000000000000000000000000000000000000000000000000000000000]
Membership proof of index 0 : 0x10ae0fdbf8c4f1f2b5e708fd7478abd2bf03b190edc878dc62ada645aa7e0310bf60c7d68271686dce02687a486ee4e339b42a6fef939c071f4578547234dc08
Member for index 0? : true
Member for index 1? : false
Member for index 2? : false
Member for index 3? : false
```
