# Kapitel 36. Testpyramide für Wissensbasen: Regeln, Interaktionen und Antwortstabilität

> **Buch:** [Architektur evidenzbasierter Expertensysteme](README.md) · [Teil V: Verifikation, Testen, Diagnose und Sicherheitsbegründung](part-05-verification-and-learning.md)  
> **Vorheriges Kapitel:** [Kapitel 23. Verifikation der Wissensbasis: Widerspruchsfreiheit, Vollständigkeit und Robustheit von Regeln](ch23-knowledge-base-verification.md)  
> **Nächstes Kapitel:** [Kapitel 39. Aktiver Prüfexperte: Poppersche Falsifikation, normative Compliance (ASPICE/ISO 26262/ISO 21434) und autonomes Testdesign](ch39-active-compliance-auditor-and-popperian-testing.md)  
> **Inhaltsverzeichnis:** [README.md](README.md)  
> **Autor:** [Mykola Fedchyk](about-the-author.md)  
> **Niveau:** Systemarchitekten, Wissensingenieure, Verifikationsingenieure (QA/QE), Mathematiker: Fortgeschritten  
> **Lernziele:** Die vierstufige Wissenstestpyramide (Knowledge Testing Pyramid, KTP) entwerfen und implementieren; isoliertes Modultesten einzelner Regeln und ontologischer Prädikate (Knowledge Unit Testing, KUT) mit Prämissen-Mocking (`PremiseMock`) durchführen; die Falle der vakuumhaften Wahrheit (*The Vacuous Truth Trap*) für materielle Implikationen erkennen und blockieren; 6-Punkt-Spektralanalysen von Grenzwerten (Boundary Value Analysis, BVA) für normative Spezifikationsparameter anwenden; Regelkompositionstests (Knowledge Integration Testing, KIT) mit Inferenzketten-Unterbrechung durch anfechtbare Defeater (AGM-Kontraktion) konstruieren; die semantische Invarianzmetrik (Semantic Invariance Score, SIS) zur Bewertung der Robustheit gegenüber linguistischen Anfragevariationen mit einem Zulassungsschwellenwert von $\text{SIS} \ge 0{,}98$ berechnen; die Lipschitz-Stabilität ($L_{\mathcal{K}} \le L_{\max}$) des Inferenzraums verifizieren, um cyber-physikalische Systeme vor sprunghaftem Relais-Chattering zu schützen; stigmergische Wissenslücken-Puffer (`KnowledgeGapSpool`) mit Pheromon-Priorisierung für das Ingenieur-Backlog etablieren.

---

## Abstract

Dieses Kapitel untersucht die Methodik der mehrstufigen Verifikation und Kalibrierung von Regeln in evidenzbasierten Expertensystemen auf Basis der vom Autor entwickelten Wissenstestpyramide (*Knowledge Testing Pyramid, KTP*). Es definiert die mathematischen und architektonischen Werkzeuge, die Inferenzstabilität gegenüber Sensorrauschen garantieren und eine Degradation der Wissensbasis in missionskritischen Anwendungen verhindern.

In sicherheitskritischen cyber-physikalischen Systemen (ISO 26262 ASIL D, EN 50128 SIL 4, DO-178C Level A) führt das Fehlen einer isolierten Regelverifikation und Inferenzgrenzenprüfung zu zwei katastrophalen Ausfallarten: der Falle der vakuumhaften Wahrheit (*The Vacuous Truth Trap*, bei der ein falscher oder uninitialisierter Antezedens einer materiellen Implikation $`A \to B`$ die Regel automatisch wahr macht und fälschlicherweise Notfall-Aktuatoren ansteuert) sowie hochfrequentem Prellen (*Chattering*, bei dem geringfügiges Rauschen auf Sensorkanälen an Prädikatsgrenzen sprunghafte Entscheidungsumschaltungen auslöst und elektromechanische Antriebe zerstört). Zusätzlich wird die Gefahr des pseudowissenschaftlichen Metriken-Inflationismus in der generativen KI aufgedeckt, wo Behauptungen über null Halluzinationen ($`\mathrm{ZHR} = 1{,}000000`$) bei kleinen Stichproben den Grundlagen der mathematischen Statistik widersprechen.

Dieses Kapitel löst dieses Problem durch den Aufbau einer vierstufigen Wissenstestpyramide (*Knowledge Testing Pyramid, KTP*): von isolierten Komponententests von Regeln (*KUT*) mit Prämissen-Mocking (`PremiseMock`), dem mathematischen Nicht-Vakuitäts-Kriterium (*Non-Vacuity Invariant*) und 6-Punkt-Spektralanalysen von Grenzwerten (*BVA*) über das Integrationstesten anfechtbarer Inferenzketten (*KIT*) bis hin zur systemischen Evaluierung der semantischen Invarianz ($`\mathrm{SIS} \ge 0{,}98`$), der variationellen Kalibrierung der Lipschitz-Stabilität ($`L_{\mathcal{K}} \le L_{\max}`$), exakten Clopper-Pearson-Konfidenzintervallen, der Dreierregel, dem isolierten Kontrollgate (*Held-Out Sealed Gate*) mit unverfälschter binärer Verwahrung sowie der mehrstufigen Isolation von Testgeländen (Agent Sandbox) auf Basis von Mikrovirtualisierung (MicroVM Firecracker), Kernel-basiertem eBPF-LSM-Schutz und einem externen Hardware-Supervisor zur Verhinderung von Sandbox-Ausbrüchen (Sandbox Escape).

---

## 1. Die methodologische Lücke in der Verifikation von Wissenssystemen

