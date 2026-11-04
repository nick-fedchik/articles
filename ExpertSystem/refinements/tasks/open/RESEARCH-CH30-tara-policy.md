# RESEARCH-CH30-tara-policy. Підтвердження шкали TARA та політики допуску

**ID:** RESEARCH-CH30-tara-policy  
**Статус:** Відкрито; локальну числову неузгодженість доведено, нормативні таблиці й policy не підтверджено.  
**Власник:** дослідник safety/cybersecurity та розробник VIDA policy engine.  
**Пріоритет:** P0 — закриття цілі без повного risk argument.  
**Підстава:** [аудит W27-30-19, W27-30-20, W27-30-21](../../audit-2026-10-28-ch27-ch30.md), [CH30 §2.1.1–§2.1.3, §2.3](../../../ch30-safety-cybersecurity-co-engineering.md).  
**Вимоги:** [L1](../../../requirements/L1-stakeholder/ch30-L1-stakeholder.md): REQ-L1-30-05, REQ-L1-30-06; [L2](../../../requirements/L2-system/ch30-L2-system.md): REQ-L2-30-06, REQ-L2-30-07, REQ-L2-30-08, REQ-L2-30-09; [L3](../../../requirements/L3-component/ch30-L3-component.md): REQ-L3-30-05, REQ-L3-30-06, REQ-L3-30-07, REQ-L3-30-08, REQ-L3-30-09, REQ-L3-30-10.

## 1. Залежності й предмет

Для власної таблиці Standard=0, тому приклад має суму 8, не 9. Окрема складніша проблема — feasibility threshold як достатня альтернатива mitigation, без impact/залишкового ризику. Таймаут автентифікації не доводить spoofing.

Передумова нормативної звірки — законний доступ уповноваженого фахівця до конкретної редакції ISO/SAE 21434. Не публікувати ліцензовані таблиці повністю; надати пункт, редакцію й дозволений короткий опис. [CH18 risk calibration](RESEARCH-CH18-risk-model-calibration.md) лишається суміжним методом, не другою реалізацією TARA. Мікротаймінги криптографії — [TASK-TST-010](TASK-TST-010-performance-evidence.md).

## 2. Спростовуваний безпечний протокол

1. Звірити фактори/дискретні значення/категорії; відокремити нормативний метод від author policy. Погодити risk-acceptance owner, impact, горизонт переоцінки й залишкові ризики.
2. На синтетичних записах порівняти поточний feasibility-only rule, незалежне експертне рішення та кандидатну policy. Контролі меж 13/14/19/20/24/25, невідомий фактор, недопустимий проміжний бал, той самий score з різним impact.
3. Перевірити F=20 без verified mitigation і high-impact сценарій: automatic closure не дозволяти без погодженого повного аргументу.
4. Локальний mock-канал без транспортного обладнання: справжній кадр з пізньою перевіркою, invalid tag, replay, новий counter без автентифікації, коректні tag+freshness, unknown key. Розділити safety refusal та причину інциденту.
5. Встановити per-frame deadline із загального budget, arrival/burst/queue і recovery; не переносити одиничні 180 мкс на всю систему. Не здійснювати реальних атак/ін'єкцій у CAN або isoUART.

## 3. Цілі та PASS/FAIL

13/19/24 і F≥20 — **наявні книжкові межі, нормативна/політична застосовність яких очікує підтвердження**. Нового порога ризику, швидкодії або криптоалгоритму не вводити. А-08 signed baseline і А-04 deterministic ASIL proof зберегти.

**PASS:** trace кожного правила до законного джерела або затвердженої локальної policy; повний acceptance argument; коректні synthetic verdicts. **FAIL:** closure лише за feasibility без дозволеної policy, невідомий score проходить або timeout стає доведеним spoofing. Публікувати негативні рішення й неоднозначності нормативної звірки.

## 4. Артефакти, команди й оновлення

Передати reference matrix без забороненого відтворення стандарту, signed policy version, fixtures, expected/actual decisions, mock-channel trace, arrival/budget assumptions і точні команди test runner фактичної реалізації.

Безпечний початковий контроль внутрішньої таблиці книги, не нормативна валідація:

```bash
python3 -c 'equipment_scores={0,4,9}; assert 1 not in equipment_scores; corrected=1+3+3+1+0; assert corrected==8; print("corrected score",corrected,"category remains High")'
```

Оновити CH30 §2.1.1: дискретні фактори й посилання; §2.1.2: достатні умови/межі closure; §2.1.3: суму, auth/freshness, типи timeout/invalid/replay і budget; §2.3: узгоджену policy. Синхронізувати всі перелічені L1–L3. Не приписувати стандарту автоматичний наказ конкретного чипа або підписування кожної runtime-дії людиною.
