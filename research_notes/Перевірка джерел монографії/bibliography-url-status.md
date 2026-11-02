# Перевірка доступності бібліографічних URL у главах 1–40

Первісний URL-знімок: 2026-11-02T20:54:00+02:00. 2026-11-02T22:12:45+02:00 повторно прочитано поточні байти глав 1–40 після виправлень бібліографії. Сфера: лише кінцеві бібліографічні розділи `## Джерела` (або однойменні еквіваленти) у 40 файлах; цей файл оновлює дослідницький запис і не змінює глави чи реєстри.

> **Поправка щодо Brier DOI:** парсер першого проходу обрізав DOI у главах 24 і 25 на кутовій дужці, перевіривши неіснуючу коротку ціль. Повний DOI з глав підтверджено через Crossref; percent-encoded варіант повертає HTTP 202. Після подальших книжкових виправлень семеро раніше несправних URL більше не є бібліографічними цілями й замінені перевіреними URL. Поточний знімок має **629 цілей 2xx та 4 невирішені 404**. Початкові дванадцять діагнозів лишено як історичні в [доповненні про збої посилань](link-follow-up.md).

## Підсумок і метод

- У поточних байтах виявлено 1 005 входжень HTTP(S)-посилань у бібліографіях; 777 різних URL-рядків, 777 різних мережевих цілей після вилучення fragment-компонента та 185 hostname. Після виправлення посилань чотири поточні 404 — Plaza DOI (ch34), Bell–LaPadula MTR-2547 DOI (ch02), ExploitGym (ch30) і PDF Rushby (ch30); сім старих 404 замінено, а Brier 404 був помилкою парсера.
- URL виділялися з бібліографічних розділів і очищалися від Markdown-розділових дужок: наприкінці знімалася лише незбалансована закривальна дужка та кінцева пунктуація; збалансовані дужки всередині шляху збережено. Для Brier DOI кутові дужки передавалися в percent-encoded формі; fragment не передавався серверу.
- Для кожної цілі першою виконано `HEAD`; редиректи пройдено вручну, максимум 9 переходів. Лише за HEAD 400/403/404/405/501 виконано `GET` з `Range: bytes=0-1023`, прочитано не більше 1 024 байтів відповіді й закрито з’єднання. Повні сторінки, RFC і стандарти не завантажувалися.
- Обмеження частоти: не більш як один запит за 1,05 с на hostname; до 8 незалежних hostname паралельно; тайм-аут одного HTTP-запиту — 12 с. Неуспішні мережеві запити автоматично не повторювалися; результати є одноразовим знімком доступності на час перевірки.
- Аудит встановлює HTTP-доступність і кінцеву адресу, але не доводить, що вміст кожного URL відповідає цитованому твердженню. Ідентичність та текст критичних сторінок перевірялися окремо в `local-source-matches.md` і `chapters-15-21.md`.

### Розподіл кінцевих статусів

- blocked/auth/policy 4xx: 126
- network error: 16
- not found/gone 4xx: 4
- other client 4xx: 1
- reachable 2xx: 629
- server 5xx: 1

Статус HTTP-методу за фактичними запитами (включно з редиректами та обмеженим GET fallback):
- `GET` → `200`: 5
- `GET` → `206`: 8
- `GET` → `302`: 1
- `GET` → `400`: 1
- `GET` → `403`: 126
- `GET` → `404`: 4
- `HEAD` → `200`: 567
- `HEAD` → `202`: 49
- `HEAD` → `301`: 191
- `HEAD` → `302`: 406
- `HEAD` → `303`: 44
- `HEAD` → `307`: 2
- `HEAD` → `308`: 2
- `HEAD` → `400`: 1
- `HEAD` → `403`: 129
- `HEAD` → `404`: 12
- `HEAD` → `405`: 2
- `HEAD` → `502`: 1
- `HEAD` → `NETWORK_ERROR`: 16

- URL із хоча б одним редиректом або зміненою кінцевою адресою: 352.
- HEAD, для яких використано обмежений GET fallback: 144.

## Перелік перевірених мережевих цілей

