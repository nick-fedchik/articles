# Part III. Input Assessment, Knowledge Acquisition, and Linguistic Formalization

[← To Part II](part-02-knowledge-models.md) | [Table of Contents](README.md) | [To Part IV →](part-04-architecture-and-inference.md)

---

## Purpose of the Part

The purpose of this part is to demonstrate how to prepare verifiable knowledge for an expert system — spanning from initial skeptical assessment of incoming information and source provenance to expert elicitation, document parsing, and formal knowledge base construction:

1. **Input Assessment:** Disentangle source reliability, claim corroboration, evidence dependence, and authorization to utilize the result prior to parsing.
2. **Sources and Lifecycle:** Govern object passports, admission status, bitemporal validity, and primary source revocation.
3. **Expert Experience:** Elicit candidates equipped with operating context, exceptions, and validation procedures.
4. **Document Linguistic Analysis:** Apply parsers and local models while preserving original byte coordinates without obscuring ambiguities and extraction errors.
5. **Requirement and Rule Formalization:** Transform deontic modalities (`MUST`, `SHALL`, `MAY`) into explicit representations of predicates, grammars, and finite-state automata.

---

## Overview of the Theme and Interconnection of Chapters

Chapter 37 establishes the entry epistemic perimeter (L0–L3 funnel): a correctly parsed message may nonetheless originate from an untrusted source or remain causally dependent upon previously processed intelligence, so evidence assessment precedes pipeline ingestion. Chapters 10–11 establish provenance governance and the elicitation of specialist expertise. Chapter 12 preserves semantic content and byte-level provenance throughout structural and linguistic document analysis. Chapters 14–15 transform normative specifications into formal requirements, invariants, grammars, and state machines.

The CommonKADS framework segregates domain knowledge from inference operations and task control flow. Linguistic analyzers propose hypotheses, whereas ontology validation and evidence assessment address distinct verification planes. The deliverable of this part is a formalized candidate knowledge unit bounded by explicit invariants and cryptographic source lineage.

Knowledge engineering does not terminate at extraction. User query compilation and logical inference over the constructed base are examined in [Part IV](part-04-architecture-and-inference.md), formal verification and regression-free continual learning in [Part V](part-05-verification-and-learning.md), and immutable binary knowledge package release in [Chapter 32](ch32-high-performance-knowledge-packs-mmap-and-harvesting.md).

```mermaid
flowchart TD
    accTitle: Knowledge Engineering Pipeline of Part III
    accDescr: From entry epistemic source assessment through knowledge acquisition and linguistic analysis to knowledge base construction.

    E["<b>Raw Sources, Signals, and Experts</b>"] --> INPUT["<b>Chapter 37</b><br/>Input Assessment, Skepticism & Evidence Dependence"]
    INPUT --> KAS["<b>Chapters 10–11</b><br/>Knowledge Acquisition & Expert Elicitation"]
    KAS --> NLP["<b>Chapter 12</b><br/>Linguistic Analysis & Coordinate Preservation"]
    NLP --> EXT["<b>Chapters 14–15</b><br/>Detection of SHALL/MUST, Grammars & FSM"]
    EXT --> REVIEW["<b>Verification & Owner Authorization</b><br/>schema · provenance · applicability · test cases"]
    REVIEW --> KB["<b>Published Knowledge Base</b><br/>Predicates and Invariants"]

    classDef src fill:#fce4ec,stroke:#c2185b,stroke-width:2px,color:#880e4f;
    classDef step fill:#e0f2f1,stroke:#00796b,stroke-width:2px,color:#004d40;
    classDef out fill:#e8f5e9,stroke:#2e7d32,stroke-width:2px,color:#1b5e20;

    class E src;
    class INPUT,KAS,NLP,EXT,REVIEW step;
    class KB out;
```

---

## Chapters in This Part

### [Chapter 37. Input Information Assessment: Sources, Evidence, and Uncertainty](ch37-input-information-assessment-and-algorithmic-skepticism.md)

* **Abstract:** From raw physical signal to admissible inference: disentangling source reliability, claim corroboration, analytical confidence, and authorization to act. Survey of nine open-source intelligence doctrines, law enforcement analytical frameworks (UNODC, FBI, INTERPOL, Europol), and formal incident response standards (FIRST, NIST). Multimodal data streams (text, audio, RF signals, imagery) bound to primary carriers and transformation lineage; dependent corroboration, interval probabilities, algorithmic skepticism, and streaming versus batch verification. Architectural analysis of i2 Analyst's Notebook and specialized intelligence platforms separates marketed capabilities from evidential admission. An educational Go reference module verifies invariance under replicated observations, dependency tracking, and principled abstention.

