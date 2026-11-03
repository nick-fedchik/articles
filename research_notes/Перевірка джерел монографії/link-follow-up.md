# Доповнення: збої бібліографічних і текстових URL

Початкове доповнення перевірено 2026-11-02 як продовження `bibliography-url-status.md`. Після наступних бібліографічних виправлень поточні байти глав повторно звірено 2026-11-02 о 22:12. Наведені нижче 12 первісних 404-діагнозів збережено як історичні, а не як перелік поточних URL: сім старих цілей замінено у главах, один Brier 404 був помилкою парсера, чотири 404 залишаються поточними. Поточний розподіл 777 цілей записано в `bibliography-url-status.md`. Повні RFC, стандарти, книги й статті не завантажувалися.

## Поточні байти, посилання та сканування тексту

У знімку 2026-11-02 було **1 005 входжень, 777 різних URL і мережевих цілей, 185 hostname**. Сім бібліографічних URL, які раніше дали 404, вилучено з тодішніх поточних глав і замінено. Після наступної цільової перевірки й виправлення Plaza DOI у ch34 залишаються три непідтверджені посилання: Bell–LaPadula MTR-2547 (ch02), ExploitGym (ch30), Rushby PDF (ch30). Старі URL з таблиці нижче не рахуються як поточні цілі.

Повторно витягнуто URL із тіла всіх 40 глав, виключивши кінцеві бібліографії: **33 входження, 25 різних URL; 24 зовнішні адреси**, бо `http://localhost:11434/api/chat` не надсилається назовні. Склад той самий, що в попередньому body-link аудиті: 20 зовнішніх цілей мали 2xx, а чотири демонстраційні `example.org` адреси — 404. Цього разу повторно перевірено наявність/кількість посилань у поточних байтах, але HTTP-статуси 24 body-цілей не повторювалися.

Повторне мережеве опитування 16 попередніх network errors у цьому оновленні не виконувалося: їхній попередній triage лишається окремо зафіксованим нижче — 15 OpenLibrary `ConnectionRefusedError`, один TLS/transport failure для ILASP. Ці помилки не є HTTP 404.

| Поточний файл | SHA-256 | `wc -l` | Поточні бібліографічні URL |
|---|---|---:|---:|
| `ExpertSystem/ch01-introduction-to-expert-systems.md` | `2f9a042071a9d9265963b7df89ac6f0463bc66adcf24b2e173ad98bc610999c5` | 637 | 23; новий occurrence — ISO 26262-1, URL `68383.html` |
| `ExpertSystem/ch05-triad-of-trust-and-corporate-memory.md` | `2e4b45e2f00d2cb3db88f218cd38c5bf6e0edc65a82629bc248d99b88199cc48` | 924 | 18 |
| `ExpertSystem/ch06-applied-mathematics-for-expert-systems.md` | `220186d2bd17df4dc4fe01298937dac1374deb70fcdddc679735085e2bc6585f` | 2474 | 59 |
| `ExpertSystem/ch18-execution-infrastructure.md` | `367d851cfcfa69e6f183802b960348c7f42435bc0dc2fa09003b2251c2200384` | 1961 | 49 |
| `ExpertSystem/ch19-from-question-to-evidence.md` | `b749d5dcf7281fbc977198f6808032ad37d4995b8d2981982f6d1c05e148ffc1` | 1236 | 12 |
| `ExpertSystem/ch26-continual-learning.md` | `9699466e3a0115c56ea76fba3bc4d1a80996bd10fe7f7e625c1bbd2e9f75afb9` | 1251 | 26 |
| `ExpertSystem/ch28-dual-mode-expert-systems.md` | `a0d666a5f82aaf47ef05169ea650c3eb1dec0f4b8ceffccfc508fc8cee2d32a5` | 1199 | 13 |
| `ExpertSystem/ch32-high-performance-knowledge-packs-mmap-and-harvesting.md` | `34a2fbad4c6bf117e7d4dcf01bb5ca89ddbf84a40dc82a38f624739e08f3c23c` | 3081 | 39 |
| `ExpertSystem/ch33-inter-system-knowledge-exchange-and-model-teaching.md` | `63066e8a2fb5e9e1fe0a394fc944f0bdb4668bf15d9e30648f8220bf461ea8e3` | 1881 | 27 |
| `ExpertSystem/ch38-curing-machine-hallucinations-and-knowledge-deficits.md` | `ced1adac4240ca268d4488d762714b4f0eb06a52f3e824ca16d423e7b6acb05e` | 1162 | 16 |

