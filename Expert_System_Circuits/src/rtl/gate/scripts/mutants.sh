#!/usr/bin/env bash
# Testbench adequacy check: each mutant of gate_core.v must produce at
# least one one-tick mismatch against the Go vectors (Verilator, full vector set).
set -euo pipefail
BUILD=${1:-build}
N=$(cat "$BUILD/vectors.count")
OTHER="rtl/prefix_carry.v rtl/gate_core_prefix.v rtl/gate_reg.v rtl/gate_pipe.v rtl/gate_top.v"
mkdir -p "$BUILD/mut"
: > "$BUILD/mutants.log"
declare -A MUT=(
  [threshold_equality]='s/assign ok\[8\] = (upper_raw <= limit_raw);/assign ok[8] = (upper_raw < limit_raw);/'
  [unsigned_guard_band]='s/wire signed \[63:0\]/wire [63:0]/g'
  [drop_future_timestamp_guard]='s/assign ok\[7\] = (timestamp <= now) && /assign ok[7] = /'
  [stale_age_strict]='s/(age <= max_age)/(age < max_age)/'
  [evidence_low_bits_only]="s/(evidence_code == 8'd1)/(evidence_code[1:0] == 2'd1)/"
  [ignore_undercut]='s/norm_permitted && !undercut/norm_permitted/'
  [revision_31_bits]='s/(data_revision == validated_rev)/(data_revision[30:0] == validated_rev[30:0])/'
  [policy_ignored]='s/assign ok\[6\] = (policy_revision == expected_policy);/assign ok[6] = 1'"'"'b1;/'
)
for name in $(printf '%s\n' "${!MUT[@]}" | sort); do
  sed "${MUT[$name]}" rtl/gate_core.v > "$BUILD/mut/core_$name.v"
  if cmp -s rtl/gate_core.v "$BUILD/mut/core_$name.v"; then
    echo "MUTANT $name not_applied" | tee -a "$BUILD/mutants.log"; exit 1
  fi
  verilator --binary --timing -Wno-fatal -Wno-lint -GN="$N" --top-module tb_gate -Mdir "$BUILD/mut/vl_$name" \
    "$BUILD/mut/core_$name.v" $OTHER tb/tb_gate.v > "$BUILD/mut/build_$name.log" 2>&1
  # A killed mutant makes tb_gate exit non-zero via $fatal; that is expected here.
  line=$("$BUILD/mut/vl_$name/Vtb_gate" 2>&1 | grep '^RESULT one_tick ' || true)
  mism=$(sed -E 's/.*mismatches=([0-9]+).*/\1/' <<< "$line")
  status=killed; [ "$mism" -gt 0 ] || status=SURVIVED
  echo "MUTANT $name mismatches=$mism $status" | tee -a "$BUILD/mutants.log"
done
