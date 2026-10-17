# Реферативний реєстр авторського наукового доробку та епістемічної чистоти монографії «Архітектура доказових експертних систем»

**Автор монографії:** Микола Федчик  
**Призначення документа:** Систематизація, наукова класифікація та верифікація академічної доброчесності (відсутності плагіату й чіткого розмежування першоджерел) усіх авторських теорій, теорем, наукових законів, гіпотез, концепцій, правил, принципів, постулатів, аксіом, математичних моделей, інженерних інваріантів та винаходів у монографії (Глави 1–40, Додатки А–Д, Специфікація MEMO-079).

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
| **`[L2: AUTHORIAL_EXTENSION]`** | Авторська формалізація / синтез на базі класики | Базується на відомій науковій теорії (Байєс, Демпстер–Шефер, Кліні $\mathbb{K}_3$, Белнап $\mathcal{B}_4$, Гаєк, Поппер, Флоріді, Данґ, ASPIC+, Z3, ISO 26262 тощо), але містить **оригінальну інженерну або математичну надбудову автора**. У реєстрі чітко фіксується межа: *«Класичне ядро (автор, рік) vs. Внесок М. Федчика»*. |
| **`[L3: IMPLICIT_UNMARKED]`** | Неявний (прихований) авторський доробок | Оригінальна авторська думка, формула, класифікація або правило вже є в тексті глави, але викладені безособово без формальної назви («Закон...», «Принцип...», «Теорема...») та без вказівки на авторство. Потребує номінації та явної атрибуції. |
| **`[L4: CLASSICAL_NEEDS_CITATION]`** | Класичний базис (контроль плагіату / цитування) | Загальновідомий науковий факт, класична теорема або стандарт, який не є авторським винаходом. Перевіряється наявність коректного посилання на першоджерело у тексті та списку `## Джерела`, щоб виключити будь-яку двозначність щодо авторства. |

### 1.3. Вісь В: Поточний стан відображення у книзі (`BOOK_STATUS`)

- **`[EXPLICIT_COMPLETE]`** — Вже явно описано в тексті та у блоці `[!NOTE]` з чіткою авторською атрибуцією; потребує лише фіксації у загальній науковій номенклатурі.
- **`[EXPLICIT_NEEDS_CLASSIFICATION]`** — Присутнє в тексті та/або в `[!NOTE]`, але сформульоване описово; потребує надання чіткого статусу (*Теорема*, *Закон*, *Аксіома*, *Постулат*, *Принцип*, *Правило*).
- **`[IMPLICIT_NEEDS_ATTRIBUTION]`** — Присутнє у тексті чи формулах безособово; потребує явного маркування авторства та/або додавання у відповідний розділ.
- **`[PROPOSED_INSERTION]`** — Логічно випливає зі змісту розділу або замикає математичну модель, але ще не сформульоване явно в тексті; рекомендовано до вставки після узгодження з автором.

---

## 2. Зведена навігаційна матриця монографії (Скелет реєстру)

### Частина I. Концептуальний та епістемічний фундамент ([part-01-foundations.md](file:///home/mfedchyk/book/ExpertSystem/part-01-foundations.md))

