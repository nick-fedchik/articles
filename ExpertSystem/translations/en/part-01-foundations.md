# Part I: Conceptual and Epistemic Foundations

[← Back to Table of Contents](README.md) | [To Part II →](part-02-knowledge-models.md)

---

## Purpose of This Part

To determine when an expert system is required and under what conditions its propositions qualify as actionable for engineering decisions. This part demarcates raw data, evidence, rules, hypotheses, and conclusions, establishing responsibility and provenance for generated answers.

---

## Abstract and Interconnection of Chapters

Chapter 1 formulates the practical problem of verifiable answers. Chapter 2 formalizes the knowledge contract: provenance, applicability bounds, temporal validity, access control, and inference mechanisms. Chapter 3 examines the fundamental distinction between this epistemic contract and standard information retrieval systems. Chapter 4 traces the historical evolution of inference methods across diverse dimensions of uncertainty. Chapter 5 unites conclusions, recommendations, and institutional decision memory.

Key Takeaway: The reader can classify required system architectures and pinpoint missing grounds necessary for verifiable decision-making. The contract cannot guarantee the factual truth of flawed sources; mathematical formalisms and knowledge representation are explored in the subsequent part.

```mermaid
flowchart LR
    accTitle: Logical Progression of Part I
    accDescr: Progression from basic concepts of expert systems through knowledge philosophy and search critique to Bayesian evolution and the memory triad.

    CH1["<b>Chapter 1</b><br/>Introduction & Core Concepts"] --> CH2["<b>Chapter 2</b><br/>Philosophy of Machine Knowledge"]
    CH2 --> CH3["<b>Chapter 3</b><br/>Limits of Search Systems"]
    CH3 --> CH4["<b>Chapter 4</b><br/>Evolution: Bayes to AI"]
    CH4 --> CH5["<b>Chapter 5</b><br/>Triad of Trust & Memory"]

    classDef step fill:#e3f2fd,stroke:#1565c0,stroke-width:2px,color:#0d47a1;
    classDef focus fill:#fff3e0,stroke:#e65100,stroke-width:2px,color:#bf360c;

    class CH1,CH3,CH4,CH5 step;
    class CH2 focus;
```

---

## Chapters in This Part

### [Chapter 1. Introduction to Expert Systems: From Chaos to Governed Knowledge](ch01-introduction-to-expert-systems.md)

* **Abstract:** One engineering query and three distinct answers: from search engines, from generative chatbots, and from an evidence-governed expert system. Why expert systems have experienced an engineering renaissance: plausible prose has become commoditized, compliance requires auditable decisions, and institutional knowledge departs teams faster than physical products. Historical engineering causes of the 1980s AI winter. Definitions of knowledge bases, inference engines, justification packages, five industrial application domains, boundaries of liability, and recommended reading pathways.

### [Chapter 2. Philosophy for the Engineer: What a Machine May Lawfully Call Knowledge](ch02-epistemology-of-machine-knowledge.md)

* **Abstract:** Seven philosophical dilemmas translated into rigorous verification checks for automated answers: evidentiary grounding, entity identity, logic, intent, quotation context, assertion category, and explicit authorization. A procedural contract separating auditable knowledge from plausible generation without turning flawed sources into empirical truth.

### [Chapter 3. How an Expert System Differs from an Information Retrieval System](ch03-beyond-reference-information-systems.md)

* **Abstract:** The boundary between retrieval systems and expert systems across three architectural dimensions: queries over document corpora versus case-specific situations, document summaries versus expert decisions with justification packages, search pipelines versus knowledge bases with inference engines. The role of Large Language Models and Retrieval-Augmented Generation (RAG), comparison of six information system classes under epistemic contracts, and three operational rules: controlled knowledge modification, confidence calibration, and autonomous action boundaries.

### [Chapter 4. Evolution of Expert Systems: From Bayes' Theorem to Evidence-Based AI](ch04-evolution-from-bayes-to-evidence-ai.md)

* **Abstract:** The history of production rules, probabilities, fuzzy logic, knowledge graphs, and case-based reasoning. For each method, we identify its specific uncertainty type and operational limits; rigorous mathematical models are detailed in Chapter 6. Crucial lessons from earlier AI winters regarding knowledge elicitation bottlenecks, fragility, maintenance burden, and total cost of ownership.

### [Chapter 5. The Triad of Trust: Expert System, Evidence-Based Recommendation, and Corporate Memory](ch05-triad-of-trust-and-corporate-memory.md)

* **Abstract:** The expert system applies knowledge, the recommendation delineates evidentiary grounds and validity bounds, while corporate memory records rationale, versioning, and decision ownership. A realistic case study of engineering team turnover illustrates why none of these three components can replace the other two.

---

[← Back to Table of Contents](README.md) | [To Part II →](part-02-knowledge-models.md)