Точково перевірені поточні заміни (HEAD, 2026-11-02): Crossref API (ch33) `200`; Bibliopolis (ch32) `200`; NVIDIA R36.2 (ch18) `200`; IJCAI (ch05) `200`; PMLR REALM (ch19/ch38) `200`; DOI Sha (ch28) `202` після двох редиректів; DOI Hansson (ch26) `200` після трьох редиректів. HTTP-доступ підтверджує лише доступність цілей, не повний текст чи всі твердження глав. Поточний ISO 26262-1 каталог у ch01, навпаки, залишається `HEAD/GET 403`; статус джерельного тексту окремий.

## Історичні 12 URL, позначених у первісному скані як 404

Один із 12 записів виявився помилкою парсера. У главах 24 і 25 DOI Brier містить кутові дужки в середині ідентифікатора; попередній URL-сканер обрізав його на першому `<`, перевіривши неіснуючу коротку адресу. Решта 11 старих URL підтвердили 404.

| Цитата й точний шлях у книзі | Повторний статус | Перевірка ідентичності / альтернативний URL | Межа або рекомендація |
|---|---|---|---|
| Per Martin-Löf, *Intuitionistic Type Theory*. `ExpertSystem/ch32-high-performance-knowledge-packs-mmap-and-harvesting.md:3075`, джерело [35]. | Історичний URL `https://archive.org/details/martin-lof-1984`: `404` (`HEAD`, bounded `GET`). Пошук Archive.org за точними назвою й автором дав 0 записів. | Поточний ch32:3075 використовує Bibliopolis URL, який зараз повернув `HEAD 200`; сторінка вказує ISBN 9788870881059, 1984 рік та 92 сторінки, але відображає ім'я Peer Martin-Löf. | Стару адресу не рахувати як поточну ціль. Нова сторінка підтверджує бібліографічні метадані, не повний текст; книга позначена як відсутня в продажу. |
| David E. Bell, Leonard J. LaPadula, *Secure Computer System: Unified Exposition and Multics Interpretation*, ESD-TR-75-306. `ExpertSystem/ch33-inter-system-knowledge-exchange-and-model-teaching.md:1860`, [10]. | Історична сторінка NIST повернула `404`. | Поточний ch33:1860 використовує Crossref API URL, `HEAD 200`; метадані підтверджують назву, авторів, березень 1976 року та зв’язок із TR-75-306. DOI веде до DTIC, але ціль не відкрилася (transport/TLS failure); прямий DTIC URL повернув `403`. | Заміна перевірена тільки як метадані, не як текст звіту; старий NIST URL більше не є бібліографічною ціллю. |
| NVIDIA, *Jetson Linux Developer Guide: Power Management and nvpmodel*, JetPack 6.2 / R36.2. `ExpertSystem/ch18-execution-infrastructure.md:1952`, [43]. | Історична ціль із `/DeveloperGuide/text/SD/...` повернула `404`. | Поточний URL без сегмента `/text/` повернув `HEAD 200`; це сторінка тієї самої редакції й лінійки Orin. | Стару адресу більше не рахувати як поточну ціль. Перевірена заміна підтверджує сторінку документа, але не конкретні виміри чи WCET. |
| Jan Plaza, *Logics of public communications*, *Synthese* 158(2), 165–179 (2007). `ExpertSystem/ch34-deterministic-relational-analysis-abduction-and-socratic-dialogue.md:1748`, [11]. | Старий DOI `10.1007/s11225-007-9067-1` повернув `404`; Crossref exact lookup не знайшов його. | Crossref запис за DOI <https://doi.org/10.1007/s11229-007-9168-7> підтверджує автора, журнал, том, випуск, сторінки й рік; зареєстрована назва — **Logics of public communications**. | Виправлено DOI й назву в бібліографії глави 34. Перевірено метадані Crossref, не повний текст статті. |
| Edward A. Feigenbaum, *The Art of Artificial Intelligence: Themes and Case Studies of Knowledge Engineering*, IJCAI-77. `ExpertSystem/ch05-triad-of-trust-and-corporate-memory.md:917`, [15]. | Історичний DOI `10.1016/B978-0-934613-03-3.50008-3`: `404`; відповідний ScienceDirect PII URL відповів `403`. | Поточний ch05:917 URL офіційного архіву IJCAI повернув `HEAD 200` і підтверджує назву, автора та сторінки **1014–1028**. | Заміна підтверджує paper; у книзі раніше зазначено сторінки 1014–1029, а опис також об'єднував окремий HPP-84-1 (1984). Розділити ці записи, якщо потрібно точне бібліографічне оформлення. |
| Lui Sha, *Using Simplicity to Control Complexity*, *IEEE Software* 18(4), 20–28 (2001). `ExpertSystem/ch28-dual-mode-expert-systems.md:1195`, [13]. | Первісний DOI `10.1109/52.936213` повернув `HEAD/GET 404`; це історичний URL, уже замінений у главі. | Поточний DOI <https://doi.org/10.1109/ms.2001.936213> повернув `HEAD 202` після редиректів до IEEE Xplore; Crossref metadata збігаються за автором, назвою, журналом, випуском і сторінками. | Нова ціль є поточною. Landing/metadata доступні, але повний текст статті не читався. |
| Glenn W. Brier, *Verification of Forecasts Expressed in Terms of Probability*, *Monthly Weather Review* 78(1), 1–3 (1950). `ExpertSystem/ch24-system-diagnosis.md:1355`, [6]; `ch25-how-expert-systems-learn.md:1617`, [18]. | Попередній 404 стосувався обрізаного `https://doi.org/10.1175/1520-0493(1950)078`, а не повного DOI в главах. | Повний DOI `10.1175/1520-0493(1950)078<0001:VOFEIT>2.0.CO;2` є в Crossref із точними назвою, автором, журналом і сторінками; percent-encoded URL <https://doi.org/10.1175/1520-0493%281950%29078%3C0001:VOFEIT%3E2.0.CO;2> перенаправив до American Meteorological Society (`202`). | **Не є несправним джерелом.** Помилка була в попередньому URL-сканері: він відкинув частину DOI на символі `<`. Для безпечного Markdown-посилання можна використовувати percent-encoded варіант; глави не редагувалися. |
| David E. Bell, Leonard J. LaPadula, *Secure Computer Systems: Mathematical Foundations*, MITRE MTR-2547, Vol. I (1973). `ExpertSystem/ch02-epistemology-of-machine-knowledge.md:1527`, [34]. | DOI `10.21236/AD0770768`: повторно `404`; Crossref exact lookup і варіант `ADA0770768` повернули `404`. Перевірені варіанти DTIC для цього ідентифікатора відповіли `403`. | **Підтверджену альтернативу не знайдено.** | Відкрита метаінформація не встановила, що запропонований DTIC-кандидат ідентифікує саме MTR-2547, Vol. I. Не замінювати URL за схожістю назви; точна прив’язка лишається неперевіреною. |
| Sven Ove Hansson, *Kernel Contraction*, *Journal of Symbolic Logic* 59(3), 845–859 (1994). `ExpertSystem/ch26-continual-learning.md:1247`, [26]. | Первісний DOI `10.2307/2275662` повернув `HEAD/GET 404`; це історичний URL, уже замінений у главі. | Поточний DOI <https://doi.org/10.2307/2275912> повернув `HEAD 200` після трьох редиректів до Cambridge; Crossref metadata збігаються за назвою, автором, журналом, томом, випуском і сторінками. | Нова ціль є поточною; доступний landing/abstract, повний текст видавця не читався. |
| ExploitGym Team, *ExploitGym: Benchmarking Frontier AI Agents on Offensive Cyber Vulnerability Discovery and Sandbox Escape* (у книзі датовано 2025). `ExpertSystem/ch30-safety-cybersecurity-co-engineering.md:1378`, [17]. | `https://github.com/exploitgym`: повторно `404`. | GitHub metadata search знайшов <https://github.com/sunblaze-ucb/exploitgym> (назва ExploitGym, опис benchmark; створено 2026-06-02). arXiv API знайшов запис 2605.11086v1 від 2026-05-11 із **іншою назвою** (*ExploitGym: Can AI Agents Turn Security Vulnerabilities into Real Attacks?*). | **Заміни для цитати 2025 не підтверджено.** Репозиторій 2026 року й арXiv-запис іншої назви не доводять тотожності з бібліографічним записом. Перевірено лише метадані репозиторію та статті; код, інструкції чи експлуатаційні матеріали не відкривалися. |
| Kelvin Guu та ін., *REALM: Retrieval-Augmented Language Model Pre-Training*, ICML 2020. `ExpertSystem/ch19-from-question-to-evidence.md:1232`, [12]; `ch38-curing-machine-hallucinations-and-knowledge-deficits.md:1150`, [13]. | Історичний Google Research URL повернув `404`. | Поточний PMLR record повернув `HEAD 200` і підтверджує точні назву, авторів та ICML 2020 metadata. | Заміна тепер є поточною бібліографічною ціллю; вона підтверджує identity/abstract, не всі прикладні висновки глав. |
| John Rushby та ін., *Understanding Assurance 2.0*, SRI, 2020. `ExpertSystem/ch30-safety-cybersecurity-co-engineering.md:1381`, [20]. | PDF URL повторно `404`. | Офіційна бібліографія автора <https://www.csl.sri.com/~rushby/biblio.html> і сторінка <https://www.csl.sri.com/users/rushby/assurance2.0> доступні (`200`); вони описують тему Assurance 2.0 та інші статті, але не показали цитований документ із цими авторами/назвою. | **Заміни для цієї конкретної праці не підтверджено.** Сторінка Assurance 2.0 є лише пов’язаним оглядом; не використовувати її як тотожну статтю 2020 року. |

