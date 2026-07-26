# Event-driven project management: чому система має реагувати на події, а не чекати звіту

У більшості організацій project management працює за календарем. Є weekly status meeting, є monthly steering committee, є quarterly portfolio review. Усе це корисно, але має одну фундаментальну проблему: проблеми проекту виникають не за розкладом цих зустрічей. Failed test з'являється у вівторок ввечері, supplier повідомляє про затримку в середу зранку, ризик переходить у червону зону у п'ятницю. Управлінська реакція ж відбувається у понеділок, після того як PM зведе звіт.

Я бачив це багато разів: команда тижнями жила з проблемою, яка була видимою у системах одразу, але не стала управлінським сигналом, бо ніхто не дивився у цей куток до підготовки звіту. Це і є типове calendar-driven status management. І воно завжди запізнюється.

Альтернатива - event-driven project management.

## Чому calendar-driven status management запізнюється

Календарна модель має кілька прихованих припущень. Вона припускає, що проблеми накопичуються рівномірно. Що PM має час між зустрічами для обробки сигналів. Що інженери самі прийдуть до PM, якщо станеться щось важливе. Що інформаційна затримка у кілька днів несуттєва.

У сучасних інженерних, software, hardware/software або regulated проектах жодне з цих припущень не витримує реальності. Швидкість змін висока, кількість сигналів велика, залежності між компонентами і командами щільні. Затримка в декілька днів означає, що до моменту обговорення проблема вже мутувала, втягнула інші ризики або стала помітною замовнику.

Calendar-driven status management придатний для стратегічного огляду. Він не придатний для управління поточним станом.

## Які project events важливі

Event-driven підхід починається з простого: визначити, що саме у проектному середовищі є event-ом, на який система має реагувати.

Типові класи проектних подій:

- failed test, failed regression suite, failed build у release-critical контексті;
- defect, що піднімається до high/critical severity;
- requirement change, що зачіпає baseline-елементи;
- overdue critical task на критичному шляху;
- залежність, що стала blocked, або була розв'язана після блокування;
- risk, що перетнув threshold (probability, impact, exposure);
- supplier event: затримка, escalation, change in scope;
- milestone у небезпечній зоні (затримка, відсутність evidence);
- safety / cybersecurity finding;
- audit finding або corrective action із простроченим терміном;
- зміна evidence status (depreciated, missing);
- зміна compliance scope.

Це не повний список, але він показує діапазон. Усі ці події можуть і повинні запускати реакцію, а не чекати звіту.

## Event → analysis → action

Просто реагувати на event недостатньо. Треба перетворити сирий сигнал на структуровану управлінську інформацію. Гарна event-driven архітектура передбачає три фази для кожної значущої події.

Подія. Структурований запис: тип, джерело, контекст, час, прив'язка до артефактів.

Аналіз. Система автоматично визначає impact: які milestones під загрозою, які залежності зачеплені, які ризики змінюються, які stakeholder-и мають бути попереджені. Тут також підраховується severity і relevance.

Дія. Створюються пропоновані варіанти: escalate, mitigation action, resource request, scope change, CCB review, additional review. PM або відповідний owner ухвалює рішення.

Без цієї послідовності event-driven перетворюється на потік notifications. Із нею - на робочий інструмент.

## Role-based notifications і severity

Одна з найбільших небезпек event-driven моделі - alert fatigue. Якщо кожна подія генерує повідомлення для всіх, користувачі швидко перестають їх читати.

Лікування - це багатошарова логіка релевантності.

Severity. Подія має класифікуватися за серйозністю, і нижні рівні не йдуть людям у real-time, лише в periodic digest.

Role-based routing. PM проекту бачить інше, ніж PgM програми. Compliance officer бачить інше, ніж QA lead. Кожен отримує те, що стосується його зони відповідальності.

Suppression logic. Подія, що вже зареєстрована або переведена у роботу, не дублюється у нових повідомленнях. Подія, що очікувалася, відрізняється від несподіваної.

Aggregation. Кілька дрібних подій того ж типу зводяться в один сигнал з підсумком, а не у потік окремих повідомлень.

Channel selection. Critical event - real-time канал, medium - email або digest, низький - dashboard.

У сумі це означає, що користувач отримує менше повідомлень, але кожне з них варте уваги.

## AI Assistant як пояснювач подій

AI у event-driven моделі грає роль контекстного пояснювача. Сама по собі подія - це структурований запис із технічними полями. Менеджеру потрібне розуміння, що це означає.

Корисні застосування AI:

- перетворити сирий event у людську мову з контекстом проекту;
- пояснити, чому ця подія важлива, які залежності зачеплені, які milestones під загрозою;
- порівняти з історичними схожими подіями і показати, як вони розв'язувалися;
- запропонувати дії з посиланнями на конкретні артефакти;
- сформувати draft escalation message або change request package;
- сформувати summary для steering committee за період.

AI тут не приймає рішення. Він робить так, щоб менеджер міг прийняти його швидше і обґрунтованіше.

## Локальна проблема vs programme-level risk

