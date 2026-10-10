# Über den Autor: Mykola Fedchyk (Nick Fedchik)

> **Buch:** [Architektur evidenzbasierter Expertensysteme](README.md)  
> **Inhaltsverzeichnis:** [README.md](README.md) · [Anhang A](appendix-a-evidence-governed-framework.md)  
> **Profile:** [LinkedIn — Nick Fedchik](https://www.linkedin.com/in/nickfedchik/) · [GitHub — nick-fedchik](https://github.com/nick-fedchik)  

---

## Berufliches Profil & technischer Hintergrund

**Mykola Fedchyk (Nick Fedchik)** ist ein ukrainischer Systemarchitekt, Informationstechnologie-Ingenieur und Forscher auf den Gebieten cyber-physikalischer Systeme, evidenzbasierter künstlicher Intelligenz und sicherheitskritischer eingebetteter Elektronik.

Er verfügt über mehr als zwanzig Jahre praktische Erfahrung in Softwaretechnik, dem Entwurf verteilter Netzwerkarchitekturen, harter Echtzeitsysteme auf Treiberebene sowie im technologischen Produktmanagement komplexer Innovationsvorhaben.

Heute ist er als **Technical Project Manager (Automotive)** beim weltweit führenden Halbleiterkonzern **Infineon Technologies** tätig und konzentriert sich auf die Entwicklung und Einführung zukunftsweisender Plattformlösungen für die Automobilindustrie:

* Hochzuverlässige Mehrkern-Mikrocontroller der Familie **Infineon AURIX (TriCore Lockstep)**;
* Erreichung höchster Standards der funktionalen Sicherheit (**ISO 26262 ASIL-D**);
* Echtzeitbetriebssysteme (RTOS, insbesondere **Zephyr RTOS**);
* Hard- und Software-Co-Design an der Schnittstelle von deterministischer Regelung, Sensornetzwerken und Edge-KI.

---

## Systemwurzeln, Linux-Kernel & Open-Source-Beiträge

Die ingenieurwissenschaftliche Denkweise des Autors wurde maßgeblich durch systemnahe Programmierung, Netzwerkarchitektur und Betriebssystem-Kernelentwicklung geprägt. Seit den frühen 2000er-Jahren war Mykola maßgeblich an Basissystemsoftware, dem Linux-Netzwerk-Stack und kritischen Open-Source-Projekten beteiligt:

### 1. Der Linux-Kernel und das Subsystem ebtables (Ethernet Bridge Tables)

* **Entwicklung des Moduls `ebt_vlan`:** Mykola ist der Hauptautor des 802.1Q-VLAN-Paketfiltermoduls für das Ethernet-Bridge-Filtersubsystem (**ebtables / Netfilter**) im Linux-Kernel (`net/bridge/netfilter/ebt_vlan.c`).
* Das Modul ermöglichte die hardware- und softwarenahe Erkennung, Inspektion und Filterung von VLAN-Tags (802.1Q Priority / VID) direkt im L2-Switching-Pfad des Linux-Kernels und wurde zum Industriestandard für Bridge-Firewalls, Router und sichere Switches.
* Die Arbeit mit Kernel-Strukturen (`sk_buff`), die strikte Vermeidung von dynamischen Speicherallokationen im kritischen Paketdurchsatz und die deterministische Interruptverarbeitung legten das Fundament für das Verständnis latenzkritischer Fail-Closed-Architekturen.

### 2. Das Projekt BusyBox & Embedded Linux

* **Portierung und Pflege von `arping`:** Mykola erstellte die offizielle Portierung des Netzwerkdiagnose-Werkzeugs `arping` in die Codebasis von **BusyBox** (`networking/arping.c` — *„The Swiss Army Knife of Embedded Linux“*).
* Das Dienstprogramm ermöglichte die Erkennung von MAC/IP-Kollisionen, Verbindungsdiagnosen und ARP-Probing unter extrem restriktiven Speicher- und Prozessorressourcen industrieller Router und Steuerungen.
* Darüber hinaus trug der Autor zur Optimierung und Fehlerbehebung systemnaher Utilities für Embedded-Linux-Plattformen bei, die in Millionen Industrieanlagen weltweit im Einsatz sind.

### 3. Die GNU C Library (glibc) & Protokollinfrastruktur

* **Erweiterung der glibc-Funktionalität:** Der Autor entwickelte Patches für die Systembibliothek **glibc** zur Implementierung der `getethertype`-Funktion auf Basis des Name Service Switch (NSS) zur Auflösung von Ethernet-Protokollnamen und -nummern über `/etc/ethertypes`. Diese Funktionalität fand breite Anwendung in Netzwerküberwachungstools und Paketfiltern.
* Upstream-Beiträge zur Unterstützung drahtloser und peripherer Schnittstellen (insbesondere im Subsystem `irda-usb` des Linux-Kernels).

### 4. Vom Kernel-Code zur Architektur evidenzbasierter Expertensysteme

Diese langjährige Erfahrung in den untersten Schichten von Betriebssystemen definierte die architektonische Philosophie dieser Monographie:

* **Determinismus statt Zufall:** Genauso wie eine `ebtables`-Paketfilterregel ein Paket nicht „mit einer gewissen Wahrscheinlichkeit“ durchleiten darf, darf ein evidenzbasiertes Expertensystem niemals auf stochastischen Mutmaßungen beruhen.
* **Das Zulassungsgateway (Admission Gateway):** Der symbolische Verifikationskern fungiert exakt wie eine Kernel-Firewall: Er erzwingt unverletzliche mathematische Invarianten und weist fehlerhafte Ausgaben von Sprachmodellen deterministisch ab, bevor sie die Laufzeitumgebung erreichen können.
* **Direkte Erdung in Silizium:** Der Weg von ersten Linux-Kernel-Treibern bis hin zu modernen Automobilplattformen (**Infineon AURIX TC4xx**) und Echtzeitsystemen (**Zephyr RTOS**) verleiht dem Autor eine durchgängige Perspektive — von Transistoren und Lockstep-Cores bis hin zu mathematischer Logik, Ontologien und künstlicher Intelligenz.

---

## Forschungsaktivitäten & Förderung evidenzbasierter KI

Mykola ist Autor einer vielbeachteten ingenieurwissenschaftlichen Publikationsreihe auf DOU, gewidmet folgenden Kernfragen:

1. **Die Krise dokumentierten Wissens:** Überwindung des „PDF-Dokumentenfriedhofs“ in großen R&D-Entwicklungsprojekten;
2. **Engineering Knowledge Graphs (EKG):** Durchgängige Verknüpfung von Normenanforderungen, Quellcode, Testfällen und Hardware-Komponenten;
3. **Neuro-symbolische Architektur (NeSy AI):** Aufbau hybrider Systeme, in denen stochastische Sprachmodelle (SLMs/LLMs) vom Sicherheitsperimeter isoliert und einem deterministischen Go- und Datalog-Kern untergeordnet sind;
4. **Maschinelle Sicherheitszertifizierung (GSN):** Algorithmische Synthese rechtssicherer und technischer Sicherheitsnachweise (*Safety Cases* via Goal Structuring Notation).

Diese Monographie ist das Resultat der Systematisierung jahrzehntelanger Industrieerfahrung zur Schließung der Kluft zwischen akademischer mathematischer Logik, generativer KI und der kompromisslosen industriellen Realität sicherheitskritischer Hardware.

---

## Ingenieurkredo

Zwei Jahrzehnte technischer Praxis — von hardwarenahen Netzwerktreibern im Linux-Kernel über ressourcenbeschränkte Mikrocontroller bis hin zu Hardware-Lockstep-Architekturen in Infineon AURIX und Sicherheitsnormen nach ISO 26262 ASIL-D — führten zu einer unverrückbaren Überzeugung:

> **„Echte Zuverlässigkeit beginnt dort, wo Annahmen enden und Determinismus einsetzt.“**

* **Im Betriebssystemkern** verzeihen Interrupthandler keine unvorhersehbaren Latenzspitzen, und der Netzwerk-Stack operiert nicht mit „Übertragungswahrscheinlichkeiten“: Ein Frame entspricht entweder exakt den Protokollinvarianten oder wird deterministisch verworfen.
* **In der Automobilelektronik und physischen Regelungssystemen** wird Sicherheit niemals in Prozentpunkten subjektiven Vertrauens gemessen. Wenn zwei Prozessorkerne im Takt-für-Takt-Lockstep arbeiten, ist eine Abweichung von einem einzigen Bit ein Hardwarefehler und keine „tolerierbare statistische Varianz“. Die physische Welt verzeiht keine Halluzinationen.
* **Im Zeitalter der künstlichen Intelligenz** wird diese Ingenieursdisziplin noch unverzichtbarer. Statistische neuronale Netze und große Sprachmodelle sind exzellente Heuristiken zur Generierung von Hypothesen, dürfen jedoch niemals direkte Kontrollgewalt über Aktuatoren erhalten. Die finale Entscheidungsgewalt gebührt formalen Verifizierern, mathematischer Logik und hardwaregeerdeten Beweisen:

> *„Das neuronale Netz schlägt vor — Der symbolische Kern verifiziert — Deterministische Hardware garantiert Sicherheit.“*

---

## Kontakt & Vernetzung

* **LinkedIn:** [https://www.linkedin.com/in/nickfedchik/](https://www.linkedin.com/in/nickfedchik/)
* **GitHub:** [https://github.com/nick-fedchik](https://github.com/nick-fedchik)
* **Monographie-Repository:** [https://github.com/nick-fedchik/articles](https://github.com/nick-fedchik/articles)
