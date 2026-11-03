# Automatic generation of candidate knowledge (rules, exceptions, laws) for predictive expert systems and its validation before admission to the knowledge base

Compiled 2026-11-03. Evidence-status tags used below:
**[ESTABLISHED]** — peer-reviewed result that is widely cited and not contradicted by later sources found here;
**[CONDITIONAL]** — peer-reviewed result that holds only under stated assumptions (data regime, domain, benchmark) or comes from a single study;
**[PROPOSAL]** — method, framework or roadmap without independent validation;
**[PREPRINT]** — not peer-reviewed as of the check date.
Access level: **(abstract-verified)** means the claim was checked against the publisher, arXiv or OpenAlex abstract; **(metadata-only)** means only the title, venue, volume, pages and DOI were verified through Crossref or OpenAlex, and the full text was not read. No claim below relies on full text unless stated. Bibliographic metadata for every DOI was checked through Crossref or OpenAlex on 2026-11-03. Citation counts are OpenAlex values retrieved on 2026-11-03.

## 1. Evolutionary computation for knowledge: LCS/XCS, genetic programming, genetic fuzzy systems — what they learn, accuracy vs interpretability, PHM applications

### Takeaway
Evolutionary methods learn explicit, inspectable knowledge: condition–action rules with accuracy-based fitness (XCS), symbolic expressions or health-indicator formulas (GP), and fuzzy rule bases (GFS). In prognostics and health management (PHM), the strongest documented use is genetic programming that constructs health indicators or explicit life expressions. Genetic fuzzy rule bases have documented use in fault diagnosis. This targeted search found no peer-reviewed application of Learning Classifier Systems to RUL or prognostics. Multi-objective search that explicitly trades accuracy against rule-base complexity is the established answer to the accuracy–interpretability tension.