Окреме питання, яке event-driven model дозволяє адресувати - різниця між локальною проблемою і ризиком програмного або портфельного рівня.

Багато подій є локальними. Інженер виправить дефект, тестувальник перезапустить регресію, supplier доставить компонент із невеликою затримкою. Це не треба ескалювати.

Інші події мають programme impact. Той самий supplier delay може зачепити критичний milestone програми, який залежить від кількох проектів. Той самий cybersecurity finding може поширюватися на сусідні системи.

Зріла event-driven архітектура дає механізми, що відрізняють ці випадки автоматично, на основі залежностей у digital thread, історичних patterns і structured rules. Завдяки цьому PgM отримує лише те, що справді стосується програмного рівня, а не дублювання всіх локальних сигналів.

## Event backlog і lifecycle

Подія не має зникати після notification. Для значущих events потрібен lifecycle: detected, classified, acknowledged, action proposed, action accepted, resolved, closed. Інакше event-driven система швидко перетвориться на стрічку повідомлень, де незрозуміло, що вже опрацьовано, а що загубилося.

Event record має містити source, timestamp, affected artifacts, severity, owner, related risk або issue, proposed action і decision. Якщо подія закрита, має бути видно, чому: false positive, already mitigated, accepted risk, fixed, escalated. Це створює audit trail і дозволяє вчитися на подіях.

Корисно також мати event backlog для PM/PgM. Не всі events потребують негайної реакції, але вони мають накопичуватися в керованому списку: open high-severity events, repeated medium events, aging unacknowledged events, events with no owner. Це вже не шум, а operational queue.

## Правила, які визначають реакцію

Event-driven модель не може триматися лише на інтуїції. Потрібні rules: які events critical, хто owner, який SLA на acknowledgement, коли ескалація, коли створюється mitigation task, коли потрібен CCB, коли event only logged. Ці правила мають бути прозорими, інакше команда не довірятиме notification logic.

Наприклад, failed test у non-release branch може піти в digest. Failed mandatory regression у release candidate - real-time alert для QA lead і PM. Supplier delay без critical path impact - monitor. Supplier delay, який зачіпає integration event, - programme-level escalation.

## З чого почати

Перший крок - вибрати 5-7 event types, які вже зараз болять. Не треба покривати все. Зазвичай це failed release tests, critical defects, baseline-impacting requirement changes, blocked dependencies, overdue mitigations, supplier delays, pending approvals. Для кожного event type треба визначити source, severity, owner, action policy і escalation rule.

Після цього можна додати AI explanation. Але foundation має бути rule-based і traceable, щоб менеджери не отримували красиві пояснення до погано визначених сигналів.

## Як не створити новий шум

Event-driven система має регулярно проходити hygiene review. Які alerts користувачі ігнорують? Які events завжди закриваються як false positive? Які notifications приходять не тим ролям? Які events повторюються, але не мають root-cause action? Без такого review event-driven підхід швидко перетвориться на ще один notification layer.

Корисна метрика - action rate. Якщо подія часто не приводить до action, можливо, вона не має бути real-time alert. Інша метрика - acknowledgement latency. Якщо high-severity events довго не визнаються owner-ом, routing або ownership налаштовані погано. Третя - repeat event rate: одна й та сама проблема повертається, бо система сигналізує, але організація не змінює причину.

AI може допомогти з цією hygiene: групувати повторювані events, пояснювати, які alerts стали noise, пропонувати зміни thresholds. Але governance має залишатися людським, бо alert policy впливає на поведінку всієї команди.

Ще один важливий принцип - кожен critical alert має мати очікувану дію. Якщо подія не може привести до decision, mitigation, escalation або conscious acceptance, вона не має бути critical. Це просте правило різко зменшує шум і змушує команду проектувати event taxonomy навколо управління, а не навколо цікавих технічних фактів.

Також варто мати quiet hours і escalation override. Не кожна подія заслуговує нічного повідомлення, але деякі release-critical або safety/security events мають проривати звичайний режим. Це має бути rule, а не імпровізація чергового менеджера.

Добра event policy заспокоює команду: люди знають, що важливе не загубиться, а неважливе не буде постійно смикати їх без причини.

Це і є різниця між alerting заради активності та управлінням подіями як частиною project governance.

## Мій висновок

Project management не може залежати лише від календаря у проектах, де події відбуваються щодня. Event-driven модель не скасовує weekly meetings і portfolio reviews, але переносить реактивну роботу з людської пам'яті у систему. Подія - аналіз - дія, із розумною severity, role-based routing і AI-поясненнями, перетворює management із "хто перший помітить і встигне сказати" на контрольовану реакцію на реальний стан. У результаті проблеми перестають накопичуватися мовчки і зустрічають організацію, поки вони ще керовані.

## Питання до читачів

- Які події у ваших проектах мали б одразу тригерити PM/PgM, а не чекати наступного статусу?
- Де у вас notifications вже перетворилися на шум, який ніхто не читає?
- Як ви зараз відрізняєте локальну проблему від programme-level risk?
- Що зміниться у вашій роботі, якщо event-driven логіка піде разом із digital thread?
