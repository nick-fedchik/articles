# Операційний Playbook комплексного аудиту вимог і реалізації (з NPU Forge)

Цей документ є покроковою інструкцією (runbook) для інженера-аудитора та оркестратора, який виконує комплексний аудит проєктної документації та коду за методикою [README.md](README.md) (Karl Wiegers, ISO/IEC/IEEE 29148, INCOSE, EARS) із залученням прискорення на Intel NPU ([RM-006](../../docs/research-memos/RM-006-NPU-Accelerated-Large-Scale-Project-Audit.md), [SWR-034](../../requirements/swr/SWR-034-NPU-Accelerated-Audit-Candidate-Matrix.md)).

---

## 1. Трирівнева модель роботи аудитора

| Рівень | Хто виконує | Інструменти | Час | Призначення |
| --- | --- | --- | --- | --- |
| **Tier 0 — Детермінована механіка** | Node.js-скрипти без LLM | `inventory.mjs`, `linkcheck.mjs`, `id-index.mjs`, `names-check.mjs`, `verify-quotes.mjs`, `trace-sync.mjs`, `status-report.mjs` | $1\text{–}5\text{ с}$ | Інвентар, биті посилання, синтаксис ID, відсутні у коді назви налаштувань/команд, точність цитат, симетрія `ADR ↔ SHR/SWR` |
| **Tier 1 — NPU-відбір кандидатів** | Локальний рушій NPU Forge (`127.0.0.1:3939`) | `NPU Forge: Index Workspace Files`, `npu-candidates.mjs` (`POST /v1/audit/candidates`) | $2\text{–}90\text{ с}$ | Пакетний ембединг корпусу на NPU ($281.52\text{ chunks/s}$) і звуження $O(N^2)$ простору пар до top-$k$ кандидатів (`conflicts`, `cascades`, `coverage`) з позначкою `"origin": "candidate"` |
| **Tier 2 — Верифікація і рішення** | Субагенти за промптами `01..19` + детермінований граф + власник | Промпти [prompts/](prompts/), `/v1/traceability/*`, `DECISIONS-1.md` | За пакетами (`≤ 15–20 хв` на пачку) | Смислова перевірка відібраних пар, відтворення дефектів у коді (`V-NN`), рішення власника «код чи текст», фінальне приймання |

---

## 2. Підготовка середовища та ізоляція (`Перед стартом`)

1. **Перевірте ізоляцію робочого каталогу `audit/work/`:**
   - У `.gitignore` проєкту має бути рядок `/audit/work/`.
   - Створіть структуру робочих каталогів:
     ```powershell
     New-Item -ItemType Directory -Force -Path audit/work/scan, audit/work/sources, audit/work/lenses, audit/work/tasks, audit/work/snapshots, audit/work/scripts, audit/work/tmp | Out-Null
     ```
   - **Жорстке правило:** жодні тимчасові файли чи пробні скрипти не створюються у `scripts/`, `docs/`, тестових каталогах або `%TEMP%` — лише в `audit/work/tmp/RUN-NNN/` та `audit/work/scripts/`.

2. **Підніміть локальний рушій NPU Forge і побудуйте секційний індекс (`Tier 1`):**
   - У VS Code виконайте **`NPU Forge: Start Local Engine`** (або переконайтеся, що `GET http://127.0.0.1:3939/health` повертає `{"status":"ok"}`).
   - Виконайте **`NPU Forge: Index Workspace Files`** (`npuForge.indexWorkspace`) для кореня репозиторію.
   - Перевірте наявність `.npu-forge/index.bin` (`NPU Forge: Show Local Index Location`). Документи Markdown індексуються за заголовками секцій (`NPUFIDX3`), код — за межами функцій/вікон.

---

## 3. Покроковий маршрут за фазами та контрольними точками

### Фаза 0 — Інвентаризація, швидкий скан і NPU-преселекція (`КТ0`)