| Бібліографічні глави | URL (фрагмент може бути опущений у запиті) | Результат | Перенаправлення / кінцева адреса |
|---|---|---|---|
| ch36 | <https://aclanthology.org/2020.acl-main.442/> | HEAD: HTTP 200 | — |
| ch08 | <https://aclanthology.org/2023.eacl-main.148/> | HEAD: HTTP 200 | — |
| ch12 | <https://aclanthology.org/2023.emnlp-main.398/> | HEAD: HTTP 200 | — |
| ch12 | <https://aclanthology.org/2023.emnlp-main.741/> | HEAD: HTTP 200 | — |
| ch13 | <https://aclanthology.org/2024.naacl-long.300/> | HEAD: HTTP 200 | — |
| ch10, ch31 | <https://aclanthology.org/C92-2082/> | HEAD: HTTP 200 | — |
| ch12 | <https://aclanthology.org/D18-1269/> | HEAD: HTTP 200 | — |
| ch06, ch08, ch12 | <https://aclanthology.org/D19-1410/> | HEAD: HTTP 200 | — |
| ch13 | <https://aclanthology.org/L16-1680/> | HEAD: HTTP 200 | — |
| ch23 | <https://aclanthology.org/P18-1128/> | HEAD: HTTP 200 | — |
| ch04, ch34 | <https://aima.cs.berkeley.edu/> | HEAD: HTTP 200 | — |
| ch23 | <https://apalache-mc.org/docs/> | HEAD: HTTP 200 | — |
| ch33 | <https://api.crossref.org/works/10.21236/ADA023588> | HEAD: HTTP 200 | — |
| ch12 | <https://apps.dtic.mil/sti/citations/ADA039324> | GET (HEAD 403 → ranged GET): HTTP 403 | — |
| ch22 | <https://archive.org/details/introductiontocy00ashb> | HEAD: HTTP 200 | — |
| ch10 | <https://arm-software.github.io/CMSIS_5/SVD/html/index.html> | HEAD: HTTP 200 | — |
| ch09 | <https://arm-software.github.io/CMSIS_6/latest/General/index.html> | HEAD: HTTP 200 | — |
| ch26 | <https://artowen.su.domains/mc/Ch-var-is.pdf> | HEAD: HTTP 200 | — |
| ch12 | <https://arxiv.org/abs/1003.4394> | HEAD: HTTP 200 | — |
| ch26 | <https://arxiv.org/abs/1103.4601> | HEAD: HTTP 200 | — |
| ch24 | <https://arxiv.org/abs/1302.3563> | HEAD: HTTP 200 | — |
| ch33 | <https://arxiv.org/abs/1503.02531> | HEAD: HTTP 200 | — |
| ch25 | <https://arxiv.org/abs/1604.06174> | HEAD: HTTP 200 | — |
| ch26 | <https://arxiv.org/abs/1606.06565> | HEAD: HTTP 200 | — |
| ch06 | <https://arxiv.org/abs/1609.02907> | HEAD: HTTP 200 | — |
| ch25 | <https://arxiv.org/abs/1705.08500> | HEAD: HTTP 200 | — |
| ch04, ch06 | <https://arxiv.org/abs/1706.03762> | HEAD: HTTP 200 | — |
| ch16 | <https://arxiv.org/abs/1706.04599> | HEAD: HTTP 200 | — |
| ch38 | <https://arxiv.org/abs/1708.06374> | HEAD: HTTP 200 | — |
| ch20 | <https://arxiv.org/abs/1711.00399> | HEAD: HTTP 200 | — |
| ch12 | <https://arxiv.org/abs/1711.10160> | HEAD: HTTP 200 | — |
| ch01, ch12 | <https://arxiv.org/abs/1802.06769> | HEAD: HTTP 200 | — |
| ch25 | <https://arxiv.org/abs/1807.06068> | HEAD: HTTP 200 | — |
| ch13 | <https://arxiv.org/abs/1901.04085> | HEAD: HTTP 200 | — |
| ch25 | <https://arxiv.org/abs/1910.02054> | HEAD: HTTP 200 | — |
| ch29 | <https://arxiv.org/abs/1912.03817> | HEAD: HTTP 200 | — |
| ch29 | <https://arxiv.org/abs/2002.06177> | HEAD: HTTP 200 | — |
| ch27 | <https://arxiv.org/abs/2004.10474> | HEAD: HTTP 200 | — |
| ch03, ch04, ch05, ch06, ch10, ch12, ch15, ch19, ch29 | <https://arxiv.org/abs/2005.11401> | HEAD: HTTP 200 | — |
| ch20 | <https://arxiv.org/abs/2009.01884> | HEAD: HTTP 200 | — |
| ch34 | <https://arxiv.org/abs/2011.10819> | HEAD: HTTP 200 | — |
| ch12 | <https://arxiv.org/abs/2012.07421> | HEAD: HTTP 200 | — |
| ch12 | <https://arxiv.org/abs/2012.15613> | HEAD: HTTP 200 | — |
| ch25 | <https://arxiv.org/abs/2106.09685> | HEAD: HTTP 200 | — |
| ch25 | <https://arxiv.org/abs/2107.06499> | HEAD: HTTP 200 | — |
| ch12 | <https://arxiv.org/abs/2107.07511> | HEAD: HTTP 200 | — |
| ch25 | <https://arxiv.org/abs/2110.01052> | HEAD: HTTP 200 | — |
| ch12 | <https://arxiv.org/abs/2110.09485> | HEAD: HTTP 200 | — |
| ch07 | <https://arxiv.org/abs/2112.01488> | HEAD: HTTP 200 | — |
| ch29 | <https://arxiv.org/abs/2201.11903> | HEAD: HTTP 200 | — |
| ch25 | <https://arxiv.org/abs/2203.14960> | HEAD: HTTP 200 | — |
| ch08 | <https://arxiv.org/abs/2204.08387> | HEAD: HTTP 200 | — |
| ch13 | <https://arxiv.org/abs/2212.10496> | HEAD: HTTP 200 | — |
| ch12 | <https://arxiv.org/abs/2302.13971> | HEAD: HTTP 200 | — |
| ch29 | <https://arxiv.org/abs/2305.04388> | HEAD: HTTP 200 | — |
| ch29 | <https://arxiv.org/abs/2305.13971> | HEAD: HTTP 200 | — |
| ch25 | <https://arxiv.org/abs/2305.14314> | HEAD: HTTP 200 | — |
| ch12 | <https://arxiv.org/abs/2305.15425> | HEAD: HTTP 200 | — |
| ch25 | <https://arxiv.org/abs/2305.18290> | HEAD: HTTP 200 | — |
| ch29 | <https://arxiv.org/abs/2305.20050> | HEAD: HTTP 200 | — |
| ch25 | <https://arxiv.org/abs/2306.05685> | HEAD: HTTP 200 | — |
| ch13, ch14, ch29, ch38 | <https://arxiv.org/abs/2307.09702> | HEAD: HTTP 200 | — |
| ch10, ch25 | <https://arxiv.org/abs/2309.15217> | HEAD: HTTP 200 | — |
| ch29 | <https://arxiv.org/abs/2310.03714> | HEAD: HTTP 200 | — |
| ch13 | <https://arxiv.org/abs/2310.05736> | HEAD: HTTP 200 | — |
| ch12 | <https://arxiv.org/abs/2310.06825> | HEAD: HTTP 200 | — |
| ch10, ch13 | <https://arxiv.org/abs/2310.11511> | HEAD: HTTP 200 | — |
| ch34 | <https://arxiv.org/abs/2310.13548> | HEAD: HTTP 200 | — |
| ch01 | <https://arxiv.org/abs/2311.04910> | HEAD: HTTP 200 | — |
| ch13 | <https://arxiv.org/abs/2312.10997> | HEAD: HTTP 200 | — |
| ch13 | <https://arxiv.org/abs/2401.15884> | HEAD: HTTP 200 | — |
| ch12 | <https://arxiv.org/abs/2403.08295> | HEAD: HTTP 200 | — |
| ch25 | <https://arxiv.org/abs/2403.13372> | HEAD: HTTP 200 | — |
| ch07, ch13 | <https://arxiv.org/abs/2404.16130> | HEAD: HTTP 200 | — |
| ch13 | <https://arxiv.org/abs/2411.15100> | HEAD: HTTP 200 | — |
| ch20, ch29 | <https://arxiv.org/abs/2501.05435> | HEAD: HTTP 200 | — |
| ch08 | <https://arxiv.org/abs/2502.13595> | HEAD: HTTP 200 | — |
| ch13 | <https://arxiv.org/abs/2507.18546> | HEAD: HTTP 200 | — |
| ch29 | <https://arxiv.org/abs/2509.04664> | HEAD: HTTP 200 | — |
| ch15 | <https://atlas.mitre.org/> | HEAD: HTTP 200 | — |
| ch06 | <https://bayes.cs.ucla.edu/BOOK-2K/> | HEAD: HTTP 200 | — |
| ch32 | <https://bibliopolis.it/shop/intuitionistic-type-theory/> | HEAD: HTTP 200 | — |
| ch18 | <https://blogs.windows.com/devices/2026/06/02/building-the-next-generation-of-devices-for-developers-surface-rtx-spark-dev-box/> | HEAD: HTTP 200 | — |
| ch18 | <https://blogs.windows.com/devices/2026/10/07/pre-order-our-most-powerful-surface-devices-ever/> | HEAD: HTTP 200 | — |
| ch05 | <https://books.google.com/books/about/The_Fifth_Discipline.html?id=bVZqAAAAMAAJ> | HEAD: HTTP 200 | 302 → https://books.google.com.ua/books/about/The_Fifth_Discipline.html?id=bVZqAAAAMAAJ&redir_esc=y |
| ch05 | <https://books.google.com/books/about/Working_Knowledge.html?id=-4-7vmCVG5cC> | HEAD: HTTP 200 | 302 → https://books.google.com.ua/books/about/Working_Knowledge.html?id=-4-7vmCVG5cC&redir_esc=y |
| ch25 | <https://buildkite.com/docs/agent/v3> | HEAD: HTTP 200 | 301 → https://buildkite.com/docs/agent |
| ch25 | <https://buildkite.com/docs/pipelines/architecture> | HEAD: HTTP 200 | — |
| ch20 | <https://cdn.aaai.org/AAAI/2004/AAAI04-027.pdf> | HEAD: HTTP 200 | — |
| ch01 | <https://ceur-ws.org/Vol-716/ICTERI-2011-CEUR-WS-paper-3-p-51-64.pdf> | HEAD: HTTP 200 | — |
| ch10 | <https://cheatsheetseries.owasp.org/cheatsheets/File_Upload_Cheat_Sheet.html> | HEAD: HTTP 200 | — |
| ch17 | <https://codeberg.org/TauCeti/mangle-go> | HEAD: HTTP 200 | — |
| ch05 | <https://cognitect.com/blog/2011/11/15/documenting-architecture-decisions> | HEAD: HTTP 200 | — |
| ch18 | <https://computeexpresslink.org/> | HEAD: HTTP 200 | — |
| ch10 | <https://csrc.nist.gov/pubs/ai/100/2/e2025/final> | HEAD: HTTP 200 | — |
| ch10 | <https://csrc.nist.gov/pubs/fips/198-1/final> | HEAD: HTTP 200 | — |
| ch10 | <https://csrc.nist.gov/pubs/sp/800/224/ipd> | HEAD: HTTP 200 | — |
| ch37 | <https://csrc.nist.gov/pubs/sp/800/61/r3/final> | HEAD: HTTP 200 | — |
| ch14 | <https://cvc5.github.io/docs/latest/> | HEAD: HTTP 200 | — |
| ch10, ch30 | <https://cyclonedx.org/specification/overview/> | HEAD: HTTP 200 | — |
| ch26 | <https://datasets-benchmarks-proceedings.neurips.cc/paper/2021/hash/33e75ff09dd601bbe69f351039152189-Abstract-round2.html> | HEAD: HTTP 200 | — |
| ch37 | <https://datawalk.com/> | HEAD: HTTP 200 | — |
| ch29 | <https://deepmind.google/discover/blog/ai-solves-imo-problems-at-silver-medal-level/> | HEAD: HTTP 200 | 301 → https://deepmind.google/blog/ai-solves-imo-problems-at-silver-medal-level/ |
| ch18 | <https://developer.nvidia.com/blog/accelerating-inference-with-sparsity-using-ampere-and-tensorrt/> | HEAD: HTTP 200 | — |
| ch18 | <https://developer.nvidia.com/embedded/learn/jetson-agx-orin-devkit-user-guide> | HEAD: HTTP 200 | 301 → https://developer.nvidia.com/embedded/learn/jetson-agx-orin-devkit-user-guide/index.html<br>301 → https://docs.nvidia.com/jetson/agx-orin-devkit/user-guide/latest/index.html |
| ch07, ch17 | <https://developers.google.com/optimization/cp/cp_solver> | HEAD: HTTP 200 | — |
| ch25 | <https://doc.dvc.org/user-guide/project-structure/dvc-files> | HEAD: HTTP 200 | — |
| ch12 | <https://docling-project.github.io/docling/> | HEAD: HTTP 200 | — |
| ch12 | <https://docling-project.github.io/docling/concepts/docling_document/> | HEAD: HTTP 200 | — |
| ch25 | <https://docs.axolotl.ai/> | HEAD: HTTP 200 | — |
| ch21 | <https://docs.cedarpolicy.com/> | HEAD: HTTP 200 | — |
| ch17 | <https://docs.drools.org/latest/drools-docs/drools/rule-engine/index.html> | HEAD: HTTP 200 | — |
| ch18 | <https://docs.nats.io/nats-concepts/jetstream> | HEAD: HTTP 200 | 307 → https://docs.nats.io/concepts/jetstream |
| ch25 | <https://docs.nvidia.com/deeplearning/transformer-engine/index.html> | HEAD: HTTP 200 | — |
| ch18 | <https://docs.nvidia.com/dgx/dgx-spark/hardware.html> | HEAD: HTTP 200 | — |
| ch18 | <https://docs.nvidia.com/jetson/archives/r36.2/DeveloperGuide/SD/PlatformPowerAndPerformance/JetsonOrinNanoSeriesJetsonOrinNxSeriesAndJetsonAgxOrinSeries.html> | HEAD: HTTP 200 | — |
| ch09 | <https://docs.oasis-open-projects.org/oslc-op/qm/v2.1/os/quality-management-vocab.html> | HEAD: HTTP 200 | — |
| ch09 | <https://docs.oasis-open-projects.org/oslc-op/rm/v2.1/os/requirements-management-vocab.html> | HEAD: HTTP 200 | — |
| ch37 | <https://docs.oasis-open.org/cti/stix/v2.1/os/stix-v2.1-os.html> | HEAD: HTTP 200 | — |
| ch04 | <https://docs.oasis-open.org/legalruleml/legalruleml-core-spec/v1.0/os/legalruleml-core-spec-v1.0-os.html> | HEAD: HTTP 200 | — |
| ch10 | <https://docs.oasis-open.org/oslc-core/oslc-core/v3.0/cs01/part1-overview/oslc-core-v3.0-cs01-part1-overview.html> | HEAD: HTTP 200 | — |
| ch29 | <https://docs.ollama.com/capabilities/structured-outputs> | HEAD: HTTP 200 | — |
| ch17 | <https://docs.openvino.ai/> | HEAD: HTTP 200 | 302 → https://docs.openvino.ai/index.html<br>302 → https://docs.openvino.ai/2026/index.html |
| ch18 | <https://docs.openvino.ai/2024/openvino-workflow/running-inference/inference-devices-and-modes/npu-device.html> | HEAD: HTTP 200 | — |
| ch10 | <https://docs.openvino.ai/2026/about-openvino/release-notes-openvino.html> | HEAD: HTTP 200 | — |
| ch16 | <https://docs.pact.io/> | HEAD: HTTP 200 | — |
| ch32 | <https://docs.python.org/3/library/mmap.html> | HEAD: HTTP 200 | — |
| ch32 | <https://docs.python.org/3/library/sqlite3.html> | HEAD: HTTP 200 | — |
| ch25 | <https://docs.pytorch.org/docs/main/distributed.fsdp.fully_shard.html> | HEAD: HTTP 200 | — |
| ch27 | <https://docs.sigstore.dev/logging/overview/> | HEAD: HTTP 200 | — |
| ch21 | <https://docs.temporal.io/activities> | HEAD: HTTP 200 | — |
| ch25 | <https://doi.org/10.1001/jama.1983.03330370053031> | GET (HEAD 403 → ranged GET): HTTP 403 | 302 → http://jama.jamanetwork.com/article.aspx?doi=10.1001/jama.1983.03330370053031<br>301 → https://jama.jamanetwork.com/article.aspx?doi=10.1001/jama.1983.03330370053031<br>301 → https://jamanetwork.com/article.aspx?doi=10.1001/jama.1983.03330370053031 |
| ch06, ch33 | <https://doi.org/10.1002/0471221279> | GET (HEAD 403 → ranged GET): HTTP 403 | 302 → https://onlinelibrary.wiley.com/doi/book/10.1002/0471221279 |
| ch24 | <https://doi.org/10.1002/047174882X> | GET (HEAD 403 → ranged GET): HTTP 403 | 302 → https://onlinelibrary.wiley.com/doi/book/10.1002/047174882X |
| ch06 | <https://doi.org/10.1002/9780470316887> | GET (HEAD 403 → ranged GET): HTTP 403 | 302 → https://onlinelibrary.wiley.com/doi/book/10.1002/9780470316887 |
| ch09 | <https://doi.org/10.1002/asi.20591> | GET (HEAD 403 → ranged GET): HTTP 403 | 302 → https://onlinelibrary.wiley.com/doi/10.1002/asi.20591 |
| ch11 | <https://doi.org/10.1002/for.3980010305> | GET (HEAD 403 → ranged GET): HTTP 403 | 302 → https://onlinelibrary.wiley.com/doi/10.1002/for.3980010305 |
| ch23 | <https://doi.org/10.1002/int.4550090804> | GET (HEAD 403 → ranged GET): HTTP 403 | 302 → https://onlinelibrary.wiley.com/doi/10.1002/int.4550090804 |
| ch11, ch24 | <https://doi.org/10.1002/j.1538-7305.1948.tb01338.x> | HEAD: HTTP 202 | 302 → https://ieeexplore.ieee.org/document/6773024 |
| ch11 | <https://doi.org/10.1006/ijhc.1994.1083> | HEAD: HTTP 200 | 302 → https://linkinghub.elsevier.com/retrieve/pii/S1071581984710834 |
| ch33 | <https://doi.org/10.1006/knac.1993.1008> | HEAD: HTTP 200 | 302 → https://linkinghub.elsevier.com/retrieve/pii/S1042814383710083 |
| ch11 | <https://doi.org/10.1006/obhd.1995.1039> | HEAD: HTTP 200 | 302 → https://linkinghub.elsevier.com/retrieve/pii/S0749597885710394 |
| ch27 | <https://doi.org/10.1007/0-387-34805-0_21> | HEAD: HTTP 200 | 302 → http://link.springer.com/10.1007/0-387-34805-0_21<br>301 → https://link.springer.com/10.1007/0-387-34805-0_21<br>301 → https://link.springer.com/chapter/10.1007/0-387-34805-0_21<br>303 → https://idp.springer.com/authorize?response_type=cookie&client_id=springerlink&redirect_uri=https%3A%2F%2Flink.springer.com%2Fchapter%2F10.1007%2F0-387-34805-0_21<br>302 → https://idp.springer.com/transit?redirect_uri=https%3A%2F%2Flink.springer.com%2Fchapter%2F10.1007%2F0-387-34805-0_21&code=cddcdc16-1322-4904-974a-adb3c0b3c6cc<br>302 → https://link.springer.com/chapter/10.1007/0-387-34805-0_21?error=cookies_not_supported&code=cddcdc16-1322-4904-974a-adb3c0b3c6cc |
| ch21 | <https://doi.org/10.1007/11609773_24> | HEAD: HTTP 200 | 302 → http://link.springer.com/10.1007/11609773_24<br>301 → https://link.springer.com/10.1007/11609773_24<br>301 → https://link.springer.com/chapter/10.1007/11609773_24<br>303 → https://idp.springer.com/authorize?response_type=cookie&client_id=springerlink&redirect_uri=https%3A%2F%2Flink.springer.com%2Fchapter%2F10.1007%2F11609773_24<br>302 → https://idp.springer.com/transit?redirect_uri=https%3A%2F%2Flink.springer.com%2Fchapter%2F10.1007%2F11609773_24&code=17595028-e758-4175-8b03-5c529c063000<br>302 → https://link.springer.com/chapter/10.1007/11609773_24?error=cookies_not_supported&code=17595028-e758-4175-8b03-5c529c063000 |
| ch19 | <https://doi.org/10.1007/11736790_9> | HEAD: HTTP 200 | 302 → http://link.springer.com/10.1007/11736790_9<br>301 → https://link.springer.com/10.1007/11736790_9<br>301 → https://link.springer.com/chapter/10.1007/11736790_9<br>303 → https://idp.springer.com/authorize?response_type=cookie&client_id=springerlink&redirect_uri=https%3A%2F%2Flink.springer.com%2Fchapter%2F10.1007%2F11736790_9<br>302 → https://idp.springer.com/transit?redirect_uri=https%3A%2F%2Flink.springer.com%2Fchapter%2F10.1007%2F11736790_9&code=6862113d-7a92-4955-bd39-891517f0027f<br>302 → https://link.springer.com/chapter/10.1007/11736790_9?error=cookies_not_supported&code=6862113d-7a92-4955-bd39-891517f0027f |
| ch33 | <https://doi.org/10.1007/3-540-48184-2_32> | HEAD: HTTP 200 | 302 → http://link.springer.com/10.1007/3-540-48184-2_32<br>301 → https://link.springer.com/10.1007/3-540-48184-2_32<br>301 → https://link.springer.com/chapter/10.1007/3-540-48184-2_32<br>303 → https://idp.springer.com/authorize?response_type=cookie&client_id=springerlink&redirect_uri=https%3A%2F%2Flink.springer.com%2Fchapter%2F10.1007%2F3-540-48184-2_32<br>302 → https://idp.springer.com/transit?redirect_uri=https%3A%2F%2Flink.springer.com%2Fchapter%2F10.1007%2F3-540-48184-2_32&code=cf98c3a3-e502-48cf-bb92-bfb8b6c7ae40<br>302 → https://link.springer.com/chapter/10.1007/3-540-48184-2_32?error=cookies_not_supported&code=cf98c3a3-e502-48cf-bb92-bfb8b6c7ae40 |
| ch32 | <https://doi.org/10.1007/3-540-62222-5_39> | HEAD: HTTP 200 | 302 → http://link.springer.com/10.1007/3-540-62222-5_39<br>301 → https://link.springer.com/10.1007/3-540-62222-5_39<br>301 → https://link.springer.com/chapter/10.1007/3-540-62222-5_39<br>303 → https://idp.springer.com/authorize?response_type=cookie&client_id=springerlink&redirect_uri=https%3A%2F%2Flink.springer.com%2Fchapter%2F10.1007%2F3-540-62222-5_39<br>302 → https://idp.springer.com/transit?redirect_uri=https%3A%2F%2Flink.springer.com%2Fchapter%2F10.1007%2F3-540-62222-5_39&code=9c9fb2fb-9310-4c43-9b12-aafeec995b04<br>302 → https://link.springer.com/chapter/10.1007/3-540-62222-5_39?error=cookies_not_supported&code=9c9fb2fb-9310-4c43-9b12-aafeec995b04 |
| ch24 | <https://doi.org/10.1007/978-0-387-68282-2> | HEAD: HTTP 200 | 302 → http://link.springer.com/10.1007/978-0-387-68282-2<br>301 → https://link.springer.com/10.1007/978-0-387-68282-2<br>301 → https://link.springer.com/book/10.1007/978-0-387-68282-2<br>303 → https://idp.springer.com/authorize?response_type=cookie&client_id=springerlink&redirect_uri=https%3A%2F%2Flink.springer.com%2Fbook%2F10.1007%2F978-0-387-68282-2<br>302 → https://idp.springer.com/transit?redirect_uri=https%3A%2F%2Flink.springer.com%2Fbook%2F10.1007%2F978-0-387-68282-2&code=f1cba43b-dde5-43f9-a7db-86a126bc4788<br>302 → https://link.springer.com/book/10.1007/978-0-387-68282-2?error=cookies_not_supported&code=f1cba43b-dde5-43f9-a7db-86a126bc4788 |
| ch29 | <https://doi.org/10.1007/978-1-4471-0211-3> | HEAD: HTTP 200 | 302 → http://link.springer.com/10.1007/978-1-4471-0211-3<br>301 → https://link.springer.com/10.1007/978-1-4471-0211-3<br>301 → https://link.springer.com/book/10.1007/978-1-4471-0211-3<br>303 → https://idp.springer.com/authorize?response_type=cookie&client_id=springerlink&redirect_uri=https%3A%2F%2Flink.springer.com%2Fbook%2F10.1007%2F978-1-4471-0211-3<br>302 → https://idp.springer.com/transit?redirect_uri=https%3A%2F%2Flink.springer.com%2Fbook%2F10.1007%2F978-1-4471-0211-3&code=dd356ae1-d2b2-4efb-9a3f-1d400931352f<br>302 → https://link.springer.com/book/10.1007/978-1-4471-0211-3?error=cookies_not_supported&code=dd356ae1-d2b2-4efb-9a3f-1d400931352f |
| ch02, ch06, ch07, ch32 | <https://doi.org/10.1007/978-1-4684-3384-5_3> | HEAD: HTTP 200 | 302 → http://link.springer.com/10.1007/978-1-4684-3384-5_3<br>301 → https://link.springer.com/10.1007/978-1-4684-3384-5_3<br>301 → https://link.springer.com/chapter/10.1007/978-1-4684-3384-5_3<br>303 → https://idp.springer.com/authorize?response_type=cookie&client_id=springerlink&redirect_uri=https%3A%2F%2Flink.springer.com%2Fchapter%2F10.1007%2F978-1-4684-3384-5_3<br>302 → https://idp.springer.com/transit?redirect_uri=https%3A%2F%2Flink.springer.com%2Fchapter%2F10.1007%2F978-1-4684-3384-5_3&code=9ef564ac-898c-4cca-9fed-2f1b0857c675<br>302 → https://link.springer.com/chapter/10.1007/978-1-4684-3384-5_3?error=cookies_not_supported&code=9ef564ac-898c-4cca-9fed-2f1b0857c675 |
| ch26 | <https://doi.org/10.1007/978-3-030-01252-6_33> | HEAD: HTTP 200 | 302 → https://link.springer.com/10.1007/978-3-030-01252-6_33<br>301 → https://link.springer.com/chapter/10.1007/978-3-030-01252-6_33<br>303 → https://idp.springer.com/authorize?response_type=cookie&client_id=springerlink&redirect_uri=https%3A%2F%2Flink.springer.com%2Fchapter%2F10.1007%2F978-3-030-01252-6_33<br>302 → https://idp.springer.com/transit?redirect_uri=https%3A%2F%2Flink.springer.com%2Fchapter%2F10.1007%2F978-3-030-01252-6_33&code=29d37b32-b22d-4dbc-99c5-2d0127ddab2f<br>302 → https://link.springer.com/chapter/10.1007/978-3-030-01252-6_33?error=cookies_not_supported&code=29d37b32-b22d-4dbc-99c5-2d0127ddab2f |
| ch07 | <https://doi.org/10.1007/978-3-030-26253-2> | HEAD: HTTP 200 | 302 → http://link.springer.com/10.1007/978-3-030-26253-2<br>301 → https://link.springer.com/10.1007/978-3-030-26253-2<br>301 → https://link.springer.com/book/10.1007/978-3-030-26253-2<br>303 → https://idp.springer.com/authorize?response_type=cookie&client_id=springerlink&redirect_uri=https%3A%2F%2Flink.springer.com%2Fbook%2F10.1007%2F978-3-030-26253-2<br>302 → https://idp.springer.com/transit?redirect_uri=https%3A%2F%2Flink.springer.com%2Fbook%2F10.1007%2F978-3-030-26253-2&code=5e44c43a-58bf-4458-91cd-8c9ee724c34d<br>302 → https://link.springer.com/book/10.1007/978-3-030-26253-2?error=cookies_not_supported&code=5e44c43a-58bf-4458-91cd-8c9ee724c34d |
| ch29 | <https://doi.org/10.1007/978-3-030-79876-5_37> | HEAD: HTTP 200 | 302 → https://link.springer.com/10.1007/978-3-030-79876-5_37<br>301 → https://link.springer.com/chapter/10.1007/978-3-030-79876-5_37<br>303 → https://idp.springer.com/authorize?response_type=cookie&client_id=springerlink&redirect_uri=https%3A%2F%2Flink.springer.com%2Fchapter%2F10.1007%2F978-3-030-79876-5_37<br>302 → https://idp.springer.com/transit?redirect_uri=https%3A%2F%2Flink.springer.com%2Fchapter%2F10.1007%2F978-3-030-79876-5_37&code=1246330d-b1e1-4180-8532-ea563a5ef374<br>302 → https://link.springer.com/chapter/10.1007/978-3-030-79876-5_37?error=cookies_not_supported&code=1246330d-b1e1-4180-8532-ea563a5ef374 |
| ch17 | <https://doi.org/10.1007/978-3-030-99524-9_24> | HEAD: HTTP 200 | 302 → https://link.springer.com/10.1007/978-3-030-99524-9_24<br>301 → https://link.springer.com/chapter/10.1007/978-3-030-99524-9_24<br>303 → https://idp.springer.com/authorize?response_type=cookie&client_id=springerlink&redirect_uri=https%3A%2F%2Flink.springer.com%2Fchapter%2F10.1007%2F978-3-030-99524-9_24<br>302 → https://idp.springer.com/transit?redirect_uri=https%3A%2F%2Flink.springer.com%2Fchapter%2F10.1007%2F978-3-030-99524-9_24&code=dd90a648-b052-4a57-ad86-27dae102a425<br>302 → https://link.springer.com/chapter/10.1007/978-3-030-99524-9_24?error=cookies_not_supported&code=dd90a648-b052-4a57-ad86-27dae102a425 |
| ch14 | <https://doi.org/10.1007/978-3-540-30206-3_12> | HEAD: HTTP 200 | 302 → https://link.springer.com/10.1007/978-3-540-30206-3_12<br>301 → https://link.springer.com/chapter/10.1007/978-3-540-30206-3_12<br>303 → https://idp.springer.com/authorize?response_type=cookie&client_id=springerlink&redirect_uri=https%3A%2F%2Flink.springer.com%2Fchapter%2F10.1007%2F978-3-540-30206-3_12<br>302 → https://idp.springer.com/transit?redirect_uri=https%3A%2F%2Flink.springer.com%2Fchapter%2F10.1007%2F978-3-540-30206-3_12&code=9511354b-f18c-46cc-a585-6ddae97bd3eb<br>302 → https://link.springer.com/chapter/10.1007/978-3-540-30206-3_12?error=cookies_not_supported&code=9511354b-f18c-46cc-a585-6ddae97bd3eb |
| ch01 | <https://doi.org/10.1007/978-3-540-44616-3_18> | HEAD: HTTP 200 | 302 → https://link.springer.com/10.1007/978-3-540-44616-3_18<br>301 → https://link.springer.com/chapter/10.1007/978-3-540-44616-3_18<br>303 → https://idp.springer.com/authorize?response_type=cookie&client_id=springerlink&redirect_uri=https%3A%2F%2Flink.springer.com%2Fchapter%2F10.1007%2F978-3-540-44616-3_18<br>302 → https://idp.springer.com/transit?redirect_uri=https%3A%2F%2Flink.springer.com%2Fchapter%2F10.1007%2F978-3-540-44616-3_18&code=bd1c0ccf-7404-4435-90ad-e06df14bafeb<br>302 → https://link.springer.com/chapter/10.1007/978-3-540-44616-3_18?error=cookies_not_supported&code=bd1c0ccf-7404-4435-90ad-e06df14bafeb |
| ch34 | <https://doi.org/10.1007/978-3-540-68856-3> | HEAD: HTTP 200 | 302 → http://link.springer.com/10.1007/978-3-540-68856-3<br>301 → https://link.springer.com/10.1007/978-3-540-68856-3<br>301 → https://link.springer.com/book/10.1007/978-3-540-68856-3<br>303 → https://idp.springer.com/authorize?response_type=cookie&client_id=springerlink&redirect_uri=https%3A%2F%2Flink.springer.com%2Fbook%2F10.1007%2F978-3-540-68856-3<br>302 → https://idp.springer.com/transit?redirect_uri=https%3A%2F%2Flink.springer.com%2Fbook%2F10.1007%2F978-3-540-68856-3&code=672ba9c9-3316-40f5-8475-c4fd86c2555e<br>302 → https://link.springer.com/book/10.1007/978-3-540-68856-3?error=cookies_not_supported&code=672ba9c9-3316-40f5-8475-c4fd86c2555e |
| ch17 | <https://doi.org/10.1007/978-3-540-74970-7_38> | HEAD: HTTP 200 | 302 → http://link.springer.com/10.1007/978-3-540-74970-7_38<br>301 → https://link.springer.com/10.1007/978-3-540-74970-7_38<br>301 → https://link.springer.com/chapter/10.1007/978-3-540-74970-7_38<br>303 → https://idp.springer.com/authorize?response_type=cookie&client_id=springerlink&redirect_uri=https%3A%2F%2Flink.springer.com%2Fchapter%2F10.1007%2F978-3-540-74970-7_38<br>302 → https://idp.springer.com/transit?redirect_uri=https%3A%2F%2Flink.springer.com%2Fchapter%2F10.1007%2F978-3-540-74970-7_38&code=f7884308-7042-4e76-b6cc-71680e6c2399<br>302 → https://link.springer.com/chapter/10.1007/978-3-540-74970-7_38?error=cookies_not_supported&code=f7884308-7042-4e76-b6cc-71680e6c2399 |
| ch03, ch14, ch17, ch23, ch33 | <https://doi.org/10.1007/978-3-540-78800-3_24> | HEAD: HTTP 200 | 302 → http://link.springer.com/10.1007/978-3-540-78800-3_24<br>301 → https://link.springer.com/10.1007/978-3-540-78800-3_24<br>301 → https://link.springer.com/chapter/10.1007/978-3-540-78800-3_24<br>303 → https://idp.springer.com/authorize?response_type=cookie&client_id=springerlink&redirect_uri=https%3A%2F%2Flink.springer.com%2Fchapter%2F10.1007%2F978-3-540-78800-3_24<br>302 → https://idp.springer.com/transit?redirect_uri=https%3A%2F%2Flink.springer.com%2Fchapter%2F10.1007%2F978-3-540-78800-3_24&code=3ef8b9d0-328c-41b2-ac17-51fc20f1a599<br>302 → https://link.springer.com/chapter/10.1007/978-3-540-78800-3_24?error=cookies_not_supported&code=3ef8b9d0-328c-41b2-ac17-51fc20f1a599 |
| ch33 | <https://doi.org/10.1007/978-3-642-24550-3_29> | HEAD: HTTP 200 | 302 → http://link.springer.com/10.1007/978-3-642-24550-3_29<br>301 → https://link.springer.com/10.1007/978-3-642-24550-3_29<br>301 → https://link.springer.com/chapter/10.1007/978-3-642-24550-3_29<br>303 → https://idp.springer.com/authorize?response_type=cookie&client_id=springerlink&redirect_uri=https%3A%2F%2Flink.springer.com%2Fchapter%2F10.1007%2F978-3-642-24550-3_29<br>302 → https://idp.springer.com/transit?redirect_uri=https%3A%2F%2Flink.springer.com%2Fchapter%2F10.1007%2F978-3-642-24550-3_29&code=97fceaa3-2b10-4eb6-b6f5-f0fa6130b053<br>302 → https://link.springer.com/chapter/10.1007/978-3-642-24550-3_29?error=cookies_not_supported&code=97fceaa3-2b10-4eb6-b6f5-f0fa6130b053 |
| ch06 | <https://doi.org/10.1007/978-3-642-48318-9> | HEAD: HTTP 200 | 302 → http://link.springer.com/10.1007/978-3-642-48318-9<br>301 → https://link.springer.com/10.1007/978-3-642-48318-9<br>301 → https://link.springer.com/book/10.1007/978-3-642-48318-9<br>303 → https://idp.springer.com/authorize?response_type=cookie&client_id=springerlink&redirect_uri=https%3A%2F%2Flink.springer.com%2Fbook%2F10.1007%2F978-3-642-48318-9<br>302 → https://idp.springer.com/transit?redirect_uri=https%3A%2F%2Flink.springer.com%2Fbook%2F10.1007%2F978-3-642-48318-9&code=308a25d4-de48-454b-9c66-ec963d7975f6<br>302 → https://link.springer.com/book/10.1007/978-3-642-48318-9?error=cookies_not_supported&code=308a25d4-de48-454b-9c66-ec963d7975f6 |
| ch04 | <https://doi.org/10.1007/978-3-642-49877-0_8> | HEAD: HTTP 200 | 302 → http://link.springer.com/10.1007/978-3-642-49877-0_8<br>301 → https://link.springer.com/10.1007/978-3-642-49877-0_8<br>301 → https://link.springer.com/chapter/10.1007/978-3-642-49877-0_8<br>303 → https://idp.springer.com/authorize?response_type=cookie&client_id=springerlink&redirect_uri=https%3A%2F%2Flink.springer.com%2Fchapter%2F10.1007%2F978-3-642-49877-0_8<br>302 → https://idp.springer.com/transit?redirect_uri=https%3A%2F%2Flink.springer.com%2Fchapter%2F10.1007%2F978-3-642-49877-0_8&code=ac506f3d-c3cc-46a4-95e2-364d097d7b42<br>302 → https://link.springer.com/chapter/10.1007/978-3-642-49877-0_8?error=cookies_not_supported&code=ac506f3d-c3cc-46a4-95e2-364d097d7b42 |
| ch06, ch31 | <https://doi.org/10.1007/978-3-642-59830-2> | HEAD: HTTP 200 | 302 → http://link.springer.com/10.1007/978-3-642-59830-2<br>301 → https://link.springer.com/10.1007/978-3-642-59830-2<br>301 → https://link.springer.com/book/10.1007/978-3-642-59830-2<br>303 → https://idp.springer.com/authorize?response_type=cookie&client_id=springerlink&redirect_uri=https%3A%2F%2Flink.springer.com%2Fbook%2F10.1007%2F978-3-642-59830-2<br>302 → https://idp.springer.com/transit?redirect_uri=https%3A%2F%2Flink.springer.com%2Fbook%2F10.1007%2F978-3-642-59830-2&code=3ea26365-cf9e-4e8e-b404-382ef1ba2524<br>302 → https://link.springer.com/book/10.1007/978-3-642-59830-2?error=cookies_not_supported&code=3ea26365-cf9e-4e8e-b404-382ef1ba2524 |
| ch06 | <https://doi.org/10.1007/978-3-642-88338-5> | HEAD: HTTP 200 | 302 → https://link.springer.com/10.1007/978-3-642-88338-5<br>301 → https://link.springer.com/book/10.1007/978-3-642-88338-5<br>303 → https://idp.springer.com/authorize?response_type=cookie&client_id=springerlink&redirect_uri=https%3A%2F%2Flink.springer.com%2Fbook%2F10.1007%2F978-3-642-88338-5<br>302 → https://idp.springer.com/transit?redirect_uri=https%3A%2F%2Flink.springer.com%2Fbook%2F10.1007%2F978-3-642-88338-5&code=24220c97-aa35-457e-a007-6f8c07282a54<br>302 → https://link.springer.com/book/10.1007/978-3-642-88338-5?error=cookies_not_supported&code=24220c97-aa35-457e-a007-6f8c07282a54 |
| ch21 | <https://doi.org/10.1007/978-3-662-46681-0_33> | HEAD: HTTP 200 | 302 → http://link.springer.com/10.1007/978-3-662-46681-0_33<br>301 → https://link.springer.com/10.1007/978-3-662-46681-0_33<br>301 → https://link.springer.com/chapter/10.1007/978-3-662-46681-0_33<br>303 → https://idp.springer.com/authorize?response_type=cookie&client_id=springerlink&redirect_uri=https%3A%2F%2Flink.springer.com%2Fchapter%2F10.1007%2F978-3-662-46681-0_33<br>302 → https://idp.springer.com/transit?redirect_uri=https%3A%2F%2Flink.springer.com%2Fchapter%2F10.1007%2F978-3-662-46681-0_33&code=6705374c-e792-40b2-95b1-7eed3e9558af<br>302 → https://link.springer.com/chapter/10.1007/978-3-662-46681-0_33?error=cookies_not_supported&code=6705374c-e792-40b2-95b1-7eed3e9558af |
| ch29 | <https://doi.org/10.1007/978-3-662-49896-5_11> | HEAD: HTTP 200 | 302 → http://link.springer.com/10.1007/978-3-662-49896-5_11<br>301 → https://link.springer.com/10.1007/978-3-662-49896-5_11<br>301 → https://link.springer.com/chapter/10.1007/978-3-662-49896-5_11<br>303 → https://idp.springer.com/authorize?response_type=cookie&client_id=springerlink&redirect_uri=https%3A%2F%2Flink.springer.com%2Fchapter%2F10.1007%2F978-3-662-49896-5_11<br>302 → https://idp.springer.com/transit?redirect_uri=https%3A%2F%2Flink.springer.com%2Fchapter%2F10.1007%2F978-3-662-49896-5_11&code=7fd701ab-9da0-4665-93c4-f1920a8b2cb0<br>302 → https://link.springer.com/chapter/10.1007/978-3-662-49896-5_11?error=cookies_not_supported&code=7fd701ab-9da0-4665-93c4-f1920a8b2cb0 |
| ch02 | <https://doi.org/10.1007/978-94-010-1161-7_2> | HEAD: HTTP 200 | 302 → http://link.springer.com/10.1007/978-94-010-1161-7_2<br>301 → https://link.springer.com/10.1007/978-94-010-1161-7_2<br>301 → https://link.springer.com/chapter/10.1007/978-94-010-1161-7_2<br>303 → https://idp.springer.com/authorize?response_type=cookie&client_id=springerlink&redirect_uri=https%3A%2F%2Flink.springer.com%2Fchapter%2F10.1007%2F978-94-010-1161-7_2<br>302 → https://idp.springer.com/transit?redirect_uri=https%3A%2F%2Flink.springer.com%2Fchapter%2F10.1007%2F978-94-010-1161-7_2&code=f12c1199-4d67-4d3e-946e-d6ee5829104b<br>302 → https://link.springer.com/chapter/10.1007/978-94-010-1161-7_2?error=cookies_not_supported&code=f12c1199-4d67-4d3e-946e-d6ee5829104b |
| ch14 | <https://doi.org/10.1007/BF01995674> | HEAD: HTTP 200 | 302 → http://link.springer.com/10.1007/BF01995674<br>301 → https://link.springer.com/10.1007/BF01995674<br>301 → https://link.springer.com/article/10.1007/BF01995674<br>303 → https://idp.springer.com/authorize?response_type=cookie&client_id=springerlink&redirect_uri=https%3A%2F%2Flink.springer.com%2Farticle%2F10.1007%2FBF01995674<br>302 → https://idp.springer.com/transit?redirect_uri=https%3A%2F%2Flink.springer.com%2Farticle%2F10.1007%2FBF01995674&code=82577824-1a51-431a-95b8-96d402763686<br>302 → https://link.springer.com/article/10.1007/BF01995674?error=cookies_not_supported&code=82577824-1a51-431a-95b8-96d402763686 |
| ch25 | <https://doi.org/10.1007/BF02295996> | HEAD: HTTP 200 | 302 → https://www.cambridge.org/core/product/identifier/S0033312300045178/type/journal_article<br>301 → https://www.cambridge.org/core/journals/psychometrika/article/note-on-the-sampling-error-of-the-difference-between-correlated-proportions-or-percentages/698C2461BE63F5848763502D54E534FD<br>302 → https://www.cambridge.org/core/journals/psychometrika/article/abs/note-on-the-sampling-error-of-the-difference-between-correlated-proportions-or-percentages/698C2461BE63F5848763502D54E534FD |
| ch26 | <https://doi.org/10.1007/BF03037089> | HEAD: HTTP 200 | 302 → http://link.springer.com/10.1007/BF03037089<br>301 → https://link.springer.com/10.1007/BF03037089<br>301 → https://link.springer.com/article/10.1007/BF03037089<br>303 → https://idp.springer.com/authorize?response_type=cookie&client_id=springerlink&redirect_uri=https%3A%2F%2Flink.springer.com%2Farticle%2F10.1007%2FBF03037089<br>302 → https://idp.springer.com/transit?redirect_uri=https%3A%2F%2Flink.springer.com%2Farticle%2F10.1007%2FBF03037089&code=4a93d8cc-ae2f-4656-a397-e1072821bfb2<br>302 → https://link.springer.com/article/10.1007/BF03037089?error=cookies_not_supported&code=4a93d8cc-ae2f-4656-a397-e1072821bfb2 |
| ch22 | <https://doi.org/10.1007/b97451> | HEAD: HTTP 200 | 302 → http://link.springer.com/10.1007/b97451<br>301 → https://link.springer.com/10.1007/b97451<br>301 → https://link.springer.com/book/10.1007/b97451<br>303 → https://idp.springer.com/authorize?response_type=cookie&client_id=springerlink&redirect_uri=https%3A%2F%2Flink.springer.com%2Fbook%2F10.1007%2Fb97451<br>302 → https://idp.springer.com/transit?redirect_uri=https%3A%2F%2Flink.springer.com%2Fbook%2F10.1007%2Fb97451&code=ca6a811e-63be-4892-8fac-4efadeadb6ed<br>302 → https://link.springer.com/book/10.1007/b97451?error=cookies_not_supported&code=ca6a811e-63be-4892-8fac-4efadeadb6ed |
| ch18 | <https://doi.org/10.1007/s00453-004-1120-3> | HEAD: HTTP 200 | 302 → http://link.springer.com/10.1007/s00453-004-1120-3<br>301 → https://link.springer.com/10.1007/s00453-004-1120-3<br>301 → https://link.springer.com/article/10.1007/s00453-004-1120-3<br>303 → https://idp.springer.com/authorize?response_type=cookie&client_id=springerlink&redirect_uri=https%3A%2F%2Flink.springer.com%2Farticle%2F10.1007%2Fs00453-004-1120-3<br>302 → https://idp.springer.com/transit?redirect_uri=https%3A%2F%2Flink.springer.com%2Farticle%2F10.1007%2Fs00453-004-1120-3&code=d0a713dd-974f-45df-b234-c562a13bb2fc<br>302 → https://link.springer.com/article/10.1007/s00453-004-1120-3?error=cookies_not_supported&code=d0a713dd-974f-45df-b234-c562a13bb2fc |
| ch10 | <https://doi.org/10.1007/s007999900023> | HEAD: HTTP 200 | 302 → http://link.springer.com/10.1007/s007999900023<br>301 → https://link.springer.com/10.1007/s007999900023<br>301 → https://link.springer.com/article/10.1007/s007999900023<br>303 → https://idp.springer.com/authorize?response_type=cookie&client_id=springerlink&redirect_uri=https%3A%2F%2Flink.springer.com%2Farticle%2F10.1007%2Fs007999900023<br>302 → https://idp.springer.com/transit?redirect_uri=https%3A%2F%2Flink.springer.com%2Farticle%2F10.1007%2Fs007999900023&code=975a09c3-3d7d-46d3-9226-65b8d9667931<br>302 → https://link.springer.com/article/10.1007/s007999900023?error=cookies_not_supported&code=975a09c3-3d7d-46d3-9226-65b8d9667931 |
| ch01 | <https://doi.org/10.1007/s10506-014-9158-2> | HEAD: HTTP 200 | 302 → http://link.springer.com/10.1007/s10506-014-9158-2<br>301 → https://link.springer.com/10.1007/s10506-014-9158-2<br>301 → https://link.springer.com/article/10.1007/s10506-014-9158-2<br>303 → https://idp.springer.com/authorize?response_type=cookie&client_id=springerlink&redirect_uri=https%3A%2F%2Flink.springer.com%2Farticle%2F10.1007%2Fs10506-014-9158-2<br>302 → https://idp.springer.com/transit?redirect_uri=https%3A%2F%2Flink.springer.com%2Farticle%2F10.1007%2Fs10506-014-9158-2&code=bb887bf9-ac2a-44bd-b217-3ac405e91239<br>302 → https://link.springer.com/article/10.1007/s10506-014-9158-2?error=cookies_not_supported&code=bb887bf9-ac2a-44bd-b217-3ac405e91239 |
| ch34 | <https://doi.org/10.1007/s11225-007-9067-1> | GET (HEAD 404 → ranged GET): HTTP 404 | — |
| ch06 | <https://doi.org/10.1007/s13171-019-00164-5> | HEAD: HTTP 200 | 302 → http://link.springer.com/10.1007/s13171-019-00164-5<br>301 → https://link.springer.com/10.1007/s13171-019-00164-5<br>301 → https://link.springer.com/article/10.1007/s13171-019-00164-5<br>303 → https://idp.springer.com/authorize?response_type=cookie&client_id=springerlink&redirect_uri=https%3A%2F%2Flink.springer.com%2Farticle%2F10.1007%2Fs13171-019-00164-5<br>302 → https://idp.springer.com/transit?redirect_uri=https%3A%2F%2Flink.springer.com%2Farticle%2F10.1007%2Fs13171-019-00164-5&code=012aa8bf-13ba-46ae-add5-8cf675782d30<br>302 → https://link.springer.com/article/10.1007/s13171-019-00164-5?error=cookies_not_supported&code=012aa8bf-13ba-46ae-add5-8cf675782d30 |
| ch05, ch33 | <https://doi.org/10.1007/s13389-012-0027-1> | HEAD: HTTP 200 | 302 → http://link.springer.com/10.1007/s13389-012-0027-1<br>301 → https://link.springer.com/10.1007/s13389-012-0027-1<br>301 → https://link.springer.com/article/10.1007/s13389-012-0027-1<br>303 → https://idp.springer.com/authorize?response_type=cookie&client_id=springerlink&redirect_uri=https%3A%2F%2Flink.springer.com%2Farticle%2F10.1007%2Fs13389-012-0027-1<br>302 → https://idp.springer.com/transit?redirect_uri=https%3A%2F%2Flink.springer.com%2Farticle%2F10.1007%2Fs13389-012-0027-1&code=7f91e9f7-769e-439b-bb23-98de5b73d051<br>302 → https://link.springer.com/article/10.1007/s13389-012-0027-1?error=cookies_not_supported&code=7f91e9f7-769e-439b-bb23-98de5b73d051 |
| ch23 | <https://doi.org/10.1007/s41468-020-00054-y> | HEAD: HTTP 200 | 302 → https://link.springer.com/10.1007/s41468-020-00054-y<br>301 → https://link.springer.com/article/10.1007/s41468-020-00054-y<br>303 → https://idp.springer.com/authorize?response_type=cookie&client_id=springerlink&redirect_uri=https%3A%2F%2Flink.springer.com%2Farticle%2F10.1007%2Fs41468-020-00054-y<br>302 → https://idp.springer.com/transit?redirect_uri=https%3A%2F%2Flink.springer.com%2Farticle%2F10.1007%2Fs41468-020-00054-y&code=d31ef8b2-a5e0-4d23-9fed-20ba4ffa2da0<br>302 → https://link.springer.com/article/10.1007/s41468-020-00054-y?error=cookies_not_supported&code=d31ef8b2-a5e0-4d23-9fed-20ba4ffa2da0 |
| ch06, ch16, ch26, ch38 | <https://doi.org/10.1016/0004-3702(79)90008-0> | HEAD: HTTP 200 | 302 → https://linkinghub.elsevier.com/retrieve/pii/0004370279900080 |
| ch01, ch03, ch04, ch06, ch07, ch16, ch17 | <https://doi.org/10.1016/0004-3702(82)90020-0> | HEAD: HTTP 200 | 302 → https://linkinghub.elsevier.com/retrieve/pii/0004370282900200 |
| ch01, ch04 | <https://doi.org/10.1016/0004-3702(82)90021-2> | HEAD: HTTP 200 | 302 → https://linkinghub.elsevier.com/retrieve/pii/0004370282900212 |
| ch20 | <https://doi.org/10.1016/0004-3702(83)90008-5> | HEAD: HTTP 200 | 302 → https://linkinghub.elsevier.com/retrieve/pii/0004370283900085 |
| ch06, ch16, ch24 | <https://doi.org/10.1016/0004-3702(86)90080-9> | HEAD: HTTP 200 | 302 → https://linkinghub.elsevier.com/retrieve/pii/0004370286900809 |
| ch24 | <https://doi.org/10.1016/0004-3702(87)90062-2> | HEAD: HTTP 200 | 302 → https://linkinghub.elsevier.com/retrieve/pii/0004370287900622 |
| ch24 | <https://doi.org/10.1016/0004-3702(87)90063-4> | HEAD: HTTP 200 | 302 → https://linkinghub.elsevier.com/retrieve/pii/0004370287900634 |
| ch24 | <https://doi.org/10.1016/0004-3702(89)90079-9> | HEAD: HTTP 200 | 302 → https://linkinghub.elsevier.com/retrieve/pii/0004370289900799 |
| ch01, ch04 | <https://doi.org/10.1016/0004-3702(93)90068-M> | HEAD: HTTP 200 | 302 → https://linkinghub.elsevier.com/retrieve/pii/000437029390068M |
| ch02, ch06, ch27, ch31 | <https://doi.org/10.1016/0004-3702(94)00041-X> | HEAD: HTTP 200 | 302 → https://linkinghub.elsevier.com/retrieve/pii/000437029400041X |
| ch11 | <https://doi.org/10.1016/0010-4485(85)90293-3> | HEAD: HTTP 200 | 302 → https://linkinghub.elsevier.com/retrieve/pii/0010448585902933 |
| ch04 | <https://doi.org/10.1016/0010-4809(75)90009-9> | HEAD: HTTP 200 | 302 → https://linkinghub.elsevier.com/retrieve/pii/0010480975900099 |
| ch04 | <https://doi.org/10.1016/0025-5564(75)90047-4> | HEAD: HTTP 200 | 302 → https://linkinghub.elsevier.com/retrieve/pii/0025556475900474 |
| ch04 | <https://doi.org/10.1016/0167-2789(90)90087-6> | HEAD: HTTP 200 | 302 → https://linkinghub.elsevier.com/retrieve/pii/0167278990900876 |
| ch10 | <https://doi.org/10.1016/0306-4379(96)00017-8> | HEAD: HTTP 200 | 302 → https://linkinghub.elsevier.com/retrieve/pii/0306437996000178 |
| ch04, ch06 | <https://doi.org/10.1016/0306-4573(88)90021-0> | HEAD: HTTP 200 | 302 → https://linkinghub.elsevier.com/retrieve/pii/0306457388900210 |
| ch06 | <https://doi.org/10.1016/0377-2217(90)90057-I> | HEAD: HTTP 200 | 302 → https://linkinghub.elsevier.com/retrieve/pii/037722179090057I |
| ch06, ch28 | <https://doi.org/10.1016/B978-0-12-417750-5.50022-1> | HEAD: HTTP 200 | 302 → https://linkinghub.elsevier.com/retrieve/pii/B9780124177505500221 |
| ch06 | <https://doi.org/10.1016/B978-0-273-08793-9.50010-8> | HEAD: HTTP 200 | 302 → https://linkinghub.elsevier.com/retrieve/pii/B9780273087939500108 |
| ch06 | <https://doi.org/10.1016/B978-0-934613-40-8.50006-3> | HEAD: HTTP 200 | 302 → https://linkinghub.elsevier.com/retrieve/pii/B9780934613408500063 |
| ch06 | <https://doi.org/10.1016/B978-1-55860-856-6.X5000-5> | HEAD: HTTP 200 | 302 → https://linkinghub.elsevier.com/retrieve/pii/B9781558608566X50005 |
| ch04, ch06, ch07 | <https://doi.org/10.1016/C2009-0-27609-4> | HEAD: HTTP 200 | 302 → https://linkinghub.elsevier.com/retrieve/pii/C20090276094 |
| ch04, ch06 | <https://doi.org/10.1016/S0019-9958(65)90241-X> | HEAD: HTTP 200 | 302 → https://linkinghub.elsevier.com/retrieve/pii/S001999586590241X |
| ch29 | <https://doi.org/10.1016/S0079-7421(08)60536-8> | HEAD: HTTP 200 | 302 → https://linkinghub.elsevier.com/retrieve/pii/S0079742108605368 |
| ch25 | <https://doi.org/10.1016/S0169-023X(97)00056-6> | HEAD: HTTP 200 | 302 → https://linkinghub.elsevier.com/retrieve/pii/S0169023X97000566 |
| ch09 | <https://doi.org/10.1016/S0378-8733(03)00009-1> | HEAD: HTTP 200 | 302 → https://linkinghub.elsevier.com/retrieve/pii/S0378873303000091 |
| ch11 | <https://doi.org/10.1016/S1042-8143(05)80010-X> | HEAD: HTTP 200 | 302 → https://linkinghub.elsevier.com/retrieve/pii/S104281430580010X |
| ch07 | <https://doi.org/10.1016/j.artint.2004.11.004> | HEAD: HTTP 200 | 302 → https://linkinghub.elsevier.com/retrieve/pii/S0004370204001778 |
| ch04 | <https://doi.org/10.1016/j.artint.2012.10.002> | HEAD: HTTP 200 | 302 → https://linkinghub.elsevier.com/retrieve/pii/S0004370212001245 |
| ch20 | <https://doi.org/10.1016/j.artint.2018.07.007> | HEAD: HTTP 200 | 302 → https://linkinghub.elsevier.com/retrieve/pii/S0004370218305988 |
| ch14 | <https://doi.org/10.1016/j.jss.2016.02.047> | HEAD: HTTP 200 | 302 → https://linkinghub.elsevier.com/retrieve/pii/S0164121216000789 |
| ch01, ch37 | <https://doi.org/10.1016/j.neunet.2014.09.003> | HEAD: HTTP 200 | 302 → https://linkinghub.elsevier.com/retrieve/pii/S0893608014002135 |
| ch15 | <https://doi.org/10.1016/j.simpa.2023.100556> | HEAD: HTTP 200 | 302 → https://linkinghub.elsevier.com/retrieve/pii/S2665963823000933 |
| ch26 | <https://doi.org/10.1017/9781108571401> | HEAD: HTTP 200 | 302 → https://www.cambridge.org/core/product/identifier/9781108571401/type/book<br>301 → https://www.cambridge.org/core/books/bandit-algorithms/8E39FD004E6CE036680F90DD0C6F09FC |
| ch07, ch24 | <https://doi.org/10.1017/CBO9780511803161> | HEAD: HTTP 200 | 302 → https://www.cambridge.org/core/product/identifier/9780511803161/type/book<br>301 → https://www.cambridge.org/core/books/causality/B0046844FAE10CBF274D4ACBDAEB5F5B |
| ch05 | <https://doi.org/10.1017/CBO9780511803932> | HEAD: HTTP 200 | 302 → https://www.cambridge.org/core/product/identifier/9780511803932/type/book<br>301 → https://www.cambridge.org/core/books/communities-of-practice/724C22A03B12D11DFC345EEF0AD3F22A<br>302 → https://www.cambridge.org/highereducation/product/724C22A03B12D11DFC345EEF0AD3F22A<br>302 → https://www.cambridge.org/highereducation/books/communities-of-practice/724C22A03B12D11DFC345EEF0AD3F22A |
| ch11 | <https://doi.org/10.1017/CBO9780511809477> | HEAD: HTTP 200 | 302 → https://www.cambridge.org/core/product/identifier/9780511809477/type/book<br>301 → https://www.cambridge.org/core/books/judgment-under-uncertainty/6F9E814794E08EC43D426E480A4B412C |
| ch02 | <https://doi.org/10.1017/CBO9780511840005> | HEAD: HTTP 200 | 302 → https://www.cambridge.org/core/product/identifier/9780511840005/type/book<br>301 → https://www.cambridge.org/core/books/uses-of-argument/26CF801BC12004587B66778297D5567C |
| ch34 | <https://doi.org/10.1017/CBO9780511974533> | HEAD: HTTP 200 | 302 → https://www.cambridge.org/core/product/identifier/9780511974533/type/book<br>301 → https://www.cambridge.org/core/books/logical-dynamics-of-information-and-interaction/C6007450DC8C99C2CD35BA50FE3ACA5F |
| ch02, ch19 | <https://doi.org/10.1017/CBO9781139173438> | HEAD: HTTP 200 | 302 → https://www.cambridge.org/core/product/identifier/9781139173438/type/book<br>301 → https://www.cambridge.org/core/books/speech-acts/D2D7B03E472C8A390ED60B86E08640E7 |
| ch21 | <https://doi.org/10.1017/CBO9781139583923> | HEAD: HTTP 200 | 302 → https://www.cambridge.org/core/product/identifier/9781139583923/type/book<br>301 → https://www.cambridge.org/core/books/automated-planning-and-acting/E6DE5715A2190651352DFB0869916BC3 |
| ch11 | <https://doi.org/10.1017/S0269888900000060> | HEAD: HTTP 200 | 302 → https://www.cambridge.org/core/product/identifier/S0269888900000060/type/journal_article<br>301 → https://www.cambridge.org/core/journals/knowledge-engineering-review/article/knowledge-acquisition-tools-based-on-personal-construct-psychology/826CBE6EE7140EAB7261A175D73D0103<br>302 → https://www.cambridge.org/core/journals/knowledge-engineering-review/article/abs/knowledge-acquisition-tools-based-on-personal-construct-psychology/826CBE6EE7140EAB7261A175D73D0103 |
| ch10 | <https://doi.org/10.1017/S1351324900000048> | HEAD: HTTP 200 | 302 → https://www.cambridge.org/core/product/identifier/S1351324900000048/type/journal_article<br>301 → https://www.cambridge.org/core/journals/natural-language-engineering/article/technical-terminology-some-linguistic-properties-and-an-algorithm-for-identification-in-text/D5F076938C4E3F24B11EDC2E831216AF<br>302 → https://www.cambridge.org/core/journals/natural-language-engineering/article/abs/technical-terminology-some-linguistic-properties-and-an-algorithm-for-identification-in-text/D5F076938C4E3F24B11EDC2E831216AF |
| ch31, ch33 | <https://doi.org/10.1017/S1471068418000054> | HEAD: HTTP 200 | 302 → https://www.cambridge.org/core/product/identifier/S1471068418000054/type/journal_article<br>301 → https://www.cambridge.org/core/journals/theory-and-practice-of-logic-programming/article/multishot-asp-solving-with-clingo/FAED3429900D84CDD5155326A36548F2<br>302 → https://www.cambridge.org/core/journals/theory-and-practice-of-logic-programming/article/abs/multishot-asp-solving-with-clingo/FAED3429900D84CDD5155326A36548F2 |
| ch24 | <https://doi.org/10.1023/A:1009748302351> | HEAD: HTTP 200 | 302 → https://link.springer.com/10.1023/A:1009748302351<br>301 → https://link.springer.com/article/10.1023/A:1009748302351<br>303 → https://idp.springer.com/authorize?response_type=cookie&client_id=springerlink&redirect_uri=https%3A%2F%2Flink.springer.com%2Farticle%2F10.1023%2FA%3A1009748302351<br>302 → https://idp.springer.com/transit?redirect_uri=https%3A%2F%2Flink.springer.com%2Farticle%2F10.1023%2FA%3A1009748302351&code=cefce55e-5b08-4af7-a804-22f0a0739820<br>302 → https://link.springer.com/article/10.1023/A:1009748302351?error=cookies_not_supported&code=cefce55e-5b08-4af7-a804-22f0a0739820 |
| ch29 | <https://doi.org/10.1037/0003-066X.58.9.697> | HEAD: HTTP 200 | 302 → https://doi.apa.org/doi/10.1037/0003-066X.58.9.697<br>301 → https://psycnet.apa.org:443/doi/10.1037/0003-066X.58.9.697<br>302 → https://psycnet.apa.org:443/doiLanding?doi=10.1037/0003-066X.58.9.697 |
| ch06, ch31 | <https://doi.org/10.1037/0033-295X.84.4.327> | HEAD: HTTP 200 | 302 → https://doi.apa.org/doi/10.1037/0033-295X.84.4.327<br>301 → https://psycnet.apa.org:443/doi/10.1037/0033-295X.84.4.327<br>302 → https://psycnet.apa.org:443/doiLanding?doi=10.1037/0033-295X.84.4.327 |
| ch11 | <https://doi.org/10.1037/0096-1523.1.3.288> | HEAD: HTTP 200 | 302 → https://doi.apa.org/doi/10.1037/0096-1523.1.3.288<br>301 → https://psycnet.apa.org:443/doi/10.1037/0096-1523.1.3.288<br>302 → https://psycnet.apa.org:443/doiLanding?doi=10.1037/0096-1523.1.3.288 |
| ch22, ch26 | <https://doi.org/10.1038/nature08227> | HEAD: HTTP 200 | 302 → https://www.nature.com/articles/nature08227<br>303 → https://idp.nature.com/authorize?response_type=cookie&client_id=grover&redirect_uri=https%3A%2F%2Fwww.nature.com%2Farticles%2Fnature08227<br>302 → https://idp.nature.com/transit?redirect_uri=https%3A%2F%2Fwww.nature.com%2Farticles%2Fnature08227&code=ea284214-16b1-4507-954d-3efbf7c39eed<br>302 → https://www.nature.com/articles/nature08227?error=cookies_not_supported&code=ea284214-16b1-4507-954d-3efbf7c39eed |
| ch22 | <https://doi.org/10.1038/nrn2787> | HEAD: HTTP 200 | 302 → https://www.nature.com/articles/nrn2787<br>303 → https://idp.nature.com/authorize?response_type=cookie&client_id=grover&redirect_uri=https%3A%2F%2Fwww.nature.com%2Farticles%2Fnrn2787<br>302 → https://idp.nature.com/transit?redirect_uri=https%3A%2F%2Fwww.nature.com%2Farticles%2Fnrn2787&code=03884420-cbe5-42cb-afc8-ae5a974efa6b<br>302 → https://www.nature.com/articles/nrn2787?error=cookies_not_supported&code=03884420-cbe5-42cb-afc8-ae5a974efa6b |
| ch19, ch38 | <https://doi.org/10.1038/s41586-024-07566-y> | HEAD: HTTP 200 | 302 → https://www.nature.com/articles/s41586-024-07566-y<br>303 → https://idp.nature.com/authorize?response_type=cookie&client_id=grover&redirect_uri=https%3A%2F%2Fwww.nature.com%2Farticles%2Fs41586-024-07566-y<br>302 → https://idp.nature.com/transit?redirect_uri=https%3A%2F%2Fwww.nature.com%2Farticles%2Fs41586-024-07566-y&code=8ebf09dc-8023-4706-af00-ae84384bb393<br>302 → https://www.nature.com/articles/s41586-024-07566-y?error=cookies_not_supported&code=8ebf09dc-8023-4706-af00-ae84384bb393 |
| ch38 | <https://doi.org/10.1038/s41586-026-10549-w> | HEAD: HTTP 200 | 302 → https://www.nature.com/articles/s41586-026-10549-w<br>303 → https://idp.nature.com/authorize?response_type=cookie&client_id=grover&redirect_uri=https%3A%2F%2Fwww.nature.com%2Farticles%2Fs41586-026-10549-w<br>302 → https://idp.nature.com/transit?redirect_uri=https%3A%2F%2Fwww.nature.com%2Farticles%2Fs41586-026-10549-w&code=dd80c0df-c765-4ddd-af98-a1048004978c<br>302 → https://www.nature.com/articles/s41586-026-10549-w?error=cookies_not_supported&code=dd80c0df-c765-4ddd-af98-a1048004978c |
| ch06 | <https://doi.org/10.1051/ro/196802v100571> | GET (HEAD 403 → ranged GET): HTTP 403 | 302 → http://www.rairo-ro.org/10.1051/ro/196802V100571<br>301 → https://www.rairo-ro.org/10.1051/ro/196802V100571 |
| ch22 | <https://doi.org/10.1057/jors.1984.2> | GET (HEAD 403 → ranged GET): HTTP 403 | 302 → https://www.tandfonline.com/doi/full/10.1057/jors.1984.2 |
| ch26 | <https://doi.org/10.1073/pnas.1611835114> | GET (HEAD 403 → ranged GET): HTTP 403 | 302 → https://pnas.org/doi/full/10.1073/pnas.1611835114 |
| ch36 | <https://doi.org/10.1073/pnas.1708274114> | GET (HEAD 403 → ranged GET): HTTP 403 | 302 → https://pnas.org/doi/full/10.1073/pnas.1708274114 |
| ch06 | <https://doi.org/10.1073/pnas.36.1.48> | GET (HEAD 403 → ranged GET): HTTP 403 | 302 → https://pnas.org/doi/full/10.1073/pnas.36.1.48 |
| ch22 | <https://doi.org/10.1080/00207727008920220> | GET (HEAD 403 → ranged GET): HTTP 403 | 302 → http://www.tandfonline.com/doi/abs/10.1080/00207727008920220<br>301 → https://www.tandfonline.com/doi/abs/10.1080/00207727008920220 |
| ch10 | <https://doi.org/10.1080/01621459.1952.10483446> | GET (HEAD 403 → ranged GET): HTTP 403 | 302 → https://www.tandfonline.com/doi/full/10.1080/01621459.1952.10483446 |
| ch37 | <https://doi.org/10.1080/13518040490450529> | GET (HEAD 403 → ranged GET): HTTP 403 | 302 → http://www.tandfonline.com/doi/abs/10.1080/13518040490450529<br>301 → https://www.tandfonline.com/doi/abs/10.1080/13518040490450529 |
| ch29, ch31 | <https://doi.org/10.1080/19462166.2013.869766> | GET (HEAD 403 → ranged GET): HTTP 403 | 302 → https://journals.sagepub.com/doi/10.1080/19462166.2013.869766 |
| ch02 | <https://doi.org/10.1093/acprof:oso/9780199232383.001.0001> | GET (HEAD 403 → ranged GET): HTTP 403 | 302 → https://academic.oup.com/book/32518 |
| ch02 | <https://doi.org/10.1093/analys/23.6.121> | GET (HEAD 403 → ranged GET): HTTP 403 | 302 → https://academic.oup.com/analysis/article-lookup/doi/10.1093/analys/23.6.121 |
| ch26, ch33 | <https://doi.org/10.1093/biomet/41.1-2.100> | GET (HEAD 403 → ranged GET): HTTP 403 | 302 → https://academic.oup.com/biomet/article-lookup/doi/10.1093/biomet/41.1-2.100 |
| ch23 | <https://doi.org/10.1093/comjnl/25.4.465> | GET (HEAD 403 → ranged GET): HTTP 403 | 302 → https://academic.oup.com/comjnl/article-lookup/doi/10.1093/comjnl/25.4.465 |
| ch34 | <https://doi.org/10.1093/logcom/2.6.719> | GET (HEAD 403 → ranged GET): HTTP 403 | 302 → https://academic.oup.com/logcom/article-lookup/doi/10.1093/logcom/2.6.719 |
| ch06, ch28 | <https://doi.org/10.1093/mind/LX.237.1> | GET (HEAD 403 → ranged GET): HTTP 403 | 302 → https://academic.oup.com/mind/article-lookup/doi/10.1093/mind/LX.237.1 |
| ch05 | <https://doi.org/10.1093/oso/9780195092691.001.0001> | GET (HEAD 403 → ranged GET): HTTP 403 | 302 → https://academic.oup.com/book/52097 |
| ch04 | <https://doi.org/10.1098/rstl.1763.0053> | GET (HEAD 403 → ranged GET): HTTP 403 | 302 → https://royalsocietypublishing.org/rstl/article/doi/10.1098/rstl.1763.0053/119736/LII-An-essay-towards-solving-a-problem-in-the |
| ch02 | <https://doi.org/10.1109/2.161279> | HEAD: HTTP 202 | 302 → http://ieeexplore.ieee.org/document/161279/<br>301 → https://ieeexplore.ieee.org/document/161279/ |
| ch16 | <https://doi.org/10.1109/2.485845> | HEAD: HTTP 202 | 302 → http://ieeexplore.ieee.org/document/485845/<br>301 → https://ieeexplore.ieee.org/document/485845/ |
| ch22 | <https://doi.org/10.1109/2.84874> | HEAD: HTTP 202 | 302 → http://ieeexplore.ieee.org/document/84874/<br>301 → https://ieeexplore.ieee.org/document/84874/ |
| ch11 | <https://doi.org/10.1109/21.31053> | HEAD: HTTP 202 | 302 → http://ieeexplore.ieee.org/document/31053/<br>301 → https://ieeexplore.ieee.org/document/31053/ |
| ch29 | <https://doi.org/10.1109/32.210303> | HEAD: HTTP 202 | 302 → http://ieeexplore.ieee.org/document/210303/<br>301 → https://ieeexplore.ieee.org/document/210303/ |
| ch21 | <https://doi.org/10.1109/3468.844354> | HEAD: HTTP 202 | 302 → http://ieeexplore.ieee.org/document/844354/<br>301 → https://ieeexplore.ieee.org/document/844354/ |
| ch06, ch24 | <https://doi.org/10.1109/5.18626> | HEAD: HTTP 202 | 302 → http://ieeexplore.ieee.org/document/18626/<br>301 → https://ieeexplore.ieee.org/document/18626/ |
| ch22 | <https://doi.org/10.1109/52.469757> | HEAD: HTTP 202 | 302 → http://ieeexplore.ieee.org/document/469757/<br>301 → https://ieeexplore.ieee.org/document/469757/ |
| ch28 | <https://doi.org/10.1109/ms.2001.936213> | HEAD: HTTP 202 | 302 → http://ieeexplore.ieee.org/document/936213/<br>301 → https://ieeexplore.ieee.org/document/936213/<br>final: https://ieeexplore.ieee.org/document/936213/ |
| ch06, ch16, ch17 | <https://doi.org/10.1109/69.43410> | HEAD: HTTP 202 | 302 → http://ieeexplore.ieee.org/document/43410/<br>301 → https://ieeexplore.ieee.org/document/43410/ |
| ch32 | <https://doi.org/10.1109/90.663936> | HEAD: HTTP 202 | 302 → http://ieeexplore.ieee.org/document/663936/<br>301 → https://ieeexplore.ieee.org/document/663936/ |
| ch37 | <https://doi.org/10.1109/ACC.1997.609105> | HEAD: HTTP 202 | 302 → http://ieeexplore.ieee.org/document/609105/<br>301 → https://ieeexplore.ieee.org/document/609105/ |
| ch30 | <https://doi.org/10.1109/ACCESS.2021.3100711> | HEAD: HTTP 202 | 302 → https://ieeexplore.ieee.org/document/9499031/ |
| ch30 | <https://doi.org/10.1109/ACCESS.2023.3284751> | HEAD: HTTP 202 | 302 → https://ieeexplore.ieee.org/document/10147201/ |
| ch23 | <https://doi.org/10.1109/C-M.1978.218136> | HEAD: HTTP 202 | 302 → http://ieeexplore.ieee.org/document/1646911/<br>301 → https://ieeexplore.ieee.org/document/1646911/ |
| ch12 | <https://doi.org/10.1109/CVPR.2018.00286> | HEAD: HTTP 202 | 302 → https://ieeexplore.ieee.org/document/8578384/ |
| ch18 | <https://doi.org/10.1109/DCS.1988.12507> | HEAD: HTTP 202 | 302 → http://ieeexplore.ieee.org/document/12507/<br>301 → https://ieeexplore.ieee.org/document/12507/ |
| ch18 | <https://doi.org/10.1109/GRID.2004.14> | HEAD: HTTP 202 | 302 → http://ieeexplore.ieee.org/document/1382809/<br>301 → https://ieeexplore.ieee.org/document/1382809/ |
| ch24 | <https://doi.org/10.1109/ICDM.2008.17> | HEAD: HTTP 202 | 302 → http://ieeexplore.ieee.org/document/4781136/<br>301 → https://ieeexplore.ieee.org/document/4781136/ |
| ch27 | <https://doi.org/10.1109/ICSE.2017.9> | HEAD: HTTP 202 | 302 → http://ieeexplore.ieee.org/document/7985645/<br>301 → https://ieeexplore.ieee.org/document/7985645/ |
| ch24 | <https://doi.org/10.1109/ICWS.2017.13> | HEAD: HTTP 202 | 302 → http://ieeexplore.ieee.org/document/8029742/<br>301 → https://ieeexplore.ieee.org/document/8029742/ |
| ch09 | <https://doi.org/10.1109/IEEESTD.2017.8055462> | HEAD: HTTP 202 | 302 → https://ieeexplore.ieee.org/document/8055462/ |
| ch33 | <https://doi.org/10.1109/IEEESTD.2020.9120376> | HEAD: HTTP 202 | 302 → https://ieeexplore.ieee.org/document/9120376/ |
| ch09 | <https://doi.org/10.1109/IEEESTD.2023.10054520> | HEAD: HTTP 202 | 302 → https://ieeexplore.ieee.org/document/10054520/ |
| ch18 | <https://doi.org/10.1109/ISCA45697.2020.00045> | HEAD: HTTP 202 | 302 → https://ieeexplore.ieee.org/document/9138989/ |
| ch04 | <https://doi.org/10.1109/JPROC.1998.658762> | HEAD: HTTP 202 | 302 → http://ieeexplore.ieee.org/document/658762/<br>301 → https://ieeexplore.ieee.org/document/658762/ |
| ch37 | <https://doi.org/10.1109/JRPROC.1949.232969> | HEAD: HTTP 202 | 302 → http://ieeexplore.ieee.org/document/1697831/<br>301 → https://ieeexplore.ieee.org/document/1697831/ |
| ch18 | <https://doi.org/10.1109/MICRO.2016.7783759> | HEAD: HTTP 202 | 302 → http://ieeexplore.ieee.org/document/7783759/<br>301 → https://ieeexplore.ieee.org/document/7783759/ |
| ch23 | <https://doi.org/10.1109/MS.1984.233702> | HEAD: HTTP 202 | 302 → http://ieeexplore.ieee.org/document/1695100/<br>301 → https://ieeexplore.ieee.org/document/1695100/ |
| ch14, ch15 | <https://doi.org/10.1109/RE.2009.9> | HEAD: HTTP 202 | 302 → http://ieeexplore.ieee.org/document/5328509/<br>301 → https://ieeexplore.ieee.org/document/5328509/ |
| ch10, ch32 | <https://doi.org/10.1109/SEQUEN.1997.666900> | HEAD: HTTP 202 | 302 → http://ieeexplore.ieee.org/document/666900/<br>301 → https://ieeexplore.ieee.org/document/666900/ |
| ch14 | <https://doi.org/10.1109/SFCS.1977.32> | HEAD: HTTP 202 | 302 → http://ieeexplore.ieee.org/document/4567924/<br>301 → https://ieeexplore.ieee.org/document/4567924/ |
| ch38 | <https://doi.org/10.1109/SP40001.2021.00019> | HEAD: HTTP 202 | 302 → https://ieeexplore.ieee.org/document/9519428/ |
| ch13 | <https://doi.org/10.1109/TIT.1956.1056813> | HEAD: HTTP 202 | 302 → http://ieeexplore.ieee.org/document/1056813/<br>301 → https://ieeexplore.ieee.org/document/1056813/ |
| ch13 | <https://doi.org/10.1109/TPAMI.2010.57> | HEAD: HTTP 202 | 302 → http://ieeexplore.ieee.org/document/5432202/<br>301 → https://ieeexplore.ieee.org/document/5432202/ |
| ch07, ch13, ch15, ch16, ch18, ch32 | <https://doi.org/10.1109/TPAMI.2018.2889473> | HEAD: HTTP 202 | 302 → https://ieeexplore.ieee.org/document/8594636/ |
| ch26 | <https://doi.org/10.1109/TPAMI.2021.3057446> | HEAD: HTTP 202 | 302 → https://ieeexplore.ieee.org/document/9349197/ |
| ch27 | <https://doi.org/10.1109/TR.2014.2335995> | HEAD: HTTP 202 | 302 → http://ieeexplore.ieee.org/document/6866924/<br>301 → https://ieeexplore.ieee.org/document/6866924/ |
| ch23 | <https://doi.org/10.1109/TSE.2010.62> | HEAD: HTTP 202 | 302 → http://ieeexplore.ieee.org/document/5487526/<br>301 → https://ieeexplore.ieee.org/document/5487526/ |
| ch01, ch37 | <https://doi.org/10.1109/TSMC.1971.4308320> | HEAD: HTTP 202 | 302 → http://ieeexplore.ieee.org/document/4308320/<br>301 → https://ieeexplore.ieee.org/document/4308320/ |
| ch11 | <https://doi.org/10.1109/TSSC.1966.300074> | HEAD: HTTP 202 | 302 → http://ieeexplore.ieee.org/document/4082064/<br>301 → https://ieeexplore.ieee.org/document/4082064/ |
| ch37 | <https://doi.org/10.1111/0162-895X.00317> | GET (HEAD 403 → ranged GET): HTTP 403 | 302 → https://onlinelibrary.wiley.com/doi/10.1111/0162-895X.00317 |
| ch11 | <https://doi.org/10.1111/1468-0394.00045> | GET (HEAD 403 → ranged GET): HTTP 403 | 302 → https://onlinelibrary.wiley.com/doi/10.1111/1468-0394.00045 |
| ch28 | <https://doi.org/10.1111/j.1469-8137.1912.tb05611.x> | GET (HEAD 403 → ranged GET): HTTP 403 | 302 → https://nph.onlinelibrary.wiley.com/doi/10.1111/j.1469-8137.1912.tb05611.x |
| ch06 | <https://doi.org/10.1111/j.2517-6161.1951.tb00088.x> | GET (HEAD 403 → ranged GET): HTTP 403 | 302 → https://academic.oup.com/jrsssb/article/13/2/238/7026675 |
| ch06 | <https://doi.org/10.1112/plms/s2-42.1.230> | GET (HEAD 403 → ranged GET): HTTP 403 | 302 → http://doi.wiley.com/10.1112/plms/s2-42.1.230<br>302 → https://onlinelibrary.wiley.com/doi/10.1112/plms/s2-42.1.230 |
| ch06, ch22 | <https://doi.org/10.1115/1.3662552> | GET (HEAD 403 → ranged GET): HTTP 403 | 302 → https://asmedigitalcollection.asme.org/fluidsengineering/article/82/1/35/397706/A-New-Approach-to-Linear-Filtering-and-Prediction |
| ch04 | <https://doi.org/10.1126/science.217.4563.927> | GET (HEAD 403 → ranged GET): HTTP 403 | 302 → https://www.science.org/doi/10.1126/science.217.4563.927 |
| ch29 | <https://doi.org/10.1126/science.ade9097> | GET (HEAD 403 → ranged GET): HTTP 403 | 302 → https://www.science.org/doi/10.1126/science.ade9097 |
| ch31 | <https://doi.org/10.1137/0201010> | GET (HEAD 403 → ranged GET): HTTP 403 | 302 → https://epubs.siam.org/doi/10.1137/0201010 |
| ch10 | <https://doi.org/10.1137/0218012> | GET (HEAD 403 → ranged GET): HTTP 403 | 302 → https://epubs.siam.org/doi/10.1137/0218012 |
| ch32 | <https://doi.org/10.1137/090752286> | GET (HEAD 403 → ranged GET): HTTP 403 | 302 → https://epubs.siam.org/doi/10.1137/090752286 |
| ch32 | <https://doi.org/10.1137/1.9780898718003> | GET (HEAD 403 → ranged GET): HTTP 403 | 302 → https://epubs.siam.org/doi/book/10.1137/1.9780898718003 |
| ch26 | <https://doi.org/10.1137/1.9781611972771.42> | GET (HEAD 403 → ranged GET): HTTP 403 | 302 → https://epubs.siam.org/doi/10.1137/1.9781611972771.42 |
| ch07 | <https://doi.org/10.1137/S1064827595287997> | GET (HEAD 403 → ranged GET): HTTP 403 | 302 → https://epubs.siam.org/doi/10.1137/S1064827595287997 |
| ch18 | <https://doi.org/10.1145/103162.103163> | GET (HEAD 403 → ranged GET): HTTP 403 | 302 → https://dl.acm.org/doi/10.1145/103162.103163 |
| ch07 | <https://doi.org/10.1145/129888.129894> | GET (HEAD 403 → ranged GET): HTTP 403 | 302 → https://dl.acm.org/doi/10.1145/129888.129894 |
| ch18 | <https://doi.org/10.1145/1347375.1347389> | GET (HEAD 403 → ranged GET): HTTP 403 | 302 → https://dl.acm.org/doi/10.1145/1347375.1347389 |
| ch32 | <https://doi.org/10.1145/143062.143090> | GET (HEAD 403 → ranged GET): HTTP 403 | 302 → https://dl.acm.org/doi/10.1145/143062.143090 |
| ch18 | <https://doi.org/10.1145/1498765.1498785> | GET (HEAD 403 → ranged GET): HTTP 403 | 302 → https://dl.acm.org/doi/10.1145/1498765.1498785 |
| ch17 | <https://doi.org/10.1145/1555349.1555372> | GET (HEAD 403 → ranged GET): HTTP 403 | 302 → https://dl.acm.org/doi/10.1145/1555349.1555372 |
| ch06, ch07, ch08, ch12, ch13, ch15 | <https://doi.org/10.1145/1571941.1572114> | GET (HEAD 403 → ranged GET): HTTP 403 | 302 → https://dl.acm.org/doi/10.1145/1571941.1572114 |
| ch22 | <https://doi.org/10.1145/1866307.1866315> | GET (HEAD 403 → ranged GET): HTTP 403 | 302 → https://dl.acm.org/doi/10.1145/1866307.1866315 |
| ch33 | <https://doi.org/10.1145/191246.191322> | GET (HEAD 403 → ranged GET): HTTP 403 | 302 → http://portal.acm.org/citation.cfm?doid=191246.191322<br>301 → https://portal.acm.org/citation.cfm?doid=191246.191322 |
| ch24 | <https://doi.org/10.1145/203330.203341> | GET (HEAD 403 → ranged GET): HTTP 403 | 302 → https://dl.acm.org/doi/10.1145/203330.203341 |
| ch12 | <https://doi.org/10.1145/219717.219748> | GET (HEAD 403 → ranged GET): HTTP 403 | 302 → https://dl.acm.org/doi/10.1145/219717.219748 |
| ch18 | <https://doi.org/10.1145/2342509.2342513> | GET (HEAD 403 → ranged GET): HTTP 403 | 302 → https://dl.acm.org/doi/10.1145/2342509.2342513 |
| ch32 | <https://doi.org/10.1145/235968.233333> | GET (HEAD 403 → ranged GET): HTTP 403 | 302 → https://dl.acm.org/doi/10.1145/235968.233333 |
| ch31, ch34 | <https://doi.org/10.1145/2488388.2488425> | GET (HEAD 403 → ranged GET): HTTP 403 | 302 → https://dl.acm.org/doi/10.1145/2488388.2488425 |
| ch26 | <https://doi.org/10.1145/2523813> | GET (HEAD 403 → ranged GET): HTTP 403 | 302 → https://dl.acm.org/doi/10.1145/2523813 |
| ch30 | <https://doi.org/10.1145/2556938> | GET (HEAD 403 → ranged GET): HTTP 403 | 302 → https://dl.acm.org/doi/10.1145/2556938 |
| ch32 | <https://doi.org/10.1145/258533.258660> | GET (HEAD 403 → ranged GET): HTTP 403 | 302 → http://portal.acm.org/citation.cfm?doid=258533.258660<br>301 → https://portal.acm.org/citation.cfm?doid=258533.258660 |
| ch09 | <https://doi.org/10.1145/2593882.2593891> | GET (HEAD 403 → ranged GET): HTTP 403 | 302 → https://dl.acm.org/doi/10.1145/2593882.2593891 |
| ch09 | <https://doi.org/10.1145/2623330.2623623> | GET (HEAD 403 → ranged GET): HTTP 403 | 302 → https://dl.acm.org/doi/10.1145/2623330.2623623 |
| ch12 | <https://doi.org/10.1145/2629489> | GET (HEAD 403 → ranged GET): HTTP 403 | 302 → https://dl.acm.org/doi/10.1145/2629489 |
| ch32 | <https://doi.org/10.1145/263699.263712> | GET (HEAD 403 → ranged GET): HTTP 403 | 302 → http://portal.acm.org/citation.cfm?doid=263699.263712<br>301 → https://portal.acm.org/citation.cfm?doid=263699.263712 |
| ch18 | <https://doi.org/10.1145/2749469.2750386> | GET (HEAD 403 → ranged GET): HTTP 403 | 302 → https://dl.acm.org/doi/10.1145/2749469.2750386 |
| ch06 | <https://doi.org/10.1145/2939672.2939778> | GET (HEAD 403 → ranged GET): HTTP 403 | 302 → https://dl.acm.org/doi/10.1145/2939672.2939778 |
| ch34 | <https://doi.org/10.1145/302979.303030> | GET (HEAD 403 → ranged GET): HTTP 403 | 302 → http://portal.acm.org/citation.cfm?doid=302979.303030<br>301 → https://portal.acm.org/citation.cfm?doid=302979.303030 |
| ch23 | <https://doi.org/10.1145/3143561> | GET (HEAD 403 → ranged GET): HTTP 403 | 302 → https://dl.acm.org/doi/10.1145/3143561 |
| ch02 | <https://doi.org/10.1145/317087.317089> | GET (HEAD 403 → ranged GET): HTTP 403 | 302 → https://dl.acm.org/doi/10.1145/317087.317089 |
| ch02 | <https://doi.org/10.1145/318898.318921> | GET (HEAD 403 → ranged GET): HTTP 403 | 302 → https://dl.acm.org/doi/10.1145/318898.318921 |
| ch01, ch04 | <https://doi.org/10.1145/321250.321253> | GET (HEAD 403 → ranged GET): HTTP 403 | 302 → https://dl.acm.org/doi/10.1145/321250.321253 |
| ch06 | <https://doi.org/10.1145/321978.321991> | GET (HEAD 403 → ranged GET): HTTP 403 | 302 → https://dl.acm.org/doi/10.1145/321978.321991 |
| ch20 | <https://doi.org/10.1145/3287560.3287566> | GET (HEAD 403 → ranged GET): HTTP 403 | 302 → https://dl.acm.org/doi/10.1145/3287560.3287566 |
| ch03, ch05 | <https://doi.org/10.1145/3287560.3287596> | GET (HEAD 403 → ranged GET): HTTP 403 | 302 → https://dl.acm.org/doi/10.1145/3287560.3287596 |
| ch37, ch38 | <https://doi.org/10.1145/3314221.3314642> | GET (HEAD 403 → ranged GET): HTTP 403 | 302 → https://dl.acm.org/doi/10.1145/3314221.3314642 |
| ch07 | <https://doi.org/10.1145/331499.331504> | GET (HEAD 403 → ranged GET): HTTP 403 | 302 → https://dl.acm.org/doi/10.1145/331499.331504 |
| ch36 | <https://doi.org/10.1145/3373376.3378519> | GET (HEAD 403 → ranged GET): HTTP 403 | 302 → https://dl.acm.org/doi/10.1145/3373376.3378519 |
| ch13 | <https://doi.org/10.1145/3404835.3463098> | GET (HEAD 403 → ranged GET): HTTP 403 | 302 → https://dl.acm.org/doi/10.1145/3404835.3463098 |
| ch38 | <https://doi.org/10.1145/3442188.3445922> | GET (HEAD 403 → ranged GET): HTTP 403 | 302 → https://dl.acm.org/doi/10.1145/3442188.3445922 |
| ch05, ch25 | <https://doi.org/10.1145/3458723> | GET (HEAD 403 → ranged GET): HTTP 403 | 302 → https://dl.acm.org/doi/10.1145/3458723 |
| ch31 | <https://doi.org/10.1145/3473582> | GET (HEAD 403 → ranged GET): HTTP 403 | 302 → https://dl.acm.org/doi/10.1145/3473582 |
| ch23 | <https://doi.org/10.1145/351240.351266> | GET (HEAD 403 → ranged GET): HTTP 403 | 302 → https://dl.acm.org/doi/10.1145/351240.351266 |
| ch03, ch11, ch38 | <https://doi.org/10.1145/3571730> | GET (HEAD 403 → ranged GET): HTTP 403 | 302 → https://dl.acm.org/doi/10.1145/3571730 |
| ch05 | <https://doi.org/10.1145/359168.359176> | GET (HEAD 403 → ranged GET): HTTP 403 | 302 → https://dl.acm.org/doi/10.1145/359168.359176 |
| ch22, ch33 | <https://doi.org/10.1145/359545.359563> | GET (HEAD 403 → ranged GET): HTTP 403 | 302 → https://dl.acm.org/doi/10.1145/359545.359563 |
| ch30 | <https://doi.org/10.1145/3600006> | GET (HEAD 403 → ranged GET): HTTP 403 | 302 → https://dl.acm.org/doi/proceedings/10.1145/3600006 |
| ch12, ch17 | <https://doi.org/10.1145/3600006.3613165> | GET (HEAD 403 → ranged GET): HTTP 403 | 302 → https://dl.acm.org/doi/10.1145/3600006.3613165 |
| ch02, ch07, ch33 | <https://doi.org/10.1145/360051.360056> | GET (HEAD 403 → ranged GET): HTTP 403 | 302 → https://dl.acm.org/doi/10.1145/360051.360056 |
| ch12, ch21 | <https://doi.org/10.1145/3605764.3623985> | GET (HEAD 403 → ranged GET): HTTP 403 | 302 → https://dl.acm.org/doi/10.1145/3605764.3623985 |
| ch38 | <https://doi.org/10.1145/3618260.3649777> | GET (HEAD 403 → ranged GET): HTTP 403 | 302 → https://dl.acm.org/doi/10.1145/3618260.3649777 |
| ch32 | <https://doi.org/10.1145/362686.362692> | GET (HEAD 403 → ranged GET): HTTP 403 | 302 → https://dl.acm.org/doi/10.1145/362686.362692 |
| ch21 | <https://doi.org/10.1145/38713.38742> | GET (HEAD 403 → ranged GET): HTTP 403 | 302 → http://portal.acm.org/citation.cfm?doid=38713.38742<br>301 → https://portal.acm.org/citation.cfm?doid=38713.38742 |
| ch07 | <https://doi.org/10.1145/503124.503150> | GET (HEAD 403 → ranged GET): HTTP 403 | 302 → https://dl.acm.org/doi/10.1145/503124.503150 |
| ch16 | <https://doi.org/10.1145/582415.582418> | GET (HEAD 403 → ranged GET): HTTP 403 | 302 → https://dl.acm.org/doi/10.1145/582415.582418 |
| ch04 | <https://doi.org/10.1145/800157.805047> | GET (HEAD 403 → ranged GET): HTTP 403 | 302 → http://portal.acm.org/citation.cfm?doid=800157.805047<br>301 → https://portal.acm.org/citation.cfm?doid=800157.805047 |
| ch30 | <https://doi.org/10.1145/800216.806586> | GET (HEAD 403 → ranged GET): HTTP 403 | 302 → http://portal.acm.org/citation.cfm?doid=800216.806586<br>301 → https://portal.acm.org/citation.cfm?doid=800216.806586 |
| ch33 | <https://doi.org/10.1145/857076.857078> | GET (HEAD 403 → ranged GET): HTTP 403 | 302 → https://dl.acm.org/doi/10.1145/857076.857078 |
| ch28 | <https://doi.org/10.11591/ijeecs.v23.i2.pp980-992> | HEAD: HTTP 200 | 302 → http://ijeecs.iaescore.com/index.php/IJEECS/article/view/23780<br>301 → https://ijeecs.iaescore.com/index.php/IJEECS/article/view/23780 |
| ch32 | <https://doi.org/10.11591/ijeecs.v40.i3.pp1430-1438> | HEAD: HTTP 200 | 302 → https://ijeecs.iaescore.com/index.php/IJEECS/article/view/41268 |
| ch23 | <https://doi.org/10.11591/ijeecs.v41.i2.pp555-563> | HEAD: HTTP 200 | 302 → https://ijeecs.iaescore.com/index.php/IJEECS/article/view/38265 |
| ch09 | <https://doi.org/10.11591/ijeecs.v41.i3.pp1040-1048> | HEAD: HTTP 200 | 302 → https://ijeecs.iaescore.com/index.php/IJEECS/article/view/44903 |
| ch21 | <https://doi.org/10.11591/ijeecs.v42.i1.pp48-61> | HEAD: HTTP 200 | 302 → https://ijeecs.iaescore.com/index.php/IJEECS/article/view/45294 |
| ch19 | <https://doi.org/10.1162/089976698300017197> | GET (HEAD 403 → ranged GET): HTTP 403 | 302 → https://direct.mit.edu/neco/article/10/7/1895-1923/6224 |
| ch10, ch12, ch13 | <https://doi.org/10.1162/coli_a_00402> | GET (HEAD 403 → ranged GET): HTTP 403 | 302 → https://direct.mit.edu/coli/article/doi/10.1162/coli_a_00402/98516/Universal-Dependencies |
| ch19 | <https://doi.org/10.1162/tacl_a_00276> | GET (HEAD 403 → ranged GET): HTTP 403 | 302 → https://direct.mit.edu/tacl/article/43518 |
| ch13 | <https://doi.org/10.1162/tacl_a_00298> | GET (HEAD 403 → ranged GET): HTTP 403 | 302 → https://direct.mit.edu/tacl/article/43535 |
| ch13 | <https://doi.org/10.1162/tacl_a_00638> | GET (HEAD 403 → ranged GET): HTTP 403 | 302 → https://direct.mit.edu/tacl/article/doi/10.1162/tacl_a_00638/119630/Lost-in-the-Middle-How-Language-Models-Use-Long |
| ch02 | <https://doi.org/10.1163/9789004368811_003> | HEAD: HTTP 202 | 302 → https://brill.com/view/book/edcoll/9789004368811/BP000003.xml |
| ch24, ch25 | <https://doi.org/10.1175/1520-0493%281950%29078%3C0001:VOFEIT%3E2.0.CO;2> | HEAD: HTTP 202 | 302 → http://journals.ametsoc.org/doi/10.1175/1520-0493(1950)078%3C0001:VOFEIT%3E2.0.CO;2<br>301 → https://journals.ametsoc.org/doi/10.1175/1520-0493(1950)078%3C0001:VOFEIT%3E2.0.CO;2 |
| ch07 | <https://doi.org/10.1175/1520-0493(1950)078%3C0001:VOFEIT%3E2.0.CO;2> | HEAD: HTTP 202 | 302 → http://journals.ametsoc.org/doi/10.1175/1520-0493(1950)078%3C0001:VOFEIT%3E2.0.CO;2<br>301 → https://journals.ametsoc.org/doi/10.1175/1520-0493(1950)078%3C0001:VOFEIT%3E2.0.CO;2 |
| ch11, ch25 | <https://doi.org/10.1177/001316446002000104> | GET (HEAD 403 → ranged GET): HTTP 403 | 302 → https://journals.sagepub.com/doi/10.1177/001316446002000104 |
| ch21 | <https://doi.org/10.1177/0018720810376055> | GET (HEAD 403 → ranged GET): HTTP 403 | 302 → https://journals.sagepub.com/doi/10.1177/0018720810376055 |
| ch11 | <https://doi.org/10.1177/1049732316654870> | GET (HEAD 403 → ranged GET): HTTP 403 | 302 → https://journals.sagepub.com/doi/10.1177/1049732316654870 |
| ch27 | <https://doi.org/10.1184/R1/6573413.v1> | HEAD: HTTP 202 | 302 → https://kilthub.cmu.edu/articles/Eliminative_Argumentation_A_Basis_for_Arguing_Confidence_in_System_Properties/6573413/1 |
| ch15 | <https://doi.org/10.1186/s12859-019-3002-3> | HEAD: HTTP 200 | 302 → https://bmcbioinformatics.biomedcentral.com/articles/10.1186/s12859-019-3002-3<br>301 → https://link.springer.com/article/10.1186/s12859-019-3002-3<br>303 → https://idp.springer.com/authorize?response_type=cookie&client_id=springerlink&redirect_uri=https%3A%2F%2Flink.springer.com%2Farticle%2F10.1186%2Fs12859-019-3002-3<br>302 → https://idp.springer.com/transit?redirect_uri=https%3A%2F%2Flink.springer.com%2Farticle%2F10.1186%2Fs12859-019-3002-3&code=15d4f9e7-f42b-444a-8d47-2236ce1fa6fd<br>302 → https://link.springer.com/article/10.1186/s12859-019-3002-3?error=cookies_not_supported&code=15d4f9e7-f42b-444a-8d47-2236ce1fa6fd |
| ch02 | <https://doi.org/10.1207/s15516709cog1104_4> | GET (HEAD 403 → ranged GET): HTTP 403 | 302 → https://onlinelibrary.wiley.com/doi/10.1207/s15516709cog1104_4 |
| ch04, ch06 | <https://doi.org/10.1214/aoms/1177698950> | HEAD: HTTP 200 | 302 → http://projecteuclid.org/euclid.aoms/1177698950<br>301 → https://projecteuclid.org/euclid.aoms/1177698950 |
| ch37 | <https://doi.org/10.1214/aoms/1177731118> | HEAD: HTTP 200 | 302 → http://projecteuclid.org/euclid.aoms/1177731118<br>301 → https://projecteuclid.org/euclid.aoms/1177731118 |
| ch12, ch19, ch25 | <https://doi.org/10.1214/aos/1176344552> | HEAD: HTTP 200 | 302 → https://projecteuclid.org/journals/annals-of-statistics/volume-7/issue-1/Bootstrap-Methods-Another-Look-at-the-Jackknife/10.1214/aos/1176344552.full |
| ch06 | <https://doi.org/10.1287/mnsc.31.6.647> | GET (HEAD 403 → ranged GET): HTTP 403 | 302 → https://pubsonline.informs.org/doi/10.1287/mnsc.31.6.647 |
| ch11 | <https://doi.org/10.1287/mnsc.9.3.458> | GET (HEAD 403 → ranged GET): HTTP 403 | 302 → https://pubsonline.informs.org/doi/10.1287/mnsc.9.3.458 |
| ch10 | <https://doi.org/10.1287/opre.9.3.383> | GET (HEAD 403 → ranged GET): HTTP 403 | 302 → https://pubsonline.informs.org/doi/10.1287/opre.9.3.383 |
| ch07 | <https://doi.org/10.14778/3402707.3402747> | GET (HEAD 403 → ranged GET): HTTP 403 | 302 → https://dl.acm.org/doi/10.14778/3402707.3402747 |
| ch04, ch06 | <https://doi.org/10.1515/9780691214696> | GET (HEAD 405 → ranged GET): HTTP 200 | 302 → https://apps.crossref.org/coaccess/coaccess.html?doi=10.1515%2F9780691214696<br>301 → https://chooser.crossref.org/?doi=10.1515%2F9780691214696 |
| ch01, ch12 | <https://doi.org/10.15407/pp2020.02-03.259> | HEAD: HTTP 502 | 302 → http://pp.isofts.kiev.ua/ojs1/article/view/417 |
| ch03, ch04, ch06, ch08, ch13, ch15, ch16 | <https://doi.org/10.1561/1500000019> | GET (HEAD 403 → ranged GET): HTTP 403 | 302 → https://www.emerald.com/ftinr/article/4/1-2/1/1326508/The-Probabilistic-Relevance-Framework-BM25-and |
| ch11 | <https://doi.org/10.1561/1500000024> | GET (HEAD 403 → ranged GET): HTTP 403 | 302 → https://www.emerald.com/ftinr/article/6/2-3/127/1356738/Expertise-Retrieval |
| ch03 | <https://doi.org/10.1609/aaai.v29i1.9602> | HEAD: HTTP 200 | 302 → https://ojs.aaai.org/index.php/AAAI/article/view/9602 |
| ch33 | <https://doi.org/10.1609/aaai.v32i1.11797> | HEAD: HTTP 200 | 302 → https://ojs.aaai.org/index.php/AAAI/article/view/11797 |
| ch25 | <https://doi.org/10.1613/jair.1.12125> | HEAD: HTTP 200 | 302 → https://jair.org/index.php/jair/article/view/12125 |
| ch21 | <https://doi.org/10.1613/jair.1141> | HEAD: HTTP 200 | 302 → https://jair.org/index.php/jair/article/view/10362 |
| ch17 | <https://doi.org/10.1613/jair.1705> | HEAD: HTTP 200 | 302 → https://jair.org/index.php/jair/article/view/10457 |
| ch07, ch33 | <https://doi.org/10.1613/jair.2375> | HEAD: HTTP 200 | 302 → https://jair.org/index.php/jair/article/view/10534 |
| ch02 | <https://doi.org/10.17487/RFC2119> | HEAD: HTTP 200 | 302 → https://www.rfc-editor.org/info/rfc2119<br>301 → https://www.rfc-editor.org/info/rfc2119/ |
| ch02 | <https://doi.org/10.17487/RFC3161> | HEAD: HTTP 200 | 302 → https://www.rfc-editor.org/info/rfc3161<br>301 → https://www.rfc-editor.org/info/rfc3161/ |
| ch02 | <https://doi.org/10.17487/RFC8785> | HEAD: HTTP 200 | 302 → https://www.rfc-editor.org/info/rfc8785<br>301 → https://www.rfc-editor.org/info/rfc8785/ |
| ch13 | <https://doi.org/10.18653/v1/2020.acl-demos.14> | HEAD: HTTP 200 | 302 → https://www.aclweb.org/anthology/2020.acl-demos.14<br>301 → https://aclanthology.org/2020.acl-demos.14<br>301 → https://aclanthology.org/2020.acl-demos.14/ |
| ch13 | <https://doi.org/10.18653/v1/2020.acl-main.442> | HEAD: HTTP 200 | 302 → https://www.aclweb.org/anthology/2020.acl-main.442<br>301 → https://aclanthology.org/2020.acl-main.442<br>301 → https://aclanthology.org/2020.acl-main.442/ |
| ch06 | <https://doi.org/10.18653/v1/2020.emnlp-main.550> | HEAD: HTTP 200 | 302 → https://aclanthology.org/2020.emnlp-main.550<br>301 → https://aclanthology.org/2020.emnlp-main.550/ |
| ch15 | <https://doi.org/10.18653/v1/2023.eacl-main.148> | HEAD: HTTP 200 | 302 → https://aclanthology.org/2023.eacl-main.148<br>301 → https://aclanthology.org/2023.eacl-main.148/ |
| ch21 | <https://doi.org/10.18653/v1/2023.emnlp-main.674> | HEAD: HTTP 200 | 302 → https://aclanthology.org/2023.emnlp-main.674<br>301 → https://aclanthology.org/2023.emnlp-main.674/ |
| ch34 | <https://doi.org/10.18653/v1/2023.findings-acl.847> | HEAD: HTTP 200 | 302 → https://aclanthology.org/2023.findings-acl.847<br>301 → https://aclanthology.org/2023.findings-acl.847/ |
| ch16 | <https://doi.org/10.18653/v1/2024.eacl-demo.16> | HEAD: HTTP 200 | 302 → https://aclanthology.org/2024.eacl-demo.16<br>301 → https://aclanthology.org/2024.eacl-demo.16/ |
| ch33 | <https://doi.org/10.18653/v1/P16-1228> | HEAD: HTTP 200 | 302 → http://aclweb.org/anthology/P16-1228<br>301 → https://aclweb.org/anthology/P16-1228<br>301 → https://aclanthology.org/P16-1228<br>301 → https://aclanthology.org/P16-1228/ |
| ch02 | <https://doi.org/10.21236/AD0770768> | GET (HEAD 404 → ranged GET): HTTP 404 | — |
| ch01, ch04 | <https://doi.org/10.21236/ADA046289> | GET (HEAD 403 → ranged GET): HTTP 403 | 302 → https://apps.dtic.mil/sti/citations/tr/ADA046289 |
| ch06 | <https://doi.org/10.2139/ssrn.3063289> | GET (HEAD 403 → ranged GET): HTTP 403 | 302 → https://www.ssrn.com/abstract=3063289<br>302 → https://papers.ssrn.com/sol3/papers.cfm?abstract_id=3063289 |
| ch06 | <https://doi.org/10.2140/pjm.1955.5.285> | HEAD: HTTP 200 | 302 → http://msp.org/pjm/1955/5-2/p11.xhtml<br>302 → https://msp.org/pjm/1955/5-2/p11.xhtml |
| ch06 | <https://doi.org/10.2307/2269326> | HEAD: HTTP 200 | 302 → https://www.cambridge.org/core/product/identifier/S0022481200038664/type/journal_article<br>301 → https://www.cambridge.org/core/journals/journal-of-symbolic-logic/article/note-on-the-entscheidungsproblem/9461BEAD94BB16D56EC78933D7D67DEF<br>302 → https://www.cambridge.org/core/journals/journal-of-symbolic-logic/article/abs/note-on-the-entscheidungsproblem/9461BEAD94BB16D56EC78933D7D67DEF |
| ch26 | <https://doi.org/10.2307/2274239> | HEAD: HTTP 200 | 302 → https://www.cambridge.org/core/product/identifier/S0022481200032849/type/journal_article<br>301 → https://www.cambridge.org/core/journals/journal-of-symbolic-logic/article/on-the-logic-of-theory-change-partial-meet-contraction-and-revision-functions/7ED837BAD5FB6D9A7C77906D73527F9C<br>302 → https://www.cambridge.org/core/journals/journal-of-symbolic-logic/article/abs/on-the-logic-of-theory-change-partial-meet-contraction-and-revision-functions/7ED837BAD5FB6D9A7C77906D73527F9C |
| ch26 | <https://doi.org/10.2307/2275912> | HEAD: HTTP 200 | 302 → https://www.cambridge.org/core/product/identifier/S0022481200019551/type/journal_article<br>301 → /core/journals/journal-of-symbolic-logic/article/kernel-contraction/C822638049E45B34EE10FBE45A8865E1<br>302 → /core/journals/journal-of-symbolic-logic/article/abs/kernel-contraction/C822638049E45B34EE10FBE45A8865E1<br>final: https://www.cambridge.org/core/journals/journal-of-symbolic-logic/article/abs/kernel-contraction/C822638049E45B34EE10FBE45A8865E1 |
| ch37 | <https://doi.org/10.2307/2346806> | HEAD: HTTP 200 | 302 → https://www.jstor.org/stable/10.2307/2346806?origin=crossref |
| ch01 | <https://doi.org/10.2307/2371809> | HEAD: HTTP 200 | 302 → https://www.jstor.org/stable/2371809?origin=crossref |
| ch25 | <https://doi.org/10.2307/3250961> | GET (HEAD 403 → ranged GET): HTTP 403 | 302 → https://misq.umn.edu/misq/article/25/1/107/1295/Review-Knowledge-Management-and-Knowledge |
| ch07 | <https://doi.org/10.24963/ijcai.2020/688> | HEAD: HTTP 200 | 302 → https://www.ijcai.org/proceedings/2020/688<br>301 → https://ijcai.org/proceedings/2020/688 |
| ch04, ch06, ch07, ch24, ch26, ch28 | <https://doi.org/10.3233/AIC-1994-7104> | GET (HEAD 403 → ranged GET): HTTP 403 | 302 → https://journals.sagepub.com/doi/full/10.3233/AIC-1994-7104 |
| ch18 | <https://doi.org/10.48550/arXiv.1701.06538> | HEAD: HTTP 200 | 302 → https://arxiv.org/abs/1701.06538 |
| ch18 | <https://doi.org/10.48550/arXiv.2508.10925> | HEAD: HTTP 200 | 302 → https://arxiv.org/abs/2508.10925 |
| ch20 | <https://doi.org/10.48550/arXiv.2603.01890> | HEAD: HTTP 200 | 302 → https://arxiv.org/abs/2603.01890 |
| ch12 | <https://doi.org/10.5281/zenodo.19332453> | HEAD: HTTP 200 | 302 → https://zenodo.org/doi/10.5281/zenodo.19332453<br>302 → https://zenodo.org/records/19332453 |
| ch12, ch13, ch14 | <https://doi.org/10.5281/zenodo.20509698> | HEAD: HTTP 200 | 302 → https://zenodo.org/doi/10.5281/zenodo.20509698<br>302 → https://zenodo.org/records/20509698 |
| ch37 | <https://doi.org/10.5604/01.3001.0012.1474> | HEAD: HTTP 200 | 302 → https://securityanddefence.pl/Open-source-intelligence-OSINT-as-an-element-of-military-recon,103337,0,2.html |
| ch02, ch37 | <https://doi.org/10.59161/JCGM100-2008E> | HEAD: HTTP 200 | 302 → https://www.bipm.org/doi/10.59161/JCGM100-2008E<br>302 → https://www.bipm.org/en/doi/10.59161/JCGM100-2008E |
| ch02 | <https://doi.org/10.59161/JCGM106-2012> | HEAD: HTTP 200 | 302 → https://www.bipm.org/doi/10.59161/JCGM106-2012<br>302 → https://www.bipm.org/en/doi/10.59161/JCGM106-2012 |
| ch03, ch05, ch10, ch16, ch25, ch26 | <https://doi.org/10.6028/NIST.AI.100-1> | GET (HEAD 404 → ranged GET): HTTP 206 | 302 → http://nvlpubs.nist.gov/nistpubs/ai/NIST.AI.100-1.pdf<br>301 → https://nvlpubs.nist.gov/nistpubs/ai/NIST.AI.100-1.pdf |
| ch02, ch10, ch25 | <https://doi.org/10.6028/NIST.AI.600-1> | GET (HEAD 404 → ranged GET): HTTP 206 | 302 → https://nvlpubs.nist.gov/nistpubs/ai/NIST.AI.600-1.pdf |
| ch05 | <https://doi.org/10.6028/NIST.FIPS.140-3> | GET (HEAD 404 → ranged GET): HTTP 206 | 302 → https://nvlpubs.nist.gov/nistpubs/FIPS/NIST.FIPS.140-3.pdf |
| ch15, ch27, ch32 | <https://doi.org/10.6028/NIST.FIPS.180-4> | GET (HEAD 404 → ranged GET): HTTP 206 | 302 → https://nvlpubs.nist.gov/nistpubs/FIPS/NIST.FIPS.180-4.pdf |
| ch32 | <https://doi.org/10.6028/NIST.FIPS.204> | GET (HEAD 404 → ranged GET): HTTP 206 | 302 → https://nvlpubs.nist.gov/nistpubs/FIPS/NIST.FIPS.204.pdf |
| ch02, ch16 | <https://doi.org/10.6028/NIST.SP.800-162> | HEAD: HTTP 200 | 302 → https://nvlpubs.nist.gov/nistpubs/SpecialPublications/NIST.SP.800-162.pdf |
| ch21 | <https://doi.org/10.6028/NIST.SP.800-207> | GET (HEAD 404 → ranged GET): HTTP 206 | 302 → https://nvlpubs.nist.gov/nistpubs/SpecialPublications/NIST.SP.800-207.pdf |
| ch23, ch32 | <https://doi.org/10.6028/NIST.SP.800-218> | GET (HEAD 404 → ranged GET): HTTP 206 | 302 → https://nvlpubs.nist.gov/nistpubs/SpecialPublications/NIST.SP.800-218.pdf |
| ch27, ch30 | <https://doi.org/10.65391/r1386> | HEAD: HTTP 200 | 302 → https://scsc.uk/resources/citation_r1386.html |
| ch22 | <https://doi.org/10.7551/mitpress/11810.001.0001> | GET (HEAD 403 → ranged GET): HTTP 403 | 302 → https://direct.mit.edu/books/book/4581/Cybernetics-or-Control-and-Communication-in-the |
| ch09 | <https://dou.ua/forums/topic/49883/> | HEAD: HTTP 200 | — |
| ch08, ch10 | <https://dou.ua/forums/topic/60526/> | HEAD: HTTP 200 | — |
| ch29 | <https://dou.ua/forums/topic/61201/> | HEAD: HTTP 200 | — |
| ch01 | <https://dspace.mit.edu/handle/1721.1/6089> | HEAD: HTTP 200 | — |
| ch06, ch28 | <https://en.wikisource.org/wiki/Popular_Science_Monthly/Volume_13/August_1878/Illustrations_of_the_Logic_of_Science_VI> | HEAD: HTTP 200 | — |
| ch03 | <https://eur-lex.europa.eu/eli/reg/2014/910/oj> | HEAD: HTTP 202 | — |
| ch01, ch23 | <https://eur-lex.europa.eu/eli/reg/2024/1689/oj> | HEAD: HTTP 202 | — |
| ch24 | <https://extapps.ksc.nasa.gov/Reliability/Documents/Fault_Tree_Handbook_with_Aerospace_Applications_August_2002.pdf> | HEAD: HTTP 200 | — |
| ch10 | <https://fmi-standard.org/docs/3.0.2/> | HEAD: HTTP 200 | — |
| ch10, ch21, ch23 | <https://genai.owasp.org/llmrisk/llm01-prompt-injection/> | HEAD: HTTP 200 | — |
| ch23 | <https://genai.owasp.org/llmrisk/llm022025-sensitive-information-disclosure/> | HEAD: HTTP 200 | — |
| ch10 | <https://genai.owasp.org/llmrisk/llm042025-data-and-model-poisoning/> | HEAD: HTTP 200 | — |
| ch29 | <https://gist.github.com/karpathy/442a6bf555914893e9891c11519de94f> | HEAD: HTTP 200 | — |
| ch17 | <https://github.com/NRules/NRules> | HEAD: HTTP 200 | — |
| ch12 | <https://github.com/NVIDIA/TensorRT-LLM> | HEAD: HTTP 200 | — |
| ch10 | <https://github.com/PaddlePaddle/PaddleOCR> | HEAD: HTTP 200 | — |
| ch08 | <https://github.com/PaddlePaddle/PaddleOCR/blob/main/docs/version3.x/pipeline_usage/PP-StructureV3.en.md> | HEAD: HTTP 200 | — |
| ch17 | <https://github.com/blevesearch/bleve> | HEAD: HTTP 200 | — |
| ch17 | <https://github.com/cayleygraph/cayley> | HEAD: HTTP 200 | — |
| ch17 | <https://github.com/cel-expr/cel-go> | HEAD: HTTP 200 | — |
| ch17 | <https://github.com/dgraph-io/dgraph> | HEAD: HTTP 200 | — |
| ch23 | <https://github.com/dig-team/AMIE> | HEAD: HTTP 200 | — |
| ch08, ch10 | <https://github.com/docling-project/docling> | HEAD: HTTP 200 | — |
| ch30 | <https://github.com/exploitgym> | GET (HEAD 404 → ranged GET): HTTP 404 | — |
| ch23 | <https://github.com/flyingmutant/rapid> | HEAD: HTTP 200 | — |
| ch14 | <https://github.com/ggerganov/llama.cpp/blob/master/grammars/README.md> | HEAD: HTTP 200 | 301 → https://github.com/ggml-org/llama.cpp/blob/master/grammars/README.md |
| ch12, ch17 | <https://github.com/ggml-org/llama.cpp> | HEAD: HTTP 200 | — |
| ch29 | <https://github.com/ggml-org/llama.cpp/blob/master/grammars/README.md> | HEAD: HTTP 200 | — |
| ch17 | <https://github.com/google/mangle> | HEAD: HTTP 200 | — |
| ch17 | <https://github.com/gorules/zen> | HEAD: HTTP 200 | — |
| ch17 | <https://github.com/hyperjumptech/grule-rule-engine> | HEAD: HTTP 200 | — |
| ch27 | <https://github.com/in-toto/attestation/blob/main/spec/README.md> | HEAD: HTTP 200 | — |
| ch17 | <https://github.com/krotik/eliasdb> | HEAD: HTTP 200 | — |
| ch17 | <https://github.com/ledongthuc/pdf> | HEAD: HTTP 200 | — |
| ch12, ch25 | <https://github.com/ml-explore/mlx-lm> | HEAD: HTTP 200 | — |
| ch17 | <https://github.com/nats-io/nats-server> | HEAD: HTTP 200 | — |
| ch22 | <https://github.com/nickovic/rtamt> | HEAD: HTTP 200 | — |
| ch17 | <https://github.com/ollama/ollama> | HEAD: HTTP 200 | — |
| ch29 | <https://github.com/ollama/ollama/blob/main/docs/modelfile.mdx> | HEAD: HTTP 200 | — |
| ch17 | <https://github.com/open-policy-agent/opa> | HEAD: HTTP 200 | — |
| ch30 | <https://github.com/openvex/spec/blob/main/OPENVEX-SPEC.md> | HEAD: HTTP 200 | — |
| ch15, ch17 | <https://github.com/pgvector/pgvector> | HEAD: HTTP 200 | — |
| ch17 | <https://github.com/pgvector/pgvector-go> | HEAD: HTTP 200 | — |
| ch17 | <https://github.com/philippgille/chromem-go> | HEAD: HTTP 200 | — |
| ch15 | <https://github.com/process-intelligence-solutions/pm4py> | HEAD: HTTP 200 | — |
| ch17 | <https://github.com/sashabaranov/go-openai> | HEAD: HTTP 200 | — |
| ch22 | <https://github.com/sidliarchukpetro/witness-integrity> | HEAD: HTTP 200 | — |
| ch30 | <https://github.com/strictdoc-project/reqif> | HEAD: HTTP 200 | — |
| ch17 | <https://github.com/tree-sitter/go-tree-sitter> | HEAD: HTTP 200 | — |
| ch25 | <https://github.com/unslothai/unsloth> | HEAD: HTTP 200 | — |
| ch13 | <https://github.com/urchade/GLiNER> | HEAD: HTTP 200 | — |
| ch17 | <https://github.com/weaviate/weaviate> | HEAD: HTTP 200 | — |
| ch17 | <https://github.com/yalue/onnxruntime_go> | HEAD: HTTP 200 | — |
| ch32 | <https://go.dev/doc/security/fuzz/> | HEAD: HTTP 200 | — |
| ch32 | <https://go.dev/ref/spec#For_range> | HEAD: HTTP 200 | — |
| ch18 | <https://hailo.ai/products/ai-accelerators/hailo-10h-ai-accelerator/> | HEAD: HTTP 200 | — |
| ch18 | <https://hailo.ai/products/ai-accelerators/hailo-8-ai-accelerator/> | HEAD: HTTP 200 | — |
| ch25 | <https://huggingface.co/docs/peft/index> | HEAD: HTTP 200 | — |
| ch25 | <https://huggingface.co/docs/transformers/installation> | HEAD: HTTP 200 | — |
| ch25 | <https://huggingface.co/docs/trl/en/sft_trainer> | HEAD: HTTP 200 | — |
| ch25 | <https://huggingface.co/docs/trl/index> | HEAD: HTTP 200 | — |
| ch12 | <https://huggingface.co/lapa-llm/lapa-v0.1.2-instruct> | HEAD: HTTP 200 | — |
| ch12 | <https://huggingface.co/lapa-llm/tokenizer> | HEAD: HTTP 200 | — |
| ch18 | <https://huggingface.co/openai/gpt-oss-120b> | HEAD: HTTP 200 | — |
| ch23 | <https://hypothesis.readthedocs.io/en/latest/stateful.html> | HEAD: HTTP 200 | — |
| ch37 | <https://i2group.com/solutions/i2-analysts-notebook> | HEAD: HTTP 200 | — |
| ch05 | <https://ijcai.org/proceedings/1977/231> | HEAD: HTTP 200 | — |
| ch23 | <https://ilasp.com/> | HEAD: HTTP NETWORK_ERROR; SSLCertVerificationError: [SSL: CERTIFICATE_VERIFY_FAILED] certificate verify failed: unable to get local issuer certificate (_ssl.c:1081) | — |
| ch06, ch31 | <https://jmlr.org/papers/v12/shervashidze11a.html> | HEAD: HTTP 200 | — |
| ch26 | <https://jmlr.org/papers/v24/23-0130.html> | HEAD: HTTP 200 | — |
| ch08 | <https://json-schema.org/draft/2020-12> | HEAD: HTTP 200 | — |
| ch16, ch19 | <https://json-schema.org/specification> | HEAD: HTTP 200 | — |
| ch37, ch38 | <https://juliamath.github.io/MeasureTheory.jl/stable/> | HEAD: HTTP 200 | — |
| ch23, ch33 | <https://lamport.azurewebsites.net/tla/book.html> | HEAD: HTTP 200 | — |
| ch37 | <https://learn.microsoft.com/en-us/azure/sentinel/overview> | HEAD: HTTP 200 | — |
| ch32 | <https://man7.org/linux/man-pages/man2/madvise.2.html> | HEAD: HTTP 200 | — |
| ch32 | <https://man7.org/tlpi/> | HEAD: HTTP 200 | — |
| ch11, ch26 | <https://minds.wisconsin.edu/handle/1793/60660> | HEAD: HTTP 200 | 301 → https://minds.wisc.edu/handle/1793/60660 |
| ch32 | <https://mirrors.edge.kernel.org/pub/linux/kernel/people/paulmck/perfbook/perfbook.html> | HEAD: HTTP 200 | — |
| ch11 | <https://mitpress.mit.edu/9780262193009/knowledge-engineering-and-management/> | HEAD: HTTP 200 | — |
| ch12 | <https://mitpress.mit.edu/9780262561167/wordnet/> | HEAD: HTTP 200 | — |
| ch31, ch34 | <https://mitpress.mit.edu/9780262661133/> | HEAD: HTTP 200 | — |
| ch25 | <https://mlflow.org/docs/latest/ml/tracking/> | HEAD: HTTP 200 | — |
| ch22 | <https://nightlies.apache.org/flink/flink-docs-stable/docs/concepts/time/> | HEAD: HTTP 200 | — |
| ch18 | <https://nvidianews.nvidia.com/news/nvidia-dgx-spark-arrives-for-worlds-ai-developers> | HEAD: HTTP 200 | — |
| ch18 | <https://nvidianews.nvidia.com/news/nvidia-dgx-station-for-windows-puts-a-trillion-parameter-ai-supercomputer-on-every-enterprise-desk> | HEAD: HTTP 200 | — |
| ch18 | <https://nvidianews.nvidia.com/news/nvidia-microsoft-windows-pcs-agents-rtx-spark> | HEAD: HTTP 200 | — |
| ch04 | <https://ojs.aaai.org/aimagazine/index.php/aimagazine/article/view/510> | HEAD: HTTP 200 | — |
| ch06 | <https://ojs.aaai.org/aimagazine/index.php/aimagazine/article/view/542> | HEAD: HTTP 200 | — |
| ch38 | <https://ojs.aaai.org/index.php/AAAI/article/view/11674> | HEAD: HTTP 200 | — |
| ch12 | <https://ollama.com/blog/structured-outputs> | HEAD: HTTP 200 | — |
| ch12 | <https://onnxruntime.ai/> | HEAD: HTTP 200 | — |
| ch08 | <https://onnxruntime.ai/docs/execution-providers/> | HEAD: HTTP 200 | — |
| ch04 | <https://op.europa.eu/en/web/eu-vocabularies/eli> | HEAD: HTTP 200 | — |
| ch09 | <https://open-cmsis-pack.github.io/svd-spec/main/index.html> | HEAD: HTTP 200 | — |
| ch07 | <https://openalex.org/W1672891595> | HEAD: HTTP 200 | 301 → https://openalex.org/works/W1672891595 |
| ch04 | <https://openlibrary.org/books/OL5615707M> | HEAD: HTTP NETWORK_ERROR; ConnectionRefusedError: [Errno 111] Connection refused | — |
| ch11 | <https://openlibrary.org/works/OL117061W> | HEAD: HTTP NETWORK_ERROR; ConnectionRefusedError: [Errno 111] Connection refused | — |
| ch06 | <https://openlibrary.org/works/OL1311228W> | HEAD: HTTP NETWORK_ERROR; ConnectionRefusedError: [Errno 111] Connection refused | — |
| ch01 | <https://openlibrary.org/works/OL16444441W> | HEAD: HTTP 200 | 303 → https://openlibrary.org/verify_human?next=/works/OL16444441W |
| ch09 | <https://openlibrary.org/works/OL16585333W> | HEAD: HTTP NETWORK_ERROR; ConnectionRefusedError: [Errno 111] Connection refused | — |
| ch06 | <https://openlibrary.org/works/OL17872278W> | HEAD: HTTP NETWORK_ERROR; ConnectionRefusedError: [Errno 111] Connection refused | — |
| ch03 | <https://openlibrary.org/works/OL1814674W> | HEAD: HTTP NETWORK_ERROR; ConnectionRefusedError: [Errno 111] Connection refused | — |
| ch15 | <https://openlibrary.org/works/OL21407442W> | HEAD: HTTP NETWORK_ERROR; ConnectionRefusedError: [Errno 111] Connection refused | — |
| ch04 | <https://openlibrary.org/works/OL24825917W> | HEAD: HTTP NETWORK_ERROR; ConnectionRefusedError: [Errno 111] Connection refused | — |
| ch15 | <https://openlibrary.org/works/OL3307771W> | HEAD: HTTP NETWORK_ERROR; ConnectionRefusedError: [Errno 111] Connection refused | — |
| ch11 | <https://openlibrary.org/works/OL3515903W> | HEAD: HTTP NETWORK_ERROR; ConnectionRefusedError: [Errno 111] Connection refused | — |
| ch11 | <https://openlibrary.org/works/OL4305747W> | HEAD: HTTP NETWORK_ERROR; ConnectionRefusedError: [Errno 111] Connection refused | — |
| ch11 | <https://openlibrary.org/works/OL5282413W> | HEAD: HTTP NETWORK_ERROR; ConnectionRefusedError: [Errno 111] Connection refused | — |
| ch01 | <https://openlibrary.org/works/OL5532645W> | HEAD: HTTP 200 | 303 → https://openlibrary.org/verify_human?next=/works/OL5532645W |
| ch02, ch03, ch06, ch31 | <https://openlibrary.org/works/OL5959470W> | HEAD: HTTP NETWORK_ERROR; ConnectionRefusedError: [Errno 111] Connection refused | 303 → https://openlibrary.org/verify_human?next=/works/OL5959470W |
| ch02 | <https://openlibrary.org/works/OL6408872W> | HEAD: HTTP NETWORK_ERROR; ConnectionRefusedError: [Errno 111] Connection refused | — |
| ch09 | <https://openlibrary.org/works/OL9203874W> | HEAD: HTTP NETWORK_ERROR; ConnectionRefusedError: [Errno 111] Connection refused | — |
| ch10 | <https://openlineage.io/docs/> | HEAD: HTTP 200 | — |
| ch16 | <https://openlineage.io/docs/spec/object-model> | HEAD: HTTP 200 | 301 → https://openlineage.io/docs/spec/object-model/ |
| ch29 | <https://openreview.net/forum?id=rJgMlhRctm> | HEAD: HTTP 200 | 307 → https://openreview.net/challenge?redirect=%2Fforum%3Fid%3DrJgMlhRctm |
| ch26 | <https://papers.neurips.cc/paper_files/paper/2023/hash/1b44b878bb782e6954cd888628510e90-Abstract-Conference.html> | HEAD: HTTP 200 | — |
| ch26 | <https://papers.neurips.cc/paper_files/paper/2023/hash/91edff07232fb1b55a505a9e9f6c0ff3-Abstract-Conference.html> | HEAD: HTTP 200 | — |
| ch02 | <https://papers.nips.cc/paper_files/paper/2017/hash/4a8423d5e91fda00bb7e46540e2b0cf1-Abstract.html> | HEAD: HTTP 200 | — |
| ch20 | <https://papers.nips.cc/paper_files/paper/2017/hash/8a20a8621978632d76c43dfd28b67767-Abstract.html> | HEAD: HTTP 200 | — |
| ch26 | <https://papers.nips.cc/paper_files/paper/2017/hash/f87522788a2be2d171666752f97ddebb-Abstract.html> | HEAD: HTTP 200 | — |
| ch32 | <https://people.freebsd.org/~lstewart/articles/cpumemory.pdf> | HEAD: HTTP 200 | — |
| ch17 | <https://peps.python.org/pep-0703/> | HEAD: HTTP 200 | — |
| ch14 | <https://pint.readthedocs.io/en/stable/user/nonmult.html> | HEAD: HTTP 200 | — |
| ch02, ch06 | <https://plato.stanford.edu/entries/abduction/> | HEAD: HTTP 200 | — |
| ch28 | <https://plato.stanford.edu/entries/abduction/peirce.html> | HEAD: HTTP 200 | — |
| ch31 | <https://plato.stanford.edu/entries/aristotle-logic/> | HEAD: HTTP 200 | — |
| ch02 | <https://plato.stanford.edu/entries/epistemology-social/> | HEAD: HTTP 200 | — |
| ch02 | <https://plato.stanford.edu/entries/hermeneutics/> | HEAD: HTTP 200 | — |
| ch02 | <https://plato.stanford.edu/entries/knowledge-analysis/> | HEAD: HTTP 200 | — |
| ch28 | <https://potassco.org/clingo/> | HEAD: HTTP 200 | — |
| ch02 | <https://press.uchicago.edu/ucp/books/book/chicago/S/bo3684126.html> | GET (HEAD 404 → ranged GET): HTTP 206 | — |
| ch02 | <https://proceedings.mlr.press/r0/pearl95a.html> | HEAD: HTTP 200 | — |
| ch26 | <https://proceedings.mlr.press/v15/ross11a.html> | HEAD: HTTP 200 | — |
| ch24 | <https://proceedings.mlr.press/v162/budhathoki22a.html> | HEAD: HTTP 200 | — |
| ch03, ch04, ch06, ch25, ch37 | <https://proceedings.mlr.press/v70/guo17a.html> | HEAD: HTTP 200 | — |
| ch33 | <https://proceedings.mlr.press/v80/xu18h.html> | HEAD: HTTP 200 | — |
| ch37, ch38 | <https://proceedings.mlr.press/v84/ge18b.html> | HEAD: HTTP 200 | — |
| ch26 | <https://proceedings.mlr.press/v97/laroche19a.html> | HEAD: HTTP 200 | — |
| ch19, ch38 | <https://proceedings.mlr.press/v119/guu20a.html> | HEAD: HTTP 200 | — |
| ch13 | <https://proceedings.neurips.cc/paper/2020/hash/6b493230205f780e1bc26945df7481e5-Abstract.html> | HEAD: HTTP 200 | — |
| ch06 | <https://proceedings.neurips.cc/paper_files/paper/2017/hash/8a20a8621978632d76c43dfd28b67767-Abstract.html> | HEAD: HTTP 200 | — |
| ch38 | <https://proceedings.neurips.cc/paper_files/paper/2023/hash/a85b405ed65c6477a4fe8302b5e06ce7-Abstract-Conference.html> | HEAD: HTTP 200 | — |
| ch15 | <https://protege.stanford.edu/> | HEAD: HTTP 200 | — |
| ch07 | <https://protege.stanford.edu/publications/ontology_development/ontology101.pdf> | HEAD: HTTP 200 | — |
| ch16 | <https://protobuf.dev/> | HEAD: HTTP 200 | — |
| ch22 | <https://ptolemy.berkeley.edu/books/leeseshia/> | HEAD: HTTP 200 | — |
| ch15, ch18 | <https://pubs.opengroup.org/onlinepubs/9799919799/functions/mmap.html> | HEAD: HTTP 200 | — |
| ch13 | <https://pydantic.dev/docs/validation/latest/concepts/strict_mode/> | HEAD: HTTP 200 | — |
| ch15 | <https://rasbt.github.io/mlxtend/user_guide/frequent_patterns/association_rules/> | HEAD: HTTP 200 | — |
| ch32 | <https://reproducible-builds.org/docs/definition/> | HEAD: HTTP 200 | — |
| ch32 | <https://reproducible-builds.org/specs/source-date-epoch/> | HEAD: HTTP 200 | — |
| ch08 | <https://research.google/pubs/data-validation-for-machine-learning/> | HEAD: HTTP 200 | — |
| ch23 | <https://research.google/pubs/hidden-technical-debt-in-machine-learning-systems/> | HEAD: HTTP 200 | — |
| ch35 | <https://research.google/pubs/in-datacenter-performance-analysis-of-a-tensor-processing-unit/> | HEAD: HTTP 200 | — |
| ch29, ch38 | <https://research.google/pubs/self-consistency-improves-chain-of-thought-reasoning-in-language-models/> | HEAD: HTTP 200 | — |
| ch22, ch37 | <https://research.google/pubs/the-dataflow-model-a-practical-approach-to-balancing-correctness-latency-and-cost-in-massive-scale-unbounded-out-of-order-data-processing/> | HEAD: HTTP 200 | — |
| ch18 | <https://riscv.org/specifications/ratified/> | HEAD: HTTP 200 | — |
| ch22 | <https://riverml.xyz/latest/> | HEAD: HTTP 200 | — |
| ch22 | <https://riverml.xyz/latest/api/drift/ADWIN/> | HEAD: HTTP 200 | — |
| ch26 | <https://riverml.xyz/latest/api/evaluate/progressive-val-score/> | HEAD: HTTP 200 | — |
| ch15 | <https://robot.obolibrary.org/> | HEAD: HTTP 200 | — |
| ch15 | <https://robot.obolibrary.org/reason> | HEAD: HTTP 200 | — |
| ch15 | <https://robot.obolibrary.org/template> | HEAD: HTTP 200 | — |
| ch15 | <https://robot.obolibrary.org/verify> | HEAD: HTTP 200 | — |
| ch25 | <https://rocm.docs.amd.com/en/docs-7.2.3/how-to/rocm-for-ai/training/> | HEAD: HTTP 200 | — |
| ch25 | <https://scikit-learn.org/stable/modules/cross_validation.html> | HEAD: HTTP 200 | — |
| ch27 | <https://scsc.uk/file/gc-main/GSNv2-to-v3_changes-1092.pdf> | HEAD: HTTP 200 | — |
| ch37 | <https://securityanddefence.pl/Open-source-intelligence-OSINT-as-an-element-of-military-recon%2C103337%2C0%2C2.html> | HEAD: HTTP 200 | — |
| ch16 | <https://semver.org/spec/v2.0.0.html> | HEAD: HTTP 200 | — |
| ch32 | <https://slsa.dev/spec/v1.0/> | HEAD: HTTP 200 | — |
| ch27 | <https://slsa.dev/spec/v1.0/provenance> | HEAD: HTTP 200 | — |
| ch10 | <https://slsa.dev/spec/v1.2/> | HEAD: HTTP 200 | — |
| ch07 | <https://smt-lib.org/> | HEAD: HTTP 200 | — |
| ch14 | <https://smtlib.cs.uiowa.edu/papers/smt-lib-reference-v2.6-r2021-05-12.pdf> | HEAD: HTTP 200 | 302 → https://smt-lib.github.io/papers/smt-lib-reference-v2.6-r2021-05-12.pdf<br>301 → https://smt-lib.org/papers/smt-lib-reference-v2.6-r2021-05-12.pdf |
| ch28 | <https://souffle-lang.github.io/docs.html> | HEAD: HTTP 200 | — |
| ch10 | <https://spdx.dev/use/specifications/> | HEAD: HTTP 200 | — |
| ch08 | <https://spec.commonmark.org/0.31.2/> | HEAD: HTTP 200 | — |
| ch16 | <https://spec.openapis.org/oas/v3.1.0.html> | HEAD: HTTP 200 | — |
| ch10 | <https://standards.ieee.org/ieee/1685/7414/> | HEAD: HTTP 200 | — |
| ch18 | <https://tenstorrent.com/hardware/blackhole> | HEAD: HTTP 200 | 308 → https://tenstorrent.com/hardware/cards |
| ch22 | <https://theupdateframework.github.io/specification/latest/> | HEAD: HTTP 200 | — |
| ch08, ch12 | <https://tika.apache.org/> | HEAD: HTTP 200 | — |
| ch19 | <https://ucum.org/ucum> | HEAD: HTTP 200 | — |
| ch10 | <https://ucum.org/ucum.html> | HEAD: HTTP 200 | 301 → https://ucum.org/ucum |
| ch30 | <https://uptane.org/docs/latest/standard/uptane-standard> | HEAD: HTTP 200 | — |
| ch25, ch30 | <https://vda-qmc.de/wp-content/uploads/2023/12/Automotive-SPICE-PAM-v40.pdf> | HEAD: HTTP 200 | — |
| ch28, ch32 | <http://webdam.inria.fr/Alice/> | HEAD: HTTP 200 | — |
| ch15 | <https://webstore.iec.ch/en/publication/22794> | HEAD: HTTP 200 | — |
| ch18 | <https://webstore.iec.ch/en/publication/2452> | HEAD: HTTP 200 | — |
| ch18 | <https://webstore.iec.ch/en/publication/33615> | HEAD: HTTP 200 | — |
| ch21 | <https://webstore.iec.ch/en/publication/5515> | HEAD: HTTP 200 | — |
| ch22 | <https://webstore.iec.ch/en/publication/7029> | HEAD: HTTP 200 | — |
| ch33 | <https://webstore.iec.ch/publication/5515> | HEAD: HTTP 200 | 301 → https://webstore.iec.ch/en/publication/5515 |
| ch18 | <https://www.amd.com/en/blogs/2026/amd-powers-next-generation-agent-computers-with-new-ryzen-ai-hal.html> | HEAD: HTTP 200 | — |
| ch18 | <https://www.amd.com/en/products/processors/desktops/ryzen/ryzen-ai-halo/ryzen-ai-max-plus-395.html> | HEAD: HTTP 200 | — |
| ch16 | <https://www.asyncapi.com/docs/reference/specification/v3.0.0> | HEAD: HTTP 200 | — |
| ch37 | <https://www.bnd.bund.de/DE/Die_Arbeit/Analyse/analyse_node.html> | GET (HEAD 400 → ranged GET): HTTP 400 | 303 → https://www.bnd.bund.de/error_path/400.html?al_req_id=aujcN9DJpXMd5cxKZs2ovwAAACE |
| ch37 | <https://www.cas.go.jp/jp/gaiyou/jimu/nic_nib.html> | HEAD: HTTP 200 | — |
| ch37 | <https://www.chapsvision.com/> | HEAD: HTTP 200 | — |
| ch37 | <https://www.chinalawtranslate.com/en/national-intelligence-law-of-the-p-r-c-2017/> | HEAD: HTTP 200 | — |
| ch37 | <https://www.cia.gov/resources/csi/books-monographs/a-cold-war-conundrum/> | HEAD: HTTP 200 | — |
| ch37 | <https://www.cia.gov/resources/csi/books-monographs/psychology-of-intelligence-analysis-2/> | HEAD: HTTP 200 | — |
| ch32 | <https://www.cidrdb.org/cidr2022/papers/p13-crotty.pdf> | HEAD: HTTP 200 | — |
| ch30 | <https://www.cl.cam.ac.uk/~rja14/book.html> | HEAD: HTTP 200 | 301 → https://www.cl.cam.ac.uk/archive/rja14/book.html |
| ch17 | <https://www.clipsrules.net/> | HEAD: HTTP 200 | — |
| ch27, ch30 | <https://www.csl.sri.com/users/rushby/abstracts/sss10> | HEAD: HTTP 200 | — |
| ch30 | <https://www.csl.sri.com/users/rushby/papers/Understanding_Assurance_2.0.pdf> | GET (HEAD 404 → ranged GET): HTTP 404 | — |
| ch27 | <https://www.dbai.tuwien.ac.at/research/argumentation/aspartix/> | HEAD: HTTP 200 | — |
| ch37 | <https://www.defense.gouv.fr/drm/nos-missions> | HEAD: HTTP 200 | — |
| ch22 | <https://www.do-mpc.com/en/latest/> | HEAD: HTTP 200 | — |
| ch37 | <https://www.esri.com/en-us/arcgis/products/arcgis-pro/overview> | HEAD: HTTP 200 | — |
| ch37 | <https://www.europol.europa.eu/operations-services-innovation/europol-analysis-projects> | HEAD: HTTP 200 | 301 → https://www.europol.europa.eu/how-we-work/europol-analysis-projects |
| ch04 | <https://www.faa.gov/aircraft/air_cert/step/roadmap_for_AI_safety_assurance> | HEAD: HTTP 200 | — |
| ch37 | <https://www.fbi.gov/how-we-investigate/intelligence> | HEAD: HTTP 200 | — |
| ch01, ch04 | <https://www.fda.gov/regulatory-information/search-fda-guidance-documents/clinical-decision-support-software> | HEAD: HTTP 200 | — |
| ch33 | <http://www.fipa.org/specs/fipa00061/SC00061G.html> | GET (HEAD 403 → ranged GET): HTTP 403 | — |
| ch37 | <https://www.first.org/cvss/v4.0/specification-document> | HEAD: HTTP 200 | — |
| ch37 | <https://www.first.org/standards/frameworks/csirts/csirt_services_framework_v2.1> | HEAD: HTTP 200 | — |
| ch37 | <https://www.first.org/tlp/> | HEAD: HTTP 200 | — |
| ch37 | <https://www.gov.uk/government/publications/explaining-uncertainty-in-uk-intelligence-assessment/explaining-uncertainty-in-uk-intelligence-assessment> | HEAD: HTTP 200 | — |
| ch04 | <https://www.gov.uk/government/publications/jsp-936-dependable-artificial-intelligence-ai-in-defence-part-1-directive> | HEAD: HTTP 200 | — |
| ch27 | <https://www.gov.uk/government/publications/the-nimrod-review> | HEAD: HTTP 200 | — |
| ch34 | <https://www.hup.harvard.edu/books/9780674138001> | HEAD: HTTP 202 | — |
| ch04 | <https://www.iec.ch/functional-safety> | GET (HEAD 403 → ranged GET): HTTP 403 | — |
| ch04 | <http://www.incompleteideas.net/IncIdeas/BitterLesson.html> | HEAD: HTTP 200 | — |
| ch18 | <https://www.intel.com/content/www/us/en/developer/articles/technical/intel-resource-director-technology-architecture.html> | HEAD: HTTP 200 | 301 → https://corpredirect.intel.com/Redirector/404Redirector.aspx?https://www.intel.com/content/www/us/en/developer/articles/technical/intel-resource-director-technology-architecture.html<br>301 → https://www.intel.com/content/www/us/en/developer/overview.html |
| ch18 | <https://www.intel.com/content/www/us/en/newsroom/news/artificial-intelligence/ces-2026-intel-core-ultra-series-3-debut-first-built-on-intel-18a.html> | HEAD: HTTP 200 | — |
| ch18 | <https://www.intel.com/content/www/us/en/newsroom/news/artificial-intelligence/intel-to-expand-ai-accelerator-portfolio-with-new-gpu.html> | HEAD: HTTP 200 | — |
| ch37 | <https://www.interpol.int/How-we-work/Criminal-intelligence-analysis> | HEAD: HTTP 200 | — |
| ch02, ch14 | <https://www.iso.org/sites/directives/current/part2/index.xhtml> | HEAD: HTTP 200 | — |
| ch05, ch18 | <https://www.iso.org/standard/27001> | GET (HEAD 403 → ranged GET): HTTP 403 | — |
| ch15 | <https://www.iso.org/standard/31898.html> | GET (HEAD 403 → ranged GET): HTTP 403 | — |
| ch08 | <https://www.iso.org/standard/64599.html> | GET (HEAD 403 → ranged GET): HTTP 403 | — |
| ch01, ch02, ch14, ch15, ch18, ch30, ch31, ch33, ch34 | <https://www.iso.org/standard/68383.html> | GET (HEAD 403 → ranged GET): HTTP 403 | — |
| ch09 | <https://www.iso.org/standard/68386.html> | GET (HEAD 403 → ranged GET): HTTP 200 | — |
| ch04 | <https://www.iso.org/standard/68387.html> | GET (HEAD 403 → ranged GET): HTTP 403 | — |
| ch03, ch05, ch25, ch29 | <https://www.iso.org/standard/68388.html> | GET (HEAD 403 → ranged GET): HTTP 200 | — |
| ch01, ch14, ch30 | <https://www.iso.org/standard/68390.html> | GET (HEAD 403 → ranged GET): HTTP 403 | — |
| ch25 | <https://www.iso.org/standard/68391.html> | GET (HEAD 403 → ranged GET): HTTP 200 | — |
| ch03, ch04, ch05, ch25, ch30 | <https://www.iso.org/standard/70918.html> | GET (HEAD 403 → ranged GET): HTTP 403 | — |
| ch08, ch14, ch25 | <https://www.iso.org/standard/72089.html> | GET (HEAD 403 → ranged GET): HTTP 403 | — |
| ch17 | <https://www.iso.org/standard/76120.html> | GET (HEAD 403 → ranged GET): HTTP 403 | — |
| ch23 | <https://www.iso.org/standard/79016.html> | GET (HEAD 403 → ranged GET): HTTP 403 | — |
| ch23 | <https://www.iso.org/standard/81291.html> | GET (HEAD 403 → ranged GET): HTTP 403 | — |
| ch06, ch31 | <https://www.iti.zcu.cz/wl2018/pdf/wl_paper_translation.pdf> | HEAD: HTTP 200 | — |
| ch24 | <https://www.jmlr.org/papers/v11/el-yaniv10a.html> | HEAD: HTTP 200 | — |
| ch24 | <https://www.jmlr.org/papers/v25/23-0487.html> | HEAD: HTTP 200 | — |
| ch24 | <https://www.jmlr.org/papers/v25/23-0970.html> | HEAD: HTTP 200 | — |
| ch37 | <https://www.maltego.com/> | HEAD: HTTP 200 | — |
| ch18 | <https://www.microsoft.com/en-us/surface/devices/surface-rtx-spark-dev-box> | HEAD: HTTP 200 | — |
| ch37 | <https://www.misp-project.org/> | HEAD: HTTP 200 | — |
| ch36 | <https://www.mountaingoatsoftware.com/books/succeeding-with-agile-software-development-using-scrum> | HEAD: HTTP 200 | — |
| ch04 | <https://www.nato.int/en/about-us/official-texts-and-resources/official-texts/2024/07/10/summary-of-natos-revised-artificial-intelligence-ai-strategy> | HEAD: HTTP 200 | — |
| ch18 | <https://www.nvidia.com/en-us/autonomous-machines/embedded-systems/jetson-orin/> | HEAD: HTTP 200 | — |
| ch18 | <https://www.nvidia.com/en-us/autonomous-machines/embedded-systems/jetson-thor/> | HEAD: HTTP 200 | — |
| ch18 | <https://www.nvidia.com/en-us/products/rtx-spark/> | HEAD: HTTP 200 | — |
| ch18 | <https://www.nvidia.com/en-us/products/workstations/professional-desktop-gpus/rtx-pro-6000/> | HEAD: HTTP 200 | — |
| ch22 | <https://www.o-ran.org/specifications> | HEAD: HTTP 200 | — |
| ch17 | <https://www.omg.org/spec/DMN/1.5/About-DMN> | HEAD: HTTP 200 | — |
| ch10 | <https://www.omg.org/spec/ReqIF/1.2> | HEAD: HTTP 200 | — |
| ch08, ch14, ch30 | <https://www.omg.org/spec/ReqIF/1.2/About-ReqIF> | HEAD: HTTP 200 | — |
| ch27 | <https://www.omg.org/spec/SACM/2.3> | HEAD: HTTP 200 | — |
| ch10 | <https://www.omg.org/spec/SysML/> | HEAD: HTTP 200 | — |
| ch16 | <https://www.openpolicyagent.org/docs/latest/policy-language/> | HEAD: HTTP 200 | 301 → https://www.openpolicyagent.org/docs/policy-language/<br>301 → https://www.openpolicyagent.org/docs/policy-language |
| ch21 | <https://www.openpolicyagent.org/docs/policy-language> | HEAD: HTTP 200 | — |
| ch32 | <https://www.os-book.com/> | HEAD: HTTP 200 | — |
| ch37 | <https://www.palantir.com/platforms/foundry/data-integration/> | HEAD: HTTP 200 | — |
| ch37 | <https://www.palantir.com/platforms/gotham/> | HEAD: HTTP 200 | — |
| ch03 | <https://www.penguin.co.uk/books/55823/the-beginning-of-infinity-by-david-deutsch/9780140278163> | HEAD: HTTP 200 | 308 → https://www.penguin.co.uk/books/24745/the-beginning-of-infinity-by-david-deutsch/9780140278163 |
| ch07 | <https://www.postgresql.org/docs/current/ddl-partitioning.html> | HEAD: HTTP 200 | — |
| ch18 | <https://www.qualcomm.com/laptops/products/snapdragon-x2-elite> | HEAD: HTTP 200 | — |
| ch18 | <https://www.qualcomm.com/news/releases/2025/09/new-snapdragon-x2-elite-extreme-and-snapdragon-x2-elite-are-the-> | HEAD: HTTP 200 | — |
| ch37 | <https://www.recordedfuture.com/> | HEAD: HTTP 200 | — |
| ch15 | <https://www.rfc-editor.org/errata/rfc5321> | GET (HEAD 405 → ranged GET): HTTP 200 | 302 → https://errata.rfc-editor.org/rfc5321<br>301 → https://errata.rfc-editor.org/rfc5321/<br>302 → https://errata.rfc-editor.org/search/?rfc_number=5321&presentation=records |
| ch08 | <https://www.rfc-editor.org/info/rfc8259/> | HEAD: HTTP 200 | — |
| ch24 | <https://www.rfc-editor.org/info/rfc9293> | HEAD: HTTP 200 | 301 → https://www.rfc-editor.org/info/rfc9293/ |
| ch32 | <https://www.rfc-editor.org/rfc/rfc1122> | HEAD: HTTP 200 | 301 → https://www.rfc-editor.org/info/rfc1122/ |
| ch06, ch14, ch31 | <https://www.rfc-editor.org/rfc/rfc2119> | HEAD: HTTP 200 | 301 → https://www.rfc-editor.org/info/rfc2119/ |
| ch15 | <https://www.rfc-editor.org/rfc/rfc2821> | HEAD: HTTP 200 | 301 → https://www.rfc-editor.org/info/rfc2821/ |
| ch05 | <https://www.rfc-editor.org/rfc/rfc3161> | HEAD: HTTP 200 | 301 → https://www.rfc-editor.org/info/rfc3161/ |
| ch13 | <https://www.rfc-editor.org/rfc/rfc3207> | HEAD: HTTP 200 | 301 → https://www.rfc-editor.org/info/rfc3207/ |
| ch29 | <https://www.rfc-editor.org/rfc/rfc3629> | HEAD: HTTP 200 | 301 → https://www.rfc-editor.org/info/rfc3629/ |
| ch15, ch32 | <https://www.rfc-editor.org/rfc/rfc5234> | HEAD: HTTP 200 | 301 → https://www.rfc-editor.org/info/rfc5234/ |
| ch13, ch15, ch20, ch28, ch31 | <https://www.rfc-editor.org/rfc/rfc5321> | HEAD: HTTP 200 | 301 → https://www.rfc-editor.org/info/rfc5321/ |
| ch31 | <https://www.rfc-editor.org/rfc/rfc5961> | HEAD: HTTP 200 | 301 → https://www.rfc-editor.org/info/rfc5961/ |
| ch28 | <https://www.rfc-editor.org/rfc/rfc6409> | HEAD: HTTP 200 | 301 → https://www.rfc-editor.org/info/rfc6409/ |
| ch27 | <https://www.rfc-editor.org/rfc/rfc6962> | HEAD: HTTP 200 | 301 → https://www.rfc-editor.org/info/rfc6962/ |
| ch15 | <https://www.rfc-editor.org/rfc/rfc7504> | HEAD: HTTP 200 | 301 → https://www.rfc-editor.org/info/rfc7504/ |
| ch13, ch19 | <https://www.rfc-editor.org/rfc/rfc768> | HEAD: HTTP 200 | 301 → https://www.rfc-editor.org/info/rfc768/ |
| ch31, ch32 | <https://www.rfc-editor.org/rfc/rfc793> | HEAD: HTTP 200 | 301 → https://www.rfc-editor.org/info/rfc793/ |
| ch37 | <https://www.rfc-editor.org/rfc/rfc7970.html> | HEAD: HTTP 200 | — |
| ch15 | <https://www.rfc-editor.org/rfc/rfc8032> | HEAD: HTTP 200 | 301 → https://www.rfc-editor.org/info/rfc8032/ |
| ch14, ch31 | <https://www.rfc-editor.org/rfc/rfc8174> | HEAD: HTTP 200 | 301 → https://www.rfc-editor.org/info/rfc8174/ |
| ch15 | <https://www.rfc-editor.org/rfc/rfc821> | HEAD: HTTP 200 | 301 → https://www.rfc-editor.org/info/rfc821/ |
| ch32 | <https://www.rfc-editor.org/rfc/rfc8554> | HEAD: HTTP 200 | 301 → https://www.rfc-editor.org/info/rfc8554/ |
| ch21, ch27, ch32 | <https://www.rfc-editor.org/rfc/rfc8785> | HEAD: HTTP 200 | 301 → https://www.rfc-editor.org/info/rfc8785/ |
| ch21 | <https://www.rfc-editor.org/rfc/rfc9110> | HEAD: HTTP 200 | 301 → https://www.rfc-editor.org/info/rfc9110/ |
| ch27 | <https://www.rfc-editor.org/rfc/rfc9162> | HEAD: HTTP 200 | 301 → https://www.rfc-editor.org/info/rfc9162/ |
| ch31, ch32 | <https://www.rfc-editor.org/rfc/rfc9293> | HEAD: HTTP 200 | 301 → https://www.rfc-editor.org/info/rfc9293/ |
| ch22 | <https://www.rfc-editor.org/rfc/rfc9334> | HEAD: HTTP 200 | 301 → https://www.rfc-editor.org/info/rfc9334/ |
| ch02 | <https://www.routledge.com/The-Logic-of-Scientific-Discovery/Popper/p/book/9780415278447> | HEAD: HTTP 200 | — |
| ch01, ch04, ch22, ch30, ch33 | <https://www.rtca.org/> | HEAD: HTTP 200 | — |
| ch09, ch14, ch15 | <https://www.rtca.org/do-178/> | HEAD: HTTP 200 | — |
| ch29 | <https://www.rtca.org/products/do-178c/> | HEAD: HTTP 200 | 301 → https://www.rtca.org/products/do-178c-electronic/ |
| ch01 | <https://www.shortliffe.net/Buchanan-Shortliffe-1984/MYCIN%20Book.htm> | HEAD: HTTP 200 | — |
| ch18 | <https://www.sqlite.org/fts5.html> | HEAD: HTTP 200 | — |
| ch17 | <https://www.sqlite.org/json1.html> | HEAD: HTTP 200 | — |
| ch17 | <https://www.sqlite.org/lang_with.html> | HEAD: HTTP 200 | — |
| ch18 | <https://www.ti.com/product/TMS570LC4357> | HEAD: HTTP 200 | — |
| ch29 | <https://www.tomzahavy.com/files/llms-cant-jump.pdf> | HEAD: HTTP 200 | — |
| ch12, ch32 | <https://www.unicode.org/reports/tr15/> | HEAD: HTTP 200 | — |
| ch12 | <https://www.unicode.org/reports/tr29/> | HEAD: HTTP 200 | — |
| ch12 | <https://www.unicode.org/reports/tr39/> | HEAD: HTTP 200 | — |
| ch12 | <https://www.unicode.org/reports/tr9/> | HEAD: HTTP 200 | — |
| ch37 | <https://www.unodc.org/documents/organized-crime/Law-Enforcement/Criminal_Intelligence_for_Analysts.pdf> | HEAD: HTTP 200 | — |
| ch33 | <https://www.usenix.org/conference/usenixsecurity19/presentation/torres-arias> | HEAD: HTTP 200 | — |
| ch25 | <https://www.usenix.org/conference/usenixsecurity21/presentation/carlini-extracting> | HEAD: HTTP 200 | — |
| ch12 | <https://www.usenix.org/conference/usenixsecurity23/presentation/boucher> | HEAD: HTTP 200 | — |
| ch37 | <https://www.vldb.org/pvldb/vol8/p1792-akidau.pdf> | HEAD: HTTP 200 | — |
| ch32 | <https://www.w3.org/TR/2003/REC-PNG-20031110/> | HEAD: HTTP 200 | — |
| ch33 | <https://www.w3.org/TR/json-ld11/> | HEAD: HTTP 200 | — |
| ch02 | <https://www.w3.org/TR/odrl-model/> | HEAD: HTTP 200 | — |
| ch02 | <https://www.w3.org/TR/owl-time/> | HEAD: HTTP 200 | — |
| ch07, ch16, ch23 | <https://www.w3.org/TR/owl2-overview/> | HEAD: HTTP 200 | — |
| ch02 | <https://www.w3.org/TR/owl2-primer/> | HEAD: HTTP 200 | — |
| ch07 | <https://www.w3.org/TR/owl2-profiles/> | HEAD: HTTP 200 | — |
| ch02, ch03, ch05, ch08, ch09, ch10, ch16, ch19 | <https://www.w3.org/TR/prov-o/> | HEAD: HTTP 200 | — |
| ch25 | <https://www.w3.org/TR/prov-overview/> | HEAD: HTTP 200 | — |
| ch06 | <https://www.w3.org/TR/rdf-schema/> | HEAD: HTTP 200 | — |
| ch04, ch07 | <https://www.w3.org/TR/rdf11-concepts/> | HEAD: HTTP 200 | — |
| ch31 | <https://www.w3.org/TR/rdf11-mt/> | HEAD: HTTP 200 | — |
| ch07 | <https://www.w3.org/TR/rdf12-concepts/> | HEAD: HTTP 200 | — |
| ch02, ch06, ch07, ch08, ch16, ch23, ch33 | <https://www.w3.org/TR/shacl/> | HEAD: HTTP 200 | — |
| ch07 | <https://www.w3.org/TR/shacl12-core/> | HEAD: HTTP 200 | — |
| ch07 | <https://www.w3.org/TR/sparql11-overview/> | HEAD: HTTP 200 | — |
| ch16, ch17, ch31 | <https://www.w3.org/TR/sparql11-query/> | HEAD: HTTP 200 | — |
| ch07 | <https://www.w3.org/TR/sparql12-query/> | HEAD: HTTP 200 | — |
| ch16 | <https://www.w3.org/TR/turtle/> | HEAD: HTTP 200 | — |
| ch10 | <https://www.wiley.com/> | HEAD: HTTP 200 | 302 → https://www.wiley.com/en-ie/<br>301 → https://www.wiley.com/en-ie |
| ch08 | <https://yaml.org/spec/1.2.2/> | HEAD: HTTP 200 | — |
| ch34 | <https://zakon.rada.gov.ua/laws/show/435-15> | HEAD: HTTP 200 | — |
| ch37 | <https://zakon.rada.gov.ua/laws/show/4651-17#Text> | HEAD: HTTP 200 | — |
| ch37 | <https://zakon.rada.gov.ua/laws/show/912-20#Text> | HEAD: HTTP 200 | — |

