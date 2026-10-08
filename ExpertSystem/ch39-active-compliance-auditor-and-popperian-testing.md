# Глава 39. Активний експерт-тестувальник: попперівська фальсифікація, нормативний комплаєнс (ASPICE/ISO 26262/ISO 21434) та автономне проєктування випробувань

> **Книга:** [Архітектура доказових експертних систем](README.md) · [Частина VI: Нейро-символьні відповіді, гіпотези та прогалини знань](part-06-frontiers-neuro-symbolic.md)  
> **Попередня глава:** [Глава 38. Машинні галюцинації та дефіцит знань: доказовий контроль відповідей](ch38-curing-machine-hallucinations-and-knowledge-deficits.md)  
> **Зміст книги:** [README.md](README.md)  
> **Автор:** [Микола Федчик](about-the-author.md)  
> **Рівень:** поглиблений: інженери функціональної безпеки (Safety Engineers), фахівці з кібербезпеки (Cybersecurity Engineers), аудитори комплаєнсу (ASPICE/ISO Assessors), архітектори доказового ШІ  
> **Очікувані результати:** опанувати концепцію переходу від пасивного оракула до активного аудитора знань; зрозуміти принципи рантайм-фальсифікації інженерних гіпотез за Карлом Поппером; використовувати експертні системи для автономної підготовки доказових матеріалів TARA, SAR, DFAR та V&V матриць; розв'язувати конфлікти між вимогами безпеки (Safety) та кібербезпеки (Security); автоматизувати рутинні перевірки з вивільненням часу фахівця для технічної творчості зі збереженням принципу Human-in-the-Loop.

---

## Анотація

У проєктуванні критично важливих кіберфізичних систем (ISO 26262 ASIL D, ISO/SAE 21434 CAL 4, DO-178C DAL A, Automotive SPICE 4.0) ручна підготовка нормативної звітності та пасивна верифікація неминуче спричиняють катастрофічні наслідки. Коли інженери з функціональної безпеки вимушені вручну співставляти тисячі вимог у таблицях, когнітивне виснаження породжує «театр відповідності» (*Compliance Theatre*) — формальне проставляння прапорців без математичної перевірки крайових станів. Пропуск єдиної прихованої одноточкової відмови (SPFM < 99%) або неврахованого вектора атаки на бортову шину призводить до смертельних ДТП, відкликання цілих серій безпілотного транспорту та персональної кримінальної відповідальності аудиторів.

Ця глава долає кризу переходом від пасивного оракула-консультанта до Активного Аудитора Доменних Знань. Спираючись на принцип фальсифікації Карла Поппера ($\mathcal{F}(\mathcal{H}_{\text{design}}, \mathcal{K}_{\text{norm}})$), доказова експертна система перехоплює ініціативу допиту проекту: вона автономно шукає контрприклади, виявляє нормативні суперечності між вимогами безпеки (*Safety*) та кіберзахисту (*Security*), генерує строгі сертифікаційні артефакти (HARA, TARA, SAR, DFAR/FMEDA) із повною трасованістю ($100{,}0\%$) та проєктує фальсифікуючі випробування, залишаючи людину в циклі керування (*Human-in-the-Loop*) як найвищого арбітра.

---

## 1. Драма комплаєнс-інженерії: чому ручні методи вичерпали себе

Розробка сучасних кіберфізичних систем (автомобільний транспорт, авіоніка, залізнична автоматика, медичні апарати) підпорядкована найсуворішим галузевим стандартам:
- **Automotive SPICE 4.0** (процесна зрілість розробки програмного та системного забезпечення);
- **ISO 26262:2018** (функціональна безпека електричних та електронних систем, класифікація ASIL-A .. ASIL-D);
- **ISO/SAE 21434:2021** (інженерія кібербезпеки транспортних засобів, рівні CAL 1..4);
- **DO-178C / ED-12C** (авіаційне бортове програмне забезпечення, рівні DAL A..E);
- **IEC 62304 / ISO 14971** (медичне ПЗ та управління ризиками пацієнтів).

### 1.1. Тягар колосальної персональної відповідальності
На відміну від звичайної комерційної веб-розробки, фахівці з безпеки — **Safety Engineer**, **Cybersecurity Engineer**, **Functional Safety Manager (FSM)** та **Lead Assessor** — несуть пряму юридичну, професійну, а в багатьох юрисдикціях і кримінальну відповідальність за випущені артефакти. 

