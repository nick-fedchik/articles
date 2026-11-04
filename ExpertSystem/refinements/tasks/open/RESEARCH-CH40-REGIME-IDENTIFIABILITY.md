# RESEARCH-CH40-REGIME-IDENTIFIABILITY. Межі діагностики за двома нев'язками

**ID:** RESEARCH-CH40-REGIME-IDENTIFIABILITY  
**Статус:** Відкрито; спостережуваність, шум і повний автомат причин не погоджені.  
**Власник:** VIDA researcher діагностики систем; співвиконавець — VIDAdev монітора режимів.  
**Пріоритет:** P0 — хибна впевненість у причині відхилення.  
**Підстава:** [W38-40-16](../../audit-2026-10-28-ch38-ch40.md#w38-40-16), [CH40 §9.4, 884–898](../../../ch40-distributed-epistemic-architectures-soa-and-cluster-scaling.md).  
**Вимоги:** [L1](../../../requirements/L1-stakeholder/ch40-L1-stakeholder.md): REQ-L1-40-04; [L2](../../../requirements/L2-system/ch40-L2-system.md): REQ-L2-40-08; [L3](../../../requirements/L3-component/ch40-L3-component.md): REQ-L3-40-09.

## 1. Мета та залежності

Встановити, які причини можна розрізнити за parity/conservation residuals, а які вимагають ambiguity/unknown або додаткового спостереження. Передумови: модель об'єкта/спостереження, одиниці, часові координати, класи причин і припущення незалежності шуму.

[CH21 physical postconditions](RESEARCH-CH21-PHYSICAL-POSTCONDITIONS.md) уже перевіряє досягнення стану; тут окрема задача ідентифікації причини. Доказ безпечної коригувальної дії не виводиться з назви діагнозу.

## 2. Спростовуваний безпечний протокол

1. Запросити наявний discriminator, профіль τ й fault dictionary. Формально задати mapping «причина→спостереження→дві нев'язки» та повну таблицю переходів.
2. Створити малий synthetic simulator із відомою прихованою причиною. Базові лінії: надрукований three-state rule, oracle множини сумісних причин і кандидатний monitor.
3. Пара нерозрізненності: однаковий bias двох давачів і фізичний зсув об'єкта дають однакові residuals. Якщо додаткового сигналу немає, кандидат не має стверджувати єдиний OOD-діагноз.
4. Перебрати всі комбінації high/low parity, high/low conservation, known/unknown fault match; додати boundary equality, stale/frozen sensor, time skew, missing value, відмову обох каналів і похибку фізичної моделі.
5. Для кожної величини визначити одиницю, σ окремого й різницевого сигналу, кореляцію та normalization. 3σ за негаусового шуму не оголошувати універсальною частотою помилок.
6. Верифікувати, що ambiguity/unknown спричиняє дозволений baseline-стан без точного вигаданого діагнозу. Лише симулятор, без зміни фізичних захисних конвертів або обладнання.

## 3. Цілі та PASS/FAIL

`3σ_sensor`/`3σ_model` походять із L3, не з доведеного шумового профілю глави. Їхню застосовність, permissible error rate, частоту дискретизації, N і ресурси **потрібно погодити до експерименту**; нових числових thresholds немає.

**PASS:** доведена обмежена ідентифікованість або чесно повернута множина причин/unknown, повна таблиця й незалежні controls. **FAIL:** нерозрізненні причини отримують удаваний однозначний статус, case випадає з автомата або σ не має одиниці. Публікувати confusion/ambiguity matrices і негативний результат про неідентифікованість.

## 4. Відтворення та оновлення

Надати observation model, fault dictionary, units/covariance, synthetic generator/seed, exact compatible-cause oracle, hashes/versions, raw traces і команди. Майбутній runner: `python3 reproduce.py --manifest manifest.json --suite regime-identifiability`; його реалізація ще не заявляється.

Оновити CH40 §9.4, пункт 4/NOTE/інваріант: достатні припущення, ambiguity, units і worked example; синхронізувати три ID, відділивши категоричну потребу від досяжного діагностичного контракту. Зберегти А-04, А-07–А-09 й А-17.
