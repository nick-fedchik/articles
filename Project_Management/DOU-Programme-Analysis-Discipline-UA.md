# Програмна аналітика: дисципліна між project management, архітектурою і AI

Є речі, які в проектах ніби всі розуміють, але довго не мають власної назви. Ми говоримо про dependencies, architecture decisions, risks, requirements, integration gaps, release compatibility, operational readiness. Ми пишемо нотатки, робимо рев'ю, створюємо epics, переносимо висновки у backlog. Але часто це живе як набір локальних документів, а не як окрема дисципліна.

Мені дедалі більше подобається назва "програмна аналітика". Важливе уточнення: тут "програмна" не в сенсі software analytics. Йдеться про programme analysis - аналітику на рівні програми, де кілька продуктів, команд, репозиторіїв, протоколів і workflow мають працювати разом.

У software/project management ми давно навчилися трекати tasks, defects, requirements, risks, tests і releases. Але коли система стає крос-продуктовою, з'являється інший шар роботи: треба аналізувати не лише окремий проект, а зв'язки між проектами. Саме тут і потрібна програмна аналітика.

## Проблема: рішення губляться між репозиторіями

Уявімо продуктову сім'ю. Один сервіс відповідає за project state і traceability. Другий - за expert-system reasoning. Третій - за knowledge acquisition. Четвертий - за GitLab automation. Окремо є shared messaging contract. Кожен репозиторій має свої docs, backlog, requirements, tests, changelog.

Локально все виглядає нормально. Але справжні питання живуть між ними:

- хто є owner-ом canonical ID;
- який компонент має відповідати за degraded-mode behavior;
- чи однаково всі трактують readiness state;
- як release evidence з одного проекту потрапляє у programme-level compatibility claim;
- чи не роз'їхалися NATS subjects, payload fields і CLI diagnostics;
- де має жити shared SDK contract;
- які backlog items блокують один одного;
- чи має KAS запускати локальний LLM, чи це відповідальність AFES;
- чи APMS має виконувати inference, чи лише показувати diagnostics.

Якщо ці питання вирішувати прямо в локальних backlog files, reasoning швидко розпадається. APMS отримує один фрагмент рішення, AFES - другий, KAS - третій. Через місяць важко пояснити, чому boundary саме такий. Через два місяці з'являється новий contributor і питає: "а чому це не в shared package?" І всі знову відновлюють дискусію з пам'яті.

Програмна аналітика потрібна саме для того, щоб не втрачати цей середній шар reasoning.

## Що таке програмна аналітика

Програмна аналітика - це керована робота з фактами, gaps, decisions і promotion plan на рівні всієї програми до того, як зміни розходяться по окремих проектах.

Вона не замінює requirements engineering, architecture, PMO, QA або product management. Вона стоїть між ними і зшиває їх.

У практичному вигляді це означає, що для кожної крос-проектної теми створюється analysis artifact:

- problem statement;
- questions to answer;
- affected projects;
- evidence sources;
- current cross-project model;
- gaps and risks;
- decisions or recommendations;
- promotion plan;
- validation plan;
- residual risks.

Це звучить формально, але насправді дуже приземлено. Ми не пишемо документ заради документа. Ми створюємо місце, де зібрана логіка рішення. Потім із цього місця вже можна промоутити конкретні SHR/SWR, ADR, backlog tasks, protocol updates, CLI requirements, tests або documentation changes.

## Чим це відрізняється від architecture review

Architecture review часто відповідає на питання "як має бути влаштована система". Програмна аналітика ширша. Вона питає: які проекти зачеплені, де ownership boundary, які контракти змінюються, які requirements треба оновити, які tests доводять інтеграцію, які release evidence fields потрібні, які residual risks залишаються.

Наприклад, якщо ми досліджуємо локальні LLM runtime-и, architecture discussion може швидко перейти в "Ollama vs llama.cpp vs ONNX". Це корисно, але недостатньо. Programme analysis додає інші питання:

- AFES чи KAS має бути owner-ом runtime profile;
- чи не суперечить це попередньому NPU/OpenVINO рішенню;
- чи треба APMS тільки показувати diagnostics;
- які runtime evidence fields потрібні для release compatibility;
- чи це має стати CLI command, docs task, SWR або лише research note;
- як довести, що fallback не був помилково названий acceleration.