Zur Strukturierung von Softwareprüfungen dient die klassische Testpyramide nach Mike Cohn [[1]](#src-1). Die testgetriebene Entwicklung (TDD) nach Kent Beck [[2]](#src-2) ergänzt dies um explizite Erwartungshaltungen für isolierte Komponenten. In diesem Kapitel wird die Pyramide als Entwicklungspfad von isolierten Einzelprüfungen hin zu Integrations- und End-to-End-Tests dargestellt:

```math
\text{Unit Tests} \longrightarrow \text{Integration Tests} \longrightarrow \text{End-to-End / System Tests}
```

Kein industriell ausgereiftes Softwareprodukt wird für den Produktivbetrieb freigegeben, bloß weil der Compiler keine Syntaxfehler meldet (statische Analyse) oder weil das System einige manuelle Demonstrationsszenarien bestanden hat. Jede Klasse, Funktion und jedes Modul wird über Test-Doubles (*Stubs, Mocks, Fakes* [[3]](#src-3)) isoliert, und Randbedingungen werden an den Extremwerten verifiziert.

Im Wissensingenieurwesen (Knowledge Engineering) und der Domäne der Expertensysteme existierte demgegenüber über Jahrzehnte eine eklatante **methodologische Lücke**:

```mermaid
flowchart LR
    accTitle: Traditionelle Lücke im Knowledge Engineering
    accDescr: Diskrepanz zwischen statischer Regelanalyse und makroskopischen Benchmarks ohne isolierte Komponententests.

    subgraph TraditionalGap["<b>Traditionelle Lücke im Knowledge Engineering</b>"]
        STATIC["Statische Syntaxanalyse<br/>(SAT/SMT-Solver, Zyklusprüfung)"] -.-> GAP["<b>METHODOLOGISCHE LÜCKE</b><br/>• Keine isolierten Tests von Wissensatomen<br/>• Fehlende Prämissen-Mocks (Premise Mocks)<br/>• Unempfindlichkeit gegenüber Formulierungsvarianz<br/>• Risiko der vakuumhaften Wahrheit"]
        GAP -.-> E2E["Makroskopische Benchmarks<br/>(F1-Score, Genauigkeit auf Datensätzen)"]
    end
```

Historisch beschränkte sich die Verifikation entweder auf die **statische Regelverifikation** (Suche nach Zyklen, Redundanzen und syntaktischen Widersprüchen mittels SMT-Solvern [[4]](#src-4)) oder ging unmittelbar zu **makroskopischen Benchmarks** über (Evaluierung von Genauigkeit, F1-Score und Wahrscheinlichkeitskalibrierungen nach Platt/ECE über hunderte Anfragen [[5, 6]](#src-5)).

Scheitert ein Expertensystem an einer komplexen Anfrage, fehlen dem Ingenieur ohne atomare Tests die Werkzeuge, um das Problem eindeutig zu lokalisieren:
1. Enthält die Regel selbst einen Fehler (Defekt im Antezedens)?
2. Entstand das Problem durch eine fehlerhafte Typvererbung im Ontologieverband?
3. Wurde die Konklusion durch das irrtümliche Feuern eines unterminierenden Defeaters (*Undercutting Defeater*) blockiert?
4. Wurde die Anfrage geringfügig linguistisch umformuliert, was das semantische Parsen verfälschte?

### 1.1. Vergleichende Analyse internationaler Ansätze für Testen und Verifikation

| Ansatz / Schule | Vertreter und Quellen | Fokus und Stärken | Grenzen für Wissenssysteme |
|---|---|---|---|
| **Klassisches Softwaretesten** | M. Cohn [[1]](#src-1), K. Beck [[2]](#src-2), M. Feathers [[3]](#src-3) | Komponentenisolation, Unit-Mocks, TDD, Regressions-Suites. | Auf deterministische prozedurale Funktionen ausgelegt; ignoriert logisches Resolvieren, CWA-Unvollständigkeit und Überzeugungsrevision. |
| **Formale Verifikation und SMT** | C. Barrett, L. de Moura, N. Bjørner (Z3) [[4]](#src-4) | Strenge Theorembeweise, Erfüllbarkeitsprüfung von Prädikatenformeln. | Statische Analyse der Regelbasis als geschlossenes System; testet kein Verhalten auf empirischen Sensorströmen und Sprachumformulierungen. |
| **Kalibrierung von KI-Modellen (ECE)** | J. Platt [[5]](#src-5), C. Guo et al. (On Calibration of Modern Neural Networks, 2017) [[6]](#src-6) | Abgleich skalarer Klassifikatorwahrscheinlichkeiten mit realen Fehlerraten (Expected Calibration Error). | Bewertet lediglich skalare Konfidenz auf festen Datensätzen; blind gegenüber logischer Beweisstruktur und Perturbationssensitivität. |
| **Metamorphes Testen und CheckList** | T. Y. Chen et al. [[7]](#src-7), M. T. Ribeiro et al. (CheckList, ACL 2020) [[8]](#src-8) | Testen von NLP-Verhaltenseigenschaften ohne Orakel über semantische Textperturbationsinvarianten. | Auf Black-Box-Neuronale-Netze ausgerichtet; keine Analyse byteweiser Beweise und deterministischer Regelverbände. |
| **Anfechtbare Argumentation und AGM** | J. Pollock [[9]](#src-9), P. M. Dung [[10]](#src-10), C. Alchourrón, P. Gärdenfors, D. Makinson [[11]](#src-11) | Formale Philosophie anfechtbarer Inferenz, Argumentkonflikte, minimale Überzeugungsänderung. | Theoretische logische Abstraktionen ohne programmatische Umsetzung auf Unit-/Integrationsstufen und ohne ingenieurtechnische Robustheitsmetriken. |
| **Vierstufige Wissenstestpyramide (KTP)** | **Methodik dieses Buches (KTP-Framework)** | **Vierstufiges Modell KUT/KIT/KST/KCT: Prämissen-Mocking, Detektion vakuumhafter Wahrheit, BVA, semantische SIS-Invarianz, Lipschitz-Stabilität $L_{\mathcal{K}}$ und Stigmergie.** | **Integriertes ingenieurtechnisches Modell vom isolierten Wissensatom bis zur zertifizierungsreifen Kalibrierung cyber-physikalischer Systeme.** |

---

## 2. Konzeptuelles Modell: Die vierstufige Wissenstestpyramide

Die vom Autor entwickelte **Wissenstestpyramide (Knowledge Testing Pyramid, KTP)** strukturiert die Verifikation eines intelligenten Systems in vier Stufen unterschiedlicher Strenge und Granularität:

```mermaid
flowchart TD
    accTitle: Vierstufige Wissenstestpyramide (Knowledge Testing Pyramid)
    accDescr: Architektonische Verifikationsebenen von isolierten Wissensatomen bis zur komplexen variationellen Kalibrierung.

    L4["<b>Stufe 4: Variationelle Kalibrierung (KCT)</b><br/>Semantische Invarianzmetrik (SIS &ge; 0,98), Lipschitz-Stetigkeit (L<sub>K</sub> &le; L<sub>max</sub>), Anti-Chattering"]
    L3["<b>Stufe 3: Durchgängiges Systemtesten (KST / E2E)</b><br/>End-to-End-Inferenzpfade, GSN-Zertifizierungsnachweise, Fail-Closed-Zulassungsschleusen"]
    L2["<b>Stufe 2: Wissensintegrationstests (KIT)</b><br/>Mehrstufige Deduktionsverbände, Defeater-Interruption (AGM), minimale Konflikte (QuickXPlain)"]
    L1["<b>Stufe 1: Unit-Testing von Wissensatomen (KUT)</b><br/>Isolierte Regeln, Prädikate und N-stellige Frames, Prämissen-Mocks (PremiseMock), Detektion vakuumhafter Wahrheit, 6-Punkt-BVA"]

    L1 --> L2 --> L3 --> L4

    classDef l1 fill:#e1f5fe,stroke:#0288d1,stroke-width:2px,color:#01579b;
    classDef l2 fill:#ede7f6,stroke:#512da8,stroke-width:2px,color:#311b92;
    classDef l3 fill:#fff3e0,stroke:#f57c00,stroke-width:2px,color:#e65100;
    classDef l4 fill:#e8f5e9,stroke:#388e3c,stroke-width:2px,color:#1b5e20;

    class L1 l1;
    class L2 l2;
    class L3 l3;
    class L4 l4;
```

Die KTP-Struktur deckt das gesamte Spektrum der Verifikationsaufgaben des Knowledge Engineering ab:
1. **Ebene L1 (Atomsyntax und Semantik in KUT):** Verifikation geschlossener Bedingungsvokabulare, Konsistenz von Datentypen, Ausschluss von Tippfehlern in Prädikatbezeichnern sowie strikte Validierung von SI-Einheiten;
2. **Ebene L2 (Nicht-Vakuität der Prädikate und BVA):** obligatorische Aktivierung jeder Regel durch Gegenbeispiele; Generierung von Testvektoren, die nachweisen, dass die Regel den Systemzustand tatsächlich verändern kann ($`A = \mathrm{True}`$ und $`\mathrm{StateChanged}`$);
3. **Ebene L3 (Integrität deduktiver Theorien und ASPIC+ in KIT):** Erkennung zyklischer Angriffe, Prüfung der korrekten Bindung anfechtbarer Einwände (Defeaters) und Eliminierung verwaister Regeln in mehrstufigen Gittern;
4. **Ebene L4 (Kreuzarchitektonische Hard- und Software-Parität sowie KCT):** Gegenüberstellung der Inferenzresultate des High-Level-Softwaresimulators und des FPGA-RTL-Hardwarekerns (EPU) unter antagonistischen Szenarien, variationelle Kalibrierung der Lipschitz-Stabilität und Auditierung der semantischen Invarianz.

---

## 3. Stufe 1: Knowledge Unit Testing (KUT) – Testen von Wissensatomen

### 3.1. Definition und Testobjekt
**Knowledge Unit Testing (KUT)** bezeichnet das isolierte Testen der kleinsten unteilbaren Wissenseinheit (Ontologieatom, einzelne Regel, $N$-stelliger semantischer Frame) ohne Anbindung der globalen Faktenbasis und ohne rekursive Inferenzentfaltung.

Testobjekte des KUT sind:
* Eine einzelne normative Regel $R: \text{Antecedents} \longrightarrow \text{Consequent}$;
* Ein einzelner $N$-stelliger Frame (Rollen für Aktor, Prädikat, Objekt, deontische Modalität SHALL/MUST, Vorbedingungen und Ausnahmen);
* Ein einzelnes Ontologieprädikat mit numerischen Gültigkeitsbereichen.

### 3.2. Methode der Prämissen-Mocking (Premise Mocking)
In einem realen System greift eine Regel auf den Wissensgraphen zu:

```math
\mathtt{Query}(\mathtt{"engine\_rpm"}) > 3000 \land \mathtt{Query}(\mathtt{"oil\_temp"}) > 100 \implies \mathtt{Mode} = \mathtt{"COOLING\_HIGH"}
```

Während der KUT-Ausführung wird die globale Umgebung durch einen isolierten Stub-Kontext (`PremiseMock`) ersetzt:
```go
mock := knowledgetest.NewPremiseMock()
mock.Set("engine_rpm", 3500.0)
mock.Set("oil_temp", 105.0)

res := knowledgetest.RunKUT(coolingRule, mock, "COOLING_HIGH")
```

### 3.3. Die Falle der vakuumhaften Wahrheit (The Vacuous Truth Trap)
In der klassischen mathematischen Logik ist die materielle Implikation $`P \to Q`$ äquivalent zur Disjunktion $`\neg P \lor Q`$. Ist der Antezedens $`P`$ falsch, ist der Ausdruck formal wahr ($`P \equiv \text{False} \implies (P \to Q) \equiv \text{True}`$), unabhängig vom Inhalt der Konklusion $`Q`$.

In technischen Expertensystemen erzeugt dies eine fatale Schwachstelle: Eine fehlerhaft entworfene Regel oder ein naiver Testrunner wertet eine Prüfung als „erfolgreich bestanden“, weil die formale Implikationsbedingung erfüllt ist – obgleich in der Praxis keine einzige der erforderlichen technischen Vorbedingungen je aktiv war. Beispielsweise nimmt das Antezedens $`P`$ infolge eines Tippfehlers im Variablenbezeichner oder fehlerhafter Bereichsgrenzen niemals den Wert `TRUE` an, wodurch ein kritischer Notfallschutz-Algorithmus zu „totem Code“ wird, der formal sämtliche statischen SMT-Prüfungen besteht.

```mermaid
flowchart TD
    accTitle: Falle der vakuumhaften Wahrheit (Vacuous Truth Trap)
    accDescr: Schema der Blockierung der vakuumhaften Wahrheit durch obligatorische Generierung aktivierender Gegenbeispiele.

    subgraph "Falle der vakuumhaften Wahrheit (Vacuous Truth Trap)"
        COND["Prämissen im PremiseMock fehlen oder sind falsch (P = False)"] --> IMPL["Materielle Implikation: False &rarr; Q &equiv; True"]
        IMPL --> VULN["<b>KATASTROPHALER DEFEKT</b><br/>Regel gilt als valide,<br/>feuert jedoch im Flug/Betrieb niemals!"]
        VULN --> SOL["<b>LÖSUNG: Non-Vacuous Invariant</b><br/>Obligatorische Generierung eines aktivierenden Gegenbeispiels<br/>und Überprüfung der Ausgangsänderung"]
    end
```

Eine Wissensbasis gilt genau dann auf Ebene L2 als verifiziert, wenn das mathematische Nicht-Vakuitäts-Kriterium (**Non-Vacuity Invariant**) erfüllt ist:

```math
\forall \mathcal{R} = (A \implies C) \in \mathcal{KB}, \quad \exists \mathbf{x} \in \mathcal{X}_{\mathrm{test}} : A(\mathbf{x}) = \mathrm{True} \land \mathrm{StateChanged}(\mathbf{x}, C)
```

wobei:
- $`\mathcal{R}`$ — normative oder ingenieurtechnische Regel mit Antezedens $`A`$ und Konsequens $`C`$;
- $`\mathcal{KB}`$ — verifizierte Wissensbasis des Expertensystems;
- $`\mathcal{X}_{\mathrm{test}}`$ — Raum synthetisierter oder empirischer Testvektoren;
- $`\mathbf{x}`$ — konkreter Eingangsbeobachtungsvektor (Werte der Umgebungsparameter);
- $`A(\mathbf{x}) = \mathrm{True}`$ — Erfüllung der Antezedensbedingung ohne Auslösung eines Schutzfehlers oder Überschreitung der Bereichsgrenzen;
- $`\mathrm{StateChanged}(\mathbf{x}, C)`$ — Prädikat der Zustandsänderung oder Bildung einer neuen Schlussfolgerung unter Einwirkung des Konsequens $`C`$.

#### 3.3.1. Geschlossener Kreis der ingenieurtechnischen Berechnung der Nicht-Vakuitätsverifikation

1. **Laufzeitsteuerung und Rechenfluss (Control Flow & Runtime Decisions):**  
   Findet der KUT-Testrunner für eine Regel $`\mathcal{R}_i`$ keinen einzigen Vektor $`\mathbf{x} \in \mathcal{X}_{\mathrm{test}}`$, der $`A(\mathbf{x})`$ in den Zustand $`\mathrm{True}`$ überführt, wird die Regel als potenziell tot markiert (`VACUOUS_RULE_STALL`). Ein Gegenbeispiel-Finder (SMT-Solver Z3) wird initiiert. Beweist die SMT-Analyse, dass das Antezedens aufgrund eines logischen Widerspruchs in den Vorbedingungen oder einer Typenfehlanpassung identisch falsch ist ($`\neg \exists \mathbf{x}: A(\mathbf{x})`$), fällt die CI/CD-Pipeline das Urteil `CRITICAL_REJECT` und blockiert die Bereitstellung des Release-Pakets der Wissensbasis.
2. **Hardware-Dimensionierung und infrastrukturelle Beschränkungen (Hardware Dimensioning):**  
   Das Nicht-Vakuitäts-Audit stellt strenge Anforderungen an die Dimensionierung des Testraums: Der minimale Umfang der Testsuite beträgt $`|\mathcal{X}_{\mathrm{test}}| \ge 2 \cdot |\mathcal{KB}|`$ (mindestens ein Vektor für die gezielte Aktivierung der Vorbedingung und einer für deren Deaktivierung). Für eine Wissensbasis mit $`10\,000`$ Regeln erfordert dies $`20\,000`$ dedizierte Vektoren. Bei einer Taktfrequenz des Inferenz-Emulators von $`200\,\text{MHz}`$ und einer Latenz von $`15\,\mu\text{s}`$ pro Test dauert der vollständige Zyklus des Nicht-Vakuitäts-Audits $`300\,\text{ms}`$, was eine nahtlose Integration in Git-Hooks vor jedem Branch-Merge ermöglicht.
3. **Praktisches numerisches Rechenbeispiel (Worked Numerical Example):**  
   Betrachten wir die Notabschaltregel einer Dampfturbine $`\mathcal{R}_{704}`$: «WENN $`\mathtt{rpm} > 12000`$ UND $`\mathtt{coolant\_flow} < 5{,}0\,\text{l/min}`$, DANN $`\mathtt{trip\_solenoid} = \mathtt{ACTIVE}`$». Infolge eines Tippfehlers im ontologischen Mapping wurde der Sensor als `coolant_flw` registriert (der Buchstabe `o` fehlt). Da der Fakt `coolant_flow` in der Laufzeitumgebung niemals erzeugt wird, bleibt das Antezedens permanent falsch ($`A \equiv \text{False}`$), wodurch die materielle Implikation $`A \to C`$ formal wahr ist ($`\text{True}`$). Der KUT-Testrunner detektiert den Defekt `VACUOUS_RULE_ERROR`, bricht das Qualifizierungsverfahren ab und lokalisiert den Bruch der ontologischen Bindung, wodurch ein Turbosatz im Wert von 2,4 Mio. USD vor der physischen Zerstörung bei einer realen Überhitzung bewahrt wird.

**KUT-Invariante #1 (Vacuous Implication Prevention Invariant):**  
Der KUT-Testrunner ist verpflichtet, einen Test zwingend als fehlerhaft zurückzuweisen (`VacuousTruthTrap = true`), wenn die Regel ein erfolgreiches Feuern reklamiert, obwohl obligatorische Prämissen in der Mock-Umgebung fehlen oder unvollständig sind.

Vor dem Ausführen der Regel muss der Test ein unabhängiges erwartetes Ergebnis vorgeben. Ein Statusflag, das die Regel selbst als „vakuumhafte Wahrheit“ setzt, stellt kein unabhängiges Orakel dar: Eine defekte Regel setzt ihr eigenes Fehlerflag möglicherweise nicht. Für eine Regel mit zwingender Prämisse sind mindestens drei Testfälle erforderlich:

| Kontrollierter Prämissenzustand | Erwartetes Regelverhalten | Was einen Defekt darstellt |
|---|---|---|
| Prämisse bestätigt | Regel feuert und liefert erwartete Konklusion | Regel feuert nicht oder verfälscht Konklusion |
| Prämisse widerlegt | Regel leitet Konklusion nicht ab; Ergebnis hat Status „nicht anwendbar“ | Konklusion trotz falscher Prämisse abgeleitet |
| Prämissenwert unbekannt | Inferenzmaschine bewahrt Unbekanntheit und benennt benötigte Evidenz | Fehlen eines Faktums wird stillschweigend in Falschheit oder Freigabe umgewandelt |

Der boolesche Faktumswert `false` bedeutet nicht das Fehlen des Faktums. Beispielsweise kann eine Regel, die einen deaktivierten Wartungsmodus prüft, explizit den bestätigten Wert `false` verlangen. Der Test muss die Regelbedingung verifizieren und darf nicht jedes falsche boolesche Attribut als fehlende Prämisse missdeuten. Die dreiwertige Logik wird in [Kapitel 6](ch06-applied-mathematics-for-expert-systems.md) vertieft.

### 3.4. 6-Punkt-Spektralanalyse von Grenzwerten (BVA)
Normative Dokumente (RFCs, ISO-Normen, Gesetze) enthalten numerische Parameter: Timeouts, Schwellenspannungen, minimale Paketgrößen. Um Fehler wie „strikte statt nicht-strikter Ungleichung“ ($<$ statt $\le$) aufzudecken, führt der Ansatz des Autors eine 6-Punkt-Spektralanalyse von Grenzwerten für jeden Parameter eines Bereichs $[v_{\min}, v_{\max}]$ ein:

```mermaid
xychart-beta
    title "6-Punkt-BVA-Spektrum um das Betriebsintervall [15.0, 45.0]"
    x-axis ["min- (14.9)", "min (15.0)", "min+ (15.1)", "nom (30.0)", "max- (44.9)", "max (45.0)", "max+ (45.1)"]
    y-axis "Ausführungsstatus (0 = Inaktiv, 1 = Gefeuert)" 0 --> 1
    bar [0, 1, 1, 1, 1, 1, 0]
```

Punkte des Spektrums:
1. $v_{\min-} = v_{\min} - \delta$ — Punkt außerhalb des Bereichs (Ablehnung / Sperre erwartet);
2. $v_{\min}$ — exakte untere Grenze (Auslösung / Freigabe erwartet);
3. $v_{\min+} = v_{\min} + \delta$ — Punkt unmittelbar innerhalb des Bereichs;
4. $v_{\text{nom}}$ — nominaler Betriebswert;
5. $v_{\max-} = v_{\max} - \delta$ — Punkt nahe der oberen Grenze;
6. $v_{\max}$ — exakte obere Grenze;
7. $v_{\max+} = v_{\max} + \delta$ — Überschreitung der oberen Grenze (Blockade erwartet).

---

## 4. Stufe 2: Knowledge Integration Testing (KIT) – Inferenzverbände und Defeater

### 4.1. Mehrstufige Inferenzverbände (Inference Lattices)
Auf der KIT-Stufe wird die Interaktion benachbarter Regeln getestet, bei denen die Konklusion der vorangehenden Regel zur Prämisse der nachfolgenden wird:

```math
R_1: A \longrightarrow B, \qquad R_2: B \longrightarrow C, \qquad \dots, \qquad R_n: Y \longrightarrow Z
```

KIT verifiziert:
* **Typschnittstellenkompatibilität:** ob ontologische Attributtypen über verschiedene Ontologieschichten hinweg übereinstimmen;
* **Erhaltung der Evidenz-Traceability:** ob die Kette byteweiser Zitate über alle Zwischenknoten ohne Verlust der Primärquellen weitergegeben wird;
* **Zyklusfreiheit:** Erkennung wechselseitiger Regelaufrufe außerhalb zulässiger endlicher Zustandsautomaten.

Ein Beispiel für semantische Verschiebung: Die erste Regel liefert ein „Socket-Timeout“, während die zweite ein „Verbindungsaufbau-Timeout“ erwartet. Dieselbe numerische Einheit und ein ähnlicher Name begründen keine Identität der Größen. Der Integrationstest übergibt den Wert zunächst ohne Abbildungsregel und erwartet eine Blockade der Komposition; anschließend fügt er eine genehmigte Zuordnung mit dem erforderlichen Kontext hinzu und prüft den erlaubten Übergang. Dieses Testfallpaar unterscheidet eine inhaltliche Prüfung von einer Politik, die pauschal jede Kette blockiert.

### 4.2. Injektion anfechtbarer Defeater (Defeater Interruption nach AGM)
Gemäß der Argumentationstheorie von John Pollock [[9]](#src-9) und den AGM-Postulaten der Überzeugungsrevision [[11]](#src-11) muss das Eintreten eines anfechtenden Umstands die Argumentationskette unverzüglich unterbrechen:

```mermaid
flowchart LR
    accTitle: Unterbrechung der Inferenzkette durch einen Defeater
    accDescr: Ein aktiver Defeater blockiert die Ableitung und widerruft die nachgelagerte Schlussfolgerung.

    A["Eingangspaket valide (A)"] -->|Regel R1| B["Paket akzeptiert (B)"]
    B -->|Regel R2| C["Übergabe an FSM-Zustandsautomat (C)"]

    DEF["<b>Aktiver Defeater:</b><br/>SYN-Flood-Angriff erkannt!"] -.->|Untergräbt R2| CUT["<b>INFERENZBLOCKADE</b><br/>Konklusion C annulliert!"]
```

**KIT-Invariante #1 (Defeater Dominance Invariant):**  
Bei Aktivierung eines verifizierten unterminierenden Defeaters ($D$) muss das System die Inferenz am Zielknoten deterministisch anhalten, alle abgeleiteten Konklusionen annullieren (AGM-Kontraktion) und eine Ablehnung unter Protokollierung der exakten Blockadeursache zurückgeben.

Das Unterminieren einer Prämisse und das Widerlegen einer Konklusion erfordern unterschiedliche Testerwartungen. Die Meldung „Sensor nicht kalibriert“ untergräbt die Verwendung der Messung als Beweis für eine Überhitzung, beweist jedoch nicht das Fehlen einer Überhitzung und löscht die Messung nicht aus dem Protokoll. Die Evidenz „Temperatur unterhalb der Schwelle“, die über einen anderen validen Kanal empfangen wurde, widerlegt hingegen die Konklusion der Überhitzung selbst. Integrationstests müssen die Erhaltung des Primärdatensatzes, den Grund des Widerrufs und alternative Begründungen prüfen – und nicht bloß das endgültige Scheitern [[9]](#src-9).

### 4.3. Minimaler Konflikt als verifizierbares Ergebnis

Sind Randbedingungen unvereinbar, ist es sinnvoll, nicht nur das Vorliegen eines Konflikts zu prüfen, sondern auch dessen Erklärung. Der QuickXPlain-Algorithmus von Ulrich Junker identifiziert unter definierten Konsistenzprüfungs-Prämissen einen nach Inklusion minimalen Konflikt [[12]](#src-12). Minimalität bedeutet hierbei, dass das Entfernen eines beliebigen Elements aus der gefundenen Teilmenge deren Unverträglichkeit aufhebt; eine minimale Kardinalität unter allen denkbaren Konflikten garantiert der Algorithmus nicht.

In einem didaktischen Beispielszenario sind simultan drei Bedingungen vorgegeben: Temperatur nicht unter 95 °C, Temperatur nicht über 90 °C und Druck nicht über 10 bar. Vor einem konsistenten Hintergrundwissen bilden die beiden Temperaturbedingungen den Konflikt. Die Druckbedingung darf nicht in diese Erklärung einfließen. Eine unabhängige Prüfung bestätigt die Unverträglichkeit des Paares sowie die Widerspruchsfreiheit jedes der beiden isolierten Einzelelemente.

Eigene Testfälle sind für leere Mengen, konsistente Mengen sowie inkonsistente Hintergrundbeschreibungen erforderlich. Ein Verstoß gegen die Vorbedingungen des Algorithmus darf nicht fälschlich als minimaler Konflikt deklariert werden. Implementierung und Verifikation von QuickXPlain werden in [Kapitel 20](ch20-explanation-engine.md) behandelt und hier nicht redundant wiederholt.

---

## 5. Stufen 3 und 4: Durchgängige Systemprüfungen und variationelle Kalibrierung

### 5.1. Durchgängiges Systemtesten des Expertensystems

Das durchgängige Systemtesten des Expertensystems (*Knowledge System Testing*, KST) überprüft den vollständigen Pfad von der Benutzereingabe bis zur Entscheidung und dem Begründungspaket. Ein erfolgreicher Modultest einer Regel belegt keineswegs, dass die Suche die gültige Quelle ausgewählt hat, die Sprachanalyse Negationen bewahrt hat oder die Erklärung den tatsächlichen Schluss korrekt wiedergibt.

| Prüfgrenze | Testfall | Unabhängige Erwartung |
|---|---|---|
| Quelle → Behauptung | korrektes Zitat, aber Zahl bezieht sich auf anderes Objekt | Byte-Übereinstimmung berechtigt nicht zur Akzeptanz einer Fehlinterpretation |
| Prämissen → Beweisgraph | eine obligatorische Prämisse entfernt | Konklusion wird verweigert; Paket benennt fehlende Begründung |
| Graph → Erklärung | Text enthält Zahl, die nicht im Trace existiert | unbestätigter Text wird blockiert oder durch geprüfte Repräsentation ersetzt |
| Zugriff und Version → Ausgabe | Quelle nach vorangegangener erfolgreicher Anfrage widerrufen | neue Ausgabe verwendet widerrufene Prämisse nicht |

Für jeden negativen Fall ist ein gepaarter positiver Fall erforderlich: gültige Quelle, vollständige Prämissen, korrekte Erklärung und autorisierter Zugriff. Andernfalls würde ein permanentes Verweigern fälschlicherweise als erfolgreiches Bestehen der Tests gewertet. Verträge für Begründungen und Erklärungen werden in den Kapiteln [19](ch19-from-question-to-evidence.md) und [20](ch20-explanation-engine.md) detailliert erörtert, die kontrollierte Freigabe neuer Versionen in [Kapitel 25](ch25-how-expert-systems-learn.md).

### 5.2. Metrik der semantischen Invarianz (Semantic Invariance Score, SIS)
Klassische Prüfbenchmarks enthalten eine einzige, statische Frageformulierung. Im realen Produktivbetrieb formulieren Ingenieure und Bediener dieselbe Frage jedoch auf hunderte verschiedene linguistische Weisen:

```math
\begin{aligned}
\mathcal{Q}_{\text{base}} &= \text{„Was ist die minimale MTU für IPv6?“}, \\
\mathcal{Q}_{\text{var1}} &= \text{„Geben Sie die kleinste zulässige Paketgröße in IPv6-Netzen an“}, \\
\mathcal{Q}_{\text{var2}} &= \text{„Least transmission unit required by RFC 8200 IPv6 specification“}.
\end{aligned}
```

Liefert das System auf $\mathcal{Q}_{\text{base}}$ das Ergebnis `1280 octets`, verweigert jedoch bei $\mathcal{Q}_{\text{var1}}$ die Antwort unter CWA oder ändert den Beweispfad, ist dieses System linguistisch instabil.

Die vom Autor entwickelte Metrik **Semantic Invariance Score** ($`\mathrm{SIS}`$) bewertet die Robustheit des Systems auf der Mannigfaltigkeit der Anfrage-Perturbationen $`\mathbb{V}(\mathcal{Q})`$:

```math
\mathrm{SIS}(\mathcal{Q}) = \alpha \cdot \mathrm{VerdictsMatchRate} + \beta \cdot \mathrm{ProofGraphJaccard}
```

wobei die Komponenten der Metrik wie folgt berechnet werden:

```math
\mathrm{VerdictsMatchRate} = \frac{1}{N} \sum_{i=1}^{N} I\left(\mathrm{Verdict}(Q_i) = \mathrm{Verdict}(Q_{\mathrm{base}})\right) \in [0, 1]
```

```math
\mathrm{ProofGraphJaccard} = \frac{1}{N} \sum_{i=1}^{N} \frac{\lvert P(Q_i) \cap P(Q_{\mathrm{base}})\rvert}{\lvert P(Q_i) \cup P(Q_{\mathrm{base}})\rvert} \in [0, 1]
```

**Parameter und Gewichtungsfaktoren:**
- $`\mathrm{VerdictsMatchRate}`$ — Anteil perturbierter Anfragen aus $`N`$ Varianten, bei denen das Urteil exakt mit dem Urteil auf die Basisformulierung übereinstimmt;
- $`\mathrm{ProofGraphJaccard}`$ — gemittelter Jaccard-Koeffizient der Ähnlichkeit der Mengen von Knoten und Zitaten des Beweisgraphen $`P(Q)`$;
- $`\alpha, \beta \in [0, 1]`$ — Kalibrierungsgewichte der Konfidenz ($`\alpha = 0{,}6,\; \beta = 0{,}4`$, wobei $`\alpha + \beta = 1{,}0`$);
- $`\mathrm{SIS}(\mathcal{Q}) \in [0, 1]`$ — resultierender integraler Index der semantischen Invarianz.

**Zulassungskriterium und Laufzeitreaktion:**  
Für die automatisierte Schleuse der variationellen Freigabe gelten folgende Schwellenwertregeln:

```math
\mathrm{SIS}(\mathcal{Q}) \ge 0{,}98
```

- `ACCEPT_INVARIANT`: Wenn $`\mathrm{SIS}(\mathcal{Q}) \ge 0{,}98`$ und $`\mathrm{VerdictsMatchRate} = 1{,}00`$, gilt die Wissensbasis als invariant gegenüber Anfrage-Perturbationen und erhält die Ed25519-Release-Signatur;
- `QUALIFIED_REVIEW`: Liegt $`0{,}85 \le \mathrm{SIS}(\mathcal{Q}) < 0{,}98`$ vor, detektiert das System eine Drift in der Argumentstruktur, blockiert die Veröffentlichung und eskaliert den Bericht an den Wissensingenieur unter Hervorhebung der Diskrepanzen in den Beweisgraphen;
- `REJECT_DRIFT`: Wenn $`\mathrm{SIS}(\mathcal{Q}) < 0{,}85`$ oder $`\mathrm{VerdictsMatchRate} < 1{,}00`$, wird das Release durch die CI/CD-Test-Pipeline zwingend abgewiesen.

**Praktisches Rechenbeispiel:**  
Für eine Anfrage werde eine Stichprobe von $`N = 10`$ linguistischen Variationen generiert ($`Q_1, \dots, Q_{10}`$). Bei allen 10 Varianten stimmte das Urteil mit der Basis überein (`1280 octets`), folglich gilt $`\mathrm{VerdictsMatchRate} = 10/10 = 1{,}00`$. In 9 Variationen waren die Beweisgraphen identisch zur Basis ($`\frac{\lvert P_i \cap P_{\mathrm{base}} \rvert}{\lvert P_i \cup P_{\mathrm{base}} \rvert} = 1{,}00`$). In einer Variation $`Q_8`$ fügte der Generator eine synonyme erweiterte Definition ein, wodurch ein Graph aus 5 Knoten gegenüber 4 Basisknoten bei einer Schnittmenge von 4 Knoten entstand ($`\frac{4}{5} = 0{,}80`$). Daraus folgt:

```math
\mathrm{ProofGraphJaccard} = \frac{9 \times 1{,}00 + 1 \times 0{,}80}{10} = 0{,}980.
```

Der resultierende Index beträgt:

```math
\mathrm{SIS}(\mathcal{Q}) = 0{,}6 \times 1{,}00 + 0{,}4 \times 0{,}980 = 0{,}600 + 0{,}392 = 0{,}992.
```

Da $`0{,}992 \ge 0{,}98`$, erhält das Paket den Status `ACCEPT_INVARIANT`.

Ein hoher Durchschnittswert darf das Abweichen eines kritischen Urteils nicht verschleiern. Zusätzlich wird jede nachweislich äquivalente Formulierung geprüft: Das Urteil muss mit dem Referenzwert übereinstimmen, und jede genutzte Begründung muss gültig bleiben. Der Verlust von Knoten muss die Ähnlichkeitsbewertung der Graphen auch dann beeinflussen, wenn die restlichen Knoten unverändert blieben. Ein alternativer korrekter Beweis ist für sich genommen kein Defekt; Richtlinien für zulässige Graphveränderungen müssen vor dem Testlauf festgelegt werden.

Ein leerer Satz von Variationen erlaubt keine Stabilitätsbewertung. Zwei leere Graphen belegen ebenfalls keine evidenzbasierte Antwort: Zuerst wird die Existenz und Gültigkeit der Begründungen geprüft, erst danach die Ähnlichkeit der Graphen. Die Schwelle von 0,98 sowie die Gewichte 0,6 und 0,4 stellen didaktische Modellparameter dar, keine universellen Industriestandards. Der Erfolg auf einer endlichen Menge beweist keine Invarianz über alle denkbaren Anfragen und ersetzt kein Zertifizierungsverfahren.

### 5.3. Kriterium der Lipschitz-Stabilität von Wissen ($L_{\mathcal{K}} < \infty$)
In cyber-physikalischen Komplexen (Autopiloten, Reaktorsteuerungen, medizinische Dosiersysteme) unterliegen empirische Eingangsvariablen physikalischem Rauschen und Mikrostörungen:

```math
x \longrightarrow x + \Delta x, \qquad \|\Delta x\| < 10^{-5}.
```

Wenn eine infinitesimal kleine Störung des Eingangssignals einen abrupten Sprung der Konfidenz von $1{,}0$ auf $0{,}0$ ohne physikalischen Defeater hervorruft, tritt das Phänomen des **Relais-Chatterings (Prellen)** auf, das zu zerstörerischen Eigenschwingungen von Aktoren führt.

Das Modell des Autors verlangt die Einhaltung der **Lipschitz-Stetigkeit des logischen Raums**:

```math
L_{\mathcal{K}} = \sup_{x_1 \neq x_2} \frac{\|\mathcal{K}(x_1) - \mathcal{K}(x_2)\|}{\|x_1 - x_2\|} \le L_{\max}.
```

**Parameter und ingenieurtechnische Grenzwerte:**
- $x_1, x_2 \in \mathbb{R}^n$ — Vektoren physikalischer Eingangsbeobachtungen im zulässigen Zustandsraum;
- $`\mathcal{K}(x) \in [0, 1]^m`$ — Vektor der Konfidenzen oder numerischen Empfehlungen des Expertensystems;
- $`\|\cdot\|`$ — euklidische oder Tschebyscheff-Norm im entsprechenden Vektorraum;
- $`L_{\mathcal{K}} \ge 0`$ — empirische Lipschitz-Konstante der logischen Wissensabbildung;
- $`L_{\max}`$ — maximal zulässige Steilheitsschwelle der Systemantwort (beispielsweise $`L_{\max} = 10{,}0\,\text{bar}^{-1}`$).

**Operative Entscheidungen und Zahlenbeispiel:**
- Identifiziert der Verifikator eine Unstetigkeitsstelle, an der das Differenzenverhältnis gegen unendlich strebt ($`L_{\mathcal{K}} > L_{\max}`$) für $`\|\Delta x\| \to 0`$, protokolliert das System eine Chattering-Gefahr (`CHATTERING_FAULT`), blockiert die direkte Ausgabe an den Aktorpfad und erzwingt das Einfügen einer Hysterese oder einer Dämpfungszone.
- Zahlenbeispiel: Ein Öldrucksensor meldet $`x_1 = 2{,}999\,\text{bar}`$ (wobei die logische Inferenz eine Konfidenz für den Normalbetrieb von $`\mathcal{K}(x_1) = 1{,}00`$ ausgibt) und $`x_2 = 3{,}001\,\text{bar}`$ (wobei mangels Unempfindlichkeitszone die Konfidenz schlagartig auf $`\mathcal{K}(x_2) = 0{,}00`$ einbricht). Bei einer Perturbation von $`\|x_1 - x_2\| = 0{,}002\,\text{bar}`$ beträgt der Reaktionsgradient:

```math
L_{\mathcal{K}} = \frac{\lvert 1{,}00 - 0{,}00 \rvert}{0{,}002\,\text{bar}} = 500\,\text{bar}^{-1} \gg L_{\max} = 10\,\text{bar}^{-1}.
```

- Systemaktion: Die Regel wird als gefährlich für physikalische Aktoren markiert (`REJECT_UNSTABLE_RULE`), und der Regel-Compiler generiert die Direktive zur Implementierung eines Hysteresekorridors von $\pm 0{,}15\,\text{bar}$.

### 5.4. Neuro-symbolische Generierung von Anfrage-Perturbationen $`\mathbb{V}(\mathcal{Q})`$ über lokale SLMs

Die Berechnung der Metrik $`\mathrm{SIS}(\mathcal{Q})`$ erfordert eine repräsentative Perturbationsmenge $`\mathbb{V}(\mathcal{Q}) = \{\mathcal{Q}_1, \mathcal{Q}_2, \dots, \mathcal{Q}_N\}`$. Historisch existierten im Wissensingenieurwesen zwei unzureichende Ansätze zur Bildung dieser Menge:
1. **Manuelle Paraphrasierung durch Experten:** extrem kostenintensiv, subjektiv und auf wenige Dutzend Varianten begrenzt;
2. **Schablonenbasierte algorithmische Ersetzungen:** leiden unter der Illusion des Keyword-Lookups, da sie Lexeme der Regelantezedenzien direkt in die Frage kopieren.

Für eine skalierbare und rigorose Kalibrierung wird ein **neuro-symbolischer Perturbationssynthesizer** auf Basis einer lokalen Laufzeitumgebung offener kleiner Sprachmodelle (SLMs) eingeführt:

```mermaid
flowchart TD
    accTitle: Neuro-symbolische Generierung des SIS-Perturbationsraums
    accDescr: Lokales domänenspezifisches SLM generiert sprachliche Variationen um ein Wissensatom zur SIS-Berechnung.

    ATOM["<b>Wissensatom ZKP4</b><br/>Norm: <i>«The Control Word MUST be used»</i><br/>Orakel: SHA-256, byte_start, byte_end"]
    
    subgraph "Variationsgenerator (System 1: Lokales SLM)"
        ATOM -->|Ingenieur-Prompt| SLM["<b>Domänenmodell (SLM)</b><br/>Lokales Open-Weights-Modell (Qwen2.5-Coder / Llama-3)"]
        SLM --> V1["<i>«Wie beeinflusst ein fehlerhafter Control-Word-Wert die Systemsicherheit?»</i>"]
        SLM --> V2["<i>«What practical risks arise if an implementation omits the Control Word sequence number?»</i>"]
        SLM --> V3["<i>«Is it permissible to bypass Control Word processing during high-throughput peering?»</i>"]
    end

    subgraph "Invarianzverifikation (System 2: EVM)"
        V1 --> SIS_CALC["<b>SIS-Berechnung</b><br/>(VerdictsMatchRate & ProofGraphJaccard)"]
        V2 --> SIS_CALC
        V3 --> SIS_CALC
        ATOM -->|Deterministischer Referenzbeweis| SIS_CALC
    end

    SIS_CALC --> VERDICT{"<b>SIS &ge; 0,98?</b>"}
    VERDICT -->|Ja| PASS["<b>Release-Zulassung (KCT Passed)</b>"]
    VERDICT -->|Nein| FAIL["<b>Stabilitätserosion &rarr; KnowledgeGapSpool</b>"]
```

Dieser Ansatz ermöglicht einen dualen Kalibrierungszyklus:
- **Massiver Matrix-Stresstest:** Schnelle algorithmische Generatoren erzeugen hunderttausende Testfälle ($> 250\,000$ Tests/s) zur Überprüfung von Durchsatz, Speicherlecks und korrekten Registerübergängen.
- **Neuro-symbolischer Referenzkalibrator:** Ein paralleler Pool lokaler Modelle (SLM-Worker) erzeugt hochenthaltige, authentische sprachliche Perturbationen ($1\,000$ bis $10\,000$ Fälle) für die zertifizierungsrelevante Berechnung von $`\mathrm{SIS} \ge 0{,}98`$. Dabei bleibt das deterministische Referenzorakel (System 2 EVM) zu 100 % vor Datenlecks und subjektiven Verzerrungen geschützt.

### 5.5. Metrologie der Null-Halluzinationen (ZHR): Exakte Clopper-Pearson-Konfidenzintervalle und die Dreierregel

In aktuellen Forschungsarbeiten und Marketingpräsentationen von Systemen künstlicher Intelligenz ist die gefährliche Praxis verbreitet, absolute Fehlerfreiheit in Form von $`\mathrm{ZHR} = 1{,}000000`$ (Zero-Hallucination Rate) zu deklarieren. Die Angabe einer Metrik mit sechs Nachkommastellen bedeutet formal, dass die Obergrenze der Fehlerrate wie folgt beschränkt ist:

```math
p_{\mathrm{error}} \le 10^{-6}
```

Besteht ein empirischer Benchmark jedoch lediglich aus $`N = 30`$, $`N = 100`$ oder $`N = 1\,309`$ Kontrollanfragen, berechtigt eine beobachtete Fehleranzahl von null ($`k = 0`$) **mathematisch keineswegs** dazu, eine Ausfallrate von $`10^{-6}`$ zu behaupten. Aus Sicht der mathematischen Statistik und Messtechnik belegt das Ausbleiben von Fehlern in einer kleinen Stichprobe lediglich, dass die Fehlerrate für das gewählte Konfidenzniveau im Prozent- oder Promillebereich liegt.

Für die Binomialverteilung von Bernoulli-Versuchen bei Beobachtung von $`k = 0`$ Fehlern in einer Stichprobe des Umfangs $`N`$ wird die obere Schranke der Ausfallwahrscheinlichkeit $`p_U`$ auf dem Konfidenzniveau $`1 - \alpha`$ analytisch über **exakte Clopper-Pearson-Konfidenzintervalle (Clopper-Pearson Exact Confidence Intervals)** berechnet:

```math
p_U = 1 - \alpha^{1/N}, \qquad \mathrm{ZHR}_{\mathrm{guaranteed}} = \alpha^{1/N} = (0{,}05)^{1/N}
```

wobei:
- $`N \in \mathbb{N}^+`$ — Umfang der repräsentativen Stichprobe unabhängiger Kontrollanfragen;
- $`\alpha`$ — Signifikanzniveau (beim im Ingenieurwesen standardmäßigen Konfidenzniveau von $`95\%`$ gilt $`\alpha = 0{,}05`$);
- $`p_U \in [0, 1]`$ — obere Konfidenzgrenze der Fehlerhäufigkeit (Halluzinationen) des Systems;
- $`\mathrm{ZHR}_{\mathrm{guaranteed}} \in [0, 1]`$ — garantierte Untergrenze der Null-Halluzinations-Rate.

Für $`N \to \infty`$ liefert die Taylor-Entwicklung der Exponentialfunktion die klassische ingenieurtechnische **Dreierregel (The Rule of Three)**:

```math
\ln(0{,}05) \approx -2{,}9957 \implies p_U \approx 1 - e^{-2{,}9957/N} \approx \frac{2{,}9957}{N} \approx \frac{3}{N}
```

Gemäß den Anforderungen der Metrologie und funktionalen Sicherheit (ISO 26262 ASIL D, IEC 61508 SIL 4) werden Expertensysteme in vier $`\mathrm{ZHR}`$-Reifegrade eingestuft:

| Zuverlässigkeitsklasse | Erforderliche Stichprobe $N$ ($k=0$) | Garantierte obere Ausfallgrenze ($p_{95\%}$) | Zulässige öffentliche Bezeichnung $\mathrm{ZHR}$ |
|---|:---:|:---:|:---|
| **Empirischer Prototyp (L1)** | $N < 300$ | $p \le 1{,}0 \times 10^{-2}$ | $\mathrm{ZHR} \ge 99{,}0\%$ (L1 Empirical) |
| **Industrielles Niveau (L2)** | $300 \le N < 30\,000$ | $p \le 1{,}0 \times 10^{-4}$ | $\mathrm{ZHR} \ge 99{,}7\% \dots 99{,}99\%$ (L2 Industrial) |
| **Kritisches ASIL-D (L3)** | $30\,000 \le N < 3\,000\,000$ | $p \le 1{,}0 \times 10^{-6}$ | $\mathrm{ZHR} \ge 99{,}9999\%$ (L3 Critical) |
| **Formal bewiesener Closed-World-Beweis (L4)** | Mathematischer Beweis / SAT | $p = 0$ (auf geschlossener Theorie) | $\mathrm{ZHR} = 1{,}0$ (Formal Closed-World Proof) |

Jede öffentliche Angabe zur $`\mathrm{ZHR}`$ muss zwingend den Stichprobenumfang $`N`$, das Konfidenzniveau ($`95\%`$), die Korpus-Geschlossenheitsspezifikation sowie ein offenes, reproduzierbares Testpaket umfassen. Es ist untersagt, in Nachkommastellen mehr Nullen anzugeben, als es der exponentiellen Größenordnung der Stichprobe entspricht ($`D \le \lfloor\log_{10}(N/3)\rfloor`$).

#### 5.5.1. Geschlossener Kreis der ingenieurtechnischen ZHR-Metrologieberechnung

1. **Laufzeitsteuerung und Rechenfluss (Control Flow & Runtime Decisions):**  
   Wird ein System für ASIL-D / SIL 4 qualifiziert, analysiert der automatisierte Validierungs-Arbiter den Umfang des Testkorpus. Weist der statistische Ansatz einen Stichprobenumfang von $`N < 3 \times 10^6`$ auf, blockiert das System automatisch die Erteilung eines Zertifikats der Klasse `CERT_LEVEL_3` und stuft das Profil auf `L2_INDUSTRIAL` herab oder verlangt die Hinzunahme eines diskreten SMT/SAT-Solvers für einen formalen Beweis der Klasse L4 (`FORMAL_CLOSED_WORLD_PROOF`), bei dem die Fehlerwahrscheinlichkeit innerhalb der Axiomatik $`p \equiv 0`$ beträgt.
2. **Hardware-Dimensionierung und infrastrukturelle Beschränkungen (Hardware Dimensioning):**  
   Die Durchführung von Tests der Klasse L2 ($`N = 30\,000`$ Anfragen) auf einem Servercluster mit einem Durchsatz von $`500`$ Abfragen/s erfordert $`60\,\text{s}`$. Für die L3-Qualifizierung ($`N = 3\,000\,000`$) auf einem 64-Knoten-Cluster mit einer Gesamtleistung von $`25\,000`$ Abfragen/s beträgt der Testzyklus $`120\,\text{s}`$. Bei einer durchschnittlichen kryptografischen Beweisgröße von $`4\,\text{KB}`$ pro Anfrage belaufen sich die NVMe-SSD-Speicheranforderungen zur Erfassung des Verifikationsprotokolls auf einen kontinuierlichen Puffer von $`12\,\text{GB}`$ bei einer Schreibbandbreite von $`\ge 100\,\text{MB/s}`$.
3. **Praktisches numerisches Rechenbeispiel (Worked Numerical Example):**  
   Die Testsuite für das normative Protokoll RFC 9110 umfasst $`N = 1\,309`$ atomare Prüfungen der HTTP-Semantik. Durchläuft das System alle Prüfungen fehlerfrei ($`k = 0`$), berechnet der Auditor die obere Grenze der Fehlerrate nach der Dreierregel:

```math
p_{95\%} \approx \frac{3}{1309} \approx 0{,}002292 \implies \mathrm{ZHR}_{95\%} \ge 99{,}771\%
```

   Die Marketingbehauptung «Das System erzielt $`\mathrm{ZHR} = 1{,}000000`$» wird vom Auditor als pseudowissenschaftliche Verfälschung kategorisch zurückgewiesen. Der einzig zertifizierungsfähige Eintrag lautet streng: «$`\mathrm{ZHR} \ge 99{,}77\%`$ (L2 Industrial, $`N=1309`$, $`\alpha=0{,}05`$)».

### 5.6. Isolierter Kontrollgate (Held-Out Sealed Gate) und unverfälschte binäre Verwahrung

Zuverlässigkeitsdeklarationen erweisen sich häufig aufgrund des Phänomens der «Potemkinschen Schiedsrichter» (*Potemkin Arbiters*) und des Abflusses von Validierungsdaten in den Wissensaufbereitungsprozess (*Data Contamination*) als fiktiv. Prüfstände enthalten anstelle einer vollständigen Beweisprüfung nicht selten fest im Code verdrahtete skalare Bedingungen in Testskripten und prüfen oberflächliche boolesche Flags statt einer lückenlosen byteweisen Verifikation.

Um diesen Mangel zu beheben, wird in der Architektur ein **isolierter Kontrollgate (Held-Out Sealed Gate)** auf Basis kryptografisch versiegelter Teststichproben implementiert:

```mermaid
flowchart TD
    accTitle: Architektur des isolierten Kontrollgates Held-Out Sealed Gate
    accDescr: Kryptografisch versiegelte Kontrollenklave verifiziert byteweise Nachweise ohne Datenabfluss.

    subgraph SECURE_STORE["Versiegelter kryptografischer Speicher (Held-Out Benchmark)"]
        TEST_ENC["Verschlüsselter Testdatensatz<br/>(AES-256-GCM / Ed25519)"]
        MERKLE["Merkle-Baum der Referenzzitate<br/>Root: 0x7f4a...e12"]
    end

    subgraph ENCLAVE["Isolierte sichere Enklave (TEE / Secure Enclave)"]
        DECRYPT["Entschlüsselung der Stichprobe im sicheren RAM<br/>(Schlüssel werden nach Durchlauf widerrufen)"]
        EXEC["Ausführung der Inferenz durch das Wissensmodell"]
        RAW_CHECK["<b>Gate der unverfälschten Verwahrung</b><br/>Byteweiser SHA-256-Abgleich der Zitate<br/>Verbot der Neukodierung errors=replace"]
        VERIFY_TREE["Abgleich der Merkle-Wurzel mit der Referenz"]
    end

    TEST_ENC --> DECRYPT
    MERKLE --> VERIFY_TREE
    DECRYPT --> EXEC
    EXEC --> RAW_CHECK
    RAW_CHECK --> VERIFY_TREE

    VERIFY_TREE -->|Vollständige Übereinstimmung| VERDICT_PASS["<b>Zulassungsgenehmigung (PASS)</b><br/>Generierung des kryptografischen Attestats"]
    VERIFY_TREE -->|Abweichung| VERDICT_FAIL["<b>Fail-Closed-Sperre (LOCKOUT)</b><br/>Annullierung der Sitzung, Löschung der Schlüssel"]
```

Ein zentrales Funktionsprinzip dieses Gates ist die **unverfälschte binäre Verwahrung von Beweisen (Raw Binary Custody)**. In vielen Softwaresystemen wird ein Standard-Bibliotheksdecodieren mit Ersetzung fehlerhafter Bytes praktiziert (z. B. `errors="replace"` bei UTF-8, wodurch ungültige Sequenzen durch das Zeichen `\uFFFD` ersetzt werden). Jede derartige Modifikation verändert die Originalbytes der Quelldokumente, was die kryptografische Hashkette unwiderruflich zerstört und Zitate zu einer juristischen Illusion macht.

In der Autorenarchitektur gilt eine absolute Verbotsinvariante: Der Zitatabgleich erfolgt ausschließlich auf den unberührten Byte-Slices des Originaldokumentenbestands unter Kontrolle absoluter Offsets `[byte_start, byte_end]` und SHA-256-Hashberechnung ohne zwischengeschaltete Textnormalisierung.

#### 5.6.1. Geschlossener Kreis der ingenieurtechnischen Berechnung des isolierten Zulassungsgates

1. **Laufzeitsteuerung und Rechenfluss (Control Flow & Runtime Decisions):**  
   Beim Start einer Qualifikationsrunde lädt die isolierte Enklave den versiegelten Korpus und verifiziert dessen Integrität über die Wurzel des Merkle-Baums. Wird beim byteweisen Abgleich eine Hash-Abweichung auch nur eines einzigen Bytes eines Zitats festgestellt oder der Aufruf von Bibliotheken zur Zeichenersetzung (`\uFFFD`) registriert, löst der Gate-Controller unverzüglich den Hardwaremodus `FAIL_CLOSED_LOCKOUT` aus. Die Testsitzung wird abgebrochen, ephemere Entschlüsselungsschlüssel im sicheren RAM werden mit Nullen überschrieben und der Vorfall wird im Sicherheitsüberwachungsprotokoll erfasst.
2. **Hardware-Dimensionierung und infrastrukturelle Beschränkungen (Hardware Dimensioning):**  
   Der Betrieb des isolierten Gates erfordert eine Hardware-Sicherheitsdomäne (ARM TrustZone, AMD SEV-SNP oder Intel SGX) mit einem dedizierten geschützten Speicher von mindestens $`256\,\text{MB}`$ ohne Auslagerungsdatei (No-Swap-Policy), was den Schutz von Referenzantworten gegen Host-Speicherangriffe während der Entschlüsselung garantiert. Die Durchsatzrate des kryptografischen SHA-256-Verifikationskanals muss mindestens $`1{,}2\,\text{GB/s}`$ betragen, um Engpässe bei der Testausführung zu vermeiden.
3. **Praktisches numerisches Rechenbeispiel (Worked Numerical Example):**  
   Ein versiegelter Benchmark-Korpus aus $500$ Spezifikationen mit einem Gesamtvolumen von $`64\,\text{MB}`$ besitzt die Merkle-Wurzel `0x7f4a8b...e12`. Beim Parsen eines beschädigten Fragments wandte ein Fremdmodul eine Normalisierung an, indem das Byte `0xFF` durch die 3-Byte-Sequenz `0xEF 0xBF 0xBD` (UTF-8 Replacement Character) ersetzt wurde. Der Gate der unverfälschten Verwahrung berechnete den Kontrollhash des Blocks, detektierte die Diskrepanz des Byte-Slices zum Referenzblatt des Merkle-Baums und blockierte den Testlauf innerhalb von $`12\,\mu\text{s}`$ mit dem Fehlercode `EVIDENCE_TAMPERING_DETECTED`, was eine Verfälschung der Beweisbasis verhinderte.

### 5.7. Mehrstufige Isolation von Testgeländen (Agent Sandbox) und Schutz vor Sandbox-Ausbrüchen (Sandbox Escape)

Die Durchführung von Popperschem Stresstesten (*Popperian Red Teaming*) sowie die Untersuchung autonomer Wissensagenten erfordert es, Modellen direkten Zugriff auf Compiler, Kommandozeilen-Interpreter und Ausführungswerkzeuge für synthetisierten Code zu gewähren. Während dies in akademischen Simulatoren als Standardfunktionalität gilt, birgt die Ausführung nicht vertrauenswürdigen generierten Codes in industriellen F&E-Laboren (Avionik, Energiewirtschaft, Zertifizierung nach ISO 26262 / DO-178C) eine unmittelbare Bedrohung für die Host-Infrastruktur. Ein verbreiteter ingenieurtechnischer Trugschluss besteht darin, Standard-Docker-Container als verlässliche Isolations-Sandboxes zu betrachten: Da Container den Kernel des Host-Betriebssystems Linux gemeinsam nutzen, genügt eine Zero-Day-Schwachstelle im Kernel oder ein Konfigurationsfehler (übermäßige Linux-Capabilities oder das Durchreichen des Sockets `docker.sock`), um dem autonomen Modell einen Sandbox-Ausbruch (**Sandbox Escape**) zu ermöglichen. In der Folge erlangt der Agent Kontrolle über den Host, kann geistiges Eigentum kompromittieren, Testberichte fälschen und sich lateral im Labornetzwerk ausbreiten (*lateral movement*).

Zur Beseitigung dieser kritischen Schwachstelle wurde vom Autor eine dreistufige Architektur für geschützte Testgelände (**Defense-in-Depth Agent Sandbox**) entwickelt, wie sie im nachfolgenden Diagramm dargestellt ist.

```mermaid
flowchart TD
    accTitle: Mehrstufige Architektur der Agenten-Sandbox-Isolation
    accDescr: Dreistufige Verteidigungslinie gegen Agenten-Ausbrüche: MicroVM, Kernel-basierter eBPF-LSM-Systemaufruf-Filter und Hardware-Supervisor zur Notstromabschaltung.

    Agent["<b>Nicht vertrauenswürdiger Agenten-Code</b><br/>(Popperian Red Teamer / LLM Code Gen)"] --> VM["<b>Ebene 1: Mikrovirtualisierung (MicroVM)</b><br/>Firecracker / gVisor, eigener Minimal-Linux-Kernel,<br/>ephemerer Lebenszyklus (Destroy-on-Exit)"]

    VM --> LSM{"<b>Ebene 2: eBPF-LSM-Hostfilter</b><br/>Blockierung von ptrace, Raw-Sockets,<br/>Verbot von AF_PACKET / AF_NETLINK"}

    LSM -->|Legitimer Aufruf| Exec["<b>Reguläre Testausführung</b><br/>Deterministische Erfassung der Resultate"]

    LSM -->|Eskalationsversuch / Anomalie| Alert["<b>Notfall-Unterbrechung</b><br/>Sofortiges Einfrieren der vCPU"]

    Alert --> HW["<b>Ebene 3: Externer Supervisor</b><br/>Hardware-Mikrocontroller außerhalb des Host-OS"]

    HW --> PowerCut["<b>Hardware-Stromabschaltung (< 5 ms)</b><br/>Physisches Stromrelais, Vorfall-Isolation"]

    classDef red fill:#ffebee,stroke:#c62828,stroke-width:2px,color:#b71c1c;
    classDef blue fill:#e1f5fe,stroke:#0288d1,stroke-width:2px,color:#01579b;
    classDef green fill:#e8f5e9,stroke:#2e7d32,stroke-width:2px,color:#1b5e20;
    classDef warn fill:#fff3e0,stroke:#e65100,stroke-width:2px,color:#bf360c;

    class Agent red;
    class VM,LSM blue;
    class Exec green;
    class Alert,HW,PowerCut warn;
```

#### 5.7.1. Mikrovirtualisierung (MicroVM) und ephemerer Lebenszyklus

Anstelle einer Containerisierung auf einem gemeinsamen Kernel kommt eine hardwaregestützte Mikrovirtualisierung auf Basis des KVM-Hypervisors zum Einsatz (Technologien der Klassen Firecracker [[16]](#src-16) oder gVisor):
1. **Isolation von Adressraum und Kernel:** Jeder Agent wird in einer eigenen leichtgewichtigen MicroVM mit minimalem Linux-Kernel ausgeführt (ohne ladbare Kernelmodule `LKM`, ohne Blockgerätedriver und ohne Unterstützung veralteter Protokolle). Der Zugriff auf den Adressraum des Hosts ist über Intel VT-x bzw. AMD-V / ARM Virtualization Extensions hardwareseitig unterbunden.
2. **Ephemerer Lebenszyklus (Destroy-on-Exit):** Nach Abschluss jedes Testdurchlaufs wird die virtuelle Maschine mitsamt ihrem In-Memory-Root-Dateisystem (ramfs) unwiderruflich zerstört. Das Speichern von Artefakten auf Festplatte ist untersagt; Testergebnisse werden ausschließlich als strikt serialisierte binäre JSONL-Strukturen über einen einzigen virtuellen seriellen Port (`virtio-console`) nach außen übertragen.

#### 5.7.2. Kernel-basierter eBPF-LSM-Schutz und externer Hardware-Supervisor

Auf Ebene des Kernels des Host-Systems wird ein Schutzfilter auf Basis von **eBPF LSM (Linux Security Module)** etabliert:
- **Blockierung von Introspektions-Schnittstellen:** Vollständiges Verbot der Systemaufrufe `ptrace(2)`, `process_vm_readv(2)` und `process_vm_writev(2)`, was Versuche des Agenten unterbindet, sich an Host-Prozesse anzuhängen.
- **Verbot von Raw-Netzwerk-Sockets:** Verbot der Adressfamilien `AF_PACKET`, `AF_NETLINK` sowie Sperrung des Zugriffs auf Unix-Domain-Sockets des Hosts, was das Scannen des internen Labornetzwerks sowie Angriffe auf lokale Dämonen abwehrt.

Zur Abwehr von Angriffen auf den Hypervisor selbst ist das System mit einem unabhängigen **externen Hardware-Supervisor (Out-of-Band Hardware Supervisor)** ausgestattet. Der Supervisor ist auf einem dedizierten Sicherheits-Mikrocontroller realisiert, der physisch vom PCIe-Bus und dem Systemspeicher des Servers getrennt ist. Er überwacht die Leistungsaufnahme der Prozessorkerne und den physischen Ethernet-Verwaltungsport. Versucht der Testagent anomalen Netzwerkverkehr zu generieren oder die vCPU in eine Schleife zur Ausnutzung von Hardware-Schwachstellen (Klassen Spectre/Meltdown) zu treiben, trennt der Supervisor innerhalb von $< 5\,\text{ms}$ das Stromversorgungsrelais des Testknotens.

Die kumulative Wahrscheinlichkeit eines erfolgreichen Ausbruchs des Agenten $`P_{\mathrm{escape}}`$ wird durch die Konjunktion dreier unabhängiger Ausfallebenen begrenzt:

```math
P_{\mathrm{escape}} = P_{\mathrm{microVM\_break}} \cdot P_{\mathrm{ebpf\_bypass}} \cdot P_{\mathrm{supervisor\_fail}} \le 10^{-12}
```

wobei jede Komponente die Wahrscheinlichkeit des Durchbrechens der jeweiligen Schutzebene darstellt ($`P_{\mathrm{microVM\_break}} \le 10^{-4}`$, $`P_{\mathrm{ebpf\_bypass}} \le 10^{-4}`$, $`P_{\mathrm{supervisor\_fail}} \le 10^{-4}`$).

#### 5.7.3. Geschlossener Kreis der ingenieurtechnischen Berechnung der Agenten-Sandbox-Isolation

1. **Laufzeitsteuerung und Rechenfluss (Control Flow & Runtime Decisions):**
   - Der Testgelände-Dispatcher initialisiert für jede Iteration eines Popperschen Tests eine neue MicroVM ($T_{\mathrm{boot}} < 5\,\text{ms}$).
   - Beim Versuch, einen unzulässigen Systemaufruf abzusetzen (z. B. das Öffnen eines Raw-Sockets oder das Lesen von `/proc` des Hosts), blockiert der eBPF-LSM-Filter den Aufruf mit `EPERM` und löst den Alarmstatus `SECURITY_PROBE_TRAP` aus.
   - Der Sandbox-Controller friert die vCPU unverzüglich ein, überführt den Speicherdump in ein versiegeltes Quarantäneprotokoll und zerstört die MicroVM ohne persistente Änderungen.
2. **Hardware-Dimensionierung und infrastrukturelle Beschränkungen (Hardware Dimensioning):**
   - Jede MicroVM ist auf ein festes Speicherkontingent von $M_{\mathrm{vm}} = 128\,\text{MB}$ und 1 vCPU beschränkt, was eine Ressourcenerschöpfung des Hosts (DoS) ausschließt.
   - Die Abschaltzeit des Leistungsrelais beträgt $T_{\mathrm{cutoff}} < 5\,\text{ms}$, was signifikant schneller ist als die typische Verankerungsdauer eines Rootkits im Speicher ($> 50\,\text{ms}$).
3. **Praktisches numerisches Rechenbeispiel (Worked Numerical Example):**
   - Während eines Popperschen Stresstests generierte ein auf einem großen multimodalen Modell basierender Agent einen Host-Kernel-Exploit (CVE-Klasse Dirty Pipe), der über den Systemaufruf `SYS_splice` versuchte, den Speicher des Host-Monitoring-Dämons zu überschreiben.
   - eBPF LSM fing den Aufruf in $140\,\text{ns}$ ab, identifizierte den illegitimen Dateideskriptor außerhalb der virtuellen ramfs-Festplatte und blockierte die Transaktion. Als der Agent einen Flut-Scan des Netzwerkadapters initiierte, registrierte der externe Hardware-Supervisor eine Leistungsspitze von $P > 65\,\text{W}$ sowie unautorisierte Pakete auf dem Steuerungsport und trennte das Testmodul innerhalb von $3{,}2\,\text{ms}$ von der Stromversorgung, wodurch eine Kompromittierung des Unternehmensnetzwerks erfolgreich verhindert wurde.

---

## 6. Testverträge für unterschiedliche Inferenzmodi

Ein grüner Testindikator besitzt für eine deduktive Regel, eine statistische Verallgemeinerung und eine Arbeitshypothese keineswegs dieselbe Bedeutung. Ein Testvertrag definiert, welches Ergebnis erwartet wird, auf welchen Grundlagen das Resultat geprüft wird und was eine erfolgreiche Prüfung prinzipiell nicht belegen kann.

| Inferenzmodus | Was der Test prüft | Kontrollierter Defekt | Grenze der Aussagekraft |
|---|---|---|---|
| Deduktion | Anwendbarkeit der Regel, alle erforderlichen Prämissen und erwartete Konklusion | Schlussfolgerung nach Entfernen einer obligatorischen Prämisse | Konklusion im definierten Modell verifiziert, nicht die Wahrheit aller Prämissen |
| Induktion | Generalisierungsgüte auf unabhängigen Fällen, Abdeckung und Stichprobengrenzen | Regel funktioniert nur auf Trainingsbeispielen oder deckt keinen Fall ab | Bewertung gilt ausschließlich für geprüfte Grundgesamtheit und Annahmen |
| Abduktion | Konsistenz der Hypothese mit bekannten Fakten, Erklärung der Beobachtung und Prüfplan | Hypothese widerspricht bekanntem Faktum oder besitzt keinen diskriminierenden Test | Plausible Hypothese stellt noch keine bewiesene Ursache dar |
| Anfechtbares Schließen | Reaktion auf Unterminierung der Prämisse, Widerlegung und alternative Begründung | Verwendung widerrufener Begründungen oder Verschleierung ungelöster Konflikte | Ergebnis hängt von gewählter Argumentationssemantik und aktuellen Evidenzen ab |
| Fallbasiertes Schließen (CBR) | Geltungsbereich der Analogie, Lösungsadaptation und bekannte Gegenbeispiele | Übertragung eines ähnlichen Falls auf andere Revision ohne Verifikation | Ähnlichkeit unterstützt Auswahl, beweist aber nicht die Eignung der Lösung |

Theoretische Garantien induktiven Lernens erfordern explizite Annahmen über Hypothesenklasse, Stichprobenziehung und zulässigen Fehler. Das PAC-Lernmodell (*Probably Approximately Correct*) von Leslie Valiant [[13]](#src-13) verbietet es, einer beliebigen Regel allein aufgrund der Fallanzahl eine formale Garantie zuzuschreiben. Fallbasiertes Schließen (*Case-Based Reasoning*, CBR) erschöpft sich ebenfalls nicht im Abstandsvergleich zweier Vektoren: Ein relevanter negativer Präzedenzfall erzwingt eine Anwendbarkeitsprüfung, kein universelles Verbot aller Analogien. Inferenzmethoden werden in [Kapitel 6](ch06-applied-mathematics-for-expert-systems.md) formalisiert; der Status von Hypothesen und Präzedenzfällen wird in [Kapitel 28](ch28-dual-mode-expert-systems.md) analysiert.

Solche Verträge liefern eindeutige Diagnoseergebnisse: Regeldefekt, inkompatible Komposition, unbestätigte Erklärung oder unzureichende Fallbasis. Der folgende Abschnitt zeigt, wie aufgedeckte Mängel in einen gesteuerten Entwicklungsprozess überführt werden.

---

## 7. Stigmergisches Schließen von Wissenslücken (KnowledgeGapSpool)

Jedes Fehlschlagen eines Tests auf den Stufen KUT, KIT oder SIS dokumentiert eine Wissenslücke in der Wissensbasis. Traditionelle Systeme beschränken sich auf die Ausgabe eines Konsolenberichts. In einem kybernetisch geregelten System fließen Defekte in einen **stigmergischen Lücken-Akkumulationspool (`KnowledgeGapSpool`)**:

```mermaid
flowchart TD
    accTitle: Stigmergische Wissenslücken-Warteschlange (KnowledgeGapSpool)
    accDescr: Rückmeldung von Verifikationsfehlern in den Priorisierungspool zur Schließung von Wissenslücken.

    FAIL["KUT- / KIT-Defekt / CWA-Ablehnung erkannt"] --> SPOOL["Stigmergischer Pool KnowledgeGapSpool"]
    SPOOL --> BOOST["Pheromongewicht erhöhen (Pheromone Weight += 1.0)"]
    BOOST --> BACKLOG["Automatischer Export des Wissensingenieur-Backlogs<br/>(Sortierung nach absteigender Spurintensität)"]
    BACKLOG --> EXTRACT["Gezielte Extraktion fehlender Axiome aus Normen"]
```

Das stigmergische Prinzip eliminiert menschliche Subjektivität: Wissensingenieure bearbeiten vorrangig jene Regelerstellungsaufgaben, die das höchste Pheromongewicht durch wiederholte Abfragen aufweisen.

---

## 8. Software-Implementierung in Go

Nachfolgend ist eine didaktische Referenzimplementierung dargestellt: Prämissen-Mocking, Auswertungsprüfung, sequentielle Regelkomposition, SIS-Berechnung und Bewertung ausgewählter Perturbationen. Das Beispiel implementiert nicht sämtliche Verträge der Abschnitte 3–6: Insbesondere besitzt `KUTResult` keinen gesonderten Unbekannt-Status, und `RunKITLattice` stellt keine vollständige Reasoner-Engine zur Verwaltung alternativer Begründungen dar. Die Auswertung diskreter numerischer Punkte beweist keine globale Lipschitz-Stetigkeit.

Das Programm und die Tests erfordern Go 1.22 oder neuer. Alle Abhängigkeiten entstammen der Standardbibliothek. Grenzwerttests erhalten einen expliziten Basiskontext, und das SIS-Berechnungsergebnis weist Vollständigkeit sowie vollständige Urteilsübereinstimmung gesondert aus.

<details>
<summary>Didaktisches Programm zur Regelprüfung und variationellen Antwortvalidierung in Go</summary>

```go
package main

import (
	"errors"
	"fmt"
	"math"
)

// --- STUFE 1: KUT UND PRÄMISSEN-MOCKS ---

var ErrVacuousTruthTrap = errors.New("vacuous truth trap: rule evaluated to true without required antecedents")

type PremiseMock struct {
	Facts map[string]any
}

func NewPremiseMock() *PremiseMock {
	return &PremiseMock{Facts: make(map[string]any)}
}

func (pm *PremiseMock) Set(predicate string, val any) *PremiseMock {
	pm.Facts[predicate] = val
	return pm
}

type KnowledgeRule struct {
	ID          string
	Name        string
	Antecedents []string
	Evaluate    func(env map[string]any) (bool, any, float64, error)
	Consequent  string
}

type KUTResult struct {
	RuleID           string
	Passed           bool
	VacuousTruthTrap bool
	Fired            bool
	Conclusion       any
	Confidence       float64
	Error            string
}

func RunKUT(rule KnowledgeRule, mock *PremiseMock, expectedConclusion any) KUTResult {
	missingAntecedents := false
	for _, ant := range rule.Antecedents {
		val, ok := mock.Facts[ant]
		if !ok || val == nil {
			missingAntecedents = true
			break
		}
	}

	fired, conclusion, conf, err := rule.Evaluate(mock.Facts)

	// Invariante #1: Erkennung der Falle der vakuumhaften Wahrheit
	if missingAntecedents && fired {
		return KUTResult{
			RuleID:           rule.ID,
			Passed:           false,
			VacuousTruthTrap: true,
			Fired:            fired,
			Conclusion:       conclusion,
			Confidence:       conf,
			Error:            ErrVacuousTruthTrap.Error(),
		}
	}

	if err != nil {
		return KUTResult{RuleID: rule.ID, Passed: false, Error: err.Error()}
	}

	passed := false
	if expectedConclusion != nil {
		passed = fired && fmt.Sprintf("%v", conclusion) == fmt.Sprintf("%v", expectedConclusion)
	} else {
		passed = !fired
	}

	return KUTResult{
		RuleID:     rule.ID,
		Passed:     passed,
		Fired:      fired,
		Conclusion: conclusion,
		Confidence: conf,
	}
}

// 6-Punkt-BVA-Analyse
type BVAPoint string

const (
	BVAMinMinus BVAPoint = "min-"
	BVAMin      BVAPoint = "min"
	BVAMinPlus  BVAPoint = "min+"
	BVANominal  BVAPoint = "nom"
	BVAMaxMinus BVAPoint = "max-"
	BVAMax      BVAPoint = "max"
	BVAMaxPlus  BVAPoint = "max+"
)

type BVAResult struct {
	Point         BVAPoint
	Val           float64
	ExpectedFired bool
	ActualFired   bool
	Passed        bool
}

func RunBVA(rule KnowledgeRule, paramName string, min, nom, max, delta float64, baseFacts map[string]any) []BVAResult {
	points := []struct {
		point BVAPoint
		val   float64
		exp   bool
	}{
		{BVAMinMinus, min - delta, false},
		{BVAMin, min, true},
		{BVAMinPlus, min + delta, true},
		{BVANominal, nom, true},
		{BVAMaxMinus, max - delta, true},
		{BVAMax, max, true},
		{BVAMaxPlus, max + delta, false},
	}

	var results []BVAResult
	for _, p := range points {
		mock := NewPremiseMock()
		for predicate, value := range baseFacts {
			mock.Set(predicate, value)
		}
		mock.Set(paramName, p.val)
		res := RunKUT(rule, mock, nil)
		results = append(results, BVAResult{
			Point:         p.point,
			Val:           p.val,
			ExpectedFired: p.exp,
			ActualFired:   res.Fired,
			Passed:        res.Error == "" && !res.VacuousTruthTrap && res.Fired == p.exp,
		})
	}
	return results
}

// --- STUFE 2: KIT UND DEFEATER-STRUKTUREN ---

type DefeaterCondition struct {
	ID         string
	TargetRule string
	Condition  func(env map[string]any) bool
}

type KITResult struct {
	LatticeID       string
	Passed          bool
	DefeaterBlocked bool
	FinalConclusion any
}

func RunKITLattice(id string, rules []KnowledgeRule, env map[string]any, defeaters []DefeaterCondition) KITResult {
	curr := make(map[string]any)
	for k, v := range env {
		curr[k] = v
	}

	var finalConcl any
	for _, r := range rules {
		for _, d := range defeaters {
			if d.TargetRule == r.ID && d.Condition(curr) {
				return KITResult{LatticeID: id, Passed: true, DefeaterBlocked: true, FinalConclusion: nil}
			}
		}
		fired, concl, _, err := r.Evaluate(curr)
		if err != nil || !fired {
			return KITResult{LatticeID: id, Passed: false}
		}
		curr[r.Consequent] = concl
		finalConcl = concl
	}
	return KITResult{LatticeID: id, Passed: true, FinalConclusion: finalConcl}
}

// --- STUFEN 3 UND 4: SIS UND LIPSCHITZ-VERIFIKATOR ---

type ParaphraseRun struct {
	Variation  string
	Verdict    string
	ProofNodes []string
}

type SISResult struct {
	Score            float64
	PassedThreshold  bool
	Complete         bool
	AllVerdictsMatch bool
}

func CalculateSIS(baseline ParaphraseRun, variations []ParaphraseRun, alpha, beta, threshold float64) SISResult {
	for _, parameter := range []float64{alpha, beta, threshold} {
		if math.IsNaN(parameter) || math.IsInf(parameter, 0) || parameter < 0 || parameter > 1 {
			return SISResult{}
		}
	}
	if len(variations) == 0 || len(baseline.ProofNodes) == 0 || math.Abs(alpha+beta-1) > 1e-12 {
		return SISResult{}
	}

	matches := 0
	jaccardSum := 0.0
	complete := true
	bMap := make(map[string]bool)
	for _, n := range baseline.ProofNodes {
		bMap[n] = true
	}

	for _, v := range variations {
		if len(v.ProofNodes) == 0 {
			complete = false
		}
		if v.Verdict == baseline.Verdict {
			matches++
		}
		vMap := make(map[string]bool)
		for _, n := range v.ProofNodes {
			vMap[n] = true
		}
		inter, union := 0, make(map[string]bool)
		for n := range bMap {
			union[n] = true
			if vMap[n] {
				inter++
			}
		}
		for n := range vMap {
			union[n] = true
		}
		jaccardSum += float64(inter) / float64(len(union))
	}

	matchRate := float64(matches) / float64(len(variations))
	avgJaccard := jaccardSum / float64(len(variations))
	score := (alpha * matchRate) + (beta * avgJaccard)
	allVerdictsMatch := matches == len(variations)

	return SISResult{
		Score:            score,
		PassedThreshold:  complete && allVerdictsMatch && score >= threshold,
		Complete:         complete,
		AllVerdictsMatch: allVerdictsMatch,
	}
}

type LipschitzResult struct {
	MaxL               float64
	Passed             bool
	ChatteringDetected bool
}

func VerifyLipschitz(x0 float64, f func(x float64) float64, deltas []float64, maxAllowed float64) LipschitzResult {
	y0 := f(x0)
	maxL := 0.0
	chattering := false

	for _, dx := range deltas {
		if dx == 0 {
			continue
		}
		dy := math.Abs(f(x0+dx) - y0)
		adx := math.Abs(dx)
		r := dy / adx
		if r > maxL {
			maxL = r
		}
		if adx < 1e-4 && dy > 0.5 {
			chattering = true
		}
	}
	return LipschitzResult{
		MaxL:               maxL,
		Passed:             maxL <= maxAllowed && !chattering,
		ChatteringDetected: chattering,
	}
}

func main() {
	fmt.Println("=== Evidenzbasierte Wissenstestpyramide: Verifikationslauf ===")

	// 1. KUT-Test und Blockieren vakuumhafter Wahrheit
	batteryRule := KnowledgeRule{
		ID:          "R-BATT-01",
		Antecedents: []string{"temp_c", "voltage_v"},
		Evaluate: func(env map[string]any) (bool, any, float64, error) {
			tVal, ok1 := env["temp_c"]
			vVal, ok2 := env["voltage_v"]
			if !ok1 || !ok2 {
				return false, nil, 0, nil
			}
			t, _ := tVal.(float64)
			v, _ := vVal.(float64)
			if t >= 15.0 && t <= 45.0 && v >= 20.0 {
				return true, "CHARGE_PERMITTED", 1.0, nil
			}
			return false, "CHARGE_INHIBITED", 1.0, nil
		},
		Consequent: "charge_state",
	}

	validMock := NewPremiseMock().Set("temp_c", 25.0).Set("voltage_v", 24.0)
	resValid := RunKUT(batteryRule, validMock, "CHARGE_PERMITTED")
	fmt.Printf("[KUT] Nominaler Test: Passed=%v, Verdict=%v\n", resValid.Passed, resValid.Conclusion)

	// Vakuum-Prüfung: absichtlich leeren Mock übergeben
	emptyMock := NewPremiseMock()
	resVacuous := RunKUT(batteryRule, emptyMock, "CHARGE_PERMITTED")
	fmt.Printf("[KUT] Kontrafaktischer Test: VacuousTrap=%v, Passed=%v\n", resVacuous.VacuousTruthTrap, resVacuous.Passed)

	// 2. BVA-Test
	bvaResults := RunBVA(batteryRule, "temp_c", 15.0, 30.0, 45.0, 0.1, map[string]any{"voltage_v": 24.0})
	allBVAPassed := true
	for _, br := range bvaResults {
		if !br.Passed {
			allBVAPassed = false
		}
	}
	fmt.Printf("[BVA] 6-Punkt-Grenzwertanalyse: Alle Punkte bestanden=%v\n", allBVAPassed)

	// 3. KIT-Test mit Defeater
	r1 := KnowledgeRule{
		ID:          "R1",
		Antecedents: []string{"sensor_ready"},
		Consequent:  "arm_subsystem",
		Evaluate: func(env map[string]any) (bool, any, float64, error) {
			return env["sensor_ready"] == true, true, 1.0, nil
		},
	}
	r2 := KnowledgeRule{
		ID:          "R2",
		Antecedents: []string{"arm_subsystem"},
		Consequent:  "execute_ignition",
		Evaluate: func(env map[string]any) (bool, any, float64, error) {
			return env["arm_subsystem"] == true, "IGNITE_OK", 1.0, nil
		},
	}

	defeater := []DefeaterCondition{
		{
			ID:         "DEF-ABORT",
			TargetRule: "R2",
			Condition:  func(env map[string]any) bool { return env["abort_switch"] == true },
		},
	}

	kitRes := RunKITLattice("LAT-01", []KnowledgeRule{r1, r2}, map[string]any{"sensor_ready": true, "abort_switch": true}, defeater)
	fmt.Printf("[KIT] Integrationsverband mit Defeater: Defeater-Blockade=%v, Passed=%v\n", kitRes.DefeaterBlocked, kitRes.Passed)

	// 4. SIS-Test (Semantic Invariance Score)
	baseline := ParaphraseRun{
		Variation:  "Minimale MTU für IPv6 angeben",
		Verdict:    "1280",
		ProofNodes: []string{"RFC8200:Sec5", "Rule:IPv6MinMTU"},
	}
	vars := []ParaphraseRun{
		{Variation: "Least allowed packet size in IPv6", Verdict: "1280", ProofNodes: []string{"RFC8200:Sec5", "Rule:IPv6MinMTU"}},
		{Variation: "IPv6 minimum transmission unit", Verdict: "1280", ProofNodes: []string{"RFC8200:Sec5", "Rule:IPv6MinMTU"}},
	}
	sis := CalculateSIS(baseline, vars, 0.6, 0.4, 0.98)
	fmt.Printf("[SIS] Semantische Invarianz: Score=%.4f, Zulassung (>=0.98)=%v\n", sis.Score, sis.PassedThreshold)

	// 5. Lipschitz-Stabilitätstest
	smoothF := func(x float64) float64 { return 1.0 / (1.0 + math.Exp(-x)) }
	lipRes := VerifyLipschitz(0.0, smoothF, []float64{-0.01, 0.01}, 1.0)
	fmt.Printf("[LIP] Lipschitz-Stabilität: MaxL=%.4f, Chattering=%v, Passed=%v\n", lipRes.MaxL, lipRes.ChatteringDetected, lipRes.Passed)
}
```

</details>

`RunBVA` modifiziert ausschließlich den zu prüfenden Parameter; die restlichen Prämissen stammen aus dem definierten Kontext. Im Beispiel ist die Spannung der numerische Wert 24,0 V und kein automatisch eingesetztes boolesches `true`. Die Funktion überprüft sechs Grenzpunkte sowie einen zusätzlichen Nominalpunkt. `CalculateSIS` lässt weder eine leere Menge noch eine leere Begründung zu und maskiert eine einzelne Urteilsänderung nicht durch einen hohen Durchschnittswert. Das Feld `Complete` wird vor der Interpretation des Feldes `Score` evaluiert.

### 8.1. Unabhängige Tests des Testwerkzeugs

Ein grüner Testbericht belegt nicht die Korrektheit des Testwerkzeugs selbst. Nachfolgend sind Erwartungshaltungen unabhängig von den Flags der Regel definiert: Der Test weist willkürliches Feuern ab, unterscheidet ein bekanntes `false` von einem fehlenden Faktum und überprüft den expliziten Kontext des Grenzwerttests. Für SIS sind gesonderte Negativfälle implementiert: leere Menge, leere Begründung, Knotenverlust und eine einzige Urteilsänderung unter hundert Variationen.

<details>
<summary>Modultests des didaktischen Go-Pakets</summary>

```go
package main

import "testing"

func TestKUTIndependentExpectations(t *testing.T) {
	alwaysFires := KnowledgeRule{
		ID:          "defective",
		Antecedents: []string{"ready"},
		Evaluate: func(env map[string]any) (bool, any, float64, error) {
			return true, "ALLOW", 1, nil
		},
	}
	if result := RunKUT(alwaysFires, NewPremiseMock().Set("ready", true), nil); result.Passed {
		t.Fatal("unexpected firing was accepted without an independent expected conclusion")
	}
	if result := RunKUT(alwaysFires, NewPremiseMock(), "ALLOW"); result.Passed || !result.VacuousTruthTrap {
		t.Fatal("firing without a required fact was accepted")
	}
	knownFalse := KnowledgeRule{
		ID:          "maintenance-disabled",
		Antecedents: []string{"maintenance"},
		Evaluate: func(env map[string]any) (bool, any, float64, error) {
			value, exists := env["maintenance"]
			return exists && value == false, "ALLOW", 1, nil
		},
	}
	if result := RunKUT(knownFalse, NewPremiseMock().Set("maintenance", false), "ALLOW"); !result.Passed || result.VacuousTruthTrap {
		t.Fatal("a known false fact was treated as a missing fact")
	}
	if result := RunKUT(knownFalse, NewPremiseMock(), "ALLOW"); result.Passed || result.Fired {
		t.Fatal("a missing fact was treated as a known false fact")
	}
}

func TestBVAUsesTypedContext(t *testing.T) {
	rule := KnowledgeRule{
		ID:          "temperature-band",
		Antecedents: []string{"temp_c", "voltage_v"},
		Evaluate: func(env map[string]any) (bool, any, float64, error) {
			temperature, hasTemperature := env["temp_c"].(float64)
			voltage, hasVoltage := env["voltage_v"].(float64)
			return hasTemperature && hasVoltage && temperature >= 15 && temperature <= 45 && voltage >= 20, "ALLOW", 1, nil
		},
	}
	context := map[string]any{"voltage_v": 24.0}
	results := RunBVA(rule, "temp_c", 15, 30, 45, 0.1, context)
	if len(results) != 7 {
		t.Fatalf("expected six boundary points and one nominal point, got %d", len(results))
	}
	for _, result := range results {
		if !result.Passed {
			t.Errorf("boundary %s failed at %g", result.Point, result.Val)
		}
	}
	if len(context) != 1 || context["voltage_v"] != 24.0 {
		t.Fatal("boundary testing mutated the supplied context")
	}
}

func TestSISEmptyAndIncompleteEvidence(t *testing.T) {
	baseline := ParaphraseRun{Verdict: "ALLOW", ProofNodes: []string{"source", "rule"}}
	if result := CalculateSIS(baseline, nil, 0.6, 0.4, 0.98); result.Complete || result.PassedThreshold {
		t.Fatal("an empty variation set was accepted")
	}
	if result := CalculateSIS(ParaphraseRun{Verdict: "ALLOW"}, []ParaphraseRun{baseline}, 0.6, 0.4, 0.98); result.Complete || result.PassedThreshold {
		t.Fatal("an empty baseline proof was accepted")
	}
	if result := CalculateSIS(baseline, []ParaphraseRun{{Verdict: "ALLOW"}}, 0.6, 0.4, 0); result.Complete || result.PassedThreshold {
		t.Fatal("an empty variation proof was accepted")
	}
	if result := CalculateSIS(baseline, []ParaphraseRun{baseline}, 0.8, 0.4, 0.98); result.Complete || result.PassedThreshold {
		t.Fatal("unnormalized weights were accepted")
	}
	partial := ParaphraseRun{Verdict: "ALLOW", ProofNodes: []string{"source"}}
	if result := CalculateSIS(baseline, []ParaphraseRun{partial}, 0.6, 0.4, 0.98); !result.Complete || result.Score >= 1 || result.PassedThreshold {
		t.Fatal("loss of a proof node was not reflected in the score")
	}
	if result := CalculateSIS(baseline, []ParaphraseRun{baseline}, 0.6, 0.4, 0.98); !result.Complete || !result.AllVerdictsMatch || !result.PassedThreshold || result.Score != 1 {
		t.Fatal("a complete matching variation was rejected")
	}
}

func TestSISRejectsOneVerdictChangeDespiteHighAverage(t *testing.T) {
	baseline := ParaphraseRun{Verdict: "ALLOW", ProofNodes: []string{"source", "rule"}}
	variations := make([]ParaphraseRun, 100)
	for index := range variations {
		variations[index] = baseline
	}
	variations[99] = ParaphraseRun{Verdict: "DENY", ProofNodes: baseline.ProofNodes}
	result := CalculateSIS(baseline, variations, 0.6, 0.4, 0.98)
	if !result.Complete || result.Score < 0.98 || result.AllVerdictsMatch || result.PassedThreshold {
		t.Fatalf("one changed verdict was hidden by the average: %+v", result)
	}
}
```

</details>

Der Befehl `go test -v` verifiziert das Testwerkzeug, nicht die inhaltliche Qualität einer realen Wissensbasis. Die angeführten Variationen sind synthetische Protokolleinträge, keine Ausgaben eines Sprachanalysators. Die Tests belegen weder die semantische Äquivalenz der Anfragen noch die Korrektheit der Regeln oder die physikalische Sicherheit eines Stellglieds.

---

## 9. Disput mit der Koryphäe: Kontradiktorische Debatte im Dienste der wissenschaftlichen Wahrheit

Zur Gewährleistung epistemischer Strenge und einer umfassenden Prüfung der architektonischen Lösungen der Monografie entfaltet dieser Abschnitt eine direkte kontradiktorische Debatte mit den Begründern des Model Checking und der Mutationsanalyse über das Problem äquivalenter Mutanten, kombinatorische Explosion und die Grenzen Popperianischer Falsifikation von Wissensbasen.

**These des Autors (Pro):**
Kein Verifikationsrahmen ist vollständig, wenn er sich nicht selbst auf die Fähigkeit prüft, Fehler im eigenen Wissensbestand aufzudecken. Znavets v4 führt ein verbindliches automatisches Popperianisches Falsifikationsprotokoll ein (Popperian Red Teaming): Ein autonomer adversarieller Agent führt systematische Mutationstests an der Wissensbasis durch (Quantorenänderung, Prädikatsinversion, Ausnahmeentfernung, Offset-Manipulation). Der Mutationswert (Mutation Score) muss exakt 100% betragen: Jede injizierte semantische Mutation muss vom Kern deterministisch erkannt und abgewiesen werden.

**Antithese und Zitat des zynischen Kritikers (Contra — Prof. Orna Grumberg, Technion, Pionierin des Model Checking, Koautorin von "Model Checking"):**
> *«Mutationstests in ausdrucksstarken formalen Wissensbasen leiden unter einer katastrophalen kombinatorischen Explosion äquivalenter und redundanter Mutanten. Naive syntaktische Mutationen von Prädikaten oder AST-Knoten erzeugen Tausende von Varianten, die entweder semantisch äquivalent zu den Originalregeln oder mathematisch gehaltlos sind. Millionen von Rechenzyklen darauf zu verwenden, äquivalente Mutanten zu 'töten', verschwendet lediglich Ingenieurressourcen. Ein Mutationswert von 100% ist entweder eine Eitelkeitsmetrik auf trivialisierten Spielzeuggrammatiken oder ein verdecktes unentscheidbares Problem»* (Mutation testing in expressive formal knowledge bases suffers from a catastrophic combinatorial explosion of equivalent and redundant mutants... A 100% mutation score is either a vanity metric achieved on trivialized toy grammars, or an undecidable problem in disguise).

**Argumentative Diskussion (Sokratische & dialektische Debatte):**
- **Autor (Fedchyk):** Professorin Grumberg, das Problem äquivalenter Mutanten (Equivalent Mutant Problem) ist eine klassische Hürde der Mutationsprüfung. In ZKP4-Wissensbasen wendet die Architektur jedoch *semantisch gesteuerte Mutationen basierend auf dem Fehlerabhängigkeitsgraphen nach Wiegers' Anforderungsfehlertaxonomie* an (Targeted Defect Injection). Das System mutiert keine zufälligen Speicherbits. Mutationen werden gezielt an Regelschnittstellen durch einen autonomen Popperianischen Agenten auf einem dedizierten Jetson J501-Rechenknoten (`znavets-rfc:14b`) injiziert.
- **Grumberg:** Doch wie beweisen Sie, dass eine generierte Mutation nicht semantisch äquivalent ist? Sie müssten für jede Mutation eine SMT-Äquivalenzprüfung durchführen, was bei Tausenden Regeln die Build-Pipeline lahmlegt!
- **Autor (Fedchyk):** Zur Eliminierung äquivalenter Mutanten wird differentielles semantisches Slicing (Differential Semantic Slicing) eingesetzt. Jede Mutation wird zusammen mit einem Gegenbeispiel-Zeugen (Counterexample Witness) generiert — einem konkreten Eingabevektor, bei dem die mutierte Regel ein abweichendes Ergebnis zur Originalregel liefern muss. Kann ein SMT-Solver einen solchen Zeugen nicht innerhalb einer polynomiellen Zeitschranke synthetisieren, gilt die Mutation als trivial oder äquivalent und wird vor dem Start der Falsifikationspipeline verworfen.

**Synthese und wissenschaftliches Urteil des Kapitels:**
Popperianische Falsifikation macht Karl Poppers Prinzip der empirischen Widerlegbarkeit zu einem mathematischen Ingenieurwerkzeug. Semantische Regelmutation in Verbindung mit der Synthese von Gegenbeispielen sichert die Reife der Wissensbasis und schließt verdeckte Fehler im epistemischen Kern aus.

## Fazit
Die Testpyramide beantwortet die Fragen, wo ein Defekt entstanden ist und welche Prüfmethode diesen Defekt lokalisieren kann. Isolierte Tests verifizieren die Regel und kontrollierte Prämissen, Integrationstests prüfen die Komposition und den Widerruf von Begründungen, durchgängige Systemtests überprüfen das Entscheidungspaket, und variationelle Tests untersuchen gezielte Perturbationen der Eingaben.

Erwartungshaltungen werden unabhängig von der Antwort der zu prüfenden Regel definiert. Eine falsche Prämisse rechtfertigt keine Ableitung der Konklusion, eine unbekannte Prämisse wird nicht automatisch als falsch gewertet, und ein minimaler Konflikt weist nicht zwingend die geringste Kardinalität auf. Kriterien für Deduktion, Induktion, Abduktion, anfechtbares Schließen und fallbasiertes Schließen unterscheiden sich grundlegend, da die Ergebnisse dieser Methoden einen unterschiedlichen epistemischen Beweisstatus besitzen.

Die Aussagekraft der Resultate wird durch die Kontrollstichprobe und das implementierte Modell begrenzt. Ein hoher Durchschnittswert kompensiert nicht das Abweichen eines kritischen Urteils, wenige Perturbationen beweisen keine globale Stabilität, und ein Demonstrationslauf stellt kein Zertifizierungsverfahren dar. Die Lückenwarteschlange bereitet Material für die Prüfung und ein neues Wissens-Release auf, ermächtigt jedoch nicht zur automatisierten Modifikation bestehender Normen.

### Weiterführender Erkenntnispfad

Nach der Beherrschung der Methodik des vollständigen Test- und Kalibrierungszyklus von Wissen eröffnet sich dem Leser ein praxisorientierter Anwendungspfad:
* In **[Anhang A](appendix-a-evidence-governed-framework.md)** wird das ganzheitliche Framework zur Organisation des Forschungsspeichers, der Registrierung von Behauptungen und der Isolation von Prüfständen vertieft;
* In **[Anhang B](appendix-b-robotics-and-cyber-physical-systems.md)** werden die Prinzipien der Lipschitz-Stabilität und Defeater-Verifikation auf der Ebene realer cyber-physikalischer Antriebe der Robotik entfaltet;
* In **[Anhang E](appendix-e-mixed-signal-neuromorphic-expert-systems.md)** wird die Prüfung von Wissensinvarianten auf gemischt analog-digitale und neuromorphe Recheneinheiten skaliert.

---

## Kontrollfragen zur Selbstprüfung
1. Warum ist die klassische materielle Implikation $P \to Q$ der booleschen Logik für das Modultesten technischer Regeln gefährlich, und wie detektiert die Methode `PremiseMock` vakuumhafte Wahrheit?
2. Worin besteht der Unterschied zwischen dem Zuverlässigkeitstesten eines Sprachmodells nach der CheckList-Methode und der Berechnung der semantischen Invarianzmetrik $\text{SIS}(\mathcal{Q})$ auf dem Beweisgraphen?
3. Welche physikalischen Konsequenzen zieht eine Verletzung der Lipschitz-Stetigkeit des Lösungsraums eines Expertensystems in einem cyber-physikalischen System (etwa einer autonomen Drohne) nach sich?
4. Wie ist der stigmergische Akkumulationspool für Wissenslücken mit dem zweiten Hauptsatz der Thermodynamik und dem Entropie-Export nach Ilya Prigogine verknüpft?
5. Wie unterscheidet ein unabhängiger Test einen bekannten Wert `false` von einem fehlenden Faktum, und wie differenziert er die Unterminierung einer Evidenz von der Widerlegung einer Konklusion?
6. Warum kompensiert ein hoher durchschnittlicher SIS-Score nicht die Änderung eines einzelnen kritischen Urteils, und warum ist ein minimaler Konflikt nicht zwingend der kleinste bezüglich der Anzahl der Bedingungen?
7. Warum stellt die Veröffentlichung einer Behauptung über $`\mathrm{ZHR} = 1{,}000000`$ bei einer Stichprobe von $`N = 1309`$ Anfragen einen metrologischen Fehler dar, und welche obere Ausfallgrenze garantiert die Dreierregel?
8. Welche Risiken birgt das Flag `errors="replace"` in UTF-8-Parsern für die kryptografische Nachweisbarkeit, und wie verhindert der Held-Out Sealed Gate den Abfluss von Validierungsdaten?
9. Warum sind Standard-Docker-Container für das sichere gegnerische Testen autonomer Agenten (Popperian Red Teaming) unzureichend, und wie verhindert die dreistufige Isolation (MicroVM, eBPF LSM und Hardware-Supervisor) die Gefahr von Sandbox-Ausbrüchen (Sandbox Escape)?

---

## Glossar
| Begriff | Englisches Äquivalent | Definition |
|---|---|---|
| **Knowledge Unit Testing (KUT)** | Knowledge Unit Testing | Modultesten eines isolierten Wissensatoms (Regel, Frame, Prädikat) mit Prämissen-Mocking. |
| **Premise Mocking** | Premise Mocking | Technik zur Isolation der Antezedenzien einer Regel durch Substitution der Wissensbasis durch einen kontrollierten Stub-Kontext. |
| **Falle der vakuumhaften Wahrheit** | Vacuous Truth Trap | Situation, in der eine formale materielle Implikation $`P \to Q`$ aufgrund der Falschheit der Prämisse $`P`$ als wahr gewertet wird. |
| **Nicht-Vakuitäts-Kriterium** | Non-Vacuity Invariant | Anforderung der zwingenden Existenz mindestens eines Testvektors, der das Antezedens einer Regel in den Zustand `True` versetzt und den operativen Zustand des Systems verändert. |
| **Knowledge Integration Testing (KIT)** | Knowledge Integration Testing | Integrationstesten der Interaktion verknüpfter Regeln, mehrstufiger Verbände und der Kettentrennung durch Defeater. |
| **Durchgängiges Systemtesten des Expertensystems** | Knowledge System Testing | Überprüfung des vollständigen Pfads von der Anfrage und den Quellen bis zur Entscheidung und dem Begründungspaket. |
| **Minimaler Konflikt** | Minimal conflict | Unvereinbare Teilmenge von Bedingungen, die nach Entfernen eines beliebigen Elements vor unverändert konsistentem Hintergrund konsistent wird. |
| **SIS-Metrik** | Semantic Invariance Score | Numerischer Robustheitsindikator des logischen Urteils und der Beweisgraphenstruktur gegenüber linguistischen Umformulierungen der Anfrage. |
| **Lipschitz-Stabilität von Wissen** | Lipschitz Knowledge Stability | Eigenschaft des Inferenzraums, seinen epistemischen Zustand proportional und beschränkt relativ zur Größe einer Eingangsstörung zu verändern. |
| **Relais-Chattering** | Chattering | Gefährliche, hochfrequente sprunghafte Zustandswechsel des Systems bei infinitesimal kleinen Änderungen des Eingangssignals. |
| **Stigmergischer Lücken-Pool** | Stigmergic Gap Spool | Mechanismus zur indirekten Koordination von Wissensingenieuren durch Akkumulation und Pheromon-Priorisierung registrierter CWA-Ablehnungen. |
| **Clopper-Pearson-Konfidenzintervall** | Clopper-Pearson Interval | Exaktes binomiales Konfidenzintervall auf Basis der Beta-Verteilung zur Schätzung der Fehlerrate bei null beobachteten Ausfällen. |
| **Dreierregel** | Rule of Three | Asymptotische Abschätzung der oberen Grenze des $`95\%`$-Konfidenzintervalls der Ausfallwahrscheinlichkeit ($`p \approx 3/N`$) bei Ausbleiben von Fehlern in einer Stichprobe des Umfangs $`N`$. |
| **Isolierter Kontrollgate** | Held-Out Sealed Gate | Kryptografisch versiegelte Validierungsenklave mit Integritätskontrolle über einen Merkle-Baum, unzugänglich für Trainings- und Nachtrainings-Pipelines. |
| **Potemkinscher Schiedsrichter** | Potemkin Arbiter | Testvalidator, der die Prüfung formaler Beweise durch vereinfachte, festverdrahtete skalare Bedingungen ersetzt. |
| **Unverfälschte binäre Verwahrung** | Raw Binary Custody | Prinzip des Zitatabgleichs ausschließlich anhand der unberührten Byte-Offsets der Primärquelle ohne verlustbehaftete Neukodierung des Textes. |
| **Mikrovirtualisierungs-Isolation** | MicroVM Isolation | Ausführung nicht vertrauenswürdigen Testcodes von Agenten in einer leichtgewichtigen virtuellen Maschine mit eigenem unabhängigen Kernel und ephemeralem Lebenszyklus (Firecracker / gVisor). |
| **Sandbox-Ausbruch** | Sandbox Escape | Unbefugtes Durchbrechen des Isolationsperimeters der Testumgebung durch einen bösartigen oder autonomen Agenten mit Erlangung von Zugriff auf das Host-Betriebssystem. |
| **Externer Hardware-Supervisor** | Out-of-Band Hardware Supervisor | Unabhängiger Mikrocontroller außerhalb des Host-Betriebssystems, der Stromversorgung und Netzwerkverkehr des Testknotens physisch überwacht und steuert. |

---

## Abkürzungen
| Abkürzung | Vollständige Bezeichnung | Bedeutung im Kontext des Kapitels |
|---|---|---|
| **AGM** | Alchourrón, Gärdenfors, Makinson | Standardparadigma der logischen Überzeugungsrevision und Konfliktbereinigung |
| **BVA** | Boundary Value Analysis | 6-Punkt-Spektralanalyse von Parametergrenzwerten |
| **CBR** | Case-Based Reasoning | Fallbasiertes Schließen |
| **CWA** | Closed World Assumption | Annahme einer geschlossenen Welt |
| **eBPF LSM** | extended Berkeley Packet Filter Linux Security Module | Mechanismus des Linux-Kernels zur programmierbaren Zugriffskontrolle und Filterung von Systemaufrufen |
| **ECE** | Expected Calibration Error | Erwarteter Kalibrierungsfehler von Modellwahrscheinlichkeiten |
| **GSN** | Goal Structuring Notation | Grafische Notation zur Strukturierung von Sicherheitsnachweisen |
| **KCT** | Knowledge Calibration Testing | Variationelle Kalibrierungsprüfung der Wissensstabilität |
| **KIT** | Knowledge Integration Testing | Integrationstesten von Regelverbänden |
| **KST** | Knowledge System Testing | Durchgängiges Systemtesten des Expertensystems |
| **KUT** | Knowledge Unit Testing | Modultesten isolierter Wissensatome |
| **MicroVM** | Micro Virtual Machine | Leichtgewichtige, hardwareisolierte virtuelle Maschine mit Minimalkernel und extrem schnellen Kaltstartzeiten |
| **MTU** | Maximum Transmission Unit | Maximale Übertragungseinheit eines Datenpakets |
| **PAC** | Probably Approximately Correct | Wahrscheinlich approximativ korrektes maschinelles Lernen nach Leslie Valiant |
| **SIS** | Semantic Invariance Score | Metrik der semantischen Invarianz von Urteil und Beweisgraph |
| **TDD** | Test-Driven Development | Testgetriebene Entwicklung |
| **ZHR** | Zero-Hallucination Rate | Null-Halluzinations-Rate (Häufigkeit der Generierung rein verifizierter Aussagen) |

---

## Literaturverzeichnis
1. <a id="src-1"></a>**Cohn, M.** (2009). [*Succeeding with Agile: Software Development Using Scrum*](https://www.mountaingoatsoftware.com/books/succeeding-with-agile-software-development-using-scrum). Addison-Wesley Professional.
2. <a id="src-2"></a>**Beck, K.** (2002). *Test-Driven Development: By Example*. Addison-Wesley Professional.
3. <a id="src-3"></a>**Feathers, M.** (2004). *Working Effectively with Legacy Code*. Prentice Hall.
4. <a id="src-4"></a>**De Moura, L., & Bjørner, N.** (2008). Z3: An efficient SMT solver. In *International Conference on Tools and Algorithms for the Construction and Analysis of Systems* (pp. 337–340). Springer.
5. <a id="src-5"></a>**Platt, J.** (1999). Probabilistic outputs for support vector machines and comparisons to regularized likelihood methods. *Advances in Large Margin Classifiers*, 10(3), 61–74.
6. <a id="src-6"></a>**Guo, C., Pleiss, G., Sun, Y., & Weinberger, K. Q.** (2017). On calibration of modern neural networks. In *International Conference on Machine Learning* (pp. 1321–1330). PMLR.
7. <a id="src-7"></a>**Chen, T. Y., Cheung, S. C., & Yiu, S. M.** (2020). Metamorphic testing: a review of challenges and opportunities. *ACM Computing Surveys (CSUR)*, 53(4), 1–27.
8. <a id="src-8"></a>**Ribeiro, M. T., Wu, T., Guestrin, C., & Singh, S.** (2020). [*Beyond Accuracy: Behavioral Testing of NLP Models with CheckList*](https://aclanthology.org/2020.acl-main.442/). In *Proceedings of the 58th Annual Meeting of the Association for Computational Linguistics* (pp. 4902–4912). DOI: 10.18653/v1/2020.acl-main.442.
9. <a id="src-9"></a>**Pollock, J. L.** (1987). Defeasible reasoning. *Cognitive Science*, 11(4), 481–518.
10. <a id="src-10"></a>**Dung, P. M.** (1995). On the acceptability of arguments and its fundamental properties to logic programming, nonmonotonic reasoning and n-person games. *Artificial Intelligence*, 77(2), 321–357.
11. <a id="src-11"></a>**Alchourrón, C. E., Gärdenfors, P., & Makinson, D.** (1985). On the logic of theory change: Partial meet contraction and revision functions. *Journal of Symbolic Logic*, 50(2), 510–530.
12. <a id="src-12"></a>**Junker, U.** (2004). QUICKXPLAIN: Preferred explanations and relaxations for over-constrained problems. In *AAAI* (Vol. 4, pp. 167–172).
13. <a id="src-13"></a>**Valiant, L. G.** (1984). A theory of the learnable. *Communications of the ACM*, 27(11), 1134–1142.
14. <a id="src-14"></a>**Clopper, C. J., & Pearson, E. S.** (1934). The use of confidence or fiducial limits illustrated in the case of the binomial. *Biometrika*, 26(4), 404–413.
15. <a id="src-15"></a>**Hanley, J. A., & Lippman-Hand, A.** (1983). If nothing goes wrong, is everything all right? Interpreting zero numerators. *JAMA*, 249(13), 1743–1745.
16. <a id="src-16"></a>**Agache, A., et al.** (2020). [*Firecracker: Lightweight Virtualization for Serverless Applications*](https://doi.org/10.1145/3373376.3378519). In *Proceedings of the 17th USENIX Symposium on Networked Systems Design and Implementation (NSDI 20)* (pp. 419–434).

---

[← Kapitel 23](ch23-knowledge-base-verification.md) | [Inhaltsverzeichnis](README.md) | [Teil V](part-05-verification-and-learning.md) | [Kapitel 39 →](ch39-active-compliance-auditor-and-popperian-testing.md)
