# Knowledge evolution in a predictive expert system: mapping local research-prototype design evidence and the author's books

Scope note for the report writer. Internal names, identifiers and paths of the research prototype were deliberately generalised in this note; the prototype's project documentation is not part of this repository. This is a local, read-only reading. No tests, builds, or hardware commands were run, and nothing was sent to external services. Labels used throughout:

- **[DESIGN]**: intent or requirement stated in a memo, ADR, requirement, or chapter.
- **[CODE-static]**: Go code that appears to implement the design, based only on static reading. It was not executed.
- **[EVIDENCE-recorded]**: a measurement or test result that a document *claims*. It was not re-run or independently verified here.
- **[INTERNAL-ONLY]**: a finding whose source is internal prototype material (code, memos, requirements, decisions, task records). In this note its internal names, paths and identifiers are already replaced by generic descriptions; only the generalised concept may move into public book text. The prototype's development rules (§9–§10) forbid internal names in the monograph, including the prototype's own name, subsystem and storage-subsystem names, the resident-host binary name, knowledge-container file extensions and format magic, concrete format-version constants, the names of the checker instruction set and the check virtual machine, and internal hex section codes. They require generic terms such as «бінарний контейнер бази знань» and «резидентний демон керування пам'яттю». All prototype material is [INTERNAL-ONLY] as a source.

Research-prototype documents are project design evidence, not independent scientific validation.

Prototype documents are cited by topic, not by identifier. Short names used below:

- *self-learning memo*: memo on expert-system self-learning and neuro-symbolic epistemic feedback loops;
- *gap-ticket memo*: memo on autonomous discovery of knowledge deficits and gap tickets;
- *self-organisation memo*: memo on self-organising expert systems and the cybernetics of knowledge management;
- *continual-learning memo*: memo on continual learning, delayed labels and off-policy evaluation;
- *telemetry memo*: memo on telemetry, state estimation and runtime assurance;
- *release-gate memo*: memo on release gates, data splits and label quality;
- *superseding memo*: memo on cross-document temporal superseding and conflict mining;
- *atomic-swap memo*: memo on the resident knowledge-host daemon and atomic swap;
- *changelog decision*: architecture decision on a bitemporal changelog and idempotent recompilation of knowledge containers;
- *format-version decision*: architecture decision on a single-integer container format version and runtime migration;
- *federation decision*: architecture decision on federation of corporate epistemic memory and ontology migration;
- *hypothesis-kind decision*: architecture decision that marks hypothesis atoms as a distinct possibility kind;
- *trust-chain decision*: architecture decision on host trust-chain verification, rollback and revocation;
- *hypothesis-atom requirement*: knowledge-compiler requirement on epistemic-possibility hypothesis atoms from non-normative sources;
- *gap-triage requirement*: requirement of the event and notification module on automated gap-ticket triage;
- *semantic-drift requirement*: system requirement on continuous semantic SRE and silent-degradation telemetry.

## 1. How the research prototype's documentation describes knowledge acquisition, compilation, versioning, provenance, hypotheses, gap queues, human gates, continual learning, delayed labels and OPE

### Takeaway
The research prototype has a rich design corpus for knowledge evolution:

- gap tickets lead to harvesting, compilation, and atomic hot-swap;
- a bitemporal changelog with SemVer content versions;
- superseded defeaters;
- typed hypothesis atoms;
- a quarantine for feedback;
- delayed-label states, OPE with a common-support guard, anchor replay, and grouped release gates.

The documents disagree on one central point. Some say knowledge admission is fully autonomous (the self-learning and gap-ticket memos). Others say it ends in a human or owner admission decision (the self-organisation, continual-learning, telemetry and release-gate memos). The continual-learning, OPE, release-gate and telemetry memos are design-only. Their linked requirements and ADRs point to a different repository location that has no matching packages.

### Cited Findings

**Continual learning, delayed labels and OPE [DESIGN]**

- **Problem statement.** The continual-learning memo (§1) names four challenges in evolving the system from operating experience (session logs, operator feedback, new specification releases):
  1. delayed labels: treating unfinished episodes as failures raises false degradation alarms;
  2. policy shift: evaluating new policy π on logs from π₀ without propensity correction or a common-support check;
  3. catastrophic forgetting on base standards when new ones are added;
  4. feedback poisoning.

  Source: continual-learning memo, §1.
- **Invariants.** The continual-learning memo (§2) states three:
  1. Zero Regression on Core Axioms: every rule or weight update is tested on a fixed anchor replay buffer, and a core regression is a zero-tolerance block.
  2. Common Support Guard: IPS/DR is allowed only where π₀(a|x) ≥ ε_min > 0; otherwise the evaluation is refused.
  3. Quarantine Feedback Invariant: no operator feedback enters the KB directly. It goes to a feedback quarantine queue and needs byte-level citation verification against the primary source.

  Source: continual-learning memo, §2.
