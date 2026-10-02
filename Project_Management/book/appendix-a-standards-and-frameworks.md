# Додаток А. Реєстр інженерних стандартів і методологій

> **Книга:** [Керування складними інженерними проєктами та програмами](README.md) · Додатки  
> **Попередній розділ:** [Глава 18. Відповідальність і нагляд: люди в циклі рішень, запобігання управлінським галюцинаціям та етика](ch18-governance-and-human-accountability.md)  
> **Наступний додаток:** [Додаток Б. Еталонна схема цифрової нитки інженерного проєкту](appendix-b-digital-thread-data-schema.md)  
> **Зміст книги:** [README.md](README.md)  
> **Автор:** [Микола Федчик](about-the-author.md)  
> **Рівень:** системні інженери, директори з якості, аудитори, керівники програм  
> **Очікувані результати:** орієнтуватися в ландшафті міжнародних інженерних і управлінських стандартів; зіставляти вимоги до процесів життєвого циклу між різними галузями; будувати інтегровану систему комплаєнсу підприємства.

---

## Призначення реєстру

Розробка складних високотехнологічних комплексів вимагає дотримання взаємопов'язаних нормативних документів міжнародних організацій (ISO, IEC, IEEE, SAE, RTCA, PMI). Кожен стандарт фокусується на своїй ділянці життєвого циклу: системній інженерії, функціональній безпеці, надійності мікропрограм, управлінні конфігураціями чи фінансовому контролі розкладу.

Цей додаток систематизує стандарти, розглянуті в книзі, та подає зведену матрицю їхніх взаємозв'язків.

---

## Зведена таблиця інженерних і управлінських стандартів

| Стандарт | Організація | Галузь застосування | Головний фокус вимог | Відповідні глави книги |
| :--- | :--- | :--- | :--- | :--- |
| **ISO/IEC/IEEE 15288** | ISO/IEC/IEEE | Загальна системна інженерія | Процеси життєвого циклу систем: від концепції до виведення з експлуатації | [Глава 3](ch03-digital-thread-and-engineering-data-model.md), [Глава 4](ch04-wbs-and-multi-cadence-delivery.md) |
| **ISO 26262** | ISO | Автомобільна промисловість | Функціональна безпека електричних та електронних систем (ASIL A-D) | [Глава 10](ch10-compliance-as-operational-state.md), [Глава 11](ch11-release-evidence-and-deterministic-gates.md) |
| **DO-178C / ED-12C** | RTCA / EUROCAE | Цивільна та військова авіоніка | Верифікація бортового програмного забезпечення (DAL A-E, покриття MC/DC) | [Глава 10](ch10-compliance-as-operational-state.md), [Глава 11](ch11-release-evidence-and-deterministic-gates.md) |
| **DO-254 / ED-80** | RTCA / EUROCAE | Авіаційна апаратура | Проєктування та верифікація електронного обладнання та програмованої логіки FPGA | [Глава 4](ch04-wbs-and-multi-cadence-delivery.md), [Глава 10](ch10-compliance-as-operational-state.md) |
| **IEC 62304** | IEC | Медичне приладобудування | Процеси життєвого циклу медичного софту (Class A, B, C) | [Глава 10](ch10-compliance-as-operational-state.md) |
| **ISO 14971** | ISO | Медичні вироби | Управління ризиками для медичних пристроїв | [Глава 10](ch10-compliance-as-operational-state.md), [Глава 14](ch14-programme-management-and-systemic-risk.md) |
| **ANSI/EIA-748-D** | SAE / EIA | Промисловість і оборона | Критерії систем управління освоєним обсягом (32 критерії EVMS) | [Глава 9](ch09-earned-value-and-probabilistic-forecasting.md) |
| **ANSI/EIA-649-C** | SAE | Системна інженерія | Стандарт управління конфігураціями та базовими лініями | [Глава 5](ch05-baselines-and-change-control-board.md) |
| **CMMI-DEV v2.0** | ISACA / CMMI | Інженерні підприємства | Модель зрілості спроможностей організації (Рівні 1-5) | [Глава 16](ch16-pmo-maturity-and-platform-strategy.md) |
| **PMI PMBOK Guide** | PMI | Загальний проєктний менеджмент | Стандарт управління проєктами, WBS, розклад, вартість | [Глава 1](ch01-engineering-management-in-crisis.md), [Глава 2](ch02-three-horizons-project-programme-portfolio.md) |
| **PMI Standard for Program Management** | PMI | Програмний менеджмент | Інтеграція підпроєктів, пулінг ризиків, реалізація вигод | [Глава 2](ch02-three-horizons-project-programme-portfolio.md), [Глава 14](ch14-programme-management-and-systemic-risk.md) |
| **PMI Standard for Portfolio Management** | PMI | Портфельне управління | Розподіл капіталу, стратегічне вирівнювання, ліміти ресурсів | [Глава 2](ch02-three-horizons-project-programme-portfolio.md), [Глава 15](ch15-portfolio-governance-and-capital-allocation.md) |
| **W3C PROV** | W3C | Інформаційні системи | Модель даних походження та історії трансформацій артефактів | [Глава 3](ch03-digital-thread-and-engineering-data-model.md), [Глава 12](ch12-audit-ready-by-design.md) |
| **OASIS OSLC** | OASIS | Інтеграція інструментів | Відкриті протоколи інтеграції життєвого циклу систем (ALM-PLM) | [Глава 3](ch03-digital-thread-and-engineering-data-model.md), [Глава 16](ch16-pmo-maturity-and-platform-strategy.md) |
| **GSN Standard v3** | SCSC | Безпека систем | Стандарт візуалізації та структурування аргументів безпеки | [Глава 11](ch11-release-evidence-and-deterministic-gates.md) |