### Cited Findings
- **XCS (Wilson, 1995) [ESTABLISHED] (abstract-verified):**
  - Each classifier stores a payoff prediction, but its GA fitness is the *accuracy* of that prediction rather than its strength.
  - The GA runs in niches defined by match sets rather than across the whole population.
  - The population therefore tends to form "a complete and accurate mapping X × A → P" from inputs and actions to payoff predictions.
  - The population also tends toward classifiers that are "maximally general, subject to an accuracy criterion".
  - Source: Evolutionary Computation 3(2):149–175. — [Wilson 1995, doi:10.1162/evco.1995.3.2.149](https://doi.org/10.1162/evco.1995.3.2.149)
- **LCS review (Urbanowicz & Moore, 2009) [ESTABLISHED] (abstract-verified):**
  - Learning Classifier Systems (LCS) are rule-based machine-learning algorithms derived from evolutionary biology and AI.
  - Documented application areas include robotics, classification, knowledge discovery and modelling.
  - Epidemiology is named as a field of growing interest because LCS can support etiological discovery.
  - The authors state that "implementation optimization is nontrivial" and provide a roadmap of algorithmic components.
  - Source: J. Artificial Evolution and Applications 2009:1–25. — [Urbanowicz & Moore 2009, doi:10.1155/2009/736398](https://doi.org/10.1155/2009/736398)
- **LCS textbook (metadata-only):** Urbanowicz & Browne, *Introduction to Learning Classifier Systems*, SpringerBriefs in Intelligent Systems, 2017. — [doi:10.1007/978-3-662-55007-6](https://doi.org/10.1007/978-3-662-55007-6)
- **Genetic programming foundations (metadata-only):** Koza, "Genetic programming as a means for programming computers by natural selection", Statistics and Computing 4 (1994). — [doi:10.1007/BF00175355](https://doi.org/10.1007/BF00175355)
- **Liao (2013), GP-discovered prognostic features [CONDITIONAL; single study] (abstract-verified):**
  - Problem addressed: RUL prediction is hard when a health-indicator trend is not visible during the whole life or appears only just before failure.
  - Method: genetic programming automatically discovers "advanced features", i.e. mathematical combinations of extracted features, that capture fault progression for RUL prediction.
  - Source: IEEE Trans. Industrial Electronics; 275 citations in OpenAlex. — [Liao 2013, doi:10.1109/tie.2013.2270212](https://doi.org/10.1109/tie.2013.2270212)
- **Ding et al. (2019), explicit RUL models for slewing bearings [CONDITIONAL] (abstract-verified):**
  - Symbolic regression with several GP variants builds explicit life expressions from degradation indicators.
  - The authors validate on three run-to-failure experiments under different operating conditions.
  - Epigenetic linear GP with hybrid and piecewise strategies plus similarity-based combination performed best.
  - Its life expressions were "more succinct and intelligible".
  - Source: IEEE Access 2019. — [Ding et al. 2019, doi:10.1109/access.2019.2919663](https://doi.org/10.1109/access.2019.2919663)
- **Toffolo & Lazzaretto (2008), genetically evolved diagnostic rules [CONDITIONAL] (abstract-verified):**
  - The knowledge base of a fuzzy diagnostic expert system is generated by an evolutionary algorithm instead of being elicited from plant personnel.
  - The search uses two objectives: minimum rule-set complexity and maximum correct fault identifications.
  - The method was applied to a two-shaft gas turbine in a combined-cycle cogeneration plant "to evaluate the potentialities and the limits".
  - Source: Int. J. of Thermodynamics. — [Toffolo & Lazzaretto 2008, doi:10.5541/ijot.1034000219](https://doi.org/10.5541/ijot.1034000219)
- **Lemos, Caminhas & Gomide, evolving fuzzy classifier for fault detection (metadata-only):**
  - Title: "Adaptive fault detection and diagnosis using an evolving fuzzy classifier".
  - Source: Information Sciences 220:64–85; OpenAlex gives 2011 as the publication year. — [doi:10.1016/j.ins.2011.08.030](https://doi.org/10.1016/j.ins.2011.08.030)
- **Genetic fuzzy systems: taxonomy and interpretability (metadata-only):**
  - Herrera, "Genetic fuzzy systems: taxonomy, current research trends and prospects", Evolutionary Intelligence 1:27–46 (2008). — [doi:10.1007/s12065-007-0001-5](https://doi.org/10.1007/s12065-007-0001-5)
  - Cordón, "A historical review of evolutionary learning methods for Mamdani-type fuzzy rule-based systems: designing interpretable genetic fuzzy systems", Int. J. Approximate Reasoning 52:894–913 (2011). — [doi:10.1016/j.ijar.2011.03.004](https://doi.org/10.1016/j.ijar.2011.03.004)
  - Gacto, Alcalá & Herrera, "Interpretability of linguistic fuzzy rule-based systems: an overview of interpretability measures", Information Sciences 181:4340–4360 (2011). — [doi:10.1016/j.ins.2011.02.021](https://doi.org/10.1016/j.ins.2011.02.021)
  - Ishibuchi & Nojima, "Analysis of interpretability-accuracy tradeoff of fuzzy systems by multiobjective fuzzy genetics-based machine learning", Int. J. Approximate Reasoning (OpenAlex year 2006). — [doi:10.1016/j.ijar.2006.01.004](https://doi.org/10.1016/j.ijar.2006.01.004)
- **Multi-objective evolutionary fuzzy systems (Fazzolari et al.) [ESTABLISHED as a research line] (abstract-verified):**
  - Genetic and evolutionary design of fuzzy systems has been extended to multi-objective algorithms that "consider multiple conflicting objectives, instead of a single one".
  - The review provides a two-level taxonomy of these systems.
  - Source: IEEE Trans. Fuzzy Systems; OpenAlex gives 2012 as the publication year. — [Fazzolari, Alcalá, Nojima, Ishibuchi, Herrera, doi:10.1109/TFUZZ.2012.2201338](https://doi.org/10.1109/TFUZZ.2012.2201338)

### Inferences
- The accuracy-based fitness and "maximally general subject to accuracy" bias in XCS correspond to a desirable admission property. Each candidate rule carries its own measured prediction-error statistic and generality, which an admission gate could threshold. **[inference]**
- PHM evidence for GP and symbolic regression mostly covers single assets or test rigs: Liao used one study and Ding used three run-to-failure experiments. This supports the claim that GP can produce intelligible health indicators or life expressions, **conditionally**. It does not yet support fleet-level generalisation. **[inference]**
- The genetic-fuzzy literature turns accuracy versus interpretability into a Pareto front with complexity as an explicit objective, as in Toffolo & Lazzaretto and Fazzolari et al. In a predictive expert system, the complexity objective can be reused as an admission criterion: do not admit a rule set that is dominated on the front. **[inference / proposal]**

### Gaps
- This bounded search (OpenAlex, about 20 queries) found **no peer-reviewed application of XCS, UCS or ExSTraCS to RUL, prognostics or condition monitoring**. "LCS fault diagnosis" queries returned generic ML classifiers. This is not proof of absence. A dedicated search of GECCO and IWLCS proceedings is needed.
- Quantitative accuracy-versus-interpretability numbers from Gacto 2011, Cordón 2011 and Ishibuchi & Nojima were not extracted; only metadata was accessible. The direction of the trade-off is supported here only by the Fazzolari and Toffolo abstracts.
- The full text of Koza 1994 was not consulted. GP's known bloat and overfitting issues are not cited here for that reason.

## 2. Inductive rule learning and inductive logic programming: RIPPER, Progol, Metagol/Popper, ∂ILP — strengths, data requirements, noise handling

### Takeaway
ILP induces first-order rules from examples plus background knowledge. Its main strengths are data efficiency and auditability.

Classical ILP and the basic learning-from-failures setting are not tolerant of noise or mislabelling. Differentiable ILP (∂ILP) was proposed specifically to add noise robustness and neural front-ends.

ILP has old reports of satellite fault-diagnosis rules. Recent PHM-specific evidence is sparse and conference-level.

### Cited Findings
- **RIPPER (metadata-only):** Cohen, "Fast Effective Rule Induction", Machine Learning Proceedings 1995 (ICML), pp. 115–123. — [doi:10.1016/B978-1-55860-377-6.50023-2](https://doi.org/10.1016/B978-1-55860-377-6.50023-2)
- **ILP overview (Muggleton, 1994) [ESTABLISHED] (abstract-verified):**
  - ILP constructs "first-order definite clause theories from examples and background knowledge".
  - Real-world successes include "fault diagnosis rules for satellites", drug-design structure–activity rules, finite-element mesh design and protein secondary structure.
  - Theoretical support includes learning-in-the-limit results and some PAC-learning results.
  - "ILP is underconstrained by its logical specification", so ILP systems need extra-logical constraints such as language bias.
  - Source: ACM SIGART Bulletin 1994. — [Muggleton 1994, doi:10.1145/181668.181671](https://doi.org/10.1145/181668.181671)
- **Progol (metadata-only):** Muggleton, "Inverse entailment and Progol", New Generation Computing 13:245–286 (1995). — [doi:10.1007/BF03037227](https://doi.org/10.1007/BF03037227)
- **Metagol / meta-interpretive learning (metadata-only):** Muggleton, Lin & Tamaddoni-Nezhad, "Meta-interpretive learning of higher-order dyadic datalog: predicate invention revisited", Machine Learning 100:49–73 (2015). — [doi:10.1007/s10994-014-5471-y](https://doi.org/10.1007/s10994-014-5471-y)
- **Popper / learning from failures (Cropper & Morel, 2021) [ESTABLISHED] (abstract-verified):**
  - The learner uses a **generate → test → constrain** loop.
  - A hypothesis *fails* if it does not entail all positive examples or entails any negative example.
  - Constraints learned from failures prune generalisations of over-general hypotheses and specialisations of over-specific hypotheses.
  - Popper combines answer-set programming and Prolog, and learns textually minimal and recursive programs.
  - Source: Machine Learning 110:801–856. — [Cropper & Morel 2021, doi:10.1007/s10994-020-05934-z](https://doi.org/10.1007/s10994-020-05934-z)
- **ILP survey (Cropper & Dumančić, 2022) [ESTABLISHED] (abstract-verified):**
  - The goal of ILP is to induce a hypothesis, "a set of logical rules", that generalises the training examples.
  - The survey compares systems and describes Aleph, TILDE, ASPAL and Metagol.
  - It also summarises current limitations.
  - Source: JAIR 74:765–850. — [Cropper & Dumančić 2022, doi:10.1613/jair.1.13507](https://doi.org/10.1613/jair.1.13507)
- **Differentiable ILP (Evans & Grefenstette, 2018) [ESTABLISHED as a proposal; CONDITIONAL for robustness] (abstract-verified):**
  - The authors characterise classical ILP as "extremely data-efficient" but "not robust to noise in or mislabelling of inputs" and unable to handle non-symbolic data such as raw pixels.
  - ∂ILP is trained by backpropagation against a likelihood objective and "shows a robustness to noise and error in the training data which ILP cannot cope with".
  - It can be combined with neural networks over ambiguous inputs.
  - Source: JAIR 61:1–64. — [Evans & Grefenstette 2018, doi:10.1613/jair.5714](https://doi.org/10.1613/jair.5714)
- **ILP for diagnosing a distributed cooperative system (Sato et al., 2020) [CONDITIONAL; conference, case study] (abstract-verified):**
  - The ACEDIA ILP system learns transition rules and same-time relationship rules from time series of individual subsystems.
  - During diagnosis, rules and specifications violated by diagnostic data are ranked by generality.
  - In a simplified automotive multi-controller case study, "the rules essentially related to the error were ranked higher".
  - Source: IEEE ICPHM 2020. — [Sato et al. 2020, doi:10.1109/icphm49022.2020.9187032](https://doi.org/10.1109/icphm49022.2020.9187032)

### Inferences
- Popper's "test" step is already an admission check of the form *entails all positives and no negatives*. In a mission-critical knowledge base, this logical-consistency check is necessary but not sufficient, because noisy sensor labels make strict consistency unattainable. A statistical tolerance step, as in Question 5, must be added. **[inference]**
- ILP needs symbolic background predicates. For PHM, continuous sensor streams must first be discretised into events or predicates such as thresholds or state transitions. Sato et al. do this through transition rules, and the quality of that discretisation is itself a failure point. **[inference]**
- ∂ILP's noise robustness is reported in its own benchmarks. Independent PHM-scale evidence was not found, so it remains **conditional**. **[inference]**

### Gaps
- RIPPER's reported accuracy and runtime relative to C4.5rules and its noise behaviour could not be verified without the full text; only metadata was confirmed.
- Scalability limits of ∂ILP, including memory growth with the template or rule space, are widely discussed but were not verified here.
- Noise-tolerant extensions of Popper exist in later literature but were not verified in this session.
- No peer-reviewed ILP application to RUL or degradation trajectories was found, only diagnosis.

## 3. Discovering physical or degradation laws from data: symbolic regression, SINDy, PySR — use for degradation and RUL, and failure modes

### Takeaway
Symbolic regression (Schmidt & Lipson 2009) and SINDy (Brunton et al. 2016) are established for recovering parsimonious governing equations from clean or moderately noisy data on canonical systems.

PHM uses exist: GP-based RUL expressions, SINDy-based battery state-of-health (SOH) models, and physics-informed SR for tool wear. These are mostly single-study and asset-specific.

Documented failure modes include sensitivity to noise and low data, dependence on the candidate library, distortion from data scaling or normalisation, and inflated benchmark results. Those results are inflated by memorisation in LLM-assisted SR, or by recovering only common equations.

### Cited Findings
- **Schmidt & Lipson (2009) [ESTABLISHED] (abstract-verified):**
  - The paper proposes "a principle for the identification of nontriviality" for correlations found in data.
  - Without prior physics knowledge, the method recovered Hamiltonians, Lagrangians and conservation laws from motion-tracking data. The systems ranged from harmonic oscillators to chaotic double pendula.
  - Laws found for simpler systems were used to bootstrap more complex ones.
  - Source: Science 324(5923):81–85. — [Schmidt & Lipson 2009, doi:10.1126/science.1165893](https://doi.org/10.1126/science.1165893)
- **SINDy (Brunton, Proctor & Kutz, 2016) [ESTABLISHED] (abstract-verified):**
  - The method assumes the dynamics are sparse "in the space of possible functions"; this assumption "holds for many physical systems in an appropriate basis".
  - Sparse regression finds "the fewest terms" and produces "parsimonious models that balance accuracy with model complexity to avoid overfitting".
  - Demonstrations include oscillators, the chaotic Lorenz system and vortex shedding.
  - The method extends to parameterised, time-varying and externally forced systems.
  - Source: PNAS 113(15):3932–3937. — [Brunton et al. 2016, doi:10.1073/pnas.1517384113](https://doi.org/10.1073/pnas.1517384113)
- **SINDy-PI (2020) [ESTABLISHED] (abstract-verified):**
  - Earlier implicit and rational SINDy extensions are "extremely sensitive to noise".
  - SINDy-PI adds principled model selection.
  - It is reported as "several orders of magnitude more noise robust than previous approaches".
  - Source: Proc. R. Soc. A 476:20200279. — [Kaheman, Kutz & Brunton 2020, doi:10.1098/rspa.2020.0279](https://doi.org/10.1098/rspa.2020.0279)
- **Ensemble-SINDy (2022) [ESTABLISHED] (abstract-verified):**
  - Sparse identification "is sensitive to noise, especially in the low-data limit".
  - Bagging gives *inclusion probabilities* for candidate terms, enabling uncertainty quantification and probabilistic forecasts.
  - E-SINDy recovers PDE models from data with "more than twice as much measurement noise as has been previously reported".
  - Source: Proc. R. Soc. A, 2022. — [Fasel, Kutz, B. W. Brunton & S. L. Brunton 2022, doi:10.1098/rspa.2021.0904](https://doi.org/10.1098/rspa.2021.0904)
- **Data-scaling failure mode (Raut, Wilke & Schmidt, 2026) [PREPRINT] (abstract-verified):**
  - Data normalisation "can severely distort the discovery of governing equations by magnitude-based sparse regression".
  - SINDy is especially vulnerable because scaling interacts with noise. The resulting models "can be dense, uninterpretable, and physically incorrect".
  - The proposed STCV method replaces magnitude-based thresholding with a dimensionless statistical presence metric.
  - Source: arXiv 2603.05201, March 2026, not peer-reviewed. — [arXiv:2603.05201](https://doi.org/10.48550/arxiv.2603.05201)
- **SRBench (La Cava et al., 2021) [ESTABLISHED as a benchmark] (abstract-verified; venue verified through OpenReview):**
  - The authors evaluate 14 symbolic-regression methods and 7 machine-learning methods on 252 problems.
  - On real-world data, the best methods "combine genetic algorithms with parameter estimation and/or semantic search drivers".
  - For exact recovery under noise, "deep learning and genetic algorithm-based approaches perform similarly".
  - Source: NeurIPS Datasets & Benchmarks 2021. — [arXiv:2107.14351](https://arxiv.org/abs/2107.14351); [OpenReview listing](https://openreview.net/)
- **PySR (Cranmer, 2023) [PREPRINT] (abstract-verified):**
  - PySR uses a multi-population evolutionary algorithm with an "evolve-simplify-optimize loop" for constants in new expressions.
  - Its Julia back end is SymbolicRegression.jl.
  - The paper introduces "EmpiricalBench" to measure recovery of historical empirical equations.
  - Source: arXiv 2305.01582; no journal reference on arXiv as of 2026-11-03. — [arXiv:2305.01582](https://arxiv.org/abs/2305.01582)
- **PHM uses:**
  - **GP features for RUL, Liao 2013:** see Question 1. — [doi:10.1109/tie.2013.2270212](https://doi.org/10.1109/tie.2013.2270212)
  - **Explicit slewing-bearing RUL models, Ding et al. 2019:** see Question 1. — [doi:10.1109/access.2019.2919663](https://doi.org/10.1109/access.2019.2919663)
  - **SINDy battery SOH, Lee et al. 2025 [CONDITIONAL] (abstract-verified):**
    - The model formulation uses domain knowledge plus statistical variables.
    - The authors describe this as the first use of statistical variables in SINDy for SOH estimation.
    - The approach is reported as "approximately 30 times faster than GPR"; accuracy figures were not in the abstract.
    - Source: J. Intelligent & Robotic Systems 2025. — [Lee et al. 2025, doi:10.1007/s10846-025-02274-9](https://doi.org/10.1007/s10846-025-02274-9)
  - **Physics-informed SR for tool wear and RUL (metadata-only):** Han, Awasthi & Bollas, J. Manufacturing Systems 80:734–748 (2025); 29 citations in OpenAlex. — [doi:10.1016/j.jmsy.2025.03.023](https://doi.org/10.1016/j.jmsy.2025.03.023)
  - **Grey model plus SR for aircraft RUL (metadata-only):** Liu et al., Engineering Applications of AI 136:108986 (2024). — [doi:10.1016/j.engappai.2024.108986](https://doi.org/10.1016/j.engappai.2024.108986)
  - **C-MAPSS run-to-failure simulation benchmark (metadata-only):** Saxena, Goebel, Simon & Eklund, PHM 2008, pp. 1–9. It is commonly used to evaluate RUL methods. — [doi:10.1109/PHM.2008.4711414](https://doi.org/10.1109/PHM.2008.4711414)

### Inferences
- For degradation modelling, SINDy's sparsity assumption holds only "in an appropriate basis". The candidate library, such as Arrhenius, power-law or Paris-law terms, is therefore domain knowledge that should be curated and versioned like any other knowledge-base artefact. **[inference]**
- E-SINDy inclusion probabilities, and the Coefficient Presence metric in the 2026 preprint, are natural **admission statistics** for terms in a discovered degradation law. A term enters only if it is selected in most bootstrap resamples. **[inference / proposal]**
- Run-to-failure data are scarce and censored. The low-data, high-noise regime is therefore the normal case in PHM, so plain SINDy or SR results should be considered **conditional** until they are checked by bagging and held-out unit or asset validation. **[inference]**
- Symbolic laws can extrapolate wrongly beyond the operating envelope used for training. A wrong extrapolation is safety-relevant for RUL. Admission should therefore bind each law to its validity domain. **[inference; the reviewed abstracts did not quantify this risk]**

### Gaps
- No peer-reviewed systematic review or benchmark specifically comparing SR or SINDy for RUL across assets was found. Abstracts for Han 2025 and Liu 2024 were unavailable.
- Quantitative RUL errors in the PHM SR and SINDy papers were not extracted because full texts were not read.

## 4. LLM-assisted extraction or synthesis of rules and laws from documents and data (2023–2026): accuracy, hallucination, need for verification; neuro-symbolic verification pipelines

### Takeaway
Peer-reviewed evidence from 2023–2025 consistently shows LLMs as useful **proposers** of hypotheses, rules, equations and failure modes. They are unreliable as **verifiers or appliers** of those proposals.

The successful pipelines pair the LLM with an external evaluator: a symbolic interpreter, program execution, or a data-fitting and scoring function.

On benchmarks designed to block memorisation, LLM-based equation discovery remains weak, at 31.5% symbolic accuracy for the best system. In reliability and FMEA work, expert review discarded more than half of LLM-generated candidates and corrected systematically inflated risk scores.

### Cited Findings
- **FunSearch (Romera-Paredes et al.) [ESTABLISHED] (abstract-verified):**
  - The authors state that LLM "confabulations (or hallucinations) … hinder the use of current large models in scientific discovery".
  - FunSearch pairs a pretrained LLM with "a systematic evaluator" in an evolutionary loop.
  - It found new cap-set constructions and bin-packing heuristics.
  - The discovered programs are "more interpretable than raw solutions, enabling feedback loops between domain experts" and the search.
  - Source: Nature 625:468–475; Crossref issue year 2023, print volume 2024. — [doi:10.1038/s41586-023-06924-6](https://doi.org/10.1038/s41586-023-06924-6)
- **Qiu et al., "Phenomenal Yet Puzzling" [CONDITIONAL; peer-reviewed] (abstract-verified):**
  - LMs are "phenomenal hypothesis proposers".
  - When coupled with "a (task-specific) symbolic interpreter that is able to systematically filter the proposed set of rules", they achieve strong results.
  - They show "notable performance gaps between rule induction … and rule application".
  - Source: ICLR 2024 oral according to OpenReview. — [arXiv:2310.08559](https://arxiv.org/abs/2310.08559)
- **Wang et al., "Hypothesis Search" [CONDITIONAL; peer-reviewed] (abstract-verified):**
  - The LLM proposes natural-language hypotheses, which are implemented as Python programs and verified by running them on observed examples.
  - Hypotheses are filtered by LLM summarisation or by human selection.
  - The automated pipeline reached 30% accuracy on a random 100-problem subset of ARC.
  - Source: ICLR 2024 poster according to OpenReview. — [arXiv:2309.05660](https://arxiv.org/abs/2309.05660)
- **Zhu et al., Hypotheses-to-Theories (HtT) [PREPRINT] (abstract-verified):**
  - The LLM generates and verifies rules on training examples.
  - Rules that "appear and lead to correct answers sufficiently often" enter a rule library.
  - The method reports an absolute accuracy gain of 10–30% on relational, numerical and concept-learning tasks, and the learned rules transfer across models.
  - Status: OpenReview lists "Submitted to ICLR 2024". Later peer-reviewed publication was not verified in this session.
  - Source: [arXiv:2310.07064](https://arxiv.org/abs/2310.07064)
- **LLM-SR (Shojaee et al.) [CONDITIONAL; peer-reviewed] (abstract-verified):**
  - The LLM proposes equation-skeleton hypotheses as programs, using its scientific priors.
  - Parameters are optimised against data within evolutionary search.
  - Source: ICLR 2025 oral according to OpenReview. — [arXiv:2404.18400](https://arxiv.org/abs/2404.18400)
- **LLM-SRBench [ESTABLISHED as a benchmark; peer-reviewed] (abstract-verified):**
  - Existing benchmarks "often rely on common equations that are susceptible to memorization by LLMs, leading to inflated performance metrics".
  - The new benchmark contains 239 problems in four domains.
  - The best system achieves "only 31.5% symbolic accuracy".
  - Source: ICML 2025 oral according to OpenReview. — [arXiv:2504.10415](https://arxiv.org/abs/2504.10415)
- **LLMs for failure-mode classification in maintenance work orders (2023) [PREPRINT] (abstract-verified):**
  - Fine-tuned GPT-3.5 reached F1 = 0.80.
  - An existing text classifier trained on the same annotated data reached F1 = 0.60.
  - Out-of-the-box GPT-3.5 reached F1 = 0.46.
  - The authors stress the need for high-quality fine-tuning data.
  - Source: [arXiv:2309.08181](https://doi.org/10.48550/arxiv.2309.08181)
- **LLM-assisted FMEA in radiation oncology (Nair et al., 2026) [CONDITIONAL; single case study, non-PHM domain] (abstract-verified):**
  - Four LLMs generated 190 candidate failure modes. After review for relevance and duplication, 79 unique and interpretable modes were retained.
  - Three retained modes exceeded the risk priority number (RPN) threshold of 125.
  - Compared with experts, the LLMs gave higher severity (7.3 vs 4.1), similar occurrence (2.8 vs 3.3) and higher detectability scores (5.4 vs 2.8). Their mean RPN was 110, compared with 36 from experts.
  - Clinical users confirmed that several AI-identified risks were plausible.
  - Source: Advances in Radiation Oncology 2026. — [doi:10.1016/j.adro.2026.102060](https://doi.org/10.1016/j.adro.2026.102060)
- **LLM-integrated FMEA in engineering design [PROPOSAL / single case study] (abstract-verified):**
  - The framework uses LLMs "with human in the loop".
  - The authors report efficiency and accuracy gains over manual FMEA, but the abstract gives no figures.
  - Source: Proc. Design Society 2024. — [El Hassani et al. 2024, doi:10.1017/pds.2024.204](https://doi.org/10.1017/pds.2024.204)
- **Foundation models for PHM survey (metadata-only):** Li, Wang & Sun, "ChatGPT-like large-scale foundation models for prognostics and health management: a survey and roadmaps", Reliability Engineering & System Safety 243:109850 (2024). — [doi:10.1016/j.ress.2023.109850](https://doi.org/10.1016/j.ress.2023.109850)
- **LLM plus knowledge graph for industrial fault diagnosis [PROPOSAL; conference] (abstract-verified):**
  - The authors note that LLMs have been "criticized by their 'hallucination'" in specialised fields.
  - They respond by building knowledge graphs and fine-tuning datasets.
  - Source: IEEE RICAI 2023. — [doi:10.1109/ricai60863.2023.10489566](https://doi.org/10.1109/ricai60863.2023.10489566)
  - A related journal paper is Liao et al., Advanced Engineering Informatics 2025 (metadata-only). — [doi:10.1016/j.aei.2025.103134](https://doi.org/10.1016/j.aei.2025.103134)

### Inferences
- These sources agree on one architectural pattern, **LLM as generator and external symbolic or data evaluator as gatekeeper**. The pattern appears in FunSearch, Qiu et al., Hypothesis Search, HtT and LLM-SR. This is direct support for a two-stage admission gate in which no LLM-proposed rule enters the knowledge base on the LLM's own judgement. **[inference from multiple peer-reviewed sources]**
- The LLM-SRBench memorisation finding means that an LLM recovering a textbook degradation law, such as Arrhenius or the Paris law, is **not evidence of discovery capability**. Validation must use held-out or transformed problems from the target asset. **[inference]**
- In Nair et al., risk-score inflation of about 3× in mean RPN, and the retention of only about 42% (79 of 190) of candidates, are the most concrete quantitative evidence found that expert rescoring is not optional. The study is medical rather than industrial, so transfer to PHM is **conditional**. **[inference]**
- Failure-mode classification reached F1 = 0.46 out of the box versus 0.80 after fine-tuning. Unadapted LLMs should therefore not be the sole extractor of maintenance knowledge from work orders. **[inference; preprint evidence]**

### Gaps
- No peer-reviewed study was found that quantifies hallucination rates in LLM extraction of **diagnostic or prognostic rules** from maintenance manuals or standards, with rule-level precision and recall against expert ground truth.
- The full content of the Li et al. 2024 RESS survey was not accessible here; only metadata was checked.
- The final peer-review status of HtT is unresolved.

## 5. Validation before admission: holdout and temporal validation, MDL, significance and multiple-testing control, stability, expert review — and the size of false-discovery risk

### Takeaway
When many candidate rules are generated and tested, chance alone produces many apparently valid rules. Webb (2007) shows empirically that standard pattern discovery finds "numerous spurious patterns" in random data.

Established safeguards include:

- experimentwise error control, either FWER or Bonferroni-style correction;
- false-discovery-rate control (Benjamini–Hochberg);
- subsampling-based stability selection with finite-sample false-discovery control;
- parsimony criteria such as MDL and sparsity;
- holdout evaluation.

For PHM, these should be complemented by temporal or unit-level holdout and documented expert review. Expert review is supported by the FMEA case evidence in Question 4.

### Cited Findings
- **Webb (2007), spurious patterns in pattern discovery [ESTABLISHED] (abstract verified on the Springer page):**
  - Pattern discovery explores large spaces and therefore suffers "an extreme risk of type-1 error".
  - The proposed techniques "allow the user to enforce a strict upper limit on the risk of experimentwise error".
  - Empirically, standard techniques "can discover numerous spurious patterns when applied to random data". On real data they "result in large numbers of patterns that are rejected when subjected to sound statistical evaluation".
  - "A number of pragmatic choices about how such tests are performed can greatly affect their power."
  - Source: Machine Learning 68:1–33. — [Webb 2007, doi:10.1007/s10994-007-5006-x](https://doi.org/10.1007/s10994-007-5006-x)
- **Hämäläinen & Webb (2019 issue), statistically sound pattern discovery [ESTABLISHED] (abstract-verified):**
  - Appropriate statistical tests allow "precise control over the risk of false discoveries".
  - Tests can also "filter out patterns that are unlikely to be useful, removing uninformative variations".
  - The tutorial covers techniques "for controlling the likelihood of spurious discoveries when multitudes of patterns are evaluated", for dependency rules and dependency sets.
  - Source: Data Mining and Knowledge Discovery 33:325–377. — [doi:10.1007/s10618-018-0590-x](https://doi.org/10.1007/s10618-018-0590-x)
- **Benjamini & Hochberg (1995), false-discovery-rate control [ESTABLISHED] (abstract-verified):**
  - The false discovery rate (FDR) is the expected proportion of falsely rejected hypotheses.
  - FDR equals the familywise error rate (FWER) when all null hypotheses are true and is smaller otherwise, which offers more power.
  - A sequential Bonferroni-type procedure is proved to control FDR for independent test statistics.
  - Source: JRSS B 57(1). — [doi:10.1111/j.2517-6161.1995.tb02031.x](https://doi.org/10.1111/j.2517-6161.1995.tb02031.x)
- **Meinshausen & Bühlmann (2010), stability selection [ESTABLISHED] (abstract-verified):**
  - Stability selection combines subsampling with selection algorithms.
  - It provides "finite sample control for some error rates of false discoveries" and a transparent principle for choosing regularisation.
  - Variable selection and structure estimation "improve markedly" across selection methods.
  - Source: JRSS B 72(4):417–473. — [doi:10.1111/j.1467-9868.2010.00740.x](https://doi.org/10.1111/j.1467-9868.2010.00740.x)
- **Minimum description length (metadata-only):** Rissanen, "Modeling by shortest data description", Automatica 14:465–471 (1978). — [doi:10.1016/0005-1098(78)90005-5](https://doi.org/10.1016/0005-1098(78)90005-5)
- **Parsimony as an overfitting safeguard in law discovery [ESTABLISHED] (abstract-verified):**
  - SINDy seeks "the fewest terms" to "balance accuracy with model complexity to avoid overfitting". — [Brunton et al. 2016](https://doi.org/10.1073/pnas.1517384113)
  - Schmidt & Lipson use a nontriviality criterion to separate important correlations from trivial ones. — [Schmidt & Lipson 2009](https://doi.org/10.1126/science.1165893)
- **Resampling stability for discovered equation terms [ESTABLISHED] (abstract-verified):** E-SINDy uses bootstrap inclusion probabilities for candidate terms. — [Fasel et al. 2022](https://doi.org/10.1098/rspa.2021.0904)
- **Evaluator-gated admission in LLM pipelines [CONDITIONAL; peer-reviewed for FunSearch, Qiu et al. and Hypothesis Search; preprint for HtT] (abstract-verified):**
  - FunSearch uses a systematic evaluator. — [FunSearch](https://doi.org/10.1038/s41586-023-06924-6)
  - Qiu et al. use a symbolic interpreter as a filter. — [Qiu et al.](https://arxiv.org/abs/2310.08559)
  - Hypothesis Search executes programs on the observed examples. — [Wang et al.](https://arxiv.org/abs/2309.05660)
  - HtT admits rules by frequency and correctness on training examples. — [Zhu et al.](https://arxiv.org/abs/2310.07064)
- **Logical-consistency gate in ILP [ESTABLISHED] (abstract-verified):** Popper rejects any hypothesis that fails to entail all positive examples or entails a negative example. — [Cropper & Morel 2021](https://doi.org/10.1007/s10994-020-05934-z)
- **Expert review as the final gate [CONDITIONAL] (abstract-verified):**
  - Expert review retained 79 of 190 LLM-generated failure modes.
  - Experts rescored LLM risk estimates downward, with a mean RPN of 36 versus 110 from the LLMs.
  - Source: [Nair et al. 2026](https://doi.org/10.1016/j.adro.2026.102060)

### Inferences
- **Scale of false-discovery risk (arithmetic, not an empirical source):**
  - Under the global null, testing m independent candidate rules at per-test level α gives an expected m·α false positives.
  - For m = 10,000 and α = 0.05, about 500 spurious rules would pass a naive significance filter.
  - Bonferroni control of the familywise error rate at 0.05 would require a per-test threshold of 0.05 / 10,000 = 5 × 10⁻⁶.
  - This illustrates Webb's empirical finding that random data yield "numerous spurious patterns". **[inference]**
- **Proposed layered admission gate for machine-proposed PHM knowledge [PROPOSAL, synthesised from the sources above]:**
  1. syntactic and logical consistency, as in the Popper-style test;
  2. parsimony and MDL or sparsity, as in SINDy and Schmidt & Lipson;
  3. statistical significance with experimentwise or FDR correction over the *whole explored candidate space*, following Webb and Benjamini–Hochberg;
  4. resampling stability through bootstrap inclusion probability, as in E-SINDy and stability selection;
  5. temporal and unit-level holdout, meaning evaluation on later time windows and unseen assets, to respect degradation's time ordering and fleet heterogeneity;
  6. domain-of-validity tagging;
  7. documented expert review with rescoring, informed by the FMEA evidence.

  Steps 5 and 6 are standard practice in reliability engineering but are **not** supported by a verified source in this session.
- In PHM, rows from the same unit's life are strongly autocorrelated. Random row-wise cross-validation would leak information and understate false-discovery risk. Holdout should therefore be by unit and by time. **[inference; no verified PHM-specific source located]**

### Gaps
- Two details of Webb (2007) are recalled from the literature but could not be verified from the abstract: the recommendation of a "direct adjustment" based on the size of the search space and a "holdout evaluation" approach. Check them against the full text before citing them as such.
- No verified PHM-specific source was found on temporal or unit-level validation leakage, such as on C-MAPSS. This should be sourced separately.
- No empirical estimate of false-discovery rates for evolutionary or LLM rule generators in PHM was found. Such estimates would come from the number of admitted rules later invalidated in operation.
- No verified source was found on formal expert-review protocols for admitting machine-generated rules into industrial knowledge bases, for example a Delphi method or two-person review with rationale capture.