1. **Запустіть механічні утиліти (`Tier 0`):**
   ```powershell
   node audit/methodology/tools/inventory.mjs README.md Project-Overvew.md CHANGELOG.md CONTRIBUTING.md requirements adr docs packaging > audit/work/scan/inventory.md
   node audit/methodology/tools/linkcheck.mjs . --tracked > audit/work/scan/linkcheck-initial.txt
   node audit/methodology/tools/id-index.mjs requirements adr docs/research-memos --prefixes SHR,SWR,ADR,RM > audit/work/scan/id-index.md
   node audit/methodology/tools/trace-sync.mjs > audit/work/scan/trace-sync-initial.txt
   node audit/methodology/tools/names-check.mjs --docs README.md,docs,requirements,adr --code engine,extension --pattern "npuForge\.[A-Za-z0-9_]+" > audit/work/scan/names-settings.md
   ```
2. **Згенеруйте NPU-матриці кандидатів (`Tier 1`):**
   ```powershell
   # Кандидати на внутрішньорівневі конфлікти вимог (лінза C)
   node audit/methodology/tools/npu-candidates.mjs --mode conflicts --sources requirements/shr,requirements/swr --targets requirements/shr,requirements/swr --top-k 5 --min-score 0.45 --out audit/work/scan/candidates-conflicts.json

   # Кандидати на міжрівневі каскади та розсинхронізацію (лінзи T, S)
   node audit/methodology/tools/npu-candidates.mjs --mode cascades --sources requirements/shr,requirements/swr --targets adr,docs --top-k 5 --min-score 0.40 --out audit/work/scan/candidates-cascades.json

   # Кандидати на покриття вимог кодом і список вимог без кандидатів (фаза 10)
   node audit/methodology/tools/npu-candidates.mjs --mode coverage --sources requirements/swr --targets engine/src,extension/src --top-k 5 --min-score 0.35 --out audit/work/scan/candidates-coverage.json
   ```
3. **Сформуйте базові файли управління ([ARTIFACTS.md](ARTIFACTS.md) §3.10–3.14):**
   - `audit/work/README.md` (базовий коміт `git rev-parse HEAD`, таблиця параметрів `{{PROJECT}}`, `{{DOCS}}`, `{{CODE}}`).
   - `audit/work/STATUS.md`, `audit/work/RUNLOG.md`, `audit/work/PROGRESS.md`.
   - `audit/work/scan/F-SCAN.md` (знахідки механічного скану `F-SCAN-NNN`).
4. **Контрольна точка `КТ0`:** покажіть власнику обсяг (повний vs за ризиком), результати `F-SCAN` і дочекайтеся рішення про старт.

---

### Фази 1–2 та пілот видобування (`КТ1`)

1. **Фаза 1 (Аудит задуму, промпт [02](prompts/02-concept-audit.md)):**
   - Запис у `audit/work/CONCEPT.md` (знахідки `PC-01..PC-NN`).
2. **Фаза 2 (Аудит реалізації, промпт [03](prompts/03-implementation-audit.md)):**
   - Аналіз архітектури, безпеки, граничних умов та тестів у `audit/work/IMPLEMENTATION.md` (`IM-01..IM-NN`).
   - Динамічне відтворення дефектів діагностичними пробами в `audit/work/VERIFICATION.md` (`V-01..V-NN`).
3. **Пілотний пакет видобування `E01` (промпт [04](prompts/04-extraction.md)):**
   - Видобування перших 4–5 джерел у `audit/work/sources/*.md`.
   - Перевірка цитат:
     ```powershell
     node audit/methodology/tools/verify-quotes.mjs audit/work/sources .
     ```
4. **Контрольна точка `КТ1`:** звіт власнику з `PC-*`, `IM-*`, `V-*` та виміряною швидкістю видобування.

---

### Фази 3–4 — Повне видобування і валідація ознак (`КТ1b`)

1. **Фаза 3 (Пакети `E02..En`, промпт [04](prompts/04-extraction.md)):**
   - Кожен пакет: до ~100 тверджень або ~1500 рядків джерел.
   - Агент створює файл джерела зі статусом `в роботі`, записує секції `RQ-NNN` порціями і ставить статус `готово` в кінці.
   - **Ритм Heartbeat:** не більше 2–4 запусків субагентів за один хід розмови, після чого оновлення `STATUS.md` і `PROGRESS.md`.