- **Specified components.** The continual-learning memo (§3) specifies a delayed-outcome manager with four states (pending, resolved as success, resolved as failure, expired without a signal), and pending cases are not counted as errors. An off-policy evaluation engine covers IPS, DM and DR, with ESS = (Σw)²/Σw² and rejection when ESS < ESS_min. An anchor replay buffer is kept per registered protocol. The planned validation is: no false alarm on day 6 when labels arrive on day 10; a reliable refusal when there is no overlap; 100% retention on base protocols after sequential updates. Source: continual-learning memo, §3–§4. [INTERNAL-ONLY: the memo names an internal continual-learning package and links to a software requirement, a task and an ADR stored in a separate repository of the query verification engine]
- **Earlier formulation.** The self-learning memo (§4) gives an earlier version: a replay buffer of core axioms that requires ∀α∈A_core: K_t ∪ ΔK ⊢ α, or else an "Atomic Fail-Closed Commit" block, plus a Doubly Robust estimator V̂_DR(π_new). It claims this "гарантує" improvement without risky experiments on users. That is an overclaim. The continual-learning memo later adds the support/ESS guard. Source: self-learning memo, §4.1–§4.2.

**Self-learning and the human gate [DESIGN]**

- The self-learning memo (§1.2) defines self-learning as an update of discrete verified assertions, K_{t+1} = K_t ∪ ΔK_verified \ ΔK_conflicted, where the proof-verification predicate for α against the raw source corpus must be TRUE. It explicitly says "без втручання людини". It describes three feedback loops:
  1. refusal-driven discovery: gap descriptor ⟨S, P, D_seed⟩, then harvester, then byte-verified fact;
  2. Socratic calibration: the user's answer becomes an ontological constraint;
  3. counterfactual self-play: predicate mutation, an SMT consistency check, and cross-RFC conflict hunting that marks old assertions `superseded`.

  Source: self-learning memo, §1.2, §2.1–§2.3.
- **Change levels by time scale.** The self-organisation memo (§2) defines five levels and who authorises each:
  - 0, reflex (µs–ms): rule verified before launch;
  - 1, state adaptation (ms–s): truth maintenance through defeaters, tombstones and retractions;
  - 2, self-organisation inside the admissible configuration set C_adm (s–min): mode estimator;
  - 3, structural learning (hours–weeks; rules, concepts, L0): **"процедура допуску й рецензент"**, via AMIE, abduction, extraction, quarantine and release;
  - 4, organisational (weeks–years; policies, domain ontology): knowledge engineers and the policy owner.

  The memo says levels 0–2 may run on board without a human, and level 3 "може починатися автоматично, але закінчується рішенням про допуск". A new rule or concept is "навчання бази знань, а не самоорганізація". Its diagram also states that no loop changes the variables of a slower loop (paraphrased). Source: self-organisation memo, §2–§3. Note: the memo's original Ukrainian wording uses a term that the book reserves for geometric and electrical contexts; book text should use «цикл».
- **Contradiction on autonomy.** The gap-ticket memo (§4) describes "автономне самонавчання без галюцинацій": gap tickets → gap prospector → autonomous harvester → compiler → hot-swap daemon attaching a new pack through `mmap(MAP_SHARED)`, with no human admission step. Source: gap-ticket memo, §2–§4. This is contradicted by the level-3 human/admission gate in the self-organisation memo (§2) and by the quarantine invariant in the continual-learning memo (§2). [INTERNAL-ONLY: the memo names a message-bus package, a NATS gap stream, a file-based gap-ticket spool, the resident-host binary and internal product names; all are generalised here]

**Gap tickets and triage [DESIGN]**

- The gap ticket is triggered by an out-of-domain verdict, a qualified non-answer verdict, or an unresolved collision between two equally ranked norms. Its payload includes missing entities, missing predicates, suggested standards, status and priority. Source: gap-ticket memo, §2–§2.1.
- Gap triage requires SimHash/embedding deduplication into meta-tickets with a hit counter, a popularity threshold that triggers a harvester directive, and a broadcast "knowledge pack updated" event followed by an RCU hot-swap. The acceptance criteria are:
  - criterion 1: ≥90% collapse of 1,000 repeated tickets;
  - criterion 2: swap < 500 µs with no interrupted request;
  - criterion 3: full version history in the audit log.

  Source: gap-triage requirement, §1–§2. [INTERNAL-ONLY]

**Atomic swap [DESIGN]**

- The new pack is mapped by a fresh reader. The header, 64-byte alignment and Merkle root are verified. The active reader is switched with `atomic.Pointer`, and the old reader is unmapped only after in-flight work ends. The claims are zero downtime and zero lost requests. Source: atomic-swap memo, §3. [INTERNAL-ONLY]

**Versioning, changelog and bitemporality [DESIGN]**

- The changelog decision (status ACCEPTED) defines:
  - **Idempotence guard.** Skip recompilation when the source SHA-256 and the container format version match and no forced-rebuild flag is given.
  - **Two-level change detection** "через порівняння AST та множини видобутих деонтичних норм". The "syntactic realignment" change class (rules unchanged, byte offsets re-aligned) bumps PATCH. The "semantic delta" change class bumps MINOR for added compatible norms and MAJOR for annulled or changed constraints.
  - **Three storage levels for the changelog:**
    1. a signed in-container manifest (Ed25519; manifest as Merkle leaf 0);
    2. a registry with supersedes / superseded-by links;
    3. a hash-chained ledger H_i = SHA-256(Tag‖Seq‖H_{i−1}‖Payload), where H_{i−1} is the previous entry's hash and Tag is a fixed ASCII domain-separation prefix.
  - **Bitemporal validity.** Valid time and transaction time follow Jensen and Snodgrass. On a new edition, the old pack is **not overwritten** but marked `Superseded`. Rebuttal defeaters in the new pack block the old norm for the present while keeping it for retrospective audit.

  Source: changelog decision, §1–§3. [INTERNAL-ONLY: the decision names the container file extension, the format magic, a registry index file, a ledger package and hex section codes; all are generalised here]