## Повторна перевірка 16 мережевих помилок

Усі 16 URL перевірено ще раз через `HEAD`. 15 сторінок OpenLibrary знову дали `ConnectionRefusedError: [Errno 111]`; окрема спроба `web_fetch` на один із них також завершилася transport-level failure. Це **мережева недоступність**, не HTTP 404. Для цих 15 книг альтернативні каталожні/видавничі URL не верифікувалися, тому не пропонуються.

| URL OpenLibrary | Книга й точний chapter/source | Повтор |
|---|---|---|
| `https://openlibrary.org/books/OL5615707M` | `ch04-evolution-from-bayes-to-evidence-ai.md` [31], M. Ross Quillian, *Semantic Memory*. | Connection refused |
| `https://openlibrary.org/works/OL117061W` | `ch11-knowledge-elicitation-from-experts.md` [1], Michael Polanyi, *The Tacit Dimension*. | Connection refused |
| `https://openlibrary.org/works/OL1311228W` | `ch06-applied-mathematics-for-expert-systems.md` [12], Nils J. Nilsson, *Problem-Solving Methods in Artificial Intelligence*. | Connection refused |
| `https://openlibrary.org/works/OL16585333W` | `ch09-engineering-knowledge-graph-traceability.md` [12], Antoniou, Groth, van Harmelen, Hoekstra, *A Semantic Web Primer*. | Connection refused |
| `https://openlibrary.org/works/OL17872278W` | `ch06-applied-mathematics-for-expert-systems.md` [45], Pearl і Mackenzie, *The Book of Why*. | Connection refused |
| `https://openlibrary.org/works/OL1814674W` | `ch03-beyond-reference-information-systems.md` [14], Bertrand Meyer, *Object-Oriented Software Construction*, 2nd ed. | Connection refused |
| `https://openlibrary.org/works/OL21407442W` | `ch15-knowledge-extraction-and-kb-construction.md` [14], Hopcroft, Motwani, Ullman, *Introduction to Automata Theory, Languages, and Computation*. | Connection refused |
| `https://openlibrary.org/works/OL24825917W` | `ch04-evolution-from-bayes-to-evidence-ai.md` [1], Alty і Coombs, *Expert Systems: Concepts and Examples*. | Connection refused |
| `https://openlibrary.org/works/OL3307771W` | `ch15-knowledge-extraction-and-kb-construction.md` [26], W. Ross Ashby, *An Introduction to Cybernetics*. | Connection refused |
| `https://openlibrary.org/works/OL3515903W` | `ch11-knowledge-elicitation-from-experts.md` [22], Nonaka і Takeuchi, *The Knowledge-Creating Company*. | Connection refused |
| `https://openlibrary.org/works/OL4305747W` | `ch11-knowledge-elicitation-from-experts.md` [7], Ericsson і Simon, *Protocol Analysis: Verbal Reports as Data*. | Connection refused |
| `https://openlibrary.org/works/OL5282413W` | `ch11-knowledge-elicitation-from-experts.md` [14], Krippendorff, *Content Analysis: An Introduction to Its Methodology*. | Connection refused |
| `https://openlibrary.org/works/OL5959470W` | `ch02` [9], `ch03` [16], `ch06` [7], `ch31` [5], Kleene, *Introduction to Metamathematics*. | Connection refused |
| `https://openlibrary.org/works/OL6408872W` | `ch02-epistemology-of-machine-knowledge.md` [17], J. L. Austin, *How to Do Things with Words*. | Connection refused |
| `https://openlibrary.org/works/OL9203874W` | `ch09-engineering-knowledge-graph-traceability.md` [3], John F. Sowa (ed.), *Principles of Semantic Networks*. | Connection refused |

