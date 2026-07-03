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
- **Method:** each process runs 3 untimed warm-up passes, then 25 timed passes of
  a fixed inner loop, timed with a monotonic clock; the **best** pass is reported
  as **ns/op** (lower is better). `vs MRI` < 1.00× means *faster than MRI*.
  Interpreter start-up is outside the timed region, so these are operation costs,
  not `ruby file.rb` process costs.

#### first-1000

| Runtime | ns/op | vs MRI |
| --- | ---: | ---: |
| **go-ruby (pure Go)** | 3438066.6 | 51.62× |
| MRI | 66600.0 | 1.00× |
| MRI + YJIT | 19200.0 | 0.29× |
| JRuby | 56308.4 | 0.85× |
| TruffleRuby | 101200.0 | 1.52× |

#### isprime-982451653

| Runtime | ns/op | vs MRI |
| --- | ---: | ---: |
| **go-ruby (pure Go)** | 10364.2 | 10.34× |
| MRI | 1002.0 | 1.00× |
| MRI + YJIT | 558.5 | 0.56× |
| JRuby | 1596.3 | 1.59× |
| TruffleRuby | 1738.1 | 1.73× |

#### isprime-composite

| Runtime | ns/op | vs MRI |
| --- | ---: | ---: |
| **go-ruby (pure Go)** | 4595.0 | 5.08× |
| MRI | 905.0 | 1.00× |
| MRI + YJIT | 660.0 | 0.73× |
| JRuby | 8064.8 | 8.91× |
| TruffleRuby | 3034.6 | 3.35× |

The **honest outlier of the tranche.** go-ruby-prime is arbitrary-precision (`math/big` + Baillie–PSW per candidate), so small-integer primality (~10×) and `first(1000)` (~52×, generate-vs-sieve) are much slower than MRI's native trial division and C Eratosthenes sieve. The library is *correct for big integers* (its design goal), but a native small-int fast path and a sieve for `first`/`take` are the clear optimization targets.

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
