// SPDX-License-Identifier: BSD-3-Clause
package main

import (
	"math/big"

	"github.com/go-ruby-prime/prime"
)

func main() {
	p15 := big.NewInt(982451653)
	c15 := new(big.Int).Mul(big.NewInt(982451651), big.NewInt(982451653))
	bench("isprime-982451653", 2000, func() { sink = prime.IsPrime(p15) })
	bench("isprime-composite", 200, func() { sink = prime.IsPrime(c15) })
	bench("first-1000", 5, func() { sink = prime.Take(1000) })
}
