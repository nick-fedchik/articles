# RESEARCH-CH35-REVOCATION-EPOCH. Скінченна лавина та атомарна заборона застарілого висновку

**ID:** RESEARCH-CH35-REVOCATION-EPOCH  
**Статус:** Відкрито; тип матриці, release protocol і межі часу не підтверджено.  
**Власник:** розробник VIDAdev JTMS/runtime; співвиконавець — дослідник VIDA конкурентних алгоритмів.  
**Пріоритет:** P0 — відкликаний засновок не повинен лишати придатний дозвіл.  
**Підстава:** [W33-35-23](../../audit-2026-10-28-ch33-ch35.md#w33-35-23), [W33-35-24](../../audit-2026-10-28-ch33-ch35.md#w33-35-24), [CH35 §3.3.2/§4/§10.1.1](../../../ch35-self-organizing-expert-systems-synergetics-and-npu-runtime.md), 221–228, 278–316, 986–991.

**Вимоги:** [L1-35](../../../requirements/L1-stakeholder/ch35-L1-stakeholder.md): `REQ-L1-35-02`, `REQ-L1-35-03`; [L2-35](../../../requirements/L2-system/ch35-L2-system.md): `REQ-L2-35-03`, `REQ-L2-35-04`, `REQ-L2-35-06`; [L3-35](../../../requirements/L3-component/ch35-L3-component.md): `REQ-L3-35-04`, `REQ-L3-35-06`, `REQ-L3-35-07`.

## 1. Мета й залежності

Довести дві різні властивості: обмеженість обробки скінченного графа при одноразовій інвалідації та атомарну логічну заборону випуску старого результату при epoch bump. Фізичне оновлення N статусів не називати O(1).

Передумови: семантика залежностей/JTMS, тип A (суміжність чи expectation operator), область епохи, memory model, всі release paths, cache/accelerator consumers, wraparound/restart policy.

[CH26-CONTRACTION-CERTIFICATE](RESEARCH-CH26-CONTRACTION-CERTIFICATE.md) і [CH20-frozen-trace](RESEARCH-CH20-frozen-trace-monotonicity.md) залишаються власниками логічної ревізії/оракулів. Тут не шукати повторно мінімальні ядра. Для мережевих споживачів використати [CH33 replication](RESEARCH-CH33-REPLICATION-CONTRACT.md) як interface assumptions; локальну модель можна почати незалежно.

## 2. Безпечний протокол

1. Отримати реальні interfaces VIDA й карту release paths; відокремити їх від навчального Go-прикладу. Зафіксувати implemented/target.
2. Побудувати exact reference invalidator із visited set. Порівняти recursive baseline, bounded worklist та bulkhead+epoch guard на однакових synthetic graphs.
3. Контролі спектральної моделі: двовузловий цикл із ρ=1 й двома одноразовими інвалідаціями; великий DAG із ρ=0 та широким каскадом; weighted nilpotent/non-normal operator; окрема branching-expectation model із явними probabilistic assumptions.
4. Перевірити, що convergence Neumann series не змішується з invertibility: A=2I має inverse I−A, але розбіжний ряд. Не відновлювати стару помилкову претензію до вже виправленого тексту.
5. Визначити linearization point epoch bump. Model checker перебирає reader start/read/validate/publish паралельно revocation, кешованому висновку й accelerator completion.
6. Контролі: reader у середині транзакції, old DMA result, пропущена epoch-check, перезапуск, epoch wrap/ABA, tombstone rollback, bulkhead boundary, альтернативне незалежне обґрунтування та revocation під час promotion.
7. Після revocation первинний запис має семантику `REVOKED`, похідні — `INVALIDATED_BY_REVOCATION`; фізична lazy materialization допустима лише за еквівалентної негайної release-заборони. Історична траса не стирається.
8. Виміряти окремо revoke→gate-close, revoke→all-materialized, cleanup та reverify. Усе локально на synthetic-графах і символічних дозволах; жодних приводів, мережевих production-вузлів чи змін реальної baseline.

## 3. Targets та критерії

**Логічна ціль:** після linearization point жоден release path не видає допустимий висновок зі старої епохи; усі finite controls відповідають reference invalidator.

`<5 мс` — наявна **проєктна** вимога, не підтверджений тут WCET. Hardware/workload/max graph/measurement endpoints та статус затвердження профілю **очікують погодження**. Кількість зразків і resource bounds затверджуються до вимірювань; максимальний спостережений час не є WCET-proof.

**PASS:** формальна linearizability/safety argument, повне покриття release paths, збіг з oracle, окремі часові метрики. **FAIL:** stale ACCEPT, пропущена epoch-check, invalidated source повертається через rollback або O(1) publication підміняє весь cascade.

Публікувати всі failing schedules, worst-case graph families та непідтверджені deadlines. Не приховувати лавину вилученням складних графів із корпусу.

## 4. Артефакти, команди та точні зміни

Передати graph/operator schema, invalidation oracle, formal protocol/memory assumptions, model-checker traces/bounds, generator/seed/manifest SHA-256, timings усіх фаз, raw failures і implemented/target inventory.

Майбутній інтерфейс; не наявна реалізація:

```text
python3 reproduce.py --manifest manifest.json --suite revocation-epoch
python3 verify_traces.py --manifest manifest.json --property no-stale-release
```

Оновити CH35 §3.3.2 модель A/finite cascade, §3.1/§4 logical/materialized statuses, §10.1.1 epoch/reader obligations. Узгодити всі перелічені вимоги, виправити parent `REQ-L3-35-04` і погодити прямі L3 на release/epoch. А-14 та незмінна історія А-25 лишаються чинними.
