# Project Instructions: Language, Style & AI Marker Inspection

## 1. Author Persona & Voice
- **Author:** Mykola Fedchyk (Микола Федчик).
- **Domain:** Mission-critical software engineering, functional safety (ISO 26262, IEC 61508, DO-178C), systems engineering, Linux kernel architecture, and evidence-governed neuro-symbolic systems.
- **Voice:** Authoritative, rigorous, engineering-grounded Ukrainian practitioner voice. Never sound like a generic AI translator.

## 2. Strict Prohibition: The «Контур» Machine Buzzword
In Ukrainian technical prose, the word **«контур»** is naturally and properly used **only** in specific physical or mathematical contexts:
- **Geometry & Image Processing:** контури геометричних фігур, оконтурення об'єктів, контурні карти рельєфу / перепадів яскравості (contour lines), контурні дескриптори.
- **Electrical & Physical Plumbing:** замкнені контури заземлення (ground loops) або фізичні гідравлічні контури замкнених трубопроводів (якщо мова суто про труби/циркуляцію рідини).

### Machine-Generated Clichés to ALWAYS Avoid:
Do **NOT** use «контур» for software architectures, cybernetic loops, security perimeters, execution environments, or pipelines. Always replace with natural Ukrainian technical terminology:
- ❌ *«контур безпеки»* → ✅ **система безпеки**, **периметр безпеки**, **межа безпеки**
- ❌ *«двоконтурний шлюз / підхід / перевірка»* → ✅ **дворівневий шлюз**, **двоетапна перевірка**, **двоканальний тракт**
- ❌ *«апаратний контур»* → ✅ **апаратне середовище**, **апаратний комплекс**, **апаратна підсистема**
- ❌ *«контур керування / регулювання»* → ✅ **цикл керування**, **ланцюг керування**
- ❌ *«людина в контурі (Human-in-the-Loop)»* → ✅ **людина в циклі керування (Human-in-the-Loop)**
- ❌ *«контур допуску»* → ✅ **шлюз допуску**, **контрольно-пропускний шлюз**
- ❌ *«контур навчання / донавчання»* → ✅ **конвеєр навчання**, **цикл навчання**, **процес донавчання**
- ❌ *«оперативний контур»* → ✅ **оперативне середовище**, **оперативний режим**
- ❌ *«виконавчий контур»* → ✅ **виконавчий механізм**, **виконавчий тракт**
- ❌ *«контур інженерного моделювання»* → ✅ **середовище інженерного моделювання**
- ❌ *«контур сприйняття / міркування»* → ✅ **тракт сприйняття**, **рівень міркування**
- ❌ *«тестування з обладнанням у контурі (HIL)»* → ✅ **моделювання / тестування з апаратурою в циклі (HIL)**

## 3. Pedagogical Quality Criteria
- **Zero-context opening:** Every chapter must start with immediate practical grounding — what mission-critical failure occurs if this topic is misunderstood, and how the chapter resolves it.
- **Fatigue relief:** Mathematical formulas must always be paired with intuitive analogies, Mermaid diagrams, and parameter definitions.
- **Concrete grounding:** Reference concrete standards (ISO 26262, ISO/SAE 21434, ASPICE, RFC), real hardware (Infineon AURIX, Xilinx FPGA, ARM Cortex), and formal methods (SMT Z3, Datalog, GSN).

## 4. Scientific Editing, Academic Voice & Authorial Balance (The Knuth/Kleppmann Paradigm)
- **80% Objective Impersonal Exposition:**
  - The primary narrative mode must be objective and depersonalized, using natural Ukrainian impersonal verbal forms on **-но / -то**: *«досліджено»*, *«запропоновано»*, *«верифіковано»*, *«побудовано»*, *«встановлено»*, *«експериментально виміряно»*.
  - Use system-driven agents: *«алгоритм гарантує»*, *«архітектура забезпечує»*, *«результати моделювання свідчать»*.