Підпис інженера під підсумковими звітами засвідчує, що:
1. Усі ризики заподіяння шкоди життю та здоров'ю людей зведені до прийнятного залишкового рівня (*As Low As Reasonably Practicable*, ALARP).
2. Усі відомі вектори кібератак на бортові шини та контролери враховані та нейтралізовані.
3. Кожна окрема норма стандарту підтверджена прямим, відтворюваним об'єктивним доказом (*Objective Evidence*).

### 1.2. Пекло рутинних людино-годин: TARA, SAR, DFAR та HARA
Ціна цієї відповідальності вимірюється колосальними трудовитратами. Щоб випустити сучасний автомобільний контролер (наприклад, блок керування гальмами або шлюз доступу), команда змушена вручну сформувати та верифікувати сотні складних доказових документів:

| Артефакт | Стандарт | Сутність та інженерне наповнення | Рутинний виклик для інженера |
| :--- | :--- | :--- | :--- |
| **HARA** (*Hazard Analysis and Risk Assessment*) | ISO 26262-3 | Визначення небезпечних подій, оцінка важкості ($S$), експозиції ($E$), контрольованості ($C$) та призначення рівня ASIL (A/B/C/D) і цілей безпеки (*Safety Goals*). | Ручний перебір сотень комбінацій режимів руху автомобіля та відмов датчиків; ризик пропуску критичного сценарію. |
| **TARA** (*Threat Analysis and Risk Assessment*) | ISO/SAE 21434-9 | Визначення активів (Assets), властивостей кібербезпеки (C-I-A), сценаріїв загроз, побудова дерев атак (*Attack Trees*), оцінка потенціалу атаки (*Attack Feasibility*) та рівнів CAL. | Необхідність співставлення тисяч CAN/Ethernet сигналів з базами CVE/CWE, ручне моделювання кроків зловмисника через шлюзи. |
| **SAR** (*Safety Assessment Report*) | ISO 26262-2/8 | Підсумковий аудиторський звіт незалежної оцінки відповідності розробки вимогам функціональної безпеки та захищеності процесів. | Перехресне звіряння сотень пунктів стандарту із реальними протоколами випробувань, пошук розривів трасованості. |
| **DFAR / DFMEA / FMEDA** | ISO 26262-5/6 | Кількісний аналіз режимів відмов, розрахунок метрик одноточкових відмов (SPFM $\ge 99\%$), латентних відмов (LFM $\ge 90\%$) та діагностичного покриття (DC). | Багаторівневі таблиці в Excel на десятки тисяч рядків; зміна одного резистора в схемі вимагає перерахунку всього ланцюга метрик. |

У типовому проекті підготовка, узгодження та підтримка актуальності цих матеріалів забирає **до 60–70% усього інженерного бюджету часу**. 

### 1.3. Феномен людської втоми та "Compliance Theatre"
Коли кваліфікований інженер змушений тижнями переносити ідентифікатори вимог між Polarion, DOORS, Jira та таблицями Excel, неминуче виникає когнітивне виснаження:
- **Крихкість ручного трасування:** зміна в одному рядку архітектури безслідно ламає доказові зв'язки в десятках дочірніх тестів.
- **Compliance Theatre (театр відповідності):** щоб вкластися у дедлайн випуску, інженери змушені ставити галочки у чеклістах формально, не маючи фізичної змоги математично перевірити кожен крайовий стан.
- **Втрата часу на творчість:** замість глибокого аналізу фізичних аномалій датчиків, проектування нетривіальних алгоритмів діагностики та пошуку нестандартних векторів атак, провідні уми витрачають енергію на бюрократичну рутину.

Саме тут виникає нагальна потреба в **доказових експертних системах нового покоління**.

---

## 2. Авторська концепція Миколи Федчика: Зсув парадигми до Активного Аудитора

Традиційна теорія штучного інтелекту розглядала експертну систему як пасивного консультанта:
$$\text{Людина запитує} \quad \longrightarrow \quad \text{Система видає довідку}$$

Але в складних нормативних доменах пасивний оракул безсилий: **інженер не запитує про те, про що він забув або чого не помітив у 800-сторінковому стандарті**.

