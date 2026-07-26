# Programme management: чому зелені проекти не гарантують здорову програму

Зелений звіт від кожного проекту ще не гарантує здорову програму. Це одна з найнеприємніших істин programme management. Project status може бути чесним у межах конкретної команди, але programme risk часто живе між проектами: у залежностях, integration events, shared resources, supplier constraints, release trains, benefits і рішеннях, які ніхто не бачить як локальну проблему. Найгірше те, що це часто стає очевидним уже після втрати integration window.

Уявімо три проекти. Перший вчасно готує hardware board revision. Другий вчасно завершує firmware scope. Третій вчасно бронює lab testing. Кожен окремо має підстави показати green. Але board revision затримується на кілька днів, firmware build не встигає потрапити в integration window, lab slot пропущено, наступний доступний тільки через три тижні. Локально всі майже праві. Програма втратила milestone.

Саме тому таблиця project statuses не витягує programme management. Programme manager потребує власної моделі: dependencies, integration points, shared risks, cross-project decisions, benefits realisation, gates, resources, scenarios. Без цього він стає редактором звітів, а не керівником системної доставки.

## Project management і programme management

Project management фокусується на delivery конкретного scope: plan, tasks, team, budget, risks, milestones, quality, evidence. Project manager питає: чи виконуємо ми затверджену роботу в межах constraints?

Programme management працює з набором пов'язаних проектів, які разом мають дати business або mission benefit. Programme manager питає: чи всі частини разом створять потрібний результат? Чи синхронізовані dependencies? Чи є shared risks? Чи не конфліктують priorities? Чи буде integration успішною? Чи benefit досі валідний?

Ця різниця здається очевидною, але в реальності programme governance часто перетворюється на aggregation of statuses. PM-и надсилають updates, programme office консолідує, leadership бачить heatmap. Агрегований статус виглядає акуратно, але системних властивостей програми він не показує.

## Чому сума статусів не працює

Status report зазвичай описує локальний стан. Проект може бути green, бо його scope під контролем. Але він може залежати від deliverable, який не контролює. Або використовувати shared expert-а, який перевантажений іншою програмою. Або мати risk, який стає critical тільки у комбінації з risk іншого проекту.

Наприклад, software team готова, hardware team майже готова, QA team готова, але supplier certification document не прийшов. Кожен локальний план виглядає прийнятно, а release gate заблокований. Або два проекти планують integration з одним lab setup в один тиждень. Окремо це не видно. Разом це неможливо.

Programme risk часто emergent. Він виникає у взаємодії, а не в одному work package. Саме тому programme management потребує dependency map разом зі статусами. Треба бачити, хто від кого залежить, які dates синхронізовані, які assumptions спільні, які ресурси shared, які рішення зачіпають кілька проектів.

## Programme objects

У програми мають бути власні об'єкти. Cross-project dependency не має ховатися як comment у статусі. Integration event не варто зводити до milestone одного проекту. Shared risk не є копією ризику в трьох risk registers. Programme gate не складається механічно з project gates. Benefit живе окремо від delivery output.

Cross-project dependency має owner-а з обох сторін, due date, acceptance criteria, impact if late, escalation path. Integration event має entry criteria: які deliverables, versions, environments, test data, people і approvals потрібні. Shared risk має programme-level owner-а і mitigation, яка може виходити за межі одного проекту. Benefit має measurement plan: як зрозуміти, що програма дала очікуваний результат.

Якщо ці об'єкти не існують, programme manager змушений тримати їх у голові або в окремій таблиці. Це працює до першого масштабу, а потім перетворюється на ручну пам'ять організації.

## Cross-project dependencies і shared resources

Найчастіші проблеми програми живуть у dependencies і shared resources. Один project планує deliverable до кінця sprint-а, інший будує на ньому integration plan. Один supplier затримує input, але impact розходиться на кілька streams. Один safety або cybersecurity expert потрібен усім одночасно. Одна lab capacity стає bottleneck для трьох release candidates.

Локальні команди часто не мають повної картини. Вони оптимізують свій план. Це нормально. Programme manager оптимізує систему. Для цього потрібна видимість status разом із contention: хто конкурує за ресурс, які dependencies critical, які delays cascade, які decisions треба прийняти на programme level.

Корисний сигнал - dependency aging. Якщо dependency довго open без movement, вона небезпечна навіть без офіційного red status. Інший сигнал - single point of expertise. Якщо кілька critical path items залежать від однієї людини або одного supplier-а, programme risk вищий, ніж показує будь-який окремий project plan.

## Benefits realisation

Programme management не закінчується delivery. Програма існує для benefit: market launch, platform capability, compliance readiness, cost reduction, operational improvement, mission outcome. Якщо всі проекти здали свої outputs, але benefit не реалізований, програма неуспішна.

Саме тому programme manager має тримати зв'язок між outputs і outcomes. Project A поставив component. Project B поставив integration. Project C поставив training. Але чи customer може реально використовувати capability? Чи operational team готова? Чи support model є? Чи compliance evidence достатня? Чи business case досі актуальний?

