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
		sha256.New(),
		2,  // depth
		32, // leaf size (bytes)
	)
	if err != nil {
		panic(err)
	}

	tree, err := merkle.NewTree(
		conf,
		[][]byte{ // leaves
			[]byte{0x01, 0x01, 0x01, 0x01, 0x01, 0x01, 0x01, 0x01, 0x01, 0x01, 0x01, 0x01, 0x01, 0x01, 0x01, 0x01, 0x01, 0x01, 0x01, 0x01, 0x01, 0x01, 0x01, 0x01, 0x01, 0x01, 0x01, 0x01, 0x01, 0x01, 0x01, 0x01},
			[]byte{0x02, 0x02, 0x02, 0x02, 0x02, 0x02, 0x02, 0x02, 0x02, 0x02, 0x02, 0x02, 0x02, 0x02, 0x02, 0x02, 0x02, 0x02, 0x02, 0x02, 0x02, 0x02, 0x02, 0x02, 0x02, 0x02, 0x02, 0x02, 0x02, 0x02, 0x02, 0x02},
			[]byte{0x03, 0x03, 0x03, 0x03, 0x03, 0x03, 0x03, 0x03, 0x03, 0x03, 0x03, 0x03, 0x03, 0x03, 0x03, 0x03, 0x03, 0x03, 0x03, 0x03, 0x03, 0x03, 0x03, 0x03, 0x03, 0x03, 0x03, 0x03, 0x03, 0x03, 0x03, 0x03},
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
Root : 0xc8d3d8d2b13f27ceeccdc699119871f9f32ea7ed86ff45d0ad11f77b28cd7568
Level1 : [0x50a27d4746f357cb700cbe9d4883b77fb64f0128828a3489dc6a6f21ddbf2414 0xa41b855d2db4de9052cd7be5ec67d6586629cb9f6e3246a4afa5ba313f07a9c5]
Level2 : [0x72cd6e8422c407fb6d098690f1130b7ded7ec2f7f5e1d30bd9d521f015363793 0x75877bb41d393b5fb8455ce60ecd8dda001d06316496b14dfa7f895656eeca4a 0x648aa5c579fb30f38af744d97d6ec840c7a91277a499a0d780f3e7314eca090b 0x0000000000000000000000000000000000000000000000000000000000000000]
Membership proof of index 0 : 0x75877bb41d393b5fb8455ce60ecd8dda001d06316496b14dfa7f895656eeca4aa41b855d2db4de9052cd7be5ec67d6586629cb9f6e3246a4afa5ba313f07a9c5
Member for index 0? : true
Member for index 1? : false
Member for index 2? : false
Member for index 3? : false
```