2. **Механічна верифікація всіх цитат:**
   ```powershell
   node audit/methodology/tools/verify-quotes.mjs audit/work/sources .
   ```
   Очікуваний результат: `in_progress=0`, `mismatch=0`.
3. **Фаза 4 (Валідація ознак на стратифікованій вибірці, промпт [05](prompts/05-flag-validation.md)):**
   - Перевірка точності ознак (`ambiguous`, `unverifiable`, `compound`, `status_claim` тощо) на вибірці $\ge 50$ вимог.
   - Запис правил винятків `EX-01..EX-NN` у `audit/work/FLAG-VALIDATION.md`, щоб лінзи не плодили хибних спрацьовувань.
4. **Проміжна точка `КТ1b`:** фіксація підсумку фаз 3–4 у `PROGRESS.md`.

---

### Фаза 5 — Шість лінз аналізу з NPU-кандидатами (`Q, C, T, D, S, N`)

Під час формування вхідних завдань для агентів лінз оркестратор додає відповідні зрізи з `audit/work/scan/candidates-*.json`:

| Лінза | Промпт | Що подається на вхід агенту | Вихідні файли |
| --- | --- | --- | --- |
| **Q** (Якість) | [06-lens-quality.md](prompts/06-lens-quality.md) | Видобуті `sources/*.md` + правила винятків `EX-*` | `lenses/L-Q*.md` |
| **C** (Конфлікти) | [07-lens-conflicts.md](prompts/07-lens-conflicts.md) | Видобуті `sources/*.md` + **пари-кандидати з `candidates-conflicts.json`** + `names-check` | `lenses/L-C*.md` |
| **T** (Трасування) | [08-lens-traceability.md](prompts/08-lens-traceability.md) | `id-index.md`, `trace-sync.txt` + **пари з `candidates-cascades.json`** | `lenses/L-T.md`, `TRACE-MATRIX.md` |
| **D** (Декларації) | [09-lens-declarations.md](prompts/09-lens-declarations.md) | Усі `RQ` з ознакою `status_claim` проти реального коду і тестів | `lenses/L-D*.md` |
| **S** (Застарілість і каскади) | [10-lens-staleness-cascade.md](prompts/10-lens-staleness-cascade.md) | Лічильники, версії, дати + **міжрівневі пари з `candidates-cascades.json`** | `lenses/L-S.md` |
| **N** (Прогалини) | [11-lens-gaps.md](prompts/11-lens-gaps.md) | Чеклист категорій (безпека, приватність, ліміти, помилки, сумісність) | `lenses/L-N.md` |

---

### Фази 6–8 — Консолідація, задачі беклогу та рішення власника (`КТ2`)

1. **Фаза 6 ([12-consolidation-report.md](prompts/12-consolidation-report.md)):**
   - Дедуплікація всіх знахідок `PC-*`, `IM-*`, `F-SCAN-*`, `F-Q/C/T/D/S/N-*` у єдиний реєстр кластерів `audit/work/REGISTER.md` (`R-01..R-NN`) з обов'язковою перевіркою доказів оркестратором.
   - Формування управлінського звіту `audit/work/REPORT.md` (включно з базовою оцінкою ДО виправлень за двома шкалами: **1. Критерії якості артефактів** та **2. Еволюційна градація практик інженерії вимог за Wiegers / CMMI / ASPICE**, див. [ARTIFACTS.md](ARTIFACTS.md) §3.6 п. 12–13) та списку відхилених гіпотез `audit/work/REJECTED.md`.
2. **Фаза 7 ([13-backlog-tasks.md](prompts/13-backlog-tasks.md)):**
   - Створення задач `audit/work/tasks/TASK-*-NN.md` з критеріями приймання та командами перевірки.
