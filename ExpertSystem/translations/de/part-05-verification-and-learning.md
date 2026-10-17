# Teil V. Verifikation, Testen, Lernen und Sicherheitsbegründung

[← Zu Teil IV](part-04-architecture-and-inference.md) | [Zum Inhaltsverzeichnis](README.md) | [Zu Teil VI →](part-06-frontiers-neuro-symbolic.md)

---

## Ziel des Teils

Bereitstellung eines lückenlosen ingenieurtechnischen Lebenszyklus für Verifikation, Validierung (V&V), kontrollierte Wissensevolution und formale Sicherheitsbegründung von Expertensystemen. Dieser Teil vereint mathematische SMT-Beweise, die mehrstufige Wissenstestpyramide, Prüfungsmatrizen für regressionsfreies Lernen, kontinuierliches Lernen unter Protokolldrift, Popper’sche Falsifikation von Compliance-Anforderungen sowie den Aufbau belastbarer Sicherheitsargumente gemäß den Normen ISO 26262 und ISO/SAE 21434.

---

## Übersicht des Themas und Zusammenhang der Kapitel

Kapitel 23 definiert Methoden der formalen Verifikation der Wissensbasis (Z3-SMT-Solver, Prüfung auf Vollständigkeit und Widerspruchsfreiheit). Kapitel 36 entfaltet die autorenseitige Wissenstestpyramide (KTP) von einzelnen Prädikaten über Integrationspakete bis hin zur variationellen Kalibrierung. Kapitel 25 führt Prüfungsmatrizen, unabhängige Wissensaudits und Regressionskontrollen vor dem Release neuer Revisionen ein. Kapitel 26 gewährleistet kontinuierliches Lernen (Continual Learning) auf Systemprotokollen und erkennt Verteilungsdrift über CUSUM-Indikatoren. Kapitel 39 führt einen aktiven Popper’schen Tester zur automatisierten Generierung von Gegenbeispielen und zur Prüfung der normativen Compliance (ASPICE, ISO 26262, ISO 21434, DO-178C) ein. Kapitel 27 synthetisiert strukturierte Sicherheitsargumente in der GSN-Notation mit lückenloser Rückverfolgbarkeitskontrolle. Kapitel 30 schließt diesen Block mit dem gemeinsamen Entwurf von funktionaler Sicherheit und Cybersicherheit (Co-Engineering) ab.

Das Ergebnis dieses Teils: mathematisch und empirisch nachgewiesene Systemzuverlässigkeit über statische Releases und kontinuierliche Lernzyklen hinweg, gestützt durch rechtsverbindliche und zertifizierungsfähige Zulassungsartefakte. Neuro-symbolische Architekturen, abduktive Schließung von Wissenslücken und Halluzinationskontrolle werden in [Teil VI](part-06-frontiers-neuro-symbolic.md) behandelt.

```mermaid
flowchart TD
    accTitle: Pipeline für Verifikation, Lernen und Sicherheitsnachweis (V&V)
    accDescr: Von der formalen Verifikation und der Testpyramide über Prüfungsmatrizen, kontinuierliches Lernen und das Popper'sche Audit bis zur Sicherheitsbegründung.

    KB["<b>Wissensbasis und Regeln</b>"] --> VER["<b>Kapitel 23</b><br/>SMT-Verifikation und Invarianten"]
    VER --> TEST["<b>Kapitel 36</b><br/>Wissenstestpyramide (KTP)"]
    TEST --> EXAM["<b>Kapitel 25</b><br/>Prüfungsmatrizen und Audit"]
    EXAM --> CL["<b>Kapitel 26</b><br/>Kontinuierliches Lernen & Driftkontrolle"]
    CL --> POPP["<b>Kapitel 39</b><br/>Popper'sche Falsifikation und Compliance"]
    POPP --> GSN["<b>Kapitel 27</b><br/>Synthese des Sicherheitsnachweises (GSN)"]
    GSN --> COENG["<b>Kapitel 30</b><br/>Co-Engineering von Safety und Security"]
    COENG --> ACC["<b>Zertifizierte Zulassung zum kritischen Betrieb</b>"]

    classDef proc fill:#fbe9e7,stroke:#d84315,stroke-width:2px,color:#bf360c;
    classDef base fill:#ede7f6,stroke:#512da8,stroke-width:2px,color:#311b92;
    classDef out fill:#e8f5e9,stroke:#2e7d32,stroke-width:2px,color:#1b5e20;

    class KB base;
    class VER,TEST,EXAM,CL,POPP,GSN,COENG proc;
    class ACC out;
```

---

## Kapitel dieses Teils

### [Kapitel 23. Verifikation der Wissensbasis: Prüfung auf Widerspruchsfreiheit, Vollständigkeit und Robustheit von Regeln](ch23-knowledge-base-verification.md)

* **Abstract:** Herkunft des Orakels, verbotene und zulässige Zustände, Grenzen der vollständigen Aufzählung und des Mutations-Scores. Rapid, Hypothesis, Z3 und TLA+-Werkzeuge verifizieren unterschiedliche formale Eigenschaften. Reduktion kontrolliert spezifische strukturelle Defekte; ein vollständiger formaler Beweis erfordert die Etablierung von Wurzelbehauptungen, Prämissen und allen Vorbedingungen der angewendeten Regel. Induktives Data-Mining schlägt Kandidaten vor, keine normative Wahrheit.

### [Kapitel 36. Wissenstestpyramide: Regeln, Interaktionen und Robustheit der Antworten](ch36-knowledge-testing-pyramid-and-variational-calibration.md)

