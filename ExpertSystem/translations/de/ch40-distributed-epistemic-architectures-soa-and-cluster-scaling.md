# Kapitel 40. Verteilte Architektur evidenzbasierter Expertensysteme: Epistemische SOA, semantisches Routing, Speicherhierarchien und mehrquellenbasiertes defeasibles Schiedsverfahren

> **Buch:** [Architektur evidenzbasierter Expertensysteme](README.md) · [Teil VII: Reaktive Ausführung, systemübergreifender Wissensaustausch und verteilte SOA](part-07-runtime-and-knowledge-exchange.md)  
> **Vorheriges Kapitel:** [Kapitel 33. Systemübergreifender Wissensaustausch: Regelbereitstellung für Drittsysteme, Modell-Training und sicheres Feedback](ch33-inter-system-knowledge-exchange-and-model-teaching.md)  
> **Inhaltsverzeichnis:** [README.md](README.md)  
> **Autor:** [Mykola Fedchyk](about-the-author.md)  
> **Niveau:** Fortgeschritten: Leitende Systemarchitekten, Wissensingenieure, Entwickler verteilter Systeme, Compliance- und Sicherheitsexperten  
> **Lernziele:** Synthese der theoretischen Grundlagen aller Teile der Monographie zu einem vollständigen industriellen Referenzarchitektur-Prototyp eines verteilten Expertensystems; Entwurf einer dreischichtigen serviceorientierten Architektur (Epistemic SOA) mit strikter Entkopplung schlanker Clients von Wissensdiensten; Implementierung des syntaktisch-lexikalischen Präprocessings und des semantischen Routings von Abfragen; Überwindung der Diskrepanz zwischen Makrovolumina des Speichers und Mikrokapazitäten des schnellen On-Chip-Speichers mittels aktiver Arbeitssätze (Active Working Sets) und doppelter Ping-Pong-Pipeline-Pufferung; Organisation der föderierten Suche (Scatter-Gather) sowie der mehrquellenbasierten defeasiblen Evidenzaggregation (ASPIC+) bei unvollständigem oder unsicherem Wissen; fundierte ingenieurtechnische Entscheidungen bezüglich Skalierung und Auswahl von Hardware- und Softwareklassen im Rahmen des Projektbudgets.

---

## Abstract

Alle vorangegangenen Kapitel dieser Monographie haben eine fundamentale These schlüssig bewiesen: Ein zuverlässiges Expertensystem darf nicht auf stochastischen Verallgemeinerungen großer Sprachmodelle oder ungeprüften Faktensammlungen beruhen. Wahre Expertise erfordert mathematischen Determinismus, formale Inferenz, bytegenaue Verwahrungskette der Primärquellen ($ZHR = 1.00$) und unabhängige Auditierbarkeit. Der Übergang des Systems vom isolierten Laborprüfstand in reale industrielle F&E-Umgebungen stellt den Softwarearchitekten jedoch vor eine neue kritische Herausforderung: **Skalierbarkeit, Verteiltheit und allgegenwärtige Verfügbarkeit (Ubiquitous Availability)**.

Dieses Kapitel schließt den theoretischen Zyklus der Monographie ab und synthetisiert für den Leser einen **industriellen Referenzarchitektur-Prototyp eines evidenzbasierten Expertensystems**. Losgelöst von spezifischen kommerziellen Chipmarken oder proprietären Schnittstellen definiert dieses Kapitel funktionale Klassen von Hardware- und Software-Recheneinheiten, die für den Aufbau eines hochskalierbaren Systems erforderlich sind. Gegenstand der Analyse ist eine dreischichtige epistemische serviceorientierte Architektur (Epistemic SOA), in der schlanke mobile Clients (Smartphones, Terminals) über einen intelligenten semantischen Broker mit einem föderierten Cluster spezialisierter Domänendienste interagieren. Eine mehrstufige Speicherhierarchie wird formalisiert: Mathematisch wird die Bedingung der Pipeline-Kontinuität und der vollständigen Latenzverbergung ($`T_{\mathrm{refill}} \ll T_{\mathrm{execute}}`$) durch das Konzept aktiver Wissensarbeitssätze (*Active Working Sets*) und symmetrischer doppelter Pufferung (*Ping-Pong Pipeline*) bewiesen. Für Szenarien mit unvollständigen lokalen Wissensbasen wird ein Modell der föderierten Suche (*Scatter-Gather*) in Verbindung mit einer mehrquellenbasierten defeasiblen Evidenzaggregation (*ASPIC+*) entwickelt, welches Korroboration, die Entkräftung schwacher Empfehlungen und die qualifizierte Erkennung normativer Widersprüche gewährleistet. Besondere Aufmerksamkeit gilt der verteilten semantischen Triage von Wissensdefiziten auf Basis von 64-Bit-SimHash und BK-Bäumen mit Schutz vor Entartung des Merkmalsraums (Degeneracy Guard) und autonomen WAL-Puffern sowie dem Aufbau eines unveränderlichen Entscheidungsbuchs (Immutable Decision Ledger) auf hierarchischen Merkle-Bäumen mit WORM-Verankerung zur Gewährleistung der rechtlichen und technischen Unabstreitbarkeit (Non-Repudiation) im ASIL-D-Audit. Abschließend bietet eine ingenieurtechnische Technologieauswahlmatrix dem Systemarchitekten konkrete Richtlinien zur Bereitstellung eines lauffähigen Prototyps auf verfügbarer Hardware im Einklang mit dem jeweiligen Projektbudget.


---

## 1. Synthese der Monographie: Vom lokalen Modell zum Referenzarchitektur-Prototyp

Neununddreißig vorangegangene Kapitel dieses Werkes haben die Anatomie einer evidenzbasierten künstlichen Intelligenz detailliert seziert. Die Monographie untersuchte die Erkenntnistheorie maschinellen Wissens ([Kapitel 2](ch02-epistemology-of-machine-knowledge.md)), Modelle und Typologien von Wissensbasen ([Kapitel 7](ch07-knowledge-base-typology.md)), die automatisierte Extraktion von Normen aus Primärquellen ([Kapitel 15](ch15-knowledge-extraction-and-kb-construction.md)), strikte deterministische Inferenz ([Kapitel 16](ch16-expert-systems-architecture.md), [Kapitel 31](ch31-syllogistic-reasoning-and-relation-lattices.md)), die Behebung maschineller Halluzinationen ([Kapitel 38](ch38-curing-machine-hallucinations-and-knowledge-deficits.md)) und Popper'sche Compliance-Audits ([Kapitel 39](ch39-active-compliance-auditor-and-popperian-testing.md)).

In der industriellen Praxis stellt sich jedoch unweigerlich die entscheidende Frage: **Wie lassen sich all diese theoretischen Konzepte zu einem einheitlichen, lauffähigen Gesamtsystem auf Unternehmensebene zusammenführen?**

In realen F&E-Projekten der Hochtechnologie (Automobilindustrie, Avionik, Medizintechnik) umfasst der regulatorische Rahmen dutzende verbindliche Standards (ISO, IEC, DO, IEEE, RFC), die gemeinsam mit wörtlichen Primärtexten, Konformitätsmatrizen und Defeatern Dutzende Gigabyte an Daten beanspruchen. Gleichzeitig benötigen Forschungsingenieure, Auditoren und Softwareentwickler permanenten, ortsunabhängigen Zugriff auf das System:
- am Büroarbeitsplatz an der Workstation während des Architekturentwurfs;
- im Prüflabor oder auf dem Testgelände über Tablets und Smartphones;
- in der automatisierten Continuous-Integration-Pipeline (CI/CD) während nächtlicher Code-Audits.

Der Versuch, diese Anforderung durch einen monolithischen lokalen „Fat Client“ zu lösen, scheitert zwangsläufig. Ein mobiles Endgerät ist weder in der Lage, mehrbändige, gigabytegroße Wissensbasen im Arbeitsspeicher vorzuhalten, noch verfügt es über die Rechenleistung für tiefe semantische Analysen. Andererseits garantiert ein gewöhnlicher Cloud-Dienst auf Basis populärer Sprachmodelle weder Determinismus noch Halluzinationsfreiheit oder eine bytegenaue Verwahrungskette (*Custody*) der Primärquellen.

Erforderlich ist daher ein **verteilter Referenzarchitektur-Prototyp**, der auf den Prinzipien der serviceorientierten Architektur (SOA), einer strikten Trennung der Zuständigkeiten und einer mathematisch fundierten Speicherhierarchie aufbaut.

```mermaid
flowchart TD
    accTitle: Referenzarchitektur-Prototyp eines verteilten Expertensystems
    accDescr: Dreischichtige Epistemic-SOA-Architektur mit schlankem Client, semantischem Broker und föderiertem Cluster von Domänenwissensdiensten.

    subgraph TIER1["Schicht 1: Allgegenwärtige Interaktions- und Thin-Client-Schicht (Interaction Tier)"]
        direction LR
        CLI["<b>Konsolenterminal des Ingenieurs</b><br>CLI / TUI (Linux-Workstation)<br><i>Code-Audit und Regelverifikation</i>"]
        WEB["<b>Browser-Client</b><br>Web / WebAssembly / PWA<br><i>Interaktive Anforderungsanalyse</i>"]
        MOB["<b>Mobiles Auditor-Terminal</b><br>Tablet / Smartphone (Field Inspection)<br><i>Verifikation auf dem Testgelände</i>"]
    end

    subgraph TIER2["Schicht 2: Semantischer Broker und Wissens-Gateway (Semantic Broker & Knowledge Gateway)"]
        direction TB
        DISPATCHER["<b>Syntaktisch-lexikalischer Diskurs-Präprozessor</b><br><i>Extraktion der ontologischen Signatur der Abfrage</i>"]
        REGISTRY["<b>Dynamisches Register der Wissensdienste</b><br><i>Berechnung der Affinität Affinity(S_query, P_service)</i>"]
        ROUTER["<b>Inhaltsbasierter semantischer Router</b><br><i>Semantische Adressierung an die Zieldienste</i>"]
        ARBITER["<b>ASPIC+-Argumentationsschiedsrichter</b><br><i>Korroboration, Entkräftung und Fail-Closed-Schutzschleuse</i>"]
        
        DISPATCHER --> REGISTRY --> ROUTER
        ROUTER -.->|Normative Kollision oder Wissensunvollständigkeit| ARBITER
    end

    subgraph TIER3["Schicht 3: Föderierter Cluster von Domänenwissensdiensten (Domain Evaluation Tier)"]
        direction LR
        subgraph DS1["Domänendienst 1: Sicherheit (ISO 26262)"]
            DS1_CORE["<b>ASIL-D-Safe-Kern</b><br>• Unveränderlicher Speicher .kp<br>• AWS-Arbeitssatz-Manager<br>• Deterministischer Regelprozessor<br>• Ping-Pong-Doppelpuffer"]
        end
        subgraph DS2["Domänendienst 2: Avionik (DO-178C)"]
            DS2_CORE["<b>Airborne-Safe-Kern</b><br>• Unveränderlicher Speicher .kp<br>• AWS-Arbeitssatz-Manager<br>• Deterministischer Regelprozessor<br>• Ping-Pong-Doppelpuffer"]
        end
        subgraph DSN["Domänendienst N: Cybersicherheit (ISO 21434)"]
            DSN_CORE["<b>Cyber-Physical-Security-Kern</b><br>• Unveränderlicher Speicher .kp<br>• AWS-Arbeitssatz-Manager<br>• Deterministischer Regelprozessor<br>• Ping-Pong-Doppelpuffer"]
        end
    end

    TIER1 ==>|"Schlanke typisierte Anfragen (bis zu 4 KB)<br>gRPC / Protobuf / TLS · Zero-Knowledge-Store auf dem Client"| TIER2
    TIER2 ==>|"Streaming-Rendering des Evidenzpakets<br>Verdict + EvidenceMap + Ed25519 (ZHR = 1.00)"| TIER1

    ROUTER ==>|"Streaming-AWS (8–32 MB)<br>Hochgeschwindigkeitsbus / 10GbE LAN"| DS1
    ROUTER ==>|"Streaming-AWS (8–32 MB)<br>Hochgeschwindigkeitsbus / 10GbE LAN"| DS2
    ROUTER ==>|"Streaming-AWS (8–32 MB)<br>Hochgeschwindigkeitsbus / 10GbE LAN"| DSN

    DS1 -.->|"Lokale Evidenzfragmente"| ARBITER
    DS2 -.->|"Lokale Evidenzfragmente"| ARBITER
    DSN -.->|"Lokale Evidenzfragmente"| ARBITER

    classDef client fill:#e1f5fe,stroke:#0288d1,stroke-width:2px,color:#01579b;
    classDef broker fill:#ede7f6,stroke:#512da8,stroke-width:2px,color:#311b92;
    classDef domain fill:#e8f5e9,stroke:#2e7d32,stroke-width:2px,color:#1b5e20;
    classDef arbiter fill:#fff3e0,stroke:#f57c00,stroke-width:2px,color:#e65100;

    class CLI,WEB,MOB client;
    class DISPATCHER,REGISTRY,ROUTER broker;
    class DS1_CORE,DS2_CORE,DSN_CORE domain;
    class ARBITER arbiter;
```