- **Version axes.** Container format version and tool/API versioning are separate axes. The binary container is identified by a single integer format-version field with no pseudo-SemVer, and any incompatible layout change bumps N→N+1. Tools, services, instruction sets, protocols and message schemas use SemVer 2.0.0. Source: the prototype's development rules, §9; format-version decision (title only read).
- **Cross-document superseding.** "Obsoletes" and "Updates" relations are mined from text patterns ("updates section X of RFC Y", "is now DEPRECATED/HISTORIC/PROHIBITED"). The delta-patch algorithm sets the original fact's valid-until date to the new document's date, creates a new fact whose supersedes link points to the original fact's identifier, and keeps both with byte-level quotes. Source: superseding memo, §2.1–§2.2.
- **Federation and migration.** Rules are indexed by a validity interval [T_valid_start, T_valid_end). Old rules get a superseded-by marker carrying the new clause's identifier in the dynamic L1 layer without mutating the immutable L0 segment. Queries with T_query < T_migration route to the previous rule version. A new-edition rule acts as an ASPIC+ undercutting defeater. The cost is 32 bytes per rule in the L1 delta index. Source: federation decision, §2–§3.

**Hypotheses versus admitted facts [DESIGN]**

- The hypothesis-atom requirement (APPROVED 2026-10-20):
  - A `MAY` trigger in sources of the informational or technical-documentation authority tier yields a distinct `HYPOTHESIS` atom, never a deontic code. The source tier decides, not the wording.
  - Unknown codes must be treated as non-assertive (fail-closed).
  - Hypotheses carry the same byte custody as norms.
  - When a hypothesis is the only evidence, the system must return a qualified non-answer verdict and never an affirmative answer. A hypothesis cannot satisfy an obligation, prohibition or permission query.
  - Hypothesis quality is measured by deterministic sampling with a Wilson 95% interval. The acceptance threshold is "fixed by the owner after the first measurement, not before".

  Source: hypothesis-atom requirement, §1–§3. The hypothesis-kind decision records the status "ACCEPTED / IMPLEMENTED" in its header.

**Release gates and label quality [DESIGN]**

- The release-gate memo requires:
  - group splits by document identifier, revision and protocol family, with n-gram/SHA leakage checks;
  - fail-closed rejection when any mandatory slice is missing;
  - zero failures on blocking cases;
  - a group bootstrap for confidence intervals;
  - a label sanitiser (Confident Learning) that escalates suspicious labels to a human "без автоматичної мутації еталону";
  - LLM-as-judge only for linguistic properties.

  Source: release-gate memo, §2–§3.

**Telemetry and drift [DESIGN]**

- The telemetry memo requires:
  - bounded observation freshness τ_fresh, then a degraded-observation state;
  - dual timestamps plus a sequence number;
  - a fail-closed safe-hold state when critical telemetry is lost;
  - STL monitors □_[0,T], ◇_[0,T] with watermarks;
  - a drift detector with an adaptive window that "без прийняття остаточних рішень про зміну бізнес-правил (тільки ескалація до аудитора)".

  Source: telemetry memo, §2–§3. Drift is a signal for escalation, never an automatic rule change.
- The semantic-drift requirement requires a semantic-drift circuit breaker: Δ_semantic(r, canonical reference) > τ_drift ⇒ circuit-breaker trip ∧ fallback. Source: semantic-drift requirement. [INTERNAL-ONLY; it also credits a named external SRE practitioner, so attribution must be checked before any public use]

**Honesty constraint [DESIGN]**

- Byte custody proves that a quote matches source bytes. It does **not** prove that the source is true, that extraction was correct, or that the answer is relevant or semantically correct. Empirical claims need evaluated answers with confidence intervals. Source: the prototype's architecture overview.

### Inferences

- The research prototype's design separates five change classes that a predictive expert system needs:
  1. new or changed facts or norms from new source editions (superseding, bitemporal);
  2. container format changes (integer migration);
  3. content-version changes (SemVer changelog);
  4. runtime state adaptation (defeaters and L1 deltas, no L0 mutation);
  5. learned components and policies (OPE, replay, quarantine).

  This decomposition fits the book's objects (R, Π, predictor, B-coverage, metrics) but is not yet stated for them.
- The most book-ready generalisable principle is the self-organisation memo's rule: faster loops may not change slower loops' variables, and new rules or concepts require an admission decision. It reconciles autonomy with the ch06 requirement of "позитивна рецензія" (positive review).
- The autonomy claims of the self-learning and gap-ticket memos ("без втручання людини", "автономне самонавчання без галюцинацій") should not be carried into the book as stated. They conflict with the self-organisation, continual-learning, telemetry and release-gate memos, and with the architecture overview's statement that custody ≠ truth.

### Gaps

- The linked requirements, ADRs and tasks for the telemetry, release-gate and continual-learning memos (three software requirements, three ADRs and three tasks) point to the documentation of a separate repository of the query verification engine. A listing of that repository's ADR directory showed none of the referenced ADRs. Their content could not be checked.
- The semantic-drift requirement cites two research memos numbered beyond the last memo present in the prototype's research-memo directory, so these references appear dangling.
- No research-prototype document found specifies:
  - recalibration of a predictor's conformal calibration set over time;
  - versioning of an action-schema registry with owner/review;
  - versioning of a policy object Π;
  - re-verification of completeness of prohibition coverage when the extraction method changes.
