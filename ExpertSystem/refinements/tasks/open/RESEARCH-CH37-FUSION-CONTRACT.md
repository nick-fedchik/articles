# RESEARCH-CH37-FUSION-CONTRACT. Передумови узгодженого злиття свідчень

**ID:** RESEARCH-CH37-FUSION-CONTRACT.
**Статус:** Відкрито; probabilistic model і consistency certificate ще не погоджено.
**Власник виконання:** дослідник VIDA числового злиття свідчень; співвиконавець — розробник VIDA provenance/fusion adapter.
**Пріоритет:** P0 — хибна надмірна впевненість після злиття.
**Підстава:** [W36-37-15](../../audit-2026-10-28-ch36-ch37.md#w36-37-15), [глава 37 §7.2/§7.4](../../../ch37-input-information-assessment-and-algorithmic-skepticism.md), рядки 391–398, 492–501, 1111.
**Вимоги:** [L1](../../../requirements/L1-stakeholder/ch37-L1-stakeholder.md#req-l1-37-03) `REQ-L1-37-03`; [L2](../../../requirements/L2-system/ch37-L2-system.md) `REQ-L2-37-05`, `REQ-L2-37-06`; [L3](../../../requirements/L3-component/ch37-L3-component.md) `REQ-L3-37-04`, `REQ-L3-37-05`.
**Залежності:** однаковий оцінюваний стан/система координат, одиниці, модель середнього й похибок, правила встановлення GroupID. [CH26 OPE](RESEARCH-CH26-OPE-VALIDITY.md) і [CH22 RF](RESEARCH-CH22-RF-EVIDENCE.md) не повторюються: вони перевіряють інші estimands.

## 1. Спростовуваний контракт

Визначити admission conditions одного fusion adapter: незалежне LR-злиття дозволене лише за потрібної факторизації під обома гіпотезами; CI повертає коректну uncertainty bound за своїх consistency assumptions, але не доводить правдивості вхідних оцінок і не виявляє конфлікт за trace автоматично.

Спершу отримати наявні код/схеми VIDA, calibration evidence і тести. Не приписувати продуктові відсутність реалізації.

## 2. Протокол, baseline та контролі

1. На скінченній synthetic binary-моделі побудувати exact joint-table oracle. Контролі: незалежність під H і ¬H; незалежність лише під H; дублікати; спільна прихована причина; різні GroupID для одного походження.
2. Порівняти exact conditional-LR, наївний добуток marginal LR та candidate adapter. Перейменування GroupID не повинне створювати доказу фізичної незалежності.
3. Для CI зафіксувати state vector, units, SPD і визначення consistency як верхньої межі похибки за названих припущень. Для trace різнорозмірного вектора задати нормування/ваги; для log-odds — основу логарифма та одиниці. Довести матричну властивість у коректному порядку, не використовувати невизначений `P1 ∩ P2`.
4. Порівняти CI з independent-fusion baseline та synthetic exact covariance oracle. Пройти допустиму сітку cross-correlation, включно з повною кореляцією й завідомо заниженими вхідними uncertainty bounds.
5. Відтворити P₁=P₂=4 → P_CI=4; потім подати x₁=0, x₂=100 з тими самими P. Trace не зміниться; потрібна окрема узгодженість середніх або чесний unknown, а не фальшиве підтвердження збігу.
6. Контролі: різні фізичні стани, зміна одиниць, singular/non-SPD covariance, stale calibration, NaN, порожній набір. Усі дані синтетичні й локальні; без випромінювання, реальної телеметрії або фізичного керування.

## 3. Пререєстрація та результати

Точність floating-point, розмір synthetic grid, емпіричний consistency/coverage test, admissible conflict і часові цілі **очікують погодження**. Універсального «порога дисперсії ISO 26262» не вигадувати. Доказ consistency відділяється від емпіричної частоти покриття.

**PASS:** кожен admissible випадок відповідає exact oracle/формальному bound у погодженій області, invalid assumptions не дають безумовного допуску. **FAIL:** marginal independence під однією H видано за достатню, GroupID використано як proof або низький trace приховав конфлікт. Публікувати контрмоделі, невиконані цілі й результати з порушеними припущеннями.

## 4. Артефакти, команди й зміни

Передати joint tables, covariance generator/seed, state/units schema, формальний consistency proof, assumption checks, хеші/версії, raw outputs і negative report. Майбутній інтерфейс пакета:

```text
python3 reproduce.py --manifest manifest.json --suite evidence-fusion
python3 verify_controls.py --manifest manifest.json --suite independence-and-consistency
```

Оновити главу 37 §7.2/§7.4 і §12: обидві умовні факторизації, математично коректну CI-гарантію та окремий conflict test. Синхронізувати перелічені L1/L2/L3; не перетворювати scalar confidence на детермінований ASIL-доказ.
