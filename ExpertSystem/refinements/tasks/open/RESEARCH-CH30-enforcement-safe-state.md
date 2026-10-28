# RESEARCH-CH30-enforcement-safe-state. Від permission gate до фізичного safe state

**ID:** RESEARCH-CH30-enforcement-safe-state  
**Статус:** Відкрито; permission API, alert-path і повний фізичний бюджет не підтверджено.  
**Власник:** розробник VIDA Linux/апаратного шлюзу та дослідник функціональної безпеки.  
**Пріоритет:** P0 — хибна гарантія блокування до небезпеки.  
**Підстава:** [аудит W27-30-17, W27-30-18](../../audit-2026-10-28-ch27-ch30.md), [CH30 §1.2–§1.2.1, §2.2](../../../ch30-safety-cybersecurity-co-engineering.md).  
**Вимоги:** [L1](../../../requirements/L1-stakeholder/ch30-L1-stakeholder.md): REQ-L1-30-02, REQ-L1-30-03, REQ-L1-30-04; [L2](../../../requirements/L2-system/ch30-L2-system.md): REQ-L2-30-03, REQ-L2-30-04, REQ-L2-30-05; [L3](../../../requirements/L3-component/ch30-L3-component.md): REQ-L3-30-01, REQ-L3-30-02, REQ-L3-30-03, REQ-L3-30-04.

## 1. Залежності та безпечні межі

Слід незалежно перевірити три твердження: заборона операції до ефекту, доставка alert, досягнення безпечного фізичного стану. `sys_enter`-телеметрія із SIGKILL не тотожна permission-hook; відсічення енергії не зупиняє інерцію миттєво.

Передумови: погоджені kernel/version/hooks, threat model, signal polarity, незалежність hardware trip і plant/FTTI. [CH18 WCET/DMA](RESEARCH-CH18-wcet-dma-proof.md) надає апаратний budget, [Simplex recoverability](RESEARCH-CH28-simplex-recoverability.md) — суміжну модель recovery; інтерфейс можна перевіряти незалежно від цих результатів.

## 2. Спростовуваний протокол без експлуатації вразливостей

1. Спершу code-review фактичних hook definitions/return contracts; зіставити observation, permission decision і hardware alert. Зафіксувати, які операції реально покрито; не називати tracepoint LSM-hook.
2. У власній одноразовій ізольованій VM запускати лише нешкідливі тестові операції проти синтетичних ресурсів. Контролі: дозволена/заборонена дія, спостереження без відмови, помилка політики, затриманий alert, недоступний consumer. Не застосовувати exploits, privilege escalation або тестів на спільному хості.
3. Hardware output замінити індикатором/logic simulator. Вимірювати окремі інтервали за спільним початком відліку; не вважати відправлення сигналу підтвердженням його фізичного ефекту.
4. У моделі механіки додати coast/braking, навантаження й похибки; перевірити nominal, граничну інерцію, втрату сигналу, відмову контактора та fail-operational режим. Не підключати небезпечний маніпулятор і не випробовувати біля людей.
5. Для SMT окремо перевірити timeout/unknown/error та відсутність виходу команди до позитивного рішення.

## 3. Targets, PASS/FAIL

1,5 мкс detection, 5 мс trip, 50 мс SMT і приклад 4,2012 проти 83,3 мс — **книжкові заяви/цілі для підтвердження**, а не виконані вимірювання. Застосовність і похибка очікують попереднього погодження; нових порогів немає.

**PASS:** кожен шар має явний контракт, негативні контролі, незалежну safe-state модель й узгоджений повний budget. **FAIL:** заборонена операція має ефект, alert губиться або safe state недосяжний до погодженої межі. Публікувати негативні результати й обмеження, не приховувати їх коротким relay timing.

## 4. Артефакти, команди й оновлення

Передати kernel/config/policy digests, таблицю hooks, synthetic harness, VM recipe, timestamps, logic traces, plant model і розрахунок фаз. Фактичні build/test-команди надати для вибраного ізольованого стенда; на робочому хості не виконувати privileged loader.

Безпечна арифметична baseline-команда не завантажує eBPF і не керує обладнанням:

```bash
python3 -c 'trip_ms=0.0012+4.2; travel_ms=1000*0.15/1.8; assert trip_ms<travel_ms; print("trip_ms",trip_ms,"travel_ms",travel_ms,"mechanical recovery still unmodelled")'
```

Оновити CH30 §1.2 (enforcement проти observation), §1.2.1 (поля сигналів, початок відліку, повна реакція), §2.2 (узгодження фаз), словник LSM. Синхронізувати всі перелічені вимоги, виправивши «замикати/розмикати» й позитивність FTTI. Не робити універсального наказу знеструмлювати кожен виріб.
