# Symbolic Knowledge Revision and Maintenance of Expert-System Knowledge Bases (rules, exceptions, ontologies)

Scope: formal theories of knowledge change, justification tracking, historical KB-maintenance experience, post-change verification, ontology/KG versioning and provenance, and defeasible formalisms for adding exceptions. The purpose is to inform evolutionary knowledge development in predictive expert systems (prognostics, condition monitoring, decision support). Research was done on 2026-11-03.

Status tags used below:
- **[ESTABLISHED]**: a proven theorem or definition, or a normative standard, in the cited source.
- **[CONDITIONAL]**: holds only under stated assumptions, such as logical closure, classical monotone consequence, Horn justifications, or a propositional language.
- **[EMPIRICAL REPORT]**: a measured or reported result from one deployment or experiment. It is not a general law.
- **[PROPOSAL/OPINION]**: an author's design argument, recommendation or prospective claim.
- **[INFERENCE]**: the note-taker's synthesis. It appears only in the Inferences subsections.

Verification note: every DOI below was resolved through the Crossref API or the OpenAlex API during this session. ScienceDirect (Elsevier), ACM DL and Springer blocked access, and Semantic Scholar rate-limited requests. For the *Artificial Intelligence* journal papers by Katsuno–Mendelzon (1991), Eiter–Gottlob (1992), Darwiche–Pearl (1997), de Kleer (1986), Politakis–Weiss (1984) and Ginsberg–Weiss–Politakis (1988), only the metadata (title, authors, volume, pages, date) was verified. Claims about their content come from accessible companion or secondary sources, which are named where used: the SEP entry, the AAAI or IJCAI versions, and the MIT memo.

## 1. What belief revision and update theories guarantee, and their limits for rule bases with priorities and exceptions

### Takeaway
AGM theory gives a representation theorem rather than an algorithm. Six basic postulates characterize exactly the family of "partial meet" revision operations. Every member of that family guarantees success, consistency, minimal change relative to expansion, and syntax-independence. The theory does **not** choose one operator, because the choice depends on an extra-logical selection mechanism or entrenchment ordering. It assumes logically closed belief sets under a classical *monotone* consequence relation, and it says nothing about iterated change unless extra postulates, such as Darwiche–Pearl, are added. Katsuno–Mendelzon showed that "the world changed" (update) and "we learned we were wrong" (revision) need different postulates.