- The search covered memos, requirements and ADRs by keyword. An architecture memo on epistemic geology, calibration and estimation matched "recalibrat" but was not read.

## 2. Static inspection: is versioning or admission of new knowledge implemented in code?

### Takeaway
Implemented, by static reading:

- an append-only changelog with SemVer bumps in the compiler manifest;
- single-integer format versioning with migration;
- anti-rollback lineage floors and signed revocation manifests in the resident host;
- an ASPIC+ superseded-defeater filter with bitemporal query dates;
- hypothesis sampling for quality labelling;
- SimHash gap-ticket clustering;
- an extraction-level quarantine with typed reasons.

The change classification in code is much weaker than the changelog decision says. No code was found for delayed labels, OPE, anchor replay, release-gate slices, label sanitising or drift detection. The hot-swap path shows no human approval step.

### Cited Findings

**Changelog and content versioning [CODE-static]**

- The changelog-entry record has a version, a Unix timestamp, a change class ("initial compilation" / "syntactic realignment" / "semantic delta" / "lexicon refinement"), the previous and current source SHA-256, a semantic-delta block (rules added, rules removed, rules modified, coordinates realigned) and a rationale. Source: changelog module of the knowledge compiler. [INTERNAL-ONLY]
- The function that computes the version-change (evolution) record classifies changes **only by comparing atom counts**:
  - same source SHA, same count and a changed lexicon filter give a PATCH "lexicon refinement";
  - **any equal atom count gives a PATCH "syntactic realignment"**, even if rule contents changed;
  - a different count gives a MINOR "semantic delta", with the rules-added/rules-removed fields set from the *net* count difference;
  - the rules-modified field is never set;
  - no MAJOR bump is produced;
  - the new entry is appended to the previous changelog.

  Source: changelog module of the knowledge compiler, evolution-record function. This contradicts the AST and norm-set diff and the MINOR/MAJOR rule in the changelog decision (§2.2).
- The compiler calls the evolution-record function at two sites, passing the previous manifest, the source SHA and a `len(...)` count. The container manifest includes the version, authority tier, source class, obsoletes/updates relations (RFC numbers), atom count, provenance, changelog, provenance sources and conditional evidence. Source: knowledge-compiler core (manifest definition and both call sites). [INTERNAL-ONLY]

**Format versioning and migration [CODE-static]**

- Comments in the format-version module state: a single unsigned integer format-version field in the container header; a machine-readable registry of versions; a readable window [Min..Current]; a test of released-version immutability fails when a released entry is edited; and a layout fingerprint. The byte-migration function migrates N−x→N and enforces quote-byte equality. Sources: format-version and migrator modules of the knowledge compiler. [INTERNAL-ONLY]

**Admission controls at load or swap time [CODE-static]**

- The check-and-advance operation of the lineage-floor store (over lineage, version and creation time) gives persistent, monotonic anti-rollback. An explicit floor-setting operation raises the floor from a signed revocation manifest. Source: anti-rollback module of the resident host daemon in the query verification engine.
- The revocation manifest is epoch-based. Checks for a revoked key, a revoked pack (by Merkle root) and staleness (fail-closed beyond a maximum age) are provided, and the manifest-application operation verifies against trust roots. Source: revocation module of the resident host daemon.
- The host module of the resident daemon calls the check-and-advance operation when mounting (two call sites) and applies revocations. The swap operation and the atomic knowledge-base namespace swap both delegate to the mount operation. The automatic pack-upgrade function verifies quote custody, performs an RCU swap, and keeps the original version if signing keys are required but missing. Source: host module of the resident host daemon. No human or owner approval call was seen on this path. [INTERNAL-ONLY]

**Retirement through superseding [CODE-static]**

- The superseded-rule defeater record holds old/new rule IDs, effective date, reason and citation. In the superseded-defeat evaluation function, if the query date ≥ the superseding date, the argument relying on the old rule is undercut; before that date the old rule keeps historical validity. The bitemporal superseded-rule filter function splits candidates into admissible and defeated. It is called from the enterprise bitemporal index. Sources: ASPIC+ superseding module of the query verification engine; its enterprise bitemporal index.

**Hypotheses and human labelling [CODE-static]**

- The hypothesis-atom sampling function selects a deterministic quota by the smallest SHA-256(seed, pack, atom index) among atoms whose deontic modality is the hypothesis kind. The hypothesis CSV-export function writes a file whose last two columns are "filled in by a human (or evaluation model)". The hypothesis-label scoring function computes precision overall and by tier. The sampling command is wired into the compiler CLI as a subcommand. Sources: hypothesis-sampling audit module of the knowledge compiler; the corresponding subcommand of the compiler CLI.

**Gap queue [CODE-static]**

- The gap-ticket triage cluster aggregates structurally similar gap tickets. The engine uses an exact-key map plus multi-index SimHash tables and persists to a NATS KV interface. Source: triage clustering module of the event and notification module. [INTERNAL-ONLY]

**Extraction quarantine [CODE-static]**

- Proposals that fail verification are quarantined with a typed quarantine reason (e.g., empty span, span outside the sentence), and "confidence is never an acceptance criterion". Source: proposal verifier of the knowledge compiler's ingestion module. This is quarantine of *extraction proposals*, not the operator-feedback quarantine queue of the continual-learning memo.

**Absent [CODE-static]**

