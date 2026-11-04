# RESEARCH-CH33-RECEIPT-CONTRACT. Перевірна квитанція федеративного доказу

**ID:** RESEARCH-CH33-RECEIPT-CONTRACT  
**Статус:** Відкрито; statement schema, trust model і часовий профіль очікують погодження.  
**Власник:** дослідник VIDA proof-carrying systems; співвиконавець — розробник VIDAdev шлюзу доказів.  
**Пріоритет:** P0 — неповна умова допуску доказового твердження.  
**Підстава:** [W33-35-04](../../audit-2026-10-28-ch33-ch35.md#w33-35-04), [CH33 §6.3](../../../ch33-inter-system-knowledge-exchange-and-model-teaching.md), 456–501.

**Вимоги:** [L1-33](../../../requirements/L1-stakeholder/ch33-L1-stakeholder.md): `REQ-L1-33-12`; [L2-33](../../../requirements/L2-system/ch33-L2-system.md): `REQ-L2-33-15`; [L3-33](../../../requirements/L3-component/ch33-L3-component.md): `REQ-L3-33-16`.

## 1. Гіпотеза та залежності

За погодженої trust model квитанція може засвідчувати **конкретну теорему про конкретну ревізію**, а не лише справжність окремого підпису й синтаксичну коректність proof blob.

До proof потрібні: канонічне типізоване кодування, ClaimID/statement/specification hash, policy version, issuer/key purpose, audience, freshness epoch/sequence, replay identity, формат proof і довірений checker. Секретний предикат, його commitment та відкритий statement мають мати явний зв'язок. Не обіцяти zero knowledge лише через наявність хешу.

Абстракцію цифрового підпису спочатку моделювати символічно; криптографічну примітивну стійкість не «доводити» тестами. Протокол доставки — [CH33-REPLICATION-CONTRACT](RESEARCH-CH33-REPLICATION-CONTRACT.md), але логічні квитанції можна перевіряти незалежно від брокера.

## 2. Безпечний протокол і базові лінії

1. Отримати наявну schema/checker API VIDA й відокремити implemented/target. Не підключатися до зовнішніх анклавів та не використовувати справжні приватні ключі/закриті специфікації.
2. Побудувати маленькі синтетичні специфікації та незалежний theorem checker із відкритим statement. Базові лінії: signature-only; signature+syntax; повне зв'язування signed statement/proof/policy.
3. Позитивний контроль: підписаний кортеж, proof саме заявленої теореми, відповідна ревізія й припустимий часовий стан.
4. Негативні контроли: доказ іншої ClaimID; правильний proof іншої ревізії; proof про слабшу властивість; змінений commitment; інший audience/policy; недовірений issuer; невідомий proof format; checker timeout/unknown.
5. Для часу відтворити майбутню мітку, прострочення, replay у межах freshness, дубль після restart, clock rollback і допустимий clock skew. Вік пакета й витрати CPU рахувати окремо. Replay-cache/sequence state мають мати визначені persistence та expiry.
6. Для односпрямованої доставки перевірити, які гарантії можливі без challenge-response, що забезпечує фізичний діод, а що — прикладний протокол. Наявність DMA-буфера не називати фізичною неможливістю витоку.
7. Навантаження й чергу спершу моделювати синтетичними розмірами/сервісним часом; 128 КБ і 2500/с — заяви книги для перевірки, не погоджені тут норми.
8. Негативні дані ніколи не допускаються до production-KB. А-10 зберігає REFUSAL+QUARANTINE+audit; нові знання затверджуються в baseline, а не через per-action інтерактивний підпис.

## 3. Targets, PASS/FAIL, публікація

**Логічний target:** прийнята квитанція має незалежно перевірений proof саме підписаного statement у допустимій версії/часовій області; повтор не створює повторного ефекту.

**Ще не погоджено:** clock-skew bound, freshness policy, storage/replay horizon, максимальний proof/час перевірки, throughput. Книжкові 5000 мс і 1,852 мс не є доказом досягнутого SLA.

**PASS:** явна theorem/trust boundary, усі controls, machine-readable відмови, перевірні traces та proof artifacts. **FAIL:** accepted wrong-statement proof, майбутня/повторна транзакція з небажаним ефектом, або timeout як acceptance.

Негативний результат, обмеження приватності й недоведений фізичний діод публікуються. Не підміняти перевірку теореми позитивним тестом підпису.

## 4. Артефакти, команди й оновлення

Передати schema/version, synthetic signed fixtures, checker/proof/statement mapping, manifest SHA-256, abstract model, часові traces, persistence controls, команди й повний звіт.

Майбутній пакет має реалізувати; команди ще не виконувалися:

```text
python3 reproduce.py --manifest manifest.json --suite receipt-contract
python3 verify_receipts.py --manifest manifest.json --check statement-binding --check freshness-replay
```

Оновити CH33 §6.3 кортеж/предикат/визначення/часовий worked-приклад та межі hardware claims. Узгодити `REQ-L1-33-12`, `REQ-L2-33-15`, `REQ-L3-33-16`; окремі L3 для binding/replay/unknown погодити з автором. Підпис джерела не є детермінованим доказом ASIL за А-04.