3. **Фаза 8 ([14-owner-decisions.md](prompts/14-owner-decisions.md)):**
   - Формування опитувальника `audit/work/DECISIONS-1.md` для всіх дилем «виправити код чи оновити документацію» з явною рекомендацією аудитора.
4. **Контрольна точка `КТ2`:** власник затверджує рішення у `DECISIONS-1.md`.

---

### Фаза 9 — Внесення рішень у документацію (`КТ2b`)

1. **Архівний знімок перед правками:**
   - Збережіть копію або хеш стану документів у `audit/work/snapshots/`.
2. **Пакети правок `APPLY-01..APPLY-NN` ([15](prompts/15-apply-decisions.md), [16](prompts/16-requirement-refinement.md)):**
   - **Жорсткий ліміт:** не більше **8–10 файлів** на один запуск агента.
   - Послідовний запуск, якщо групи мають спільні індексні файли (`requirements/README.md`, `adr/README.md`).
3. **Автоматична синхронізація та перевірка після кожного пакета (`Tier 0`):**
   ```powershell
   node audit/methodology/tools/trace-sync.mjs --write
   node audit/methodology/tools/linkcheck.mjs . --tracked
   node audit/methodology/tools/id-index.mjs requirements adr docs/research-memos --prefixes SHR,SWR,ADR,RM
   ```
4. **Проміжна точка `КТ2b`:** збереження `git diff` у `audit/work/snapshots/` і запис у `PROGRESS.md`.

---

### Фаза 10 — Покриття вимог кодом і тестами (`КТ3`)

1. **Пакети покриття `COV-01..COV-NN` ([17-code-coverage.md](prompts/17-code-coverage.md)):**
   - **Жорсткий ліміт:** не більше **10–12 документів вимог** або **60–70 критеріїв приймання** на один запуск.
   - Використовуйте згенерований на NPU файл `audit/work/scan/candidates-coverage.json` як стартову карту `source_path -> target_path:target_line`.
   - Агент створює вихідний файл першим кроком і дописує таблицю кожні 3–4 вимоги (**Chunked Flush**).
2. **Вибіркова перевірка 40 вердиктів реальним запуском тестів:**
   - Відберіть стратифіковану вибірку (40 пунктів із фіксованим seed) і запустіть названі C++/TypeScript тести, щоб переконатися, що вони існують і проходять.
3. **Контрольна точка `КТ3`:** звірка `COVERAGE-CODE.md` із власником.

---

### Фаза 11 — Незалежне рев'ю, виправлення коду та архівація (`КТ4`)

1. **Незалежне рев'ю ([18-independent-review.md](prompts/18-independent-review.md)):**
   - Перевірка оновленої документації, матриць трасування, інсталяторів та релізних скриптів свіжим проходом (`lenses/L-REVIEW.md`).
2. **Виконання задач із виправлення коду (`TASK-*-01..NN`) та релізний гейт:**
   - Після внесення виправлень у C++ рушій та VS Code розширення запустіть повний гейт:
     ```powershell
     .\scripts\release-gate.ps1
     node audit/methodology/tools/trace-sync.mjs
     node audit/methodology/tools/linkcheck.mjs . --tracked
     ```
3. **Архівація завершеного аудиту в `audit/done/<YYYY-MM-DD>/` (`КТ4`):**
   - Заповніть у `REPORT.md` підсумкові колонки **«ПІСЛЯ виправлень (КТ4)»** для обох обов'язкових таблиць оцінювання:
     1. **Критерії якості артефактів (Karl Wiegers / ISO/IEC/IEEE 29148)**;
     2. **Еволюційна градація практик інженерії вимог (Wiegers / CMMI / ASPICE)**.
   - Скопіюйте підсумкові артефакти (`REPORT.md`, `REGISTER.md`, `DECISIONS-1.md`, `CONCEPT.md`, `IMPLEMENTATION.md`, `VERIFICATION.md`, `TRACE-MATRIX.md`, `COVERAGE-CODE.md`, `STATUS.md`, `RUNLOG.md`, `PROGRESS.md`) у `audit/done/<YYYY-MM-DD>/`.

