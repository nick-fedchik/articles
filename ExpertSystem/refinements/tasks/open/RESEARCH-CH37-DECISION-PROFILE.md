# RESEARCH-CH37-DECISION-PROFILE. Рішення обмеженого монітора за ризиком і timeout

**ID:** RESEARCH-CH37-DECISION-PROFILE.
**Статус:** Відкрито; loss/hazard profile і правила багаторазового моніторингу ще не погоджено.
**Власник виконання:** дослідник VIDA послідовного оцінювання; співвиконавець — розробник VIDA шлюзу дозволу дій.
**Пріоритет:** P0 — недостатність доказу може помилково блокувати потрібний захист.
**Підстава:** [W36-37-16 / W36-37-18](../../audit-2026-10-28-ch36-ch37.md), [глава 37 §7.1/§8/§9.1](../../../ch37-input-information-assessment-and-algorithmic-skepticism.md), рядки 377, 480, 595, 675–696.
**Вимоги:** [L1](../../../requirements/L1-stakeholder/ch37-L1-stakeholder.md) `REQ-L1-37-04`, `REQ-L1-37-11`; [L2](../../../requirements/L2-system/ch37-L2-system.md) `REQ-L2-37-08`, `REQ-L2-37-19`; [L3](../../../requirements/L3-component/ch37-L3-component.md) `REQ-L3-37-08`, `REQ-L3-37-09`, `REQ-L3-37-19`.
**Залежності:** зафіксовані H0/H1, loss matrix, спостережуваність, sample period, horizon/reset, N_max та baseline-політика fallback. Обов'язки фізичного safe state лишаються поза числовим монітором; не дублювати [CH26 early warning](RESEARCH-CH26-EARLY-WARNING.md).

## 1. Гіпотеза та обмеження

Перевірити один формальний контракт bounded monitor: він видає оцінку/timeout і дозволяє лише дію, передбачену затвердженим hazard-профілем. `DEFER` не означає бездіяльності фізичного захисту, а статистичне рішення не доводить істинності H.

Запросити наявні алгоритми/політики VIDA, а не вважати їх відсутніми. Усі дії в досліді — mock-команди локальному simulator, без SCADA/ECU/актуаторів.

## 2. Протокол та контролі

1. Відтворити P(H)=0,01, sensitivity=0,90, false alarm=0,10 → posterior=1/12. Для C_FP=1 і C_FN=99 exact oracle дає втрату дії 11/12 проти 8,25 без дії; для зворотних втрат рішення змінюється.
2. Порівняти blanket probability gate, очікувану втрату за погодженим профілем та детермінований safety fallback. Окремо журналювати «не доведено як факт», «дозволено запобіжну дію» і «потрібне уточнення».
3. Для SPRT побудувати просту скінченну synthetic модель з відомими f₀/f₁, включивши i.i.d. baseline і корельовані відліки. Перевірити ±4,595 при alpha=beta=0,01 як наближення, не точну універсальну межу.
4. Незалежним dynamic-programming/exhaustive oracle перевірити ймовірності трьох результатів `H0/H1/INCONCLUSIVE` до N_max. Окремо оцінити багаторазовий reset протягом фіксованого горизонту; per-test alpha не видавати за lifetime risk.
5. Негативні контролі: невідомі/нульові likelihood, модель поза областю, missing label, накопичення без перетину порогів, запізнілі дані, високий hazard при низькому posterior, інтервал через поріг. Врахувати час здобуття відліків, не тільки CPU.
6. За timeout перевірити попередньо дозволений degraded/fallback state. Не вимагати інтерактивного підпису на кожну дію в затвердженому baseline.

## 3. Цілі, PASS/FAIL та негативна публікація

Loss matrix, N_max, lifetime false-alarm/miss/inconclusive budgets і допустима затримка **очікують пререєстрації** власником. 0,95/0,05, alpha/beta=0,01 та 20 мс у книзі не стають цим дослідом універсальними targets. Рішення А-04/А-08/А-16 не змінювати.

**PASS:** exact controls відтворені, типи оцінки/дії розділені, error bounds відповідають horizon/reset, timeout не активує випадкової або недозволеної дії. **FAIL:** lower-risk action заблоковано без hazard-підстави, перетин порога названо детермінованою істиною або приховано inconclusive/miss. Публікувати всі помилки, невигідні профілі й неможливість виконати одночасні цілі.

## 4. Артефакти, команди й оновлення

Передати policy/state machine, f₀/f₁/loss definitions, synthetic generator/seed, exact oracle, version/hash manifest, raw trajectories, time/error matrices та negative report. Майбутній пакет:

```text
python3 reproduce.py --manifest manifest.json --suite bounded-monitor-policy
python3 verify_controls.py --manifest manifest.json --suite risk-timeout-reset
```

Оновити главу 37 §7.1/§7.3/§8/§9.1, приклад 8,3%, усічення й reset. Узгодити зазначені вимоги, розділивши document S_doc, оцінку свідчень і actuator authorization; конкретні safety targets затверджує власник до фінального тесту.