---

## 2. Epistemische serviceorientierte Architektur (Epistemic SOA) und schlanke Clients

Die serviceorientierte Architektur adressiert im klassischen Software-Engineering primär die lose Kopplung von Komponenten, die unabhängige Bereitstellung und die Wiederverwendbarkeit von Geschäftsfunktionen [[4]](#src-4). In **epistemischen evidenzbasierten Systemen** erhält das SOA-Paradigma eine qualitativ neue Dimension: Es gewährleistet die **Isolation von Wahrheitskontexten** und den **allgegenwärtigen Benutzerzugriff**, ohne die lückenlose Verwahrungskette der Evidenzen zu gefährden.

### 2.1. Das Konzept des schlanken Interaktions-Clients (Thin Interaction Client)

Das Endgerät des Benutzers (Qualitätsingenieur, Entwickler oder Auditor) wird als **ultraschlanker Interaktionsagent** definiert:
1. **Null-Wissensspeicher ($`S_{\mathrm{client}} = 0`$):** Die Client-Anwendung speichert keine Wissensbasen und hält keine lokalen Kopien von Normen im Cache. Dies eliminiert das Risiko versioneller Desynchronisation von Vorschriften innerhalb des Ingenieurteams vollständig.
2. **Abstraktion von Recheneinheiten:** Der Client benötigt weder neuronale Hardwarebeschleuniger (NPU/GPU) noch spezialisierte Inferenzprozessoren. Dadurch lässt sich die Client-Schnittstelle des Expertensystems problemlos auf Standard-Smartphones unter mobilen Betriebssystemen, in Webbrowsern via WebAssembly/PWA oder im Terminal einer Entwickler-Workstation betreiben.
3. **Strikte Typisierung von Anfragen und Antworten:** Der Datenaustausch erfolgt über hocheffiziente binäre RPC-Verträge (beispielsweise Protocol Buffers über gRPC oder abgesicherte WebSockets):
   - Client-Anfrage: natürlichsprachliche Zeichenkette oder strukturierter Identifikator eines Entwicklungsartefakts (Paketgröße $\le 2\text{--}4\,\text{KB}$).
   - Client-Antwort: strukturiertes Evidenztupel, bestehend aus dem logischen Urteil (`ACCEPT`, `REFUSAL`, `QUALIFIED`), den Identifikatoren der angewendeten Regeln, wörtlichen Zitaten mit Byte-Koordinaten in der Primärquelle sowie 256-Bit-Kryptohashes zur Sicherung der Verwahrungskette.

### 2.2. Wahrung der Invariante $ZHR = 1.00$ auf dem Client

Von essenzieller Bedeutung ist, dass der schlanke Client nicht autorisiert ist, die empfangene Antwort eigenständig zu interpretieren oder zusammenzufassen. Er erfüllt eine rein repräsentative Funktion:

```math
\text{Display}(\text{Response}) = \text{Render}(\text{Verdict}, \text{EvidenceMap}, \text{Signature})
```

wobei:
- $\text{Verdict} \in \{\text{ACCEPT}, \text{REFUSAL}, \text{QUALIFIED}\}$ den formalen logischen Inferenzstatus bezeichnet, der vom Server generiert wurde;
- $\text{EvidenceMap}$ die bytegenaue Abbildung der Inferenzterme auf Koordinaten in den Primärquellen darstellt;
- $\text{Signature}$ die digitale Ed25519-Signatur des Evidenzpakets zur Gewährleistung der Unveränderlichkeit ist.

Selbst wenn die Netzwerkverbindung zwischen Client und Broker bei Feldversuchen temporär abreißt, verliert das vorliegende Ergebnis weder seine technische noch seine rechtliche Beweiskraft: Jede dem Ingenieur präsentierte Tatsache bleibt durch die kryptographische Signatur und den Hash der Primärquelle verifiziert.

---

## 3. Semantischer Broker, Diskurs-Präprozessor und inhaltsbasiertes Routing

Arbeitet im Ingenieur-Ökosystem nicht nur eine einzelne Wissensbasis, sondern ein Cluster hochspezialisierter Dienste, entsteht eine zentrale Koordinationsaufgabe: **Wer bestimmt auf welche Weise, an welche Zieldienste eine Ingenieurabfrage geleitet wird?**

In klassischen Websystemen erfolgt das Routing anhand von Netzwerkadressen oder URL-Präfixen (beispielsweise `/api/v1/orders`). In evidenzbasierten Expertensystemen erfolgt die Adressierung hingegen rein **semantisch (Content-Based Semantic Routing)** und stützt sich auf die Ergebnisse eines linguistischen Präprozessors.

### 3.1. Syntaktisch-lexikalischer Diskurs-Präprozessor

Sobald eine Abfrage eintrifft, leitet der semantische Broker diese durch eine mehrstufige semantische Analysepipeline:

```mermaid
flowchart TD
    accTitle: Pipeline der syntaktisch-lexikalischen Diskursanalyse im Semantischen Broker
    accDescr: Schrittweise Transformation einer natürlichsprachlichen Ingenieurabfrage in eine strukturierte ontologische Signatur für das Routing von Domänenregeln.

    IN["&emsp;&emsp;&emsp;&emsp; <b>EINGEHENDE INGENIEURABFRAGE ZUM COMPLIANCE-AUDIT (Natural Language Query)</b> &emsp;&emsp;&emsp;&emsp;<br><i>«Ist die Verwendung dynamischer Speicherallokation für ein ASIL-D-Modul gemäß ISO 26262-6 zulässig?»</i><br><i>Herkunftsquelle: Ingenieur-Workstation / Continuous-Integration-Pipeline (CI/CD Quality Gate)</i>"]

    STEP1["&emsp;&emsp;&emsp;&emsp; <b>1. DISKURSKLASSIFIKATION (Discourse Qualification)</b> &emsp;&emsp;&emsp;&emsp;<br>• Pragmatischer Typ: <code>NORMATIVE_COMPLIANCE_QUERY</code> &emsp;·&emsp; • Zielmodalität: <code>DEONTIC_PROHIBITION_OR_PERMISSION</code><br><i>Semantische Intention: formale Verifikation der Zulässigkeit einer Architekturentscheidung gegen den Zertifizierungsstandard</i>"]

    STEP2["&emsp;&emsp;&emsp;&emsp; <b>2. SYNTAKTISCH-FRAME-BASIERTES PARSING (Syntactic Frame Parsing)</b> &emsp;&emsp;&emsp;&emsp;<br>• Akteur: <i>«Softwaremodul»</i> &emsp;·&emsp; • Aktion: <i>«dynamische Speicherallokation» (Dynamic Allocation)</i> &emsp;·&emsp; • Attribut: <code>ASIL_Level = 'D'</code><br><i>Extraktion der Prädikat-Argument-Struktur der Anforderung für den deterministischen Abgleich mit ontologischen Axiomen</i>"]

    STEP3["&emsp;&emsp;&emsp;&emsp; <b>3. EXTRAKTION DER ONTOLOGISCHEN SIGNATUR (Domain Signature Extraction)</b> &emsp;&emsp;&emsp;&emsp;<br>• Domänenmarker: <code>{ISO-26262, ASIL-D, Software-Architecture, Memory-Management}</code><br>• Abfrageprofil: <code>TargetDomain = 'Automotive-Safety'</code> &emsp;·&emsp; • Normativ: <code>Standard = 'ISO-26262-6:2018'</code><br><i>Bildung des Signaturvektors der Abfrage zur Berechnung der Affinitätskoeffizienten im dynamischen Register der Dienste</i>"]

    OUT["&emsp;&emsp;&emsp;&emsp; <b>ZUWEISUNG DES PRIMÄREN EVALUATORS (Primary Evaluator Dispatching)</b> &emsp;&emsp;&emsp;&emsp;<br>Zugewiesener Dienst: <b>Domänendienst für funktionale Sicherheit im Automobilbereich (ISO 26262-6:2018)</b> &emsp;·&emsp; <code>Affinity(S_query, P_service) = 1.00</code><br><i>Das Active-Working-Set-Paket (AWS) wird über den lokalen Hochgeschwindigkeitsbus an den systolischen deterministischen Regelkern übertragen</i>"]

    IN ==> STEP1
    STEP1 ==> STEP2
    STEP2 ==> STEP3
    STEP3 ==> OUT

    classDef query fill:#e3f2fd,stroke:#1565c0,stroke-width:2px,color:#0d47a1;
    classDef step1 fill:#ede7f6,stroke:#5e35b1,stroke-width:2px,color:#311b92;
    classDef step2 fill:#fff8e1,stroke:#ffa000,stroke-width:2px,color:#e65100;
    classDef step3 fill:#e8f5e9,stroke:#2e7d32,stroke-width:2px,color:#1b5e20;
    classDef target fill:#f3e5f5,stroke:#7b1fa2,stroke-width:2px,color:#4a148c;

    class IN query;
    class STEP1 step1;
    class STEP2 step2;
    class STEP3 step3;
    class OUT target;
```

### 3.2. Dynamisches Register der Wissensdienste (Knowledge Service Registry)

Der semantische Broker verwaltet im flüchtigen Speicher ein aktuelles Verzeichnis der im Cluster verfügbaren Domänendienste. Jeder Dienst publiziert bei seiner Registrierung ein formales Manifest:
- **Identifikator und Endpunkt:** Netzwerkadresse (IP/Port oder Nachrichtenwarteschlange);
- **Wissensraum (Knowledge Domain):** Menge der abgedeckten Standards und Ontologieversionen (gemäß SemVer 2.0.0);
- **Klasse der Recheneinheit:** Software-CPU, Vektorprozessor-Array oder dedizierter Hardware-Inferenzkern für deterministische Ausführung;
- **Durchsatzmetriken:** Aktuelle Warteschlangentiefe und geschätzte Inferenzlatenz.

Der Broker berechnet eine Affinitätsfunktion zwischen der ontologischen Signatur der Abfrage $`\mathcal{S}_{\mathrm{query}}`$ und dem Wissensraum $`\mathcal{P}_{\mathrm{domain}}`$ im Dienstpass $`\mathcal{P}_{\mathrm{service}}`$:

```math
\mathrm{Affinity}(\mathcal{S}_{\mathrm{query}}, \mathcal{P}_{\mathrm{service}}) = \frac{\lvert\mathcal{S}_{\mathrm{query}} \cap \mathcal{P}_{\mathrm{domain}}\rvert}{\lvert\mathcal{S}_{\mathrm{query}}\rvert}
```

wobei:
- $`\mathcal{S}_{\mathrm{query}}`$ die Menge der ontologischen Konzepte und Domänenmarker bezeichnet, die aus der Ingenieurabfrage extrahiert wurden;
- $`\mathcal{P}_{\mathrm{domain}}`$ die Menge der unterstützten Standards und Ontologieversionen im Dienstpass $`\mathcal{P}_{\mathrm{service}}`$ darstellt;
- $`\lvert \mathcal{S}_{\mathrm{query}} \cap \mathcal{P}_{\mathrm{domain}} \rvert`$ die Mächtigkeit der Schnittmenge der ontologischen Signaturen angibt (Anzahl übereinstimmender Konzepte);
- $`\lvert \mathcal{S}_{\mathrm{query}} \rvert`$ die Gesamtzahl der ontologischen Marker in der Eingabeabfrage beziffert;
- $`\mathrm{Affinity}(\mathcal{S}_{\mathrm{query}}, \mathcal{P}_{\mathrm{service}}) \in [0, 1]`$ den normierten Koeffizienten der semantischen Affinität darstellt (wobei der Wert $1{,}0$ einer vollständigen Abdeckung der Abfragedomäne entspricht).

Der Dienst mit dem höchsten Affinitätswert wird als primäre Ausführungseinheit (*Primary Evaluator*) zugewiesen.

---

## 4. Theoretisches Modell der mehrstufigen Speicherhierarchie und doppelten Pufferung

Jedes Rechensystem unterliegt den physikalischen Gesetzmäßigkeiten der Signalübertragung und der Organisation von Speicherhierarchien [[3]](#src-3): **Je größer das Speichervolumen, desto höher ist die Zugriffslatenz; umgekehrt verfügt extrem schneller Speicher, der direkt in die Logik der Recheneinheit integriert ist, über eine stark begrenzte Kapazität**.

In industriellen Expertensystemen erreicht diese Diskrepanz ein kritisches Ausmaß:
1. **Makrovolumen des Wissensspeichers ($`V_{\mathrm{corpus}}`$):** Das vollständige Normen-Repository umfasst einschließlich Primärtexten, Syntaxbäumen und Rückverfolgbarkeitsmatrizen $4\text{--}32\,\text{GB}$.
2. **Mikrokapazität des schnellen Speichers des Regelprozessors ($`V_{\mathrm{local}}`$):** Der lokale Speicher eines spezialisierten Regelprozessors (Block-RAM in Siliziumbeschleunigern oder schneller L1/L2-Cache) überschreitet selten wenige hundert Kilobyte ($256\text{--}512\,\text{KB}$).
3. **Bandbreite des Übertragungskanals ($`B_{\mathrm{channel}}`$):** Die Übertragung des gesamten Wissenskorpus an die Recheneinheit bei jeder Abfrage würde zu unzulässigen Stillstandszeiten des Systems führen ($> 1$ Minute), was einen interaktiven Dialog mit dem Ingenieur ($0{,}5\text{--}2{,}0\,\text{s}$) verunmöglicht.

```mermaid
flowchart TD
    accTitle: Theoretische Speicherhierarchie und doppelte Pufferung deterministischer Systeme
    accDescr: Mehrstufiges Datenzugriffsmodell vom Universalspeicher (Universe Storage) bis zum Hardware-BRAM-Bank mit Latenzverbergung.

    subgraph L0["Ebene 0: Universeller Wissensspeicher (Universe Storage) · 4–32 GB"]
        direction TB
        L0_DATA["• Vollständiger Korpus unveränderlicher Quellen (Texte, Diagramme, Byte-Grenzen)<br>• Vollständige DWARF-ähnliche Evidenzabbildungstabelle<br>• Globale ontologische Graphen und Indizes<br><i>Medium: Nichtflüchtiger NVMe-Speicher / Großvolumiges System-RAM</i>"]
    end

    subgraph L1["Ebene 1: Cache des aktiven Standards (Board Domain Cache) · 512 MB – 4 GB"]
        direction TB
        L1_DATA["• Vorab geladene Regeln, Defeater und endliche Automaten des Zielstandards<br>• Keine Textzitate (ausschließlich normalisierte numerische Prädikate und Hashes)<br><i>Medium: Lokales schnelles RAM des Rechenknotens</i>"]
    end

    subgraph L2["Ebene 2: Symmetrische doppelte Pufferung (Ping-Pong Local BRAM) · 2 × 256 KB"]
        direction LR
        BANK_A["<b>BANK A (256 KB)</b><br><i>Aktive systolische Ausführung</i><br>Deterministischer Inferenzkern<br><b>0 Wartezyklen (0 wait-states)</b>"]
        SWITCH{{"<b>Hardware-Umschalter</b><br>Umschaltung in 1 Takt"}}
        BANK_B["<b>BANK B (256 KB)</b><br><i>Hintergrund-Pipeline-Befüllung</i><br>Autonomer DMA-Controller<br><b>T_refill ≈ 218 µs</b>"]
    end

    L0 ==>|"Projektion des aktiven Arbeitssatzes (AWS: 8–32 MB)<br>Standard-Datenübertragungsbusse"| L1
    L1 ==>|"Autonome Direct-Memory-Access-Pipeline (DMA Pipeline: 1.2 GB/s)<br>Bedingung der Latenzverbergung: T_refill ≪ T_execute"| L2

    BANK_A <--> SWITCH
    SWITCH <--> BANK_B

    classDef l0 fill:#eceff1,stroke:#546e7a,stroke-width:2px,color:#263238;
    classDef l1 fill:#e0f2f1,stroke:#00897b,stroke-width:2px,color:#004d40;
    classDef active fill:#e8f5e9,stroke:#2e7d32,stroke-width:2px,color:#1b5e20;
    classDef bg fill:#fff3e0,stroke:#f57c00,stroke-width:2px,color:#e65100;
    classDef sw fill:#ede7f6,stroke:#512da8,stroke-width:2px,color:#311b92;

    class L0_DATA l0;
    class L1_DATA l1;
    class BANK_A active;
    class BANK_B bg;
    class SWITCH sw;
```

### 4.1. Konzept des aktiven Arbeitssatzes (Active Working Set — AWS)

Anstatt den gesamten Wissenskorpus zu transferieren, erzeugt der semantische Dispatcher einen **aktiven Arbeitssatz (Active Working Set — AWS)**:

```math
\mathrm{AWS} = \mathrm{Closure}(\mathcal{R}_{\mathrm{target}}) = \mathcal{R}_{\mathrm{target}} \cup \mathrm{Prerequisites}(\mathcal{R}_{\mathrm{target}}) \cup \mathrm{Defeaters}(\mathcal{R}_{\mathrm{target}})
```

wobei:
- $`\mathcal{R}_{\mathrm{target}}`$ die Menge der Zielregeln bezeichnet, die durch die ontologische Signatur der Abfrage aktiviert wurden;
- $`\mathrm{Prerequisites}(\mathcal{R}_{\mathrm{target}})`$ die Voraussetzungsprädikate darstellt, die für die deterministische Auflösung von Abhängigkeiten erforderlich sind;
- $`\mathrm{Defeaters}(\mathcal{R}_{\mathrm{target}})`$ die Menge der assoziierten Defeater (Ausnahmen und Gegenregeln) bildet;
- $`\mathrm{AWS}`$ den transitiven logischen Abschluss der Regeln ohne Textzitate der Primärquellen darstellt, optimiert für das unmittelbare Laden in den lokalen Cache der Recheneinheit.

Das AWS enthält ausschließlich das mathematische Skelett der Regeln: numerische Grenzwerte, logische Bedingungen, Defeater-Identifikatoren und 256-Bit-Zitathashes. Die wörtlichen Primärtexte verbleiben im Speicher der Ebene 0. Das Volumen des AWS reduziert sich hierdurch auf **$8\text{--}32\,\text{MB}$**, was über standardisierte Industriebusse in Bruchteilen einer Sekunde übertragen werden kann.

### 4.2. Mathematische Bedingung der Pipeline-Kontinuität bei doppelter Pufferung (Ping-Pong Pipeline)

In deterministischen Expertensystemen ist es unzulässig, dass ein hochperformanter Inferenzkern anhält, während er auf das Eintreffen des nächsten Regelblocks aus einem langsameren externen Speicher wartet. Zur Eliminierung von Wartezeiten wird der Speicher des Regelprozessors als zwei symmetrische, voneinander unabhängige Speicherbänke gleicher Größe organisiert: Bank A und Bank B. Während der Prozessor die Regeln in der einen Bank verarbeitet, befüllt ein Direct-Memory-Access-Controller (DMA) im Hintergrund die andere Bank mit dem nächsten Wissensblock.

Die Dauer für das Befüllen einer Bank mit Daten aus dem lokalen System-RAM ($`T_{\mathrm{refill}}`$) wird durch das Verhältnis der Bankkapazität zur effektiven Busbandbreite bestimmt:

```math
T_{\mathrm{refill}} = \frac{V_{\mathrm{bank}}}{B_{\mathrm{bus}}}
```

wobei:
- $`V_{\mathrm{bank}}`$ die Kapazität einer Bank des schnellen On-Chip-Speichers bezeichnet (beispielsweise $256\,\text{KB}$ Block-RAM auf einem FPGA oder dediziertes statisches On-Chip-SRAM);
- $`B_{\mathrm{bus}}`$ die effektive Bandbreite des internen Wissenstransferbusses angibt (beispielsweise $1{,}2\,\text{GB/s}$ für einen AXI4-Stream-Bus oder PCIe-DMA);
- $`T_{\mathrm{refill}}`$ die berechnete Zeitspanne für das autonome Hintergrundbefüllen der Bank ohne Beteiligung des Regelprozessors darstellt;
- $`T_{\mathrm{execute}}`$ die Zeitspanne beziffert, die der Regelprozessor für den vollständigen systolischen Abgleich des eingehenden Anforderungspakets mit dem Wissensblock der Kapazität $`V_{\mathrm{bank}}`$ benötigt.

**Theorem der Latenzverbergung (Latency Hiding Invariant):**
Die Inferenz-Pipeline arbeitet genau dann mit einem maximalen Auslastungsgrad der Rechenressourcen ($\eta = 1{,}00$) ohne einen einzigen Wartezyklus des Prozessors (Zero Wait-States), wenn die Dauer des Hintergrundbefüllens die Rechenzeit für den aktuellen Regelblock nicht überschreitet:

```math
T_{\mathrm{refill}} \le T_{\mathrm{execute}}
```

**Praktische ingenieurtechnische Anwendung des Berechnungsergebnisses:**

Der ermittelte Wert $`T_{\mathrm{refill}}`$ und dessen Verhältnis zu $`T_{\mathrm{execute}}`$ stellen nicht bloß eine theoretische Abschätzung dar, sondern dienen als **operatives Werkzeug für Systementwurf und dynamisches Dispatching**:

1. **Hardware-Dimensionierung (Hardware Dimensioning):**
   Die Formel determiniert die maximal zulässige Bankgröße des schnellen Speichers $`V_{\mathrm{bank}}`$ beim Entwurf des Hardwarekomplexes:

```math
V_{\mathrm{bank}} \le B_{\mathrm{bus}} \cdot T_{\mathrm{execute}}
```

   Dimensioniert der Architekt die Bank fehlerhaft übermäßig groß (beispielsweise $`V_{\mathrm{bank}} = 16\,\text{MB}`$), benötigt ein Bus mit einer Bandbreite von $1{,}2\,\text{GB/s}$ für deren Befüllung:

```math
T_{\mathrm{refill}} = \frac{16\,\text{MB}}{1.2\,\text{GB/s}} \approx 13.3\,\text{ms}
```

   Schließt der Prozessor die Regelprüfung bereits in $`T_{\mathrm{execute}} = 2{,}0\,\text{ms}`$ ab, tritt eine gravierende Pipeline-Blockade (*Pipeline Stall*) auf: Der Prozessor wird gezwungen, für $11{,}3\,\text{ms}$ im Leerlauf auf den Bus zu warten, und der Wirkungsgrad des Inferenzkerns bricht ein auf:

```math
\eta = \frac{T_{\mathrm{execute}}}{T_{\mathrm{refill}}} = \frac{2.0\,\text{ms}}{13.3\,\text{ms}} \approx 0.15
```

   (85 % der Taktzyklen des teuren Siliziumkerns gehen ungenutzt verloren). Die Berechnung beweist: Um die Pipeline-Kontinuität zu wahren, darf die Größe eines Regelblocks $256\text{--}512\,\text{KB}$ nicht überschreiten.

2. **Dynamische Regelblockplanung (Adaptive Chunk Sizing) in der Software-Laufzeitumgebung:**
   Während der Erstellung des aktiven Arbeitssatzes (AWS) evaluiert der semantische Dispatcher dynamisch die Komplexität der Regeln des nachfolgenden Blocks. Sind die Regeln einfach strukturiert und sinkt die prognostizierte Ausführungszeit ($`T_{\mathrm{execute}} \to 0{,}5\,\text{ms}`$), paketiert der Dispatcher die Regeln automatisch in kleinere Fragmente oder aktiviert den Burst-Modus des DMA-Busses, um einen Stillstand des Rechenkerns zu verhindern.

3. **Betriebsmodi und ingenieurtechnische Entscheidungen auf Basis des Laufzeitvergleichs:**
   - **Optimaler deterministischer Modus ($`T_{\mathrm{refill}} \ll T_{\mathrm{execute}}`$, Sicherheitsmarge $\ge 2\times$):**
     Für einen typischen Wissensblock von $256\,\text{KB}$ beträgt die Ladezeit $`T_{\mathrm{refill}} = \frac{256\,\text{KB}}{1.2\,\text{GB/s}} \approx 218\,\mu\text{s}`$. Der systolische Abgleich eines komplexen Anforderungspakets gegen diese Regeln beansprucht $`T_{\mathrm{execute}} \approx 1{,}5\text{--}5{,}0\,\text{ms}`$. Die Zeitreserve liegt beim 7- bis 23-Fachen.
     **Systemaktion:** Die passive Speicherbank wird garantiert lange vor dem Abschluss der Berechnungen befüllt. Sobald der Prozessor bereit ist, schaltet der Hardware-Umschalter die Bänke in **1 Taktzyklus** um ($A \leftrightarrow B$), was eine durchgehende Effizienz von $\eta = 1{,}00$ und eine strikte Vorhersehbarkeit der Antwortzeit (Worst-Case Execution Time, WCET) garantiert.
   - **Grenzbereich der Buskonkurrenz ($`0{,}8 \cdot T_{\mathrm{execute}} < T_{\mathrm{refill}} \le T_{\mathrm{execute}}`$):**
     Die zeitliche Marge ist minimal. Jede temporäre Buskollision (externe Interrupts, Peripheriezugriffe) droht die Pipeline-Frist zu verletzen.
     **Systemaktion:** Der Dispatcher aktiviert eine prioritäre Bus-Arbitrierung (*DMA Channel Priority / QoS*) und weist dem Wissenstransferstrom die höchste Priorität gegenüber anderen Systemprozessen zu.
   - **Notfallmodus bei Bandbreitendefizit ($`T_{\mathrm{refill}} > T_{\mathrm{execute}}`$):**
     **Systemaktion:** Ein Hardware-Watchdog-Timer registriert den Pipeline-Stillstand. Der semantische Dispatcher aktiviert einen Degradationsmodus: Er drosselt temporär die Parallelität des Vektorkerns oder schneidet unkritische Defeater zweiter Ordnung ab, wodurch der stabile Zustand $`T_{\mathrm{refill}} \le T_{\mathrm{execute}}`$ wiederhergestellt wird.

### 4.3. Systolische Batch-Compliance-Audits (Batch Compliance Sweeps)

Für die Massenverifikation (beispielsweise die simultane Prüfung von $M = 1\,500$ Sicherheitsanforderungen gegen $K = 500$ normative Regeln) gilt:
- Anforderungen werden als gestreamter Vektorstrom eingespeist;
- elementare atomare Vergleiche werden in $1\text{--}4$ Hardware-Taktzyklen abgearbeitet;
- die Gesamtdauer der Hardwareprüfung von $750\,000$ Paaren beträgt weniger als **$30\text{--}50\,\text{ms}$**, was eine End-to-End-Antwortzeit von unter $1{,}0\,\text{s}$ für das Gesamtsystem sicherstellt.

---

## 5. Mehrquellenbasierte defeasible Wissensaggregation (ASPIC+ / Dung)

In hochkomplexen F&E-Ingenieurprojekten besitzt keine einzelne Wissensbasis ein absolutes Vollständigkeitsmonopol. Ein einzelner Domänendienst liefert unter Umständen Ergebnisse, die nicht vollständig eindeutig sind.

### 5.1. Trigger für Wissensunvollständigkeit (Uncertainty Triggers)

Der semantische Broker identifiziert die Notwendigkeit der Einbindung zusätzlicher Dienste anhand von drei strikten Kriterien:
1. **Niedriger diskreter Konfidenzwert ($Score < 70/100$):** Der Beleg entstammt einer nichtharmonisierten oder sekundären Quelle;
2. **Nichtimperativer normativer Status (Informative vs. Normative):** Die Regel wurde einem informativen Anhang oder Leitfaden entnommen (*Guideline / May / Informative Annex*), während der Ingenieur ein verbindliches Zertifizierungsurteil anfordert;
3. **Offener Defeasibility-Zustand (Ungrounded Defeater):** Die Regel enthält eine Ausnahmebedingung `unless condition_X`, in der lokalen Wissensbasis fehlt jedoch der Sachverhalt zu `condition_X`.

### 5.2. Föderierte Abfrage (Scatter-Gather) und argumentatives Schiedsverfahren

Wird ein Wissensdefizit detektiert, initiiert der Broker eine parallele Abfrage (**Scatter**) benachbarter Dienste im Cluster (beispielsweise eines Dienstes für Industriepräzedenzfälle oder allgemeine Zuverlässigkeitsnormen). Die von den verschiedenen Diensten zurückgelieferten Evidenzfragmente werden im Anschluss (**Gather**) in das formale Argumentationsframework **ASPIC+** [[2]](#src-2) auf Basis abstrakter Argumentationsstrukturen nach Dung [[1]](#src-1) überführt:

```mermaid
flowchart TD
    accTitle: Föderiertes argumentatives Evidenz-Schiedsverfahren auf Basis von ASPIC+
    accDescr: Drei Verarbeitungszweige föderierter Argumente: Korroboration, Entkräftung und Fail-Closed-Schutzschleuse bei normativer Kollision.

    ENGINE["<b>ASPIC+-ARGUMENTATIONSMASCHINE</b><br><i>Föderiertes Evidenz-Schiedsverfahren des Clusters (Scatter-Gather)</i>"]

    BRANCH1["<b>1. Korroboration von Evidenzen (Corroboration)</b><br>Unabhängige Dienste bestätigen übereinstimmende Schlussfolgerung<br><code>Conf(A₁ ∪ A₂) = 1 - (1 - Cred(A₁)) · (1 - Cred(A₂))</code><br><b>CONF_RES &gt; MAX(CONF_i) · Vertrauen wächst</b>"]

    BRANCH2["<b>2. Entkräftung von Behauptungen (Undercutting Defeat)</b><br>Verbindliche Norm eines höherrangigen Standards entkräftet Empfehlung<br><code>A_strong (SHALL) ➔ DEFEATS ➔ A_weak (MAY/SHOULD)</code><br><b>Eindeutiges normatives Verbot mit lückenlosem Beweis</b>"]

    BRANCH3["<b>3. Normative Kollision (Qualified Conflict)</b><br>Gleichrangige Standards enthalten widersprüchliche Imperative<br><b>FAIL-CLOSED GATE (ZHR = 1.00)</b><br><i>Verbot stochastischer Mutmaßungen · Eskalation an den Architekten</i>"]

    ENGINE ==>|Übereinstimmende unabhängige Argumente| BRANCH1
    ENGINE ==>|Konflikt von Normen unterschiedlicher Priorität| BRANCH2
    ENGINE ==>|Widersprüchliche gleichrangige Normen| BRANCH3

    classDef root fill:#ede7f6,stroke:#512da8,stroke-width:2px,color:#311b92;
    classDef corrob fill:#e8f5e9,stroke:#2e7d32,stroke-width:2px,color:#1b5e20;
    classDef defeat fill:#e1f5fe,stroke:#0288d1,stroke-width:2px,color:#01579b;
    classDef conflict fill:#ffebee,stroke:#c62828,stroke-width:2px,color:#b71c1c;

    class ENGINE root;
    class BRANCH1 corrob;
    class BRANCH2 defeat;
    class BRANCH3 conflict;
```

#### 5.2.1. Korroboration von Evidenzen (Corroboration)

Beim Audit sicherheitskritischer Systeme (beispielsweise Anforderungen nach ISO 26262 ASIL D [[5]](#src-5) oder DO-178C DAL A [[6]](#src-6)) kann die Evidenz eines einzelnen Domänendienstes ein nicht zu vernachlässigendes epistemisches Defizit oder ein unzureichendes Konfidenzniveau für ein bedingungsloses automatisches Urteil aufweisen. Der Korroborationsmechanismus eliminiert Subjektivität durch die Fusion unabhängiger Nachweise aus benachbarten Clusterknoten. Gelangen zwei oder mehr Dienste unabhängig voneinander zu einer übereinstimmenden Feststellung, steigt die aggregierte probabilistische Konfidenz in die Schlussfolgerung mathematisch gemäß dem Additionssatz für unabhängige Ereignisse an:

```math
\mathrm{Conf}(A_1 \cup A_2) = 1 - (1 - \mathrm{Cred}(A_1)) \cdot (1 - \mathrm{Cred}(A_2))
```

wobei:
- $`\mathrm{Cred}(A_i) \in [0, 1)`$ den diskreten Glaubwürdigkeitsgrad der Evidenz bezeichnet, die vom Domänendienst $`i`$ bereitgestellt wurde (bestimmt durch die Autorität der Wissensquelle und die Rückverfolgbarkeitstiefe);
- $`\mathrm{Conf}(A_1 \cup A_2) \in [0, 1)`$ das aggregierte probabilistische Vertrauen in die korroborierte Schlussfolgerung nach Zusammenführung der unabhängigen Argumente $`A_1`$ und $`A_2`$ darstellt.

**Praktische Anwendung und ingenieurtechnische Schlussfolgerungen:**
1. **Aufrufpunkt im Lebenszyklus:** Die Berechnung erfolgt durch den semantischen Broker in der Konsolidierungsphase (Gather-Phase), sobald die von mehreren Knoten gelieferten Nachweise bezüglich eines Zielprädikats übereinstimmen.
2. **Entscheidungskriterien anhand numerischer Schwellenwerte:**
   - **Schwelle für automatische Anforderungsfreigabe ($`\mathrm{Conf} \ge \tau_{\mathrm{accept}} = 0{,}95`$):** Der Nachweis gilt als erschöpfend. Das System vergibt automatisch das Urteil `ACCEPT`, bringt eine kryptographische Ed25519-Signatur an und überführt das verifizierte Artefakt ohne manuellen Prüfaufwand in das Zertifizierungspaket.
   - **Bereich für Expertenbegutachtung ($`0{,}80 \le \mathrm{Conf} < 0{,}95`$):** Das System generiert das qualifizierte Urteil `QUALIFIED`. Das Evidenzpaket wird als bedingt akzeptabel markiert und automatisch an den leitenden Architekten eskaliert (Mensch im Regelkreis, *Human-in-the-Loop*), versehen mit einer bytegenauen Kennzeichnung der Schwachstellen in der Argumentationskette.
   - **Bereich unzureichender Evidenz ($`\mathrm{Conf} < 0{,}80`$):** Das System blockiert die Freigabe mit dem Status `REFUSAL`. Die Anfrage wird aufgrund fehlender theoretischer Fundierung in der Wissensbasis zurückgewiesen.
3. **Numerisches Ingenieurbeispiel:**
   Ein Knoten für Architekturmuster-Analysen liefert die Aussage $`A_1`$ mit einer Glaubwürdigkeit von $`\mathrm{Cred}(A_1) = 0{,}85`$ (isoliert betrachtet reicht dies nicht für eine automatische ASIL-D-Zertifizierung aus, die einen Schwellenwert von $0{,}95$ erfordert). Ein benachbarter Knoten für Speicheraudits bestätigt die Sicherheit desselben Codefragments unabhängig durch das Argument $`A_2`$ mit der Glaubwürdigkeit $`\mathrm{Cred}(A_2) = 0{,}80`$. Das aggregierte Vertrauen berechnet sich zu:

```math
\mathrm{Conf}(A_1 \cup A_2) = 1 - (1 - 0.85) \cdot (1 - 0.80) = 1 - (0.15 \cdot 0.20) = 1 - 0.03 = 0.97
```

   Da $`0{,}97 \ge \tau_{\mathrm{accept}} = 0{,}95`$ gilt, übersteigt die wechselseitige Bestätigung zweier unabhängiger Quellen die Zulassungsschwelle und erübrigt eine zeitaufwendige manuelle Begutachtung.

#### 5.2.2. Entkräftung schwacher Behauptungen (Defeat / Undercutting)

In technischen Regelwerken werden normative Vorgaben nach ihrer Verbindlichkeit differenziert: zwingende Vorschriften (*Shall / Must*), Empfehlungen (*Should / Recommended*) und informative Hinweise (*May / Informative*).

Tritt während der föderierten Suche ein Widerspruch zwischen den Resultaten verschiedener Dienste auf, wendet das System keine probabilistische Mittelung an. Gemäß den Regeln der nichtmonotonen Inferenz in ASPIC+ entkräftet (*undercuts*) eine zwingende Norm höherer Priorität ($\mathrm{Rank} = \text{Shall}$) ein argumentatives Übergewicht einer Empfehlung oder eines Hinweises niedrigeren Ranges ($\mathrm{Rank} = \text{Should/May}$) bedingungslos. Das schwächere Argument wird aus dem finalen Begründungsgraphen eliminiert, während das Auditprotokoll die exakte Fundstelle des derogierenden Standardabschnitts festhält.

#### 5.2.3. Qualifizierte Verweigerung bei normativer Kollision (Qualified Conflict)

Ein kritischer Zustand tritt ein, wenn zwei voneinander unabhängige, harmonisierte Standards identischer Rechtskraft (beispielsweise Anforderungen an die funktionale Sicherheit nach ISO 26262 und Cybersicherheitsspezifikationen nach ISO/SAE 21434 auf der modalen Ebene *Shall*) zueinander inkompatible, sich gegenseitig ausschließende Imperative für dieselbe Architekturoperation formulieren.

In einer solchen Konstellation untersagt das Expertensystem gemäß den Richtlinien für vertrauenswürdige künstliche Intelligenz (insbesondere NIST AI RMF 1.0 [[7]](#src-7)) **kategorisch jedes stochastische Erraten eines Kompromisses** oder die Auswahl einer „am wahrscheinlichsten passenden“ Variante durch ein neuronales Netz. Das System aktiviert die **Fail-Closed-Schutzschleuse**:
1. Es generiert das typisierte Urteil `QUALIFIED_CONFLICT`;
2. Es konstruiert einen zweiseitigen Konfliktbaum unter Angabe der exakten Absätze beider Standards;
3. Es blockiert die automatische Artefaktfreigabe in der CI/CD-Pipeline und eskaliert den Konflikt an den leitenden Systemarchitekten zur Erstellung einer formellen technischen Begründung für eine Sicherheitsabweichung (*Safety Deviation Justification*).

---

## 6. Verteilte Triage von Wissensdefiziten und Härtung von Nachrichtenwarteschlangen

In der verteilten Architektur evidenzbasierter Expertensysteme, in der jeder Knoten dem kompromisslosen Prinzip des Blockierens bei unvollständigem Wissen folgt (Fail-Closed, $`\mathrm{ZHR} = 1{,}00`$), führt jeder fehlende Fakt, jede unaufgelöste normative Ambivalenz oder Regelkollision zur automatischen Erzeugung eines Wissensdefizit-Tickets (*Knowledge Gap Ticket*). Bei großflächigen industriellen Bereitstellungen (Hunderte von Bordsteuergeräten oder Tausende von kontinuierlichen CI/CD-Audit-Agenten) entsteht die systemische Bedrohung eines Ticketsturms bei Wissenslücken (*Gap-Ticket Storm*): Eine naive lineare Aggregation von Anfragen unter globalen Sperren führt zur Erschöpfung des Thread-Pools des Nachrichtenbrokers, zum Überlauf der Arbeitsspeicher-Warteschlangen und zu Denial-of-Memory-Zuständen. Zur Bewältigung dieser Krise stützt sich das verteilte Wissenssystem auf ein dreistufiges Subsystem der semantischen Triage.

### 6.1. SimHash-Fingerabdrücke und metrische Burkhard-Keller-Bäume (BK-Tree)

Um redundante Duplikate zu eliminieren und gleichartige Defizite zu einheitlichen Engineering-Aufgaben zusammenzufassen, wird jedes eingehende Defizit-Ticket tokenisiert und mittels des SimHash-Algorithmus [[8]](#src-8) zu einem 64-Bit-lokalitätssensitiven Fingerabdruck (*Locality-Sensitive Hash, LSH*) komprimiert. Im Gegensatz zu kryptografischen Hashfunktionen (SHA-256), bei denen die Änderung eines einzelnen Bits das Ergebnis vollständig verändert (Lawineneffekt), bewahrt SimHash die topologische Nähe der Eingangsvektoren im Hamming-Raum: Ähnliche Defizitformulierungen erhalten nahe beieinander liegende Binärvektoren.

Jedes $`i`$-te Bit des 64-Bit-Fingerabdrucks $`h(\mathbf{t}) \in \{0, 1\}^{64}`$ des Tickets $`\mathbf{t}`$ wird als gewichtete Summe der Merkmalsprojektionen berechnet:

```math
h_i(\mathbf{t}) = \frac{1 + \mathrm{sign}\left(\sum_{w \in \mathbf{t}} \mathrm{weight}(w) \cdot \bigl(2 \cdot b_i(w) - 1\bigr)\right)}{2}, \quad i \in [0, 63]
```

wobei:
- $`\mathbf{t} = \{w_1, w_2, \dots, w_k\}`$ die Menge der Token (Prädikate, Normenbezeichner, semantische Lexeme des Defizitkontexts) darstellt;
- $`\mathrm{weight}(w)`$ den Gewichtungskoeffizienten des Tokens bezeichnet (beispielsweise TF-IDF-Score oder Informationsgewicht des ontologischen Konzepts, $`\mathrm{weight}(w) > 0`$);
- $`b_i(w) \in \{0, 1\}`$ das $`i`$-te Bit des 64-Bit-Basishashes des Tokens $`w`$ angibt (generiert durch die schnelle nicht-kryptografische Funktion xxHash64);
- $`\mathrm{sign}(x)`$ die Vorzeichenfunktion ist ($`+1`$ für $`x \ge 0`$ und $`-1`$ für $`x < 0`$).

Zwei Defizit-Tickets $`\mathbf{t}_1`$ und $`\mathbf{t}_2`$ werden demselben strukturellen Wissensdefizit-Cluster zugeordnet, wenn ihr bitweiser Hamming-Abstand den Cluster-Schwellenwert nicht überschreitet:

```math
D_{\mathrm{Hamming}}\bigl(h(\mathbf{t}_1), h(\mathbf{t}_2)\bigr) = \sum_{i=0}^{63} \Bigl(h_i(\mathbf{t}_1) \oplus h_i(\mathbf{t}_2)\Bigr) \le \tau_{\mathrm{gap\_cluster}} = 3
```

wobei $`\oplus`$ die bitweise XOR-Operation bezeichnet und der Schwellenwert $`\tau_{\mathrm{gap\_cluster}} = 3`$ einer inhaltlichen Übereinstimmung von über $95\,\%$ entspricht.

Um eine lineare Durchmusterung bestehender Cluster ($`O(N)`$) zu vermeiden, die während eines Sturms zum Stillstand der Verarbeitung führt, werden die gespeicherten Cluster-Fingerabdrücke in einem metrischen Burkhard-Keller-Baum (*BK-Tree*) [[9]](#src-9) organisiert. Da der Hamming-Abstand die Metrik-Axiome (insbesondere die Dreiecksungleichung) erfüllt, erfolgt die Suche nach dem nächstgelegenen Cluster in logarithmischer Zeit $`O(\log N)`$, was den Rechenaufwand um drei Größenordnungen senkt.

#### 6.1.1. Geschlossener Kreis der ingenieurtechnischen Berechnung der Wissensdefizit-Triage

1. **Laufzeitsteuerung und Rechenfluss (Control Flow & Runtime Decisions):**
   - Stellt ein Agent das Fehlen einer Regel fest, generiert er ein Defizit-Ticket und übergibt dieses an das lokale Triage-Modul. Gilt $`D_{\mathrm{Hamming}}(h(\mathbf{t}), h_{\mathrm{existing}}) \le 3`$, wird das neue Ticket nicht im zentralen NATS-Ereignisbus publiziert, sondern lokal als Ereigniszähler des Clusters aggregiert (`IncClusterHit(ClusterID)`);
   - Ist der Hamming-Abstand zu allen bekannten Knoten des BK-Baums $`D_{\mathrm{Hamming}} > 3`$, registriert das System einen neuen eindeutigen Wissensdefizit-Cluster, fügt ihn in den BK-Baum ein und publiziert das Ereignis `vida.gap.discovered` auf dem NATS-Bus zur weiteren Untersuchung durch Wissensingenieure.
2. **Hardware-Dimensionierung (Hardware Dimensioning):**
   - Jeder Knoten des BK-Baums belegt $`64\,\text{Byte}`$ im Arbeitsspeicher (64-Bit-Schlüssel, Trefferzähler, Zeiger auf Kindknoten von $`0`$ bis $`64`$). Für einen Pool von $`100\,000`$ eindeutigen Defizit-Clustern benötigt die gesamte BK-Baumstruktur lediglich $`6{,}4\,\text{MB}`$ RAM. Dies erlaubt es, die Struktur vollständig im L3-Cache des Prozessors zu halten und Zugriffsverzögerungen auf den externen DRAM-Speicher zu eliminieren.
3. **Praktisches numerisches Rechenbeispiel (Worked Numerical Example):**
   - Während eines nächtlichen Testlaufs eines autonomen Manövriersystems generierten $`1\,000`$ virtuelle Agenten $`50\,000`$ Defizit-Tickets infolge einer fehlenden Spezifikation für das Bremsverhalten bei Gierschleudern auf Glatteis.
   - Durch die Auswertung von $`D_{\mathrm{Hamming}} \le 3`$ im BK-Baum fasste das System $`50\,000`$ gleichartige Anfragen in exakt $`2`$ Wurzel-Defizit-Cluster zusammen (Deduplizierungsrate $`99{,}996\,\%`$). Die Last auf dem Kommunikationsbus sank von $`50\,000`$ auf $`2`$ Nachrichten, wodurch ein Zusammenbruch der Aufgabenwarteschlange abgewendet wurde.

### 6.2. Ingress-Validierung: Schutz vor Entartung des Merkmalsraums und strikte Deserialisierung

Im realen verteilten Betrieb sehen sich asynchrone Message-Broker (NATS JetStream, Apache Kafka) mit einem gravierenden latenten Fehler konfrontiert: dem **schleichenden Kollaps des metrischen Hashraums** (*Contract Drift Collapse*) [[7]](#src-7). Übermittelt ein sendender Agent infolge eines Serialisierer-Versionswechsels ein Ticket mit abweichenden Feldbezeichnern (beispielsweise `ticket_id` anstelle des vom Empfänger erwarteten `TicketID` oder `context` statt `Query`), belässt ein toleranter JSON-Deserialisierer diese Felder stillschweigend als leere Zeichenketten (`""`).

Der Text-Vektorisierer erhält eine leere Eingabe. Der SimHash-Algorithmus über einer leeren Tokenmenge liefert deterministisch einen konstanten Nullvektor:

```math
\mathrm{SimHash}(\emptyset) = \mathrm{0x0000000000000000}
```

Infolgedessen erhalten grundverschiedene Wissensdefizite aus völlig unterschiedlichen Subsystemen (vom Kalibrierungsfehler eines Stromsensors bis zur juristischen Kollision in Richtlinien) den identischen Fingerabdruck `0x0000000000000000`. Da der Hamming-Abstand zwischen ihnen $`D_{\mathrm{Hamming}} = 0 \le 3`$ beträgt, verschmilzt das System völlig unzusammenhängende Fehler zu einem gigantischen, entarteten Pseudocluster und lähmt damit die Arbeit der Wissensingenieure.

Zur Beseitigung dieser Schwachstelle wird ein kombinierter architektonischer Schutzmechanismus (*Degeneracy Guard*) etabliert:

1. **Invariante minimaler lexikalischer Komplexität (Distinct Tokens Invariant):**
   Vor der Berechnung der SimHash-Faltung wird die hinreichende Entropie des Merkmalsraums überprüft:

```math
\mathrm{DistinctTokens}(T) \ge K_{\mathrm{min}} \quad (K_{\mathrm{min}} = 2)
```

   Ist die Anzahl eindeutiger semantischer Token $`\mathrm{DistinctTokens}(T) < 2`$ oder die Länge des Eingabetextes null, bricht der Aufruf unmittelbar mit dem Fehler `ErrDegenerateFeatureSpace` ab. Die Werte `0x0000000000000000` und `0xFFFFFFFFFFFFFFFF` werden als Anomalie-Marker reserviert und hardwareseitig gegen die Aufnahme in den BK-Baum gesperrt;

2. **Strikte zweiseitige Vertragsabgleichung (Strict Unmarshal):**
   Systemnachrichten-Deserialisierer im verteilten Bus werden auf das absolute Verbot unbekannter oder fehlender Pflichtfelder konfiguriert:
   ```go
   decoder := json.NewDecoder(streamReader)
   decoder.DisallowUnknownFields()
   ```
   Nachrichten mit fehlenden oder korrumpierten Feldern werden unverzüglich abgewiesen und in eine Quarantäne-Warteschlange für unzustellbare Nachrichten (*Dead-Letter Queue, DLQ*) umgeleitet;

3. **Standardisierung hierarchischer Multi-Token-Busmasken:**
   Im NATS-Broker filtert eine Sternchen-Maske (`vida.tasks.*`) **strikt genau ein Token** der Hierarchie. Um den Nachrichtendurchlauf von hierarchisch tief verschachtelten Knoten zu garantieren, werden alle System-Abonnements auf den Operator für vollständige Teilbäume `>` standardisiert:

```math
\mathrm{SubSubject} = \text{"vida.tasks.>"} \quad (\text{anstelle des fehlerhaften } \text{"vida.tasks.*"})
```

#### 6.2.1. Geschlossener Kreis der ingenieurtechnischen Ingress-Vertragsüberprüfung

1. **Laufzeitsteuerung und Rechenfluss (Control Flow & Runtime Decisions):**
   - Enthält eine eingehende Nachricht nicht die geforderten Schema-Pflichtfelder oder gilt $`\mathrm{DistinctTokens} < 2`$, erzeugt das Triage-Modul das Ereignis `ALERT_CONTRACT_DEGENERACY` und verwirft das Ticket in die DLQ. Der Ausführungsfluss verhindert, dass korrumpierte Daten in den BK-Baum gelangen;
   - Ist der Vertrag valide und gilt $`\mathrm{DistinctTokens} \ge 2`$, berechnet das System den regulären 64-Bit-SimHash-Fingerabdruck und leitet ihn an den Clustering-Algorithmus weiter.
2. **Hardware-Dimensionierung (Hardware Dimensioning):**
   - Der DLQ-Quarantänepuffer wird auf einem separaten SSD-Speicher mit einer Größenbegrenzung von $`10\,\text{GB}`$ bereitgestellt, gekoppelt mit einer Alarmierung des Bereitschaftsingenieurs bei einer Belegung von über $`50\,\%`$;
   - Die Prüfung mittels `DisallowUnknownFields` erfordert nicht mehr als $`1{,}2\,\mu\text{s}`$ pro Nachricht, wodurch der Knotendurchsatz von über $`80\,000`$ Nachrichten/s auf einem einzelnen CPU-Kern aufrechterhalten wird.
3. **Praktisches numerisches Rechenbeispiel (Worked Numerical Example):**
   - Nach einem Firmware-Update des Sensormoduls eines Flugreglers wurde das Kontextfeld versehentlich von `context_spec` in `ctx_spec` umbenannt.
   - Das Schutz-Gateway `decoder.DisallowUnknownFields()` erkannte das unbekannte Feld `ctx_spec` innerhalb von $`0{,}8\,\mu\text{s}`$ und leitete die Anfrage unter Ausgabe eines Notfallprotokolls in die DLQ um. Dadurch wurde die fehlerhafte Entartung von $`15\,000`$ Tickets in den Null-Fingerabdruck `0x0` verhindert und die Reinheit des semantischen Defizit-Index im Unternehmens-Repository bewahrt.

### 6.3. Autonomer lokaler WAL-Puffer und Lame-Duck-Verhalten bei Netztrennung

In industriellen verteilten Systemen ist die Verbindung zwischen Edge-Knoten (Workstations, Bordcomputern) und dem zentralen Wissens-Broker nicht ausfallfrei: Elektromagnetische Störungen, Netzwerkumschaltungen oder Serverwartungen führen zu temporären Knotenisolationen (*Network Partitions* gemäß CAP-Theorem). Erkennt ein Edge-Agent in einem solchen Moment ein Wissensdefizit, führt ein Fehlen von Pufferungsmechanismen entweder zum Blockieren des Arbeitsprozesses oder zum unbemerkten Verlust sicherheitskritischer Wissenslücken.

Zur Gewährleistung einer lückenlosen Auditierung wird der Edge-Knoten mit einem autonomen lokalen Write-Ahead-Log (*WAL*) ausgestattet:
1. **Verschlüsselter WAL-Ringpuffer:** Defizit-Tickets werden unverzüglich in einem lokalen Festplattenjournal protokolliert, das mit einem asymmetrischen Sitzungsschlüssel geschützt ist. Das Journal unterliegt einer festen Festplattenquote (beispielsweise $`500\,\text{MB}`$). Bei Erreichen der Quote wird eine Ringrotation aktiviert, bei der ältere redundante Einträge verdrängt werden, eindeutige kritische Urteile jedoch erhalten bleiben;
2. **Wiederverbindung und FIFO-Drainage:** Sobald der Netzwerk-Stack die Wiederherstellung der Verbindung zum Wissens-Broker registriert, liest ein Hintergrund-Agent das lokale WAL sequenziell aus und publiziert die akkumulierten Tickets nach dem FIFO-Prinzip paketweise, ohne den primären Invarianten-Prüfzyklus zu blockieren;
3. **Geregeltes Außerbetriebnahme-Protokoll (Lame-Duck-Modus):**
   Bei einem geplanten Update oder Herunterfahren eines Clusterknotens wird ein zweistufiges Sicherheitsstopp-Protokoll aktiviert:
   - Der Knoten deregistriert sich umgehend im dynamischen Register des Brokers (`UNREGISTER_PENDING`) und stoppt die Annahme neuer Anfragen von Thin Clients;
   - Alle aktiven Inferenzsitzungen werden regulär zu Ende geführt, und intermediäre Defizit-Warteschlangen werden vollständig auf die Festplatte in das WAL synchronisiert und an einen benachbarten Standby-Knoten übergeben;
   - Erst nach Bestätigung des vollständigen Pufferabgleichs beendet der Prozess seine Ausführung ordnungsgemäß (`SHUTDOWN_CLEAN`).

## 7. Unveränderliches Entscheidungsbuch (Immutable Decision Ledger) auf Merkle-Bäumen und WORM-Speichern

In sicherheitskritischen Anwendungsfeldern von Expertensystemen – wie der funktionalen Sicherheit autonomer Fahrzeuge (ISO 26262 ASIL D), der Zertifizierung von Bordsoftware (DO-178C DAL A) und der medizinischen Diagnostik – ist die **Nichtabstreitbarkeit von Entscheidungen (Non-Repudiation)** eine zentrale rechtliche und technische Kernanforderung. Herkömmliche relationale Datenbanken (PostgreSQL, SQLite) oder textuelle Betriebssystemprotokolle sind für diese Aufgabe ungeeignet: Jeder privilegierte Serveradministrator könnte nachträglich Datensätze manipulieren, Urteile fälschen oder das Vorliegen von Sicherheitsverstößen bei retrospektiven Unfalluntersuchungen verschleiern. Andererseits verursacht der naive Versuch, das gesamte Gigabyte-Volumen an Beweisketten in eine öffentliche Blockchain zu schreiben, unzulässige Latenzen (Sekunden statt Mikrosekunden), immense finanzielle Kosten und die Offenlegung vertraulichen Firmenwissens.

Die evidenzbasierte Architektur löst dieses Dilemma durch den Aufbau eines **unveränderlichen Entscheidungsbuchs (Immutable Decision Ledger)** auf Basis lokaler Merkle-Bäume und einer verteilten Verankerung von Kontrollwurzeln.

```mermaid
flowchart TD
    accTitle: Architektur des unveränderlichen Entscheidungsbuchs auf Basis von Merkle-Bäumen
    accDescr: Ein vierkomponentiges Entscheidungs-Beweisbündel wird gehasht, mit Ed25519 signiert, in einem Merkle-Baum zusammengefasst und periodisch in einem WORM-Speicher verankert.

    subgraph Decision["Entscheidungsfindung des Expertensystems"]
        INP["<b>Eingangsdaten</b><br/>Sensorvektor I"]
        KBV["<b>Version der Wissensbasis</b><br/>Hash des Regelkorpus V_KB"]
        PRF["<b>Beweisgraph</b><br/>ASPIC+ / GSN-Zielbaum"]
        VRD["<b>Systemurteil</b><br/>ACCEPT / REFUSAL"]
    end

    subgraph Bundle["Kryptografisches Entscheidungsbündel (Evidence Bundle)"]
        H_INP["H(InputData)"]
        H_KBV["H(KB_Version)"]
        H_PRF["H(ProofGraph)"]
        H_VRD["H(Verdict)"]
        SIG["<b>Digitale Signatur Ed25519</b><br/>Sign(H_INP || H_KBV || H_PRF || H_VRD)"]
        
        H_INP --- H_KBV --- H_PRF --- H_VRD --> SIG
    end

    subgraph MerkleLayer["Hierarchischer Merkle-Entscheidungsbaum"]
        LEAF1["Entscheidungsbündel 1"]
        LEAF2["Entscheidungsbündel 2"]
        LEAF3["Entscheidungsbündel 3"]
        LEAF4["Entscheidungsbündel 4"]
        
        NODE12["Knoten Hash(1 || 2)"]
        NODE34["Knoten Hash(3 || 4)"]
        M_ROOT["<b>Merkle-Wurzel (Merkle Root)</b><br/>Hash(Node12 || Node34)"]
        
        LEAF1 --> NODE12
        LEAF2 --> NODE12
        LEAF3 --> NODE34
        LEAF4 --> NODE34
        NODE12 --> M_ROOT
        NODE34 --> M_ROOT
    end

    subgraph Anchoring["Unveränderlicher Speicher (Anchoring Layer)"]
        WORM["<b>Hardware-WORM-Speicher</b><br/>Write Once Read Many<br/>oder Konsortial-Distributed-Ledger"]
    end

    INP --> H_INP
    KBV --> H_KBV
    PRF --> H_PRF
    VRD --> H_VRD

    SIG --> LEAF1
    M_ROOT -->|"periodische Epochen-Verankerung"| WORM

    classDef dec fill:#e3f2fd,stroke:#1565c0,stroke-width:2px,color:#0d47a1;
    classDef bun fill:#fff8e1,stroke:#f57f17,stroke-width:2px,color:#e65100;
    classDef mer fill:#e8f5e9,stroke:#2e7d32,stroke-width:2px,color:#1b5e20;
    classDef anc fill:#ede7f6,stroke:#512da8,stroke-width:2px,color:#311b92;

    class INP,KBV,PRF,VRD dec;
    class H_INP,H_KBV,H_PRF,H_VRD,SIG bun;
    class LEAF1,LEAF2,LEAF3,LEAF4,NODE12,NODE34,M_ROOT mer;
    class WORM anc;
```

### 7.1. Decision Evidence Bundle und kryptografische Bindung

Jedes Inferenzurteil wird vor der Übermittlung an den Anfragenden in ein kryptografisch integres, binäres Beweisbündel (*Decision Evidence Bundle*) verpackt, welches vier miteinander verknüpfte Hash-Deskriptoren umfasst:

```math
\mathrm{EvidenceBundle} = \mathrm{Sign}_{\mathrm{Ed25519}}\Bigl(H(I) \mathbin{\Vert} H(V_{\mathrm{KB}}) \mathbin{\Vert} H(G_{\mathrm{proof}}) \mathbin{\Vert} H(\mathrm{Verdict})\Bigr)
```

wobei:
- $`H(I)`$ den kryptografischen SHA-256-Hash des Eingangsvektors der Beobachtungen, der Telemetrie oder des normativen Abfragetextes $`I`$ bezeichnet;
- $`H(V_{\mathrm{KB}})`$ den Hash der aktuellen Version der Wissensbasis, der normativen Ontologien und Defeater darstellt (gewährleistet die bytegenaue Fixierung des Regelzustands zum Inferenzzeitpunkt);
- $`H(G_{\mathrm{proof}})`$ den Hash der Struktur des ASPIC+-Argumentationsgraphen oder des GSN-Zielbeweisbaums angibt (einschließlich aller aktivierten Prädikate und entkräfteten Zweige);
- $`H(\mathrm{Verdict})`$ den Hash des resultierenden Urteils (`ACCEPT`, `REFUSAL`, `QUALIFIED`) bildet;
- $`\mathbin{\Vert}`$ die binäre Konkatenationsoperation darstellt;
- $`\mathrm{Sign}_{\mathrm{Ed25519}}`$ die digitale Signatur mit dem asymmetrischen Schlüssel des Rechenknotens nach dem Ed25519-Verfahren (RFC 8032) ist.

Diese Struktur liefert den mathematischen Beweis: Weder die Eingangsparameter noch die angewandten Regelversionen oder das Urteil selbst können nachträglich manipuliert werden, ohne die Gültigkeit der digitalen Signatur zu zerstören.

### 7.2. Hierarchischer Merkle-Baum und periodische WORM-Verankerung

Um einen hohen Durchsatz sicherzustellen, werden Entscheidungsbündel nicht einzeln in langsame externe Register geschrieben. Stattdessen werden sie im Arbeitsspeicher des Knotens in Echtzeit zu lokalen binären Merkle-Bäumen [[10]](#src-10) aggregiert:

```math
\mathrm{Root}_{\mathrm{Merkle}} = \mathrm{Hash}\bigl(\mathrm{Node}_{\mathrm{left}} \mathbin{\Vert} \mathrm{Node}_{\mathrm{right}}\bigr)
```

Alle $`T_{\mathrm{epoch}}`$ Sekunden (oder nach der Akkumulation von $`K_{\mathrm{batch}}`$ Entscheidungen, beispielsweise $`1\,000`$ Bündeln) wird die berechnete Baumwurzel $`\mathrm{Root}_{\mathrm{Merkle}}`$ in einem geschützten WORM-Speicher (*Write Once Read Many*, optisches Medium oder Cloud-Speicher mit S3 Object Lock) abgelegt oder in einem Konsortial-Ledger registriert.

Um die Authentizität einer einzelnen historischen Entscheidung gegenüber einem Auditor nachzuweisen, genügt die Bereitstellung eines kompakten Audit-Pfades (*Merkle Audit Trail*) der Länge $`\lceil \log_2 N \rceil`$ Hashes. Der Auditor verifiziert die Kette vom Blatt bis zur öffentlich verankerten Wurzel in wenigen Mikrosekunden, ohne den gesamten Datenbestand einlesen zu müssen.

#### 7.2.1. Geschlossener Kreis der ingenieurtechnischen Entscheidungsbuch-Verifikation

1. **Laufzeitsteuerung und Rechenfluss (Control Flow & Runtime Decisions):**
   - Nach der Urteilsfindung erzeugt der Wissensdienst ein `EvidenceBundle` und fügt dieses als Blatt in den aktuellen Epochen-Baum ein;
   - Weicht bei der Validierung einer historischen Abfrage der berechnete Hash ab ($`\mathrm{Root}_{\mathrm{computed}} \ne \mathrm{Root}_{\mathrm{WORM}}`$), löst das System einen Alarm höchster Priorität `AUDIT_CORRUPTION_DETECTED` aus und sperrt den Zertifizierungsstatus des Knotens bis zur Klärung der Datenmanipulation;
   - Stimmen die Wurzeln überein, stellt der Audit-Inspektor ein Zertifikat der rechtlichen und technischen Nichtabstreitbarkeit des Urteils aus (`NON_REPUDIATION_VERIFIED`).
2. **Hardware-Dimensionierung (Hardware Dimensioning):**
   - Die Berechnung eines Baumknotens auf Basis von SHA-256 oder BLAKE3 beansprucht auf einem modernen Prozessor mit SHA-NI-Hardwarebeschleunigung weniger als $`0{,}15\,\mu\text{s}`$;
   - Für einen Cluster mit einem Durchsatz von $`10\,000`$ Entscheidungen/s erzeugt eine Epoche von $`T_{\mathrm{epoch}} = 10\,\text{s}`$ einen Merkle-Baum mit $`100\,000`$ Blättern. Die Größe des Audit-Pfades beträgt lediglich $`\lceil \log_2 100\,000 \rceil = 17`$ Hashes ($`17 \times 32 = 544\,\text{Byte}`$), was den Netzwerkverkehr bei Zertifizierungsaudits minimiert.
3. **Praktisches numerisches Rechenbeispiel (Worked Numerical Example):**
   - Während eines ISO-26262-Audits eines automobilen ASIL-D-Systems entstand ein Verdacht hinsichtlich eines vor 6 Monaten gefällten Notbremsurteils.
   - Das Expertensystem lieferte das kryptografische Entscheidungsbündel zusammen mit den 17 Hashes des Merkle-Pfades.
   - Ein unabhängiger Verifizierer rekonstruierte die Baumwurzel in 17 Hash-Schritten ($`2{,}5\,\mu\text{s}`$) und glich sie mit der im WORM-Register hinterlegten Wurzel ab. Die bitgenaue Übereinstimmung bewies zweifelsfrei, dass weder die Sensordaten des Schleudervorgangs noch die Version der Bremsregeln nachträglich modifiziert wurden, wodurch der Hersteller vor Haftungsansprüchen geschützt wurde.

---

## 8. Ingenieurtechnische Entscheidungsmatrix: Skalierung nach Projektbudget

Der wesentliche Nutzen des theoretischen Referenzprototyps besteht darin, dass er einen **Entwicklungspfad** vorgibt und nicht eine starre Einkaufsliste proprietärer Hardware vorschreibt. Der Entwicklungsingenieur kann das System schrittweise auf derjenigen Hardware implementieren, die im Rahmen des aktuellen Budgets seines Forschungs- oder Industrieprojekts zur Verfügung steht:

| Infrastrukturstufe | Hardwareklasse | Semantischer Broker und System 1 | Deterministischer Inferenzkern (System 2) | Speicher und Speicherhierarchie | Budgetkontext und Anwendungsbereich |
| :--- | :--- | :--- | :--- | :--- | :--- |
| **Stufe 1: Laborprüfstand (Lab Prototype)** | Einzelner PC oder Entwickler-Workstation | Schlankes lokales SLM (3B–7B) auf CPU / Consumer-GPU | Software-Kern (Software EVM) als Bibliothek im selben Prozess | Dateispeicher mit direktem `mmap` in den Arbeitsspeicher des OS | **Minimales Budget.** Machbarkeitsstudien, Master- und Doktorarbeiten, Regelprüfung auf kleinen Datensätzen. |
| **Stufe 2: Verteilte Arbeitsgruppe (Team R&D Cluster)** | Dedizierter lokaler Server + schlanke Clients (PC/Smartphones) | Dedizierter Server mit 14B-Modell; Zugriff via REST/gRPC-Broker | Software-Inferenz-Daemon im Speicher mit IPC-Socket | Caching im Server-RAM (32–64 GB), Streaming-Bereitstellung an Clients | **Mittleres Budget.** Entwicklungsteams von 10–50 Ingenieuren, Integration in Corporate CI/CD, tägliche Code-Audits. |
| **Stufe 3: Industrieller Großkomplex (Enterprise / ASIL-D Safe)** | Server-Cluster + Hardware-Karten für deterministischen Schutz | Heterogener Modellpool (14B–32B) auf GPU/NPU mit Circuit Breaker | Spezialisierte Steckkarten mit deterministischem Silizium-Rechenkern | Hardware-Doppelpufferung (Ping-Pong BRAM) + Hochgeschwindigkeits-DMA-Bus | **Hohes / branchenspezifisches Budget.** Autonomes Fahren, Avionik, kritische Energieinfrastruktur, TÜV-Zertifizierungslabore. |

---

## 9. Pro et Contra — Virtueller Dialog

Zu dieser virtuellen intellektuellen Diskussion wurden herausragende Architekten verteilter Systeme eingeladen, deren grundlegende Theoreme und Standardwerke die moderne Landschaft globaler IT-Infrastrukturen prägen: Martin Kleppmann (Wissenschaftler an der University of Cambridge, Verfasser der wegweisenden Monografie *Designing Data-Intensive Applications* und führender Forscher im Bereich konfliktfreier replizierter Datentypen — CRDTs sowie lokaler Anwendungsarchitekturen — Local-First Software) und Eric Brewer (Professor Emeritus an der University of California, Berkeley, Vice President of Infrastructure bei Google und Begründer des fundamentalen CAP-Theorems). Ihre bahnbrechenden Erkenntnisse über die physikalischen Grenzen verteilter Konsensbildung, die unvermeidlichen Kompromisse zwischen Konsistenz und Verfügbarkeit bei Netzwerktrennungen sowie die mathematischen Gesetze der Zustandsreplikation bilden den optimalen Prüfstein für die vom Autor konzipierte epistemische serviceorientierte Architektur (Epistemic SOA) und das unveränderliche Merkle-Baum-Entscheidungsprotokoll.

*Disclaimer:* Diese Diskussion ist eine heuristische virtuelle Simulation im wissenschaftlich-pädagogischen Kontext der Monografie. Die Positionen und Argumente der Kontrahenten wurden auf der Basis ihrer Veröffentlichungen, Fachaufsätze und Monografien synthetisiert; sie erheben keinen Anspruch auf wörtliche Zitate und dienen der kritischen Rationalisierung (Popper'schen Falsifikation) der ingenieurtechnischen Hypothesen des Autors.

**Runde 1: Martin Kleppmann kritisiert Latenzen und Skalierbarkeit verteilter Merkle-Baum-Entscheidungsprotokolle**

- **Martin Kleppmann:** Herr Fedchyk, bei der Erforschung verteilter Datenspeicher und lokalfokussierter Softwaresysteme (Local-First Software) zeigt sich immer wieder, dass der Versuch, ein verteiltes unveränderliches Entscheidungsprotokoll (Immutable Decision Ledger) auf Basis von Merkle-Bäumen und lückenlosem kryptografischem Konsens zu errichten, unweigerlich an massiven Skalierungsengpässen, Koordinationslatenzen und Schema-Migrationsproblemen scheitert. Wenn jeder Dienst in Ihrer Epistemic SOA verpflichtet ist, jedes Zwischenurteil, jedes Beweis-Tupel und jede Statusänderung im ASPIC+-Argumentationsgraphen in einem gemeinsamen kryptografischen Ledger zu protokollieren, schaffen Sie einen globalen Flaschenhals. Unter hoher Last führen das ständige Neuberechnen der Merkle-Wurzeln, die Verifikation von Ed25519-Signaturen und die netzwerkweite Synchronisation zu drastischen Latenzspitzen (Latency Spikes), die das System für harte Echtzeitanforderungen vollständig disqualifizieren. Wie verhindert Ihre Architektur, dass das unveränderliche Audit-Log das gesamte verteilte System lahmlegt?

- **Mykola Fedchyk:** Herr Kleppmann, Ihre Abhandlung über datenintensive Anwendungen warnt zu Recht vor der naiven Einführung globaler synchroner Transaktions-Logs, und wir haben diese Warnung bei der Hardware- und Pipeline-Dimensionierung konsequent berücksichtigt. In der Epistemic SOA ist das Entscheidungs-Ledger weder eine monolithische synchrone Blockchain noch ein zentraler Engpass. Wir setzen auf eine lokal-autonome, asynchrone Pufferarchitektur: Jeder Dienst schreibt Urteile zunächst in einen lokalen, ringförmig organisierten Speicherpuffer (Ring Buffer) in geschützter Hardware und liefert das signierte Beweistupel mit Sub-Millisekunden-Latenz ($\le 200\,\mu\text{s}$) an den Aufrufer zurück. Die Aggregation der Merkle-Wurzeln erfolgt vollständig entkoppelt über einen asynchronen Hintergrund-Batcher in festen Zeitfenstern (z. B. $100\,\text{ms}$). Überdies nutzen wir für die Synchronisation zwischen Knoten exakt Ihren methodischen Ansatz — zustandsbasierte CRDT-Halbverbände: Die Zusammenführung von epistemischen Defizitmeldungen und Argumentationsständen erfolgt deterministisch und kommutativ ohne zentralen Koordinator, was eine monotone Zustandskonvergenz ohne Blockierung des Laufzeitkerns garantiert.

**Runde 2: Eric Brewer attackiert das Überleben der epistemischen Architektur unter dem CAP-Theorem und Netzwerktrennungen**

- **Eric Brewer:** Herr Fedchyk, mit der Formulierung des CAP-Theorems habe ich ein unumstößliches Naturgesetz verteilter Systeme bewiesen: Tritt eine unvermeidliche Netzwerktrennung (Network Partition) auf, muss das System zwingend zwischen starker Konsistenz (Consistency) und Verfügbarkeit (Availability) wählen. In Ihrer Epistemic SOA propagieren Sie eine Aufteilung in Thin Clients, semantische Broker und föderierte Domänendienste über heterogene Netze. Doch was geschieht, wenn die Verbindung zwischen dem Broker und einem Domänencluster abreißt? Wählen Sie Verfügbarkeit (AP), fällen die isolierten Dienste Entscheidungen auf der Grundlage veralteter oder fragmentierter Daten, was Ihre fundamentale Invariante $ZHR = 1{,}00$ zerstört und zu widersprüchlichen juristischen Urteilen führt. Wählen Sie strikte Konsistenz (CP), legt die Netzwerktrennung die gesamte Anlage lahm. Behauptet Ihre epistemische SOA nicht insgeheim, die unnachgiebigen Gesetze des CAP-Theorems auf magische Weise außer Kraft setzen zu können?

- **Mykola Fedchyk:** Herr Brewer, das CAP-Theorem ist für unseren verteilten Laufzeitstapel ein physikalisches Gesetz, und wir versuchen keineswegs, es zu umgehen — wir treffen eine bewusste, sicherheitsgerichtete Architekturentscheidung im Sinne funktionaler Sicherheit (Safety-First Partitioning). Kommt es zu einer Netzwerktrennung, erteilt das vorgeschlagene System einer naiven Verfügbarkeit (AP) bei allen sicherheitskritischen Entscheidungen (ISO 26262 ASIL C/D) eine strikte Absage. Das System agiert kompromisslos CP-orientiert: Reißt die Verbindung zu den autoritativen Domänendiensten ab, rät der semantische Broker niemals heuristisch, sondern aktiviert deterministisch das Fail-Closed-Protokoll `DISP_REFUSE` mit dem Fehlercode `EP_NET_PARTITION`. Für unkritische Leseanfragen existiert jedoch ein hermetischer Ausweichmodus: Thin Clients oder Edge-Gateways operieren lokal autonom auf Basis kryptografisch signierter ZKP4-Binärpakete. Solange das hardwaregebundene Ticket innerhalb seiner definierten epistemischen Gültigkeitsepoche (Epistemic Epoch) verbleibt, verfügt der lokale Knoten über absolute Konsistenz und arbeitet ohne das geringste Risiko semantischer Widersprüche weiter.

## Fazit

1. **Isolation von Zuständigkeitsebenen in der Epistemic SOA:** Die verteilte Architektur evidenzbasierter Expertensysteme überwindet das fundamentale Skalierungsdilemma. Die Aufteilung in eine Thin-Client-Schicht, einen semantischen Koordinationsbroker und föderierte Domänendienste isoliert Wahrheitskontexte und entkoppelt das System vollständig von den Hardwaregrenzen des Benutzerendgeräts.
2. **Garantie der Invariante $`ZHR = 1.00`$ auf dem Thin Client:** Der Interaktions-Client hält keine lokalen Wissensbasen vor ($`S_{\mathrm{client}} = 0`$) und führt keine logische Inferenz durch. Seine Aufgabe beschränkt sich auf das Rendering des kryptographisch signierten Evidenztupels $\langle\text{Verdict}, \text{EvidenceMap}, \text{Signature}\rangle$, wodurch Rechtsverbindlichkeit und absolute Halluzinationsfreiheit selbst auf Mobilgeräten gewahrt bleiben.
3. **Ontologiebasiertes Abfrage-Routing:** Der semantische Broker ersetzt naives URL-Routing durch eine tiefgehende Inhaltsanalyse. Der syntaktisch-lexikalische Präprozessor extrahiert Domänenmarker, und das dynamische Dienstregister weist den primären Evaluator auf Basis des normierten semantischen Affinitätskoeffizienten $`\mathrm{Affinity}(\mathcal{S}_{\mathrm{query}}, \mathcal{P}_{\mathrm{service}})`$ zu.
4. **Überwindung der Speicherkapazitätskluft:** Die strikte Trennung statischer Textzitate vom mathematischen Regelskelett reduziert den Arbeitsbereich auf das kompakte Active Working Set ($`\mathrm{AWS} \le 32\,\text{MB}`$). Eine zweibankige Pipelined-Pufferung (Ping-Pong Pipeline) garantiert eine vollständige Latenzverbergung durch Erfüllung der Ungleichung $`T_{\mathrm{refill}} \ll T_{\mathrm{execute}}`$ und verhindert Leerlaufzyklen des Inferenzkerns ($\eta = 1{,}00$).
5. **Mehrquellenbasiertes argumentatives Schiedsverfahren:** Bei unvollständigem Wissen aggregiert das Scatter-Gather-Muster Aussagen benachbarter Clusterdienste über den ASPIC+-Formalismus. Das System verstärkt die Verlässlichkeit durch Evidenz-Korroboration, entkräftet schwache Empfehlungen deterministisch durch höherrangige Imperative und aktiviert bei unauflösbaren Normenkollisionen eine kompromisslose *Fail-Closed*-Schutzschleuse.
6. **Praxisnahe Evolution nach Projektbudget:** Der vorgestellte Referenzprototyp fungiert als modularer Entwicklungsleitfaden, der stufenlos vom Einplatz-Laborprüfstand (Stufe 1) über den R&D-Team-Cluster (Stufe 2) bis hin zum industriellen, hardware-deterministischen ASIL-D-System (Stufe 3) skaliert.
7. **Triage von Wissensdefiziten und SimHash-Clustering:** Die dezentrale Verarbeitung identifizierter Wissenslücken mittels 64-Bit-SimHash und BK-Tree unterbindet lawinenartige Anfrage-Stürme und eine Überlastung der Dienste. Ein lokaler WAL-Puffer mit strikter Festplattenquote garantiert die Erhaltung ermittelter Wissensdefizite selbst bei Verbindungsausfällen zum zentralen Broker, während der Übergang in den Lame-Duck-Modus die Ausgabe unbewiesener Aussagen verhindert.
8. **Schutz des Merkmalsraums und unveränderliches Entscheidungsbuch:** Token-Entartungen in Abfragen und schematische Anomalien werden an der Ingress-Grenze durch strikte Deserialisierung (`DisallowUnknownFields`) und die Invariante $`\mathrm{DistinctTokens}(T) \ge 2`$ abgewehrt. Alle erzeugten Entscheidungs-Beweisbündel werden irreversibel in Merkle-Bäumen und WORM-Medien verankert, wodurch ein manipulationssicherer Audit-Trail mit kryptografischer Bindung von Eingangsbedingungen, Urteil und Regelversion entsteht.

---

## Fragen zur Selbstüberprüfung

1. Warum verletzt der Versuch, ein Expertensystem als monolithischen „Fat Client“ auf dem Mobilgerät eines Ingenieurs auszuführen, fundamentale Anforderungen an Skalierbarkeit und Wissenssynchronisation?
2. Auf welche Weise garantiert ein schlanker Client die Evidenzinvariante $ZHR = 1.00$, obwohl er selbst keine logische Inferenz ausführt?
3. Worin unterscheidet sich semantisches Routing (Content-Based Semantic Routing) vom klassischen Routing auf Ebene gängiger Web-API-Gateways?
4. Wie lautet die mathematische Bedingung für eine vollständige Latenzverbergung in einer Ping-Pong-Doppelpufferungs-Pipeline, und warum muss die Hintergrund-Nachladezeit kleiner sein als die Regelausführungszeit?
5. Welche drei epistemischen Trigger signalisieren dem semantischen Broker die Notwendigkeit, von einer lokalen Auswertung zu einer föderierten Abfrage nach dem Scatter-Gather-Muster überzugehen?
6. Wie löst die ASPIC+-Argumentationsmaschine eine normative Kollision auf, wenn zwei gleichrangige, harmonisierte Standards einander widersprechende Anforderungen stellen?
7. Beschreiben Sie anhand der Technologieauswahlmatrix den Evolutionspfad eines Expertensystems vom Ein-Rechner-Laborprototypen bis zum heterogenen Industriekomplex.
8. Auf welche Weise lösen 64-Bit-SimHash-Fingerabdrücke und metrische BK-Bäume das Problem von Lawinenstürmen bei Wissensdefizit-Tickets in einem verteilten Cluster?
9. Warum kann das Ersetzen des Wildcard-Operators `vida.tasks.>` durch `vida.tasks.*` in einer NATS-Nachrichtenwarteschlange zu einem systemischen Verlust von Wissenskontext-Teilbäumen führen?
10. Welche vier kryptografischen Hashes bilden die Struktur des Decision Evidence Bundle, und wie garantiert der Merkle-Baum die Unveränderlichkeit und Nichtabstreitbarkeit gefällter Urteile?

---

## Glossar

| Deutscher Begriff | Englische Entsprechung | Kurzbeschreibung |
|---|---|---|
| Epistemische SOA | Epistemic Service-Oriented Architecture | Verteilte Expertensystem-Architektur mit Isolation von Wahrheitskontexten zwischen unabhängigen Wissensdiensten |
| Schlanker Interaktions-Client | Thin Interaction Client | Benutzeroberfläche ohne lokalen Wissensspeicher, ausschließlich für das Rendering kryptographisch signierter Nachweise |
| Semantischer Broker | Semantic Broker | Zentraler Cluster-Koordinator für syntaktisches Präprocessing und inhaltsbasiertes Abfrage-Routing |
| Inhaltsbasiertes Routing | Content-Based Semantic Routing | Weiterleitung von Abfragen an Wissensdienste auf Basis ontologischer Signaturen von Konzepten und Zielstandards |
| Aktiver Arbeitssatz | Active Working Set (AWS) | Kompakte transitive Hülle relevanter Regeln und Defeater ohne Textkorpus der Primärquellen |
| Ping-Pong-Doppelpufferung | Ping-Pong Double Buffering | Symmetrische Speicherarchitektur mit alternierendem Hintergrundladen einer Bank und paralleler Regelausführung aus der anderen |
| Latenzverbergung | Latency Hiding | Zustand einer Inferenz-Pipeline ($`T_{\mathrm{refill}} \le T_{\mathrm{execute}}`$), in dem Speicherladezeiten vollständig durch Rechenzeiten maskiert werden |
| Scatter-Gather-Muster | Scatter-Gather Pattern | Parallele Abfrage verteilter Clusterdienste mit anschließender Zusammenführung der Ergebnisse zu einer einheitlichen Schlussfolgerung |
| Evidenz-Korroboration | Corroboration | Mathematische Verstärkung des aggregierten Vertrauens in ein Urteil bei Vorliegen übereinstimmender Nachweise unabhängiger Quellen |
| Regelentkräftung | Undercutting Defeat | Mechanismus der nichtmonotonen Argumentation, bei dem eine höherrangige verbindliche Norm eine kollidierende Empfehlung annulliert |
| Fail-Closed-Kollisionsschleuse | Fail-Closed Qualified Conflict Gate | Sicherheitsmechanismus zur Verhinderung automatischer Scheinkompromisse bei unvereinbaren gleichrangigen Normen mit Eskalation an den Architekten |
| Triage von Wissensdefiziten | Knowledge Gap Triage | Algorithmische Clusterung und Deduplizierung ermittelter fehlender Fakten in verteilten Umgebungen |
| SimHash-Fingerabdruck | SimHash Fingerprint | 64-Bit-Hash mit lokaler Sensitivität zur Messung der semantischen Hamming-Distanz |
| Burkhard-Keller-Baum | Burkhard-Keller Tree (BK-Tree) | Metrischer Baum zur diskreten logarithmischen Suche $`O(\log N)`$ nach nächstgelegenen Wissensvektoren |
| Schutz vor Merkmalsentartung | Feature Degeneracy Guard | Ingress-Validierungsmechanismus zur Filterung trivialer und anomaler Vektoren anhand von Entropieschwellen |
| Entscheidungs-Beweisbündel | Decision Evidence Bundle | Kryptografisches Tupel aus Hashes der Eingangsdaten, der Wissensbasis-Version, des Beweisgraphen und einer Ed25519-Signatur |
| Unveränderliches Entscheidungsbuch | Immutable Decision Ledger | Verteilte Struktur auf Basis von Merkle-Bäumen und WORM-Speichern für das manipulationssichere Audit von Urteilen |

---

## Abkürzungen

| Abkürzung | Vollständige Bezeichnung | Bedeutung |
|---|---|---|
| ASIL | Automotive Safety Integrity Level | Sicherheitsanforderungsstufe für Automobilelektronik nach ISO 26262 (Stufen A bis D) |
| ASPIC+ | Argumentation Service Platform with Incomplete and Conflicting Information | Formaler Rahmen zur Strukturierung von Argumenten, Annahmen und Angriffs-/Entkräftungsregeln |
| AWS | Active Working Set | Aktiver Arbeitssatz von Wissensregeln, der minimal zur Berechnung einer Abfrage benötigt wird |
| BK-Tree | Burkhard-Keller Tree | Metrischer Baum zur Suche nächster Nachbarn in Räumen mit diskreter Metrik |
| BRAM | Block Random Access Memory | Schneller On-Chip-Speicher eines Hardware-Regelprozessors mit deterministischer Zugriffszeit |
| DLQ | Dead-Letter Queue | Quarantäne-Warteschlange für unzustellbare oder anomale Nachrichten zur Isolation fehlerhafter Pakete |
| DMA | Direct Memory Access | Direkter Speicherzugriff für den autonomen Hardware-Datentransfer ohne CPU-Beteiligung |
| EVM | Epistemic Virtual Machine | Deterministische virtuelle Maschine zur Ausführung logischer Resolutionen und Beweisführungen |
| gRPC | Google Remote Procedure Call | Leistungsfähiges binäres RPC-Framework auf Basis von HTTP/2 und Protocol Buffers |
| IPC | Inter-Process Communication | Betriebssystemmechanismen zum schnellen Datenaustausch zwischen Prozessen auf demselben Host |
| KTP | Knowledge Testing Pyramid | Mehrstufige Methodik zum Testen von Wissensbasen von atomaren Fakten bis zu Integrationsszenarien |
| LSH | Locality-Sensitive Hashing | Lokalitätssensitives Hashing zur schnellen Gruppierung ähnlicher Objekte |
| NPU | Neural Processing Unit | Spezialisierter Hardwarebeschleuniger für Tensor- und Matrixoperationen künstlicher neuronaler Netze |
| PWA | Progressive Web Application | Webanwendungs-Technologie mit Unterstützung für Offline-Caching der Benutzeroberfläche |
| RPC | Remote Procedure Call | Paradigma für die Interaktion verteilter Softwarekomponenten über entfernte Prozeduraufrufe |
| SLM | Small Language Model | Kompaktes spezialisiertes Sprachmodell (1B–7B Parameter) für linguistische Diskursanalyse |
| SOA | Service-Oriented Architecture | Serviceorientierte Softwarearchitektur mit lose gekoppelten Dienstkomponenten |
| TARA | Threat Analysis and Risk Assessment | Methodik zur Bedrohungsanalyse und Risikobewertung im Fahrzeugbau nach ISO/SAE 21434 |
| WAL | Write-Ahead Log | Vorauseilendes Festplattenprotokoll zur Gewährleistung von Atomarität und Persistenz |
| WORM | Write Once, Read Many | Speicher-Medium mit einmaligem Schreibvorgang und unveränderlichem Mehrfach-Lesezugriff |
| ZHR | Zero-Hallucination Rate | Metrik für das vollständige Fehlen unbegründeter Modellausgaben ($ZHR = 1.00$) in evidenzbasierten ES |

---

## Literaturverzeichnis

<a id="src-1"></a>[1] Dung, P. M. (1995). On the acceptability of arguments and its fundamental role in nonmonotonic reasoning, logic programming and n-person games. *Artificial Intelligence*, 77(2), 321–357.  
<a id="src-2"></a>[2] Prakken, H. (2010). An abstract framework for argumentation with structured arguments. *Argument & Computation*, 1(2), 93–124.  
<a id="src-3"></a>[3] Hennessy, J. L., & Patterson, D. A. (2019). *Computer Architecture: A Quantitative Approach* (6th ed.). Morgan Kaufmann. (Abschnitte zu Speicherhierarchien und Latenzverbergung).  
<a id="src-4"></a>[4] Fielding, R. T. (2000). *Architectural Styles and the Design of Network-based Software Architectures*. Doctoral dissertation, University of California, Irvine.  
<a id="src-5"></a>[5] ISO 26262:2018. *Road vehicles — Functional safety* (Parts 1–12). International Organization for Standardization.  
<a id="src-6"></a>[6] RTCA DO-178C / EUROCAE ED-12C (2011). *Software Considerations in Airborne Systems and Equipment Certification*. RTCA.  
<a id="src-7"></a>[7] NIST AI 100-1 (2023). *Artificial Intelligence Risk Management Framework (AI RMF 1.0)*. National Institute of Standards and Technology.  
<a id="src-8"></a>[8] Charikar, M. S. (2002). Similarity estimation techniques from rounding algorithms. *Proceedings of the thirtieth annual ACM symposium on Theory of computing*, 380–388.  
<a id="src-9"></a>[9] Burkhard, W. A., & Keller, R. M. (1973). Some approaches to best-match file searching. *Communications of the ACM*, 16(4), 230–236.  
<a id="src-10"></a>[10] Merkle, R. C. (1987). A digital signature based on a conventional encryption function. *Advances in Cryptology — CRYPTO ’87*, LNCS 293, 369–378.

---

[← Kapitel 33](ch33-inter-system-knowledge-exchange-and-model-teaching.md) | [Inhaltsverzeichnis](README.md) | [Teil VII](part-07-runtime-and-knowledge-exchange.md) | [Zu den Anhängen →](appendix-a-evidence-governed-framework.md)