---

## 4. Швидке відновлення після обриву сесії або таймауту

Якщо сесія обірвалася або агент перевищив ліміт часу:

1. **Перша команда в новій сесії (1 секунда):**
   ```powershell
   node audit/methodology/tools/status-report.mjs audit/work
   ```
   Утиліта покаже:
   - відсоток завершення кожної фази;
   - останній завершений `RUN-NNN` і запуски, що зависли у стані `в роботі`;
   - файли у `audit/work/sources/`, де залишився статус `в роботі`;
   - залишки в `audit/work/tmp/` та змінені файли в `git status --short`.
2. **Правило Resume without Redo:**
   - Ніколи не видаляйте вже записані та валідні файли перерваного пакета.
   - Перевірте, які саме файли з пакета вже відредаговані або записані в журнал, закрийте `RUN-NNN` як `перервано (частково виконано)` і запустіть дочірній `RUN-NNNb` **лише для решти не опрацьованих файлів**.
   - Якщо агент упав до створення журналу, розділіть пакет навпіл (`RUN-NNNa` та `RUN-NNNb`).

---

## 5. Чеклист фінальної здачі аудиту (`КТ4 Gate`)

- [ ] `node audit/methodology/tools/verify-quotes.mjs audit/work/sources .` → `in_progress=0`, `mismatch=0`.
- [ ] `node audit/methodology/tools/trace-sync.mjs` → `asymmetries=0`.
- [ ] `node audit/methodology/tools/linkcheck.mjs . --tracked` → `broken=0` (жодного посилання з відстежуваних документів на ігноровані файли `audit/work/`).
- [ ] `node audit/methodology/tools/id-index.mjs requirements adr docs/research-memos --prefixes SHR,SWR,ADR,RM` → 0 битих ID, 0 дублікатів.
- [ ] `.\scripts\release-gate.ps1` → `All suites passed` (усі 6 C++ тестових наборів + усі TypeScript тести).
- [ ] У `REPORT.md` заповнені обидві таблиці підсумкової оцінки (ДО і ПІСЛЯ): **1) Критерії якості артефактів** та **2) Еволюційна градація практик інженерії вимог (Wiegers / CMMI / ASPICE)**.
- [ ] Підсумкові звіти скопійовані в `audit/done/<YYYY-MM-DD>/`.

---

## 6. Очищення після завершення аудиту

| Каталог | Статус у Git | Що робити після закриття `КТ4` |
| --- | --- | --- |
| `audit/methodology/` | Відстежується | **Не видаляти ніколи.** Містить методику, промпти `prompts/`, спільні утиліти `tools/*.mjs` і цей `PLAYBOOK.md`. |
| `audit/done/<YYYY-MM-DD>/` | Відстежується | **Залишити назавжди.** Незмінний історичний архів завершеного аудиту (`REPORT.md`, `REGISTER.md`, `DECISIONS-1.md`, `CONCEPT.md`, `IMPLEMENTATION.md`, `VERIFICATION.md`, `TRACE-MATRIX.md`, `COVERAGE-CODE.md`, `STATUS.md`, `RUNLOG.md`, `PROGRESS.md`). |
| `audit/work/` | Ігнорується (`/audit/work/`) | **Можна видалити повністю** після копіювання підсумкових звітів у `audit/done/<YYYY-MM-DD>/` та закриття задач `tasks/TASK-*`. |

Команди для очищення:

- **Варіант А — Повне очищення робочої теки** (рекомендовано, коли аудит прийнято і всі задачі `TASK-*` виконані):
  ```powershell
  Remove-Item -Recurse -Force audit/work
  ```
- **Варіант Б — Легке очищення** (якщо планується повторний прохід фази 12 і потрібно тимчасово залишити видобуті твердження `requirements/` та `tasks/`):
  ```powershell
  Remove-Item -Recurse -Force audit/work/tmp, audit/work/scripts, audit/work/snapshots, audit/work/diffs, audit/work/scan -ErrorAction SilentlyContinue
  ```

