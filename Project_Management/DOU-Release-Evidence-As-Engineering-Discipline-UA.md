# Release evidence: чому реліз має доводити себе сам

У багатьох командах реліз досі доводиться словами. Є release notes, changelog, status meeting, кілька посилань на pipeline, скриншот тестів, коментар QA, короткий summary від PM. Усе це може бути чесним. Але коли продукт складний, крос-командний або regulated, слів уже недостатньо.

Реліз має доводити себе сам.

Під цим я маю на увазі release evidence - структурований набір доказів, який показує, що саме було зібрано, з яких версій, якими тестами перевірено, з якими протоколами сумісно, які gaps лишилися, які security/redaction checks пройдені, куди evidence опубліковано і чи можна цю публікацію відтворити.

Це не просто ще один файл у релізі. Це дисципліна.

## Release notes не є release evidence

Release notes пояснюють, що змінилося. Release evidence доводить, що це перевірено.

Release notes можуть сказати: «додано сумісність із новою версією протоколу». Release evidence має сказати:

- які версії компонентів були перевірені разом;
- яка версія спільного контракту використовувалася;
- чи були локальні заміни залежностей;
- який protocol version активний;
- який fixture corpus прогнано;
- яка smoke command виконана;
- який exit code;
- який trace/correlation sample підтверджує шлях;
- чи були residual gaps;
- чи був repository dirty;
- чи є evidence artifact, який можна відкрити пізніше.

Це зовсім інший рівень. Release notes працюють для комунікації. Release evidence працює для довіри.

## Чому це важливо у крос-продуктових системах

У моноліті release evidence може бути простішим. Є один repository, одна pipeline, один artifact, один set of tests. У продуктовій сім'ї все складніше.

Один компонент випускає protocol constant. Другий використовує його через shared SDK. Третій очікує payload field. Четвертий публікує Wiki report. П'ятий збирає evidence з GitLab. Кожен окремо може бути green. Але programme compatibility залежить від того, чи вони перевірені разом.

Саме тут release evidence стає programme-level об'єктом.

Якщо один компонент заявляє, що інтеграція з іншим працює, доказ має показати, на яких версіях це перевірено. Якщо можливість позначена як чернетка або відкладена, реліз не має називати її активною сумісністю. Якщо інструмент автоматизації публікує звіт, доказ має містити попередній перегляд, перевірку цілі та зворотне читання. Якщо система перейшла на резервний режим, реліз не має називати це апаратним прискоренням.

Без такого evidence сумісність стає narrative. З evidence вона стає перевірюваним фактом.

## Інструмент публікації не має підміняти управлінське рішення

Інструмент автоматизації може збирати факти з конвеєра, формувати попередній перегляд публікації, перевіряти ціль, виконувати знеособлення чутливих даних, зчитувати опублікований вміст і порівнювати хеш. Це корисна, але вузька відповідальність.

Важливо не плутати межі: інструмент публікації не є рушієм політик для всієї програми. Він не вирішує, чи компоненти сумісні, чи доказ достатній, або чи згенерований AI підсумок схвалено для публікації.

Правильна модель така:

- компоненти продукують власні докази;
- спільна документація визначає схему й словник;
- інструмент автоматизації перевіряє, публікує та фіксує доказ публікації;
- відповідальний за реліз схвалює чутливі твердження й твердження про деградований режим.

Тобто інструмент автоматизації робить докази придатними до публікації та повторюваними. Але самі докази мають народжуватися в компонентах.

## Як застосувати підхід до інших компонентів

Так, але координовано. Кожен компонент має продукувати докази у своїй зоні відповідальності: стан проєкту, сумісність інтеграції, перевірки, діагностику виконання, якість отримання знань або дотримання контракту інтеграції. Спільна бібліотека контрактів має підтверджувати сумісність інтерфейсів і контрактні тести.

Інструмент публікації відповідає лише за перевірку цілі, попередній перегляд, знеособлення, результат запису, зворотне читання та хеш вмісту. Це не означає, що всі компоненти мають копіювати один інструмент: кожен формує власні докази, а публікаційний процес лише збирає й перевіряє їх.

## Мінімальний release evidence record

Мінімальний evidence record для крос-продуктового релізу має відповідати на кілька простих питань.

Що саме релізимо? Component name, version, git SHA, tag або branch, dirty state.

З ким сумісно? Protocol version, shared SDK version, required component versions, compatibility window.

Що перевірили? Fixture corpus, smoke command, test summary, coverage або relevant CI job IDs, exit codes.

Який runtime/context був активний? CPU/GPU/NPU, local LLM runtime, provider, fallback, model/runtime evidence where relevant.

Який trace sample доводить шлях? Request ID, correlation ID, traceparent, artifact references.

Які security boundaries зачеплено? Classification, movement policy, redaction status, generated content review state.

Що лишилося невирішеним? Residual gaps, skipped checks, degraded components, draft/deferred capabilities.

Де це опубліковано? Target, URL/path, content hash, readback status.

Якщо evidence record не може відповісти на ці питання, release compatibility claim слабкий.

## Dry-run і readback як культура безпечної автоматизації

Мені дуже подобається пара dry-run/readback. Це проста, але сильна ідея.

Dry-run відповідає на питання: "що саме ми збираємося зробити, куди, з яким content hash, після яких checks, але без мутації target system".

Readback відповідає на питання: "чи справді те, що ми опублікували, лежить там, де очікували, і чи збігається content".

Це особливо важливо для GitLab Wiki, issue/MR comments, release pages, generated Markdown reports. Automation, яка пише у GitLab без dry-run і readback, може бути корисною, але не audit-grade. Вона може помилитися target-ом, опублікувати не той content, втратити форматування, пропустити redaction gap або залишити команду без доказу, що publication успішна.

