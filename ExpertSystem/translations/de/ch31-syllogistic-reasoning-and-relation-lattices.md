# Kapitel 31. Normenbasierte Inferenz: Prädikatenhierarchien, Ausnahmen und Geltung

> **Buch:** [Architektur evidenzbasierter Expertensysteme](README.md) · [Teil IV: Architektur, Technologie-Stack, Inferenz und Aktion](part-04-architecture-and-inference.md)  
> **Vorheriges Kapitel:** [Kapitel 19. Von der Frage zum Beweis: Suche, Bindung und Prüfung von Behauptungen](ch19-from-question-to-evidence.md)  
> **Nächstes Kapitel:** [Kapitel 20. Erklärungskomponente: Entscheidungen, Ablehnung und Kompetenzgrenzen](ch20-explanation-engine.md)  
> **Inhaltsverzeichnis:** [README.md](README.md)  
> **Autor:** [Mykola Fedchyk](about-the-author.md)  
> **Niveau:** Fortgeschritten: Entwickler von Inferenzmaschinen, Systemarchitekten, Knowledge Engineers, Spezialisten für formale Methoden  
> **Lernziele:** Quellenabruf (*Retrieval*) strikt von logischer Inferenz unterscheiden; Prädikatenhierarchien ohne unzulässige Rückwärtsspezialisierung konstruieren; bei unbekanntem Ausnahmestatus deterministisch verweigern (*Fail-Closed*); Rebutting Defeater von Undercutting Defeater differenzieren; die Geltung einer Norm im gegebenen Kontext mittels eines ASP-Programms mit Fail-Closed-Verhalten verifizieren; den Nutzen von RDFS, SPARQL, ASPIC+ und Knowledge Mining für Hierarchien und Ausnahmen fundiert bewerten; Inferenzpfade auditieren und die Grenzen des didaktischen Inferenzschritts in Go verstehen.

---

## Abstract

Man stelle sich ein industrielles SCADA-System oder das Bordnetz eines autonomen Transportsystems vor (IEC 62443 / ISO 26262), in dem ein automatisierter Netzwerkagent ein unerwartetes Reset-Paket (TCP RST) analysiert. Ein Vektor-Retrieval findet augenblicklich die einschlägige Klausel aus RFC 9293: „Bei Empfang eines RST muss die Verbindung unverzüglich getrennt werden.“ Akzeptiert der Agent diese allgemeine Regel als direkte Handlungsanweisung, ohne den Verbindungszustand (`SYN-SENT` oder `ESTABLISHED`), die Sequenznummer sowie das spezifische Sicherheitsprofil (beispielsweise die Schutzanforderungen gegen Session-Hijacking nach RFC 5961) zu verifizieren, trennt er eine aktive Notfalltelemetriesitzung oder die Bremsverriegelung (*Interlock*). Die Interpretation fehlender Informationen über Ausnahmen als Handlungsfreigabe („da in der Anfrage keine Sequenznummer erwähnt wurde, muss sie folglich nicht geprüft werden“) führt zur Notabschaltung der Fertigungslinie oder zur Deaktivierung der Fahrstabilisierung bei voller Fahrt.

Die Leitfrage dieses Kapitels lautet: **Wie lässt sich eine deterministische Schlussfolgerung aus mehreren miteinander verknüpften normativen Grundlagen ableiten, eine Prädikatenhierarchie entfalten und mathematisch garantieren, dass ein unbekannter Ausnahmestatus niemals in eine gefährliche Freigabe umschlägt?**

Dieses Kapitel trennt den Quellenabruf, Prädikatenverbände, die formale Ausnahmebehandlung (anfechtbares Schließen über Horn-Klauseln und ASPIC+-Logik) sowie die Auswahl der geltenden Norm strikt voneinander. Ein Suchprozess kann mehrstufig oder graphbasiert operieren, sein Ergebnis liefert jedoch stets nur Rohmaterial für den logischen Inferenzkern. Untersucht werden die Grounded-Semantik nach Dung ($`\mathrm{GE}`$), die topologische Zyklusauflösung mittels Tarjan-Algorithmus ($`O(V+E)`$), die automatische epistemische Argumentrehabilitation bei Eintreffen neuer Fakten sowie der Operator transitiver Kaskadenfalsifikation in hierarchischen Sicherheitszielbäumen gemäß ISO 26262 ASIL D. Ein didaktisches Go-Beispiel implementiert die Validierung des Inferenzpfads, während ein deklaratives Modell in Answer Set Programming (ASP / clingo) die Normenauswahl für neun Grenzszenarien der Netzwerkpaketverarbeitung formal überprüft. Der Determinismus logischer Inferenz garantiert jedoch keineswegs die Faktizität der Eingangsannahmen, weshalb die strikte Grenzziehung zwischen Beweisführung und Tatsachenbasis unantastbar bleibt.

---

## 1. Warum einstufiges Retrieval keine normative Inferenz liefert

Moderne semantische Suchpipelines stützen sich auf die Optimierung der Kosinusähnlichkeit hochdimensionaler Vektoreinbettungen (*Embeddings*):

```math
\text{Query} \xrightarrow{\text{Embed}} \mathbf{v}_q \implies \arg\max_k \cos(\mathbf{v}_q, \mathbf{v}_k).
```

Bezeichnungen des Vektor-Retrievals:

- $\text{Query}$ ist die textuelle Anfrage des Ingenieurs oder Operators;
- $\text{Embed}$ bezeichnet die Abbildungsfunktion von Text in einen reellwertigen Vektorraum;
- $\mathbf{v}_q$ ist der Anfragevektor im latenten Einbettungsraum;
- $\mathbf{v}_k$ bezeichnet die indizierten Vektoren gespeicherter Dokumentfragmente (*Chunks*);
- $\cos(\mathbf{v}_q, \mathbf{v}_k)$ ist die Kosinusähnlichkeit der Vektoren im Intervall $[-1, 1]$;
- $\arg\max_k$ selektiert den Index des Chunks mit der höchsten geometrischen Nähe.

Diese Formulierung beschreibt ein Ranking nach inhaltlicher Vektorähnlichkeit, keineswegs jedoch die Verifikation der logischen Geltung einer Norm. Geometrische Nähe beschränkt sich zwar nicht auf Wortübereinstimmungen, liefert aber keinen formalen Beweis der normativen Anwendbarkeit.

Das folgende Diagramm veranschaulicht das typische Versagen einer Pipeline, die Anwendungsbedingungen der aufgefundenen Norm nicht formal verifiziert:

```mermaid
flowchart TD
    accTitle: Falle des einstufigen Retrievals gegenüber syllogistischer Deduktion
    accDescr: Warum die Vektorsuche bei mehrstufigen normativen Schlussfolgerungen scheitert.

    subgraph RAGFail["Falle des einstufigen Vektor-Retrievals"]
        Q["Anfrage: Verbindung nach Reset-Nachricht trennen?<br/>Zustand und Sequenznummer nicht angegeben"]
        EMB["Vektorsuche im RFC-Korpus"]
        CH1["Veraltete Beschreibung der Reset-Verarbeitung gefunden<br/>ohne Prüfung von Zustand und anwendbarem Profil"]
        CH2["RFC 5961 (2010) übersehen:<br/>Schutz vor Blind Resets via Challenge-ACK"]
        CH3["RFC 9293 (2022) übersehen:<br/>RFC 793 obsoletiert"]
        CH1 --> BAD_ANS["<b>Unbestätigte Schlussfolgerung</b>:<br/>Verbindung ohne Bedingungsprüfung trennen"]
    end

    subgraph SyllogisticSuccess["Deterministische syllogistische Pipeline"]
        S_TERM["1. Begriffe: TCP-Zustandsautomat, RST-Verarbeitung"]
        S_LATTICE["2. Prädikatenhierarchie und Revisionsgraph"]
        S_CONFLICT["3. Konfliktdetektor: RFC 9293 obsoletiert RFC 793"]
        S_DEFEAT["4. Defeater: RFC-5961-Schutz ersetzt bedingungslosen Reset"]
        S_SYLLOGISM["5. Zustand, Sequenznummer<br/>und unterstütztes Profil präzisieren"]
        S_SYLLOGISM --> GOOD_ANS["<b>Konditionales Urteil</b>:<br/>Reset, Challenge-ACK<br/>oder Verwerfen je nach Eingaben"]
    end

    classDef fail fill:#ffebee,stroke:#d32f2f,stroke-width:2px,color:#b71c1c;
    classDef succ fill:#e8f5e9,stroke:#2e7d32,stroke-width:2px,color:#1b5e20;

    class Q,EMB,CH1,CH2,CH3,BAD_ANS fail;
    class S_TERM,S_LATTICE,S_CONFLICT,S_DEFEAT,S_SYLLOGISM,GOOD_ANS succ;
```

### 1.1. Wo die Ähnlichkeitssuche die Logik der Normen verfehlt

