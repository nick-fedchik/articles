# Kapitel 22. Kybernetischer Regelkreis: Sensoren, Peripherie und Rückkopplung

> **Buch:** [Architektur evidenzbasierter Expertensysteme](README.md) · [Teil VII: Technologie-Stack, Hardware-Infrastruktur, kybernetische Ausführung und verteilte SOA](part-07-runtime-and-knowledge-exchange.md)  
> **Vorheriges Kapitel:** [Kapitel 21. Von der Empfehlung zur Aktion: Autoritätskontrolle und sichere Ausführung in Produktionsumgebungen](ch21-from-recommendation-to-action.md)  
> **Nächstes Kapitel:** [Kapitel 35. Reaktives Expertensystem: Ereignisse, Widerruf und Wissensadaptation](ch35-self-organizing-expert-systems-synergetics-and-npu-runtime.md)  
> **Inhaltsverzeichnis:** [README.md](README.md)  
> **Autor:** [Mykola Fedchyk](about-the-author.md)  
> **Lernziele:** Das Gesetz der erforderlichen Varietät nach Ashby und das Good-Regulator-Theorem auf die Architektur von Expertensystemen anwenden; den Objektzustand über Filter statt Rohsignale schätzen; Funktionen über Zeithorizonte der Ebenen L0–L4 verteilen; kryptografisch signierte Merkmale statt bloßer Urteile von Sensoren übertragen; Degradierungsmodi bei Verbindungsabbruch und Zustandsabgleich nach Wiederherstellung entwerfen; Wissensbestände auf Edge-Geräten über signierte Bundles aktualisieren.  
> **Niveau:** Fortgeschrittene Ingenieurpraxis: Entwickler cyber-physischer Systeme, Embedded-Ingenieure, Systemarchitekten

## Abstract

Auf einer automatisierten Fertigungslinie registriert ein hochfrequenter Schwingungssensor eine atypische Vibrationsspitze in einer Lagerbaugruppe. Wie ist darauf zu reagieren? Es bestünde die Option, den Antriebsmotor unverzüglich abzuschalten, dem Wartungsingenieur eine Inspektionsbenachrichtigung für die nächste turnusmäßige Pause zuzustellen oder das Signal als kurzzeitige Störung durch eine benachbarte Presse zu verwerfen. Die numerische Ausgabe eines neuronalen Klassifikators liefert für sich genommen keine Handlungsanweisung. Die Entscheidung hängt vom Gesamtkontext ab: der Phase des technologischen Prozesses, der Verlässlichkeit der Sensorkalibrierung, der Schmiermitteltemperatur, den Kosten eines Bandstillstands, der Verfügbarkeit eines redundanten Kanals sowie den Befugnissen des entscheidenden Knotens.