- No Go identifiers matching the component names specified in the memos were found for: the delayed-outcome manager, the off-policy evaluation engine, the anchor replay buffer, the feedback quarantine queue, the grouped data splitter, the release-gate slice matrix, the drift detector, the temporal assurance monitor, the label sanitiser, or any propensity-related identifier (prefix match). The search covered the Go modules of the knowledge compiler, query verification engine, event and notification module, source harvester, hardware check executor and language module.
- The internal package directories for continual learning, release gates, KB verification and telemetry exist neither in the prototype's query-verification-engine module nor in the separate repository of that engine. This was checked by grep and path existence only. Specs: continual-learning memo, §3; release-gate memo, §3; telemetry memo, §3.

**Recorded evidence (claims in project documents; not re-run) [EVIDENCE-recorded]**

- **Hypothesis atoms.** A measurement dated 2026-10-20 used 180 atoms (30 per modality), labels from `qwen2.5:32b`, and 93% agreement with independent labelling by the requirement's author. `MAY` atoms were norms in 3/30 cases (10%, 95% CI 3.5–25.6%). The share of genuine hypotheses "не проводився" (was not measured). Source: hypothesis-atom requirement, §0.
- **Sampling tool.** The task on hypothesis sampling and quality measurement is `done`, but its text only restates the method ("поріг фіксує власник після першого виміру"). No recorded hypothesis-precision result was found in the prototype's report and audit directories by grep. Source: completed hypothesis-sampling task record.
- **Anti-rollback and revocation.** The completed task on host-socket authorisation, anti-rollback and revocation states that rolled-back and revoked packs are rejected (with refusal verdicts for a rolled-back pack and for a revoked pack) and that its anti-rollback-protection and revocation-enforcement tests pass "100% pass під race detector". Source: that task record. The trust-chain decision is nonetheless labelled "PARTIALLY IMPLEMENTED / BLOCKED", citing a knowledge-compiler task held in the blocked queue (decision header).
- **Not found:** measured evidence for criterion 1 (≥90% dedup) or criterion 2 (<500 µs swap) of the gap-triage requirement. Test files exist but were not run, e.g. the triage spike tests of the event and notification module and the runtime-upgrade tests of the resident host daemon.

### Inferences

- Code covers the integrity and order of knowledge versions: idempotence, a monotonic lineage, revocation, signatures, and a superseded-by-date filter. It does not cover the *epistemic quality* of a change. No content diff of rules, no regression gate on anchors, and no OPE or delayed-label machinery were seen.
- The changelog decision calls the changelog "tamper-evident" and "semantic". Because classification is by atom count, a same-count rewrite of norms would be recorded as a whitespace realignment. The audit trail could therefore mislabel exactly the changes that matter for permissions. This is a concrete argument for the book: version records need a content-addressed diff per rule, not counts.
- The swap path verifies integrity and lineage but appears to have no review gate. That matches the gap-ticket memo's autonomy model, not the self-organisation memo's level-3 admission model.

### Gaps

- Static reading cannot confirm runtime behaviour or test outcomes. No tests were run, by instruction.
- I did not inspect how superseded-rule defeater records are created from the manifest's obsoletes/updates relations, i.e. whether compile-time superseding generation exists. A grep of call sites showed only the enterprise bitemporal index as a consumer.
- I did not inspect whether the knowledge compiler's inspect command prints the changelog (acceptance criterion 3 of the changelog decision).

## 3. What the predictive book (ch06, ch05 §1) says about how R, Π, predictors, B-coverage and acceptance metrics change over time — and what is missing

### Takeaway
Chapter 6 defines *static* admissibility objects with lifecycle hooks but no lifecycle:

- an action schema has an owner, a version and a positive review;
- Reg(r) = T only if the schema is valid and not revoked;
- V(a,t) handles replaced or expired source editions;
- B requires declared-complete sources and verified detection completeness "на активному способі витягу";
- section 8 metrics serve as acceptance criteria for a new KB version.

Chapter 5 §1 warns that retraining on logs shaped by the system's own recommendations can erode sensitivity (self-defeating prophecy). It gives only a sufficient convergence condition for repeated risk minimisation.

Missing: version selection at time t for R and Π; schema validity intervals; re-verification triggers for B-coverage and for predictor calibration; regression gates across versions; delayed and censored outcomes for maintenance actions; OPE under policy shift; retirement or deprecation semantics; hypotheses about evolution.

### Cited Findings

**Chapter 6 objects and their change hooks (book text)**

