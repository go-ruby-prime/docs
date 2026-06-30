# Usage & API

The public API lives at the module root (`github.com/go-ruby-prime/prime`). It is **Ruby-shaped but Go-idiomatic**: `Take` / `IsPrime` / `PrimeDivision` mirror `Prime.take` / `Prime.prime?` / `Prime.prime_division`, while the surface follows Go conventions — explicit `*big.Int` values, no global state.

!!! success "Status: implemented"
    The library is built and importable as `github.com/go-ruby-prime/prime`, bound into
    `rbgo` as a native module; see [Roadmap](roadmap.md).

## Install

```sh
go get github.com/go-ruby-prime/prime
```

## Worked example

```go
package main

import (
	"fmt"
	"math/big"

	"github.com/go-ruby-prime/prime"
)

func main() {
	fmt.Println(prime.Take(5))                       // [2 3 5 7 11]   (Prime.take 5)
	fmt.Println(prime.IsPrime(big.NewInt(561)))      // false          (Carmichael)
	fmt.Println(prime.IsPrime(big.NewInt(7919)))     // true
	fmt.Println(prime.PrimeDivision(big.NewInt(12))) // [[2 2] [3 1]]  (prime_division)
	fmt.Println(prime.PrimeDivision(big.NewInt(-12)))// [[-1 1] [2 2] [3 1]]

	// Reconstruct the integer from its factorisation.
	n := prime.Int(prime.PrimeDivision(big.NewInt(360)))
	fmt.Println(n) // 360

	// Bounded enumeration (Prime.each(11)).
	prime.Each(11, func(p *big.Int) bool { fmt.Print(p, " "); return true })
	fmt.Println() // 2 3 5 7 11
}
```

## Shape

```go
// IsPrime reports whether n is prime (Prime.prime? / Integer#prime?).
func IsPrime(n *big.Int) bool

// Take / First return the first n primes (Prime.take / Prime.first).
func Take(n int) []*big.Int
func First(n int) []*big.Int

// Each yields every prime p <= ubound, stopping early if yield returns false
// (Prime.each(ubound) { |p| ... }).
func Each(ubound int64, yield func(p *big.Int) bool)

// EachPrime returns a stateful generator: each call returns the next prime,
// starting at 2 and continuing forever (the unbounded Prime.each enumerator).
func EachPrime() func() *big.Int

// PrimeDivision returns the [prime, exponent] pairs of n in ascending order
// (Prime.prime_division / Integer#prime_division); a leading [-1, 1] carries a
// negative sign. It panics with ZeroError for n == 0 (MRI's ZeroDivisionError);
// PrimeDivisionErr is the non-panicking form.
func PrimeDivision(n *big.Int) [][2]*big.Int
func PrimeDivisionErr(n *big.Int) ([][2]*big.Int, error)

// Int reconstructs the integer from a prime-division slice
// (Prime.int_from_prime_division), the inverse of PrimeDivision.
func Int(pairs [][2]*big.Int) *big.Int

// Next / Prev return the adjacent prime (Prev returns nil when none exists).
func Next(n *big.Int) *big.Int
func Prev(n *big.Int) *big.Int

type ZeroError struct{} // mirrors Ruby's ZeroDivisionError
```

## Ruby ↔ Go value model

Every integer flows through `*big.Int`, so a host can map its own `Integer` to and
from this package without precision loss.

| Ruby                                 | Go                            |
| ------------------------------------ | ----------------------------- |
| `Prime.prime?(n)` / `n.prime?`       | `IsPrime(n)`                  |
| `Prime.take(n)` / `Prime.first(n)`   | `Take(n)` / `First(n)`        |
| `Prime.each(ubound) { ... }`         | `Each(ubound, yield)`         |
| `Prime.each` (enumerator)            | `EachPrime()`                 |
| `Prime.prime_division(n)` / `n.prime_division` | `PrimeDivision(n)`  |
| `Prime.int_from_prime_division(ps)`  | `Int(ps)`                     |
| `ZeroDivisionError`                  | `ZeroError` (panic)           |

## Algorithm

Primality is exact, not probabilistic, over the range Ruby programs use:

- **Small inputs** (`< 1000²`) are settled by **trial division** against the
  primes up to 1000.
- **Larger inputs** use **Baillie–PSW** — a strong base-2 Miller–Rabin test
  combined with a strong Lucas test (Selfridge parameters). BPSW has **no known
  counterexample** and is proven to have none below 2⁶⁴, so the result is exact
  across the entire 64-bit range and a sound probable-prime test beyond it.
- **Factorisation** strips small primes, then splits the cofactor with **Pollard's
  rho** (Brent's variant), recursing to proven primes.

## MRI conformance

Correctness is defined by reference Ruby. A **differential oracle** runs a corpus
through both the system `ruby` (`Prime.prime?`, `Prime.take`,
`Prime.prime_division`, …) and this library and compares the results — not
approximated from memory. The oracle tests skip themselves where `ruby` is not on
`PATH` (e.g. the qemu arch lanes), so the cross-arch builds still validate the
library.

## Relationship to Ruby

`go-ruby-prime/prime` is **standalone and reusable**, and is the `prime` backend
bound into [go-embedded-ruby](https://github.com/go-embedded-ruby/ruby) by `rbgo`
as a native module — the same way [go-ruby-regexp](https://github.com/go-ruby-regexp),
[go-ruby-yaml](https://github.com/go-ruby-yaml) and
[go-ruby-marshal](https://github.com/go-ruby-marshal) are bound. The dependency
runs the other way: this library has no dependency on the Ruby runtime.
