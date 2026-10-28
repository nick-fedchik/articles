# RESEARCH-CH36-SEALED-GATE-PROTOCOL. Незалежність запечатаного приймального набору

**ID:** RESEARCH-CH36-SEALED-GATE-PROTOCOL.
**Статус:** Відкрито; split, ролі та політика повторного використання ще не погоджені.
**Власник виконання:** дослідник VIDA експериментальної методології; співвиконавець — розробник VIDA приймального шлюзу.
**Пріоритет:** P1 — залежна оцінка може створити хибну видимість надійності.
**Підстава:** [W36-37-09](../../audit-2026-10-28-ch36-ch37.md#w36-37-09), [глава 36 §5.4/§5.6](../../../ch36-knowledge-testing-pyramid-and-variational-calibration.md), рядки 440, 518–556.
**Вимоги:** [L1](../../../requirements/L1-stakeholder/ch36-L1-stakeholder.md#req-l1-36-13) `REQ-L1-36-13`; [L2](../../../requirements/L2-system/ch36-L2-system.md#req-l2-36-15) `REQ-L2-36-15`; [L3](../../../requirements/L3-component/ch36-L3-component.md#req-l3-36-18) `REQ-L3-36-18`.
**Залежності:** визначена одиниця незалежного сценарію, ролі developer/evaluator, policy доступу/звітів, модель TEE. Саму побайтову кустодію, [SRR](TASK-TST-008-srr-benchmark.md) та [CH16 calibration](RESEARCH-CH16-quality-gates.md) не досліджувати повторно.

## 1. Гіпотеза

Перевірити, що конкретний протокол запобігає перенесенню оцінювальних випадків у підбір системи та явно обмежує feedback від gate. AES-GCM/Merkle забезпечують власні властивості цілісності, але не доводять відсутності сценарного overlap.

Спершу отримати наявний gate/manifest VIDA та попередні рішення. Робота не припускає відсутності коду або апаратного анклаву.

## 2. Спростовуваний синтетичний протокол

1. Створити вигаданий корпус сценаріїв і похідних парафраз із відомими family IDs. Порівняти random-row split з family-level split; відомі похідні копії мають бути виявлені до оприлюднення результату.
2. Описати інформаційні потоки до developer, evaluator і tested model: хто бачить запит, gold, per-case verdict, aggregate та журнал. Зафіксувати кількість дозволених rounds і умову заміни sealed-набору.
3. Негативні контролі: зашифрований, але дубльований development-case; прихований парафраз тієї ж родини; підбір за багаторазовим aggregate feedback; зміна manifest після першого прогону; неавторизований читач.
4. Позитивні baseline: family-disjoint frozen corpus і дозволений одноразовий evaluator. Контроль цілісності: легітимна перевірка та навмисна зміна одного synthetic байта мають різні вердикти.
5. Локальний simulator перевіряє policy/state transitions, але **не доводить TEE confidentiality**. Якщо потрібна апаратна заява, окремо надати attestation/config/threat-model evidence без доступу до сторонніх секретів або пошуку escape.
6. Не вивантажувати приватних тестів, не підключати production-моделі до held-out відповідей; тестові ключі вигадані й не є реальними credentials.

## 3. Цілі та PASS/FAIL

Заявлені «100% захищений» і «недоступний» потребують обмеженої моделі. Допустимі feedback, число rounds, обсяг корпусу й емпіричні leakage targets **очікують пререєстрації**. Нового універсального порога не задавати.

**PASS:** frozen split/roles/feedback policy відтворюються, всі зареєстровані overlap/tamper/access controls правильно блокуються, hardware та simulator evidence чітко розділені. **FAIL:** залежний випадок зараховано як незалежний, developer може підбирати за прихованими відповідями або policy змінено після результату. Публікувати overlap, невдалі rounds і непідтверджені hardware claims; не називати успішне шифрування доказом відсутності contamination.

## 4. Артефакти, команди й оновлення

Передати synthetic families/seed/хеші, train/development/calibration/held-out manifests, role/access matrix, transcript усіх rounds, контрприклади й мінімальні конфігурації. Майбутній runner:

```text
python3 reproduce.py --manifest manifest.json --suite sealed-gate-protocol
python3 verify_controls.py --manifest manifest.json --property family-disjoint-and-frozen
```

Оновити главу 36 §5.4/§5.6, схему шлюзу та перелічені вимоги: розділити custody, access control, split independence і hardware assumptions. Відсутність незалежної оцінки позначити явно; не підміняти її статистичним ZHR=1.