---

## Крос-мапінг процесів життєвого циклу

Для побудови єдиної корпоративної системи управління важливо розуміти, як процеси різних стандартів відображаються один на один:

1. **Управління вимогами:** ISO 15288 (п. 6.4.2) $\leftrightarrow$ ISO 26262-8 (п. 6) $\leftrightarrow$ DO-178C (п. 5.1). Усі стандарти вимагають однозначної ідентифікації кожної вимоги та її двоспрямованої простежуваності до архітектури й тестів.
2. **Конфігураційне управління:** EIA-649 $\leftrightarrow$ ISO 15288 (п. 6.3.5) $\leftrightarrow$ DO-178C (п. 7.0). Вимагає наявності формальних базових ліній (Baselines) і функціонування ради змін (CCB).
3. **Верифікація та валідація:** ISO 15288 (п. 6.4.9, 6.4.11) $\leftrightarrow$ ISO 26262-4 $\leftrightarrow$ DO-178C (п. 6.0). Вимагає підтвердження виконання вимог об'єктивними доказами, отриманими на цільовому апаратному забезпеченні.

---

## Словник

- **Матриця відповідності (Compliance Matrix):** документ, що зіставляє вимоги міжнародного стандарту з внутрішніми регламентами, інструментами та артефактами підприємства.
- **Стандарт життєвого циклу:** нормативний документ, що визначає обов'язкові процеси, контрольні точки та вимоги до результатів розробки від початку проєкту до утилізації виробу.

---

## Абревіатури

- **ASIL** (*Automotive Safety Integrity Level*): рівень повноти безпеки в автомобілебудуванні.
- **CMMI** (*Capability Maturity Model Integration*): модель зрілості процесів розробки.
- **DAL** (*Development Assurance Level*): рівень гарантії розробки в авіації.
- **EVMS** (*Earned Value Management System*): система управління освоєним обсягом.
- **GSN** (*Goal Structuring Notation*): нотація структурування аргументів безпеки.
- **INCOSE** (*International Council on Systems Engineering*): міжнародна рада з системної інженерії.
- **OSLC** (*Open Services for Lifecycle Collaboration*): відкритий стандарт інтеграції систем життєвого циклу.
- **PMI** (*Project Management Institute*): міжнародний інститут проєктного менеджменту.
- **PROV** (*Provenance Model*): стандарт походження даних W3C.
- **SAE** (*Society of Automotive Engineers*): міжнародне товариство автомобільних інженерів.

---

## Джерела

1. <a id="src-1"></a>ISO/IEC/IEEE. [*ISO/IEC/IEEE 15288: Systems and software engineering - System life cycle processes*](https://www.iso.org/standard/63711.html). ISO/IEC/IEEE, 2015.
2. <a id="src-2"></a>International Organization for Standardization. [*ISO 26262: Road vehicles - Functional safety*](https://www.iso.org/standard/68383.html). ISO, 2018.
3. <a id="src-3"></a>RTCA / EUROCAE. [*DO-178C / ED-12C: Software Considerations in Airborne Systems and Equipment Certification*](https://www.rtca.org/). RTCA, 2011.
4. <a id="src-4"></a>International Electrotechnical Commission. [*IEC 62304: Medical device software - Software life cycle processes*](https://webstore.iec.ch/publication/6923). IEC, 2006/AMD1:2015.
5. <a id="src-5"></a>SAE International. [*EIA-748-D: Earned Value Management Systems*](https://www.sae.org/standards/content/eia748d/). SAE, 2019.
6. <a id="src-6"></a>CMMI Institute. [*CMMI for Development, Version 2.0*](https://cmmiinstitute.com/). ISACA, 2018.

---

[← Глава 18](ch18-governance-and-human-accountability.md) | [Зміст книги](README.md) | [Додаток Б →](appendix-b-digital-thread-data-schema.md)
