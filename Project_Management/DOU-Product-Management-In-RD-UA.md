# Продуктовий менеджмент в R&D: де ринок зустрічається з інженерною реальністю

Коли говорять про product management, часто уявляють roadmap, backlog, customer interviews, prioritisation і нескінченні дискусії про features. Це правильна, але неповна картина. В R&D, особливо в hardware/software, systems engineering або regulated проектах, продуктовий менеджмент починається набагато раніше за backlog і має значно ширшу відповідальність.

Product manager у такому середовищі не зводиться до власника списку побажань. Він має розуміти market need, customer context, business case, technical feasibility, lifecycle, suppliers, certification windows, safety або cybersecurity impact. І водночас не підміняти engineering governance. Це складний баланс: тримати product intent і чесно бачити, що product decision не розв'язує engineering reality автоматично.

Особливо цікаво це видно в українському MilTech під час війни. Там product management - не абстрактні framework-и з книжок. Це короткий цикл: потреба з поля, технічна гіпотеза, прототип, перевірка, доставка, feedback, доопрацювання. Час реакції стає продуктовой характеристикою. Але саме тому дисципліна важлива ще більше: швидкість без traceability, версій, тестів і контрольованих рішень швидко перетворюється на хаос.

## Backlog з'являється пізніше

Backlog - це вже наслідок. До нього має бути ринок, customer problem, value hypothesis, constraints, positioning, business case. Якщо product management стартує з backlog, команда ризикує будувати те, що хтось уже сформулював як feature, але не перевірила, чи це справді розв'язує важливу проблему.

У R&D це особливо небезпечно. Engineering team часто вміє зробити дуже складні речі. Але "ми можемо" не означає "це потрібно". Технічна можливість легко зачаровує. З'являється prototype, потім roadmap, потім команда шукає customer value під уже зроблене рішення. Сильний product management має ставити незручне питання раніше: яку проблему ми закриваємо і чому саме цю?

Marketing у широкому сенсі - перше джерело product intent, а не реклама. Market segmentation, customer interviews, competitive analysis, pricing assumptions, go-to-market constraints, buyer journey, field feedback - усе це формує product direction до появи детальних requirements.

## Marketing Product Manager і Technical Product Manager

У складних продуктах часто потрібно розрізняти Marketing Product Manager і Technical Product Manager. Назви можуть бути іншими, але різниця в оптиці реальна.

Marketing PM ближче до ринку, customer segments, positioning, pricing, sales enablement, competitive landscape, customer success. Він питає: хто користувач, хто buyer, яку проблему розв'язуємо, чому зараз, як це продавати або впроваджувати, які alternatives має customer.

Technical PM ближче до architecture, system requirements, feasibility, technical trade-offs, dependencies, lifecycle, integration і constraints. Він питає: що це означає для system design, firmware, hardware, cloud, security, manufacturing, verification, support. Він не обов'язково замінює architect-а або systems engineer-а, але має достатньо глибини, щоб не продавати roadmap, відірваний від reality.

Ці ролі можна поєднувати в одній людині, але треба чесно бачити ризик. Якщо PM занадто market-only, продукт може стати красивою vision-презентацією без реалізованості. Якщо PM занадто technical-only, продукт може стати набором інженерно цікавих features без ринкової логіки. У складному R&D найкраще працює тандем або хоча б явне усвідомлення двох перспектив.

## System, hardware і software живуть у різних ритмах

У чисто software product scope можна змінювати відносно швидко. Release cadence гнучкіший, rollback часто можливий, telemetry швидко дає feedback. Software від цього не стає простим, просто темп змін інший.

Hardware живе довше. Є bill of materials, supplier lead times, manufacturing constraints, tooling, lab tests, certification, physical failures. Зміна connector-а, sensor-а або enclosure може зачепити supply chain, mechanical design, thermal profile, EMC testing, regulatory documentation. Це місяці, а не два sprint-и.

System/HW/SW продукт має одночасно кілька ритмів. Firmware може оновлюватися частіше за electronics. Cloud service ще частіше. Mechanical design повільніше. Certification має свої windows. Product manager пріоритезує features і водночас узгоджує ці швидкості. Інакше roadmap буде написаний у темпі software, а deliverability визначатиметься hardware.

## Regulated domain змінює правила гри

У regulated галузях product decision може змінити не лише scope, а й evidence burden. Automotive працює зі стандартами Automotive SPICE, ISO 26262, ISO/SAE 21434, multi-tier supplier ecosystem і багаторівневою відповідальністю. Aviation - це DO-178C/DO-254, ARP4754A і тривалі certification cycles. MedTech - IEC 62304, ISO 14971, FDA 21 CFR Part 11 і clinical evidence. Defense-related продукти додають контрактні вимоги, security controls і export constraints. Стандарти різні, але спільна риса одна: вартість зміни перевищує вартість самого implementation.

У такому контексті "невелика feature" може змінити safety analysis, cybersecurity scope, verification strategy, supplier qualification або audit package. Product manager не має сам приймати safety acceptance або змінювати technical baseline. Але він має бачити, що його product intent запускає controlled process.

