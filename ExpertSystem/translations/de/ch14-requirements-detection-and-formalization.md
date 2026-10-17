# Kapitel 14. Anforderungsextraktion und Modalitäten: Vom normativen Text zu Invarianten

> **Buch:** [Architektur evidenzbasierter Expertensysteme](README.md) · [Teil III: Eingangsdatenbewertung, Wissensakquisition und linguistische Analyse](part-03-knowledge-engineering-nlp.md)  
> **Vorheriges Kapitel:** [Kapitel 12. Linguistische Analyse und lokale Modelle: Erhalt von Semantik und Herkunftsnachweis](ch12-linguistic-analysis-and-local-models.md)  
> **Nächstes Kapitel:** [Kapitel 15. Wissensextraktion und Aufbau der Wissensbasis: Fakten, Grammatiken und Automaten](ch15-knowledge-extraction-and-kb-construction.md)  
> **Inhaltsverzeichnis:** [README.md](README.md)  
> **Autor:** [Mykola Fedchyk](about-the-author.md)  
> **Niveau:** Fortgeschritten: Systemingenieure, Architekten, Entwickler, V&V-Spezialisten (Verifikation und Validierung)  
> **Lernziele:** Normative Sätze anhand der Konventionen des Quelldokuments von deskriptiven unterscheiden; Anforderungen auf EARS-Schablonen abbilden; Schablonen in Prädikatenlogik, temporale Logik und SMT-LIB-Formeln übersetzen; Widersprüche zwischen Anforderungen mittels SMT-Solver aufdecken; einen prüffähigen Auditbericht zur Anforderungsqualität mit lückenloser Herkunftsnachweis-Bindung erstellen.

## Abstract

Dieses Kapitel untersucht das Problem der automatisierten Erkennung von Anforderungen und deontischen Modalitäten in technischen Spezifikationen, Industriestandards (ISO 26262, DO-178C) und RFCs zur Synthese formaler Invarianten für evidenzbasierte Expertensysteme. Es analysiert die syntaktischen Fallstricke natürlicher Sprache, entwickelt einen Linter für semantische Mehrdeutigkeit auf Basis eines Lexikons toxischer Deskriptoren und der lexikografischen Vorzugsordnung deontischer Präferenzen sowie Schutzalgorithmen gegen die Polaritätsinversion modaler Quantoren (`MUST NOT`) durch zusammengesetzte Atome, eine Spans erhaltende Tokenisierung (`GroundedToken`) und eine UTF-8-Analyse von Interpunktionsgrenzen. Behandelt wird die neuro-symbolische Gültigkeitsschranke zwischen Relevanz- und Validitätszone mittels grammatikgesteuerter Decodierung (GBNF) und standardisierter EARS-Schablonen (Easy Approach to Requirements Syntax). Untersucht wird die Übersetzung formalisierter Anforderungen in Prädikatenlogik erster Stufe, lineare temporale Logik (LTL) und das SMT-LIB-Format zur automatisierten Konsistenzprüfung durch den SMT-Solver Z3, wodurch der deterministische Regelkern des Expertensystems formiert wird. Abschließend werden eine deterministische Pipeline-Implementierung in Go sowie eine durchgängige Engineering-Fallstudie zum Audit von Anforderungen an ein Batteriemanagementsystem (BMS) dargelegt.

