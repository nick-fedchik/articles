# Реферативний реєстр авторського наукового доробку та епістемічної чистоти монографії «Архітектура доказових експертних систем»

**Автор монографії:** Микола Федчик  
**Призначення документа:** Систематизація, наукова класифікація та верифікація академічної доброчесності (відсутності плагіату й чіткого розмежування першоджерел) усіх авторських теорій, теорем, наукових законів, гіпотез, концепцій, правил, принципів, постулатів, аксіом, математичних моделей, інженерних інваріантів та винаходів у монографії (глави 1–40).

**Редакція 2026-10-28:** за рішенням автора додатки та їхні окремі реєстри вилучено з активної української редакції. Збережені положення обліковано в реєстрах глав 16, 18, 36 і 37; решта матеріалу не переноситься й не додається до поточного доробку. Таблиці первинного аудиту нижче є історичною базовою лінією сорока глав, а не повторною перевіркою нової редакції. Їхні відсотки відображають класифікацію записів, не доводять оригінальності, правильності теорем або відсутності плагіату.

---

## 1. Методологія та двовимірна таксономія класифікації

Кожна наукова одиниця або твердження у поглавних розборах (`registry/chXX-registry.md`) класифікується за трьома осями:

### 1.1. Вісь А: Наукова форма твердження або винаходу (`SCIENTIFIC_FORM`)

| Код категорії | Українська назва | Критерій віднесення |
| :--- | :--- | :--- |
| **`AXIOM` / `POSTULATE`** | Аксіома / Постулат | Вихідне бездоказове положення або фундаментальне обмеження теорії доказових експертних систем, на якому будуються подальші формальні викладки. |
| **`LAW` / `REGULARITY`** | Науковий закон / Закономірність | Об'єктивний, стійкий, математично або емпірично виражений зв'язок між параметрами системи (наприклад, між ентропією, повнотою контексту, пропускною здатністю шини чи ймовірністю хибного допуску). |
| **`THEOREM` / `LEMMA`** | Теорема / Лема / Наслідок | Строго доведене (або таке, що прямо випливає з формальної побудови в тексті) математичне твердження про властивості решіток відношень, графів знань, бітемпоральних зрізів чи меж часу виконання (`WCET`). |
| **`PRINCIPLE`** | Науковий / Інженерний принцип | Керівна норма побудови архітектури системи, порушення якої призводить до втрати доказовості або безпеки (`Fail-Closed`, збереження походження тощо). |
| **`CONCEPT` / `THEORY`** | Концепція / Парадигма / Теорія | Цілісна система поглядів, теоретична рамка або нова методологічна парадигма (наприклад, *Епістемічний контракт*, *Дворежимна нейро-символьна система*). |
| **`HYPOTHESIS`** | Наукова гіпотеза | Обґрунтоване наукове припущення щодо поведінки складних гібридних/нейроморфних/апаратних систем, сформульоване для експериментальної перевірки та попперівської фальсифікації. |
| **`MATH_MODEL` / `METRIC`** | Формальна математична модель / Оператор / Метрика | Конкретна розрахункова формула, функціонал, булевий предикат або кількісна міра, інтегрована у прийняття рішень рантайму (*Actionable Closed Loop*). |
| **`INVARIANT`** | Інженерний інваріант | Непорушне співвідношення параметрів стану або апаратних ресурсів, яке має зберігатися на кожному такті/транзакції. |
| **`RULE`** | Правило / Евристика / Регламент | Прикладне правило сайзингу, маршрутизації, екстракції знань, калібрування порогів або верифікації. |
| **`INVENTION` / `ALGORITHM`** | Архітектурний винахід / Патерн / Протокол / Алгоритм | Оригінальна конструкція обчислювального тракту, структури даних у пам'яті (`Zero-Copy mmap`), апаратного вузла (`EPU`, аналоговий шлюз Махаланобіса) або протоколу взаємодії. |
| **`EMPIRICAL_FINDING`** | Емпіричний результат / Дослідницький стенд | Авторський експериментальний стенд, методика випробувань або кількісний результат бенчмаркінгу. |

### 1.2. Вісь Б: Рівень оригінальності та перевірка відсутності плагіату (`ORIGINALITY_LEVEL`)

| Код рівня | Статус оригінальності та чистоти | Що перевіряється та фіксується в реєстрі |
| :--- | :--- | :--- |
| **`[L1: PURE_ORIGINAL]`** | Чиста авторська розробка / винахід | Концепція, формула, теорема, закон, правило чи архітектура повністю розроблені автором монографії (Миколою Федчиком). |
| **`[L2: AUTHORIAL_EXTENSION]`** | Авторська формалізація / синтез на базі класики | Базується на відомій науковій теорії (Байєс, Демпстер–Шефер, Кліні $`\mathbb{K}_3`$, Белнап $`\mathcal{B}_4`$, Гаєк, Поппер, Флоріді, Данґ, ASPIC+, Z3, ISO 26262 тощо), але містить **оригінальну інженерну або математичну надбудову автора**. У реєстрі чітко фіксується межа: *«Класичне ядро (автор, рік) vs. Внесок М. Федчика»*. |
| **`[L3: IMPLICIT_UNMARKED]`** | Неявний (прихований) авторський доробок | Оригінальна авторська думка, формула, класифікація або правило вже є в тексті глави, але викладені безособово без формальної назви («Закон...», «Принцип...», «Теорема...») та без вказівки на авторство. Потребує номінації та явної атрибуції. |
| **`[L4: CLASSICAL_NEEDS_CITATION]`** | Класичний базис (контроль плагіату / цитування) | Загальновідомий науковий факт, класична теорема або стандарт, який не є авторським винаходом. Перевіряється наявність коректного посилання на першоджерело у тексті та списку `## Джерела`, щоб виключити будь-яку двозначність щодо авторства. |

