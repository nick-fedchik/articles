# Реферативний реєстр авторського наукового доробку та епістемічної чистоти монографії «Архітектура доказових експертних систем»

**Автор монографії:** Микола Федчик  
**Призначення документа:** Систематизація, наукова класифікація та верифікація академічної доброчесності (відсутності плагіату й чіткого розмежування першоджерел) усіх авторських теорій, теорем, наукових законів, гіпотез, концепцій, правил, принципів, постулатів, аксіом, математичних моделей, інженерних інваріантів та винаходів у монографії (Глави 1–40, Додатки А–Д).

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

## 2. Зведена навігаційна матриця монографії (Скелет реєстру за канонічною структурою Частин I–VII)

### Частина I. Концептуальний та епістемічний фундамент ([part-01-foundations.md](file:///home/mfedchyk/book/ExpertSystem/part-01-foundations.md))

| Глава | Файл глави | Файл поглавного реєстру | Статус аудиту | Кількість виявлених авторських положень (L1–L3) | Перевірка цитувань (L4) |
| :--- | :--- | :--- | :--- | :--- | :--- |
| **Глава 01** | [ch01-introduction-to-expert-systems.md](file:///home/mfedchyk/book/ExpertSystem/ch01-introduction-to-expert-systems.md) | [registry/ch01-registry.md](file:///home/mfedchyk/book/ExpertSystem/registry/ch01-registry.md) | ✅ Завершено | 33 положення (L1: 23, L2: 10; `100% EXPLICIT_COMPLETE`, 6 канонічних блоків `[!NOTE]`, L4: 4, разом: 37) | 21/21 підтверджено (`100% PASS`; додано Post, Robinson, Forgy, Minsky, DO-178C, ASPIC+) |
| **Глава 02** | [ch02-epistemology-of-machine-knowledge.md](file:///home/mfedchyk/book/ExpertSystem/ch02-epistemology-of-machine-knowledge.md) | [registry/ch02-registry.md](file:///home/mfedchyk/book/ExpertSystem/registry/ch02-registry.md) | ✅ Завершено | 72 положення (L1: 51, L2: 21; `100% EXPLICIT_COMPLETE`, 11 канонічних блоків `[!NOTE]`, L4: 10, разом: 82) | 34/34 підтверджено (`100% PASS`; додано Floridi, Bell-LaPadula, прив'язано Geifman) |
| **Глава 03** | [ch03-beyond-reference-information-systems.md](file:///home/mfedchyk/book/ExpertSystem/ch03-beyond-reference-information-systems.md) | [registry/ch03-registry.md](file:///home/mfedchyk/book/ExpertSystem/registry/ch03-registry.md) | ✅ Завершено | 38 положень (L1: 18, L2: 11, L3: 9; `100% EXPLICIT_COMPLETE`, 10 канонічних блоків `[!NOTE]`, L4: 0, разом: 38) | 16/16 підтверджено (`100% PASS`; додано Deutsch, Meyer, de Moura & Bjørner Z3) |
| **Глава 04** | [ch04-evolution-from-bayes-to-evidence-ai.md](file:///home/mfedchyk/book/ExpertSystem/ch04-evolution-from-bayes-to-evidence-ai.md) | [registry/ch04-registry.md](file:///home/mfedchyk/book/ExpertSystem/registry/ch04-registry.md) | ✅ Завершено | 40 положень (L1: 24, L2: 16; `100% EXPLICIT_COMPLETE`, 6 канонічних блоків `[!NOTE]`, L4: 4, разом: 44) | 31/31 підтверджено (`100% PASS`; додано Quillian 1968, Forgy 1982) |
| **Глава 05** | [ch05-triad-of-trust-and-corporate-memory.md](file:///home/mfedchyk/book/ExpertSystem/ch05-triad-of-trust-and-corporate-memory.md) | [registry/ch05-registry.md](file:///home/mfedchyk/book/ExpertSystem/registry/ch05-registry.md) | ✅ Завершено | 39 положень (L1: 29, L2: 10; `100% EXPLICIT_COMPLETE`, 6 канонічних блоків `[!NOTE]`/`[!WARNING]`, L4: 0, разом: 39) | 18/18 підтверджено (`100% PASS`; додано Feigenbaum 1984, Ed25519, FIPS 140-3, ISO 27001) |

> **Підсумок аудиту Частини I (Глави 1–5):**  
> Проаналізовано **240 наукових одиниць**, з яких **222 одиниці (92,5 %)** становлять авторський науковий та інженерний доробок Миколи Федчика (`L1: 145`, `L2/L3: 77`), а **18 одиниць (7,5 %)** — класичний науковий фундамент (`L4`).

### Частина II. Математичні моделі, подання та зберігання знань ([part-02-knowledge-models.md](file:///home/mfedchyk/book/ExpertSystem/part-02-knowledge-models.md))

| Глава | Файл глави | Файл поглавного реєстру | Статус аудиту | Кількість виявлених авторських положень (L1–L3) | Перевірка цитувань (L4) |
| :--- | :--- | :--- | :--- | :--- | :--- |
| **Глава 06** | [ch06-applied-mathematics-for-expert-systems.md](file:///home/mfedchyk/book/ExpertSystem/ch06-applied-mathematics-for-expert-systems.md) | [registry/ch06-registry.md](file:///home/mfedchyk/book/ExpertSystem/registry/ch06-registry.md) | ✅ Завершено | 54 положення (L1: 36, L2: 18; `100% EXPLICIT_COMPLETE`, 8 канонічних блоків `[!NOTE]`, L4: 44, разом: 98) | 52/52 підтверджено (`100% PASS`; узгоджено якорі Peirce, Douven та Aamodt & Plaza) |
| **Глава 07** | [ch07-knowledge-base-typology.md](file:///home/mfedchyk/book/ExpertSystem/ch07-knowledge-base-typology.md) | [registry/ch07-registry.md](file:///home/mfedchyk/book/ExpertSystem/registry/ch07-registry.md) | ✅ Завершено | 81 положення (L1: 68, L2: 13; `100% EXPLICIT_COMPLETE`, 5 канонічних блоків `[!NOTE]`, L4: 11, разом: 92) | 38/38 підтверджено (`100% PASS`; додано Lenat, Minsky, Z3, EARS, GBNF) |
| **Глава 08** | [ch08-engineering-artifacts-as-data.md](file:///home/mfedchyk/book/ExpertSystem/ch08-engineering-artifacts-as-data.md) | [registry/ch08-registry.md](file:///home/mfedchyk/book/ExpertSystem/registry/ch08-registry.md) | ✅ Завершено | 73 положення (L1: 59, L2: 14; `100% EXPLICIT_COMPLETE`, 5 канонічних блоків `[!NOTE]`, L4: 2, разом: 75) | 24/24 підтверджено (`100% PASS`; додано Parnas 1995, Brooks 1987, UCUM) |
| **Глава 09** | [ch09-engineering-knowledge-graph-traceability.md](file:///home/mfedchyk/book/ExpertSystem/ch09-engineering-knowledge-graph-traceability.md) | [registry/ch09-registry.md](file:///home/mfedchyk/book/ExpertSystem/registry/ch09-registry.md) | ✅ Завершено | 38 положень (L1: 33, L2: 5; `100% EXPLICIT_COMPLETE`, 5 канонічних блоків `[!NOTE]`, L4: 4, разом: 42) | 12/12 підтверджено (`100% PASS`; прив'язано інлайн-посилання Knowledge Vault) |
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
| **Глава 18** | [ch18-execution-infrastructure.md](file:///home/mfedchyk/book/ExpertSystem/ch18-execution-infrastructure.md) | [registry/ch18-registry.md](file:///home/mfedchyk/book/ExpertSystem/registry/ch18-registry.md) | ✅ Завершено | 67 положень (L1: 55, L2: 12; `100% EXPLICIT_COMPLETE`, 7 канонічних блоків `[!NOTE]`, L4: 4, разом: 71) | 46/46 підтверджено (`100% PASS`; 0 залишків стоп-слова) |
| **Глава 21** | [ch21-from-recommendation-to-action.md](file:///home/mfedchyk/book/ExpertSystem/ch21-from-recommendation-to-action.md) | [registry/ch21-registry.md](file:///home/mfedchyk/book/ExpertSystem/registry/ch21-registry.md) | ✅ Завершено | 67 положень (L1: 55, L2: 12; `100% EXPLICIT_COMPLETE`, 7 канонічних блоків `[!NOTE]`, L4: 8, разом: 75) | 18/18 підтверджено (`100% PASS`; 0 залишків стоп-слова) |
| **Глава 22** | [ch22-cybernetics-edge-to-backend.md](file:///home/mfedchyk/book/ExpertSystem/ch22-cybernetics-edge-to-backend.md) | [registry/ch22-registry.md](file:///home/mfedchyk/book/ExpertSystem/registry/ch22-registry.md) | ✅ Завершено | 63 положення (L1: 55, L2: 8; `100% EXPLICIT_COMPLETE`, 7 канонічних блоків `[!NOTE]`, L4: 21, разом: 84) | 27/27 підтверджено (`100% PASS`; 0 залишків стоп-слова) |
| **Глава 35** | [ch35-self-organizing-expert-systems-synergetics-and-npu-runtime.md](file:///home/mfedchyk/book/ExpertSystem/ch35-self-organizing-expert-systems-synergetics-and-npu-runtime.md) | [registry/ch35-registry.md](file:///home/mfedchyk/book/ExpertSystem/registry/ch35-registry.md) | ✅ Завершено | 66 положень (L1: 47, L2: 19; `100% EXPLICIT_COMPLETE`, 8 канонічних блоків `[!NOTE]`, L4: 4, разом: 70) | 15/15 підтверджено (`100% PASS`; 0 залишків стоп-слова) |
| **Глава 33** | [ch33-inter-system-knowledge-exchange-and-model-teaching.md](file:///home/mfedchyk/book/ExpertSystem/ch33-inter-system-knowledge-exchange-and-model-teaching.md) | [registry/ch33-registry.md](file:///home/mfedchyk/book/ExpertSystem/registry/ch33-registry.md) | ✅ Завершено | 80 положень (L1: 64, L2: 16; `100% EXPLICIT_COMPLETE`, 5 канонічних блоків `[!NOTE]`, L4: 9, разом: 89) | 13/13 підтверджено (`100% PASS`; 0 залишків стоп-слова) |
| **Глава 40** | [ch40-distributed-epistemic-architectures-soa-and-cluster-scaling.md](file:///home/mfedchyk/book/ExpertSystem/ch40-distributed-epistemic-architectures-soa-and-cluster-scaling.md) | [registry/ch40-registry.md](file:///home/mfedchyk/book/ExpertSystem/registry/ch40-registry.md) | ✅ Завершено | 57 положень (L1: 39, L2: 18; `100% EXPLICIT_COMPLETE`, 8 канонічних блоків `[!NOTE]`, L4: 0, разом: 57) | 13/13 підтверджено (`100% PASS`; 100% положень мають авторську інженерну надбудову) |

> **Підсумок аудиту Частини VII (Глави 17, 18, 21, 22, 35, 33, 40):**  
> Проаналізовано **504 наукові одиниці**, з яких **453 одиниці (89,9 %)** становлять авторський науковий та інженерний доробок Миколи Федчика (`L1: 362`, `L2: 91`), а **51 одиниця (10,1 %)** — класичний базис (`L4`).

### Додатки (А–Д)

| Документ | Файл документа | Файл реєстру | Статус аудиту | Кількість виявлених авторських положень (L1–L3) | Перевірка цитувань (L4) |
| :--- | :--- | :--- | :--- | :--- | :--- |
| **Додаток А** | [appendix-a-evidence-governed-framework.md](file:///home/mfedchyk/book/ExpertSystem/appendix-a-evidence-governed-framework.md) | [registry/appendix-a-registry.md](file:///home/mfedchyk/book/ExpertSystem/registry/appendix-a-registry.md) | ✅ Завершено | 56 положень (L1: 46, L2: 10; `100% EXPLICIT_COMPLETE`, 7 канонічних блоків `[!NOTE]`, L4: 0, разом: 56) | 15/15 підтверджено (`100% PASS`; додано ISO 26262, DO-178C, GSN, Datalog, Peirce) |
| **Додаток Б** | [appendix-b-robotics-and-cyber-physical-systems.md](file:///home/mfedchyk/book/ExpertSystem/appendix-b-robotics-and-cyber-physical-systems.md) | [registry/appendix-b-registry.md](file:///home/mfedchyk/book/ExpertSystem/registry/appendix-b-registry.md) | ✅ Завершено | 49 положень (L1: 27, L2: 22; `100% EXPLICIT_COMPLETE`, 5 канонічних блоків `[!NOTE]`, L4: 1, разом: 50) | 20/20 підтверджено (`100% PASS`; прив'язано TLE90xx [15], Johnson-Nyquist [16], Popper [17], DARPA ASOR [18], Ames CBF [20]) |
| **Додаток В** | [appendix-c-autonomous-navigation-and-geosearch.md](file:///home/mfedchyk/book/ExpertSystem/appendix-c-autonomous-navigation-and-geosearch.md) | [registry/appendix-c-registry.md](file:///home/mfedchyk/book/ExpertSystem/registry/appendix-c-registry.md) | ✅ Завершено | 43 положення (L1: 32, L2: 11; `100% EXPLICIT_COMPLETE`, 5 канонічних блоків `[!NOTE]`, L4: 1, разом: 44) | 13/13 підтверджено (`100% PASS`; чітко розмежовано Kalman, Thrun, Guttman, H3/S2) |
| **Додаток Г** | [appendix-d-analog-expert-systems-and-neuromorphic-computing.md](file:///home/mfedchyk/book/ExpertSystem/appendix-d-analog-expert-systems-and-neuromorphic-computing.md) | [registry/appendix-d-registry.md](file:///home/mfedchyk/book/ExpertSystem/registry/appendix-d-registry.md) | ✅ Завершено | 27 положень (L1: 21, L2: 6; `100% EXPLICIT_COMPLETE`, 5 канонічних блоків `[!NOTE]`, L4: 17, разом: 44) | 32/32 підтверджено (`100% PASS`; розмежовано Mead, Chua, Strukov, Рачковський–Куссуль) |
| **Додаток Д** | [appendix-e-mixed-signal-neuromorphic-expert-systems.md](file:///home/mfedchyk/book/ExpertSystem/appendix-e-mixed-signal-neuromorphic-expert-systems.md) | [registry/appendix-e-registry.md](file:///home/mfedchyk/book/ExpertSystem/registry/appendix-e-registry.md) | ✅ Завершено | 43 положення (L1: 35, L2: 8; `100% EXPLICIT_COMPLETE`, 5 канонічних блоків `[!NOTE]`, L4: 18, разом: 61) | 38/38 підтверджено (`100% PASS`; розмежовано Mead, Maass, Davies Loihi, Indiveri, Kanerva) |

> **Підсумок аудиту Додатків А–Д (5 документів):**  
> Проаналізовано **255 наукових одиниць**, з яких **218 одиниць (85,5 %)** становлять авторський науковий та інженерний доробок Миколи Федчика (`L1: 161`, `L2: 57`), а **37 одиниць (14,5 %)** — класичний науковий базис і стандарти (`L4`).

---

## 3. Зведений підсумок науково-експертного аудиту всієї монографії (45 документів)

### 3.1. Генеральна кількісна матриця наукового доробку Миколи Федчика (за канонічною структурою Частин I–VII)

| Частина монографії | К-сть документів | Усього наукових одиниць | Чисті авторські розробки (`L1`) | Авторські розширення та синтез (`L2/L3`) | Сукупний авторський доробок (`L1–L3`) | Класичний базис і стандарти (`L4`) | Частка авторського доробку |
| :--- | :---: | :---: | :---: | :---: | :---: | :---: | :---: |
| **Частина I (Глави 1–5)** | 5 | 240 | 145 | 77 | **222** | 18 | **92,5 %** |
| **Частина II (Глави 6–9, 32)** | 5 | 364 | 239 | 62 | **301** | 63 | **82,7 %** |
| **Частина III (Глави 37, 10–12, 14–15)** | 6 | 449 | 305 | 89 | **394** | 55 | **87,8 %** |
| **Частина IV (Глави 16, 13, 19, 31, 24, 20)** | 6 | 394 | 273 | 82 | **355** | 39 | **90,1 %** |
| **Частина V (Глави 23, 36, 25, 26, 39, 27, 30)** | 7 | 521 | 362 | 75 | **437** | 84 | **83,9 %** |
| **Частина VI (Глави 28, 34, 29, 38)** | 4 | 264 | 191 | 39 | **230** | 34 | **87,1 %** |
| **Частина VII (Глави 17, 18, 21, 22, 35, 33, 40)** | 7 | 504 | 362 | 91 | **453** | 51 | **89,9 %** |
| **Додатки А–Д** | 5 | 255 | 161 | 57 | **218** | 37 | **85,5 %** |
| **РАЗОМ ПО МОНОГРАФІЇ** | **45** | **2 991** | **2 038 (68,1 %)** | **572 (19,1 %)** | **2 610 (87,3 %)** | **381 (12,7 %)** | **87,3 %** |

---

### 3.2. Головні висновки епістемічного аудиту та перевірки на відсутність плагіату

1. **Висока концентрація оригінального наукового та інженерного доробку (`87,3 %`):**
   - Із **2 991** паспортизованої наукової одиниці **2 038 одиниць (`68,1 %`)** є чистими авторськими розробками (`[L1: PURE_ORIGINAL]`): авторські закони, теореми, аксіоми, постулати, математичні оператори з замкненим циклом рантайму (*Actionable Closed Loop*), архітектурні патерни (`Zero-Copy mmap`, `Double-Bank BRAM Buffering`, `Fail-Closed Gates`, `Cognitive Shims`, `SimHash/BK-Tree Triage`, `VAIM`), синергетичні фазові переходи NPU-середовища, синхронізація Курамото у розподілених кластерах знань, трирівнева ієрархія складності експертних систем, шестирольова таксономія стейкхолдерів, інженерні інваріанти та експериментальні стенди.
   - Ще **572 одиниці (`19,1 %`)** становлять авторські формалізації та інженерні розширення класичних теорій (`[L2: AUTHORIAL_EXTENSION]` / `[L3: IMPLICIT_UNMARKED]`), де в кожному поглавному реєстрі (`registry/*.md`) явно зафіксовано демаркаційну лінію: *«Класичне ядро першоджерела vs. Інженерно-математична надбудова М. Федчика»*.
2. **Канонічна топологічна структура монографії (`Частини I–VII`):**
   - Усі 40 глав організовано в суворому дедуктивному причинно-наслідковому порядку без інверсій залежностей: **Частина I** (епістемічний фундамент) $\to$ **Частина II** (математичні моделі, схеми та `mmap`-пакети знань) $\to$ **Частина III** (вхідний скептицизм та здобуття знань) $\to$ **Частина IV** (архітектура ядра, детерміноване виведення, решітки відношень, діагностика та пояснення) $\to$ **Частина V** (верифікація, піраміда випробувань, індукція, неперервне навчання, аудит та сертифікація безпеки) $\to$ **Частина VI** (дворежимні системи, абдуктивне відкриття та нейро-символьні архітектури) $\to$ **Частина VII** (стек реалізації, інфраструктура виконання, безпечна дія, кібернетика Edge-to-Backend, синергетика NPU, міжсистемний обмін та розподілена SOA).
3. **Академічна чистота та відсутність плагіату (`100 % верифіковано`):**
   - Усі **381 класична одиниця (`12,7 %`)** коректно спираються на світову наукову спадщину та індустріальні стандарти без жодного привласнення чужих результатів.
   - В **усіх 45 документах монографії (усі 40 глав та всі 5 додатків А–Д)** досягнуто статусу **`100% EXPLICIT_COMPLETE`** та **`100% PASS`** бібліографічної звірки (усунуто всі «осиротілі» джерела, оформлено канонічні авторські блоки `> [!NOTE]` та підтверджено повне покриття інлайн-якорями для 100% бібліографічних позицій).
4. **Результати повної зачистки забороненого машинного штампа (за п. 2 `AGENTS.md`) (`100% PASS`):**
   - Усі поодинокі залишки забороненого машинного штампа у главах `03, 18, 21, 22, 23, 24, 26, 33, 35, 39` повністю усунуто й замінено на канонічну інженерну термінологію (*тракт*, *цикл керування*, *оперативне середовище*, *ланцюг зворотного зв'язку*).
   - За результатами глобального автоматичного сканування `scripts/inspect_buzzwords.py` усі файли монографії та усі **45 файлів реєстру** мають **0 зауважень (`100% OK`)**.