- **Action schema.** Definition 6.2 makes the action schema r ∈ R a production rule π₁∧…∧π_k ⇒ α_r "яке має власника, версію та позитивну рецензію". Source: [ch06 §4.1, Def. 6.2](../../ch06-from-prediction-to-admissible-action.md); glossary entry "Схема дії … має власника, версію та рецензію" — [ch06 Словник](../../ch06-from-prediction-to-admissible-action.md).
- **Reg(r).** Reg(r_a) = T if the schema is valid in the registry ("має власника, версію, рецензію і не відкликана"), otherwise F. It enters the grounding predicate G (formula 6.11). If G ∈ {F, U}, the candidate is dropped with an explanation of what is missing. Ideas proposed by the language model without a registered schema "зберігаються лише як гіпотези для людини й не отримують статусів". Source: [ch06 §4.1, formula (6.11)](../../ch06-from-prediction-to-admissible-action.md). The only lifecycle states are therefore "registered and valid" versus "revoked or incomplete". There are no draft, quarantine, deprecated or superseded states, and no rule picks *which* version of r applies at time t.
- **Schemas that use the predictor's state.** In §4.2, r₂ fires when the premise state is U because the interval is wide, producing "збільшити частоту вимірювань". r₃ fires on C and escalates to a reliability engineer. Source: [ch06 §4.2](../../ch06-from-prediction-to-admissible-action.md). These are hooks for an information-gathering path (U-driven data acquisition) but not for knowledge revision.
- **Policy Π.** Π is "явно задана політика: ієрархія джерел, правила розв'язання колізій, таблиця повноважень і тягар доведення для винятків". β(d) ∈ {strict, presumed} "визначає політика системи, а не текст документа". A conflict C is escalated to the policy owner. Sources: [ch06 §5.1](../../ch06-from-prediction-to-admissible-action.md); [ch06 §2.3, formula (6.5)](../../ch06-from-prediction-to-admissible-action.md); [ch06 §2.2 and §5.2 item 4](../../ch06-from-prediction-to-admissible-action.md). No version, owner-review or audit-trail requirement is stated for Π itself. Table 6.18 places Π in the knowledge base but attaches "власник, версія, рецензія" only to R.
- **Temporal validity V(a,t).** V(a,t) = F "якщо редакцію явно замінено або строк її дії сплив". Source: [ch06 §5.1](../../ch06-from-prediction-to-admissible-action.md). This is the only explicit time-indexed validity in the chapter. It applies to *source editions*, not to schemas, policy or predictors.
- **Burden of proof and B-coverage.** B = T only if all three hold:
  1. the authoritative source set is declared complete for the query scope;
  2. completeness of detecting prohibitions and exceptions has been verified (§7.2);
  3. no applicable prohibition is unresolved under Π.

  If (1) or (2) fails, B = U; contradictions give B = C. "Дія не стає дозволеною лише тому, що пошук не знайшов заборони." `presumed` is allowed only where detection completeness has been verified. Source: [ch06 §5.1 and §2.3](../../ch06-from-prediction-to-admissible-action.md). §7.2 adds that for document types where completeness "не перевірено на активному способі витягу", B and `presumed` are downgraded to U, so a model failure raises UNKNOWN rather than false permissions — [ch06 §7.2](../../ch06-from-prediction-to-admissible-action.md). This implies that coverage evidence is tied to an extraction-method version. The chapter does not say how that evidence is versioned, when it expires, or what re-verification is required.
- **Two-level citation check.** Level 1 is custody: the fragment exists verbatim in the stated document version (file hash, bounds, text match). Level 2 is semantic correctness of modality, negation, actor and condition. PROV-O lineage is suggested, "але не посвідчує істинності джерела". Source: [ch06 §7.2](../../ch06-from-prediction-to-admissible-action.md).
- **Acceptance metrics.** FDR_perm (6.20) and FPR_perm (6.21) are defined with different denominators. There is a worked example with a Clopper–Pearson CI. Action rule: if the upper CI bound exceeds the risk owner's pre-registered threshold, PERMITTED may only be used with human review, and the sample should grow or B be tightened. A matching protocol is fixed before test labels are revealed, and annotators are independent. Sources: [ch06 §8.1–§8.2](../../ch06-from-prediction-to-admissible-action.md). Table 6.18 says: "Набуття знань: Метрики розділу 8 стають критеріями приймання нової версії бази знань до розгортання" — [ch06 §7.1, Table 6.18](../../ch06-from-prediction-to-admissible-action.md). This is the chapter's only explicit link between metrics and versioned KB releases. No regression rule between versions (old versus new on the same cases) is given.
- **Limits.** §9.1 lists the knowledge-acquisition bottleneck: R and Π "мають писати й підтримувати фахівці … модель її не усуває". Usefulness depends on how fast gaps (U) are closed. B = T is often unreachable in open domains. Guarantees are marginal. Source: [ch06 §9.1 items 1–4](../../ch06-from-prediction-to-admissible-action.md).
- **Hypotheses 6.1–6.3.**
  - 6.1: lower FPR_perm than the naive rule within a missed-action tolerance.
  - 6.2: the interval premise gives fewer false-safe recommendations at comparable UNKNOWN share.
  - 6.3: U-premise schemas shorten time to a grounded decision.

  All thresholds are pre-registered. Source: [ch06 §9.2](../../ch06-from-prediction-to-admissible-action.md). All three are single-version, cross-sectional comparisons. None concerns stability across knowledge versions, drift, or the effect of revisions.

**Chapter 5 §1: performative prediction (book text)**

