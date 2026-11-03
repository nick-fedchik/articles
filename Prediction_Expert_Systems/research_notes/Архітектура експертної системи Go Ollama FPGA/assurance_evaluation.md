# Assurance boundaries and evaluation protocol

**Research note — 2 November 2026.** The domain remains open, so this is a reference decomposition and falsifiable test plan, not a safety classification or certification opinion. “Documented evidence” below is distinguished from recommended engineering practice.

## Assurance boundary and component responsibilities

**Recommended decomposition.** Separate the assurance-bearing decision path from probabilistic assistance and acceleration:

1. **Evidence and knowledge:** record source/device identity, event and receipt times, units, quality, calibration status, and immutable provenance. Version rules and knowledge; cite exact premises and versions in each result. Represent missing, stale, contradictory, or invalid-unit inputs explicitly.
2. **Deterministic Go core:** normalize validated input, evaluate versioned rules, and emit typed decisions with rule IDs, premises, and reason codes. Restrict side effects; use race detection and concurrency tests. Go’s memory model guarantees sequential consistency for data-race-free programs, not deterministic application logic, scheduling, I/O, or forecasts [5].
3. **Temporal monitor:** check freshness, ordering, rates, persistence, deadlines, and invariants over timestamped events; report the checked window and missing-data behavior. A runtime-verification study demonstrates monitor-triggered switching in one simulation, not proof for this system or domain [1].
4. **Forecast/uncertainty service:** return point estimates and horizon-indexed prediction intervals/sets, model version, and abstention/OOD flags. Keep forecasts advisory absent a domain safety case. Conformal coverage relies on assumptions; exchangeability guarantees do not automatically survive temporal or unit-level shift [3].
5. **Ollama interface:** limit it to summaries, operator questions, or candidate explanations. Treat output as untrusted: validate schema and cited evidence IDs, and forbid changes to facts, thresholds, rule outcomes, or actuators. Ollama documents schema-constrained local JSON but says cloud currently lacks structured-output support; valid JSON proves neither factuality nor semantic validity [4].
6. **FPGA accelerator:** confine it to bounded kernels. Specify widths, rounding, saturation/overflow, endianness, and numeric tolerance; keep a software reference. Bind the bitstream to a version/hash and device; validate DMA bounds, completion identity, watchdog, reset, and fallback. FPGA execution alone ensures neither determinism nor correctness.
7. **Human escalation:** route disagreement, OOD, critical evidence gaps, or high-consequence/near-threshold cases to a qualified reviewer under predefined procedures. Record cited inputs and rationale; an override does not rewrite source evidence or provenance.

**Boundary hazards.** Test API timeouts, partial streams, malformed output, model changes, prompt injection, duplicate retries, and latency spikes; advisory failure must not block or alter the deterministic path. Test Go races, event-order ambiguity, parsing, unit conversion, and stale caches; sensor/knowledge clock skew, gaps, duplicates, calibration drift, conflicting revisions, and untrusted updates; and FPGA bitstream mismatch, DMA/cache coherency, truncation, numeric mismatch, timeout, stale interrupts, and reset faults. These are hypotheses, not claims about an implementation.

## Falsifiable evaluation protocol

**Recommended experiments (pre-register thresholds before testing).**

- **Data and forecast validity:** split chronologically and by independent asset/unit/site; prevent overlap or near-duplicate units across partitions. Report coverage, interval width, calibration, and point error by horizon, site, regime, and missingness. Test drift, outages, and unseen units. Marginal coverage can conceal subgroup or horizon failures; the cited drift-robust conformal methods do not establish conditional coverage at every horizon [3].
- **Decision quality:** on labeled episodes and fault injection, report false alarms, misses, confidence intervals, lead time, alarm persistence, abstention/escalation, and severity-weighted misses. Lock the test set; include conflicting rules, unavailable evidence, delayed/out-of-order events, and monitor boundaries. Synthetic results are not field validation.
- **Runtime and service boundaries:** inject slow responses, disconnects, false-but-schema-valid output, prompt attacks, upgrades, and duplicate requests. Prove advisory output cannot change deterministic verdicts and timeout invokes fallback/escalation. Measure end-to-end p50/p95/p99 and observed maximum under cold/warm starts, load, contention, and dependency failure; predefine a deadline and count violations.
- **FPGA correctness and costs:** compare against an independent software reference using golden/random vectors, extrema, overflow/rounding boundaries, and legal modes; require bit equality or a predeclared decision-relevant tolerance. Where feasible, prove equivalence under explicit input/state assumptions and verify them independently. Test bitstream rejection, timeout, reset, DMA errors, and fallback. Compare throughput, latency, memory, utilization, and board power with software; record device, tools, constraints, and bitstream digest. AMD documents target-specific logic equivalence, not whole-system assurance [8].
- **Evidence reproducibility:** freeze the code, rules/knowledge snapshot, model tag/weights, prompts, FPGA bitstream, driver, and datasets; retain raw inputs, timestamps, configuration, and outputs. A separate evaluator should reproduce the headline metrics from this package.

## Claims and assurance limits

**Documented evidence.** NIST describes its AI Risk Management Framework as voluntary risk-management guidance, not a product certificate [2]. ISO/SAE 21434 specifies cybersecurity engineering for road-vehicle E/E systems across their lifecycle; ISO 26262 is likewise vehicle-specific functional-safety standardization [6][7]. Neither scope statement certifies this architecture generically. Open research and vendor documentation support methods and tool capabilities only within their stated assumptions [1][3][8].

**Claims not yet supportable.** Until domain, intended function, hazard class, jurisdiction, hardware, environment, and conformity route are selected, claim no SIL/ASIL/performance level, compliance, certification, calibrated operational forecast, or safety from determinism/FPGA alone. LLM explanations, provenance, simulation scores, and equivalence checks are not by themselves validated knowledge or certification evidence. An assurance case still needs domain hazards/requirements, traceability, independence, configuration and verification strategy, fault analysis, and independent review.

## Central sources

1. Yang et al., “Case Study: Runtime Safety Verification of Neural Network Controlled System,” 2024 preprint: https://arxiv.org/abs/2408.08592
2. NIST, AI Risk Management Framework overview: https://www.nist.gov/itl/ai-risk-management-framework
3. Barber et al., “Conformal prediction beyond exchangeability,” 2023: https://arxiv.org/abs/2202.13415
4. Ollama, Structured Outputs documentation: https://docs.ollama.com/capabilities/structured-outputs
5. Go, The Go Memory Model: https://go.dev/ref/mem
6. ISO, ISO/SAE 21434:2021: https://www.iso.org/standard/70918.html
7. ISO, ISO 26262-1:2018: https://www.iso.org/standard/68383.html
8. AMD, Vivado Design Analysis, Logic Equivalence Checking (UG906): https://docs.amd.com/r/en-US/ug906-vivado-design-analysis/Logic-Equivalence-Checking
