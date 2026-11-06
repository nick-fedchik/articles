# Architektur evidenzbasierter Expertensysteme

## Zusammenfassung

Dieses Buch untersucht, wie technische Dokumente, Anforderungen, Quellcode und Prüfergebnisse in strukturiertes Wissen für nachvollziehbares maschinelles Schlussfolgern überführt werden können. In sieben Teilen und vierzig Kapiteln behandelt es die Formalisierung von Begriffen und Quellen, Wissensmodelle und Wissenserwerb, Rückverfolgbarkeit, die Architektur von Expertensystemen, Erklärungen, Verifikation, Sicherheit und neuro-symbolische Ansätze.

Im Mittelpunkt stehen die Herkunft und Gültigkeit von Belegen sowie die nachvollziehbaren Beziehungen zwischen Anforderungen, Regeln, Entwurfsentscheidungen und Tests. Das Buch unterscheidet zwischen einer vorgeschlagenen Hypothese, einer innerhalb eines formalen Modells abgeleiteten Schlussfolgerung und einem experimentell bestätigten Ergebnis; bloße Plausibilität gilt nicht als Beleg. Weitere Themen sind logisches Schließen, Unsicherheit, kausale und zeitliche Modelle, Erklärungen, kontrollierte Aktionen, kontinuierliches Lernen und Sicherheitsargumentationen. Sprachmodelle dienen als Hilfsmittel; ihre Vorschläge müssen unabhängig geprüft werden.

Das Buch richtet sich an Studierende, Ingenieurinnen und Ingenieure sowie Forschende, die Entscheidungsunterstützungssysteme in komplexen technischen Umgebungen entwerfen oder bewerten. Es zeigt, wie Systeme gestaltet werden können, die ihre Schlussfolgerungen begründen, deren Grenzen erläutern und sich bei unzureichender Beweislage einer Antwort enthalten.