### 1.3. Вісь В: Поточний стан відображення у книзі (`BOOK_STATUS`)

- **`[EXPLICIT_COMPLETE]`** — Вже явно описано в тексті та у блоці `[!NOTE]` з чіткою авторською атрибуцією; потребує лише фіксації у загальній науковій номенклатурі.
- **`[EXPLICIT_NEEDS_CLASSIFICATION]`** — Присутнє в тексті та/або в `[!NOTE]`, але сформульоване описово; потребує надання чіткого статусу (*Теорема*, *Закон*, *Аксіома*, *Постулат*, *Принцип*, *Правило*).
- **`[IMPLICIT_NEEDS_ATTRIBUTION]`** — Присутнє у тексті чи формулах безособово; потребує явного маркування авторства та/або додавання у відповідний розділ.
- **`[PROPOSED_INSERTION]`** — Логічно випливає зі змісту розділу або замикає математичну модель, але ще не сформульоване явно в тексті; рекомендовано до вставки після узгодження з автором.

---

## 2. Зведена навігаційна матриця монографії (Скелет реєстру за канонічною структурою Частин I–VII)

### Частина I. Концептуальний та епістемічний фундамент ([part-01-foundations.md](file:///home/mfedchyk/book/ExpertSystem/part-01-foundations.md))