- **20% Authorial Practical Voice (Findings, Benches, Caveats & Guidelines):**
  - When presenting experimental test benches, empirical results, mission-critical production caveats, or architectural recommendations, employ a restrained, authoritative practitioner voice:
    - *«Автором розгорнуто дослідницький стенд...»*
    - *«Автором запропоновано підхід / доведено, що...»*
    - *«З інженерного досвіду автора випливає, що...»*
    - *«Автор рекомендує для прототипів...»*
  - **Strictly eliminate artificial collective «ми» (Pluralis Auctoris):** Never write *«нами розроблено»*, *«наша група розгортає»* (this is a single-author monograph). Formulas of joint pedagogical reasoning with the reader (*«розглянемо приклад»*, *«звернемо увагу на»*) remain valid.
  - **Eliminate colloquial «я-стиль»:** Avoid informal blogging phrasing (e.g. *«я написав скриптик»* $\rightarrow$ *«Автором реалізовано програмний модуль...»*).
- **Zero Breaking Changes to Book Structure & Code:**
  - Never mechanically rename `# Глава N` to `РОЗДІЛ N` or alter file names / anchors — preserve all cross-references and table-of-contents integrity.
  - Never tear code out of chapters into separate detached files if it serves a direct pedagogical purpose; format long code listings (> 45 lines) inside collapsible `<details><summary>` blocks.
  - Preserve LaTeX equations ($...$ and $$...$$) and ensure all variables are defined immediately below.
  - Always verify with `python3 scripts/inspect_buzzwords.py` after editing.

## 5. Structural Hierarchy, Numbering & Heading Semantics
- **Strict 4-Level Numbering Hierarchy:**
  - **Level 1 (H1):** `# Глава N. Назва глави` або `# Додаток X. Назва додатка` (preserves IDs for anchors and cross-references).
  - **Level 2 (H2) — Розділ:** `## 1. Назва розділу`, `## 2. Назва розділу` (numbered sequentially).
    - **Structural Exceptions (Unnumbered):** Meta and closing sections MUST remain unnumbered: `## Анотація`, `## Висновки`, `## Глосарій`, `## Абревіатури` (або `## Скорочення та терміни`), `## Джерела` (або `## Література` / `## Посилання`).
  - **Level 3 (H3) — Підрозділ:** `### 1.1. Назва підрозділу`, `### 1.2. Назва підрозділу`.
  - **Level 4 (H4) — Пункт:** `#### 1.1.1. Назва пункту`, `#### 1.1.2. Назва пункту`.
  - No headers deeper than H4; use numbered or bulleted lists inside paragraphs.
- **Subject Engineering Semantics in Headings:**
  - Titles must strictly convey engineering and scientific substance (e.g. *«## 1. Історичні межі першого покоління експертних систем та передумови ренесансу»* instead of *«## Хіба вони не померли у 1980-х?»*).
  - **Relocation of Provocative Questions & Slogans:** Motivational slogans and provocative questions are NEVER used as titles; they MUST be moved into the body text, primarily into the first paragraph where the engineering problem is introduced.
- **Fatigue-Free Mathematics for Practicing Engineers:**
  - Every standalone formula must be isolated in `$$...$$` or ````math`.
  - Immediately below the formula, provide a comprehensive breakdown of every parameter AND its physical/intuitive action, tailored for practicing systems engineers who haven't dealt with pure academic math syntax daily.
- **Rich Diagrams & Alerts:**
  - Leverage diverse GitHub-supported Mermaid diagrams (`flowchart`, `sequenceDiagram`, `stateDiagram-v2`, `classDiagram`, `xychart-beta`) with `accTitle` / `accDescr`.
  - Use GitHub-style callouts (`> [!NOTE]`, `> [!IMPORTANT]`, `> [!WARNING]`, `> [!TIP]`).
- **Batch Processing & Quality Gate:**
  - Process chapters in structured batches by parts (Part 1, Part 2, etc.).
  - Run `python3 scripts/inspect_buzzwords.py` after editing each batch.
  - Commit with clear semantic commit messages and push to `origin/main` after each batch.


