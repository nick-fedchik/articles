# Teil I. Konzeptionelle und epistemische Grundlagen

[← Zurück zum Inhaltsverzeichnis](README.md) | [Zu Teil II →](part-02-knowledge-models.md)

---

## Ziel des Teils

Klar definieren, wann ein Expertensystem erforderlich ist und unter welchen Bedingungen seine Aussagen für ingenieurtechnische Entscheidungen tragfähig sind. Dieser Teil unterscheidet zwischen Daten, Evidenzen, Regeln, Hypothesen und Schlussfolgerungen und regelt die Verantwortlichkeit für Begründungen generierter Antworten.

---

## Abstract und Verknüpfung der Kapitel

Kapitel 1 stellt die praktische Frage nach einer verifizierbaren Antwort. Kapitel 2 definiert den epistemischen Wissenskontrakt: Herkunft, Anwendungsbereich, Gültigkeitsdauer, Zugriffskontrolle und Inferenzmethodik. Kapitel 3 analysiert die Abgrenzung dieses Kontrakts von reinen Auskunftssystemen. Kapitel 4 erläutert die historische Entwicklung der Inferenzverfahren und die Dimensionen von Unsicherheit. Kapitel 5 verknüpft Schlussfolgerungen, Empfehlungen und das institutionelle Gedächtnis konkreter Entscheidungen.

Zentrales Ergebnis: Der Leser kann erforderliche Werkzeugklassen bestimmen und fehlende Begründungsgrundlagen präzise identifizieren. Der Kontrakt garantiert keine inhaltliche Wahrheit fehlerhafter Quellen; mathematische Formalismen und Wissensrepräsentationen werden im folgenden Teil vertieft.

```mermaid
flowchart LR
    accTitle: Logische Abfolge von Teil I
    accDescr: Übergang vom Grundbegriff des Expertensystems über Wissensphilosophie und Suchkritik zur bayesschen Evolution und Gedächtnistriade.

    CH1["<b>Kapitel 1</b><br/>Einführung & Grundbegriffe"] --> CH2["<b>Kapitel 2</b><br/>Philosophie maschinellen Wissens"]
    CH2 --> CH3["<b>Kapitel 3</b><br/>Grenzen von Suchsystemen"]
    CH3 --> CH4["<b>Kapitel 4</b><br/>Evolution: Bayes bis KI"]
    CH4 --> CH5["<b>Kapitel 5</b><br/>Triade des Vertrauens & Gedächtnis"]

    classDef step fill:#e3f2fd,stroke:#1565c0,stroke-width:2px,color:#0d47a1;
    classDef focus fill:#fff3e0,stroke:#e65100,stroke-width:2px,color:#bf360c;

    class CH1,CH3,CH4,CH5 step;
    class CH2 focus;
```

---

## Kapitel des Teils

### [Kapitel 1. Einführung in Expertensysteme: Vom Chaos zu kontrolliertem Wissen](ch01-introduction-to-expert-systems.md)

* **Abstract:** Eine Ingenieursfrage und drei unterschiedliche Antworten: von einer Suchmaschine, von einem generativen Chatbot und von einem evidenzbasierten Expertensystem. Warum Expertensysteme eine Renaissance erleben: Plausibler Text ist zur Massenware geworden, Regulierungen erfordern lückenlose Nachvollziehbarkeit, und Erfahrung verlässt Entwicklungsteams schneller als Produkte. Technische Ursachen des KI-Winters der 1980er-Jahre. Definition von Wissensbasis, Inferenzmaschine, Begründungspaket, fünf Anwendungsdomänen, Haftungsgrenzen und Lesepfade.

### [Kapitel 2. Philosophie für Ingenieure: Was eine Maschine rechtmäßig Wissen nennen darf](ch02-epistemology-of-machine-knowledge.md)

* **Abstract:** Sieben philosophische Fragen als konkrete Prüfungen für automatisierte Aussagen: Begründungsfundament, Objektidentität, Logik, Absicht, Zitatkontext, Behauptungskategorie und explizite Vollmacht. Ein verfahrensbasierter Kontrakt trennt verifizierbares Wissen von plausibler Prosa, ohne fehlerhafte Quellen nachträglich wahr zu machen.

### [Kapitel 3. Wie sich ein Expertensystem von einem Informations- und Auskunftssystem unterscheidet](ch03-beyond-reference-information-systems.md)

* **Abstract:** Die Grenze zwischen Auskunftssystemen und Expertensystemen auf drei Ebenen: Abfragen über Dokumentenkorpora versus Einzelfallentscheidungen, Dokumentauszüge versus Expertenurteile mit Begründungspaket, Such-Pipelines versus Wissensbasen mit Inferenzmaschine. Die Rolle großer Sprachmodelle und Retrieval-Augmented Generation (RAG), Vergleich von sechs Informationssystemklassen und drei Betriebsregeln: kontrollierte Wissensänderung, Konfidenzkalibrierung und autonome Aktionsgrenzen.

### [Kapitel 4. Evolution der Expertensysteme: Vom Satz von Bayes zu evidenzbasierten KI-Entscheidungen](ch04-evolution-from-bayes-to-evidence-ai.md)

* **Abstract:** Die Geschichte von Produktionsregeln, Wahrscheinlichkeiten, Fuzzy-Logik, Wissensgraphen und fallbasiertem Schließen. Für jedes Verfahren werden Unsicherheitstyp und Einsatzgrenzen herausgearbeitet; die mathematischen Details folgen in Kapitel 6. Zentrale Lehren aus früheren KI-Wintern bezüglich Wissenserfassungsengpässen, Brüchigkeit, Wartungsaufwand und Gesamtbetriebskosten.

### [Kapitel 5. Die Triade des Vertrauens: Expertensystem, evidenzbasierte Empfehlung und Unternehmensgedächtnis](ch05-triad-of-trust-and-corporate-memory.md)

* **Abstract:** Das Expertensystem wendet Wissen an, die Empfehlung legt Begründungen und Gültigkeitsgrenzen offen, und das Unternehmensgedächtnis sichert Entscheidungsmotive, Versionen und Verantwortlichkeiten. Eine durchgängige Fallstudie über Teamwechsel veranschaulicht, warum keine der drei Komponenten die beiden anderen ersetzen kann.

---

[← Zurück zum Inhaltsverzeichnis](README.md) | [Zu Teil II →](part-02-knowledge-models.md)
