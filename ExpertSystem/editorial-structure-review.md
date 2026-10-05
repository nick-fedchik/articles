# Редакторська карта структури книги

Дата перегляду: 2026-10-05. Об'єкт: 38 основних глав. Додатки враховано як прикладні маршрути, але не як основні глави. Дослідницькі нотатки не входять до змісту книги.

Це структурна рецензія: зіставлено назви всіх розділів і підрозділів поза блоками коду, вступні цілі та висновки. Програмні приклади, математичні доведення й бібліографічні твердження не проходили нового технічного випробування. Редакторська оцінка цілісності не є оцінкою наукової правильності.

## Таксономія заголовків

До редагування основні глави містили 992 заголовки всіх рівнів, включно з назвами глав і довідковим апаратом. Заголовки описують різні ознаки: предмет знання, задачу, метод, інструмент, контекст застосування або етап перевірки. Тому перелік заголовків не утворював одного дерева взаємовиключних наукових дисциплін.

Запропоновано дві незалежні осі. **Тематична вісь** визначає основну частину глави за її головним питанням. **Композиційна вісь** визначає роль розділу в аргументі: проблема; об'єкт і модель; метод і процедура; реалізація; перевірка; висновок із межами. Назва країни, галузі чи продукту є контекстом, а не рівнем тематичної ієрархії.

| Тематичний клас | Основне питання | Послідовність глав |
|---|---|---|
| I. Основи й контракт знання | Коли потрібна експертна система і що вона може стверджувати? | 1, 2, 3, 4, 5 |
| II. Подання та зберігання | Як описати й зберегти знання разом із залежностями? | 6, 7, 8, 9, 32 |
| III. Здобуття й оцінювання входу | Як отримати придатні кандидати з документів, людей і спостережень? | 10, 11, 12, 13, 14, 15, 37 |
| IV. Архітектура рішення | Як перейти від запиту до висновку, пояснення й дозволеної дії? | 16, 19, 31, 20, 21 |
| V. Перевірка й кероване вдосконалення | Як перевіряти правила, діагностувати, змінювати знання й обґрунтовувати безпеку? | 23, 36, 24, 25, 26, 27, 30 |
| VI. Гібридні відповіді й прогалини | Як поєднати мовну модель, строгий результат і непідтверджену гіпотезу? | 28, 29, 34, 38 |
| VII. Реалізація й експлуатація | Як обрати засоби, розгорнути виконання й обмінюватися знаннями? | 17, 18, 22, 35, 33 |

Номери глав залишаються ідентифікаторами. Тематичний порядок задано змістом і переходами, а не перейменуванням файлів. Довідкові глави 6 і 7 та великі оглядові глави не потрібно читати повністю перед першим програмним прикладом.

## Цілісність кожної глави

Позначка **цілісна** означає, що основні розділи працюють на одне питання, а висновок повертається до питання. **Оглядова** означає навмисне порівняння методів за одним критерієм, не випадкове змішування тем. **Потребує уточнення** означає проблему меж, обіцянок або композиції, яку нова анотація не усуває всередині глави.

| Глава | Головна тема й зв'язок мети з висновком | Редакторська оцінка та рішення |
|---|---|---|
| [1](ch01-introduction-to-expert-systems.md) | Потреба в перевірюваній відповіді; висновок повертається до пакета обґрунтування | Цілісна вступна; галузі й перший приклад залишаються орієнтирами, не окремими темами |
| [2](ch02-epistemology-of-machine-knowledge.md) | Право називати твердження знанням; сім перевірок відповідають заявленій меті | Потребує композиційного уточнення: матеріал про результати випробувань розміщений після висновку; його варто розташувати перед висновком або відокремити як приклад |
| [3](ch03-beyond-reference-information-systems.md) | Межа між довідкою й експертним рішенням на трьох рівнях | Цілісна порівняльна; опис експлуатації є перевіркою цієї межі |
| [4](ch04-evolution-from-bayes-to-evidence-ai.md) | Історія методів і вибір за типом невизначеності | Цілісна оглядова; історія й одна ідея методу належать тут, повний розрахунок у главі 6 |
| [5](ch05-triad-of-trust-and-corporate-memory.md) | Збереження підстав конкретного організаційного рішення | Цілісна; три складники об'єднані одним прикладом корпоративної пам'яті |
| [6](ch06-applied-mathematics-for-expert-systems.md) | Вибір математичної моделі за видом запитання | Цілісна довідкова; численні методи є навмисною картою, потрібні вибіркові маршрути |
| [7](ch07-knowledge-base-typology.md) | Вибір типу бази знань за семантикою результату | Цілісна оглядова, але велика; фізичне розбиття має спиратися на семантику, а реалізація пакетів належить главі 32 |
| [8](ch08-engineering-artifacts-as-data.md) | Перетворення артефакта на версійований об'єкт із походженням | Цілісна з широким практичним охопленням; апаратний вибір є допоміжним, не головною темою |
| [9](ch09-engineering-knowledge-graph-traceability.md) | Типізовані зв'язки для покриття та впливу змін | Цілісна; прогнозовані ребра відокремлено від підтверджених |
| [10](ch10-knowledge-acquisition-systems.md) | Контрольований конвеєр постачання знань | Цілісна архітектурна; назва не повинна обіцяти абсолютної відсутності втрат чи витоку |
| [11](ch11-knowledge-elicitation-from-experts.md) | Перетворення досвіду людини на перевірюване правило | Цілісна; інтерв'ю й цифровий слід є різними джерелами кандидатів |
| [12](ch12-linguistic-analysis-and-local-models.md) | Мовний аналіз зі збереженням зв'язку з джерелом | Цілісна; локальні моделі й токенізація підпорядковані контракту мовної підсистеми |
| [13](ch13-language-variability-vs-determinism.md) | Перевірка відповідності різних формулювань одній логічній формі | Цілісна; обмежена граматика не підміняє перевірки змісту |
| [14](ch14-requirements-detection-and-formalization.md) | Виявлення модальностей і формалізація умов вимоги | Цілісна; зберегти межу між кандидатом вимоги та схваленою формулою |
| [15](ch15-knowledge-extraction-and-kb-construction.md) | Побудова фактів, граматик і автоматів із простежуваних джерел | Цілісна розгорнута; спеціалізовані екстрактори є підзадачами одного конвеєра |
| [16](ch16-expert-systems-architecture.md) | Розподіл відповідальності за допустимість рішення | Цілісна; логічна архітектура відокремлена від вибору бібліотек |
| [17](ch17-implementation-stack.md) | Вибір інструмента за семантикою й поведінкою за відмов | Цілісна оглядова; перенесено до реалізації, не до теорії виведення |
| [18](ch18-execution-infrastructure.md) | Обґрунтування апаратури, затримки та автономності | Цілісна; перспективні процесори є окремо позначеними гіпотезами |
| [19](ch19-from-question-to-evidence.md) | Перевірка твердження після пошуку й побайтової прив'язки | Цілісна; джерело й семантична підтримка є різними перевірками |
| [20](ch20-explanation-engine.md) | Пояснення з фактичної траси рішення | Цілісна; контрфакти, відмови й доступ є різними видами пояснювального запиту |
| [21](ch21-from-recommendation-to-action.md) | Допуск і виконання дії з перевіреною післяумовою | Цілісна; правильність рекомендації не надає повноважень |
| [22](ch22-cybernetics-edge-to-backend.md) | Зворотний зв'язок за шуму, запізнення й втрати зв'язку | Цілісна широка; кібернетичні аналогії не повинні підміняти перевірку фізичної моделі |
| [23](ch23-knowledge-base-verification.md) | Сукупність незалежних перевірок бази правил | Цілісна; загальна методологія, на відміну від конкретної піраміди глави 36 |
| [24](ch24-system-diagnosis.md) | Сумісні пояснення несправності та вибір розрізнювальної перевірки | Цілісна; діагностика зовнішнього об'єкта відокремлена від перевірки власної бази правил |
| [25](ch25-how-expert-systems-learn.md) | Навчання як керований випуск нової версії | Цілісна, але найбільша за обсягом; іспит, життєвий цикл і доналаштування моделі потребують окремих внутрішніх маршрутів, не автоматичного розрізання |
| [26](ch26-continual-learning.md) | Незалежна перевірка кандидатів із власного досвіду експлуатації | Цілісна; основна відмінність від глави 25 полягає у зсуві журналу, дрейфі й забуванні |
| [27](ch27-safety-case-gsn-synthesis.md) | Синтез і перевірка структури аргументу безпеки | Цілісна; назва має говорити про обґрунтування, не автоматичну сертифікацію |
| [28](ch28-dual-mode-expert-systems.md) | Розділення строгого результату та дорадчої гіпотези | Цілісна; одна відповідь може мати два статуси, які не змішуються |
| [29](ch29-neuro-symbolic-architecture.md) | Інтеграція мовної моделі через перевірку кандидатів | Цілісна; вилучити застаріле твердження, що глава завершує всю книгу |
| [30](ch30-safety-cybersecurity-co-engineering.md) | Узгодження вимог функціональної безпеки й кібербезпеки | Цілісна спеціалізована; не повторює загальний синтез аргументу глави 27 |
| [31](ch31-syllogistic-reasoning-and-relation-lattices.md) | Застосування норм за ієрархією, винятками й часом чинності | Цілісна; навчальний однокроковий приклад не варто називати повним багатоходовим рушієм |
| [32](ch32-high-performance-knowledge-packs-mmap-and-harvesting.md) | Незмінний пакет знань і його фізичні індекси | Цілісна розгорнута; формат, перевірка читача й шардування є властивостями одного артефакта |
| [33](ch33-inter-system-knowledge-exchange-and-model-teaching.md) | Передача знань без втрати семантики, доступу й походження | Цілісна, але межова; навчання чужих моделей є одним способом використання експорту, не новою головною темою |
| [34](ch34-deterministic-relational-analysis-abduction-and-socratic-dialogue.md) | Перетворення прогалини на гіпотезу й запит уточнення | Потребує звуження: індуктивний пошук правил і міжпредметне узагальнення мають лишатися засобами роботи з прогалиною; уникати повторення глави 31 |
| [35](ch35-self-organizing-expert-systems-synergetics-and-npu-runtime.md) | Реактивне виконання та перегляд підстав за подіями | Потребує уточнення: фізична синергетика, еволюція онтологій і прискорювачі розширюють заявлену тему; приклад виконання не доводить автономної еволюції всієї онтології |
| [36](ch36-knowledge-testing-pyramid-and-variational-calibration.md) | Авторська піраміда перевірки правил і стійкості відповідей | Тематично цілісна; висновки про сертифікацію та абсолютні гарантії потребують окремого змістового перегляду |
| [37](ch37-input-information-assessment-and-algorithmic-skepticism.md) | Придатність вхідного повідомлення як підстави висновку | Цілісна широка; національні традиції та платформи є прикладами застосування спільного контракту, не каталогом дев'яти універсальних алгоритмів |
| [38](ch38-curing-machine-hallucinations-and-knowledge-deficits.md) | Контроль непідтверджених відповідей і керування прогалинами | Потребує уточнення: точна цитата не доводить істинності чи правильного тлумачення; назва й анотація не повинні обіцяти універсального лікування |

## Межі між сусідніми темами

| Пара або група | Де зберігати розгорнутий виклад | Що лишати в сусідній главі |
|---|---|---|
| 4 / 6 / 7 | 4: історія; 6: математична операція; 7: семантика подання | Коротку ідею, приклад межі й посилання на основний виклад |
| 8 / 10 / 15 / 37 | 8: носій; 10: керування постачанням; 15: змістове вилучення; 37: якість свідчення | Не повторювати весь конвеєр у кожній главі |
| 13 / 19 / 31 / 34 | 13: форма запиту; 19: підтримка твердження; 31: застосування норми; 34: прогалина | Розрізняти пошук, прив'язку, виведення й абдукцію |
| 23 / 36 / 25 / 26 | 23: загальні перевірки; 36: рівні тестування; 25: випуск; 26: досвід експлуатації | Не змішувати результат тесту, критерій допуску й навчальний сигнал |
| 22 / 35 | 22: фізичний цикл керування; 35: реактивний перегляд знань | Апаратні виміри винести до 18; аналогію самоорганізації позначати як аналогію |
| 28 / 29 / 38 | 28: статус відповіді; 29: архітектура інтеграції; 38: клас помилок і контроль | Не подавати побайтову прив'язку як універсальний доказ істинності |
| 27 / 30 | 27: структура аргументу; 30: спільний предметний ризик | Не ототожнювати машинну перевірку з регуляторною сертифікацією |

## Виконані та відкладені зміни

Виконано тематичне групування, уточнення назв, оновлення описів частин і наскрізних переходів. Службові заголовки висновку, словника, абревіатур і джерел приведено до спільної форми. У главах код, формули й нумерацію бібліографічних записів збережено. В оглядових анотаціях звужено абсолютні обіцянки, яких локальний приклад не доводить.

Масове переписування всіх розділів не виконувалося. Для глав 2, 34, 35, 36 і 38 залишено конкретні редакторські зауваження вище. До окремої змістової перевірки не слід вважати твердження про стовідсоткове усунення помилок, мілісекундну безпеку чи автономну еволюцію онтологій доведеними лише на підставі назви глави або зеленого демонстраційного тесту.

## Покажчик усіх заголовків

Нижче наведено поточні заголовки кожної глави в їхній ієрархії. Рівні позначено як H1–H6; заголовки всередині блоків коду не враховано. Покажчик дозволяє зіставити фактичну композицію з оцінкою глави, не підміняючи тематичну класифікацію автоматичним добором ключових слів.

<details>
<summary>Глава 1. Вступ до експертних систем: від хаосу до керованих знань</summary>