Die Entwicklung eines komplexen Systems – vom Steuergerät einer Traktionsbatterie eines Elektrofahrzeugs bis hin zur bordseitigen Flugzeugsoftware – beginnt nicht mit Code, sondern mit Anforderungen. Anforderungen entstammen drei wesentlichen Quellen: internationalen und branchenspezifischen Normen, dem Lastenheft des Auftraggebers sowie der Komponentendokumentation der Hardwarebasis, also Datenblättern von Mikrocontrollern und Listen bekannter Hardware-Fehler (Errata). Spezifikationen liegen häufig als PDF-, DOCX- oder ReqIF-Dateien (*Requirements Interchange Format*) vor [[1]](#src-1). Standards für funktionale Sicherheit wie die ISO 26262 für Straßenfahrzeuge [[2]](#src-2) oder die DO-178C für flugzeuggestützte Software [[3]](#src-3) fordern eine lückenlose Rückverfolgbarkeit (Traceability) jeder einzelnen Anforderung zu Architekturentscheidungen, Quellcode und Verifikationstests.

Alle nachgelagerten Phasen – vom Systementwurf bis zur Zulassung – verlangen mathematische Präzision; die Primärquellen sind jedoch in natürlicher Sprache verfasst: geprägt von vagen Formulierungen, Passivkonstruktionen und Ausnahmeregelungen innerhalb desselben Satzes. Vertraut ein Entwicklungsteam die Analyse derartiger Dokumente einer Vektorähnlichkeitssuche oder einem großen Sprachmodell (*Large Language Model*, LLM) an, treten systematische Fehlerbilder auf: Das Modell verliert den Allquantor („für alle“), übergeht Zeitgrenzen, verwechselt Verbote mit Empfehlungen und kann die vollständige Erfassung aller normativen Klauseln nicht garantieren. Für ein Zertifizierungsaudit ist die Auskunft „die Anforderung ist höchstwahrscheinlich erfüllt“ wertlos. Warum textuelle Ähnlichkeit keineswegs semantischem Verständnis entspricht, wurde in [Kapitel 13](ch13-language-variability-vs-determinism.md) dargelegt.

Daraus resultiert die zentrale Fragestellung dieses Kapitels: **Wie lassen sich Anforderungen in normativen Dokumenten automatisiert identifizieren und in verifizierbare formale Invarianten überführen, ohne den Herkunftsnachweis zur Quelle zu verlieren?** Die Grundthese lautet: Eine Anforderung wird nicht allein durch das Signalwort SHALL konstituiert, sondern durch die Modalität gemäß den Konventionen des konkreten Dokuments, die syntaktische Satzstruktur und operationell prüfbare Parameter. Das Expertensystem behandelt eine Spezifikation wie unkompilierten Quelltext: Es bestimmt deterministisch die Modalität, normalisiert den Satz auf eine Schablone, übersetzt diese in logische Formeln und verifiziert deren Widerspruchsfreiheit; ein Sprachmodell liefert hierbei lediglich Hypothesenkandidaten, die durch deterministischen Host-Code validiert werden.

```mermaid
flowchart TD
    accTitle: Pipeline zur Anforderungserkennung und Invariantenkompilierung
    accDescr: Normativer Text durchläuft Satzextraktion, Modalitätsbestimmung, syntaktisches Parsing, EARS-Typisierung, Erkennung von Spezifikationsmängeln und Formalisierung bis zum Auditbericht.

    Doc["<b>Eingangsdokumente</b><br/>PDF, DOCX, ReqIF, RFC"] --> Extr["<b>Textextraktion</b><br/>Sätze mit Byte-Offset-Bindung an Quelle"]
    Extr --> Modal["<b>Modalitätsbestimmung</b><br/>nach RFC 2119 oder ISO/IEC-Direktiven"]
    Modal --> Syntax["<b>Syntaktisches Parsing</b><br/>Subjekt, Aktion, Trigger, Zustand, Bedingung"]
    Syntax --> EARS["<b>EARS-Typisierung</b><br/>ubiquitär, ereignisgesteuert, zustandsbasiert, Fehlerreaktion"]
    EARS --> Smell["<b>Anforderungsmängel (Smells)</b><br/>Vagheit, Nicht-Atomarität, fehlender Zweig"]
    Smell --> Formal["<b>Formalisierung</b><br/>Prädikatenlogik, temporale Logik, SMT-LIB"]
    Formal --> Audit["<b>Auditbericht</b><br/>Bericht, Metriken, Export in Wissensgraphen"]

    classDef src fill:#fce4ec,stroke:#c2185b,stroke-width:2px,color:#880e4f;
    classDef proc fill:#e1f5fe,stroke:#0288d1,stroke-width:2px,color:#01579b;
    classDef check fill:#fff3e0,stroke:#f57c00,stroke-width:2px,color:#e65100;
    classDef ok fill:#e8f5e9,stroke:#388e3c,stroke-width:2px,color:#1b5e20;

    class Doc src;
    class Extr,Modal,Syntax,EARS proc;
    class Smell,Formal check;
    class Audit ok;
```

Dieses Schema strukturiert den Ablauf des Kapitels. Die ersten beiden Schritte entscheiden, ob ein Satz überhaupt normativen Charakter besitzt. Die nachfolgenden zwei Schritte überführen die Anforderung in eine kanonische Form. Die abschließenden drei Schritte prüfen die Spezifikationsqualität, übersetzen die Anforderungen in formale Logik und erzeugen ein prüffähiges Gutachten für Fachexperten.

## 1. Deontische Modalität und Abgrenzungskriterien für normative Anforderungen

Im Lebenszyklus der funktionalen Sicherheit eingebetteter und cyber-physischer Systeme (gemäß ISO 26262-8 und DO-178C Abschnitt 5) muss die Wissensbasis eines Expertensystems zwingend zwischen verbindlichen normativen Vorgaben und rein informativen Ingenieurkommentaren unterscheiden. Werden ungefiltert sämtliche Sätze einer technischen Spezifikation in die formale Regelbasis übernommen, droht epistemische Kontamination: Die Inferenzmaschine interpretiert unverbindliche Empfehlungen oder illustrative Beispiele fälschlich als fundamentale Sicherheitsinvarianten. Dies führt zur kombinatorischen Explosion des Zustandsraums, zu Fehlalarmen oder zur künstlichen Inkonsistenz des Gesamtsystems. Das Fundament der ingenieurmäßigen Anforderungsanalyse bildet daher die mathematische Erfassung der **deontischen Modalität** – der Modallogik von Gebot, Verbot und Erlaubnis.

Ein technisches Dokument besteht keineswegs ausschließlich aus Anforderungen. Es enthält Begründungen von Entwurfsentscheidungen (*rationale*), Anmerkungen, Anwendungsbeispiele und Hinweise zur Benutzerfreundlichkeit. Die Modalität eines Satzes definiert, ob er verpflichtet, verbietet, empfiehlt, erlaubt oder lediglich den Systemkontext beschreibt.

Die Modalität wird maßgeblich durch die redaktionellen Richtlinien des Dokuments bestimmt, wobei unterschiedliche Dokumentenfamilien divergierende Konventionen anwenden. RFC 2119 definiert die Schlüsselwörter für Internetspezifikationen: MUST, SHALL und REQUIRED bezeichnen verbindliche Anforderungen, MUST NOT und SHALL NOT Verbote, SHOULD und RECOMMENDED Empfehlungen, MAY und OPTIONAL Erlaubnisse [[4]](#src-4). RFC 8174 präzisiert, dass diese Wörter nur dann normative Bedeutung besitzen, wenn sie vollständig in Großbuchstaben gesetzt sind [[5]](#src-5). ISO- und IEC-Normen folgen anderen Vorgaben: Die ISO/IEC-Direktiven, Teil 2, legen Verbformen in Kleinbuchstaben fest, wobei *shall* eine Anforderung, *should* eine Empfehlung, *may* eine Erlaubnis, *can* eine Möglichkeit oder Fähigkeit und *must* eine externe Randbedingung bezeichnet, die keine Anforderung des Dokuments selbst darstellt [[6]](#src-6).

Der Klassifikator muss daher stets die Dokumentenkonvention sowie die funktionale Rolle des jeweiligen Abschnitts berücksichtigen. Zwar schreibt RFC 8174 Großschreibung für Schlüsselwörter vor, doch bedeutet das Fehlen eines solchen Markers nicht zwingend, dass der Satz frei von Anforderungen ist. Eine Verpflichtung kann in natürlicher Sprache formuliert oder durch einen Verweis auf einen anderen Normenteil ausgedrückt sein. Daher stellt das Auffinden eines Signalworts zunächst nur einen Modalitätskandidaten dar; Zitate, Beispiele und informative Notizen bedürfen einer gesonderten Verifikation. Eine unbekannte Konvention darf keinesfalls stillschweigend als „rein deskriptiver Satz“ interpretiert werden.

| Normative Verbindlichkeit | Signalwörter nach RFC 2119 und RFC 8174 | Signalwörter nach ISO/IEC-Direktiven | Aktion des Expertensystems |
|---|---|---|---|
| **Verbindliche Anforderung** | MUST, SHALL, REQUIRED | shall | Erzeugt eine verbindliche Invariante; erfordert rückverfolgbare Verifikation |
| **Verbot** | MUST NOT, SHALL NOT | shall not | Erzeugt eine Sicherheitsinvariante „Zustand tritt niemals ein“; plant Negativtests |
| **Empfehlung** | SHOULD, RECOMMENDED | should | Registriert ein Soft-Constraint; Abweichungen erfordern dokumentierte Begründung |
| **Unerwünschte Handlung** | SHOULD NOT, NOT RECOMMENDED | should not | Generiert eine Warnung für das Architektur-Review |
| **Erlaubnis** | MAY, OPTIONAL | may, need not | Erfasst eine optionale Funktion; kein Kriterium zur Ablehnung des Releases |
| **Möglichkeit oder Fähigkeit** | keine | can, cannot | Dokumentiert eine Systemeigenschaft, keine Anforderung |
| **Externe Randbedingung** | keine | must | Hinterlegt eine Umgebungsvorbedingung, z. B. physikalische Gesetze oder Normen |
| **Tatsachenfeststellung** | Kleinbuchstaben, is, will | is, will | Erfasst Kontextinformation; stellt keine Anforderung dar |

Die Gegenüberstellung verdeutlicht, dass identische Wörter in unterschiedlichen Dokumentenkontexten eine grundlegend andere normative Verbindlichkeit besitzen. Folglich speichert das System die Modalität stets zusammen mit dem Bezeichner der Dokumentenkonvention, sodass ein Auditor jederzeit nachvollziehen kann, nach welchen Regeln ein Satz klassifiziert wurde.

### 1.1. Syntaktische Fallstricke und semantische Ambiguität natürlicher Sprache

Das Auffinden des Signalworts SHALL mittels regulärer Ausdrücke deckt nur einen Bruchteil realer Anforderungen ab. Technische Dokumente werden von heterogenen Autoren verfasst – häufig Nicht-Muttersprachlern –, wodurch sich drei charakteristische Fehlermuster wiederholen.

**Passivkonstruktionen ohne Akteur.** Im Satz *„Data shall be validated before transmission“* bleibt undefiniert, welche Komponente für die Prüfung verantwortlich ist: der Sensortreiber, der Kommunikationscontroller oder die Anwendungsebene. Eine Anforderung ohne zugewiesenen Akteur kann keiner Architekturkomponente zugeordnet und somit nicht auf Code-Ebene rückverfolgt werden. Ein syntaktisches Dependenzparsing (wie in [Kapitel 13](ch13-language-variability-vs-determinism.md#31-syntaktisches-dependenzparsing-und-aktantenextraktion) beschrieben) identifiziert das fehlende Subjekt und liefert die formale Grundlage zur Zurückweisung an den Autor.

**Implizite Normativität.** Formulierungen wie *„The ECU is responsible for monitoring battery voltage“* oder *„The firmware needs to reboot if a watchdog timeout occurs“* enthalten kein formales *shall*, stellen inhaltlich jedoch zwingende funktionale Anforderungen dar. Das Expertensystem führt daher ein Verzeichnis quasimodaler Konstruktionen (*is responsible for*, *has to*, *needs to*, *is required to*) und markiert derartige Sätze als Kandidaten für eine manuelle Begutachtung.

**Eingebettete Ausnahmeregelungen.** Der Satz *„The system shall maintain 50 Hz PWM frequency under all load conditions, except during initial power-up calibration where 20 Hz is permitted for a maximum of 200 ms“* vereint eine globale Anforderung, eine Ausnahmebedingung, ein alternatives Systemverhalten und eine Zeitrestriktion. Das Expertensystem muss einen solchen Satz in atomare logische Zweige zerlegen, da die Ausnahmebedingung andernfalls bei der Formalisierung verloren geht.

Modalität und Syntax müssen demnach zwingend gekoppelt analysiert werden: Die Modalität bestimmt, ob eine Verpflichtung vorliegt, während die Syntax aufdeckt, wer wann welche Aktion unter welchen Randbedingungen ausführen muss. Im nächsten Schritt werden die heterogenen Satzstrukturen auf eine standardisierte Grammatik normiert.

### 1.2. Linter für semantische Mehrdeutigkeit und lexikografische deontische Vorzugsordnung

In der Architektur des Systems zur automatisierten Wissensakquisition ([Kapitel 10](ch10-knowledge-acquisition-systems.md)) bildet die Phase des vorbereitenden lexikalischen Lintings die erste deterministische Barriere vor dem Aufbau des formalen abstrakten Syntaxbaums (AST). In der Praxis industrieller Systementwicklung wimmeln technische Spezifikationen von nicht formalisierter, vager Semantik: Superlativen und relativen Adjektiven („akzeptable Latenz“, „ausreichende Kühlleistung“, „hohe Zuverlässigkeit“, „minimaler Energieverbrauch“). Gelangt ein derart unscharfes Prädikat ohne strikte numerische Verankerung in die Regelbasis, ist der symbolische SMT-Solver außerstande, dafür eine Grenzinvariante zu synthetisieren, während ein probabilistisches Sprachmodell beginnt, halluzinierte Annahmen über zulässige Ansprechschwellen zu generieren. Die Konsequenz ist entweder eine Blockade des formalen Beweises (Deadlock des Verifizierers) oder das katastrophale Übersehen einer Gefahrensituation während des Betriebs des Hardwarekomplexes.

Um diesen Mangel systematisch zu unterbinden, wurde vom Autor ein spezialisierter Linter für semantische Mehrdeutigkeit (*Ambiguity Linter*) entwickelt, der auf einem geschlossenen Lexikon toxischer Deskriptoren $\mathcal{D}_{\mathrm{toxic}}$ (definiert gemäß den Richtlinien von ISO/IEC/IEEE 29148 [[13]](#src-13)) operiert. Jede Anforderung, die ein qualitatives Adjektiv ohne Bindung an ein numerisches Intervall des SI-Einheitensystems enthält, wird durch den Compiler unmittelbar zurückgewiesen:

```math
\text{Inkorrekt: } \mathrm{FastResponse}(\mathrm{ECU}) \quad \longrightarrow \quad \text{Korrekt: } t_{\mathrm{resp}} \le 150\,\mu\text{s}
```

Gleichzeitig tritt in realen Regelwerken das Problem deontischer Kollisionen auf: wenn eine internationale Norm eine Operation zwingend untersagt (`PROHIBITION`), während die internen technischen Lieferbedingungen des Werks diese temporär gestatten (`PERMISSION`). Ohne deterministische Konfliktlösung verharrt die Inferenzmaschine im Stillstand.

Zur deterministischen Auflösung von Normenkonflikten wird eine lexikografische Vorzugsordnungsfunktion (*Lexicographical Preference Ordering*) eingeführt:

```math
\mathrm{Rank}(\mathcal{R}) = \langle \mathrm{Tier}(\mathrm{Source}),\ \mathrm{ModalityRank}(M),\ \mathrm{Specificity}(\mathcal{R}) \rangle
```

Parameter und Komponenten des lexikografischen Rangs:
- $\mathrm{Tier}(\mathrm{Source}) \in \{1, 2, 3\}$ — hierarchischer Rang der normativen Wissensquelle: $\mathrm{Tier}(\text{Harmonisierte internationale Norm ISO/IEC}) = 3 \succ \mathrm{Tier}(\text{Unternehmensinterne Werksnorm}) = 2 \succ \mathrm{Tier}(\text{Ingenieurheuristik}) = 1$;
- $\mathrm{ModalityRank}(M) \in \{1, 2\}$ — imperatives Gewicht der deontischen Modalitätskategorie: $\mathrm{ModalityRank}(\text{PROHIBITION}) = \mathrm{ModalityRank}(\text{OBLIGATION}) = 2 \succ \mathrm{ModalityRank}(\text{PERMISSION}) = 1$. Ein striktes Verbot und eine kategorische Verpflichtung besitzen absoluten lexikografischen Vorrang vor einer dispositiven Erlaubnis;
- $\mathrm{Specificity}(\mathcal{R}) \in \mathbb{N}$ — Länge der Kette prädikativer Kontextbedingungen (eine hochspezialisierte Regel sticht eine allgemeine Regel nach dem Rechtsprinzip *Lex Specialis Derogat Legi Generali* aus).

#### 1.2.1. Geschlossener Kreis der ingenieurtechnischen Berechnung des deontischen Rangs
1. **Laufzeitsteuerung und Rechenfluss (Control Flow & Runtime Decisions):**
   Beim Aufeinandertreffen zweier zueinander im Widerspruch stehender Regeln $\mathcal{R}_1 \land \neg \mathcal{R}_2$ berechnet der Wissens-Arbitrator den Differenzvektor $\Delta \mathbf{r} = \mathrm{Rank}(\mathcal{R}_1) - \mathrm{Rank}(\mathcal{R}_2)$. Gilt $\Delta \mathbf{r} > \mathbf{0}$, setzt sich Regel $\mathcal{R}_1$ deterministisch durch (*Strict Undercut*), während $\mathcal{R}_2$ als temporär außer Kraft gesetzt markiert wird; gilt $\Delta \mathbf{r} = \mathbf{0}$ (Kollision gleichrangiger Normen), schaltet das System in den Schutzmodus *Fail-Closed* mit dem Urteil `QUALIFIED_CONFLICT` und eskaliert den Konflikt an den Chefarchitekten.
2. **Hardware-Dimensionierung und Infrastrukturbeschränkungen (Hardware Dimensioning):**
   Die Rangberechnung nach obiger Formel erfolgt in deterministischer Zeit $O(1)$ über einer gepackten 64-Bit-Ganzzahl `uint64` (wobei 16 Bit für den Tier, 16 Bit für die Modalität und 32 Bit für die Spezifität reserviert sind). Dies erfordert lediglich 8 Byte Speicher pro Regel und ermöglicht das lineare Sortieren von 100.000 Regeln im SRAM-Speicher des Controllers in weniger als $1{,}8\,\text{ms}$.
3. **Praktisches numerisches Rechenbeispiel (Worked Numerical Example):**
   Angenommen, die ISO 26262 ASIL D untersagt den Eintritt in den Rekuperationsmodus bei einem Fehler auf dem CAN-Bus ($\mathrm{Tier} = 3, M = \text{PROHIBITION} \to 2, \mathrm{Specificity} = 4$), während eine interne Werkstattvorschrift dies zu Kalibrierzwecken gestattet ($\mathrm{Tier} = 2, M = \text{PERMISSION} \to 1, \mathrm{Specificity} = 2$). Der lexikografische Rang der Sicherheitsnorm $\langle 3, 2, 4 \rangle$ dominiert die Werksvorschrift $\langle 2, 1, 2 \rangle$ bereits in der ersten Koordinate ($3 > 2$). Das System blockiert die sicherheitskritische Aktion deterministisch und verhindert eine Havarie.

> [!NOTE]
> **Drei Koordinaten der lexikografischen Anforderungsrangfolge $\mathrm{Rank}(\mathcal{R})$**
> 1. Hierarchie der Wissensquelle ($\mathrm{Tier}(\mathrm{Source})$): Vorrang harmonisierter internationaler Normen (ISO/IEC) vor internen Werksnormen und Ingenieurheuristiken.
> 2. Imperativität der deontischen Modalität ($\mathrm{ModalityRank}(M)$): absoluter Vorrang von Verboten und Verpflichtungen (`PROHIBITION`, `OBLIGATION`) vor dispositiven Erlaubnissen (`PERMISSION`).
> 3. Spezifität von Kontextbedingungen ($\mathrm{Specificity}(\mathcal{R})$): Vorrang fachspezifischer Ausnahmeregeln vor allgemeinen Normen nach dem Prinzip *Lex Specialis*.
> *Laufzeitsteuerung:* Der Vergleich des 64-Bit gepackten Deskriptors erfolgt in $O(1)$, und bei einer Kollision gleichrangiger Normen ($\Delta\mathbf{r} = \mathbf{0}$) greift das Fail-Closed-Gateway.

---

### 1.3. Unteilbarkeit deontischer Quantoren, Multibyte-Unicode-Grenzen und Span-Erhaltung

Im Subsystem zum Parsing technischer Regelwerke (RFC, ISO, AUTOSAR) stellt die Fragilität der Tokenisierung an den Nahtstellen modaler Quantoren und Multibyte-Zeichen der UTF-8-Codierung eine kritische Quelle semantischer Defekte dar. Ein naiver sequenzieller Lexer interpretiert den Eingabestrom als lineare Folge von Lexemen, die durch Leerzeichen oder Zeilenumbrüche getrennt sind. Im industriellen Layout normativer Dokumente treten Zeilenumbrüche jedoch nicht selten unmittelbar innerhalb eines zusammengesetzten modalen Quantors auf:

```text
The client MUST
NOT generate a 100-continue header if...
```

Trifft ein naiver Tokenizer auf das Wort `MUST`, registriert er eine strikte Verpflichtung (`Obligation = MUST`), interpretiert den Zeilenumbruch `\n` fälschlich als Satzende oder Klauselgrenze und verwirft die nachfolgende Negationspartikel `NOT`. Infolgedessen wird ein kategorisches Verbot in sein direktes Gegenteil – eine Erlaubnis oder gar Pflicht – invertiert. In Systemen der funktionalen Sicherheit nach ASIL D führt dies zur fatalen Zerstörung von Sicherheitsperimetern.

Darüber hinaus treten zwei flankierende Systemprobleme auf:
- **Zerstörung der Start-Byte-Offsets (Start-Offset Obliteration):** Entfernt der Tokenizer bei der Textnormalisierung führende Leer- oder Steuerzeichen, ohne das physische Offset vom Beginn der Eingabedatei beizubehalten, erzeugt das Register der Hardware-Beweiskette (`%ebx` register custody) ein fehlerhaftes Intervall `[byte_start, byte_end]`. Der byteweise SHA-256-Hash dieses Segments stimmt nicht mehr mit der Primärquelle überein, was eine lückenlose Verifikation verunmöglicht.
- **Toter Winkel des Linters bei Nicht-ASCII-Interpunktion (UTF-8 Boundary Blindness):** Standardmäßige reguläre Ausdrücke mit Wortgrenzen `\b` sind ausschließlich für ASCII konzipiert. Bei der Analyse deutscher oder europäischer Dokumente mit Gedankenstrichen (— U+2014, – U+2013), typografischen Anführungszeichen („ “ « »), Apostrophen (’ U+2019) oder geschützten Leerzeichen versagen `\b`-Grenzen. Verbotene oder toxische Deskriptoren (z. B. *„verboten—“* oder *„gefährlich…“*) werden vom Linter schlicht übersehen.

Zur Gewährleistung vollständiger deontischer und kryptografischer Integrität wurde vom Autor ein obligatorisches Vorab-Falten zusammengesetzter Modalphrasen in **unteilbare atomare deontische Lexeme (Composite Deontic Atoms)** eingeführt:

```math
\mathrm{Tok}_{\mathrm{MUST\_NOT}} = \mathrm{Regex}\bigl(\text{"(?i)\b(must|shall|is\s+required\s+to)\s*[\r\n\t ]+\s*not\b"}\bigr)
```

Parameter und Scanregeln:
- Die Grammatik behandelt das Paar $\langle\mathrm{MODAL}, \mathrm{NEGATION}\rangle$ als unteilbares binäres Atom;
- Beliebige Trennzeichen aus Leerzeichen oder Zeilenumbrüchen innerhalb eines maximalen linguistischen Fensters von $W \le 3$ Token dürfen den deontischen Kontext nicht aufbrechen;
- Die vollständige Bewahrung der physischen Koordinaten wird durch die Struktur `GroundedToken` garantiert:

```go
type GroundedToken struct {
    Lexeme      string
    ByteStart   uint64 // Absolutes physisches Offset in der Datei
    ByteEnd     uint64 // Physisches Ende des Lexems
    Line        uint32
    Column      uint32
}
```

Jegliche Normalisierung des Textes ist **ausschließlich in der virtuellen Repräsentation des Tokens** zulässig; die physischen Offsets `ByteStart` und `ByteEnd` bleiben strikt unverändert.

Zur Beseitigung toter Winkel des Lintings wird ein Scan-Algorithmus über die Runen-Decodierfunktion `utf8.DecodeRuneInString` mit Klassifikation anhand von Unicode-Tabellen (`unicode.IsLetter`, `unicode.IsDigit`, `unicode.IsPunct`, `unicode.IsSpace`) eingesetzt. Dies garantiert 0 % falsch-negative Auslassungen toxischer Begriffe an Nahtstellen mit Nicht-ASCII-Interpunktion.

#### 1.3.1. Geschlossener Kreis der ingenieurtechnischen Berechnung der Span-Integrität
1. **Laufzeitsteuerung und Rechenfluss (Control Flow & Runtime Decisions):**
   Wird eine Faltung `MUST_NOT` erkannt, setzt der Parser das Verbots-Flag `InvertDeonticPolarity = true`. Stimmt der berechnete SHA-256-Hash des Byte-Fragments `RawFile[ByteStart:ByteEnd]` nicht mit dem im Wissensbasis-Eintrag hinterlegten Hash überein, wird die Translation der Anforderung sofort mit dem Fehler `ERR_CUSTODY_HASH_MISMATCH` abgebrochen und der Regelstatus in den sicheren Zustand Fail-Closed überführt.
2. **Hardware-Dimensionierung und Infrastrukturbeschränkungen (Hardware Dimensioning):**
   Die Struktur `GroundedToken` belegt exakt 32 Byte Speicher (16 Byte `string`-Header, 8 Byte `ByteStart`, 8 Byte `ByteEnd`, 4 Byte `Line`, 4 Byte `Column`). Ein Puffer für 1.000.000 Token eines umfangreichen Industriestandards beansprucht lediglich 32 MB Arbeitsspeicher, was problemlos im RAM der Ingenieurs-Workstation Platz findet.
3. **Praktisches numerisches Rechenbeispiel (Worked Numerical Example):**
   Beim Parsen der Zeichenkette `The ECU shall\nnot disable cooling` mit einer Länge von 33 Byte und einem Start-Offset von $14\,200$: Ein naiver Tokenizer erzeugte $\text{SHALL}$ mit dem Offset-Bereich $[14\,208, 14\,213]$ (fataler Invarianteffekt: Pflicht zur Abschaltung der Kühlung!). Der Spans erhaltende Tokenizer formiert den atomaren Quantor `shall not` mit dem einheitlichen Intervall $[14\,208, 14\,217]$ und dem berechneten SHA-256-Hash $H = \text{0x9f8a\dots}$, wodurch die Überhitzung der Batterie zuverlässig verhindert wird.

---

## 2. EARS-Schablonen: Strukturelle Standardisierung von Engineering-Anforderungen

Freitext in natürlicher Sprache widersetzt sich direkter logischer Inferenz: Autoren verschachteln Bedingungen, verkomplizieren Satzgefüge und verwischen den Gültigkeitsbereich von Negationen. Versucht ein automatisierter Parser, aus beliebiger Prosa unmittelbar einen abstrakten Syntaxbaum (AST) abzuleiten, übersteigt die Fehlerquote bei der Rollenzuweisung 40 %. Für eine zuverlässige semantische Translation müssen Anforderungen daher zunächst in ein kontrolliertes grammatikalisches Gerüst überführt werden.

Als de-facto-Standard hat sich hierbei der von Alistair Mavin und Kollegen bei Rolls-Royce im Jahr 2009 für Triebwerksregelsysteme entwickelte Ansatz *Easy Approach to Requirements Syntax* (EARS) etabliert [[7]](#src-7). EARS beschränkt Anforderungen auf eine überschaubare Menge standardisierter Schablonen, wobei jede Schablone exakt einem Vorbedingungstyp entspricht.

```mermaid
flowchart TD
    accTitle: Auswahl der EARS-Schablone
    accDescr: Entscheidungsbaum zur Bestimmung des Anforderungstyps anhand des Vorhandenseins und der Art von Vorbedingungen: Zustand, Ereignis, Fehlerfall, optionale Funktion oder deren Kombination.

    Start["Analyse der Satzstruktur"] --> Q_Pre{"Vorbedingung, Zustand oder Trigger vorhanden?"}
    Q_Pre -- "Nein" --> Ubi["<b>Ubiquitär</b><br/><i>The system shall...</i>"]
    Q_Pre -- "Ja" --> Q_Cond{"Art der Vorbedingung?"}

    Q_Cond -- "Zeitlich andauernder Zustand" --> State["<b>Zustandsgesteuert</b><br/><i>While in state..., the system shall...</i>"]
    Q_Cond -- "Diskretes Ereignis" --> Event["<b>Ereignisgesteuert</b><br/><i>When trigger..., the system shall...</i>"]
    Q_Cond -- "Fehler oder unerwünschtes Ereignis" --> Fault["<b>Reaktion auf unerwünschtes Verhalten</b><br/><i>If fault..., then the system shall...</i>"]
    Q_Cond -- "Vorhandensein einer optionalen Funktion" --> Opt["<b>Optionale Funktion</b><br/><i>Where feature..., the system shall...</i>"]

    State & Event & Fault --> Complex["<b>Komplex</b><br/><i>While..., when..., the system shall...</i>"]

    classDef root fill:#e1f5fe,stroke:#0288d1,stroke-width:2px,color:#01579b;
    classDef choice fill:#fff3e0,stroke:#f57c00,stroke-width:2px,color:#e65100;
    classDef leaf fill:#e8f5e9,stroke:#388e3c,stroke-width:2px,color:#1b5e20;
    classDef comp fill:#f3e5f5,stroke:#7b1fa2,stroke-width:2px,color:#4a148c;

    class Start root;
    class Q_Pre,Q_Cond choice;
    class Ubi,State,Event,Fault,Opt leaf;
    class Complex comp;
```

Der Entscheidungsbaum stellt lediglich zwei Fragen: Enthält der Satz eine Vorbedingung und welcher Art ist diese? Die Antwort bestimmt die Schablone und legt die logische Repräsentation für den nachfolgenden Formalisierungsschritt fest.

**Ubiquitäre Anforderung (*Ubiquitous*):** Gilt permanent, ohne Bindung an Ereignisse oder spezifische Betriebszustände. Schablone: `The <system name> shall <system response>.` Beispiel: *„The CAN controller shall support extended 29-bit identifiers“*.

**Ereignisgesteuerte Anforderung (*Event-driven*):** Definiert die Systemreaktion auf ein diskretes Startereignis. Schablone: `When <trigger>, the <system name> shall <system response>.` Beispiel: *„When the E-STOP button is pressed, the motor driver shall disable gate drive outputs within 5 ms“*.

**Zustandsgesteuerte Anforderung (*State-driven*):** Gilt ausschließlich während der Dauer eines definierten Systemzustands. Schablone: `While <in a specific state>, the <system name> shall <system response>.` Beispiel: *„While in PRE-CHARGE mode, the BMS shall limit the pre-charge resistor current to 10 A“*.

**Reaktion auf unerwünschtes Verhalten (*Unwanted behaviour*):** Beschreibt Maßnahmen bei Fehlfunktionen, Protokollverletzungen oder Grenzwertüberschreitungen. Schablone: `If <trigger>, then the <system name> shall <system response>.` Beispiel: *„If the cell temperature exceeds 65 °C, then the cooling controller shall activate the refrigerant pump at 100% duty cycle“*.

**Optionale Funktion (*Optional feature*):** Greift nur, wenn eine bestimmte Hardwarekomponente oder Lizenz vorhanden ist. Schablone: `Where <feature is included>, the <system name> shall <system response>.` Beispiel: *„Where the hardware watchdog is populated, the CPU supervisor shall toggle the WDI pin every 50 ms“*.

**Komplexe Anforderung (*Complex*):** Kombiniert mehrere Bedingungen, beispielsweise Zustand und Trigger: *„While in CHARGING state, when the charge plug is unlocked, the charger shall open the high-voltage interlock loop within 10 ms“*.

EARS strukturiert Vorbedingungen und Reaktionen, legt jedoch keine eindeutige formale Semantik fest. Eine zustandsbasierte Anforderung kann je nach Vollverb und zeitlichem Kontext eine permanente Invariante, eine begrenzte Übergangszeit oder eine Eintrittsbedingung erfordern. Abweichungen von den EARS-Mustern indizieren Prüfbedarf, stellen jedoch keinen automatischen Beweis für einen Mangel dar. Formale Vorbedingungen, Maßeinheiten und Ausnahmebehandlungen müssen explizit spezifiziert werden.

### 2.1. Grammatikgesteuerte Decodierung (GBNF) und die neuro-symbolische Gültigkeitsschranke

Bei der Interaktion zwischen Fachingenieur und Expertensystem über natürlichsprachliche Anfragen entsteht eine fundamentale Gefahr der Vermischung zweier sich gegenseitig ausschließender Funktionszonen:
1. **Relevanzzone (*Relevance Zone*):** heuristische Erfassung natürlicher Sprache, Tippfehlerkorrektur, Synonymanreicherung mittels lokaler kompakter Sprachmodelle (*Small Language Models*, SLM);
2. **Validitäts- bzw. Gültigkeitszone (*Eligibility Zone*):** deterministische logische Überprüfung normativer Toleranzen, mathematische Berechnung von Degressionen und Sicherheitsbewertungen nach ISO 26262 und DO-178C.

Erhält ein probabilistisches Sprachmodell die Freiheit, logische Ausgangsaussagen unbeschränkt zu generieren, entsteht die Schwachstelle einer unkontrollierten Prädikatinjektion (*Wildcard Prompt Injection*): Eine freie Grammatik erlaubt es dem Modell, willkürliche Entitäten zu erfinden (`trigger ::= [a-zA-Z0-9]+`) und damit die formale Wissensbasis zu umgehen. Mehr noch: Bei erkannter Unvollständigkeit des Wissens beginnt das System, scheinbar plausible, jedoch völlig frei erfundene Pseudoerklärungen für Sperrentscheidungen zu fabrizieren.

Zur Errichtung einer unüberwindbaren neuro-symbolischen Gültigkeitsschranke wird die Sprachgenerierung in der Validitätszone durch den Mechanismus der **grammatikgesteuerten Decodierung (*GBNF Grammar-Constrained Sampling*)** auf Basis eines geschlossenen Lexikons des aktiven Wissenspakets maskiert:

```math
\mathcal{V}_{\mathrm{valid}} = \mathrm{Entities}(\mathcal{K}) \cup \mathrm{Predicates}(\mathcal{K}) \cup \mathrm{Units}(\mathcal{K}) \cup \mathrm{Keywords}_{\mathrm{EARS}}
```

Parameter und Bestandteile des geschlossenen Decodierungsraums:
- $\mathcal{K}$ — aktives, verifiziertes Wissenspaket;
- $\mathrm{Entities}(\mathcal{K})$ — endliche Menge verifizierter Komponenten- und Signalbezeichner (`BMS_Relay`, `Cell_Voltage`, `PreCharge_Resistor`);
- $\mathrm{Predicates}(\mathcal{K})$ — Menge formalisierter Relationsprädikate und Toleranzoperatoren (`exceeds`, `is_equal`, `activates`);
- $\mathrm{Units}(\mathcal{K})$ — standardisierte physikalische Maßeinheiten des SI-Systems (`ms`, `V`, `A`, `deg_C`);
- $\mathrm{Keywords}_{\mathrm{EARS}} = \{\text{WHEN}, \text{WHILE}, \text{IF}, \text{THEN}, \text{WHERE}, \text{THE}, \text{SHALL}, \text{SHALL NOT}\}$ — reservierte normative Schlüsselwörter der EARS-Syntax.

Während der autoregressiven Auswahl jedes Tokens legt der Syntax-Decoder eine binäre Maske über die Ausgabeschicht der Modell-Logits:

```math
\mathrm{Logit}'(w) = \begin{cases} \mathrm{Logit}(w), & \text{wenn } w \in \mathcal{V}_{\mathrm{valid}} \text{ und der Übergang im GBNF-Automaten zulässig ist}, \\ -\infty, & \text{sonst.} \end{cases}
```

Jeder Versuch, ein Token außerhalb des Lexikons oder entgegen den Übergangsregeln der EARS-Grammatik zu generieren, wird mit einer Wahrscheinlichkeit von null ($P(w) = 0$) deterministisch abgeschnitten.

Führt die Verifikation zu einem Negativurteil (`REFUSAL`), formuliert das System die Begründung keineswegs über freie Textgenerierung, sondern durch direkten Durchlauf des Beweisgraphen (kontrafaktische Erklärung): Es benennt exakt die verletzte Vorbedingung (z. B. $T_{\mathrm{sensor}} = 135^\circ\text{C} > T_{\mathrm{max}} = 125^\circ\text{C}$), zitiert den normativen Textabschnitt mit physischem Byte-Offset und generiert die minimale Korrekturmaßnahme in strikter deontischer Form.

#### 2.1.1. Geschlossener Kreis der ingenieurtechnischen Berechnung der Grammatikmaskierung
1. **Laufzeitsteuerung und Rechenfluss (Control Flow & Runtime Decisions):**
   Beim Laden der Regelbasis kompiliert der Syntax-Compiler das GBNF-Skript in einen deterministischen endlichen Automaten (DFA). Trifft eine Benutzeranfrage ein, wird die Eingabe auf eine der 4 kanonischen EARS-Schablonen normalisiert. Versucht das Modell, ein nach dem Automatenzustand unzulässiges Token auszuwählen, stoppt das semantische Gateway die Generierung unverzüglich mit dem Fehler `EARS_SYNTAX_REJECT`, ohne Rechenressourcen von GPU oder NPU für die weitere Kontextentfaltung zu verschwenden.
2. **Hardware-Dimensionierung und Infrastrukturbeschränkungen (Hardware Dimensioning):**
   Die Bitmasken-Tabelle für ein Lexikon mit $\lvert \mathcal{V}_{\mathrm{valid}} \rvert = 4096$ Token belegt exakt 512 Byte (64 Wörter zu je `uint64`). Dies passt garantiert vollständig in den ultraschnellen L1-Befehlscache ($32\,\text{KB}$) eines NPU-Beschleunigers oder CPU-Kerns, wodurch die Maskierung in einem einzigen Taktzyklus ($< 1\,\text{ns}$) angewendet wird.
3. **Praktisches numerisches Rechenbeispiel (Worked Numerical Example):**
   Versucht das Sprachmodell, einen nicht existierenden Parameter `battery_fluid_viscosity` zu generieren (nicht im Vokabular $\mathcal{V}_{\mathrm{valid}}$ von 1850 Entitäten enthalten), wird der Logit des ersten Worttokens sofort auf $-\infty$ maskiert. Das alternative valide Token `cell_temperature` mit einem Ausgangs-Logit von $14{,}2$ erhält die resultierende Wahrscheinlichkeit von $1.00$. Die Anfrage wird eindeutig in die EARS-Form *„If cell temperature exceeds 65 deg_C, then BMS shall open main relay“* überführt – unter vollständiger Eliminierung jeglichen Halluzinationsrisikos.

> [!NOTE]
> **Fünf kanonische EARS-Schablonen zur Beseitigung syntaktischer Ambiguität**
> 1. Allgegenwärtige Anforderungen (Ubiquitous): bedingungslose Systeminvarianten ohne zeitliche oder ereignisbezogene Auslöser.
> 2. Ereignisgesteuerte Anforderungen (Event-driven): zwingend deterministische Systemreaktion auf ein diskretes externes Ereignis („Wenn...“).
> 3. Zustandsgesteuerte Anforderungen (State-driven): Parametergrenzen, die ausschließlich während des Verweilens im definierten Zustand gelten („Solange...“).
> 4. Reaktionen auf unerwünschtes Verhalten (Unwanted behavior): Notfall-Schutzprotokolle zur Rückkehr in den sicheren Zustand bei Abweichungen („Falls...“).
> 5. Optionale Funktionsanforderungen (Optional features): Vorgaben, die nur bei hardwareseitiger oder konfigurativer Verfügbarkeit greifen („Wo vorhanden...“).
> *Laufzeitsteuerung:* Der Anforderungs-Compiler weist jede Konstruktion zurück, die die deterministische EARS-Grammatik verletzt, und blockiert unstrukturierten Text.

---

## 3. Mathematische Formalisierung: Von Schablonen zur formalen Logik

Eine EARS-Schablone strukturiert den Text, versetzt Maschinen jedoch noch nicht in die Lage, Anforderungen automatisiert zu verifizieren. Hierzu muss die Schablone in eine formale Zielsprache übersetzt werden: Prädikatenlogik erster Stufe, temporale Logik oder das standardisierte SMT-LIB-Format für automatische Theorembeweiser und Solver.

### 3.1. Prädikatenlogik erster Stufe

In evidenzbasierten Expertensystemen wird die Anforderung an einen zustandsbehafteten Steuerungsautomaten als strikter Zusicherungskontrakt (*Assume-Guarantee Contract*) über einem endlichen Zustandsraum $\mathcal{S}$, einem Vektor gemessener Eingangssignale $\mathbf{x}\in\mathcal{X}$ und einem Vektor von Steuerausgängen $\mathbf{y}\in\mathcal{Y}$ formalisiert:

```math
\forall s\in\mathcal{S},\ \forall\mathbf{x}\in\mathcal{X}:\quad \Phi_{\mathrm{pre}}(s,\mathbf{x})\Rightarrow\exists s'\in\mathcal{S},\ \exists\mathbf{y}\in\mathcal{Y}:\ \bigl(\Phi_{\mathrm{post}}(s',\mathbf{y})\land\mathcal{T}(s,s')\bigr)
```

Parameter und mathematische Komponenten des Kontrakts:

- $\mathcal{S}$ ist der endliche Raum diskreter Automatenzustände (z. B. $\mathcal{S} = \{\text{INIT}, \text{STANDBY}, \text{CHARGE}, \text{DISCHARGE}, \text{FAULT}\}$);
- $\mathbf{x} \in \mathcal{X} \subseteq \mathbb{R}^n$ ist der Vektor kontinuierlicher und diskreter Telemetriedaten (Zellspannungen, Ströme, Temperaturen, Steuerflags);
- $\mathbf{y} \in \mathcal{Y} \subseteq \mathbb{R}^m$ ist der Vektor aktorischer Stellgrößen (PWM-Tastgrade, Schaltsignale der Relais-Treiber);
- $s, s' \in \mathcal{S}$ bezeichnen den aktuellen und den Folgezustand des Automaten;
- $\Phi_{\mathrm{pre}}(s,\mathbf{x}): \mathcal{S} \times \mathcal{X} \to \{\text{True}, \text{False}\}$ ist die prädikative Vorbedingung, die Aktivitätszustand (While), Startereignis (When) und Fehlersignalisierung (If) verknüpft;
- $\Phi_{\mathrm{post}}(s',\mathbf{y}): \mathcal{S} \times \mathcal{Y} \to \{\text{True}, \text{False}\}$ ist die prädikative Nachbedingung, die die geforderte Systemreaktion (shall) festlegt;
- $\mathcal{T}(s,s'): \mathcal{S} \times \mathcal{S} \to \{\text{True}, \text{False}\}$ ist die Übergangsrelation des Automaten, die zulässige Zustandssprünge sowie das Übergangszeitbudget $\Delta t \le t_{\mathrm{timeout}}$ beschränkt.

Praktische Anwendung und Engineering-Entscheidungen:
- **Vollständigkeitsprüfung (Non-blocking):** Der SMT-Solver verifiziert, dass für jede Kombination $(s, \mathbf{x})$, die $\Phi_{\mathrm{pre}}$ erfüllt, mindestens ein gültiger Folgezustand existiert. Wird eine Eingabe $\mathbf{x}^*$ identifiziert, für die kein Übergang möglich ist, wird der Mangel `REQ_DEFECT_DEADLOCK` (Deadlock / Spezifikationslücke) protokolliert.
- **Prüfung auf Determinismus:** Existieren für ein Tupel $(s, \mathbf{x})$ divergierende Folgezustände $(s'_1, \mathbf{y}_1) \ne (s'_2, \mathbf{y}_2)$, meldet der Verifizierer eine Ambiguität `REQ_DEFECT_AMBIGUITY`. Eine nichtdeterministische Spezifikation wird für die Baseline-Freigabe bis zur Korrektur durch den Autor gesperrt.

### 3.2. Lineare temporale Logik (LTL) für Zeitinvarianten

Eingebettete Systeme agieren in der Zeitdimension; statische Prädikate reichen zur Abbildung dynamischen Verhaltens nicht aus. Amir Pnueli etablierte die lineare temporale Logik (*Linear Temporal Logic*, LTL) für den formalen Nachweis von Systemeigenschaften [[8]](#src-8). Darin bezeichnet der Modaloperator $\Box$ „immer (in allen zukünftigen Zuständen)“ und der Operator $\Diamond$ „irgendwann in der Zukunft“. Die metrische temporale Logik (*Metric Temporal Logic*, MTL), die Ron Koymans für Echtzeitanforderungen formulierte, erweitert diese Operatoren um explizite Zeitintervalle, sodass $\Diamond_{\le\tau}$ für „spätestens nach Ablauf der Zeitspanne $\tau$“ steht [[9]](#src-9). Die Signal-Temporallogik (*Signal Temporal Logic*, STL) von Maler und Nickovic überträgt diesen Formalismus auf kontinuierliche physikalische Signale wie Strom oder Temperatur [[10]](#src-10).

In der Praxis dominieren zwei Klassen temporallogischer Eigenschaften. Eine **Sicherheitsinvariante (*Safety Invariant*)** garantiert, dass ein unzulässiger, kritischer Systemzustand zu keinem Zeitpunkt eintritt:

```math
\Box\,\neg\bigl(\mathit{Current}>I_{\max}\land\mathit{ContactorState}=\mathit{CLOSED}\bigr)
```

Parameter und physikalische Dimensionen:

- $\mathit{Current} \in \mathbb{R}_{\ge 0}$ ist der gemessene physikalische Laststrom des Leistungspfads in Ampere ($\text{A}$);
- $I_{\max} \in \mathbb{R}_{> 0}$ ist der thermisch zulässige Abschaltstrom in Ampere ($\text{A}$), dimensioniert nach der Stromtragfähigkeit der Stromschiene (z. B. $I_{\max} = 450\,\text{A}$);
- $\mathit{ContactorState} \in \{\mathit{OPEN}, \mathit{CLOSED}\}$ repräsentiert den diskreten Schaltzustand des Hochvoltschützes;
- $\Box$ ist der LTL-Temporaloperator „immer“, der die Gültigkeit der Formel für alle Zeitschritte $t \in [0, \infty)$ einfordert;
- $\neg$ und $\land$ bezeichnen die klassischen logischen Operatoren für Negation und Konjunktion.

Praktische Anwendung und Engineering-Entscheidungen:
- **Hardwarenahes Monitoring:** Diese Bedingung wird durch einen dedizierten Hardware-Supervisor mit einer Abtastrate von mindestens 10 kHz autark vom Lastzustand des Haupt-Mikrocontrollers überwacht.
- **Fail-Safe-Aktion:** Wird $\mathit{Current} > I_{\max}$ bei geschlossenem Schütz ($\mathit{ContactorState} = \mathit{CLOSED}$) detektiert, löst eine Hardwareschaltung innerhalb von $t_{\mathrm{reaction}} \le 2\,\text{ms}$ eine pyrotechnische Trennstelle (Pyro-Fuse) aus und schaltet den Hochvoltpfad galvanisch spannungsfrei.

Die **beschränkte Lebendigkeit (*Bounded Liveness*)** fordert dagegen, dass das System nach Auftreten eines Fehlers garantiert innerhalb eines definierten Zeitintervalls in den sicheren Zustand übergeht:

```math
\Box\Bigl(\mathit{FaultTriggered}\Rightarrow\Diamond_{\le\tau}\,\mathit{SafeStateAchieved}\Bigr)
```

Parameter und zeitliche Toleranzen:

- $\mathit{FaultTriggered} \in \{\text{True}, \text{False}\}$ ist ein diskretes Fehlerflag (z. B. Zellübertemperatur $> 60^\circ\text{C}$ oder CAN-Kommunikationsausfall);
- $\mathit{SafeStateAchieved} \in \{\text{True}, \text{False}\}$ signalisiert das verifizierte Erreichen des sicheren Zustands (Schütze geöffnet, Traktionspfad spannungsfrei);
- $\tau \in \mathbb{R}_{> 0}$ ist die maximal zulässige Reaktionszeit in Millisekunden ($\text{ms}$), die in Sicherheitsnormen durch das Fehlertoleranzzeitintervall (*Fault Tolerant Time Interval*, FTTI) strikt limitiert ist: $\tau \le \text{FTTI} - \Delta t_{\mathrm{margin}}$ (z. B. $\tau = 50\,\text{ms}$ bei $\text{FTTI} = 100\,\text{ms}$);
- $\Diamond_{\le\tau}$ ist der metrische MTL-Operator „irgendwann innerhalb der nächsten $\tau$ Zeiteinheiten“.

Praktische Anwendung und Engineering-Entscheidungen:
- **HIL-Prüfstandsverifikation:** Im Hardware-in-the-Loop-Test (HIL) wird durch Fehlerinjektion die exakte Latenz $\Delta t_{\mathrm{meas}} = t_{\mathrm{safe}} - t_{\mathrm{fault}}$ ermittelt.
- **Abnahmekriterium:** Ein Test gilt nur dann als bestanden, wenn $\max(\Delta t_{\mathrm{meas}}) \le \tau$. Überschreitet auch nur ein Durchlauf diese Grenze ($\Delta t_{\mathrm{meas}} > \tau$), emittiert das Testsystem den Fehler `E_FTTI_TIMEOUT_BREACH`, und die Anforderung wird als auf der gewählten Hardwareplattform nicht einhaltbar klassifiziert.

```mermaid
flowchart LR
    accTitle: Zeitstrahl für Sicherheitsinvariante und beschränkte Lebendigkeit
    accDescr: Visualisierung temporallogischer Operatoren auf der Zeitachse vom Auftreten des Fehlers bis zum garantierten Erreichen des sicheren Zustands innerhalb der Toleranzzeit Tau.

    subgraph TimeLine["Physikalische Zeitachse (ms)"]
        direction LR
        T0["<b>t = t₀ (Ereignis)</b><br/><code>FaultTriggered == true</code><br/><i>(Übertemperatur, CAN-Timeout)</i>"]
        Delta["<b>Reaktionsfenster Δt ≤ τ</b><br/><i>(Zeitbudget: z. B. 50 ms)</i><br/>Interrupt-Bearbeitung, Diagnose"]
        T_Safe["<b>t ≤ t₀ + τ (Sicherer Zustand)</b><br/><code>SafeStateAchieved == true</code><br/><i>(Schütz geöffnet)</i>"]
        T0 -->|"Operator ◊≤τ"| Delta
        Delta --> T_Safe
    end

    subgraph Invariant["Sicherheitsinvariante (Operator □)"]
        direction TB
        SafeBox["<b>Permanente Invariante:</b><br/><code>□ ¬(Current > I_max ∧ Closed)</code><br/>Gilt an jedem Zeitpunkt t ∈ [0, ∞)"]
    end

    classDef ev fill:#ffebee,stroke:#c62828,stroke-width:2px,color:#b71c1c;
    classDef dt fill:#fff3e0,stroke:#f57c00,stroke-width:2px,color:#e65100;
    classDef ok fill:#e8f5e9,stroke:#388e3c,stroke-width:2px,color:#1b5e20;
    classDef inv fill:#e1f5fe,stroke:#0288d1,stroke-width:2px,color:#01579b;

    class T0 ev;
    class Delta dt;
    class T_Safe ok;
    class SafeBox inv;
```

> [!TIP] Die Falle der vacuosen Wahrheit (Vacuous Truth)
> In der formalen Logik ist die Implikation $A \Rightarrow B$ trivial wahr, wenn die Prämisse $A = \text{False}$ ist (aus Falschem folgt Beliebiges). Wurde in einem Testlauf die Vorbedingung $\mathit{FaultTriggered}$ durch den Teststand niemals stimuliert (weil der Prüfadapter beispielsweise einen Temperatursprung über 100 °C nicht simulieren konnte), meldet ein unbedarfter Verifizierer: *„Anforderung zu 100 % erfolgreich verifiziert!“*. In der Realität wurde der Fehlerbehandlungscode kein einziges Mal ausgeführt. Zur Verhinderung vacuoser Verifikationen muss die Audit-Pipeline jede Implikation mit einer Erreichbarkeitsprüfung der Prämisse koppeln: Es muss nachgewiesen sein, dass $\Diamond\,\Phi_{\mathrm{pre}}$ in mindestens einem Testlauf erfüllt war.

Auf dieser Basis lassen sich die EARS-Muster in logische Ausdrücke überführen. Die nachfolgende Tabelle veranschaulicht didaktische Interpretationen: Die Wahl des Zeitmodells, der Schrittweite, des Nullpunkts und des Quantorenbereichs bleibt stets eine explizite Modellierungsentscheidung.

| EARS-Schablone | Beispiel | Logische Form |
|---|---|---|
| Ubiquitär | Der CAN-Controller unterstützt 29-Bit-Identifier | $\Box\,\mathit{Supports}(\mathit{CAN},\mathit{Ext29})$ |
| Ereignisgesteuert | Nach Betätigung des Not-Aus schaltet der Treiber die Ausgänge innerhalb von 5 ms ab | $\Box\bigl(\mathit{Pressed}\Rightarrow\Diamond_{\le5\,\mathrm{ms}}\,\mathit{Disabled}\bigr)$ |
| Zustandsbasiert | Im PRE-CHARGE-Modus übersteigt der Strom nicht 10 A | $\Box\bigl(\mathit{State}=\mathit{PRECHARGE}\Rightarrow\mathit{Current}\le10\,\mathrm{A}\bigr)$ |
| Reaktion auf unerwünschtes Verhalten | Bei Temperaturen über 65 °C läuft die Pumpe mit 100 % Tastgrad | $\Box\bigl(\mathit{Temp}>65\Rightarrow\mathit{PumpDuty}=1{,}0\bigr)$ |
| Optionale Funktion | Ist ein Watchdog vorhanden, wechselt das WDI-Signal alle 50 ms den Pegel | $\mathit{HasWatchdog}\Rightarrow\Box\,(\mathit{ToggleInterval}=50\,\mathrm{ms})$ |

Die Tabelle illustriert exemplarische Abbildungen und keinen universellen Compiler. Der Verifizierungsingenieur muss Quantoren, Zeitgrenzen, Gleichzeitigkeit, Rücksetzbedingungen und Signalzuordnungen gesondert prüfen. Es muss strikt zwischen der bloßen Existenz einer zulässigen Systemreaktion und der Einhaltung der Anforderung durch sämtliche Ausführungspfade der Implementierung differenziert werden.

### 3.3. Automatisierte Konsistenzprüfung von Anforderungen mittels SMT-Solver

Formeln verifizieren sich nicht selbst; diese Aufgabe übernehmen **SMT-Solver** (*Satisfiability Modulo Theories*). Ein SMT-Solver prüft, ob eine Belegung von Variablen existiert, die ein System von Formeln unter Berücksichtigung von Hintergrundtheorien (wie linearer Arithmetik, Bitvektoren und Datentypen) erfüllt. Ein weit verbreiteter Solver ist Z3 von Leonardo de Moura und Nikolaj Bjørner [[11]](#src-11); als standardisierte Eingabesprache dient SMT-LIB [[12]](#src-12).

Für ein Expertensystem besteht der primäre Nutzen des Solvers im automatisierten Nachweis von Widersprüchen (*Inconsistencies*) zwischen Anforderungen. Betrachten wir zwei Anforderungen an ein Batteriemanagementsystem (*Battery Management System*, BMS): REQ-BMS-042 verlangt, dass bei einer Zelltemperatur über 60 °C das Hauptschütz geöffnet werden muss. REQ-BMS-077 fordert, dass im Ladebetrieb (*CHARGE*) das Schütz geschlossen sein muss. Isoliert betrachtet sind beide Anforderungen plausibel. Das folgende SMT-LIB-Modell prüft, ob beide Anforderungen gleichzeitig erfüllbar sind, wenn im Ladebetrieb eine Überhitzung auftritt. Jede Assertion wird benannt, damit der Solver im Konfliktfall die Ursachen benennen kann.

<details>
<summary>Formales SMT-LIB-Modell</summary>

```lisp
(set-option :produce-unsat-cores true)
(declare-datatype BmsState ((INIT) (STANDBY) (CHARGE) (DISCHARGE) (FAULT)))
(declare-const state BmsState)
(declare-const cell_temp_c Real)
(declare-const contactor_closed Bool)

; REQ-BMS-042 (If): Übersteigt die Zelltemperatur 60 °C, öffnet das BMS das Schütz.
(assert (! (=> (> cell_temp_c 60.0) (not contactor_closed)) :named REQ_BMS_042))

; REQ-BMS-077 (While): Im Zustand CHARGE hält das BMS das Schütz geschlossen.
(assert (! (=> (= state CHARGE) contactor_closed) :named REQ_BMS_077))

; Szenario: Überhitzung während des Ladevorgangs.
(assert (! (and (= state CHARGE) (> cell_temp_c 60.0)) :named SCENARIO))

(check-sat)
(get-unsat-core)
```

</details>

Das Skript kann über das Python-Paket `z3-solver` (`pip install z3-solver`) mittels der Funktion `Z3_eval_smtlib2_string` ausgeführt werden. Der Solver Z3 (Version 5.1.0) liefert:

<details>
<summary>Beispieldaten oder Ausführungsergebnis</summary>

```text
unsat
(REQ_BMS_042 REQ_BMS_077 SCENARIO)
```

</details>

Die Antwort `unsat` belegt die logische Unverträglichkeit der formalisierten Aussagen im definierten Szenario. Der *Unsat Core* liefert eine hinreichende (nicht zwingend minimale) Teilmenge widersprüchlicher Aussagen. Das Ergebnis entscheidet nicht autonom, welche Anforderung modifiziert werden muss. Die Ergänzung der Ladebedingung um die Klausel `(<= cell_temp_c 60.0)` stellt eine fachliche Lösung dar, die der Freigabe durch den Systemverantwortlichen bedarf. Ein nach Korrektur erzieltes `sat` beweist die Erfüllbarkeit des Modells, garantiert jedoch keineswegs die funktionale Gesamtsicherheit des physischen Produkts. Ein Timeout oder der Status `unknown` darf niemals mit `sat` oder `unsat` gleichgesetzt werden.

Dieses Beispiel offenbart die Grenzen formaler Methoden: Der Solver identifiziert Widersprüche ausschließlich innerhalb des formalisierten Modells und exakt so, wie der Ingenieur die Abbildung vorgenommen hat. Korrespondiert die Variable `cell_temp_c` fehlerhaft mit dem physikalischen Sensorregister, bleibt dies für den Solver unsichtbar. Die Formalisierung muss daher selbst einem Peer-Review unterzogen werden; der Beweiswert entsteht erst im Verbund mit dem Herkunftsnachweis über Autor, Modellversion und Anforderungsstand.

> [!NOTE]
> **Vier Invarianten der SMT-Konsistenz- und Erreichbarkeitsprüfung von Anforderungen**
> 1. Symbolische Formelkodierung: Übersetzung temporaler Prädikate in quantifizierte SMT-LIB-Zusicherungen für den Theorembeweiser Z3.
> 2. Widerspruchserkennung (Unsat Core): automatische Extraktion der minimalen inkonsistenten Axiomenmenge bei kollidierenden normativen Bedingungen.
> 3. Schutz vor vacuoser Wahrheit: obligatorische Erreichbarkeitsprüfung der Prämissen ($\Diamond\Phi_{\mathrm{pre}}$) zum Ausschluss fiktiver Validierungen.
> 4. Deterministisches Zeitlimit (Timeout Gate): Überführung des Solvers in einen sicheren Haltezustand bei Überschreitung des Rechenbudgets.
> *Laufzeitsteuerung:* Das Auffinden eines nicht-leeren Unsat Core blockiert sofort die Erstellung der Wissensbasis und erzeugt einen Konfliktdiagnosebericht.

## 4. Metriken und automatisierte Qualitätsprüfung von Anforderungen (Requirements Smells)

Eine Formalisierung setzt qualitativ hochwertige Anforderungen voraus. Der Standard ISO/IEC/IEEE 29148:2018 definiert die Qualitätsmerkmale professioneller Anforderungen: Sie müssen notwendig, eindeutig, vollständig, atomar, realisierbar, prüfbar und korrekt sein [[13]](#src-13). Henning Femmer und Kollegen schlugen vor, Verstöße gegen diese Kriterien automatisiert als Spezifikationsmängel (*Requirements Smells*) zu detektieren: subjektive Formulierungen, vage Adverbien und Adjektive, Schlupflöcher, offene unprüfbare Begriffe, Superlative, uneindeutige Negationen, unklare Pronomina und unvollständige Verweise [[14]](#src-14).

```mermaid
flowchart LR
    accTitle: Kategorien von Spezifikationsmängeln (Requirements Smells)
    accDescr: Anforderungstexte werden auf fünf Mängeltypen geprüft: Ambiguität, mangelnde Prüfbarkeit, Nicht-Atomarität, fehlender Alternativzweig und Passivkonstruktion ohne Akteur.

    Text["Anforderungstext"] --> Amb["<b>Ambiguität</b><br/>vage qualitative Begriffe"]
    Text --> NonTest["<b>Mangelnde Prüfbarkeit</b><br/>keine Toleranzen und Metriken"]
    Text --> NonAtom["<b>Nicht-Atomarität</b><br/>mehrere Aktionen unter einem shall"]
    Text --> MissingElse["<b>Fehlender Alternativzweig</b><br/>keine Reaktion bei Wegfall der Bedingung"]
    Text --> Passive["<b>Passivkonstruktion</b><br/>kein Akteur definiert"]

    classDef root fill:#e1f5fe,stroke:#0288d1,stroke-width:2px,color:#01579b;
    classDef smell fill:#ffebee,stroke:#c62828,stroke-width:2px,color:#b71c1c;

    class Text root;
    class Amb,NonTest,NonAtom,MissingElse,Passive smell;
```

Das Diagramm klassifiziert Mängel nach ihren technischen Konsequenzen: Eine ambige Anforderung wird von Entwicklern unterschiedlich interpretiert, eine unprüfbare lässt sich nicht durch Tests verifizieren, eine nicht-atomare verhindert eine 1:1-Rückverfolgbarkeit auf Testfälle, und ein fehlender Alternativzweig überlässt das Systemverhalten nach einer Störung der Willkür des Programmierers.

**Ambiguität und vage Begriffe.** Wörter wie *fast*, *promptly*, *immediately*, *as soon as possible*, *user-friendly*, *robust*, *adequate*, *approximately* oder *etc.* besitzen keine physikalische Messgröße. Das Expertensystem führt ein Lexikon derartiger Begriffe und fordert quantitative Präzisierungen: Statt *promptly* muss eine Zeitgrenze (z. B. $t\le15\,\text{ms}$), statt *approximately* müssen Nennwert und Toleranzband definiert werden.

**Nicht-Atomarität.** Der Satz *„The gateway shall parse the incoming CAN message, verify the CRC, update the internal state machine, and transmit an acknowledgment frame“* bündelt vier eigenständige Operationen. Schlägt lediglich die CRC-Prüfung fehl, verharrt der Anforderungsstatus in einem undefinierten Zustand („teilweise erfüllt“). Für die lückenlose Nachverfolgbarkeit nach ISO 26262 sind separate Anforderungen mit jeweils eigenen Testfällen zwingend erforderlich.

**Fehlender Alternativzweig.** Die Anforderung *„If the battery temperature exceeds 55 °C, the cooling fan shall turn ON“* spezifiziert nicht, wann der Lüfter wieder deaktiviert wird. Erfolgt die Abschaltung bei exakt denselben 55 °C, führt thermisches Rauschen an der Schaltschwelle zu hochfrequentem Flattern des Relais (*chattering*). Die regelungstechnische Lösung verlangt eine Schalthysterese: Einschalten bei $T>55\,^{\circ}\mathrm{C}$, Ausschalten erst bei $T\le48\,^{\circ}\mathrm{C}$. Ohne korrespondierende Gegenanforderung wählt der Entwickler Schwellwerte nach eigenem Ermessen, was die funktionale Verifikation untergräbt.

**Passivkonstruktionen ohne Akteur.** Wie in Abschnitt 1 dargelegt, kann eine Anforderung ohne explizites grammatikalisches Subjekt keiner Architekturkomponente zugewiesen werden. Das Expertensystem markiert solche Sätze und erzwingt die Benennung der zuständigen Komponente.

Das Qualitätsaudit generiert entweder einen bereinigten Kandidaten für die Formalisierung oder einen Mängelbericht an den Autor. Anzeichen mehrerer Aktionen, eines fehlenden Zweigs oder unbestimmter Fristen stellen Prüfindikatoren dar, keine universellen Fehler: Mehrere Aktionen können eine unteilbare Transaktion bilden, und unerwähntes Verhalten kann bewusst undefiniert bleiben. Die Freigabe für eine Baseline setzt fachliche Begutachtung voraus und nicht allein das Fehlen detektierter Smells.

### 4.1. Qualifikationsgateway für normative Anforderungen auf Basis externer Zulässigkeitskriterien

Bei der Übersetzung und Formalisierung normativer Anforderungen darf sich die linguistische Pipeline nicht auf weiche probabilistische Konfidenzwerte („Klassifikatorkonfidenz 87 %“) verlassen. In sicherheitskritischen Systemen (ISO 26262 ASIL D, DO-178C) muss das Verifikationsergebnis deterministisch und typisiert sein. In Untersuchungen zu SIF 1.0 [[20]](#src-20) wurde ein begrenzter Mechanismus der externen Registry-Kontrolle über kanonische metrologische Termini im Ausgabetext der neuronalen Übersetzung empirisch verifiziert. In Analogie zur architektonischen Trennung von Generierung und unabhängiger Qualifikation schlägt der Autor der Monografie ein eigenes Qualifikationsgateway für normative Anforderungen vor, dessen Regeln und operative Ergebnisse einer gesonderten Begründung bedürfen. Die Audit-Pipeline für Anforderungen fällt eines von vier diskreten operativen Verdikten:

```mermaid
flowchart TD
    accTitle: Vierstufiges Qualifikationsgateway für Anforderungen
    accDescr: Entscheidungsprotokoll des Gateways: REQUIREMENT_ADMITTED, VERIFY_REPAIRED, REVISION_REQUESTED, FAIL_CLOSED_REJECT.

    REQ["Anforderungskandidat<br/>(aus NLP-Pipeline)"] --> GATE{"Qualifikationsgateway<br/>für normative Anforderungen"}
    GATE -->|"Invarianten gewahrt,<br/>Toleranzen gebunden"| RA["REQUIREMENT_ADMITTED<br/>(Zulassung zu SMT / KB)"]
    GATE -->|"Zulässige Verschiebung behoben,<br/>Anker rekonstruiert"| VR["VERIFY_REPAIRED<br/>(erfordert Bestätigung)"]
    GATE -->|"Semantischer Konflikt,<br/>unauflösbare Kollision"| RR["REVISION_REQUESTED<br/>(Überarbeitungsanforderung)"]
    GATE -->|"Toxische Deskriptoren,<br/>Inversion von Verboten (MUST NOT)"| FR["FAIL_CLOSED_REJECT<br/>(Blockierung, Fail-Closed)"]

    classDef ok fill:#e8f5e9,stroke:#388e3c,stroke-width:2px,color:#1b5e20;
    classDef check fill:#fff3e0,stroke:#f57c00,stroke-width:2px,color:#e65100;
    classDef warn fill:#fce4ec,stroke:#c2185b,stroke-width:2px,color:#880e4f;
    classDef err fill:#ffebee,stroke:#d32f2f,stroke-width:2px,color:#b71c1c;

    class RA ok;
    class VR check;
    class RR warn;
    class FR err;
```

Bedeutung der operativen Qualifikationszustände:

1. **`REQUIREMENT_ADMITTED`:** Die syntaktische Struktur entspricht vollständig den EARS-Schablonen, die deontische Modalität stimmt mit der Konvention des Quelldokuments überein, und alle numerischen Toleranzen sowie physikalischen Einheiten sind an registrierte Konstanten und Konfigurationsbereiche gebunden. Die Anforderung kann unter dem Vorbehalt der Erfüllung der definierten Prüfungen an die nachfolgende Formalisierungsstufe übergeben werden. Ein solcher Status bestätigt nicht automatisch die semantische Übereinstimmung des formalen Modells mit der Primärquelle und ersetzt keineswegs die verfahrensmäßig vorgesehene fachliche Freigabe.
2. **`VERIFY_REPAIRED`:** Unkritische Mängel (z. B. geringfügige Wortstellungsanomalien, Passivkonstruktionen bei eindeutigem Kontextsubjekt oder Quasimodalität) wurden während der Normalisierung deterministisch korrigiert. Der Status erfordert die Freigabe durch den zuständigen Prüfingenieur.
3. **`REVISION_REQUESTED`:** Es wurden logische Widersprüche, nicht-atomare Mehrfachaktionen oder fehlende Alternativzweige festgestellt. Der Autor erhält eine strukturierte Diagnosekarte mit der Empfehlung zur Aufteilung der Anforderung oder zur Definition einer Hysterese.
4. **`FAIL_CLOSED_REJECT`:** Vage subjektive Deskriptoren, undefinierte Toleranzen im Sicherheitskreis oder das Risiko einer Inversion von Verboten (`MUST NOT`) wurden detektiert. Es greift die Fail-Closed-Invariante: Das System blockiert die Aufnahme der Regel in die Wissensbasis bedingungslos.

> [!NOTE]
> **Fünf Autoren-Metriken zur Erkennung sprachlicher Anforderungsdefekte (Requirements Smells)**
> 1. Indikator unqualifizierter Qualität: Erkennung subjektiver Superlative ohne numerische Toleranzen des SI-Einheitensystems.
> 2. Indikator für Passivkonstruktionen: Blockierung von Anforderungsklauseln ohne explizit benanntes ausführendes Systemsubjekt.
> 3. Indikator verdeckter Normativität: Identifikation umgangssprachlicher Quasimodalitäten anstelle standardisierter deontischer Operatoren.
> 4. Indikator syntaktischer Überlastung: Abwertung verschachtelter Sätze mit multiplen Nebensätzen, die Ausführungsbedingungen verschleiern.
> 5. Indikator struktureller EARS-Unvollständigkeit: Prüfung auf das Vorhandensein aller Pflichtfelder entsprechend dem gewählten Schablonentyp.
> *Laufzeitsteuerung:* Das Überschreiten toxischer Mängelgrenzwerte (Smells) führt zur Ablehnung der Spezifikation und Rückgabe an den Ingenieur.

## 5. Software-Implementierung: Eine deterministische Analyse-Pipeline in Go

Die vorangegangenen Abschnitte behandelten die Analyseschritte isoliert. Das folgende Programm führt sie zu einer deterministischen Pipeline zusammen: Es bestimmt die Modalität anhand der Dokumentenkonvention, parst den Satz gegen EARS-Schablonen und prüft auf vage Begriffe, Nicht-Atomarität, generische Subjekte sowie fehlende Reaktionsfristen. Die Implementierung stützt sich ausschließlich auf die Go-Standardbibliothek und wird mittels `go run main.go` ausgeführt.

<details>
<summary>Go-Beispiel: Deterministische Pipeline für das Anforderungs-Audit</summary>

```go
package main

import (
	"crypto/sha256"
	"fmt"
	"regexp"
	"strconv"
	"strings"
)

// Convention definiert die Formatierungsregeln des Quelldokuments.
type Convention int

const (
	RFC2119       Convention = iota // BCP 14: normativ nur Wörter in GROSSBUCHSTABEN (RFC 8174)
	ISODirectives                   // ISO/IEC Directives, Part 2: Verbformen in Kleinbuchstaben
)

type Modality string

const (
	Requirement        Modality = "REQUIREMENT"
	Prohibition        Modality = "PROHIBITION"
	Recommendation     Modality = "RECOMMENDATION"
	NotRecommended     Modality = "NOT_RECOMMENDED"
	Permission         Modality = "PERMISSION"
	Capability         Modality = "CAPABILITY"          // ISO: can, cannot
	ExternalConstraint Modality = "EXTERNAL_CONSTRAINT" // ISO: must
	Statement          Modality = "STATEMENT"
	UnknownConvention  Modality = "UNKNOWN_CONVENTION"
)

type rule struct {
	re *regexp.Regexp
	m  Modality
}

// Verneinende Formen werden vor bejahenden geprüft.
var rules = map[Convention][]rule{
	RFC2119: {
		{regexp.MustCompile(`\b(MUST NOT|SHALL NOT)\b`), Prohibition},
		{regexp.MustCompile(`\b(MUST|SHALL|REQUIRED)\b`), Requirement},
		{regexp.MustCompile(`\b(SHOULD NOT|NOT RECOMMENDED)\b`), NotRecommended},
		{regexp.MustCompile(`\b(SHOULD|RECOMMENDED)\b`), Recommendation},
		{regexp.MustCompile(`\b(MAY|OPTIONAL)\b`), Permission},
	},
	ISODirectives: {
		{regexp.MustCompile(`(?i)\bshall not\b`), Prohibition},
		{regexp.MustCompile(`(?i)\bshall\b`), Requirement},
		{regexp.MustCompile(`(?i)\bshould not\b`), NotRecommended},
		{regexp.MustCompile(`(?i)\bshould\b`), Recommendation},
		{regexp.MustCompile(`(?i)\bmay\b`), Permission},
		{regexp.MustCompile(`(?i)\bmust\b`), ExternalConstraint},
		{regexp.MustCompile(`(?i)\bcan(not)?\b`), Capability},
	},
}

func ClassifyModality(text string, c Convention) Modality {
	patterns, known := rules[c]
	if !known {
		return UnknownConvention
	}
	for _, r := range patterns {
		if r.re.MatchString(text) {
			return r.m
		}
	}
	return Statement
}

type EARS string

// EARS-Schablonen (Mavin et al., 2009); erster Treffer gewinnt.
var earsPatterns = []struct {
	kind EARS
	re   *regexp.Regexp
}{
	{"COMPLEX", regexp.MustCompile(`(?i)^While (?P<state>.+?), when (?P<trigger>.+?), the (?P<subject>[\w -]+?) shall (?P<action>.+?)\.?$`)},
	{"UNWANTED_BEHAVIOUR", regexp.MustCompile(`(?i)^If (?P<trigger>.+?), then the (?P<subject>[\w -]+?) shall (?P<action>.+?)\.?$`)},
	{"STATE_DRIVEN", regexp.MustCompile(`(?i)^While (?P<state>.+?), the (?P<subject>[\w -]+?) shall (?P<action>.+?)\.?$`)},
	{"EVENT_DRIVEN", regexp.MustCompile(`(?i)^When (?P<trigger>.+?), the (?P<subject>[\w -]+?) shall (?P<action>.+?)\.?$`)},
	{"OPTIONAL_FEATURE", regexp.MustCompile(`(?i)^Where (?P<feature>.+?), the (?P<subject>[\w -]+?) shall (?P<action>.+?)\.?$`)},
	{"UBIQUITOUS", regexp.MustCompile(`(?i)^The (?P<subject>[\w -]+?) shall (?P<action>.+?)\.?$`)},
}

func ParseEARS(text string) (EARS, map[string]string) {
	for _, p := range earsPatterns {
		m := p.re.FindStringSubmatch(text)
		if m == nil {
			continue
		}
		g := map[string]string{}
		for i, name := range p.re.SubexpNames() {
			if name != "" {
				g[name] = m[i]
			}
		}
		return p.kind, g
	}
	return "NON_CONFORMANT", nil
}

var (
	reTiming   = regexp.MustCompile(`(?i)\b(?:within|no later than|in less than)\s+(\d+(?:\.\d+)?)\s*(ms|s)\b`)
	reVague    = regexp.MustCompile(`(?i)\b(promptly|quickly|as soon as possible|user[- ]friendly|robust|adequate|sufficient(?:ly)?|approximately|etc)\b`)
	reTwoVerbs = regexp.MustCompile(`(?i)\band (?:then )?(open|close|set|send|transmit|notify|flash|update|verify|validate|disable|enable|activate|trigger|log|store)\b`)
)

type Finding struct{ Severity, Code, Detail string }

func Audit(text string, c Convention) (Modality, EARS, float64, []Finding) {
	mod := ClassifyModality(text, c)
	if mod == UnknownConvention {
		return mod, "UNASSESSED", 0, []Finding{{"CRITICAL", "UNKNOWN_CONVENTION", "document convention is not defined"}}
	}
	if mod != Requirement && mod != Prohibition {
		return mod, "", 0, nil // Satz wird nicht als Anforderung geprüft
	}
	kind, g := ParseEARS(text)
	var f []Finding
	if kind == "NON_CONFORMANT" {
		f = append(f, Finding{"CRITICAL", "NON_CONFORMANT_EARS", "no EARS template matches"})
	}
	if w := reVague.FindString(text); w != "" {
		f = append(f, Finding{"CRITICAL", "VAGUE_TERM", fmt.Sprintf("%q is not objectively verifiable", w)})
	}
	if v := reTwoVerbs.FindStringSubmatch(g["action"]); v != nil {
		f = append(f, Finding{"WARNING", "NON_ATOMIC", fmt.Sprintf("second action %q under one shall", v[1])})
	}
	if s := strings.ToLower(g["subject"]); s == "system" || s == "software" {
		f = append(f, Finding{"WARNING", "GENERIC_SUBJECT", fmt.Sprintf("subject %q names no component", s)})
	}
	var ms float64
	if t := reTiming.FindStringSubmatch(g["action"]); t != nil {
		ms, _ = strconv.ParseFloat(t[1], 64)
		if strings.EqualFold(t[2], "s") {
			ms *= 1000
		}
	} else if kind == "EVENT_DRIVEN" || kind == "UNWANTED_BEHAVIOUR" {
		f = append(f, Finding{"CRITICAL", "MISSING_TIME_BOUND", "reaction without a deadline"})
	}
	return mod, kind, ms, f
}

func main() {
	samples := []struct {
		id   string
		conv Convention
		text string
	}{
		{"REQ-1", ISODirectives, "When the cell temperature exceeds 60 °C, the BMS shall open the main contactor within 100 ms."},
		{"REQ-2", ISODirectives, "The system shall promptly validate incoming CAN messages and flash the status LED."},
		{"REQ-3", ISODirectives, "If the CAN bus is lost, then the gateway shall enter the SAFE state."},
		{"REQ-4", ISODirectives, "The installer must follow local electrical codes."},
		{"REQ-5", RFC2119, "The client should retry the connection."},
	}
	for _, s := range samples {
		mod, kind, ms, findings := Audit(s.text, s.conv)
		sum := sha256.Sum256([]byte(s.text))
		if kind == "" {
			fmt.Printf("%s sha256:%x %s (not audited as a requirement)\n", s.id, sum[:4], mod)
			continue
		}
		fmt.Printf("%s sha256:%x %s %s deadline=%gms parser_checks_passed=%v\n",
			s.id, sum[:4], mod, kind, ms, !hasCritical(findings))
		for _, f := range findings {
			fmt.Printf("    [%s] %s: %s\n", f.Severity, f.Code, f.Detail)
		}
	}
}

func hasCritical(fs []Finding) bool {
	for _, f := range fs {
		if f.Severity == "CRITICAL" {
			return true
		}
	}
	return false
}
```

Die Negativtests in `audit_test.go` verifizieren den Gültigkeitsbereich von Reaktionsfristen sowie die Erkennung unbekannter Konventionen; Aufruf: `go test -v main.go audit_test.go`.

```go
package main

import "testing"

func TestDeadlineScopeAndConvention(testCase *testing.T) {
	text := "When the input arrives within 2 s, the gateway shall respond within 100 ms."
	_, _, deadline, findings := Audit(text, ISODirectives)
	if deadline != 100 || hasCritical(findings) {
		testCase.Fatalf("trigger time became action deadline: %g %v", deadline, findings)
	}
	text = "When the input arrives within 2 s, the gateway shall respond."
	_, _, deadline, findings = Audit(text, ISODirectives)
	if deadline != 0 || !hasCritical(findings) {
		testCase.Fatal("trigger time concealed missing response deadline")
	}
	modality, _, _, findings := Audit("The gateway shall respond.", Convention(99))
	if modality != UnknownConvention || !hasCritical(findings) {
		testCase.Fatal("unknown convention silently accepted")
	}
	if ClassifyModality("The client MUST NOT retry.", RFC2119) != Prohibition {
		testCase.Fatal("negative normative marker lost")
	}
}
```

Ausgabe des Programms:

```text
REQ-1 sha256:0d8de9d2 REQUIREMENT EVENT_DRIVEN deadline=100ms parser_checks_passed=true
REQ-2 sha256:96ae8a8f REQUIREMENT UBIQUITOUS deadline=0ms parser_checks_passed=false
    [CRITICAL] VAGUE_TERM: "promptly" is not objectively verifiable
    [WARNING] NON_ATOMIC: second action "flash" under one shall
    [WARNING] GENERIC_SUBJECT: subject "system" names no component
REQ-3 sha256:7ba442f9 REQUIREMENT UNWANTED_BEHAVIOUR deadline=0ms parser_checks_passed=false
    [CRITICAL] MISSING_TIME_BOUND: reaction without a deadline
REQ-4 sha256:611be070 EXTERNAL_CONSTRAINT (not audited as a requirement)
REQ-5 sha256:cdf96e2a STATEMENT (not audited as a requirement)
```

</details>

Die Programmausgabe illustriert die Wirkungsweise der implementierten Heuristiken. REQ-1 weist eine erkannte ereignisgesteuerte Struktur mit einer Frist von 100 ms auf, stellt jedoch noch keine formal freigegebene Anforderung dar. Bei REQ-2 schlagen die Mängelprüfungen für *promptly*, Mehrfachaktionen und die generische Komponentenbezeichnung an. Für REQ-3 meldet die Pipeline das Fehlen einer Reaktionsfrist. REQ-4 wird gemäß den ISO/IEC-Direktiven korrekt als externe Randbedingung identifiziert. In REQ-5 wird das kleingeschriebene *should* gemäß RFC 8174 nicht als normative Vorgabe interpretiert; der normative Charakter des Gesamtsatzes bedarf der Kontextprüfung. Der vorangestellte SHA-256-Hash bindet das Prüfergebnis an den exakten Wortlaut, liefert jedoch keinen Beweis für die inhaltliche Richtigkeit der Klassifikation.

Das Flag `parser_checks_passed` belegt lediglich das Passieren der implementierten Heuristiken, nicht die Reife für eine freigegebene Anforderungs-Baseline. `STATEMENT` signalisiert im Code das Fehlen eines erkannten Schlüsselworts, nicht die bewiesene Absenz von Normativität. Reguläre Ausdrücke können Modalitätsbereiche, Zitate, multiple Fristen oder die Semantik von Begriffen wie „sofort“ nicht auflösen. Zeitgrenzen, Zahlenwerte und Operatoren müssen in einer Zwischenrepräsentation mit eigenem Quellnachweis persistiert werden; geringe Klassifikationssicherheit erfordert ein Review durch den Fachexperten. Keine Häufung erfolgreicher syntaktischer Prüfungen ersetzt die fachliche Freigabe.

## 6. Praxis-Fallstudie: Audit eines Batteriemanagementsystems (BMS) nach ISO 26262

Betrachten wir einen Auszug aus einem Lastenheft für ein Traktionsbatterie-Managementsystem mit dem Sicherheitsintegritätslevel ASIL C nach ISO 26262.

<details>
<summary>Beispieldaten oder Ausführungsergebnis</summary>

```text
Document: SRS_HV_Battery_Management_v2.4.docx
Section: 5.3 Safety Mechanisms and Thermal Runaway Prevention

Paragraph 1:
"If an over-temperature condition (cell temperature > 60°C) is detected by the analog front-end,
the BMS controller shall promptly open the pyrotechnic switch, set the fault register to 0xEF,
and notify the vehicle VCU via CAN message within 50 ms."

Paragraph 2:
"The battery status should be robust and user-friendly under normal driving states."
```

</details>

Der erste Absatz stellt eine verbindliche Anforderung dar, der zweite eine vage Empfehlung mit unpräzisen Begriffen. Das nachfolgende Sequenzdiagramm illustriert den Durchlauf des ersten Absatzes durch die Audit-Pipeline.

```mermaid
sequenceDiagram
    accTitle: Ablauf des Audits einer ASIL-C-Anforderung
    accDescr: Ingestionsmodul, linguistischer Parser, EARS-Klassifizierer und SMT-Formelgenerator verarbeiten die Anforderung sequenziell und erstellen den Auditbericht.

    autonumber
    participant Ingest as Ingestionsmodul
    participant Parser as Linguistischer Parser
    participant EARS as EARS-Klassifizierer
    participant SMT as SMT-Formelgenerator
    participant Audit as Auditbericht

    Ingest->>Parser: Absatz 1 und Quellbindung (SHA-256)
    Parser->>EARS: Modalität: verbindliche Anforderung (shall)
    EARS->>EARS: Schablone: Reaktion auf unerwünschtes Verhalten (If...)
    EARS->>Audit: gefundene Mängel
    Note over Audit: 1. Kritisch: vager Begriff „promptly“<br/>2. Warnung: drei Aktionen unter einem shall
    EARS->>SMT: Konstruktion von Sicherheitsprädikaten
    Audit->>Audit: Fazit: Anforderung nicht für Baseline freigegeben
```

Das Audit liefert vier Kernbefunde:

1. **Quellenbindung.** Das Prüfergebnis referenziert die Datei `SRS_HV_Battery_Management_v2.4.docx`, deren SHA-256-Prüfsumme, Abschnitt 5.3 und Absatz 1.
2. **Dekomposition.** Die Vorbedingung umfasst eine Zelltemperatur über 60 °C sowie deren messtechnische Erfassung. Die geforderten Aktionen lauten: Trennschalter öffnen, Registerwert 0xEF setzen, Fahrzeugsteuergerät benachrichtigen. Der Wirkungsbereich von `within 50 ms` ist ambig: Die Frist kann sich ausschließlich auf die CAN-Nachricht oder auf die gesamte Aktionskette beziehen. Der Parser darf diese Entscheidung nicht eigenmächtig treffen.
3. **Audit-Fazit.** Fristen, Referenzzeitpunkte und Abhängigkeiten zwischen den Teilaktionen müssen präzisiert werden. ISO 26262 definiert das Fehlertoleranzzeitintervall (*Fault Tolerant Time Interval*, FTTI) [[15]](#src-15); der Nachweis einer ausreichenden Reaktionszeit verlangt eine Gefahren- und Risikoanalyse sowie die Berücksichtigung von Sensor- und Aktorverzögerungen. Die syntaktische Prüfung allein erbringt keinen ASIL-C-Konformitätsnachweis.
4. **Korrekturvorschlag.** Der Autor muss den Gültigkeitsbereich der 50-ms-Frist bestätigen und Zeitbudgets für die übrigen Aktionen festlegen. Bis zu dieser Bestätigung verbleibt die Anforderung im Entwurfsstatus. Bei einer Aufspaltung müssen der ursprüngliche Erkennungstrigger, die Schreibaktion (statt einer angenommenen Zustandshaltung) und unbekannte Parameter exakt erhalten bleiben; fiktive Werte wie 10 ms oder Nachrichtenkennungen dürfen nicht hinzuerfunden werden. Eine Dekomposition ist erst zulässig, wenn sichergestellt ist, dass transaktionale und temporale Bindungen nicht verloren gehen.

Für Continuous-Integration-Pipelines wird das Auditergebnis in ein maschinenlesbares JSON-Format serialisiert. Der folgende Auszug zeigt den Aufbau; Statistikfelder sind mit Beispieldaten belegt.

<details>
<summary>Strukturierte JSON-Daten</summary>

```json
{
  "audit_run_id": "audit-run-2026-09-20-0042",
  "source_document": {
    "filename": "SRS_HV_Battery_Management_v2.4.docx",
    "hash_sha256": "8f4c2e17a3b94091d32a0fbc556281e0129a00b21a8f902345bc7981ef40a012",
    "section": "5.3 Safety Mechanisms and Thermal Runaway Prevention"
  },
  "verdict": "REJECTED_FOR_BASELINE",
  "statistics": {
    "total_sentences_scanned": 142,
    "requirements_detected": 38,
    "ears_compliant_ratio": 0.763,
    "critical_smells_count": 3,
    "warnings_count": 8
  },
  "findings": [
    {
      "requirement_id": "DRAFT_REQ_042",
      "modality": "REQUIREMENT",
      "convention": "ISO_IEC_DIRECTIVES_PART2",
      "ears_pattern": "UNWANTED_BEHAVIOUR",
      "formal_invariants": {
        "fault_condition": "cell_temperature > 60.0",
        "action_predicates": [
          "open(pyro_switch)",
          "set(fault_register, 0xEF)",
          "broadcast_can(VCU_ALARM)"
        ]
      },
      "blocking_issues": [
        {
          "code": "VAGUE_TERM",
          "severity": "CRITICAL",
          "target_token": "promptly",
          "explanation": "A qualitative time bound cannot be verified against the FTTI."
        },
        {
          "code": "NON_ATOMIC",
          "severity": "WARNING",
          "target_token": "and",
          "explanation": "Three actions need separate traceability identifiers."
        }
      ]
    }
  ]
}
```

</details>

Die JSON-Struktur spiegelt die Pipeline-Schritte exakt wider: Quelldokument mit kryptografischem Hash, Modalität mit Dokumentenkonvention, EARS-Muster, formale Prädikate sowie Mängel, die eine Übernahme in die Baseline blockieren. Ein derartiges Artefakt lässt sich programmatisch weiterverarbeiten, während Prüfer die Herleitung jeder Entscheidung bis zum Quelltext auditieren können.

## 7. Intermediäre Repräsentation von Anforderungen (IR): Vereinheitlichter Verifizierer-Vertrag

Der direkte Übergang vom Satz zur logischen Formel birgt das Risiko stillschweigend hinzugedichteter Annahmen. Eine intermediäre Repräsentation (*Intermediate Representation*, IR) macht jede Entwurfsentscheidung explizit sichtbar, bevor der Solver aufgerufen wird. Für eine präzise Formalisierung genügen Zahlenwerte allein nicht; erforderlich sind semantische Rollen und Quelltextpositionen der Prämissen.

| IR-Felder | Verifikation vor der Kompilierung |
|---|---|
| Akteur, Aktion, Objekt | Zulässige Typen und Identifikatoren; Textfragment für jede semantische Rolle |
| Modalität, Konvention, Abschnittsrolle | Gebot, Verbot, Empfehlung oder Erlaubnis; Zitate werden nicht automatisch zu Anforderungen |
| Vorbedingung, Trigger, Ausnahme | Getrennte Bedingungsbäume; definierter Skopus jeder Bedingung |
| Quantor und Negation | „für alle“, „existiert“, „nicht“ mit explizit festgelegtem Wirkungsbereich |
| Operator, Zahlenwert, Einheit | Striktheit von Ungleichungen, Dimensionen, Originalschreibweise und Konvertierungsregeln |
| Zeitursprung, Frist, Aktionsfolge | Zeit der Vorbedingung wird nicht zur Reaktionsfrist; unbekannte Reihenfolge bleibt unbestimmt |
| Quelle, Revisionsstand, Status | Byte-Offsets, Anwendungsbereich und explizite Freigabe der Interpretation |

Der Compiler arbeitet **partiell**: Er verarbeitet ausschließlich unterstützte Formulierungen und generiert für alle übrigen offene Verifikationsverpflichtungen. Ist unklar, ob eine Frist von 50 ms sämtliche Aktionen umfasst, erzeugt der Compiler keine willkürliche temporale Formel. Das beschränkte Decodieren aus [Kapitel 13](ch13-language-variability-vs-determinism.md) garantiert zwar die grammatikalische Struktur des IR-Kandidaten, vermag jedoch fachliche Lücken nicht zu schließen.

Für Zahlenwerte und Einheiten eignen sich Bibliotheken wie Pint [[16]](#src-16) mit domänenspezifischen Validierungstests. Die Umrechnung von Grad Celsius in Kelvin erfordert einen Offset und nicht nur einen Skalierungsfaktor; *Octet* bezeichnet exakt acht Bit, während die Bitbreite von *Byte* plattformabhängig definiert sein kann. Das Dezimalpräfix Mega (M) darf nicht unbesehen durch das Binärpräfix Mebi (Mi) ersetzt werden. „Kleiner als“ und „nicht größer als“ dürfen nicht auf denselben Operator abgebildet werden. Die normalisierte Größe wird stets zusammen mit der ursprünglichen Notation und dem Quelltextfragment gespeichert.

Z3 verifiziert das formale Modell, nicht den Prosatext. Für eine gewählte Logik-Teilmenge empfiehlt sich ein Kreuzvergleich von Z3 mit cvc5 [[17]](#src-17), indem SMT-LIB-Logik, Versionen, Ressourcenlimits und Soll-Ergebnisse fixiert werden. Die Übereinstimmung zweier Solver heilt jedoch keine fehlerhafte Formalisierung. Erforderlich sind dedizierte Testfälle für unerreichbare Trigger, Grenzzeitwerte, Widersprüche und unvollständig belegte IR-Felder.

Das Ergebnis dieser Zwischenrepräsentation ist ein verifizierbarer Formelkandidat mit transparenten Grenzen. Exakt dieses Artefakt kann – verknüpft mit Herkunftsnachweisen und Genehmigungsvermerken – in die Konstruktion der Wissensbasis überführt werden.

> [!NOTE]
> **Vier Invarianten des einheitlichen Anforderungs-Zwischenrepräsentations-Vertrags (IR Contract)**
> 1. Kanonische Entitätstypisierung: Bindung von Subjekt, deontischer Modalität, Prädikat und Bedingungstriggern in einer stabilen AST-Struktur.
> 2. Byte-Kustodie der Primärquelle: zwingende Erfassung exakter `[byte_start..byte_end]`-Offsets und des SHA-256-Kryptohashes des Quellenzitats.
> 3. Isolation von der Quellsyntax: vollständige semantische Unabhängigkeit der Zwischenrepräsentation von der ursprünglichen natürlichen Sprache.
> 4. Deterministische Übersetzung in formale Regeln: direkte Synthese von EVM-Mikrocode und SMT-Zusicherungen ohne probabilistische Modelle.
> *Laufzeitsteuerung:* Jede Zwischenrepräsentation ohne validierte Byte-Kustodie wird durch das Wissensbasis-Integritätsgateway zurückgewiesen.

## 8. Dialog Pro et Contra

> **Epistemische Präambel und akademischer Disclaimer:**  
> Dieser Abschnitt modelliert eine virtuelle intellektuelle Debatte, um die architektonischen Methoden des Autors zur automatisierten Erkennung, zum semantischen Linting und zur Formalisierung normativer technischer Anforderungen einer kompromisslosen wissenschaftlichen Überprüfung zu unterziehen. Zu diesem Fachdialog laden wir international renommierte Pioniere des Requirements Engineerings und formaler Spezifikationen ein: den Schöpfer der EARS-Methodik (Easy Approach to Requirements Syntax) **Alistair Mavin**, dessen strukturierte Muster zum Industriestandard in Luftfahrt und Automobiltechnik wurden, sowie die herausragende Forscherin für formale Methoden und ACM-Fellow **Pamela Zave**, Mitbegründerin des grundlegenden Vier-Variablen-Modells (Four-Variable Model) für Anforderungen an eingebettete Systeme.  
> *Disclaimer:* Diese Diskussion ist eine simulierte heuristische Synthese, die ausschließlich im wissenschaftlich-pädagogischen Rahmen dieser Monografie auf Basis der veröffentlichten Werke, Monografien und Vorträge der eingeladenen Forscher generiert wurde. Sie stellt weder ein wörtliches Zitat noch eine reale Äußerung lebender Personen dar, sondern dient als Instrument des kritischen Rationalismus zur Falsifikation der Hypothesen des Autors. Jeder Gesprächspartner debattiert direkt und ausschließlich mit dem Autor der Monografie (Mykola Fedchyk).

**These des Kapitels (Mykola Fedchyk):**  
Die Transformation unstrukturierter technischer Standards in verifizierbare mathematische Invarianten gelingt durch syntaktische Kanalisierung natürlicher Sprache mittels EARS-Schablonen, semantisches Linting vager Adjektive unter Erzwingung physikalischer SI-Toleranzen und grammatikgesteuerte Kompilierung in geschlossene SMT-LIB-Formeln zur automatisierten Erkennung von Spezifikationsblockaden mittels Z3.

**Dialog mit Alistair Mavin: Ausdruckskraft von EARS-Syntaxmustern in komplexen normativen Systemen**

- **Alistair Mavin:**  
  *„Mykola Fedchyk, als wir die EARS-Methodik bei Rolls-Royce entwickelten, war unser primäres Ziel, technische Anforderungen von der Mehrdeutigkeit der natürlichen Sprache durch intuitive Strukturmuster (Ubiquitous, Event-driven, State-driven usw.) zu befreien. EARS wurde jedoch als praktische Kommunikationsbrücke zwischen Menschen — Anforderungsautoren und Systemingenieuren — konzipiert. Der Versuch, die Extraktion von Anforderungen aus beliebigen Standards vollständig zu automatisieren und vielschichtige Prozesse der funktionalen Sicherheit in starre Formelkäfige zu zwingen, birgt das Risiko, implizite Annahmen abzuschneiden, die Ingenieure im Kopf tragen, aber selten in Schablonenslots explizieren.“*

- **Mykola Fedchyk:**  
  *„Herr Alistair Mavin, Ihr Beitrag zur Standardisierung von Anforderungen ist unschätzbar, aber in sicherheitskritischen Systemen nach ASIL D oder DO-178C sind implizite Annahmen von Ingenieuren die tödlichste Ursache katastrophaler Ausfälle!  
  Wenn eine technische Anforderung nicht in eine disziplinierte EARS-Schablone mit messbaren Triggern, harten Reaktionsfristen (`innerhalb von X ms`) und erschöpfenden Alternativzweigen passt, ist sie keine Anforderung, sondern toxische Mehrdeutigkeit.  
  In meiner Architektur klassifiziert die Go-Audit-Pipeline nicht bloß deontische Modalitäten (`SHALL`, `MUST NOT`), sondern erzwingt vollständige strukturelle Integrität nach EARS. Formuliert eine Norm 'das System muss unverzüglich reagieren', erzeugt unser semantischer Linter sofort den Fehler `MISSING_TIME_BOUND`. Wir schneiden keinen Kontext ab — wir zwingen Fachexperten dazu, Annahmen in verifizierbare physikalische Invarianten zu überführen.“*

**Fortsetzung des Dialogs mit Alistair Mavin: Auflösung ebenenübergreifender Konflikte in hierarchischen EARS-Spezifikationen**

- **Alistair Mavin:**  
  *„In komplexen Großsystemen verteilen sich EARS-Anforderungen über mehrere Architekturschichten: System-EARS, Teilsystem-EARS und Software-EARS. Häufig fügt sich eine lokale Sicherheitsanforderung eines Subsystems perfekt in ein EARS-Muster ein, kollidiert jedoch diametral mit einer übergeordneten Systemverfügbarkeits-Anforderung. Wie können Ihr syntaktischer Linter und die Formalisierungs-Engine solche vertikalen Hierarchiekonflikte automatisiert aufdecken und harmonisieren, ohne dass Systemarchitekten permanent manuell eingreifen müssen?“*

- **Mykola Fedchyk:**  
  *„Herr Mavin, wir lösen dies durch eine deterministische lexikografische Rangordnung $\mathrm{Rank}(\mathcal{R})$ und die automatisierte Extraktion minimaler unerfüllbarer Kerne (*Minimal Unsatisfiable Cores*, MUC) in Z3.  
  Sicherheitsanforderungen höchster Integrität (ASIL D / DO-178C Level A) dominieren funktionale Bequemlichkeits- oder Energieoptimierungsziele strikt.  
  Wenn die Pipeline eine ebenenübergreifende Konjunktion in eine SMT-LIB-Formel übersetzt und der Solver `unsat` meldet, erzeugt das System nicht bloß einen Abbruch, sondern lokalisiert das exakte konfliktverursachende Erdungstupel $\mathcal{G}$ über beide Spezifikationsebenen und konstruiert einen Gegenbeispiel-Trace (*Counter-Example Trace*): 'Sicherheitsanforderung SYS-REQ-104 erzwingt die Aktorabschaltung binnen $15\,\text{ms}$, während funktionale Anforderung SW-REQ-210 eine Diagnoseeinschwingzeit von $40\,\text{ms}$ verlangt'. Dies liefert Architekten eine mathematisch exakte Diagnose anstelle wochenlanger Abstimmungsrunden.“*

**Dialog mit Pamela Zave: Formalisierung in SMT-Prädikaten und Vollständigkeit des Four-Variable-Domänenmodells**

- **Pamela Zave:**  
  *„Mykola Fedchyk, die Kompilierung von EARS-Anforderungen in SMT-LIB-Prädikate und die Widerspruchssuche via Z3-Solver ist faszinierend, birgt jedoch eine klassische Falle formaler Methoden. In unserem Vier-Variablen-Modell (Jackson & Zave) erfordert Korrektheit den Beweis von $S, W \vdash R$: Die Spezifikationsverhaltensweisen $S$ müssen zusammen mit den physikalischen Umweltannahmen $W$ die Anforderungen $R$ erfüllen. Wenn Sie den Text $R$ isoliert formalisieren, ohne eine vollständige mathematische Theorie des physikalischen Umfelds $W$ (wie elektrochemische Übergangsprozesse einer Batteriezelle) einzubinden, liefert ein SMT-`sat` von Z3 eine trügerische und gefährliche Illusion von Sicherheit.“*

- **Mykola Fedchyk:**  
  *„Frau Dr. Pamela Zave, ich teile Ihre fundamentale Erkenntnis über die Unverzichtbarkeit des Umgebungsmodells $W$ uneingeschränkt. Genau deshalb evaluiert unser Requirements-IR Prädikate niemals isoliert, sondern als Randbedingungen über geschlossenen Zustandsvektoren.  
  Erstens formalisiert unsere Kompilierungs-Pipeline Umweltannahmen $W$ über explizite physikalische SI-Intervalle (z. B. $T \in [-40, +85]^\circ\text{C}$, $V \in [2.5, 4.2]\,\text{V}$).  
  Zweitens zielt die SMT-Prüfung nicht auf ein bequemes `sat`, sondern auf das Auffinden unerfüllbarer Kerne (`unsat core`): Führt die Schnittmenge zweier Anforderungen unter gültigen Umgebungsbedingungen $W$ zu einem logischen Widerspruch (wie im Fallbeispiel der gleichzeitigen Notkühlung und Schützabschaltung), stoppt die Pipeline die Architekturfreigabe per Fail-Closed vor jeder Codeerstellung. SMT ersetzt nicht die Physik, sondern dient als beweisbarer Schutzwall gegen Spezifikationskonflikte.“*

**Fortsetzung des Dialogs mit Pamela Zave: Hardwareverzerrungen, Sensorquantisierung und Zeitverzögerungen bei der Anforderungsverifikation**

- **Pamela Zave:**  
  *„Im Vier-Variablen-Modell stützt sich die Brücke zwischen kontrollierten Variablen $C$ und beobachteten Variablen $M$ auf imperfekte physikalische Eingabe- $I$ und Ausgabegeräte $O$. Reale Sensoren weisen Temperaturdrift, ADC-Quantisierungsrauschen und Busübertragungslatenzen auf. Wenn Ihre SMT-Übersetzung mit idealisierten mathematischen Variablen $x \in \mathbb{R}$ rechnet, wird Ihre verifizierte Spezifikation auf realem Silizium versagen. Kann Ihr SMT-Modell Abtastverzögerungen $\tau_{\mathrm{delay}}$ und Sensortoleranzen $\pm \delta$ berücksichtigen, ohne den Solver in einen kombinatorischen Kollaps zu treiben?“*

- **Mykola Fedchyk:**  
  *„Frau Dr. Zave, genau um diese Diskrepanz zu schließen, formulieren wir Umweltbeschränkungen in quantorenfreier linearer reeller Arithmetik (QF_LRA) und Intervallanalysen.  
  Jede überwachte Größe $m(t)$ wird nicht als idealer Skalar modelliert, sondern als Unsicherheitsintervall auf Basis realer Registerdaten aus CMSIS-SVD: $m(t) \in [s(t) - \delta_{\mathrm{adc}}, s(t) + \delta_{\mathrm{adc}}]$ mit beschränkter Abtastlatenz $\tau_{\mathrm{sample}} \le \Delta t_{\mathrm{bus}}$.  
  Sicherheitsanforderungen werden so verifiziert, dass die Invariante auch unter ungünstigstem Rauschen (*Worst-Case Noise Margin*) garantiert wahr bleibt. Z3 beweist in polynomieller Zeit, dass selbst an den Extremgrenzen von Hardwarefehlern und Übertragungslatenzen das Software-Regelgesetz $S$ die Anforderung $R$ niemals verletzt. Damit verschmilzt Ihr Vier-Variablen-Modell direkt mit der industriellen eingebetteten Mikroelektronik.“*

**Synthese und wissenschaftliches Urteil des Kapitels:**  
Die deterministische Erkennung und Formalisierung normativer Anforderungen schafft durchgängige Nachverfolgbarkeit vom Standard-Byteoffset bis zur Ausführungslogik: EARS-Muster beseitigen sprachliche Unschärfen, semantisches Linting verbietet unbestimmte Adjektive ohne SI-Einheiten, und die Transformation in SMT-LIB garantiert die mathematische Widerspruchsfreiheit technischer Spezifikationen vor Implementierungsbeginn.

## Fazit

1. **Abgrenzungskriterien und Modalkonventionen:** Kandidaten für normative Anforderungen werden anhand der Konventionen und des Kontexts des Quelldokuments (RFC 2119, ISO/IEC-Direktiven) identifiziert, in eine typisierte Zwischenrepräsentation überführt, auditiert und erst anschließend kompiliert. Unbekannte Felder und ambige zeitliche Wirkungsbereiche blockieren die Übersetzung, anstatt durch Mutmaßungen gefüllt zu werden.
2. **Semantisches Linting und deontische Vorzugsordnung:** Der Linter für semantische Mehrdeutigkeit sperrt toxische, vage Deskriptoren und erzwingt physikalische SI-Toleranzen, während die lexikografische Vorzugsordnung $\mathrm{Rank}(\mathcal{R})$ normative Normenkollisionen deterministisch zugunsten höherrangiger Quellen sowie imperativer Verbote und Gebote auflöst ($\text{PROHIBITION} \succ \text{PERMISSION}$).
3. **Unteilbarkeit modaler Quantoren und Span-Erhaltung:** Der Einsatz zusammengesetzter deontischer Atome (`MUST NOT`) und eines Spans erhaltenden Tokenizers `GroundedToken` eliminiert das Risiko fataler Polaritätsinversionen von Verboten durch Zeilenumbrüche, sichert die byteweise SHA-256-Integrität der Beweiskette zu Primärquellen und schließt tote Winkel an Nicht-ASCII-Interpunktionsgrenzen via UTF-8-Analyse aus.
4. **Neuro-symbolische Gültigkeitsschranke (GBNF):** Die grammatikgesteuerte Decodierung auf Basis eines geschlossenen Lexikons $\mathcal{V}_{\mathrm{valid}}$ und der EARS-Syntax isoliert die heuristische Relevanzzone strikt von der beweisbaren Validitätszone, maskiert unzulässige Logits mit dem Wert $-\infty$ und liefert direkte kontrafaktische Erklärungen für Ablehnungen über den Beweisgraphen ohne generative Halluzinationen.
5. **Synchronisation der Formalisierung mit SMT-Solvern:** Die Überführung strukturierter EARS-Anforderungen in Prädikatenlogik erster Stufe, LTL und SMT-LIB ermöglicht das automatisierte Aufdecken verborgener Widersprüche, Deadlocks und fehlender Zweige mittels Z3 noch vor dem Verfassen der ersten Codezeile.
6. **Praktische Verifizierbarkeit der Pipeline:** Eine deterministische Go-Pipeline gewährleistet ein lückenloses End-to-End-Audit von Anforderungen, validiert die Wahrung von Trigger- und Reaktionsfristen sowie Sicherheitsinvarianten und erstellt prüffähige Auditberichte mit direkter Rückverfolgbarkeit zu den Quelltext-Bytes.

---

## Fragen zur Selbstprüfung

1. Warum muss der Modalitätsklassifikator die Konvention des Quelldokuments kennen? Nennen Sie ein Beispiel für ein Wort, das in RFCs und ISO-Normen eine unterschiedliche normative Verbindlichkeit besitzt.
2. Welche drei syntaktischen Fallstricke übersieht eine Suche nach dem Wort SHALL mittels regulärer Ausdrücke und wie deckt das Expertensystem diese auf?
3. Welche EARS-Schablone eignet sich für die Anforderung „Im Ladebetrieb darf der Strom 20 A nicht überschreiten“ und welche logische Formel entspricht ihr?
4. Was bedeutet die Solver-Antwort `unsat` zusammen mit dem *Unsat Core* im Batteriebeispiel und warum beweist die Solver-Antwort allein noch nicht die Korrektheit der Anforderungen?
5. Warum ist eine Anforderung ohne Alternativzweig bei Relaissteuerungen gefährlich und wie löst eine Schalthysterese dieses Problem?
6. Wie blockiert der Linter für semantische Mehrdeutigkeit vage Adjektive und wie löst die lexikografische Vorzugsordnung Kollisionen zwischen Normen unterschiedlichen Rangs?
7. Warum kann ein Zeilenumbruch zwischen `MUST` und `NOT` zu einer katastrophalen Polaritätsinversion der Anforderung in ASIL-D-Systemen führen und wie sichern `GroundedToken`-Strukturen die Beweiskette (Custody)?
8. Auf welche Weise errichten grammatikgesteuerte Decodierung (GBNF) und das geschlossene Lexikon $\mathcal{V}_{\mathrm{valid}}$ eine undurchdringliche Schranke zwischen der heuristischen Relevanzzone und der symbolischen Validitätszone?

---

## Glossar

| Deutscher Begriff | Englisches Äquivalent | Kurzbeschreibung |
|---|---|---|
| Modalität | *modality* | Normative Verbindlichkeit eines Satzes: Anforderung, Verbot, Empfehlung, Erlaubnis, Fähigkeit oder Feststellung |
| Dokumentenkonvention | *drafting convention* | Redaktionelle Richtlinien für den Gebrauch von Modalverben, z. B. RFC 2119 oder ISO/IEC-Direktiven |
| Quasimodale Konstruktion | *quasi-modal construction* | Ausdruck, der ohne das Wort *shall* verpflichtet, z. B. *is responsible for* |
| EARS-Schablone | *EARS template* | Standardisierte Satzform des *Easy Approach to Requirements Syntax* |
| Invariante | *invariant* | Bedingung, die in allen zulässigen Systemzuständen zwingend erfüllt sein muss |
| Sicherheitsinvariante | *safety invariant* | Zusicherung, dass ein kritischer Fehlerzustand niemals eintritt |
| Beschränkte Lebendigkeit | *bounded liveness* | Zusicherung, dass ein gefordertes Ereignis spätestens nach Ablauf einer definierten Frist eintritt |
| Temporale Logik | *temporal logic* | Modallogik mit Operatoren wie „immer“, „irgendwann“ und expliziten Zeitgrenzen |
| SMT-Solver | *SMT solver* | Programm zur Prüfung der Erfüllbarkeit von Formeln unter Berücksichtigung arithmetischer Theorien und Datentypen |
| Unsat Core | *unsat core* | Teilmenge von Aussagen, die für die Inkonsistenz hinreichend ist; Minimalität erfordert gesonderte Prüfung |
| Anforderungsmangel | *requirement smell* | Syntaktischer oder semantischer Indikator für einen potenziellen Formulierungsfehler in einer Anforderung |
| Nicht-atomare Anforderung | *non-atomic requirement* | Anforderung, die mehrere voneinander unabhängige Aktionen zusammenfasst |
| Hysterese | *hysteresis* | Unterschiedliche Schwellwerte für Ein- und Ausschalten zur Vermeidung von hochfrequentem Flattern |
| Anforderungs-Baseline | *requirements baseline* | Freigegebener Stand von Anforderungen als verbindliche Basis für Entwicklung und Traceability |
| Fehlertoleranzzeitintervall | *fault tolerant time interval* | Kürzeste Zeitspanne vom Auftreten eines Fehlers bis zum Eintritt eines potenziellen Schadens ohne Sicherheitsmaßnahme |
| Linter für semantische Mehrdeutigkeit | *ambiguity linter* | Automatisierter Anforderungsfilter auf Basis eines Lexikons toxischer Deskriptoren und relativer Adjektive |
| Lexikografische Vorzugsordnung | *lexicographical preference ordering* | Mathematischer Formalismus zur deterministischen Auflösung deontischer Kollisionen nach Quellen- und Modalitätsrängen |
| Zusammengesetztes deontisches Atom | *composite deontic atom* | Unteilbare lexikalische Einheit eines Verbotsquantors (z. B. `MUST NOT`), resistent gegen Zeilenumbrüche |
| Spans erhaltender Tokenizer | *span-preserving tokenizer* | Lexikalischer Scanner, der absolute physische Byte-Offsets der Primärquelle zur Integritätssicherung des Custody-Hashes fixiert |
| Grammatikgesteuerte Decodierung GBNF | *GBNF grammar-constrained decoding* | Mechanismus zur dynamischen Maskierung von Sprachmodell-Logits anhand einer kontextfreien EARS-Grammatik |

---

## Abkürzungen

| Abkürzung | Vollständige Bezeichnung | Bedeutung |
|---|---|---|
| ASIL | Automotive Safety Integrity Level | Sicherheitsintegritätslevel für Kraftfahrzeuge nach ISO 26262 |
| AST | Abstract Syntax Tree | Abstrakter Syntaxbaum eines Programms oder einer Anforderung |
| BMS | Battery Management System | Batteriemanagementsystem |
| CAN | Controller Area Network | Feldbus zur Kommunikation zwischen Steuergeräten im Fahrzeug |
| CRC | Cyclic Redundancy Check | Zyklische Redundanzprüfung (Prüfsumme) |
| EARS | Easy Approach to Requirements Syntax | Kontrollierte Grammatikschablonen für Anforderungen |
| ECU | Electronic Control Unit | Elektronisches Steuergerät |
| FTTI | Fault Tolerant Time Interval | Fehlertoleranzzeitintervall nach ISO 26262 |
| GBNF | GGML / GGUF Backus-Naur Form | Kontextfreies Grammatikformat zur Einschränkung des Token-Abtastbereichs |
| IR | Intermediate Representation | Typisierte intermediäre Repräsentation einer Anforderung |
| LLM | Large Language Model | Großes Sprachmodell |
| LTL | Linear Temporal Logic | Lineare temporale Logik |
| MTL | Metric Temporal Logic | Metrische temporale Logik |
| PWM | Pulse-Width Modulation | Pulsweitenmodulation |
| ReqIF | Requirements Interchange Format | Standardisiertes XML-Austauschformat für Anforderungen |
| RFC | Request for Comments | Dokumentenreihe von Standards und Richtlinien der Internet Engineering Task Force |
| SHA-256 | Secure Hash Algorithm, 256 bits | Kryptografische Hashfunktion |
| SI | Système International d'Unités | Internationales Einheitensystem für physikalische Größen |
| SLM | Small Language Model | Kompaktes spezialisiertes Sprachmodell für das lokale Deployment |
| SMT | Satisfiability Modulo Theories | Erfüllbarkeit modulo Theorien |
| STL | Signal Temporal Logic | Signal-Temporallogik für kontinuierliche Signale |
| VCU | Vehicle Control Unit | Übergeordnetes Fahrzeugsteuergerät |

---

## Literaturverzeichnis

1. <a id="src-1"></a>Object Management Group. [*Requirements Interchange Format (ReqIF), Version 1.2*](https://www.omg.org/spec/ReqIF/1.2/About-ReqIF). OMG, 2016.
2. <a id="src-2"></a>ISO. [*ISO 26262-8:2018. Road vehicles: Functional safety: Part 8: Supporting processes*](https://www.iso.org/standard/68390.html). 2018.
3. <a id="src-3"></a>RTCA. [*DO-178C: Software Considerations in Airborne Systems and Equipment Certification*](https://www.rtca.org/do-178/). RTCA, 2011.
4. <a id="src-4"></a>S. Bradner. [*RFC 2119: Key Words for Use in RFCs to Indicate Requirement Levels*](https://www.rfc-editor.org/rfc/rfc2119). IETF, 1997.
5. <a id="src-5"></a>B. Leiba. [*RFC 8174: Ambiguity of Uppercase vs Lowercase in RFC 2119 Key Words*](https://www.rfc-editor.org/rfc/rfc8174). IETF, 2017.
6. <a id="src-6"></a>ISO, IEC. [*ISO/IEC Directives, Part 2: Principles and Rules for the Structure and Drafting of ISO and IEC Documents*](https://www.iso.org/sites/directives/current/part2/index.xhtml). 9th edition, 2021.
7. <a id="src-7"></a>Alistair Mavin, Philip Wilkinson, Adrian Harwood, Mark Novak. [*Easy Approach to Requirements Syntax (EARS)*](https://doi.org/10.1109/RE.2009.9). *2009 17th IEEE International Requirements Engineering Conference*, 317–322, 2009.
8. <a id="src-8"></a>Amir Pnueli. [*The Temporal Logic of Programs*](https://doi.org/10.1109/SFCS.1977.32). *18th Annual Symposium on Foundations of Computer Science*, 46–57, 1977.
9. <a id="src-9"></a>Ron Koymans. [*Specifying Real-Time Properties with Metric Temporal Logic*](https://doi.org/10.1007/BF01995674). *Real-Time Systems*, 2(4), 255–299, 1990.
10. <a id="src-10"></a>Oded Maler, Dejan Nickovic. [*Monitoring Temporal Properties of Continuous Signals*](https://doi.org/10.1007/978-3-540-30206-3_12). *Formal Techniques, Modelling and Analysis of Timed and Fault-Tolerant Systems*, LNCS, 152–166, 2004.
11. <a id="src-11"></a>Leonardo de Moura, Nikolaj Bjørner. [*Z3: An Efficient SMT Solver*](https://doi.org/10.1007/978-3-540-78800-3_24). *Tools and Algorithms for the Construction and Analysis of Systems (TACAS)*, LNCS, 337–340, 2008.
12. <a id="src-12"></a>Clark Barrett, Pascal Fontaine, Cesare Tinelli. [*The SMT-LIB Standard: Version 2.6*](https://smtlib.cs.uiowa.edu/papers/smt-lib-reference-v2.6-r2021-05-12.pdf). 2021.
13. <a id="src-13"></a>ISO, IEC, IEEE. [*ISO/IEC/IEEE 29148:2018. Systems and Software Engineering: Life Cycle Processes: Requirements Engineering*](https://www.iso.org/standard/72089.html). 2018.
14. <a id="src-14"></a>Henning Femmer, Daniel Méndez Fernández, Stefan Wagner, Sebastian Eder. [*Rapid Quality Assurance with Requirements Smells*](https://doi.org/10.1016/j.jss.2016.02.047). *Journal of Systems and Software*, 123, 190–213, 2017.
15. <a id="src-15"></a>ISO. [*ISO 26262-1:2018. Road vehicles: Functional safety: Part 1: Vocabulary*](https://www.iso.org/standard/68383.html). 2018.
16. <a id="src-16"></a>Mitwirkende an Pint. [*Non-Multiplicative Units*](https://pint.readthedocs.io/en/stable/user/nonmult.html). Dokumentation zu Temperaturskalen und Offsets.
17. <a id="src-17"></a>Mitwirkende an cvc5. [*cvc5 Documentation*](https://cvc5.github.io/docs/latest/). Dokumentation des Solvers und der unterstützten Logiken.
18. <a id="src-18"></a>G. Gerganov et al. [*GBNF (GGML BNF) Grammar Specification and Constrained Sampling*](https://github.com/ggerganov/llama.cpp/blob/master/grammars/README.md). llama.cpp repository, 2023.
19. <a id="src-19"></a>B. T. Willard, R. Louf. [*Efficient Guided Generation for Large Language Models*](https://arxiv.org/abs/2307.09702). arXiv:2307.09702, 2023.
20. <a id="src-20"></a>Yurii Kuzmenko, Olena Trofimova. [*Structural Invariant Framework (SIF 1.0): A Registry-Based Framework for Empirically Verified Terminology Preservation in Neural EN→UA Translation of Metrological Standards*](https://doi.org/10.5281/zenodo.20509698). Zenodo, 2026.

---

[← Kapitel 12](ch12-linguistic-analysis-and-local-models.md) | [Inhaltsverzeichnis](README.md) | [Teil III](part-03-knowledge-engineering-nlp.md) | [Kapitel 15 →](ch15-knowledge-extraction-and-kb-construction.md)