### [Chapter 10. Knowledge Acquisition Systems: Sources, Admission, and Lifecycle](ch10-knowledge-acquisition-systems.md)

* **Abstract:** The contract between sources and consumers: object passports, independent verification and approval statuses, access control, bitemporal timelines, release gating, and revocation. Rolling back a database snapshot does not restore revoked authority for a compromised source.

### [Chapter 11. Knowledge Elicitation from Experts: Interviewing, Cognitive Maps, and Experience Formalization](ch11-knowledge-elicitation-from-experts.md)

* **Abstract:** Protocols for empirical knowledge elicitation, uncovering latent boundary exceptions, independent rule verification, and boundaries of statistical heuristics. CommonKADS cleanly partitions domain models, inference schemas, and task control; digital audit footprints identify candidates for expert consultation rather than constituting automated competency rankings.

### [Chapter 12. Linguistic Analysis and Local Models: Preserving Content and Provenance](ch12-linguistic-analysis-and-local-models.md)

* **Abstract:** Decoupling coordinates of the original artifact, derived text, and tokenizer spans. Structural document parsing, failure boundaries of statistical NLP, weak supervision, active learning strategies, and decoupled uncertainty estimation.

### [Chapter 14. Requirements and Modality Extraction: From Normative Text to Invariants](ch14-requirements-detection-and-formalization.md)

* **Abstract:** Document conventions and regulatory context, structural limitations of EARS templates, typed intermediate representations, and partial compilation. Trigger timestamps must not be confused with operational reaction deadlines; formula satisfiability does not establish physical product compliance.

### [Chapter 15. Knowledge Extraction and KB Construction: Facts, Grammars, and Automata](ch15-knowledge-extraction-and-kb-construction.md)

* **Abstract:** Facts, grammars, and deterministic state automata backed by verifiable provenance; citation span boundaries, physical units, and multi-edition document reconciliation. Decoupling admission, authorization, and formal release. Protégé and ROBOT for ontology validation; process mining and association rule discovery as generators of candidate hypotheses rather than normative mandates.

## Applied Analysis of Chapters 12, 14–15

Epistemic errors compound causally: a document parser drops a table header, an extractor misbinds a numerical threshold to an unrelated operational mode, a compiler masks an undefined temporal scope, and an admission gateway approves the record solely due to an authenticated citation. Substituting a larger language model does not eradicate these architectural failure modes. Rigorous evaluation therefore isolates processing stages while holding the remaining pipeline constant.

| Boundary | Methodological Challenge | Low-Cost Control Case |
|---|---|---|
| File → Text | Normalization can be non-invertible; PDF coordinates are derived | Unicode ligatures, repeated citations, truncated multi-byte UTF-8 sequences |
| Text → Candidate | Statistical classification does not constitute formal proof | Syntactic negation, swapped semantic roles, distinct numerals in conflicting roles |
| Candidate → Formula | Formal grammars and EARS templates cannot verify semantic correctness | Latency embedded in trigger condition, ambiguous deadlines, unreachable preconditions |
| Formula → Admission | Logical consistency does not guarantee operational applicability or authorization | Valid citation possessing revoked credentials or mismatched operational profile |
| Change → New Pack | An unchanged verbatim citation may acquire conflicting semantics in a revised document | Modified section header or newly introduced exception; delta vs. clean rebuild diff |

For each structural boundary, the chapters present targeted experimental designs:

1. [Chapter 37](ch37-input-information-assessment-and-algorithmic-skepticism.md): L0–L3 cascading screening, echo-chamber detection, and LOGO sensitivity analysis.
2. [Chapter 12](ch12-linguistic-analysis-and-local-models.md): Docling, Apache Tika, and OCR extraction; hierarchical document schemas, byte-granular source maps, and active learning for annotation.
3. [Chapters 14](ch14-requirements-detection-and-formalization.md) and [15](ch15-knowledge-extraction-and-kb-construction.md): Partial compilation, constraint satisfaction solvers, process event logs, and bit-level equivalence between incremental and clean builds.

These experimental configurations represent proposed benchmark protocols rather than pre-computed results. Independent ground-truth baselines, document family cross-validation, and segregated hold-out sets prevent evaluating models on their own generated labels. Synthetic test suites verify operational invariants, but cannot guarantee zero-defect execution across all arbitrary regulatory corpora.

---

[← To Part II](part-02-knowledge-models.md) | [Table of Contents](README.md) | [To Part IV →](part-04-architecture-and-inference.md)