Dry-run/readback змінює характер автоматизації. Вона перестає бути "скрипт щось написав" і стає "workflow має доказ наміру і доказ результату".

## Redaction як частина release evidence

Release evidence часто містить sensitive context: internal URLs, logs, job traces, snippets, generated summaries, provider diagnostics, local paths, tokens by accident. Якщо evidence автоматично публікується у GitLab, redaction не може бути optional afterthought.

Добрий release evidence workflow має явно фіксувати redaction status:

- passed;
- blocked;
- manual review required;
- not required.

Також має бути видно redaction profile. Одне діло - публікувати open-source release notes. Інше - публікувати cross-product engineering evidence з logs, trace IDs, generated summaries і internal component names.

AI тут додає ще один ризик. Generated content може випадково переказати restricted details або змішати source contexts. Тому AI summaries у release evidence мають бути або human-approved, або marked as manual review required.

## Release evidence має бути machine-readable

Markdown report корисний для людей. Але source of truth має бути machine-readable JSON або інший стабільний structured format. Інакше evidence неможливо автоматично перевіряти.

Machine-readable evidence дозволяє:

- перевірити schema;
- порівняти component versions;
- побачити dirty repos;
- знайти skipped checks;
- gate-ити реліз за required fields;
- будувати programme compatibility matrix;
- публікувати human-readable report з того самого source;
- робити readback hash;
- накопичувати history для trend analysis.

Якщо CLI має `--json`, stdout має бути чистим JSON. Без banner, progress text, color codes, warning messages. Diagnostics можуть іти у stderr. Це дрібниця, яка відрізняє automation-friendly CLI від красивого human CLI.

## Release evidence як operational memory

Через пів року після релізу команда часто не пам'ятає деталей. Яка версія спільного контракту була в компоненті? Чи була локальна заміна залежності? Чи проганяли набір перевірок? Чи можливість була активною або лише чернеткою? Чи апаратне прискорення справді працювало, чи система перейшла на резервний режим? Чому реліз схвалили попри залишкову прогалину?

Release evidence відповідає на ці питання без розкопок.

Це важливо не лише для audit. Це важливо для debugging, incident review, migration planning, customer support, compatibility matrix, onboarding і future architecture decisions.

У цьому сенсі release evidence - це пам'ять релізу.

## Зв'язок із програмною аналітикою

Release evidence і програмна аналітика дуже близькі. Programme analysis формулює, що має бути правдою між проектами. Release evidence доводить, що це було правдою для конкретного релізу.

Наприклад, аналітика може вимагати, щоб поля трасування й кореляції однаково проходили через усі компоненти. Доказ релізу має показати відповідний приклад трасування.

Analysis каже: "classification and movement policy metadata must be preserved without raw restricted content". Release evidence має показати metadata sample і redaction status.

Analysis каже: "local LLM runtime claims need model hash, quantization, backend and fallback". Release evidence має містити ці fields, якщо runtime був частиною compatibility claim.

Без release evidence programme analysis залишається policy. З release evidence вона стає перевірюваною дисципліною.

## З чого почати команді

Не треба одразу будувати ідеальний release evidence platform. Почати можна з малого.

Перший крок - визначити canonical evidence schema для одного release flow.

Другий - змусити кожен компонент виводити свою частину evidence у machine-readable format.

Третій - додати smoke command, який створює evidence artifact.

Четвертий - додати redaction check перед publication.

П'ятий - додати dry-run для GitLab publication.

Шостий - додати readback hash після publication.

Сьомий - зробити skipped checks явними. Не "ми не встигли", а structured skipped check з owner, reason і risk.

Це вже дає величезну різницю.

## Чи є тут тема для статті

Так, і дуже сильна. Release evidence — це міст між релізною інженерією, відповідністю вимогам, DevOps, проєктним менеджментом і урядуванням програми. Майже кожна інженерна організація має проблему: релізи є, але докази релізів живуть у журналах CI, сторінках документації, коментарях, таблицях, пошті й пам'яті людей.

Стаття на цю тему може бути навіть практичнішою за багато абстрактних розмов про DevOps maturity. Бо вона ставить просте питання: якщо завтра треба довести, що реліз був сумісний, перевірений і безпечно опублікований, де цей доказ?

Якщо відповідь - "пошукаємо в pipeline і спитаємо QA", release evidence discipline ще не побудована.

## Мій висновок

Так, підхід із доказами релізу варто застосовувати до інших компонентів, але з правильними межами відповідальності. Компоненти мають продукувати власні докази, а інструмент публікації — безпечно збирати, перевіряти, публікувати й зчитувати назад зведений доказ. Спільна документація має тримати схему, словник і правила перенесення результатів у роботу.

Release evidence - це не бюрократія. Це спосіб зробити реліз доказовим. Команда менше сперечається про те, що було перевірено, і більше працює з фактами. А коли з'являється audit, incident або compatibility question, відповідь уже існує.

Реліз, який не може показати evidence, просить довіри. Реліз, який має evidence, доводить себе сам.

## Питання до читачів

- Чи можете ви зараз відкрити evidence record для релізу піврічної давності?
- Чи видно там component versions, commit SHAs, fixture set, smoke command і residual gaps?
- Чи має ваша release automation dry-run і readback, чи просто пише у target system?
- Який evidence у вашій організації найважче відтворити після релізу?

## Посилання

- [NIST SP 800-218 — Secure Software Development Framework](https://csrc.nist.gov/pubs/sp/800/218/final)
- [ISO 21502:2020 — guidance on project management](https://www.iso.org/standard/74947.html)