| Рівень | Заголовок |
|---|---|
| H1 | [Глава 1. Вступ до експертних систем: від хаосу до керованих знань](ch01-introduction-to-expert-systems.md#L1) |
| H2 | [Одне запитання, три відповіді](ch01-introduction-to-expert-systems.md#L15) |
| H2 | [Чому ця тема повернулася саме зараз](ch01-introduction-to-expert-systems.md#L31) |
| H2 | [Хіба експертні системи не померли у 1980-х?](ch01-introduction-to-expert-systems.md#L45) |
| H2 | [Хто такий експерт і що з його знань можна передати програмі](ch01-introduction-to-expert-systems.md#L60) |
| H2 | [Що таке експертна система](ch01-introduction-to-expert-systems.md#L78) |
| H3 | [Що таке база знань](ch01-introduction-to-expert-systems.md#L128) |
| H3 | [Що таке машина виведення](ch01-introduction-to-expert-systems.md#L140) |
| H3 | [Перший виконуваний приклад: перевірка зміни інтерфейсу](ch01-introduction-to-expert-systems.md#L157) |
| H3 | [Які бувають експертні системи](ch01-introduction-to-expert-systems.md#L295) |
| H2 | [Вхід, міркування і вихід](ch01-introduction-to-expert-systems.md#L312) |
| H3 | [Вхід](ch01-introduction-to-expert-systems.md#L367) |
| H3 | [Міркування](ch01-introduction-to-expert-systems.md#L388) |
| H3 | [Вихід](ch01-introduction-to-expert-systems.md#L399) |
| H2 | [Експертна відповідь і пакет обґрунтування](ch01-introduction-to-expert-systems.md#L409) |
| H2 | [Простежуваність: чому один таймаут пов'язаний із шістьма артефактами](ch01-introduction-to-expert-systems.md#L434) |
| H2 | [Документ не дорівнює знанню](ch01-introduction-to-expert-systems.md#L472) |
| H3 | [Мінімальний об'єкт знань](ch01-introduction-to-expert-systems.md#L517) |
| H2 | [Як знання потрапляють в експертну систему і повертаються до людини](ch01-introduction-to-expert-systems.md#L572) |
| H2 | [Експертна система серед пошуку, чатботів і RAG](ch01-introduction-to-expert-systems.md#L620) |
| H3 | [Де корисні машинне навчання й мовні моделі](ch01-introduction-to-expert-systems.md#L643) |
| H2 | [П'ять галузей, до яких повертається книга](ch01-introduction-to-expert-systems.md#L655) |
| H2 | [Межі: коли експертна система не допоможе](ch01-introduction-to-expert-systems.md#L683) |
| H2 | [Як довести, що експертна система корисна](ch01-introduction-to-expert-systems.md#L706) |
| H3 | [Чи бачить експертна система всі потрібні зв'язки](ch01-introduction-to-expert-systems.md#L716) |
| H3 | [Чи спираються відповіді на докази](ch01-introduction-to-expert-systems.md#L733) |
| H3 | [Чи не ховається експертна система за відмовою](ch01-introduction-to-expert-systems.md#L750) |
| H3 | [Що вимірюють додатково](ch01-introduction-to-expert-systems.md#L769) |
| H2 | [З чого почати в команді](ch01-introduction-to-expert-systems.md#L781) |
| H2 | [Як читати цю книгу далі](ch01-introduction-to-expert-systems.md#L794) |
| H2 | [Висновки](ch01-introduction-to-expert-systems.md#L818) |
| H2 | [Запитання до читачів](ch01-introduction-to-expert-systems.md#L827) |
| H2 | [Словник](ch01-introduction-to-expert-systems.md#L834) |
| H2 | [Абревіатури](ch01-introduction-to-expert-systems.md#L857) |
| H2 | [Джерела](ch01-introduction-to-expert-systems.md#L875) |

</details>

<details>
<summary>Глава 2. Філософія для інженера: що машина має право називати знанням</summary>

| Рівень | Заголовок |
|---|---|
| H1 | [Глава 2. Філософія для інженера: що машина має право називати знанням](ch02-epistemology-of-machine-knowledge.md#L1) |
| H2 | [Практичний вхід: контракт відповіді](ch02-epistemology-of-machine-knowledge.md#L15) |
| H2 | [Одне запитання і три знайдені документи](ch02-epistemology-of-machine-knowledge.md#L21) |
| H2 | [Що таке знання: від Платона до Геттієра](ch02-epistemology-of-machine-knowledge.md#L84) |
| H2 | [Робоче визначення знання експертної системи](ch02-epistemology-of-machine-knowledge.md#L98) |
| H2 | [Епістемологія: твердження, доказ і походження](ch02-epistemology-of-machine-knowledge.md#L126) |
| H2 | [Онтологія: однакова назва не означає той самий об'єкт](ch02-epistemology-of-machine-knowledge.md#L198) |
| H3 | [Два часи замість одного поля `updated_at`](ch02-epistemology-of-machine-knowledge.md#L318) |
| H3 | [Чи застосовне твердження до запиту](ch02-epistemology-of-machine-knowledge.md#L352) |
| H2 | [Логіка: висновок має називати спосіб, яким його отримано](ch02-epistemology-of-machine-knowledge.md#L479) |
| H3 | [Відкритий і замкнений світ: небезпечне «не знайдено»](ch02-epistemology-of-machine-knowledge.md#L496) |
| H3 | [Суперечність не повинна породжувати довільну відповідь](ch02-epistemology-of-machine-knowledge.md#L510) |
| H2 | [Філософія мови: що саме запитує користувач](ch02-epistemology-of-machine-knowledge.md#L588) |
| H2 | [Герменевтика: цитата без контексту не є доказом](ch02-epistemology-of-machine-knowledge.md#L617) |
| H2 | [Філософія науки: вимірювання, прогноз і норма не взаємозамінні](ch02-epistemology-of-machine-knowledge.md#L681) |
| H2 | [Соціальна епістемологія: авторитет, доступ і відповідальність](ch02-epistemology-of-machine-knowledge.md#L707) |
| H2 | [Умова відповіді: коли експертна система має право стверджувати](ch02-epistemology-of-machine-knowledge.md#L863) |
| H2 | [Як перевірити експертну систему: види відмов і швидкий аудит](ch02-epistemology-of-machine-knowledge.md#L1007) |
| H2 | [Висновки](ch02-epistemology-of-machine-knowledge.md#L1050) |
| H2 | [Тріада заземленої фактуальності та емпіричне калібрування: результати W3C та RFC-1000](ch02-epistemology-of-machine-knowledge.md#L1069) |
| H3 | [Емпіричний протокол калібрування (W3C-150 та RFC-1000)](ch02-epistemology-of-machine-knowledge.md#L1080) |
| H2 | [Запитання до читачів](ch02-epistemology-of-machine-knowledge.md#L1096) |
| H2 | [Словник](ch02-epistemology-of-machine-knowledge.md#L1102) |
| H2 | [Абревіатури](ch02-epistemology-of-machine-knowledge.md#L1136) |
| H2 | [Джерела](ch02-epistemology-of-machine-knowledge.md#L1160) |

</details>

<details>
<summary>Глава 3. Чим експертна система відрізняється від інформаційно-довідкової системи</summary>

| Рівень | Заголовок |
|---|---|
| H1 | [Глава 3. Чим експертна система відрізняється від інформаційно-довідкової системи](ch03-beyond-reference-information-systems.md#L1) |
| H2 | [Чотири класи інформаційних систем](ch03-beyond-reference-information-systems.md#L15) |
| H2 | [Рівень запитання: про корпус документів чи про конкретний випадок](ch03-beyond-reference-information-systems.md#L34) |
| H2 | [Де місце великої мовної моделі](ch03-beyond-reference-information-systems.md#L48) |
| H2 | [Рівень документа: довідка і експертний висновок](ch03-beyond-reference-information-systems.md#L69) |
| H3 | [Коли текст стає документом](ch03-beyond-reference-information-systems.md#L73) |
| H3 | [Чим експертний висновок відрізняється від довідки](ch03-beyond-reference-information-systems.md#L119) |
| H3 | [Приклад: рішення про допуск прошивки до випробувань](ch03-beyond-reference-information-systems.md#L167) |
| H2 | [Рівень архітектури: пошуковий конвеєр і машина виведення](ch03-beyond-reference-information-systems.md#L322) |
| H3 | [Як влаштована генерація з пошуком](ch03-beyond-reference-information-systems.md#L326) |
| H3 | [Як влаштована експертна система](ch03-beyond-reference-information-systems.md#L359) |
| H2 | [Чого бракує кожному класу: порівняння за епістемічним контрактом](ch03-beyond-reference-information-systems.md#L401) |
| H2 | [Експертна система в експлуатації: навчання, впевненість і дії](ch03-beyond-reference-information-systems.md#L426) |
| H3 | [Навчання лише через перевірену нову версію бази знань](ch03-beyond-reference-information-systems.md#L430) |
| H3 | [Чи можна довіряти числовій впевненості](ch03-beyond-reference-information-systems.md#L442) |
| H3 | [Від дорадчого висновку до автоматичної дії](ch03-beyond-reference-information-systems.md#L452) |
| H2 | [Висновки](ch03-beyond-reference-information-systems.md#L464) |
| H2 | [Запитання до читачів](ch03-beyond-reference-information-systems.md#L498) |
| H2 | [Словник](ch03-beyond-reference-information-systems.md#L504) |
| H2 | [Абревіатури](ch03-beyond-reference-information-systems.md#L544) |
| H2 | [Джерела](ch03-beyond-reference-information-systems.md#L566) |

</details>

<details>
<summary>Глава 4. Еволюція експертних систем: від теореми Байєса до доказових рішень ШІ</summary>

| Рівень | Заголовок |
|---|---|
| H1 | [Глава 4. Еволюція експертних систем: від теореми Байєса до доказових рішень ШІ](ch04-evolution-from-bayes-to-evidence-ai.md#L1) |
| H2 | [Історія методу й інструкція до обчислення](ch04-evolution-from-bayes-to-evidence-ai.md#L15) |
| H2 | [Тип невизначеності визначає метод](ch04-evolution-from-bayes-to-evidence-ai.md#L19) |
| H2 | [Чотири хвилі: чому експертні системи злетіли, впали й повернулися](ch04-evolution-from-bayes-to-evidence-ai.md#L81) |
| H3 | [Перша хвиля: знання важливіші за універсальний алгоритм](ch04-evolution-from-bayes-to-evidence-ai.md#L113) |
| H3 | [Друга хвиля: промисловий бум](ch04-evolution-from-bayes-to-evidence-ai.md#L123) |
| H3 | [Зима ШІ: чотири інженерні причини](ch04-evolution-from-bayes-to-evidence-ai.md#L129) |
| H3 | [Третя хвиля: рушії правил і семантичний веб](ch04-evolution-from-bayes-to-evidence-ai.md#L140) |
| H3 | [Четверта хвиля: мовні моделі разом із символьним виведенням](ch04-evolution-from-bayes-to-evidence-ai.md#L146) |
| H2 | [Логіка й продукційні правила](ch04-evolution-from-bayes-to-evidence-ai.md#L152) |
| H3 | [Дві культури раннього ШІ: LISP і PROLOG](ch04-evolution-from-bayes-to-evidence-ai.md#L237) |
| H2 | [Теорема Байєса: як оновлювати впевненість після свідчення](ch04-evolution-from-bayes-to-evidence-ai.md#L319) |
| H2 | [Коли повної ймовірнісної моделі немає](ch04-evolution-from-bayes-to-evidence-ai.md#L472) |
| H3 | [Коефіцієнти впевненості MYCIN](ch04-evolution-from-bayes-to-evidence-ai.md#L476) |
| H3 | [Нечітка логіка](ch04-evolution-from-bayes-to-evidence-ai.md#L505) |
| H3 | [Теорія Демпстера–Шафера](ch04-evolution-from-bayes-to-evidence-ai.md#L534) |
| H2 | [Зв'язки, джерела й досвід](ch04-evolution-from-bayes-to-evidence-ai.md#L581) |
| H3 | [Від семантичних мереж до графів знань](ch04-evolution-from-bayes-to-evidence-ai.md#L585) |
| H3 | [Байєсівські мережі](ch04-evolution-from-bayes-to-evidence-ai.md#L597) |
| H3 | [Пошук інформації: від зважування слів до векторів](ch04-evolution-from-bayes-to-evidence-ai.md#L615) |
| H3 | [Міркування за прецедентами](ch04-evolution-from-bayes-to-evidence-ai.md#L635) |
| H2 | [Порівняння альтернатив і точні розрахунки](ch04-evolution-from-bayes-to-evidence-ai.md#L643) |
| H3 | [Зважене оцінювання альтернатив](ch04-evolution-from-bayes-to-evidence-ai.md#L647) |
| H3 | [Детерміновані обчислювальні ядра](ch04-evolution-from-bayes-to-evidence-ai.md#L680) |
| H2 | [Що успадкувала доказова експертна система](ch04-evolution-from-bayes-to-evidence-ai.md#L717) |
| H2 | [Доказовість як головна вимога](ch04-evolution-from-bayes-to-evidence-ai.md#L777) |
| H2 | [Чотири практичні уроки](ch04-evolution-from-bayes-to-evidence-ai.md#L840) |
| H3 | [Урок 1. Сміття не повинно потрапити в індекс](ch04-evolution-from-bayes-to-evidence-ai.md#L844) |
| H3 | [Урок 2. Між пошуком і мовною моделлю має стояти фільтр](ch04-evolution-from-bayes-to-evidence-ai.md#L864) |
| H3 | [Урок 3. Відтворюваність означає знімок і версії](ch04-evolution-from-bayes-to-evidence-ai.md#L956) |
| H3 | [Урок 4. Впевненість є набором індикаторів, а не одним числом](ch04-evolution-from-bayes-to-evidence-ai.md#L973) |
| H2 | [Висновки](ch04-evolution-from-bayes-to-evidence-ai.md#L991) |
| H2 | [Запитання до читачів](ch04-evolution-from-bayes-to-evidence-ai.md#L1004) |
| H2 | [Словник](ch04-evolution-from-bayes-to-evidence-ai.md#L1010) |
| H2 | [Абревіатури](ch04-evolution-from-bayes-to-evidence-ai.md#L1056) |
| H2 | [Джерела](ch04-evolution-from-bayes-to-evidence-ai.md#L1092) |

</details>

<details>
<summary>Глава 5. Тріада довіри: експертна система, доказова рекомендація та корпоративна пам'ять</summary>

| Рівень | Заголовок |
|---|---|
| H1 | [Глава 5. Тріада довіри: експертна система, доказова рекомендація та корпоративна пам'ять](ch05-triad-of-trust-and-corporate-memory.md#L1) |
| H2 | [Пам'ять конкретного рішення](ch05-triad-of-trust-and-corporate-memory.md#L15) |
| H2 | [Трикутник довіри](ch05-triad-of-trust-and-corporate-memory.md#L29) |
| H2 | [Перший кут: експертна система не дорівнює чат-боту](ch05-triad-of-trust-and-corporate-memory.md#L70) |
| H3 | [Що вміє велика мовна модель](ch05-triad-of-trust-and-corporate-memory.md#L74) |
| H3 | [Що додає експертна система](ch05-triad-of-trust-and-corporate-memory.md#L82) |
| H3 | [Поділ ролей і рівні знань](ch05-triad-of-trust-and-corporate-memory.md#L90) |
| H3 | [Як виміряти якість експертної системи](ch05-triad-of-trust-and-corporate-memory.md#L98) |
| H2 | [Другий кут: коли рекомендація ШІ стає доказовою](ch05-triad-of-trust-and-corporate-memory.md#L112) |
| H3 | [Приклад: драйвер пристрою й фраза «ризик прийнятний»](ch05-triad-of-trust-and-corporate-memory.md#L116) |
| H3 | [Що вважається доказом в інженерії](ch05-triad-of-trust-and-corporate-memory.md#L122) |
| H3 | [Мінімальний доказовий запис](ch05-triad-of-trust-and-corporate-memory.md#L137) |
| H3 | [Дисципліна запиту на прикладі перевірки коду](ch05-triad-of-trust-and-corporate-memory.md#L277) |
| H3 | [Журнал аудиту](ch05-triad-of-trust-and-corporate-memory.md#L317) |
| H2 | [Третій кут: корпоративна пам'ять досліджень і розробки](ch05-triad-of-trust-and-corporate-memory.md#L325) |
| H3 | [Де живуть знання про дослідження й розробку](ch05-triad-of-trust-and-corporate-memory.md#L331) |
| H3 | [Чому пошук не розв'язує проблеми](ch05-triad-of-trust-and-corporate-memory.md#L380) |
| H3 | [Адаптація новачків як тест пам'яті](ch05-triad-of-trust-and-corporate-memory.md#L388) |
| H3 | [Активні вивчені уроки](ch05-triad-of-trust-and-corporate-memory.md#L394) |
| H3 | [Мовна модель, граф знань і експертна система разом](ch05-triad-of-trust-and-corporate-memory.md#L427) |
| H3 | [Власники знань і неявне знання](ch05-triad-of-trust-and-corporate-memory.md#L435) |
| H3 | [Чутливість і межі доступу](ch05-triad-of-trust-and-corporate-memory.md#L445) |
| H3 | [Як виміряти корпоративну пам'ять](ch05-triad-of-trust-and-corporate-memory.md#L449) |
| H2 | [Тріада в п'яти галузях](ch05-triad-of-trust-and-corporate-memory.md#L459) |
| H2 | [Висновки](ch05-triad-of-trust-and-corporate-memory.md#L473) |
| H2 | [Підсумки Частини I](ch05-triad-of-trust-and-corporate-memory.md#L482) |
| H2 | [Подальший шлях пізнання](ch05-triad-of-trust-and-corporate-memory.md#L494) |
| H2 | [Запитання до читачів](ch05-triad-of-trust-and-corporate-memory.md#L504) |
| H2 | [Словник](ch05-triad-of-trust-and-corporate-memory.md#L511) |
| H2 | [Абревіатури](ch05-triad-of-trust-and-corporate-memory.md#L559) |
| H2 | [Джерела](ch05-triad-of-trust-and-corporate-memory.md#L580) |

</details>

<details>
<summary>Глава 6. Прикладна математика експертних систем: правила, ймовірності, графи та причинність</summary>

| Рівень | Заголовок |
|---|---|
| H1 | [Глава 6. Прикладна математика експертних систем: правила, ймовірності, графи та причинність](ch06-applied-mathematics-for-expert-systems.md#L1) |
| H2 | [Мінімум для першої програмної перевірки](ch06-applied-mathematics-for-expert-systems.md#L15) |
| H2 | [Карта: дев'ять запитань і дев'ять інструментів](ch06-applied-mathematics-for-expert-systems.md#L21) |
| H2 | [Логіка й правила: з чого випливає висновок](ch06-applied-mathematics-for-expert-systems.md#L39) |
| H3 | [Продукційне правило](ch06-applied-mathematics-for-expert-systems.md#L43) |
| H3 | [Пряме й зворотне виведення](ch06-applied-mathematics-for-expert-systems.md#L63) |
| H3 | [Логіка предикатів: одне правило для всіх компонентів](ch06-applied-mathematics-for-expert-systems.md#L69) |
| H3 | [Datalog і найменша нерухома точка](ch06-applied-mathematics-for-expert-systems.md#L91) |
| H3 | [Відкритий світ і тризначна логіка Кліні](ch06-applied-mathematics-for-expert-systems.md#L246) |
| H3 | [Деонтична логіка: обов'язок, заборона, дозвіл](ch06-applied-mathematics-for-expert-systems.md#L287) |
| H2 | [Дедукція, індукція, абдукція й зведення задач](ch06-applied-mathematics-for-expert-systems.md#L328) |
| H3 | [Дедукція: застосування правила](ch06-applied-mathematics-for-expert-systems.md#L360) |
| H3 | [Індукція: узагальнення спостережень](ch06-applied-mathematics-for-expert-systems.md#L374) |
| H3 | [Абдукція: гіпотеза про причину](ch06-applied-mathematics-for-expert-systems.md#L389) |
| H3 | [Зведення складної задачі до простіших](ch06-applied-mathematics-for-expert-systems.md#L408) |
| H2 | [Рушій правил: Rete, черга правил і розв'язання конфліктів](ch06-applied-mathematics-for-expert-systems.md#L444) |
| H2 | [Підтримання істинності: відкликання висновку без підстави](ch06-applied-mathematics-for-expert-systems.md#L456) |
| H2 | [Невизначеність: ймовірність, процеси в часі, нечіткість і свідчення](ch06-applied-mathematics-for-expert-systems.md#L477) |
| H3 | [Чому формула Байєса правильна](ch06-applied-mathematics-for-expert-systems.md#L481) |
| H3 | [Байєсівська мережа: менше параметрів і пояснення через іншу причину](ch06-applied-mathematics-for-expert-systems.md#L512) |
| H3 | [Марковські моделі: процеси в часі](ch06-applied-mathematics-for-expert-systems.md#L555) |
| H3 | [Нечітка логіка: ступінь замість різкої межі](ch06-applied-mathematics-for-expert-systems.md#L588) |
| H3 | [Теорія свідчень: коли джерела суперечать одне одному](ch06-applied-mathematics-for-expert-systems.md#L605) |
| H2 | [Досвід і зв'язки: прецеденти, графи й онтології](ch06-applied-mathematics-for-expert-systems.md#L732) |
| H3 | [Міркування за прецедентами](ch06-applied-mathematics-for-expert-systems.md#L736) |
| H3 | [Графи й онтології](ch06-applied-mathematics-for-expert-systems.md#L755) |
| H3 | [Ієрархія відношень](ch06-applied-mathematics-for-expert-systems.md#L775) |
| H2 | [Вибір і дії: багатокритеріальні рішення, оптимізація, планування](ch06-applied-mathematics-for-expert-systems.md#L805) |
| H3 | [Багатокритеріальний вибір](ch06-applied-mathematics-for-expert-systems.md#L809) |
| H3 | [Оптимізація, обмеження й планування](ch06-applied-mathematics-for-expert-systems.md#L901) |
| H2 | [Пошук і мовні моделі: від слів до векторів](ch06-applied-mathematics-for-expert-systems.md#L920) |
| H3 | [Зважування слів: TF-IDF і BM25](ch06-applied-mathematics-for-expert-systems.md#L924) |
| H3 | [Векторні представлення й контрастне навчання](ch06-applied-mathematics-for-expert-systems.md#L954) |
| H3 | [Гібридний пошук і злиття рангів](ch06-applied-mathematics-for-expert-systems.md#L986) |
| H3 | [Мовна модель і генерація з пошуком](ch06-applied-mathematics-for-expert-systems.md#L1045) |
| H2 | [Перевірка вимірювань: відстань Махаланобіса](ch06-applied-mathematics-for-expert-systems.md#L1101) |
| H2 | [Синергетика: редукція розмірності через параметри порядку та принцип підпорядкування Хакена](ch06-applied-mathematics-for-expert-systems.md#L1132) |
| H2 | [Причинність: кореляція не є причиною](ch06-applied-mathematics-for-expert-systems.md#L1158) |
| H3 | [Драбина причинності](ch06-applied-mathematics-for-expert-systems.md#L1169) |
| H3 | [Втручання й коригування за обхідними шляхами](ch06-applied-mathematics-for-expert-systems.md#L1196) |
| H3 | [Структурна причинна модель і контрфакти](ch06-applied-mathematics-for-expert-systems.md#L1212) |
| H2 | [Пояснюваність і калібрування](ch06-applied-mathematics-for-expert-systems.md#L1229) |
| H2 | [Два уроки математичного шару](ch06-applied-mathematics-for-expert-systems.md#L1237) |
| H3 | [Урок 1. Оцінка без порогу не є доказом](ch06-applied-mathematics-for-expert-systems.md#L1239) |
| H3 | [Урок 2. Модель без сценарію відмови говорить зайве](ch06-applied-mathematics-for-expert-systems.md#L1243) |
| H2 | [Висновки](ch06-applied-mathematics-for-expert-systems.md#L1247) |
| H2 | [Запитання до читачів](ch06-applied-mathematics-for-expert-systems.md#L1258) |
| H2 | [Словник](ch06-applied-mathematics-for-expert-systems.md#L1265) |
| H2 | [Абревіатури](ch06-applied-mathematics-for-expert-systems.md#L1351) |
| H2 | [Джерела](ch06-applied-mathematics-for-expert-systems.md#L1386) |

</details>

<details>
<summary>Глава 7. Типологія баз знань: правила, онтології, прецеденти та вектори</summary>

| Рівень | Заголовок |
|---|---|
| H1 | [Глава 7. Типологія баз знань: правила, онтології, прецеденти та вектори](ch07-knowledge-base-typology.md#L1) |
| H2 | [Типові збої: коли тип знань не відповідає запитанню](ch07-knowledge-base-typology.md#L15) |
| H2 | [База знань і база даних](ch07-knowledge-base-typology.md#L31) |
| H3 | [Коли набір файлів стає базою знань](ch07-knowledge-base-typology.md#L44) |
| H3 | [Вікі як корпоративна база знань](ch07-knowledge-base-typology.md#L56) |
| H3 | [Як користувач ставить запити до бази знань](ch07-knowledge-base-typology.md#L71) |
| H3 | [Вимоги, код, тести й задачі як розподілена база знань](ch07-knowledge-base-typology.md#L120) |
| H2 | [Відкриті бази знань, які можна дослідити](ch07-knowledge-base-typology.md#L136) |
| H2 | [Наскрізний приклад: випуск BrakeController 3.2](ch07-knowledge-base-typology.md#L151) |
| H2 | [Порівняння типів баз знань](ch07-knowledge-base-typology.md#L227) |
| H2 | [База правил, база фактів і робоча пам'ять](ch07-knowledge-base-typology.md#L244) |
| H3 | [Не зливати твердження голосуванням більшості](ch07-knowledge-base-typology.md#L311) |
| H3 | [База фактів і база прецедентів не одне й те саме](ch07-knowledge-base-typology.md#L328) |
| H2 | [Продукційна база знань: правила «якщо, то»](ch07-knowledge-base-typology.md#L344) |
| H2 | [Фреймова база знань: об'єкти з ролями й типовими значеннями](ch07-knowledge-base-typology.md#L352) |
| H2 | [Семантичні мережі й онтології: типи й зв'язки](ch07-knowledge-base-typology.md#L398) |
| H3 | [Запитання до компетентності: що має дозволяти модель](ch07-knowledge-base-typology.md#L466) |
| H3 | [OntoClean: тип об'єкта, роль і змінний стан](ch07-knowledge-base-typology.md#L484) |
| H2 | [База прецедентів: пам'ять про те, що вже траплялося](ch07-knowledge-base-typology.md#L503) |
| H2 | [База обмежень: які конфігурації допустимі](ch07-knowledge-base-typology.md#L575) |
| H2 | [Ймовірнісна база знань: наскільки правдоподібна гіпотеза](ch07-knowledge-base-typology.md#L617) |
| H2 | [Нечітка база знань: робота з градаціями](ch07-knowledge-base-typology.md#L670) |
| H2 | [Документний і векторний індекси: пошук кандидатів, а не істини](ch07-knowledge-base-typology.md#L691) |
| H2 | [Гібридна архітектура: нейронний пошук, граф і символьна перевірка](ch07-knowledge-base-typology.md#L728) |
| H2 | [Відсутній факт: відкритий, закритий і частково закритий світ](ch07-knowledge-base-typology.md#L805) |
| H2 | [Розбиття бази знань: фрагментація, шардування й групування за ознаками](ch07-knowledge-base-typology.md#L918) |
| H3 | [Три рішення, які називають одним словом](ch07-knowledge-base-typology.md#L922) |
| H3 | [Фрагментація: види і правила коректності](ch07-knowledge-base-typology.md#L932) |
| H3 | [Ключ шардування: рівномірно, але не ціною розриву виведення](ch07-knowledge-base-typology.md#L953) |
| H3 | [Ознаки розбиття в експертній системі](ch07-knowledge-base-typology.md#L981) |
| H3 | [Шість правил, що відрізняють розбиття знань від розбиття даних](ch07-knowledge-base-typology.md#L996) |
| H3 | [Що розбиття коштує й коли його не потрібно](ch07-knowledge-base-typology.md#L1048) |
| H2 | [Як обрати тип бази знань](ch07-knowledge-base-typology.md#L1056) |
| H2 | [Висновки](ch07-knowledge-base-typology.md#L1086) |
| H2 | [Запитання до читачів](ch07-knowledge-base-typology.md#L1098) |
| H2 | [Словник](ch07-knowledge-base-typology.md#L1109) |
| H2 | [Абревіатури](ch07-knowledge-base-typology.md#L1176) |
| H2 | [Джерела](ch07-knowledge-base-typology.md#L1203) |

</details>

<details>
<summary>Глава 8. Інженерні артефакти як дані експертної системи</summary>

| Рівень | Заголовок |
|---|---|
| H1 | [Глава 8. Інженерні артефакти як дані експертної системи](ch08-engineering-artifacts-as-data.md#L1) |
| H2 | [Карта глави: проблема, метод і вимірювання](ch08-engineering-artifacts-as-data.md#L15) |
| H2 | [Чому тексту недостатньо](ch08-engineering-artifacts-as-data.md#L29) |
| H2 | [Артефакт як об'єкт даних, а не документ](ch08-engineering-artifacts-as-data.md#L37) |
| H2 | [Вимога як структурований об'єкт](ch08-engineering-artifacts-as-data.md#L47) |
| H2 | [Автоматична перевірка якості вимог](ch08-engineering-artifacts-as-data.md#L91) |
| H2 | [Простежуваність як головна причина структури](ch08-engineering-artifacts-as-data.md#L156) |
| H2 | [Плани, ризики й контрольні точки як об'єкти](ch08-engineering-artifacts-as-data.md#L175) |
| H2 | [Формат носія визначає, скільки структури збережеться](ch08-engineering-artifacts-as-data.md#L183) |
| H2 | [Первинна обробка: від байтів до канонічного документа](ch08-engineering-artifacts-as-data.md#L202) |
| H2 | [Від файла до фрагмента: базове нарізання](ch08-engineering-artifacts-as-data.md#L287) |
| H2 | [Первинне сортування фрагментів: знання чи сміття](ch08-engineering-artifacts-as-data.md#L399) |
| H2 | [Векторизація й два індекси](ch08-engineering-artifacts-as-data.md#L554) |
| H3 | [Що таке токен](ch08-engineering-artifacts-as-data.md#L558) |
| H3 | [Від токенів до вектора фрагмента](ch08-engineering-artifacts-as-data.md#L587) |
| H3 | [Мітка доступу рухається разом із фрагментом](ch08-engineering-artifacts-as-data.md#L616) |
| H3 | [Як довести, що пошук справді покращився](ch08-engineering-artifacts-as-data.md#L646) |
| H2 | [Яке обчислення на якому обладнанні](ch08-engineering-artifacts-as-data.md#L664) |
| H2 | [Від фрагмента до заповненого об'єкта](ch08-engineering-artifacts-as-data.md#L678) |
| H2 | [Добрий артефакт схожий на добрий код](ch08-engineering-artifacts-as-data.md#L686) |
| H2 | [Як структура допомагає штучному інтелекту](ch08-engineering-artifacts-as-data.md#L702) |
| H3 | [Велика генеративна чи компактна спеціалізована модель](ch08-engineering-artifacts-as-data.md#L706) |
| H3 | [Відповіді над графом об'єктів](ch08-engineering-artifacts-as-data.md#L747) |
| H2 | [Аудит і відповідність](ch08-engineering-artifacts-as-data.md#L753) |
| H2 | [Ризики надмірної формалізації](ch08-engineering-artifacts-as-data.md#L759) |
| H2 | [Як почати без великої перебудови](ch08-engineering-artifacts-as-data.md#L765) |
| H2 | [Висновки](ch08-engineering-artifacts-as-data.md#L771) |
| H2 | [Запитання до читачів](ch08-engineering-artifacts-as-data.md#L781) |
| H2 | [Словник](ch08-engineering-artifacts-as-data.md#L790) |
| H2 | [Абревіатури](ch08-engineering-artifacts-as-data.md#L841) |
| H2 | [Джерела](ch08-engineering-artifacts-as-data.md#L874) |

</details>

<details>
<summary>Глава 9. Інженерний граф знань: простежуваність від вимог до апаратури</summary>

| Рівень | Заголовок |
|---|---|
| H1 | [Глава 9. Інженерний граф знань: простежуваність від вимог до апаратури](ch09-engineering-knowledge-graph-traceability.md#L1) |
| H2 | [Розрив простежуваності](ch09-engineering-knowledge-graph-traceability.md#L63) |
| H2 | [Формальна модель інженерного графа знань](ch09-engineering-knowledge-graph-traceability.md#L82) |
| H3 | [Типи вузлів](ch09-engineering-knowledge-graph-traceability.md#L157) |
| H3 | [Типи ребер](ch09-engineering-knowledge-graph-traceability.md#L166) |
| H2 | [Автоматичне вилучення й побудова графа](ch09-engineering-knowledge-graph-traceability.md#L183) |
| H3 | [Код і вимоги через синтаксичне дерево](ch09-engineering-knowledge-graph-traceability.md#L228) |
| H3 | [Код і апаратні регістри](ch09-engineering-knowledge-graph-traceability.md#L254) |
| H3 | [Тести й вимоги](ch09-engineering-knowledge-graph-traceability.md#L296) |
| H3 | [Сутності, кореферентність, розрідженість і часова чинність](ch09-engineering-knowledge-graph-traceability.md#L317) |
| H2 | [Аналіз на інженерному графі](ch09-engineering-knowledge-graph-traceability.md#L328) |
| H3 | [Пошук прогалин і покриття вимог](ch09-engineering-knowledge-graph-traceability.md#L353) |
| H3 | [Прогнозування пропущених зв'язків: кандидати, а не факти](ch09-engineering-knowledge-graph-traceability.md#L371) |
| H3 | [Непростежуваний і мертвий код](ch09-engineering-knowledge-graph-traceability.md#L397) |
| H3 | [Аналіз впливу змін](ch09-engineering-knowledge-graph-traceability.md#L414) |
| H2 | [Архітектура промислової платформи](ch09-engineering-knowledge-graph-traceability.md#L462) |
| H3 | [Незмінні версійовані факти](ch09-engineering-knowledge-graph-traceability.md#L516) |
| H3 | [Контрольна точка якості в конвеєрі змін](ch09-engineering-knowledge-graph-traceability.md#L533) |
| H2 | [Холодний старт: автоматична початкова версія графа](ch09-engineering-knowledge-graph-traceability.md#L543) |
| H2 | [Таблиці простежуваності проти інженерного графа](ch09-engineering-knowledge-graph-traceability.md#L865) |
| H2 | [Висновки](ch09-engineering-knowledge-graph-traceability.md#L881) |
| H2 | [Підсумки маршруту від математичної моделі до графа](ch09-engineering-knowledge-graph-traceability.md#L892) |
| H2 | [Подальший шлях пізнання](ch09-engineering-knowledge-graph-traceability.md#L903) |
| H2 | [Запитання до читачів](ch09-engineering-knowledge-graph-traceability.md#L909) |
| H2 | [Словник](ch09-engineering-knowledge-graph-traceability.md#L916) |
| H2 | [Абревіатури](ch09-engineering-knowledge-graph-traceability.md#L953) |
| H2 | [Джерела](ch09-engineering-knowledge-graph-traceability.md#L977) |

</details>

<details>
<summary>Глава 10. Системи здобуття знань: джерела, допуск і життєвий цикл</summary>

| Рівень | Заголовок |
|---|---|
| H1 | [Глава 10. Системи здобуття знань: джерела, допуск і життєвий цикл](ch10-knowledge-acquisition-systems.md#L1) |
| H2 | [Що таке здобуття знань і система здобуття знань](ch10-knowledge-acquisition-systems.md#L15) |
| H2 | [Чому небезпечно «дати ШІ всі документи»](ch10-knowledge-acquisition-systems.md#L25) |
| H3 | [Три рівні зрілості корпоративних систем знань](ch10-knowledge-acquisition-systems.md#L33) |
| H2 | [Що має забезпечити здобуття знань](ch10-knowledge-acquisition-systems.md#L43) |
| H2 | [Етапи конвеєра здобуття знань](ch10-knowledge-acquisition-systems.md#L60) |
| H2 | [Курування: які рішення залишаються за людиною](ch10-knowledge-acquisition-systems.md#L120) |
| H2 | [Різна довіра до джерел знань](ch10-knowledge-acquisition-systems.md#L126) |
| H2 | [Класифікація, доступ і знеособлення](ch10-knowledge-acquisition-systems.md#L132) |
| H2 | [Недовірений вхід: файл атакує аналізатор, текст атакує рішення](ch10-knowledge-acquisition-systems.md#L191) |
| H2 | [Збирання документів і контроль якості розбору](ch10-knowledge-acquisition-systems.md#L233) |
| H2 | [Що показав експеримент із виявлення знань](ch10-knowledge-acquisition-systems.md#L251) |
| H2 | [Як математика допомагає відбирати знання](ch10-knowledge-acquisition-systems.md#L308) |
| H2 | [Операційні метрики здобуття знань](ch10-knowledge-acquisition-systems.md#L382) |
| H2 | [Походження й родовід знань](ch10-knowledge-acquisition-systems.md#L433) |
| H2 | [Стандарти імпорту структурованих інженерних даних](ch10-knowledge-acquisition-systems.md#L439) |
| H2 | [KAS як програмна система](ch10-knowledge-acquisition-systems.md#L450) |
| H2 | [Розподіл відповідальності між KAS і генерацією з пошуком](ch10-knowledge-acquisition-systems.md#L497) |
| H2 | [Розподіл навантаження між CPU, GPU і NPU](ch10-knowledge-acquisition-systems.md#L503) |
| H2 | [Паспорт об'єкта: не змішувати різні види перевірки](ch10-knowledge-acquisition-systems.md#L536) |
| H2 | [Життєвий цикл і власник об'єкта знань](ch10-knowledge-acquisition-systems.md#L555) |
| H2 | [Випуск бази знань](ch10-knowledge-acquisition-systems.md#L714) |
| H2 | [Перелік перевірок перед промисловою експлуатацією](ch10-knowledge-acquisition-systems.md#L767) |
| H2 | [Здобуття знань 2.0: від пасивного індексування до активного виявлення знань](ch10-knowledge-acquisition-systems.md#L781) |
| H2 | [Висновки](ch10-knowledge-acquisition-systems.md#L812) |
| H2 | [Запитання до читачів](ch10-knowledge-acquisition-systems.md#L823) |
| H2 | [Словник](ch10-knowledge-acquisition-systems.md#L830) |
| H2 | [Абревіатури](ch10-knowledge-acquisition-systems.md#L873) |
| H2 | [Джерела](ch10-knowledge-acquisition-systems.md#L907) |

</details>

<details>
<summary>Глава 11. Вилучення знань у експертів: інтерв'ювання, когнітивні карти та формалізація досвіду</summary>

| Рівень | Заголовок |
|---|---|
| H1 | [Глава 11. Вилучення знань у експертів: інтерв'ювання, когнітивні карти та формалізація досвіду](ch11-knowledge-elicitation-from-experts.md#L1) |
| H2 | [Чому запис розмови ще не є знанням](ch11-knowledge-elicitation-from-experts.md#L17) |
| H2 | [Спершу визначити рішення, а не призначати інтерв'ю](ch11-knowledge-elicitation-from-experts.md#L59) |
| H3 | [CommonKADS: предметні знання, виведення й задача](ch11-knowledge-elicitation-from-experts.md#L75) |
| H2 | [Як вибрати спосіб роботи з експертом](ch11-knowledge-elicitation-from-experts.md#L93) |
| H3 | [Тріада замість запитання «що важливо?»](ch11-knowledge-elicitation-from-experts.md#L110) |
| H3 | [Критичний випадок замість загального правила](ch11-knowledge-elicitation-from-experts.md#L114) |
| H2 | [Як перетворити фразу на кандидата знання](ch11-knowledge-elicitation-from-experts.md#L140) |
| H2 | [Як ставити запитання, що виявляють межі](ch11-knowledge-elicitation-from-experts.md#L179) |
| H2 | [Як відрізнити впевненість від частоти](ch11-knowledge-elicitation-from-experts.md#L196) |
| H2 | [Яке запитання поставити наступним](ch11-knowledge-elicitation-from-experts.md#L243) |
| H2 | [Кілька експертів: розбіжності не усереднюють автоматично](ch11-knowledge-elicitation-from-experts.md#L314) |
| H2 | [Роль мовних моделей у вилученні знань](ch11-knowledge-elicitation-from-experts.md#L335) |
| H2 | [Як перевірити отримане правило](ch11-knowledge-elicitation-from-experts.md#L373) |
| H2 | [Типові помилки вилучення знань](ch11-knowledge-elicitation-from-experts.md#L400) |
| H2 | [Перший цикл вилучення знань](ch11-knowledge-elicitation-from-experts.md#L418) |
| H2 | [Неявні знання: від моделі SECI до цифрового сліду інженерної роботи](ch11-knowledge-elicitation-from-experts.md#L428) |
| H3 | [Модель SECI у цифрову добу](ch11-knowledge-elicitation-from-experts.md#L432) |
| H3 | [Пасивне здобуття знань із цифрового сліду роботи](ch11-knowledge-elicitation-from-experts.md#L475) |
| H2 | [Висновки](ch11-knowledge-elicitation-from-experts.md#L486) |
| H2 | [Запитання до читачів](ch11-knowledge-elicitation-from-experts.md#L493) |
| H2 | [Словник](ch11-knowledge-elicitation-from-experts.md#L502) |
| H2 | [Абревіатури](ch11-knowledge-elicitation-from-experts.md#L522) |
| H2 | [Джерела](ch11-knowledge-elicitation-from-experts.md#L530) |

</details>

<details>
<summary>Глава 12. Лінгвістичний аналіз і локальні моделі: збереження змісту та джерела</summary>

| Рівень | Заголовок |
|---|---|
| H1 | [Глава 12. Лінгвістичний аналіз і локальні моделі: збереження змісту та джерела](ch12-linguistic-analysis-and-local-models.md#L1) |
| H2 | [Від людської репліки до права стверджувати](ch12-linguistic-analysis-and-local-models.md#L19) |
| H2 | [Як не загубити зв'язок між текстом і джерелом](ch12-linguistic-analysis-and-local-models.md#L60) |
| H3 | [Чому нормалізація може звести різні записи до одного](ch12-linguistic-analysis-and-local-models.md#L107) |
| H2 | [Як текст може приховати підміну](ch12-linguistic-analysis-and-local-models.md#L161) |
| H2 | [Як розподілити відповідальність між складниками](ch12-linguistic-analysis-and-local-models.md#L210) |
| H2 | [Бюджет контексту, токенізація й пам'ять моделі](ch12-linguistic-analysis-and-local-models.md#L260) |
| H3 | [Токенізатор і українська мова](ch12-linguistic-analysis-and-local-models.md#L300) |
| H3 | [Межа мовної моделі: зсув розподілу даних](ch12-linguistic-analysis-and-local-models.md#L320) |
| H3 | [Пам'ять кешу ключів і значень](ch12-linguistic-analysis-and-local-models.md#L332) |
| H2 | [Чому схожість тексту ще не є доказом](ch12-linguistic-analysis-and-local-models.md#L355) |
| H2 | [Парсер документа й навчальні дані: де шукати втрати](ch12-linguistic-analysis-and-local-models.md#L428) |
| H2 | [Профілі розгортання локальних моделей](ch12-linguistic-analysis-and-local-models.md#L444) |
| H2 | [Контроль регресій під час квантування й оновлення моделей](ch12-linguistic-analysis-and-local-models.md#L459) |
| H2 | [Як перевіряти мовну підсистему](ch12-linguistic-analysis-and-local-models.md#L525) |
| H2 | [Висновки](ch12-linguistic-analysis-and-local-models.md#L606) |
| H2 | [Запитання до читачів](ch12-linguistic-analysis-and-local-models.md#L611) |
| H2 | [Словник](ch12-linguistic-analysis-and-local-models.md#L618) |
| H2 | [Абревіатури](ch12-linguistic-analysis-and-local-models.md#L639) |
| H2 | [Джерела](ch12-linguistic-analysis-and-local-models.md#L655) |

</details>

<details>
<summary>Глава 13. Варіативність природної мови проти детермінізму: компіляція сенсу запитання</summary>

| Рівень | Заголовок |
|---|---|
| H1 | [Глава 13. Варіативність природної мови проти детермінізму: компіляція сенсу запитання](ch13-language-variability-vs-determinism.md#L1) |
| H2 | [Як еволюціонував RAG: три парадигми та їхні межі](ch13-language-variability-vs-determinism.md#L44) |
| H3 | [Наївний RAG: «знайти і прочитати»](ch13-language-variability-vs-determinism.md#L50) |
| H3 | [Просунутий RAG: оптимізація до і після пошуку](ch13-language-variability-vs-determinism.md#L99) |
| H3 | [Модульний RAG: компоненти замість конвеєра](ch13-language-variability-vs-determinism.md#L109) |
| H2 | [Чому схожість тексту не є доказом розуміння](ch13-language-variability-vs-determinism.md#L173) |
| H2 | [Конвеєр семантичної компіляції запитання](ch13-language-variability-vs-determinism.md#L218) |
| H3 | [Крок 1. Синтаксичний розбір залежностей](ch13-language-variability-vs-determinism.md#L244) |
| H3 | [Крок 2. Доменне тегування та ізоляція констант](ch13-language-variability-vs-determinism.md#L248) |
| H3 | [Крок 3. Семантичний розбір малою мовною моделлю](ch13-language-variability-vs-determinism.md#L252) |
| H3 | [Крок 4. Онтологічна перевірка верифікатором хоста](ch13-language-variability-vs-determinism.md#L273) |
| H3 | [Умовні запитання: гіпотези окремо від фактів](ch13-language-variability-vs-determinism.md#L293) |
| H2 | [Обмежене декодування за граматикою](ch13-language-variability-vs-determinism.md#L433) |
| H2 | [Перевірка стійкості групами еквівалентних запитань](ch13-language-variability-vs-determinism.md#L495) |
| H2 | [Засоби реалізації: від сутності до типізованого запиту](ch13-language-variability-vs-determinism.md#L564) |
| H2 | [Розподіл відповідальності між компонентами](ch13-language-variability-vs-determinism.md#L581) |
| H2 | [Висновки](ch13-language-variability-vs-determinism.md#L622) |
| H2 | [Запитання до читачів](ch13-language-variability-vs-determinism.md#L631) |
| H2 | [Словник](ch13-language-variability-vs-determinism.md#L642) |
| H2 | [Абревіатури](ch13-language-variability-vs-determinism.md#L669) |
| H2 | [Джерела](ch13-language-variability-vs-determinism.md#L696) |

</details>

<details>
<summary>Глава 14. Виявлення вимог і модальностей: від нормативного тексту до інваріантів</summary>

| Рівень | Заголовок |
|---|---|
| H1 | [Глава 14. Виявлення вимог і модальностей: від нормативного тексту до інваріантів](ch14-requirements-detection-and-formalization.md#L1) |
| H2 | [Не кожне речення зі словом shall є вимогою](ch14-requirements-detection-and-formalization.md#L43) |
| H3 | [Синтаксичні пастки природної мови](ch14-requirements-detection-and-formalization.md#L64) |
| H2 | [Шаблони EARS: шість форм замість вільного тексту](ch14-requirements-detection-and-formalization.md#L76) |
| H2 | [Формалізація: від шаблону до логіки](ch14-requirements-detection-and-formalization.md#L123) |
| H3 | [Логіка предикатів першого порядку](ch14-requirements-detection-and-formalization.md#L127) |
| H3 | [Темпоральна логіка](ch14-requirements-detection-and-formalization.md#L143) |
| H3 | [Перевірка суперечностей SMT-розв'язувачем](ch14-requirements-detection-and-formalization.md#L185) |
| H2 | [Аудит якості вимог](ch14-requirements-detection-and-formalization.md#L234) |
| H2 | [Детермінований конвеєр на Go](ch14-requirements-detection-and-formalization.md#L268) |
| H2 | [Інженерний кейс: аудит вимоги до системи керування батареєю](ch14-requirements-detection-and-formalization.md#L505) |
| H2 | [Типізоване проміжне подання: вимога як вхід компілятора](ch14-requirements-detection-and-formalization.md#L617) |
| H2 | [Висновки](ch14-requirements-detection-and-formalization.md#L639) |
| H2 | [Запитання до читачів](ch14-requirements-detection-and-formalization.md#L646) |
| H2 | [Словник](ch14-requirements-detection-and-formalization.md#L653) |
| H2 | [Абревіатури](ch14-requirements-detection-and-formalization.md#L672) |
| H2 | [Джерела](ch14-requirements-detection-and-formalization.md#L694) |

</details>

<details>
<summary>Глава 15. Вилучення знань і побудова бази знань: факти, граматики та автомати</summary>

| Рівень | Заголовок |
|---|---|
| H1 | [Глава 15. Вилучення знань і побудова бази знань: факти, граматики та автомати](ch15-knowledge-extraction-and-kb-construction.md#L1) |
| H2 | [Джерела інженерних знань: чотири класи документів](ch15-knowledge-extraction-and-kb-construction.md#L23) |
| H2 | [Інваріант доказовості: кожен факт має адресу в першоджерелі](ch15-knowledge-extraction-and-kb-construction.md#L79) |
| H2 | [Конвеєр екстракції: від байтів до кандидатів у факти](ch15-knowledge-extraction-and-kb-construction.md#L100) |
| H3 | [Фаза 1. Структурна декомпозиція](ch15-knowledge-extraction-and-kb-construction.md#L154) |
| H3 | [Фази 2 і 3. Розбір речень і спеціалізовані екстрактори](ch15-knowledge-extraction-and-kb-construction.md#L164) |
| H4 | [Параметри й реляційні факти](ch15-knowledge-extraction-and-kb-construction.md#L168) |
| H4 | [Нормативні вимоги](ch15-knowledge-extraction-and-kb-construction.md#L184) |
| H4 | [Формальні граматики](ch15-knowledge-extraction-and-kb-construction.md#L188) |
| H3 | [Фаза 4. Побайтова прив'язка](ch15-knowledge-extraction-and-kb-construction.md#L205) |
| H3 | [Структурне розбиття для пошуку](ch15-knowledge-extraction-and-kb-construction.md#L374) |
| H2 | [Динаміка: вилучення скінченних автоматів](ch15-knowledge-extraction-and-kb-construction.md#L396) |
| H3 | [Математична модель](ch15-knowledge-extraction-and-kb-construction.md#L400) |
| H3 | [Приклад: сеанс SMTP](ch15-knowledge-extraction-and-kb-construction.md#L417) |
| H3 | [Як автомат вилучається з тексту](ch15-knowledge-extraction-and-kb-construction.md#L442) |
| H3 | [Негативні тести з автомата](ch15-knowledge-extraction-and-kb-construction.md#L466) |
| H3 | [Чотири стратегії виявлення автоматних переходів (FSM Patterns)](ch15-knowledge-extraction-and-kb-construction.md#L480) |
| H2 | [Формальні граматики: ABNF-екстракція та компіляція синтаксичних дерев](ch15-knowledge-extraction-and-kb-construction.md#L524) |
| H3 | [Елементи вилученого правила ABNF:](ch15-knowledge-extraction-and-kb-construction.md#L554) |
| H2 | [Параметри та структуровані фізичні величини: нормалізація одиниць СІ та розв'язання нерівностей](ch15-knowledge-extraction-and-kb-construction.md#L582) |
| H3 | [Математична нормалізація структурованих параметрів](ch15-knowledge-extraction-and-kb-construction.md#L588) |
| H2 | [Деонтичні винятки та умови скасування (Defeaters)](ch15-knowledge-extraction-and-kb-construction.md#L618) |
| H2 | [Версії документів: заміщення, оновлення та виправлення](ch15-knowledge-extraction-and-kb-construction.md#L652) |
| H3 | [Секційні патчі спадкоємності: точкова резолюція оновлень замість повторного перепарсингу](ch15-knowledge-extraction-and-kb-construction.md#L704) |
| H2 | [Сховище бази знань: три моделі для трьох типів запитів](ch15-knowledge-extraction-and-kb-construction.md#L739) |
| H2 | [Шлюз допуску: що потрапляє до бази знань](ch15-knowledge-extraction-and-kb-construction.md#L789) |
| H3 | [Мовна модель пропонує, шлюз вирішує](ch15-knowledge-extraction-and-kb-construction.md#L823) |
| H3 | [Безперервне поповнення бази знань](ch15-knowledge-extraction-and-kb-construction.md#L847) |
| H2 | [Практичний цикл онтології: Protégé й ROBOT](ch15-knowledge-extraction-and-kb-construction.md#L853) |
| H3 | [Мінімальна перевірка запитання до компетентності](ch15-knowledge-extraction-and-kb-construction.md#L872) |
| H2 | [Сертифікат доведення: відповідь, яку може перевірити програма](ch15-knowledge-extraction-and-kb-construction.md#L921) |
| H2 | [Що запозичити з аналізу даних і процесів](ch15-knowledge-extraction-and-kb-construction.md#L961) |
| H2 | [Порівняння: база знань із доказами та пошук із генерацією](ch15-knowledge-extraction-and-kb-construction.md#L973) |
| H2 | [Висновки](ch15-knowledge-extraction-and-kb-construction.md#L990) |
| H2 | [Підсумки маршруту здобуття та формалізації](ch15-knowledge-extraction-and-kb-construction.md#L999) |
| H2 | [Подальший шлях пізнання](ch15-knowledge-extraction-and-kb-construction.md#L1012) |
| H2 | [Запитання до читачів](ch15-knowledge-extraction-and-kb-construction.md#L1018) |
| H2 | [Словник](ch15-knowledge-extraction-and-kb-construction.md#L1028) |
| H2 | [Абревіатури](ch15-knowledge-extraction-and-kb-construction.md#L1059) |
| H2 | [Джерела](ch15-knowledge-extraction-and-kb-construction.md#L1097) |

</details>

<details>
<summary>Глава 16. Архітектура експертної системи: від формального знання до доказового рішення</summary>

| Рівень | Заголовок |
|---|---|
| H1 | [Глава 16. Архітектура експертної системи: від формального знання до доказового рішення](ch16-expert-systems-architecture.md#L1) |
| H2 | [Архітектурне ядро і межі відповідальності](ch16-expert-systems-architecture.md#L19) |
| H2 | [Знання в інженерному сенсі](ch16-expert-systems-architecture.md#L43) |
| H2 | [Фізичні форми знань](ch16-expert-systems-architecture.md#L58) |
| H2 | [Чому векторна база не може бути джерелом істини](ch16-expert-systems-architecture.md#L121) |
| H2 | [Референсна архітектура: шлях наповнення і шлях читання](ch16-expert-systems-architecture.md#L145) |
| H2 | [Життєвий цикл запиту: від запитання до доказового пакета](ch16-expert-systems-architecture.md#L204) |
| H2 | [Підтримання істинності: що відбувається з висновком, коли змінюється факт](ch16-expert-systems-architecture.md#L248) |
| H2 | [Узгоджений контекст читання й повторення рішення](ch16-expert-systems-architecture.md#L434) |
| H2 | [Контракти даних між підсистемами](ch16-expert-systems-architecture.md#L444) |
| H2 | [Режими прихованих відмов і протидія](ch16-expert-systems-architecture.md#L492) |
| H2 | [Калібрування якості та регресійний контроль](ch16-expert-systems-architecture.md#L507) |
| H2 | [Версіонування та походження рішень](ch16-expert-systems-architecture.md#L561) |
| H2 | [Сім архітектурних правил](ch16-expert-systems-architecture.md#L646) |
| H2 | [Висновки](ch16-expert-systems-architecture.md#L658) |
| H2 | [Запитання до читачів](ch16-expert-systems-architecture.md#L667) |
| H2 | [Словник](ch16-expert-systems-architecture.md#L676) |
| H2 | [Абревіатури](ch16-expert-systems-architecture.md#L700) |
| H2 | [Джерела](ch16-expert-systems-architecture.md#L743) |

</details>

<details>
<summary>Глава 17. Технологічний стек: критерії вибору інструментів, мов програмування та рушіїв правил</summary>

| Рівень | Заголовок |
|---|---|
| H1 | [Глава 17. Технологічний стек: критерії вибору інструментів, мов програмування та рушіїв правил](ch17-implementation-stack.md#L1) |
| H2 | [Мінімальна реалізація перед каталогом технологій](ch17-implementation-stack.md#L15) |
| H2 | [Мови програмування за шарами відповідальності](ch17-implementation-stack.md#L67) |
| H3 | [Go: сервісний шар і координація](ch17-implementation-stack.md#L71) |
| H3 | [C і C++: обчислювальні ядра](ch17-implementation-stack.md#L83) |
| H3 | [Rust: безпека пам'яті на межі з ненадійними даними](ch17-implementation-stack.md#L87) |
| H3 | [Python: офлайнова лабораторія знань](ch17-implementation-stack.md#L91) |
| H3 | [TypeScript: інтерфейси експертного аудиту](ch17-implementation-stack.md#L95) |
| H2 | [Мови запитів до знань](ch17-implementation-stack.md#L101) |
| H3 | [SQL і рекурсивні запити](ch17-implementation-stack.md#L105) |
| H3 | [Cypher, GQL і SPARQL](ch17-implementation-stack.md#L177) |
| H3 | [Datalog](ch17-implementation-stack.md#L181) |
| H2 | [Рушії правил і таблиці рішень](ch17-implementation-stack.md#L198) |
| H3 | [Класичні рушії правил](ch17-implementation-stack.md#L202) |
| H3 | [DMN і FEEL](ch17-implementation-stack.md#L206) |
| H2 | [Сховища за фізичною формою даних](ch17-implementation-stack.md#L221) |
| H2 | [Оптимізація, планування та розв'язання обмежень](ch17-implementation-stack.md#L236) |
| H2 | [Семантична перевірка кандидата рушія](ch17-implementation-stack.md#L247) |
| H2 | [Виконання моделей і апаратні прискорювачі](ch17-implementation-stack.md#L263) |
| H2 | [Карта п'ятнадцяти класів технологій](ch17-implementation-stack.md#L301) |
| H2 | [Компоненти Go за шарами можливостей](ch17-implementation-stack.md#L360) |
| H3 | [Правило: спершу порівняльний вимір, потім прийняття](ch17-implementation-stack.md#L376) |
| H2 | [Три уроки вибору технологій](ch17-implementation-stack.md#L389) |
| H2 | [Висновки](ch17-implementation-stack.md#L395) |
| H2 | [Запитання до читачів](ch17-implementation-stack.md#L402) |
| H2 | [Словник](ch17-implementation-stack.md#L411) |
| H2 | [Абревіатури](ch17-implementation-stack.md#L433) |
| H2 | [Джерела](ch17-implementation-stack.md#L474) |

</details>

<details>
<summary>Глава 18. Інфраструктура виконання: локальні моделі, апаратні прискорювачі, Edge та On-Premise</summary>

| Рівень | Заголовок |
|---|---|
| H1 | [Глава 18. Інфраструктура виконання: локальні моделі, апаратні прискорювачі, Edge та On-Premise](ch18-execution-infrastructure.md#L1) |
| H2 | [Профіль обчислень: чому правила й нейромережі потребують різної апаратури](ch18-execution-infrastructure.md#L57) |
| H2 | [Експертна система як координатор виконання](ch18-execution-infrastructure.md#L104) |
| H2 | [Ізоляція середовищ виконання і керування пам'яттю](ch18-execution-infrastructure.md#L141) |
| H3 | [Відображення бази фактів у пам'ять](ch18-execution-infrastructure.md#L154) |
| H2 | [Топології розгортання: від дата-центру до сенсора](ch18-execution-infrastructure.md#L192) |
| H2 | [Фонове індексування на простійних NPU робочих станцій](ch18-execution-infrastructure.md#L226) |
| H3 | [Які задачі підходять для мережі NPU](ch18-execution-infrastructure.md#L232) |
| H3 | [Як влаштована мережа](ch18-execution-infrastructure.md#L245) |
| H3 | [Скільки дасть мережа](ch18-execution-infrastructure.md#L281) |
| H2 | [Апаратна база й фізичні обмеження периферії](ch18-execution-infrastructure.md#L301) |
| H2 | [Бюджет затримки й числовий дрейф](ch18-execution-infrastructure.md#L360) |
| H3 | [Бюджет затримки](ch18-execution-infrastructure.md#L362) |
| H3 | [Договір вимірювання й автономності](ch18-execution-infrastructure.md#L380) |
| H3 | [Числовий дрейф](ch18-execution-infrastructure.md#L386) |
| H2 | [Перспективні напрями апаратури для систем знань](ch18-execution-infrastructure.md#L482) |
| H3 | [Концептуальна модель процесора знань](ch18-execution-infrastructure.md#L496) |
| H2 | [Мінімальна автономна конфігурація](ch18-execution-infrastructure.md#L573) |
| H2 | [Висновки](ch18-execution-infrastructure.md#L594) |
| H2 | [Запитання до читачів](ch18-execution-infrastructure.md#L601) |
| H2 | [Словник](ch18-execution-infrastructure.md#L611) |
| H2 | [Абревіатури](ch18-execution-infrastructure.md#L634) |
| H2 | [Джерела](ch18-execution-infrastructure.md#L665) |

</details>

<details>
<summary>Глава 19. Від запитання до доказу: пошук, прив'язка та перевірка твердження</summary>

| Рівень | Заголовок |
|---|---|
| H1 | [Глава 19. Від запитання до доказу: пошук, прив'язка та перевірка твердження](ch19-from-question-to-evidence.md#L1) |
| H2 | [Місце задачі в архітектурі](ch19-from-question-to-evidence.md#L54) |
| H2 | [Чому людське запитання слабо формалізоване](ch19-from-question-to-evidence.md#L71) |
| H2 | [Два часові шляхи здобуття знань](ch19-from-question-to-evidence.md#L106) |
| H2 | [Конвеєр перетворення запиту](ch19-from-question-to-evidence.md#L140) |
| H2 | [Вилучення та детерміноване зв'язування з байтами](ch19-from-question-to-evidence.md#L184) |
| H3 | [Повнота прив'язки](ch19-from-question-to-evidence.md#L216) |
| H3 | [Приклад: прив'язка твердження з RFC 768](ch19-from-question-to-evidence.md#L232) |
| H2 | [Типізовані атоми знань](ch19-from-question-to-evidence.md#L458) |
| H2 | [Походження: від файлу до атома](ch19-from-question-to-evidence.md#L469) |
| H2 | [Проміжне представлення твердження з доказом](ch19-from-question-to-evidence.md#L510) |
| H2 | [Перевірка: чи доказ справді підтримує твердження](ch19-from-question-to-evidence.md#L553) |
| H2 | [Повнота пошуку, безпечна відмова й рецензія](ch19-from-question-to-evidence.md#L561) |
| H2 | [Регресійне тестування конвеєра](ch19-from-question-to-evidence.md#L573) |
| H2 | [Життєвий цикл об'єктів знань: від відповіді до бази](ch19-from-question-to-evidence.md#L610) |
| H2 | [Висновки](ch19-from-question-to-evidence.md#L634) |
| H2 | [Запитання до читачів](ch19-from-question-to-evidence.md#L641) |
| H2 | [Словник](ch19-from-question-to-evidence.md#L650) |
| H2 | [Абревіатури](ch19-from-question-to-evidence.md#L669) |
| H2 | [Джерела](ch19-from-question-to-evidence.md#L689) |

</details>

<details>
<summary>Глава 20. Рушій пояснень: рішення, відмова та межі компетентності</summary>

| Рівень | Заголовок |
|---|---|
| H1 | [Глава 20. Рушій пояснень: рішення, відмова та межі компетентності](ch20-explanation-engine.md#L1) |
| H2 | [Пояснення як окремий інженерний артефакт](ch20-explanation-engine.md#L23) |
| H2 | [П'ять типів пояснювальних запитів](ch20-explanation-engine.md#L83) |
| H2 | [Граф доведення на прикладі](ch20-explanation-engine.md#L126) |
| H2 | [WHY NOT: мінімальний конфлікт і алгоритм QuickXPlain](ch20-explanation-engine.md#L190) |
| H2 | [WHAT IF і WHAT MUST CHANGE: контрфакти й допустимі дії](ch20-explanation-engine.md#L417) |
| H2 | [Захист від витоку через пояснення](ch20-explanation-engine.md#L441) |
| H2 | [Межі компетентності й самоперевірка](ch20-explanation-engine.md#L460) |
| H2 | [Природномовне пояснення поверх перевіреного дерева](ch20-explanation-engine.md#L497) |
| H2 | [Протокол перевірки вірності пояснення](ch20-explanation-engine.md#L557) |
| H2 | [Висновки](ch20-explanation-engine.md#L563) |
| H2 | [Запитання до читачів](ch20-explanation-engine.md#L570) |
| H2 | [Словник](ch20-explanation-engine.md#L579) |
| H2 | [Абревіатури](ch20-explanation-engine.md#L596) |
| H2 | [Джерела](ch20-explanation-engine.md#L609) |

</details>

<details>
<summary>Глава 21. Від рекомендації до дії: контроль повноважень і безпечне виконання у виробничому середовищі</summary>

| Рівень | Заголовок |
|---|---|
| H1 | [Глава 21. Від рекомендації до дії: контроль повноважень і безпечне виконання у виробничому середовищі](ch21-from-recommendation-to-action.md#L1) |
| H2 | [Рекомендація, план і дія: три рівні відповідальності](ch21-from-recommendation-to-action.md#L17) |
| H2 | [П'ять рівнів автономності](ch21-from-recommendation-to-action.md#L62) |
| H2 | [Типізовані контракти дій](ch21-from-recommendation-to-action.md#L76) |
| H2 | [Покрокове виконання плану](ch21-from-recommendation-to-action.md#L125) |
| H3 | [Ієрархічний план і локальне перепланування](ch21-from-recommendation-to-action.md#L160) |
| H2 | [Людина в циклі та втома від погоджень](ch21-from-recommendation-to-action.md#L196) |
| H2 | [Ідемпотентність і компенсації](ch21-from-recommendation-to-action.md#L224) |
| H2 | [Ізоляція: дані ніколи не стають командами](ch21-from-recommendation-to-action.md#L426) |
| H2 | [Засоби реалізації та аналіз журналів дій](ch21-from-recommendation-to-action.md#L467) |
| H2 | [Висновки](ch21-from-recommendation-to-action.md#L482) |
| H2 | [Запитання до читачів](ch21-from-recommendation-to-action.md#L489) |
| H2 | [Словник](ch21-from-recommendation-to-action.md#L499) |
| H2 | [Абревіатури](ch21-from-recommendation-to-action.md#L522) |
| H2 | [Джерела](ch21-from-recommendation-to-action.md#L544) |

</details>

<details>
<summary>Глава 22. Кібернетичний цикл керування: сенсори, периферія та зворотний зв'язок</summary>

| Рівень | Заголовок |
|---|---|
| H1 | [Глава 22. Кібернетичний цикл керування: сенсори, периферія та зворотний зв'язок](ch22-cybernetics-edge-to-backend.md#L1) |
| H2 | [Кібернетика й закон необхідної різноманітності](ch22-cybernetics-edge-to-backend.md#L17) |
| H2 | [Паралельний програмний приклад: черга обробки завдань](ch22-cybernetics-edge-to-backend.md#L83) |
| H2 | [Оцінювання стану: фільтр замість сирого сигналу](ch22-cybernetics-edge-to-backend.md#L99) |
| H2 | [Кібернетика другого порядку, модель життєздатної системи й активний висновок](ch22-cybernetics-edge-to-backend.md#L290) |
| H2 | [Синергетична стійкість: біфуркації, критичне уповільнення (Critical Slowing Down) та передбачення фазових переходів](ch22-cybernetics-edge-to-backend.md#L348) |
| H2 | [Ієрархія часових горизонтів L0–L4](ch22-cybernetics-edge-to-backend.md#L393) |
| H2 | [Контракти сенсорних подій](ch22-cybernetics-edge-to-backend.md#L465) |
| H3 | [Підписувати ознаки, а не вердикти](ch22-cybernetics-edge-to-backend.md#L507) |
| H3 | [Нефункціональні інтерфейси й неявні припущення](ch22-cybernetics-edge-to-backend.md#L521) |
| H2 | [Режими деградації за обриву зв'язку](ch22-cybernetics-edge-to-backend.md#L525) |
| H2 | [Звірка станів після відновлення зв'язку](ch22-cybernetics-edge-to-backend.md#L565) |
| H2 | [Безпечне оновлення знань: підписані пакети](ch22-cybernetics-edge-to-backend.md#L599) |
| H2 | [Сучасні засоби й навчання з часових даних](ch22-cybernetics-edge-to-backend.md#L629) |
| H2 | [Наскрізні приклади](ch22-cybernetics-edge-to-backend.md#L650) |
| H2 | [Висновки](ch22-cybernetics-edge-to-backend.md#L660) |
| H2 | [Архітектура рішення та фізичний зворотний зв'язок](ch22-cybernetics-edge-to-backend.md#L667) |
| H2 | [Подальший шлях пізнання](ch22-cybernetics-edge-to-backend.md#L671) |
| H2 | [Запитання до читачів](ch22-cybernetics-edge-to-backend.md#L675) |
| H2 | [Словник](ch22-cybernetics-edge-to-backend.md#L685) |
| H2 | [Абревіатури](ch22-cybernetics-edge-to-backend.md#L710) |
| H2 | [Джерела](ch22-cybernetics-edge-to-backend.md#L732) |

</details>

<details>
<summary>Глава 23. Верифікація бази знань: як перевірити несуперечливість, повноту та надійність правил</summary>

| Рівень | Заголовок |
|---|---|
| H1 | [Глава 23. Верифікація бази знань: як перевірити несуперечливість, повноту та надійність правил](ch23-knowledge-base-verification.md#L1) |
| H2 | [Верифікація і валідація: правило може бути правильно виконаним, але хибним за змістом](ch23-knowledge-base-verification.md#L17) |
| H2 | [Перевіряють не файл правил, а всю систему рішення](ch23-knowledge-base-verification.md#L53) |
| H2 | [Аномалії бази знань](ch23-knowledge-base-verification.md#L72) |
| H2 | [Від простих перевірок до випробування в реальних умовах](ch23-knowledge-base-verification.md#L83) |
| H2 | [Кількість перевірених правил ще не доводить якості](ch23-knowledge-base-verification.md#L131) |
| H2 | [Інваріанти й контрприклади](ch23-knowledge-base-verification.md#L157) |
| H2 | [Мутаційне тестування: чи помітять тести зламане правило](ch23-knowledge-base-verification.md#L212) |
| H3 | [Скорочення набору зі збереженням контрольних перевірок](ch23-knowledge-base-verification.md#L473) |
| H2 | [Метаморфні й диференційні перевірки](ch23-knowledge-base-verification.md#L492) |
| H2 | [Формальні методи: пошук забороненого стану](ch23-knowledge-base-verification.md#L549) |
| H2 | [Відтворення пояснення разом із рішенням](ch23-knowledge-base-verification.md#L672) |
| H3 | [Маніфест запуску відповіді](ch23-knowledge-base-verification.md#L692) |
| H2 | [Правила разом із машинним навчанням](ch23-knowledge-base-verification.md#L722) |
| H2 | [Перевірка на минулих випадках без самообману](ch23-knowledge-base-verification.md#L742) |
| H2 | [Керовані зміни бази правил](ch23-knowledge-base-verification.md#L754) |
| H2 | [Засоби генерації перевірок і аналізу тестових даних](ch23-knowledge-base-verification.md#L790) |
| H2 | [Хибні висновки, які має зупинити перевірка](ch23-knowledge-base-verification.md#L808) |
| H3 | [Рубрикатор технічного боргу бази знань](ch23-knowledge-base-verification.md#L825) |
| H2 | [З чого почати перевірку правил](ch23-knowledge-base-verification.md#L834) |
| H2 | [Висновки](ch23-knowledge-base-verification.md#L849) |
| H2 | [Запитання до читачів](ch23-knowledge-base-verification.md#L856) |
| H2 | [Словник](ch23-knowledge-base-verification.md#L868) |
| H2 | [Абревіатури](ch23-knowledge-base-verification.md#L896) |
| H2 | [Джерела](ch23-knowledge-base-verification.md#L917) |

</details>

<details>
<summary>Глава 24. Технічна діагностика: як не сплутати симптом із першопричиною в умовах неповноти</summary>

| Рівень | Заголовок |
|---|---|
| H1 | [Глава 24. Технічна діагностика: як не сплутати симптом із першопричиною в умовах неповноти](ch24-system-diagnosis.md#L1) |
| H2 | [Симптом, діагноз і першопричина](ch24-system-diagnosis.md#L17) |
| H2 | [Які дані потрібні для надійного висновку](ch24-system-diagnosis.md#L62) |
| H2 | [Модель допомагає не загубити альтернативи](ch24-system-diagnosis.md#L110) |
| H2 | [Конфлікти звужують коло причин](ch24-system-diagnosis.md#L149) |
| H3 | [Конфлікти між специфікаціями](ch24-system-diagnosis.md#L215) |
| H2 | [Сумісний діагноз ще не найімовірніший](ch24-system-diagnosis.md#L231) |
| H3 | [Коли помиляється сам вимірювальний канал](ch24-system-diagnosis.md#L258) |
| H2 | [Правила, прецеденти й машинне навчання мають різні ролі](ch24-system-diagnosis.md#L290) |
| H2 | [Часова послідовність може змінити пояснення](ch24-system-diagnosis.md#L327) |
| H2 | [Наступна перевірка має розділяти версії](ch24-system-diagnosis.md#L346) |
| H2 | [Результат втручання не доводить причини](ch24-system-diagnosis.md#L691) |
| H2 | [Яке пояснення має отримати інженер](ch24-system-diagnosis.md#L706) |
| H2 | [Коли правильна відповідь «даних недостатньо»](ch24-system-diagnosis.md#L739) |
| H2 | [Як перевірити, чи діагностика справді допомагає](ch24-system-diagnosis.md#L755) |
| H2 | [Що можуть дати сучасні засоби й інтелектуальний аналіз даних](ch24-system-diagnosis.md#L772) |
| H2 | [Помилки діагностики й запобіжники](ch24-system-diagnosis.md#L794) |
| H2 | [З чого почати діагностику в невеликому об'єкті](ch24-system-diagnosis.md#L811) |
| H2 | [Висновки](ch24-system-diagnosis.md#L826) |
| H2 | [Запитання до читачів](ch24-system-diagnosis.md#L833) |
| H2 | [Словник](ch24-system-diagnosis.md#L845) |
| H2 | [Абревіатури](ch24-system-diagnosis.md#L869) |
| H2 | [Джерела](ch24-system-diagnosis.md#L882) |

</details>

<details>
<summary>Глава 25. Як навчати експертну систему: екзаменаційні матриці, аудит знань та контроль регресій</summary>

| Рівень | Заголовок |
|---|---|
| H1 | [Глава 25. Як навчати експертну систему: екзаменаційні матриці, аудит знань та контроль регресій](ch25-how-expert-systems-learn.md#L1) |
| H2 | [Три окремі маршрути зміни](ch25-how-expert-systems-learn.md#L31) |
| H2 | [Життєвий цикл знань: від сигналу до версії](ch25-how-expert-systems-learn.md#L43) |
| H2 | [Наскрізний приклад: прошивка автомобільного блоку керування](ch25-how-expert-systems-learn.md#L119) |
| H3 | [Зміна вимоги й відкликання свідчення: чотири знімки знань](ch25-how-expert-systems-learn.md#L131) |
| H2 | [П'ять шарів, які можна поліпшувати](ch25-how-expert-systems-learn.md#L192) |
| H2 | [Матеріал для навчання: результати роботи проєкту](ch25-how-expert-systems-learn.md#L249) |
| H2 | [Як трасування підтверджує відповідність стандартам](ch25-how-expert-systems-learn.md#L274) |
| H3 | [Як експертна система використовує таблиці TARA, SAR і DFAR](ch25-how-expert-systems-learn.md#L284) |
| H2 | [Іспит важливіший за враження](ch25-how-expert-systems-learn.md#L336) |
| H3 | [Ролі наборів даних](ch25-how-expert-systems-learn.md#L353) |
| H3 | [Екзаменаційна матриця](ch25-how-expert-systems-learn.md#L372) |
| H3 | [Правило допуску: кандидат проти чинної версії](ch25-how-expert-systems-learn.md#L431) |
| H2 | [Калібрування: чи можна довіряти впевненості](ch25-how-expert-systems-learn.md#L641) |
| H2 | [Керована петля зворотного зв'язку](ch25-how-expert-systems-learn.md#L722) |
| H2 | [Коли й чим доналаштовувати мовну модель](ch25-how-expert-systems-learn.md#L793) |
| H3 | [Який метод обрати](ch25-how-expert-systems-learn.md#L799) |
| H3 | [Засоби доналаштування](ch25-how-expert-systems-learn.md#L867) |
| H3 | [Пам'ять рахують до запуску](ch25-how-expert-systems-learn.md#L884) |
| H3 | [Навчання без доступу до інтернету](ch25-how-expert-systems-learn.md#L933) |
| H3 | [Гібридний конвеєр перевірок: хмарне керування й локальні виконавці](ch25-how-expert-systems-learn.md#L950) |
| H2 | [Засоби обліку іспиту й пошуку слабких місць](ch25-how-expert-systems-learn.md#L988) |
| H2 | [Перший керований цикл для команди](ch25-how-expert-systems-learn.md#L1009) |
| H2 | [Висновки](ch25-how-expert-systems-learn.md#L1024) |
| H2 | [Запитання до читачів](ch25-how-expert-systems-learn.md#L1033) |
| H2 | [Словник](ch25-how-expert-systems-learn.md#L1047) |
| H2 | [Абревіатури](ch25-how-expert-systems-learn.md#L1094) |
| H2 | [Джерела](ch25-how-expert-systems-learn.md#L1157) |

</details>

<details>
<summary>Глава 26. Неперервне навчання (Continual Learning) на досвіді та подолання зсуву системного журналу</summary>

| Рівень | Заголовок |
|---|---|
| H1 | [Глава 26. Неперервне навчання (Continual Learning) на досвіді та подолання зсуву системного журналу](ch26-continual-learning.md#L1) |
| H2 | [Як читати главу за об'єктом зміни](ch26-continual-learning.md#L17) |
| H2 | [Що саме може навчатися](ch26-continual-learning.md#L23) |
| H2 | [Як фіксувати навчальний випадок](ch26-continual-learning.md#L72) |
| H2 | [Чому власний журнал не показує всієї картини](ch26-continual-learning.md#L195) |
| H2 | [Як відрізнити нову проблему від зміни даних](ch26-continual-learning.md#L397) |
| H3 | [Випереджальні індикатори: деградація до першої хибної відповіді](ch26-continual-learning.md#L474) |
| H2 | [Як виміряти забування](ch26-continual-learning.md#L505) |
| H2 | [Як додати нове знання й не забути старого](ch26-continual-learning.md#L526) |
| H2 | [Як оновлювати правила й прецеденти](ch26-continual-learning.md#L549) |
| H2 | [Коли просити людину перевірити випадок](ch26-continual-learning.md#L599) |
| H2 | [Чому рекомендації змінюють дані, на яких експертна система вчиться](ch26-continual-learning.md#L630) |
| H2 | [Мовна модель може запропонувати, але не підтвердити](ch26-continual-learning.md#L672) |
| H3 | [Як хибний зворотний зв'язок спотворює навчання](ch26-continual-learning.md#L685) |
| H2 | [Як кандидат доходить до робочої версії](ch26-continual-learning.md#L728) |
| H2 | [Як перевірити, що нова версія справді краща](ch26-continual-learning.md#L772) |
| H2 | [Які помилки треба зупинити до випуску](ch26-continual-learning.md#L790) |
| H2 | [Засоби послідовного навчання й оцінювання політик](ch26-continual-learning.md#L808) |
| H2 | [Строгий аудит і дорадчий режим](ch26-continual-learning.md#L828) |
| H2 | [З чого почати безпечне навчання](ch26-continual-learning.md#L834) |
| H2 | [Висновки](ch26-continual-learning.md#L840) |
| H2 | [Підсумки маршруту перевірки й навчання](ch26-continual-learning.md#L847) |
| H2 | [Подальший шлях пізнання](ch26-continual-learning.md#L851) |
| H2 | [Запитання до читачів](ch26-continual-learning.md#L855) |
| H2 | [Словник](ch26-continual-learning.md#L868) |
| H2 | [Абревіатури](ch26-continual-learning.md#L902) |
| H2 | [Джерела](ch26-continual-learning.md#L927) |

</details>

<details>
<summary>Глава 27. Обґрунтування безпеки: синтез і перевірка аргументів</summary>

| Рівень | Заголовок |
|---|---|
| H1 | [Глава 27. Обґрунтування безпеки: синтез і перевірка аргументів](ch27-safety-case-gsn-synthesis.md#L1) |
| H2 | [Чому обґрунтування безпеки перетворюється на паперову вправу](ch27-safety-case-gsn-synthesis.md#L15) |
| H2 | [Нотація GSN](ch27-safety-case-gsn-synthesis.md#L29) |
| H2 | [Синтез аргументу з інженерного графа знань](ch27-safety-case-gsn-synthesis.md#L79) |
| H2 | [Формальні правила перевірки аргументу](ch27-safety-case-gsn-synthesis.md#L136) |
| H2 | [Цілісність аргументу: дерево Меркла й вибіркове розкриття](ch27-safety-case-gsn-synthesis.md#L201) |
| H2 | [Заперечення як частина аргументу: рамка Зунга](ch27-safety-case-gsn-synthesis.md#L247) |
| H2 | [Програмна перевірка: від правил до доказу включення](ch27-safety-case-gsn-synthesis.md#L278) |
| H2 | [Сучасні засоби й аналіз даних для аргументів гарантування](ch27-safety-case-gsn-synthesis.md#L491) |
| H2 | [Висновки](ch27-safety-case-gsn-synthesis.md#L506) |
| H2 | [Запитання до читачів](ch27-safety-case-gsn-synthesis.md#L513) |
| H2 | [Словник](ch27-safety-case-gsn-synthesis.md#L525) |
| H2 | [Абревіатури](ch27-safety-case-gsn-synthesis.md#L551) |
| H2 | [Джерела](ch27-safety-case-gsn-synthesis.md#L571) |

</details>

<details>
<summary>Глава 28. Дворежимні експертні системи: строгий висновок і дорадча гіпотеза</summary>

| Рівень | Заголовок |
|---|---|
| H1 | [Глава 28. Дворежимні експертні системи: строгий висновок і дорадча гіпотеза](ch28-dual-mode-expert-systems.md#L1) |
| H2 | [Два режими й один інваріант](ch28-dual-mode-expert-systems.md#L15) |
| H2 | [Строге ядро: правила Горна й найменша нерухома точка](ch28-dual-mode-expert-systems.md#L60) |
| H2 | [Три форми міркування Пірса і редукція](ch28-dual-mode-expert-systems.md#L132) |
| H2 | [Дорадчий режим: прецеденти й гіпотези](ch28-dual-mode-expert-systems.md#L158) |
| H2 | [Відмова як джерело знань](ch28-dual-mode-expert-systems.md#L176) |
| H3 | [Ескалація до експерта, коли прецедента немає](ch28-dual-mode-expert-systems.md#L202) |
| H2 | [Програмна перевірка: два режими на одному прикладі](ch28-dual-mode-expert-systems.md#L245) |
| H2 | [Сучасні засоби строгого й дорадчого виведення](ch28-dual-mode-expert-systems.md#L423) |
| H2 | [Порівняння режимів](ch28-dual-mode-expert-systems.md#L435) |
| H2 | [Висновки](ch28-dual-mode-expert-systems.md#L450) |
| H2 | [Запитання до читачів](ch28-dual-mode-expert-systems.md#L457) |
| H2 | [Словник](ch28-dual-mode-expert-systems.md#L470) |
| H2 | [Абревіатури](ch28-dual-mode-expert-systems.md#L494) |
| H2 | [Джерела](ch28-dual-mode-expert-systems.md#L504) |

</details>

<details>
<summary>Глава 29. Нейро-символьна архітектура: мовні моделі та перевірка доказових підстав</summary>

| Рівень | Заголовок |
|---|---|
| H1 | [Глава 29. Нейро-символьна архітектура: мовні моделі та перевірка доказових підстав](ch29-neuro-symbolic-architecture.md#L1) |
| H2 | [Чому правдоподібність не є істинністю](ch29-neuro-symbolic-architecture.md#L15) |
| H2 | [Дві системи: аналогія, а не модель мозку](ch29-neuro-symbolic-architecture.md#L42) |
| H2 | [Модель пропонує, ядро затверджує](ch29-neuro-symbolic-architecture.md#L88) |
| H2 | [Шлюз допуску: байти, а не символи](ch29-neuro-symbolic-architecture.md#L143) |
| H2 | [Мала мовна модель як генератор кандидатів](ch29-neuro-symbolic-architecture.md#L157) |
| H2 | [Практична реалізація на Go з Ollama](ch29-neuro-symbolic-architecture.md#L199) |
| H2 | [Дворежимне виконання](ch29-neuro-symbolic-architecture.md#L455) |
| H3 | [Відновлення структури запиту з перевіркою хостом](ch29-neuro-symbolic-architecture.md#L459) |
| H2 | [Сучасні засоби обмеженого виходу й перевірки](ch29-neuro-symbolic-architecture.md#L471) |
| H2 | [Відкриті напрями](ch29-neuro-symbolic-architecture.md#L484) |
| H2 | [Висновки](ch29-neuro-symbolic-architecture.md#L504) |
| H2 | [Запитання до читачів](ch29-neuro-symbolic-architecture.md#L513) |
| H2 | [Словник](ch29-neuro-symbolic-architecture.md#L525) |
| H2 | [Абревіатури](ch29-neuro-symbolic-architecture.md#L549) |
| H2 | [Джерела](ch29-neuro-symbolic-architecture.md#L562) |

</details>

<details>
<summary>Глава 30. Спільне проєктування функціональної безпеки та кібербезпеки</summary>

| Рівень | Заголовок |
|---|---|
| H1 | [Глава 30. Спільне проєктування функціональної безпеки та кібербезпеки](ch30-safety-cybersecurity-co-engineering.md#L1) |
| H2 | [Анотація](ch30-safety-cybersecurity-co-engineering.md#L13) |
| H2 | [1. Проблема роздільних інженерних культур](ch30-safety-cybersecurity-co-engineering.md#L21) |
| H3 | [1.1. Чотири архетипові міжкатегоріальні колізії](ch30-safety-cybersecurity-co-engineering.md#L53) |
| H2 | [2. Математичний та онтологічний апарат ко-інженерії](ch30-safety-cybersecurity-co-engineering.md#L68) |
| H3 | [2.1. Відображення TARA на HARA](ch30-safety-cybersecurity-co-engineering.md#L111) |
| H3 | [2.2. Часовий бюджет ко-інженерії: FTTI проти ARTI](ch30-safety-cybersecurity-co-engineering.md#L140) |
| H3 | [2.3. Загальна матриця сумісного ризику](ch30-safety-cybersecurity-co-engineering.md#L186) |
| H2 | [3. Обробка та верифікація вимог у форматі ReqIF (ASPICE 4.0)](ch30-safety-cybersecurity-co-engineering.md#L207) |
| H3 | [3.1. Структура ReqIF та нормативні атрибути](ch30-safety-cybersecurity-co-engineering.md#L246) |
| H3 | [3.2. Автоматична перевірка метрик повноти ASPICE 4.0](ch30-safety-cybersecurity-co-engineering.md#L292) |
| H2 | [4. Синтез доказів безпеки за стандартом Goal Structuring Notation (GSN)](ch30-safety-cybersecurity-co-engineering.md#L346) |
| H3 | [4.1. Підписаний запис свідчення](ch30-safety-cybersecurity-co-engineering.md#L390) |
| H2 | [5. Навчальні перевірки в експертній системі мовою Go](ch30-safety-cybersecurity-co-engineering.md#L428) |
| H2 | [6. Практичний кейс: оновлення складського контролера](ch30-safety-cybersecurity-co-engineering.md#L765) |
| H3 | [6.1. Погоджені входи](ch30-safety-cybersecurity-co-engineering.md#L769) |
| H3 | [6.2. Висновок правила](ch30-safety-cybersecurity-co-engineering.md#L780) |
| H3 | [6.3. Альтернативи й перевірки](ch30-safety-cybersecurity-co-engineering.md#L786) |
| H2 | [7. Кваліфікація самої експертної системи як програмного інструмента (ISO 26262-8, розділ 11)](ch30-safety-cybersecurity-co-engineering.md#L798) |
| H3 | [7.1. Вплив інструмента, виявлення помилки й рівень довіри](ch30-safety-cybersecurity-co-engineering.md#L802) |
| H3 | [7.2. Методи кваліфікації й особливість експертної системи](ch30-safety-cybersecurity-co-engineering.md#L832) |
| H3 | [7.3. Комплект документів кваліфікації](ch30-safety-cybersecurity-co-engineering.md#L845) |
| H2 | [8. Проєктний перелік для підготовки перегляду](ch30-safety-cybersecurity-co-engineering.md#L857) |
| H2 | [9. Сучасні засоби й аналіз даних для спільної інженерії](ch30-safety-cybersecurity-co-engineering.md#L870) |
| H2 | [Висновки](ch30-safety-cybersecurity-co-engineering.md#L898) |
| H2 | [Запитання до читачів](ch30-safety-cybersecurity-co-engineering.md#L907) |
| H2 | [Словник](ch30-safety-cybersecurity-co-engineering.md#L921) |
| H2 | [Абревіатури](ch30-safety-cybersecurity-co-engineering.md#L949) |
| H2 | [Джерела](ch30-safety-cybersecurity-co-engineering.md#L987) |

</details>

<details>
<summary>Глава 31. Виведення за нормами: ієрархії предикатів, винятки та чинність</summary>

| Рівень | Заголовок |
|---|---|
| H1 | [Глава 31. Виведення за нормами: ієрархії предикатів, винятки та чинність](ch31-syllogistic-reasoning-and-relation-lattices.md#L1) |
| H2 | [Анотація](ch31-syllogistic-reasoning-and-relation-lattices.md#L13) |
| H2 | [1. Чому однокроковий пошук не дає нормативного висновку](ch31-syllogistic-reasoning-and-relation-lattices.md#L21) |
| H3 | [1.1. Де саме пошук за подібністю втрачає логіку норм](ch31-syllogistic-reasoning-and-relation-lattices.md#L72) |
| H2 | [2. Поняття, судження й висновок](ch31-syllogistic-reasoning-and-relation-lattices.md#L80) |
| H3 | [2.1. Поняття](ch31-syllogistic-reasoning-and-relation-lattices.md#L120) |
| H3 | [2.2. Судження](ch31-syllogistic-reasoning-and-relation-lattices.md#L151) |
| H3 | [2.3. Висновок](ch31-syllogistic-reasoning-and-relation-lattices.md#L168) |
| H2 | [3. Ієрархія предикатів і семантика субсумції](ch31-syllogistic-reasoning-and-relation-lattices.md#L194) |
| H3 | [3.1. Математичні закони спрямованої субсумції](ch31-syllogistic-reasoning-and-relation-lattices.md#L245) |
| H3 | [3.2. Зіставлення фрази з предикатом](ch31-syllogistic-reasoning-and-relation-lattices.md#L279) |
| H3 | [3.3. Розрізнення однойменних документів за лінією чинності](ch31-syllogistic-reasoning-and-relation-lattices.md#L291) |
| H2 | [4. Тризначна логіка Кліні та спростувачі Поллока](ch31-syllogistic-reasoning-and-relation-lattices.md#L304) |
| H3 | [4.1. Таблиця істинності сильної логіки Кліні](ch31-syllogistic-reasoning-and-relation-lattices.md#L321) |
| H3 | [4.2. Порядок інформованості й закрита відмова](ch31-syllogistic-reasoning-and-relation-lattices.md#L335) |
| H3 | [4.3. Спростувачі за Поллоком](ch31-syllogistic-reasoning-and-relation-lattices.md#L354) |
| H2 | [5. Дослідження міждокументних конфліктів: еволюція протоколу TCP](ch31-syllogistic-reasoning-and-relation-lattices.md#L410) |
| H3 | [5.1. Виконувана модель вибору норми мовою ASP](ch31-syllogistic-reasoning-and-relation-lattices.md#L426) |
| H2 | [6. Умовні запитання з кількома кроками](ch31-syllogistic-reasoning-and-relation-lattices.md#L571) |
| H2 | [7. Перегляд сліду виведення](ch31-syllogistic-reasoning-and-relation-lattices.md#L608) |
| H2 | [8. Навчальний крок виведення на Go](ch31-syllogistic-reasoning-and-relation-lattices.md#L632) |
| H2 | [9. Сучасні засоби й видобування знань для ієрархій і винятків](ch31-syllogistic-reasoning-and-relation-lattices.md#L895) |
| H2 | [Висновки](ch31-syllogistic-reasoning-and-relation-lattices.md#L916) |
| H2 | [Запитання до читачів](ch31-syllogistic-reasoning-and-relation-lattices.md#L923) |
| H2 | [Словник](ch31-syllogistic-reasoning-and-relation-lattices.md#L935) |
| H2 | [Абревіатури](ch31-syllogistic-reasoning-and-relation-lattices.md#L961) |
| H2 | [Джерела](ch31-syllogistic-reasoning-and-relation-lattices.md#L983) |

</details>

<details>
<summary>Глава 32. Незмінні пакети знань: побайтовий допуск, індекси та відображення в пам'ять</summary>

| Рівень | Заголовок |
|---|---|
| H1 | [Глава 32. Незмінні пакети знань: побайтовий допуск, індекси та відображення в пам'ять](ch32-high-performance-knowledge-packs-mmap-and-harvesting.md#L1) |
| H2 | [Анотація](ch32-high-performance-knowledge-packs-mmap-and-harvesting.md#L13) |
| H2 | [1. Пакет знань: як база знань стає артефактом випуску](ch32-high-performance-knowledge-packs-mmap-and-harvesting.md#L21) |
| H3 | [1.1. Пакет знань як продукт збирання](ch32-high-performance-knowledge-packs-mmap-and-harvesting.md#L30) |
| H3 | [1.2. Склад пакета: канонічний шар, похідні шари й маніфест](ch32-high-performance-knowledge-packs-mmap-and-harvesting.md#L36) |
| H3 | [1.3. Інваріанти пакета](ch32-high-performance-knowledge-packs-mmap-and-harvesting.md#L92) |
| H3 | [1.4. Життєвий цикл пакета](ch32-high-performance-knowledge-packs-mmap-and-harvesting.md#L105) |
| H3 | [1.5. Два покоління формату](ch32-high-performance-knowledge-packs-mmap-and-harvesting.md#L140) |
| H3 | [1.6. Матеріалізовані подання: що обчислювати під час збирання](ch32-high-performance-knowledge-packs-mmap-and-harvesting.md#L158) |
| H2 | [2. Вибір сховища: порівнювати однакову роботу](ch32-high-performance-knowledge-packs-mmap-and-harvesting.md#L186) |
| H3 | [2.1. Які витрати й гарантії потрібно розділити](ch32-high-performance-knowledge-packs-mmap-and-harvesting.md#L218) |
| H3 | [2.2. Архівні вимірювання та протокол порівняння](ch32-high-performance-knowledge-packs-mmap-and-harvesting.md#L227) |
| H3 | [2.3. Відкритий навчальний стенд SQLite і `mmap`](ch32-high-performance-knowledge-packs-mmap-and-harvesting.md#L248) |
| H2 | [3. Анатомія нульової десеріалізації: системні механізми `mmap(2)`](ch32-high-performance-knowledge-packs-mmap-and-harvesting.md#L360) |
| H3 | [3.1. Механіка сторінкового підкачування та системний виклик `madvise(2)`](ch32-high-performance-knowledge-packs-mmap-and-harvesting.md#L398) |
| H3 | [3.2. Коли відображення файла в пам'ять не підходить](ch32-high-performance-knowledge-packs-mmap-and-harvesting.md#L407) |
| H2 | [4. Специфікація бінарного формату ZNAV-INDEX v2](ch32-high-performance-knowledge-packs-mmap-and-harvesting.md#L415) |
| H3 | [4.1. Що змінилося порівняно з ZNAV-INDEX v1](ch32-high-performance-knowledge-packs-mmap-and-harvesting.md#L446) |
| H2 | [5. Алгоритм безвтратного інвертованого кластерування фактів](ch32-high-performance-knowledge-packs-mmap-and-harvesting.md#L465) |
| H3 | [5.1. Ключ кластера: однозначне кодування полів](ch32-high-performance-knowledge-packs-mmap-and-harvesting.md#L503) |
| H2 | [6. Паралельний стрімінговий конвеєр компіляції та побітова відтворюваність](ch32-high-performance-knowledge-packs-mmap-and-harvesting.md#L657) |
| H3 | [6.1. Алгоритм впорядкованого буфера](ch32-high-performance-knowledge-packs-mmap-and-harvesting.md#L701) |
| H2 | [7. Неперервне нейро-символьне збагачення бази знань](ch32-high-performance-knowledge-packs-mmap-and-harvesting.md#L857) |
| H3 | [7.1. Математичні умови шлюзу допуску](ch32-high-performance-knowledge-packs-mmap-and-harvesting.md#L888) |
| H2 | [8. Інженерія знаннєвої щільності: індекс $\mathrm{KDI}$ та епістемічне профілювання](ch32-high-performance-knowledge-packs-mmap-and-harvesting.md#L910) |
| H3 | [8.1. Емпіричний профіль знаннєвої щільності за галузевими корпусами](ch32-high-performance-knowledge-packs-mmap-and-harvesting.md#L936) |
| H3 | [8.2. Щільність не є повнотою: оцінка методом повторного вилову](ch32-high-performance-knowledge-packs-mmap-and-harvesting.md#L949) |
| H2 | [9. Еталонна реалізація читача mmap-індексу на Go](ch32-high-performance-knowledge-packs-mmap-and-harvesting.md#L972) |
| H2 | [10. Шардування незмінного пакета знань](ch32-high-performance-knowledge-packs-mmap-and-harvesting.md#L1257) |
| H3 | [10.1. Що додається до маніфесту](ch32-high-performance-knowledge-packs-mmap-and-harvesting.md#L1261) |
| H3 | [10.2. Розміщення за найбільшою вагою](ch32-high-performance-knowledge-packs-mmap-and-harvesting.md#L1292) |
| H3 | [10.3. Побудова шардів і перевірка розбиття](ch32-high-performance-knowledge-packs-mmap-and-harvesting.md#L1306) |
| H3 | [10.4. Маршрутизація й чотири підсумки запиту](ch32-high-performance-knowledge-packs-mmap-and-harvesting.md#L1604) |
| H3 | [10.5. Результати й межі](ch32-high-performance-knowledge-packs-mmap-and-harvesting.md#L1879) |
| H2 | [Висновки](ch32-high-performance-knowledge-packs-mmap-and-harvesting.md#L1902) |
| H2 | [Запитання до читачів](ch32-high-performance-knowledge-packs-mmap-and-harvesting.md#L1914) |
| H2 | [Словник](ch32-high-performance-knowledge-packs-mmap-and-harvesting.md#L1933) |
| H2 | [Абревіатури](ch32-high-performance-knowledge-packs-mmap-and-harvesting.md#L1975) |
| H2 | [Джерела](ch32-high-performance-knowledge-packs-mmap-and-harvesting.md#L2012) |

</details>

<details>
<summary>Глава 33. Міжсистемний обмін знаннями: постачання правил стороннім системам, навчання моделей і захищений зворотний зв'язок</summary>

| Рівень | Заголовок |
|---|---|
| H1 | [Глава 33. Міжсистемний обмін знаннями: постачання правил стороннім системам, навчання моделей і захищений зворотний зв'язок](ch33-inter-system-knowledge-exchange-and-model-teaching.md#L1) |
| H2 | [Анотація](ch33-inter-system-knowledge-exchange-and-model-teaching.md#L13) |
| H2 | [1. Експертна система як постачальник знань і вчитель сторонніх систем](ch33-inter-system-knowledge-exchange-and-model-teaching.md#L21) |
| H2 | [2. Постачання знань обробникам фізичних сигналів та сенсорної інформації](ch33-inter-system-knowledge-exchange-and-model-teaching.md#L58) |
| H3 | [2.1. Переклад логічних норм у числові конверти валідності](ch33-inter-system-knowledge-exchange-and-model-teaching.md#L62) |
| H3 | [2.2. Предикатні правила детекції завад та сенсорного спуфінгу](ch33-inter-system-knowledge-exchange-and-model-teaching.md#L76) |
| H2 | [3. Навчання сторонніх систем: символьна дистиляція та верифікований курікулум](ch33-inter-system-knowledge-exchange-and-model-teaching.md#L86) |
| H3 | [3.1. Символьна дистиляція знань через семантичні функції втрат](ch33-inter-system-knowledge-exchange-and-model-teaching.md#L90) |
| H3 | [3.2. Формальні щити безпеки (Shielding) у навчанні агентів](ch33-inter-system-knowledge-exchange-and-model-teaching.md#L108) |
| H2 | [4. Формати, способи та технології міжсистемного обміну знаннями](ch33-inter-system-knowledge-exchange-and-model-teaching.md#L116) |
| H3 | [4.1. Семантичний контракт повідомлення на базі JSON-LD](ch33-inter-system-knowledge-exchange-and-model-teaching.md#L129) |
| H3 | [4.2. Суб'єктна маршрутизація в NATS JetStream](ch33-inter-system-knowledge-exchange-and-model-teaching.md#L159) |
| H2 | [5. Умови, правила збереження семантики та онтологічне узгодження](ch33-inter-system-knowledge-exchange-and-model-teaching.md#L172) |
| H3 | [5.1. Консервативні розширення та локальність модулів](ch33-inter-system-knowledge-exchange-and-model-teaching.md#L176) |
| H3 | [5.2. Правила збереження семантичних типів](ch33-inter-system-knowledge-exchange-and-model-teaching.md#L186) |
| H2 | [6. Розмежування доступу та безпека знань: решітка міток і селективне розкриття](ch33-inter-system-knowledge-exchange-and-model-teaching.md#L194) |
| H3 | [6.1. Багаторівнева безпека на базі решітки міток Дена](ch33-inter-system-knowledge-exchange-and-model-teaching.md#L198) |
| H3 | [6.2. Селективне розкриття доказів через дерева Меркла](ch33-inter-system-knowledge-exchange-and-model-teaching.md#L212) |
| H2 | [7. Підтвердження фактів (Fact Attestation) та захист від отруєння знань](ch33-inter-system-knowledge-exchange-and-model-teaching.md#L228) |
| H3 | [7.1. Криптографічна атестація фактів за схемою Ed25519](ch33-inter-system-knowledge-exchange-and-model-teaching.md#L232) |
| H3 | [7.2. Таксономія загроз: форми атак отруєння знань (Knowledge Poisoning)](ch33-inter-system-knowledge-exchange-and-model-teaching.md#L240) |
| H3 | [7.3. Шлюз вхідного контролю та життєвий цикл карантину знань](ch33-inter-system-knowledge-exchange-and-model-teaching.md#L257) |
| H2 | [8. Зворотний зв'язок (Feedback Loop) та донавчання експертної системи](ch33-inter-system-knowledge-exchange-and-model-teaching.md#L288) |
| H3 | [8.1. Черга контрприкладів та агрегація інцидентів](ch33-inter-system-knowledge-exchange-and-model-teaching.md#L292) |
| H3 | [8.2. Живлення контуру неперервного навчання](ch33-inter-system-knowledge-exchange-and-model-teaching.md#L298) |
| H2 | [9. Програмна реалізація мовою Go: модуль xchange](ch33-inter-system-knowledge-exchange-and-model-teaching.md#L307) |
| H2 | [Висновки](ch33-inter-system-knowledge-exchange-and-model-teaching.md#L755) |
| H2 | [Запитання до читачів](ch33-inter-system-knowledge-exchange-and-model-teaching.md#L767) |
| H2 | [Словник](ch33-inter-system-knowledge-exchange-and-model-teaching.md#L781) |
| H2 | [Абревіатури](ch33-inter-system-knowledge-exchange-and-model-teaching.md#L800) |
| H2 | [Джерела](ch33-inter-system-knowledge-exchange-and-model-teaching.md#L823) |

</details>

<details>
<summary>Глава 34. Прогалини знань: реляційний пошук, абдукція та діалог уточнення</summary>

| Рівень | Заголовок |
|---|---|
| H1 | [Глава 34. Прогалини знань: реляційний пошук, абдукція та діалог уточнення](ch34-deterministic-relational-analysis-abduction-and-socratic-dialogue.md#L1) |
| H2 | [Анотація](ch34-deterministic-relational-analysis-abduction-and-socratic-dialogue.md#L13) |
| H2 | [1. Проблема неповноти знань та епістемічна дилема експертних систем](ch34-deterministic-relational-analysis-abduction-and-socratic-dialogue.md#L25) |
| H2 | [2. Епістемічна тріада класичної логіки в архітектурі системи](ch34-deterministic-relational-analysis-abduction-and-socratic-dialogue.md#L66) |
| H3 | [2.1. Поняття (Concepts / Terms)](ch34-deterministic-relational-analysis-abduction-and-socratic-dialogue.md#L110) |
| H3 | [2.2. Судження (Judgments / Propositions)](ch34-deterministic-relational-analysis-abduction-and-socratic-dialogue.md#L117) |
| H3 | [2.3. Висновки (Inferences / Conclusions)](ch34-deterministic-relational-analysis-abduction-and-socratic-dialogue.md#L124) |
| H2 | [3. Детерміністичний багатоходовий реляційний аналіз](ch34-deterministic-relational-analysis-abduction-and-socratic-dialogue.md#L132) |
| H3 | [3.1. Двонаправлений обмежений пошук у ширину (Bidirectional Bounded BFS)](ch34-deterministic-relational-analysis-abduction-and-socratic-dialogue.md#L136) |
| H3 | [3.2. Композитний ланцюг свідчень (Composite Path Evidence)](ch34-deterministic-relational-analysis-abduction-and-socratic-dialogue.md#L169) |
| H2 | [4. Автономний майнінг асоціативних правил (AMIE PCA)](ch34-deterministic-relational-analysis-abduction-and-socratic-dialogue.md#L179) |
| H3 | [4.1. Формальні правила та припущення PCA](ch34-deterministic-relational-analysis-abduction-and-socratic-dialogue.md#L185) |
| H3 | [4.2. Метрики якості правил](ch34-deterministic-relational-analysis-abduction-and-socratic-dialogue.md#L197) |
| H2 | [5. Символьне абдуктивне виведення робочих гіпотез](ch34-deterministic-relational-analysis-abduction-and-socratic-dialogue.md#L217) |
| H3 | [5.1. Абдукція за Чарльзом Пірсом](ch34-deterministic-relational-analysis-abduction-and-socratic-dialogue.md#L221) |
| H3 | [5.2. Топологічні евристики генерації гіпотез](ch34-deterministic-relational-analysis-abduction-and-socratic-dialogue.md#L251) |
| H3 | [5.3. Інваріант епістемічної гігієни (Hypothesis Isolation Invariant)](ch34-deterministic-relational-analysis-abduction-and-socratic-dialogue.md#L270) |
| H2 | [6. Сократівський діалог та змішана ініціатива](ch34-deterministic-relational-analysis-abduction-and-socratic-dialogue.md#L284) |
| H3 | [6.1. Структура фрейму уточнення (Clarification Frame)](ch34-deterministic-relational-analysis-abduction-and-socratic-dialogue.md#L288) |
| H3 | [6.2. Взаємодія з інтерфейсами людини-оператора](ch34-deterministic-relational-analysis-abduction-and-socratic-dialogue.md#L319) |
| H2 | [7. Крос-доменна генералізація знань](ch34-deterministic-relational-analysis-abduction-and-socratic-dialogue.md#L328) |
| H2 | [8. Програмна реалізація: повний Go-модуль реляційного пошуку, абдукції та сократівського опитувача](ch34-deterministic-relational-analysis-abduction-and-socratic-dialogue.md#L346) |
| H2 | [Висновки](ch34-deterministic-relational-analysis-abduction-and-socratic-dialogue.md#L775) |
| H2 | [Запитання до читачів](ch34-deterministic-relational-analysis-abduction-and-socratic-dialogue.md#L785) |
| H2 | [Словник](ch34-deterministic-relational-analysis-abduction-and-socratic-dialogue.md#L799) |
| H2 | [Абревіатури](ch34-deterministic-relational-analysis-abduction-and-socratic-dialogue.md#L817) |
| H2 | [Джерела](ch34-deterministic-relational-analysis-abduction-and-socratic-dialogue.md#L834) |

</details>

<details>
<summary>Глава 35. Реактивна експертна система: події, відкликання й адаптація знань</summary>

| Рівень | Заголовок |
|---|---|
| H1 | [Глава 35. Реактивна експертна система: події, відкликання й адаптація знань](ch35-self-organizing-expert-systems-synergetics-and-npu-runtime.md#L1) |
| H2 | [1. Обмеження статичних експертних систем: від припущення замкненого світу до живої сенсорики та синергетики знань](ch35-self-organizing-expert-systems-synergetics-and-npu-runtime.md#L11) |
| H3 | [Різниця між простою реактивністю та синергетичною самоорганізацією](ch35-self-organizing-expert-systems-synergetics-and-npu-runtime.md#L33) |
| H2 | [2. Реактивне програмування та потоки подій як епістемічні імпульси](ch35-self-organizing-expert-systems-synergetics-and-npu-runtime.md#L49) |
| H2 | [3. Багаторівнева гібридна пам'ять L0 / L1: незмінний золотий базис та динамічний дельта-граф](ch35-self-organizing-expert-systems-synergetics-and-npu-runtime.md#L74) |
| H3 | [Алгоритм розв'язання фактів (Evidence Resolution Priority):](ch35-self-organizing-expert-systems-synergetics-and-npu-runtime.md#L86) |
| H3 | [3.1. Синергетика динамічної бази знань: самоорганізація, дисипативні структури та принцип підпорядкування](ch35-self-organizing-expert-systems-synergetics-and-npu-runtime.md#L94) |
| H4 | [1. Закон нерівноважного інформаційного припливу (Non-Equilibrium Influx)](ch35-self-organizing-expert-systems-synergetics-and-npu-runtime.md#L106) |
| H4 | [2. Принцип підпорядкування Хакена та семантичні параметри порядку (Slaving Principle & Order Parameters)](ch35-self-organizing-expert-systems-synergetics-and-npu-runtime.md#L109) |
| H4 | [3. Обмежена самоорганізація під доказовими атракторами (Constrained Self-Organization under Truth Attractors)](ch35-self-organizing-expert-systems-synergetics-and-npu-runtime.md#L151) |
| H4 | [4. Апаратний енергетичний метаболізм (Edge Hardware Metabolism)](ch35-self-organizing-expert-systems-synergetics-and-npu-runtime.md#L157) |
| H2 | [4. Спростовне супроводження істинності (Truth Maintenance Systems): динамічні активні дефітери](ch35-self-organizing-expert-systems-synergetics-and-npu-runtime.md#L162) |
| H2 | [5. Синергетичний рушій самоорганізації: мікросервіси Re-Search, Re-Ranking та Re-Thinking](ch35-self-organizing-expert-systems-synergetics-and-npu-runtime.md#L190) |
| H3 | [5.1. Re-Search: еволюційний добір знань та закон необхідної різноманітності Ешбі](ch35-self-organizing-expert-systems-synergetics-and-npu-runtime.md#L217) |
| H3 | [5.2. Re-Ranking: адаптивна реструктуризація решіток переваг](ch35-self-organizing-expert-systems-synergetics-and-npu-runtime.md#L223) |
| H3 | [5.3. Re-Thinking: дисипація протиріч та спростовний резолвінг за AGM](ch35-self-organizing-expert-systems-synergetics-and-npu-runtime.md#L231) |
| H2 | [6. Енергоефективне апаратне прискорення на Edge: NPU, DSP проти серверних GPU](ch35-self-organizing-expert-systems-synergetics-and-npu-runtime.md#L239) |
| H2 | [7. Прецеденти у світовій практиці та індустрії](ch35-self-organizing-expert-systems-synergetics-and-npu-runtime.md#L257) |
| H2 | [8. Програмна реалізація мовою Go: пакет `reactive`](ch35-self-organizing-expert-systems-synergetics-and-npu-runtime.md#L268) |
| H2 | [9. Модульні тести: верифікація реактивного рантайму](ch35-self-organizing-expert-systems-synergetics-and-npu-runtime.md#L566) |
| H2 | [Висновки](ch35-self-organizing-expert-systems-synergetics-and-npu-runtime.md#L676) |
| H2 | [Реактивне виконання в маршруті експлуатації](ch35-self-organizing-expert-systems-synergetics-and-npu-runtime.md#L683) |
| H2 | [Подальший шлях пізнання](ch35-self-organizing-expert-systems-synergetics-and-npu-runtime.md#L687) |
| H2 | [Запитання до читачів](ch35-self-organizing-expert-systems-synergetics-and-npu-runtime.md#L691) |
| H2 | [Словник](ch35-self-organizing-expert-systems-synergetics-and-npu-runtime.md#L703) |
| H2 | [Абревіатури](ch35-self-organizing-expert-systems-synergetics-and-npu-runtime.md#L719) |
| H2 | [Джерела](ch35-self-organizing-expert-systems-synergetics-and-npu-runtime.md#L735) |

</details>

<details>
<summary>Глава 36. Піраміда тестування знань: правила, взаємодії та стійкість відповідей</summary>

| Рівень | Заголовок |
|---|---|
| H1 | [Глава 36. Піраміда тестування знань: правила, взаємодії та стійкість відповідей](ch36-knowledge-testing-pyramid-and-variational-calibration.md#L1) |
| H2 | [1. Методологічний розрив у верифікації систем знань](ch36-knowledge-testing-pyramid-and-variational-calibration.md#L13) |
| H3 | [Порівняльний аналіз світових підходів до тестування та верифікації](ch36-knowledge-testing-pyramid-and-variational-calibration.md#L38) |
| H2 | [2. Концептуальна модель: Чотирирівнева піраміда тестування знань](ch36-knowledge-testing-pyramid-and-variational-calibration.md#L51) |
| H2 | [3. Рівень 1: Knowledge Unit Testing (KUT) — Тестування атомів знань](ch36-knowledge-testing-pyramid-and-variational-calibration.md#L80) |
| H3 | [3.1. Визначення та об'єкт тестування](ch36-knowledge-testing-pyramid-and-variational-calibration.md#L82) |
| H3 | [3.2. Метод фіктивних передумов (Premise Mocking)](ch36-knowledge-testing-pyramid-and-variational-calibration.md#L90) |
| H3 | [3.3. Пастка вакуумної істинності (The Vacuous Truth Trap)](ch36-knowledge-testing-pyramid-and-variational-calibration.md#L103) |
| H3 | [3.4. 6-точковий спектральний аналіз граничних значень (BVA)](ch36-knowledge-testing-pyramid-and-variational-calibration.md#L120) |
| H2 | [4. Рівень 2: Knowledge Integration Testing (KIT) — Решітки виведення та дефітери](ch36-knowledge-testing-pyramid-and-variational-calibration.md#L142) |
| H3 | [4.1. Багатоходові решітки виведення (Inference Lattices)](ch36-knowledge-testing-pyramid-and-variational-calibration.md#L144) |
| H3 | [4.2. Інжекція спростовних дефітерів (Defeater Interruption за AGM)](ch36-knowledge-testing-pyramid-and-variational-calibration.md#L153) |
| H2 | [5. Рівні 3 і 4: Варіативне комплексне калібрування — Метрики SIS та Ліпшицева стійкість](ch36-knowledge-testing-pyramid-and-variational-calibration.md#L169) |
| H3 | [5.1. Метрика семантичної інваріантності (Semantic Invariance Score, SIS)](ch36-knowledge-testing-pyramid-and-variational-calibration.md#L171) |
| H3 | [5.2. Критерій Ліпшицевої стійкості знань ($L_{\mathcal{K}} < \infty$)](ch36-knowledge-testing-pyramid-and-variational-calibration.md#L193) |
| H2 | [6. Стигмергічне закриття прогалин знань (`KnowledgeGapSpool`)](ch36-knowledge-testing-pyramid-and-variational-calibration.md#L206) |
| H2 | [7. Програмна реалізація на мові Go](ch36-knowledge-testing-pyramid-and-variational-calibration.md#L222) |
| H2 | [Висновки](ch36-knowledge-testing-pyramid-and-variational-calibration.md#L592) |
| H2 | [Запитання до читачів](ch36-knowledge-testing-pyramid-and-variational-calibration.md#L601) |
| H2 | [10. Подальший шлях пізнання](ch36-knowledge-testing-pyramid-and-variational-calibration.md#L609) |
| H2 | [Словник](ch36-knowledge-testing-pyramid-and-variational-calibration.md#L618) |
| H2 | [Абревіатури](ch36-knowledge-testing-pyramid-and-variational-calibration.md#L632) |
| H2 | [Джерела](ch36-knowledge-testing-pyramid-and-variational-calibration.md#L651) |

</details>

<details>
<summary>Глава 37. Оцінювання вхідної інформації: джерела, свідчення та невизначеність</summary>

| Рівень | Заголовок |
|---|---|
| H1 | [Глава 37. Оцінювання вхідної інформації: джерела, свідчення та невизначеність](ch37-input-information-assessment-and-algorithmic-skepticism.md#L1) |
| H2 | [Анотація](ch37-input-information-assessment-and-algorithmic-skepticism.md#L11) |
| H2 | [1. Від повідомлення до знання](ch37-input-information-assessment-and-algorithmic-skepticism.md#L17) |
| H2 | [2. Дев'ять аналітичних традицій: метод, свідчення й межі перенесення](ch37-input-information-assessment-and-algorithmic-skepticism.md#L34) |
| H3 | [2.1. Радянська традиція: індикатори й небезпека замкненої картини](ch37-input-information-assessment-and-algorithmic-skepticism.md#L38) |
| H3 | [2.2. Американська традиція: аналіз конкуруючих гіпотез](ch37-input-information-assessment-and-algorithmic-skepticism.md#L44) |
| H3 | [2.3. Британська традиція: дві мови невизначеності](ch37-input-information-assessment-and-algorithmic-skepticism.md#L50) |
| H3 | [2.4. Німецька традиція: перевірка й предметне ущільнення](ch37-input-information-assessment-and-algorithmic-skepticism.md#L56) |
| H3 | [2.5. Французька традиція: цикл запиту й взаємодоповнювальні сенсори](ch37-input-information-assessment-and-algorithmic-skepticism.md#L62) |
| H3 | [2.6. Японська традиція: міжвідомче узагальнення та координація](ch37-input-information-assessment-and-algorithmic-skepticism.md#L68) |
| H3 | [2.7. Китайська традиція: відбір, інтеграція й технологічна обробка](ch37-input-information-assessment-and-algorithmic-skepticism.md#L74) |
| H3 | [2.8. Ізраїльська традиція: критика незмінної концепції](ch37-input-information-assessment-and-algorithmic-skepticism.md#L80) |
| H3 | [2.9. Українська традиція: аналітичне опрацювання та зовнішня перевірка моделей](ch37-input-information-assessment-and-algorithmic-skepticism.md#L86) |
| H2 | [3. Правоохоронна аналітика: зв'язок не дорівнює доказу вини](ch37-input-information-assessment-and-algorithmic-skepticism.md#L94) |
| H3 | [3.1. Національні та міжнаціональні структури](ch37-input-information-assessment-and-algorithmic-skepticism.md#L98) |
| H2 | [4. Кіберзахист: від спрацювання детектора до підтвердженого інциденту](ch37-input-information-assessment-and-algorithmic-skepticism.md#L116) |
| H3 | [4.1. Валідація, верифікація й повторна перевірка](ch37-input-information-assessment-and-algorithmic-skepticism.md#L122) |
| H3 | [4.2. Не плутати достовірність, тяжкість і межі поширення](ch37-input-information-assessment-and-algorithmic-skepticism.md#L137) |
| H2 | [5. Текст, звук, радіохвилі та зображення: контракт первинного спостереження](ch37-input-information-assessment-and-algorithmic-skepticism.md#L145) |
| H2 | [6. Категорії достовірності: багатовісна таксономія](ch37-input-information-assessment-and-algorithmic-skepticism.md#L165) |
| H2 | [7. Ймовірнісна оцінка: сила свідчення, залежність і невизначеність](ch37-input-information-assessment-and-algorithmic-skepticism.md#L196) |
| H3 | [7.1. Від базової частоти до відношення правдоподібностей](ch37-input-information-assessment-and-algorithmic-skepticism.md#L198) |
| H3 | [7.2. Десять копій не створюють десяти підтверджень](ch37-input-information-assessment-and-algorithmic-skepticism.md#L211) |
| H3 | [7.3. Інтервал оцінки замість удаваної точності](ch37-input-information-assessment-and-algorithmic-skepticism.md#L226) |
| H3 | [7.4. Конфлікт і невідома кореляція](ch37-input-information-assessment-and-algorithmic-skepticism.md#L240) |
| H2 | [8. Експертний скепсис як алгоритм](ch37-input-information-assessment-and-algorithmic-skepticism.md#L253) |
| H2 | [9. Реальний час і пакетна обробка без перегляду кожного повідомлення](ch37-input-information-assessment-and-algorithmic-skepticism.md#L277) |
| H3 | [9.1. Послідовне рішення зі скінченним бюджетом](ch37-input-information-assessment-and-algorithmic-skepticism.md#L308) |
| H3 | [9.2. Запізнення, перевантаження й часові мітки](ch37-input-information-assessment-and-algorithmic-skepticism.md#L322) |
| H2 | [10. Аналітичні програми: що вони дають доказовій ЕС](ch37-input-information-assessment-and-algorithmic-skepticism.md#L336) |
| H3 | [10.1. Чого бракує графу для доказового висновку](ch37-input-information-assessment-and-algorithmic-skepticism.md#L360) |
| H3 | [10.2. Стандартизований обмін не означає стандартизованої істинності](ch37-input-information-assessment-and-algorithmic-skepticism.md#L366) |
| H2 | [11. Навчальний алгоритм: інтервал, копії та відкладене рішення](ch37-input-information-assessment-and-algorithmic-skepticism.md#L374) |
| H2 | [Висновок](ch37-input-information-assessment-and-algorithmic-skepticism.md#L589) |
| H2 | [Словник](ch37-input-information-assessment-and-algorithmic-skepticism.md#L597) |
| H2 | [Абревіатури](ch37-input-information-assessment-and-algorithmic-skepticism.md#L611) |
| H2 | [Джерела](ch37-input-information-assessment-and-algorithmic-skepticism.md#L648) |

</details>

<details>
<summary>Глава 38. Машинні галюцинації та дефіцит знань: доказовий контроль відповідей</summary>

| Рівень | Заголовок |
|---|---|
| H1 | [Глава 38. Машинні галюцинації та дефіцит знань: доказовий контроль відповідей](ch38-curing-machine-hallucinations-and-knowledge-deficits.md#L1) |
| H2 | [1. Анатомія машинної галюцинації: чому мовна модель не здатна вилікувати себе сама](ch38-curing-machine-hallucinations-and-knowledge-deficits.md#L11) |
| H2 | [2. Детермінована доказова терапія: двоконтурний шлюз нульових галюцинацій](ch38-curing-machine-hallucinations-and-knowledge-deficits.md#L33) |
| H3 | [Контур 1: Граматично кероване декодування](ch38-curing-machine-hallucinations-and-knowledge-deficits.md#L60) |
| H3 | [Контур 2: Побайтовий шлюз допуску](ch38-curing-machine-hallucinations-and-knowledge-deficits.md#L66) |
| H2 | [3. Епістемічний дефіцит: відкритий світ і шлюз безпечної відмови](ch38-curing-machine-hallucinations-and-knowledge-deficits.md#L89) |
| H2 | [4. Абдуктивні замикання за Пірсом: подолання неповноти без конфабуляцій](ch38-curing-machine-hallucinations-and-knowledge-deficits.md#L113) |
| H2 | [5. Предикатне екранування нейромереж: предикатні щити](ch38-curing-machine-hallucinations-and-knowledge-deficits.md#L141) |
| H2 | [6. Зворотне лікування моделей: навчання та машинне забування на верифікованих знаннях](ch38-curing-machine-hallucinations-and-knowledge-deficits.md#L166) |
| H3 | [Пряма оптимізація переваг за доказовими парами (DPO)](ch38-curing-machine-hallucinations-and-knowledge-deficits.md#L170) |
| H3 | [Машинне забування скомпрометованих фактів (Machine Unlearning)](ch38-curing-machine-hallucinations-and-knowledge-deficits.md#L187) |
| H3 | [Статистична самоузгодженість (Self-Consistency)](ch38-curing-machine-hallucinations-and-knowledge-deficits.md#L193) |
| H2 | [7. Програмна реалізація мовою Go: модуль `antihallucination`](ch38-curing-machine-hallucinations-and-knowledge-deficits.md#L201) |
| H2 | [Висновок](ch38-curing-machine-hallucinations-and-knowledge-deficits.md#L457) |
| H2 | [Словник](ch38-curing-machine-hallucinations-and-knowledge-deficits.md#L470) |
| H2 | [Абревіатури](ch38-curing-machine-hallucinations-and-knowledge-deficits.md#L484) |
| H2 | [Джерела](ch38-curing-machine-hallucinations-and-knowledge-deficits.md#L505) |

</details>
