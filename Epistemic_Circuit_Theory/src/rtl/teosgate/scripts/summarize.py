#!/usr/bin/env python3
"""Collect H5 results from build/ into testdata/teos-rtl/ (summary.json, key_lines.txt).

Fmax values are nextpnr-ice40 static timing estimates after place and route,
not silicon measurements. One-tick claim at f holds only if Fmax >= f.
"""
import argparse
import glob
import json
import os
import re
import statistics

DEVICES = {
    "up5k-sg48": {"device": "iCE40UP5K", "package": "sg48", "nextpnr_flags": "--up5k --package sg48"},
    "hx8k-ct256": {"device": "iCE40HX8K", "package": "ct256", "nextpnr_flags": "--hx8k --package ct256"},
}
VARIANTS = {
    "p0": {"name": "one_tick_direct", "rtl": "teos_gate (teos_gate_core)", "latency_cycles": 1,
           "role": "primary H5 design: direct transcription of PrevalidatedGate"},
    "p1": {"name": "pipelined_2stage", "rtl": "teos_gate_pipe", "latency_cycles": 2,
           "role": "separately labelled exploration, NOT one-tick; initiation interval 1"},
    "p2": {"name": "one_tick_prefix", "rtl": "teos_gate #(PREFIX=1) (teos_gate_core_prefix)", "latency_cycles": 1,
           "role": "separately labelled exploration: same predicate, parallel-prefix comparators"},
}


def read(path):
    with open(path, encoding="utf-8") as handle:
        return handle.read()


def results(path):
    out = {}
    for match in re.finditer(r"^RESULT (\S+) checked=(\d+) mismatches=(\d+)", read(path), re.M):
        out[match.group(1)] = {"checked": int(match.group(2)), "mismatches": int(match.group(3))}
    return out


def yosys_cells(path):
    text = read(path)
    if "=== design hierarchy ===" in text:
        text = text.split("=== design hierarchy ===")[-1]
    cells = {}
    for match in re.finditer(r"^\s+(SB_\w+|\$\w+)\s+(\d+)\s*$", text, re.M):
        cells[match.group(1)] = int(match.group(2))
    lut = cells.get("SB_LUT4", 0)
    ff = sum(v for k, v in cells.items() if k.startswith("SB_DFF"))
    return {"SB_LUT4": lut, "flip_flops": ff, "SB_CARRY": cells.get("SB_CARRY", 0), "cells": cells}