1. **Prämissen sind über Dokumente fragmentiert.** Eine Grundregel ist in einem Standard fixiert, die Ausnahme in einem zweiten, und die Außerkraftsetzung (*Obsoletierung*) der Altversion in einem dritten Dokument, das Jahrzehnte später publiziert wurde. Die Ähnlichkeitssuche bewertet jeden Chunk isoliert. Sie liefert womöglich alle drei Fragmente zurück, stellt jedoch nicht fest, dass das dritte Dokument das erste aufhebt (*Lex Posterior*) und das zweite dessen Gültigkeitsbereich einschränkt (*Lex Specialis*). Diese Abhängigkeiten konstituieren explizite Dokumenten- und Normenrelationen, die strukturiert im Wissensgraphen abgebildet werden müssen.
2. **Verbindlichkeitsgrade spiegeln sich nicht im Vektor wider.** Schlüsselwörter wie `MUST`, `SHOULD`, `RECOMMENDED` und `MAY` treten in strukturell identischen Kontexten auf, weshalb ihre Einbettungsvektoren nahezu ununterscheidbar nah beieinander liegen. Vektornähe garantiert keineswegs die Bewahrung der normativen Verbindlichkeitsstufen, wie sie in RFC 2119 und RFC 8174 definiert sind [[1]](#src-1) [[2]](#src-2). RFC 8174 führt eine weitere formale Feinheit ein: Ausschlaggebende normative Bindekraft besitzen ausschließlich großgeschriebene Schlüsselwörter; statistische Sprachmodelle ignorieren Groß- und Kleinschreibung jedoch häufig. Daher muss die deontische Modalität als eigenständiges Feld des Urteilstupels persistiert und darf nicht aus Textähnlichkeiten extrapoliert werden.
3. **Abwesenheit eines Eintrags wird fälschlich als Negation interpretiert.** Fehlt in der Wissensbasis ein explizites Handlungsverbot, folgert eine Inferenzmaschine mit Negation als Fehlschlag (*Negation as Failure*, NAF, wie im Standard-Prolog), dass kein Verbot vorliegt. Dies ist ausschließlich in explizit geschlossenen Wissensdomänen zulässig. In allen übrigen Fällen bedeutet das Fehlen von Daten den epistemischen Status „unbekannt“ (*Unknown*), woraufhin das System nach dem Prinzip des geschlossenen Ausfalls (*Fail-Closed*) blockieren muss.

---

## 2. Begriff, Urteil und Schluss

Für das didaktische Modell trennen wir Begriff, Urteil und Inferenzschritt strikt. Die klassischen kategorischen Syllogismen des Aristoteles bilden das historische Fundament formaler Logik [[3]](#src-3); das nachfolgend eingeführte normative Tupel stellt jedoch eine moderne ingenieurtechnische Abstraktion dieses Buches dar und ist keine buchstäbliche Rekonstruktion aristotelischer oder peircescher Notationen.

```mermaid
flowchart LR
    accTitle: Triade des Denkens: Begriff, Urteil, Schluss
    accDescr: Struktureller Übergang von Entitäten und Prädikaten zu bewiesenen Syllogismen.

    subgraph Level1["1. Begriffsebene"]
        direction TB
        TERM["<b>Terminus (Begriff)</b><br/>• Intension (Merkmale)<br/>• Extension (Umfang)"]
        LATTICE["<b>Prädikatenhierarchie</b><br/>Subsumtion: Sub ⊑ Super"]
    end

    subgraph Level2["2. Urteilsebene"]
        direction TB
        PROP["<b>Urteil (Proposition)</b><br/>Subjekt, Relation, Objekt"]
        MODAL["<b>Deontische Modalität</b><br/>MUST, SHOULD, MAY, MUST_NOT"]
        DEFEAT["<b>Defeater</b><br/>Ausnahmebedingungen (UNLESS)"]
    end

    subgraph Level3["3. Schlussebene"]
        direction TB
        PREM["<b>Prämissen des Syllogismus</b><br/>Obersatz und Untersatz"]
        RULES["<b>Inferenzregeln</b><br/>Modus Ponens, Subsumtion"]
        DAG["<b>Inferenzpfad (Proof Trace)</b><br/>Regel · Begründung · Ergebnis"]
    end

    Level1 ==> Level2 ==> Level3

    classDef c1 fill:#e1f5fe,stroke:#0288d1,stroke-width:2px,color:#01579b;
    classDef c2 fill:#fff3e0,stroke:#e65100,stroke-width:2px,color:#bf360c;
    classDef c3 fill:#e8f5e9,stroke:#2e7d32,stroke-width:2px,color:#1b5e20;

    class TERM,LATTICE c1;
    class PROP,MODAL,DEFEAT c2;
    class PREM,RULES,DAG c3;
```

### 2.1. Der Begriff

Ein Begriff formalisiert eine abstrakte Entität oder ein physikalisches Objekt der Ingenieurdomäne. Er wird mathematisch durch ein Mengenpaar beschrieben:
- **Intension (Begriffsinhalt):** Die Menge aller wesentlichen Merkmale, Eigenschaften und Invarianten, die den Begriff eindeutig von anderen Entitäten abgrenzen:

```math
\text{Intension}(C) = \{ P_1, P_2, \dots, P_k \}.
```

Bestandteile der Intension:

- $C$ ist der Begriff der Fachdomäne;
- $\text{Intension}(C)$ ist die Menge der konstitutiven Merkmale und Invarianten des Begriffs;
- $P_1, \dots, P_k$ sind Eigenschaftsprädikate, die das Wesen der Entität determinieren;
- $k$ ist die Anzahl der obligatorischen Merkmale der Definition.

- **Extension (Begriffsumfang):** Die Menge aller konkreten Instanzen oder Spezialisierungen, die sämtliche Merkmale der Intension erfüllen:

```math
\text{Extension}(C) = \{ x \mid \forall P \in \text{Intension}(C) : P(x) = \text{True} \}.
```

Symbole der Extension:

- $\text{Extension}(C)$ bezeichnet den Begriffsumfang, also die Menge aller konkreten Objekte;
- $x$ steht für eine konkrete Ingenieurentität oder Instanz;
- $\forall P$ fordert die Gültigkeit jedes Merkmals $P$ der Intension für das gegebene Objekt;
- $\text{True}$ belegt die Erfüllung der Bedingung durch das Objekt.

Zwischen Begriffen gilt das fundamentale **Gesetz des umgekehrten Verhältnisses von Begriffsinhalt und Begriffsumfang**: Je reichhaltiger die Intension (je mehr einschränkende Merkmale einer Definition hinzugefügt werden), desto enger wird die Extension (desto weniger Entitäten erfüllen diese Definition).

### 2.2. Das Urteil

Ein Urteil affirmiert oder negiert das Bestehen einer Relation zwischen Begriffen. In einem beweisbasierten Expertensystem wird ein Urteil als erweitertes semantisches Tupel formalisiert:

```math
\mathcal{J} = \langle \text{Subject}, \; \mathcal{R}, \; \text{Object}, \; \mathcal{M}, \; \mathcal{D}, \; \text{Provenance} \rangle.
```

Elemente des normativen Urteils:

- $\mathcal{J}$ ist das formale Tupel des normativen Urteils;
- $\text{Subject}$ und $\text{Object}$ sind Konzepte im ingenieurtechnischen Domänenraum;
- $\mathcal{R}$ ist ein Prädikat aus der Relationen- und Prädikatenhierarchie;
- $`\mathcal{M} \in \{ \text{MUST}, \text{MUST-NOT}, \text{SHOULD}, \text{SHOULD-NOT}, \text{MAY} \}`$ bezeichnet die deontische normative Modalität;
- $`\mathcal{D} = \{ d_1, d_2, \dots \}`$ ist die Menge anfechtender Ausnahmebedingungen (*Defeater*);
- $\text{Provenance}$ ist der Herkunftsnachweis der Primärquelle (Dokumenten-ID, kryptografischer SHA-256-Hash und Bytestellen-Offsets).

### 2.3. Der Schluss

Ein Syllogismus ist ein deterministischer Inferenzschritt, bei dem aus zwei Prämissen (dem Obersatz und dem Untersatz), die über einen gemeinsamen Mittelbegriff ($M$) verfügen, mit logischer Notwendigkeit ein drittes Urteil (die Konklusion bzw. der Schluss) folgt:

```math
\frac{\text{Major Premise: } \forall x : M(x) \xrightarrow{\mathcal{M}} P(x), \quad \text{Minor Premise: } M(S)}{\text{Conclusion: } S \xrightarrow{\mathcal{M}} P}.
```

Konstitutive Elemente des Syllogismus:

- $\text{Major Premise}$ ist der Obersatz der allgemeinen Regel für alle Entitäten der Klasse $M$;
- $\text{Minor Premise}$ ist der Untersatz, der die Zugehörigkeit des Subjekts $S$ zur Klasse $M$ feststellt;
- $M$ ist der Mittelbegriff (*terminus medius*), der beide Prämissen relational bindet;
- $\text{Conclusion}$ ist die logisch zwingende Konklusion über die Gültigkeit der Eigenschaft $P$ für Subjekt $S$;
- $\mathcal{M}$ bezeichnet die strikt bewahrte deontische Modalität der Norm.

**Anwendungsbeispiel im automatisierten Protokollaudit:** Das Simple Mail Transfer Protocol (SMTP) trennt strikt zwischen Client und Server. Gemäß Abschnitt 4.1.1.1 von RFC 5321 muss ein SMTP-Client vor Beginn einer Mail-Transaktion den Befehl `HELO` oder `EHLO` absetzen [[4]](#src-4).

1. Obersatz: Für einen SMTP-Client ist vor Beginn einer Mail-Transaktion der Befehl `HELO` oder `EHLO` obligatorisch (`MUST`).
2. Untersatz: In einem verifizierten oder explizit hypothetischen Kontext ist `MailClient-01` ein solcher Client und leitet eine Transaktion ein.
3. Konditionale Konklusion: Für `MailClient-01` gilt exakt diese Anforderung, keineswegs jedoch die isolierte Forderung nach ausschließlicher Nutzung von `EHLO`.

Wurde die Client-Rolle lediglich im Rahmen einer kontrafaktischen „Was-wäre-wenn“-Anfrage postuliert, verbleibt auch der Schluss rein hypothetisch. Für ein operatives Sicherheitsaudit ist ein unabhängig verifizierter Zustand zwingend erforderlich. Ein erfolgreicher TCP-Verbindungsaufbau garantiert keineswegs das Zustandekommen einer SMTP-Sitzung: Abschnitt 3.1 von RFC 5321 gestattet dem Server eine initiale Verweigerung mittels Statuscode 554.

---

## 3. Prädikatenhierarchie und Semantik der Subsumtion

Stellt ein Ingenieur eine generalisierte Anfrage: *„Welche Sicherheitsanforderungen gelten für Protokoll X?“*, darf sich das System keinesfalls auf die Suche nach Prädikaten mit dem wörtlichen Bezeichner `security_requirement` beschränken. In realen Wissensbasen stammen Fakten aus disparaten Abschnitten und tragen Bezeichner wie `must_encrypt_channel`, `authenticate_peer_certificate` oder `validate_sequence_number`.

Um generalisierte Abfragen deterministisch aufzulösen, werden Prädikate in einer azyklischen Hierarchie strukturiert. Die Erreichbarkeit im gerichteten Graphen definiert eine Halbordnung der Spezialisierung:

```math
\mathcal{L} = \langle \mathcal{R}, \sqsubseteq \rangle.
```

Notationen der Hierarchie:

- $\mathcal{L}$ ist die partiell geordnete Menge (*Poset*) von Prädikaten;
- $\mathcal{R}$ ist die Menge aller ingenieurtechnischen Relationstypen;
- $\sqsubseteq$ bezeichnet die Halbordnungsrelation der Subsumtion (Spezialisierung).

Die Relation $R_1 \sqsubseteq R_2$ besagt, dass das Prädikat $R_1$ das Prädikat $R_2$ spezialisiert. So spezialisiert `must_requirement` das allgemeinere `normative_requirement`. Eine Hierarchie darf erst dann als Verband (*Lattice*) bezeichnet werden, wenn für jedes Paar von Elementen ein eindeutiges Infimum (größte untere Schranke, Meet $\sqcap$) und Supremum (kleinste obere Schranke, Join $\sqcup$) existiert. Der nachfolgende Graph und Beispielcode verzichten auf diese Verbandsoperationen, da für die Pfadprüfung eine reine Baum- bzw. DAG-Struktur genügt.

```mermaid
flowchart TD
    accTitle: Prädikatenhierarchie von der Spezialisierung zur Generalisierung
    accDescr: Baumartiger DAG von Relationen ausgehend von generischen Kategorien bis zu spezifischen Normen.

    TOP["Generisches Prädikat"]
    
    DEF["definition"]
    NORM["normative_requirement"]
    PROP["protocol_property"]
    DIAG["diagnostic_rule"]

    DEF & NORM & PROP & DIAG --> TOP

    DEF_PURP["defines_purpose"] --> DEF
    ABBR["abbreviation_expansion"] --> DEF

    MUST_REQ["must_requirement"] --> NORM
    PROHIB["prohibited_requirement"] --> NORM
    RECOM["recommended_practice"] --> NORM

    PORT["default_port"] --> TRANSPORT["network_transport"] --> PROP
    SEQ_ERR["sequence_error_cause"] --> ERR_CAUSE["error_cause"] --> DIAG

    classDef top fill:#ede7f6,stroke:#512da8,stroke-width:2px,color:#311b92;
    classDef mid fill:#e1f5fe,stroke:#0288d1,stroke-width:2px,color:#01579b;
    classDef leaf fill:#e8f5e9,stroke:#2e7d32,stroke-width:2px,color:#1b5e20;

    class TOP top;
    class DEF,NORM,PROP,DIAG,TRANSPORT,ERR_CAUSE mid;
    class DEF_PURP,ABBR,MUST_REQ,PROHIB,RECOM,PORT,SEQ_ERR leaf;
```

### 3.1. Mathematische Gesetze gerichteter Subsumtion

Die Halbordnung erfüllt folgende Axiome:
1. **Reflexivität:** $\forall R \in \mathcal{R} : R \sqsubseteq R$.
2. **Antisymmetrie:** $`\forall R_1, R_2 \in \mathcal{R} : (R_1 \sqsubseteq R_2 \land R_2 \sqsubseteq R_1) \implies R_1 = R_2`$.
3. **Transitivität:** $`\forall R_1, R_2, R_3 \in \mathcal{R} : (R_1 \sqsubseteq R_2 \land R_2 \sqsubseteq R_3) \implies R_1 \sqsubseteq R_3`$.

**Regel der Anfragegeneralisierung.**  
Eine in der Wissensbasis persistierte Relation $`R_{\text{fact}}`$ erfüllt eine Benutzeranfrage mit Prädikat $`R_{\text{query}}`$ genau dann, wenn sie sich im Teilbaum der Nachkommen des angefragten Prädikats befindet:

```math
\text{Matches}(R_{\text{query}}, R_{\text{fact}}) \iff R_{\text{fact}} \sqsubseteq^* R_{\text{query}}.
```

In der Regel gerichteter Subsumtion:

- $\text{Matches}$ ist die boolesche Validierungsfunktion der Faktenanwendbarkeit auf die Anfrage;
- $`R_{\text{query}}`$ ist das Prädikat aus der Benutzeranfrage;
- $`R_{\text{fact}}`$ ist das in der Wissensbasis explizit hinterlegte Faktenprädikat;
- $\sqsubseteq^*$ bezeichnet die reflexiv-transitive Hülle der Subsumtion (Existenz eines gerichteten Pfades zur Verallgemeinerung).

**Verbot der Rückwärtsspezialisierung.**  
Fordert die Anfrage eines Sicherheitsprüfers ein konkretes striktes Verbot $`R_{\text{specific}}`$ (etwa `prohibited_requirement`), ist es kategorisch unzulässig, ein allgemeines übergeordnetes Faktum $`R_{\text{general}}`$ (`normative_requirement`) zurückzugeben, da eine allgemeine normative Vorgabe keineswegs die Erfüllung eines spezifischen Verbots garantiert:

```math
R_{\text{general}} \not\sqsubseteq^* R_{\text{specific}} \quad \text{wenn } R_{\text{specific}} \sqsubset R_{\text{general}}.
```

Bezeichnungen des Spezialisierungsverbots:

- $`R_{\text{general}}`$ ist das allgemeinere Prädikat auf höherer Hierarchiestufe;
- $`R_{\text{specific}}`$ ist das strikte Prädikat einer untergeordneten Ebene;
- $\not\sqsubseteq^*$ verbietet den Match: Die Existenz einer allgemeinen Anforderung beweist nicht die Erfüllung einer spezifischen Ausnahme- oder Verbotsnorm.

### 3.2. Abgleich natürlicher Phrasen mit Prädikaten

In einer industriellen Pipeline formulieren Ingenieure ihre Anfragen selten in den formalen Termini der internen Ontologie. Ein Entwickler fragt: *„Welcher Standardport gehört zu BGP?“*, *„Wurde RFC 821 für ungültig erklärt?“* oder *„Durch welchen Standard wird dieses Dokument ersetzt?“*.

Ein kontrolliertes Vokabular oder ein syntaktischer Sprachanalysator gleicht die natürlichsprachliche Phrase zunächst mit einem registrierten kanonischen Prädikat ab. Die Hierarchie selbst interpretiert keine natürliche Sprache und erzeugt keine heuristischen Synonyme.

1. Ein eindeutiges Mapping verknüpft die Phrase mit einem konkreten Prädikat und einer kontrollierten Wörterbuchversion.
2. Für eine generalisierte Abfrage wird der Pfad von diesem Prädikat zu den Vorfahren validiert; eine Rückwärtsspezialisierung ist verboten.
3. Unbekannte oder mehrdeutige Phrasen erzwingen eine Klärungsinteraktion (*Clarification Dialog*). Eine Phrase erhält nicht allein deshalb ein bestimmtes Ahnenprädikat, weil im Text das Wort „veraltet“ auftaucht.

Latenz und Präzision dieses Abgleichs müssen an konkreten Testkorpora evaluiert werden. Auch ein deterministisches Wörterbuch kann fehlerhafte Zuordnungen aufweisen; Reproduzierbarkeit darf nicht mit inhaltlicher Korrektheit verwechselt werden.

### 3.3. Disambiguierung gleichnamiger Dokumente entlang der Geltungslinie

Eine typische Fehlerquelle technischer Korpora ist die **Polysemie von Entitätsnamen**: Ein und dieselbe Dokument- oder Protokollbezeichnung tritt in Dutzenden Spezifikationen unterschiedlicher Dekaden auf. So referenziert der Name *„Simple Mail Transfer Protocol“* gleichermaßen RFC 821 (1982), RFC 2821 (2001) und RFC 5321 (2008).

Wird ein Dokument nach seinem Titel ohne explizite Normennummer angefragt, liefert ein naives Retrieval drei konfligierende Dokumentfragmente, was zum stillschweigenden Rückgriff auf eine veraltete Spezifikation führen kann.

Um dieses Versagen auszuschließen, nutzt der Inferenzkern eine **Disambiguierung über den Geltungsgraphen**:
1. Der Resolver evaluiert Revision, Domäne, Veröffentlichungsdatum und das Ziel der Anfrage: operative Sicherheitsprüfung des Ist-Zustands oder historische Rekonstruktion.
2. Für eine aktuelle Konformitätsprüfung kann die einzig gültige Revision gewählt werden, keineswegs jedoch für Fragen zu Altsystemen allein aufgrund eines neueren Datums.
3. Liegen mehrere geltende Revisionen vor oder ist der Status unbestimmt, initiiert das System einen Klärungsdialog. Die Abstammungslinie eines Dokuments determiniert nicht automatisch, welche Spezifikationsprofile ein konkretes industrielles Erzeugnis tatsächlich implementiert.

---

## 4. Kleenes dreiwertige Logik und Pollocks Defeater

Zwei boolesche Wahrheitswerte implizieren für sich genommen keineswegs die Annahme einer geschlossenen Welt (*Closed-World Assumption*, CWA). Ob fehlende Fakten als falsch gewertet werden, ist eine bewusste epistemische Architekturentscheidung. Im hier dargestellten Modell besitzt ein fehlendes Prüfergebnis grundsätzlich den Status „unbekannt“ (*Unknown*). Nur für explizit als vollständig deklarierte Datenbereiche darf von diesem Fail-Closed-Verhalten abgewichen werden ([Kapitel 7](ch07-knowledge-base-typology.md)).

Zur didaktischen Formalisierung von Inferenz unter unvollständiger Information nutzen wir die **starke dreiwertige Logik nach Kleene** (*Strong Kleene 3-Valued Logic*, 3VL) [[5]](#src-5):

```math
\mathcal{V}_3 = \{ \text{True}, \; \text{False}, \; \text{Unknown} \}.
```

Wahrheitswerte der Kleene-Logik:

- $`\mathcal{V}_3`$ ist der dreiwertige Wahrheitsraum nach Strong Kleene;
- $\text{True}$ repräsentiert die formal bewiesene Wahrheit einer Aussage;
- $\text{False}$ repräsentiert die formal bewiesene Falschheit;
- $\text{Unknown}$ kennzeichnet fehlende Daten oder Unbestimmtheit, ohne dass Falschheit unterstellt werden darf.

### 4.1. Wahrheitstafel der starken Kleene-Logik

| $A$ | $B$ | $A \land B$ | $A \lor B$ | $\neg A$ | $A \to B$ |
| :---: | :---: | :---: | :---: | :---: | :---: |
| $\text{True}$ | $\text{True}$ | $\text{True}$ | $\text{True}$ | $\text{False}$ | $\text{True}$ |
| $\text{True}$ | $\text{False}$ | $\text{False}$ | $\text{True}$ | $\text{False}$ | $\text{False}$ |
| $\text{True}$ | $\text{Unknown}$ | $\mathbf{Unknown}$ | $\text{True}$ | $\text{False}$ | $\mathbf{Unknown}$ |
| $\text{False}$ | $\text{True}$ | $\text{False}$ | $\text{True}$ | $\text{True}$ | $\text{True}$ |
| $\text{False}$ | $\text{False}$ | $\text{False}$ | $\text{False}$ | $\text{True}$ | $\text{True}$ |
| $\text{False}$ | $\text{Unknown}$ | $\text{False}$ | $\mathbf{Unknown}$ | $\text{True}$ | $\text{True}$ |
| $\text{Unknown}$ | $\text{True}$ | $\mathbf{Unknown}$ | $\text{True}$ | $\mathbf{Unknown}$ | $\text{True}$ |
| $\text{Unknown}$ | $\text{False}$ | $\text{False}$ | $\mathbf{Unknown}$ | $\mathbf{Unknown}$ | $\mathbf{Unknown}$ |
| $\text{Unknown}$ | $\text{Unknown}$ | $\mathbf{Unknown}$ | $\mathbf{Unknown}$ | $\mathbf{Unknown}$ | $\mathbf{Unknown}$ |

### 4.2. Informationsordnung und Fail-Closed-Verhalten

In der Kleene-Logik wird eine Halbordnung nach Informationsgehalt ($`\le_i`$) definiert:

```math
\text{Unknown} \le_i \text{True}, \quad \text{Unknown} \le_i \text{False}.
```

Informationsordnung:

- $`\le_i`$ ist die Halbordnung wachsenden Wissens (*Information Ordering*);
- Die Ungleichungen besagen, dass der Zustand $\text{Unknown}$ minimale Information trägt, während der Übergang zu $\text{True}$ oder $\text{False}$ den Informationsgehalt monoton steigert.

Eine Funktion $f$ ist informationsmonoton, wenn ein Erkenntnisgewinn (der Übergang von $\text{Unknown}$ zu $\text{True}$ oder $\text{False}$) ein bereits ermitteltes deterministisches Ergebnis niemals ins Gegenteil verkehrt.

> [!IMPORTANT]
> **Das Prinzip des geschlossenen Ausfalls (Fail-Closed).**  
> Evaluiert eine zwingend vorgeschriebene Bedingung zu $\text{Unknown}$, verweigert die Sicherheitsarchitektur die Ausführung und benennt die fehlenden Primärfakten. Dies ist ein deterministisches Schutzverhalten des Verifikators, keine Zertifizierungsfreigabe.

### 4.3. Defeater nach Pollock

In technischen Standards besitzen Normen überwiegend anfechtbaren (*defeasible*) Charakter: Sie beanspruchen Gültigkeit, solange keine Ausnahmebedingung greift. Der Philosoph John Pollock differenzierte zwei fundamentale Klassen von Einwänden (*Defeaters*) [[6]](#src-6):

1. **Rebutting Defeater (Widerlegender Einwand)**  
   Greift unmittelbar die Konklusion an, indem eine diametral entgegengesetzte Aussage abgeleitet wird:

```math
A \implies P, \quad B \implies \neg P.
```

Parameter des direkten Widerlegers:

- $A$ und $B$ sind die Prämissen zweier konkurrierender Regeln;
- $P$ ist die behauptete Konklusion und $\neg P$ deren logische Negation;
- Ein Rebutting Defeater erzeugt zunächst eine Antinomie zweier Argumente; welches davon obsiegt, entscheidet eine explizite Vorrangrelation (etwa „speziellere Norm bricht allgemeinere Norm“).

2. **Undercutting Defeater (Untergrabender Einwand)**  
   Greift den Geltungszusammenhang zwischen Prämisse und Konklusion an. Er besagt, dass die Regel unter den gegebenen Bedingungen suspendiert ist, ohne notwendigerweise die Falschheit der Konklusion zu behaupten:

```math
U \implies \neg (A \hookrightarrow P).
```

Komponenten des untergrabenden Einwands:

- $U$ ist eine kontextuelle Ausnahmebedingung (etwa der Zustand absoluter Funkstille);
- $A \hookrightarrow P$ bezeichnet die normative Verknüpfung zwischen Tatbestand und Handlungspflicht;
- $\neg (A \hookrightarrow P)$ annulliert die Handlungspflicht, ohne das Gegenteil zu affirmieren.

```mermaid
flowchart LR
    accTitle: Schema der Aktivierung eines Undercutting Defeaters
    accDescr: Eine Ausnahmebedingung neutralisiert den Verknüpfungszusammenhang zwischen Regel und Handlungspflicht.

    RULE["Grundregel: Alle 10 s Kontrollsignal senden (MUST)"]
    DEF["Undercutting Defeater: Sofern Funkstille nicht aktiviert ist"]
    ENV["Sensorstatus: Modus Funkstille ist aktiv"]

    ENV -->|"Aktiviert Bedingung"| DEF
    DEF -->|"Untergräbt Verknüpfung"| RULE
    RULE --> RES["<b>Ergebnis</b>: Signalforderung entfällt, Alarm wegen fehlenden Signals wird unterdrückt"]

    classDef r fill:#e1f5fe,stroke:#0288d1,stroke-width:2px,color:#01579b;
    classDef d fill:#ffebee,stroke:#d32f2f,stroke-width:2px,color:#b71c1c;
    classDef res fill:#e8f5e9,stroke:#2e7d32,stroke-width:2px,color:#1b5e20;

    class RULE r;
    class DEF,ENV d;
    class RES res;
```

Diese theoretische Unterscheidung hat gravierende ingenieurtechnische Konsequenzen: Ein Undercutting Defeater hinterlässt das Expertensystem ohne Handlungsanweisung. Die Regel ist unanwendbar; die Systemantwort muss die suspendierte Regel nebst Ausnahmebedingung benennen, ohne eine Ersatzaktion zu fingieren. Ein Rebutting Defeater hingegen liefert eine alternative Konklusion, sodass lediglich ein Vorrangkonflikt zwischen zwei wohlgeformten Handlungsoptionen zu lösen ist. Sanjay Modgil und Henry Prakken formalisierten im ASPIC+-Framework für strukturierte Argumentation beide Angriffsarten zusammen mit einer dritten: dem Untergraben von Prämissen (*Undermining*) [[7]](#src-7). In ASPIC+ ist ein Undercutting-Angriff auf eine Inferenzregel unabhängig von Prioritätsordnungen stets erfolgreich, während Angriffe auf Konklusionen (Rebutting) oder Prämissen (Undermining) von der Präferenzrelation zwischen Argumenten abhängen. Das Go-Programm in Abschnitt 8 implementiert exakt diese Differenzierung: Ein Undercutting Defeater führt zur fundierten Verweigerung ohne Handlungsvorschlag; ein Rebutting Defeater liefert eine alternative Konklusion samt eigener Primärquelle. Der Vorrang wird durch den Wissensingenieur bei der Regelkonstruktion explizit festgelegt; zwei unpriorisierte, konfligierende Regeln führen deterministisch zur Verweigerung.

---

### 4.4. Dung-Grounded-Semantik und deterministische Auflösung zyklischer Argumentationsfallen (Grounded Extension)

Die Formalisierung von Defeatern nach Pollock löst den lokalen Konflikt zwischen einem Paar zueinander in Widerspruch stehender Regeln auf. In realen industriellen Expertensystemen bilden Normen jedoch engmaschige Netze wechselseitiger Restriktionen. Wirken im Beweisraum hunderte Standards simultan, erzeugen lokale Falsifikationsrelationen komplexe Angriffsgraphen, in denen wechselseitige Zyklen oder fehlerhafte Verdrängungskaskaden von Wissen entstehen können. Ohne eine globale mathematische Akzeptanzsemantik droht dem System eine unendliche Rekursion oder die Ausgabe nichtdeterministischer Verdikte, die von der zufälligen Reihenfolge der Regelauswertung im Speicher abhängen.

Zur rigorosen Formalisierung solcher Netze dient Dungs System der abstrakten Argumentation (*Argumentation Framework*, $`\mathcal{AF}`$) [[17]](#src-17), definiert als Tupel $`\mathcal{AF} = \langle \mathcal{A}, \mathcal{R}_{\mathrm{defeat}} \rangle`$, worin $`\mathcal{A}`$ eine endliche Menge von Argumenten darstellt und $`\mathcal{R}_{\mathrm{defeat}} \subseteq \mathcal{A} \times \mathcal{A}`$ die binäre Falsifikations- bzw. Angriffsrelation (*Attack/Defeat*) beschreibt. Im strukturierten ASPIC+-Modell [[7]](#src-7) manifestiert sich eine Relation $`(B, A) \in \mathcal{R}_{\mathrm{defeat}}`$ über drei distinkte Mechanismen:
1. **Untergraben von Prämissen (*Undermining*):** Angriff auf anfechtbare Präsuppositionen oder primäre Beobachtungsprämissen des Arguments $`A`$ ($`\mathrm{Prem}(A)`$);
2. **Widerlegung der Konklusion (*Rebutting*):** Direkte Kollision entgegengesetzter Konklusionen ($`\mathrm{Conc}(A) \equiv \neg \mathrm{Conc}(B)`$), die nur dann zu einer erfolgreichen Falsifikation führt, sofern das angreifende Argument dem attackierten gemäß der definierten Präferenzordnung nicht unterlegen ist ($`B \not< A`$);
3. **Untergraben der Inferenzregel (*Undercutting*):** Angriff auf die epistemische Legitimität des Inferenzschritts nach Regel $`r`$, der unabhängig von Prioritätsordnungen stets bedingungslos erfolgreich ist ($`B \to \mathrm{defeat}(A)`$).

Zur Bestimmung der Menge jener Argumente, die der logische Inferenzkern als unanfechtbar akzeptabel anerkennen muss, wird **Dungs charakteristischer Operator** $`\Gamma_{\mathcal{AF}} : 2^{\mathcal{A}} \to 2^{\mathcal{A}}`$ eingeführt:

```math
\Gamma_{\mathcal{AF}}(S) = \{ A \in \mathcal{A} \mid \forall B \in \mathcal{A} \, ((B, A) \in \mathcal{R}_{\mathrm{defeat}} \implies \exists C \in S : (C, B) \in \mathcal{R}_{\mathrm{defeat}}) \}.
```

Parameter des charakteristischen Dung-Operators:
- $`S \subseteq \mathcal{A}`$ ist eine beliebige Argumentmenge, die als potentielle Verteidigerbasis fungiert;
- $`A \in \mathcal{A}`$ ist das auf Akzeptanz geprüfte Argument;
- $`(B, A) \in \mathcal{R}_{\mathrm{defeat}}`$ bezeichnet den Angriff des Arguments $`B`$ auf das Argument $`A`$;
- $`(C, B) \in \mathcal{R}_{\mathrm{defeat}}`$ ist der verteidigende Gegenangriff des Arguments $`C \in S`$ auf den Angreifer $`B`$; der Operator selektiert exakt jene Argumente, deren Angreifer allesamt durch Elemente aus $`S`$ erfolgreich widerlegt werden ($`A`$ wird durch $`S`$ verteidigt).

Der Operator $`\Gamma_{\mathcal{AF}}`$ ist bezüglich der Mengeninklusion monoton ($`S_1 \subseteq S_2 \implies \Gamma_{\mathcal{AF}}(S_1) \subseteq \Gamma_{\mathcal{AF}}(S_2)`$). Gemäß dem Fixpunktsatz von Tarski besitzt er einen eindeutigen kleinsten Fixpunkt (*least fixed point*, lfp), der die **fundierte Erweiterung** (*Grounded Extension*, GE) konstituiert:

```math
\mathrm{GE}(\mathcal{AF}) = \mathrm{lfp}(\Gamma_{\mathcal{AF}}) = \bigcup_{k=0}^{\infty} \Gamma_{\mathcal{AF}}^k(\emptyset).
```

Parameter der fundierten Erweiterung:
- $`\emptyset`$ ist die initiale leere Menge der Verteidiger (Basisschritt $`\Gamma_{\mathcal{AF}}^0(\emptyset) = \emptyset`$);
- $`\Gamma_{\mathcal{AF}}^1(\emptyset)`$ ist die Menge unangefochtener Argumente, die keinerlei Angreifern ausgesetzt sind (unwiderlegbare Axiome und Primärfakten);
- $`\Gamma_{\mathcal{AF}}^k(\emptyset)`$ ist die $`k`$-te Iteration der Verteidigung, die jene Argumente restituiert, deren Angreifer in vorangegangenen Schritten eliminiert wurden;
- $`\mathrm{GE}(\mathcal{AF})`$ ist die deterministische Erweiterung, berechenbar in polynomieller Zeit $`O(\lvert \mathcal{A} \rvert \cdot \lvert \mathcal{R}_{\mathrm{defeat}} \rvert)`$.

Zur vollständigen Statuszuweisung an sämtliche Knoten des Graphen dient das dreiwertige Caminada-Labeling (*Caminada Tri-State Labeling*):
- `IN` ($`A \in \mathrm{GE}`$): Das Argument ist vom System vollständig akzeptiert (alle seine Angreifer tragen das Label `OUT`);
- `OUT`: Das Argument ist verworfen (es existiert mindestens ein Angreifer mit dem Status `IN`);
- `UNDEC`: Das Argument verharrt im Schwebezustand ungelöster symmetrischer Konflikte oder unpriorisierter Zyklen.

```mermaid
flowchart TB
    accTitle: Algorithmus der Grounded-Evaluation und Caminada-Kennzeichnung
    accDescr: Topologische Tarjan-Analyse und iterative Berechnung des kleinsten Fixpunkts des Dung-Operators.

    G["Argumentationsgraph AF = (A, R)"] --> SCC["Tarjan-SCC-Analyse starker Zusammenhangskomponenten O(V+E)"]
    SCC --> CHK{"Symmetrischer Zyklus<br/>ohne Präferenz erkannt?"}
    
    CHK -- "Ja" --> UND["Status UNDEC für<br/>Zyklusknoten zuweisen"]
    UND --> REF["Sicherheitsperimeter:<br/>REFUSAL_CYCLIC_DEPENDENCY<br/>(Fail-Closed)"]
    
    CHK -- "Nein" --> FP["Dung-Operator-Iteration:<br/>GE = lfp(Gamma)"]
    FP --> LAB["Caminada-Labeling:<br/>IN / OUT / UNDEC"]
    LAB --> VER{"Status des Zielarguments"}
    
    VER -- "IN" --> ACC["Verdikt: ACCEPT<br/>(Aktion zulässig)"]
    VER -- "OUT" --> BLK["Verdikt: REJECT<br/>(Aktion unzulässig)"]
    VER -- "UNDEC" --> ESC["Verdikt: REFUSAL<br/>(Eskalation an Ingenieur)"]

    classDef proc fill:#e1f5fe,stroke:#0288d1,stroke-width:2px,color:#01579b;
    classDef dec fill:#fff9c4,stroke:#fbc02d,stroke-width:2px,color:#f57f17;
    classDef safe fill:#e8f5e9,stroke:#2e7d32,stroke-width:2px,color:#1b5e20;
    classDef err fill:#ffebee,stroke:#d32f2f,stroke-width:2px,color:#b71c1c;

    class G,SCC,FP,LAB proc;
    class CHK,VER dec;
    class ACC safe;
    class UND,REF,BLK,ESC err;
```

#### 4.4.1. Geschlossener Kreis der ingenieurtechnischen Berechnung der Grounded-Semantik

1. **Laufzeitsteuerung und Rechenfluss (Control Flow & Runtime Decisions):**
   Die Berechnung der fundierten Erweiterung $`\mathrm{GE}(\mathcal{AF})`$ determiniert unmittelbar den Zustand des endlichen Entscheidungsautomaten:
   - Erhält das Argument $`A_{\mathrm{action}}`$ das Label `IN`, erteilt das System die Freigabe für das Stellglied (`ACCEPT`) unter Generierung eines kryptographischen Konformitätsnachweises;
   - Wird das Argument als `OUT` klassifiziert, ergeht eine begründete Zurückweisung mit Auflistung der aktiven Falsifikatoren (`REJECT`);
   - Verbleibt das Argument infolge eines symmetrischen Angriffs oder unpriorisierten Zyklus im Zustand `UNDEC`, verwirft die Laufzeitumgebung beide Alternativen kategorisch und wechselt in den Sicherheitsmodus `REFUSAL_UNRESOLVED_SYMMETRIC_CONFLICT` (`Fail-Closed`), wodurch potenziell gefahrbringende Operationen gesperrt und Operatoren alarmiert werden.
2. **Hardware-Dimensionierung und infrastrukturelle Beschränkungen (Hardware Dimensioning):**
   Der Tarjan-Algorithmus [[19]](#src-19) und der Iterator des Operators $`\Gamma_{\mathcal{AF}}`$ verlangen ein strikt fixes Speicherbudget:
   - Die Tiefe des DFS-Rekursionsstacks zur SCC-Ermittlung ist durch die maximale Kardinalität der Argumente $`\lvert \mathcal{A} \rvert`$ beschränkt: Die Stack-Allokation beträgt $`M_{\mathrm{stack}} = \lvert \mathcal{A} \rvert \times 32\,\text{Byte}`$;
   - Bitmasken für die dreiwertigen Caminada-Zustände erfordern lediglich $`2 \times \lvert \mathcal{A} \rvert`$ Bit im statischen BRAM des Mikrocontrollers oder FPGA;
   - Die zeitliche Komplexität ist mit $`O(\lvert \mathcal{A} \rvert + \lvert \mathcal{R} \rvert)`$ für die Zykluserkennung und $`O(k \cdot \lvert \mathcal{R} \rvert)`$ für $`k \le \mathrm{diam}(\mathcal{AF})`$ Relaxationsschritte streng deterministisch, was die Einhaltung der maximalen Reaktionszeit (Worst-Case Execution Time, WCET) ohne exponentielle Explosion $`2^N`$ garantiert.
3. **Praktisches numerisches Rechenbeispiel (Worked Numerical Example):**
   Betrachten wir das Subsystem einer automatischen Notbremsung eines fahrerlosen U-Bahn-Zugs mit vier interagierenden Argumenten:
   - $`A_1`$: „Betriebliche elektrodynamische Bremsung im Gefälle aktivieren“ (gestützt auf Geschwindigkeitssensorik);
   - $`A_2`$: „Elektrodynamische Bremsung wegen Überhitzung der Bremswiderstände sperren“ (Rebutting Defeater $`A_2 \xrightarrow{\mathrm{rebut}} A_1`$);
   - $`A_3`$: „Anomale Temperaturmessung: Leitungsbruch an Thermistor TB-12“ (Undercutting Defeater $`A_3 \xrightarrow{\mathrm{undercut}} A_2`$);
   - $`A_4`$: „Messwert von Sensor TB-12 durch redundantes Pyrometer verifiziert“ ($`A_4 \xrightarrow{\mathrm{rebut}} A_3`$).

   Im ersten Schritt gilt $`\Gamma_{\mathcal{AF}}^1(\emptyset) = \{ A_4 \}`$, da $`A_4`$ keinen Gegenargumenten ausgesetzt ist. Im zweiten Schritt verteidigt $`A_4 \to \mathrm{IN}`$ den Gegenangriff gegen $`A_3`$, woraus $`A_3 \to \mathrm{OUT}`$ resultiert. Im dritten Schritt annulliert der verworfene Status von $`A_3`$ (`OUT`) dessen Angriff auf $`A_2`$, sodass $`A_2`$ das Label `IN` erhält. Im vierten Schritt falsifiziert $`A_2 \to \mathrm{IN}`$ das Argument $`A_1`$ erfolgreich, was $`A_1 \to \mathrm{OUT}`$ erzwingt.

   Das System konvergiert deterministisch nach 4 Iterationen ($`42\,\mu\text{s}`$ auf einem ARM Cortex-R5F bei $`\lvert \mathcal{A} \rvert = 4, \lvert \mathcal{R} \rvert = 3`$) zur eindeutigen Erweiterung $`\mathrm{GE} = \{ A_4, A_2 \}`$. **Systementscheidung:** Argument $`A_1`$ ist verworfen (`OUT`), die elektrodynamische Betriebsbremse wird gesperrt, und das System aktiviert gemäß Ausfallsicherheitsprotokoll die pneumatische mechanische Notbremse.

---

### 4.5. Epistemische Argumentrehabilitation und transitive Falsifikationsausbreitung in Sicherheitszielbäumen

Der Praxiseinsatz von Argumentationssystemen in sicherheitskritischen Domänen (ISO 26262 ASIL D für Kraftfahrzeuge [[18]](#src-18), DO-178C für die Avionik) ist mit der inhärent dynamischen Natur des Systemwissens konfrontiert. Werden durch die Umwelt neue Fakten generiert oder Sensor-Firmwares aktualisiert, muss die Inferenzmaschine zwei polare Anforderungen erfüllen: Regeln, deren Falsifikatoren sich als fehlerhaft erwiesen haben, epistemisch rehabilitieren, und gleichzeitig jede neu verifizierte Widerlegung unmittelbar bis zur obersten Ebene des Sicherheitszielbaums kaskadieren.

#### 4.5.1. Mechanismus der epistemischen Argumentrehabilitation (Reinstatement)
In naiven Inferenzmaschinen führt die Widerlegung eines Arguments $`A`$ durch den Gegenspieler $`B`$ zur irreversiblen Löschung von $`A`$ aus dem Arbeitsspeicher (epistemische Amnesie). Trifft im nachfolgenden Abtastzyklus die Evidenz $`C`$ ein, die den Angreifer $`B`$ diskreditiert, vermag das naive System $`A`$ ohne vollständigen Neustart der Wissensbasis nicht wieder in den aktiven Inferenzprozess einzugliedern.

In Dungs Grounded-Semantik ist die Argumentrehabilitation ein fundamentaler Lehrsatz: Gilt $`C \to \mathrm{IN}`$ und $`(C, B) \in \mathcal{R}_{\mathrm{defeat}}`$, so wird $`B`$ als `OUT` gekennzeichnet. Folglich wird die Menge aktiver Angreifer gegen $`A`$ leer, und $`A`$ wird auf demselben Fixpunkt $`\mathrm{lfp}(\Gamma_{\mathcal{AF}})`$ automatisch rehabilitiert ($`A \to \mathrm{IN}`$). Dies eliminiert epistemische Amnesie und sichert die ununterbrochene Kontinuität des Inferenzzyklus im operativen Betrieb.

#### 4.5.2. Transitive Falsifikationsausbreitung in ISO-26262-Sicherheitszielbäumen
In der Goal Structuring Notation (GSN) zerlegt ein Argumentationsbaum $`\mathcal{T} = \langle \mathcal{G}, \mathcal{E}, \mathcal{R} \rangle`$ ein übergeordnetes Sicherheitsziel (*Safety Goal* der Einstufung ASIL D, etwa $`G_0`$: „Vermeidung unkontrollierter Blockaden des Lenkaktuators“) über intermediäre Sicherheitsanforderungen $`G_i`$ bis hinab zu hardwarenahen Evidenzblättern $`e \in \mathcal{E}`$.

Ein Sicherheitsaudit offenbart eine typische Schwachstelle: Registriert ein Lenkwinkelsensor eine Spannungsabweichung außerhalb des Toleranzbands, markiert eine isolierte Regel lediglich das lokale Sensorfaktum als ungültig. Das übergeordnete Ziel $`G_0`$ verbleibt in naiven Architekturen fälschlich als valide (`VALID`), da seine logische Klausel keine explizite Bedingung zur Sensorspannung enthält (Defekt des „stillen Übergehens“ bzw. *Silent Failure Pass*).

Zur Gewährleistung funktionaler Sicherheit nach ISO 26262 wird der **Operator transitiver Kaskadenfalsifikation** $`\mathrm{Status}(G_i)`$ definiert:

```math
\mathrm{Status}(G_i) = \begin{cases} 
\mathrm{DIRECTLY\_FALSIFIED}, & \text{falls } \exists d \in \mathrm{Defeaters}(G_i) : \mathrm{Evaluated}(d) = \mathrm{True}, \\
\mathrm{TRANSITIVELY\_FALSIFIED}, & \text{falls } \exists G_j \in \mathrm{Children}(G_i) : \mathrm{Status}(G_j) \in \{\mathrm{DIRECTLY\_FALSIFIED}, \mathrm{TRANSITIVELY\_FALSIFIED}\}, \\
\mathrm{REFUSAL}, & \text{falls } \exists G_j \in \mathrm{Children}(G_i) : \mathrm{Status}(G_j) = \mathrm{REFUSAL}, \\
\mathrm{VALID}, & \text{sonst}.
\end{cases}
```

Parameter des Operators transitiver Falsifikation:
- $`G_i \in \mathcal{G}`$ ist ein Knoten des Argumentationsbaums (Sicherheitsziel oder Argumentationsstrategie);
- $`\mathrm{Defeaters}(G_i)`$ bezeichnet die Menge lokaler Falsifikatoren des Ziels $`G_i`$;
- $`\mathrm{Children}(G_i)`$ ist die Menge direkter Unterziele oder beweisführender Evidenzblätter;
- $`\mathrm{Evaluated}(d)`$ ist das Prädikat zur Auswertung der Ausnahmebedingung des Defeaters $`d`$;
- $`\mathrm{DIRECTLY\_FALSIFIED}`$ kennzeichnet die unmittelbare lokale Falsifikation durch einen aktiven Defeater;
- $`\mathrm{TRANSITIVELY\_FALSIFIED}`$ ist der vererbte Falsifikationsstatus, der bis zum Wurzelziel $`G_0`$ propagiert wird;
- $`\mathrm{REFUSAL}`$ signalisiert die Verweigerung infolge Datenunvollständigkeit oder Ungewissheit eines Kindknotens;
- $`\mathrm{VALID}`$ repräsentiert den bestätigten Status bei vollständiger Abwesenheit von Falsifikatoren im gesamten Teilbaum.

#### 4.5.3. Strikte Fail-Closed-Verweigerung bei unbekannten Entitäten und deontische Kontrolle (ModalityMust)
1. **Strikte Fail-Closed-Verweigerung bei mangelnder Evidenz (Strict Fail-Closed on Unknowns):** Besitzt eine Aktion oder ein physikalischer Parameter $`X`$ keine normativen Relationen in der Wissensbasis, verweigert das System standardmäßig jegliche positive Ableitung und generiert das Verdikt `REFUSAL_INSUFFICIENT_EVIDENCE` mit dem Statusflag `EvidenceVerified = false`.
2. **Deontischer Vorrang zwingender Pflichten (`ModalityMust`):** Jeder Aktionsvorschlag oder Ausnahmezweig, der eine zwingende Sicherheitsvorschrift zu unterdrücken, umgehen oder aufzuschieben versucht (beispielsweise das Überspringen zyklischer BIST-Speichertests oder das Ignorieren eines Türschalter-Interlocks einer Hochspannungsschaltzelle), generiert automatisch einen hochprioritären Undercutting Defeater, der das betreffende Sicherheitsziel unmittelbar in den Status `DIRECTLY_FALSIFIED` überführt.

```mermaid
flowchart TD
    accTitle: Kaskade transitiver Falsifikation im GSN-Sicherheitszielbaum
    accDescr: Falsifikation einer hardwarenahen Evidenz löst eine sofortige Fehlerkaskade zum ASIL-D-Wurzelziel aus.

    G0["<b>Sicherheits-Wurzelziel G0 (ASIL D)</b><br/>Status: TRANSITIVELY_FALSIFIED"]
    G1["Teilziel G1: Lenkwinkelregelung<br/>Status: TRANSITIVELY_FALSIFIED"]
    G2["Teilziel G2: Servostromüberwachung<br/>Status: VALID"]
    
    E1["Evidenzblatt E1: Hall-Winkelsensor<br/>Status: DIRECTLY_FALSIFIED"]
    E2["Evidenzblatt E2: Redundanter Encoder<br/>Status: VALID"]
    E3["Evidenzblatt E3: Phasenstrom-Shunt A<br/>Status: VALID"]
    
    DEF["<b>Untergrabender Defeater d_fault:</b><br/>Offsetspannung > 450 mV (Nullpunktdrift)"]

    DEF -->|"Widerlegt Prämisse"| E1
    E1 -->|"Kaskadierende Propagation"| G1
    G2 --> G0
    G1 -->|"Induzierte Falsifikation"| G0
    
    E2 --> G1
    E3 --> G2

    G0 ==>|"Auslösung des Notfallprotokolls"| INTERLOCK["Hardware-Interlock:<br/>Übergang in sicheren Zustand (Safe State)<br/>Sperrung der Leistungsendstufe"]

    classDef valid fill:#e8f5e9,stroke:#2e7d32,stroke-width:2px,color:#1b5e20;
    classDef fals fill:#ffebee,stroke:#d32f2f,stroke-width:2px,color:#b71c1c;
    classDef alert fill:#fff3e0,stroke:#e65100,stroke-width:2px,color:#bf360c;

    class G2,E2,E3 valid;
    class G0,G1,E1 fals;
    class DEF,INTERLOCK alert;
```

#### 4.5.4. Geschlossener Kreis der ingenieurtechnischen Berechnung der transitiven Falsifikation

1. **Laufzeitsteuerung und Rechenfluss (Control Flow & Runtime Decisions):**
   Der berechnete Status des Wurzelziels $`\mathrm{Status}(G_0)`$ steuert unmittelbar das Schutzsystem:
   - Der Status `VALID` gestattet die Ausstellung eines Vertrauenszertifikats und die Ausgabe von Ansteuerimpulsen an Leistungsaktoren;
   - Der Status `TRANSITIVELY_FALSIFIED` auf beliebiger Hierarchieebene annulliert das Zertifikat augenblicklich, deaktiviert die Pulsweitenmodulation (PWM) der Leistungsendstufen und initiiert das Notfallprotokoll gemäß ISO 26262;
   - Der Status `REFUSAL` versetzt das System in einen sicheren Wartezustand auf fehlende Evidenz unter Blockierung des Übergangs in den aktiven Betriebsmodus.
2. **Hardware-Dimensionierung und infrastrukturelle Beschränkungen (Hardware Dimensioning):**
   - Die Baumstruktur $`\mathcal{T}`$ wird in einem einzigen Post-Order-Durchlauf von den Blättern zur Wurzel (Bottom-Up Post-Order Traversal) mit linearer Komplexität $`O(\lvert \mathcal{G} \rvert + \lvert \mathcal{E} \rvert)`$ verarbeitet;
   - Für einen Baum mit einer Tiefe von $`H = 6`$ Ebenen und $`N = 64`$ Knoten beansprucht die Kaskadenberechnung auf einem ARM Cortex-M7 ($`400\,\text{MHz}`$) maximal $`2.8\,\mu\text{s}`$. Dies unterschreitet das vorgeschriebene Fehlertoleranzzeitintervall (Fault Tolerant Time Interval, FTTI), das im Automobilbereich typischerweise zwischen $`10`$ und $`50\,\text{ms}`$ liegt, um Größenordnungen;
   - Der Bitvektor der Knotenstatus belegt exakt $`2 \times 64 = 128`$ Bit, was parallele bitweise Reduktionsoperationen über SIMD-Instruktionen des Prozessors in minimalen Taktzyklen ermöglicht.
3. **Praktisches numerisches Rechenbeispiel (Worked Numerical Example):**
   Betrachten wir die Schubvektorsteuerung eines unbemannten Fluggeräts mit dem ASIL-D-Wurzelziel $`G_0`$, unterteilt in 3 Teilziele $`G_1, G_2, G_3`$ und 9 Evidenzblätter. Während eines Flugmanövers detektiert der Drucksensor der Hydraulikleitung $`E_7`$ einen Druckabfall unter $`12\,\text{MPa}`$ (Nennwert: $`21\,\text{MPa}`$), was den Defeater $`d_{\mathrm{hydr}}`$ aktiviert.

   Berechnung gemäß Formel:
   1. $`\mathrm{Status}(E_7) = \mathrm{DIRECTLY\_FALSIFIED}`$ (Prädikat $`\mathrm{Evaluated}(d_{\mathrm{hydr}}) = \mathrm{True}`$);
   2. $`\mathrm{Status}(G_3) = \mathrm{TRANSITIVELY\_FALSIFIED}`$, da $`E_7 \in \mathrm{Children}(G_3)`$;
   3. $`\mathrm{Status}(G_0) = \mathrm{TRANSITIVELY\_FALSIFIED}`$, da $`G_3 \in \mathrm{Children}(G_0)`$.

   **Systementscheidung:** Obwohl die verbleibenden 8 Sensoren den Status `VALID` aufweisen, verifiziert das System binnen $`3.2\,\mu\text{s}`$ die Kompromittierung des Wurzelziels $`G_0`$, entzieht die Flugfreigabe, drosselt die Turbine und wirft den Notfallschirm aus.

---

## 5. Untersuchung von Konflikten zwischen Dokumenten: Evolution des TCP-Protokolls

Die Evolution des Transmission Control Protocol (TCP) demonstriert mustergültig, warum eine neuere Revision keineswegs das pauschale Verbot einer früheren Verhaltensweise bedeutet. RFC 793 definierte das ursprüngliche Protokollverhalten [[8]](#src-8); RFC 5961 führte Schutzmechanismen gegen Blind-Reset-Angriffe ein [[9]](#src-9); RFC 9293 schließlich ersetzte RFC 793 vollständig und differenziert in Abschnitt 3.10.7.4 explizit zwischen Implementierungen mit und ohne Unterstützung dieses Schutzes [[10]](#src-10).

Für die in Abschnitt 3.10.7.4 spezifizierten synchronisierten Zustände existieren bei aktiver Unterstützung von RFC 5961 drei disjunkte Prüffälle bei Empfang eines Reset-Segments (RST):

| Bedingung der Sequenznummer | Aktion | Unzulässige Fehlinterpretation |
|---|---|---|
| Außerhalb des aktuellen Empfangsfensters | Segment verwerfen, keine Antwort | Dass Resets generell unzulässig seien |
| Exakt gleich `RCV.NXT` (nächste erwartete Nummer) | Verbindung gemäß Zustand sofort zurücksetzen | Dass stets ein Challenge-ACK erforderlich sei |
| Im Fenster, aber ungleich `RCV.NXT` | Challenge-ACK senden und Segment verwerfen | Dass das bloße Treffen des Fensters zum Reset genügt |

Ohne Unterstützung von RFC 5961 bewahrt RFC 9293 das historische Verhalten von RFC 793: Nach Abschnitt 3.5.3 ist ein Reset verbindlich, sobald seine Sequenznummer in das Empfangsfenster fällt; Segmente außerhalb werden stillschweigend verworfen.

Sind der Verbindungszustand oder das unterstützte Sicherheitsprofil unbekannt, verweigert das Expertensystem die Entscheidung und fordert eine Präzisierung. Eine Dokumenten-Ablösungsrelation ersetzt weder diese Eingangsdaten noch die Prüfung der Geltungsbedingungen der anwendbaren Norm.

### 5.1. Ausführbares Normenauswahlmodell in ASP

Die obige Entscheidungstabelle lässt sich als Programm in Answer Set Programming (ASP) abbilden und mit dem Solver clingo ausführen, der von Martin Gebser, Roland Kaminski, Benjamin Kaufmann und Torsten Schaub entwickelt wurde [[11]](#src-11). ASP eignet sich aus zwei Gründen hervorragend für diese Aufgabe: Erstens erlaubt Negation as Failure die elegante Kodierung der Regel „ein Dokument gilt, sofern es nicht durch ein anderes obsoletiert wurde“ in einer einzigen Zeile. Zweitens garantiert ein stratifiziertes Programm (ein Programm ohne Rekursion über Negation) genau ein stabiles Modell (*Answer Set*), wodurch die Schlussfolgerung unabhängig von der Auswertungsreihenfolge ist. Negation as Failure birgt jedoch das inhärente Risiko, dass ein fehlender Eintrag über eine Obsoletierung die veraltete Norm fälschlich als aktiv einstuft. Deshalb implementiert das Modell Fail-Closed-Regeln für sechs Szenarien: fehlendes Eingangsfaktum; Verbindungszustand außerhalb des Modells; Norm referenziert nicht registriertes Dokument; Registervollständigkeit nicht bestätigt; zwei aktive Normen konfligieren; keine Norm anwendbar.

<details>
<summary>Normenauswahlmodell in ASP (Datei norms.lp)</summary>

```prolog
% Dokumentenregister und Ablösungsfolge (Geltungslinie) nach Angaben des RFC Editor.
document(rfc793). document(rfc9293).
obsoletes(rfc9293, rfc793).
% Register mit dem Index des RFC Editor für den Bereich „RST-Verarbeitung in TCP“ abgeglichen.
complete(tcp_rst).

% norm(Dokument, Schutz RFC 5961, Lage der Sequenznummer, Aktion).
% RFC 9293, Abschnitt 3.10.7.4: Mit Schutzunterstützung greifen drei Prüfungen.
norm(rfc9293, yes, out_of_window, drop).
norm(rfc9293, yes, exact, reset).
norm(rfc9293, yes, in_window_not_exact, challenge_ack).
% Ohne Schutz ist der Reset gültig, falls die Nummer ins Empfangsfenster fällt (Abschnitt 3.5.3).
norm(rfc9293, no, out_of_window, drop).
norm(rfc9293, no, exact, reset).
norm(rfc9293, no, in_window_not_exact, reset).
% RFC 793 kannte den Schutz nach RFC 5961 nicht und setzte die Verbindung bei jeder Nummer im Fenster zurück.
norm(rfc793, P, out_of_window, drop) :- protection(P).
norm(rfc793, P, exact, reset) :- protection(P).
norm(rfc793, P, in_window_not_exact, reset) :- protection(P).
protection(yes; no).

% Das Modell umfasst nur synchronisierte Zustände, in denen ein RST die Verbindung nach CLOSED überführt.
synchronized(established; fin_wait_1; fin_wait_2; close_wait; closing; last_ack; time_wait).

% Geltende Revision: ein Dokument aus dem Register, das nicht abgelöst (obsoletiert) wurde.
superseded(D) :- obsoletes(_, D).
active(D) :- document(D), not superseded(D).

% Fail-Closed: Jede Verweigerungsursache blockiert die Antwort und benennt das fehlende Faktum.
required(state; protection; seq).
known(K) :- input(K, _).
clarify(missing(K)) :- required(K), not known(K).
clarify(out_of_model(S)) :- input(state, S), not synchronized(S).
clarify(unregistered_source(D)) :- norm(D, _, _, _), not document(D).
clarify(registry_incomplete) :- not complete(tcp_rst).

candidate(D, A) :- active(D), input(protection, P), input(seq, Q), norm(D, P, Q, A).
clarify(conflict(A1, A2)) :- candidate(_, A1), candidate(_, A2), A1 < A2.
has_candidate :- candidate(_, _).
clarify(no_applicable_norm) :- known(protection), known(seq), not has_candidate.

blocked :- clarify(_).
answer(A, D) :- candidate(D, A), not blocked.

#show answer/2.
#show clarify/1.
```

</details>

Das folgende Python-Skript injiziert Fakten für neun Testfälle in das ASP-Modell. In drei Szenarien wird gezielt ein Registerfaktum entfernt, um Defekte der Wissensbasis selbst zu simulieren. Zur Ausführung ist das Paket clingo erforderlich (`pip install clingo`); das Skript wurde unter Python 3.14 mit clingo 5.8.2 verifiziert. Beide Dateien werden im selben Verzeichnis abgelegt und mittels `python check_norms.py` ausgeführt.

<details>
<summary>Python-Testskript: Validierung der Modell-Szenarien (Datei check_norms.py)</summary>

```python
"""Führt das Normenauswahlmodell norms.lp für Testfälle aus und verifiziert die erwarteten Antworten."""
from pathlib import Path

import clingo

PROGRAM = Path(__file__).with_name("norms.lp").read_text(encoding="utf-8")


def solve(facts, drop=()):
    program = PROGRAM
    for fact in drop:
        assert fact in program, fact
        program = program.replace(fact, "", 1)  # entfernt nur das Faktum, nicht die Erwähnung in der Regel
    ctl = clingo.Control(["--warn=none"])
    ctl.add("base", [], program + "\n" + facts)
    ctl.ground([("base", [])])
    models = []
    ctl.solve(on_model=lambda m: models.append(sorted(str(s) for s in m.symbols(shown=True))))
    assert len(models) == 1, models  # ein stratifiziertes Programm besitzt genau ein stabiles Modell
    return models[0]


SYNC = "input(state, established). "
SCENARIOS = [
    ("Mit Schutz, Nummer im Fenster, aber nicht exakt", SYNC + "input(protection, yes). input(seq, in_window_not_exact).",
     (), ["answer(challenge_ack,rfc9293)"]),
    ("Mit Schutz, Nummer exakt", SYNC + "input(protection, yes). input(seq, exact).",
     (), ["answer(reset,rfc9293)"]),
    ("Ohne Schutz, Nummer im Fenster", SYNC + "input(protection, no). input(seq, in_window_not_exact).",
     (), ["answer(reset,rfc9293)"]),
    ("Nummer außerhalb des Fensters", SYNC + "input(protection, yes). input(seq, out_of_window).",
     (), ["answer(drop,rfc9293)"]),
    ("Schutzstatus unbekannt", SYNC + "input(seq, in_window_not_exact).",
     (), ["clarify(missing(protection))"]),
    ("Zustand LISTEN außerhalb des Modells", "input(state, listen). input(protection, yes). input(seq, exact).",
     (), ["clarify(out_of_model(listen))"]),
    ("Ablösungsfaktum fehlt", SYNC + "input(protection, yes). input(seq, in_window_not_exact).",
     ("obsoletes(rfc9293, rfc793).",), ["clarify(conflict(challenge_ack,reset))"]),
    ("RFC 9293 nicht im Register", SYNC + "input(protection, yes). input(seq, in_window_not_exact).",
     ("document(rfc9293).",), ["clarify(no_applicable_norm)", "clarify(unregistered_source(rfc9293))"]),
    ("Registervollständigkeit unbestätigt", SYNC + "input(protection, yes). input(seq, exact).",
     ("complete(tcp_rst).",), ["clarify(registry_incomplete)"]),
]

for name, facts, drop, expected in SCENARIOS:
    got = solve(facts, drop)
    assert got == expected, (name, got, expected)
    print(f"{name:<45} -> {', '.join(got)}")
print("alle Szenarien stimmen mit den Erwartungen überein")
```

</details>

Die Ausführung von `python check_norms.py` liefert folgende Ausgabe:

<details>
<summary>Programmausgabe</summary>

```text
Mit Schutz, Nummer im Fenster, aber nicht exakt -> answer(challenge_ack,rfc9293)
Mit Schutz, Nummer exakt                      -> answer(reset,rfc9293)
Ohne Schutz, Nummer im Fenster                -> answer(reset,rfc9293)
Nummer außerhalb des Fensters                 -> answer(drop,rfc9293)
Schutzstatus unbekannt                        -> clarify(missing(protection))
Zustand LISTEN außerhalb des Modells          -> clarify(out_of_model(listen))
Ablösungsfaktum fehlt                         -> clarify(conflict(challenge_ack,reset))
RFC 9293 nicht im Register                    -> clarify(no_applicable_norm), clarify(unregistered_source(rfc9293))
Registervollständigkeit unbestätigt           -> clarify(registry_incomplete)
alle Szenarien stimmen mit den Erwartungen überein
```

</details>

Die ersten vier Szenarien bilden die Spezifikation von RFC 9293 präzise ab: Bei aktivem Schutz bewirkt eine Sequenznummer im Fenster, die nicht exakt übereinstimmt, ein Challenge-ACK, während ohne Schutz dieselbe Nummer zum sofortigen Reset führt. Das fünfte und sechste Szenario demonstrieren das Fail-Closed-Verhalten bei unvollständigen Eingaben: Das System unterstellt keine standardmäßige Schutzaktivierung und wendet die Regeln synchronisierter Zustände nicht auf den Zustand LISTEN an. Die letzten drei Szenarien simulieren Fehlerzustände der Wissensbasis: Ohne das Obsoletierungsfaktum gelten beide Revisionen simultan, woraufhin das Programm einen Konflikt meldet, statt stillschweigend eine Norm zu wählen. Fehlt RFC 9293 im Register, können dessen Regeln nicht feuern; da RFC 793 obsoletiert ist, meldet das Programm beide Verweigerungsgründe. Ist die Vollständigkeit des Registers nicht deklariert, verweigert das System selbst bei validen Eingaben: Ohne diese Zusicherung lässt sich das Fehlen eines Datensatzes nicht von Unkenntnis unterscheiden.

Die Modellgrenzen sind klar umrissen: Es behandelt ausschließlich das RST-Bit in synchronisierten Zuständen, empfängt die Lage der Sequenznummer als vorberechnetes Eingangsfaktum (statt sie aus `RCV.NXT` und `RCV.WND` zu berechnen) und ersetzt keine Netzwerk-Stack-Tests. In Produktionssystemen wird das Dokumentenregister aus maschinenlesbaren Indizes gespeist, und das Faktum `complete(tcp_rst)` wird erst nach kryptografischer Signaturprüfung des Indexstandes gesetzt.

Wird das Dokumentenregister über Shards verteilt, bezieht sich die Vollständigkeitserklärung `complete(tcp_rst)` auf die Gesamtheit aller Shards. Die Revisionen RFC 793, 5961 und 9293 bilden eine untrennbare Dokumentenfamilie und müssen im selben Shard partitioniert sein. Antwortet dieser Shard nicht (*Shard Silence*), verweigert das System mit `registry_incomplete`, anstatt unvollständige Normen verbliebener Shards anzuwenden. Die Sharding-Integritätsregeln wurden in [Kapitel 7](ch07-knowledge-base-typology.md) hergeleitet; ihre hochperformante Implementierung wird in Abschnitt 10 von [Kapitel 32](ch32-high-performance-knowledge-packs-mmap-and-harvesting.md) demonstriert.

---

## 6. Mehrstufige hypothetische Anfragen

Ingenieurtechnische Fragestellungen sind selten einfache Tatsachenabfragen. Häufiger treten komplexe hypothetische Konstruktionen auf:  
*„Falls unser Client das SMTP-Protokoll gemäß RFC 5321 implementiert: Ist er verpflichtet, den Befehl EHLO vor MAIL FROM zu senden, und wie muss er reagieren, wenn der Server ausschließlich das veraltete RFC 821 unterstützt?“*

Das nachfolgende Diagramm skizziert den Plan einer solchen mehrstufigen Inferenz. Die Abschnitte 3.2 und 4.1.1.1 von RFC 5321 differenzieren zwischen der präferierten Nutzung von `EHLO` mit Modalität `SHOULD` und der strikten Vorschrift (`MUST`), vor Beginn einer Transaktion entweder `HELO` oder `EHLO` auszuführen [[4]](#src-4). Die synoptische Prüfung mit Abschnitt 4.1.4 stellt sicher, dass Befehlsabfolgen und Abwärtskompatibilitätsmodi gewahrt bleiben.

```mermaid
flowchart TD
    accTitle: Ausführungsplan einer konditionalen syllogistischen Anfrage
    accDescr: Zerlegung der Anfrage in Hypothesen, Zielprädikate und Ausnahmezweige.

    Q["Komplexe hypothetische Anfrage"] --> PARSE["Semantischer Parser"]
    
    PARSE --> H["Hypothetische Prämissen:<br/>• Implements(Client, RFC-5321)<br/>• Step = PriorTo(MAIL_FROM)"]
    PARSE --> T["Zielbehauptung:<br/>Welcher Befehl und welche Modalität gelten?"]
    PARSE --> B["Kontextueller Ausnahmezweig:<br/>• SupportsOnly(Server, RFC-821)"]

    H & T --> DED1["<b>Anwendbare Norm</b><br/>Vorzugsweise EHLO: SHOULD;<br/>vor Transaktion HELO oder EHLO: MUST"]
    
    B --> DEF_EVAL["<b>Kompatibilitätsbedingung</b><br/>Wird EHLO nicht erkannt,<br/>SHOULD der Client auf HELO zurückfallen"]

    DED1 & DEF_EVAL --> FINAL["<b>Konditionale Antwort</b><br/>Modalität und Rückfallbefehl<br/>hängen von der geltenden Norm ab"]

    classDef q fill:#fff3e0,stroke:#e65100,stroke-width:2px,color:#bf360c;
    classDef p fill:#e1f5fe,stroke:#0288d1,stroke-width:2px,color:#01579b;
    classDef d fill:#ede7f6,stroke:#512da8,stroke-width:2px,color:#311b92;
    classDef r fill:#e8f5e9,stroke:#2e7d32,stroke-width:2px,color:#1b5e20;

    class Q q;
    class PARSE,H,T,B p;
    class DED1,DEF_EVAL d;
    class FINAL r;
```

---

## 7. Inspektion des Inferenzpfads

Ein Revisions- und Audit-Inspektor muss Regel, Begründungen, unbestimmte Eingaben und den Status der Konklusion lückenlos visualisieren. Das folgende Protokoll zeigt ein didaktisches Mockup eines Inferenzpfads. Byte-Offsets und Hashes sind hier bewusst abstrahiert: In realen Systemen werden sie direkt über den Artefakt-Bytes berechnet.

<details>
<summary>Didaktisches Mockup eines konditionalen Inferenzpfads</summary>

```text
Ziel: Ermittlung des Befehls zur Initiierung einer Mail-Transaktion
    Begründung: RFC 5321, Abschnitte 3.2, 4.1.1.1 und 4.1.4
    Bedingung: Entität ist ein SMTP-Client
        Status: Benutzerhypothese, kein verifiziertes Faktum
    Norm: Vorzugsweise EHLO (SHOULD), vor Transaktion HELO oder EHLO (MUST)
    Kompatibilitätsbedingung: Server erkennt EHLO nicht
        Status: Unbekannt, Antwort des Servers erforderlich
    Ergebnis: Lediglich bedingter Schluss, keine faktische Handlungsautorisierung
```

</details>

In einer Produktionsarchitektur verlinkt jede Begründung direkt auf die unveränderliche Primärquelle. Die Zitatanzeige belegt die Provenienz; die Korrektheit des Inferenzschritts wird unabhängig anhand formaler Regeln überprüft. Dieses Text-Mockup implementiert noch keine interaktive Navigation oder kryptografische Verifikation.

---

## 8. Didaktischer Inferenzschritt in Go

Das folgende Beispielprogramm implementiert einen elementaren Inferenzschritt über einer expliziten Regel und einer verifizierten Klassenzugehörigkeit. Es verzichtet bewusst auf eine vollständige mehrstufige Inferenzmaschine oder automatische Konfliktresolution. Ein unbekannter Ausnahmestatus erzwingt eine Verweigerung; Zykluseinschlüsse in die Hierarchie werden abgewiesen. Die beiden Defeater-Klassen führen zu distinkten Resultaten: Ein Rebutting Defeater liefert eine alternative Konklusion samt Quellenbeleg; ein Undercutting Defeater blockiert die Regelanwendung vollständig. Zur Ausführung ist Go 1.20 oder neuer erforderlich; externe Abhängigkeiten bestehen nicht.

<details>
<summary>Go-Referenzimplementierung: Prädikatenhierarchie und elementarer Inferenzschritt</summary>

```go
package syllogism

import (
	"errors"
	"fmt"
)

// Kleene3VL modelliert die dreiwertige Logik nach Kleene
type Kleene3VL int

const (
	Unknown Kleene3VL = 0
	True    Kleene3VL = 1
	False   Kleene3VL = -1
)

func (v Kleene3VL) And(other Kleene3VL) Kleene3VL {
	if v == False || other == False {
		return False
	}
	if v == True && other == True {
		return True
	}
	return Unknown
}

func (v Kleene3VL) Or(other Kleene3VL) Kleene3VL {
	if v == True || other == True {
		return True
	}
	if v == False && other == False {
		return False
	}
	return Unknown
}

func (v Kleene3VL) Not() Kleene3VL {
	return -v
}

type RelationHierarchy struct {
	parentMap map[string]string // child -> parent
}

func NewRelationHierarchy() *RelationHierarchy {
	return &RelationHierarchy{parentMap: map[string]string{
		"must_requirement":       "normative_requirement",
		"prohibited_requirement": "normative_requirement",
		"normative_requirement":  "concept",
		"obsoleted_by":           "lineage_relation",
	}}
}

func (hierarchy *RelationHierarchy) AddRelation(child, parent string) error {
	if child == "" || parent == "" || hierarchy.Subsumes(child, parent) {
		return errors.New("invalid relation or cycle")
	}
	if previous, exists := hierarchy.parentMap[child]; exists && previous != parent {
		return errors.New("this example supports one parent per predicate")
	}
	hierarchy.parentMap[child] = parent
	return nil
}

// Subsumes prüft, ob subRelation ein Spezialfall von superRelation ist
func (hierarchy *RelationHierarchy) Subsumes(superRelation, subRelation string) bool {
	seen := make(map[string]bool)
	for curr := subRelation; curr != ""; curr = hierarchy.parentMap[curr] {
		if seen[curr] {
			return false
		}
		seen[curr] = true
		if curr == superRelation {
			return true
		}
	}
	return false
}

// Provenance beschreibt die physikalische Primärquelle des Zitats
type Provenance struct {
	DocumentID string
	ByteStart  int
	ByteEnd    int
	SHA256     string
}

// Premise repräsentiert eine Prämisse des Urteils
type Premise struct {
	ID         string
	Subject    string
	Predicate  string
	Object     string
	Modality   string // MUST, MUST_NOT, SHOULD
	Provenance Provenance
}

// Defeater beschreibt eine vom Wissensingenieur an die Regel gekoppelte Ausnahmebedingung.
// Ein Undercutting Defeater (IsRebutting == false) annulliert lediglich die Anwendbarkeit der Regel.
// Ein Rebutting Defeater (IsRebutting == true) liefert eine abweichende Schlussfolgerung FallbackAction
// aus der eigenen Quelle Source; die Kopplung an die Regel kodiert ihren Vorrang.
type Defeater struct {
	ConditionPredicate string
	FallbackAction     string
	IsRebutting        bool
	Source             string
}

// Rule repräsentiert eine deduktive Norm
type Rule struct {
	MajorPremise Premise
	Defeaters    []Defeater
}

// Engine bildet den Kern des syllogistischen Deduktors
type Engine struct {
	hierarchy *RelationHierarchy
	rules     []Rule
}

func NewEngine(hierarchy *RelationHierarchy) *Engine {
	return &Engine{hierarchy: hierarchy}
}

// InferConclusion führt einen syllogistischen Inferenzschritt aus
func (e *Engine) InferConclusion(minorSubject, minorClass, requestedRelation string, envConditions map[string]Kleene3VL) (string, error) {
	if minorSubject == "" || minorClass == "" {
		return "", errors.New("confirmed object and class are required")
	}
	var matched *Rule
	for _, rule := range e.rules {
		if rule.MajorPremise.Subject != minorClass || !e.hierarchy.Subsumes(requestedRelation, rule.MajorPremise.Predicate) {
			continue
		}
		if matched != nil {
			return "", errors.New("multiple applicable rules require conflict resolution")
		}
		candidate := rule
		matched = &candidate
	}
	if matched == nil {
		return "", errors.New("insufficient facts for inference")
	}
	// Zunächst müssen alle Ausnahmebedingungen bekannt sein: Der Schluss darf nicht von der Prüfreihenfolge abhängen.
	var active []Defeater
	for _, exception := range matched.Defeaters {
		state, exists := envConditions[exception.ConditionPredicate]
		if !exists || (state != True && state != False) {
			return "", fmt.Errorf("exception state is unknown: %s", exception.ConditionPredicate)
		}
		if state == True {
			active = append(active, exception)
		}
	}
	premise := matched.MajorPremise
	switch {
	case len(active) > 1:
		return "", errors.New("several active defeaters require conflict resolution")
	case len(active) == 1 && !active[0].IsRebutting:
		return "", fmt.Errorf("rule from %s is not applicable while %s holds", premise.Provenance.DocumentID, active[0].ConditionPredicate)
	case len(active) == 1:
		return fmt.Sprintf("%s %s; source: %s; defeats rule from %s", minorSubject,
			active[0].FallbackAction, active[0].Source, premise.Provenance.DocumentID), nil
	}
	return fmt.Sprintf("%s %s %s; source: %s", minorSubject,
		premise.Modality, premise.Object, premise.Provenance.DocumentID), nil
}
```

Speichern Sie das Modul als `syllogism.go` und den nachfolgenden Testblock als `syllogism_test.go`; der Befehl `go test syllogism.go syllogism_test.go` führt die Verifikation aus.

```go
package syllogism

import (
	"strings"
	"testing"
)

func TestHierarchyAndInference(t *testing.T) {
	hierarchy := NewRelationHierarchy()
	if !hierarchy.Subsumes("concept", "must_requirement") {
		t.Fatal("generalization failed")
	}
	if hierarchy.Subsumes("must_requirement", "concept") {
		t.Fatal("reverse specialization accepted")
	}
	if hierarchy.AddRelation("concept", "must_requirement") == nil {
		t.Fatal("cycle accepted")
	}
	engine := NewEngine(hierarchy)
	engine.rules = []Rule{{
		MajorPremise: Premise{Subject: "Client", Predicate: "must_requirement", Object: "authenticate", Modality: "MUST"},
		Defeaters:    []Defeater{{ConditionPredicate: "exception"}},
	}}
	if _, err := engine.InferConclusion("client-1", "Client", "normative_requirement", map[string]Kleene3VL{"exception": False}); err != nil {
		t.Fatal(err)
	}
	if _, err := engine.InferConclusion("client-1", "Client", "normative_requirement", nil); err == nil {
		t.Fatal("missing exception treated as false")
	}
	if _, err := engine.InferConclusion("client-1", "Client", "normative_requirement", map[string]Kleene3VL{"exception": True}); err == nil {
		t.Fatal("active exception ignored")
	}
	if _, err := engine.InferConclusion("server-1", "Server", "normative_requirement", map[string]Kleene3VL{"exception": False}); err == nil {
		t.Fatal("unrelated class accepted")
	}
	engine.rules = append(engine.rules, engine.rules[0])
	if _, err := engine.InferConclusion("client-1", "Client", "normative_requirement", map[string]Kleene3VL{"exception": False}); err == nil {
		t.Fatal("competing rules ignored")
	}
	if Unknown.Not() != Unknown || Unknown.And(False) != False || Unknown.Or(True) != True {
		t.Fatal("three-valued logic failed")
	}
}

func TestRebuttingAndUndercuttingDefeaters(t *testing.T) {
	engine := NewEngine(NewRelationHierarchy())
	rst := Rule{
		MajorPremise: Premise{Subject: "SynchronizedConnection", Predicate: "must_requirement",
			Object: "reset_connection", Modality: "MUST", Provenance: Provenance{DocumentID: "rfc9293"}},
		Defeaters: []Defeater{{ConditionPredicate: "rfc5961_and_seq_not_exact",
			FallbackAction: "MUST send_challenge_ack", IsRebutting: true, Source: "rfc9293#3.10.7.4"}},
	}
	engine.rules = []Rule{rst}
	conditions := map[string]Kleene3VL{"rfc5961_and_seq_not_exact": False}
	got, err := engine.InferConclusion("conn-1", "SynchronizedConnection", "normative_requirement", conditions)
	if err != nil || got != "conn-1 MUST reset_connection; source: rfc9293" {
		t.Fatalf("base rule: %q, %v", got, err)
	}
	conditions["rfc5961_and_seq_not_exact"] = True
	got, err = engine.InferConclusion("conn-1", "SynchronizedConnection", "normative_requirement", conditions)
	if err != nil || !strings.Contains(got, "send_challenge_ack") || !strings.Contains(got, "rfc9293#3.10.7.4") {
		t.Fatalf("rebutting defeater must yield the contrary conclusion with its source: %q, %v", got, err)
	}
	engine.rules[0].Defeaters[0].IsRebutting = false
	got, err = engine.InferConclusion("conn-1", "SynchronizedConnection", "normative_requirement", conditions)
	if err == nil || got != "" || strings.Contains(err.Error(), "challenge") {
		t.Fatalf("undercutting defeater must block without proposing an action: %q, %v", got, err)
	}
	engine.rules[0].Defeaters = append(engine.rules[0].Defeaters, Defeater{ConditionPredicate: "radio_silence"})
	conditions["radio_silence"] = Unknown
	conditions["rfc5961_and_seq_not_exact"] = False
	if _, err = engine.InferConclusion("conn-1", "SynchronizedConnection", "normative_requirement", conditions); err == nil {
		t.Fatal("unknown second exception ignored")
	}
}
```

</details>

Die hier implementierte Hierarchie erlaubt lediglich einen Elternknoten pro Prädikat und realisiert keine vollständigen Verbandsoperationen. Die übergebenen Klassen- und Umgebungszustände müssen aus unabhängig verifizierten Primärfakten stammen; ein unvalidierter String aus einem Benutzer-Prompt genügt dieser Anforderung nicht. Der erste Test simuliert gezielt fehlende Ausnahmestatus, inkompatible Klassen, Zyklen und konkurrierende Regeln: In allen Grenzfällen darf kein unbegründeter Schluss aufgrund zufälliger Iterationsreihenfolgen entstehen. Der zweite Test validiert die Differenzierung der Defeater-Klassen an der Norm aus RFC 9293. Bei aktivem Rebutting Defeater liefert das Programm das Challenge-ACK samt Quellenverweis auf Abschnitt 3.10.7.4; wird derselbe Defeater als Undercutting deklariert, verweigert das System deterministisch ohne Handlungsvorschlag. Ein unbestimmter Status einer zweiten Ausnahmebedingung blockiert den Schluss selbst dann, wenn die erste Ausnahme bekannt und inaktiv ist.

---

## 9. Zeitgenössische Werkzeuge und Knowledge Mining für Hierarchien und Ausnahmen

Die didaktischen Programme dieses Kapitels demonstrieren Inferenzmechanismen an kompakten Regelsätzen. In einer industriellen Wissensbasis mit Tausenden Prädikaten und Hunderten Spezifikationen stellen sich zwei ingenieurtechnische Kernfragen: Welche Standardwerkzeuge implementieren diese Mechanismen und wie lassen sich Kandidaten für Hierarchien und Ausnahmen automatisiert gewinnen, ohne sie vollständig manuell zu kuratieren? Die nachfolgende Tabelle stellt Mechanismen des Kapitels, einschlägige Werkzeuge, deren Mehrwert und Grenzen gegenüber.

| Mechanismus des Kapitels | Werkzeug oder Methode | Ingenieurtechnischer Ertrag | Systemgrenze |
|---|---|---|---|
| Prädikatenhierarchie und Anfragegeneralisierung | `rdfs:subPropertyOf` und Inferenzregel rdfs7 [[12]](#src-12); SPARQL-1.1-Eigenschaftspfade [[13]](#src-13) | Standardisiertes Format für Hierarchien und dynamische Generalisierung bei Abfragen | RDFS ist rein monoton und kann keine Ausnahmen abbilden |
| Revisionsauswahl und Fail-Closed-Verhalten | ASP und Solver clingo [[11]](#src-11) | Deklarative Regeln mit Negation as Failure; genau ein stabiles Modell für stratifizierte Programme | Registervollständigkeit bleibt eine externe Prämisse |
| Defeater und Vorrangrelationen | Strukturierte Argumentation ASPIC+ [[7]](#src-7) | Formale Semantik von Angriffen und Präferenzen zur Verifikation der Inferenzmaschine | Vorrangrelationen müssen vom Fachexperten definiert werden |
| Normen mit Ausnahmen als Programm | Domänenspezifische Sprache Catala [[14]](#src-14) | Ausnahmen als native Sprachkonstrukte; ausführbare Spezifikation direkt neben dem Normentext | Für Rechtsnormen entworfen; technische Standards erfordern eigene Domänenmodelle |
| Hierarchiekandidaten aus Rohtext | Hearst-Muster (*Hearst Patterns*) [[15]](#src-15) | Automatisierte Extraktion von Hyponym-Hypernym-Paaren aus Dokumentenkorpora | Musterbasierte Kandidaten sind fehlerträchtig; Expertenprüfung zwingend |
| Regel- und Eigenschaftskandidaten aus Faktenbasen | Assoziationsregel-Mining AMIE [[16]](#src-16) | Induktion von Regeln der Form „wenn $`r_1(x,y)`$, dann $`r_2(x,y)`$“ mit Konfidenzwerten | Statistische Signifikanz begründet keine normative Verbindlichkeit |

Die erste Zeile der Tabelle schlägt die Brücke zu Abschnitt 3: Die RDFS-Inferenzregel rdfs7 besagt, dass aus dem Bestehen eines Fakts $x\,a\,y$ und der Relation $a \text{ rdfs:subPropertyOf } b$ zwingend das Faktum $x\,b\,y$ folgt. Dies entspricht exakt der Bedingung $\text{Matches}$ aus Abschnitt 3.1 und dem Verbot der Rückwärtsspezialisierung: Aus einem Faktum über die Ober-Eigenschaft lässt sich niemals ein Faktum über die Unter-Eigenschaft ableiten. Ein SPARQL-Eigenschaftspfad der Form `?p rdfs:subPropertyOf* :normative_requirement` expandiert zur Abfragezeit alle Sub-Properties dynamisch, sodass abgeleitete Fakten nicht statisch materialisiert werden müssen. Die fundamentale Grenze von RDFS liegt in seiner strikten Monotonie: Ein neu hinzugefügtes Faktum kann einen bestehenden Schluss niemals ungültig machen. Daher empfiehlt es sich, die statische Prädikatenhierarchie in RDF/RDFS zu verwalten, während Ausnahmen, Geltungsprüfung und Revisionsentscheidungen in ASP oder einer dedizierten Inferenzmaschine abgebildet werden — exakt analog zur Trennung zwischen Abschnitt 3 und den Abschnitten 4 und 5 dieses Kapitels.

Die beiden letzten Zeilen der Tabelle entstammen dem Text- und Datamining. Marti Hearst wies nach, dass lexico-syntaktische Muster wie „X, wie etwa Y und Z“ verlässliche Begriffshierarchien ohne überwachtes Training erschließen. In Normen treten solche Wendungen gehäuft auf, beispielsweise „Transportprotokolle wie etwa TCP und UDP“. AMIE wiederum operiert nicht auf Rohtext, sondern auf Faktenbasen: Der Algorithmus induziert häufig auftretende Assoziationsregeln unter der Annahme partieller Vollständigkeit (*Partial Completeness Assumption*, PCA). Nach der PCA gilt: Kennt die Wissensbasis zu einem Subjekt mindestens ein Objekt für eine Relation, wird sie für dieses Subjekt-Relations-Paar als vollständig angenommen; über Subjekte ohne jeden Eintrag wird keine Annahme getroffen. Dies markiert einen wohlüberlegten Mittelweg zwischen geschlossener und offener Welt (Abschnitt 4). In einem Expertensystem darf AMIE ausschließlich zur Generierung von Hypothesen genutzt werden: Eine Regel wie „wenn `must_encrypt_channel(x, y)`, dann `security_requirement(x, y)`“ mit hoher Konfidenz wird als Kandidat für eine Kante im Graphen vorgeschlagen. Die Aufnahme in die autoritative Hierarchie erfolgt erst nach Begutachtung durch einen Wissensingenieur und automatischer Zyklusprüfung, wie sie in der Methode `AddRelation` aus Abschnitt 8 implementiert ist.

Fertige Softwarewerkzeuge decken mithin einzelne Teilaspekte ab, kein einzelnes Werkzeug integriert jedoch Prädikatenhierarchien, Ausnahmen, Revisionswahl und Fail-Closed-Verhalten autonom. Mining-Verfahren reduzieren den manuellen Konstruktionsaufwand, liefern jedoch ausschließlich Kandidaten, die dasselbe strenge Zulassungsgateway durchlaufen müssen wie die Fakten aus [Kapitel 15](ch15-knowledge-extraction-and-kb-construction.md).

---

## 10. Pro et Contra — Virtueller Dialog

Zu diesem virtuellen intellektuellen Fachdialog sind die Begründer der Formalen Begriffsanalyse (Formal Concept Analysis, FCA) und der Theorie der Relationsverbände eingeladen: Bernhard Ganter (Professor an der Technischen Universität Dresden, Koautor der grundlegenden Monografie *Formal Concept Analysis: Mathematical Foundations*) und Rudolf Wille (verstorbener Professor an der Technischen Universität Darmstadt, Schöpfer der Begriffsverbandstheorie und der philosophischen Grundlagen der algebraischen Wissensmodellierung). Ihr mathematischer Apparat der Halbordnungen, Begriffsverbände $\mathfrak{B}(G, M, I)$ und syllogistischen Hierarchien bietet ein maßgebliches Fundament zur kritischen Überprüfung des vom Autor vorgeschlagenen Modells der deterministischen syllogistischen Inferenz und Ausnahmebehandlung in normativen Wissensbasen.

*Disclaimer:* Diese Diskussion ist eine heuristische virtuelle Simulation im wissenschaftlich-didaktischen Kontext der Monografie. Die Positionen und Argumente der Diskutanten wurden auf der Grundlage ihrer wissenschaftlichen Publikationen und Monografien synthetisiert; sie stellen keine wörtlichen Zitate dar und dienen der kritischen Rationalisierung (Popperschen Falsifikation) der ingenieurtechnischen Hypothesen des Autors.

**Runde 1: Bernhard Ganter hinterfragt das Risiko einer kombinatorischen Verbandsexplosion bei dynamischen Ausnahmen**

- **Bernhard Ganter:** Herr Fedchyk, bei der Ausarbeitung der Formalen Begriffsanalyse (FCA) haben wir mathematisch bewiesen, dass für jeden formalen Kontext $(G, M, I)$ die Menge aller formalen Begriffe einen vollständigen Verband $\mathfrak{B}(G, M, I)$ bildet. In ingenieurtechnischen normativen Wissensbasen führt jedoch der Versuch, einen vollständigen Verband über Tausende vernetzter Prädikate, Objekte und Attribute zu spannen, unweigerlich zu einer kombinatorischen Verbandsexplosion (Lattice Explosion), bei der die Anzahl intermediärer Begriffe exponentiell zur Vokabulargröße wächst. In Ihrem Modell koppeln Sie eine partielle Prädikatenordnung (Relation Lattice) mit Kleenes dreiwertiger Logik und Pollocks nichtmonotonen Defeatern (Rebutting und Undermining Defeater). Zerstört eine solche Überlagerung dynamischer Ausnahmen auf einem statischen Begriffsverband nicht die algebraische Verbandsstruktur selbst, wodurch Suchoperationen zeitlich nichtdeterministisch werden?

- **Mykola Fedchyk:** Herr Ganter, Ihre Monografie ist die Bibel der Verbandstheorie, und Sie verweisen vollkommen zu Recht auf die Gefahren einer naiven Materialisierung des gesamten Begriffsverbands. Um eine kombinatorische Explosion in unserer Architektur auszuschließen, wird der Relationsverband $\mathcal{L} = (P, \sqsubseteq)$ niemals als vollständiger extensionaler Verband über alle Domänenobjekte materialisiert. Wir beschränken uns auf eine *intensionale Halbordnung von Prädikaten*, bei der die Spezialisierungsprüfung `Matches(p_child, p_parent)` über den transitiven Abschluss elementarer `subPropertyOf*`-Kanten dynamisch zur Abfragezeit erfolgt. Nichtmonotone Ausnahmen verändern die Verbandstopologie keineswegs: Wie in unseren TCP-RST-Szenarien und der ASP-Formalisierung gezeigt, dient die Prädikatenhierarchie ausschließlich der deduktiven Generalisierung, während Defeater auf der Filterebene des Ableitungsbaums wirken. Dies garantiert lineare Prüfkomplexität und wahrt strikten Determinismus der Worst-Case Execution Time (WCET).

**Runde 2: Rudolf Wille hinterfragt die Reduktion des philosophischen Begriffskontexts auf diskrete Fail-Closed-Verweigerungen**

- **Rudolf Wille:** Herr Fedchyk, der philosophische Kern der Formalen Begriffsanalyse besteht in der Erfassung des menschlichen Begreifens der Wirklichkeit durch die dialektische Einheit von Begriffsumfang (Extension) und Begriffsinhalt (Intension): Wissen entfaltet sich in einem offenen Kontext, in dem jedes neue Objekt oder Merkmal das Beziehungsgefüge auf natürliche Weise transformiert. In Ihrem Modell führen Sie eine rigide Invariante der typisierten Verweigerung (Closed Refusal) ein: Ist ein Faktum nicht registriert oder die Vollständigkeit des Dokumentenregisters unbestätigt, verweigert das System die Inferenz vollständig mit dem Urteil `DISP_REFUSE`. Verwandelt ein solch überzogener operativer Rigorismus das lebendige philosophische Begriffsgefüge nicht in einen sterilen syntaktischen Filter, der das Expertensystem jeder heuristischen Wissenserweiterung in der realen offenen Welt beraubt?

- **Mykola Fedchyk:** Herr Wille, Ihr philosophisches Konzept der Begriffsverbände als Medium menschlicher Verständigung inspiriert unsere beratenden Teilsysteme, in denen Ingenieure nach Analogien forschen. Geht es jedoch um die Steuerung kritischer Infrastrukturen — Eisenbahnsignaltechnik, Bahnstromnetze oder Hochleistungs-Netzwerkstacks nach RFC 5961 —, dürfen Begriffe unter keinen Umständen verschwimmen. In der missionskritischen Technik ist die stillschweigende Erweiterung des Kontexts ohne normative Registrierung gleichbedeutend mit einer Katastrophe: Die Anwendung einer veralteten TCP-Normrevision oder die Annahme, ein RST-Segment sei legitim, bloß weil kein explizites Verbot vorliegt, öffnet Tür und Tor für Blind-Reset-Angriffe. Eine geschlossene typisierte Verweigerung in der vorgeschlagenen Architektur ist keine Abkehr von der Erkenntnis, sondern ein Schutzschild funktionaler Sicherheit: Das System deklariert transparent die Grenzen des Bewiesenen und überlässt jede ontologische Kontexterweiterung ausschließlich dem autorisierten Fachexperten.

## Fazit

Eine Prädikatenhierarchie gestattet die Generalisierung von Anfragen, beweist jedoch weder die Vollständigkeit des Wissens noch konstituiert sie automatisch einen mathematischen Verband. Die gültige Revision einer Norm, die bestätigte Objektklasse und verifizierte Ausnahmestatus bilden drei unabhängige Prämissen einer fundierten Inferenz. Das deklarative ASP-Modell zeigte anhand von neun Grenzszenarien, dass die Reaktion auf ein TCP-RST-Segment deterministisch vom Verbindungszustand, der Unterstützung des Schutzes nach RFC 5961 und der Sequenznummer abhängt. Ein fehlendes Ablösungsfaktum, ein nicht registriertes Dokument oder eine unbestätigte Registervollständigkeit lösen ein Fail-Closed-Verhalten aus, statt stillschweigend eine veraltete Norm anzuwenden. Die Go-Referenztests verifizierten die Verweigerung bei unbekannten Ausnahmen, Zyklen, fremden Klassen und konkurrierenden Regeln sowie den fundamentalen Unterschied zwischen Rebutting und Undercutting Defeatern.

Die didaktische Referenzimplementierung beschränkt sich auf einen einzelnen Inferenzschritt und prüft weder die Vollständigkeit des Quellenkorpus noch mehrstufige Ausführungspläne. Der Inferenzpfad muss diese Grenzen transparent dokumentieren, statt sie hinter der Bezeichnung „Beweis“ zu verschleiern. Zeitgenössische Werkzeuge von RDFS über SPARQL bis hin zu ASP und ASPIC+ decken isolierte Teilmechanismen ab; automatisches Knowledge Mining aus Texten und Faktenbasen liefert lediglich Kandidaten zur Begutachtung durch den Wissensingenieur. Für den industriellen Produktiveinsatz sind separate Verifikationsschritte für Faktenintegrität, Konfliktdetektion, Terminierung und Reproduzierbarkeit unverzichtbar.

---

## Fragen zur Selbstprüfung

1. Warum stellt einstufiges Vektor-Retrieval nicht fest, dass ein Standard $A$ durch einen Standard $B$ obsoletiert wurde, selbst wenn beide Dokumente im Ranking zurückgegeben werden?
2. Worin besteht der fundamentale Unterschied zwischen gerichteter Subsumtion bei generalisierten Anfragen und dem strikten Verbot der Rückwärtsspezialisierung bei Verboten?
3. Wie verhindert die Wahrheitstafel der starken Kleene-Logik die unzulässige Anwendung der Closed-World Assumption bei unvollständigen Sensordaten?
4. Erläutern Sie die Differenzierung zwischen Rebutting Defeater und Undercutting Defeater nach Pollock. Wie manifestiert sich dieser Unterschied in der Systemantwort und welcher Angriff ist im ASPIC+-Framework unabhängig von Präferenzrelationen erfolgreich?
5. Auf welche Weise ermöglicht die Verknüpfung von Obersatz, Untersatz und Ausnahmezweigen in einem mehrstufigen Syllogismus die Synthese einer korrekten Antwort auf komplexe hypothetische Anfragen?
6. Warum erzeugt ein fehlendes Obsoletierungsfaktum im ASP-Modell einen Konflikt anstelle eines stillschweigenden Rückgriffs auf die Altversion, und welche Sicherheitsfunktion erfüllt das Faktum `complete(tcp_rst)`?
7. Inwiefern entspricht die Inferenzregel rdfs7 in RDFS der Subsumtionsprüfung `Subsumes` und welche logischen Grenzen weist RDFS bei Ausnahmen auf?
8. Wie lassen sich durch AMIE oder Hearst-Muster entdeckte Regeln im Knowledge Engineering einbinden, ohne statistische Korrelation fälschlich mit normativer Verbindlichkeit gleichzusetzen?
9. Wie eliminieren Dungs charakteristischer Operator $`\Gamma_{\mathcal{AF}}`$ und Tarjans Algorithmus für stark zusammenhängende Komponenten (SCC) Nichtdeterminismus und Endlosschleifen bei symmetrischen Angriffen gleicher Stärke?
10. Auf welche Weise schützt der Operator transitiver Falsifikation in GSN-Sicherheitszielbäumen (ISO 26262 ASIL D) vor dem Defekt des „stillen Übergehens“ (*Silent Failure Pass*), und unter welchen Bedingungen erfolgt eine epistemische Rehabilitation einer zuvor widerlegten Regel?

---

## Glossar

| Deutscher Begriff | Englisches Äquivalent | Kurzbeschreibung |
|---|---|---|
| Syllogismus | Syllogism | Deduktiver Schluss, bei dem aus zwei Prämissen eine Konklusion abgeleitet wird |
| Deduktion | Deduction | Inferenzmethode vom allgemeinen Prinzip auf das konkrete Einzelurteil |
| Prädikatenhierarchie | Predicate hierarchy | Azyklische Spezialisierungsstruktur von Prädikaten, die eine Halbordnung definiert |
| Verband | Lattice | Partiell geordnete Menge mit eindeutigem Supremum und Infimum für jedes Paar |
| Subsumtion | Subsumption | Unterordnung eines Begriffs oder Prädikats unter einen allgemeineren Begriff |
| Defeater | Defeater | Ausnahmebedingung oder Einwand, der einer Prämisse oder Regel die Beweiskraft entzieht |
| Rebutting Defeater | Rebutting defeater | Widerlegender Einwand, der eine der Regelkonklusion entgegengesetzte Aussage ableitet |
| Undercutting Defeater | Undercutting defeater | Untergrabender Einwand, der den Verknüpfungszusammenhang zwischen Prämisse und Konklusion annulliert |
| Strukturierte Argumentation | Structured argumentation | Formales Modell von Argumenten, Angriffen und Präferenzen zur Ermittlung akzeptabler Konklusionen |
| Answer Set Programming | Answer set programming | Deklaratives Programmierparadigma, bei dem Lösungen stabilen Modellen entsprechen |
| Stratifiziertes Programm | Stratified program | Logisches Programm ohne Rekursion über Negation mit genau einem stabilen Modell |
| Negation als Fehlschlag | Negation as failure | Schluss auf „nicht P“ allein aus der Unbeweisbarkeit von P |
| Annahme partieller Vollständigkeit | Partial completeness assumption | Wissensbasis gilt als vollständig für ein Subjekt-Relations-Paar, wenn mindestens ein Faktum existiert |
| Hearst-Muster | Hearst patterns | Lexico-syntaktische Muster zur Extraktion von Hyponym-Hypernym-Paaren aus Texten |
| Starke Kleene-Logik | Strong Kleene 3VL | Dreiwertiges Logiksystem mit den Wahrheitswerten Wahr, Falsch und Unbekannt |
| Closed-World Assumption | Closed-World Assumption | Epistemische Annahme, dass jede nicht bewiesene Aussage falsch ist |
| Inferenzpfad | Proof DAG / Trace | Gerichteter azyklischer Graph deduktiver Schritte mit Verweisen auf Primärquellen |
| Deontische Modalität | Deontic modality | Normativer Charakter eines Urteils (verpflichtend, verboten, freigestellt) |
| Dokumentenfamilie | Document family | Gesamtheit aller durch Revisionen und Ausnahmen verbundenen Fassungen einer Quelle; Sharding-Einheit |
| Schweigen eines Shards | Shard silence | Ausbleiben der Rückmeldung eines Shards; darf nicht mit der Nichtexistenz eines Fakts gleichgesetzt werden |
| Fundierte Erweiterung (Grounded Extension) | Grounded extension | Eindeutiger kleinster Fixpunkt des charakteristischen Dung-Operators, der skeptisch akzeptable Argumente enthält |
| Caminada-Labeling | Caminada labeling | Dreiwertige Kennzeichnungsfunktion zur Zuweisung disjunkter Zustände IN, OUT oder UNDEC an Argumente |
| Epistemische Argumentrehabilitation | Argument reinstatement | Mathematische Wiederherstellung des IN-Status eines attackierten Arguments bei erfolgreicher Widerlegung seiner Angreifer |
| Transitive Falsifikation | Transitive falsification | Bottom-up-Propagation des Falsifikationsstatus von hardwarenahen Telemetrieblättern bis zu den Sicherheits-Wurzelzielen |
| Stark zusammenhängende Komponente | Strongly connected component | Maximaler gerichteter Teilgraph, in dem zwischen beliebigen zwei Knoten ein gerichteter Pfad existiert |

---

## Abkürzungen

| Abkürzung | Vollständige Bezeichnung | Bedeutung |
|---|---|---|
| 3VL | Three-Valued Logic | Dreiwertige Logik |
| ACK | Acknowledgment | Bestätigungsquittung in Netzwerkprotokollen |
| AMIE | Association Rule Mining under Incomplete Evidence | Mining von Assoziationsregeln aus unvollständigen Faktenbasen |
| ASIL | Automotive Safety Integrity Level | Sicherheitsanforderungsstufe für Straßenfahrzeuge nach ISO 26262 |
| ASP | Answer Set Programming | Deklaratives Paradigma logischer Programmierung in stabilen Modellen |
| ASPIC+ | Argumentation Service Platform with Integrated Components | Formales Framework für strukturierte Argumentation nach Modgil und Prakken |
| CWA | Closed-World Assumption | Annahme einer geschlossenen Welt |
| DAG | Directed Acyclic Graph | Gerichteter azyklischer Graph |
| GE | Grounded Extension | Fundierte Erweiterung in Dungs abstrakter Argumentation |
| GSN | Goal Structuring Notation | Grafische Modellierungssprache für Sicherheitsnachweise |
| PCA | Partial Completeness Assumption | Annahme partieller Vollständigkeit |
| RDF | Resource Description Framework | Datenmodell aus Tripeln (Subjekt, Prädikat, Objekt) |
| RDFS | RDF Schema | Vokabular zur Definition von Klassen und Eigenschaften für RDF mit Inferenzregeln |
| RFC | Request for Comments | Serie technischer Spezifikationen und Standards der Internet Engineering Task Force |
| RST | Reset | Steuerbit zum abrupten Zurücksetzen einer TCP-Verbindung |
| SCC | Strongly Connected Component | Stark zusammenhängende Komponente eines gerichteten Graphen |
| SMTP | Simple Mail Transfer Protocol | Standardprotokoll zur Übertragung elektronischer Post |
| SPARQL | SPARQL Protocol and RDF Query Language | Deklarative Abfragesprache für RDF-Graphen |
| TCP | Transmission Control Protocol | Verbindungsorientiertes Transportprotokoll des Internets |

---

## Literaturhinweise

1. <a id="src-1"></a>Scott Bradner. [*RFC 2119: Key words for use in RFCs to Indicate Requirement Levels*](https://www.rfc-editor.org/rfc/rfc2119). IETF, 1997.
2. <a id="src-2"></a>Barry Leiba. [*RFC 8174: Ambiguity of Uppercase vs Lowercase in RFC 2119 Key Words*](https://www.rfc-editor.org/rfc/rfc8174). IETF, 2017.
3. <a id="src-3"></a>Robin Smith. [*Aristotle's Logic*](https://plato.stanford.edu/entries/aristotle-logic/). *Stanford Encyclopedia of Philosophy*. Übersichtsdarstellung, keine Übersetzung der „Ersten Analytiken“.
4. <a id="src-4"></a>John C. Klensin. [*RFC 5321: Simple Mail Transfer Protocol*](https://www.rfc-editor.org/rfc/rfc5321). IETF, 2008.
5. <a id="src-5"></a>Stephen Kleene. [*Introduction to Metamathematics*](https://openlibrary.org/works/OL5959470W). D. Van Nostrand Co., Inc., New York, 1952.
6. <a id="src-6"></a>John L. Pollock. [*Cognitive Carpentry: A Blueprint for How to Build a Person*](https://mitpress.mit.edu/9780262661133/). MIT Press, Cambridge, MA, 1995.
7. <a id="src-7"></a>Sanjay Modgil, Henry Prakken. [*The ASPIC+ Framework for Structured Argumentation: A Tutorial*](https://doi.org/10.1080/19462166.2013.869766). *Argument & Computation*, 5(1), 31–62, 2014.
8. <a id="src-8"></a>Jon Postel. [*RFC 793: Transmission Control Protocol*](https://www.rfc-editor.org/rfc/rfc793). IETF, 1981.
9. <a id="src-9"></a>Anantha Ramaiah, Randall Stewart, Michael Dalal. [*RFC 5961: Improving TCP's Robustness to Blind In-Window Attacks*](https://www.rfc-editor.org/rfc/rfc5961). IETF, 2010.
10. <a id="src-10"></a>Wesley Eddy. [*RFC 9293: Transmission Control Protocol (TCP)*](https://www.rfc-editor.org/rfc/rfc9293). IETF, 2022.
11. <a id="src-11"></a>Martin Gebser, Roland Kaminski, Benjamin Kaufmann, Torsten Schaub. [*Multi-shot ASP Solving with clingo*](https://doi.org/10.1017/S1471068418000054). *Theory and Practice of Logic Programming*, 19(1), 27–82, 2019.
12. <a id="src-12"></a>Patrick J. Hayes, Peter F. Patel-Schneider (Hrsg.). [*RDF 1.1 Semantics*](https://www.w3.org/TR/rdf11-mt/). W3C Recommendation, 2014.
13. <a id="src-13"></a>Steve Harris, Andy Seaborne (Hrsg.). [*SPARQL 1.1 Query Language*](https://www.w3.org/TR/sparql11-query/). W3C Recommendation, 2013.
14. <a id="src-14"></a>Denis Merigoux, Nicolas Chataing, Jonathan Protzenko. [*Catala: A Programming Language for the Law*](https://doi.org/10.1145/3473582). *Proceedings of the ACM on Programming Languages*, 5(ICFP), 2021.
15. <a id="src-15"></a>Marti A. Hearst. [*Automatic Acquisition of Hyponyms from Large Text Corpora*](https://aclanthology.org/C92-2082/). COLING 1992, Volume 2.
16. <a id="src-16"></a>Luis Galárraga, Christina Teflioudi, Katja Hose, Fabian M. Suchanek. [*AMIE: Association Rule Mining under Incomplete Evidence in Ontological Knowledge Bases*](https://doi.org/10.1145/2488388.2488425). Proceedings of WWW 2013, 413–422.
17. <a id="src-17"></a>Phan Minh Dung. [*On the Acceptability of Arguments and its Fundamental Role in Nonmonotonic Reasoning, Logic Programming and n-Person Games*](https://doi.org/10.1016/0004-3702(94)00041-X). *Artificial Intelligence*, 77(2), 321–357, 1995.
18. <a id="src-18"></a>ISO. [*ISO 26262:2018 Road vehicles — Functional safety*](https://www.iso.org/standard/68383.html). International Organization for Standardization, 2018.
19. <a id="src-19"></a>Robert E. Tarjan. [*Depth-First Search and Linear Graph Algorithms*](https://doi.org/10.1137/0201010). *SIAM Journal on Computing*, 1(2), 146–160, 1972.

---

[← Kapitel 19](ch19-from-question-to-evidence.md) | [Inhaltsverzeichnis](README.md) | [Teil IV](part-04-architecture-and-inference.md) | [Kapitel 20 →](ch20-explanation-engine.md)
