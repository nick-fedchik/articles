# RESEARCH-CH33-REPLICATION-CONTRACT. Відновлення й атомарне застосування оновлень

**ID:** RESEARCH-CH33-REPLICATION-CONTRACT  
**Статус:** Відкрито; контракт зберігання/відновлення та профіль часу не затверджено.  
**Власник:** розробник VIDAdev обміну знаннями; співвиконавець — дослідник VIDA розподілених систем.  
**Пріоритет:** P1 — необмежена обіцянка доставки знань.  
**Підстава:** [W33-35-03](../../audit-2026-10-28-ch33-ch35.md#w33-35-03), [CH33 §4.2–§4.3](../../../ch33-inter-system-knowledge-exchange-and-model-teaching.md), 251–328.

**Вимоги:** [L1-33](../../../requirements/L1-stakeholder/ch33-L1-stakeholder.md): `REQ-L1-33-08`; [L2-33](../../../requirements/L2-system/ch33-L2-system.md): `REQ-L2-33-09`; [L3-33](../../../requirements/L3-component/ch33-L3-component.md): `REQ-L3-33-09`, `REQ-L3-33-10`.

## 1. Мета, межі й залежності

Перевірити умовну гарантію: вузол застосовує цілісну допустиму послідовність оновлень або переходить у явно визначений стан gap/snapshot recovery, не оголошуючи застарілий пакет актуальним.

Потрібні погоджені retention/storage, sequence domain, durable cursor, ack/retry, deduplication key, snapshot format і політика актуальності. Новий `T_safe` не вигадувати: 20 мс у draft-вимогах має проблему походження W10 й не є затвердженим тут target.

Схему mount/OpSet взяти з наявного завантажувача. Зв'язування квитанцій — у [RESEARCH-CH33-RECEIPT-CONTRACT](RESEARCH-CH33-RECEIPT-CONTRACT.md); abstraction квитанції дозволяє незалежно почати модель доставки. Локальний release після revocation — у [CH35-REVOCATION-EPOCH](RESEARCH-CH35-REVOCATION-EPOCH.md), а не другий proof цього самого gate.

## 2. Безпечний символічний протокол

1. Зафіксувати наявну реалізацію й конфігурацію брокера; окремо описати stream order і application order.
2. Побудувати скінченний автомат producer/broker/consumer/loader з поколінням, cursor, pending batch, mounted generation та validity status. Перебір запускати локально над синтетичними записами, без production-брокера.
3. Базові лінії: ідеальний надійний FIFO; at-least-once з наївним негайним застосуванням; запропонований versioned recovery protocol.
4. Контролі: втрата зв'язку коротша й довша retention; видалена дельта; дубль; переставлене доставлення; crash до/після ack; replay після перезапуску; одночасні revisions; помилка підпису/OpSet; неповний snapshot; зміна checksum.
5. Перевіряти ідемпотентність застосування, відсутність частково змонтованого покоління, виявлення пропуску, неможливість rollback на відкликану версію та повноту переходу до snapshot. Не вимагати безумовного liveness за нескінченного partition.
6. Для local CRDT назвати join і довести потрібні властивості; converge реплік не називати автоматично global safety, особливо для конкурентних відкликань.
7. Стан за прострочення lease/пакета визначати затвердженою політикою: safe profile або stop. Модель не має інструментів приводу; замість дії записує символічний вердикт.
8. Лише після проходження автомата відтворити ті самі контролі в ізольованому тестовому брокері з вигаданими subjects. Реальні строки/обсяги затверджуються до інтеграційного прогону.

## 3. Цілі та приймання

**Логічні targets:** жодний дубль не повторює зміну; відсутня дельта не маскується як complete; publication атомарна; invalid/expired пакет не отримує актуальний статус. Liveness заявляється лише за явно названих припущень fair delivery/retention.

**Очікує погодження:** максимальний offline interval, storage, retry/timeout, ресурсний профіль, допустима recovery latency. Приклади книги не стають новими нормативами.

**PASS:** інваріанти доведені в моделі й відтворені незалежними інтеграційними controls; межі й unsupported випадки видимі.

**FAIL:** silent gap, подвійне застосування, частковий mount, rollback revocation або необмежена гарантія «будь-який вузол завжди». Опублікувати всі негативні traces; зміну конфігурації після результату вважати новим експериментом.

## 4. Артефакти й точні зміни

Передати state machine, model-checking bounds, configuration manifest/SHA-256, seed/schedules, fake corpus, baseline traces, integration scripts, повний звіт failures і implemented/target.

Необхідний майбутній інтерфейс, а не вже наявні команди:

```text
python3 reproduce.py --manifest manifest.json --suite replication
python3 verify_traces.py --manifest manifest.json --property atomic-ordered-application
```

Оновити CH33 §4.2 умовами доставки, §4.3 sequenceDiagram gap/snapshot/atomic mount, §10 обмеженнями local safety і висновок 4. Синхронізувати `REQ-L1-33-08`, `REQ-L2-33-09`, `REQ-L3-33-09/10`; додаткові L3 на recovery/дедуплікацію — після погодження власником, не автоматично.
