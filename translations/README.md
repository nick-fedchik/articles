# Багатомовна архітектура монографії «Архітектура доказових експертних систем»
# Multilingual Architecture: "Architecture of Evidence-Governed Expert Systems"

> **Автор (Author):** Микола Федчик (Mykola Fedchyk)  
> **Оригінал (Original):** Українська мова (`ExpertSystem/`)  
> **Рік:** 2026

---

## 1. Концепція та інженерні принципи перекладу

Ця монографія присвячена місійно-критичним системам, функціональній безпеці (ISO 26262 ASIL D, DO-178C DAL A), формальній верифікації та доказовому штучному інтелекту. Переклад такого тексту вимагає суворої інженерної точності:

1. **Збереження авторського голосу (Authorial Practitioner Voice):**
   - Переклад зберігає авторитетний, академічно зважений і водночас практичний голос Миколи Федчика.
   - Співвідношення: 80% об'єктивного системного викладу (безособові інженерні конструкції: *«it is proven that»*, *«the algorithm guarantees»*, *«wykazano, że»*, *«es wird bewiesen»*) та 20% практичного авторського голосу при описі стендів, бенчмарків і застережень (*«The author has deployed a research test bench...»*, *«Z doświadczenia inżynierskiego autora wynika...»*).
   - Сувора відмова від штучного колективного «ми» (це праця одного автора).
2. **Незмінність артефактів (Invariant Preservation):**
   - **Математичні вирази:** усі формули у `$$...$$`, ````math```` та `` $`...`$ `` переносяться строго байт-в-байт без зміни символів, дужок чи індексів.
   - **Вихідний код:** програми на Go, Rust, C, Python та специфікації Verilog залишаються незмінними.
   - **Діаграми Mermaid:** структура діаграм (`flowchart TD`, `sequenceDiagram`) та ідентифікатори вузлів не змінюються; перекладаються лише текстові мітки у лапках.
   - **Якірні посилання (Cross-chapter Anchors):** збереження відносної структури каталогів гарантує працездатність усіх перехресних посилань.
3. **Сувора заборона російської мови:**
   - Категорично виключено російськомовний контент у будь-яких проявах.

---

## 2. Структура каталогів та мовний пул (30+ мов світу)

| Код (ISO 639-1) | Мова / Регіон | Статус Маніфесту | Цільовий домен / Технологічний фокус |
| :---: | :--- | :---: | :--- |
| **uk** | **Українська** | **Master (Оригінал)** | Фронтир системної інженерії, кібернетики та оборонного MilTech (кластер Brave1) |
| **en** | **English** | **Global Reference** | Міжнародний академічний та інженерний стандарт (IEEE, ACM, NATO, Silicon Valley) |
| **de** | **Deutsch** | **Verified** | DACH: ISO 26262, AUTOSAR, ASPICE (Infineon, Bosch, Siemens, BMW, Mercedes) |
| **fr** | **Français** | **Verified** | Франкофонія: Аерокосмічна індустрія, критика безпеки, defense, AI ethics |
| **it** | **Italiano** | **Verified** | Промислова автоматизація, робототехніка, прецизійні системи |
| **nl** | **Nederlands** | **Verified** | Напівпровідники, літографія (ASML), голландський High-Tech кластер |
| **es** | **Español** | **Verified** | Глобальний технологічний простір Іспанії та Латинської Америки |
| **pt** | **Português** | **Verified** | Аерокосмічний кластер (Embraer), інженерний хаб Португалії та Бразилії |
| **pl** | **Polski** | **Verified** | Східний фланг НАТО, оборонні комплекси (PGZ, WB Group), хаб автомотови |
| **cs** | **Čeština** | **Verified** | Машинобудування, оборонна промисловість (CSG, Aero Vodochody) |
| **el** | **Ελληνικά** | **Verified** | Колиска логіки, світовий морський флот, оборонка НАТО, FORTH |
| **sr** | **Srpski** | **Verified** | Балканський інженерний та промисловий кластер |
| **ro** | **Română** | **Verified** | Автомобільні центри R&D, східний фланг НАТО, IT-хаб |
| **hu** | **Magyar** | **Verified** | Центральноєвропейський автомобільний кластер та розробка ПЗ |
| **bg** | **Български** | **Verified** | Балканський сектор оборонних технологій та аерокосмічні розробки |
| **hr** | **Hrvatski** | **Verified** | Адріатичний хаб автономних систем, електромобілі високого класу |
| **sv** | **Svenska** | **Verified** | Скандинавська інженерія, телеком, авіакосмічні комплекси (Saab) |
| **fi** | **Suomi** | **Verified** | Високонадійні вбудовані системи, зв'язок 6G, кіберзахист |
| **no** | **Norsk** | **Verified** | Енергетичні комплекси, морські автономні системи (Kongsberg) |
| **da** | **Dansk** | **Verified** | Квантові обчислення, інженерні вбудовані платформи, відновлювана енергія |
| **et** | **Eesti** | **Verified** | GovTech, автономна робототехніка (Milrem Robotics), кібероборона НАТО |
| **lt** | **Lietuvių** | **Verified** | Лазерні технології, напівпровідники, безпека Балтійського регіону |
| **lv** | **Latviešu** | **Verified** | Радіоелектроніка, телекомунікації, безпілотні комплекси |
| **ka** | **ქართული** | **Verified** | Кавказький регіон, стратегічне партнерство, кібернетичні комплекси |
| **hy** | **Հայերեն** | **Verified** | Математична школа, апаратні розробки, мікроелектроніка, DeepTech |
| **he** | **עברית** | **Verified** | Кібербезпека, аерокосмічні системи, Silicon Wadi, високонадійні чіпи |
| **ar** | **العربية** | **Verified** | MENA: Розумні міста майбутнього, енергетика, передові AI-ініціативи |
| **sk** | **Slovenčina** | **Verified** | Центральноєвропейський автомобільний кластер, промислова робототехніка |
| **sl** | **Slovenščina** | **Verified** | Високотехнологічний альпійський хаб, передові матеріали та мехатроніка |
| **ga** | **Gaeilge** | **Verified** | Європейська штаб-квартира глобальних технологічних гігантів, суверенний захист мови |
| **mt** | **Malti** | **Verified** | Середземноморський хаб морської логістики, фінтех та кібербезпека |
| **lb** | **Lëtzebuergesch** | **Verified** | Європейський фінансовий хаб, супутниковий зв'язок (SES) та суперкомп'ютери |
| **ca** | **Català** | **Verified** | Барселонський суперкомп'ютерний центр (BSC), мікроелектроніка та біотех |
| **ja** | **日本語** | **Verified** | Робототехніка, функціональна безпека, мікроконтролери та вбудовані системи |
| **zh / zh-TW** | **繁體中文** | **Verified** | Напівпровідникова столиця світу (TSMC), апаратні архітектури, електроніка |
| **ko** | **한국어** | **Verified** | Пам'ять, напівпровідники, автономний транспорт, авангард оборонної індустрії |
| **tr** | **Türkçe** | **Verified** | Безпілотні авіаційні комплекси (Baykar, TAI), оборонний MilTech |
| **hi** | **हिन्दी** | **Verified** | Індійський технологічний субрегіон, космічні програми (ISRO), Software R&D |
| **vi** | **Tiếng Việt** | **Verified** | Південно-Східна Азія: Високотехнологічне виробництво та електроніка |
| **id** | **Bahasa Indonesia** | **Verified** | Цифрова економіка АСЕАН, розумна інфраструктура |

