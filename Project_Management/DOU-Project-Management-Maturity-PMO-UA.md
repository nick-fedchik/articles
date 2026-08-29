# Модель зрілості систем управління проектами: як з'являється управлінська пам'ять

Зрілість системи управління проектами не вимірюється кількістю дашбордів. Можна мати десятки графіків, status colors, automated reports і все одно не бачити головного: які рішення треба прийняти, які ризики ростуть, які докази відсутні, які залежності стають критичними і де організація повторює старі помилки.

Трекер задач - хороший початок. Jira, GitLab Issues, GitHub Issues, Azure DevOps, YouTrack або подібні інструменти можуть бути міцною основою для execution. Вони допомагають бачити задачі, owners, statuses, deadlines, milestones. Але складні engineering projects дуже швидко виходять за межі задач. Там є вимоги, baselines, risks, tests, defects, suppliers, release gates, compliance evidence, architecture decisions, CCB, programme dependencies і lessons learned.

У якийсь момент організація відкриває неприємну правду: процеси ніби описані, інструменти ніби є, PMO ніби працює, але реальний стан проекту все одно збирається вручну перед review або аудитом. Це і є низька зрілість інформаційної системи, навіть якщо process maturity на папері виглядає пристойно.

## Процесна зрілість і зрілість системи

Корисно розрізняти дві речі. CMMI, Automotive SPICE, ISO/IEC 330xx та подібні моделі говорять про зрілість процесів: чи процес визначений, керований, вимірюваний, повторюваний, оптимізований. Це важливо. Але навіть хороший процес може погано працювати, якщо інформаційна система не підтримує його у щоденній роботі.

Наприклад, risk management process може бути описаний чудово. Є roles, escalation rules, review cadence, risk matrix. Але якщо ризики живуть в Excel, mitigation actions у задачах, metrics у dashboard, а decisions у листуванні, система не допомагає процесу. PM мусить вручну зводити картину. PMO мусить питати статус. Leadership бачить summary, але не завжди бачить джерела.

Модель зрілості project/programme management system ставить практичніше питання: чи підтримує система процес без зайвої ручної реконструкції. Чи збираються метрики автоматично? Чи видно відхилення вчасно? Чи пов'язані вимоги, задачі, ризики, тести і evidence? Чи може PMO бачити реальний стан без прохання зробити ще одну презентацію?

## PMMM як карта, а не відповідь

Ця тема добре перетинається з напрямом, який Jana Kostalova і John McGrath окреслюють у дослідницькій пропозиції `Available Project Management Maturity Models Worldwide and its Comparison` у збірнику 13th IPMA Research Conference. Їхній фокус - як орієнтуватися у великій кількості Project Management Maturity Models, які критерії вони оцінюють і як ці критерії змінюються під впливом AI adoption, digital transformation, sustainability та project complexity.

Для практики це важлива рамка, але вона не закриває всю проблему. PMMM може показати рівень процесної зрілості, допомогти знайти прогалини і запустити періодичне reassessment. Але в engineering organization зрілість проявляється не тільки в тому, що процес описаний і оцінений. Вона проявляється в тому, чи може система за кілька хвилин відповісти на питання: що змінилося після baseline, який ризик виріс, які докази відсутні, хто прийняв рішення і які залежності стали критичними.

Тому в цій статті я дивлюся не на PMMM як на каталог моделей, а на практичний наступний шар: зрілість інформаційної системи, яка щодня несе PM, PgM і PMO процес. Якщо PMMM - це карта розвитку capability, то project/programme management system - це місце, де ця capability або працює, або розпадається на Excel, чати і презентації.

## Рівень 0: індивідуальне управління

На нульовому рівні проект тримається на особистій дисципліні PM, engineering lead-а і кількох сильних людей. Плани живуть у таблицях, рішення у чатах, ризики в голові, статуси в усних домовленостях. Це може працювати для маленької команди, де всі поруч і контекст спільний.

Проблема проявляється після першого масштабу: хтось іде у відпустку, змінюється команда, з'являється audit question, приходить новий PM. Раптом виявляється, що організація не має пам'яті. Вона має людей, які пам'ятають. Це різні речі.