Benefits realisation часто губиться, бо delivery легше виміряти. Зроблено або ні. Benefit виміряти складніше. Але без цього programme governance стає delivery factory без стратегічного сенсу.

## Programme-level metrics

Programme dashboard має починатися не з project colors. Йому потрібні programme-level signals: critical cross-project dependencies, integration readiness, shared resource conflicts, benefit confidence, cumulative risk exposure, decision latency, unresolved escalations, gate readiness, supplier bottlenecks, scope interactions.

Наприклад, корисний індикатор - integration readiness by dependency: "для integration event потрібні deliverables A, B, C; A ready, B delayed, C unverified". Це значно точніше, ніж загальне "проект 1 green". Інший індикатор - decision aging: скільки programme decisions pending і який impact затримки. Ще один - benefits confidence: чи залишаються assumptions business case валідними.

Такі metrics не завжди легко зібрати, бо вони потребують зв'язків між project systems. Але без них programme governance бачить симптоми пізно.

## AI для programme analysis

AI може допомогти programme manager-у знаходити cross-project patterns. Повторювані delays у одного supplier-а. Bottleneck у певній discipline. Dependencies, які згадуються в різних status reports різними словами. Risks, які локально low, але разом створюють programme-level exposure. Зміни requirements, які зачіпають кілька release trains.

Добрий assistant не має казати "program healthy" без пояснення. Він має показати sources: project updates, dependency records, open risks, integration plan, test status, decisions. Він має формулювати assumptions і confidence. Наприклад: "programme integration risk high because deliverable D-14 delayed, lab slot fixed, firmware build unverified, and decision on supplier waiver pending; confidence medium because hardware status update is 5 days old".

AI особливо корисний для language normalization. Один PM пише "minor delay", інший "dependency at risk", третій "waiting for supplier". Assistant може побачити, що це один pattern, але лише якщо має доступ до structured context і не вигадує зв'язки.

## Що має показувати programme dashboard

Я б хотів бачити dashboard, побудований навколо programme questions. Чи буде integration event готовий? Які dependencies можуть зірвати milestone? Які shared resources перевантажені? Які decisions pending? Які benefits під ризиком? Які assumptions змінилися? Які risks живуть між проектами?

Project statuses теж потрібні, але як нижчий рівень деталізації. Спершу programme health, потім project drill-down. Інакше leadership бачить десять окремих дерев і не бачить лісу.

Важливо також мати dissenting signals. Якщо всі проекти green, але dependency map показує critical open links, dashboard має це сказати. Якщо local owners не погоджуються з programme risk assessment, це теж треба фіксувати. Programme governance має витримувати конфлікт даних, навіть коли дуже хочеться згладити його в красивий summary.

## Programme governance cadence

Programme review швидко втрачає сенс, коли перетворюється на тур по проектах: "проект A, ваш статус; проект B, ваш статус". Краще будувати cadence навколо programme questions. Чи готове наступне integration event? Які dependencies змінили confidence? Які shared resources перевантажені? Які decisions aging? Які benefits під ризиком? Які assumptions треба переглянути?

Окремий формат корисний для dependency review. Там не треба обговорювати весь scope кожного проекту. Потрібно пройти critical cross-project links: owner, due date, evidence of readiness, impact if late, escalation. Це коротша, але набагато корисніша розмова для programme health.

Ще один формат - decision review. Programme manager має бачити risks разом із завислими decisions. Нерідко програма запізнюється не через engineering work, а через те, що ніхто не прийняв trade-off: cut scope, додати ресурс, прийняти waiver, змінити sequence, перенести integration. Decision latency - це programme metric, яку недооцінюють.

Мінімальний корисний артефакт для такої роботи - programme dependency register: живий список critical links із owners, dates, acceptance criteria і impact. Якщо він оновлюється перед programme review і пов'язаний з integration events, розмова одразу стає конкретнішою. Команда менше сперечається про загальний колір статусу і більше про те, який link треба закрити першим.

Це простий інструмент, але він змінює поведінку: залежність із фрази в статусі перетворюється на керований об'єкт.

А керований об'єкт уже можна вимірювати, ескалювати і закривати.

## Висновок

Programme management починається там, де сума project status reports уже не пояснює реальність. Це управління системою залежностей, рішень, ризиків, integration events і benefits. Локально успішні проекти можуть разом створити невдалу програму, якщо programme-level об'єкти не видно.

AI та integrated systems можуть допомогти, але тільки якщо програма описана як модель, а не як колекція презентацій. Programme manager-у потрібна карта залежностей, shared resources, gates, benefits і decision latency. Тоді зелений статус проектів уже не заспокійливий шум, а лише один із сигналів у ширшій картині.

Питання до читачів: де у ваших програмах найчастіше губляться cross-project dependencies - між командами, supplier-ами, labs, architecture, release trains чи leadership decisions?