**Kapitelzusammenfassungen:** [Alle 40 Kapitel](#kapitelzusammenfassungen).

**Schlüsselwörter:** Expertensysteme, Wissensengineering, Herkunft von Belegen, Rückverfolgbarkeit, formale Verifikation, funktionale Sicherheit, neuro-symbolische KI.

[Ukrainische Masterausgabe](../../README.md)

---

## Kapitelzusammenfassungen

Kurze Zusammenfassungen aller 40 Kapitel der ukrainischen Masterausgabe. Jede Kapitelnummer führt zum Ausgangstext.

### 01 · [Kapitel 1](../../ch01-introduction-to-expert-systems.md)

Das Kapitel definiert ein Expertensystem durch explizites Wissen, nachvollziehbares Schließen und begründete Enthaltung. Es skizziert Geschichte und Architektur des Fachgebiets und grenzt verantwortbares Schlussfolgern von Dokumentensuche und freier Textgenerierung ab.

### 02 · [Kapitel 2](../../ch02-epistemology-of-machine-knowledge.md)

Untersucht werden die Grundlagen maschineller Aussagen: Quelle, Kontext, Geltungsbereich, Unsicherheit und Gültigkeit. Widersprüchliche Temperaturgrenzen zeigen, warum eine korrekte Zahl allein keine Entscheidung rechtfertigt.

### 03 · [Kapitel 3](../../ch03-beyond-reference-information-systems.md)

Das Kapitel unterscheidet Referenz-, Such-, Dialog-, Empfehlungs- und Expertensysteme nach Aufgabe und Ergebnisvertrag. Eine Regel zu finden ist nicht dasselbe, wie sie auf einen konkreten Fall korrekt anzuwenden.

### 04 · [Kapitel 4](../../ch04-evolution-from-bayes-to-evidence-ai.md)

Es zeichnet die Entwicklung der Expertensysteme nach und vergleicht Logik, Wahrscheinlichkeit, Fuzzy-Mengen, Graphen und fallbasiertes Schließen. Die Methodenwahl richtet sich nach Fragestellung, erforderlichen Belegen und Schlussgrenzen.

### 05 · [Kapitel 5](../../ch05-triad-of-trust-and-corporate-memory.md)

Vorgestellt wird ein Vertrauensdreiklang: Das System wendet Wissen an, die Empfehlung legt ihre Gründe offen und das organisationale Gedächtnis bewahrt Kontext und Folgen. Behandelt werden Herkunft, Audit, Zugriffsrechte und Wissensprüfung.

### 06 · [Kapitel 6](../../ch06-applied-mathematics-for-expert-systems.md)

Anhand technischer Aufgaben werden logisches und probabilistisches Schließen, Graphmethoden, Optimierung und Kausalanalyse eingeführt. Reproduzierbare Berechnungen und genaue Zitate setzen dennoch gültige Eingaben und eine korrekte Auslegung voraus.

### 07 · [Kapitel 7](../../ch07-knowledge-base-typology.md)

Verglichen werden Regeln, Objektmodelle, Ontologien, Fälle, Nebenbedingungen, probabilistische Modelle sowie Dokument- und Vektorindizes. Das Kapitel erläutert, was die jeweiligen Repräsentationen speichern und warum fehlendes Wissen ausdrücklich sichtbar bleiben muss.

### 08 · [Kapitel 8](../../ch08-engineering-artifacts-as-data.md)

Anforderungen, Entscheidungen, Pläne und Testergebnisse werden als strukturierte Artefakte mit typisierten Feldern und Beziehungen modelliert. Provenienz ermöglicht Abfragen und Analysen, belegt aber weder Wahrheit noch rechtliche Ausreichendheit.

### 09 · [Kapitel 9](../../ch09-engineering-knowledge-graph-traceability.md)

Ein typisierter Wissensgraph verfolgt Änderungen von Anforderungen über Software und Hardware bis zu Tests. Abfragen unterstützen die Auswirkungsanalyse; Graphverbindungen allein beweisen jedoch keine Korrektheit.

### 10 · [Kapitel 10](../../ch10-knowledge-acquisition-systems.md)

Beschrieben wird eine Wissenserfassungspipeline, die Quellen findet, Inhalt und Herkunft bewahrt, Wissenskandidaten prüft und ihren Geltungsbereich festlegt. Zitattreue, Detektor-Konfidenz und Dokumentfreigabe beantworten unterschiedliche Fragen.

### 11 · [Kapitel 11](../../ch11-knowledge-elicitation-from-experts.md)

Interviews, Fallvergleiche, die Rekonstruktion von Entscheidungen, kognitive Karten und Gegenbeispiele machen implizites Expertenwissen zugänglich. Häufigkeit, Schätzunsicherheit und Übereinstimmung zwischen Fachleuten gelten als unterschiedliche Messgrößen.

### 12 · [Kapitel 12](../../ch12-linguistic-analysis-and-local-models.md)

Das Kapitel zeigt, wie linguistische Analyse und lokale Modelle Regelkandidaten bilden und dabei Quellbytes, Bedingungen, Negationen und Datenschutz bewahren. Ein genaues Zitat allein garantiert noch kein korrektes Verständnis.

### 13 · [Kapitel 13](../../ch13-language-variability-vs-determinism.md)

Die Suche nach Kandidaten wird von der deterministischen semantischen Prüfung getrennt. Unterschiedliche Formulierungen werden typisierten Abfragen zugeordnet, wobei Einheiten, Negation, Modalität, Quellversion und Belege geprüft werden.

### 14 · [Kapitel 14](../../ch14-requirements-detection-and-formalization.md)

Behandelt werden Modalität und Geltungsbereich von Anforderungen, ihre Überführung in formale Kandidaten sowie die Suche nach Widersprüchen. EARS-Muster und begrenzte Werkzeuge unterstützen die Prüfung, ersetzen sie aber nicht.

### 15 · [Kapitel 15](../../ch15-knowledge-extraction-and-kb-construction.md)

Das Kapitel verfolgt den Weg von der Quelle über Kandidaten und freigegebene Wissensdatenbankeinträge bis zum überprüfbaren Schluss. Es behandelt Extraktionsschritte, Annotationskoordinaten, Einheiten und Importprüfungen; vorgeschlagene Methoden müssen noch evaluiert werden.

### 16 · [Kapitel 16](../../ch16-expert-systems-architecture.md)

Erläutert werden die Grenzen zwischen Wissensspeichern, Inferenz, Verifikation, Zugriffsrichtlinien und Provenienzprotokollen. Behandelt werden Aktualisierungen und Widerrufe; Suchbewertungen gelten nicht als Belege.

### 17 · [Kapitel 17](../../ch17-implementation-stack.md)

Die Auswahl des Technologie-Stacks beginnt bei benötigten Operationen, Korrektheitskriterien und Fehlerverhalten. Verglichen werden Sprachen, Regelwerke, Speicher und Bereitstellung; Fehler müssen ausdrücklich behandelt werden.

### 18 · [Kapitel 18](../../ch18-execution-infrastructure.md)

Die Architektur wird mit physischen Ausführungsgrenzen verknüpft: Speicher, CPUs und Beschleuniger, Busse, Latenz und Platzierung der Komponenten. Ressourcenabschätzungen hängen von den Bedingungen ab und sind keine universellen Leistungszusagen.

### 19 · [Kapitel 19](../../ch19-from-question-to-evidence.md)

Eine Anfrage wird zu einer geprüften Aussage, indem Quellen gefunden und an die richtige Version, das richtige Objekt, die Einheit und die Bedingung gebunden werden. Suchrelevanz allein stützt keinen Schluss.

### 20 · [Kapitel 20](../../ch20-explanation-engine.md)

Aus Entscheidungsspuren werden Erklärungen zu „wie“, „warum“, „warum nicht“, „was wäre wenn“ und erforderlichen Änderungen abgeleitet. Behandelt werden minimale Konflikte und rollenbasierte Offenlegung; Sprachmodelle dürfen keine Gründe erfinden.

### 21 · [Kapitel 21](../../ch21-from-recommendation-to-action.md)

Eine Empfehlung wird von der Befugnis zur Handlung getrennt. Ein vorgeschlagenes Aktions-Gateway prüft Handlung, Akteur, Kontext und Richtlinie und protokolliert die Ausführung; das Design wird nicht als validiert dargestellt.

### 22 · [Kapitel 22](../../ch22-cybernetics-edge-to-backend.md)

Ein industrielles Sensorszenario erläutert kontextabhängige kybernetische Rückkopplung von der Erfassung über die Edge-Analyse bis zur Reaktion. Modellwerte allein bestimmen keine Handlung; Kalibrierung und Prozesszustand sind entscheidend.

### 23 · [Kapitel 23](../../ch23-knowledge-base-verification.md)

Unterschieden werden die Verifikation der implementierten Wissensbasis und die Validierung der Richtlinie anhand der realen Aufgabe. Endliche Modelle, Gegenbeispiele und Freigabeschranken decken Auslassungen auf, beweisen aber keine universelle Korrektheit.

### 24 · [Kapitel 24](../../ch24-system-diagnosis.md)

Modellbasierte Diagnose vergleicht Fehlerhypothesen mit Beobachtungen bei unvollständigen Daten. Konsistenz schränkt Kandidaten ein, belegt aber ohne zusätzliche Evidenz keine physische Grundursache.

### 25 · [Kapitel 25](../../ch25-how-expert-systems-learn.md)

Prüfmatrizen, Wissensaudits, Regressionstests und Kalibrierungsprüfungen unterstützen kontrollierte Aktualisierungen. Sie erkennen Verschlechterungen in geprüften Szenarien, garantieren jedoch keine Korrektheit in allen künftigen Fällen.

### 26 · [Kapitel 26](../../ch26-continual-learning.md)

Das Kapitel erklärt Rückkopplungsverzerrungen, wenn Systemempfehlungen die zum Lernen verwendeten Protokolle prägen. Überwachung und kontrollierte Aktualisierungen verhindern, dass wiederholte Entscheidungen als unabhängige Evidenz missverstanden werden.

### 27 · [Kapitel 27](../../ch27-safety-case-gsn-synthesis.md)

Nachvollziehbare Belege werden mit einer strukturierten Sicherheitsargumentation zu einem prüfbaren Fall verbunden. Das Ergebnis ist ein Argumentationskandidat, kein automatisches Zertifikat; Annahmen und Lücken müssen geprüft werden.

### 28 · [Kapitel 28](../../ch28-dual-mode-expert-systems.md)

Strenge, auf verifiziertem Wissen beruhende Schlüsse werden von beratenden Hypothesen über nächste Prüfungen getrennt. Explizite Modi, Kennzeichnungen und Zulassungskontrollen verhindern, dass Vorschläge als Beweise erscheinen.

### 29 · [Kapitel 29](../../ch29-neuro-symbolic-architecture.md)

Ein Sprachmodell interpretiert Anfragen und schlägt Textstellen vor, während ein symbolischer Kern Struktur, Herkunft, Regeln und Ausnahmen prüft. Beleg für einen Schluss und Handlungsbefugnis bleiben getrennte Fragen.

### 30 · [Kapitel 30](../../ch30-safety-cybersecurity-co-engineering.md)

Funktionale Sicherheit und Cybersicherheit werden gemeinsam analysiert, wenn eine Schutzmaßnahme selbst ein physisches Risiko erzeugen kann. Gemeinsamer Kontext und Folgenanalyse machen Zielkonflikte sichtbar, ohne eine Disziplin über die andere zu stellen.

### 31 · [Kapitel 31](../../ch31-syllogistic-reasoning-and-relation-lattices.md)

Behandelt wird normatives Schließen mit Prädikathierarchien, Ausnahmen, Widerlegungen und zeitlicher Gültigkeit. Das Kapitel unterscheidet die Suche nach relevanten Textstellen von einem kontextbezogenen, regelgeleiteten Schluss.

### 32 · [Kapitel 32](../../ch32-high-performance-knowledge-packs-mmap-and-harvesting.md)

Vorgestellt werden unveränderliche Wissenspakete, bytegenaue Zulassung, Indizes und Memory-Mapping für reproduzierbaren, effizienten Zugriff. Hashes schützen die Integrität, belegen aber weder Wahrheit noch Anwendbarkeit des Inhalts.

### 33 · [Kapitel 33](../../ch33-inter-system-knowledge-exchange-and-model-teaching.md)

Sicherer Wissensaustausch überträgt mit jedem Wert auch Einheiten, Geltungsbereich, Regelversionen, Quellen und Berechtigungen. Signaturen schützen die Herkunft, nicht die Korrektheit oder Anwendbarkeit.

### 34 · [Kapitel 34](../../ch34-deterministic-relational-analysis-abduction-and-socratic-dialogue.md)

Relationale Suche und Abduktion liefern mögliche Erklärungen bei fehlenden Fakten; gezielte Rückfragen sollen anschließend unterscheidende Belege einholen. Hypothesen bleiben von festgestellten Tatsachen getrennt.

### 35 · [Kapitel 35](../../ch35-self-organizing-expert-systems-synergetics-and-npu-runtime.md)

Reaktive Systeme überarbeiten Schlussfolgerungen, wenn stützende Beobachtungen widerrufen werden. Unveränderliche Regeln werden von veränderlichen Daten getrennt; adaptive Parameter sind Vorschläge, die Grenzen und Prüfung benötigen.

### 36 · [Kapitel 36](../../ch36-knowledge-testing-pyramid-and-variational-calibration.md)

Das Testpyramidenkonzept wird auf Wissenssysteme übertragen: geprüft werden Regeln, ihr Zusammenspiel, der vollständige Antwortpfad und Eingabeänderungen. Das Kapitel erläutert, was jede Ebene zeigen kann und was sie nicht beweist.

### 37 · [Kapitel 37](../../ch37-input-information-assessment-and-algorithmic-skepticism.md)

Eingehende Belege werden vor der Regelanwendung nach Objektidentität, Zeit, Herkunft, Unabhängigkeit und Unsicherheit bewertet. Ein Verdacht löst eine Prüfung aus, ist aber kein Fälschungsbeweis.

### 38 · [Kapitel 38](../../ch38-curing-machine-hallucinations-and-knowledge-deficits.md)

Vorschläge eines Sprachmodells werden durch Prüfungen von Quelle, Bedeutung, Anwendbarkeit und Voraussetzungen von akzeptierten Schlüssen getrennt. Authentische Zitate sind für die Rückverfolgbarkeit nötig, beweisen aber keinen Schluss.

### 39 · [Kapitel 39](../../ch39-active-compliance-auditor-and-popperian-testing.md)

Ein aktiver Auditor sucht Gegenbeispiele und fehlende Grenzfälle zu Anforderungen. Testgenerierung wird von der Festlegung erwarteter Ergebnisse getrennt; Compliance-Abdeckung ist prüfbar, aber kein automatischer Beweis.

### 40 · [Kapitel 40](../../ch40-distributed-epistemic-architectures-soa-and-cluster-scaling.md)

Beschrieben werden verteilte Expertensysteme mit semantischem Routing, vollständigem aktiven Regelsatz, abgestimmtem Speicher und Entscheidungsprotokollen. Fehlende Shards oder veraltete Versionen müssen einen ungeklärten Zustand statt stiller Freigabe ergeben.