## Broken, blocked, and inconclusive links

Наведено лише цілі, які завершилися не-2xx статусом або мережевою помилкою; 403/429 та тайм-аут не прирівнюються до остаточно видаленого документа.

| URL | Статус / помилка | Кінцева адреса | Глави |
|---|---|---|---|
| <http://www.fipa.org/specs/fipa00061/SC00061G.html> | 403 | http://www.fipa.org/specs/fipa00061/SC00061G.html | ch33 |
| <https://apps.dtic.mil/sti/citations/ADA039324> | 403 | https://apps.dtic.mil/sti/citations/ADA039324 | ch12 |
| <https://doi.org/10.1001/jama.1983.03330370053031> | 403 | https://jamanetwork.com/article.aspx?doi=10.1001/jama.1983.03330370053031 | ch25 |
| <https://doi.org/10.1002/0471221279> | 403 | https://onlinelibrary.wiley.com/doi/book/10.1002/0471221279 | ch06, ch33 |
| <https://doi.org/10.1002/047174882X> | 403 | https://onlinelibrary.wiley.com/doi/book/10.1002/047174882X | ch24 |
| <https://doi.org/10.1002/9780470316887> | 403 | https://onlinelibrary.wiley.com/doi/book/10.1002/9780470316887 | ch06 |
| <https://doi.org/10.1002/asi.20591> | 403 | https://onlinelibrary.wiley.com/doi/10.1002/asi.20591 | ch09 |
| <https://doi.org/10.1002/for.3980010305> | 403 | https://onlinelibrary.wiley.com/doi/10.1002/for.3980010305 | ch11 |
| <https://doi.org/10.1002/int.4550090804> | 403 | https://onlinelibrary.wiley.com/doi/10.1002/int.4550090804 | ch23 |
| <https://doi.org/10.1007/s11225-007-9067-1> | 404 | https://doi.org/10.1007/s11225-007-9067-1 | ch34 |
| <https://doi.org/10.1051/ro/196802v100571> | 403 | https://www.rairo-ro.org/10.1051/ro/196802V100571 | ch06 |
| <https://doi.org/10.1057/jors.1984.2> | 403 | https://www.tandfonline.com/doi/full/10.1057/jors.1984.2 | ch22 |
| <https://doi.org/10.1073/pnas.1611835114> | 403 | https://pnas.org/doi/full/10.1073/pnas.1611835114 | ch26 |
| <https://doi.org/10.1073/pnas.1708274114> | 403 | https://pnas.org/doi/full/10.1073/pnas.1708274114 | ch36 |
| <https://doi.org/10.1073/pnas.36.1.48> | 403 | https://pnas.org/doi/full/10.1073/pnas.36.1.48 | ch06 |
| <https://doi.org/10.1080/00207727008920220> | 403 | https://www.tandfonline.com/doi/abs/10.1080/00207727008920220 | ch22 |
| <https://doi.org/10.1080/01621459.1952.10483446> | 403 | https://www.tandfonline.com/doi/full/10.1080/01621459.1952.10483446 | ch10 |
| <https://doi.org/10.1080/13518040490450529> | 403 | https://www.tandfonline.com/doi/abs/10.1080/13518040490450529 | ch37 |
| <https://doi.org/10.1080/19462166.2013.869766> | 403 | https://journals.sagepub.com/doi/10.1080/19462166.2013.869766 | ch29, ch31 |
| <https://doi.org/10.1093/acprof:oso/9780199232383.001.0001> | 403 | https://academic.oup.com/book/32518 | ch02 |
| <https://doi.org/10.1093/analys/23.6.121> | 403 | https://academic.oup.com/analysis/article-lookup/doi/10.1093/analys/23.6.121 | ch02 |
| <https://doi.org/10.1093/biomet/41.1-2.100> | 403 | https://academic.oup.com/biomet/article-lookup/doi/10.1093/biomet/41.1-2.100 | ch26, ch33 |
| <https://doi.org/10.1093/comjnl/25.4.465> | 403 | https://academic.oup.com/comjnl/article-lookup/doi/10.1093/comjnl/25.4.465 | ch23 |
| <https://doi.org/10.1093/logcom/2.6.719> | 403 | https://academic.oup.com/logcom/article-lookup/doi/10.1093/logcom/2.6.719 | ch34 |
| <https://doi.org/10.1093/mind/LX.237.1> | 403 | https://academic.oup.com/mind/article-lookup/doi/10.1093/mind/LX.237.1 | ch06, ch28 |
| <https://doi.org/10.1093/oso/9780195092691.001.0001> | 403 | https://academic.oup.com/book/52097 | ch05 |
| <https://doi.org/10.1098/rstl.1763.0053> | 403 | https://royalsocietypublishing.org/rstl/article/doi/10.1098/rstl.1763.0053/119736/LII-An-essay-towards-solving-a-problem-in-the | ch04 |
| <https://doi.org/10.1111/0162-895X.00317> | 403 | https://onlinelibrary.wiley.com/doi/10.1111/0162-895X.00317 | ch37 |
| <https://doi.org/10.1111/1468-0394.00045> | 403 | https://onlinelibrary.wiley.com/doi/10.1111/1468-0394.00045 | ch11 |
| <https://doi.org/10.1111/j.1469-8137.1912.tb05611.x> | 403 | https://nph.onlinelibrary.wiley.com/doi/10.1111/j.1469-8137.1912.tb05611.x | ch28 |
| <https://doi.org/10.1111/j.2517-6161.1951.tb00088.x> | 403 | https://academic.oup.com/jrsssb/article/13/2/238/7026675 | ch06 |
| <https://doi.org/10.1112/plms/s2-42.1.230> | 403 | https://onlinelibrary.wiley.com/doi/10.1112/plms/s2-42.1.230 | ch06 |
| <https://doi.org/10.1115/1.3662552> | 403 | https://asmedigitalcollection.asme.org/fluidsengineering/article/82/1/35/397706/A-New-Approach-to-Linear-Filtering-and-Prediction | ch06, ch22 |
| <https://doi.org/10.1126/science.217.4563.927> | 403 | https://www.science.org/doi/10.1126/science.217.4563.927 | ch04 |
| <https://doi.org/10.1126/science.ade9097> | 403 | https://www.science.org/doi/10.1126/science.ade9097 | ch29 |
| <https://doi.org/10.1137/0201010> | 403 | https://epubs.siam.org/doi/10.1137/0201010 | ch31 |
| <https://doi.org/10.1137/0218012> | 403 | https://epubs.siam.org/doi/10.1137/0218012 | ch10 |
| <https://doi.org/10.1137/090752286> | 403 | https://epubs.siam.org/doi/10.1137/090752286 | ch32 |
| <https://doi.org/10.1137/1.9780898718003> | 403 | https://epubs.siam.org/doi/book/10.1137/1.9780898718003 | ch32 |
| <https://doi.org/10.1137/1.9781611972771.42> | 403 | https://epubs.siam.org/doi/10.1137/1.9781611972771.42 | ch26 |
| <https://doi.org/10.1137/S1064827595287997> | 403 | https://epubs.siam.org/doi/10.1137/S1064827595287997 | ch07 |
| <https://doi.org/10.1145/103162.103163> | 403 | https://dl.acm.org/doi/10.1145/103162.103163 | ch18 |
| <https://doi.org/10.1145/129888.129894> | 403 | https://dl.acm.org/doi/10.1145/129888.129894 | ch07 |
| <https://doi.org/10.1145/1347375.1347389> | 403 | https://dl.acm.org/doi/10.1145/1347375.1347389 | ch18 |
| <https://doi.org/10.1145/143062.143090> | 403 | https://dl.acm.org/doi/10.1145/143062.143090 | ch32 |
| <https://doi.org/10.1145/1498765.1498785> | 403 | https://dl.acm.org/doi/10.1145/1498765.1498785 | ch18 |
| <https://doi.org/10.1145/1555349.1555372> | 403 | https://dl.acm.org/doi/10.1145/1555349.1555372 | ch17 |
| <https://doi.org/10.1145/1571941.1572114> | 403 | https://dl.acm.org/doi/10.1145/1571941.1572114 | ch06, ch07, ch08, ch12, ch13, ch15 |
| <https://doi.org/10.1145/1866307.1866315> | 403 | https://dl.acm.org/doi/10.1145/1866307.1866315 | ch22 |
| <https://doi.org/10.1145/191246.191322> | 403 | https://portal.acm.org/citation.cfm?doid=191246.191322 | ch33 |
| <https://doi.org/10.1145/203330.203341> | 403 | https://dl.acm.org/doi/10.1145/203330.203341 | ch24 |
| <https://doi.org/10.1145/219717.219748> | 403 | https://dl.acm.org/doi/10.1145/219717.219748 | ch12 |
| <https://doi.org/10.1145/2342509.2342513> | 403 | https://dl.acm.org/doi/10.1145/2342509.2342513 | ch18 |
| <https://doi.org/10.1145/235968.233333> | 403 | https://dl.acm.org/doi/10.1145/235968.233333 | ch32 |
| <https://doi.org/10.1145/2488388.2488425> | 403 | https://dl.acm.org/doi/10.1145/2488388.2488425 | ch31, ch34 |
| <https://doi.org/10.1145/2523813> | 403 | https://dl.acm.org/doi/10.1145/2523813 | ch26 |
| <https://doi.org/10.1145/2556938> | 403 | https://dl.acm.org/doi/10.1145/2556938 | ch30 |
| <https://doi.org/10.1145/258533.258660> | 403 | https://portal.acm.org/citation.cfm?doid=258533.258660 | ch32 |
| <https://doi.org/10.1145/2593882.2593891> | 403 | https://dl.acm.org/doi/10.1145/2593882.2593891 | ch09 |
| <https://doi.org/10.1145/2623330.2623623> | 403 | https://dl.acm.org/doi/10.1145/2623330.2623623 | ch09 |
| <https://doi.org/10.1145/2629489> | 403 | https://dl.acm.org/doi/10.1145/2629489 | ch12 |
| <https://doi.org/10.1145/263699.263712> | 403 | https://portal.acm.org/citation.cfm?doid=263699.263712 | ch32 |
| <https://doi.org/10.1145/2749469.2750386> | 403 | https://dl.acm.org/doi/10.1145/2749469.2750386 | ch18 |
| <https://doi.org/10.1145/2939672.2939778> | 403 | https://dl.acm.org/doi/10.1145/2939672.2939778 | ch06 |
| <https://doi.org/10.1145/302979.303030> | 403 | https://portal.acm.org/citation.cfm?doid=302979.303030 | ch34 |
| <https://doi.org/10.1145/3143561> | 403 | https://dl.acm.org/doi/10.1145/3143561 | ch23 |
| <https://doi.org/10.1145/317087.317089> | 403 | https://dl.acm.org/doi/10.1145/317087.317089 | ch02 |
| <https://doi.org/10.1145/318898.318921> | 403 | https://dl.acm.org/doi/10.1145/318898.318921 | ch02 |
| <https://doi.org/10.1145/321250.321253> | 403 | https://dl.acm.org/doi/10.1145/321250.321253 | ch01, ch04 |
| <https://doi.org/10.1145/321978.321991> | 403 | https://dl.acm.org/doi/10.1145/321978.321991 | ch06 |
| <https://doi.org/10.1145/3287560.3287566> | 403 | https://dl.acm.org/doi/10.1145/3287560.3287566 | ch20 |
| <https://doi.org/10.1145/3287560.3287596> | 403 | https://dl.acm.org/doi/10.1145/3287560.3287596 | ch03, ch05 |
| <https://doi.org/10.1145/3314221.3314642> | 403 | https://dl.acm.org/doi/10.1145/3314221.3314642 | ch37, ch38 |
| <https://doi.org/10.1145/331499.331504> | 403 | https://dl.acm.org/doi/10.1145/331499.331504 | ch07 |
| <https://doi.org/10.1145/3373376.3378519> | 403 | https://dl.acm.org/doi/10.1145/3373376.3378519 | ch36 |
| <https://doi.org/10.1145/3404835.3463098> | 403 | https://dl.acm.org/doi/10.1145/3404835.3463098 | ch13 |
| <https://doi.org/10.1145/3442188.3445922> | 403 | https://dl.acm.org/doi/10.1145/3442188.3445922 | ch38 |
| <https://doi.org/10.1145/3458723> | 403 | https://dl.acm.org/doi/10.1145/3458723 | ch05, ch25 |
| <https://doi.org/10.1145/3473582> | 403 | https://dl.acm.org/doi/10.1145/3473582 | ch31 |
| <https://doi.org/10.1145/351240.351266> | 403 | https://dl.acm.org/doi/10.1145/351240.351266 | ch23 |
| <https://doi.org/10.1145/3571730> | 403 | https://dl.acm.org/doi/10.1145/3571730 | ch03, ch11, ch38 |
| <https://doi.org/10.1145/359168.359176> | 403 | https://dl.acm.org/doi/10.1145/359168.359176 | ch05 |
| <https://doi.org/10.1145/359545.359563> | 403 | https://dl.acm.org/doi/10.1145/359545.359563 | ch22, ch33 |
| <https://doi.org/10.1145/3600006> | 403 | https://dl.acm.org/doi/proceedings/10.1145/3600006 | ch30 |
| <https://doi.org/10.1145/3600006.3613165> | 403 | https://dl.acm.org/doi/10.1145/3600006.3613165 | ch12, ch17 |
| <https://doi.org/10.1145/360051.360056> | 403 | https://dl.acm.org/doi/10.1145/360051.360056 | ch02, ch07, ch33 |
| <https://doi.org/10.1145/3605764.3623985> | 403 | https://dl.acm.org/doi/10.1145/3605764.3623985 | ch12, ch21 |
| <https://doi.org/10.1145/3618260.3649777> | 403 | https://dl.acm.org/doi/10.1145/3618260.3649777 | ch38 |
| <https://doi.org/10.1145/362686.362692> | 403 | https://dl.acm.org/doi/10.1145/362686.362692 | ch32 |
| <https://doi.org/10.1145/38713.38742> | 403 | https://portal.acm.org/citation.cfm?doid=38713.38742 | ch21 |
| <https://doi.org/10.1145/503124.503150> | 403 | https://dl.acm.org/doi/10.1145/503124.503150 | ch07 |
| <https://doi.org/10.1145/582415.582418> | 403 | https://dl.acm.org/doi/10.1145/582415.582418 | ch16 |
| <https://doi.org/10.1145/800157.805047> | 403 | https://portal.acm.org/citation.cfm?doid=800157.805047 | ch04 |
| <https://doi.org/10.1145/800216.806586> | 403 | https://portal.acm.org/citation.cfm?doid=800216.806586 | ch30 |
| <https://doi.org/10.1145/857076.857078> | 403 | https://dl.acm.org/doi/10.1145/857076.857078 | ch33 |
| <https://doi.org/10.1162/089976698300017197> | 403 | https://direct.mit.edu/neco/article/10/7/1895-1923/6224 | ch19 |
| <https://doi.org/10.1162/coli_a_00402> | 403 | https://direct.mit.edu/coli/article/doi/10.1162/coli_a_00402/98516/Universal-Dependencies | ch10, ch12, ch13 |
| <https://doi.org/10.1162/tacl_a_00276> | 403 | https://direct.mit.edu/tacl/article/43518 | ch19 |
| <https://doi.org/10.1162/tacl_a_00298> | 403 | https://direct.mit.edu/tacl/article/43535 | ch13 |
| <https://doi.org/10.1162/tacl_a_00638> | 403 | https://direct.mit.edu/tacl/article/doi/10.1162/tacl_a_00638/119630/Lost-in-the-Middle-How-Language-Models-Use-Long | ch13 |
| <https://doi.org/10.1177/001316446002000104> | 403 | https://journals.sagepub.com/doi/10.1177/001316446002000104 | ch11, ch25 |
| <https://doi.org/10.1177/0018720810376055> | 403 | https://journals.sagepub.com/doi/10.1177/0018720810376055 | ch21 |
| <https://doi.org/10.1177/1049732316654870> | 403 | https://journals.sagepub.com/doi/10.1177/1049732316654870 | ch11 |
| <https://doi.org/10.1207/s15516709cog1104_4> | 403 | https://onlinelibrary.wiley.com/doi/10.1207/s15516709cog1104_4 | ch02 |
| <https://doi.org/10.1287/mnsc.31.6.647> | 403 | https://pubsonline.informs.org/doi/10.1287/mnsc.31.6.647 | ch06 |
| <https://doi.org/10.1287/mnsc.9.3.458> | 403 | https://pubsonline.informs.org/doi/10.1287/mnsc.9.3.458 | ch11 |
| <https://doi.org/10.1287/opre.9.3.383> | 403 | https://pubsonline.informs.org/doi/10.1287/opre.9.3.383 | ch10 |
| <https://doi.org/10.14778/3402707.3402747> | 403 | https://dl.acm.org/doi/10.14778/3402707.3402747 | ch07 |
| <https://doi.org/10.15407/pp2020.02-03.259> | 502 | http://pp.isofts.kiev.ua/ojs1/article/view/417 | ch01, ch12 |
| <https://doi.org/10.1561/1500000019> | 403 | https://www.emerald.com/ftinr/article/4/1-2/1/1326508/The-Probabilistic-Relevance-Framework-BM25-and | ch03, ch04, ch06, ch08, ch13, ch15, ch16 |
| <https://doi.org/10.1561/1500000024> | 403 | https://www.emerald.com/ftinr/article/6/2-3/127/1356738/Expertise-Retrieval | ch11 |
| <https://doi.org/10.21236/AD0770768> | 404 | https://doi.org/10.21236/AD0770768 | ch02 |
| <https://doi.org/10.21236/ADA046289> | 403 | https://apps.dtic.mil/sti/citations/tr/ADA046289 | ch01, ch04 |
| <https://doi.org/10.2139/ssrn.3063289> | 403 | https://papers.ssrn.com/sol3/papers.cfm?abstract_id=3063289 | ch06 |
| <https://doi.org/10.2307/3250961> | 403 | https://misq.umn.edu/misq/article/25/1/107/1295/Review-Knowledge-Management-and-Knowledge | ch25 |
| <https://doi.org/10.3233/AIC-1994-7104> | 403 | https://journals.sagepub.com/doi/full/10.3233/AIC-1994-7104 | ch04, ch06, ch07, ch24, ch26, ch28 |
| <https://doi.org/10.7551/mitpress/11810.001.0001> | 403 | https://direct.mit.edu/books/book/4581/Cybernetics-or-Control-and-Communication-in-the | ch22 |
| <https://github.com/exploitgym> | 404 | https://github.com/exploitgym | ch30 |
| <https://ilasp.com/> | NETWORK_ERROR (SSLCertVerificationError: [SSL: CERTIFICATE_VERIFY_FAILED] certificate verify failed: unable to get local issuer certificate (_ssl.c:1081)) | https://ilasp.com/ | ch23 |
| <https://openlibrary.org/books/OL5615707M> | NETWORK_ERROR (ConnectionRefusedError: [Errno 111] Connection refused) | https://openlibrary.org/books/OL5615707M | ch04 |
| <https://openlibrary.org/works/OL117061W> | NETWORK_ERROR (ConnectionRefusedError: [Errno 111] Connection refused) | https://openlibrary.org/works/OL117061W | ch11 |
| <https://openlibrary.org/works/OL1311228W> | NETWORK_ERROR (ConnectionRefusedError: [Errno 111] Connection refused) | https://openlibrary.org/works/OL1311228W | ch06 |
| <https://openlibrary.org/works/OL16585333W> | NETWORK_ERROR (ConnectionRefusedError: [Errno 111] Connection refused) | https://openlibrary.org/works/OL16585333W | ch09 |
| <https://openlibrary.org/works/OL17872278W> | NETWORK_ERROR (ConnectionRefusedError: [Errno 111] Connection refused) | https://openlibrary.org/works/OL17872278W | ch06 |
| <https://openlibrary.org/works/OL1814674W> | NETWORK_ERROR (ConnectionRefusedError: [Errno 111] Connection refused) | https://openlibrary.org/works/OL1814674W | ch03 |
| <https://openlibrary.org/works/OL21407442W> | NETWORK_ERROR (ConnectionRefusedError: [Errno 111] Connection refused) | https://openlibrary.org/works/OL21407442W | ch15 |
| <https://openlibrary.org/works/OL24825917W> | NETWORK_ERROR (ConnectionRefusedError: [Errno 111] Connection refused) | https://openlibrary.org/works/OL24825917W | ch04 |
| <https://openlibrary.org/works/OL3307771W> | NETWORK_ERROR (ConnectionRefusedError: [Errno 111] Connection refused) | https://openlibrary.org/works/OL3307771W | ch15 |
| <https://openlibrary.org/works/OL3515903W> | NETWORK_ERROR (ConnectionRefusedError: [Errno 111] Connection refused) | https://openlibrary.org/works/OL3515903W | ch11 |
| <https://openlibrary.org/works/OL4305747W> | NETWORK_ERROR (ConnectionRefusedError: [Errno 111] Connection refused) | https://openlibrary.org/works/OL4305747W | ch11 |
| <https://openlibrary.org/works/OL5282413W> | NETWORK_ERROR (ConnectionRefusedError: [Errno 111] Connection refused) | https://openlibrary.org/works/OL5282413W | ch11 |
| <https://openlibrary.org/works/OL5959470W> | NETWORK_ERROR (ConnectionRefusedError: [Errno 111] Connection refused) | https://openlibrary.org/verify_human?next=/works/OL5959470W | ch02, ch03, ch06, ch31 |
| <https://openlibrary.org/works/OL6408872W> | NETWORK_ERROR (ConnectionRefusedError: [Errno 111] Connection refused) | https://openlibrary.org/works/OL6408872W | ch02 |
| <https://openlibrary.org/works/OL9203874W> | NETWORK_ERROR (ConnectionRefusedError: [Errno 111] Connection refused) | https://openlibrary.org/works/OL9203874W | ch09 |
| <https://www.bnd.bund.de/DE/Die_Arbeit/Analyse/analyse_node.html> | 400 | https://www.bnd.bund.de/error_path/400.html?al_req_id=aujcN9DJpXMd5cxKZs2ovwAAACE | ch37 |
| <https://www.csl.sri.com/users/rushby/papers/Understanding_Assurance_2.0.pdf> | 404 | https://www.csl.sri.com/users/rushby/papers/Understanding_Assurance_2.0.pdf | ch30 |
| <https://www.iec.ch/functional-safety> | 403 | https://www.iec.ch/functional-safety | ch04 |
| <https://www.iso.org/standard/27001> | 403 | https://www.iso.org/standard/27001 | ch05, ch18 |
| <https://www.iso.org/standard/31898.html> | 403 | https://www.iso.org/standard/31898.html | ch15 |
| <https://www.iso.org/standard/64599.html> | 403 | https://www.iso.org/standard/64599.html | ch08 |
| <https://www.iso.org/standard/68383.html> | 403 | https://www.iso.org/standard/68383.html | ch01, ch02, ch14, ch15, ch18, ch30, ch31, ch33, ch34 |
| <https://www.iso.org/standard/68387.html> | 403 | https://www.iso.org/standard/68387.html | ch04 |
| <https://www.iso.org/standard/68390.html> | 403 | https://www.iso.org/standard/68390.html | ch01, ch14, ch30 |
| <https://www.iso.org/standard/70918.html> | 403 | https://www.iso.org/standard/70918.html | ch03, ch04, ch05, ch25, ch30 |
| <https://www.iso.org/standard/72089.html> | 403 | https://www.iso.org/standard/72089.html | ch08, ch14, ch25 |
| <https://www.iso.org/standard/76120.html> | 403 | https://www.iso.org/standard/76120.html | ch17 |
| <https://www.iso.org/standard/79016.html> | 403 | https://www.iso.org/standard/79016.html | ch23 |
| <https://www.iso.org/standard/81291.html> | 403 | https://www.iso.org/standard/81291.html | ch23 |