* **Abstract:** Autorenseitige Testpyramide zur Überprüfung einzelner Regeln, der Regelkomposition, des End-to-End-Entscheidungspakets und von Eingangsvariationen. Unabhängige Orakelerwartungen trennen falsche von unbekannten Prämissen; Integrationsfälle validieren inhaltliche Kompatibilität, Untergrabung der Begründung und minimale Konfliktmengen. Testverträge unterscheiden sich für Deduktion, Induktion, Abduktion, anfechtbares Schließen und Fallbasiertes Schließen. Go-Beispiele und Modultests prüfen vorgegebene Randfallkontexte, leere Mengen, den Verlust von Begründungen und das Verbot, Urteilsänderungen hinter aggregierten Durchschnittswerten zu verschleiern. Metriken und Lückenwarteschlangen stellen keine universelle Zertifizierung einer Ontologie oder eines physikalischen Reglers dar.

### [Kapitel 25. Wie Expertensysteme lernen: Prüfungsmatrizen, Wissensaudits und Regressionskontrolle](ch25-how-expert-systems-learn.md)

* **Abstract:** Gesteuerte Versionsfreigabe von Expertensystemen: Prüfungsmatrizen, paarweiser Vergleich und Zurückweisung unvollständiger Prüfungen. Stratifizierte und zeitliche Aufteilung, Label-Validierungswarteschlangen, Weak-Slice-Mining und LLM-as-a-Judge-Verfahren bereiten Prüfmaterial vor, treffen jedoch keine autonomen Beförderungsentscheidungen. Ein durchgängiges Fallbeispiel von Anforderungsänderung und Berichtswiderruf demonstriert abhängige Fakten, Regeln, Indizes, Caches und Urteile.

### [Kapitel 26. Kontinuierliches Lernen (Continual Learning) aus Erfahrung und Überwindung von Systemprotokoll-Drift](ch26-continual-learning.md)

* **Abstract:** Selbstlernen auf Basis operativer Ablehnungen: Ereignisprotokolle, Label-Verzögerungszeiten und Propensity-Scores. Off-Policy-Evaluation verweigert numerische Bewertungen bei fehlender Überdeckung von Handlungsräumen, während unvollendete Episoden nicht als Fehler gewertet werden. River, Avalanche und Open Bandit Pipeline unterstützen Experimente, doch Kandidaten-Updates durchlaufen dieselben Zulassungsschleusen. Frühindikatoren und CUSUM-Detektoren signalisieren Verteilungsdrift vor der ersten fehlerhaften Antwort.

### [Kapitel 39. Aktiver Compliance-Auditor: Popper’sche Falsifikation, normative Compliance (ASPICE/ISO 26262/ISO 21434) und autonomes Testdesign](ch39-active-compliance-auditor-and-popperian-testing.md)

* **Abstract:** Autonomes Testdesign nach dem Prinzip der Popper’schen Falsifikation: Generierung härtester Gegenbeispiele, Auditierung von Ausführungsspuren auf Konformität mit Standards für funktionale Sicherheit und Cybersicherheit sowie automatische Erkennung blinder Flecken in normativen Anforderungen. Dialektische Verifikation kombinatorischer Synthese: Widerlegung der fünf Thesen des wissenschaftlichen Opponenten (analytische Tautologie, Frame-Problem, Parsing-Degradation $(1-p)^d$, Scheintrugschlüsse und kombinatorische Wand $\mathcal{O}(b^d)$) durch ingenieurtechnische Schutzbarrieren.

### [Kapitel 27. Sicherheitsbegründung: Synthese und Verifikation von Argumenten](ch27-safety-case-gsn-synthesis.md)

* **Abstract:** Sicherheitsbegründung in der Goal Structuring Notation (GSN), synthetisiert aus dem ingenieurtechnischen Wissensgraphen. Formale Regeln für Argumentvollständigkeit und Link-Integrität, Prüfung von Hardware-Revisionen in Kontexten und Evidenzen, selektive Offenlegung von Nachweisen über Merkle-Bäume mit Salt für Blätter sowie Berücksichtigung von Einwänden im Rahmen von Dungs Argumentationsrahmen mit unaufgelösten wechselseitigen Angriffen. Das SACM-Metamodell, in-toto- und SLSA-Attestierungen, Transparenzprotokolle und die Extraktion von Rückverfolgbarkeitsbeziehungen ergänzen das Lehrbeispiel, ersetzen jedoch nicht den Fachexperten, der die Hinlänglichkeit des Arguments beurteilt.

### [Kapitel 30. Co-Engineering von funktionaler Sicherheit und Cybersicherheit](ch30-safety-cybersecurity-co-engineering.md)

* **Abstract:** Harmonisierung von funktionaler Sicherheit und Cybersicherheit anhand verifizierter kausaler Beziehungen und konkreter Produktanforderungen. Lehrreiche Go-Validatoren prüfen Sicherheitskategorien, Zeitbudgets, genehmigte übergeordnete Artefakte, dreiwertige Aktualisierungsregeln (Konflikt, unzureichende Evidenz, kein Konflikt) und die Bindung digitaler Signaturen an Nutzdaten und Metadaten. STPA-Sec, ReqIF, CycloneDX, OpenVEX und Uptane werden als maschinenlesbare Evidenzquellen mit expliziten Grenzen vorgestellt. Das Projektprofil ersetzt keinen Industriestandard; die Qualifizierung von Softwarewerkzeugen bleibt ein eigenständiges Verfahren.

---

[← Zu Teil IV](part-04-architecture-and-inference.md) | [Zum Inhaltsverzeichnis](README.md) | [Zu Teil VI →](part-06-frontiers-neuro-symbolic.md)
