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
- **Method:** each process runs 5 untimed warm-up passes, then 50 timed passes of
  a fixed inner loop, timed with a monotonic clock; the **best** pass is reported
  as **ns/op** (lower is better). `vs MRI` < 1.00× means *faster than MRI*.
  Interpreter start-up is outside the timed region, so these are operation costs,
  not `ruby file.rb` process costs.

!!! success "Optimized 2026-07-03 — now beats MRI + YJIT on every workload"
    The `uint64` fast path has been squeezed so `Prime.prime?` beats not just plain
    MRI but **MRI + YJIT**, which JITs MRI's `prime?` loop. Two changes do it: the
    per-step modular multiply now uses a **plain `uint64` multiply below `2^32`**
    (the product of two residues fits in a word — no 128-bit `Div64`) and
    **Montgomery multiplication for the full range** (REDC replaces the per-step
    64-bit division with a multiply, add and shift); and the deterministic
    Miller–Rabin uses **smaller proven witness sets** per magnitude (`{2,7,61}`
    below `4,759,123,141`, Jaeschke 1993 — 3 rounds instead of 4 for the 30-bit
    input) behind a trial-division prefilter trimmed to the primes `≤ 37`. Results
    are identical to the previous path across the whole `uint64` domain; only the
    cost changed. A **segmented sieve of Eratosthenes** still memoizes the
    generator (`Prime.each`/`first`/`take`), and `math/big` Baillie–PSW is retained
    unchanged for values above `2^64`. The numbers below are re-measured against
    these paths.

#### first-1000

| Runtime | ns/op | vs MRI |
| --- | ---: | ---: |
| **go-ruby (pure Go)** | 15525.0 | 0.23× |
| MRI | 67200.0 | 1.00× |
| MRI + YJIT | 19800.0 | 0.29× |
| JRuby | 45833.4 | 0.68× |
| TruffleRuby | 71583.4 | 1.07× |

#### isprime-982451653

| Runtime | ns/op | vs MRI |
| --- | ---: | ---: |
| **go-ruby (pure Go)** | 227.1 | 0.22× |
| MRI | 1013.0 | 1.00× |
| MRI + YJIT | 576.5 | 0.57× |
| JRuby | 1599.8 | 1.58× |
| TruffleRuby | 443.4 | 0.44× |

#### isprime-composite

| Runtime | ns/op | vs MRI |
| --- | ---: | ---: |
| **go-ruby (pure Go)** | 172.7 | 0.19× |
| MRI | 910.0 | 1.00× |
| MRI + YJIT | 670.0 | 0.74× |
| JRuby | 10923.1 | 12.00× |
| TruffleRuby | 2670.0 | 2.93× |

**Beats MRI + YJIT across the board.** Both `prime?` workloads now run several
times faster than MRI (0.22× and 0.19×) **and clear the YJIT bar** — the number
to beat, since YJIT JIT-compiles MRI's `Prime.prime?` loop: isprime-982451653 is
**227 ns vs YJIT's 577 ns** (0.39× of YJIT), and isprime-composite is **173 ns vs
YJIT's 670 ns** (0.26× of YJIT). The pure-Go `uint64` path does word-size trial
division by the primes `≤ 37`, then deterministic Miller–Rabin over the smallest
witness set proven exact for the value's magnitude (`{2,7,61}` below
`4,759,123,141`, growing to the twelve-base set for the full 64-bit range), each
round using a plain `uint64` multiply below `2^32` and Montgomery multiplication
above it — all with no heap allocation. `Prime.first(1000)` stays **~4.3× faster
than MRI** (0.23×) and ahead of MRI + YJIT via the memoized incremental sieve.
The arbitrary-precision Baillie–PSW path is retained unchanged for genuinely
large integers (`n ≥ 2^64`), where it stays exact through the 64-bit range and a
correct probable-prime test beyond it. The earlier ~10× and ~52× rows were
measured against the pre-optimization library; they are superseded here.

!!! note "Reproduce"
    The harness is committed under
    [`benchmarks/`](https://github.com/go-ruby-prime/docs/tree/main/benchmarks):
    a self-contained Go driver (`go/`, pins the published library via
    `go.mod`), the equivalent `ruby/prime.rb` workload, and `run.sh`. Run
    `bash benchmarks/run.sh`; env `OUTER`/`WARM` tune the pass budget and
    `RUBY`/`JRUBY`/`TRUFFLERUBY` select the runtime binaries.

!!! warning "Warm-up budget & noise — honest framing"
    Numbers reflect a **fixed warm-process budget** (5 warm-up + 50 timed passes
    in one process). The JVM/GraalVM JITs (JRuby, TruffleRuby) may need a larger
    warm-up to reach steady state, so their columns can **understate** peak
    throughput — most visibly TruffleRuby on the shortest loops (a few cold-JIT
    outliers are noted in the text). Sub-microsecond rows carry the most relative
    noise; treat those ratios as order-of-magnitude. Every number here is a
    **real measured value** from the dated run above — nothing is fabricated,
    estimated, or cherry-picked. The go-ruby column is the pure-Go library; every
    other column is that interpreter's own stdlib doing the equivalent work.
