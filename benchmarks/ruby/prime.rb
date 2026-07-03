# frozen_string_literal: true
# SPDX-License-Identifier: BSD-3-Clause
require "prime"
require_relative "_harness"
p15 = 982_451_653          # a prime
c15 = 982_451_651 * 982_451_653 # a hard-ish semiprime-ish composite (odd)
bench("isprime-982451653", 2000) { p15.prime? }
bench("isprime-composite",  200) { c15.prime? }
bench("first-1000",           5) { Prime.first(1000) }