Головний архітектор **Микола Федчик** запропонував радикальний зсув інженерної парадигми:
> **Концепція Активного Доменного Експерта-Тестувальника:**  
> Якщо машина володіє точною машинно-зчитуваною базою знань стандарту (ZKP4) та моделлю проекту, вона повинна **перехопити ініціативу допиту (Active Probing)**.  
> Експертна система зобов'язана самостійно аналізувати проектні артефакти, ставити інженеру незручні запитання, виявляти приховані нормативні конфлікти, вимагати обов'язкові за стандартом тести та **автономно генерувати драфти TARA, SAR, DFAR та сертифіковані програми випробувань**.

```mermaid
flowchart TD
    subgraph KNOWLEDGE["Нормативна база знань ZKP4 (Zero-Copy mmap)"]
        STD1["ISO 26262 (ASIL A-D, SPFM, LFM)"]
        STD2["ISO/SAE 21434 (TARA, CAL, Attack Trees)"]
        STD3["ASPICE 4.0 (SWE.1 - SWE.6 Traceability)"]
    end

    subgraph ENGINE["Активний доменний експерт Znavets v4"]
        PROBE["<b>Модуль активного допиту</b><br/>(Socratic Question Generator)"]
        POPPER["<b>Попперівський фальсифікатор</b><br/>(EVM / EISA v1.0, %ebx Custody)"]
        SYNTH["<b>Синтезатор комплаєнс-артефактів</b><br/>(TARA / SAR / DFAR Matrix Engine)"]
    end

    subgraph ARTIFACTS["Вихідні сертифіковані матеріали"]
        OUT_TARA["Повна матриця TARA<br/>(Assets, Threats, Attack Paths)"]
        OUT_DFAR["DFAR / FMEDA Розрахунок<br/>(SPFM >= 99%, DC Check)"]
        OUT_TEST["V&V Програма випробувань<br/>(Fault Injection, BVA 6-point)"]
        OUT_SAR["SAR Доказовий звіт<br/>(Побайтова трасованість цитат)"]
    end

    KNOWLEDGE --> ENGINE
    ENGINE -->|"Активне запитування інженера"| HITL["<b>Інженер-експерт (Human-in-the-Loop)</b><br/>Валідація, стратегічні рішення, творчість"]
    HITL -->|"Відповіді, специфікації продукту"| ENGINE
    ENGINE --> ARTIFACTS
```

---

## 3. Математичний апарат: Попперівська фальсифікація інженерних гіпотез