---

## 3. Онтологічні глосарії (Glossaries & Termbases)

Усі переклади суворо спираються на стандартизовані глосарії в каталозі `translations/glossaries/`:
- `termbase-en.json` (English)
- `termbase-de.json` (German)
- `termbase-pl.json` (Polish)
- `termbase-fr.json` (French)
- `termbase-es.json` (Spanish)
- `termbase-it.json` (Italian)
- `termbase-nl.json` (Dutch)
- `termbase-pt.json` (Portuguese)
- `termbase-cs.json` (Czech)
- `termbase-el.json` (Greek)
- `termbase-sr.json` (Serbian)
- `termbase-ro.json` (Romanian)
- `termbase-hu.json` (Hungarian)
- `termbase-bg.json` (Bulgarian)
- `termbase-hr.json` (Croatian)
- `termbase-sk.json` (Slovak)
- `termbase-sl.json` (Slovenian)
- `termbase-sv.json` (Swedish)
- `termbase-fi.json` (Finnish)
- `termbase-no.json` (Norwegian)
- `termbase-da.json` (Danish)
- `termbase-et.json` (Estonian)
- `termbase-lt.json` (Lithuanian)
- `termbase-lv.json` (Latvian)
- `termbase-ga.json` (Irish)
- `termbase-mt.json` (Maltese)
- `termbase-lb.json` (Luxembourgish)
- `termbase-ca.json` (Catalan)
- `termbase-ka.json` (Georgian)
- `termbase-hy.json` (Armenian)
- `termbase-he.json` (Hebrew)
- `termbase-ar.json` (Arabic)
- `termbase-ja.json` (Japanese)
- `termbase-zh.json` (Traditional Chinese)
- `termbase-ko.json` (Korean)
- `termbase-tr.json` (Turkish)
- `termbase-hi.json` (Hindi)
- `termbase-vi.json` (Vietnamese)
- `termbase-id.json` (Indonesian)

---

## 4. Контроль якості (Quality Gate)

Кожен перекладений файл проходить обов'язкову верифікацію:
- `python3 scripts/verify_translation.py ExpertSystem/README.md <path_to_translation>` — автоматична перевірка тотожності AST (математичні блоки, діаграми Mermaid, Go-код, ієрархія заголовків H1–H4).
- `python3 scripts/inspect_buzzwords.py` — інспекція на неприпустимі машинні кліше (зокрема заборона вживання слова «контур» у нефізичних значеннях для українських текстів).