Тобто програмна аналітика не зупиняється на технологічному виборі. Вона доводить вибір до governance, traceability і validation.

## Чим це відрізняється від status reporting

Status reporting каже: де ми зараз. Програмна аналітика каже: як пов'язані рішення, gaps і work packages між проектами.

У status report можна написати "release evidence automation in progress". Але programme analysis має відповісти:

- який evidence schema є canonical;
- які поля має продукувати APMS;
- які поля має продукувати AFES;
- що має продукувати KAS;
- що саме робить GLTK;
- хто відповідає за redaction;
- чи publication dry-run є обов'язковим;
- що робити, якщо readback hash не збігся;
- які release claims не можна робити без smoke evidence.

Це різні рівні мислення. Status reporting описує activity. Programme analysis формує operating model.

## Життєвий цикл аналізу

Мені подобається простий lifecycle.

Intake. Ми формулюємо тему. Наприклад: "як AutoForge має використовувати локальні LLM runtime-и" або "як release evidence має працювати між APMS, AFES, KAS і GLTK".

Evidence collection. Читаємо код, docs, requirements, backlog, protocol specs, CLI behavior, tests, changelog. Важливо не писати opinion до того, як подивилися на реальність.

Synthesis. Розділяємо facts, gaps, risks і proposals. Це критично. Якщо змішати facts і recommendations, документ швидко стає політичним текстом, а не аналітикою.

Decisions. Формулюємо ownership boundary і recommended direction. Наприклад: "Ollama remains default; llama.cpp becomes advanced AFES profile; ONNX/OpenVINO remains embed/rerank path; APMS reports diagnostics only".

Promotion. Перетворюємо analysis у project-local work: APMS SWR, AFES task, KAS architecture note, GLTK guide, afnats contract, central Wiki update.

Validation. Перевіряємо links, diagnostics, tests, smoke commands, fixtures, JSON schemas, release evidence. Якщо щось не перевіряли, чесно пишемо чому.

Closed or residual risks. Analysis не має вдавати, що все вирішено. Якщо лишився risk, він має бути видимим.

## Promotion queue як операційний інструмент

Найсильніша частина програмної аналітики - promotion queue. Це не backlog у звичайному сенсі. Це список packages, які переводять cross-project reasoning у конкретні зміни.

Наприклад:

- shared contract governance;
- release compatibility evidence schema;
- diagnostics and degraded-mode taxonomy;
- trace/correlation propagation;
- classification and movement policy;
- CLI functional update requirements;
- local LLM runtime profile;
- APMS mathematical substrate.

Кожен package має source analysis, owners, priority, state, expected output, validation expectation. Це рятує від типового chaos: "ми ж це вже обговорювали". Так, обговорювали. Де source? Який decision? Куди promote? Яка validation?

Promotion queue робить аналітику action-oriented. Документ не залишається красивою думкою. Він має шлях у project-local artifacts.

## Де тут AI

AI дуже корисний у програмній аналітиці, але не як генератор фінальних рішень. Його сила - швидко читати багато context-у, знаходити inconsistencies, порівнювати patterns, формулювати hypotheses, готувати synthesis drafts.

Наприклад, AI може помітити, що APMS уже має release evidence fields, AFES має smoke evidence JSON, KAS має evidence matrix, а GLTK має publication dry-run/readback. Людина з цього робить governance decision: кожен компонент має продукувати evidence, GLTK має публікувати і перевіряти publication, central docs мають тримати schema.

AI також допомагає з drift detection. Якщо PA-012 каже, що DirectML не є NPU evidence, а KAS docs десь назвали DirectML NPU support, це треба підсвітити. Якщо APMS diagnostics використовують `ok/warning`, а programme taxonomy вимагає `ready/degraded/unavailable`, це gap. Якщо release evidence claims active chat relay, а protocol doc каже draft/deferred, це risk.

Але AI має працювати з evidence. Добра програмна аналітика вимагає source links, confidence, assumptions і residual risks. Без цього AI просто прискорює chaos.

## Програмна аналітика як пам'ять програми

У великих програмах знання часто живе в людях. Це небезпечно. Не тому, що люди погані, а тому, що контекст надто великий. Через кілька місяців навіть сильна команда забуває, чому певне рішення було прийняте.

Програмна аналітика створює пам'ять програми:

