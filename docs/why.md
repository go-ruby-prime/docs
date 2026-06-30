# Why pure Go

`go-ruby-prime/prime` reimplements Ruby's `Prime` in **pure Go, with cgo disabled**. The
slice of Ruby it covers is **deterministic and interpreter-independent**: sieving
the primes, testing primality, factorising an integer — given the integer, the
result is a pure function of that input, with no live binding and no evaluation of
arbitrary Ruby. That is exactly the part that can — and should — live as a
standalone Go library, separate from the interpreter.

## Extracted from rbgo, reusable by anyone

This library began life inside [go-embedded-ruby](https://github.com/go-embedded-ruby/ruby)'s
`rbgo`. It has been **extracted into a reusable standalone library** so that:

- any Go program can import `github.com/go-ruby-prime/prime` directly, with no Ruby runtime;
- the dependency runs the *other* way — `rbgo` binds this module as a native
  module (the same pattern as [go-ruby-regexp](https://github.com/go-ruby-regexp),
  [go-ruby-yaml](https://github.com/go-ruby-yaml) and
  [go-ruby-marshal](https://github.com/go-ruby-marshal)), rather than this module
  depending on the interpreter;
- the behaviour is pinned by a **differential oracle** against the system
  `ruby`, independent of any one consumer.

## What it is — and isn't

The number theory behind `Prime` — sieving the primes, testing primality,
factorising an integer — is fully deterministic and needs **no interpreter**, so
it lives here as pure Go. Binding it to Ruby objects (the `Prime.each`
enumerator, `Integer#prime?`) is the host's job; this library hands back
`*big.Int` values the host wraps in its own `Integer`. Every integer flows
through `*big.Int`, so a host can map its own `Integer` to and from this package
without precision loss.

## Exact, not probabilistic

Primality is exact over the range Ruby programs use. Small inputs are settled by
**trial division**; larger inputs use **Baillie–PSW** — a strong base-2
Miller–Rabin test combined with a strong Lucas test (Selfridge parameters). BPSW
has **no known counterexample** and is proven to have none below 2⁶⁴, so the
result is exact across the entire 64-bit range and a sound probable-prime test
beyond it. Factorisation strips small primes, then splits the cofactor with
**Pollard's rho** (Brent's variant), recursing to proven primes.

## Why pure Go matters here

Because the library is CGO-free and depends only on `math/big`, it:

- cross-compiles to every Go target with no C toolchain, and links into a single
  static binary;
- has **no dependency on the Ruby runtime** — the dependency runs the other way;
- can be differentially tested against the `ruby` binary wherever one is on
  `PATH`, while the cross-arch lanes (where `ruby` is absent) still validate the
  library itself.

See [Usage & API](api.md) for the surface and [Roadmap](roadmap.md) for what is
in scope.
