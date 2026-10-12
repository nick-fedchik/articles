# Appendix B. Evidence-Governed Expert Systems in Autonomous Robotics and Cyber-Physical Systems

> **Book:** [Architecture of Evidence-Governed Expert Systems](README.md) · Appendices  
> **Table of Contents:** [README.md](README.md)  
> **Author:** [Mykola Fedchyk](about-the-author.md)  
> **Target Audience:** Autonomous platform architects, embedded systems engineers, robotics specialists (ROS 2, Zephyr RTOS), safety-critical real-time system developers  
> **Expected Learning Outcomes:** Decouple probabilistic perception from deterministic action admission; construct an onboard knowledge base without dynamic memory allocation; structure physical actions as formal contracts with precondition and postcondition verification; audit subsystem interfaces and physical models prior to field trials.

---

## Abstract

In safety-critical autonomous systems and mobile robotics (ISO 26262 ASIL D, IEC 61508 SIL 3, DO-178C DAL A), directly coupling neural network perception models or multimodal AI to physical actuators inevitably leads to catastrophic accidents. The physical consequences of a robot's actions (emergency braking, payload release, sharp rudder deflection) are irreversible and expend real kinetic energy. At a control loop frequency of $f = 50\ \text{Hz}$, even a negligible classification or detection error probability of $p = 10^{-3}$ guarantees a hazardous command within an average of $T_{\text{err}} = 1/(p \cdot f) = 20\ \text{seconds}$ of flight or motion.

This appendix resolves this dilemma by introducing a fundamental principle of cyber-physical safety: the strict separation of probabilistic sensory intelligence from deterministic action admission. The author refutes the engineering illusion of the "Human-on-the-Loop" (HOTL) paradigm due to the physical latency of the radio link and human neurophysiology ($T_{\Sigma} = 390 - 780\,\text{ms}$, $S_{\mathrm{react}} \approx 19.5 - 39.0\,\text{m}$), necessitating a transition to a strictly onboard Safe State ontology and emergency sagas (*Ballistic Jettison*, *Dead Reckoning RTB*, *Safe Zeroize*).

A three-tier architecture featuring an **Epistemic Governor** is proposed and examined in detail, wherein action admission is formalized as a deductive proof over Datalog rules in static memory without dynamic `malloc` allocation and backed by defense against GNSS spoofing; a research testbed for synchronous precision physical state spectroscopy based on Infineon TLE90xx and AURIX TC49x (isoUART 2.0); an autonomous Popperian Red Teamer deployed on a Seeed Studio reServer Industrial J501 (NVIDIA Jetson AGX Orin 64GB) and Xilinx Virtex FPGA with a hardware disposition gate ($T_{\mathrm{disp}} < 100\,\text{ns}$, $\mathrm{ZHR} = 1.00$); hardware and software protection of the sensor front-end against Fail-Open anomalies (filtering jagged matrices, IEEE 754 NaN/Inf anomalies, and EMC noise floor bounds), formal action contracts with precondition and postcondition verification, and methodologies for detecting interface defects prior to field trials at the proving ground.

---

## 1. Physical Irreversibility and Cybernetic Limits: Fundamental Difference Between Cyber-Physical Systems and Generative AI

Attempts to transfer architectural patterns from consumer generative artificial intelligence directly into autonomous robotics and cyber-physical systems (*cyber-physical systems*, CPS) represent one of the most perilous engineering fallacies of the modern era. The widespread media metaphor of an autonomous robot as a "language model on wheels" or the concept of end-to-end control (*End-to-End Vision-Language-Action, VLA*)—where a stochastic neural network synthesizes actuator control currents directly from camera pixels—guarantees catastrophic failures in mission-critical deployments.

The underlying cause lies in the fundamental ontological distinction between the virtual space of text token generation and physical reality. This introductory section articulates the foundational physical and cybernetic constraints that mandate a fundamentally different approach to designing onboard intelligence, serving as an engineering preamble to all subsequent sections of this appendix.

### 1.1. Thermodynamic and Material Irreversibility of Physical Action

In general-purpose software systems (search engines, generative chatbots, recommendation engines), a model hallucination or error poses no immediate physical peril: its consequence is merely an erroneous text string rendered on a display, an automated retry, or a page reload. Even in enterprise financial or transaction databases, an automated transactional rollback (`ROLLBACK`) mechanism reliably reverts the system to its initial consistent state.

In cyber-physical systems, software code translates computational verdicts into the redistribution of physical energy:
* modulating the gate drive voltage and pulse widths of power inverter transistors for traction electric motors;
* regulating hydraulic line pressures and modulating brake system valves;
* altering aerodynamic surface deflection angles or adjusting UAV thrust outputs;
* supplying firing current to pyrotechnic payload release actuators.

Every executed physical action—such as `EmergencyBraking()`, `DropPayload()`, or an abrupt $90^\circ$ servo deflection at a velocity of 80 km/h—releases kinetic or thermal energy and irreversibly modifies the operational environment.

In the physical world, no undo operator exists: metal deformation from an impact, stripped gearbox teeth, a burned-out power inverter switch, or an aerodynamic stall cannot be "rectified with a subsequent prompt." Entropy increases irreversibly. Consequently, any action that departs from the platform's safe physical operating envelope must be blocked prior to physical actuation, rather than mitigated after a catastrophic failure has already transpired.

### 1.2. Cybernetic Accumulation of Failure Probability in High-Frequency Loops

While a consumer chatbot user submits approximately one query per minute ($f \approx 0.016\ \text{Hz}$), the onboard control loop of a robotic platform executes computations at tens to hundreds of hertz: $f = 50\ \text{Hz}$ for trajectory planning, $f = 200\ \text{Hz}$ for attitude estimation, and $f \ge 1000\ \text{Hz}$ for inner actuator current loops.

Under such high operational frequency, even an error rate that machine learning benchmarks consider stellar inevitably precipitates rapid failure. The mean time to first critical error ($T_{\text{err}}$) is governed by the relation:

```math
T_{\text{err}} = \frac{1}{p\,f}.
```

Where:
- $T_{\text{err}}$ is the mean operational time preceding the first perception or planning error, in seconds;
- $p$ is the error probability of the neural network model on a single discrete time step (a dimensionless quantity in the interval $[0, 1]$);
- $f$ is the execution frequency of the control loop, in hertz (Hz).

Because the system completes $f$ steps per second, the expected number of generated errors per unit time equals the product $p\,f$. Consequently, the expected waiting time until the first hazardous command is the reciprocal $1/(p\,f)$.

```math
\begin{aligned}
\text{At } p = 10^{-3}\ (99.9\%\ \text{accuracy}),\ f = 50\ \text{Hz} &\implies T_{\text{err}} = \frac{1}{10^{-3} \cdot 50} = 20\ \text{seconds}; \\
\text{At } p = 10^{-3}\ (99.9\%\ \text{accuracy}),\ f = 200\ \text{Hz} &\implies T_{\text{err}} = \frac{1}{10^{-3} \cdot 200} = 5\ \text{seconds}; \\
\text{At } p = 10^{-4}\ (99.99\%\ \text{accuracy}),\ f = 50\ \text{Hz} &\implies T_{\text{err}} = \frac{1}{10^{-4} \cdot 50} = 200\ \text{seconds}\ (3.3\ \text{minutes}).
\end{aligned}
```

In operational field environments, the hypothesis of error independence across consecutive cycles collapses: direct sun glare blinding a camera lens, a dust cloud, raindrops on optical elements, or LiDAR multi-path reflections induce burst errors, wherein erroneous perception persists across dozens of consecutive cycles. This collapses the effective interval to catastrophe to mere fractions of a second.

### 1.3. Insurmountability of the Long-Tail Distribution (Out-of-Distribution Hallucinations)

The common engineering instinct to resolve reliability challenges solely through accumulating larger training datasets, supervised fine-tuning, or reinforcement learning (*RLHF/RLAIF*) encounters a fundamental barrier: a neural network is an empirical statistical interpolator across a high-dimensional feature space.

The physical operating environment of autonomous systems exhibits infinite variability ("long-tail distribution"). A robot inevitably encounters sensory signal combinations that were absent from its training distribution (out-of-distribution scenarios): rare optical glints, anomalous structural obstacles, airframe battle damage, or non-linear ground dynamics. At these critical points, the stochastic model does not signal epistemic uncertainty; instead, it synthesizes arbitrary and frequently hazardous control vectors with high internal confidence (overconfident hallucinations).

Functional safety standards governing the highest safety integrity levels (ISO 26262 ASIL D, IEC 61508 SIL 3, DO-178C DAL A) impose stringent limits on permissible catastrophic failure rates:

```math
\Lambda_{\text{catastrophic}} < 10^{-9}\ \text{failures per operational hour}\ (1\ \text{FIT} = 10^{-9}\ \text{h}^{-1}).
```

No modern neural network architecture can provide a mathematical proof of compliance with this criterion. Physical road testing over $10^9$ operational hours (exceeding 114,000 years of continuous operation) without a single catastrophic fault for every updated model weight release is physically impossible.

### 1.4. Physical and Cognitive Infeasibility of the "Human-on-the-loop" Approach in Terminal Phase

In contemporary debates concerning autonomous unmanned systems, mobile robotics, and critical cyber-physical complexes, a compromise paradigm known as "Human-on-the-loop" (HOTL) is frequently advocated. According to this concept, the platform may navigate, maneuver, and track targets autonomously using stochastic neural network models, while a remote operator retains supervisory control over a telemetric radio datalink, capable of aborting the action via an emergency abort button (*Abort/Cancel Button*) upon detecting an error.

However, from the perspectives of fundamental physics, kinematics, and neurophysiology, this assumption represents a perilous engineering delusion. In the terminal phase of motion (a drone diving at high velocity, high-speed obstacle closing, or emergency evasive maneuvering of a ground vehicle), video compression latency, radio propagation delays, and human cognitive reaction times render external intervention physically too late.

The aggregate system reaction time $T_{\Sigma}$ under a human-on-the-loop architecture comprises five serial physical delays:

```math
T_{\Sigma} = t_{\mathrm{cam\_encode}} + t_{\mathrm{link\_down}} + t_{\mathrm{cogn}} + t_{\mathrm{link\_up}} + t_{\mathrm{actuator}}
```

where:
- $T_{\Sigma}$ is the total reaction time of the supervisory system from the optical manifestation of a hazardous anomaly until physical displacement of aerodynamic control surfaces or brake pads, in seconds ($\text{s}$);
- $t_{\mathrm{cam\_encode}}$ is the frame capture, exposure, and hardware video compression latency (H.264/H.265/AV1), typically $20 - 40\,\text{ms}$;
- $t_{\mathrm{link\_down}}$ is the digital radio downlink transmission latency (accounting for packet buffering, repeaters, and electronic warfare interference), typically $60 - 150\,\text{ms}$;
- $t_{\mathrm{cogn}}$ is the human operator cognitive latency to visually discern the AI hallucination and actuate the emergency abort button, typically $250 - 450\,\text{ms}$;
- $t_{\mathrm{link\_up}}$ is the command uplink transmission latency across jammed RF channels, typically $40 - 100\,\text{ms}$;
- $t_{\mathrm{actuator}}$ is the actuator mechanical response time to deflect rudders or energize braking valves, typically $20 - 40\,\text{ms}$.

Summing these discrete components yields the representative range of total supervisory latency:

```math
T_{\Sigma} \in [390\,\text{ms}, \ 780\,\text{ms}].
```

The kinematic consequence of this latency determines the uncontrolled traverse distance of the platform $S_{\mathrm{react}}$ during the terminal approach:

```math
S_{\mathrm{react}} = v \cdot T_{\Sigma}
```

where $v$ is the closing velocity of the platform toward the target or terrain surface, in meters per second ($\text{m/s}$). At a representative terminal dive or maneuvering velocity $v = 180\,\text{km/h} = 50\,\text{m/s}$:

```math
S_{\mathrm{react}} = 50\,\text{m/s} \times [0.39 \dots 0.78]\,\text{s} \approx 19.5 \dots 39.0\,\text{meters}.
```

From this follows an unequivocal engineering conclusion: if a computer vision neural network hallucinates at a distance of 30 meters from the ground—mistaking a civilian structure or friendly asset for an authorized target—the remote human operator is **physically incapable of preventing impact**. Any genuine safety guarantees can be realized exclusively via an autonomous, onboard, deterministic admission gate (*Onboard Deterministic Safety Gate*).

#### 1.4.1. Actionable Closed Loop of Supervisory Latency Dimensioning

1. **Control Flow & Runtime Decisions:**
   The onboard flight dispatcher executing within the safety microcontroller continuously evaluates the current distance to destination $d_{\mathrm{target}}$ against a dynamic autonomy threshold:
   $d_{\mathrm{crit}} = v \cdot T_{\Sigma}^{\max} + \delta_{\mathrm{margin}}$. As soon as $d_{\mathrm{target}} \le d_{\mathrm{crit}}$, the system unconditionally disables dependency on the radio datalink and transitions into autonomous-veto-only mode (`AUTONOMOUS_VETO_ONLY`). Remote operator commands can no longer serve as an assurance mechanism: the onboard logic core unilaterally vetoes execution in the absence of a cryptographic proof of target or upon violation of deontic prohibitions.
2. **Hardware Dimensioning:**
   The onboard deterministic verifier on the safety microcontroller (Zephyr RTOS) or FPGA must reach admission verdicts with latency $T_{\mathrm{onboard}} \le 2\,\text{ms}$, two orders of magnitude faster than the RF link. This mandates guaranteed sensor sampling at $f_{\mathrm{sensor}} \ge 100\,\text{Hz}$, microsecond-accurate hardware timer allocation, and local residency of spatial geofence ontologies directly in on-chip SRAM of no less than $128\,\text{KB}$.
3. **Worked Numerical Example:**
   A high-speed unmanned aerial platform navigates at $v = 50\,\text{m/s}$. Under electronic warfare interference, downlink video delay is $t_{\mathrm{link\_down}} = 110\,\text{ms}$, camera compression $t_{\mathrm{cam\_encode}} = 25\,\text{ms}$, human cognitive latency $t_{\mathrm{cogn}} = 300\,\text{ms}$, uplink command latency $t_{\mathrm{link\_up}} = 65\,\text{ms}$, and actuator deflection $t_{\mathrm{actuator}} = 25\,\text{ms}$. The total latency is:
   $T_{\Sigma} = 25 + 110 + 300 + 65 + 25 = 525\,\text{ms} = 0.525\,\text{s}$. The uncontrolled travel distance is $S_{\mathrm{react}} = 50 \cdot 0.525 = 26.25\,\text{m}$. If a hallucination occurs at altitude $h = 24\,\text{m} < 26.25\,\text{m}$, the remote operator presses the abort switch after impact. In contrast, the onboard Datalog engine derives a maneuver prohibition within $45\,\mu\text{s}$, the hardware logic gate suppresses the command in $5\,\text{ns}$, and control surfaces are forcefully deflected into a climb at $23.8\,\text{m}$, completely averting disaster.

### 1.5. Deterministic Safe State Ontology and Emergency Saga Selection Tree

In automotive engineering (ISO 26262 ASIL D), entering a safe state (*Safe State*) has a well-defined physical manifestation: controlled deceleration and pulling over to the shoulder (*Safe Stop*). In aerial, underwater, or high-speed field platforms, however, cutting motor power mid-flight or in open water causes an immediate unrecoverable crash or loss of buoyancy.

