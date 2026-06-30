# go-ruby-prime documentation

**Ruby's `Prime` generation, primality testing and factorisation in pure Go — MRI-compatible, no cgo.**

`go-ruby-prime/prime` is a faithful, pure-Go (zero cgo) reimplementation of the deterministic,
interpreter-independent core of MRI 4.0.5's `Prime` class and the
`Integer#prime?` / `Integer#prime_division` refinements, matching reference Ruby
(MRI) byte-for-byte on the integer value model. The module path is
`github.com/go-ruby-prime/prime`.

It generates the primes, tests primality, factorises an integer and reconstructs
it — **without any Ruby runtime**. Primality is **exact, not probabilistic**: small
inputs use trial division; everything larger uses a deterministic **Baillie–PSW**
test (strong base-2 Miller–Rabin + strong Lucas), which has no counterexample below
2⁶⁴, so every Carmichael number and strong pseudoprime is correctly rejected.

It was **extracted from rbgo into a reusable standalone
library**: the module is standalone and importable by any Go program, and it is
the `prime` backend bound into
[go-embedded-ruby](https://github.com/go-embedded-ruby/ruby)
by `rbgo` as a native module — just like
[go-ruby-regexp](https://github.com/go-ruby-regexp),
[go-ruby-yaml](https://github.com/go-ruby-yaml) and
[go-ruby-marshal](https://github.com/go-ruby-marshal). The dependency runs the
other way: this library has **no dependency on the Ruby runtime**.

!!! success "Status: complete — MRI byte-exact"
    Faithful port of Ruby's `prime`: the **generator** (`Take` / `First` / `Each` / `EachPrime`), exact **primality** (`IsPrime`, rejecting every Carmichael number and strong pseudoprime), **factorisation** (`PrimeDivision`, with a leading `[-1, 1]` for negatives and a `ZeroError` panic for 0), **reconstruction** (`Int`), and the **cursor** (`Next` / `Prev`). Validated by a **differential oracle** against the system `ruby` / `Prime` at 100% coverage, `gofmt` + `go vet` clean, CI green across the six 64-bit Go targets.

## Quick taste

```go
fmt.Println(prime.Take(5))                        // [2 3 5 7 11]   (Prime.take 5)
fmt.Println(prime.IsPrime(big.NewInt(561)))       // false          (Carmichael)
fmt.Println(prime.PrimeDivision(big.NewInt(12)))  // [[2 2] [3 1]]  (prime_division)

n := prime.Int(prime.PrimeDivision(big.NewInt(360)))
fmt.Println(n)                                    // 360
```

## Repositories

| Repo | What it is |
| --- | --- |
| [`prime`](https://github.com/go-ruby-prime/prime) | the library — Ruby's `Prime` in pure Go |
| [`docs`](https://github.com/go-ruby-prime/docs) | this documentation site (MkDocs Material, versioned with mike) |
| [`go-ruby-prime.github.io`](https://github.com/go-ruby-prime/go-ruby-prime.github.io) | the organization landing page (Hugo) |
| [`brand`](https://github.com/go-ruby-prime/brand) | logo and brand assets |

## Principles

- **Pure Go, `CGO_ENABLED=0`** — trivial cross-compilation, a single static
  binary, no C toolchain; the only dependency is `math/big`.
- **Exact, not probabilistic.** Primality is settled by deterministic Baillie–PSW,
  exact across the whole 64-bit range; validated by a differential oracle against
  the `ruby` binary.
- **Standalone & reusable.** Extracted from rbgo's internals; no dependency on
  the Ruby runtime — the dependency runs the other way.
- **100% test coverage** is the target, enforced as a CI gate, across 6 arches.

## Where to go next

- [Why pure Go](why.md) — why prime number theory is deterministic enough to live
  as a standalone, interpreter-independent Go library.
- [Usage & API](api.md) — the public surface and worked examples.
- [Roadmap](roadmap.md) — what is done and what is downstream by design.

Source lives at [github.com/go-ruby-prime/prime](https://github.com/go-ruby-prime/prime).