- чому AFES owns model-runtime PoC;
- чому KAS не має ставати chat provider;
- чому APMS reports acceleration but does not run inference;
- чому release evidence має redaction і readback;
- чому trace fields мають canonical names;
- чому shared SDK має бути owner-ом NATS subjects;
- чому programme-level math substrate має бути deterministic before ML.

Це особливо важливо для InnerSource. Новий contributor має бачити не лише "що зробити", а і "чому boundary такий". Інакше він чесно запропонує refactor, який команда вже відхилила після трьох днів аналізу.

## Що варто трекати як дисципліну

Я б виділив кілька постійних напрямків програмної аналітики.

Capability and ownership map. Хто owns яку capability, хто consumes, де shared contract.

Data ownership and identifiers. Які IDs canonical, які aliases, які legacy fields ще підтримуються.

End-to-end value streams. Як APMS, AFES, KAS, GLTK і messaging разом виконують реальний workflow.

Requirements traceability audit. Чи всі cross-project flows мають SHR/SWR coverage.

Release compatibility and evidence. Які component versions сумісні, який smoke це довів, які residual gaps лишилися.

Operational readiness. Які стани ready/degraded/unavailable, які reason codes, які fallback-и.

Security and trust boundaries. Який content можна рухати між локальним AI, GitLab, NATS, reports і Wiki.

Observability and trace propagation. Чи один request можна простежити через всю систему.

Backlog dependency analysis. Які epics/tasks блокують один одного між репозиторіями.

Runtime and compute research. Де Ollama, llama.cpp, ONNX, OpenVINO, GPU/NPU, CPU fallback.

Mathematical substrate. Trace matrices, risk algebra, compliance coverage, evidence graphs.

Це вже не випадкові notes. Це programme analysis portfolio.

## Ризики самої дисципліни

Будь-яка дисципліна може стати бюрократією. Програмна аналітика теж.

Перший ризик - over-analysis. Якщо кожне дрібне рішення перетворюється на PA-document, команда зупиниться. Потрібен поріг: analysis потрібен там, де є cross-project boundary, shared contract, release evidence, compliance impact або ownership uncertainty.

Другий ризик - stale documents. Якщо analysis не має status, updated date, promotion log і residual risks, він швидко старіє. Старий analysis небезпечніший за відсутній, бо створює фальшиву впевненість.

Третій ризик - змішування facts і бажань. У документі має бути видно, що ми знаємо, що припускаємо, що пропонуємо і що вже вирішено.

Четвертий ризик - відсутність promotion. Якщо analysis не веде до APMS/AFES/KAS/GLTK docs, requirements, tests або backlog, він стає intellectual exercise. Це може бути цікаво, але programme не рухається.

## Як почати без великої методології

Не треба починати з великого framework. Достатньо чотирьох артефактів.

Analysis register. Список тем, статусів, owners і primary outputs.

Analysis template. Одна структура для problem, evidence, gaps, decisions, promotion, validation.

Promotion queue. Список packages, які переводять analysis у project-local зміни.

Validation habit. Link checks, diagnostics, tests або явне пояснення, що не перевірялося.

Цього вже достатньо, щоб команда почала думати системніше.

## Мій висновок

Програмна аналітика - це дисципліна для систем, де один репозиторій уже не пояснює реальність. Вона потрібна там, де project management, architecture, requirements, compliance, AI і release engineering перетинаються між кількома продуктами.

Її цінність не в тому, щоб писати більше документів. Її цінність у тому, щоб рішення мали пам'ять, ownership, validation і шлях у виконання. Коли analysis register, evidence, promotion queue і release evidence працюють разом, програма перестає бути набором локальних правд. Вона стає керованою системою.

Для мене це один із найважливіших шарів у сучасному engineering management. Особливо там, де AI вже входить у роботу. Бо AI без програмної аналітики швидко стане генератором красивих локальних відповідей. А з програмною аналітикою він може стати помічником у системному reasoning.

## Питання до читачів

- Де у ваших програмах живе reasoning між проектами - у документах, backlog, головах людей чи чатах?
- Чи маєте ви окремий register для cross-project analyses?
- Як ви переводите architecture/research findings у project-local requirements і tests?
- Які теми у вашій організації найбільше просяться в programme analysis: release compatibility, traceability, security boundaries, AI runtime, dependencies чи readiness?
