# RESEARCH-CH40-CLUSTER-CONSISTENCY. Фаза як телеметрія, версія як умова допуску

**ID:** RESEARCH-CH40-CLUSTER-CONSISTENCY  
**Статус:** Відкрито; достатність фазового індикатора спростовується контрприкладом, замінний протокол не доведено.  
**Власник:** VIDA researcher розподілених протоколів; співвиконавець — VIDAdev кластера.  
**Пріоритет:** P0 — суперечливі критичні вердикти.  
**Підстава:** [W38-40-14](../../audit-2026-10-28-ch38-ch40.md#w38-40-14), [CH40 §7.3/§10, 728–755, 914–932](../../../ch40-distributed-epistemic-architectures-soa-and-cluster-scaling.md).  
**Вимоги:** [L1](../../../requirements/L1-stakeholder/ch40-L1-stakeholder.md): REQ-L1-40-09; [L2](../../../requirements/L2-system/ch40-L2-system.md): REQ-L2-40-15; [L3](../../../requirements/L3-component/ch40-L3-component.md): REQ-L3-40-16.

## 1. Мета та справжні залежності

Довести для скінченного протоколу, які content-bound epoch/quorum/fencing передумови виключають суперечливий критичний допуск. R_sync залишити вимірювачем фаз, не тотожності KB. До моделювання потрібні failure model, membership, allowable operations, trusted clock/lease assumptions і mapping канонічних вердиктів.

Довговічність квитанцій перевіряє [DURABILITY-CONTRACT](RESEARCH-CH40-DURABILITY-CONTRACT.md). [CH29 update](RESEARCH-CH29-update-noninterference.md) надає snapshot semantics; [CH35 revocation epoch](RESEARCH-CH35-REVOCATION-EPOCH.md) — локальний no-stale-release, [CH33 replication](RESEARCH-CH33-REPLICATION-CONTRACT.md) — ordered update/recovery. Повторно ці механізми не доводити: інтеграційний висновок залежить від їхніх interface assumptions; початковий phase/content контрприклад незалежний. Новий предмет — міжвузловий content-bound admission при однакових фазах, не локальний epoch bump. А-11 зберігає закритий ASIL-профіль: усі мережеві досліди тут — лише симулятор.

## 2. Спростовуваний безпечний протокол

1. Запросити фактичний protocol/state schema. Задати epoch ID, KB digest, membership revision, monotonic fencing token, quorum й запропонований інваріант узгодженого допуску.
2. Порівняти phase-only baseline, точний state-machine oracle та кандидатний протокол. Обов'язковий контроль: однакові фази, різні KB roots, R_sync=1; phase-only не повинен вважатися доказом узгодженості.
3. Зворотний контроль: той самий root, різні фази; не вигадувати змістовну колізію з самого малого R_sync. Перевірити partition, stale heartbeat, restart, duplicate messages, clock skew, membership change та delayed revocation.
4. У малих моделях перебрати reorder/drop/delay і знайти мінімальні unsafe traces. Окремо формулювати safety та liveness за припущення eventual delivery; timeout не доводить відмови вузла.
5. CRDT дозволяти для подій/метаданих лише з визначеним join; поточний немонотонний ASPIC+-вердикт переобчислювати за узгодженим snapshot. Після recovery старий запис лишається, але current admissibility відкликаного висновку блокується.
6. Ніяких ECU, фізичних приводів, production-кворумів або зовнішньої мережі; підписані синтетичні пакети й контрольований симулятор.

## 3. Цілі й PASS/FAIL

50 мс, 0,85/0,50, 16/64 вузлів і 1,2 мс — **книжкові приклади/кандидати**. Межі моделі, вид кворуму, clock/lease budgets і часові targets **очікують погодження**; фази не «калібрувати» так, щоб приховати content mismatch.

**PASS:** явний перевірний safety-invariant для погодженої моделі й окрема умовна liveness; content mismatch не дає критичного ACCEPT. **FAIL:** conflicting admission, stale fencing, revocation bypass або підміна proof фазовим threshold. Публікувати countermodel, відсутність liveness і межі доведення.

## 4. Відтворення й оновлення

Передати formal model, checker/version, config bounds, генератор/seed, proof або countertraces, digests/epochs і всі результати. Майбутній runner: `python3 reproduce.py --manifest manifest.json --suite cluster-consistency`; фактичну команду model checker додати після вибору інструмента.

Оновити CH40 §7.3, §10, NOTE й висновки: діагностика проти admission, припущення Курамото, CRDT scope та поточне відкликання. Синхронізувати три ID; зберегти А-04, А-07–А-09, А-11, А-13/А-14 та А-25.