| Глава | Файл глави | Файл поглавного реєстру | Статус аудиту | Кількість виявлених авторських положень (L1–L3) | Перевірка цитувань (L4) |
| :--- | :--- | :--- | :--- | :--- | :--- |
| **Глава 01** | [ch01-introduction-to-expert-systems.md](file:///home/mfedchyk/book/ExpertSystem/ch01-introduction-to-expert-systems.md) | [registry/ch01-registry.md](file:///home/mfedchyk/book/ExpertSystem/registry/ch01-registry.md) | ✅ Завершено (Пілот) | 31 положення (L1: 19, L2: 12; з них 12 потребують маркування) | 15/15 підтверджено (рекомендовано додати Post, Robinson, Forgy, Minsky, DO-178C, ASPIC+) |
| **Глава 02** | [ch02-epistemology-of-machine-knowledge.md](file:///home/mfedchyk/book/ExpertSystem/ch02-epistemology-of-machine-knowledge.md) | [registry/ch02-registry.md](file:///home/mfedchyk/book/ExpertSystem/registry/ch02-registry.md) | ✅ Завершено | 72 положення (L1: 51, L2: 21; 10 канонічних блоків `[!NOTE]`, L4: 10, разом: 82) | 31/32 підтверджено (джерело [30] Geifman не процитовано в тілі; додати Floridi, Bell-LaPadula) |
| **Глава 03** | [ch03-beyond-reference-information-systems.md](file:///home/mfedchyk/book/ExpertSystem/ch03-beyond-reference-information-systems.md) | [registry/ch03-registry.md](file:///home/mfedchyk/book/ExpertSystem/registry/ch03-registry.md) | ✅ Завершено | 33 положення (L1: 14, L2: 9, L3: 10; 4 канонічні блоки `[!NOTE]`, L4: 2, разом: 35) | 12/12 підтверджено (рекомендовано додати Deutsch, Meyer, Z3) |
| **Глава 04** | [ch04-evolution-from-bayes-to-evidence-ai.md](file:///home/mfedchyk/book/ExpertSystem/ch04-evolution-from-bayes-to-evidence-ai.md) | [registry/ch04-registry.md](file:///home/mfedchyk/book/ExpertSystem/registry/ch04-registry.md) | ✅ Завершено (Пілот) | 40 положень (L1: 14, L2: 16, L3: 10) | 29/29 підтверджено (рекомендовано додати Quillian 1968, Forgy 1982) |
| **Глава 05** | [ch05-triad-of-trust-and-corporate-memory.md](file:///home/mfedchyk/book/ExpertSystem/ch05-triad-of-trust-and-corporate-memory.md) | [registry/ch05-registry.md](file:///home/mfedchyk/book/ExpertSystem/registry/ch05-registry.md) | ✅ Завершено (Пілот) | 39 положень (L1: 28, L2: 11; з них 11 потребують маркування) | 14/14 підтверджено (рекомендовано додати Feigenbaum 1984, Ed25519, FIPS 140-3, ISO 27001) |

> **Підсумок аудиту Частини I (Глави 1–5):**  
> Проаналізовано **237 наукових одиниць**, з яких **215 одиниць (90,7 %)** становлять авторський науковий та інженерний доробок Миколи Федчика (`L1: 126`, `L2/L3: 89`), а **22 одиниці (9,3 %)** — класичний науковий фундамент (`L4`).

### Частина II. Математичні моделі, подання та зберігання знань ([part-02-knowledge-models.md](file:///home/mfedchyk/book/ExpertSystem/part-02-knowledge-models.md))

| Глава | Файл глави | Файл поглавного реєстру | Статус аудиту | Кількість виявлених авторських положень (L1–L3) | Перевірка цитувань (L4) |
| :--- | :--- | :--- | :--- | :--- | :--- |
| **Глава 06** | [ch06-applied-mathematics-for-expert-systems.md](file:///home/mfedchyk/book/ExpertSystem/ch06-applied-mathematics-for-expert-systems.md) | [registry/ch06-registry.md](file:///home/mfedchyk/book/ExpertSystem/registry/ch06-registry.md) | ✅ Завершено | 53 положення (L1: 35, L2: 18; з них 20 потребують номінації/маркування) | 49/52 підтверджено ([10] Peirce, [11] Douven без інлайн-лінків; [12a]=[25] дубль) |
| **Глава 07** | [ch07-knowledge-base-typology.md](file:///home/mfedchyk/book/ExpertSystem/ch07-knowledge-base-typology.md) | [registry/ch07-registry.md](file:///home/mfedchyk/book/ExpertSystem/registry/ch07-registry.md) | ✅ Завершено | 70 положень (L1: 49, L2: 21; з них 34 потребують номінації/маркування) | 33/33 підтверджено (рекомендовано додати Lenat, Minsky, Z3, EARS, GBNF) |
| **Глава 08** | [ch08-engineering-artifacts-as-data.md](file:///home/mfedchyk/book/ExpertSystem/ch08-engineering-artifacts-as-data.md) | [registry/ch08-registry.md](file:///home/mfedchyk/book/ExpertSystem/registry/ch08-registry.md) | ✅ Завершено | 48 положень (L1: 35, L2: 13; з них 20 потребують номінації/маркування) | 21/21 підтверджено (рекомендовано додати Parnas 1995, Brooks 1987, UCUM) |
| **Глава 09** | [ch09-engineering-knowledge-graph-traceability.md](file:///home/mfedchyk/book/ExpertSystem/ch09-engineering-knowledge-graph-traceability.md) | [registry/ch09-registry.md](file:///home/mfedchyk/book/ExpertSystem/registry/ch09-registry.md) | ✅ Завершено | 38 положень (L1: 33, L2: 5; з них 26 потребують номінації/маркування) | 11/12 підтверджено (джерело [11] Knowledge Vault без інлайн-посилання в тексті) |
| **Глава 31** | [ch31-syllogistic-reasoning-and-relation-lattices.md](file:///home/mfedchyk/book/ExpertSystem/ch31-syllogistic-reasoning-and-relation-lattices.md) | [registry/ch31-registry.md](file:///home/mfedchyk/book/ExpertSystem/registry/ch31-registry.md) | ✅ Завершено | 61 положення (L1: 43, L2: 18; 7 канонічних блоків `[!NOTE]`, L4: 10, разом: 71) | 18/19 підтверджено ([18] ISO 26262 без якірного лінка; додати Ganter & Wille, Tarski) |
| **Глава 32** | [ch32-high-performance-knowledge-packs-mmap-and-harvesting.md](file:///home/mfedchyk/book/ExpertSystem/ch32-high-performance-knowledge-packs-mmap-and-harvesting.md) | [registry/ch32-registry.md](file:///home/mfedchyk/book/ExpertSystem/registry/ch32-registry.md) | ✅ Завершено | 52 положення (L1: 40, L2: 12; з них 11 потребують номінації/маркування) | 32/32 підтверджено (100% PASS; рекомендовано додати FIPS 204 ML-DSA, RFC 8554) |

> **Підсумок аудиту Частини II (Глави 6, 7, 8, 9, 31, 32):**  
> Проаналізовано **396 наукових одиниць**, з яких **322 одиниці (81,3 %)** становлять авторський науковий та інженерний доробок Миколи Федчика (`L1: 235`, `L2: 87`), а **74 одиниці (18,7 %)** — класичний математичний та алгоритмічний фундамент (`L4`).

### Частина III. Здобуття знань, мовний аналіз та оцінювання входу ([part-03-knowledge-engineering-nlp.md](file:///home/mfedchyk/book/ExpertSystem/part-03-knowledge-engineering-nlp.md))

| Глава | Файл глави | Файл поглавного реєстру | Статус аудиту | Кількість виявлених авторських положень (L1–L3) | Перевірка цитувань (L4) |
| :--- | :--- | :--- | :--- | :--- | :--- |
| **Глава 10** | [ch10-knowledge-acquisition-systems.md](file:///home/mfedchyk/book/ExpertSystem/ch10-knowledge-acquisition-systems.md) | [registry/ch10-registry.md](file:///home/mfedchyk/book/ExpertSystem/registry/ch10-registry.md) | ✅ Завершено | 85 положень (L1: 69, L2: 16; виявлено дубль номера 16.1 на рядку 958) | 37/37 підтверджено (100% PASS; рекомендовано додати MITRE ATLAS, RFC 2119) |
| **Глава 11** | [ch11-knowledge-elicitation-from-experts.md](file:///home/mfedchyk/book/ExpertSystem/ch11-knowledge-elicitation-from-experts.md) | [registry/ch11-registry.md](file:///home/mfedchyk/book/ExpertSystem/registry/ch11-registry.md) | ✅ Завершено | 67 положень (L1: 50, L2: 17; з них 66 потребують номінації/маркування) | 23/23 підтверджено (рекомендовано додати Kelly 1955, Brier 1950, Klein 1998, Vapnik) |
| **Глава 12** | [ch12-linguistic-analysis-and-local-models.md](file:///home/mfedchyk/book/ExpertSystem/ch12-linguistic-analysis-and-local-models.md) | [registry/ch12-registry.md](file:///home/mfedchyk/book/ExpertSystem/registry/ch12-registry.md) | ✅ Завершено | 30 положень (L1: 21, L2: 9) | 39/39 підтверджено (100% PASS; рекомендовано додати Chomsky, Searle, Tarski, UCUM) |
| **Глава 13** | [ch13-language-variability-vs-determinism.md](file:///home/mfedchyk/book/ExpertSystem/ch13-language-variability-vs-determinism.md) | [registry/ch13-registry.md](file:///home/mfedchyk/book/ExpertSystem/registry/ch13-registry.md) | ✅ Завершено | 41 положення (L1: 28, L2: 13; з них 22 потребують номінації/маркування) | 30/30 підтверджено (100% PASS; рекомендовано додати Frege, Montague, Wittgenstein, Roth) |
| **Глава 14** | [ch14-requirements-detection-and-formalization.md](file:///home/mfedchyk/book/ExpertSystem/ch14-requirements-detection-and-formalization.md) | [registry/ch14-registry.md](file:///home/mfedchyk/book/ExpertSystem/registry/ch14-registry.md) | ✅ Завершено | 55 положень (L1: 42, L2: 13; з них 20 потребують номінації/маркування) | 18/20 підтверджено ([18] GBNF та [19] Willard без інлайн-лінків; додати Zave & Jackson) |
| **Глава 15** | [ch15-knowledge-extraction-and-kb-construction.md](file:///home/mfedchyk/book/ExpertSystem/ch15-knowledge-extraction-and-kb-construction.md) | [registry/ch15-registry.md](file:///home/mfedchyk/book/ExpertSystem/registry/ch15-registry.md) | ✅ Завершено | 56 положень (L1: 38, L2: 18; з них 9 потребують номінації/маркування) | 27/27 підтверджено (100% PASS; рекомендовано додати HermiT, WordNet, Etzioni) |
| **Глава 37** | [ch37-input-information-assessment-and-algorithmic-skepticism.md](file:///home/mfedchyk/book/ExpertSystem/ch37-input-information-assessment-and-algorithmic-skepticism.md) | [registry/ch37-registry.md](file:///home/mfedchyk/book/ExpertSystem/registry/ch37-registry.md) | ✅ Завершено | 68 положень (L1: 43, L2: 25; з них 44 потребують номінації/маркування) | 35/37 підтверджено ([22] Guo, [25] Akidau без інлайн-лінків; додати Schum 1994) |

> **Підсумок аудиту Частини III (Глави 10–15, 37):**  
> Проаналізовано **455 наукових одиниць**, з яких **402 одиниці (88,4 %)** становлять авторський науковий та інженерний доробок Миколи Федчика (`L1: 291`, `L2: 111`), а **53 одиниці (11,6 %)** — класичний науковий фундамент (`L4`).

### Частина IV. Архітектура, технологічний стек, виведення та дія ([part-04-architecture-and-inference.md](file:///home/mfedchyk/book/ExpertSystem/part-04-architecture-and-inference.md))

| Глава | Файл глави | Файл поглавного реєстру | Статус аудиту | Кількість виявлених авторських положень (L1–L3) | Перевірка цитувань (L4) |
| :--- | :--- | :--- | :--- | :--- | :--- |
| **Глава 16** | [ch16-expert-systems-architecture.md](file:///home/mfedchyk/book/ExpertSystem/ch16-expert-systems-architecture.md) | [registry/ch16-registry.md](file:///home/mfedchyk/book/ExpertSystem/registry/ch16-registry.md) | ✅ Завершено | 69 положень (L1: 46, L2: 23; з них 30 потребують номінації/маркування) | 25/25 підтверджено (100% PASS; рекомендовано додати Fowler CQRS, McKenney RCU) |
| **Глава 17** | [ch17-implementation-stack.md](file:///home/mfedchyk/book/ExpertSystem/ch17-implementation-stack.md) | [registry/ch17-registry.md](file:///home/mfedchyk/book/ExpertSystem/registry/ch17-registry.md) | ✅ Завершено | 53 положення (L1: 45, L2: 8; з них 39 потребують номінації/маркування) | 38/38 підтверджено (100% PASS; рекомендовано додати Lattner, Stonebraker, LMDB) |
| **Глава 18** | [ch18-execution-infrastructure.md](file:///home/mfedchyk/book/ExpertSystem/ch18-execution-infrastructure.md) | [registry/ch18-registry.md](file:///home/mfedchyk/book/ExpertSystem/registry/ch18-registry.md) | ✅ Завершено | 57 положень (L1: 42, L2: 15; з них 25 потребують номінації/маркування) | 46/46 підтверджено (100% PASS; в тексті глави 2 залишки стоп-слова на рядках 525, 1140) |
| **Глава 19** | [ch19-from-question-to-evidence.md](file:///home/mfedchyk/book/ExpertSystem/ch19-from-question-to-evidence.md) | [registry/ch19-registry.md](file:///home/mfedchyk/book/ExpertSystem/registry/ch19-registry.md) | ✅ Завершено | 48 положень (L1: 38, L2: 10; з них 26 потребують номінації/маркування) | 12/12 підтверджено (100% PASS; рекомендовано додати Liang HELM, FIPS 180-4) |
| **Глава 20** | [ch20-explanation-engine.md](file:///home/mfedchyk/book/ExpertSystem/ch20-explanation-engine.md) | [registry/ch20-registry.md](file:///home/mfedchyk/book/ExpertSystem/registry/ch20-registry.md) | ✅ Завершено | 46 положень (L1: 35, L2: 11; з них 23 потребують номінації/маркування) | 9/10 підтверджено (джерело [10] не процитовано в тілі; додати Swartout, Pearl, Reiter) |
| **Глава 21** | [ch21-from-recommendation-to-action.md](file:///home/mfedchyk/book/ExpertSystem/ch21-from-recommendation-to-action.md) | [registry/ch21-registry.md](file:///home/mfedchyk/book/ExpertSystem/registry/ch21-registry.md) | ✅ Завершено | 64 положення (L1: 53, L2: 11; з них 23 потребують номінації/маркування) | 16/16 підтверджено (100% PASS; в тексті глави 1 залишок стоп-слова на рядку 244) |
| **Глава 22** | [ch22-cybernetics-edge-to-backend.md](file:///home/mfedchyk/book/ExpertSystem/ch22-cybernetics-edge-to-backend.md) | [registry/ch22-registry.md](file:///home/mfedchyk/book/ExpertSystem/registry/ch22-registry.md) | ✅ Завершено | 61 положення (L1: 52, L2: 9; з них 32 потребують номінації/маркування) | 24/24 підтверджено (100% PASS; в тексті глави 3 залишки стоп-слова на рядках 364, 519, 710) |

> **Підсумок аудиту Частини IV (Глави 16–22):**  
> Проаналізовано **455 наукових одиниць**, з яких **398 одиниць (87,5 %)** становлять авторський науковий та інженерний доробок Миколи Федчика (`L1: 311`, `L2: 87`), а **57 одиниць (12,5 %)** — класичний науковий фундамент (`L4`).

### Частина V. Верифікація, тестування, діагностика та обґрунтування безпеки ([part-05-verification-and-learning.md](file:///home/mfedchyk/book/ExpertSystem/part-05-verification-and-learning.md))

| Глава | Файл глави | Файл поглавного реєстру | Статус аудиту | Кількість виявлених авторських положень (L1–L3) | Перевірка цитувань (L4) |
| :--- | :--- | :--- | :--- | :--- | :--- |
| **Глава 23** | [ch23-knowledge-base-verification.md](file:///home/mfedchyk/book/ExpertSystem/ch23-knowledge-base-verification.md) | [registry/ch23-registry.md](file:///home/mfedchyk/book/ExpertSystem/registry/ch23-registry.md) | ✅ Завершено | 76 положень (L1: 66, L2: 10; з них 24 потребують номінації/маркування) | 25/25 підтверджено (100% PASS; в тексті глави 1 залишок стоп-слова на рядку 957) |
| **Глава 24** | [ch24-system-diagnosis.md](file:///home/mfedchyk/book/ExpertSystem/ch24-system-diagnosis.md) | [registry/ch24-registry.md](file:///home/mfedchyk/book/ExpertSystem/registry/ch24-registry.md) | ✅ Завершено | 64 положення (L1: 43, L2: 21; з них 32 потребують номінації/маркування) | 22/22 підтверджено (100% PASS; в тексті глави 2 залишки стоп-слова на рядках 858, 887) |
| **Глава 25** | [ch25-how-expert-systems-learn.md](file:///home/mfedchyk/book/ExpertSystem/ch25-how-expert-systems-learn.md) | [registry/ch25-registry.md](file:///home/mfedchyk/book/ExpertSystem/registry/ch25-registry.md) | ✅ Завершено | 74 положення (L1: 57, L2: 17; з них 24 потребують номінації/маркування) | 47/47 підтверджено (100% PASS; додати Compton RDR, Mitchell, Valiant PAC) |
| **Глава 27** | [ch27-safety-case-gsn-synthesis.md](file:///home/mfedchyk/book/ExpertSystem/ch27-safety-case-gsn-synthesis.md) | [registry/ch27-registry.md](file:///home/mfedchyk/book/ExpertSystem/registry/ch27-registry.md) | ✅ Завершено | 48 положень (L1: 38, L2: 10; 100% EXPLICIT_COMPLETE, 7 канонічних блоків `[!NOTE]`, L4: 12) | 19/19 підтверджено (100% PASS; додати Leveson STPA, Bloomfield CAE, Avizienis) |
| **Глава 30** | [ch30-safety-cybersecurity-co-engineering.md](file:///home/mfedchyk/book/ExpertSystem/ch30-safety-cybersecurity-co-engineering.md) | [registry/ch30-registry.md](file:///home/mfedchyk/book/ExpertSystem/registry/ch30-registry.md) | ✅ Завершено | 71 положення (L1: 53, L2: 18; з них 20 потребують номінації/маркування) | 17/18 підтверджено (джерело [18] не процитовано в тілі; додати Anderson, Rushby MILS) |
| **Глава 36** | [ch36-knowledge-testing-pyramid-and-variational-calibration.md](file:///home/mfedchyk/book/ExpertSystem/ch36-knowledge-testing-pyramid-and-variational-calibration.md) | [registry/ch36-registry.md](file:///home/mfedchyk/book/ExpertSystem/registry/ch36-registry.md) | ✅ Завершено | 44 положення (L1: 35, L2: 9; з них 25 потребують номінації/маркування) | 13/16 підтверджено (джерела [14]–[16] без інлайн-якорів у тексті; додати Vovk, Forgy) |
| **Глава 39** | [ch39-active-compliance-auditor-and-popperian-testing.md](file:///home/mfedchyk/book/ExpertSystem/ch39-active-compliance-auditor-and-popperian-testing.md) | [registry/ch39-registry.md](file:///home/mfedchyk/book/ExpertSystem/registry/ch39-registry.md) | ✅ Завершено | 51 положення (L1: 37, L2: 14; 7 канонічних блоків `[!NOTE]`, L4: 5, 100% EXPLICIT_COMPLETE) | 8/10 підтверджено (джерела [9]–[10] без інлайн-якорів; 1 залишок стоп-слова на рядку 357) |

> **Підсумок аудиту Частини V (Глави 23, 24, 25, 27, 30, 36, 39):**  
> Проаналізовано **507 наукових одиниць**, з яких **428 одиниць (84,4 %)** становлять авторський науковий та інженерний доробок Миколи Федчика (`L1: 329`, `L2: 99`), а **79 одиниць (15,6 %)** — класичний науковий фундамент і стандарти безпеки (`L4`).

### Частина VI. Нейро-символьні моделі, когнітивні фронтири та неперервне навчання ([part-06-frontiers-neuro-symbolic.md](file:///home/mfedchyk/book/ExpertSystem/part-06-frontiers-neuro-symbolic.md))

| Глава | Файл глави | Файл поглавного реєстру | Статус аудиту | Кількість виявлених авторських положень (L1–L3) | Перевірка цитувань (L4) |
| :--- | :--- | :--- | :--- | :--- | :--- |
| **Глава 26** | [ch26-continual-learning.md](file:///home/mfedchyk/book/ExpertSystem/ch26-continual-learning.md) | [registry/ch26-registry.md](file:///home/mfedchyk/book/ExpertSystem/registry/ch26-registry.md) | ✅ Завершено | 56 положень (L1: 47, L2: 9; з них 20 потребують номінації/маркування) | 24/24 підтверджено (100% PASS; додати McCloskey-Cohen, Grossberg ART, Pearl) |
| **Глава 28** | [ch28-dual-mode-expert-systems.md](file:///home/mfedchyk/book/ExpertSystem/ch28-dual-mode-expert-systems.md) | [registry/ch28-registry.md](file:///home/mfedchyk/book/ExpertSystem/registry/ch28-registry.md) | ✅ Завершено | 44 положення (L1: 36, L2: 8; 6 канонічних блоків `[!NOTE]`, L4: 9, 100% EXPLICIT_COMPLETE) | 13/13 підтверджено (100% PASS; додати Kahneman, Kolodner, ISO 26262, DO-178C) |
| **Глава 29** | [ch29-neuro-symbolic-architecture.md](file:///home/mfedchyk/book/ExpertSystem/ch29-neuro-symbolic-architecture.md) | [registry/ch29-registry.md](file:///home/mfedchyk/book/ExpertSystem/registry/ch29-registry.md) | ✅ Завершено | 62 положення (L1: 53, L2: 9; 6 канонічних блоків `[!NOTE]`, L4: 13, 100% EXPLICIT_COMPLETE) | 29/29 підтверджено (100% PASS; додати Kautz NeSy Taxonomy, LeCun JEPA, DeepProbLog) |
| **Глава 34** | [ch34-deterministic-relational-analysis-abduction-and-socratic-dialogue.md](file:///home/mfedchyk/book/ExpertSystem/ch34-deterministic-relational-analysis-abduction-and-socratic-dialogue.md) | [registry/ch34-registry.md](file:///home/mfedchyk/book/ExpertSystem/registry/ch34-registry.md) | ✅ Завершено | 48 положень (L1: 37, L2: 11; 6 канонічних блоків `[!NOTE]`, L4: 6, 100% EXPLICIT_COMPLETE) | 9/9 підтверджено (100% PASS; додати Pearl, McCarthy, Codd, Shannon, Popper) |
| **Глава 38** | [ch38-curing-machine-hallucinations-and-knowledge-deficits.md](file:///home/mfedchyk/book/ExpertSystem/ch38-curing-machine-hallucinations-and-knowledge-deficits.md) | [registry/ch38-registry.md](file:///home/mfedchyk/book/ExpertSystem/registry/ch38-registry.md) | ✅ Завершено | 61 положення (L1: 47, L2: 14; з них 23 потребують номінації/маркування) | 16/17 підтверджено (джерело [17] без інлайн-якоря в тілі; додати Bender 2021) |

> **Підсумок аудиту Частини VI (Глави 26, 28, 29, 34, 38):**  
> Проаналізовано **334 наукові одиниці**, з яких **271 одиниця (81,1 %)** становлять авторський науковий та інженерний доробок Миколи Федчика (`L1: 220`, `L2: 51`), а **63 одиниці (18,9 %)** — класичний науковий базис (`L4`).

### Частина VII. Реактивне виконання, міжсистемний обмін знаннями та розподілена SOA ([part-07-runtime-and-knowledge-exchange.md](file:///home/mfedchyk/book/ExpertSystem/part-07-runtime-and-knowledge-exchange.md))

| Глава | Файл глави | Файл поглавного реєстру | Статус аудиту | Кількість виявлених авторських положень (L1–L3) | Перевірка цитувань (L4) |
| :--- | :--- | :--- | :--- | :--- | :--- |
| **Глава 33** | [ch33-inter-system-knowledge-exchange-and-model-teaching.md](file:///home/mfedchyk/book/ExpertSystem/ch33-inter-system-knowledge-exchange-and-model-teaching.md) | [registry/ch33-registry.md](file:///home/mfedchyk/book/ExpertSystem/registry/ch33-registry.md) | ✅ Завершено | 78 положень (L1: 52, L2: 26; з них 31 потребує номінації/маркування) | 13/13 підтверджено (100% PASS; в тексті глави 2 залишки стоп-слова на рядках 100, 411) |
| **Глава 35** | [ch35-self-organizing-expert-systems-synergetics-and-npu-runtime.md](file:///home/mfedchyk/book/ExpertSystem/ch35-self-organizing-expert-systems-synergetics-and-npu-runtime.md) | [registry/ch35-registry.md](file:///home/mfedchyk/book/ExpertSystem/registry/ch35-registry.md) | ✅ Завершено | 59 положень (L1: 48, L2: 11; з них 6 потребують номінації) | 15/15 підтверджено (100% PASS; в тексті глави 1 залишок стоп-слова на рядку 238) |
| **Глава 40** | [ch40-distributed-epistemic-architectures-soa-and-cluster-scaling.md](file:///home/mfedchyk/book/ExpertSystem/ch40-distributed-epistemic-architectures-soa-and-cluster-scaling.md) | [registry/ch40-registry.md](file:///home/mfedchyk/book/ExpertSystem/registry/ch40-registry.md) | ✅ Завершено | 45 положень (L1: 31, L2: 14; 7 канонічних блоків `[!NOTE]`, L4: 0, 100% EXPLICIT_COMPLETE) | 13/13 підтверджено (100% PASS; 100% положень мають авторську інженерну надбудову) |

> **Підсумок аудиту Частини VII (Глави 33, 35, 40):**  
> Проаналізовано **197 наукових одиниць**, з яких **182 одиниці (92,4 %)** становлять авторський науковий та інженерний доробок Миколи Федчика (`L1: 131`, `L2: 51`), а **15 одиниць (7,6 %)** — класичний базис (`L4`).

### Додатки та Специфікації апаратних досліджень

| Документ | Файл документа | Файл реєстру | Статус аудиту | Кількість виявлених авторських положень (L1–L3) | Перевірка цитувань (L4) |
| :--- | :--- | :--- | :--- | :--- | :--- |
| **Додаток А** | [appendix-a-evidence-governed-framework.md](file:///home/mfedchyk/book/ExpertSystem/appendix-a-evidence-governed-framework.md) | [registry/appendix-a-registry.md](file:///home/mfedchyk/book/ExpertSystem/registry/appendix-a-registry.md) | ✅ Завершено | 56 положень (L1: 45, L2: 11; рекомендовано 7 канонічних блоків `[!NOTE]`) | 10/10 підтверджено (100% PASS; додати ISO 26262, DO-178C, GSN, Datalog, Peirce) |
| **Додаток Б** | [appendix-b-robotics-and-cyber-physical-systems.md](file:///home/mfedchyk/book/ExpertSystem/appendix-b-robotics-and-cyber-physical-systems.md) | [registry/appendix-b-registry.md](file:///home/mfedchyk/book/ExpertSystem/registry/appendix-b-registry.md) | ✅ Завершено | 49 положень (L1: 27, L2: 22; з них 32 потребують номінації/маркування) | 15/19 підтверджено ([15]–[17] без інлайн-якорів; [18] DARPA ASOR сирітське; додати Ames CBF) |
| **Додаток В** | [appendix-c-autonomous-navigation-and-geosearch.md](file:///home/mfedchyk/book/ExpertSystem/appendix-c-autonomous-navigation-and-geosearch.md) | [registry/appendix-c-registry.md](file:///home/mfedchyk/book/ExpertSystem/registry/appendix-c-registry.md) | ✅ Завершено | 41 положення (L1: 26, L2: 15; з них 19 потребують номінації/маркування) | 13/13 підтверджено (100% PASS; чітко розмежовано Kalman, Thrun, Guttman, H3/S2) |
| **Додаток Г** | [appendix-d-analog-expert-systems-and-neuromorphic-computing.md](file:///home/mfedchyk/book/ExpertSystem/appendix-d-analog-expert-systems-and-neuromorphic-computing.md) | [registry/appendix-d-registry.md](file:///home/mfedchyk/book/ExpertSystem/registry/appendix-d-registry.md) | ✅ Завершено | 31 положення (L1: 24, L2: 7; з них 23 потребують номінації/маркування) | 32/32 підтверджено (100% PASS; розмежовано Mead, Chua, Strukov, Рачковський–Куссуль) |
| **Додаток Д** | [appendix-e-mixed-signal-neuromorphic-expert-systems.md](file:///home/mfedchyk/book/ExpertSystem/appendix-e-mixed-signal-neuromorphic-expert-systems.md) | [registry/appendix-e-registry.md](file:///home/mfedchyk/book/ExpertSystem/registry/appendix-e-registry.md) | ✅ Завершено | 31 положення (L1: 24, L2: 7; з них 14 потребують номінації/маркування) | 38/38 підтверджено (100% PASS; розмежовано Mead, Maass, Davies Loihi, Indiveri, Kanerva) |
| **MEMO-079** | [fpga-epu-hardware-experiment-specification.md](file:///home/mfedchyk/book/ExpertSystem/fpga-epu-hardware-experiment-specification.md) | [registry/memo-079-registry.md](file:///home/mfedchyk/book/ExpertSystem/registry/memo-079-registry.md) | ✅ Завершено | 62 положення (L1: 51, L2: 11; з них 30 потребують номінації/маркування) | 12/12 стандартів + 16/16 глав монографії підтверджено (100% PASS) |

> **Підсумок аудиту Додатків А–Д та Специфікації MEMO-079 (6 документів):**  
> Проаналізовано **309 наукових одиниць**, з яких **270 одиниць (87,4 %)** становлять авторський науковий та інженерний доробок Миколи Федчика (`L1: 197`, `L2: 73`), а **39 одиниць (12,6 %)** — класичний науковий базис і стандарти (`L4`).

---

## 3. Зведений підсумок науково-експертного аудиту всієї монографії (46 документів)

### 3.1. Генеральна кількісна матриця наукового доробку Миколи Федчика

| Частина монографії | К-сть документів | Усього наукових одиниць | Чисті авторські розробки (`L1`) | Авторські розширення та синтез (`L2/L3`) | Сукупний авторський доробок (`L1–L3`) | Класичний базис і стандарти (`L4`) | Частка авторського доробку |
| :--- | :---: | :---: | :---: | :---: | :---: | :---: | :---: |
| **Частина I (Глави 1–5)** | 5 | 237 | 126 | 89 | **215** | 22 | **90,7 %** |
| **Частина II (Глави 6–9, 31, 32)** | 6 | 396 | 235 | 87 | **322** | 74 | **81,3 %** |
| **Частина III (Глави 10–15, 37)** | 7 | 455 | 291 | 111 | **402** | 53 | **88,4 %** |
| **Частина IV (Глави 16–22)** | 7 | 455 | 311 | 87 | **398** | 57 | **87,5 %** |
| **Частина V (Глави 23–25, 27, 30, 36, 39)** | 7 | 507 | 329 | 99 | **428** | 79 | **84,4 %** |
| **Частина VI (Глави 26, 28, 29, 34, 38)** | 5 | 334 | 220 | 51 | **271** | 63 | **81,1 %** |
| **Частина VII (Глави 33, 35, 40)** | 3 | 197 | 131 | 51 | **182** | 15 | **92,4 %** |
| **Додатки А–Д та MEMO-079** | 6 | 309 | 197 | 73 | **270** | 39 | **87,4 %** |
| **РАЗОМ ПО МОНОГРАФІЇ** | **46** | **2 890** | **1 840 (63,7 %)** | **648 (22,4 %)** | **2 488 (86,1 %)** | **402 (13,9 %)** | **86,1 %** |

---

### 3.2. Головні висновки епістемічного аудиту та перевірки на відсутність плагіату

1. **Висока концентрація оригінального наукового та інженерного доробку (`86,1 %`):**
   - Із **2 890** паспортизованих наукових одиниць **1 840 одиниць (`63,7 %`)** є чистими авторськими розробками (`[L1: PURE_ORIGINAL]`): авторські закони, теореми, аксіоми, постулати, математичні оператори з замкненим циклом рантайму (*Actionable Closed Loop*), архітектурні патерни (`Zero-Copy mmap`, `Double-Bank BRAM Buffering`, `Fail-Closed Gates`, `Cognitive Shims`, `SimHash/BK-Tree Triage`, `VAIM`, `EPU`), інженерні інваріанти та експериментальні стенди.
   - Ще **648 одиниць (`22,4 %`)** становлять авторські формалізації та інженерні розширення класичних теорій (`[L2: AUTHORIAL_EXTENSION]`), де в кожному поглавному реєстрі (`registry/*.md`) явно зафіксовано демаркаційну лінію: *«Класичне ядро першоджерела vs. Інженерно-математична надбудова М. Федчика»*.
2. **Академічна чистота та відсутність плагіату (`100 % верифіковано`):**
   - Усі **402 класичні одиниці (`13,9 %`)** коректно спираються на світову наукову спадщину та індустріальні стандарти без жодного привласнення чужих результатів.
   - Із понад **1 050 бібліографічних позицій** у розділах `## Джерела` всієї книги понад **98 %** мають точні активні інлайн-посилання `[[N]](#src-N)` у тілі глав.
3. **Виявлені технічні розбіжності цитування («осиротілі» джерела або пропущені інлайн-якорі у списку `## Джерела`), які легко усуваються точковою правкою:**
   - **Глава 02:** джерело `[30]` (*Geifman & El-Yaniv, 2017*) є у списку `## Джерела`, але не має інлайн-посилання в тексті.
   - **Глава 06:** джерела `[10]` (*Peirce*) та `[11]` (*Douven*) згадані за прізвищами в тексті, але без гіперпосилань `[[10]](#src-10)`, `[[11]](#src-11)`; позиція `[12a]` дублює `[25]` (*Aamodt & Plaza, 1994*).
   - **Глава 09:** джерело `[11]` (*Dong et al., Knowledge Vault*) не має інлайн-посилання в тілі глави.
   - **Глава 14:** джерела `[18]` (*GBNF*) та `[19]` (*Willard & Louf, Outlines*) згадані в тексті, але без якірних посилань `[[18]](#src-18)`, `[[19]](#src-19)`.
   - **Глава 20:** джерело `[10]` (*Fedchyk, JSAKE 2026*) присутнє у списку джерел без інлайн-посилання в тілі.
   - **Глава 30:** джерело `[18]` (*Collin et al., eBPF-LSM, 2024*) присутнє у списку джерел без інлайн-посилання в тілі (рекомендовано прив'язати до підрозділу 1.2).
   - **Глава 31:** джерело `[18]` (*ISO 26262:2018*) згадано в тексті без якірного посилання `[[18]](#src-18)`.
   - **Глава 36:** джерела `[14]` (*Clopper & Pearson*), `[15]` (*Hanley & Lippman-Hand*), `[16]` (*Agache, Firecracker*) згадані за суттю в тексті, але без інлайн-якорів `[[14]]–[[16]]`.
   - **Глава 37:** джерела `[22]` (*Guo et al.*) та `[25]` (*Akidau et al.*) не мають інлайн-якорів у тілі.
   - **Глава 38:** джерело `[17]` (*Fedchyk, Cognitive Shims*) присутнє у списку джерел без інлайн-якоря в підрозділі 1.2.
   - **Глава 39:** джерела `[9]` (*Popper, 1963*) та `[10]` (*Wiegers & Beatty, 2013*) згадані в тексті без інлайн-якорів `[[9]]`, `[[10]]`.
4. **Результати повної зачистки машинного штампа «контур» (100% PASS):**
   - Усі 16 поодиноких залишків забороненого машинного штампа у главах 18, 21, 22, 23, 24, 26, 33, 35, 39 повністю усунуто й замінено на канонічну інженерну термінологію (*тракт*, *цикл керування*, *оперативне середовище*, *ланцюг зворотного зв'язку*).
   - За результатами глобального автоматичного сканування `scripts/inspect_buzzwords.py` усі **104 файли** монографії та 46 файлів реєстру мають **0 зауважень (100% OK)**.
