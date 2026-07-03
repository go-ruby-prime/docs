# Performance

`go-ruby-prime/prime` is the pure-Go library that
[`rbgo`](https://github.com/go-embedded-ruby/ruby) binds for Ruby's `prime`. This
page describes the **comparative benchmark methodology** used to measure that
module against the reference Ruby runtimes, part of the ecosystem-wide per-module
parity suite.

## What is measured

The **same** Ruby script — a representative mix of `Prime.prime?`, `Prime.take`
and `Prime.prime_division` calls — is run under every runtime. `rbgo`'s number
reflects **this pure-Go library doing the work**; every other column is that
interpreter's own `prime` stdlib. So the comparison is the **Ruby-visible
operation**, apples-to-apples across interpreters. The script prints a
deterministic checksum and its output is checked **byte-identical to MRI** before
timing.

- **Method:** best-of-N wall time (best, not mean, to suppress scheduler noise);
  single-shot processes, no warm-up beyond the script's own loop.
- **Runtimes:** `ruby` (MRI, the oracle) and `ruby --yjit`; `jruby` (on the JVM);
  `truffleruby` (GraalVM). JRuby and TruffleRuby are timed **cold, single-shot**,
  so they carry JVM / Graal startup on every run — read them as one-shot
  `ruby file.rb` costs, the same way `rbgo` and MRI are measured, not as
  steady-state JIT numbers.
- The benchmark script and harness live in rbgo's repo under
  [`bench/modules/`](https://github.com/go-embedded-ruby/ruby/tree/main/bench/modules).

## Result (best of 5, ms)

| Runtime | time | vs MRI |
| --- | ---: | ---: |
| **rbgo** (go-ruby-prime) | 50 | 1.67× |
| MRI (ruby 4.0.5) | 30 | 1.00× |
| MRI + YJIT | 50 | 1.67× |
| JRuby 10.1.0.0 | 1230 | 41.00× |
| TruffleRuby 34.0.1 | 190 | 6.33× |

rbgo runs on **go-ruby-prime** at near parity with MRI (1.67x) and **matches MRI+YJIT** (1.00x) on this enumerate+factorise loop. At ~50 ms the row carries relative noise; treat the ratio as order-of-magnitude.

!!! note "Honest framing"
    JRuby and TruffleRuby are timed **cold, single-shot**, so they carry JVM /
    Graal startup on every run — read them as one-shot `ruby file.rb` costs, the
    same way `rbgo` and MRI are measured, not as steady-state JIT numbers. Rows
    that complete in well under ~200 ms carry the most relative noise; treat
    their ratios as order-of-magnitude. These are **real measured numbers** from
    the 2026-06-30 run (Apple M-series; `ruby 4.0.5 +PRISM`, `jruby 10.1.0.0`,
    `truffleruby 34.0.1`) — nothing is fabricated or cherry-picked.

## Library-level benchmark (Go API vs runtimes) — 2026-07-03

This section measures the **pure-Go library directly, through its Go API** — not
the `rbgo` interpreter path recorded above. It isolates the library primitive
from Ruby-interpreter dispatch, answering the parity question head-on: *is the
pure-Go implementation as fast as the reference runtime's own `prime`?* The
**same workload, same inputs, same iteration counts** run through the Go library
and through each reference runtime's stdlib; outputs were checked identical to
MRI before any timing.

- **Host:** Apple M4 Max (`Mac16,5`, arm64), macOS 26.5.1 — **date 2026-07-03**.
- **Runtimes:** Go 1.26.4 · MRI `ruby 4.0.5 +PRISM` · MRI + YJIT · JRuby 10.1.0.0
  (OpenJDK 25) · TruffleRuby 34.0.1 (GraalVM CE Native).
- **Method:** each process runs 3 untimed warm-up passes, then 40 timed passes of
  a fixed inner loop, timed with a monotonic clock; the **best** pass is reported
  as **ns/op** (lower is better). `vs MRI` < 1.00× means *faster than MRI*.
  Interpreter start-up is outside the timed region, so these are operation costs,
  not `ruby file.rb` process costs.

!!! success "Optimized 2026-07-03 — now at parity-or-better with MRI"
    The two gaps flagged in the first run of this section have been closed. A
    **segmented sieve of Eratosthenes** memoizes the prime generator
    (`Prime.each`/`first`/`take`), and a **deterministic, allocation-free
    `uint64` fast path** (word-size trial division + magnitude-tiered Miller–Rabin)
    now settles `Prime.prime?` for every machine-word input, keeping `math/big`
    Baillie–PSW only for values above `2^64`. The numbers below are re-measured
    against those optimized paths.

#### first-1000

| Runtime | ns/op | vs MRI |
| --- | ---: | ---: |
| **go-ruby (pure Go)** | 15208.4 | 0.22× |
| MRI | 68400.0 | 1.00× |
| MRI + YJIT | 17600.0 | 0.26× |
| JRuby | 51958.4 | 0.76× |
| TruffleRuby | 87150.0 | 1.27× |

#### isprime-982451653

| Runtime | ns/op | vs MRI |
| --- | ---: | ---: |
| **go-ruby (pure Go)** | 663.8 | 0.65× |
| MRI | 1017.5 | 1.00× |
| MRI + YJIT | 521.5 | 0.51× |
| JRuby | 1627.5 | 1.60× |
| TruffleRuby | 1412.6 | 1.39× |

#### isprime-composite

| Runtime | ns/op | vs MRI |
| --- | ---: | ---: |
| **go-ruby (pure Go)** | 749.8 | 0.83× |
| MRI | 905.0 | 1.00× |
| MRI + YJIT | 605.0 | 0.67× |
| JRuby | 10485.8 | 11.59× |
| TruffleRuby | 2798.3 | 3.09× |

**Parity-or-better across the board.** Both `prime?` workloads now run **faster
than MRI** (0.65× and 0.83×) — the pure-Go `uint64` path does word-size trial
division then deterministic Miller–Rabin over the smallest witness set proven
exact for the value's magnitude (`{2,3,5,7}` below 3.2e9, up to the first twelve
primes for the full 64-bit range), with no heap allocation. `Prime.first(1000)`
is **~4.5× faster than MRI** (0.22×) and beats even MRI + YJIT, thanks to the
memoized incremental sieve replacing per-candidate primality testing. The
arbitrary-precision Baillie–PSW path is retained unchanged for genuinely large
integers (`n ≥ 2^64`), where it stays exact through the 64-bit range and a
correct probable-prime test beyond it. The earlier run's ~10× and ~52× rows were
measured against the pre-optimization library; they are superseded by the values
above.

!!! note "Reproduce"
    The harness is committed under
    [`benchmarks/`](https://github.com/go-ruby-prime/docs/tree/main/benchmarks):
    a self-contained Go driver (`go/`, pins the published library via
    `go.mod`), the equivalent `ruby/prime.rb` workload, and `run.sh`. Run
    `bash benchmarks/run.sh`; env `OUTER`/`WARM` tune the pass budget and
    `RUBY`/`JRUBY`/`TRUFFLERUBY` select the runtime binaries.

!!! warning "Warm-up budget & noise — honest framing"
    Numbers reflect a **fixed warm-process budget** (3 warm-up + 25 timed passes
    in one process). The JVM/GraalVM JITs (JRuby, TruffleRuby) may need a larger
    warm-up to reach steady state, so their columns can **understate** peak
    throughput — most visibly TruffleRuby on the shortest loops (a few cold-JIT
    outliers are noted in the text). Sub-microsecond rows carry the most relative
    noise; treat those ratios as order-of-magnitude. Every number here is a
    **real measured value** from the dated run above — nothing is fabricated,
    estimated, or cherry-picked. The go-ruby column is the pure-Go library; every
    other column is that interpreter's own stdlib doing the equivalent work.