Сильний PM у regulated engineering не сприймає compliance як ворога roadmap. Він сприймає його як constraint, який впливає на feasibility, time-to-market, risk і value. Якщо compliance impact видно рано, product decision може бути розумним. Якщо його видно після реалізації, це вже expensive surprise.

## Український MilTech під час війни

Окремий випадок - сучасний український MilTech. Про нього треба говорити обережно і без чутливих деталей, але управлінський урок очевидний: product management тут дуже близький до реального використання. Feedback loop може бути коротким і жорстким. Те, що в мирному enterprise product зайняло б квартал discovery, тут іноді проходить шлях від потреби до prototype дуже швидко.

Це створює унікальну силу: висока мотивація, прямий контакт із реальністю, швидке навчання, готовність імпровізувати. Але є й ризик: героїчний менеджмент замість сталого процесу. Якщо все тримається на кількох людях, які пам'ятають контекст, продукт стає крихким. Якщо feedback не структурується, команда повторює ті самі помилки. Якщо versions і evidence не ведуться, складно зрозуміти, що саме працювало і чому.

Product manager у такому середовищі має бути мостом між людьми, які бачать реальне використання, і людьми, які проектують. Він має швидко перетворювати feedback на requirements, пріоритезувати без романтики, захищати команду від scope chaos і водночас не втратити дисципліну: versions, traceability, test evidence, security boundaries.

## Процеси вчора і сьогодні

Раніше product management часто будувався навколо довгих roadmaps, статичних product requirement documents, market research reports і великих review. Це працювало там, де ринок стабільний, продукт змінюється повільно, а feedback приходить прогнозовано.

Сьогодні процеси стають більш continuous. Discovery не закінчується перед roadmap. Product analytics, telemetry, support tickets, field reports, customer interviews, competitive intelligence і engineering feedback постійно змінюють картину. Product tools допомагають збирати і пріоритезувати signals, але самі по собі не створюють product judgment.

У R&D важливо з'єднати product tools з engineering systems. Roadmap має бачити requirements, risks, defects, tests, architecture constraints, supplier status і release readiness. Інакше product management знову стане окремим світом, де value є, а feasibility з'являється занадто пізно.

## AI у продуктовому менеджменті

AI може сильно допомогти там, де PM витрачає час на ручний синтез інформації. Він може підсумовувати customer feedback, групувати support tickets, знаходити повторювані field patterns, порівнювати конкурентів, готувати draft requirements, складати release notes, пояснювати engineering constraints, будувати чернетку business case, перевіряти consistency між roadmap і implementation status.

Але AI не визначає product vision. Він не відповідає за value, market fit, risk appetite або ethical trade-offs. Він може показати options, але не має authority вибрати direction. У regulated або defense-related контексті AI-рекомендації мають бути reviewable: sources, assumptions, confidence, human approval. Особливо якщо в аналіз входять чутливі дані.

Найкраща роль AI - звільнити PM від частини інформаційної роботи, щоб більше часу залишилося на рішення. Для цього AI має бачити product backlog разом з engineering reality. Інакше він підсилить стару проблему: красиві product summaries без зв'язку з тим, що реально можна поставити.

## Практичний чекліст для R&D product manager-а

Перед тим як перетворювати product idea на roadmap item, я б перевіряв кілька речей. Чи зрозуміла customer або field problem? Чи є value hypothesis, яку можна перевірити? Чи описані target users і operating context? Чи відомі технічні constraints? Чи є impacted architecture areas? Чи зачіпає зміна hardware, firmware, cloud, manufacturing або support? Чи потрібен compliance impact analysis? Чи є supplier або lab dependency? Чи відомо, як ми будемо перевіряти успіх після release?

Окремо варто питати: що ми не робимо. Product management часто страждає від браку відмов більше, ніж від браку ідей. Якщо roadmap містить усе, він не має пріоритетів. Якщо PM не може пояснити, чому feature відкладена або відхилена, команда не бачить strategy.

Для R&D корисна ще одна звичка: фіксувати assumptions. На ранніх етапах багато рішень приймаються на неповній інформації. Це нормально. Ненормально - забути, що це були припущення. Коли assumption змінюється, roadmap має змінюватися свідомо, без накопичення випадкових задач.

Ще один простий артефакт - product decision note. Коротко: яка проблема, які options, який вибір, які assumptions, які ризики, які engineering constraints, хто погодив. Для R&D цього часто достатньо, щоб не втратити логіку рішення через два місяці, коли prototype вже змінився, а люди пам'ятають різні версії домовленостей.

Це також допомагає новим людям у команді швидше зрозуміти, що будується і чому саме так.

## Висновок

Product management в R&D починається там, де roadmap і prioritisation треба з'єднати з інженерною перевіркою. Це здатність перетворювати market або field need на реалізований, перевірений і корисний продукт. У software це один темп. У hardware/software systems - інший. У regulated domains - ще складніший. У сучасному українському MilTech - максимально практичний і швидкий, але від цього не менш вимогливий до дисципліни.

Сильний product manager не має контролювати все. Але він має бачити достатньо: value, feasibility, constraints, risk, compliance impact, delivery reality. І має вміти тримати міст між людьми, які потребують продукт, і людьми, які роблять його реальним.

Питання до читачів: у ваших R&D-командах product manager ближче до ринку, до engineering чи справді працює як міст між цими світами?
