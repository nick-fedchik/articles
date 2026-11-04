# RESEARCH-CH40-DURABILITY-CONTRACT. WAL, незавершена епоха та перевірна квитанція

**ID:** RESEARCH-CH40-DURABILITY-CONTRACT  
**Статус:** Відкрито; модель аварій, квитанцій і квот не погоджена.  
**Власник:** VIDAdev журналювання; співвиконавець — VIDA researcher розподілених сховищ.  
**Пріоритет:** P0 — втрата підтверджених подій або хибна незмінність.  
**Підстава:** [W38-40-12](../../audit-2026-10-28-ch38-ch40.md#w38-40-12), [W38-40-13](../../audit-2026-10-28-ch38-ch40.md#w38-40-13), [CH40 §6.3/§7.1–§7.2.1/§10, 604–724, 914, 932](../../../ch40-distributed-epistemic-architectures-soa-and-cluster-scaling.md).  
**Вимоги:** [L1](../../../requirements/L1-stakeholder/ch40-L1-stakeholder.md): REQ-L1-40-03, REQ-L1-40-09, REQ-L1-40-12; [L2](../../../requirements/L2-system/ch40-L2-system.md): REQ-L2-40-02, REQ-L2-40-16, REQ-L2-40-17, REQ-L2-40-18, REQ-L2-40-20; [L3](../../../requirements/L3-component/ch40-L3-component.md): REQ-L3-40-01, REQ-L3-40-17, REQ-L3-40-18, REQ-L3-40-19.

## 1. Мета та залежності

Задати окремі стани accepted/signed/persisted/anchored для двох типів записів — gap tickets і decision bundles. Квитанція має чесно вказувати досягнутий стан. Передумови: bounded arrival/outage, quota policy, failure model, canonical encoding, retention й payload-deletion policy.

[CH27 validator](RESEARCH-CH27-validator-contract.md) залишається власником inclusion vectors, [CH33 receipt](RESEARCH-CH33-RECEIPT-CONTRACT.md) — statement-binding, [CH33 replication](RESEARCH-CH33-REPLICATION-CONTRACT.md) — ordered delta/snapshot recovery. Їхні контракти повторно не реалізовувати; використати як interfaces/fixtures. Тут нова перевірка — скінченна локальна квота, durable acknowledgement рішення й crash-window **до WORM-анкорингу**, а не другий proof доставки нормативних дельт. Узгодженість вмісту кластера — [окреме завдання](RESEARCH-CH40-CLUSTER-CONSISTENCY.md), не prerequisite для малого crash-оракула.

## 2. Безпечний протокол і контролі

1. Запросити наявні WAL, epoch schema й storage adapter; визначити durable ACK, fsync/atomic rename, tail corruption, retry identity та межі гарантії ОС/носія.
2. Створити локальний synthetic storage simulator з точним журналом подій. Порівняти volatile ring baseline, reference durable log і кандидатний writer/recovery.
3. Ін'єктувати process crash до/після запису, fsync, підпису, ACK, додавання листа та anchor receipt. Втрата непроанкореної епохи не може повертати клієнту anchored-status.
4. Контролі квоти: лише дублікати; лише унікальні критичні події; заповнення точно на межі; один наступний запис; тривалий offline; повторна доставка; переставлені ACK; partial/corrupt WAL.
5. Перевірити shutdown під partition і відмовою резервного вузла: стан не стає SHUTDOWN_CLEAN без визначеної передумови; час очікування обмежений погодженим протоколом.
6. Для ledger перевірити epoch ID, sequence, tree size, підписані canonical bytes, duplicate/omitted leaf, два різні checkpoints, збереження proof після видалення payload. Inclusion не оголошувати доказом повноти чи відсутності equivocation.
7. Лише синтетичні файли у власному дослідницькому checkout й loopback-адаптери; без аварійного вимкнення спільної машини або руйнування реального носія. Ключі одноразові; персональних даних немає.

## 3. Цілі, PASS/FAIL і негативна публікація

500 МБ, 10 с, 1000 пакетів, 200 мкс/100 мс діалогового прикладу — **наявні параметри різних профілів**, не суперечність і не затверджений загальний SLA. Arrival/outage, loss budget, durability point і shutdown deadline **очікують погодження**.

**PASS:** кожна квитанція відновлюється згідно з названим станом, квота викликає визначену backpressure/refusal, history append-only, payload-deletion не фальсифікує proof availability. **FAIL:** acknowledged-loss у заявленій моделі, хибний anchored/clean, прихована нескінченна черга або rewrite history. Опублікувати всі crash traces, втрати й невиконані бюджети.

## 4. Артефакти й точні зміни

Передати state machine, canonical vectors, storage/failure manifest, simulator/seed, crash schedules, raw ACK/WAL/anchor traces, test keys, hashes й команди. Майбутній інтерфейс: `python3 reproduce.py --manifest manifest.json --suite durability-recovery`, не наявний API.

Оновити CH40 §6.3, §7.1–§7.2.1, §10 й висновки: quota exhaustion, exact receipt, anchoring window та recoverability; синхронізувати всі перелічені ID. Зберегти А-14 та А-25: історичні метадані незмінні, актуальний допуск може анулюватися, персональний payload видаляється окремо.