| URL | Chapter/citation | Повтор і альтернатива | Межа |
|---|---|---|---|
| <https://ilasp.com/> | `ch23-knowledge-base-verification.md` [24], ILASP, *Logic-Based Machine Learning*. | Повторна `HEAD` і `web_fetch` — TLS certificate/transport failure. Перевірений офіційний посібник <https://doc.ilasp.com/> відповідає за призначенням і має `200`. | Посібник називається *The ILASP manual* і позначений як under construction; це альтернативна документація інструмента, не незалежний науковий результат. |

## Зовнішні URL у тілі глав, поза кінцевою бібліографією

У 40 главах знайдено 33 входження, 25 різних HTTP(S) URL. Один URL — loopback `http://localhost:11434/api/chat` (`ch29:646`), тому він не надсилався назовні. Решта 24 зовнішні мережеві цілі перевірені окремим HEAD-first проходом; у разі помилки виконано тільки bounded GET. Результат: 20 цілей 2xx (включно з редиректами), 4 цілі `example.org` — 404.

| Глава/рядок | URL | Статус |
|---|---|---|
| `ch30:547` | <http://www.omg.org/spec/ReqIF/20110401/reqif.xsd> | 200 |
| `ch33:314` | <http://www.w3.org/2001/XMLSchema#> | 200 після 301 до HTTPS |
| `ch07:578`, `ch15:1679,1690` | <http://www.w3.org/2000/01/rdf-schema#> | 200 після 301 до HTTPS |
| `ch15:1678` | <http://www.w3.org/2002/07/owl#> | 200 після 301 до HTTPS |
| `ch16:871,907` | <http://www.w3.org/ns/prov#> | 200 після 301 до HTTPS |
| `ch07:564` | <http://www.w3.org/ns/shacl#> | 200 після 2 редиректів до HTTPS |
| `ch07:175` | <https://conceptnet.io/> | 200 |
| `ch07:176` | <https://geneontology.org/docs/ontology-documentation/> | 200 |
| `ch08:242` | <https://github.com/jdkato/prose> | 200 |
| `ch08:242` | <https://github.com/knights-analytics/hugot> | 200 |
| `ch08:242` | <https://github.com/nlpodyssey/cybertron> | 200 |
| `ch08:242` | <https://github.com/urchade/GLiNER> | 200 |
| `ch08:242` | <https://github.com/vale-cli/vale> | 200 |
| `ch07:179` | <https://incidentdatabase.ai/> | 200 |
| `ch10:1100` | <https://nats.io/> | 200 |
| `ch07:177` | <https://openfisca.org/en/> | 200 |
| `ch07:178` | <https://www.bnlearn.com/bnrepository/> | 200 |
| `ch07:1136` | <https://www.postgresql.org/docs/18/functions-logical.html> | 200 |
| `ch19:344,357` | <https://www.rfc-editor.org/rfc/rfc768.txt> | 200 |
| `ch07:174` | <https://www.wikidata.org/wiki/Wikidata:Introduction> | 200 |
| `ch16:872,908` | <https://example.org/kb/> | 404 |
| `ch15:1677,1689` | <https://example.org/protocol/> | 404 |
| `ch07:116,563,577` | <https://example.org/rd/> | 404 |
| `ch33:313` | <https://example.org/teaching/envelope#> | 404 (fragment не передається) |

