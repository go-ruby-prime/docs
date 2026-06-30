# Roadmap

`go-ruby-prime/prime` is grown **test-first**, each capability differential-tested against MRI
rather than built in isolation. Ruby's `Prime` — the deterministic,
interpreter-independent number-theory core extracted from rbgo's internals — is
**complete**.

| Stage | What | Status |
| --- | --- | --- |
| The generator | `Take(n)` / `First(n)` return the first *n* primes (`Prime.take` / `Prime.first`); `Each(ubound, yield)` enumerates every prime `p <= ubound`; `EachPrime()` is the unbounded cursor. | **Done** |
| Primality — exact | `IsPrime(n)` mirrors `Prime.prime?` / `Integer#prime?` exactly: numbers `< 2` are not prime, every Carmichael number and strong pseudoprime is rejected. Trial division for small inputs, deterministic Baillie–PSW beyond, exact across the whole 64-bit range. | **Done** |
| Factorisation | `PrimeDivision(n)` returns `[[p, exp], …]` in ascending prime order (`Prime.prime_division` / `Integer#prime_division`), with a leading `[-1, 1]` for negatives and a `ZeroError` panic for 0. Large cofactors fall back to Pollard's rho. | **Done** |
| Reconstruction | `Int(pairs)` multiplies `prime**exp` back to the integer (`Prime.int_from_prime_division`), the inverse of `PrimeDivision`. | **Done** |
| Cursor | `Next(n)` / `Prev(n)` step to the adjacent prime. | **Done** |
| Differential oracle & coverage | Deterministic golden tables plus a differential oracle: a corpus computed both here and by the system `ruby` (`Prime.prime?`, `Prime.take`, `Prime.prime_division`, …) and compared; 100% coverage, gofmt + go vet clean, green across all six 64-bit Go arches. | **Done** |

## Documented out-of-scope boundaries

These are **deliberate**, recorded so the module's surface is unambiguous:

- **No interpreter.** The library implements the deterministic algorithm; it
  never runs arbitrary Ruby. Binding it to Ruby objects (the `Prime.each`
  enumerator, `Integer#prime?`) is the consumer's job — that is why `rbgo` binds
  this module rather than the reverse.
- **Reference is reference Ruby (MRI).** Byte-for-byte conformance on the integer
  value model targets MRI's `Prime` behaviour; differences across MRI releases
  are matched to the reference used by the differential oracle.
- **Standalone & reusable.** The module has no dependency on the Ruby runtime;
  the dependency runs the other way.

See [Usage & API](api.md) for the surface and [Why pure Go](why.md) for the
deterministic/interpreter split.
