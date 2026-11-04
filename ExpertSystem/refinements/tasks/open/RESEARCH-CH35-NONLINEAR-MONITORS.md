# RESEARCH-CH35-NONLINEAR-MONITORS. Спостережувані моделі й захисні переходи

**ID:** RESEARCH-CH35-NONLINEAR-MONITORS  
**Статус:** Відкрито; моделі, identification і цілі false-alarm/lead-time не затверджено.  
**Власник:** дослідник VIDA динамічних систем; співвиконавець — розробник VIDAdev runtime-моніторів.  
**Пріоритет:** P1 — математична аналогія не доводить випереджальний захист.  
**Підстава:** [W33-35-22](../../audit-2026-10-28-ch33-ch35.md#w33-35-22), [CH35 §3.3/§5.2](../../../ch35-self-organizing-expert-systems-synergetics-and-npu-runtime.md), 203–277, 350–399.

**Вимоги:** [L1-35](../../../requirements/L1-stakeholder/ch35-L1-stakeholder.md): `REQ-L1-35-09`; [L2-35](../../../requirements/L2-system/ch35-L2-system.md): `REQ-L2-35-12`, `REQ-L2-35-17`, `REQ-L2-35-18`; [L3-35](../../../requirements/L3-component/ch35-L3-component.md): `REQ-L3-35-11`, `REQ-L3-35-15`, `REQ-L3-35-16`, `REQ-L3-35-17`.

## 1. Мета, негайні правки та залежності

Перевірити вузькі моделі hysteresis/CSD/blow-up, а не універсальну «синергетичну гарантію». До досліду погодити числове відображення параметра порядку, state dynamics, одиниці, window/estimators, observation noise та допустимий fallback.

Умова L3-35-16 про ширину замість поточного b виправляється документально. Дискримінант/ширина cusp вже правильні; відхилену Б1-20 про знак Ляпунова не відкривати.

[CH26-EARLY-WARNING](RESEARCH-CH26-EARLY-WARNING.md) лишається власником CUSUM-протоколу. Використати спільні правила оцінювання false alarms/misses, але тут предмет — CSD, identification blow-up і конкретні state guards. [CH18-WCET/DMA](RESEARCH-CH18-wcet-dma-proof.md) — власник апаратних ресурсних гарантій, не цей квиток.

## 2. Безпечний synthetic-протокол

1. Отримати наявний код/профіль VIDA; відокремити implemented, target і аналогію. Дати рівняння стану/вимірювань, ranges, SI/nondimensional mapping та незмінні safety invariants.
2. **Hysteresis:** порівняти single threshold, fixed Schmitt baseline та погоджений cusp-derived guard. Відтворити a=−0,36 і ширину ≈0,166; при достатній налаштованій ширині, але b=−0,25 вихід із Fail-Closed заборонений. Перевірити current b, memory попереднього стану, boundary equality, noisy crossing і invalid a.
3. **CSD:** порівняти variance-only ratio, autocorrelation/relaxation estimator і exact known-dynamics oracle. Контролі: зростання амплітуди білого шуму без slow-down; змінений relaxation time при сталій variance; сезонність, jump без precursor, нульова baseline variance й enum без числової карти.
4. **Blow-up:** на додатній synthetic ODE з відомими c/β/початковим станом порівняти exact time, three-sample estimator і простий trend/threshold baseline. Додати noise, quantization, irregular sampling, β≤1, ψ≤0, saturation і невідповідну модель. Невизначений прогноз не отримує точного `T_TTF`.
5. Для entropy/Lyapunov claims спочатку визначити конкретну функцію на станах KB та transition law. Довести умови для цього піднабору або чесно лишити аналогією; не підміняти доказ знаками фізичної ентропії.
6. Pre-register train/evaluation split, seeds, вікна, error bars й multiplicity. Вимірювати false alarms, misses, lead-time і decision chattering, включно з невдалими/невизначеними серіями.
7. Всі переходи — символічні записи, без реле/силових кіл, теплового навантаження чи fault injection реального обладнання. Жорсткі інваріанти active за будь-якого γ; оптимізаційні пріоритети не можуть їх послабити.

## 3. Targets, PASS/FAIL і негативний результат

Пороги γ=3, a=−0,20, множник 3·T_brake і 120 мс — **наявні проєктні числа**, не відкалібровані цим завданням. Їхнє схвалення для конкретної системи, noise bounds, допустимі false alarms/misses, мінімальний lead-time і ресурсні maxima **очікують погодження до прогону**.

**Логічний PASS:** guard не допускає виходу лише через configured width; invalid model/state дає явну відмову; proofs обмежені названою динамікою.

**Емпіричний PASS:** виконано наперед погоджені цілі на незалежних synthetic-потоках. **FAIL:** false CSD подано як доведене наближення аварії, нестійкий estimator дає гарантований час, або safe policy вимкнена до γ≥3.

Публікувати всі negative/no-precursor series й випадки, де простий baseline кращий. Відсутність універсальної гарантії — припустимий науковий результат.

## 4. Артефакти, команди й зміни

Передати модель/одиниці, formal guards, oracle/generator/seed, frozen manifest SHA-256, estimation uncertainty, alarm timelines, chattering counts, proof/контрприклади й implemented/target.

Майбутній пакет, не вже виконані команди:

```text
python3 reproduce.py --manifest manifest.json --suite nonlinear-monitors
python3 verify_guards.py --manifest manifest.json --include-invalid-models
```

Оновити CH35 §3.3.1–§3.3.4 assumptions/змінні/інженерні рішення, §5.2 numeric mapping/window, приклади й dashboard. Синхронізувати всі перелічені ID, особливо state-dependent `REQ-L3-35-16`; thresholds baseline змінюються лише після А-09-погодження, без нових per-action підписів.
