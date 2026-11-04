# RESEARCH-CH35-ENERGY-METROLOGY. Вимірюваний estimand Дж на верифікований факт

**ID:** RESEARCH-CH35-ENERGY-METROLOGY  
**Статус:** Відкрито; estimand, метрологічна невизначеність і профіль порівняння не погоджено.  
**Власник:** дослідник VIDA метрології; співвиконавець — розробник VIDAdev телеметрії/обчислювальних профілів.  
**Пріоритет:** P1 — невідтворювана основа енергетичного планувальника.  
**Підстава:** [W33-35-25](../../audit-2026-10-28-ch33-ch35.md#w33-35-25), [CH35 §6.1–§6.2.1](../../../ch35-self-organizing-expert-systems-synergetics-and-npu-runtime.md), 431–528.

**Вимоги:** [L1-35](../../../requirements/L1-stakeholder/ch35-L1-stakeholder.md): `REQ-L1-35-08`, `REQ-L1-35-11`; [L2-35](../../../requirements/L2-system/ch35-L2-system.md): `REQ-L2-35-15`, `REQ-L2-35-16`, `REQ-L2-35-20`; [L3-35](../../../requirements/L3-component/ch35-L3-component.md): `REQ-L3-35-14`, `REQ-L3-35-19`.

## 1. Мета й недубльовані залежності

Розділити енергію на вхідну подію, на спробу, на виданий верифікований факт і повну енергію системи. Довести правильність estimator на synthetic power trace; потім перевірити claimed фізичні значення за контрольованим протоколом.

Передумови: визначення вихідного факту/відмови, batch accounting, idle/active boundaries, sensor rails/units/дискретизація, калібрування та однакова робота CPU/NPU/FPGA.

Спільний опис Orin/DVFS/thermal limits брати з [CH18-jetson-thermal-reproducibility](RESEARCH-CH18-jetson-thermal-reproducibility.md), не виконувати вдруге загальний thermal/throughput benchmark. [TST-010](TASK-TST-010-performance-evidence.md) не є окремим INA/Дж-per-output протоколом. Арифметичні помилки W26 виправляються документально без очікування вимірювання.

## 2. Безпечний протокол

1. Отримати первинні артефакти вже заявлених автором випробувань, якщо вони доступні: дата/код/model versions, workload, nvpmodel, rails, sensor configuration, raw timestamps і лічильники фактів. Відсутність їх у главі не доводить, що дослідів не було.
2. Побудувати аналітичні synthetic traces із відомим інтегралом: constant idle, відомий burst, короткий імпульс між polls, overlapping bursts, clock jitter, missing sample та drift idle.
3. Базові лінії: формула per-input із §6.1; точний інтеграл синтетичного сигналу; discrete integrator із idle subtraction; повна system energy без віднімання. Не змішувати ці різні estimands.
4. На однаковому вхідному потоці варіювати prune, частку успішного verification та batching. Контролі `N_verified=0`, усі входи відхилені, один вхід породжує кілька фактів, повторний факт і proof failure. Дж/вихід має мати правильний знаменник та undefined-case.
5. Для INA3221 перевірити sensor conversion/averaging/update period і фактичні rails у конкретній документації платформи; 100-Гц polling не називати 100 незалежними миттєвими вимірами без перевірки конфігурації. Оцінити quantization, drift, aliasing і синхронізацію burst.
6. Фізична фаза — тільки після окремого погодження власником обладнання/режиму: пасивне вимірювання безпечного synthetic workload у дозволеному power mode. Не змінювати firmware, напруги, thermal protection, керування приводами або реальну базу знань; не проводити stress/fault injection.
7. CPU/NPU/FPGA порівнювати тільки за однаковим перевірним входом/виходом/сертифікатом і scope power. Якщо workloads різні, публікувати окремі результати, не таблицю «переваги» з одним знаменником.
8. Реєструвати total та incremental power, energy/attempt/output, тривалість, лічильники відмов, confidence/uncertainty і всі повтори. Для автономності додати idle/інших споживачів; не виводити температуру лише з Вт.

## 3. Targets, PASS/FAIL та негативна публікація

**Логічний target:** estimator збігається з відомим synthetic integral в заявленому чисельному допуску; denominator і scope видимі. **Допуск, sampling/repeat plan та потрібна uncertainty очікують погодження до прогону.**

12,4 мкДж, 25 мкДж, 1,85 Дж, idle 3,4/5,2 Вт, 15 мкДж-тригер і 48,5 °C — **заяви/проєктні параметри книги**, не підтверджені targets цього завдання. Для нормативного використання потрібне окреме затвердження hardware/workload profile та невизначеності.

**PASS:** відтворюваний estimator, первинні фізичні traces або чесний статус «поки тільки модель», коректні units/denominators і порівнюваність. **FAIL:** відкинуті входи рахують як видані факти, аліасинг приховано, різну роботу ділять як однакову, або incremental power видають за повне споживання.

Опублікувати missing evidence, негативний результат щодо енергетичної переваги, всі повтори та невизначеність. Невдалий benchmark не приховувати новим постфактум порогом.

## 4. Артефакти, команди та точні оновлення

Передати estimator spec, event/fact accounting schema, synthetic traces/інтеграли, calibration/sensor config, hardware/workload manifest SHA-256, raw rails/timestamps, всі repeats і аналіз похибки; жодних секретних графів.

Команди — інтерфейс майбутнього пакета, ще не виконані:

```text
python3 reproduce.py --manifest manifest.json --suite energy-estimator
python3 integrate_energy.py --manifest manifest.json --report total-incremental --denominator verified-output
```

Оновити CH35 §6.1 estimand/формулу/приклад, §6.2 метод/одиниці/невизначеність, §6.2.1 таблицю та runtime trigger; розвести оцінку й вимірювання за А-24. Узгодити всі перелічені ID, особливо polling/denominator `REQ-L3-35-19` і trigger `REQ-L3-35-14`; нові thresholds — лише після погодження А-09.
