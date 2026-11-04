# RESEARCH-CH40-GAP-TRIAGE-CONTRACT. Безпечна дедуплікація та обмежений індекс дефіцитів

**ID:** RESEARCH-CH40-GAP-TRIAGE-CONTRACT  
**Статус:** Відкрито; семантичну безпечність злиття та ресурсний профіль не валідовано.  
**Власник:** VIDAdev тріажу; співвиконавець — VIDA researcher алгоритмів пошуку.  
**Пріоритет:** P0 — приховане злиття різних критичних прогалин.  
**Підстава:** [W38-40-10](../../audit-2026-10-28-ch38-ch40.md#w38-40-10), [W38-40-11](../../audit-2026-10-28-ch38-ch40.md#w38-40-11), [CH40 §6.1–§6.2.1, 516–598](../../../ch40-distributed-epistemic-architectures-soa-and-cluster-scaling.md).  
**Вимоги:** [L1](../../../requirements/L1-stakeholder/ch40-L1-stakeholder.md): REQ-L1-40-08; [L2](../../../requirements/L2-system/ch40-L2-system.md): REQ-L2-40-14; [L3](../../../requirements/L3-component/ch40-L3-component.md): REQ-L3-40-15. Окремих strict-schema/DLQ ID у наборі немає.

## 1. Мета та залежності

Перевірити, що LSH лише добирає кандидатів, а об'єднання не приховує іншу сутність, модальність, редакцію чи причину прогалини. До запуску потрібні canonical ticket schema, критичні поля, empty/tie convention і дозволена bounded-memory модель. Спершу запросити реалізацію, не припускати її відсутності.

[ZNA-169](TASK-ZNA-169-wl-tversky-contract.md) та [ZNA-170](TASK-ZNA-170-traversal-budgets.md) стосуються іншого графового індексу; повторно їхній benchmark не проводити.

## 2. Безпечний synthetic-протокол

1. Зафіксувати tokenizer/normalization, feature weights, token hash, sign(0), lexical eligibility і формат node/edge. Формула й реалізація повинні давати однаковий empty/tie результат до guard.
2. Порівняти exact typed-key baseline, SimHash-only, SimHash+semantic discriminator й exhaustive distance oracle.
3. Позитивні контроли — перефразування тієї самої прогалини. Негативні — заперечення, інша сутність, число/одиниця, інша редакція, два різні виключення зі схожими словами, колізія відбитка, zero/all-one та менше двох distinct tokens.
4. На loopback JSON-повідомленнях перевірити `{}`, відсутнє/порожнє/невідоме поле, дубль ключа, неправильний тип, другий об'єкт та version drift. DisallowUnknownFields не вважати presence-check. Некоректне не входить до індексу; DLQ має обмежений ресурс.
5. Для BK-дерева варіювати порядок вставки, щільні/розріджені відбитки, дублікати й adversarial геометрію. Вимірювати visited nodes, фактичні bytes/node+edges+allocator, time/allocations, не підміняти worst-case типовим O(log N).
6. Локальний синтетичний брокер/імітатор із наперед обмеженими повідомленнями, часом і пам'яттю; без production-підписок або шторму в спільній мережі.

## 3. Цілі та приймання

64 біти, Hamming≤3, 64 байти/вузол, 100000 кластерів, 80000 повідомлень/с — **книжкові параметри/заяви**. Допустимі false merges, false splits, max visited nodes і ресурси **очікують погодження**; «95% змісту» не виводити з 61/64 бітів.

**PASS:** негативні семантичні пари не зливаються як одна прогалина, exact oracle збігається з distance search, schema guard і bounded ресурс працюють у профілі. **FAIL:** критична відмінність загублена, malformed прийнято або ресурсна межа прихована. Публікувати всі колізії, worst orders, false splits та відсутність виграшу.

## 4. Артефакти та точні зміни

Передати schema/tokenizer/hash versions, pair corpus/labels/seed, node layout, exact oracle, raw search/memory traces, DLQ records і команди. Майбутній пакет: `python3 reproduce.py --manifest manifest.json --suite gap-triage`; це не заявлений production API.

Оновити CH40 §6.1/§6.1.1/§6.2/§6.2.1: claims схожості/складності/пам'яті, empty convention, Decode+presence приклад і controls. Синхронізувати три ID; окрему алокацію schema/DLQ погодити перед додаванням ID. Не вбудовувати доменні правила в runtime-код.