- When predictions drive actions that change future data, the stationarity assumption may fail (Perdomo et al.). Repeated retraining *may* converge "за певних умов". This "не означає, що будь-яка система зворотного зв'язку нестійка". Source: [ch05 §1](../../ch05-open-scientific-problems-and-technological-frontiers.md).
- The teaching example of a self-defeating prophecy: after preventive interventions, "при наївному перенавчанні моделі на нових логах без каузальної розмітки система … знижує чутливість правил на наступній ітерації". A grid-rerouting example illustrates prediction-induced load shift. Both are labelled teaching examples, not measured accidents. Source: [ch05 §1](../../ch05-open-scientific-problems-and-technological-frontiers.md).
- Figure 5.1 models the "оператор перекалібрування системи" θ_{k+1} = argmin E_{Z∼D(θ_k)} ℓ(Z;θ), where θ_k is a "вектор параметрів правил і порогів". Definition 5.1 gives the W₁ sensitivity ε_perf. Theorem 5.1 and formulas (5.3)–(5.4) give the contraction q_perf = ε_perf·β/α_str < 1 as a *sufficient* condition. q_perf ≥ 1 "не є доказом розбіжності", and q_perf "не є показником фізичної безпеки". Source: [ch05 §1, Def. 5.1, Thm 5.1](../../ch05-open-scientific-problems-and-technological-frontiers.md).
- §4.2 notes that split conformal coverage holds under exchangeability. EnbPI and ACI need mixing or online-feedback assumptions. Evaluation should use chronological calibration and splits by unit and time. §4.3 defines recommendation utility and abstention as an open hypothesis, with risk–coverage and escalation share reported. Source: [ch05 §4.2–§4.3](../../ch05-open-scientific-problems-and-technological-frontiers.md).

### Inferences

**Objects that must evolve: book coverage versus local design evidence**

| Object | What the book already says | Local design evidence (research prototype, generalised) | Code (static) | Missing |
|---|---|---|---|---|
| Action schemas R (rules) | owner, version, positive review; Reg(r) = F if revoked (ch06 §4.1) | content SemVer changelog; supersedes / superseded-by links; bitemporal validity; superseding as undercutting defeater (changelog and federation decisions; superseding memo) | changelog by atom count; superseded-defeater date filter; anti-rollback/revocation | lifecycle states (proposal → quarantine → active → deprecated/superseded → archived); per-schema validity interval; version pinning in explanation traces; content-level diff |
| Policy Π (source hierarchy, conflicts, authority, β) | defined; owned by "власник політики"; escalation on C | level 4 (weeks–years), changed by engineers and the policy owner (self-organisation memo) | none found specific to Π | Π versioning, review and audit; impact analysis of changing β from strict to presumed; retro-evaluation of past decisions under the old Π |
| Predictors and calibration | conformal interval → state (6.7); marginal guarantee; performative risk (ch05 §1); ACI and chronological calibration (ch05 §4.2) | delayed-label states; OPE with support/ESS guard; anchor replay; drift → escalation only (continual-learning and telemetry memos) | none found | recalibration trigger and schedule; calibration-set versioning; censoring of outcomes after preventive action; causal labelling of interventions |
| Coverage of prohibitions (B) | declared-complete source set + verified detection completeness on the active extraction method (ch06 §5.1, §7.2) | hypothesis atoms never satisfy prohibition or permission queries (hypothesis-atom requirement); extraction quarantine | typed extraction quarantine | coverage certificate with scope, method version and expiry; re-verification when sources, extractor or lexicon change |
| Acceptance metrics and thresholds | FDR/FPR with CI; pre-registered thresholds; matching protocol; metrics as KB-version acceptance (Table 6.18) | grouped splits; mandatory slices; zero blocking failures; group bootstrap; label sanitising with human escalation (release-gate memo) | none found | cross-version regression gate (paired old vs new); threshold governance over time; reference-set (G) versioning and label-noise handling |
| Hypotheses and gaps | LLM ideas without schema → "гіпотези для людини"; U → information-gathering schemas | typed HYPOTHESIS atoms; gap tickets → dedup → harvesting (hypothesis-atom requirement, gap-ticket memo, gap-triage requirement) | hypothesis sampler; SimHash triage | promotion criteria from hypothesis to admitted schema; gap-closure-rate metric (ch06 §9.1 item 2 depends on it) |

- ch05's self-defeating prophecy directly threatens schema thresholds and predictor recalibration in ch06. Successful preventive actions censor failures, so naive outcome logging biases both FDR/FPR re-estimation and calibration. The continual-learning memo's delayed-label states plus a common-support guard are the closest local design answer, but they are unimplemented.
- A book-safe synthesis consistent with all sources: fast loops (runtime state, defeaters) act automatically, while changes to R, Π, B-coverage certificates, predictor versions and thresholds pass an admission gate. That gate is a paired regression on frozen anchors plus the section 8 metrics, and ends with an owner decision. This merges ch06 Def. 6.2 and Table 6.18 with the self-organisation memo (§2) and the continual-learning memo (§2).

### Gaps

- ch06 says nothing about delayed or censored outcomes of recommended actions, OPE, drift detection, recalibration cadence, retirement semantics for schemas, or version pinning in explanations. ch05 discusses performative retraining only abstractly.
- No hypothesis in ch06 is about evolution (e.g., "after k revisions FPR_perm does not increase beyond δ"). Possible additions: hypotheses on revision regression, coverage-certificate decay, and recalibration under intervention censoring.
- I did not read ch06 §§1–3.2 and §6 in full, nor ch05 §§2–3, §5 and §6. Any further statements there about change over time are not captured.

## 4. Where the general book (ExpertSystem/) covers knowledge maintenance, acquisition, versioning, belief revision, drift, continual learning and KB evolution

### Takeaway
The general book already covers knowledge evolution extensively:

- ch10: acquisition, admission, lifecycle and provenance;
- ch23: KB verification and historical regression;
- ch25: change routes, knowledge lifecycle, exam matrices, calibration, feedback control;
- ch26: continual learning, delayed feedback, log bias, OPE, drift types, forgetting, rule-update procedures, active learning, performativity, admission conveyor;
- ch35: reactive L0/L1 memory, truth maintenance, retraction;
- ch38: knowledge deficits and gaps.