`example.org` URL-и — демонстраційні namespace-и в коді/прикладах, не джерельні вебсторінки; заміни не потрібні без рішення автора про використання реального namespace. HTTP → HTTPS редиректи W3C namespace не є підставою переписувати namespace IRIs: схема входить до ідентифікатора. `localhost` — локальний приклад, не зовнішня ціль.

## Межі

Повторні HTTP-коди є знімком доступності на момент перевірки. `403`, connection refusal і TLS-помилки не означають, що документ не існує; лише альтернативи, для яких наведено точне зіставлення метаданих і доступну ціль, рекомендовано як заміну.

## Цільове зіставлення 2026-11-03

- **Plaza (ch34 [11]):** Crossref запис для DOI `10.1007/s11229-007-9168-7` підтвердив Jan Plaza, *Logics of public communications*, *Synthese* 158(2), 165–179 (2007). Бібліографічний DOI й назву в главі виправлено. Перевірено метадані, не повний текст.
- **Bell–LaPadula (ch02 [34]):** повторний Crossref lookup DOI `10.21236/AD0770768` дав 404; пошук DTIC за `MTR-2547` не надав видимого збігу, а сторінка каталогу містила лише загальне повідомлення. Підтвердженої заміни немає.
- **ExploitGym (ch30 [17]):** GitHub API підтвердив однойменний репозиторій `sunblaze-ucb/exploitgym` від 2026-06-02; arXiv API підтвердив 2605.11086v1 (2026-05-11) із назвою *ExploitGym: Can AI Agents Turn Security Vulnerabilities into Real Attacks?*. Метадані пов'язують записи з benchmark-проєктом, але не підтверджують цитований у книзі запис 2025 року з іншою назвою. Код та експлуатаційні матеріали не відкривалися.
- **Rushby (ch30 [20]):** офіційна сторінка SRI підтверджує напрям Assurance 2.0 та перелік пов'язаних робіт, але не засвідчує звіт 2020 року з назвою *Understanding Assurance 2.0* й наведеним авторським складом. Сторінку не вважаємо тотожною працею.

У такий спосіб один із чотирьох бібліографічних випадків одержав метаданими підтверджену заміну; три лишаються непідтвердженими. Поточний повний інвентар URL і HTTP-коди після виправлення Plaza не перераховувалися.