The onboard expert system utilizes a formal safe state ontology that deterministically maps failure modes and geospatial context into the selection of a transactional emergency saga:

```mermaid
flowchart TD
    accTitle: Logical Decision Tree for Safe State Emergency Saga Selection
    accDescr: Deterministic deduction of the expert core for safe state transition based on spatial map ontology and threat levels.
    FAIL["<b>Invariant Violation Detected / System Fault</b><br/>(Sensor loss, tracking divergence, knowledge deficit)"] --> D1{"Platform Location<br/>per Map Ontology?"}
    D1 -->|"Friendly Zone"| S1["<b>Saga 1: Ballistic Jettison</b><br/>1. Hardware warhead disarm (Disarm)<br/>2. Glide to designated safe ditch sector<br/>3. Activate emergency recovery beacon"]
    D1 -->|"Hostile Zone"| D2{"Risk of Platform Capture<br/>by Adversary?"}
    D2 -->|"Low Risk / Altitude h >= h_safe"| S2["<b>Saga 2: Dead Reckoning RTB</b><br/>1. Disconnect compromised sensors<br/>2. Gyrocompass flight to evacuation point"]
    D2 -->|"High Risk / Altitude h < h_safe"| S3["<b>Saga 3: Safe Zeroize</b><br/>1. Zeroize cryptographic keys<br/>2. Overwrite microcontroller flash memory<br/>3. Controlled impact into uninhabited terrain"]
    classDef err fill:#ffebee,stroke:#c62828,stroke-width:2px,color:#b71c1c;
    classDef dec fill:#fff3e0,stroke:#e65100,stroke-width:2px,color:#bf360c;
    classDef safe fill:#e8f5e9,stroke:#2e7d32,stroke-width:2px,color:#1b5e20;
    class FAIL err;
    class D1,D2 dec;
    class S1,S2,S3 safe;
```

The selection of a safety saga is a logical deduction performed by the expert core, where each step is an atomic transaction verified by postconditions (actuating the detonator interlock line, verifying capacitor discharge, or blowing eFuses).

#### 1.5.1. Actionable Closed Loop of Emergency Saga Selection

1. **Control Flow & Runtime Decisions:**
   The expert system continuously evaluates the spatial predicate `current_zone(ZoneID, HostilityLevel)` and remaining energy budget $E_{\mathrm{rem}}$. When an irreconcilable epistemic deficit arises (such as irreconcilable divergence between GNSS and visual odometry), the Datalog engine fires the transaction corresponding to the selected saga. Every saga step represents an inviolable contract with postcondition verification: when transitioning to Saga 3, the controller validates the status of pin `ZEROIZE_DONE` before cutting power regulator rails.
2. **Hardware Dimensioning:**
   Retention of cryptographic credentials and critical ontologies is isolated in Secure Storage (hardware-accelerated AES-GCM in Secure Enclave), wired to an emergency hardware zeroization bus that discharges in $T_{\mathrm{zeroize}} \le 50\,\mu\text{s}$ powered by a dedicated $100\,\mu\text{F}$ backup capacitor.
3. **Worked Numerical Example:**
   A UAV suffers total camera blinding and odometry divergence at altitude $h = 45\,\text{m}$. The map ontology classifies the coordinate sector as `HostileZone` in critical proximity to the forward line of own troops ($d_{\mathrm{front}} = 1200\,\text{m}$, safe altitude threshold $h_{\mathrm{safe}} = 80\,\text{m}$). Because $h < h_{\mathrm{safe}}$, the decision tree triggers Saga 3: assertion of `DISARM`, EEPROM routing table scrub within $18\,\mu\text{s}$, eFuse blow in $40\,\mu\text{s}$, and servo command to $-15^\circ$ pitch for controlled descent into an empty sector without risk of key or software extraction.

### 1.6. Architectural Roadmap: Deterministic Admission Gate

This reality yields a foundational engineering deduction that underpins the theoretical and practical framework of this entire work:

> [!IMPORTANT]
> **Foundational Principle of Evidence-Governed Cyber-Physical Safety:**
> No probabilistic neural network model (computer vision, SLAM, or VLM) may possess direct write access to actuator control buses. Between stochastic intelligence and the physical world must stand a deterministic **admission gate**: an evidence-governed expert system that formally verifies every proposed action against rigorous physical and normative safety invariants in hard real time.

This appendix serves as a comprehensive engineering manual for constructing such an admission gate. Each subsequent section unpacks a specific tier of defense:

1. **Heterogeneous Onboard Computer and Research Platforms:** Segregating operational workloads across physically isolated computational tiers—a high-throughput NPU for probabilistic perception (NVIDIA Jetson AGX Orin) and a deterministic evidence verifier (Xilinx Virtex UltraScale+ FPGA / Zephyr RTOS lockstep MCU) linked via PCIe DMA with a hardware emergency cut-off latency of $< 5\ \text{ns}$.
2. **Simplex Architecture:** A mathematically rigorous dynamic supervision mechanism wherein an admission arbiter built on Control Barrier Functions instantaneously transfers control from an advanced neural trajectory planner to a formally verified baseline safe controller whenever a safety invariant is threatened.
3. **Onboard Knowledge Base Without Dynamic Memory Allocation:** Designing a Datalog symbolic inference engine operating exclusively in static memory without invocations of heap allocation (`malloc`), eliminating memory fragmentation and guaranteeing bounded worst-case inference latency.
4. **Action as a Contract (Preconditions, Postconditions, Invariants):** Structuring physical actuations as distributed transactional sagas equipped with precondition checking, deterministic runtime monitoring, and unconditional compensating actions upon contract breach.
5. **Pre-Deployment Interface and Physical Model Verification:** A formal mathematical verification methodology for subsystem integration (MIL/SIL/HIL) that exposes kinematic flaws, bus latency jitter, and interface collisions in simulation before deploying physical hardware to proving grounds.

## 2. Heterogeneous Onboard Computer

A mobile robot, UAV, or expeditionary platform is constrained by size, weight, power, and cost (*size, weight, power and cost*, SWaP-C). Mounting an unconstrained general-purpose server onboard rapidly depletes battery reserves, induces severe thermal throttling, and encroaches upon useful payload capacity. Consequently, as established in [Chapter 18](ch18-execution-infrastructure.md), the onboard computational architecture is partitioned across multiple heterogeneous processors fulfilling distinct functional roles.

```mermaid
flowchart TD
    accTitle: Heterogeneous Onboard Robotic Computer Architecture
    accDescr: Workload partitioning across a neural perception processor, a Zephyr RTOS safety microcontroller, and FPGA hardware logic.

    subgraph Perception["<b>1. Perception Pipeline: Soft Real-Time</b>"]
        CAM["Cameras, Stereo Vision"] --> NPU["<b>NPU or SoC</b><br/>• Obstacle detection<br/>• Visual odometry and SLAM<br/>• Environmental semantic predicates"]
        LIDAR["LiDAR, Radar"] --> NPU
    end

    subgraph Reasoning["<b>2. Evidence Reasoning Level: Hard Real-Time</b>"]
        LKB["<b>Onboard Knowledge Base</b><br/>Platform state, spatial geofences,<br/>energy budget"] <--> MCU["<b>Safety Microcontroller</b><br/>• Zephyr RTOS<br/>• Compiled Datalog rules<br/>• Action contract verification"]
    end

    subgraph Arbiter["<b>3. Hardware Veto Unit</b>"]
        FPGA["<b>FPGA or CPLD</b><br/>• Hardware watchdog timers<br/>• PWM signal interlock<br/>• Emergency mask matching"]
    end

    subgraph Physical["<b>4. Actuators</b>"]
        MOT["Traction Motors"]
        STEER["Steering Servos"]
        PAYLOAD["Payload Release"]
    end

    NPU -->|"Semantic facts and proposed actions"| MCU
    MCU -->|"Authorized commands"| FPGA
    FPGA -->|"PWM or CAN FD"| MOT
    FPGA -->|"PWM or CAN FD"| STEER
    FPGA --> PAYLOAD

    classDef perc fill:#fff3e0,stroke:#e65100,stroke-width:2px,color:#bf360c;
    classDef safe fill:#e8f5e9,stroke:#2e7d32,stroke-width:2px,color:#1b5e20;
    classDef gate fill:#ede7f6,stroke:#512da8,stroke-width:2px,color:#311b92;
    classDef act fill:#e1f5fe,stroke:#0288d1,stroke-width:2px,color:#01579b;

    class CAM,LIDAR,NPU perc;
    class LKB,MCU safe;
    class FPGA gate;
    class MOT,STEER,PAYLOAD act;
```

This architectural topology reads strictly top-down: the perception tier merely proposes, the safety microcontroller decides, and programmable logic retains the authority to instantly inhibit signals routed to physical actuators. The computational roles are allocated as follows:

