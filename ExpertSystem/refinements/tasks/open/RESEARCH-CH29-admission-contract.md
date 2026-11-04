# RESEARCH-CH29-admission-contract. Семантичний допуск, replay та діагностичний дедлайн

**ID:** RESEARCH-CH29-admission-contract  
**Статус:** Відкрито; контрприклади встановлено аналізом Go-коду, production-перевірки не виконано.  
**Власник:** розробник VIDA шлюзу знань та дослідник верифікації семантики.  
**Пріоритет:** P0 — допуск хибної трійки та невизначений deadline.  
**Підстава:** [аудит W27-30-13, W27-30-14, W27-30-16](../../audit-2026-10-28-ch27-ch30.md), [CH29 §5–§7, §11](../../../ch29-neuro-symbolic-architecture.md).  
**Вимоги:** [L1](../../../requirements/L1-stakeholder/ch29-L1-stakeholder.md): REQ-L1-29-01, REQ-L1-29-02, REQ-L1-29-04; [L2](../../../requirements/L2-system/ch29-L2-system.md): REQ-L2-29-01, REQ-L2-29-02, REQ-L2-29-04, REQ-L2-29-05, REQ-L2-29-06, REQ-L2-29-08, REQ-L2-29-13, REQ-L2-29-16; [L3](../../../requirements/L3-component/ch29-L3-component.md): REQ-L3-29-03, REQ-L3-29-05, REQ-L3-29-08, REQ-L3-29-13.

## 1. Залежності та межі

Потрібні погоджені subject/relation/type schema, режим первинного затвердження фактів, execution manifest і профіль діагностичного циклу. Для вільного NLG використати [CH20 fidelity](RESEARCH-CH20-nlg-fidelity.md), не створювати другий NLI-корпус. Загальні benchmarks — [TASK-TST-010](TASK-TST-010-performance-evidence.md).

Кандидат із автентичною цитатою про довжину 1000 октетів, але Subject «гальмо» й Relation `timeout`, проходить lexical checks. Нульова температура не фіксує model/runtime digest; 60-секундний HTTP timeout не забезпечує 120 мс end-to-end.

## 2. Безпечний протокол і контролі

1. Використати синтетичний локальний документ і recorded proposals. Порівняти чинний lexical baseline, незалежний typed-fact oracle й кандидатний semantic gate.
2. Незалежно мутувати subject, relation, одиницю, полярність, модальність, scope/редакцію; залишати число в цитаті. Контролі — справжня трійка, 100/1000, порожнє value, неоднозначна quote, коректне перетворення одиниці лише правилом ядра.
3. Довести допустимий semantic піднабір або залишати факт кандидатом до expert baseline approval. Не використовувати NLI-score як детермінований ASIL-доказ.
4. Для replay фіксувати corpus/KB/model/adapter/grammar/runtime/seed/options і самі proposals. Окремо перевіряти deterministic core replay та live generation; не змішувати їхні гарантії.
5. Замість реального Ollama застосувати loopback stub із детермінованою затримкою/помилками. Перевірити cancellation, slow headers/body, unavailable runtime, malformed JSON, recorded fallback і збереження походження. Жодних інженерних документів назовні.

## 3. Цілі та PASS/FAIL

120 мс — наявна L3-ціль; її належність до конкретного профілю **очікує підтвердження**, а не мовчазного послаблення. Обсяг/рівень довіри semantic evaluation погодити до запуску; нових числових порогів немає.

**PASS:** semantic admissibility визначена й контролі не дають false admission; replay має точний об'єкт; deadline для погодженого профілю включає всі очікування й fallback. **FAIL:** прийнята хибна трійка, замаскована live-відмова або перевищення бюджету. Публікувати всі негативні випадки, зокрема коли безпечний результат — refusal.

## 4. Артефакти, команди й оновлення

Передати schema/oracle, fixtures, model manifests, stub, raw timings і admission traces, paired results. У дослідницькому checkout з виділеними навчальними блоками виконати `go test .` та `go run .`; команду production runner додати після встановлення фактичної структури.

Оновити CH29 §5 (semantic vs lexical checks), §6 (повний manifest), §7 (межі прикладу, timeout/fallback), §11 (профіль 100–120 мс). Синхронізувати всі перелічені вимоги; не забороняти дорадчий нейронний фронтенд і не змінювати А-11/А-22.