def main():
    parser = argparse.ArgumentParser()
    parser.add_argument("--build", default="build")
    parser.add_argument("--freq", type=float, default=100.0)
    parser.add_argument("--out", required=True)
    args = parser.parse_args()
    b = args.build
    key = []

    versions = [line for line in read(os.path.join(b, "versions.txt")).splitlines() if line.strip()]
    vectors = json.loads(read(os.path.join(b, "vectors_summary.json")))

    sim = {}
    for tool in ("iverilog", "verilator"):
        sim[tool] = {}
        for log in sorted(glob.glob(os.path.join(b, f"sim_{tool}_*.log"))):
            sim[tool].update(results(log))
            key += [f"[{tool}] {line}" for line in read(log).splitlines() if line.startswith("RESULT")]
    total_mismatches = sum(r["mismatches"] for t in sim.values() for r in t.values())

    lint = [line for line in read(os.path.join(b, "lint.log")).splitlines() if line.startswith("verilator lint")]
    lint_warnings = len(re.findall(r"^%Warning", read(os.path.join(b, "lint.log")), re.M))
    key += [f"[lint] {line}" for line in lint] + [f"[lint] %Warning count: {lint_warnings}"]

    mutants = []
    for line in read(os.path.join(b, "mutants.log")).splitlines():
        match = re.match(r"MUTANT (\S+) mismatches=(\d+) (\S+)", line)
        if match:
            mutants.append({"mutant": match.group(1), "mismatches": int(match.group(2)), "status": match.group(3)})
            key.append(f"[mutant] {line}")

    synth = {}
    for v in VARIANTS:
        synth[v] = yosys_cells(os.path.join(b, f"yosys_stat_{v}.txt"))
    gate_only = {}
    for name in ("direct", "prefix", "pipe"):
        gate_only[name] = yosys_cells(os.path.join(b, f"yosys_gate_stat_{name}.txt"))

    runs = []
    for path in sorted(glob.glob(os.path.join(b, "pnr_*_s*.json"))):
        match = re.match(r"pnr_(.+)_(p\d)_s(\d+)\.json", os.path.basename(path))
        dev, var, seed = match.group(1), match.group(2), int(match.group(3))
        report = json.loads(read(path))
        (clock, fmax), = report["fmax"].items()
        util = report["utilization"]
        crit = next(p for p in report["critical_paths"] if p["from"].startswith("posedge") and p["to"].startswith("posedge"))
        delays = {}
        for step in crit["path"]:
            delays[step["type"]] = round(delays.get(step["type"], 0.0) + step["delay"], 3)
        runs.append({
            "device_tag": dev, "variant": var, "seed": seed, "clock": clock,
            "fmax_mhz": round(fmax["achieved"], 2), "constraint_mhz": fmax["constraint"],
            "meets_constraint": fmax["achieved"] >= fmax["constraint"],
            "ICESTORM_LC": util["ICESTORM_LC"], "SB_IO": util["SB_IO"], "SB_GB": util["SB_GB"],
            "critical_reg_to_reg": {
                "source": crit["path"][0]["from"]["cell"], "sink": crit["path"][-1]["to"]["cell"],
                "steps": len(crit["path"]), "delay_ns": delays,
                "total_ns": round(sum(s["delay"] for s in crit["path"]), 3),
            },
        })
        log = read(path[:-5] + ".log")
        line = re.findall(r"^.*Max frequency for clock.*$", log, re.M)[-1].strip()
        key.append(f"[nextpnr {dev} {var} seed {seed}] {line}")
        if seed == 1:
            block = re.search(r"Info: Device utilisation:\n((?:Info: \t.*\n)+)", log)
            key += [f"[nextpnr {dev} {var} seed 1] {l.strip()}" for l in block.group(1).splitlines()
                    if re.search(r"ICESTORM_LC|SB_IO|SB_GB", l)]

    aggregate = []
    for dev, info in DEVICES.items():
        for var, vinfo in VARIANTS.items():
            sel = [r for r in runs if r["device_tag"] == dev and r["variant"] == var]
            if not sel:
                continue
            f = [r["fmax_mhz"] for r in sel]
            lc = sel[0]["ICESTORM_LC"]
            aggregate.append({
                "device": info["device"], "package": info["package"], "variant": vinfo["name"],
                "latency_cycles": vinfo["latency_cycles"], "seeds": [r["seed"] for r in sel],
                "fmax_mhz_min": min(f), "fmax_mhz_median": statistics.median(f), "fmax_mhz_max": max(f),
                "target_mhz": args.freq, "seeds_meeting_target": sum(x >= args.freq for x in f),
                "claim_holds_at_target_all_seeds": min(f) >= args.freq,
                "lc_used": lc["used"], "lc_available": lc["available"],
                "io_used": sel[0]["SB_IO"]["used"], "io_available": sel[0]["SB_IO"]["available"],
            })

    summary = {
        "hypothesis": "H5 (PREREGISTRATION-2026-11-03.md): register-to-register one-tick PREVALIDATED gate predicate",
        "reference": "teoshardware.PrevalidatedGate (src/teoshardware/calculate.go), vectors by cmd/teosgatevectors",
        "tool_versions": versions,
        "target_frequency_mhz": args.freq,
        "interpretation": [
            "Fmax is nextpnr-ice40 post-route static timing estimate (tool model), not a silicon or HIL measurement.",
            "One-tick claim at f on a device holds only if post-route Fmax >= f; reported per seed.",
            "RTL implements PrevalidatedGate (the 9-condition postcondition comparator over pre-computed flags), NOT the Horn/normative closure engine of kernel K.",
            "Excluded: kernel K Horn/defeasible/normative closure, SHA-256/custody verification, bus transfer of the word, actuation, serial loading time.",
            "The pin wrapper (462-bit serial shift loader in teos_gate_top) adds 462 FFs (934 FFs in top_with_wrapper vs 472 FFs in gate_only for one-tick variants; 1071 vs 609 FFs for pipelined_2stage) and routing pressure; it is not the gate.",
            "Pins are unconstrained (--pcf-allow-unconstrained): placer chooses IO locations.",
        ],
        "scope_and_wrapper_overhead": {
            "implemented_function": "teoshardware.PrevalidatedGate (9-condition postcondition comparator over pre-computed flags, revisions, timestamps, and signed fixed-point bounds)",
            "not_implemented_in_rtl": [
                "Kernel K Horn/defeasible/normative closure engine (EvaluateKernel)",
                "SHA-256 and Ed25519 custody/signature verification",
                "Bus transfer and 462-clock serial shift-register loading time",
                "Physical actuator transport and response time",
            ],
            "flip_flop_breakdown": {
                "serial_shift_wrapper_ffs": 462,
                "gate_word_input_ffs": 462,
                "gate_result_output_ffs": 10,
                "one_tick_direct_gate_only_ffs": gate_only["direct"]["flip_flops"],
                "one_tick_direct_top_with_wrapper_ffs": synth["p0"]["flip_flops"],
                "one_tick_prefix_gate_only_ffs": gate_only["prefix"]["flip_flops"],
                "one_tick_prefix_top_with_wrapper_ffs": synth["p2"]["flip_flops"],
                "pipelined_2stage_gate_only_ffs": gate_only["pipe"]["flip_flops"],
                "pipelined_2stage_top_with_wrapper_ffs": synth["p1"]["flip_flops"],
            },
        },
        "word_bits": 462, "result_bits": 10,
        "vectors": vectors,
        "simulation": sim, "simulation_total_mismatches": total_mismatches,
        "lint": {"lines": lint, "warnings": lint_warnings},
        "mutation_adequacy": mutants,
        "yosys_synth_ice40_top_with_wrapper": synth,
        "yosys_synth_ice40_gate_only_no_pnr": gate_only,
        "nextpnr_aggregate": aggregate,
        "nextpnr_runs": runs,
    }
    os.makedirs(args.out, exist_ok=True)
    with open(os.path.join(args.out, "summary.json"), "w", encoding="utf-8") as handle:
        json.dump(summary, handle, indent=2, ensure_ascii=False)
        handle.write("\n")
    with open(os.path.join(args.out, "key_lines.txt"), "w", encoding="utf-8") as handle:
        handle.write("\n".join(["# Tool versions"] + versions + ["# Key result lines"] + key) + "\n")
    for row in aggregate:
        print(f'{row["device"]:10} {row["variant"]:18} Fmax {row["fmax_mhz_min"]:.2f}/{row["fmax_mhz_median"]:.2f}/'
              f'{row["fmax_mhz_max"]:.2f} MHz  meets {args.freq:g}: {row["seeds_meeting_target"]}/{len(row["seeds"])}  '
              f'LC {row["lc_used"]}/{row["lc_available"]}  IO {row["io_used"]}/{row["io_available"]}')
    for gname, tvar in (("direct", "p0"), ("prefix", "p2"), ("pipe", "p1")):
        g, t = gate_only[gname], synth[tvar]
        print(f'cells {gname:6} gate_only: LUT4={g["SB_LUT4"]} FF={g["flip_flops"]} CARRY={g["SB_CARRY"]} | '
              f'top_with_wrapper ({tvar}): LUT4={t["SB_LUT4"]} FF={t["flip_flops"]} (+{t["flip_flops"] - g["flip_flops"]} wrapper FFs) CARRY={t["SB_CARRY"]}')
    print(f"simulation mismatches total: {total_mismatches}")


if __name__ == "__main__":
    main()