1. **Neural Processing Unit (*neural processing unit*, NPU) or System-on-Chip (SoC).** Representative platforms include NVIDIA Jetson Orin modules, Rockchip RK3588 processors, or Hailo acceleration modules paired with single-board computers. Software environment: embedded Linux running the ROS 2 robotics framework [[1]](#src-1) or compiled standalone binaries implemented in Go or C++. Operational scope: ingesting raw sensor streams, obstacle detection, point cloud construction, and semantic feature extraction. Reliability posture: advanced untrusted controller operating without hard real-time guarantees; it may crash, experience thermal throttling, or produce erroneous classifications without jeopardizing baseline platform stabilization.
2. **Safety Microcontroller.** Representative implementations employ dual-core lockstep microcontrollers with hardware fault containment. The Zephyr RTOS real-time operating system [[2]](#src-2) provides deterministic preemptive thread scheduling and operates entirely without dynamic heap allocation. Operational scope: maintaining the local fact base, executing compiled Datalog rules, and formally validating action preconditions and postconditions. Target integrity level: dictated by hazard analysis and risk assessment, targeting ISO 26262 ASIL D [[3]](#src-3) or IEC 61508 SIL 3 [[4]](#src-4).
3. **Hardware Arbiter on Programmable Logic (FPGA or CPLD).** Operational scope: power bus telemetry monitoring, instantaneous current limit enforcement, and hardware gating of motor PWM drive signals upon exceeding critical roll/pitch thresholds or breaching spatial geofences (*geofencing*). Because programmable logic operates purely at the gate level without software intervention, its response latency is governed strictly by propagation delay and sampling clock rates, completely independent of operating system task schedulers.

### 2.1. Computer Vision Pipeline in Go and NPU

For onboard obstacle detection, eliminating interpreted intermediate runtime layers that introduce nondeterministic latency jitter is essential. The architectural pattern recommended by the author for functional prototypes couples a single-board host computer, an NPU connected via PCIe, and a single statically compiled binary in Go:

```mermaid
flowchart LR
    accTitle: Asynchronous Frame Processing Pipeline in Go and NPU
    accDescr: Decoupling frame acquisition, NPU model inference, and verdict dispatch via bounded channels in Go.

    CAM["V4L2 Camera"] -->|"DMA"| CAP["<b>Goroutine 1: Capture</b><br/>zero-copy buffer management"]
    CAP -->|"Ring Buffer"| NPU_EXEC["<b>Goroutine 2: NPU Model Inference</b><br/>vendor runtime library"]
    NPU_EXEC -->|"Detection Predicates"| BUS["<b>Goroutine 3: Dispatch to Admission Gate</b><br/>micro-ROS, CAN FD, or UART"]

    classDef stage fill:#e1f5fe,stroke:#0288d1,stroke-width:2px,color:#01579b;
    class CAM,CAP,NPU_EXEC,BUS stage;
```

This pattern exhibits three critical design properties. Raw image frames are streamed from the V4L2 kernel driver into direct memory access (DMA) buffers without superfluous memory copies. The neural model executes on the NPU via vendor C/C++ runtimes linked directly via cgo. Goroutines communicate across bounded, non-blocking channels configured with drop-oldest overflow policies, preventing latency queues from accumulating during transient processing spikes. Latency distributions and CPU utilization figures must always be benchmarked under representative operational loads on physical target hardware, rather than extrapolated from vendor promotional datasheets.

### 2.2. Epistemic Governor and Formal Proof of Physical Action

In safety-critical cyber-physical systems and autonomous platforms, the most hazardous architectural flaw is the conflation of control tiers: low-level flight autopilots (PID controllers stabilizing pitch, roll, and yaw on microcontrollers running PX4 or ArduPilot firmware) and neural network vision models (YOLO, CNNs, VLMs) are erroneously granted decision-making authority over targets and kinetic maneuvers.

A classical flight controller governs purely the spatial dynamics of motion, remaining completely oblivious to mission Rules of Engagement (RoE), geospatial exclusion zones, friendly unit locations, and sensory integrity. Conversely, a neural network synthesizes a stochastic hypothesis (e.g., "pixel bounding box classified as target with 86% probability"), yet lacks causal reasoning and cannot bear responsibility for the irreversible consequences of physical action.

To bridge this chasm, the author introduces the **Epistemic Governor**—a deterministic onboard expert system serving as the supreme arbiter of admission for any physical actuation.

```mermaid
flowchart TD
    accTitle: Architecture of the Epistemic Governor for Cyber-Physical Systems
    accDescr: Three-tier topology isolating stochastic perception, deterministic Datalog expert reasoning, and physical actuator dynamics.
    subgraph S1["<b>Tier 1. Stochastic Perception (ML / NPU)</b>"]
        CV["<b>Neural Vision & Direction Finding</b><br/>Probabilistic candidate hypotheses<br/>Stochastic bounding boxes"]
    end
    subgraph S2["<b>Tier 2. Epistemic Governor (Expert System)</b>"]
        KB["<b>Onboard Knowledge Base ZKP4</b><br/>Mission ontology, RoE norms, physical invariants"]
        EVM["<b>Deterministic Inference Engine (EVM)</b><br/>Proof tree synthesis<br/>Datalog in static memory without malloc"]
        KB <--> EVM
    end
    subgraph S3["<b>Tier 3. Mechanical Dynamics (RTOS / Actuators)</b>"]
        FC["<b>PID Regulators & Hardware Interlock</b><br/>Physical actuation permitted ONLY with valid Proof Token"]
    end
    S1 -->|"Candidate Facts"| S2
    S2 -->|"Deterministic Admission (Eligibility & Token)"| S3
    classDef s1 fill:#fff3e0,stroke:#e65100,stroke-width:2px,color:#bf360c;
    classDef s2 fill:#e8f5e9,stroke:#2e7d32,stroke-width:3px,color:#1b5e20;
    classDef s3 fill:#e1f5fe,stroke:#0288d1,stroke-width:2px,color:#01579b;
    class S1,CV s1;
    class S2,KB,EVM s2;
    class S3,FC s3;
```

#### 2.2.1. Mathematical Formalism of Proof of Physical Action

A physical action $A$ (initiating a terminal dive, releasing a payload, energizing an ignition circuit, or abrupt evasive maneuvering) is executed **if and only if** the Expert System synthesizes a formal deductive proof over the active fact base:

```math
\mathrm{PermitAction}(A) \iff \mathrm{KB} \cup \mathrm{ObservedFacts} \vdash_{\mathrm{Datalog}} \mathrm{Eligible}(A) \wedge \neg\mathrm{Conflict}(A)
```

where:
- $\mathrm{PermitAction}(A)$ is the boolean admission predicate governing physical execution of action $A$;
- $\mathrm{KB}$ is the normative mission knowledge base defining geospatial geofences, Rules of Engagement (RoE), and functional safety mandates;
- $\mathrm{ObservedFacts}$ is the set of validated sensory facts that have passed hardware plausibility audits;
- $\vdash_{\mathrm{Datalog}}$ is the deductive entailment operator in Datalog with negation stratified across cycles;
- $\mathrm{Eligible}(A)$ is the predicate certifying satisfaction of all positive qualification criteria;
- $\neg\mathrm{Conflict}(A)$ is the invariant asserting the absence of active defeaters (exceptions or deontic prohibitions).

The deductive admission rule for kinetic or high-energy physical action is formalized as a closed Horn clause in Datalog:

```prolog
permit_kinetic_engagement(TargetID) :-
    target_candidate(TargetID, Class),
    target_hostile(Class),
    inside_geofence(TargetID, AuthorizedZone),
    not friendly_unit_in_radius(TargetID, SafeRadius),
    sensor_consistency(Verified),
    altitude_above_ground(Alt), Alt >= MinSafeAlt,
    mission_rules_valid(Token), token_active(Token).
```

* **Fail-Closed Invariant:** If even a single atomic predicate in the rule body cannot be proven or contains a contradiction (epistemic incompleteness), the system returns an irrevocable `DENY / REFUSAL` verdict. The hardware safety gate implemented in FPGA logic instantly inhibits current to power drives or detonator lines.

#### 2.2.2. Deterministic Sensor Audit Against GNSS Spoofing (Epistemic Consistency)

When hostile electronic warfare systems execute satellite navigation spoofing (GPS Spoofing), conventional flight controllers crash while attempting to compensate for nonexistent spatial drift.

The expert system continuously cross-examines sensory facts against fundamental kinematic axioms:

```math
\left\lVert \Delta \vec{r}_{\mathrm{GPS}} - \int_{t}^{t+\Delta t} \vec{v}_{\mathrm{IMU}} \, dt \right\rVert \le \epsilon_{\mathrm{drift}} + \kappa \cdot \Delta t
```

where:
- $\Delta \vec{r}_{\mathrm{GPS}}$ is the spatial displacement vector reported by the satellite receiver over interval $\Delta t$, in meters ($\text{m}$);
- $\vec{v}_{\mathrm{IMU}}$ is the true platform velocity vector integrated from Inertial Measurement Unit (IMU) gyroscopes and accelerometers, in meters per second ($\text{m/s}$);
- $\epsilon_{\mathrm{drift}}$ is the permissible static sensor channel error, typically $1.5 - 3.0\,\text{m}$;
- $\kappa$ is the dead-reckoning drift divergence coefficient, in meters per second ($\text{m/s}$).

If the observed divergence exceeds the kinematic consistency threshold, the Epistemic Governor executes the following transactional actions:
1. Deterministically derives the fact `sensor_status(gps, falsified)` and evicts GPS coordinates from working memory.
2. Deductively activates fallback navigation rules: inertial dead reckoning and visual flow odometry.

#### 2.2.3. Actionable Closed Loop of Epistemic Action Admission

1. **Control Flow & Runtime Decisions:**
   Upon detecting a violation of kinematic consistency, the Datalog engine asserts `GPS_INVALID` in a single resolution step. The power gate suppresses generation of the `ProofToken`. The flight controller is transitioned into heading-and-optical-flow mode, barring stochastic dive commands.
2. **Hardware Dimensioning:**
   Deterministic Datalog inference in static microcontroller RAM (without `malloc`) completes in $T_{\mathrm{infer}} \le 50\,\mu\text{s}$ for a knowledge base of 150 rules and 500 facts. This is three orders of magnitude faster than NPU neural inference ($15 - 30\,\text{ms}$), guaranteeing response within a single hard real-time timer cycle.
3. **Worked Numerical Example:**
   A UAV cruises at $v = 40\,\text{m/s}$. Across interval $\Delta t = 0.1\,\text{s}$, the satellite receiver reports a coordinate displacement jump $\Delta \vec{r}_{\mathrm{GPS}} = (18.0; \ 0.0; \ 0.0)\,\text{m}$ (an apparent velocity of $180\,\text{m/s}$ caused by spoofing). The IMU integral yields a displacement of $\int \vec{v}_{\mathrm{IMU}} dt = (4.02; \ 0.01; \ -0.05)\,\text{m}$.
   Residual norm:
   $\left\lVert \Delta \vec{r}_{\mathrm{GPS}} - \int \vec{v}_{\mathrm{IMU}} dt \right\rVert = \sqrt{(18.0 - 4.02)^2 + (0.0 - 0.01)^2 + (0.0 + 0.05)^2} \approx 13.98\,\text{m}$.
   With allowable threshold $\epsilon_{\mathrm{drift}} + \kappa \cdot \Delta t = 2.0 + 5.0 \cdot 0.1 = 2.5\,\text{m}$, the disparity $13.98\,\text{m} \gg 2.5\,\text{m}$ triggers fault handling.
   The Epistemic Governor asserts `sensor_status(gps, falsified)` within $18\,\mu\text{s}$, revokes the GPS flight path, and inhibits payload release, preventing a strike on spoofed false coordinates.

## 3. Research Hardware Platforms: Test Benches Based on NVIDIA Jetson AGX Orin and Xilinx Virtex FPGA

To experimentally validate the theoretical models presented in this monograph, the author deployed two complementary industrial-grade hardware test benches that physically segregate statistical perception from deterministic evidence-governed control.

```mermaid
flowchart TD
    accTitle: Topology of the Jetson Orin and Xilinx Virtex Heterogeneous Research Test Bench
    accDescr: Interconnection of the Jetson AGX Orin neural module and Xilinx Virtex hardware EPU array via PCIe Gen4.

    subgraph Host["<b>NVIDIA Jetson AGX Orin (System 1: Neural Perception)</b>"]
        CAM_IN["GMSL2 Cameras / LiDAR / Radar"] --> DLA["Hardware Accelerators NVDLA v2 / PVA"]
        GPU["Ampere GPU (2048 Cores, Tensor Cores)"] --> VLM["Local SLMs / VLMs (TensorRT-LLM)"]
        DLA & VLM --> FACT_GEN["Semantic Fact & Action Proposal Generator"]
        CPU_ORIN["ARM Cortex-A78AE (Split-Lock ASIL-D, PREEMPT_RT)"]
    end

    subgraph Interconnect["<b>High-Speed PCIe Gen4 x8/x16 Bus Interconnect</b>"]
        FACT_GEN ==>|"Direct Memory Access (DMA / Ring Buffers)"| XDMA["Xilinx DMA Controller (XDMA / QDMA)"]
    end

    subgraph FPGA["<b>Xilinx Virtex FPGA (System 2: Hardware Deterministic EPU)</b>"]
        XDMA --> NOC["AXI4 Crossbar Matrix (350–400 MHz)"]
        URAM["UltraRAM / BRAM: On-Chip ZKP4 Knowledge Base (1 cycle)"] <--> NOC
        NOC <--> EPU_ARRAY["EPU Core Array (16..256 Parallel EVM Cores)"]
        EPU_ARRAY --> SHA_PIPE["Pipelined Hardware SHA-256 (%ebx Custody)"]
        SHA_PIPE --> GATE_LOGIC{"Hardware Gate Interlock:<br/>HF & ~ZF & ~DF & CF & ~UF"}
    end

    subgraph Interlock["<b>Physical Safety Interlock (Trigger Latency < 5 ns)</b>"]
        GATE_LOGIC -->|"disp_accept (1)"| PWM_EN["PWM Enable / Power Inverter Gate Drivers"]
        GATE_LOGIC -->|"disp_refuse (0)"| HARD_ESTOP["Hardware Emergency Relay (E-Stop)"]
    end

    classDef host fill:#e8f5e9,stroke:#2e7d32,stroke-width:2px,color:#1b5e20;
    classDef fpga fill:#ede7f6,stroke:#512da8,stroke-width:2px,color:#311b92;
    classDef inter fill:#e1f5fe,stroke:#0288d1,stroke-width:2px,color:#01579b;
    classDef lock fill:#fce4ec,stroke:#c2185b,stroke-width:2px,color:#880e4f;

    class CAM_IN,DLA,GPU,VLM,FACT_GEN,CPU_ORIN host;
    class XDMA,NOC,URAM,EPU_ARRAY,SHA_PIPE,GATE_LOGIC fpga;
    class Interconnect inter;
    class PWM_EN,HARD_ESTOP lock;
```

---

### 3.1. Research Program for NVIDIA Jetson AGX Orin (Perception and Proposal Pipeline — System 1)

The baseline hardware platform for research into the perception pipeline is the industrial edge computer **Seeed Studio reServer Industrial J501**, built around the **NVIDIA Jetson AGX Orin** module (delivering up to 275 TOPS of sparse tensor computing, a 2048-core Ampere-architecture GPU with 64 Tensor Cores, dual NVDLA v2 deep learning accelerators, a PVA v2 programmable vision accelerator, a 12-core ARM Cortex-A78AE CPU cluster, and up to 64 GB of unified LPDDR5 memory with 204,8 GB/s bandwidth housed in a rugged fanless enclosure rated for $-20^\circ\text{C}$ to $+60^\circ\text{C}$).

The platform software stack is configured on **NVIDIA JetPack 6.2** (Linux kernel 5.15 with PREEMPT_RT real-time patching, CUDA 12.6, cuDNN 9, TensorRT 10, Ubuntu 22.04 LTS `aarch64`) and an integrated local runtime environment powered by **Ollama / TensorRT-LLM**.

The perception pipeline possesses no direct write access to physical actuator drivers and functions strictly as a stochastic hypothesis generator: it parses unstructured sensory streams into formalized environmental state predicates and formulates candidate actions $U_{\text{adv}}$. This test bench explores the following applied engineering problems:

1. **Perception-to-Symbolic Translation:**
   * Investigating structural weight quantization (FP8, INT4 AWQ) of local open-weights models (Qwen2.5, Phi-4, Gemma-2) alongside vision-language architectures accelerated via TensorRT-LLM and NVDLA for robust spatiotemporal fact extraction.
   * Leveraging the benefits of Unified Memory Architecture (UMA): co-locating models in the 14B to 32B parameter range directly within shared system VRAM without inter-processor bus copying overhead.
   * Profiling latency percentiles (P50/P90/P99) for transforming unstructured sensor data into typed symbolic fact tuples inside a fully autonomous, air-gapped operational envelope without external network connectivity.
2. **Candidate Action Synthesis and Contract Descriptor Serialization:**
   * Automated synthesis of proposed candidate actions $U_{\text{adv}}$ (target acceleration vectors, waypoint trajectories, requested power allocations) along with associated operational preconditions subject to formal verification.
   * Compact serialization of extracted facts and action arguments into 64-byte aligned binary descriptors structured for direct zero-copy ingestion by the hardware admission gate across PCIe.
3. **Evaluation of Temporal Determinism Under Linux PREEMPT_RT:**
   * Empirical profiling of scheduling latency jitter and Worst-Case Execution Time (*Worst-Case Execution Time*, WCET) during fact synthesis and serialization on ARM Cortex-A78AE cores.
   * Quantifying LPDDR5 memory bus contention between concurrent GPU/DLA inference kernels and real-time CPU threads; verifying that total inference latency $T_{\text{infer}}$ remains strictly within the Fault Tolerant Time Interval (FTTI).
   * Assessing the efficiency of ARM Cortex-A78AE cores operating in **Split-Lock** hardware redundancy mode to validate compliance with ISO 26262 ASIL D requirements.
4. **Low-Latency Ingestion Across PCIe DMA Interconnect:**
   * Managing DMA ring buffers within pinned host shared memory mapped directly into PCIe Base Address Register (BAR) address space.
   * Direct streaming of fact descriptors to the hardware verifier without intermediate kernel system calls or third-party protocol encapsulation overhead.
   * Ingesting primary sensor streams across GMSL2 camera serializers and LiDAR interfaces using zero-copy DMA ring buffers (V4L2 / NVMM), curtailing raw frame ingestion latency into the inference engine to $< 2\ \text{ms}$.

> [!NOTE]
> **Theoretical and Engineering Foundations of Perception Pipelines in Monograph Chapters:**
> - [Chapter 12. Linguistic Analysis and Local Models](ch12-linguistic-analysis-and-local-models.md) — Local language model inference, structured fact extraction, and linguistic templates for closed-loop control systems.
> - [Chapter 18. Execution Infrastructure](ch18-execution-infrastructure.md) — Roofline performance modeling, numerical drift analysis, latency budgeting, and Linux PREEMPT_RT kernel tuning.
> - [Chapter 28. Dual-Mode Expert Systems: Strict Derivation and Advisory Hypothesis](ch28-dual-mode-expert-systems.md) and [Chapter 29. Neuro-Symbolic Architecture](ch29-neuro-symbolic-architecture.md) — Formal division of responsibilities between probabilistic hypothesis generators (System 1) and deterministic evidence verifiers (System 2).

---

### 3.2. Research Program for Xilinx Virtex FPGA (Deterministic Evidence-Governed Control Level — System 2)

Implementing a hardware admission gate mandates selecting a programmable logic family capable of providing deterministic execution with nanosecond-level fixed latency, zero reliance on dynamic memory allocation, and intrinsic hardware fault containment.

#### 3.2.1. Justification of FPGA Selection and Resource Sufficiency

The research test bench employs the **AMD / Xilinx UltraScale+ family (16 nm FinFET process)**—specifically **Virtex UltraScale+ (VU9P, VU13P)** or **Kintex UltraScale+ (KU11P, KU15P)** devices available as PCIe accelerator cards, as well as **AMD Alveo (U50, U200)** computing platforms.

This selection is grounded in rigorous hardware engineering criteria:

1. **Rejection of Legacy Architectures (Virtex-7):**
   The legacy Virtex-7 family (28 nm) is technologically obsolete and unsuitable for mission-critical real-time verification: it is restricted to PCIe Gen3 (creating severe throughput bottlenecks when paired with the Jetson Orin PCIe Gen4 root complex), lacks on-chip UltraRAM (URAM) memory blocks, and necessitates round-trips to external DDR3/DDR4 DRAM, introducing unacceptable stochastic access latency jitter.
2. **Logic Resource Sizing (LUTs and Flip-Flops):**
   * A single deterministic symbolic inference core (Datalog unification engine / hardware deontic invariant evaluator) consumes approximately **2,500 to 4,500 LUTs** and **3,000 to 5,000 flip-flops (FF)**.
   * A research array of 16 to 64 parallel cores for concurrent proof-tree exploration requires **80,000 to 280,000 LUTs**.
   * On-chip infrastructure: an integrated PCIe Gen4 DMA endpoint (XDMA/QDMA IP subsystem) consumes ~25,000 LUTs; an AXI4 Crossbar Interconnect matrix requires ~15,000 to 20,000 LUTs; a 64-stage fully pipelined SHA-256 hash verifier utilizes ~12,000 to 15,000 LUTs and dedicated DSP48E2 slices; combinational safety interlock logic and hardware watchdog timers consume ~2,000 LUTs.
   * The aggregate design requirement spans **150,000 to 350,000 LUTs**.
   * Deploying Kintex UltraScale+ (KU15P: 523k LUTs) or Virtex UltraScale+ (VU9P: 1,182k LUTs; VU13P: 1,728k LUTs) devices is **necessary and sufficient**: it maintains logic utilization within 30% to 55%. This completely prevents routing congestion, guarantees timing closure at the target clock frequency of **300 to 400 MHz** ($T_{\text{clk}} = 2.5–3.3\ \text{ns}$), and reserves routing headroom for experimental architectural extensions.
3. **Internal Memory Architecture (UltraRAM and BRAM vs. External DRAM):**
   * Relying on off-chip dynamic RAM (DDR4/DDR5) for real-time safety invariant verification is unacceptable: access latencies of 50 to 100 ns, periodic row-refresh pauses ($t_{\text{RFC}}$), and bus arbitration contention introduce stochastic jitter incompatible with hard real-time guarantees.
   * UltraScale+ devices incorporate **UltraRAM (URAM)** blocks—synchronous dual-port static memory macros of 288 kbits each (organized as 4K × 72 bits) that cascade into unified arrays without external logic overhead. Specifically, the VU9P device provides 36 MB (270 Mbits) of URAM and 75.9 Mbits of Block RAM (BRAM 36K).
   * This capacity is sufficient to **store the entire compiled normative knowledge base, ontological lattices, defeater tables, and action invariants directly in on-chip memory**. Any rule lookup or term unification completes in **1 to 2 clock cycles ($2.5–5.7\ \text{ns}$)** with absolute temporal determinism (Zero Wait-States) and zero cache misses.

#### 3.2.2. Applied Research Tasks for System 2 on FPGA

1. **Multi-Core Symbolic Inference Array:**
   * Synthesizing a deterministic multi-core logic processor (scaling from 16 to 64 parallel cores) executing compiled Datalog rules and evaluating deontic prohibitions (`MUST_NOT`) concurrently.
   * Implementing a Hardware Task Dispatcher that parallelizes action precondition audits and proof-tree expansions, achieving a throughput exceeding $> 10^7$ rule evaluations per second at a P99 latency of $< 100\ \text{ns}$.
2. **Zero-Jitter On-Chip Knowledge Store (UltraRAM):**
   * Loading binary knowledge packages into the static UltraRAM address space.
   * Investigating contention-free multi-port access: every core in the array retrieves norms, predicates, or defeaters in 1 clock cycle ($2.5–3.3\ \text{ns}$) without off-chip memory transactions.
3. **Pipelined Cryptographic Source Custody (%ebx Hardware SHA-256 Custody):**
   * Synthesizing a 64-stage fully pipelined SHA-256 engine utilizing logic gates and DSP48E2 slices. The engine verifies the cryptographic authenticity of cited normative source texts (matching the hash of the governing standard) in exactly 64 clock cycles ($< 180\ \text{ns}$ at 350 MHz) in parallel with rule execution.
   * Hardware Single Event Upset (SEU) mitigation: activating built-in hardware Error Correction Code (ECC) across BRAM/URAM blocks to ensure knowledge base integrity against ionizing radiation and electromagnetic interference.
4. **Hardware Discrete Admission Gate and Physical Interlock:**
   * Driving discrete admission signal lines `disp_accept` and `disp_refuse` directly to general-purpose I/O pins (LVCMOS / GPIO) on the FPGA package.
   * Constructing emergency cutoff logic purely out of combinational gates with an actuation latency of $< 5\ \text{ns}$, directly disabling optical gate drivers for inverter power switches or opening the coil circuit of a physical emergency stop relay (E-Stop).
   * Enforcing complete physical independence from software: upon any invariant violation, actuator power is severed at the silicon gate level without operating system intervention.

> [!NOTE]
> **Theoretical and Engineering Foundations of Hardware Symbolic Inference in Monograph Chapters:**
> - [Chapter 16. Expert Systems Architecture](ch16-expert-systems-architecture.md) and [Chapter 17. Implementation Stack](ch17-implementation-stack.md) — Inference engine design, register file organization, and formal decoupling of immutable rules from dynamic facts.
> - [Chapter 18. Execution Infrastructure](ch18-execution-infrastructure.md) (Section "Hardware Foundations and Physical Peripheral Limits") — Rationale for FPGA fabrics in deterministic microsecond- and nanosecond-scale processing pipelines.
> - [Chapter 21. From Recommendation to Action: Authority Control and Safe Execution in Production Environments](ch21-from-recommendation-to-action.md) — Hardware admission gates, power circuit interlocks, and deterministic inhibition of unsafe actions.
> - [Chapter 31. Syllogistic Reasoning and Relation Lattices](ch31-syllogistic-reasoning-and-relation-lattices.md) — Deontic norms of obligation, permission, and prohibition (`MUST`, `PROHIBITED`) translated into logic core microinstructions.
> - [Chapter 32. High-Performance Knowledge Packs: Memory Mapping (mmap) and Harvesting Pipeline](ch32-high-performance-knowledge-packs-mmap-and-harvesting.md) — Binary knowledge container layout, 64-byte memory alignment, and zero-copy hardware access inside UltraRAM.
> - [Appendix D. Analog Expert Systems, Neuromorphic Computing, and Hardware Inference](appendix-d-analog-expert-systems-and-neuromorphic-computing.md) — Energy economics of memory access and hardware execution of formal logic.

---

### 3.3. Investigation of the Heterogeneous Tandem "Jetson AGX Orin + Xilinx Virtex PCIe"

The primary research objective centers on integrating both platforms into a unified test apparatus: the **NVIDIA Jetson AGX Orin** carrier board provides a physical **PCIe Gen4 x16 slot (with 8 active lanes)**, enabling direct insertion of the **Xilinx Virtex UltraScale+** accelerator card configured as an Endpoint device.

```mermaid
sequenceDiagram
    accTitle: Interaction Protocol Between Jetson AGX Orin and Xilinx Virtex UltraScale+ via PCIe
    accDescr: Sequence of action descriptor transmission from NPU to EPU with hardware actuator cutoff in 5 ns.
    autonumber
    participant Sensor as Sensors (GMSL2 Cameras / LiDAR)
    participant Orin as Jetson AGX Orin (System 1: NPU/GPU)
    participant PCIe as PCIe Gen4 x8 (DMA Ring Buffer)
    participant Virtex as Xilinx Virtex EPU (System 2)
    participant Actuator as Hardware Actuator Interlock Relay

    Sensor->>Orin: Raw frame stream (GMSL2 DMA, 60 fps)
    Note over Orin: TensorRT / VLM: object detection, trajectory proposal U_adv
    Orin->>PCIe: Write fact descriptor to host DMA ring buffer
    PCIe->>Virtex: Dispatch descriptor to AXI4 mailbox registers (< 1.5 μs)
    Note over Virtex: EPU array: evaluate deontic rules, verify invariants in URAM (40–100 ns)
    alt Invariant satisfied (Valid Action)
        Virtex-->>Actuator: disp_accept = 1 (PWM / actuators enabled)
        Virtex->>PCIe: Write acknowledgment status to response buffer
        PCIe->>Orin: Host notification (MSI-X IRQ): action approved
    else Invariant violated or timeout (Safety Violation)
        Virtex->>Actuator: disp_refuse = 0 (Instantaneous hardware E-Stop in < 5 ns!)
        Virtex->>PCIe: Write diagnostic fault code (0x0E01 / 0x0409)
        PCIe->>Orin: Fault IRQ interrupt: transition to fail-safe mode
    end
```

#### 3.3.1. Joint Measurement Protocol and Popperian Falsification Criterion

1. **Zero-Copy Exchange via DMA Ring Buffers (Low-Latency Mailbox):**
   * Fact descriptors are mapped into pinned shared host memory on the Orin host. The XDMA controller on the FPGA streams descriptors directly into internal AXI4 mailbox registers within $T_{\text{dma}} < 1.0–1.5\ \mu\text{s}$.
   * Return telemetry channel: the FPGA writes verification status into a host-side response buffer and asserts an MSI-X interrupt (or updates a pollable completion flag).
2. **Hardware Safety Watchdog:**
   * The FPGA arbiter incorporates an independent hardware timer counting down an unalterable period $T_{\text{wdg}}$ (e.g., 10 to 20 ms).
   * Should System 1 (NPU/GPU/Linux) freeze, suffer memory bus starvation, or fail to transmit a valid fact descriptor within $T_{\text{wdg}}$, the FPGA watchdog asserts an emergency fault, drops `disp_accept = 0`, and de-energizes actuators without waiting for software recovery.
3. **End-to-End Admission Loop Latency ($T_{\mathrm{loop}}$):**
   The aggregate elapsed time from raw sensor frame arrival to assertion of the physical `disp_accept` signal at the FPGA pin is characterized by:

```math
T_{\mathrm{loop}} = T_{\mathrm{infer}} + T_{\mathrm{dma}} + T_{\mathrm{fpga}}
```

   where:
   * $T_{\mathrm{infer}}$ is the time required for predicate extraction and candidate action synthesis on the NPU/GPU ($10–30\ \text{ms}$);
   * $T_{\mathrm{dma}}$ is the descriptor transmission latency across PCIe Gen4 ($1.0–1.5\ \mu\text{s}$);
   * $T_{\mathrm{fpga}}$ is the formal verification latency across the FPGA core array ($40–100\ \text{ns}$).
4. **Falsification Gate Criterion:**
   In compliance with functional safety mandates (ISO 26262 ASIL D, IEC 61508 SIL 3), the hypothesis affirming the safety and viability of the heterogeneous architecture is **falsified** (rejected) if:

```math
\max(T_{\mathrm{loop}}) > T_{\mathrm{FTTI}}
```

   where $T_{\mathrm{FTTI}}$ is the Fault Tolerant Time Interval defined for the target physical platform,
   or if the probability of missing a hardware verdict within the watchdog deadline satisfies $P(T_{\mathrm{loop}} > T_{\mathrm{wdg}}) > 10^{-9}\ \text{failures per operational hour}$,
   or if a single occurrence of `disp_accept = 1` is recorded when an active defeater exists or a formal safety invariant is violated.

> [!NOTE]
> **Theoretical and Engineering Foundations of Heterogeneous Tandems in Monograph Chapters:**
> - [Chapter 18. Execution Infrastructure](ch18-execution-infrastructure.md) — Latency budgeting for cross-platform exchange via PCIe DMA and queue overflow prevention.
> - [Chapter 21. From Recommendation to Action: Authority Control and Safe Execution in Production Environments](ch21-from-recommendation-to-action.md) — Hardware admission gates, emergency relays, and direct actuator authority control.
> - [Chapter 27. Safety Case: GSN Argument Synthesis and Verification](ch27-safety-case-gsn-synthesis.md) and [Chapter 30. Safety and Cybersecurity Co-Engineering](ch30-safety-cybersecurity-co-engineering.md) — Structuring Goal Structuring Notation (GSN) safety cases to demonstrate compliance with ISO 26262 ASIL D.
> - [Chapter 39. Active Compliance Auditor: Popperian Falsification, Regulatory Compliance (ASPICE/ISO 26262/ISO 21434), and Autonomous Test Generation](ch39-active-compliance-auditor-and-popperian-testing.md) — Formal Popperian falsification frameworks and automated ASPICE V-model traceability across safety requirements and physical hardware tests.

This heterogeneous tandem resolves the structural conflict between the raw computational throughput of neural models and the strict determinism demanded by functional safety standards, transforming stochastic AI into a dependable instrument of cyber-physical control.

### 3.4. Autonomous Popperian Red Teamer Paradigm and Hardware Disposition Gate

Transitioning from academic software simulations of neuro-symbolic artificial intelligence to operational cyber-physical platforms necessitates uncompromising defense against Intellectual Property Leakage (*IP Leakage*). Utilizing commercial public cloud APIs to inspect proprietary engineering blueprints, avionics source code, or ASIC netlists under strict Non-Disclosure Agreements (NDAs) is unacceptable. Research and stress-testing require an entirely autonomous, air-gapped hardware environment (*Air-Gapped R&D*), wherein stochastic AI latencies and faults are arrested by an impenetrable hardware safety interlock.

To address this challenge, the author deployed an autonomous adversarial testing testbed based on a heterogeneous dual-channel topology:

1. **Heuristic Path (System 1 / Advisory Node):**
   * Implemented on an industrial embedded computer, the **Seeed Studio reServer Industrial J501**, equipped with an **NVIDIA Jetson AGX Orin 64GB** module (unified LPDDR5 memory, 275 TOPS tensor compute).
   * Executes local domain SLMs around the clock for Project Knowledge Reconnaissance (PKR), multimodal parsing of schematics and manuals (VLM), and grammar-guided constrained generation (GBNF).
2. **Deterministic Admission Gate (System 2 / Epistemic Processor):**
   * Implemented on a reconfigurable FPGA logic matrix (Xilinx Virtex UltraScale+) or a dedicated safety microcontroller.
   * Executes strictly deterministic microcode for byte-level citation provenance verification, interval bounding of physical parameters, and conflict detection across rule systems (ASPIC+).
3. **Hardware Emergency Interlock (Hardware Disposition Gate):**
   * The interconnect between System 1 and System 2 is established not merely across PCIe packet buses, but also via **isolated physical digital input/output (DI/DO) lines** and a safety-certified industrial **CAN FD** bus.
   * Upon detecting the slightest deviation or fallacy, the hardware fault line drives an interlock signal into the host controller input in $T_{\mathrm{disp}} < 100\,\text{ns}$, inhibiting hazardous actuation without operating system mediation.

This hardware foundation implements the methodology of the autonomous Popperian Red Teamer. The local System 1 tensor processor generates millions of adversarial counterexamples, false premises, and boundary conditions (Boundary Value Analysis, BVA). The System 2 hardware verifier mathematically inspects each query, validating the Zero Hallucination Rate invariant:

```math
\mathrm{ZHR} = 1 - \frac{N_{\mathrm{leak}}}{N_{\mathrm{adversarial}}} = 1.00
```

where:
- $\mathrm{ZHR}$ is the Zero Hallucination Rate metric;
- $N_{\mathrm{leak}}$ is the count of stochastic hallucinations or errors that penetrated the admission gate (governed by the engineering invariant $N_{\mathrm{leak}} \equiv 0$);
- $N_{\mathrm{adversarial}}$ is the total volume of generated adversarial counterexamples and boundary mutations ($N_{\mathrm{adversarial}} \ge 10^7$).

Any breakthrough of a hallucination past the symbolic filters is automatically logged as a critical defect of the knowledge base or linter.

#### 3.4.1. Actionable Closed Loop of Autonomous Knowledge Stress-Testing

1. **Control Flow & Runtime Decisions:**
   When the System 1 adversarial generator synthesizes an edge case triggering a `CONFLICT` or `DEFIEVER_ACTIVE` signal in the Datalog engine, the hardware disposition gate confirms interception by latching the `REDTEAM_CAUGHT` flag. If a stochastic hypothesis induces a processing stall $T > T_{\mathrm{wdg}}$, the hardware watchdog timer asserts `TEST_FAIL` and routes the mutation vector into the linter defect quarantine log.
2. **Hardware Dimensioning:**
   The Jetson Orin mutation generator delivers an adversarial descriptor fuzzing throughput of $R_{\mathrm{fuzz}} \ge 5\,000\,\text{packets/s}$. The FPGA hardware disposition gate clocked at $350\,\text{MHz}$ evaluates up to $10^7$ assertions per second, consuming less than 4% of the CAN FD bus bandwidth and guaranteeing a physical fault assertion latency of $T_{\mathrm{disp}} \le 85\,\text{ns}$.
3. **Worked Numerical Example:**
   During an overnight autonomous stress run, System 1 generated $N_{\mathrm{adversarial}} = 20,000,000$ boundary vectors of angular velocities and spoofed targets. The System 2 admission gate intercepted $19,999,998$ violations via deontic rules and 2 deviations due to watchdog timeouts ($T > 10\,\text{ms}$). The number of escapes was $N_{\mathrm{leak}} = 0$, yielding $\mathrm{ZHR} = 1.00$. The hardware disposition line latency measured $72\,\text{ns}$, confirming full compliance with ISO 26262 ASIL D.

### 3.5. Research Testbed "Generalized Synchronous Measurement Platform" Based on Infineon TLE90xx and AURIX TC49x

In expert systems designed for the diagnostics of cyber-physical complexes, autonomous platforms, and power electronics, the most acute vulnerability of the primary sensory layer stems from phase jitter and inter-channel temporal desynchronization. During sequential sensor polling by standard microcontroller ADCs, a time skew inevitably emerges between sampling instances of voltages, currents, and mechanical vibrations. This skew irreversibly corrupts spectral analysis, precluding precision early-stage detection of physical degradation in materials and load-bearing structures.

To overcome this limitation, the author deployed a Generalized Synchronous Measurement Platform research testbed based on a tandem of specialized **Infineon TLE90xx** sensing integrated circuits and an **Infineon AURIX TC49x** automotive safety microcontroller (a hexacore TriCore running at 1.6 GHz with three pairs of ISO 26262 ASIL D lockstep cores).

```mermaid
flowchart TD
    accTitle: Multi-Channel Synchronous Measurement Platform TLE90xx and AURIX TC49x
    accDescr: Architecture of synchronous sampling, isoUART 2.0 galvanic isolation, and translation into Datalog facts on the AURIX TC49x core.

    subgraph HardwareSensors["<b>Physical Object Under Observation</b>"]
        CELLS["Power Busbars / Battery Cells / Strain Gauges"]
    end

    subgraph FrontEnd["<b>Infineon TLE90xx (Synchronous Measurement Front-End)</b>"]
        ADC_ARRAY["16-bit Delta-Sigma ADCs<br/>Simultaneous strobe across all channels U_n+1 - U_n"]
        SEC_PATH["Secondary Measurement Path<br/>(Secondary Path for ASIL D)"]
        ISO_TX["Isolated isoUART 2.0 Transceiver<br/>(Galvanic isolation up to 2 Mbps)"]
        ADC_ARRAY --> ISO_TX
        SEC_PATH --> ISO_TX
    end

    subgraph HostMCU["<b>Infineon AURIX TC49x (Safety Lockstep MCU)</b>"]
        ISO_RX["isoUART Controller / DMA"]
        CALC_CORE["Digital Signal Processing Core<br/>Impedance Calculation Z(ω)"]
        DATALOG_ENGINE["Onboard Datalog Engine (Zephyr RTOS)<br/>Fact generation: busbar_degraded(ID, Z)"]
        ISO_RX --> CALC_CORE
        CALC_CORE --> DATALOG_ENGINE
    end

    subgraph Actions["<b>Admission Gate and Protection</b>"]
        INTERLOCK["Hardware Emergency Interlock Circuit (Fail-Closed)"]
        DATALOG_ENGINE --> INTERLOCK
    end

    CELLS --> ADC_ARRAY
    CELLS --> SEC_PATH
    ISO_TX ==>|"isoUART differential pair"| ISO_RX

    classDef phys fill:#e1f5fe,stroke:#0288d1,stroke-width:2px,color:#01579b;
    classDef fe fill:#fff3e0,stroke:#f57c00,stroke-width:2px,color:#e65100;
    classDef mcu fill:#e8f5e9,stroke:#2e7d32,stroke-width:2px,color:#1b5e20;
    classDef act fill:#fce4ec,stroke:#c2185b,stroke-width:2px,color:#880e4f;

    class CELLS phys;
    class ADC_ARRAY,SEC_PATH,ISO_TX fe;
    class ISO_RX,CALC_CORE,DATALOG_ENGINE mcu;
    class INTERLOCK act;
```

#### 3.5.1. Hardware Decomposition and Physical Condition Spectroscopy

The architectural approach reinterprets specialized battery management integrated circuits (BMS ICs) as a universal high-precision measurement front-end:

1. **Simultaneous Hardware Channel Strobing:**
   A multi-channel 16-bit delta-sigma ADC executes fully synchronous sampling across all differential inputs driven by a single hardware strobe. This eliminates phase drift between instantaneous voltage $U_k$ and current $I_k$ measurements.
2. **Galvanically Isolated isoUART 2.0 Interface:**
   Digitized telemetry frames are transmitted to the host microcontroller across a differential capacitive or transformer isolation barrier at up to 2 Mbps. This provides the measurement pipeline with robust immunity against high-voltage common-mode electromagnetic transients (up to 1500 V) typical of robotic high-power drives.
3. **Secondary Measurement Path:**
   To fulfill ISO 26262 ASIL D mandates, the TLE90xx die integrates an independent second ADC that redundantly monitors voltages and temperatures in parallel. Any discrepancy between primary and secondary measurements exceeding $\epsilon_{\mathrm{diag}} = 3\,\text{mV}$ triggers a hardware fault detection.
4. **Electrochemical Impedance Spectroscopy (EIS):**
   Employing a variable-frequency excitation probe signal $\omega \in [10\,\text{Hz}, 1\,\text{kHz}]$, the system computes the complex transfer impedance function of high-voltage busbars and contact junctions:

```math
Z(\omega) = \frac{V(\omega)}{I(\omega)} = R_{\mathrm{real}}(\omega) + j\,X_{\mathrm{imag}}(\omega)
```

where:
- $Z(\omega)$ is the complex impedance of the electrical joint or electrochemical cell at angular frequency $\omega$, in ohms ($\Omega$);
- $V(\omega)$ is the amplitude-phase spectrum of the AC voltage ripple across the contacts, in volts ($\text{V}$);
- $I(\omega)$ is the spectrum of the excitation probe current through the conductor, in amperes ($\text{A}$);
- $R_{\mathrm{real}}(\omega)$ is the real ohmic resistance of the conductor, indicating mechanical degradation, oxidation, and micro-cracking of terminals;
- $X_{\mathrm{imag}}(\omega)$ is the reactive impedance component, sensitive to structural shifts in internal capacitance and inductance.

#### 3.5.2. Actionable Closed Loop of Synchronous Spectroscopy TLE90xx/AURIX

1. **Control Flow & Runtime Decisions:**
   The AURIX TC49x core computes the impedance spectral vector in real time. If the real resistance of a power busbar $R_{\mathrm{real}}$ exceeds the calibrated degradation threshold $R_{\mathrm{degrade}} = R_0 \cdot (1 + \delta_{\mathrm{crit}})$, the Datalog engine asserts the symbolic fact `busbar_degraded(bus_id, r_current)`. A deontic prohibition rule immediately halts peak current draw across the traction motors, forcing the power inverter into a degraded fail-safe state with current clamped to 40% of nominal (*Fail-Closed*).
2. **Hardware Dimensioning:**
   The ADC sampling rate is set to $f_s = 10\,\text{kHz}$ per channel. Across 16 measurement channels, the raw telemetry bandwidth amounts to $16 \times 16\,\text{bits} \times 10\,\text{kHz} = 2.56\,\text{Mbps}$. Applying hardware differential delta compression compresses the stream to $1.28\,\text{Mbps}$, safely within the $2.0\,\text{Mbps}$ throughput limit of the isoUART 2.0 interface with over 35% bus bandwidth headroom. The DMA ring buffer is allocated in the AURIX core's local static DLMU memory ($4\,\text{KB}$), avoiding traffic contention on the global SRI crossbar bus.
3. **Worked Numerical Example:**
   Consider monitoring a traction inverter power busbar with a nominal resistance $R_0 = 0.250\,\text{m}\Omega$ ($250\,\mu\Omega$). The critical degradation threshold is specified as $\delta_{\mathrm{crit}} = 0.20$ ($+20\%$), yielding $R_{\mathrm{degrade}} = 0.300\,\text{m}\Omega$.
   At an interrogation frequency of $f = 100\,\text{Hz}$, the synchronous TLE90xx ADC measures an AC voltage ripple $V = 48.0\,\text{mV}$ at an alternating current of $I = 120.0\,\text{A}$ with a phase angle $\phi = 12.0^\circ$.
   Complex impedance:
   $|Z| = \frac{48.0 \cdot 10^{-3}}{120.0} = 0.400\,\text{m}\Omega$.
   $R_{\mathrm{real}} = |Z| \cdot \cos(12.0^\circ) = 0.400 \cdot 0.9781 = 0.391\,\text{m}\Omega$.
   $X_{\mathrm{imag}} = |Z| \cdot \sin(12.0^\circ) = 0.400 \cdot 0.2079 = 0.083\,\text{m}\Omega$.
   Because $R_{\mathrm{real}} = 0.391\,\text{m}\Omega > 0.300\,\text{m}\Omega$, the system identifies contact degradation exceeding nominal baseline by $+56.4\%$. Within $42\,\mu\text{s}$, the Datalog engine derives the failure predicate, and the AURIX controller throttles contactor current, averting thermal runaway of the power junction.

---

## 4. Simplex Architecture: Complex Controller Supervised by a Simple Controller

To reconcile a high-performance, unverifiable controller with strict safety guarantees, systems leverage the **Simplex architecture** (*Simplex architecture*). Lui Sha articulated its core principle as using simplicity to control complexity: a simple, formally verified controller supervises and backs up an advanced, unverified one [[5]](#src-5). The architecture comprises three primary functional blocks:

1. **Advanced Controller:** Executes on the NPU (for instance, within the Nav2 navigation stack for ROS 2 [[6]](#src-6)), employing deep neural networks and complex trajectory optimization algorithms.
2. **Baseline Controller:** Executes on the safety microcontroller, encapsulating elementary, deterministically provable control laws: deceleration, hover-in-place, or emergency landing.
3. **Admission Arbiter:** A symbolic expert engine that compares proposed commands from the advanced controller against a formal model of the robot's physical operating envelope in real time.

```mermaid
flowchart LR
    accTitle: Simplex Architecture for an Autonomous Robot
    accDescr: Switching between advanced and baseline controllers via an expert system admission arbiter.

    SENS["Sensory State<br/>(odometry, IMU, battery)"] --> ADV["<b>Advanced Controller (NPU)</b><br/>trajectory planner"]
    SENS --> BASE["<b>Baseline Controller (MCU)</b><br/>deterministic emergency stop"]
    SENS --> ARB{"<b>Admission Arbiter</b><br/>Does action violate safety invariants?"}

    ADV -->|"Candidate command U_adv"| ARB

    ARB -->|"Yes: invariant violated or fault"| SW["<b>Hardware Switch</b>"]
    ARB -->|"No: command permissible"| SW

    BASE -.->|"Safe command U_safe"| SW
    ADV -.->|"Nominal command U_adv"| SW

    SW --> ACT["Actuators and Motors"]

    classDef sens fill:#e1f5fe,stroke:#0288d1,stroke-width:2px,color:#01579b;
    classDef adv fill:#fff3e0,stroke:#e65100,stroke-width:2px,color:#bf360c;
    classDef safe fill:#e8f5e9,stroke:#2e7d32,stroke-width:2px,color:#1b5e20;
    classDef arb fill:#ede7f6,stroke:#512da8,stroke-width:2px,color:#311b92;
    classDef act fill:#fce4ec,stroke:#c2185b,stroke-width:2px,color:#880e4f;

    class SENS sens;
    class ADV adv;
    class BASE safe;
    class ARB,SW arb;
    class ACT act;
```

### 4.1. Formalization of the Safety Invariant

Let the physical state of the robot at time $t$ be represented by the state vector $X(t)$: spatial coordinates, linear velocity $v$, angular velocity $\omega$, and clearance distance to the nearest obstacle $d_{\mathrm{obs}}$. The neural trajectory planner proposes a control vector $U_{\mathrm{adv}}=(a,\alpha)$, where $a$ denotes proposed linear acceleration and $\alpha$ denotes angular acceleration. The arbiter authorizes $U_{\mathrm{adv}}$ if and only if the following safety invariant holds:

```math
d_{\mathrm{obs}}-S_{\mathrm{stop}}(v,a_{\mathrm{max}})\ge D_{\mathrm{margin}}+v\,\tau_{\mathrm{reaction}}.
```

Where:

- $d_{\mathrm{obs}}$ is the instantaneous distance from the robot to the nearest detected obstacle, in meters;
- $S_{\mathrm{stop}}(v,a_{\mathrm{max}})$ is the stopping distance—the clearance traversed while braking to a complete standstill from velocity $v$ under maximum permissible braking deceleration $a_{\mathrm{max}}$; assuming uniform deceleration, $S_{\mathrm{stop}}=v^{2}/(2a_{\mathrm{max}})$;
- $v$ is the linear velocity of the platform, in meters per second;
- $a_{\mathrm{max}}$ is the maximum guaranteed braking deceleration achievable on the current physical surface, in meters per second squared;
- $D_{\mathrm{margin}}$ is the mandatory safety margin that must remain between the halted robot and the obstacle, e.g., 0.5 m;
- $\tau_{\mathrm{reaction}}$ is the worst-case hardware actuation latency from hazard detection to physical brake engagement, in seconds, measured directly on target hardware;
- $v\,\tau_{\mathrm{reaction}}$ is the distance traversed during this reaction latency before braking deceleration begins.

The invariant condition reads: the clearance to the obstacle, minus the stopping distance, must cover both the safety margin and the distance traversed during actuation delay. Suppose, for example, $v=3$ m/s, $a_{\mathrm{max}}=2$ m/s², and $\tau_{\mathrm{reaction}}=0.02$ s. Then $S_{\mathrm{stop}}=3^{2}/(2\cdot2)=2.25$ m, $v\,\tau_{\mathrm{reaction}}=0.06$ m, and the arbiter permits forward motion only if the obstacle is farther than $2.25+0.5+0.06=2.81$ m. These parameters serve an illustrative purpose. If $U_{\mathrm{adv}}$ violates this inequality, or if the NPU fails to publish the next command within a watchdog deadline (e.g., 50 ms), the arbiter deasserts the $U_{\mathrm{adv}}$ channel and passes control to the baseline controller, which initiates maximum braking at $-a_{\mathrm{max}}$.

## 5. Onboard Knowledge Base: Compiled Datalog Without Dynamic Memory Allocation

Embedded microcontrollers operate under strict memory limits ranging from hundreds of kilobytes to several megabytes of SRAM. Running interpreted Prolog engines or rule processors that allocate dynamic objects on the heap (`malloc`, `new`) within safety-critical control loops is dangerous: heap fragmentation and unpredictable garbage collection pauses violate real-time guarantees. Consequently, knowledge base rules are compiled ahead-of-time into compact C structures, bitmasks, and static transition tables.

Datalog is an optimal formalism: it represents a subset of logic programming devoid of complex functional symbols, guaranteeing that inference evaluation always terminates [[7]](#src-7). Consider the following rule fragment evaluating UAV mission flight safety:

```prolog
% Onboard flight safety rules
hazard(critical_battery) :-
    telemetry(battery_voltage, V), V < 21.0.

hazard(geofence_breach) :-
    position(_Alt, Dist), Dist > 5000.

hazard(sensor_blindness) :-
    sensor_health(lidar, failed),
    sensor_health(optical_flow, degraded).

% Fail-safe mode transition decisions
action(emergency_landing) :-
    hazard(critical_battery).

action(return_to_home) :-
    hazard(geofence_breach),
    not hazard(critical_battery).

action(hold_position) :-
    hazard(sensor_blindness),
    not hazard(critical_battery).
```

### 5.1. Code Generation for Zephyr RTOS

The knowledge compiler translates these rules into a standalone C evaluation function completely devoid of dynamic memory allocation. The listing below illustrates the generated output; to compile outside Zephyr RTOS, one simply removes the `#include <zephyr/kernel.h>` directive.

<details>
<summary>C Code Example: Generated Onboard Safety Rule Evaluation Function for Zephyr RTOS</summary>

```c
/* Generated by the knowledge compiler for Zephyr RTOS */
#include <zephyr/kernel.h>
#include <stdint.h>
#include <stdbool.h>

typedef struct {
    float battery_voltage;
    float distance_from_home;
    float altitude;
    uint8_t lidar_status;       /* 0 = OK, 1 = Degraded, 2 = Failed */
    uint8_t opt_flow_status;
} RobotTelemetry_t;

typedef enum {
    ACTION_CONTINUE = 0,
    ACTION_HOLD_POSITION,
    ACTION_RETURN_TO_HOME,
    ACTION_EMERGENCY_LANDING
} SafetyVerdict_t;

SafetyVerdict_t evaluate_safety_rules(const RobotTelemetry_t *const telem) {
    const bool critical_battery = (telem->battery_voltage < 21.0f);
    if (critical_battery) {
        return ACTION_EMERGENCY_LANDING; /* highest priority */
    }

    const bool geofence_breach = (telem->distance_from_home > 5000.0f);
    if (geofence_breach) {
        return ACTION_RETURN_TO_HOME;
    }

    const bool sensor_blind = (telem->lidar_status == 2) &&
                              (telem->opt_flow_status >= 1);
    if (sensor_blind) {
        return ACTION_HOLD_POSITION;
    }

    return ACTION_CONTINUE;
}
```

</details>

Comparing the Datalog rules with the compiled C code reveals two fundamental architectural characteristics. First, the generated function performs zero memory allocation, contains no unbounded loops, and avoids recursion; its execution time is strictly bounded, with upper bounds verifiable via Worst-Case Execution Time (*worst-case execution time*, WCET) static analysis or target hardware profiling. Second, whereas raw Datalog evaluates all matching rules concurrently, the generated code resolves rule competition through explicit, deterministic prioritization: when a geofence breach coincides with sensor blindness, return-to-home overrides hover-in-place. This priority hierarchy constitutes a formal facet of the knowledge base that must be audited and verified, rather than left to implicit compiler ordering.

## 6. Action Contracts and Transactions in the Physical World

As demonstrated in [Chapter 21](ch21-from-recommendation-to-action.md), executing an action in the physical world requires closed-loop verification: an autonomous system cannot simply write command frames to a CAN bus and assume the task completed successfully. Complex physical commands must be structured as formal contracts specifying preconditions, execution monitoring, postconditions, and compensation procedures. This paradigm extends the database concept of sagas, where long-running transactions are composed of discrete steps, each linked to a dedicated compensating action [[8]](#src-8).

```mermaid
sequenceDiagram
    autonumber
    accTitle: Action Contract for Manipulator Joint Rotation
    accDescr: Arbiter validates preconditions, initiates actuation, monitors dynamic invariants via independent sensors, and verifies postconditions.

    participant High as ROS 2 Action Client on NPU
    participant Expert as Admission Arbiter on MCU
    participant Act as Actuator on CAN Bus
    participant Sens as Independent Sensor (Encoder, IMU)

    High->>Expert: Action request: rotate joint by 45°, torque limit 12 N·m

    Note over Expert: 1. Preconditions:<br/>battery charge > 15%<br/>workspace clear<br/>safety interlock switch closed

    alt Precondition violated
        Expert-->>High: Refusal: safety interlock switch open
    else Preconditions satisfied
        Expert->>Act: Command: angle 45°, current limit 12 A
        Expert->>Expert: Arm watchdog timer (2000 ms)

        loop Verification every 20 ms
            Sens-->>Expert: Angle and angular velocity
            Note over Expert: Dynamic invariant:<br/>angular velocity within allowable bounds
        end

        Sens-->>Expert: Final state: angle 45.1°, no stall detected
        Note over Expert: 2. Postcondition:<br/>|angle − target| ≤ 0.5°

        Expert-->>High: Action completed, receipt with SHA-256 digest
    end
```

### 6.1. Action Contract Structure in C

<details>
<summary>C Code Example: Action Contract Data Structure</summary>

```c
struct ActionContract {
    uint32_t action_id;
    uint32_t timeout_ms;

    /* Preconditions: returns true if hardware is ready */
    bool (*check_preconditions)(void);

    /* Execution: write to registers or CAN bus */
    int (*execute_command)(void *params);

    /* Postconditions: independent sensor verification */
    bool (*verify_postconditions)(void *params);

    /* Compensation: safe fallback upon contract breach */
    void (*compensate_failure)(void);
};
```

</details>

Consider an operational scenario where, during joint rotation, an independent current sensor detects an overcurrent draw exceeding 15 A while the motor encoder registers zero angular displacement—indicating a mechanical transmission lockup. The contract is breached: `compensate_failure()` de-energizes the motor windings, engages the mechanical holding brake, and transmits a high-priority diagnostic alert to the mission controller. Both the cryptographic receipt of successful actions and failure audit trails are logged immutably to the onboard evidence ledger.

## 7. Autonomy Under Loss of Communication and in Electronic Warfare Environments

In field robotics and defense deployments ([Military Expert Systems](../../../MilTech/Military-Expert-Systems-UAS-AD-ELINT-EW-UA.md)), communication channels to base stations are routinely jammed, intercepted, or completely severed. The onboard control architecture must therefore implement a deterministic degradation state machine aligned with [Chapter 22](ch22-cybernetics-edge-to-backend.md).

```mermaid
stateDiagram-v2
    accTitle: Robot Autonomy State Machine Under Communication Loss
    accDescr: Transitions between normal connected mode, local degraded mode, autonomous safe hold, and post-link reconciliation.

    [*] --> ConnectedNormal: Radio link nominal

    ConnectedNormal --> DegradedLocalOnly: Telemetry loss exceeding 500 ms
    note right of DegradedLocalOnly
        • Actions with irreversible consequences prohibited
        • Inertial odometry engaged
        • Velocity clamp: 50%
    end note

    DegradedLocalOnly --> AutonomousSafeHold: Local autonomy timer expired
    note right of AutonomousSafeHold
        • Platform standstill or hover
        • Awaiting authorized beacon
        • 360° perimeter monitoring
    end note

    DegradedLocalOnly --> ConnectedNormal: Link restored, handshake successful

    AutonomousSafeHold --> Reconciliation: Authorized base beacon acquired
    Reconciliation --> ConnectedNormal: Onboard audit log successfully verified
```

The state machine defines four operational states, each successively curtailing platform authority. Thresholds such as the 500 ms link-loss timeout and the 50% velocity clamp are representative engineering parameters calibrated through rigorous platform safety analysis.

### 7.1. Behavioral Rules in the Absence of Reliable GNSS

When operating in Global Navigation Satellite System (*Global Navigation Satellite System*, GNSS) spoofing environments, the onboard expert system continuously monitors for divergence across independent odometry sources. The comprehensive survey by Psiaki and Humphreys details GNSS spoofing mechanics and countermeasure verification techniques, particularly cross-validation against inertial navigation sensors [[9]](#src-9).

**Odometry Consistency Invariant.** The onboard expert system evaluates whether the divergence between independent velocity estimates exceeds a critical threshold:

```math
\left\lVert\mathbf{v}_{\mathrm{GNSS}}-\mathbf{v}_{\mathrm{Wheel}/\mathrm{IMU}}\right\rVert_2>\Delta V_{\mathrm{threshold}}.
```

Where:

- $\mathbf{v}_{\mathrm{GNSS}}$ is the velocity vector reported by the GNSS receiver, in meters per second;
- $\mathbf{v}_{\mathrm{Wheel}/\mathrm{IMU}}$ is the velocity vector obtained through integrating wheel encoders and Inertial Measurement Unit (*inertial measurement unit*, IMU) telemetry;
- $\lVert\cdot\rVert_2$ is the standard Euclidean norm representing vector magnitude, quantifying the degree of divergence between independent sensor modalities;
- $\Delta V_{\mathrm{threshold}}$ is the discrepancy threshold, calibrated strictly above the combined sensor noise floor measured on the physical platform to avoid nuisance alarms.

If the GNSS receiver indicates a velocity of 80 km/h while integrated wheel odometry and accelerometer readings indicate 10 km/h, the satellite navigation channel is flagged as compromised (`Status = SPOOFED`).

**Safe Reaction.** Compromised GNSS coordinates are stripped from the active Kalman filter [[10]](#src-10), the platform transitions to dead-reckoning inertial-optical navigation, and the last verified origin waypoint is preserved in non-volatile memory. Denied-environment navigation architectures are explored in depth in [Appendix C](appendix-c-autonomous-navigation-and-geosearch.md).

## 8. Evidentiary Layer of Sensory Measurements: Feature Signing and Sensor Frontend Protection

### 8.1. Signing Features Instead of Verdicts and Hardware Root of Trust

Distributed sensor networks (acoustic direction finders, RF scanners, seismic arrays) face a subtle architectural vulnerability: **a cryptographically signed packet proves the authenticity of the sender, but does not prove the truth of its verdict**. Sidliarchuk articulates this principle in the witness-integrity architecture: digital signatures prove authorship, not empirical truth [[11]](#src-11).

This distinction exposes two critical failure modes. If a captured or malfunctioning sensor transmits a signed status claiming "target detected," the receiving system can only accept or reject the packet; it cannot audit the underlying evidence. Conversely, if a sensor generates an error induced by physical environment artifacts, the signature remains mathematically valid. During field trials documented by the repository author, an acoustic detector locked onto a stable 84 Hz fundamental frequency with harmonics and emitted a positive target verdict, although the source was merely an overhead high-voltage power line located hundreds of meters away [[11]](#src-11).

```mermaid
flowchart LR
    accTitle: Comparison of Verdict Signing Versus Feature Signing
    accDescr: A signed verdict cannot be audited; signed features enable independent verdict recomputation on the receiver side.

    subgraph Naive["<b>Verdict Signing: Blind Trust</b>"]
        S1["Sensor"] -->|"Signed 'target detected' verdict"| R1["Receiver"]
        R1 -->|"Signature valid, execute action"| D1["False alarm or hazardous action"]
    end

    subgraph Evidence["<b>Feature Signing: Evidentiary Layer</b>"]
        S2["Sensor with Secure Element"] -->|"Signed feature payload"| R2["Receiver Verification"]
        R2 -->|"1. Verify signature"| H1{"Signature valid?"}
        H1 -->|"Yes"| H2{"Recomputed verdict matches?"}
        H2 -->|"Yes"| ADM["Admitted to reasoning pipeline"]
        H2 -->|"No"| REJ["Refusal: evidence conflict"]
    end

    classDef bad fill:#ffebee,stroke:#c62828,stroke-width:2px,color:#b71c1c;
    classDef good fill:#e8f5e9,stroke:#2e7d32,stroke-width:2px,color:#1b5e20;
    classDef proc fill:#e1f5fe,stroke:#0288d1,stroke-width:2px,color:#01579b;

    class D1,REJ bad;
    class ADM good;
    class S1,S2,R1,R2,H1,H2 proc;
```

The witness-integrity implementation codifies an evidentiary sensor contract built upon four principles [[11]](#src-11):

1. **Sensors Sign Features, Not Verdicts:** Feature payloads adhere to a canonical binary layout with fixed field ordering and integer-only representations. Floating-point serialization discrepancies across heterogeneous hardware architectures inevitably break cryptographic signature matching.
2. **Private Keys Never Leave Secure Silicon:** Cryptographic signatures are minted directly within a dedicated secure hardware element (such as an ATECC608B), wherein private keys remain permanently shielded from host extraction.
3. **Physically Impossible Payloads Are Rejected at the Source:** Consistency checks execute prior to signature generation: a frame claiming more confirmed detections than processed time samples is physically invalid; the secure element aborts signing rather than endorsing corrupt data.
4. **Receivers Recompute Verdicts Independently:** The receiving controller verifies the signature and recomputes the logical verdict against identical evidentiary rules; any discrepancy between the claimed status and recomputed verdict immediately collapses epistemic trust.

For an onboard expert system, this design ensures that facts supplied by external sensors enter the reasoning base only when accompanied by signed feature records and successful independent recomputation. Hardware roots of trust and anti-spoofing architectures are analyzed in detail in the author's dedicated treatise [Hardware Root of Trust and Anti-Spoofing in Military IoT Systems](../../../MilTech/Hardware-Root-Of-Trust-And-Anti-Spoofing-Military-IoT-UA.md).

### 8.2. Hardware and Software Protection of Sensor Front-End Against Fail-Open Anomalies: Filtration of Jagged Matrices, Floating Noise, and Sensor Failures

In automotive and aerospace cyber-physical systems requiring the highest integrity levels (ISO 26262 ASIL D / DO-178C DAL A), the epistemic core receives raw observation vectors from sensor matrices (specifically TLE90xx, inductive rotor position sensors, radar front-ends, and CAN FD buses). Under real operational environments, physical transmission channels are vulnerable to degradation and signal anomalies capable of triggering a catastrophic silent admission failure (*Fail-Open*):

1. **Fatal Fail-Open on Floating-Point Anomalies (`NaN` / `Inf` Noise Floor):**
   If the sensory input pipeline yields non-numeric IEEE 754 anomalies (`math.NaN` or `math.Inf`) instead of valid electrical telemetry, standard programming language comparison operators (across Go, C++, Verilog) exhibit non-trivial semantics:

```math
(\mathrm{NaN} \le \mathrm{Threshold}) \equiv \mathrm{False}, \qquad (\mathrm{NaN} > \mathrm{Threshold}) \equiv \mathrm{False}
```

   If defensive validation logic is naively structured (for example, checking merely `if noise > max_allowed { reject() }`), the presence of `NaN` bypasses the rejection condition. The system mistakenly treats corrupted noise as acceptable and issues an unsafe admission decision (`ACCEPT`), permitting autonomous platform motion with a failed sensor.

2. **Core Panic on Asynchronous Non-Rectangular ("Jagged") Matrices:**
   During streamed sensor reads across Direct Memory Access (DMA), hardware packet drops can truncate individual rows of the observation matrix (*Jagged Matrix*). An algorithm's subsequent attempt to access a fixed column index $M[i][j]$ triggers an unhandled runtime fault (`panic: index out of range`), instantly freezing the control loop.

3. **Dead-Sensor Blindness:**
   A complete line severance or ADC lockup at zero ground potential causes the measurement vector to collapse to all zeros ($V = [0.0, 0.0, \dots, 0.0]$). If zero falls within the sensor's nominal operational range, the logical reasoning core naively interprets this catastrophic open circuit as a quiescent, perfectly stable platform state.

To neutralize these vulnerabilities at the interface between low-level device drivers and hardware accelerators, a **sensor-frontline fail-closed gate** (*Sensor-Frontline Fail-Closed Gate*) is integrated into the data ingestion pipeline.

```mermaid
flowchart TD
    accTitle: Sensor-Frontline Fail-Closed Gate
    accDescr: Four-tier deterministic filtration of raw sensor streams prior to admission into the logical reasoning core.

    RAW["<b>Raw Sensor / ADC Stream</b><br/>Observation DMA Buffers"] --> G1

    subgraph Gate1["<b>1. Geometric Rectangularity Check</b>"]
        G1{"∀ i, len(Row[i]) == ExpectedCols?"}
    end
    G1 -->|"No"| TRAP1["<b>HALT_REFUSE (0x422)</b><br/>Hardware interlock on Jagged Matrix"]
    G1 -->|"Yes"| G2

    subgraph Gate2["<b>2. Physical Bounds and IEEE 754 Validation</b>"]
        G2{"!IsNaN(v) && !IsInf(v)<br/>&& v >= NoiseFloor_min?"}
    end
    G2 -->|"No"| TRAP2["<b>HALT_REFUSE (0x400)</b><br/>Trap NaN / Inf anomaly"]
    G2 -->|"Yes"| G3

    subgraph Gate3["<b>3. Dead Sensor Detector (L2-norm / variance)</b>"]
        G3{"σ²(V) >= ε_dead?"}
    end
    G3 -->|"No"| TRAP3["<b>HALT_REFUSE (0x404)</b><br/>Sensor open circuit or lockup"]
    G3 -->|"Yes"| G4

    subgraph Gate4["<b>4. Epistemic Datalog / EVM Core</b>"]
        G4["<b>Safe Deterministic Inference</b><br/>Evidentiary fact generation"]
    end

    classDef raw fill:#e1f5fe,stroke:#0288d1,stroke-width:2px,color:#01579b;
    classDef chk fill:#fff3e0,stroke:#f57c00,stroke-width:2px,color:#e65100;
    classDef trap fill:#ffebee,stroke:#c62828,stroke-width:2px,color:#b71c1c;
    classDef ok fill:#e8f5e9,stroke:#2e7d32,stroke-width:2px,color:#1b5e20;

    class RAW raw;
    class G1,G2,G3 chk;
    class TRAP1,TRAP2,TRAP3 trap;
    class G4 ok;
```

#### 8.2.1. Algorithm for Deterministic Structural Validation and Noise Filtering

Before executing any numerical computation, the adapter must enforce a strict four-stage pipeline:

1. **Geometric Invariance:** Every row length is compared against the expected dimension in $O(R)$ time. Any discrepancy triggers immediate rejection prior to indexed element dereferencing.
2. **IEEE 754 Anomaly Filter:**

```math
\forall x \in M: \quad (\mathrm{math.IsNaN}(x) \lor \mathrm{math.IsInf}(x, 0)) \implies \mathrm{TRAP\_INVALID\_ARITHMETIC}
```

3. **Activity Energy Gate:**
   To prevent blindness to disconnected or locked-up sensors, the signal variance is evaluated across the observation window:

```math
\sigma^2(V) = \frac{1}{N}\sum_{i=1}^N (v_i - \bar{v})^2
```

   If $\sigma^2(V) < \epsilon_{\mathrm{dead}}$, the protective status "Sensor de-energized or locked" is triggered, inhibiting transition to the armed operational state.
4. **Electromagnetic Noise Floor Bounds (EMC Floor Clamp):**
   The thermal and electromagnetic noise floor of the sensor is bounded bilaterally:

```math
\mathrm{NoiseFloor}_{\mathrm{min}} \le \eta \le \mathrm{NoiseFloor}_{\mathrm{max}}
```

   If observed noise collapses to absolute zero (which is physically impossible in an operational analog circuit subject to Johnson-Nyquist thermal noise $v_n = \sqrt{4 k_B T R \Delta f}$), this signals an artificial replay simulation (spoofing) or physical disconnection of the analog front-end.

#### 8.2.2. Actionable Closed Loop of Sensor Front-End Protection

1. **Control Flow & Runtime Decisions:**
   Upon detecting any structural or arithmetic anomaly, the sensor admission gate halts fact synthesis, sets operational status to `HALT_REFUSE`, asserts the hardware flag `FLAG_SENSOR_FAULT`, and commands the actuator subsystem into safe degraded mode. No element of an anomalous frame is ever passed to the inference engine.
2. **Hardware Dimensioning:**
   Array structure validation is executed directly in the controller's DMA pipeline within $T_{\mathrm{val}} < 15\,\mu\text{s}$ for matrices up to $64 \times 64$. The lower thermal noise limit $\mathrm{NoiseFloor}_{\mathrm{min}}$ is calculated from analog physical parameters: for a resistive divider $R = 10\,\text{k}\Omega$ at $T = 300\,\text{K}$ and bandwidth $\Delta f = 10\,\text{kHz}$, the minimum thermal noise voltage is $v_n \approx 1.28\,\mu\text{V}$.
3. **Worked Numerical Example:**
   A $16 \times 32$ tactile pressure sensor matrix streams an observation buffer. Geometric verification validates a row width of 32 for all 16 rows ($O(R) = 16$ checks, $0.4\,\mu\text{s}$).
   Evaluating array elements reveals cell $M[4][12]$ contains $\mathrm{0x7FF8000000000000}$ (quiet `NaN`), induced by a transient brownout on the SPI transceiver power rail.
   A naive comparator would pass this corrupted sample, since $(\mathrm{NaN} > 50.0) \equiv \mathrm{False}$.
   The specialized filter traps `math.IsNaN(x) == True` within $12\,\text{ns}$, aborts ingestion, asserts error code `0x400`, and halts platform maneuvers, preventing uncontrolled acceleration.

## 9. Non-Functional Interfaces and Assumption Boundaries Prior to Integration

Complex cyber-physical systems frequently exhibit a well-known systems engineering paradox: two fully conforming subsystems combine into a non-functional system. Component A (e.g., a thermal camera or RF transceiver) and Component B (e.g., a GNSS receiver or flight controller) may independently meet specifications and pass unit test suites, yet fail catastrophically upon physical integration. Consequently, the NASA Systems Engineering Handbook designates interface management as an autonomous, core lifecycle process [[12]](#src-12).

Such failures predominantly stem from unmodeled non-functional interfaces:

- **Shared Power Buses and Ground Loops:** High-current motor transient spikes induce voltage drops across shared digital sensor power rails;
- **Electromagnetic and RF Interference (EMI/RFI):** High-speed digital camera serialization lanes radiate EMI that degrades reception at adjacent GPS/GNSS patch antennas, mandating integrated physical shielding and spatial isolation audits;
- **Thermal and Vibration Coupling:** Heat dissipated by high-power inverter stages alters the temperature calibration curves of adjacent IMU gyroscopes.

### 9.1. Four Steps for Verifying Engineering Assumptions

Prior to embarking on expensive field trials, the system engineering design process must enforce completing an assumption verification matrix across every subsystem boundary:

```mermaid
flowchart LR
    accTitle: Four Steps for Verifying Subsystem Interface Assumptions
    accDescr: Progression from assumption formulation through disturbance injection and expected reaction to architecture revision criteria.

    A["<b>1. Assumption</b><br/>Which environmental parameters does the module presume stable?"] --> B["<b>2. Disturbance</b><br/>Injected latency, noise, 20% voltage drop"]
    B --> C["<b>3. Expected Reaction</b><br/>Transition to deterministic fail-safe mode"]
    C --> D["<b>4. Revision Criterion</b><br/>Under what outcome is the architectural decision revoked?"]

    classDef step fill:#e1f5fe,stroke:#0288d1,stroke-width:2px,color:#01579b;
    class A,B,C,D step;
```

Special emphasis is placed on characterizing the "gray zone" of partial hardware degradation. Subsystems must be stressed not merely against clean open-circuit faults, but against progressive calibration drift across thermal cycles, burst bus delays, and clock oscillator instability.

## 10. Physical Modeling: Vibration, Heat, and Mesh Convergence

Software reliability is immaterial if an onboard compute unit overheats inside a sealed chassis or if IMU sensors drown in structural vibration harmonics. Physical hardware configurations must be validated via Computer-Aided Engineering (*computer-aided engineering*, CAE) simulation combined with physical test bench characterization. Oberkampf and Roy distinguish two core disciplines: verification ensures the numerical model equations are solved correctly, while validation establishes whether the chosen equations correctly represent physical reality [[13]](#src-13).

```mermaid
flowchart TD
    accTitle: Triad of Physical Validation for Onboard Computers
    accDescr: Vibration modal analysis, conjugate heat transfer, and mesh convergence proceed through verification and validation to flight certification.

    subgraph CAE["<b>Computer-Aided Engineering (CAE) Environment</b>"]
        M1["<b>1. Modal Analysis</b><br/>natural frequencies vs. motor rotor harmonics"]
        M2["<b>2. Conjugate Heat Transfer</b><br/>NPU thermal dissipation in sealed enclosures"]
        M3["<b>3. Mesh Convergence</b><br/>eliminating artificial stress singularities"]
    end

    subgraph VV["<b>Verification & Validation</b>"]
        VER["<b>Verification</b><br/>Are the equations solved correctly?"]
        VAL["<b>Validation</b><br/>Are the correct equations being solved?"]
    end

    CAE --> VER
    VER -->|"Mathematical convergence"| VAL
    VAL -->|"Shaker table & thermal chamber tests"| PROD["<b>Flight Certification</b>"]

    classDef cae fill:#e1f5fe,stroke:#0288d1,stroke-width:2px,color:#01579b;
    classDef vv fill:#fff3e0,stroke:#e65100,stroke-width:2px,color:#bf360c;
    classDef prod fill:#e8f5e9,stroke:#2e7d32,stroke-width:2px,color:#1b5e20;

    class M1,M2,M3 cae;
    class VER,VAL vv;
    class PROD prod;
```

**Modal Analysis.** Drone propulsion systems induce severe mechanical vibrations at blade pass frequencies:

```math
f=\frac{\text{RPM}\cdot N_{\text{blades}}}{60}.
```

Where:

- $f$ is the blade pass frequency at which rotor blades pass a fixed airframe reference, in hertz (Hz);
- $\text{RPM}$ (*revolutions per minute*) is the motor shaft rotational speed, in revolutions per minute;
- $N_{\text{blades}}$ is the number of propeller blades;
- The constant 60 normalizes minutes to seconds ($1\ \text{Hz} = 1\ \text{cycle/s}$).

For example, a brushless motor spinning at 6000 RPM with a 2-blade propeller excites a fundamental vibration frequency of $f=6000\cdot2/60=200$ Hz, with upper harmonics at 400 Hz and 600 Hz. If the natural resonant frequency of the structural frame or computer mounting bracket coincides with these frequencies, resonance ensues. Induced vibrations corrupt accelerometer and gyroscope measurements, leading Kalman filters to interpret pseudo-accelerations as physical translation. Cyclic fatigue stresses introduce microcracks in BGA solder balls and loosen wiring harness terminals. Finite element modal analysis allows structural engineers to tune rib stiffeners and elastomer dampers, ensuring airframe resonant modes lie entirely outside normal motor operating RPM ranges.

**Conjugate Heat Transfer (*conjugate heat transfer*, CHT).** In sealed IP67 enclosures devoid of forced-air cooling, thermal energy from the NPU must conduct through thermal interface materials (TIM) into the chassis and dissipate via natural convection. Thermal dissipation budgets must be profiled under target computational workloads. Mismatched coefficients of thermal expansion (CTE) between PCB laminates and aluminum mounting lugs induce thermo-mechanical fatigue during rapid temperature transitions (e.g., $+40^\circ\text{C}$ ground launch transitioning to $-20^\circ\text{C}$ high-altitude cruise). CHT simulation must verify junction temperatures remain safely below thermal throttling thresholds to prevent frame drops in detection pipelines.

**Mesh Convergence.** Finite element analysis (FEA) results cannot be trusted without rigorous mesh convergence verification. Sharp internal fillet radii on structural brackets inevitably trigger mathematical stress singularities: refining the mesh produces unbounded stress spikes, rendering colorized stress contours numerical artifacts rather than physical reality. Roache established the standard for reporting grid refinement studies via the Grid Convergence Index (GCI) [[14]](#src-14). Permissible variations under mesh doubling must be established prior to computation, and localized stresses must be cross-checked against material plasticity limits.

## 11. Step-by-Step Prototype Deployment

A functional prototype of an evidence-governed robotic expert system can be constructed across four structured engineering steps.

### 11.1. Step 1: Development Environment

Initialize the Zephyr RTOS cross-compilation toolchain following official project documentation [[2]](#src-2). Under Linux, execute the following foundational sequence:

```bash
python3 -m venv ~/zephyrproject/.venv
source ~/zephyrproject/.venv/bin/activate
pip install west
west init ~/zephyrproject
cd ~/zephyrproject && west update
west zephyr-export
```

Subsequently, install required Python dependencies and the appropriate Zephyr SDK toolchain following the active release guidelines, as dependency requirements evolve across Zephyr versions.

### 11.2. Step 2: Safety Rules

Create the declarative rule file `safety_rules.dl` in your project root:

```prolog
% Input predicates: distance(SensorID, Centimeters)
unsafe_zone(SensorID) :- distance(SensorID, D), D < 30.
emergency_stop :- unsafe_zone(_).
```

### 11.3. Step 3: Supervisor Thread on Zephyr RTOS

Implement a supervisor thread executing at a fixed 10 ms period. The snippet below assumes the board device tree defines a `motor_switch` GPIO node, and that `read_sonar_distance()` is implemented by the ultrasonic rangefinder driver:

<details>
<summary>C Code Example: Supervisor Thread on Zephyr RTOS</summary>

```c
#include <zephyr/kernel.h>
#include <zephyr/drivers/gpio.h>

#define SUPERVISOR_PERIOD_MS 10

static const struct gpio_dt_spec motor_en =
    GPIO_DT_SPEC_GET(DT_NODELABEL(motor_switch), gpios);

uint16_t read_sonar_distance(void); /* implemented by rangefinder driver */

void safety_supervisor_thread(void *arg1, void *arg2, void *arg3) {
    gpio_pin_configure_dt(&motor_en, GPIO_OUTPUT_ACTIVE);

    while (1) {
        uint16_t front_distance = read_sonar_distance();

        if (front_distance < 30) {
            gpio_pin_set_dt(&motor_en, 0); /* de-energize actuator drivers */
            printk("SAFETY INTERLOCK: obstacle at %d cm, motors disabled\n",
                   front_distance);
        } else {
            gpio_pin_set_dt(&motor_en, 1);
        }

        k_msleep(SUPERVISOR_PERIOD_MS);
    }
}

K_THREAD_DEFINE(safety_thread_id, 1024, safety_supervisor_thread,
                NULL, NULL, NULL, 1, 0, 0);
```

</details>

During each cycle, the thread polls rangefinder telemetry and evaluates the compiled `unsafe_zone` invariant: if an obstacle is within 30 cm, motor power is cut immediately. In production builds, `k_msleep` is replaced with hardware kernel timers to eliminate scheduling jitter.

### 11.4. Step 4: Hardware-in-the-Loop (HIL) Simulation

Prior to flashing physical microcontrollers, validate the system via Hardware-in-the-Loop (*hardware-in-the-loop*, HIL) simulation:

1. Launch a physics simulator (such as Gazebo Sim or Isaac Sim) on the engineering workstation.
2. Interface the physical microcontroller board to the simulated virtual CAN bus (`vcan0`) via a USB-CAN adapter.
3. Inject synthetic noise and corrupted commands into the `/cmd_vel` ROS topic, verifying that the microcontroller intercepts authority within the allotted safety deadline.

## 12. Pre-Deployment Safety Checklist for Field Trials

Before deploying an autonomous cyber-physical platform to the proving ground, complete the following formal safety audit checklist:

| № | Safety Verification Item | Verification Mechanism | Status |
| :-: | :--- | :--- | :-: |
| 1 | **Physical Actuator Isolation** | Hardware emergency stop switch physically opens battery supply lines, bypassing microcontroller firmware | [ ] |
| 2 | **Hardware Safety Watchdog** | Independent microcontroller watchdog timer enabled with timeout budget derived from hazard analysis | [ ] |
| 3 | **Zero Dynamic Memory Allocation** | All safety supervision code audited: zero calls to `malloc()`, `free()`, no dynamic structures, no recursion | [ ] |
| 4 | **Velocity Clamping Limits** | Maximum angular and linear velocities clamped via hard saturation limits at the motor driver level | [ ] |
| 5 | **Closed-Loop Feedback Verification** | Every command contract specifies a postcondition timeout polling independent encoder or current telemetry | [ ] |
| 6 | **Sensor Blindness Failsafe** | Loss of camera or LiDAR telemetry triggers automatic transition to `AutonomousSafeHold` | [ ] |
| 7 | **NPU-to-MCU Admission Gateway** | Communication between Linux planner and safety microcontroller strictly framed with CRC32 packet integrity checks | [ ] |
| 8 | **GNSS Spoofing Rejection** | Multi-sensor cross-validation detects and rejects discontinuous coordinate jumps inconsistent with IMU odometry | [ ] |
| 9 | **Cryptographically Signed Config** | Operational rule packs and geofence coordinates authenticated via developer digital signature ([Chapter 22](ch22-cybernetics-edge-to-backend.md)) | [ ] |
| 10 | **Immutable Onboard Flight Recorder** | Last 10 minutes of raw telemetry, state transitions, and action contracts recorded to non-volatile flash | [ ] |

Each item in this checklist maps directly to an architectural safeguard detailed above; the completed checklist constitutes an indispensable element of the platform's formal safety case.

## Conclusions

This appendix addressed a central engineering challenge: how to design an onboard expert system that permits only mathematically verifiable, safe actions to reach physical actuators. The solution comprises multi-tiered defense: probabilistic neural models on the NPU merely generate candidate proposals; a safety microcontroller evaluates proposals against compiled Datalog rules without dynamic memory allocation; and programmable logic enforces a hardware veto. The Simplex architecture guarantees seamless fail-over to a baseline safe controller whenever advanced proposals violate safety boundaries. Physical actions are structured as contracts with explicit preconditions, postconditions, and compensation logic, while external sensor feeds are admitted only when accompanied by cryptographically signed feature payloads that allow independent verification.

The boundaries of these techniques are clearly defined. Numerical parameters in this text (thresholds, timeouts, margins) are illustrative and must be calibrated through empirical safety analysis. Integrity ratings such as ASIL or SIL require formal certification processes rather than architectural declarations alone. Subsystem interfaces and physical models require dedicated validation, as non-functional interactions can cause fully conforming components to fail in aggregate. GPS-denied navigation is explored further in [Appendix C](appendix-c-autonomous-navigation-and-geosearch.md), while authority control over autonomous actions is established in [Chapter 21](ch21-from-recommendation-to-action.md).

## Glossary

| Term | English Equivalent | Definition |
|---|---|---|
| Cyber-Physical System | *cyber-physical system* | System in which computational algorithms directly monitor and control physical processes |
| Admission Gate | *admission gate* | Deterministic architectural component that authorizes or rejects proposed actions |
| Simplex Architecture | *Simplex architecture* | System design wherein a simple verified controller supervises a complex unverified controller |
| Safety Invariant | *safety invariant* | A formal condition that must remain true across all valid operational states |
| Watchdog Timer | *watchdog timer* | Hardware timer that triggers an emergency fallback if not periodically refreshed by software |
| Lockstep Core Execution | *lockstep* | Architecture where redundant processor cores execute identical instructions to detect faults |
| Action Contract | *action contract* | Formal action descriptor specifying preconditions, execution logic, postconditions, and compensation |
| Saga | *saga* | Distributed transaction composed of discrete operations, each backed by a compensating action |
| Compensating Action | *compensating action* | Operation designed to return the system to a safe state upon failure of a transaction step |
| GNSS Spoofing | *GNSS spoofing* | Malicious broadcast of false satellite navigation signals to manipulate receiver positioning |
| Feature Signing | *feature signing* | Cryptographic signing of raw evidentiary features rather than high-level classification verdicts |
| Non-Functional Interface | *non-functional interface* | Unintended physical coupling between subsystems via power lines, EMI, heat, or vibration |
| Mesh Convergence | *mesh convergence* | Demonstration that numerical FEA solutions stabilize asymptotically under grid refinement |
| Hardware-in-the-Loop Simulation | *hardware-in-the-loop* | Testing physical embedded controllers against real-time simulated physical environments |
| Synchronous Measurement Front-End | *synchronous measurement front-end* | Hardware subsystem for synchronous ADC channel strobing dedicated to spectral state analysis |
| Electrochemical Impedance Spectroscopy | *electrochemical impedance spectroscopy (EIS)* | Method for diagnosing conductor and battery cell degradation via complex frequency response |
| Secondary Measurement Path | *secondary measurement path* | Independent redundant on-chip measurement path enabling hardware plausibility cross-checks |
| Sensor-Frontline Fail-Closed Gate | *sensor-frontline fail-closed gate* | Four-tier frontline filter shielding reasoning engines from IEEE 754 anomalies, jagged matrices, and sensor lockup |
| Johnson-Nyquist Noise | *Johnson-Nyquist noise* | Inherent thermodynamic electronic conductor noise establishing the physical lower bound of valid sensory signals |
| Epistemic Governor | *epistemic governor* | Deterministic onboard expert system serving as supreme arbiter of admission for any physical action |
| Human-on-the-Loop | *human-on-the-loop* | Concept of passive external supervision over an autonomous system with emergency abort authority |
| Safe State Ontology | *safe state ontology* | Formal model for deterministic selection of an emergency saga based on geospatial context and faults |
| Hardware Disposition Gate | *hardware disposition gate* | Hardware combinational circuit for instantaneous actuator cut-off (latency $< 100\,\text{ns}$) using DI/DO and CAN FD |
| Popperian Red Teamer | *Popperian red teamer* | Adversarial testing engine generating boundary vectors and falsifying safety hypotheses targeting $\mathrm{ZHR} = 1.00$ |
| Ballistic Jettison | *ballistic jettison* | Emergency saga deactivating the warhead and executing controlled gliding into a designated safe ditching sector |

## Abbreviations

| Abbreviation | Expansion | Meaning |
|---|---|---|
| ADC | Analog-to-Digital Converter | Device converting continuous analog physical signals into discrete digital representations |
| ASIL | Automotive Safety Integrity Level | Risk classification scheme defined by ISO 26262 for automotive systems |
| BGA | Ball Grid Array | Surface-mount packaging for integrated circuits utilizing solder spheres |
| BVA | Boundary Value Analysis | Software and knowledge testing technique evaluating boundary conditions |
| CAE | Computer-Aided Engineering | Software tools utilized for analyzing and simulating engineering designs |
| CAN FD | Controller Area Network Flexible Data-Rate | CAN communication bus protocol supporting higher data throughput |
| CHT | Conjugate Heat Transfer | Coupled simulation of thermal conduction in solids and fluid convection |
| CPLD | Complex Programmable Logic Device | Programmable logic device simpler and lower-power than an FPGA |
| CPS | Cyber-Physical System | System tightly integrating physical mechanisms with software computation |
| CRC32 | Cyclic Redundancy Check, 32 bits | Error-detecting code algorithm producing a 32-bit checksum |
| DLA | Deep Learning Accelerator | Dedicated hardware block for accelerating neural network inference (NVDLA) |
| DMA | Direct Memory Access | Feature allowing hardware subsystems to access system memory without CPU intervention |
| EIS | Electrochemical Impedance Spectroscopy | Diagnostic method evaluating material impedance spectra across varying frequencies |
| EPU | Epistemic Processing Unit | Dedicated hardware coprocessor for deterministic rule evaluation and fact verification |
| FIT | Failures in Time | Standard unit of failure rate ($10^{-9}$ component failures per operating hour) |
| FPGA | Field-Programmable Gate Array | Semiconductor device based around a matrix of configurable logic blocks |
| GMSL | Gigabit Multimedia Serial Link | High-speed serial interconnect protocol for automotive cameras and sensors |
| GNSS | Global Navigation Satellite System | Satellite-based radio navigation systems (GPS, GLONASS, Galileo, BeiDou) |
| HIL | Hardware-in-the-Loop | Test architecture connecting physical electronic controllers to a virtual plant model |
| HOTL | Human-on-the-Loop | Passive remote supervisory control over an autonomous system |
| IMU | Inertial Measurement Unit | Sensor package combining accelerometers and gyroscopes to measure movement |
| isoUART | Isolated Universal Asynchronous Receiver-Transmitter | Galvanically isolated asynchronous serial communication protocol |
| MCU | Microcontroller Unit | Compact integrated circuit designed to govern an embedded system |
| NPU | Neural Processing Unit | Specialized microprocessor tailored for accelerating machine learning algorithms |
| RoE | Rules of Engagement | Operational directives delineating permissible actions during a mission |
| RPM | Revolutions Per Minute | Measure of the frequency of rotation around an axis |
| RTOS | Real-Time Operating System | Operating system providing deterministic response times to external events |
| SIL | Safety Integrity Level | Relative level of safety risk reduction defined by IEC 61508 |
| SLAM | Simultaneous Localization and Mapping | Computational problem of constructing a map of an unknown environment while localizing |
| SoC | System on Chip | Integrated circuit that integrates all components of an electronic system |
| SWaP-C | Size, Weight, Power and Cost | Engineering constraints governing embedded and deployed equipment |
| URAM | UltraRAM | High-density synchronous static memory macro on Xilinx FPGAs with single-cycle latency |
| WCET | Worst-Case Execution Time | Maximum execution duration required for a computational task on a given architecture |
| PWM | Pulse-Width Modulation | Technique for controlling power supplied to electrical loads via variable-width pulsing |
| ZHR | Zero Hallucination Rate | Metric denoting complete absence of hallucinations ($\mathrm{ZHR} \equiv 1.00$) |

## References

1. <a id="src-1"></a>ROS 2. [*The Robot Operating System*](https://github.com/ros2/ros2). GitHub repository.
2. <a id="src-2"></a>Zephyr Project. [*Getting Started Guide*](https://docs.zephyrproject.org/latest/develop/getting_started/index.html). Zephyr Project Documentation.
3. <a id="src-3"></a>ISO. [*ISO 26262-1:2018. Road vehicles: Functional safety: Part 1: Vocabulary*](https://www.iso.org/standard/68383.html). 2018.
4. <a id="src-4"></a>IEC. [*IEC 61508-1:2010. Functional safety of electrical/electronic/programmable electronic safety-related systems: Part 1: General requirements*](https://webstore.iec.ch/en/publication/5515). 2010.
5. <a id="src-5"></a>Lui Sha. [*Using Simplicity to Control Complexity*](https://doi.org/10.1109/MS.2001.936213). *IEEE Software*, 18(4), 20–28, 2001.
6. <a id="src-6"></a>ROS Navigation. [*Navigation2: ROS 2 Navigation Framework and System*](https://github.com/ros-navigation/navigation2). GitHub repository.
7. <a id="src-7"></a>S. Ceri, G. Gottlob, L. Tanca. [*What You Always Wanted to Know About Datalog (and Never Dared to Ask)*](https://doi.org/10.1109/69.43410). *IEEE Transactions on Knowledge and Data Engineering*, 1(1), 146–166, 1989.
8. <a id="src-8"></a>Hector Garcia-Molina, Kenneth Salem. [*Sagas*](https://doi.org/10.1145/38713.38742). *Proceedings of the 1987 ACM SIGMOD International Conference on Management of Data*, 249–259, 1987.
9. <a id="src-9"></a>Mark L. Psiaki, Todd E. Humphreys. [*GNSS Spoofing and Detection*](https://doi.org/10.1109/JPROC.2016.2526658). *Proceedings of the IEEE*, 104(6), 1258–1270, 2016.
10. <a id="src-10"></a>R. E. Kalman. [*A New Approach to Linear Filtering and Prediction Problems*](https://doi.org/10.1115/1.3662552). *Journal of Basic Engineering*, 82(1), 35–45, 1960.
11. <a id="src-11"></a>Petro Sidliarchuk. [*Witness Integrity for Measurement Devices*](https://github.com/sidliarchukpetro/witness-integrity). GitHub repository, InfraVeritas project.
12. <a id="src-12"></a>NASA. [*NASA Systems Engineering Handbook*](https://www.nasa.gov/reference/systems-engineering-handbook/). NASA SP-2016-6105 Rev2.
13. <a id="src-13"></a>William L. Oberkampf, Christopher J. Roy. [*Verification and Validation in Scientific Computing*](https://doi.org/10.1017/CBO9780511760396). Cambridge University Press, 2010.
14. <a id="src-14"></a>P. J. Roache. [*Perspective: A Method for Uniform Reporting of Grid Refinement Studies*](https://doi.org/10.1115/1.2910291). *Journal of Fluids Engineering*, 116(3), 405–413, 1994.
15. <a id="src-15"></a>Infineon Technologies. [*TLE9012DQU / TLE9015DQU: Multi-Channel Battery Monitoring & UART Transceiver ICs*](https://www.infineon.com/). User Manual, Rev 1.2, 2023.
16. <a id="src-16"></a>J. B. Johnson. [*Thermal Agitation of Electricity in Conductors*](https://doi.org/10.1103/PhysRev.32.97). *Physical Review*, 32(1), 97–109, 1928; H. Nyquist. [*Thermal Agitation of Electric Charge in Conductors*](https://doi.org/10.1103/PhysRev.32.110). *Physical Review*, 32(1), 110–113, 1928.
17. <a id="src-17"></a>Karl R. Popper. [*The Logic of Scientific Discovery*](https://doi.org/10.4324/9780203994627). Routledge, 1959 / 2002.
18. <a id="src-18"></a>Defense Advanced Research Projects Agency (DARPA). [*Assurance Solutions for Operational Resilience (ASOR)*](https://www.darpa.mil/). Technical Report, 2023.

---

[← Appendix A](appendix-a-evidence-governed-framework.md) | [Table of Contents](README.md) | [Appendix C →](appendix-c-autonomous-navigation-and-geosearch.md)