The predictive book can cross-reference these chapters instead of re-deriving them.

### Cited Findings

- **ch10** «Системи здобуття знань: джерела, допуск і життєвий цикл». Relevant sections: §5 «Етапи конвеєра здобуття знань: від джерел до сертифікованих фактів», §6 «Експертне курування: розподіл повноважень між оператором та алгоритмом», §7 «Моделювання довіри та рівні достовірності джерел знань», §14 «Моделювання походження й родоводу знань (Data Provenance & Lineage)». Source: [ExpertSystem/ch10-knowledge-acquisition-systems.md](../../../ExpertSystem/ch10-knowledge-acquisition-systems.md).
- **ch23** «Верифікація бази знань…». Relevant sections: §3 typology of anomalies, §5 rule coverage metrics, §7 mutation testing, §9 formal methods, §12 «Регресійне тестування на історичних прецедентах та контроль витоку». Source: [ExpertSystem/ch23-knowledge-base-verification.md](../../../ExpertSystem/ch23-knowledge-base-verification.md).
- **ch25** «Як навчати експертну систему: екзаменаційні матриці, аудит знань та контроль регресій».
  - §1 separates three change routes: fact (assertion with source), rule (condition for a conclusion), and model parameter. "Отримання нового факту ще не схвалює правило", and approving a rule is "окремим рішенням відповідального фахівця".
  - §2 «Життєвий цикл знань: від первинного сигналу до випуску версії» presents the author's "концепцію життєвого циклу знань" with **proposal / quarantine / verified baseline version**. A version is a full dependency snapshot.
  - Further sections: §7 independent exam methodology; §8 calibration; §9 «Керований цикл зворотного зв'язку та запобігання деградації знань»; §12 rollout of the first controlled cycle.

  Source: [ExpertSystem/ch25-how-expert-systems-learn.md §1–§2, §7–§9, §12](../../../ExpertSystem/ch25-how-expert-systems-learn.md).
- **ch26** «Неперервне навчання (Continual Learning) на досвіді та подолання зсуву системного журналу». Sections:
  - §3 delayed feedback;
  - §4 log bias and OPE;
  - §5 diagnosing types of drift (6 hits for concept drift);
  - §6 catastrophic-forgetting metrics;
  - §7 retention strategies;
  - §8 «Процедури оновлення правил та бібліотеки прецедентів»;
  - §9 active learning and criteria for querying the expert;
  - §10 «Перформативність рекомендацій та замикання циклу даних», using regret (formula 26.25) for a fixed action set;
  - §11 limits of LLM-generated knowledge;
  - §12 «Кваліфікаційний конвеєр допуску кандидата до робочої версії». A developer prepares the candidate, independent data check it, and the process owner decides; "будь-який ранній провал залишає кандидата в карантині". The release unit is a compatible snapshot of knowledge, indexes and models.
  - §13 validating superiority of a new version on control samples.

  Source: [ExpertSystem/ch26-continual-learning.md](../../../ExpertSystem/ch26-continual-learning.md).
- **ch35** «Реактивна експертна система: події, відкликання й адаптація знань». Sections: §3 L0/L1 hybrid memory (immutable golden base plus delta graph), §4 truth maintenance with active defeaters, §5 Re-Search/Re-Ranking/Re-Thinking. Belief-revision terms appear 5 times and deprecation/retraction terms 41 times. Source: [ExpertSystem/ch35-self-organizing-expert-systems-synergetics-and-npu-runtime.md](../../../ExpertSystem/ch35-self-organizing-expert-systems-synergetics-and-npu-runtime.md).
- **ch38** «Машинні галюцинації та дефіцит знань…». Sections: §3 epistemic deficit and safe refusal, §4 abductive closures, §6 unlearning on verified knowledge. Source: [ExpertSystem/ch38-curing-machine-hallucinations-and-knowledge-deficits.md](../../../ExpertSystem/ch38-curing-machine-hallucinations-and-knowledge-deficits.md).
- **Other hits from a keyword count, not read in full:**
  - AGM or belief revision: ch07, ch08, ch26, ch33, ch35, ch36, ch38.
  - Concept drift: ch09, ch26.
  - Deprecation, revocation or superseded: heaviest in ch10 (49), ch35 (41), ch16 (27), ch25 (27), ch26 (17), ch40 (12), ch23 (11).
  - ch21 «Від рекомендації до дії…» (autonomy levels A0–A4, typed action contracts, human in the control loop) is the general-book counterpart of predictive ch06.

  Source: [ExpertSystem/ directory](../../../ExpertSystem/).

### Inferences

- ch25 §1–§2 and ch26 §12 already set out, in public-book language, the human or owner admission gate and the proposal/quarantine/baseline lifecycle. They are consistent with the self-organisation memo, not with the gap-ticket memo's autonomy claim. The predictive book can apply them to R, Π, predictors and B-coverage by reference.
- ch26 §10 treats performativity with regret over a fixed action set, while predictive ch05 §1 uses performative-stability contraction. A short bridge would avoid two disconnected treatments.

### Gaps

- I listed headings and counted keywords only. I did not verify how deeply ch26 §§3–5 and §13 treat delayed labels, OPE support guards, drift types or version superiority tests, nor whether they cover predictor recalibration or prohibition-coverage re-certification.
- I did not check whether ExpertSystem chapters use the forbidden internal names. That falls under the separate leakage gate (`scripts/inspect_project_leakage.py`), which was not run here.
