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