## Рівень 1: task tracking

На першому рівні з'являється backlog: задачі, assignees, statuses, labels, deadlines, milestones. Це великий крок уперед. Команда бачить поточну роботу, може координувати execution, контролювати basic flow.

Але task tracking не відповідає на багато engineering questions. Чи закрита вимога, якщо закрита задача? Чи пройдені потрібні тести? Чи змінився risk після failed regression? Чи вплинув defect на release gate? Чи є approval для waiver? Task tracker сам по собі не тримає повну модель проекту.

Тому перший рівень часто створює ілюзію контролю. Задачі рухаються, dashboard оновлюється, але evidence chain живе десь поруч.

## Рівень 2: integrated project execution

На другому рівні задачі починають зв'язуватися з реальними engineering artifacts: repositories, merge requests, CI/CD, requirements, tests, defects, documentation, release artifacts. Система бачить work item разом із частиною доказів виконання.

Це момент, коли status report можна частково збирати з джерел, а не з пам'яті PM. Наприклад, задача закрита, але пов'язаний test failed. Requirement implemented, але verification не завершена. Merge request merged, але security scan має finding. Такі сигнали мають з'являтися автоматично.

Integrated execution наближає project management до engineering reality. Але programme-level dependencies, governance, portfolio choices і PMO rules ще лишаються окремою роботою. Це сильна основа, після якої починається наступний рівень зрілості.

## Рівень 3: programme coordination

На третьому рівні кілька проектів об'єднуються в programme view. З'являються cross-project dependencies, shared risks, integration events, release trains, programme gates, benefits і decision records. Це важливо, бо programme management не є сумою project statuses.

Кілька проектів можуть бути green локально і red системно. Один deliverable залежить від board revision, другий від firmware build, третій від lab slot. Якщо ці залежності не змодельовані, programme manager бачить красиву таблицю і не бачить наближення проблеми.

Programme coordination потребує власних об'єктів: dependency, integration event, shared risk, benefit, programme decision. Якщо зробити тільки dashboard над project tasks, рівень 3 буде декоративним.

## Рівень 4: PMO governance layer

На четвертому рівні PMO-практики стають частиною системи. Templates, WBS patterns, stage gates, risk taxonomy, metrics catalogue, CCB workflow, resource governance, audit evidence model, lessons learned і methodology checks живуть не в окремих документах, а в робочому процесі.

Це змінює роль PMO. Замість офісу, який просить презентації, він працює як governance engine: визначає, які артефакти canonical, які metrics важливі, які thresholds запускають escalation, які gates блокують release, які AI-рекомендації можна використовувати як draft, а які потребують human approval.

На цьому рівні PMO працює ближче до systems engineering, quality, safety, cybersecurity, architecture і data governance. Його завдання - зробити governance видимою і виконуваною в системі, не додаючи зайвої бюрократії.

## Рівень 5: decision intelligence

Найвищий рівень - decision intelligence. Система має knowledge graph, AI Assistant, expert rules, scenario planning, automated risk monitoring, explainable recommendations і decision dossiers. Рішення все одно приймають люди. Система знімає туман навколо альтернатив, припущень, наслідків, sources і confidence.

Наприклад, change request може автоматично показати impacted requirements, tests, risks, suppliers, release gates, related historical decisions і similar incidents. AI може підготувати summary, а expert rules можуть перевірити formal gaps: missing approval, broken traceability, outdated baseline. Результатом стає decision dossier, який людина може review, approve або reject.

Це важливо: зріла система не забирає відповідальність. Вона прибирає туман навколо відповідальності.

## AI працює після дисципліни даних

AI на низькому рівні зрілості часто просто пришвидшує хаос. Якщо requirements не пов'язані з tests, risks живуть у таблиці, decisions у чатах, а evidence не має versions, assistant може написати красивий summary, але не зможе дати надійний висновок.