### Cited Findings
- **[ESTABLISHED]** AGM (Alchourrón, Gärdenfors, Makinson 1985, *J. Symbolic Logic* 50(2):510–530) defines *partial meet contraction*: the intersection of some nonempty family of maximal subsets of the theory that fail to imply the proposition being removed. The paper proves that these functions satisfy the Gärdenfors postulates and gives a representation theorem for those postulates. It also gives a further representation theorem for the "relational" and "transitively relational" subclasses together with the supplementary postulates. — [Crossref record and abstract, DOI 10.2307/2274239](https://api.crossref.org/works/10.2307/2274239)
- **[ESTABLISHED]** AGM has three change operations. Expansion $`K+p=\mathrm{Cn}(K\cup\{p\})`$ removes nothing. Contraction $K\div p$ removes $p$. Revision $`K*p`$ adds $p$ and removes whatever must go to keep consistency. Revision is built from contraction by the *Levi identity* (contract by $\neg p$, then expand by $p$). — [SEP, "Logic of Belief Revision", §§1.2, 3](https://plato.stanford.edu/entries/logic-belief-revision/)
- **[ESTABLISHED]** An operation $`*`$ is a partial meet revision **if and only if** it satisfies the six basic Gärdenfors postulates:
  - **Closure:** $`K*p=\mathrm{Cn}(K*p)`$.
  - **Success:** $`p\in K*p`$.
  - **Inclusion:** $`K*p\subseteq K+p`$.
  - **Vacuity:** if $\neg p\notin K$, then $`K*p=K+p`$.
  - **Consistency:** $`K*p`$ is consistent if $p$ is consistent.
  - **Extensionality:** logically equivalent inputs give the same result.

  Two supplementary postulates, Superexpansion and Subexpansion, are standard additions. — [SEP §3](https://plato.stanford.edu/entries/logic-belief-revision/)
- **[ESTABLISHED]** Grove's (1988) sphere-based revision, a nested system of plausibility "spheres" around $[K]$, corresponds exactly to transitively relational partial meet revision. It therefore also corresponds exactly to entrenchment-based revision. — [SEP §4](https://plato.stanford.edu/entries/logic-belief-revision/)
- **[CONDITIONAL]** The AGM assumptions are restrictive:
  - Belief sets are closed under a Tarskian consequence operator $`\mathrm{Cn}`$ that satisfies Inclusion, **Monotony** and Iteration, and $`\mathrm{Cn}`$ is supraclassical.
  - The SEP calls this closure "clearly an unrealistic idealization", because it means that the agent is taken to be "logically omniscient".
  - The models are "input-assimilating", and "no explicit representation of time is included".

  — [SEP §§1.2–1.3](https://plato.stanford.edu/entries/logic-belief-revision/)
- **[OPINION, recorded in SEP]** Success is "the most controversial" postulate. It forces any input to be accepted, and an inconsistent input yields an inconsistent belief set. *Non-prioritized* belief change, such as semi-revision, accepts a contradicting input only "if it has more epistemic value than the original beliefs that contradict it". — [SEP §§3, 6, 6.1](https://plato.stanford.edu/entries/logic-belief-revision/)
- **[ESTABLISHED] Syntax matters for finite bases.** Belief *bases*, which are not logically closed, can be "statically equivalent" (the same closure) but not "dynamically equivalent". Bases $`\{p,q\}`$ and $`\{p, p\leftrightarrow q\}`$ diverge after revision by $\neg p$: one keeps $q$ and the other derives $\neg q$. Further, Recovery does not hold for partial meet contraction of bases. Base revision splits into *internal* revision $(A\div\neg p)+'p$ and *external* revision $(A+'p)\div\neg p$, and "neither of them can be subsumed under the other". — [SEP §§5.1–5.4](https://plato.stanford.edu/entries/logic-belief-revision/)
- **[ESTABLISHED] Revision differs from update.** Katsuno and Mendelzon "argue that no such set of postulates will be adequate for every application". They distinguish *update*, which brings the KB up to date "when the world described by it changes", from revision, which incorporates new information about an unchanged world. — [Katsuno & Mendelzon 1992, in Gärdenfors (ed.) *Belief Revision*, CUP, pp. 183–203, DOI 10.1017/CBO9780511526664.007 (OpenAlex abstract)](https://api.openalex.org/works/doi:10.1017/cbo9780511526664.007)

  The SEP illustrates the difference with an example. Suppose the agent knows that exactly one of a book or a magazine is on the table. If told "book", the agent should conclude "no magazine": this is revision. If told "a book has since been put on the table", the agent should not conclude that: this is update. One approach adds time indices to sentences. — [SEP §7.1](https://plato.stanford.edu/entries/logic-belief-revision/)

  The model-theoretic AGM/KM paper for propositional KBs is verified as Katsuno & Mendelzon, "Propositional knowledge base revision and minimal change", *Artificial Intelligence* 52(3):263–294, 1991 (metadata only). — [Crossref](https://api.crossref.org/works/10.1016/0004-3702(91)90069-V)
- **[ESTABLISHED] Iterated revision is under-determined by AGM.** A global AGM operation that uses one selection function for all belief sets "does not impose any new properties on contraction or revision". Most work on iterated change therefore adds postulates, centrally the Darwiche–Pearl postulates DP1–DP4:
  - **DP1:** if $q\vdash p$, then $`(K*p)*q=K*q`$.
  - **DP2:** if $q\vdash\neg p$, then $`(K*p)*q=K*q`$.
  - **DP3:** if $`K*q\vdash p`$, then $`(K*p)*q\vdash p`$.
  - **DP4:** if $`K*q\nvdash\neg p`$, then $`(K*p)*q\nvdash\neg p`$.

  These postulates keep the ordering *within* the $p$-worlds and *within* the $\neg p$-worlds and shift only their relative position. The SEP adds: "A large number of proposals for such shifts have been put forward, but opinions differ widely on the adequacy of these proposals." — [SEP §8](https://plato.stanford.edu/entries/logic-belief-revision/). The primary source is Darwiche & Pearl, "On the logic of iterated belief revision", *Artificial Intelligence* 89(1–2):1–29, 1997 (metadata verified). — [Crossref](https://api.crossref.org/works/10.1016/S0004-3702(96)00038-0)
- **[ESTABLISHED] Change can be non-deterministic.** Indeterministic belief change allows more than one admissible outcome, for example a set of operators $`\{*,*',*''\}`$. — [SEP §7.4](https://plato.stanford.edu/entries/logic-belief-revision/)
- **[CONDITIONAL] Resources.** "Operations on infinite sets can require infinite resources." The use of finite belief bases can be seen as a way to avoid this. — [SEP §7.7](https://plato.stanford.edu/entries/logic-belief-revision/)
- **[ESTABLISHED, metadata only]** The computational complexity of propositional KB revision, update and counterfactuals was analysed in Eiter & Gottlob, *Artificial Intelligence* 57(2–3):227–270, 1992. — [Crossref](https://api.crossref.org/works/10.1016/0004-3702(92)90018-S). The text and abstract were not accessible: ScienceDirect and ACM DL (for the PODS 1992 version, DOI 10.1145/137097.137886) returned 403 errors. See Gaps.

### Inferences
- **[INFERENCE]** AGM is a *specification* of admissible change, not a procedure. For an engineered rule base, the "selection function" or entrenchment ordering corresponds to explicit, auditable rule priorities: which rules or facts are given up first. The representation theorem says that any admissible operator *is* such a choice. The postulates alone cannot make revision unique.
- **[INFERENCE]** A rule base is a finite *base*, not a closed belief set. The SEP base results therefore apply directly: two rule bases with identical current conclusions can react differently to the same correction. Change control has to version the rule text (the base), not only the derived conclusions.
- **[INFERENCE]** In prognostics, the revision/update distinction maps onto two different maintenance events:
  - **Update:** the asset changed (wear progressed, a component was replaced). New sensor facts update the world model and should not falsify the rules.
  - **Revision:** the degradation rule or threshold was wrong. A correction revises the knowledge.

  Mixing the two leads to wrong retractions, for example deleting a valid rule because a component was replaced.
- **[INFERENCE]** Classical AGM assumes monotone $`\mathrm{Cn}`$ ([SEP §1.3](https://plato.stanford.edu/entries/logic-belief-revision/)), so it does not directly govern rule bases with defeasible rules, priorities and exceptions. In those bases, adding a rule can *remove* conclusions (see question 6). Guarantees for such bases must come from the nonmonotonic formalism plus regression testing, not from AGM postulates.

### Gaps
- The specific complexity classes from Eiter & Gottlob (1992) were not verified because the full text and abstract were inaccessible. The commonly cited result that inference from a revised or updated KB is complete for the second level of the polynomial hierarchy for most operators is **not** asserted here and needs checking against the paper.
- Primary texts of Katsuno–Mendelzon 1991 (the KM postulates R1–R6 and U1–U8) and Darwiche–Pearl 1997 were not read directly. The postulate statements above come from the SEP.
- Belief revision specifically for logic programs or defeasible logic (revision of nonmonotonic rule bases) was not researched within this time budget.

## 2. What truth maintenance systems (Doyle 1979 TMS; de Kleer 1986 ATMS) provide

### Takeaway
A TMS records *why* each belief is held, as justifications, so that retracting a premise automatically retracts the conclusions it supported, with dependency-directed backtracking and explanations. Doyle's JTMS keeps one current context. The ATMS keeps all consistent contexts at once: every conclusion carries a label of minimal assumption sets, together with a database of minimal inconsistent sets ("nogoods"). The cost is possible exponential growth of labels and nogoods. The basic ATMS is restricted to Horn justifications.

### Cited Findings
- **[ESTABLISHED, Doyle's own description]** Doyle describes the system as follows:
  - Reasoning programs "must be able to make assumptions and subsequently revise their beliefs when discoveries contradict these assumptions".
  - The TMS performs this "by recording and maintaining the reasons for program beliefs".

  The paper covers:
  1. the TMS representations and structure;
  2. the mechanisms used to revise the current set of beliefs;
  3. how "dependency-directed backtracking changes the current set of assumptions";
  4. techniques for summarizing explanations of beliefs;
  5. "dialectically arguing" modules;
  6. how to revise models of other agents' belief systems;
  7. methods for embedding control structures in patterns of assumptions.

  — [MIT AI Memo AIM-521 abstract, hdl 1721.1/5733](https://dspace.mit.edu/handle/1721.1/5733). The journal version is Doyle, "A truth maintenance system", *Artificial Intelligence* 12(3):231–272, 1979. — [Crossref](https://api.crossref.org/works/10.1016/0004-3702(79)90008-0)
- **[ESTABLISHED, historical]** The SEP credits Doyle's TMS as important in the computer-science line of research that led to belief revision theory. — [SEP §1.1](https://plato.stanford.edu/entries/logic-belief-revision/)
- **[ESTABLISHED] ATMS definitions, following de Kleer 1988.**
  - An environment is a set of assumptions.
  - A node holds in environment $E$ if it can be propositionally derived from $E\cup C$, where $C$ is the set of justification clauses.
  - An environment is a *nogood* if $\bot$ holds in it.
  - The ATMS is incremental: it receives a stream of nodes, assumptions and justifications.
  - For each node it maintains a *label*, a set of environments with four properties:
    1. Soundness: the node holds in each environment.
    2. Consistency: no environment is nogood.
    3. Completeness: every environment in which the node holds is a superset of one of them.
    4. Minimality: no environment in the label is a proper subset of another.
  - "Whether n holds in environment E" is then answered by checking whether $E$ is a superset of some label environment.
  - The ATMS "also maintains a data-base of all minimal nogoods".

  — [de Kleer, "A General Labeling Algorithm for Assumption-Based Truth Maintenance", AAAI-88, p. 188](https://cdn.aaai.org/AAAI/1988/AAAI88-034.pdf). The original is de Kleer, "An assumption-based TMS", *Artificial Intelligence* 28(2):127–162, 1986 (metadata verified). — [Crossref](https://api.crossref.org/works/10.1016/0004-3702(86)90080-9)
- **[CONDITIONAL]** "The basic ATMS is restricted to accepting only horn clause justifications." Earlier generalizations needed "complex and difficult to integrate hyperresolution rules". — [de Kleer AAAI-88](https://cdn.aaai.org/AAAI/1988/AAAI88-034.pdf)
- **[EMPIRICAL REPORT]** ATMSs were used for qualitative reasoning, diagnosis, visual scene interpretation and evidential reasoning. They were also "incorporated into a variety of commercially available expert system shells". — [de Kleer AAAI-88](https://cdn.aaai.org/AAAI/1988/AAAI88-034.pdf)
- **[ESTABLISHED limitation]** There is a risk of exponential blowup:
  - "The addition of any justification can result in the discovery of any number of nogoods."
  - Each new nogood of size $k$ "can cause k new label propagations".
  - A problem "may have an exponential number of nogoods".

  — [de Kleer AAAI-88](https://cdn.aaai.org/AAAI/1988/AAAI88-034.pdf)
- **[EMPIRICAL REPORT / PROPOSAL]** Forbus & de Kleer report that "existing ATMS techniques generally perform badly" on problems with enormous search spaces. They propose an "implied-by" strategy that uses focus environments, which "prevents combinatorial explosions in the ATMS" by returning control to the inference engine. — [Forbus & de Kleer, "Focusing the ATMS", AAAI-88](https://cdn.aaai.org/AAAI/1988/AAAI88-035.pdf)

### Inferences
- **[INFERENCE]** For a predictive expert system, TMS-style justification links are the mechanism that makes knowledge change *safe at runtime*. When a sensor reading is invalidated, or a rule is retracted or revised, every derived prognosis that depended on it is identified and withdrawn. No fresh inference over everything is needed. The same links give explanations for audit.
- **[INFERENCE]** ATMS labels map naturally onto "under which assumptions does this prognosis hold?". Examples of such assumptions: sensor *k* is healthy, or the operating regime is nominal. Minimal nogoods then act as recorded conflict sets, consistent with the ATMS's documented use in diagnosis. Label growth must be bounded, for example by focusing, when many sensor-health assumptions exist.
- **[INFERENCE]** A TMS maintains *consistency of derived beliefs given the justifications*. It does not decide *which* premise to give up. That choice is the AGM-style selection or priority problem from question 1. The two layers are complementary.

### Gaps
- The text of de Kleer 1986 was not read directly. Its specific complexity statements and its comparison with the JTMS were not verified beyond the AAAI-88 companion paper.
- No quantitative field data on TMS or ATMS performance in deployed industrial expert systems was found in this session.

## 3. Documented experience with KB refinement and maintenance (TEIRESIAS, SEEK/SEEK2, RDR, XCON) and the knowledge acquisition bottleneck

### Takeaway
The successful systems tied each knowledge change to a concrete failing case. TEIRESIAS acquired rules "in the context of" an error. SEEK/SEEK2 used statistics over a case library to propose refinements. RDR added an exception that was justified by the features of the misclassified case and checked against stored cornerstone cases. The cautionary case is XCON. It grew from about 700 to about 6,200 rules, with about 50% of rules changing each year, and became hard to change because control and intent were implicit in the rules. DEC's response was a re-architecture (RIME) that made problem-solving structure explicit.

### Cited Findings
- **[EMPIRICAL REPORT] The knowledge acquisition bottleneck, in practitioners' words.**
  - Davis notes that a knowledge base "is often assembled by hand, an ongoing task that may involve several person-years of effort". Because "the domain expert often knows nothing about programming", the interaction "usually requires the mediation of a human programmer". — [Davis, "Interactive Transfer of Expertise", ch. 9 in Buchanan & Shortliffe (eds.), *Rule-Based Expert Systems* (1984), reprinted from *AI* 12:121–157 (1979)](https://www.shortliffe.net/Buchanan-Shortliffe-1984/Chapter-09.pdf)
  - Ginsberg et al. write that "the most difficult aspect of expert system design is the construction of the knowledge base". — [Ginsberg, Weiss & Politakis, IJCAI-85](https://ijcai.org/proceedings/1985/0069.pdf)
- **[ESTABLISHED design] TEIRESIAS** was designed as follows:
  - It supports two information flows: *explanation*, which shows the expert how the program reached its result, and *knowledge transfer*, by which the expert adds to or modifies the knowledge base.
  - Its goals were tools "for knowledge base construction and maintenance" and an "intelligent assistant".
  - It uses *meta-level knowledge* ("a program can 'know what it knows'"), including *rule models*.
  - Acquiring a rule in the context of a specific error "provides a number of expectations about the content of the rule". These expectations are used to interpret the expert's input and to check "that it does indeed solve the problem at hand".

  — [Davis, ch. 9, MYCIN book](https://www.shortliffe.net/Buchanan-Shortliffe-1984/Chapter-09.pdf). The journal article is Davis, *Artificial Intelligence* 12(2):121–157, 1979. — [Crossref](https://api.crossref.org/works/10.1016/0004-3702(79)90014-6)
- **[ESTABLISHED, metadata]** SEEK: Politakis & Weiss, "Using empirical analysis to refine expert system knowledge bases", *Artificial Intelligence* 22(1):23–48, 1984. — [Crossref](https://api.crossref.org/works/10.1016/0004-3702(84)90024-9). SEEK2 journal version: Ginsberg, Weiss & Politakis, "Automatic knowledge base refinement for classification systems", *AI* 35(2):197–226, 1988. — [Crossref](https://api.crossref.org/works/10.1016/0004-3702(88)90012-4)
- **[ESTABLISHED definition]** Knowledge-base refinement is "the addition, deletion, and alteration of rule-components in an existing knowledge base, in an attempt to improve an expert system's performance". It is distinct from the first phase, the acquisition of whole rules. Performance is evaluated by running the KB on each stored case and "comparing the knowledge base's conclusion with the stored expert's conclusion". The evaluation gives an overall score and a breakdown by diagnostic category. — [Ginsberg, Weiss & Politakis, "SEEK2: A generalized approach to automatic knowledge base refinement", IJCAI-85](https://ijcai.org/proceedings/1985/0069.pdf)
- **[EMPIRICAL REPORT] SEEK2 quantitative result.** The setup used 10 statistical concepts and 9 heuristics for generating refinements. On a rheumatology KB of about 140 rules with 5 final diagnostic categories and 121 cases, the automatic pilot:
  - raised performance from **73% (88/121) to 98% (119/121)**;
  - used about **18 min** of VAX-785 CPU time;
  - tried **106** experiments and accepted **8**.

  Caveats from the authors:
  - The gain is "assured" only "with respect to the given data base of cases", with a "reasonable expectation" for new cases.
  - "The expert is still the final judge" and "may agree with only a subset" of the refinements.

  — [IJCAI-85, p. 373](https://ijcai.org/proceedings/1985/0069.pdf)
- **[EMPIRICAL REPORT] XCON at DEC** (Soloway, Bachant & Jensen, AAAI-87, pp. 824–829):
  - Over 7 years XCON grew to **6,200 rules**, "of which approximately 50% change every year".
  - It draws on a database of about **20,000 parts** and has been in production use since January 1980.
  - It "started as a relatively small, rule-based system (about 700 rules)".
  - The authors write that the problems of updating "do not grow linearly".

  — [AAAI-87 paper](https://cdn.aaai.org/AAAI/1987/AAAI87-148.pdf)
- **[EMPIRICAL REPORT] XCON failure modes** (same source):
  - *Dynamic:* as rules multiply and different programmers edit them, "rules that are not appropriate become triggered", and control is implicit in OPS5 conflict resolution or hidden in "special tricks".
  - *Static:* a single rule performs several functions, and the typical extension practice is to copy the rules for a similar device and edit them.
  - These produce a "degradation in integrity", a "rat's nest of special rules, tightly coupled rules".
  - Developers "are afraid to modify the rules". A developer "needs to understand at least a major portion of the rules before he can effectively make some change".
  - The needed knowledge lives in "institutional memory".
  - Performance stayed satisfactory. The degradation concerned changeability.

  — [AAAI-87](https://cdn.aaai.org/AAAI/1987/AAAI87-148.pdf)
- **[PROPOSAL]** RIME, the language of XCON-in-RIME, makes the following explicit:
  - problem spaces;
  - method steps (propose, eliminate, apply, evaluate);
  - rule types;
  - a classification scheme intended to help ensure "one function, one rule".

  The paper argues that XCON-in-RIME "should be more maintainable". This is a prospective assessment, not a measured outcome in that paper. — [AAAI-87](https://cdn.aaai.org/AAAI/1987/AAAI87-148.pdf). A follow-up retrospective by the DEC team covers "a decade's worth of lessons learned": Barker, O'Connor, Bachant & Soloway, "Expert systems for configuration at Digital: XCON and beyond", *CACM* 32(3):298–318, 1989 (metadata and abstract only). — [Crossref](https://api.crossref.org/works/10.1145/62065.62067)
- **[ESTABLISHED design + historical claim] Ripple-Down Rules.**
  - RDR "were developed in answer to the problem of maintaining medium to large rule-based knowledge systems".
  - Traditional approaches "gave little thought to maintenance".
  - RDR "made knowledge acquisition and maintenance one and the same thing", with experts interacting directly with naturally occurring cases.
  - Cases motivate the capture of new knowledge and frame its context. They also ensure "that previously correctly classified cases remained so", because the new rule must distinguish the new case and be justified by its features. This is the cornerstone-case principle.
  - RDR later extended to multiple conclusions and to help-desk, email classification, RoboCup, configuration, simulation, planning and NLP tasks.

  — [Richards, "Two decades of Ripple Down Rules research", *Knowledge Engineering Review* 24(2):159–184, 2009](https://api.crossref.org/works/10.1017/S0269888909000241). Foundational paper (metadata verified via OpenAlex): Compton & Jansen, "A philosophical basis for knowledge acquisition", *Knowledge Acquisition* 2(3):241–258, 1990, DOI 10.1016/S1042-8143(05)80017-2. — [OpenAlex](https://api.openalex.org/works/doi:10.1016/S1042-8143(05)80017-2)
- **[EMPIRICAL REPORT] PEIRS (RDR in pathology).**
  - The authors created "over 950 rules" for thyroid function tests, arterial blood gases and other test groups "in 9 mths of operation".
  - "A staff pathologist performed all maintenance tasks as part of his routine duties without any need for computer programming skills."
  - No clerical staff were needed.
  - Stated limitation: "duplication of rule addition for reports requiring multiple comments".

  — [Edwards, Compton, Malor, Srinivasan & Lazarus, *Pathology* 25(1):27–34, 1993 (PubMed abstract, PMID 8316495)](https://pubmed.ncbi.nlm.nih.gov/8316495)

### Inferences
- **[INFERENCE]** Successful maintenance regimes share one pattern: each knowledge change is *case-anchored*. It is triggered by a concrete misclassified case, justified by that case's features, and checked against the stored case library. This is the pattern in TEIRESIAS, SEEK2 and RDR. For prognostics, the natural "case" is an asset history with a confirmed outcome, such as a failure or a maintenance finding.
- **[INFERENCE]** XCON shows that *runtime correctness and maintainability can diverge*. The system kept performing, but change cost grew superlinearly as implicit control and multi-purpose rules accumulated. Structural rules make knowledge structure explicit: one function per rule, an explicit method or phase, and explicit priorities instead of tricks to force firing order. These rules are maintainability controls that behaviour tests alone will not enforce.
- **[INFERENCE]** SEEK2's 73%→98% was measured on the same 121 cases used to drive refinement. It is evidence of fit, not of generalization. A modern process would require held-out cases or cross-validation, plus expert acceptance, before a refinement is admitted. SEEK2's 8 accepted out of 106 trials already shows that most candidate edits are rejected.

### Gaps
- Quantitative long-term RDR maintenance data (for example time per rule or KB size after years of use) was not verified. Compton et al., "Experience with Ripple-Down Rules", *Knowledge-Based Systems* 19(5):356–362, 2006 (DOI 10.1016/j.knosys.2005.11.022, metadata verified via OpenAlex) is likely the right source, but its text was inaccessible (ScienceDirect and Springer blocked).
- The text of Compton & Jansen 1990 was not read directly. Its specific claims about "knowledge in context" are therefore taken only through Richards 2009.
- The SEEK (1984) and SEEK2 (1988) journal papers were not read. Only the IJCAI-85 SEEK2 paper was.
- Barker et al. 1989 (CACM) was not read beyond its abstract. No verified later figures on XCON's size, change rate or costs are included.

## 4. How rule bases are verified and validated after a change, and what these checks cannot detect

### Takeaway
Static anomaly detection checks the rule base for four classes of anomaly: redundancy (including subsumption), ambivalence or conflict, circularity, and deficiency (missing coverage). It is a well-founded verification method, but it finds *symptoms of probable errors*, not proofs of error or of correctness. Dynamic validation re-runs the KB against a case library with known expert conclusions, as in SEEK2's scoring and RDR's cornerstone cases. It catches regressions only on the cases in the library. Neither check detects wrong-but-consistent knowledge outside the stored cases, or loss of structural maintainability.

### Cited Findings
- **[ESTABLISHED taxonomy]** Preece, Shinghal & Batarekh describe four types of anomaly: **redundancy, ambivalence, circularity and deficiency**. They explain them for rule bases based on first-order logic and review five programs that detect various subsets of these anomalies. They state that anomalies are "highly indicative of errors". — [*Knowledge Engineering Review* 7(2):115–141, 1992, DOI 10.1017/S026988890000624X (Crossref abstract)](https://api.crossref.org/works/10.1017/S026988890000624X)
- **[ESTABLISHED + EMPIRICAL REPORT]** "Anomalies such as redundant, contradictory, and deficient knowledge in a knowledge base are symptoms of probable errors." The paper gives "a theoretical foundation of anomaly detection methods" and discusses their "computational limitations". It reports empirical results from applying one anomaly-detection tool to **five real-world knowledge-based systems**. The techniques "apply specifically to verifying rule-based knowledge bases without numerical certainty measures." — [Preece & Shinghal, "Foundation and application of knowledge base verification", *Int. J. Intelligent Systems* 9(8):683–701, 1994, DOI 10.1002/int.4550090804 (OpenAlex abstract)](https://api.openalex.org/works/doi:10.1002/int.4550090804)
- **[ESTABLISHED, metadata]** An integrated verification–validation–refinement framework was published in the same issue: Zlatareva, "A framework for verification, validation, and refinement of knowledge bases: The VVR system", *IJIS* 9(8):703–737, 1994. — [Crossref](https://api.crossref.org/works/10.1002/int.4550090805)
- **[ESTABLISHED practice] Case-library regression.** SEEK2 evaluates a KB by running it on every stored case and comparing its conclusions with the expert's. — [IJCAI-85](https://ijcai.org/proceedings/1985/0069.pdf). RDR requires every new rule to keep "previously correctly classified cases" correct. — [Richards 2009](https://api.crossref.org/works/10.1017/S0269888909000241)
- **[EMPIRICAL REPORT, limit]** Improvement is "assured" only on the given case base, and the expert remains the final judge. — [SEEK2, IJCAI-85](https://ijcai.org/proceedings/1985/0069.pdf)
- **[EMPIRICAL REPORT, limit]** XCON kept functioning "quite productively" while its maintainability degraded. Degradation "does not mean that the performance of the system is necessarily impaired". — [Soloway et al., AAAI-87](https://cdn.aaai.org/AAAI/1987/AAAI87-148.pdf)

### Inferences
- **[INFERENCE]** The checks have the following blind spots:
  1. Anomaly detection is syntactic or logical. It cannot detect a rule that is consistent, non-redundant and complete but *wrong about the world*, for example a wrong degradation threshold.
  2. Detection results were established for rule bases *without* numerical certainty measures, so they do not transfer automatically to rules with confidences or probabilities.
  3. Regression on a case library covers only the regions of the input space that the library samples. Rare failure modes, which are often what prognostics most needs, may be absent.
  4. Neither check measures maintainability decay of the XCON kind.
- **[INFERENCE]** A safe admission process for a knowledge change in a predictive ES would combine four steps:
  1. Static anomaly checks on the changed rule base (redundancy and subsumption, conflict, circularity, deficiency).
  2. Full regression on cornerstone and historical cases, with a zero-regression rule for previously correct cases, as in RDR.
  3. Held-out evaluation of any claimed improvement, because of the SEEK2 caveat.
  4. Expert sign-off, recorded with provenance (see question 5).

  This is a synthesis, not a single cited method.
- **[INFERENCE]** Anomaly detection in a nonmonotonic rule base with priorities must treat an intended exception as *not* an "ambivalence" anomaly when a priority resolves it. Classical conflict checks would otherwise flag every exception.

### Gaps
- The specific formal definitions of the anomaly types (for example the exact subsumption criterion), the computational-limit results and the per-system empirical counts in Preece & Shinghal 1994 were not read (full text inaccessible). Only the abstract claims are reported.
- No verified quantitative data on how many real-world errors anomaly detection finds compared with case testing was obtained.
- Mutation-testing or coverage criteria for rule bases were not researched.

## 5. How ontologies and knowledge graphs are versioned and their changes tracked

### Takeaway
Noy & Klein argued that ontology evolution is *not* database schema evolution, for three reasons:
- Ontologies carry explicit semantics and are reused differently.
- The versioning/evolution distinction collapses.
- Compatibility has several dimensions, and diffing must be done semantically.

User-driven ontology evolution management (Stojanovic et al. 2002) and ontology change detection (EKAW 2002) form the classic line of work. W3C PROV-O (2013) is the standard vocabulary for recording which agent produced which knowledge version through which activity. Recent work (TGDK 2023) treats KG evolution as a research problem in observability, metrics and storage that remains open.

### Cited Findings
- **[ESTABLISHED position paper / PROPOSAL]** Noy & Klein, "Ontology Evolution: Not the Same as Schema Evolution", *Knowledge and Information Systems* 6(4):428–440, 2004. The abstract says:
  - The differences from schema evolution stem from "different usage paradigms, the presence of explicit semantics, and different knowledge models".
  - "The traditional distinction between versioning and evolution is not applicable to ontologies."
  - "There are several dimensions along which compatibility between versions must be considered."
  - The set of change operations is different.
  - "We must develop automatic techniques for finding similarities and differences between versions."

  — [VU research portal record (abstract), DOI 10.1007/s10115-003-0137-2](https://research.vu.nl/en/publications/ontology-evolution-not-the-same-as-schema-evolution/)
- **[ESTABLISHED, metadata only]** Stojanovic, Maedche, Motik & Stojanovic, "User-Driven Ontology Evolution Management", EKAW 2002, LNCS 2473, pp. 285–300, DOI 10.1007/3-540-45810-7_27. — [Crossref](https://api.crossref.org/works/10.1007/3-540-45810-7_27). The same volume contains "Ontology Versioning and Change Detection on the Web", pp. 197–212, DOI 10.1007/3-540-45810-7_20. — [OpenAlex](https://api.openalex.org/works/doi:10.1007/3-540-45810-7_20)
- **[ESTABLISHED standard]** PROV-O, "The PROV Ontology", is a W3C Recommendation dated 30 April 2013. It "expresses the PROV Data Model [PROV-DM] using the OWL2 Web Ontology Language". It provides "a set of classes, properties, and restrictions that can be used to represent and interchange provenance information generated in different systems and under different contexts". It "can also be specialized to create new classes and properties to model provenance information for different applications and domains". The namespace is `http://www.w3.org/ns/prov#`. Companion documents include PROV-DM, PROV-N and PROV-CONSTRAINTS (Recommendations), and PROV-DICTIONARY and PROV-SEM (Notes). — [W3C PROV-O](https://www.w3.org/TR/prov-o/)
- **[RECENT, 2023]** "How Does Knowledge Evolve in Open Knowledge Graphs?", *Transactions on Graph Data and Knowledge* 1(1):11:1–11:59, 2023, DOI 10.4230/TGDK.1.1.11. It discusses:
  1. the dimensions of KG evolution;
  2. the observability of evolution in open, collaboratively edited KGs over time;
  3. possible metrics, including methods from network science;
  4. technical approaches "and their current limitations" for storing, analysing and processing large evolving KGs.

  — [OpenAlex abstract](https://api.openalex.org/works/doi:10.4230/TGDK.1.1.11)
- **[RECENT, 2024, orientation only]** "Unifying Large Language Models and Knowledge Graphs: A Roadmap", *IEEE TKDE*, 2024, DOI 10.1109/TKDE.2024.3352100. The abstract states that "KGs are difficult to construct and evolve" and that KGs "can enhance LLMs by providing external knowledge for inference and interpretability". — [OpenAlex record](https://api.openalex.org/works/doi:10.1109/tkde.2024.3352100)

### Inferences
- **[INFERENCE]** For a predictive ES, each knowledge version (a rule set or ontology snapshot) can be represented as a `prov:Entity` and each change as a `prov:Activity`. The engineer, the expert or an automated miner is a `prov:Agent`. Derivation links (`prov:wasDerivedFrom`, `prov:wasGeneratedBy`, `prov:wasAttributedTo`) would then answer audit questions such as "which version and which evidence produced this prognosis?". This is a design mapping onto the standard vocabulary. PROV-O itself does not prescribe it.
- **[INFERENCE]** Noy & Klein's "several dimensions of compatibility" implies that the *ontology* version and the *rule base* version must be tracked jointly. A rule can stay textually unchanged while its meaning shifts because a class it references was redefined.
- **[INFERENCE]** The 2023 TGDK discussion indicates that measuring and storing evolution at scale is still unsolved in general. A bounded industrial ES therefore gains from small, explicitly versioned bases rather than open-ended graph growth.

### Gaps
- The contents of Stojanovic et al. 2002, such as their change-process phases and their handling of change propagation, and the authors and contents of the EKAW 2002 change-detection paper were not verified (Springer blocked access).
- The author list of the TGDK 2023 paper and of the 2024 TKDE roadmap was not verified in this session, so the papers are cited by title and DOI only.
- **LLM-assisted KB or ontology maintenance (2023–2026):** no primary empirical study of LLM-*assisted* rule-base or ontology *maintenance* with measured outcomes was verified within the time budget. Only the generic 2024 roadmap statement above was found. This should be treated as an open gap, not as evidence that LLM-assisted maintenance works.
- OWL-specific versioning mechanisms (for example `owl:versionIRI`, `owl:priorVersion`, `owl:backwardCompatibleWith`) were not verified from the W3C OWL 2 specification in this session.

## 6. How defeasible and nonmonotonic formalisms allow adding exceptions without rewriting existing rules, and the implications for incremental knowledge growth

### Takeaway
Defeasible logic separates the following kinds of knowledge:
- facts;
- strict rules;
- defeasible rules ("typically");
- defeaters, which can only block conclusions;
- an explicit superiority relation between rules.

An exception is added as a new rule plus a priority, and the general rule is left untouched. This is McCarthy's "elaboration tolerance" in practice. Propositional defeasible logic keeps inference linear-time, unlike most nonmonotonic logics. The price is nonmonotonicity: adding knowledge can remove previous conclusions, so every addition needs regression checking. RDR implements the same additive-exception idea operationally.

### Cited Findings
- **[ESTABLISHED definition]** "A formalism is elaboration tolerant to the extent that it is convenient to modify a set of facts expressed in the formalism to take into account new phenomena." McCarthy states that human-level AI will require it. — [McCarthy, "Elaboration Tolerance" (Common Sense 98; revised 1998)](http://jmc.stanford.edu/articles/elaboration.html)
- **[ESTABLISHED definition]** A defeasible theory consists of five kinds of knowledge:
  1. **Facts.**
  2. **Strict rules:** if the premises are indisputable, so is the conclusion.
  3. **Defeasible rules**, which "can be defeated by contrary evidence". Example: `bird(X) ⇒ flies(X)` licenses the conclusion "unless there is other, not inferior, evidence suggesting that it may not fly".
  4. **Defeaters**, which "cannot be used to draw any conclusions"; they only prevent conclusions. Example: `heavy(X) ~> ¬flies(X)`.
  5. **A superiority relation.** Example: for conflicting rules `r: bird(X) ⇒ flies(X)` and `r': brokenWing(X) ⇒ ¬flies(X)`, declaring `r' > r` resolves the conflict.

  — [Maher, "Propositional defeasible logic has linear complexity", arXiv cs/0405090 version of *TPLP* 1(6):691–711, 2001](https://arxiv.org/pdf/cs/0405090)
- **[ESTABLISHED theorem, CONDITIONAL on the propositional case]** "Inference in the propositional form of the logic can be performed in linear time. This contrasts markedly with most other propositional nonmonotonic logics, in which inference is intractable." Full first-order defeasible logic has a recursively enumerable inference problem. — [Maher 2001, DOI 10.1017/S1471068401001168 (OpenAlex abstract)](https://api.openalex.org/works/doi:10.1017/s1471068401001168); [arXiv text](https://arxiv.org/pdf/cs/0405090)
- **[ESTABLISHED]** Defeasible logic generalizes inheritance networks with exceptions under the directly sceptical semantics. It "replaces the implicit specificity relation by an explicit, programmable priority relation" and adds explicit defeaters. — [Maher 2001 (arXiv)](https://arxiv.org/pdf/cs/0405090)
- **[ESTABLISHED]** Defeasible logic is "a simple but efficient formalism for nonmonotonic reasoning based on rules and priorities". Transformations and normal forms for it were used "in the development of an efficient implementation". — [Antoniou, Billington, Governatori & Maher, "Representation results for defeasible logic", *ACM TOCL* 2(2):255–287, 2001 (Crossref abstract)](https://api.crossref.org/works/10.1145/371316.371517)
- **[ESTABLISHED design]** In RDR, a new rule must distinguish the new case from the system's current classification and be justified by the new case's features, while previously correct cases must stay correct. — [Richards 2009](https://api.crossref.org/works/10.1017/S0269888909000241)

### Inferences
- **[INFERENCE]** Defeasible rules with an explicit superiority relation are a good fit for incremental growth of prognostic knowledge. Here is an example:
  1. A general degradation rule ("bearing temperature trend ⇒ elevated wear risk") stays unchanged.
  2. A newly learned exception is added as a separate rule ("after scheduled re-greasing, the temperature transient ⇒ ¬elevated wear risk").
  3. A priority `r_exception > r_general` is added.

  The change is *additive and localized*: one new rule plus one priority edge. Version control can review and roll back such diffs.
- **[INFERENCE]** Nonmonotonicity means that "additive in text" is not "additive in conclusions". A new exception can withdraw conclusions on past cases. Therefore every admitted exception needs the case-library regression from question 4. Explicit priorities are preferable to implicit specificity or rule order, which was the XCON failure mode, because they are inspectable.
- **[INFERENCE]** Linear-time propositional inference is useful for bounded-latency runtime evaluation. It holds for the propositional case only, so first-order rule templates must be grounded, or their size bounded, to keep the guarantee.
- **[INFERENCE]** Defeasible logic (how to *represent* exceptions) and AGM or KM theory (how to *choose* what to give up) answer different questions. A practical system can use defeasible logic for representation, priorities as its entrenchment ordering, and a TMS or ATMS for propagating retractions. This layering is a synthesis, not a single published architecture verified here.

### Gaps
- Formal results on *revising* defeasible theories, that is AGM-style postulates for rule bases with priorities, were not researched or verified in this session.
- Nute's original defeasible logic papers and the SEP entry on non-monotonic logic were not consulted directly. The definitions above come from Maher 2001 and Antoniou et al. 2001.
- No empirical study was found that compares the maintenance cost of exception-based rule bases with that of rewritten monolithic rule bases, apart from the RDR/PEIRS reports in question 3.
