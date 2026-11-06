# Architecture of Evidence-Based Expert Systems

## Abstract

This book examines how engineering documents, requirements, source code, and test results can be transformed into structured knowledge for verifiable machine reasoning. Across seven parts and forty chapters, it covers the formalization of concepts and sources, knowledge models and acquisition, traceability, expert-system architecture, explanations, verification, safety, and neuro-symbolic approaches.

The book emphasizes evidence provenance, applicability limits, and explicit links among requirements, rules, design decisions, and tests. It distinguishes a proposed hypothesis, a conclusion derived within a formal model, and an experimentally confirmed result; plausibility alone is not treated as evidence. Topics include logical inference, uncertainty, causal and temporal models, explanations, controlled actions, continual learning, and safety arguments. Language models are treated as assistants whose proposals require independent validation.

The work is intended for students, engineers, and researchers who design or assess decision-support systems in complex technical environments. Its goal is to show how to build systems that not only produce conclusions, but also preserve their grounds, explain their limitations, and abstain when the available evidence is insufficient.

**Chapter abstracts:** [All 40 chapters](#chapter-abstracts).

**Keywords:** expert systems, knowledge engineering, evidence provenance, traceability, formal verification, functional safety, neuro-symbolic AI.

[Ukrainian master edition](../../README.md)

---

## Chapter abstracts

Concise abstracts for all 40 chapters of the Ukrainian master edition. Each chapter number links to its source text.

### 01 · [Chapter 1](../../ch01-introduction-to-expert-systems.md)

Defines an expert system through explicit knowledge, explainable inference, and qualified abstention. Surveys the field's history and architecture, distinguishing accountable reasoning from document retrieval and unconstrained text generation.

### 02 · [Chapter 2](../../ch02-epistemology-of-machine-knowledge.md)

Examines the grounds for a machine claim: source, context, scope, uncertainty, and validity. Conflicting temperature limits show why a correct number alone cannot justify a decision.

### 03 · [Chapter 3](../../ch03-beyond-reference-information-systems.md)

Distinguishes reference, search, dialogue, recommendation, and expert systems by their tasks and output contracts. Finding a rule is not the same as applying it correctly to a specific situation.

### 04 · [Chapter 4](../../ch04-evolution-from-bayes-to-evidence-ai.md)

Traces the evolution of expert systems and compares logic, probability, fuzzy sets, graphs, and case-based reasoning. Method selection is tied to the question, required evidence, and limits of inference.

### 05 · [Chapter 5](../../ch05-triad-of-trust-and-corporate-memory.md)

Presents a trust triad: the system applies knowledge, a recommendation exposes its grounds, and organizational memory preserves context and outcomes. It covers provenance, audit, access control, and knowledge review.

### 06 · [Chapter 6](../../ch06-applied-mathematics-for-expert-systems.md)

Introduces logical and probabilistic inference, graph methods, optimization, and causal analysis through engineering tasks. Reproducible calculations and accurate quotations still depend on valid inputs and interpretation.

### 07 · [Chapter 7](../../ch07-knowledge-base-typology.md)

Compares rules, object models, ontologies, cases, constraints, probabilistic models, and document or vector indexes. It explains what each representation stores and why missing knowledge must remain explicit.

### 08 · [Chapter 8](../../ch08-engineering-artifacts-as-data.md)

Models requirements, decisions, plans, and test results as structured artifacts with typed fields and links. Provenance supports querying and analysis but does not establish truth or legal sufficiency.

### 09 · [Chapter 9](../../ch09-engineering-knowledge-graph-traceability.md)

Builds a typed knowledge graph to trace changes from requirements through software, hardware, and tests. Queries support impact analysis, while graph links alone do not prove correctness.

### 10 · [Chapter 10](../../ch10-knowledge-acquisition-systems.md)

Describes a knowledge-acquisition pipeline that locates sources, preserves content and provenance, checks candidate knowledge, and defines its scope. Citation fidelity, detector confidence, and document approval answer different questions.

### 11 · [Chapter 11](../../ch11-knowledge-elicitation-from-experts.md)

Shows how interviews, case comparisons, decision reconstruction, cognitive maps, and counterexamples elicit tacit expertise. Frequency, uncertainty, and agreement among experts are treated as distinct measurements.

### 12 · [Chapter 12](../../ch12-linguistic-analysis-and-local-models.md)

Explains how linguistic analysis and local models can form rule candidates while preserving source bytes, conditions, negation, and privacy. An exact quotation alone does not guarantee correct interpretation.

### 13 · [Chapter 13](../../ch13-language-variability-vs-determinism.md)

Separates candidate retrieval from deterministic semantic validation. It maps varied phrasing to typed queries while checking units, negation, modality, source version, and evidence.

### 14 · [Chapter 14](../../ch14-requirements-detection-and-formalization.md)

Shows how to identify requirement modality and scope, transform normative text into formal candidates, and inspect contradictions. EARS patterns and limited tools aid review but do not replace it.

### 15 · [Chapter 15](../../ch15-knowledge-extraction-and-kb-construction.md)

Follows knowledge from source to candidate, approved knowledge-base record, and verifiable inference. It covers extraction stages, annotation coordinates, units, and import checks; proposed methods still require evaluation.

### 16 · [Chapter 16](../../ch16-expert-systems-architecture.md)

Explains boundaries among knowledge stores, inference, verification, access policy, and provenance logs. It covers updates and revocations, while retrieval scores are not treated as evidence.

### 17 · [Chapter 17](../../ch17-implementation-stack.md)

Frames technology selection around required operations, correctness criteria, and failure behavior. It compares languages, rule engines, storage, and deployment choices, emphasizing explicit handling of failure.

### 18 · [Chapter 18](../../ch18-execution-infrastructure.md)

Connects architecture to physical execution constraints: memory, CPUs and accelerators, buses, latency, and deployment placement. Resource estimates are conditional, not universal performance guarantees.

### 19 · [Chapter 19](../../ch19-from-question-to-evidence.md)

Turns a query into a checked claim by finding sources and binding the correct version, object, unit, and condition. Search relevance alone is insufficient to support a conclusion.

### 20 · [Chapter 20](../../ch20-explanation-engine.md)

Builds explanations from decision traces for how, why, why-not, what-if, and change questions. It covers minimal conflicts and role-based disclosure, and prevents language models from inventing reasons.

### 21 · [Chapter 21](../../ch21-from-recommendation-to-action.md)

Separates a recommendation from authority to act. A proposed action gateway checks the action, actor, context, and policy, then records execution; the design is not presented as validated.

### 22 · [Chapter 22](../../ch22-cybernetics-edge-to-backend.md)

Uses an industrial sensor scenario to explain context-aware cybernetic feedback across sensing, edge analysis, and response. Model scores alone cannot determine action; calibration and process state matter.

### 23 · [Chapter 23](../../ch23-knowledge-base-verification.md)

Distinguishes verification of the encoded knowledge base from validation of the policy against the real task. Finite models, counterexamples, and release gates reveal omissions but do not prove universal correctness.

### 24 · [Chapter 24](../../ch24-system-diagnosis.md)

Uses model-based diagnosis to compare fault hypotheses with observations under incomplete data. Consistency narrows candidates but does not establish a physical root cause without separate evidence.

### 25 · [Chapter 25](../../ch25-how-expert-systems-learn.md)

Describes examination matrices, knowledge audits, regression tests, and calibration checks for controlled updates. They detect degradation in tested scenarios but cannot guarantee correctness in every future case.

### 26 · [Chapter 26](../../ch26-continual-learning.md)

Explains feedback bias when a system's recommendations shape the logs used for learning. Monitoring and controlled updates help prevent repeated choices from being mistaken for independent evidence.

### 27 · [Chapter 27](../../ch27-safety-case-gsn-synthesis.md)

Combines traceable evidence with structured safety arguments to create a reviewable case. The result is a candidate argument, not an automatic certificate; assumptions and gaps require scrutiny.

### 28 · [Chapter 28](../../ch28-dual-mode-expert-systems.md)

Separates strict conclusions grounded in verified knowledge from advisory hypotheses about what to inspect next. Explicit modes, labels, and admission controls prevent suggestions from masquerading as proof.

### 29 · [Chapter 29](../../ch29-neuro-symbolic-architecture.md)

A language model interprets a query and proposes passages, while a symbolic core checks structure, provenance, rules, and exceptions. Evidence for a conclusion and authority to act remain separate questions.

### 30 · [Chapter 30](../../ch30-safety-cybersecurity-co-engineering.md)

Coordinates functional-safety and cybersecurity analysis when a protective action can itself create physical risk. Shared context and impact analysis expose trade-offs without letting one discipline override the other.

### 31 · [Chapter 31](../../ch31-syllogistic-reasoning-and-relation-lattices.md)

Examines normative inference using predicate hierarchies, exceptions, rebuttals, and temporal validity. It distinguishes retrieving a relevant passage from deriving a context-specific, rule-governed conclusion.

### 32 · [Chapter 32](../../ch32-high-performance-knowledge-packs-mmap-and-harvesting.md)

Presents immutable knowledge packs, byte-level admission, indexes, and memory mapping for reproducible and efficient access. Hashes protect integrity but do not establish the truth or applicability of content.

### 33 · [Chapter 33](../../ch33-inter-system-knowledge-exchange-and-model-teaching.md)

Defines safe knowledge exchange by carrying units, scope, rule versions, sources, and permissions with each value. Signatures protect origin, not correctness or applicability.

### 34 · [Chapter 34](../../ch34-deterministic-relational-analysis-abduction-and-socratic-dialogue.md)

Uses relational search and abduction to surface candidate explanations when facts are missing, then asks targeted questions for discriminating evidence. Hypotheses remain distinct from established facts.

### 35 · [Chapter 35](../../ch35-self-organizing-expert-systems-synergetics-and-npu-runtime.md)

Explains reactive systems that revise conclusions when supporting observations are withdrawn. It separates immutable rules from changing evidence and treats adaptive parameters as proposals requiring bounds and review.

### 36 · [Chapter 36](../../ch36-knowledge-testing-pyramid-and-variational-calibration.md)

Adapts the testing-pyramid idea to knowledge systems: rule, interaction, end-to-end, and input-perturbation tests. It clarifies what each level can reveal and what it cannot prove.

### 37 · [Chapter 37](../../ch37-input-information-assessment-and-algorithmic-skepticism.md)

Evaluates incoming evidence by object identity, time, provenance, independence, and uncertainty before applying rules. Suspicion triggers verification; it is not proof of falsification.

### 38 · [Chapter 38](../../ch38-curing-machine-hallucinations-and-knowledge-deficits.md)

Separates language-model proposals from accepted conclusions through checks of source, meaning, applicability, and premises. Authentic citations are necessary for traceability but are not sufficient to prove a conclusion.

### 39 · [Chapter 39](../../ch39-active-compliance-auditor-and-popperian-testing.md)

Adds an active auditor that seeks counterexamples and missing boundary cases for requirements. It separates test generation from authority over expected results and treats compliance coverage as auditable, not automatic proof.

### 40 · [Chapter 40](../../ch40-distributed-epistemic-architectures-soa-and-cluster-scaling.md)

Describes distributed expert systems with semantic routing, active rule-set closure, memory coordination, and decision logs. Missing shards or stale versions should yield unresolved states, not silent approval.