Тому AI варто впроваджувати як capability, яка росте разом із data discipline, а не як магічний шар поверх усього. На рівні task tracking він може підсумувати backlog. На рівні integrated execution - знайти розбіжності між задачами, merge requests і tests. На programme рівні - побачити cross-project patterns. На PMO governance рівні - перевіряти completeness evidence. На decision intelligence рівні - готувати пояснювані рекомендації.

## Хто має оцінювати зрілість

Зрілість такої системи не має оцінювати тільки IT або тільки PMO. IT бачить архітектуру, інтеграції, безпеку і підтримуваність. PMO бачить methodology. Але реальна якість проявляється в тому, як система працює для PM, programme managers, engineers, QA, safety, cybersecurity, requirements, architecture і leadership.

Практична оцінка має включати maturity workshop, audit dry-run, traceability check, metric freshness review, change-control path, risk escalation path і friction analysis. Де люди ведуть паралельний Excel? Де PMO просить слайд, хоча дані вже є? Де AI відповідає без sources? Де зміна methodology потребує місяців?

Окремий критерій - швидкість модернізації. Якщо templates, metrics, gates, risk rules або AI policies змінюються повільніше, ніж реальні проекти, люди знову підуть в обходи. Зріла система має бути контрольовано змінною: versioning, sandbox, pilot rollout, audit trail для змін governance.

## Купити, розробити чи поєднати

Для таких систем рідко працює просте "купити або написати". Vendor-рішення дає готові capabilities і підтримку, але може погано лягати на локальну модель governance. Внутрішня розробка дає контроль, але потребує сильної product ownership і довгострокової підтримки.

Найчастіше потрібна hybrid model: компанія володіє доменною моделлю, process rules, data semantics, evidence model і integration architecture, а реалізація може комбінувати готові tools, custom integrations і internal platform. Навіть якщо vendor робить більшість work, організація не має віддавати йому право визначати, що для неї означає project health.

## Практична перевірка зрілості

Щоб оцінити систему без великої методологічної програми, можна пройти кілька реальних сценаріїв. Візьміть один change request і перевірте, чи можна від нього перейти до impacted requirements, tasks, tests, risks, approvals і release notes. Візьміть один failed test і подивіться, чи видно requirement, component, owner, defect і release impact. Візьміть один milestone shift і перевірте, чи система показує supplier dependencies, lab booking, impacted gates і stakeholder communication.

Потім варто зробити audit dry-run. Попросіть систему відповісти: яка вимога змінилася після baseline, хто approve, які tests виконані, які waivers відкриті, які risks прийняті. Якщо відповідь збирається вручну з п'яти місць, maturity нижча, ніж здається.

Окремо запитайте людей, де вони ведуть shadow records. Паралельний Excel, приватна таблиця ризиків, personal notes по dependencies, локальний список approvals - це сигнал, що офіційна система не закриває реальний workflow. Не поспішайте називати це порушенням дисципліни. Саме такі обходи найкраще показують, куди треба модернізувати систему.

Ще один корисний сигнал - час відповіді на управлінське питання. Якщо "чи готові ми до gate" потребує двох днів ручного збору, система не на рівні decision intelligence. Якщо відповідь формується за хвилини і веде до джерел, maturity стала відчутною.

## Висновок

Зріла project/programme management system впізнається за здатністю підтримувати реальний процес: execution, governance, evidence, decisions, risks, dependencies і organisational memory. Красивий dashboard тут другорядний. Така система допомагає PMO бути архітектором управління, а не збирачем статусів.

Task tracker - хороший старт. Але шлях до decision intelligence проходить через traceability, integrated execution, programme objects, PMO governance і data discipline. AI має сенс тоді, коли система вже має що пояснювати. Красивий переказ хаотичних даних зрілістю не стає.

Питання до читачів: на якому рівні зрілості, на вашу думку, перебувають системи управління проектами у ваших організаціях - і де вони найчастіше змушують людей повертатися до Excel, чатів або ручних презентацій?

## Посилання

- [ISO 21502:2020 — guidance on project management](https://www.iso.org/standard/74947.html)
- [ISO 21505:2017 — guidance on governance](https://committee.iso.org/sites/tc258/home/projects/published/iso-21503.html)