Reißt die Netzwerkverbindung zum zentralen Server ab, muss die Schutzlogik unmittelbar an der Anlage auf der Peripherie (*Edge*) greifen. Gleichzeitig besitzt die lokale Steuerung weder Einblick in die vollständige Wartungshistorie des gesamten Maschinenparks noch das Recht, den technologischen Ablaufplan des Werks eigenmächtig zu modifizieren. Dieses Verhalten ist typisch für cyber-physische Systeme, in denen Rechenoperationen untrennbar mit physikalischen Prozessen verwoben sind und deren Dynamik in harter Echtzeit folgen müssen [[1]](#src-1).

Dieses Kapitel beantwortet die zentrale Frage: **Wie lässt sich die Rückkopplung zwischen dem physikalischen Prozess, den Edge-Knoten und dem zentralen Expertensystem so schließen, dass die Steuerung auch bei Latenzen, Sensorrauschen und Verbindungsabbrüchen funktional sicher bleibt?** Die Kernthese lautet: **Der Regelkreis wird nach Zeithorizonten modularisiert. Schnelle Schutzreaktionen werden von deterministischen Steuerungen direkt an der Anlage über Hardware-Begrenzer des Betriebsbereichs (Operating Envelope) und Festkommaarithmetik ausgeführt, die softwaredefinierte physikalische Schicht (AI-RAN PHY) liefert sensorische Nachweise direkt an die Inferenzlogik, Bytecode-Interpreter der EISA garantieren einen bitäquivalenten kreuzarchitektonischen Isomorphismus von Entscheidungen auf heterogenen Chips, der Edge-Knoten schätzt den Objektzustand und wendet gecachte Regeln an, während das zentrale Expertensystem Wissensbestände ausschließlich über kryptografisch signierte Bundles aktualisiert. Jede Ebene verfügt über eine ausreichende Reaktionsvarietät für ihre spezifischen Störungen, ein internes Modell des Objekts, einen eigenen Degradierungsmodus sowie formale Regeln für den Zustandsabgleich nach Wiederherstellung der Verbindung.**

## 1. Kybernetik und Ashbys Gesetz der erforderlichen Varietät

Norbert Wiener definierte die Kybernetik als Wissenschaft von Steuerung und Kommunikation im Lebewesen und in der Maschine, fundiert auf dem Konzept der geschlossenen Rückkopplung [[2]](#src-2). Für ein Expertensystem ergibt sich aus diesem Bezugsrahmen eine fundamentale praktische Fragestellung: Ist das System in der Lage, alle Zustände des Objekts zu differenzieren, auf die es mit unterschiedlichen Aktionen reagieren muss? Die quantitative Antwort liefert William Ross Ashbys Gesetz der erforderlichen Varietät (*Law of Requisite Variety*) [[3]](#src-3).

Unter der Varietät einer endlichen Zustandsmenge $\Omega$ verstand Ashby die Anzahl ihrer unterscheidbaren Elemente; in logarithmischer Form ausgedrückt:

```math
V(\Omega) = \log_2 \lvert \Omega \rvert
```

- Zur Messung der Varietät ist $\Omega$ eine endliche Zustandsmenge und $\lvert \Omega \rvert$ die Anzahl ihrer unterscheidbaren Elemente;
- $\log_2$ ist der Logarithmus zur Basis 2, und $V(\Omega)$ wird in Bit gemessen;
- Die Gleichung spezifiziert die Anzahl an Bits, die erforderlich ist, um die Zustände der Menge innerhalb dieses Modells eindeutig zu differenzieren.

Umfasst eine Menge beispielsweise 1.024 Zustände, beträgt $V(\Omega) = \log_2(1024) = 10\text{ Bit}$. Diese Metrik erfasst ausschließlich die Anzahl diskreter Zustände, nicht jedoch deren interne Komplexität oder Auftretenswahrscheinlichkeit.

Das Gesetz der erforderlichen Varietät postuliert in dieser Form, dass die Varietät der resultierenden Ausgänge $V_O$ nicht geringer sein kann als die Varietät der Störungen $V_D$ abzüglich der Varietät der Reaktionen des Reglers $V_R$:

```math
V_O \geq V_D - V_R
```

Gesetzliche Begrenzung:

- $V_O$ bezeichnet die Varietät der Ergebnisse, $V_D$ die Varietät der Störungen und $V_R$ die Varietät der Reglerreaktionen;
- Jede Größe wird in Bit innerhalb desselben Zustandsmodells gemessen;
- $\geq$ bedeutet „nicht kleiner als“, und das Minuszeichen subtrahiert das Absorptionsvermögen des Reglers von der Varietät der Störungen selbst.

Betragen in einem didaktischen Szenario beispielsweise $V_D = 10\text{ Bit}$ und $V_R = 1\text{ Bit}$, so muss die unregulierte Restvarietät der Ausgänge mindestens $9\text{ Bit}$ betragen. Diese Ungleichung definiert eine theoretische Untergrenze und legt für sich genommen noch nicht fest, welche konkreten Reaktionen ein sicheres Systemverhalten garantieren.

Unter der Annahme, dass 1.024 Störfälle jeweils distinkte Gegenmaßnahmen erfordern, ist ein Regler mit lediglich zwei Reaktionen unzureichend. Die Zahl der Klassenbezeichnungen in einem Klassifikator ist jedoch nicht identisch mit der Varietät verfügbarer Steuerbefehle. Die Ungleichung limitiert die Anzahl unterscheidbarer Resultate, zählt aber nicht automatisch kritische oder „unkontrollierbare“ Zustände ab. Eine einzige robuste Schutzmaßnahme kann für eine Vielzahl unterschiedlicher Fehlzustände angemessen sein. Eine Ontologie unterstützt zwar die semantische Kontextdifferenzierung, erhöht jedoch ohne entsprechende Beobachtbarkeit und reale Stellglieder die Regelungskapazität des Systems nicht.

Das Good-Regulator-Theorem von Conant und Ashby [[4]](#src-4) wurde für ein präzise definiertes Regelungsmodell und spezifische Optimalitätsbedingungen formuliert. Es erzwingt nicht, dass jeder industrielle Regler zwingend über einen vollumfänglichen digitalen Zwilling oder exakt ein Kalman-Filter verfügen muss [[5]](#src-5). In der hier vorgestellten Architektur stellt ein Zustandsschätzer jedoch eine fundierte Wahl dar, wenn verborgene Variablen aus verrauschten Sensordaten rekonstruiert werden müssen. Seine Eignung wird anhand des Modells, der Fehlertoleranz und des Echtzeit-Zeitbudgets verifiziert. Das nachfolgende Diagramm veranschaulicht diesen Regelkreis.

```mermaid
flowchart LR
    accTitle: Regelkreis mit Rückkopplung
    accDescr: Der physische Prozess liefert Signale an die Sensoren, Signale werden gefiltert, der Zustandsschätzer bildet das Objektmodell, das Expertensystem fällt eine autorisierte Entscheidung, der Aktor modifiziert den Prozess und der Operator überwacht die Befugnisse.

    ENV["<b>Physischer Prozess</b><br/>Mechanik, Thermodynamik"] --> SEN["<b>Sensoren und Telemetrie</b><br/>Vibration, Druck, Strom"]
    SEN --> OBS["<b>Filterung und Validierung</b><br/>Zeitstempel, Rauschunterdrückung"]
    OBS --> EST["<b>Zustandsschätzer</b><br/>Kalman-Filter, digitaler Zwilling"]
    EST --> EXP["<b>Expertensystem</b><br/>Regeln, Normen, Diagnostik"]
    EXP --> DEC["<b>Autorisierte Entscheidung</b>"]
    DEC --> ACT["<b>Aktor / Stellglied</b><br/>Frequenzumrichter, Ventil"]
    ACT --> ENV
    ENV -.->|"physische Rückmeldung"| SEN

    HUM["<b>Operator</b><br/>Autoritätskontrolle"] --> DEC
    EXP --> HUM

    classDef env fill:#fff3e0,stroke:#e65100,stroke-width:2px,color:#bf360c;
    classDef sens fill:#e1f5fe,stroke:#0288d1,stroke-width:2px,color:#01579b;
    classDef est fill:#ede7f6,stroke:#512da8,stroke-width:2px,color:#311b92;
    classDef core fill:#e8f5e9,stroke:#2e7d32,stroke-width:2px,color:#1b5e20;
    classDef hum fill:#fce4ec,stroke:#c2185b,stroke-width:2px,color:#880e4f;

    class ENV,ACT env;
    class SEN,OBS sens;
    class EST est;
    class EXP,DEC core;
    class HUM hum;
```

Das Diagramm unterstreicht, dass die Regeln des Expertensystems keine Rohsignale verarbeiten, sondern auf einer fundierten Zustandsschätzung aufsetzen. Der folgende Abschnitt belegt anhand konkreter Berechnungen, warum diese Differenzierung von prinzipieller Bedeutung ist.

## 2. Parallele Software-Pipeline: Aufgaben-Warteschlange

Rückkopplungsschleifen sind nicht auf physikalische Sensoren beschränkt. Betrachten wir einen Softwaredienst, der Aufgaben entgegennimmt, Worker-Prozesse, die diese abarbeiten, sowie einen Autoscaling-Controller, der die Worker-Anzahl anhand der Telemetrie dynamisch skaliert. Metriklatenzen, Wiederholungsversuche (*retries*) und Kapazitätsgrenzen nachgelagerter Dienste übernehmen hier die Rolle von Latenzen, Rauschen und Stellgliedbeschränkungen.

| Regelkreiselement | Software-Äquivalent | Prüfung durch das Expertensystem |
|---|---|---|
| Regelstrecke (Objekt) | Warteschlange und Worker | Lastanstieg von Durchsatzverlust differenzieren |
| Beobachtung | Alter der ältesten Aufgabe, Eingangsrate, Fertigstellungen und Retries | Zeitstempel und Vollständigkeit der Metriken prüfen |
| Zustandsschätzung | Hypothese über Worker-Mangel oder Überlastung abhängiger Dienste | Eine einzelne Warteschlangenlänge nicht als fertige Diagnose werten |
| Zulässige Aktion | Skalierung der Worker, Annahmebegrenzung oder Benachrichtigung des Operators | Vereinbarte Grenzwerte und Autorisierungen durchsetzen |
| Rückmeldung | Dynamik von Aufgabenalter, Fehlerraten und Retries | Verifizieren, ob die Aktion den Systemzustand tatsächlich verbessert hat |

Arbeitet der Skalierungs-Controller mit veralteten Metriken, kann er mehrfach hintereinander Worker hinzuschalten, bevor der Effekt der ersten Skalierung überhaupt sichtbar wird. Liegt der Engpass in der Datenbank, führen zusätzliche Worker und Retries lediglich zu erhöhter Ressourcenkonkurrenz (*lock contention*). Das Expertensystem kann diese Kausalhypothesen explizieren und Handlungsregeln validieren, ersetzt jedoch nicht den primären Regelungsalgorithmus.

Vor jeder Automatisierung müssen Beobachtungs- und Ausführungslatenzen, Fluktuationen der Worker-Anzahl, Aufgabenalter und Wiederholungsraten systematisch erfasst werden. Zeitfenster für das Eintreten von Effekten, Obergrenzen für Steueraktionen und das Verhalten bei fehlender Telemetrie sind vertraglich zu fixieren. Dieses didaktische Beispiel etabliert keine universellen Schwellenwerte: Für jeden Produktionsdienst sind empirische Lasttests unerlässlich. Das Problem der Zustandsschätzung wird im Folgenden anhand eines physikalischen Signals vertieft.

## 3. Zustandsschätzung: Filterung statt Rohsignal

Betrachten wir eine thermische Schutzregel für Leistungstransistoren in einer Stromverteilungseinheit (*Power Distribution Unit*, PDU): Die Leistung soll halbiert werden, wenn die Temperatur $85^\circ\mathrm{C}$ übersteigt und schneller als $2{,}5^\circ\mathrm{C/s}$ ansteigt:

```math
T > 85^\circ\mathrm{C} \;\land\; \frac{dT}{dt} > 2{,}5^\circ\mathrm{C/s} \;\Rightarrow\; \mathrm{DeratePower}(50\,\%)
```

- Die Schutzregel nutzt $T$, die aktuelle Transistortemperatur in Grad Celsius, mit dem Schwellenwert $85^\circ\mathrm{C}$;
- $\frac{dT}{dt}$ repräsentiert die Temperaturänderungsrate in Grad Celsius pro Sekunde, mit $2{,}5^\circ\mathrm{C/s}$ als Grenzwert der Erwärmungsgeschwindigkeit;
- $>$ bezeichnet „größer als“, $\land$ erfordert das gleichzeitige Erfülltsein beider Bedingungen;
- $\Rightarrow$ signalisiert, dass bei Erfüllung der Prämissen die Aktion $\mathrm{DeratePower}(50\,\%)$ ausgelöst wird (Leistungsreduktion auf 50 %).

Die Regel reduziert die Leistung somit ausschließlich dann, wenn die Temperatur bereits $85^\circ\mathrm{C}$ überschritten hat und weiterhin mit einer Rate oberhalb der Warnschwelle ansteigt. Die Zeitableitung muss zwingend unter Berücksichtigung des Sensorrauschens und des Abtastintervalls geschätzt werden, da hochfrequentes Rauschen andernfalls massive Fehlauslösungen provoziert.

Die Regel erscheint elementar, doch die Ableitung eines verrauschten Signals birgt erhebliche Risiken. Wird der Sensor alle $50\,\text{ms}$ abgetastet und weist das Messrauschen eine Standardabweichung von $0{,}5^\circ\mathrm{C}$ auf, so besitzt die Differenz zweier aufeinanderfolgender Messwerte, dividiert durch $0{,}05\,\text{s}$, ein Rauschen in der Größenordnung von $0{,}5 \cdot \sqrt{2} / 0{,}05 \approx 14^\circ\mathrm{C/s}$. Dieser Rauschpegel übersteigt den Regelschwellenwert um ein Vielfaches. Das folgende Go-Programm simuliert $60\,\text{Sekunden}$ Betrieb unter zwei Szenarien: Der Transistor ist heiß ($86^\circ\mathrm{C}$), erwärmt sich jedoch nicht weiter; im zweiten Szenario erwärmt er sich ausgehend von $80^\circ\mathrm{C}$ mit einer Rate von $3^\circ\mathrm{C/s}$. Die Ableitung wird auf zwei Arten ermittelt: naiv aus zwei aufeinanderfolgenden Messpunkten sowie über ein Alpha-Beta-Filter – den einfachsten Zustandsschätzer mit konstantem Geschwindigkeitsmodell, der Zustands- und Ratenschätzungen bei jedem Zeitschritt iterativ korrigiert.

<details>
<summary>Beispiel in Go: Naive Ableitung und Alpha-Beta-Filter für Temperaturüberwachungsregeln</summary>

```go
package main

import (
	"errors"
	"fmt"
	"math"
	"math/rand"
)

const dt = 0.05 // Abtastperiode des Sensors, s

type Estimate struct {
	Temperature   float64
	Rate          float64
	RateAvailable bool
}

type AlphaBeta struct {
	Alpha, Beta, MaxGap float64
	lastTime            float64
	initialized         bool
	estimate            Estimate
}

func finite(value float64) bool {
	return !math.IsNaN(value) && !math.IsInf(value, 0)
}

func (filter *AlphaBeta) Update(at, measurement float64) (Estimate, error) {
	if !finite(at) || !finite(measurement) || !finite(filter.Alpha) || !finite(filter.Beta) ||
		!finite(filter.MaxGap) || filter.Alpha <= 0 || filter.Alpha > 1 ||
		filter.Beta <= 0 || filter.Beta > 1 || filter.MaxGap <= 0 {
		return Estimate{}, errors.New("invalid measurement or filter configuration")
	}
	if !filter.initialized {
		filter.estimate = Estimate{Temperature: measurement}
		filter.lastTime, filter.initialized = at, true
		return filter.estimate, nil
	}
	interval := at - filter.lastTime
	if interval <= 0 || interval > filter.MaxGap {
		return Estimate{}, errors.New("duplicate, reordered, or stale measurement")
	}
	predicted := filter.estimate.Temperature + filter.estimate.Rate*interval
	residual := measurement - predicted
	next := Estimate{
		Temperature:   predicted + filter.Alpha*residual,
		Rate:          filter.estimate.Rate + filter.Beta/interval*residual,
		RateAvailable: true,
	}
	if !finite(next.Temperature) || !finite(next.Rate) {
		return Estimate{}, errors.New("non-finite estimate")
	}
	filter.estimate, filter.lastTime = next, at
	return next, nil
}

// rule ist die Regel zur Leistungsreduktion: Temperatur über 85 °C und Anstieg schneller als 2,5 °C/s.
func rule(temp, rate float64) bool { return temp > 85 && rate > 2.5 }

// run simuliert 60 s Betrieb und liefert die Anzahl der Regelauslösungen sowie den Zeitpunkt des ersten Ansprechens.
// Wenn alpha == 0, wird die Ableitung „naiv“ aus zwei aufeinanderfolgenden Messungen berechnet.
func run(trueTemp func(t float64) float64, alpha, beta float64) (fires int, first float64) {
	r := rand.New(rand.NewSource(7))
	first = -1
	var prev float64
	filter := AlphaBeta{Alpha: alpha, Beta: beta, MaxGap: 4 * dt}
	for k := 0; k <= int(60/dt); k++ {
		t := float64(k) * dt
		z := trueTemp(t) + 0.5*r.NormFloat64() // Messung mit Rauschen σ = 0,5 °C
		var temp, rate float64
		if k == 0 {
			prev = z
		}
		if alpha == 0 {
			temp, rate = z, (z-prev)/dt
		} else {
			estimate, err := filter.Update(t, z)
			if err != nil {
				panic(err)
			}
			temp, rate = estimate.Temperature, estimate.Rate
		}
		prev = z
		if k > 0 && rule(temp, rate) {
			fires++
			if first < 0 {
				first = t
			}
		}
	}
	return fires, first
}

func main() {
	stable := func(t float64) float64 { return 86 } // heiß, erwärmt sich jedoch nicht weiter
	ramp := func(t float64) float64 {               // Erwärmung um 3 °C/s über 10 s ab 80 °C
		if t < 10 {
			return 80 + 3*t
		}
		return 110
	}
	for _, f := range []struct {
		name        string
		alpha, beta float64
	}{{"naive Ableitung", 0, 0}, {"Alpha-Beta-Filter", 0.1, 0.005}} {
		falseFires, _ := run(stable, f.alpha, f.beta)
		_, first := run(ramp, f.alpha, f.beta)
		fmt.Printf("%-18s | Fehlauslösungen in 60 s: %4d | Erste Auslösung bei Erwärmung: %.2f s\n",
			f.name, falseFires, first)
	}
}
```

Synthetische Tests verifizieren den Zeitbereich und den Filterzustand. Ausführung: `go test -v main.go main_test.go`.

```go
package main

import (
	"math"
	"testing"
)

func TestMeasurementValidity(testCase *testing.T) {
	filter := AlphaBeta{Alpha: 0.1, Beta: 0.005, MaxGap: 0.2}
	initial, err := filter.Update(7, 25)
	if err != nil || initial.RateAvailable {
		testCase.Fatal("first sample incorrectly supplied a measured rate")
	}
	updated, err := filter.Update(7.1, 26)
	if err != nil || !updated.RateAvailable || math.Abs(updated.Rate-0.05) > 1e-10 {
		testCase.Fatalf("wrong actual-time update: %+v, %v", updated, err)
	}
	before := filter
	for _, input := range []struct{ at, value float64 }{
		{7.1, 26}, {7, 26}, {9, 26}, {7.15, math.NaN()}, {7.15, math.Inf(1)},
	} {
		if _, err := filter.Update(input.at, input.value); err == nil {
			testCase.Fatal("invalid event accepted")
		}
		if filter != before {
			testCase.Fatal("invalid event changed filter state")
		}
	}
	if _, err := (&AlphaBeta{Alpha: 0.1, Beta: 0.005}).Update(0, 25); err == nil {
		testCase.Fatal("missing gap policy accepted")
	}
}
```

</details>

Das Programm liefert folgende Ausgabe:

<details>
<summary>Daten und Ausführungsergebnis des Beispiels</summary>

```text
naive Ableitung    | Fehlauslösungen in 60 s:  520 | Erste Auslösung bei Erwärmung: 1.55 s
Alpha-Beta-Filter  | Fehlauslösungen in 60 s:    0 | Erste Auslösung bei Erwärmung: 1.80 s
```

</details>

Bei einem festen Seed des Zufallsgenerators generierte die naive Ableitung im stationären Zustand 520 positive Regelauswertungen, während das Filter keine einzige Fehlauslösung verzeichnete. Dies entspricht nicht der Zahl separater Schutzaktionen und stellt keinen formalen Beweis einer Null-Fehlalarmrate dar: Das Programm modelliert nicht die thermische Gegenkopplung nach Leistungsdrosselung. Während der Erwärmungsphase wird die Temperaturschwelle von $85^\circ\mathrm{C}$ nach ca. $1{,}67\,\text{s}$ geschnitten; das Filter lieferte die erste Auslösung bei $1{,}80\,\text{s}$. Latenzen und Schätzfehler müssen über eine Vielzahl von Rauschrealisierungen, Signalsprüngen und Messausfällen verifiziert werden.

Das Flag `RateAvailable` indiziert lediglich, dass die Rate nach zwei gültigen Ereignissen berechnet werden konnte, garantiert jedoch nicht deren hinreichende Genauigkeit für Schutzfunktionen. Das Filter verarbeitet reale Zeitabstände, verwirft ungültige Zahlenwerte, Duplikate, zeitlich vertauschte Ereignisse sowie unzulässig große Messlücken ohne Statusänderung. Eine verworfene Messung versetzt die Beobachtung in einen degradierten Status; sie wird keinesfalls stillschweigend durch Null oder das Fehlen einer Bedrohung substituiert. Grenzintervalle und Filterkoeffizienten stellen Parameter dieses didaktischen Experiments dar und bilden kein validiertes Profil eines realen Leistungstransistors. Das Beispiel illustriert den Nutzen expliziter Rauschmodelle, beweist jedoch weder das Good-Regulator-Theorem noch die absolute Sicherheit eines geschlossenen Regelkreises.

### 3.1. Betriebshüllkurve (Operating Envelope) und thermisches Derating über Festkommaarithmetik (Q16.16)

Das Filtern des Rauschens eines einzelnen Temperatursignals eliminiert zwar Fehlalarme, doch in sicherheitskritischen leistungselektronischen Systemen (Traktionswechselrichter in Elektrofahrzeugen, Batterie-Management-Systeme BMS, Steuergeräte auf Basis von Infineon AURIX gemäß ISO 26262 ASIL D) lassen sich physikalische Sicherheitsgrenzen niemals durch isolierte, eindimensionale Schwellenwerte abbilden. Der sichere Arbeitsbereich eines Systems bildet vielmehr eine mehrdimensionale gekoppelte Hyperfläche (*Operating Envelope*), in der die maximal zulässige Kernspannung $V_{\mathrm{core}}$ streng von der aktuellen Sperrschichttemperatur $T_j$, der Taktfrequenz $F_{\mathrm{cpu}}$ und den Injektionsströmen der I/O-Pins $I_{\mathrm{inj}}$ abhängt.

Versuche, solche nichtlinearen Sicherheitshüllkurven rein softwareseitig auf dem Hauptprozessor zu berechnen, scheitern an zwei gravierenden Barrieren:
1. **Inakzeptabler, nicht-deterministischer Jitter:** Gleitkommaoperationen (FP32/FP64) verursachen Rechenzeit-Varianzen und Latenzen im zweistelligen Mikrosekundenbereich, was ein rechtzeitiges Eingreifen im geforderten Sub-Mikrosekunden-Schutzzyklus ($< 1\,\mu\text{s}$) verunmöglicht.
2. **Katastrophaler Hüllkurven-Inversionsfehler (Envelope Inversion Fault):** Bei der hardwarenahen Umstellung auf Ganzzahl-Festkommaarithmetik im Format Q16.16 führt die Multiplikation eines großen Temperaturhubs ($\Delta T = 125^\circ\text{C} - 25^\circ\text{C} = 100^\circ\text{C}$) mit einer steilen Derating-Steigung zu einem Überlauf der oberen Bits des 64-Bit-Produkts. Wird das Produkt ohne Sättigungslogik mechanisch trunkiert (`prod[47:16]`), invertiert das Vorzeichenbit: Die berechnete Leistungsreduktion schlägt schlagartig in einen negativen Wert um. Dadurch wird der zulässige Leistungsbereich im thermischen Notfall fatalerweise künstlich **erweitert** statt protektiv verengt – was unmittelbar zur thermischen Zerstörung der Leistungshalbleiter führt.

Um diese Verwundbarkeit auf Hardwareebene kompromisslos auszuschließen, wird ein pipelinierter Coprozessor zur Hüllkurvenbegrenzung mit nativer Sättigungsarithmetik im Format Q16.16 realisiert:

```math
\Delta T = \max(0, T_{\mathrm{sensor}} - T_{\mathrm{ref}}).
```

```math
P_{\mathrm{derated}} = P_{\mathrm{max\_nominal}} - \mathrm{SaturateQ16}(\Delta T \cdot K_{\mathrm{derate}}).
```

```math
\mathrm{SaturateQ16}(X) = \begin{cases} 
\text{0x7FFF\_FFFF}, & \text{wenn } X[63:47] \ne \{17\{X[47]\}\} \land X[63] = 0, \\
\text{0x8000\_0000}, & \text{wenn } X[63:47] \ne \{17\{X[47]\}\} \land X[63] = 1, \\
X[47:16], & \text{sonst}.
\end{cases}
```

Parameter der Formeln zur Hardware-Degression:
- $T_{\mathrm{sensor}}$ — aktuelle Sensortemperatur im Q16.16-Format, Wertebereich $[-40{,}0, +175{,}0]^\circ\text{C}$;
- $T_{\mathrm{ref}}$ — Referenztemperatur für den Beginn des thermischen Deratings (z. B. $+85{,}0^\circ\text{C}$ bzw. $+45{,}0^\circ\text{C}$);
- $\Delta T$ — positive Überschreitung der Temperaturschwelle;
- $K_{\mathrm{derate}}$ — Derating-Koeffizient (W/°C im Format Q16.16);
- $P_{\mathrm{max\_nominal}}$ — nominale Nennleistung im ungedrosselten Kaltbereich (W im Format Q16.16);
- $P_{\mathrm{derated}}$ — dynamisch berechnete momentane Obergrenze der zulässigen Leistung;
- $\mathrm{SaturateQ16}$ — Hardware-Sättigungsoperator, der die oberen 17 Bits des 64-Bit-Produkts $X$ auswertet und Vorzeichenüberläufe deterministisch abfängt.

```mermaid
flowchart TB
    accTitle: Vierstufige Pipeline des Hardware-Hyper-Limiters mit Sättigung
    accDescr: Einspeisung des Arbeitspunktvektors in Takt 1, Multiplikation mit Überlauferkennung in Takt 2, Q16.16-Sättigung in Takt 3 und direkter Hardware-Interlock in Takt 4.

    T_IN["T_sensor, T_ref (Q16.16)"] --> S1["<b>Takt 1</b>:<br/>ΔT = max(0, T - T_ref)"]
    S1 --> S2["<b>Takt 2</b>:<br/>64-Bit-Produkt:<br/>X = ΔT × K_derate"]
    S2 --> S3["<b>Takt 3</b>:<br/>Erkennung X[63:47]<br/>SaturateQ16(X)"]
    S3 --> S4["<b>Takt 4</b>:<br/>P_derated = P_nom - Sat(X)<br/>Vergleich mit P_actual"]
    
    S4 --> CHK{"P_actual ><br/>P_derated?"}
    CHK -- "Nein" --> OK["Schleuse: PASS<br/>flag_zf = 1"]
    CHK -- "Ja" --> INT["Hardware-Interlock:<br/>flag_uf = 1<br/>HARD_INTERLOCK (Fail-Closed)"]

    classDef proc fill:#e1f5fe,stroke:#0288d1,stroke-width:2px,color:#01579b;
    classDef dec fill:#fff9c4,stroke:#fbc02d,stroke-width:2px,color:#f57f17;
    classDef safe fill:#e8f5e9,stroke:#2e7d32,stroke-width:2px,color:#1b5e20;
    classDef err fill:#ffebee,stroke:#d32f2f,stroke-width:2px,color:#b71c1c;

    class S1,S2,S3,S4 proc;
    class CHK dec;
    class OK safe;
    class INT err;
```

#### 3.1.1. Geschlossener Kreis der ingenieurtechnischen Berechnung des Betriebshüllkurven-Deratings

1. **Laufzeitsteuerung und Rechenfluss (Control Flow & Runtime Decisions):**
   In Takt 4 der Pipeline gleicht ein digitaler Komparator den aktuellen Lastvektor mit der dynamischen Obergrenze $P_{\mathrm{derated}}$ ab. Gilt $P_{\mathrm{actual}} \le P_{\mathrm{derated}}$, wird das Statusbit $\mathrm{flag\_zf}$ auf 1 gesetzt und das PWM-Steuersignal verzögerungsfrei an die Leistungsendstufe durchgeschaltet. Wird die Leistungsgrenze überschritten oder kollabiert das Gültigkeitsfenster ($P_{\mathrm{derated}} < P_{\mathrm{min\_safe}}$), setzt die Logik das Fehlerflag $\mathrm{flag\_uf} = 1$ und erzeugt unmittelbar das Hardwaresignal `HARD_INTERLOCK`. Dies schaltet die Treiber der Leistungstransistoren innerhalb eines einzigen Taktes ($5\,\text{ns}$) ohne jede Betriebssystemintervention in den sicheren Zustand (Fail-Closed).
2. **Hardware-Dimensionierung und infrastrukturelle Schranken (Hardware Dimensioning):**
   - Die Pipeline wird über vier synchrone Registerstufen in FPGA/ASIC implementiert und belegt maximal 380 LUTs sowie einen DSP48E2-Multiplizierer (z. B. auf Xilinx UltraScale+);
   - Bei einem Systemtakt von $200\,\text{MHz}$ beträgt die Latenz exakt $4 \times 5\,\text{ns} = 20\,\text{ns}$ bei null Jitter ($\mathrm{Jitter} = 0$), was um drei Größenordnungen schneller ist als Software-Timer typischer Echtzeitbetriebssysteme ($10\text{–}50\,\mu\text{s}$).
3. **Praktisches numerisches Rechenbeispiel (Worked Numerical Example):**
   Betrachten wir einen industriellen Leistungswechselrichter mit einer Nennleistung $P_{\mathrm{max\_nominal}} = 2000\,\text{W}$ (im Format Q16.16: `0x07D0_0000`), einer Referenztemperatur für den Derating-Beginn $T_{\mathrm{ref}} = 45^\circ\text{C}$ und einem Derating-Koeffizienten $K_{\mathrm{derate}} = 40\,\text{W/}^\circ\text{C}$ (im Format Q16.16: `0x0028_0000`).
   Steigt die Umgebungs- und Sensortemperatur im Schaltschrank auf $T_{\mathrm{amb}} = T_{\mathrm{sensor}} = 65^\circ\text{C}$:
   - $\Delta T = 65 - 45 = 20^\circ\text{C}$ (`0x0014_0000`);
   - Produkt $\Delta T \cdot K_{\mathrm{derate}} = 20 \times 40 = 800\,\text{W}$;
   - Dynamische Leistungsgrenze $P_{\mathrm{derated}} = 2000 - 800 = 1200\,\text{W}$ (`0x04B0_0000`).

   Fordert die übergeordnete Steuerung eine Last von $P_{\mathrm{req}} = 1800\,\text{W}$ an, detektiert die Pipeline nach $20\,\text{ns}$ die Schrankenverletzung ($1800\,\text{W} > 1200\,\text{W}$). Sie setzt unverzüglich $\mathrm{flag\_uf} = 1$, blockiert die Anforderung per `HARD_INTERLOCK Refusal` (Fail-Closed) und zwingt die Endstufe in eine sichere Leistungsbegrenzung.
   Auf dieselbe Weise schützt die Schaltung automobile Traktionswechselrichter ($P_{\mathrm{max\_nominal}} = 150\,\text{kW}$, $T_{\mathrm{ref}} = 85^\circ\text{C}$, $K_{\mathrm{derate}} = 2{,}5\,\text{kW/}^\circ\text{C}$): Steigt $T_{\mathrm{sensor}}$ auf $135^\circ\text{C}$, verbleibt $P_{\mathrm{derated}} = 25\,\text{kW}$; eine Motorlast von $80\,\text{kW}$ löst binnen $20\,\text{ns}$ den Interlock aus und schützt die IGBT-Module vor dem Durchbruch.

---

## 4. Kybernetik zweiter Ordnung, Viable System Model und Active Inference

Die klassische Kybernetik betrachtete das zu regelnde Objekt als strikt extern gegenüber dem Regler. In komplexen soziotechnischen Systemen greift diese Annahme zu kurz. Drei theoretische Konzepte der zweiten Hälfte des 20. und des frühen 21. Jahrhunderts liefern das Rüstzeug, um diesen Wechselwirkungen architektonisch zu begegnen.

**Kybernetik zweiter Ordnung.** Heinz von Foerster postulierte, dass der Beobachter als integraler Bestandteil des von ihm beobachteten Systems modelliert werden muss [[6]](#src-6). Für ein Expertensystem ist dies keineswegs eine philosophische Abstraktion: Sobald das System eine Diagnose oder Handlungsempfehlung ausgibt, modifizieren menschliche Operatoren ihr Verhalten. Dadurch verändern sich exakt jene Datenströme, auf deren Basis das Expertensystem trainiert und verifiziert wird. Das Entscheidungsprotokoll des Expertensystems dient somit nicht allein Revisionszwecken, sondern bildet einen integralen Bestandteil der Prozessdaten; Evaluierungsmetriken müssen den Rückkopplungseffekt vergangener Empfehlungen explizit einbeziehen.

**Viable System Model (VSM).** Stafford Beer formulierte das Modell lebensfähiger Systeme (*Viable System Model*), wonach jede überlebensfähige Organisation zwingend fünf funktionale Subsysteme aufweisen muss [[7]](#src-7). Das folgende Diagramm stellt diese Subsysteme dar.

```mermaid
flowchart TB
    accTitle: Viable System Model nach Stafford Beer
    accDescr: Fünf Subsysteme von primären Operationen bis zur Politik: Operationen, Koordination gegen Oszillationen, interne Kontrolle, Anpassung an die Zukunft sowie Identität und Politik.

    subgraph VSM["<b>Viable System Model</b>"]
        direction TB
        S5["<b>System 5: Politik und Identität</b><br/>Sicherheits-Grenzregeln, Mission"]
        S4["<b>System 4: Anpassung und Zukunft</b><br/>Szenariomodellierung, Bedrohungsprognose"]
        S3["<b>System 3: Interne Kontrolle</b><br/>Ressourcenaudit, Invariantenprüfung"]
        S2["<b>System 2: Koordination</b><br/>Dämpfung von Schwingungen zwischen Einheiten"]
        S1["<b>System 1: Primäre Operationen</b><br/>Technologische Prozesse, Sensoren, Aktoren"]
        S5 <--> S4
        S4 <--> S3
        S3 <--> S2
        S2 <--> S1
    end

    classDef s5 fill:#fce4ec,stroke:#c2185b,stroke-width:2px,color:#880e4f;
    classDef s4 fill:#ede7f6,stroke:#512da8,stroke-width:2px,color:#311b92;
    classDef s3 fill:#e8f5e9,stroke:#2e7d32,stroke-width:2px,color:#1b5e20;
    classDef s2 fill:#fff3e0,stroke:#e65100,stroke-width:2px,color:#bf360c;
    classDef s1 fill:#e1f5fe,stroke:#0288d1,stroke-width:2px,color:#01579b;

    class S5 s5;
    class S4 s4;
    class S3 s3;
    class S2 s2;
    class S1 s1;
```

In diesem Ordnungsrahmen positioniert sich das Expertensystem typischerweise zwischen System 3 (Einhaltung technologischer Vorschriften im Hier und Jetzt) und System 4 (Modellierung von Verschleißprozessen oder Störfallszenarien), eingebettet in die normativen Leitplanken von System 5. Der architektonische Mehrwert des VSM liegt in einer diagnostischen Kontrollfrage: Fehlt in einer Anlage ein explizites System 2, existiert keine Instanz, um destruktive Schwingungen und Ressourcenkonflikte zwischen autonomen Edge-Knoten (etwa die gleichzeitige Lastabsenkung zweier Einheiten auf derselben Versorgungsleitung) aktiv zu dämpfen.

**Active Inference.** Karl Friston formulierte das Free-Energy-Prinzip als einheitliche Theorie adaptiver Systeme: Ein adaptives System minimiert die variationelle freie Energie $F$, die eine obere Schranke der „Überraschung“ (*Surprise*) bezüglich der Beobachtungen $\tilde{y}$ darstellt [[8]](#src-8):

```math
F = \mathbb{E}_{q(\vartheta)}\left[\ln q(\vartheta) - \ln p(\tilde{y}, \vartheta)\right] = \mathrm{KL}\left[q(\vartheta) \,\|\, p(\vartheta \mid \tilde{y})\right] - \ln p(\tilde{y})
```

- In der Gleichung der freien Energie ist $F$ ein Skalar und $\vartheta$ bezeichnet den verborgenen Zustand des Objekts;
- $q(\vartheta)$ ist die approximative Verteilung interner Zustandsüberzeugungen des Modells und $p(\tilde{y},\vartheta)$ die Verbundverteilung von Beobachtung $\tilde{y}$ und Zustand $\vartheta$ gemäß dem generativen Modell $p$;
- $\mathbb{E}_{q(\vartheta)}$ bezeichnet den Erwartungswert bezüglich der Verteilung $q$, und $\ln$ ist der natürliche Logarithmus;
- $p(\vartheta\mid\tilde{y})$ ist die A-posteriori-Verteilung des verborgenen Zustands nach der Beobachtung, und $\mathrm{KL}[q\|p]$ steht für die Kullback-Leibler-Divergenz zwischen den Verteilungen;
- $p(\tilde{y})$ bezeichnet die Randwahrscheinlichkeit der Beobachtung; die Gleichung zeigt die Äquivalenz zwischen der Erwartungswert-Darstellung und der KL-Form mit dem additiven Term $-\ln p(\tilde{y})$;
- Bei Verwendung des natürlichen Logarithmus wird die Informationsgröße in Nat gemessen.

Die Minimierung von $F$ kann über zwei distinkte Pfade erfolgen: Entweder passt das System seine internen Überzeugungen $q$ an, um die Beobachtungen besser zu erklären (Wahrnehmung), oder es verändert die physische Umwelt durch gezielte Aktionen so, dass künftige Beobachtungen den Erwartungen entsprechen (Active Inference). Für ein Expertensystem fungiert diese Theorie als wertvolle konzeptionelle Analogie: Der Zustandsschätzer übernimmt die Funktion der Wahrnehmung, während die in [Kapitel 21](ch21-from-recommendation-to-action.md) beschriebenen Aktoren die Aktionen vollziehen. Das Free-Energy-Prinzip entstammt den Neurowissenschaften; ein cyber-physisches Expertensystem muss es nicht buchstabengetreu nachbilden, um Regelkreise deterministisch zu schließen.

> [!NOTE]
> **Zwei Facetten ingenieurtechnischer Active Inference nach dem Autormodell**
> 1. Wahrnehmung (Perception): Das System aktualisiert seine interne A-posteriori-Verteilung $q(\vartheta)$ (Zustandsvektor des Kalman-Filters oder Bayes-Netzes) zur Minimierung der Divergenz $\mathrm{KL}[q \| p]$ als passive Modellanpassung an Sensordaten.
> 2. Aktion (Action): Bei Detektion kritischer Abweichungen generiert die Steuerung deterministische Stellsignale an Servomotoren oder Ventile und zwingt die Regelstrecke in sichere Sollwerte zur Verringerung der Überraschung $-\ln p(\tilde{y})$.
> *Runtime-Steuerung:* Die Minimierung variationeller freier Energie übersetzt sich direkt in die hardwarenahe Taktung von Zustandsschätzung und kompensatorischer Regelung.

Diese drei theoretischen Ansätze konvergieren in einem organisatorischen Grundsatz: Wissen zirkuliert in technischen Infrastrukturen über zwei Schleifen unterschiedlicher Taktung. Der schnelle operative Zyklus umfasst operative Ad-hoc-Entscheidungen, Sensordatenanalysen, Hypothesen autonomer Agenten und probabilistische Schätzungen; er besitzt jedoch keine normative Zertifizierungsautorität. Der langsame, gesteuerte Zyklus verwaltet den kanonischen Wissensgraphen, validierte Architekturrichtlinien, geprüfte Ontologien und rechtliche Randbedingungen. Wissen darf aus dem schnellen in den langsamen Zyklus ausschließlich über strukturierte Promotionsverfahren mit Invariantenprüfung und kryptografischer Signatur autorisierter Ingenieure überführt werden, wie in [Kapitel 19](ch19-from-question-to-evidence.md) dargelegt.

## 5. Synergetische Stabilität: Bifurkationen, Critical Slowing Down und Vorhersage von Phasenübergängen

Die klassische Regelungstheorie stützt sich überwiegend auf lineare Rückkopplungen: Eine Regelabweichung vom Sollwert erzeugt ein proportionales Korrektursignal. In realen cyber-physischen Komplexen (autonome Fluggeräte, Industrieturbinen, HGÜ-Stromnetze) ist die Anlagendynamik jedoch hochgradig nichtlinear. Unter steigender Belastung oder thermischem Verschleiß kann sich die Regelstrecke einem **Bifurkationspunkt** nähern: einer kritischen Schwelle qualitativer Änderungen des dynamischen Regimes, an der lineare Rückkopplungen ihre Stabilitätsreserve einbüßen und das Gesamtsystem in chaotische Oszillationen oder katastrophales Versagen abgleitet.

Die **Synergetik nach Haken und Prigogine** zeigt, dass komplexe dynamische Systeme unmittelbar vor einem destruktiven Phasenübergang ein fundamentales Phänomen manifestieren: **Critical Slowing Down (CSD)** [[8a]](#src-8a).

Beschreibt man die Dynamik der Zustandsabweichung $x$ in der Nähe eines stationären Gleichgewichtspunkts durch eine linearisierte Differentialgleichung mit stochastischem Rauschen:

```math
\frac{dx}{dt} = -\lambda x + \sigma \eta(t)
```

wobei:
- $x$ die Auslenkung des Zustands aus der Gleichgewichtslage darstellt;
- $\lambda > 0$ die Erholungsrate (*Recovery Rate*, Rückkehrgeschwindigkeit ins Gleichgewicht) bezeichnet;
- $\sigma \eta(t)$ das stochastische Umweltrauschen mit der Rauschintensität $\sigma$ ist.

Nähert sich das System einem Bifurkationspunkt, strebt der dominierende Eigenwert gegen null: $\lambda \to 0$. Dies induziert zwei statistische Indikatoren, die im Sensorstrom in Software präzise berechnet werden können:

1. **Anstieg der Autokorrelation (Memory Effect):**
Für diskrete Abtastungen im Zeitintervall $\Delta t$ konvergiert der Autokorrelationskoeffizient erster Ordnung gegen eins:

```math
\rho_1 = e^{-\lambda \Delta t} \xrightarrow{\lambda \to 0} 1
```

Das System „erinnert“ sich signifikant länger an zufällige Auslenkungen; das Einschwingverhalten wird zäh und stark verzögert.

2. **Divergenz der Varianz (Variance Inflation):**
Die Varianz der Messungen wächst sprunghaft an, da das System nicht mehr in der Lage ist, stochastische Störungen zügig zu dissipieren:

```math
\mathrm{Var}(x) = \frac{\sigma^2}{2\lambda} \xrightarrow{\lambda \to 0} \infty
```

**Laufzeitsteuerung und präventive Arbitrierung:**
- **CSD-Warnschwelle (CSD Alarm):** Werden in einem gleitenden Beobachtungsfenster simultan die Bedingungen $\rho_1 \ge \tau_{\rho} = 0{,}85$ und $\mathrm{Var}(x) / \sigma_0^2 \ge 3{,}0$ erfüllt (die Varianz hat sich gegenüber dem Basisrauschen $\sigma_0^2$ verdreifacht), signalisiert das System den drohenden Verlust der dynamischen Stabilität;
- **Aktion des Expertenkerns:** Die symbolische Inferenzmaschine entzieht dem linearen Regler die Steuerautorität und initiiert das Protokoll zur präventiven Überführung in einen sicheren Zustand (`PREEMPTIVE_FAIL_SAFE`), indem Drehzahl oder elektrische Last vor Eintritt der irreversiblen Bifurkation aktiv gedrosselt werden.

**Numerisches Berechnungsbeispiel zur Bifurkationsdetektion:**
Ein Temperatursensor eines Wechselrichters mit Abtastintervall $\Delta t = 0{,}5\,\text{s}$ registriert eine schleichende Degradation des Kühlsystems. Die Erholungsrate ist auf $\lambda = 0{,}1\,\text{s}^{-1}$ eingebrochen, bei einer Rauschintensität von $\sigma = 0{,}2^\circ\text{C}$ (Basisvarianz $\sigma_0^2 = 0{,}04$):

```math
\rho_1 = e^{-0{,}1 \cdot 0{,}5} = e^{-0{,}05} \approx 0{,}951 \ge 0{,}85
```

```math
\mathrm{Var}(x) = \frac{0{,}2^2}{2 \cdot 0{,}1} = \frac{0{,}04}{0{,}2} = 0{,}20\,(^\circ\text{C})^2 \implies \frac{\mathrm{Var}(x)}{\sigma_0^2} = \frac{0{,}20}{0{,}04} = 5{,}0 \ge 3{,}0
```

Beide Schwellenwerte sind überschritten: Das System wechselt unverzüglich in den präventiven Schutzbetrieb zur Lastabsenkung und bannt den thermischen Durchschlag der Halbleiter.

```mermaid
flowchart LR
    accTitle: Vorhersage von Bifurkationen anhand synergetischer Indikatoren
    accDescr: Von der Beobachtung des Sensorrauschens über die CSD-Berechnung bis zum präventiven Umschalten durch das Expertensystem vor einem Störfall.

    RAW["Sensorstrom<br/>(verrauschte Zeitreihe)"] --> CSD["CSD-Berechner<br/>(Autokorrelation ρ₁ und Varianz σ²)"]
    CSD --> COND{"ρ₁ > 0.85 und<br/>Var steigt?"}
    COND -->|"Nein: stabiler Modus"| REG["Regulärer kybernetischer Regler<br/>(PID / MPC)"]
    COND -->|"Ja: Annäherung an Bifurkation"| EXP["Symbolischer Expertenkern<br/>(Präventive Schutzarbitrierung)"]
    EXP --> ACT["Sicherer Zustand (Fail-Safe)<br/>Drehzahlreduktion, Spannungsabschaltung"]
```

Für Expertensysteme begründet dies eine neue Klasse von Regeln: **präventive Vor-Bifurkations-Schutzregeln**. Anstatt passiv abzuwarten, bis eine Temperatur den Notabschaltwert reißt oder ein Fluggerät unsteuerbar taumelt, detektiert das System synergetische CSD-Frühwarnindikatoren ($\rho_1 \uparrow, \mathrm{Var} \uparrow$) als aktiven Defeater für den Normalbetrieb und schaltet die Anlage noch im Bereich linearer Beherrschbarkeit präventiv in einen geschützten Modus.

## 6. Hierarchie der Zeithorizonte L0–L4

Unterschiedliche Entscheidungen erfordern grundlegend verschiedene Zeitbudgets. Die Ebenen L0–L4 und die genannten Intervalle bilden die didaktische Strukturierung dieses Kapitels und stellen keine starre Industrienorm dar. Die Zulässigkeit einer Softwarekomponente auf einer Ebene bemisst sich an Reaktionsfristen, Worst-Case Execution Time (WCET), Netzwerkanforderungen und den spezifischen Sicherheitszielen der Anlage; weder eine bestimmte Programmiersprache noch Hardwarebeschleuniger beweisen für sich genommen die Einhaltung dieser Grenzen.

| Ebene | Zeithorizont | Funktion und Zweck | Auf dieser Ebene unzulässig |
|:---:|---|---|---|
| **L0** | Kontinuierlich | Physisches Objekt: Mechanik, Werkstoffe, Thermodynamik, Energiewandlung | Softwareannahmen ohne sensorische Bestätigung |
| **L1** | 1 µs–10 ms | Harte lokale Regelung: Servoregler, SPS, FPGA, Hardware-Verriegelungen | Nicht-deterministische neuronale Netze, LLMs, Netzwerkaufrufe |
| **L2** | 10 ms–1 s | Wahrnehmung: Digitale Signalfilterung, Sensorfusion, Klassifikation auf NPU | Fakten ohne Zeitstempel und Unsicherheitsabschätzung |
| **L3** | 1 s–1 min | Edge-Aufsicht: Lokale Regeln, Risikobewertung, Notbetriebskoordination | Befugnis zur eigenmächtigen Änderung globaler Richtlinien |
| **L4** | 1 min–Tage | Zentrales Expertensystem: Globaler Wissensgraph, Flottenanalyse, Modellpflege | Annahme einer permanenten, latenzfreien Edge-Verbindung |

Das nachfolgende Diagramm veranschaulicht das Zusammenspiel dieser fünf Ebenen.

```mermaid
flowchart TB
    accTitle: Verteilte Architektur der Ebenen L0–L4
    accDescr: Physische Sensoren, Aktoren und schnelle Controller auf den Ebenen L0–L1, Signalverarbeitung und Zustandsschätzung auf L2, Edge-Überwachung auf L3 und zentrales Expertensystem auf L4.

    subgraph L4_Backend["L4: Zentrales Expertensystem und Fachexperten"]
        direction LR
        HUM["Prozessingenieur"]
        KG["Globaler Wissens- und Vorfalls-Graph"]
        DT["Digitaler Zwilling des Unternehmens"]
        TRAIN["Modellverifikation und Nachtraining"]
    end

    subgraph L3_Edge["L3: Edge-Überwachungsknoten"]
        direction LR
        E_RULE["Lokale Regel-Engine"]
        E_DIAG["Fehlerdiagnose"]
        E_SAFE["Sicherheitsüberwachung und Fallback-Modus"]
    end

    subgraph L2_Perception["L2: Signalverarbeitung und Zustandsschätzung"]
        direction LR
        DSP["Digitale Filterung"]
        NPU_INF["NPU-Klassifikation"]
        FUSION["Datenfusion, Kalman-Filter"]
    end

    subgraph L1_L0_Plant["L1–L0: Schnelle Steuerung und physisches Objekt"]
        direction LR
        PLC["Speicherprogrammierbare Steuerung SPS (L1)"]
        INTER["Hardware-Sicherheitsverriegelung (L1)"]
        ACT_PHYS["Physischer Aktor (L0)"]
        SENS_PHYS["Physische Sensoren (L0)"]
    end

    SENS_PHYS --> DSP --> NPU_INF --> FUSION
    FUSION --> E_DIAG --> E_RULE --> E_SAFE
    E_SAFE --> PLC --> ACT_PHYS
    INTER -.->|"Hardware-Notabschaltung"| ACT_PHYS

    E_RULE <-->|"Fakten und Richtlinien"| KG
    FUSION -->|"Aggregierte Telemetrie"| DT
    HUM -->|"Direktive Befehle"| E_SAFE
    TRAIN -->|"Signierte Wissens-Bundles"| E_RULE

    classDef l4 fill:#fce4ec,stroke:#c2185b,stroke-width:2px,color:#880e4f;
    classDef l3 fill:#e8f5e9,stroke:#2e7d32,stroke-width:2px,color:#1b5e20;
    classDef l2 fill:#ede7f6,stroke:#512da8,stroke-width:2px,color:#311b92;
    classDef l1 fill:#e1f5fe,stroke:#0288d1,stroke-width:2px,color:#01579b;

    class HUM,KG,DT,TRAIN l4;
    class E_RULE,E_DIAG,E_SAFE l3;
    class DSP,NPU_INF,FUSION l2;
    class PLC,INTER,ACT_PHYS,SENS_PHYS l1;
```

Die gestrichelte Linie von der Hardware-Verriegelung direkt zum Stellglied ist von zentraler sicherheitstechnischer Relevanz: Die Notabschaltung auf Ebene L1 durchläuft keinerlei Software höherer Ebenen. Das zentrale Expertensystem wirkt auf die Peripherie ausschließlich über zwei Kanäle ein: die Synchronisation aggregierter Fakten und Richtlinien sowie die Bereitstellung kryptografisch signierter Wissens-Bundles. Beide Kanäle operieren asynchron und verifizierbar; keiner von ihnen ist Bestandteil der zeitkritischen Schutzschleife.

## 7. Verträge für Sensorereignisse

Ob die vollständige Übertragung sämtlicher Rohabtastungen sinnvoll ist, hängt von Abtastfrequenzen, Kanalzahlen, Busbandbreiten und der konkreten Diagnoseaufgabe ab. Für die kontinuierliche Inferenz überträgt der Edge-Knoten typischerweise aggregierte Fenster-Features inklusive Qualitätsmetadaten, während autorisierte Rohdatenfragmente für tiefgehende Audits und Nachprüfungen separat vorgehalten werden. Eine Aggregation darf jedoch niemals transiente Ereignisse verschleifen, die für Schutzregeln entscheidend sind. Das folgende Listing demonstriert das strukturierte Schema eines solchen Sensorereignisses:

<details>
<summary>Strukturierte JSON-Repräsentation eines Sensorereignisses</summary>

```json
{
  "event_id": "evt:vibr:sensor-42:2026-09-27T10:14:02.105Z",
  "source_urn": "urn:plant:sector-b:motor-07:vibr-probe-1",
  "boot_id": "synthetic-boot-17",
  "sequence": 1842,
  "event_time": "2026-09-27T10:14:02.105Z",
  "clock_uncertainty_ms": 5,
  "measurement_window_ms": 100,
  "sampling_rate_hz": 20000,
  "features": {
    "peak_acceleration_g": 4.82,
    "rms_velocity_mms": 11.4,
    "dominant_frequency_hz": 142.5,
    "kurtosis": 5.12
  },
  "signal_integrity": {
    "snr_db": 28.4,
    "sensor_status": "CALIBRATED_VALID",
    "clipping_detected": false
  },
  "local_classification": {
    "anomaly_score": 0.89,
    "score_semantics": "uncalibrated_model_output",
    "model_fingerprint": "sha256:<Klassifikatormodell-Hash>",
    "inference_device": "NPU"
  }
}
```

</details>

Die Inferenzregel verifiziert Signal-Rausch-Verhältnis (*Signal-to-Noise Ratio*, SNR), Kalibrierstatus, Amplitudenbegrenzung (*clipping*), Aktualität und Messüberdeckung anhand eines versionierten Profils. Mangelhafte Signalintegrität stellt ein eigenständiges Überwachungsereignis dar: Sie darf nicht verworfen und als Abwesenheit von Gefahr fehlinterpretiert werden. Die Felder `boot_id` und `sequence` dienen der Erkennung von Duplikaten und Neustarts, setzen jedoch eine vertrauenswürdige Herkunftsüberprüfung voraus. Der Ereigniszeitstempel wird mit der Uhrenunsicherheit und der Eingangszeit abgeglichen. Ein unkalibrierter `anomaly_score` darf keinesfalls mit der physischen Ausfallwahrscheinlichkeit gleichgesetzt werden.

### 7.1. Signieren von Merkmalen statt Urteilen

Petro Sidliarchuk formulierte ein Prinzip, das für verteilte Sensornetze wegweisend ist: Eine digitale Signatur beweist die Urheberschaft, beweist jedoch nicht die materielle Wahrheit [[9]](#src-9). Signiert ein Sensor ein vorgefertigtes Urteil („Ziel detektiert“, „Schwellenwert überschritten“), kann die Empfängerseite das Urteil lediglich akzeptieren oder verwerfen, vermag jedoch nicht nachzuvollziehen, auf welchen Rohdaten es beruhte. In den im selben Repository dokumentierten Feldversuchen hielt ein akustischer Detektor eine stabile Grundfrequenz von $84\,\text{Hz}$ samt Oberschwingungen und meldete ein Ziel, obwohl die Schallquelle eine mehrere hundert Meter entfernte Hochspannungsleitung war. Der Detektor funktionierte spezifikationsgemäß; das Merkmal, auf das er reagierte, gehörte an diesem Einsatzort jedoch zu einem völlig anderen Objekt.

Ein evidenzbasierter Vertrag auf Sensorebene erfordert daher drei Garantien:

1. **Messungen und Kontext werden signiert.** Die signierte Nutzlast umfasst Rohmerkmale, physikalische Einheiten, Zeitfenster, Quellen-URN, Sequenznummer, Verarbeitungspipeline-Version sowie Integritätsmetadaten. Das abgeleitete Urteil kann ergänzend beigefügt werden, ersetzt jedoch niemals die Primärdaten.
2. **Bekannte physikalische Schranken werden geprüft.** Firmware und Empfänger validieren ein explizites Wertebereichsmodell (*Validity Envelope*). Eine Plausibilitätsprüfung kann ohne ein vollständiges physikalisches Modell nicht alle unmöglichen Kombinationen eliminieren; eine Signatur macht eine Fehlmessung nicht wahr.
3. **Versionierte Regeln werden reproduziert.** Der empfangende Knoten kann die Schlussfolgerung nur dann deterministisch nachvollziehen, wenn Modellversionen, Regelsätze und Kalibrierparameter bekannt sind. Eine Diskrepanz ist Anlass für eine Systemdiagnose, kein automatischer Beweis einer Kompromittierung: Ursache können unterschiedliche Revisionsstände, numerische Rundungsdrifts oder verloren gegangener Kontext sein.

Dieses Vorgehen harmoniert mit der RATS-Architektur (*Remote ATtestation procedureS*), in der ein attestierendes Gerät Nachweise (*Evidence*) liefert und ein separater Verifizierer diese anhand eigener Richtlinien bewertet [[10]](#src-10). Das Beispiel der Hochspannungsleitung demonstriert eine weitere fundamentale Grenze: Selbst authentische, signierte Merkmale können Fehlinterpretationen unterliegen; Regeln der Ebene L3 müssen daher stets den räumlichen und betrieblichen Gesamtkontext einbeziehen.

Die Bereitstellung prädikativer Schranken und Gültigkeitsgrenzen (*Validity Envelopes*) für digitale Signalprozessoren (DSP) über NATS-Broker und Ingress-Gateways wird in [Kapitel 33](ch33-inter-system-knowledge-exchange-and-model-teaching.md) vertieft: Dort wird der formale Wissenstransfer vom Expertensystem zu Signalprozessoren sowie die Behandlung von Gegenbeispiel-Warteschlangen systematisch beschrieben.

### 7.2. Nicht-funktionale Schnittstellen und implizite Annahmen

Ein signifikanter Teil von Ausfällen in verteilten cyber-physischen Komplexen entsteht nicht an funktionalen Schnittstellen (Nachrichtenformate, Busprotokolle), sondern durch unzureichend dokumentierte, implizite Annahmen einzelner Module übereinander. Garlan, Allen und Ockerbloom zeigten am Beispiel von Softwarekomponenten, dass Wiederverwendbarkeit regelmäßig an solchen architektonischen Fehlanpassungen (*Architectural Mismatch*) scheitert [[11]](#src-11). In cyber-physischen Systemen treten zu den Softwareannahmen physikalische Kopplungen hinzu: gemeinsame Stromschienen, auf denen der Anlaufstrom eines Motors die Versorgungsspannung empfindlicher Sensoren einbrechen lässt; elektromagnetische Einkopplungen benachbarter Leistungsteile; schleichende Sensordrifts nach thermischen Zyklen oder Bus-Jitter anstelle eines Totalausfalls. Die Annahmenmatrix jedes Moduls muss daher neben Schnittstellenformaten auch jene physikalischen Umgebungsbedingungen explizieren, die das Modul als stabil voraussetzt, damit Abweichungen als eigenständige Anomalieereignisse erfasst werden können.

### 7.3. Epistemische AI-RAN-Bitübertragungsschicht (Epistemic PHY): Softwaredefinierte Signalebene als Generator sensorischer Evidenz

Die in den Abschnitten 7.1 und 7.2 formalisierten Verträge setzen grundsätzlich einen transparenten Kommunikationskanal voraus. Unter realen industriellen Einsatzbedingungen (Fabrikhallen mit massiven elektromagnetischen Störungen, urbane Mehrwegeausbreitung, Zonen elektronischer Kampfführung mit aktivem Jamming und GNSS-Spoofing) stellt der Funkkanal jedoch eine hochgradig volatile und feindselige Umgebung dar. Klassische Expertensysteme betrachten drahtlose Netzwerke als intransparente Blackbox, die im Störfall schlagartig abreißt. In einer modernen kybernetischen Architektur transformiert sich die physikalische Schicht (Layer 1) des Hochfrequenzpfads dagegen in einen primären, evidenzbasierten Sensor des physikalischen Umgebungszustands [[22]](#src-22).

Dieses Paradigma wird auf Basis von Software-Defined Radio (SDR) und der offenen **AI-RAN-Architektur (Open RAN mit GPU-beschleunigtem cuPHY)** implementiert. Im Gegensatz zu anwendungsspezifischen ASICs verlagert die cuPHY-Pipeline die Verarbeitung von I/Q-Abtastwerten auf hochparallele Tensor-Kerne. Dies ermöglicht es, simultan zur Nutzdatendemodulation radiophysikalische Spektralmetriken in harter Echtzeit zu erfassen und direkt in symbolische Datalog-Prädikate zu überführen.

Als Schlüsselmetrik für die Kanalintegrität dient der Jamming-Index (*Jamming Index*, $\mathrm{JammingIndex}$):

```math
\mathrm{JammingIndex} = \frac{P_{\mathrm{total}} - P_{\mathrm{pilot}}}{N_0 \cdot B_{\mathrm{channel}}}.
```

Parameter der Formel zur Störsendererkennung:
- $P_{\mathrm{total}}$ — gesamte integrierte Hochfrequenzleistung im Nutzkanalband (in dBm oder Watt);
- $P_{\mathrm{pilot}}$ — kohärent akkumulierte Nutzleistung der Pilotsignale der Basisstation (DMRS / SSB);
- $N_0$ — spektrale Rauschleistungsdichte des thermischen Eigenrauschens der Empfangsstufe (W/Hz);
- $B_{\mathrm{channel}}$ — effektive Kanalbandbreite (Hz);
- $\mathrm{JammingIndex}$ — dimensionslose Kennzahl für das Vielfache, um welches die Leistung künstlicher Störsignale den natürlichen Rauschteppich übersteigt.

```mermaid
flowchart LR
    accTitle: Architektur der Überführung der AI-RAN-Bitübertragungsschicht in Evidenzfakten
    accDescr: SDR-Antenne, cuPHY-Basisbandverarbeitung, Spektralmetrik-Extraktion, Übersetzung in Datalog-Fakten und situationsbezogene Schutzaktion.

    ANT["<b>Funkkanal / SDR</b><br/>I/Q-Abtastung 100 MHz"] --> CUPHY["<b>AI-RAN L1 (cuPHY)</b><br/>FFT, DMRS-Korrelation,<br/>CUDA-Pipeline"]
    CUPHY --> EXT["<b>Metrik-Extraktor</b><br/>SNR, Doppler, JammingIndex"]
    EXT --> TRANS["<b>Symbolischer Übersetzer</b><br/>Datalog-Prädikate"]
    
    TRANS --> D1["channel_jammed(probe_1)."]
    TRANS --> D2["gnss_spoofed(sat_cluster)."]
    
    D1 & D2 --> EXP["<b>Expertensystem L3</b><br/>Evaluator & Schutzschleuse"]
    EXP --> ACT["<b>Situative Aktion:</b><br/>FHSS-Frequenzsprung /<br/>Optischer Fluss (Fail-Closed GNSS)"]

    classDef rf fill:#e1f5fe,stroke:#0288d1,stroke-width:2px,color:#01579b;
    classDef soft fill:#ede7f6,stroke:#512da8,stroke-width:2px,color:#311b92;
    classDef fact fill:#fff9c4,stroke:#fbc02d,stroke-width:2px,color:#f57f17;
    classDef act fill:#e8f5e9,stroke:#2e7d32,stroke-width:2px,color:#1b5e20;

    class ANT,CUPHY rf;
    class EXT,TRANS soft;
    class D1,D2 fact;
    class EXP,ACT act;
```

#### 7.3.1. Geschlossener Kreis der ingenieurtechnischen Berechnung des epistemischen Funkkanals

1. **Laufzeitsteuerung und Rechenfluss (Control Flow & Runtime Decisions):**
   Kontinuierlich erfasste Funkmetriken werden durch formale Sicherheitsregeln interpretiert:
   - Liegt $\mathrm{JammingIndex} \le 3{,}0$ (normales Hintergrundrauschen), gilt der Übertragungskanal als vertrauenswürdig; die Telemetrie wird als `PHYS_TRANSPORT_VALID` klassifiziert;
   - Bei $3{,}0 < \mathrm{JammingIndex} \le 12{,}0$ initiiert das System einen pseudozufälligen Frequenzsprungbetrieb (FHSS) oder schaltet das Modulations- und Codierschema (MCS) auf orthogonale Spreizcodes mit maximaler Bandspreizung (DSSS) um;
   - Übersteigt $\mathrm{JammingIndex} > 12{,}0$ oder treten Anzeichen inkohärenter Doppler-Drifts im GNSS-Signal auf, leitet die Inferenzmaschine sofort das Prädikat `gnss_unreliable` ab und sperrt Satellitennavigationsdaten. Der Autopilot autonomer Agenten wechselt deterministisch auf Trägheitsnavigation (IMU) und optischen Fluss (Fail-Closed bezüglich GNSS), wodurch ein Abfangen des Systems über manipulierte Raumkoordinaten zuverlässig verhindert wird.
2. **Hardware-Dimensionierung und infrastrukturelle Schranken (Hardware Dimensioning):**
   - Die Verarbeitung von 100-MHz-I/Q-Abtastwerten in cuPHY auf NVIDIA Jetson AGX Orin oder RTX 6000 Ada erfordert eine Speicherbandbreite von $B_{\mathrm{mem}} \ge 1{,}6\,\text{GB/s}$ bei 16-Bit-Quantisierung;
   - Die Berechnung einer 4096-Punkte-FFT und die Bestimmung von $\mathrm{JammingIndex}$ beanspruchen auf Tensor-Kernen lediglich $18\,\mu\text{s}$. Damit wird eine Störsenderattacke detektiert, lange bevor höhere Protokollschichten des TCP/IP-Stacks in Pufferüberläufe laufen.
3. **Praktisches numerisches Rechenbeispiel (Worked Numerical Example):**
   An Bord eines fahrerlosen Transportsystems registriert der Funkempfänger einen Anstieg der Gesamtleistung im 2,4-GHz-Band auf $P_{\mathrm{total}} = -45\,\text{dBm}$ ($31{,}6\,\text{nW}$), während die Nutzleistung der Pilotsignale der Basisstation auf $P_{\mathrm{pilot}} = -82\,\text{dBm}$ ($6{,}3\,\text{pW}$) einbricht. Bei einer Kanalbandbreite von $B = 20\,\text{MHz}$ und einem thermischen Rauschteppich von $N_0 \cdot B = -101\,\text{dBm}$ ($0{,}08\,\text{pW}$):

```math
\mathrm{JammingIndex} = \frac{31{,}6 \times 10^{-9} - 0{,}0063 \times 10^{-9}}{0{,}08 \times 10^{-12}} \approx 395 \gg 12{,}0.
```

   **Ingenieurtechnische Entscheidung:** Das Expertensystem erkennt unverzüglich einen gezielten elektronischen Störangriff, verwirft alle Reconnect-Versuche, leitet einen atomaren Übergang in den geschützten Autonomie-Modus (`AutonomousSafeHold`) ein und aktiviert optische Odometrie zur sicheren Rückkehr zur Basis.

---

## 8. Degradierungsmodi bei Verbindungsabbruch

Die fundamentale Sicherheitsanforderung an ein cyber-physisches Expertensystem lautet: **Funktionale Sicherheit darf niemals von der permanenten Verfügbarkeit einer Netzwerkverbindung abhängen**. Bricht die Verbindung zur Zentrale ab, durchläuft der Edge-Knoten die im nachfolgenden Zustandsdiagramm spezifizierten Übergänge.

```mermaid
stateDiagram-v2
    accTitle: Degradierungsmodi eines Edge-Knotens bei Verbindungsabbruch
    accDescr: Der Knoten wechselt bei Verbindungsverlust vom Normalbetrieb in den lokalen Modus, nach Ablauf der Autorisierungsfristen in den sicheren Halt und kehrt über den Zustandsabgleich in den Normalbetrieb zurück.

    [*] --> ConnectedNormal: Verbindung stabil
    ConnectedNormal --> DegradedLocalOnly: Verbindungsverlust zum Zentrum
    DegradedLocalOnly --> AutonomousSafeHold: Autorisierungsfrist abgelaufen
    DegradedLocalOnly --> Reconciliation: Verbindung wiederhergestellt
    AutonomousSafeHold --> Reconciliation: Verbindung wiederhergestellt und autorisiert
    Reconciliation --> ConnectedNormal: Abgleich abgeschlossen

    state ConnectedNormal {
        [*] --> FullCoordination
        FullCoordination: Lokale Regeln
        FullCoordination: Streaming-Ereignissynchronisation
        FullCoordination: Globale Sperren
    }

    state DegradedLocalOnly {
        [*] --> IsolatedEnforcement
        IsolatedEnforcement: Inferenz nach gecachten Regeln
        IsolatedEnforcement: Aktionen nur nach gültiger lokaler Richtlinie
        IsolatedEnforcement: Append-Only-Protokollpuffer
    }

    state AutonomousSafeHold {
        [*] --> FailsafeMode
        FailsafeMode: Verifizierter domänenspezifischer Degradierungsmodus
        FailsafeMode: Autonome Verfahrbewegungen gesperrt
        FailsafeMode: Warten auf menschlichen Befehl
    }
```

Der autonome Modus wird durch domänenspezifische Sicherheitsrichtlinien definiert. Ein Verbindungsverlust erzeugt nicht automatisch Autonomie-Rechte der Stufe A2: Erforderlich sind gültige lokale Vollmachten, eine vertrauenswürdige Zeitbasis sowie ein bestätigter Betriebskontext. Ohne Netzwerk kann der Knoten von einem im Zentrum soeben widerrufenen Zertifikat keine Kenntnis haben; zeitliche Befristungen autonomer Aktionen begrenzen dieses Restrisiko, eliminieren es jedoch nicht vollständig. Zeitfenster (wie etwa 30 Minuten) sind systemspezifische Auslegungsparameter. Der sichere Zustand muss für jedes Objekt individuell definiert werden: Ein unkonditioniertes Stillsetzen kann unter Umständen gefährlicher sein als ein geregelter Weiterbetrieb. Ein manipulationssicheres lokales Audit-Protokoll ist zwingende Voraussetzung für den späteren Abgleich; eine kryptografische Signatur allein beweist jedoch noch keine Lückenlosigkeit des Protokolls.

## 9. Zustandsabgleich nach Wiederherstellung der Verbindung

Nach dem Wiederaufbau der Netzwerkverbindung müssen Zentrale und Edge-Knoten die während der Trennung unabhängig akkumulierten Zustände abgleichen: lokale Ereignisse, Alarme und Schutzentscheidungen des Knotens einerseits sowie globale Richtlinienänderungen, Zertifikatswiderrufe und Rollenmodifikationen des Zentrums andererseits. Das folgende Diagramm zeigt diesen Abgleichprozess.

```mermaid
flowchart LR
    accTitle: Zustandsabgleich nach Wiederherstellung der Verbindung
    accDescr: Der lokale Puffer des Edge-Knotens und Modifikationen des zentralen Expertensystems fließen in den Abgleichdienst; die Konfliktrichtlinie priorisiert funktionale Sicherheit, Versionskonflikte werden an den Ingenieur eskaliert.

    EDGE_BUF["<b>Lokaler Knotenpuffer</b><br/>Ereignisse, ausgelöste Regeln,<br/>Alarme, logische Zeitstempel"] --> REC["<b>Abgleichdienst (Reconciliation)</b>"]
    CENTRAL_KB["<b>Zentrales Expertensystem</b><br/>Neue Vorschriften, widerrufene Zertifikate,<br/>Rollenänderungen"] --> REC

    REC --> POL{"<b>Konfliktrichtlinie</b>"}
    POL -->|"Sicherheitspriorität"| MERGE["Einheitlicher Audit-Graph"]
    POL -->|"Versionskonflikt"| ESC["Eskalation an den Chefingenieur"]
    MERGE --> OUT_ACK["Abgleich abgeschlossen"]

    classDef src fill:#e1f5fe,stroke:#0288d1,stroke-width:2px,color:#01579b;
    classDef proc fill:#ede7f6,stroke:#512da8,stroke-width:2px,color:#311b92;
    classDef gate fill:#fff8e1,stroke:#e65100,stroke-width:2px,color:#bf360c;
    classDef ok fill:#e8f5e9,stroke:#2e7d32,stroke-width:2px,color:#1b5e20;
    classDef err fill:#ffebee,stroke:#c62828,stroke-width:2px,color:#b71c1c;

    class EDGE_BUF,CENTRAL_KB src;
    class REC proc;
    class POL gate;
    class MERGE,OUT_ACK ok;
    class ESC err;
```

Logische Lamport-Uhren garantieren die Eigenschaft: Geht ein Ereignis einem anderen kausal voraus, ist sein Zeitstempel kleiner [[12]](#src-12). Der Umkehrschluss ist für zwei beliebige Zeitstempel unzulässig: Numerische Werte allein beweisen keine Kausalität. Vektoruhren können Nebenläufigkeit im Rahmen eines definierten Nachrichtenaustauschs auflösen [[13]](#src-13), messen jedoch keine physikalischen Zeitintervalle, wie sie für Zeitableitungen oder Autorisierungsfristen zwingend benötigt werden. Kausale Nachrichtenordnung und die simultane Einwirkung einer physikalischen Erschütterung auf zwei Sensoren sind kategorial verschiedene Phänomene. Für einen belastbaren Abgleich werden beide Zeitbegriffe mit explizit spezifizierten Unsicherheitsgrenzen benötigt.

Die Konfliktlösungsrichtlinie folgt dem fundamentalen Invariantenprinzip „Safety First“. Hat der Knoten während des Inselbetriebs einen Antrieb aus Sicherheitsgründen verriegelt und sendet die Zentrale nach Wiederverbindung einen planmäßigen Startbefehl für eine neue Produktionscharge, darf der zentrale Befehl die Sicherheitsabschaltung keinesfalls rückwirkend überschreiben. Die Entriegelung darf ausschließlich durch eine explizite Bestätigung des verantwortlichen Ingenieurs vor Ort erfolgen.

## 10. Sichere Wissensaktualisierung: Kryptografisch signierte Bundles

Aktualisierungen von Regeln, Modellen und Ontologien auf Edge-Geräten dürfen niemals über ungesicherte Skripte erfolgen. Ein Wissens-Bundle (*Knowledge Bundle*) ist ein vorkompiliertes, typisiertes und kryptografisch signiertes Software-Artefakt:

```math
\mathrm{KnowledgeBundle} = \{\mathrm{Rules},\; \mathrm{Models},\; \mathrm{Schemas},\; \mathrm{Signatures},\; \mathrm{ValidFrom},\; \mathrm{ValidTo}\}
```

Bestandteile des Wissens-Bundles:

- $\mathrm{Rules}$ bezeichnet den Regelsatz, $\mathrm{Models}$ die Modelle und $\mathrm{Schemas}$ die Datenstrukturschemata;
- $\mathrm{Signatures}$ enthält die digitalen Signaturen des Bundles; $\mathrm{ValidFrom}$ und $\mathrm{ValidTo}$ definieren Beginn und Ende der zeitlichen Gültigkeit;
- $\mathrm{KnowledgeBundle}$ bezeichnet ein einzelnes Bundle; die geschweiften Klammern gruppieren dessen Attribute als Verbunddatentyp.

Das Bundle umfasst somit nicht allein Regelsätze und Modelle, sondern bindet Schemata, Signaturen und Gültigkeitsgrenzen untrennbar aneinander. Die Aufzählung der Felder determiniert noch nicht den Algorithmus der Signaturprüfung oder die Behandlung des Ablaufs; diese Logik wird durch das standardisierte Aktualisierungsprotokoll festgelegt.

Das Aktualisierungsverfahren umfasst fünf Schritte:

1. Der neue Regelsatz wird im Simulator der CI/CD-Pipeline verifiziert, insbesondere durch Regressionstests gemäß [Kapitel 16](ch16-expert-systems-architecture.md).
2. Aus Regeln, Modellen und Schemata wird ein binäres Bundle mit einer deterministischen SHA-256-Prüfsumme kompiliert.
3. Das Bundle wird von mehreren autorisierten Rolleninhabern digital signiert; das Zielgerät akzeptiert das Bundle nur bei Erreichen des definierten Signatur-Schwellenwerts (*Quorum*).
4. Der Edge-Knoten verifiziert die Signaturen gestützt auf einen Hardware Root of Trust (z. B. TPM-Modul und Secure Boot).
5. Der Knoten führt die Umschaltung nach einem atomaren Aktivierungsprotokoll durch; scheitert die Validierung, greift das System auf das vorherige, nach aktueller Richtlinie noch zulässige Bundle zurück.

TUF (*The Update Framework*) trennt kryptografische Vertrauensrollen und unterstützt Schwellenwert-Signaturen [[14]](#src-14). Die Sicherheit hängt von Schlüsselhierarchien, Versionsnummern und Metadaten-Ablauffristen ab; eine beliebige Anzahl von Signaturen schützt nicht vor Replay-Angriffen veralteter Bundles. Ein Rollback darf niemals bereits widerrufene Berechtigungen oder verwundbare Artefakte reaktivieren. Das Prüfmanifest muss Abhängigkeiten zwischen Regeln, Modellen, Schemata und Laufzeitumgebung deklarieren, und die Aktivierung muss die Zustandskompatibilität des Schätzers sicherstellen. Die Normenreihe IEC 62443 [[15]](#src-15) definiert grundlegende Anforderungen an die industrielle Cybersicherheit; das Bundle-Verfahren bedarf stets einer domänenspezifischen Verifikation.

Die aktuelle TUF-Spezifikation [[16]](#src-16) trennt die Integritätsprüfung einer Datei strikt von deren Installation und semantischer Eignung. Zwei Signaturen mit demselben privaten Schlüssel stellen keine zwei unabhängigen Stimmen dar. Metadatenfristen erfordern eine vertrauenswürdige Zeitbasis, und das atomare Umschalten des Bundles garantiert noch nicht die interne Zustandskompatibilität eines bereits laufenden Reglers.

Ein Edge-Knoten benötigt selten die vollständige Wissensbasis des Gesamtbetriebs. Ein Knoten mit fokussiertem Aufgabenprofil erhält eine gezielte Teilmenge: Shards der relevanten Dokumentenfamilien oder ein Modul, das auf das terminologische Profil des Geräts zugeschnitten ist ([Kapitel 7](ch07-knowledge-base-typology.md)). Diese Teilmenge bildet ein eigenständiges Bundle mit eigener Signatur und Datei-Hashes; das Knotenprofil ist im Manifest verankert. Ein reduziertes Profil führt dazu, dass nicht enthaltene Fakten als „unbekannt“ (*unknown*) und nicht als „falsch“ gewertet werden: Betrifft eine Anfrage eine Domäne außerhalb des Profils, antwortet der Knoten mit `UNKNOWN` und postuliert keine Abwesenheit. Ein Update, das einen Shard unberührt lässt, verändert dessen Hash nicht; der Knoten lädt ausschließlich modifizierte Shards nach. Die Bedingungen hierfür sowie die Shard-Aufteilung werden in Abschnitt 10 von [Kapitel 32](ch32-high-performance-knowledge-packs-mmap-and-harvesting.md) formalisiert.

### 10.1. Kreuzarchitektonischer Isomorphismus von EISA-Bytecode-Interpretern

Die kryptografische Signatur eines Wissens-Bundles garantiert die Authentizität des Urhebers und die Integrität des Binärartefakts. In heterogenen industriellen Großanlagen wird der Regel-Mikrocode jedoch auf grundverschiedenen Prozessorarchitekturen ausgeführt: Server-CPUs für Vorab-Analysen (x86-64), zentrale Bordcomputer (ARM64 Cortex-A), isolierte Sicherheitsinseln (RISC-V) und dedizierte Inferenz-Coprozessoren (FPGA/EPU) [[23]](#src-23).

Ein identisch kompilierter Binärregelsatz der EISA (*Epistemic Instruction Set Architecture*) muss unabhängig von der zugrunde liegenden Mikroarchitektur zwingend ein bitgenau identisches Ergebnis liefern. Kreuzarchitektonische Audits decken regelmäßig drei kritische Abweichungsursachen auf:
1. **Endianness-Kollisionen:** Eine Serialisierung ohne feste Festlegung der Byte-Reihenfolge verfälscht numerische Gewichte und Entitäts-IDs auf Big-Endian- oder umschaltbaren Plattformen;
2. **Hardware-Traps bei unaligned Access:** Das Lesen von 64-Bit-Ganzzahlen an ungeraden Speicheradressen löst auf strikten ARM Cortex-M- oder RISC-V-Kernen ein Hardware-Ausnahmesignal (`Alignment Fault`) aus, während x86-64-CPUs den Zugriff unbemerkt mit Bus-Stalls emulieren, was den Fehler in der Simulation maskiert;
3. **FPU-Optimierungsdrift:** Abweichungen in Hardware-Gleitkommaeinheiten (FMAC-Instruktionen, Rundungsmodi und Registerbreiten) führen zu Drifts in den niederwertigsten Bits, was in Randfällen konträre logische Urteile hervorruft.

Um vollständigen Determinismus zu erzwingen, wird das **Kriterium strikter Ausführungs-Isomorphie** eingeführt:

```math
\forall \mathcal{P}, \mathbf{X}: \quad \mathrm{State}(\mathcal{A}_1, \mathcal{P}, \mathbf{X}) \equiv \mathrm{State}(\mathcal{A}_2, \mathcal{P}, \mathbf{X}).
```

wobei der Zustand der Inferenzmaschine durch folgendes Tupel formal definiert ist:

```math
\mathrm{State} = \langle R_0 \dots R_7, \mathrm{EIR}, \mathrm{EFLAGS}, \mathrm{ProofTreeHash} \rangle.
```

Parameter des Kriteriums für kreuzarchitektonischen Isomorphismus:
- $\mathcal{P}$ — kompiliertes EISA-Regelprogramm oder Wissenspaket;
- $\mathbf{X}$ — normalisierter Eingangsvektor von Telemetrie- und Faktenwerten;
- $\mathcal{A}_1, \mathcal{A}_2$ — beliebige physische Hardwareplattformen (z. B. x86-64 und RISC-V);
- $R_0 \dots R_7$ — 64-Bit-Universalregister des EISA-Interpreters;
- $\mathrm{EIR}$ — Epistemic Instruction Register (Befehlszeiger);
- $\mathrm{EFLAGS}$ — Statusregister (Flags `flag_zf`, `flag_uf`, `flag_cf`);
- $\mathrm{ProofTreeHash}$ — kryptografischer SHA-256-Digest des Beweisbaums, der die vollständige Kette angewendeter Inferenzregeln festschreibt.

```mermaid
flowchart TB
    accTitle: Schema zur Verifikation des kreuzarchitektonischen EISA-Isomorphismus
    accDescr: Wissens-Bundle wird auf heterogene Zielknoten übertragen; Registerzustände und ProofTreeHash werden bitgenau verglichen.

    BUNDLE["<b>Signiertes EISA-Wissens-Bundle</b><br/>Little-Endian, 8-Byte-Ausrichtung"] --> V1["x86-64-Knoten<br/>(Server-Backend)"]
    BUNDLE --> V2["ARM64-Knoten<br/>(Bordcomputer)"]
    BUNDLE --> V3["RISC-V-Knoten<br/>(ASIL-D-Sicherheitsinsel)"]
    BUNDLE --> V4["FPGA/EPU-Knoten<br/>(Hardware-Pipeline)"]

    V1 --> S1["State(x86) = &lang;R, EFLAGS, Hash&rang;"]
    V2 --> S2["State(ARM) = &lang;R, EFLAGS, Hash&rang;"]
    V3 --> S3["State(RISC-V) = &lang;R, EFLAGS, Hash&rang;"]
    V4 --> S4["State(EPU) = &lang;R, EFLAGS, Hash&rang;"]

    S1 & S2 & S3 & S4 --> CMP{"Bitgenauer Vergleich<br/>State(A1) ≡ State(A2)"}
    CMP -- "100% Übereinstimmung" --> PASS["Verifikations-Gate: PASS<br/>Isomorphismus bewiesen"]
    CMP -- "Diskrepanz" --> FAIL["Sicherheits-Gate: FAIL<br/>REFUSAL_ARITHMETIC_DIVERGENCE<br/>(Fail-Closed)"]

    classDef src fill:#e1f5fe,stroke:#0288d1,stroke-width:2px,color:#01579b;
    classDef node fill:#ede7f6,stroke:#512da8,stroke-width:2px,color:#311b92;
    classDef check fill:#fff9c4,stroke:#fbc02d,stroke-width:2px,color:#f57f17;
    classDef ok fill:#e8f5e9,stroke:#2e7d32,stroke-width:2px,color:#1b5e20;
    classDef err fill:#ffebee,stroke:#d32f2f,stroke-width:2px,color:#b71c1c;

    class BUNDLE src;
    class V1,V2,V3,V4,S1,S2,S3,S4 node;
    class CMP check;
    class PASS ok;
    class FAIL err;
```

#### 10.1.1. Geschlossener Kreis der ingenieurtechnischen Berechnung der EISA-Zustandsäquivalenz

1. **Laufzeitsteuerung und Rechenfluss (Control Flow & Runtime Decisions):**
   Beim Laden eines Regel-Bundles erzwingt der Boot-Loader drei strukturelle Invarianten:
   - Alle 64-Bit-Felder und Instruktionen sind zwingend an 8-Byte-Grenzen ausgerichtet (`offset % 8 == 0`), 32-Bit-Werte an 4-Byte-Grenzen; direktes Dereferenzieren unausgerichteter Pointer ist strikt untersagt;
   - Die Byte-Reihenfolge ist ausnahmslos Little-Endian (LE); auf Big-Endian-Zielsystemen erfolgt der Speicherzugriff ausschließlich über Hardware-Intrinsics für Byte-Reversal (`__builtin_bswap64`);
   - Sämtliche logischen Gewichte werden in den ganzzahligen Wertebereich $[0, 65535]$ (oder Q16.16) skaliert, was den Einsatz von Gleitkommaoperationen in der sicherheitskritischen Inferenzlogik vollständig eliminiert.
   Weicht im automatisierten Regressionstest auch nur ein einziges Bit im $\mathrm{ProofTreeHash}$ ab, wird das Vertrauenszertifikat des Knotens widerrufen und der Befehl `REFUSAL_ARITHMETIC_DIVERGENCE` (Fail-Closed) ausgelöst.
2. **Hardware-Dimensionierung und infrastrukturelle Schranken (Hardware Dimensioning):**
   - Die natürliche Ausrichtung verhindert kostspielige Split-Lock-Zyklen auf x86 sowie Hardware-Traps auf ARM/RISC-V, wodurch der Speicherzugriff auf exakt einen AXI-Bustakt ($5\,\text{ns}$ bei $200\,\text{MHz}$) gedrückt wird;
   - Der Speicher-Overhead durch Füllbytes (*Padding*) beträgt weniger als $3{,}2\,\%$ des gesamten Paketvolumens.
3. **Praktisches numerisches Rechenbeispiel (Worked Numerical Example):**
   Während einer Konformitätsprüfung nach DO-178C Level A wird ein Paket aus 1.024 Sicherheitsregeln parallel auf x86-64 (Intel Xeon), ARM64 (Cortex-A78AE) und RISC-V (SiFive U74) ausgeführt. Bei identischen Eingangsvektoren $\mathbf{X}$ aus 64 Telemetrieparametern beenden alle drei Plattformen die Inferenz nach exakt 256 EISA-Befehlen mit dem identischen SHA-256-Beweisbaum-Hash `0x8F4A2C19...` und gesetztem Flag `flag_zf = 1`. **Ingenieurtechnische Entscheidung:** Der kreuzarchitektonische Isomorphismus ist zu 100 % bewiesen; das Bundle wird für den Einsatz im heterogenen Produktionscluster freigegeben.

---

## 11. Moderne Werkzeuge zur Verarbeitung und zum Lernen aus Zeitreihendaten

Für das Softwarebeispiel einer Aufgaben-Warteschlange sind zunächst exakte Ereigniszeitstempel und ein deterministischer Simulator erforderlich. Darauf aufbauend können Latenzbehandlung, Signalprüfungen temporal-logischer Eigenschaften oder modellprädiktive Regler integriert werden; kein Werkzeug vermag eine unzureichende Eingangsdatenqualität zu kompensieren.

| Werkzeug | Rolle im Systementwurf | Nicht garantiert |
|---|---|---|
| Apache Flink [[17]](#src-17) | Event-Time-Fensterung, Behandlung verspäteter Datensätze | Ein Watermark beweist nicht, dass alle physischen Ereignisse empfangen wurden |
| RTAMT [[18]](#src-18) | Signalprüfung mittels Signal Temporal Logic (STL) | Der Monitor bewertet gegebene Traces; Zukunftsoperatoren erfordern Wartezeit |
| do-mpc [[19]](#src-19) | Simulation, Zustandsschätzung und modellprädiktive Regelung (MPC) | Modellzulässigkeit beweist nicht die Sicherheit der realen Anlage |
| Sequenzielle Statistik- und Modellanpassung | Sensordrift-Erkennung, Regimewechsel, Lastprognosen | Gelernte Muster beweisen keine Kausalität und ändern Schutzregeln nicht autonom |

Ein *Watermark* in Flink signalisiert den Verarbeitungsfortschritt gemäß einer definierten Latenzrichtlinie; reale physikalische Ströme können diese Schranke verletzen. Für RTAMT müssen diskrete oder kontinuierliche Zeitsemantik, Abtastperiode, Zeithorizont und Lückenbehandlung explizit fixiert werden. Zudem unterliegt die Python-Umgebung durch ANTLR-Abhängigkeiten Versionsbeschränkungen, die vorab zu prüfen sind. Für do-mpc sind dynamische Prozessmodelle, Stellgrößenbeschränkungen und ein deterministisches Ausfallverhalten bei Nichtkonvergenz des Optimierers erforderlich.

Im Bereich der Datenstromanalyse sind Changepoint-Erkennung, Residuenanalyse von Prognosen und Systemidentifikation von hohem Nutzen. Streaming-Verfahren führen inkrementelle Modellaktualisierungen durch und vergleichen Vorhersagen mit nachfolgenden Messungen. Regelaktionen beeinflussen jedoch künftige Messdaten: Das Ausbleiben eines Schadens nach einem Eingriff beweist nicht, dass der Ausgangszustand unkritisch war. Das Audit-Log muss Aktion, Betriebsmodus, Modellversion und Zeitstempellatenz lückenlos festhalten. Neue Parameter und Schwellenwerte erfordern separate Genehmigungsverfahren; Active Learning berechtigt das System keinesfalls zu gefährdenden Erkundungsmanövern an der realen Anlage.

Die Bibliothek River stellt Werkzeuge für kontinuierliches Lernen bereit; ihr ADWIN-Detektor (*ADaptive WINdowing*) vergleicht statistische Momente über adaptive Fensterteilungen [[20]](#src-20). Signifikante Änderungen im Prognoseresiduum können Indikatoren für eine notwendige Neukalibrierung oder einen Regimewechsel sein. Sie beweisen jedoch keine Ursachen und autorisieren keine autonome Anpassung sicherheitskritischer Schwellen. Zunächst wird die gespeicherte Prognose evaluiert, anschließend am realen Label nachgelernt; Labels aus der Zukunft dürfen niemals rückwirkend für vergangene Entscheidungen herangezogen werden.

Für verteilte Datenströme, die Edge-Sensoren mit dem Backend verbinden, formulierten Akidau et al. in den Modellen MillWheel und Dataflow das klassische theoretische Fundament der Stromverarbeitung [[21]](#src-21): Die fundamentale Differenzierung zwischen Ereigniszeit (*Event Time*) und Verarbeitungszeit (*Processing Time*), Berechnungen über gleitenden Fenstern und die Festlegung von Low Watermarks ermöglichen es, verspätete und umgeordnete Pakete deterministisch zu handhaben, ohne die kausale Konsistenz logischer Schlüsse zu verletzen.

Zustandsschätzer müssen in der Simulation unter Injektion von Zeitfehlern und Kommunikationsausfällen verifiziert werden. Eine belastbare Evaluierung muss Fehlalarme, unerkannte Gefahren, Reaktionslatenzen und den zulässigen Zustandsraum (*State Space*) simultan quantifizieren.

## 12. Durchgängige Ingenieurbeispiele

Die folgenden Szenarien stellen didaktische Architekturmodelle zur Verantwortungsverteilung dar und erheben keinen Anspruch auf verifizierte Hardware-Freigaben. Unabhängigkeit der Schutzebenen, Wahrnehmungspräzision, Zeitschranken und das Verhalten bei Common-Cause-Fehlern müssen für jedes Zielsystem gesondert nachgewiesen werden. Auch ein Hardware-Schaltwerk kann mit fehlerhaften Eingangsdaten beaufschlagt werden; der bloße Einsatz eines FPGAs oder von Lockstep-Kernen beweist noch kein physikalisches Invariantenverhalten.

**Stromverteilungseinheit (Power Distribution Unit, PDU).** Betrachten wir das Zusammenwirken aller Ebenen in einer mikrocontrollergesteuerten PDU. Auf den Ebenen L0–L1 werden ein Shunt-Stromsensor und ein Transistor-Thermistor über integrierte Analog-Digital-Wandler abgetastet; ein analoger Hardware-Komparator öffnet bei einem Kurzschluss unverzüglich ein Schutzrelais – vollständig autark und ohne Mitwirkung der Firmware. Auf den Ebenen L2–L3 glättet ein Zustandsschätzer das Temperatursignal, schätzt die Erwärmungsrate und übermittelt bei drohender Überhitzung gemäß der Schutzregel aus Abschnitt 3 einen Drosselungsbefehl über den CAN-Bus an den Antriebswechselrichter. Wird die CAN-Leitung unterbrochen, versetzt ein Watchdog-Timer die Einheit in einen sicheren Zustand: Der Strom wird hardwarenah begrenzt und das Ereignis im geschützten Flash-Speicher protokolliert. Nach Wiederherstellung der Kommunikation wird das signierte Vorfallsprotokoll an das zentrale Expertensystem (L4) übertragen, das den Vorfall mit dem globalen Wissensgraphen abgleicht. Ergibt die Analyse, dass die Transistoren dieser Platinenrevision aus einer Charge mit erhöhtem thermischen Übergangswiderstand stammen (dokumentiert in einem Errata-Report des Herstellers), generiert das Zentrum ein signiertes Wissens-Bundle mit angepassten Abschaltschwellen für die gesamte Flotte.

**Sicherheitsknoten mit getrennter Wahrnehmungs- und Logikebene.** In sicherheitskritischen Systemen höchster Integrität nach ISO 26262 (ASIL D) dürfen Brems- und Lenkaktoren niemals ungeprüft auf probabilistische Wahrnehmungsmodelle vertrauen. Ein Edge-KI-Beschleuniger verarbeitet Kamera- und Lidardatenströme und extrahiert semantische Fakten wie „Hindernis in 12 m Entfernung mit kollidierendem Trajektorienvektor“. Ein funktional sicherer Mikrocontroller mit Lockstep-Kernen (wie in [Kapitel 18](ch18-execution-infrastructure.md) detailliert) empfängt diese Fakten und evaluiert innerhalb deterministischer Zeitschranken die Regeln raum-zeitlicher Sicherheitskorridore. Zwischen dem Mikrocontroller und den Leistungstreibern ist ein FPGA mit einer fest verdrahteten Hardware-Zustandsmaschine geschaltet, die Steuerimpulse, welche die minimale Sicherheitsdistanz verletzen würden, unabhängig von Softwareabstürzen physikalisch blockiert. Fällt die KI-basierte Wahrnehmung aus oder bricht die Ethernet-Verbindung ab, führen Mikrocontroller und FPGA einen kontrollierten Nothalt durch und sichern einen kryptografisch signierten Telemetrie-Snapshot für das Audit.

Beide Architekturen demonstrieren denselben kybernetischen Grundsatz: Jede Ebene verfügt über ein eigenes Modell der Regelstrecke, eine autarke Antwort auf Ausfälle benachbarter Ebenen und verlässt sich für lebenswichtige Schutzfunktionen niemals auf das Netzwerk.

## 13. Dialog Pro et Contra

Zur virtuellen intellektuellen Fachdiskussion sind der Begründer der Kybernetik Norbert Wiener (Professor am MIT, Schöpfer der mathematischen Theorie der Regelung und Signalübertragung in Mensch und Maschine) sowie der führende Forscher cyber-physischer Systeme Edward A. Lee (emeritierter Professor der UC Berkeley, Initiator zeitpräziser Architekturen PRET und Ptolemy) eingeladen. Ihre grundlegenden Paradigmen dynamischer Stabilität und zeitlichen Determinismus in verteilten Systemen bilden ein kompromissloses Fundament für die kritische Überprüfung der vom Autor vorgeschlagenen Ebenenhierarchie L0–L4.

*Disclaimer:* Diese Diskussion stellt eine heuristische virtuelle Modellierung im wissenschaftlich-didaktischen Kontext der Monografie dar. Die Standpunkte und Argumente der Kontrahenten wurden auf der Basis ihrer öffentlich zugänglichen Monografien, Fachaufsätze und Vorlesungen synthetisiert; sie stellen keine wörtlichen Zitate dar und dienen der kritischen Rationalisierung (im Sinne der Popper'schen Falsifikation) der Architekturentscheidungen des Autors.

**Kapitelthese:** Der kybernetische Regelkreis eines Expertensystems lässt sich über ein verteiltes Netzwerk niemals monolithisch und deterministisch schließen. Er muss zwingend nach Zeithorizonten L0–L4 stratifiziert werden: Harte reflexive Begrenzungen des Arbeitsbereichs (Operating Envelope) und Schutzverriegelungen laufen lokal im Submikrosekunden- und Millisekundenbereich (L0–L1) auf FPGAs und RTOSs mit Festkomma-Arithmetik; der Anlagenzustand wird am Edge (L2) durch Filterung von Beobachtungsinnovationen geschätzt; symbolische Regeln werden vom Edge-Knoten (L3) in isomorpher EISA-Bytecode-Ausführung verarbeitet; und das zentrale Expertensystem (L4) aktualisiert Ontologien und Regelsätze ausschließlich über asynchrone, kryptografisch signierte Pakete, was Stabilität trotz Latenz, Jitter, Rauschen und Verbindungsabbrüchen garantiert.

**Dialog mit Norbert Wiener: Rückkopplungsstabilität in verteilten Regelkreisen mit Latenz**

- **Norbert Wiener:** Herr Fedchyk, das Herzstück der Kybernetik ist die negative Rückkopplung. Wird jedoch ein Regelkreis über ein asynchrones, jitterbehaftetes Netzwerk zwischen Edge-Knoten und zentralem Backend geschlossen, treten unvermeidlich Phasenverzögerungen auf. In dynamischen physikalischen Systemen verwandelt eine Phasenverschiebung stabilisierende Gegenkopplung in destruktive Mitkopplung — es entstehen verheerende Schwingungen und Systemaufschaukelungen. Wenn logische Schlüsse auf den Ebenen L3 oder L4 variable Rechenzeiten beanspruchen, bricht das physikalische System zusammen, lange bevor das zentrale Expertensystem seinen Weltzustand konsolidiert hat.

- **Mykola Fedchyk:** Herr Wiener, Ihre mathematischen Theoreme über Phasenverschiebung und dynamische Instabilität sind unumstößlich. Genau deshalb ist es in der vorgeschlagenen Architektur strikt verboten, harte physikalische Stabilisierungsschleifen über ein Netzwerk zu schließen. Die Hierarchie der Zeithorizonte isoliert die Ebenen funktional: Ebene L0 (Sub-Mikrosekunden-Hardwareverriegelungen auf FPGA/ASIC) und Ebene L1 (Millisekunden-PID-Regler unter RTOS) sind vollständig autark auf dem Aktor geschlossen und agieren netzwerkunabhängig. Das zentrale L4-Expertensystem erteilt niemals unmittelbare Steuerbefehle in Echtzeit; es distribuiert asynchron kryptografisch signierte Regel- und Richtlinienpakete, die der L3-Knoten zwischenspeichert und ausschließlich innerhalb deterministisch vorausberechneter Stabilitätskorridore ausführt. Bei Netzwerkabriss oder Latenzüberschreitung geht der Edge-Knoten autonom in einen sicheren degradierten Zustand über, was die asymptotische Stabilität garantiert.

**Fortsetzung des Dialogs mit Norbert Wiener: Entropie, Rauschen und nichtlineare Informationsdämpfung in kybernetischen Regelkreisen**

- **Norbert Wiener:** Betrachten wir das interne Modell der Regelstrecke. Nach dem Conant-Ashby-Theorem muss jeder gute Regler ein Modell des zu regelnden Systems sein. In Ihrem Abschnitt 7 wird kontinuierliche Sensortelemetrie in diskrete semantische Ereignisverträge quantisiert. Doch jede grobe Diskretisierung eines kontinuierlichen Phasenraums zerstört unvermeidlich Shannon-Wiener-Transinformation. Nähert sich ein physikalisches System einer kritischen Verlangsamung oder einem Bifurkationspunkt (Abschnitt 5), divergieren Varianzfluktuationen und Autokorrelation 1. Ordnung exponentiell. Wenn der Edge-Knoten L2/L3 auf verarmten diskreten Kategorien operiert, verliert sein Modell im Augenblick des Phasenübergangs nicht jede Passgenauigkeit und führt zu regulatorischer Blindheit?

- **Mykola Fedchyk:** Herr Wiener, das Conant-Ashby-Theorem wird in unserer Architektur durch eine mehrstufige Zustandsbeobachtung realisiert, bei der Quantisierung niemals ein naiver skalarer Schwellenwertentscheid ist. Auf Ebene L2 verfolgt der Zustandsschätzer (Kalman-Filter oder Alpha-Beta-Beobachter) kontinuierlich den latenten dynamischen Zustandsvektor $\mathbf{x}_k$ und die Fehlerkovarianzmatrix $\mathbf{P}_k$. Die Quantisierung in Ereignisverträge erfolgt nicht anhand von Rohsignalen, sondern auf Basis von Kalman-Innovationsresiduen und Chi-Quadrat-Hypothesentests ($\chi^2$). Zudem werden synergetische Indikatoren kritischer Verlangsamung — ansteigende gleitende Varianz $\sigma^2$ und Lag-1-Autokorrelation $\mathrm{AC}(1)$ — als kontinuierliche Zeitreihenmetriken direkt in Hardwareregistern berechnet. Das diskrete semantische Faktum „Phasenübergang steht bevor“ wird erst gefeuert, wenn eine informationstheoretische Entropieschranke mathematisch verifiziert überschritten ist. Dadurch verliert das interne Modell des Edge-Knotens auch an kritischen Bifurkationspunkten nie den kontinuierlichen Kontakt zur physikalischen Realität.

**Dialog mit Edward Lee: Determinismus multithreaded Software-Pipelines und Hardware-Zeitunterstützung**

- **Edward A. Lee:** Herr Fedchyk, die standardmäßigen Abstraktionen moderner Betriebssysteme — Threads, Mutexe, Queues und Go-Channels — sind bezüglich der physikalischen Zeit von Natur aus nicht-deterministisch. In Ihrem Entwurf werden Sensoraufnahme, Filterung und Zustandsschätzung über Software-Goroutinen abgebildet. Standard-Scheduler und Garbage-Kollektoren gewähren jedoch keinerlei Worst-Case-Execution-Time-Garantien (WCET). Ohne zeitpräzise Prozessorarchitekturen (Precision-Timed Architectures, PRET) und synchron-reaktive Semantik sind Ihre temporalen Verträge und logischen Uhren bloße Scheingenauigkeit, die vor Zertifizierungsstellen für cyber-physische Systeme nicht bestehen kann.

- **Mykola Fedchyk:** Herr Lee, Ihre fundamentale Kritik an threadbasierter Nebenläufigkeit und fehlender Zeitsemantik in Allzwecksprachen ist vollkommen berechtigt. Daher ist die Go-basierte Datenpipeline in unserer Monografie ausdrücklich auf die analytischen Ebenen L2–L3 (Zustandsschätzung, Telemetrie-Logging und Regelauswertung) beschränkt, wo das System mit Zeitstempeln und gleitenden Ereignisfenstern operiert, statt Mikrosekunden-Deadlines zu garantieren. Harte Zeitanforderungen auf den Ebenen L0–L1 werden dedizierten Mikrocontrollern (ARM Cortex-R, Infineon AURIX Lockstep) unter formal verifizierten Echtzeit-Mikrokernen (wie seL4) mit hardwaregebundenen Zeitgebern übertragen. An der Schnittstelle zur Sensordatenaufnahme fungiert ein Zulassungstor: Weicht der Jitter eines Telemetrierahmens über den Grenzwert $\Delta t_{\mathrm{max}}$ ab, wird der Frame nicht in den Modellzustand eingerechnet, sondern in ein Fehlertrace umgeleitet, um veraltete Messwerte aus dem Regelkreis fernzuhalten.

**Fortsetzung des Dialogs mit Edward Lee: Semantik der physikalischen Zeit bei der Zustandsabstimmung und thermodynamische Irreversibilität**

- **Edward A. Lee:** Wenden wir uns Abschnitt 9 zu — dem Zustandsabgleich nach Wiederherstellung abgebrochener Netzwerkverbindungen. Sie schlagen Vektoruhren oder konfliktfreie replizierte Datentypen (CRDTs) vor. Doch logische Lamport-Uhren erfassen Kausalität *innerhalb von Softwarenachrichten*, nicht physikalische Gleichzeitigkeit in der relativistischen Raumzeit. Wenn während einer Netzwerkpartition der lokale Knoten L3 und das zentrale Backend L4 unabhängig voneinander physikalische Aktionen autorisiert haben, lassen sich diese physischen Eingriffe nicht über einen Drei-Wege-Merge zusammenführen oder durch ein Rollback rückgängig machen. Mechanischer Verschleiß, Reaktoraufheizung oder Materialverformungen sind thermodynamisch irreversibel. Wie stellt Ihre Architektur Systemintegrität sicher, wenn logische Konfliktfreiheit auf die physikalische Irreversibilität der realen Welt prallt?

- **Mykola Fedchyk:** Herr Lee, Sie berühren die gefährlichste Fehlannahme in verteilten cyber-physischen Systemen: die unkritische Übertragung von Web-Datenbankabstraktionen auf physikalische Industrieanlagen. Unsere Architektur erzwingt eine unverletzliche Invariante: monotonen physikalischen Konservatismus. Im autonomen Betrieb im degradierten Modus ist es dem L3-Knoten kategorisch untersagt, irreversible Aktionen auszuführen, die den Arbeitsbereich der Anlage erweitern; erlaubt sind ausschließlich restriktive Eingriffe, die den Arbeitsbereich (Operating Envelope) entlang bewiesener Ljapunow-stabiler Funktionen verengen (Drehzahlreduktion, Strombegrenzung, Übergang in sicheren Stillstand). Vektoruhren und CRDTs werden bei Reconnection ausschließlich auf Informations- und Diagnosezustände angewendet (Fehlerzähler, Abgleich kryptografischer Telemetrie-Hashketten). Der physikalische Zustand wird niemals zurückgerollt: Das zentrale L4-System führt einen vorwärtsgerichteten epistemischen Abgleich durch, berechnet den akkumulierten Restlebensdauerverbrauch (RUL) aus dem manipulationssicheren Ereignisprotokoll und gibt die Anlage erst nach formaler Sicherheitsverifikation wieder für den voll koordinierten Betrieb frei.

**Synthese:** Kybernetische Zuverlässigkeit eines Expertensystems entsteht nicht durch den vergeblichen Versuch, ein verteiltes Netzwerk in einen verzögerungsfreien Idealkanal zu verwandeln, sondern durch strikte zeitliche Stratifikation von Schutz- und Entscheidungsbefugnissen. Indem schnelle Schutzreflexe in deterministischer L0–L1-Hardware isoliert werden, kontinuierliche mathematische Zustandsschätzung auf L2 erfolgt, byteäquivalenter EISA-Bytecode auf L3 ausgeführt wird und zentrale L4-Aktualisierungen auf asynchrone kryptografische Autorisierungen beschränkt bleiben, wahrt das System seine Regelungsstabilität über beliebige Netzwerkausfälle hinweg und garantiert physikalische Sicherheitsinvarianten selbst bei vollständiger Isolation der Knoten.

## Fazit

Ein sicherer kybernetischer Regelkreis eines Expertensystems wird strikt nach Zeithorizonten modularisiert: Hardware-Verriegelungen und deterministische L1-Controller reagieren im Mikro- und Millisekundenbereich; Ebene L2 schätzt den internen Zustand; der Edge-Knoten L3 wendet gecachte Regelsätze an; das zentrale Expertensystem L4 aktualisiert Wissensbestände ausschließlich über kryptografisch signierte Bundles. Ashbys Gesetz der erforderlichen Varietät begründet, warum ein industrielles Expertensystem eine differenzierte Fehlerontologie anstelle binärer Statusflags benötigt; das Good-Regulator-Theorem von Conant und Ashby belegt die Notwendigkeit expliziter Objektmodelle zwischen Sensoren und Regelwerken.

Das didaktische Softwarebeispiel demonstrierte den fundamentalen Unterschied zwischen einer naiven numerischen Ableitung und einer modellbasierten Zustandsschätzung für ein konkretes Rauschprofil; es ersetzt jedoch keinen formalen Sicherheitsnachweis des geschlossenen Regelkreises. Negative Tests verifizierten das Verhalten bei Zeitlücken, Umordnungen und nicht-finiten Zahlenwerten. Kryptografisch signierte Merkmale garantieren verifizierbaren Kontext bei gesicherter Schlüsselbindung, beweisen jedoch nicht die materielle Richtigkeit physikalischer Messungen. Degradierungsmodi und Zustandsabgleich erfordern domänenspezifische Richtlinien; logische Uhren können physikalische Zeitintervalle nicht substituieren.

Die Grenzen dieses Kapitels sind präzise gesteckt: Filterkoeffizienten und Schwellenwerte dienen didaktischen Zwecken; reale Betriebsparameter müssen aus physikalischen Rauschmodellen und Sicherheitsanalysen abgeleitet werden. Das Free-Energy-Prinzip wurde als theoretische Analogie, nicht als verbindliche Softwareanforderung eingeführt. Notfall- und Schutzfunktionen müssen stets durch zertifizierte Sicherheitsarchitekturen implementiert werden, nicht durch das Expertensystem allein.

### Architektonisches Fazit und physische Rückkopplung

Die Kapitel 16 bis 22 formen die geschlossene ingenieurtechnische Vertikale von Teil IV. Die Gesamtreferenzarchitektur ([Kapitel 16](ch16-expert-systems-architecture.md)), der Technologie-Stack ([Kapitel 17](ch17-implementation-stack.md)) und die Hardware-Infrastruktur ([Kapitel 18](ch18-execution-infrastructure.md)) schaffen das materielle Fundament; der Weg von der Frage zum Beweis ([Kapitel 19](ch19-from-question-to-evidence.md)), normatives Schließen ([Kapitel 31](ch31-syllogistic-reasoning-and-relation-lattices.md)) und die Erklärungskomponente ([Kapitel 20](ch20-explanation-engine.md)) gewährleisten beweisbare Inferenzketten; die Aktionsautorisierung ([Kapitel 21](ch21-from-recommendation-to-action.md)) und der kybernetische Regelkreis mit Rückkopplung ([Kapitel 22](ch22-cybernetics-edge-to-backend.md)) überführen analytische Empfehlungen schließlich deterministisch in die physische Realität.

### Weiterer Erkenntnisweg

Das thematisch folgende [Kapitel 23](ch23-knowledge-base-verification.md) eröffnet [Teil V](part-05-verification-and-learning.md), der der formalen Verifikation von Regeln, der Testpyramide für Wissensbasen und dem Aufbau von Safety Cases gewidmet ist. Jeder industrielle Regelkreis erfordert den mathematisch fundierten Nachweis der Verlässlichkeit und Widerspruchsfreiheit seiner Wissensbasis, bevor er für den Produktivbetrieb zugelassen werden kann.

## Fragen zur Selbstprüfung

1. Wie lässt sich anhand des Gesetzes der erforderlichen Varietät formal bewerten, ob ein Klassifikator mit lediglich zwei Ausgangsstatus für eine Anlage mit 1.024 distinkten Fehlerzuständen ausreichend ist?
2. Was postuliert das Good-Regulator-Theorem von Conant und Ashby, und welche architektonische Konsequenz ergibt sich daraus für den Einsatz von Zustandsschätzern im Regelkreis?
3. Warum generierte die naive Zeitableitung der Temperatur 520 Fehlauslösungen, und welchen systematischen Preis (in Bezug auf Phasenverschiebung/Latenz) zahlt das Alpha-Beta-Filter für deren Eliminierung?
4. Warum ist die Ausführung eines großen Sprachmodells (LLM) auf Ebene L1 (beispielsweise im Notabschaltkreis einer Industrieturbine) sicherheitstechnisch unzulässig?
5. Warum müssen Edge-Sensoren primäre Merkmale statt fertiger Urteile signieren, und welche fundamentale Lehre liefert das Fallbeispiel des akustischen Detektors an der Hochspannungsleitung?
6. Welche Mechanismen verhindern, dass ein Edge-Knoten im Zustand `DegradedLocalOnly` unkontrolliert Fehlentscheidungen akkumuliert?
7. Wie löst die Konfliktrichtlinie den Fall, dass ein Edge-Knoten während eines Netzwerkausfalls einen Motor aus Sicherheitsgründen sperrt, die Zentrale jedoch nach Wiederverbindung einen planmäßigen Startbefehl erteilt?
8. Warum erfordert die Aktualisierung von Gewichten eines neuronalen Netzes auf einem Edge-Controller dasselbe kryptografische Signatur- und Quorumverfahren wie die Modifikation symbolischer Regeln?
9. Wie ermöglichen synergetische Indikatoren des Critical Slowing Down (CSD) – Autokorrelation erster Ordnung $\rho_1$ und Varianz $\mathrm{Var}(x)$ – die frühzeitige Detektion drohender Bifurkationen vor dem physischen Erreichen von Alarmgrenzen?
10. Worin unterscheidet sich das Konzept logischer Lamport-Uhren von physikalischer astronomischer Zeit, und warum sind für einen sicheren Zustandsabgleich nach Verbindungsabbrüchen zwingend beide Zeitstempelarten erforderlich?
11. Wie verhindert ein Hardware-Begrenzer des Betriebsbereichs (*Operating Envelope Limiter*) mit ganzzahligem Derating im Format Q16.16 den thermischen Durchbruch eines Antriebs und garantiert die Einhaltung der Sicherheitsinvariante $P_{\mathrm{actual}} \le P_{\mathrm{derated}}$ in $O(1)$ Takten?
12. Auf welche Weise überführt eine softwaredefinierte Bitübertragungsschicht (AI-RAN / cuPHY) das rohe Hochfrequenzspektrum und Interferenzmetriken (Jamming Index $\mathrm{JammingIndex}$) in strukturierte sensorische Datalog-Evidenzen und eliminiert damit den Nichtdeterminismus probabilistischer Klassifikatoren?
13. Warum ist für Expertensysteme der kreuzarchitektonische Isomorphismus von EISA-Interpretern kritisch, und wie garantieren eine natürliche Speicherausrichtung sowie der Verzicht auf Hardware-`float` bitgenaue Inferenzäquivalenz $\mathrm{State}(\mathcal{A}_1) \equiv \mathrm{State}(\mathcal{A}_2)$ zwischen x86-64, ARM64 und RISC-V?

## Glossar

| Deutscher Begriff | Englischer Begriff | Kurzerläuterung |
|---|---|---|
| Kybernetik | Cybernetics | Wissenschaft von Steuerung, Regelung und Kommunikation in Systemen mit Rückkopplung |
| Rückkopplung | Feedback | Rückführung des Ausgangs- oder Wirkungssignals auf den Eingang des Reglers |
| Varietät | Variety | Anzahl unterscheidbarer Zustände einer Menge, vorzugsweise im logarithmischen Maß |
| Gesetz der erforderlichen Varietät | Law of requisite variety | Nur Varietät im Regler kann die Varietät von Störungen absorbieren (Ashby) |
| Good-Regulator-Theorem | Good regulator theorem | Jeder gute Regler eines Systems muss ein Modell dieses Systems sein (Conant & Ashby) |
| Zustandsschätzer | State estimator | Algorithmus zur Rekonstruktion verborgener Zustandsvariablen aus verrauschten Messungen |
| Alpha-Beta-Filter | Alpha-beta filter | Einfacher Schätzer für Position und Änderungsrate mit konstanten Filtergewichten |
| Peripherie / Edge | Edge | Rechenknoten in unmittelbarer Nähe von Sensoren und physikalischen Aktoren |
| Viable System Model | Viable System Model | Beers kybernetisches Modell aus fünf Subsystemen zur Aufrechterhaltung der Lebensfähigkeit |
| Active Inference | Active inference | Minimierung freier Energie wahlweise über Wahrnehmungsanpassung oder gezielte Aktion |
| Kybernetik zweiter Ordnung | Second-order cybernetics | Kybernetik, die den Beobachter als integralen Teil des beobachteten Systems modelliert |
| Degradierungsmodus | Degraded mode | Eingeschränkter Betriebsmodus bei Ausfall von Teilkomponenten zur Aufrechterhaltung der Sicherheit |
| Zustandsabgleich | Reconciliation | Konsolidierung von Zuständen, die während einer Netzwerktrennung unabhängig entstanden sind |
| Logische Uhr | Logical clock | Monotoner Zähler zur Ordnung von Ereignissen nach der Kausalrelation „happened-before“ |
| Wissens-Bundle | Knowledge bundle | Kryptografisch signiertes Artefakt aus Regeln, Modellen und Schemata zur Knotenaktualisierung |
| Remote Attestation | Remote attestation | Sicherheitsverfahren, bei dem ein Verifizierer Integritätsnachweise eines Geräts prüft |
| Watermark (Zeitstempel) | Watermark | Fortschrittsmarke der Ereigniszeit gemäß Latenzrichtlinie; kein Beweis vollständiger Daten |
| Modellprädiktive Regelung | Model predictive control | Optimierungsbasierte Wahl von Stellgrößen anhand von Modellprognosen über einen Zeithorizont |
| Shard | Shard | Partitioniertes Teilsegment einer Wissensbasis, das von einem Knoten verwaltet wird (Kapitel 7) |
| Knotenprofil | Node profile | Spezifikation der Dokumentenfamilien oder Terme, die ein Knoten im Bundle erhält |
| Gültigkeitskorridor | Validity envelope | Numerische Werte- und Ratenschranken, innerhalb derer ein physikalisches Signal als valide gilt |
| Betriebshüllkurven-Begrenzer | Operating envelope limiter | Hardwarebasierte oder deterministische Schranke, die Stellsignale zwingend auf den physikalisch sicheren Betriebsbereich einschränkt |
| Sättigungs-Derating | Saturation derating | Berechnung von Leistungs- oder Stromgrenzen mit Sättigungsarithmetik (Q16.16) zur Verhinderung von Überhitzung und Zerstörung |
| Epistemische PHY-Schicht | Epistemic PHY | Softwaredefinierte Signalverarbeitungsschicht (AI-RAN), die verifizierte physikalische Evidenzen für die logische Inferenz generiert |
| Kreuzarchitektonischer Isomorphismus | Cross-architectural isomorphism | Eigenschaft des EISA-Interpreters, bitäquivalente Zustände und identische Urteile über heterogene Prozessorarchitekturen hinweg zu garantieren |
| Natürliche Speicherausrichtung | Natural alignment | Platzierung von Datenfeldern an Speicheradressen, die Vielfache ihrer Bytegröße sind (8 Byte für 64-Bit-Worte), zur Vermeidung von Hardware-Traps |

## Abkürzungen

| Abkürzung | Vollständige Bezeichnung | Bedeutung |
|---|---|---|
| AI-RAN | Artificial Intelligence Radio Access Network | Funkzugangsnetzarchitektur mit KI-gestützter Signalverarbeitung auf der Bitübertragungsschicht |
| CAN | Controller Area Network | Robustes serielles Bussystem für Steuergeräte |
| CI/CD | Continuous Integration / Continuous Delivery | Automatisierte Software-Integrations- und Bereitstellungspipeline |
| cuPHY | CUDA Physical Layer | GPU-beschleunigte Software-Bibliothek für 5G/6G-Bitübertragungsschichten auf NVIDIA CUDA |
| DSP | Digital Signal Processor | Digitaler Signalprozessor zur hochfrequenten Erfassung und Filterung von Signalen |
| EISA | Epistemic Instruction Set Architecture | Befehlssatz- und Bytecode-Spezifikation für deterministische Expertensysteme |
| IEC | International Electrotechnical Commission | Internationale Elektrotechnische Kommission |
| KL | Kullback-Leibler | Kullback-Leibler-Divergenz zur Quantifizierung des Unterschieds zweier Verteilungen |
| LE | Little-Endian | Byte-Reihenfolge mit dem niederwertigsten Byte an der niedrigsten Speicheradresse |
| MPC | Model Predictive Control | Modellprädiktive Regelung |
| NATS | Neural Autonomic Transport System (Plattformname) | Hochperformantes Publish-Subscribe-Nachrichtensystem mit subjektbasierter Adressierung |
| ADWIN | ADaptive WINdowing | Adaptives Fensterungsverfahren zur Detektion von Konzeptdrifts in Datenströmen |
| NPU | Neural Processing Unit | Hardwarebeschleuniger für neuronale Netze und Matrixoperationen |
| RATS | Remote ATtestation procedureS | IETF-Architekturstandard für sichere Remote-Attestierungsverfahren |
| SHA-256 | Secure Hash Algorithm, 256 bits | Kryptografische Hashfunktion mit 256-Bit-Ausgabe |
| SNR | Signal-to-Noise Ratio | Signal-Rausch-Verhältnis als Maß der Signalqualität |
| STL | Signal Temporal Logic | Signal-Temporallogik zur formalen Spezifikation von Echtzeitsignaleigenschaften |
| TPM | Trusted Platform Module | Hardware-Sicherheitschip als hardwarebasierter Vertrauensanker |
| TUF | The Update Framework | Sicherheitsarchitektur für manipulationssichere Software-Aktualisierungen |
| URN | Uniform Resource Name | Eindeutiger, persistenter Ressourcenname |
| VSM | Viable System Model | Kybernetisches Modell lebensfähiger Systeme nach Stafford Beer |
| SPS | Speicherprogrammierbare Steuerung | Industrieller Digitalrechner zur Steuerung von Maschinen und Anlagen (engl. PLC) |

## Quellen

1. <a id="src-1"></a>Edward A. Lee, Sanjit A. Seshia. [*Introduction to Embedded Systems: A Cyber-Physical Systems Approach*](https://ptolemy.berkeley.edu/books/leeseshia/). 2nd edition, MIT Press, 2017.
2. <a id="src-2"></a>Norbert Wiener. [*Cybernetics or Control and Communication in the Animal and the Machine*](https://doi.org/10.7551/mitpress/11810.001.0001). MIT Press, Neuauflage 2019 (Erstausgabe 1948).
3. <a id="src-3"></a>W. Ross Ashby. [*An Introduction to Cybernetics*](https://archive.org/details/introductiontocy00ashb). Chapman & Hall, 1956.
4. <a id="src-4"></a>Roger C. Conant, W. Ross Ashby. [*Every Good Regulator of a System Must Be a Model of That System*](https://doi.org/10.1080/00207727008920220). *International Journal of Systems Science*, 1(2), 89–97, 1970.
5. <a id="src-5"></a>R. E. Kalman. [*A New Approach to Linear Filtering and Prediction Problems*](https://doi.org/10.1115/1.3662552). *Journal of Basic Engineering*, 82(1), 35–45, 1960.
6. <a id="src-6"></a>Heinz von Foerster. [*Understanding Understanding: Essays on Cybernetics and Cognition*](https://doi.org/10.1007/b97451). Springer, 2003.
7. <a id="src-7"></a>Stafford Beer. [*The Viable System Model: Its Provenance, Development, Methodology and Pathology*](https://doi.org/10.1057/jors.1984.2). *Journal of the Operational Research Society*, 35(1), 7–25, 1984.
8. <a id="src-8"></a>Karl Friston. [*The Free-Energy Principle: A Unified Brain Theory?*](https://doi.org/10.1038/nrn2787). *Nature Reviews Neuroscience*, 11(2), 127–138, 2010.
8a. <a id="src-8a"></a>Marten Scheffer, Jordi Bascompte, William A. Brock et al. [*Early-warning signals for critical transitions*](https://doi.org/10.1038/nature08227). *Nature*, 461, 53–59, 2009; Hermann Haken. *Synergetic Computers and Cognition: A Top-Down Approach to Neural Nets*, Springer, 2004.
9. <a id="src-9"></a>Petro Sidliarchuk. [*witness-integrity: Witness Integrity for Measurement Devices*](https://github.com/sidliarchukpetro/witness-integrity). GitHub.
10. <a id="src-10"></a>H. Birkholz, D. Thaler, M. Richardson, N. Smith, W. Pan. [*RFC 9334: Remote ATtestation procedureS (RATS) Architecture*](https://www.rfc-editor.org/rfc/rfc9334). IETF, 2023.
11. <a id="src-11"></a>D. Garlan, R. Allen, J. Ockerbloom. [*Architectural Mismatch: Why Reuse Is So Hard*](https://doi.org/10.1109/52.469757). *IEEE Software*, 12(6), 17–26, 1995.
12. <a id="src-12"></a>Leslie Lamport. [*Time, Clocks, and the Ordering of Events in a Distributed System*](https://doi.org/10.1145/359545.359563). *Communications of the ACM*, 21(7), 558–565, 1978.
13. <a id="src-13"></a>C. Fidge. [*Logical Time in Distributed Computing Systems*](https://doi.org/10.1109/2.84874). *Computer*, 24(8), 28–33, 1991.
14. <a id="src-14"></a>Justin Samuel, Nick Mathewson, Justin Cappos, Roger Dingledine. [*Survivable Key Compromise in Software Update Systems*](https://doi.org/10.1145/1866307.1866315). *Proceedings of the 17th ACM Conference on Computer and Communications Security (CCS)*, 61–72, 2010.
15. <a id="src-15"></a>IEC. [*IEC TS 62443-1-1:2009. Industrial Communication Networks: Network and System Security: Part 1-1: Terminology, Concepts and Models*](https://webstore.iec.ch/en/publication/7029). 2009.
16. <a id="src-16"></a>TUF Contributors. [*The Update Framework Specification*](https://theupdateframework.github.io/specification/latest/). Offizielle Spezifikation von Rollen, Schwellenwerten, Versionen und Metadaten-Ablauffristen.
17. <a id="src-17"></a>Apache Software Foundation. [*Flink: Timely Stream Processing*](https://nightlies.apache.org/flink/flink-docs-stable/docs/concepts/time/). Offizielle Dokumentation zu Event Time, Watermarks und Latenzbehandlung.
18. <a id="src-18"></a>Dejan Nickovic, Tomoya Yamaguchi et al. [*RTAMT: Runtime Monitoring Library*](https://github.com/nickovic/rtamt). Dokumentation zur Zeitsemantik, Zeithorizonten und unterstützten Umgebungen.
19. <a id="src-19"></a>do-mpc Contributors. [*Model Predictive Control Python Toolbox*](https://www.do-mpc.com/en/latest/). Dokumentation zu Modellierung, Schätzung und Regelung.
20. <a id="src-20"></a>River Contributors. [*River*](https://riverml.xyz/latest/) und [*ADWIN*](https://riverml.xyz/latest/api/drift/ADWIN/). Dokumentation zu kontinuierlichem Lernen und Changepoint-Detektion in Datenströmen.
21. <a id="src-21"></a>Tyler Akidau et al. [*The Dataflow Model: A Practical Approach to Balancing Correctness, Latency, and Cost in Massive-Scale, Unbounded, Out-of-Order Data Processing*](https://research.google/pubs/the-dataflow-model-a-practical-approach-to-balancing-correctness-latency-and-cost-in-massive-scale-unbounded-out-of-order-data-processing/). *Proceedings of the VLDB Endowment*, 8(12), 1792–1803, 2015.
22. <a id="src-22"></a>O-RAN Alliance. [*O-RAN Working Group 4: Open Fronthaul Interface Specification*](https://www.o-ran.org/specifications). 2023; NVIDIA Corporation. *Aerial SDK cuPHY: High-Throughput 5G/6G Physical Layer Processing on CUDA*, Architecture Whitepaper, 2024.
23. <a id="src-23"></a>RTCA / EUROCAE. [*DO-178C / ED-12C: Software Considerations in Airborne Systems and Equipment Certification*](https://www.rtca.org/). 2011; Mykola Fedchyk. *Epistemic Instruction Set Architecture (EISA): Deterministic Bytecode Specification for Provable Expert Systems*, Technical Report EISA-TR-2024-01, 2024.

---

[← Kapitel 21](ch21-from-recommendation-to-action.md) | [Inhaltsverzeichnis](README.md) | [Teil VII](part-07-runtime-and-knowledge-exchange.md) | [Kapitel 35 →](ch35-self-organizing-expert-systems-synergetics-and-npu-runtime.md)