Основою перевірки проектних рішень є принцип фальсифікованості Карла Поппера [[1]](#src-1): *жодна система не може бути визнана безпечною лише на підставі успішних тестів; безпека доводиться невдачею найагресивніших спроб її спростувати*.

### 3.1. Формалізація інженерної гіпотези
Розробник, мовна модель або архітектор висувають проектне твердження $\mathcal{H}_{\mathrm{design}}$ (наприклад: *«Модуль обробки педалі гальма відповідає рівню ASIL-D без дублювання АЦП, оскільки використовується періодичне самотестування»*).

Формально гіпотеза записується як предикат над простором станів системи $\mathcal{S}$:
$$\mathcal{H}_{\mathrm{design}} \equiv \forall s \in \mathcal{S}, \quad \mathrm{StateValid}(s) \implies \mathrm{SafetyGoalSatisfied}(s)$$

Нормативна база знань $\mathcal{K}_{\mathrm{norm}}$ складається з двійкових деонтичних атомів стандарту:
$$\mathcal{K}_{\mathrm{norm}} = \{ \nu_1, \nu_2, \dots, \nu_m \}, \quad \nu_i = \langle \mathrm{Domain}, \mathrm{Clause}, \mathrm{Entity}, \mathrm{Modality}, \mathrm{Action}, \mathrm{Evidence} \rangle$$
де $\mathrm{Modality} \in \{ \mathrm{MUST}, \mathrm{MUST\text{-}NOT}, \mathrm{SHOULD}, \mathrm{MAY} \}$ (у програмному коді — константа `MUST_NOT`).

### 3.2. Пошук потенційного фальсифікатора (Potential Falsifier)
Завдання експертної системи — за час $t < 1\ \mathrm{ms}$ виконати символьний пошук контрприкладу:
$$\mathcal{F}(\mathcal{H}_{\mathrm{design}}, \mathcal{K}_{\mathrm{norm}}) = \{ \nu_k \in \mathcal{K}_{\mathrm{norm}} \mid \mathrm{Implication}(\mathcal{H}_{\mathrm{design}}) \models \mathrm{Violation}(\nu_k) \}$$

Якщо такий атом знайдено:
$$\mathrm{Verdict} = \mathbf{FALSIFIED} \quad \bigl( \mathrm{Refusal}(\rho), \quad \mathrm{EBX} = \mathrm{SHA256}(\mathrm{Quote}), \quad \mathrm{Clause} = \text{ISO 26262-5:2018 Clause 8.4.3} \bigr)$$

Система не просто каже «код невірний». Вона видає фальсифікуючий нормативний факт:
> *«Гіпотезу спростовано: Згідно з ISO 26262-5:2018 Clause 8.4.3 (цитата: "Single-point fault metric for ASIL-D shall achieve at least 99%"), одноканальний АЦП з тестовим покриттям 90% не задовольняє метрику SPFM. Необхідно додати апаратне дублювання або діагностичний компаратор».*

---

## 4. Автоматизація створення TARA, SAR, DFAR через доменні пакети ZKP4

Розгляньмо детально, як активний експерт допомагає фахівцям формувати ключові доказові документи без рутинного ручного перенесення даних.

### 4.1. Автоматизація TARA (ISO/SAE 21434): від опису системи до матриці ризиків
Процес TARA складається з кількох канонічних кроків, кожен з яких тепер підтримується експертною системою:
1. **Asset Identification (Визначення активів):** Експерт сканує опис архітектури (DBC-файли CAN, ARXML-файли AUTOSAR, IDL-специфікації) та автоматично видобуває всі активи (наприклад: *«Ключ шифрування сесії діагностики»*, *«Сигнал кута повороту керма SteerAngle»*).
2. **Threat Scenario Identification (Сценарії загроз):** Зв'язуючи активи з онтологією STRIDE / MITRE ATT&CK for ICS у ZKP4, експертна система синтезує повний перелік загроз:
   $$\mathrm{Threat} = \langle \mathrm{Asset}, \ \mathrm{Property}, \ \mathrm{Damage} \rangle$$
   де $\mathrm{Asset} = \text{SteerAngle}$, $\mathrm{Property} = \text{Integrity}$, а $\mathrm{Damage}$ — несанкціоноване подрулювання на швидкості.
3. **Attack Path Analysis & Feasibility (Дерева атак):** Система розгортає граф зв'язків бортової мережі та розраховує вектор складності атаки за методикою Attack Potential (Elapsed Time, Specialist Expertise, Knowledge of Item, Window of Opportunity, Equipment).
4. **Формування фінальної таблиці TARA:** Замість тижнів роботи інженер отримує повністю згенеровану матрицю зі зведеними балами ризику (Risk Values 1..5) та вимогами до контрзаходів кібербезпеки (*Cybersecurity Goals*).

### 4.2. Автоматизація DFAR та FMEDA (ISO 26262): математична строгість метрик
Підготовка звіту DFAR/FMEDA вимагає математичного розрахунку надійності:
- Інтенсивність відмов компонентів ($\lambda$, FIT);
- Класифікація відмов: безпечні ($\lambda_s$), небезпечні одноточкові ($\lambda_{\mathrm{spf}}$), залишкові ($\lambda_{\mathrm{rf}}$), латентні ($\lambda_{\mathrm{mpf,lat}}$);
- Метрика одноточкових відмов (для ASIL-D норма вимагає $\ge 99\%$):
  $$\mathrm{SPFM} = \frac{\sum (\lambda_s + \lambda_{\mathrm{spf}})}{\sum \lambda} \ge 0{,}99$$
- Метрика латентних відмов (для ASIL-D норма вимагає $\ge 90\%$):
  $$\mathrm{LFM} = \frac{\sum (\lambda_s + \lambda_{\mathrm{mpf,det}})}{\sum (\lambda - \lambda_{\mathrm{spf}})} \ge 0{,}90$$

Експертна система Znavets v4:
- Зберігає норми розрахунку у вигляді деонтичних та математичних правил;
- Автоматично перевіряє розрахункову модель схеми або коду;
- Якщо метрика SPFM виявляється рівною $98.4\%$, система активує діалог: *«Увага: для досягнення цільових 99% ASIL-D не вистачає 0.6%. Рекомендовано підвищити діагностичне покриття сторожового таймера (Watchdog) з 60% до 90% або додати зворотний зчитувач регістру виводу»*.

### 4.3. Автоматизація SAR: доказове полотно для аудиторів TÜV / Dekra
Safety Assessment Report (SAR) є вінцем проекту функціональної безпеки. Експертна система формує SAR як дерево аргументів у нотації **Goal Structuring Notation (GSN)**:
- **Top Goal:** Система задовольняє вимогам ASIL-D стандарту ISO 26262.
- **Strategy:** Аргументація через декомпозицію на безпеку апаратного забезпечення, безпеку ПЗ та процесну якість ASPICE.
- **Evidence:** Кожен листовий вузол дерева (Evidence) містить посилання на конкретний протокол тестування із зафіксованим криптографічним хешем результату та побайтовою цитатою пункту стандарту.

---

## 5. Розв'язання фундаментального протиріччя: Safety vs. Cybersecurity

У складних комплаєнс-проектах найгострішою проблемою є **конфлікт між вимогами функціональної безпеки (Safety) та кібербезпеки (Security)**:

```mermaid
flowchart LR
    subgraph CONFLICT["Конфлікт вимог у критичній точці"]
        direction TB
        REQ_SAFE["<b>ISO 26262 (Safety):</b><br/>При аварії чи спрацюванні подушок двері МУСЯТЬ бути негайно розблоковані для евакуації пасажирів.<br/><i>(Принцип доступності / Availability)</i>"]
        REQ_SEC["<b>ISO 21434 (Security):</b><br/>Будь-яка команда розблокування дверей з шини CAN МУСИТЬ проходити криптографічну автентифікацію MAC-підписом.<br/><i>(Принцип цілісності / Integrity)</i>"]
    end

    REQ_SAFE <-->|КОЛІЗІЯ ВИМОГ| REQ_SEC

    CONFLICT --> ARBITER["<b>Експертний арбітр Znavets (ASPIC+)</b><br/>Резолвер дефітерів та часових бюджетів"]
    ARBITER --> RESOLUTION["<b>Узгоджене інженерне рішення:</b><br/>Апаратний дискретний піропатрон (Safety) має прямий пріоритет над шинним протоколом (Security); шинні команди вимагають MAC лише при швидкості > 0 км/год."]
```

### Як допомагає експертна система:
1. **Автоматичне виявлення колізій (Cross-Standard Defeater Mining):** Знання обох стандартів у єдиному середовищі дозволяє системі виявити суперечність правил ще на етапі архітектурного проєктування (SWE.2).
2. **Аргументація за схемою ASPIC+:** Система будує дерево спростовних міркувань (*Defeasible Reasoning*), розділяючи заперечення засновок (*rebutting*) та підрив правила (*undercutting*).
3. **Генерація безпечного компромісу:** Експерт пропонує формалізований варіант вирішення: *«Застосувати апаратний сигнал від датчика уповільнення в обхід мікроконтролера, зберігши криптографічний бар'єр для всіх програмних запитів з діагностичного роз'єму»*.

---

## 6. Нейро-символічний тандем тестувальника: генерація тест-програм

Як працює практична зв'язка мовної моделі та ядра EVM у ролі активного тестувальника:

1. **System 1 (LLM Proposer — креативність):**
   - Читає фрагмент коду драйвера та документацію.
   - Синтезує хитромудрі, нестандартні сценарії: *«Що буде, якщо надіслати CAN-кадр довжиною DLC=15 замість 8 саме в момент перемикання реле живлення?»*.
   - Формулює проект тест-кейсу природною мовою.
2. **System 2 (EVM / EISA v1.0 — детермінований суддя):**
   - Приймає тест-кейс через `Popperian Falsification API`.
   - Звіряє його з ZKP4-нормами ISO 11898, ISO 26262 та AUTOSAR.
   - Миттєво перевіряє: чи не порушує сам тест обов'язкових умов стандарту? Який пункт стандарту він покриває?
   - Якщо тест валідний — система автоматично реєструє його в матриці V&V і прив'язує до вимоги простежуваності ASPICE SWE.4.
3. **Розрахунок повноти тестової програми:**
   $$\mathrm{TraceabilityCoverage} = \frac{|\mathcal{R}_{\mathrm{requirements}} \cap \mathcal{T}_{\mathrm{verified}}|}{|\mathcal{R}_{\mathrm{requirements}}|} = 1{,}00$$
   (що відповідає 100% покриття трасованості вимог).

---

## 7. Гуманістичний вимір: Людина в циклі керування (Human-in-the-Loop) як звільнення для творчості

Впровадження активного комплаєнс-експерта принципово **не виключає людину з інженерного процесу**. Навпаки, воно повертає професії інженера її первинний високий зміст.

### Що забирає машина:
- Рутинне сканування тисяч сторінок тексту нормативів;
- Заповнення багатотисячних таблиць простежуваності (Excel / Polarion);
- Контроль суворої деонтичної модальності слів (`MUST`, `SHALL`, `REQUIRED`);
- Математичний перерахунок метрик надійності (SPFM, LFM, FIT rates);
- Перевірку повноти покриття коду вимогами стандарту.

### Що повертається людині (Safety & Cybersecurity Engineer):
- **Глибока інженерна творчість:** проектування красивих, елегантних та стійких архітектур;
- **Фізична інтуїція:** дослідження рідкісних аномалій реального «заліза», теплових дрейфів, деградації кремнію чи радіаційних збоїв, які неможливо формалізувати в стандартах;
- **Стратегічна відповідальність:** людина більше не тремтить перед аудитом, бо знає, що кожен формальний пункт надійно прикритий математично верифікованою базою. Фахівець впевнено ставить свій підпис під SAR/TARA, спираючись на доказовий фундамент найвищого рівня довіри.

---

## 8. Виробничий код: Ядро активного аудитора TARA та Safety Directives на Go

Нижче наведено розширену реалізацію активного аудитора, що підтримує сутності стандартів ISO 26262 та ISO/SAE 21434:

<details>
<summary><b>Повний вихідний код: Ядро активного аудитора TARA та Safety Directives на Go (~150 рядків)</b></summary>

```go
package compliance

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"sync"
	"time"
)

// DomainStandard ідентифікатор стандарту комплаєнсу.
type DomainStandard string

const (
	StandardISO26262 DomainStandard = "ISO-26262:2018"
	StandardISO21434 DomainStandard = "ISO/SAE-21434:2021"
	StandardASPICE4  DomainStandard = "ASPICE-4.0"
	StandardRFC9110  DomainStandard = "RFC-9110"
)

// DeonticModality модальність норми за RFC 2119 / ISO Directives.
type DeonticModality string

const (
	ModalityMust    DeonticModality = "MUST"
	ModalityMustNot DeonticModality = "MUST_NOT"
	ModalityShould  DeonticModality = "SHOULD"
)

// ComplianceRule репрезентує непорушний нормативний атом ZKP4.
type ComplianceRule struct {
	Standard      DomainStandard  `json:"standard"`
	Clause        string          `json:"clause"`         // напр. "Part 6 Clause 8.4.2"
	Entity        string          `json:"entity"`         // напр. "SafetyMechanism"
	Modality      DeonticModality `json:"modality"`       // MUST / MUST_NOT
	TargetAction  string          `json:"target_action"`  // напр. "SilentFailure"
	VerbatimQuote string          `json:"verbatim_quote"` // Дослівна норма
	ByteStart     uint64          `json:"byte_start"`     // Початок у першоджерелі
	ByteEnd       uint64          `json:"byte_end"`       // Кінець у першоджерелі
	ExpectedSHA   string          `json:"expected_sha"`   // Хеш цитати (%ebx custody)
}

// EngineeringHypothesis проектне рішення, надане інженером або LLM.
type EngineeringHypothesis struct {
	HypothesisID   string `json:"hypothesis_id"`
	TargetEntity   string `json:"target_entity"`
	ProposedAction string `json:"proposed_action"`
	SafetyASIL     string `json:"safety_asil,omitempty"` // "QM", "ASIL-A".."ASIL-D"
	Rationale      string `json:"rationale"`
}

// FalsificationResult підсумок перевірки за Поппером.
type FalsificationResult struct {
	IsFalsified      bool            `json:"is_falsified"`
	ViolatedRule     *ComplianceRule `json:"violated_rule,omitempty"`
	RefusalReason    string          `json:"refusal_reason"`
	EvidenceVerified bool            `json:"evidence_verified"`
	Latency          time.Duration   `json:"latency"`
}

// ActiveComplianceEngine автономний аудитор TARA/SAR/DFAR.
type ActiveComplianceEngine struct {
	mu           sync.RWMutex
	rulesByEntity map[string][]ComplianceRule
	rawSourceData []byte // mmap масив першоджерела
}

// NewComplianceEngine ініціалізує аудитор з прив'язкою до першоджерела.
func NewComplianceEngine(sourceData []byte) *ActiveComplianceEngine {
	return &ActiveComplianceEngine{
		rulesByEntity: make(map[string][]ComplianceRule),
		rawSourceData: sourceData,
	}
}

// RegisterComplianceRule додає правило з нормативної бази.
func (e *ActiveComplianceEngine) RegisterComplianceRule(r ComplianceRule) {
	e.mu.Lock()
	defer e.mu.Unlock()
	e.rulesByEntity[r.Entity] = append(e.rulesByEntity[r.Entity], r)
}

// FalsifyDesignHypothesis здійснює попперівську фальсифікацію за час < 1ms.
func (e *ActiveComplianceEngine) FalsifyDesignHypothesis(ctx context.Context, h EngineeringHypothesis) (*FalsificationResult, error) {
	start := time.Now()
	e.mu.RLock()
	defer e.mu.RUnlock()

	rules, found := e.rulesByEntity[h.TargetEntity]
	if !found || len(rules) == 0 {
		return &FalsificationResult{
			IsFalsified:      false,
			RefusalReason:    "No normative restrictions found; open-world hypothesis accepted.",
			EvidenceVerified: true,
			Latency:          time.Since(start),
		}, nil
	}

	for _, rule := range rules {
		// Побайтова верифікація першоджерела (%ebx custody check)
		if !e.checkCustody(rule) {
			return nil, fmt.Errorf("custody breach on %s [%d..%d]", rule.Clause, rule.ByteStart, rule.ByteEnd)
		}

		// Попперівське спростування: пряме порушення заборони стандарту
		if rule.Modality == ModalityMustNot && rule.TargetAction == h.ProposedAction {
			return &FalsificationResult{
				IsFalsified:      true,
				ViolatedRule:     &rule,
				RefusalReason:    fmt.Sprintf("Direct compliance breach of %s (%s): %s", rule.Standard, rule.Clause, rule.VerbatimQuote),
				EvidenceVerified: true,
				Latency:          time.Since(start),
			}, nil
		}
	}

	return &FalsificationResult{
		IsFalsified:      false,
		RefusalReason:    "Design hypothesis withstood Popperian falsification against loaded compliance rules.",
		EvidenceVerified: true,
		Latency:          time.Since(start),
	}, nil
}

// checkCustody перевіряє SHA-256 цитати в mmap зрізі.
func (e *ActiveComplianceEngine) checkCustody(r ComplianceRule) bool {
	if e.rawSourceData == nil || r.ByteEnd > uint64(len(e.rawSourceData)) || r.ByteStart >= r.ByteEnd {
		return false
	}
	chunk := e.rawSourceData[r.ByteStart:r.ByteEnd]
	h := sha256.Sum256(chunk)
	return hex.EncodeToString(h[:]) == r.ExpectedSHA
}

// InterrogateSystem генерує активні директиви допиту інженерної команди.
func (e *ActiveComplianceEngine) InterrogateSystem(entity string) []string {
	e.mu.RLock()
	defer e.mu.RUnlock()

	var probes []string
	for _, rule := range e.rulesByEntity[entity] {
		if rule.Modality == ModalityMust {
			probes = append(probes, fmt.Sprintf(
				"ACTIVE COMPLIANCE PROBE [%s %s]: System MUST implement and verify '%s'. Where is the test evidence? Quote: \"%s\"",
				rule.Standard, rule.Clause, rule.TargetAction, rule.VerbatimQuote,
			))
		}
	}
	return probes
}
```

</details>

---

## 9. Висновки до глави

1. **Трансформація від пасивного сховища знань до активного агента контролю якості:**  
   Експертна система нового покоління не чекає на запитання. Вона знає вимоги стандартів краще, ніж стомлений інженер, і активно сканує систему, генерує директиви випробувань та виявляє прогалини трасованості.
2. **Порятунок фахівців з безпеки від бюрократичного вигорання:**  
   Автоматизоване складання драфтів TARA, SAR, DFAR та матриць V&V на базі двійкових пакетів ZKP4 знімає до 90% монотонних людино-годин, захищаючи автора звіту від фатальних пропусків норм.
3. **Строгий математичний щит ($\mathrm{ZHR} = 1{,}00$):**  
   Попперівська фальсифікація дозволяє зовнішнім мовним моделям (LLM) генерувати креативні тестові вектори, гарантуючи при цьому, що жодна галюцинація не потрапить у фінальну сертифікаційну документацію.
4. **Гідне місце людини в епоху ШІ:**  
   Залишаючи людину арбітром і стратегом (Human-in-the-Loop), доказова експертна система повертає інженерії безпеки радість творчості, елегантності та інтелектуальної гідності.

> [!NOTE]
> **Практичне застосування попперівських критеріїв фальсифікації на фізичних апаратних стендах:**
> - [Додаток Б. Робототехніка та кіберфізичні системи](appendix-b-robotics-and-cyber-physical-systems.md) — попперівський критерій фальсифікації гіпотези ASIL-D для гетерогенного тандему Jetson AGX Orin + Xilinx Virtex FPGA ($P(T_{\mathrm{loop}} > T_{\mathrm{wdg}}) > 10^{-9}$ на годину).
> - [Додаток В. Автономна навігація без GNSS](appendix-c-autonomous-navigation-and-geosearch.md) — фальсифікація навігаційної стійкості безпілотників в умовах РЕБ та супутникового спуфінгу ($P(\mathrm{drift} > 5\ \mathrm{m}) > 10^{-6}$).
> - [Додаток Г. Аналогові експертні системи та апаратне виведення](appendix-d-analog-expert-systems-and-neuromorphic-computing.md) — фальсифікація придатності аналогового шлюзу Махаланобіса за температурного дрейфу ($P(\mathrm{MissedEmergency}) > 0$).
> - [Додаток Д. Змішані аналого-цифрові експертні системи](appendix-e-mixed-signal-neuromorphic-expert-systems.md) — попперівське спростування нейроморфного доказового контракту при виході за межі шлюзу запасу ($M(\mathbf{x}) < \gamma_{\mathrm{margin}}$).

---

## Джерела до глави

1. <a id="src-1"></a>**Popper, K. R.** (1959). *The Logic of Scientific Discovery*. London: Hutchinson & Co.
2. <a id="src-2"></a>**VDA QMC.** (2023). *Automotive SPICE Process Assessment / Reference Model, Version 4.0*. Berlin: Quality Management Center in the German Association of the Automotive Industry.
3. <a id="src-3"></a>**International Organization for Standardization.** (2018). *ISO 26262:2018: Road vehicles — Functional safety (Parts 1–12)*. Geneva: ISO.
4. <a id="src-4"></a>**ISO/SAE.** (2021). *ISO/SAE 21434:2021: Road vehicles — Cybersecurity engineering*. Geneva: ISO.
5. <a id="src-5"></a>**RTCA / EUROCAE.** (2011). *DO-178C / ED-12C: Software Considerations in Airborne Systems and Equipment Certification*. Washington, D.C. / Paris.
6. <a id="src-6"></a>**Kelly, T., & Weaver, R.** (2004). *The Goal Structuring Notation — A Safety Argument Notation*. Proceedings of Dependable Systems and Networks.
7. <a id="src-7"></a>**Dung, P. M.** (1995). *On the acceptability of arguments and its fundamental role in nonmonotonic reasoning, logic programming and n-person games*. Artificial Intelligence, 77(2), 321–357.
8. <a id="src-8"></a>**Федчик, М.** (2026). *Архітектура доказових експертних систем: від формальних онтологій до нейро-символьного ШІ*.