| Глава | Файл глави | Файл поглавного реєстру | Статус аудиту | Кількість виявлених авторських положень (L1–L3) | Перевірка цитувань (L4) |
| :--- | :--- | :--- | :--- | :--- | :--- |
| **Глава 01** | [ch01-introduction-to-expert-systems.md](file:///home/mfedchyk/book/ExpertSystem/ch01-introduction-to-expert-systems.md) | [registry/ch01-registry.md](file:///home/mfedchyk/book/ExpertSystem/registry/ch01-registry.md) | ✅ Завершено | 33 положення (L1: 23, L2: 10; `100% EXPLICIT_COMPLETE`, 6 канонічних блоків `[!NOTE]`, L4: 4, разом: 37) | 21/21 підтверджено (`100% PASS`; додано Post, Robinson, Forgy, Minsky, DO-178C, ASPIC+) |
| **Глава 02** | [ch02-epistemology-of-machine-knowledge.md](file:///home/mfedchyk/book/ExpertSystem/ch02-epistemology-of-machine-knowledge.md) | [registry/ch02-registry.md](file:///home/mfedchyk/book/ExpertSystem/registry/ch02-registry.md) | ✅ Завершено | 72 положення (L1: 51, L2: 21; `100% EXPLICIT_COMPLETE`, 11 канонічних блоків `[!NOTE]`, L4: 10, разом: 82) | 34/34 підтверджено (`100% PASS`; додано Floridi, Bell-LaPadula, прив'язано Geifman) |
| **Глава 03** | [ch03-beyond-reference-information-systems.md](file:///home/mfedchyk/book/ExpertSystem/ch03-beyond-reference-information-systems.md) | [registry/ch03-registry.md](file:///home/mfedchyk/book/ExpertSystem/registry/ch03-registry.md) | ✅ Завершено | 38 положень (L1: 18, L2: 11, L3: 9; `100% EXPLICIT_COMPLETE`, 10 канонічних блоків `[!NOTE]`, L4: 0, разом: 38) | 16/16 підтверджено (`100% PASS`; додано Deutsch, Meyer, de Moura & Bjørner Z3) |
| **Глава 04** | [ch04-evolution-from-bayes-to-evidence-ai.md](file:///home/mfedchyk/book/ExpertSystem/ch04-evolution-from-bayes-to-evidence-ai.md) | [registry/ch04-registry.md](file:///home/mfedchyk/book/ExpertSystem/registry/ch04-registry.md) | ✅ Завершено | 40 положень (L1: 24, L2: 16; `100% EXPLICIT_COMPLETE`, 6 канонічних блоків `[!NOTE]`, L4: 4, разом: 44) | 31/31 підтверджено (`100% PASS`; додано Quillian 1968, Forgy 1982) |
| **Глава 05** | [ch05-triad-of-trust-and-corporate-memory.md](file:///home/mfedchyk/book/ExpertSystem/ch05-triad-of-trust-and-corporate-memory.md) | [registry/ch05-registry.md](file:///home/mfedchyk/book/ExpertSystem/registry/ch05-registry.md) | ✅ Завершено | 40 положень (L1: 30, L2: 10; `100% EXPLICIT_COMPLETE`, 7 канонічних блоків `[!NOTE]`/`[!WARNING]`, L4: 0, разом: 40) | 18/18 підтверджено (`100% PASS`; додано Feigenbaum 1984, Ed25519, FIPS 140-3, ISO 27001) |

> **Підсумок аудиту Частини I (Глави 1–5):**  
> Проаналізовано **241 наукову одиницю**, з яких **223 одиниці (92,5 %)** становлять авторський науковий та інженерний доробок Миколи Федчика (`L1: 146`, `L2/L3: 77`), а **18 одиниць (7,5 %)** — класичний науковий фундамент (`L4`).

### Частина II. Математичні моделі, подання та зберігання знань ([part-02-knowledge-models.md](file:///home/mfedchyk/book/ExpertSystem/part-02-knowledge-models.md))

| Глава | Файл глави | Файл поглавного реєстру | Статус аудиту | Кількість виявлених авторських положень (L1–L3) | Перевірка цитувань (L4) |
| :--- | :--- | :--- | :--- | :--- | :--- |
| **Глава 06** | [ch06-applied-mathematics-for-expert-systems.md](file:///home/mfedchyk/book/ExpertSystem/ch06-applied-mathematics-for-expert-systems.md) | [registry/ch06-registry.md](file:///home/mfedchyk/book/ExpertSystem/registry/ch06-registry.md) | ✅ Завершено | 54 положення (L1: 36, L2: 18; `100% EXPLICIT_COMPLETE`, 8 канонічних блоків `[!NOTE]`, L4: 44, разом: 98) | 52/52 підтверджено (`100% PASS`; узгоджено якорі Peirce, Douven та Aamodt & Plaza) |
| **Глава 07** | [ch07-knowledge-base-typology.md](file:///home/mfedchyk/book/ExpertSystem/ch07-knowledge-base-typology.md) | [registry/ch07-registry.md](file:///home/mfedchyk/book/ExpertSystem/registry/ch07-registry.md) | ✅ Завершено | 81 положення (L1: 68, L2: 13; `100% EXPLICIT_COMPLETE`, 5 канонічних блоків `[!NOTE]`, L4: 11, разом: 92) | 38/38 підтверджено (`100% PASS`; додано Lenat, Minsky, Z3, EARS, GBNF) |
| **Глава 08** | [ch08-engineering-artifacts-as-data.md](file:///home/mfedchyk/book/ExpertSystem/ch08-engineering-artifacts-as-data.md) | [registry/ch08-registry.md](file:///home/mfedchyk/book/ExpertSystem/registry/ch08-registry.md) | ✅ Завершено | 73 положення (L1: 59, L2: 14; `100% EXPLICIT_COMPLETE`, 5 канонічних блоків `[!NOTE]`, L4: 2, разом: 75) | 24/24 підтверджено (`100% PASS`; додано Parnas 1995, Brooks 1987, UCUM) |
| **Глава 09** | [ch09-engineering-knowledge-graph-traceability.md](file:///home/mfedchyk/book/ExpertSystem/ch09-engineering-knowledge-graph-traceability.md) | [registry/ch09-registry.md](file:///home/mfedchyk/book/ExpertSystem/registry/ch09-registry.md) | ✅ Завершено | 38 положень (L1: 33, L2: 5; `100% EXPLICIT_COMPLETE`, 5 канонічних блоків `[!NOTE]`, L4: 4, разом: 42) | 12/12 підтверджено (`100% PASS`; прив'язано інлайн-посилання Knowledge Vault; чітко демарковано Закон деградації простежуваності Федчика на основі ймовірнісної моделі Клеланд-Гуанг) |
| **Глава 32** | [ch32-high-performance-knowledge-packs-mmap-and-harvesting.md](file:///home/mfedchyk/book/ExpertSystem/ch32-high-performance-knowledge-packs-mmap-and-harvesting.md) | [registry/ch32-registry.md](file:///home/mfedchyk/book/ExpertSystem/registry/ch32-registry.md) | ✅ Завершено | 55 положень (L1: 43, L2: 12; `100% EXPLICIT_COMPLETE`, 10 канонічних блоків `[!NOTE]`, L4: 2, разом: 57) | 32/32 підтверджено (`100% PASS`; додано FIPS 204 ML-DSA, RFC 8554 LMS) |

> **Підсумок аудиту Частини II (Глави 6–9, 32):**  
> Проаналізовано **364 наукові одиниці**, з яких **301 одиниця (82,7 %)** становлять авторський науковий та інженерний доробок Миколи Федчика (`L1: 239`, `L2: 62`), а **63 одиниці (17,3 %)** — класичний науковий фундамент (`L4`).

### Частина III. Оцінювання входу, здобуття знань та інженерія онтологій ([part-03-knowledge-engineering-nlp.md](file:///home/mfedchyk/book/ExpertSystem/part-03-knowledge-engineering-nlp.md))

| Глава | Файл глави | Файл поглавного реєстру | Статус аудиту | Кількість виявлених авторських положень (L1–L3) | Перевірка цитувань (L4) |
| :--- | :--- | :--- | :--- | :--- | :--- |
| **Глава 37** | [ch37-input-information-assessment-and-algorithmic-skepticism.md](file:///home/mfedchyk/book/ExpertSystem/ch37-input-information-assessment-and-algorithmic-skepticism.md) | [registry/ch37-registry.md](file:///home/mfedchyk/book/ExpertSystem/registry/ch37-registry.md) | ✅ Завершено | 65 положень (L1: 47, L2: 18; `100% EXPLICIT_COMPLETE`, 6 канонічних блоків `[!NOTE]`, L4: 23, разом: 88) | 38/38 підтверджено (`100% PASS`; прив'язано [22] Guo, [25] Akidau, додано Schum 1994 [38]) |
| **Глава 10** | [ch10-knowledge-acquisition-systems.md](file:///home/mfedchyk/book/ExpertSystem/ch10-knowledge-acquisition-systems.md) | [registry/ch10-registry.md](file:///home/mfedchyk/book/ExpertSystem/registry/ch10-registry.md) | ✅ Завершено | 87 положень (L1: 69, L2: 18; `100% EXPLICIT_COMPLETE`, 6 канонічних блоків `[!NOTE]`, L4: 4, разом: 91) | 39/39 підтверджено (`100% PASS`; додано MITRE ATLAS, RFC 2119) |
| **Глава 11** | [ch11-knowledge-elicitation-from-experts.md](file:///home/mfedchyk/book/ExpertSystem/ch11-knowledge-elicitation-from-experts.md) | [registry/ch11-registry.md](file:///home/mfedchyk/book/ExpertSystem/registry/ch11-registry.md) | ✅ Завершено | 66 положень (L1: 53, L2: 13; `100% EXPLICIT_COMPLETE`, 6 канонічних блоків `[!NOTE]`, L4: 22, разом: 88) | 27/27 підтверджено (`100% PASS`; додано Kelly 1955, Brier 1950, Klein 1998, Vapnik) |
| **Глава 12** | [ch12-linguistic-analysis-and-local-models.md](file:///home/mfedchyk/book/ExpertSystem/ch12-linguistic-analysis-and-local-models.md) | [registry/ch12-registry.md](file:///home/mfedchyk/book/ExpertSystem/registry/ch12-registry.md) | ✅ Завершено | 46 положень (L1: 32, L2: 14; `100% EXPLICIT_COMPLETE`, 9 канонічних блоків `[!NOTE]`, L4: 1, разом: 47) | 43/43 підтверджено (`100% PASS`; додано Chomsky, Searle, Tarski, UCUM) |
| **Глава 14** | [ch14-requirements-detection-and-formalization.md](file:///home/mfedchyk/book/ExpertSystem/ch14-requirements-detection-and-formalization.md) | [registry/ch14-registry.md](file:///home/mfedchyk/book/ExpertSystem/registry/ch14-registry.md) | ✅ Завершено | 55 положень (L1: 42, L2: 13; `100% EXPLICIT_COMPLETE`, 7 канонічних блоків `[!NOTE]`, L4: 5, разом: 60) | 21/21 підтверджено (`100% PASS`; прив'язано GBNF, Willard; додано Zave & Jackson) |
| **Глава 15** | [ch15-knowledge-extraction-and-kb-construction.md](file:///home/mfedchyk/book/ExpertSystem/ch15-knowledge-extraction-and-kb-construction.md) | [registry/ch15-registry.md](file:///home/mfedchyk/book/ExpertSystem/registry/ch15-registry.md) | ✅ Завершено | 75 положень (L1: 62, L2: 13; `100% EXPLICIT_COMPLETE`, 8 канонічних блоків `[!NOTE]`, L4: 0, разом: 75) | 30/30 підтверджено (`100% PASS`; додано HermiT, WordNet, Etzioni) |

> **Підсумок аудиту Частини III (Глави 37, 10–12, 14–15):**  
> Проаналізовано **449 наукових одиниць**, з яких **394 одиниці (87,8 %)** становлять авторський науковий та інженерний доробок Миколи Федчика (`L1: 305`, `L2: 89`), а **55 одиниць (12,2 %)** — класичний науковий фундамент (`L4`).

### Частина IV. Архітектура ядра, детерміноване виведення, діагностика та пояснення ([part-04-architecture-and-inference.md](file:///home/mfedchyk/book/ExpertSystem/part-04-architecture-and-inference.md))

| Глава | Файл глави | Файл поглавного реєстру | Статус аудиту | Кількість виявлених авторських положень (L1–L3) | Перевірка цитувань (L4) |
| :--- | :--- | :--- | :--- | :--- | :--- |
| **Глава 16** | [ch16-expert-systems-architecture.md](file:///home/mfedchyk/book/ExpertSystem/ch16-expert-systems-architecture.md) | [registry/ch16-registry.md](file:///home/mfedchyk/book/ExpertSystem/registry/ch16-registry.md) | ✅ Завершено | 71 положення (L1: 54, L2: 17; `100% EXPLICIT_COMPLETE`, 7 канонічних блоків `[!NOTE]`, L4: 4, разом: 75) | 27/27 підтверджено (`100% PASS`; додано Fowler CQRS, McKenney RCU) |
| **Глава 13** | [ch13-language-variability-vs-determinism.md](file:///home/mfedchyk/book/ExpertSystem/ch13-language-variability-vs-determinism.md) | [registry/ch13-registry.md](file:///home/mfedchyk/book/ExpertSystem/registry/ch13-registry.md) | ✅ Завершено | 42 положення (L1: 30, L2: 12; `100% EXPLICIT_COMPLETE`, 5 канонічних блоків `[!NOTE]`, L4: 3, разом: 45) | 34/34 підтверджено (`100% PASS`; додано Frege, Montague, Wittgenstein, Roth) |
| **Глава 19** | [ch19-from-question-to-evidence.md](file:///home/mfedchyk/book/ExpertSystem/ch19-from-question-to-evidence.md) | [registry/ch19-registry.md](file:///home/mfedchyk/book/ExpertSystem/registry/ch19-registry.md) | ✅ Завершено | 63 положення (L1: 51, L2: 12; `100% EXPLICIT_COMPLETE`, 6 канонічних блоків `[!NOTE]`, L4: 3, разом: 66) | 14/14 підтверджено (`100% PASS`; додано Liang HELM, FIPS 180-4) |
| **Глава 31** | [ch31-syllogistic-reasoning-and-relation-lattices.md](file:///home/mfedchyk/book/ExpertSystem/ch31-syllogistic-reasoning-and-relation-lattices.md) | [registry/ch31-registry.md](file:///home/mfedchyk/book/ExpertSystem/registry/ch31-registry.md) | ✅ Завершено | 62 положення (L1: 47, L2: 15; `100% EXPLICIT_COMPLETE`, 7 канонічних блоків `[!NOTE]`, L4: 9, разом: 71) | 19/19 підтверджено (`100% PASS`; прив'язано ISO 26262, додано Ganter & Wille, Tarski) |
| **Глава 24** | [ch24-system-diagnosis.md](file:///home/mfedchyk/book/ExpertSystem/ch24-system-diagnosis.md) | [registry/ch24-registry.md](file:///home/mfedchyk/book/ExpertSystem/registry/ch24-registry.md) | ✅ Завершено | 64 положення (L1: 47, L2: 17; `100% EXPLICIT_COMPLETE`, 8 канонічних блоків `[!NOTE]`, L4: 16, разом: 80) | 22/22 підтверджено (`100% PASS`; 0 залишків стоп-слова) |
| **Глава 20** | [ch20-explanation-engine.md](file:///home/mfedchyk/book/ExpertSystem/ch20-explanation-engine.md) | [registry/ch20-registry.md](file:///home/mfedchyk/book/ExpertSystem/registry/ch20-registry.md) | ✅ Завершено | 53 положення (L1: 44, L2: 9; `100% EXPLICIT_COMPLETE`, 7 канонічних блоків `[!NOTE]`, L4: 4, разом: 57) | 13/13 підтверджено (`100% PASS`; додано Swartout, Pearl, Reiter, прив'язано JSAKE 2026) |

> **Підсумок аудиту Частини IV (Глави 16, 13, 19, 31, 24, 20):**  
> Проаналізовано **394 наукові одиниці**, з яких **355 одиниць (90,1 %)** становлять авторський науковий та інженерний доробок Миколи Федчика (`L1: 273`, `L2: 82`), а **39 одиниць (9,9 %)** — класичний науковий фундамент (`L4`).

### Частина V. Верифікація, випробування, контрольована еволюція та гарантування безпеки ([part-05-verification-and-learning.md](file:///home/mfedchyk/book/ExpertSystem/part-05-verification-and-learning.md))

| Глава | Файл глави | Файл поглавного реєстру | Статус аудиту | Кількість виявлених авторських положень (L1–L3) | Перевірка цитувань (L4) |
| :--- | :--- | :--- | :--- | :--- | :--- |
| **Глава 23** | [ch23-knowledge-base-verification.md](file:///home/mfedchyk/book/ExpertSystem/ch23-knowledge-base-verification.md) | [registry/ch23-registry.md](file:///home/mfedchyk/book/ExpertSystem/registry/ch23-registry.md) | ✅ Завершено | 81 положення (L1: 71, L2: 10; `100% EXPLICIT_COMPLETE`, 9 канонічних блоків `[!NOTE]`, L4: 25, разом: 106) | 25/25 підтверджено (`100% PASS`; 0 залишків стоп-слова) |
| **Глава 36** | [ch36-knowledge-testing-pyramid-and-variational-calibration.md](file:///home/mfedchyk/book/ExpertSystem/ch36-knowledge-testing-pyramid-and-variational-calibration.md) | [registry/ch36-registry.md](file:///home/mfedchyk/book/ExpertSystem/registry/ch36-registry.md) | ✅ Завершено | 48 положень (L1: 38, L2: 10; `100% EXPLICIT_COMPLETE`, 9 канонічних блоків `[!NOTE]`, L4: 1, разом: 49) | 16/16 підтверджено (`100% PASS`; прив'язано інлайн-якорі [14] Clopper-Pearson, [15] Hanley, [16] Firecracker) |
| **Глава 25** | [ch25-how-expert-systems-learn.md](file:///home/mfedchyk/book/ExpertSystem/ch25-how-expert-systems-learn.md) | [registry/ch25-registry.md](file:///home/mfedchyk/book/ExpertSystem/registry/ch25-registry.md) | ✅ Завершено | 74 положення (L1: 57, L2: 17; `100% EXPLICIT_COMPLETE`, 8 канонічних блоків `[!NOTE]`, L4: 8, разом: 82) | 50/50 підтверджено (`100% PASS`; додано Compton RDR, Mitchell, Valiant PAC) |
| **Глава 26** | [ch26-continual-learning.md](file:///home/mfedchyk/book/ExpertSystem/ch26-continual-learning.md) | [registry/ch26-registry.md](file:///home/mfedchyk/book/ExpertSystem/registry/ch26-registry.md) | ✅ Завершено | 59 положень (L1: 50, L2: 9; `100% EXPLICIT_COMPLETE`, 6 канонічних блоків `[!NOTE]`, L4: 26, разом: 85) | 27/27 підтверджено (`100% PASS`; додано McCloskey-Cohen, Grossberg ART, Pearl) |
| **Глава 39** | [ch39-active-compliance-auditor-and-popperian-testing.md](file:///home/mfedchyk/book/ExpertSystem/ch39-active-compliance-auditor-and-popperian-testing.md) | [registry/ch39-registry.md](file:///home/mfedchyk/book/ExpertSystem/registry/ch39-registry.md) | ✅ Завершено | 52 положення (L1: 41, L2: 11; `100% EXPLICIT_COMPLETE`, 7 канонічних блоків `[!NOTE]`, L4: 3, разом: 55) | 10/10 підтверджено (`100% PASS`; 0 залишків стоп-слова) |
| **Глава 27** | [ch27-safety-case-gsn-synthesis.md](file:///home/mfedchyk/book/ExpertSystem/ch27-safety-case-gsn-synthesis.md) | [registry/ch27-registry.md](file:///home/mfedchyk/book/ExpertSystem/registry/ch27-registry.md) | ✅ Завершено | 49 положень (L1: 41, L2: 8; `100% EXPLICIT_COMPLETE`, 5 канонічних блоків `[!NOTE]`, L4: 11, разом: 60) | 22/22 підтверджено (`100% PASS`; додано Leveson STPA, Bloomfield CAE, Avizienis) |
| **Глава 30** | [ch30-safety-cybersecurity-co-engineering.md](file:///home/mfedchyk/book/ExpertSystem/ch30-safety-cybersecurity-co-engineering.md) | [registry/ch30-registry.md](file:///home/mfedchyk/book/ExpertSystem/registry/ch30-registry.md) | ✅ Завершено | 74 положення (L1: 64, L2: 10; `100% EXPLICIT_COMPLETE`, 6 канонічних блоків `[!NOTE]`, L4: 10, разом: 84) | 18/18 підтверджено (`100% PASS`; прив'язано [18] eBPF LSM, Ross Anderson, Rushby MILS) |

> **Підсумок аудиту Частини V (Глави 23, 36, 25, 26, 39, 27, 30):**  
> Проаналізовано **521 наукову одиницю**, з яких **437 одиниць (83,9 %)** становлять авторський науковий та інженерний доробок Миколи Федчика (`L1: 362`, `L2: 75`), а **84 одиниці (16,1 %)** — класичний науковий фундамент і стандарти безпеки (`L4`).

### Частина VI. Дворежимні системи, абдуктивне відкриття та нейро-символьні архітектури ([part-06-frontiers-neuro-symbolic.md](file:///home/mfedchyk/book/ExpertSystem/part-06-frontiers-neuro-symbolic.md))

| Глава | Файл глави | Файл поглавного реєстру | Статус аудиту | Кількість виявлених авторських положень (L1–L3) | Перевірка цитувань (L4) |
| :--- | :--- | :--- | :--- | :--- | :--- |
| **Глава 28** | [ch28-dual-mode-expert-systems.md](file:///home/mfedchyk/book/ExpertSystem/ch28-dual-mode-expert-systems.md) | [registry/ch28-registry.md](file:///home/mfedchyk/book/ExpertSystem/registry/ch28-registry.md) | ✅ Завершено | 49 положень (L1: 40, L2: 9; `100% EXPLICIT_COMPLETE`, 6 канонічних блоків `[!NOTE]`, L4: 4, разом: 53) | 13/13 підтверджено (`100% PASS`; додано Kahneman, Kolodner, ISO 26262, DO-178C) |
| **Глава 34** | [ch34-deterministic-relational-analysis-abduction-and-socratic-dialogue.md](file:///home/mfedchyk/book/ExpertSystem/ch34-deterministic-relational-analysis-abduction-and-socratic-dialogue.md) | [registry/ch34-registry.md](file:///home/mfedchyk/book/ExpertSystem/registry/ch34-registry.md) | ✅ Завершено | 56 положень (L1: 45, L2: 11; `100% EXPLICIT_COMPLETE`, 10 канонічних блоків `[!NOTE]`, L4: 7, разом: 63) | 9/9 підтверджено (`100% PASS`; додано Pearl, McCarthy, Codd, Shannon, Popper) |
| **Глава 29** | [ch29-neuro-symbolic-architecture.md](file:///home/mfedchyk/book/ExpertSystem/ch29-neuro-symbolic-architecture.md) | [registry/ch29-registry.md](file:///home/mfedchyk/book/ExpertSystem/registry/ch29-registry.md) | ✅ Завершено | 62 положення (L1: 53, L2: 9; `100% EXPLICIT_COMPLETE`, 7 канонічних блоків `[!NOTE]`, L4: 13, разом: 75) | 29/29 підтверджено (`100% PASS`; додано Kautz NeSy Taxonomy, LeCun JEPA, DeepProbLog) |
| **Глава 38** | [ch38-curing-machine-hallucinations-and-knowledge-deficits.md](file:///home/mfedchyk/book/ExpertSystem/ch38-curing-machine-hallucinations-and-knowledge-deficits.md) | [registry/ch38-registry.md](file:///home/mfedchyk/book/ExpertSystem/registry/ch38-registry.md) | ✅ Завершено | 63 положення (L1: 53, L2: 10; `100% EXPLICIT_COMPLETE`, 6 канонічних блоків `[!NOTE]`, L4: 10, разом: 73) | 18/18 підтверджено (`100% PASS`; прив'язано інлайн-якір [17] Azaria, додано [18] Bender 2021) |

> **Підсумок аудиту Частини VI (Глави 28, 34, 29, 38):**  
> Проаналізовано **264 наукові одиниці**, з яких **230 одиниць (87,1 %)** становлять авторський науковий та інженерний доробок Миколи Федчика (`L1: 191`, `L2: 39`), а **34 одиниці (12,9 %)** — класичний науковий базис (`L4`).

### Частина VII. Інфраструктура виконання, кіберфізичне керування та розподілені системи ([part-07-runtime-and-knowledge-exchange.md](file:///home/mfedchyk/book/ExpertSystem/part-07-runtime-and-knowledge-exchange.md))

| Глава | Файл глави | Файл поглавного реєстру | Статус аудиту | Кількість виявлених авторських положень (L1–L3) | Перевірка цитувань (L4) |
| :--- | :--- | :--- | :--- | :--- | :--- |
| **Глава 17** | [ch17-implementation-stack.md](file:///home/mfedchyk/book/ExpertSystem/ch17-implementation-stack.md) | [registry/ch17-registry.md](file:///home/mfedchyk/book/ExpertSystem/registry/ch17-registry.md) | ✅ Завершено | 53 положення (L1: 47, L2: 6; `100% EXPLICIT_COMPLETE`, 6 канонічних блоків `[!NOTE]`, L4: 5, разом: 58) | 41/41 підтверджено (`100% PASS`; додано Lattner, Stonebraker, LMDB) |
| **Глава 18** | [ch18-execution-infrastructure.md](file:///home/mfedchyk/book/ExpertSystem/ch18-execution-infrastructure.md) | [registry/ch18-registry.md](file:///home/mfedchyk/book/ExpertSystem/registry/ch18-registry.md) | ✅ Завершено | 68 положень (L1: 56, L2: 12; `100% EXPLICIT_COMPLETE`, 8 канонічних блоків `[!NOTE]`, L4: 4, разом: 72) | 46/46 підтверджено (`100% PASS`; 0 залишків стоп-слова) |
| **Глава 21** | [ch21-from-recommendation-to-action.md](file:///home/mfedchyk/book/ExpertSystem/ch21-from-recommendation-to-action.md) | [registry/ch21-registry.md](file:///home/mfedchyk/book/ExpertSystem/registry/ch21-registry.md) | ✅ Завершено | 67 положень (L1: 55, L2: 12; `100% EXPLICIT_COMPLETE`, 7 канонічних блоків `[!NOTE]`, L4: 8, разом: 75) | 18/18 підтверджено (`100% PASS`; 0 залишків стоп-слова) |
| **Глава 22** | [ch22-cybernetics-edge-to-backend.md](file:///home/mfedchyk/book/ExpertSystem/ch22-cybernetics-edge-to-backend.md) | [registry/ch22-registry.md](file:///home/mfedchyk/book/ExpertSystem/registry/ch22-registry.md) | ✅ Завершено | 63 положення (L1: 55, L2: 8; `100% EXPLICIT_COMPLETE`, 7 канонічних блоків `[!NOTE]`, L4: 21, разом: 84) | 27/27 підтверджено (`100% PASS`; 0 залишків стоп-слова) |
| **Глава 35** | [ch35-self-organizing-expert-systems-synergetics-and-npu-runtime.md](file:///home/mfedchyk/book/ExpertSystem/ch35-self-organizing-expert-systems-synergetics-and-npu-runtime.md) | [registry/ch35-registry.md](file:///home/mfedchyk/book/ExpertSystem/registry/ch35-registry.md) | ✅ Завершено | 67 положень (L1: 48, L2: 19; `100% EXPLICIT_COMPLETE`, 9 канонічних блоків `[!NOTE]`, L4: 4, разом: 71) | 15/15 підтверджено (`100% PASS`; 0 залишків стоп-слова) |
| **Глава 33** | [ch33-inter-system-knowledge-exchange-and-model-teaching.md](file:///home/mfedchyk/book/ExpertSystem/ch33-inter-system-knowledge-exchange-and-model-teaching.md) | [registry/ch33-registry.md](file:///home/mfedchyk/book/ExpertSystem/registry/ch33-registry.md) | ✅ Завершено | 81 положення (L1: 65, L2: 16; `100% EXPLICIT_COMPLETE`, 6 канонічних блоків `[!NOTE]`, L4: 9, разом: 90) | 13/13 підтверджено (`100% PASS`; 0 залишків стоп-слова) |
| **Глава 40** | [ch40-distributed-epistemic-architectures-soa-and-cluster-scaling.md](file:///home/mfedchyk/book/ExpertSystem/ch40-distributed-epistemic-architectures-soa-and-cluster-scaling.md) | [registry/ch40-registry.md](file:///home/mfedchyk/book/ExpertSystem/registry/ch40-registry.md) | ✅ Завершено | 60 положень (L1: 42, L2: 18; `100% EXPLICIT_COMPLETE`, 9 канонічних блоків `[!NOTE]`, L4: 0, разом: 60) | 13/13 підтверджено (`100% PASS`; 100% положень мають авторську інженерну надбудову) |

> **Підсумок аудиту Частини VII (Глави 17, 18, 21, 22, 35, 33, 40):**  
> Проаналізовано **510 наукових одиниць**, з яких **459 одиниць (90,0 %)** становлять авторський науковий та інженерний доробок Миколи Федчика (`L1: 368`, `L2: 91`), а **51 одиниця (10,0 %)** — класичний базис (`L4`).

### Збережені положення після скорочення редакції

| Положення | Канонічне місце | Поглавний облік |
|---|---|---|
| Вердикти `ACCEPT`, `QUALIFIED`, `REFUSAL`, підстатуси й межі повноважень | [Глава 16, §1.1](ch16-expert-systems-architecture.md#11-канонічний-словник-вердиктів-рантайму-та-межі-повноважень) | [Реєстр глави 16](registry/ch16-registry.md) |
| Запас до порога, чинна межа сумарної похибки та цифрове перерахування | [Глава 18, §8.3.1](ch18-execution-infrastructure.md#831-запас-до-порога-та-чинна-межа-сумарної-похибки) | [Реєстр глави 18](registry/ch18-registry.md) |
| Попередня реєстрація, незалежний еталон і облік негативних результатів | [Глава 36, §5.9](ch36-knowledge-testing-pyramid-and-variational-calibration.md#59-протокол-доказового-експерименту-та-попередня-реєстрація) | [Реєстр глави 36](registry/ch36-registry.md) |
| Незалежний аудит мовчазних пропусків до цифрового шлюзу | [Глава 37, §5.1](ch37-input-information-assessment-and-algorithmic-skepticism.md#51-незалежний-аудит-мовчазних-пропусків) | [Реєстр глави 37](registry/ch37-registry.md) |

Перенесення не створює другого авторського результату. Ці положення мають одне канонічне місце; нові поглавні записи й змінені рядкові координати наведено у відповідних реєстрах. Наведені нижче первинні кількісні зрізи не збільшено механічно на число перенесених блоків.

---

## 3. Історична кількісна базова лінія сорока глав

### 3.1. Генеральна кількісна матриця наукового доробку Миколи Федчика (за канонічною структурою Частин I–VII)

| Частина монографії | К-сть документів | Усього наукових одиниць | Чисті авторські розробки (`L1`) | Авторські розширення та синтез (`L2/L3`) | Сукупний авторський доробок (`L1–L3`) | Класичний базис і стандарти (`L4`) | Частка авторського доробку |
| :--- | :---: | :---: | :---: | :---: | :---: | :---: | :---: |
| **Частина I (Глави 1–5)** | 5 | 241 | 146 | 77 | **223** | 18 | **92,5 %** |
| **Частина II (Глави 6–9, 32)** | 5 | 364 | 239 | 62 | **301** | 63 | **82,7 %** |
| **Частина III (Глави 37, 10–12, 14–15)** | 6 | 449 | 305 | 89 | **394** | 55 | **87,8 %** |
| **Частина IV (Глави 16, 13, 19, 31, 24, 20)** | 6 | 394 | 273 | 82 | **355** | 39 | **90,1 %** |
| **Частина V (Глави 23, 36, 25, 26, 39, 27, 30)** | 7 | 521 | 362 | 75 | **437** | 84 | **83,9 %** |
| **Частина VI (Глави 28, 34, 29, 38)** | 4 | 264 | 191 | 39 | **230** | 34 | **87,1 %** |
| **Частина VII (Глави 17, 18, 21, 22, 35, 33, 40)** | 7 | 510 | 368 | 91 | **459** | 51 | **90,0 %** |
| **РАЗОМ: ПЕРВИННИЙ ЗРІЗ ГЛАВ 1–40** | **40** | **2 743** | **1 884 (68,7 %)** | **515 (18,8 %)** | **2 399 (87,5 %)** | **344 (12,5 %)** | **87,5 %** |

---

### 3.2. Головні висновки епістемічного аудиту та перевірки на відсутність плагіату

1. **Класифікація первинного зрізу сорока глав:**
   - Із **2 743** паспортизованих наукових одиниць **1 884 (`68,7 %`)** у первинному аудиті віднесено до `[L1: PURE_ORIGINAL]`.
   - Ще **515 (`18,8 %`)** віднесено до авторських формалізацій та інженерних розширень класичних теорій (`[L2: AUTHORIAL_EXTENSION]` / `[L3: IMPLICIT_UNMARKED]`). Демаркацію між першоджерелом та внеском автора зафіксовано в поглавних реєстрах; після змістових правок вона підлягає перевірці, а не автоматичному підтвердженню цим відсотком.
2. **Канонічна топологічна структура монографії (`Частини I–VII`):**
   - Усі 40 глав організовано в суворому дедуктивному причинно-наслідковому порядку без інверсій залежностей: **Частина I** (епістемічний фундамент) $\to$ **Частина II** (математичні моделі, схеми та `mmap`-пакети знань) $\to$ **Частина III** (вхідний скептицизм та здобуття знань) $\to$ **Частина IV** (архітектура ядра, детерміноване виведення, решітки відношень, діагностика та пояснення) $\to$ **Частина V** (верифікація, піраміда випробувань, індукція, неперервне навчання, аудит та сертифікація безпеки) $\to$ **Частина VI** (дворежимні системи, абдуктивне відкриття та нейро-символьні архітектури) $\to$ **Частина VII** (стек реалізації, інфраструктура виконання, безпечна дія, кібернетика Edge-to-Backend, синергетика NPU, міжсистемний обмін та розподілена SOA).
3. **Межі бібліографічного висновку після зміни редакції:**
   - Первинний зріз сорока глав містить **344 класичні одиниці (`12,5 %`)**. Його статуси не є новою бібліографічною атестацією після перенесення матеріалу.
   - Для актуальної редакції окремо перевіряються підтримка твердження джерелом, коректність атрибуції та посилання. Синтаксичний інспектор і кількість записів не доводять відсутності плагіату або правильності кожної теореми.
4. **Результати повної зачистки забороненого машинного штампа (за п. 2 `AGENTS.md`) (`100% PASS`):**
   - Усі поодинокі залишки забороненого машинного штампа у главах `03, 18, 21, 22, 23, 24, 26, 33, 35, 39` повністю усунуто й замінено на канонічну інженерну термінологію (*тракт*, *цикл керування*, *оперативне середовище*, *ланцюг зворотного зв'язку*).
   - Після скорочення редакції перевіряються сорок глав і сорок поглавних реєстрів; результати первинного сканування не підміняють нового прогону `scripts/inspect_buzzwords.py`.
