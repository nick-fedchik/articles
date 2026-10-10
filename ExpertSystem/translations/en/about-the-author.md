# About the Author: Mykola Fedchyk (Nick Fedchik)

> **Book:** [Architecture of Evidence-Governed Expert Systems](README.md)  
> **Table of Contents:** [README.md](README.md) · [Appendix A](appendix-a-evidence-governed-framework.md)  
> **Profiles:** [LinkedIn — Nick Fedchik](https://www.linkedin.com/in/nickfedchik/) · [GitHub — nick-fedchik](https://github.com/nick-fedchik)  

---

## Professional Profile & Engineering Background

**Mykola Fedchyk (Nick Fedchik)** is a Ukrainian systems architect, information technology engineer, and researcher in cyber-physical systems, expert artificial intelligence, and safety-critical embedded electronics.

With more than two decades of hands-on experience in software engineering, distributed network architecture design, low-level hard real-time systems, and technological product management of complex high-tech innovations.

Today, he serves as **Technical Project Manager (Automotive)** at the global semiconductor leader **Infineon Technologies**, focusing on the development and deployment of next-generation automotive platform solutions:

* High-reliability multi-core microcontrollers of the **Infineon AURIX (TriCore lockstep)** family;
* Attainment of the most stringent functional safety certifications (**ISO 26262 ASIL-D**);
* Real-Time Operating Systems (RTOS, specifically **Zephyr RTOS**);
* Hardware-software co-design at the intersection of deterministic control, sensor networks, and edge artificial intelligence (Edge AI).

---

## Systems Roots, Linux Kernel & Open Source Contributions

The author's engineering mindset was forged in low-level systems programming, networking architecture, and operating system kernel engineering. Beginning in the early 2000s, Mykola directly contributed to foundational systems software, the Linux kernel networking stack, and critical open-source infrastructure:

### 1. Linux Kernel and the ebtables (Ethernet Bridge Tables) Subsystem

* **Development of the `ebt_vlan` module:** Mykola is the original author of the 802.1Q VLAN packet filtering module for the Ethernet bridge filtering subsystem (**ebtables / Netfilter**) in the Linux kernel (`net/bridge/netfilter/ebt_vlan.c`).
* The module introduced hardware- and software-level inspection, matching, and filtering of VLAN tags (802.1Q Priority / VID) directly within the L2 Ethernet bridge switching path of the Linux kernel, becoming an industry standard for bridge firewalls, carrier-grade routers, and secure managed switches.
* Engineering low-level kernel data structures (`sk_buff`), enforcing zero memory allocations in the packet-forwarding fast path, and deterministic interrupt handling laid the foundational principles for how latency-critical, fail-closed systems must be architected.

### 2. The BusyBox Project & Embedded Linux

* **Porting and maintenance of `arping`:** Mykola authored the official port of the `arping` network diagnostics utility to the **BusyBox** codebase (`networking/arping.c` — *"The Swiss Army Knife of Embedded Linux"*).
* The tool enabled MAC/IP address collision detection, link diagnostics, and ARP probing under strictly bounded storage and RAM constraints typical of embedded controllers, routers, and industrial RTUs.
* In addition, the author contributed optimization patches and bug fixes to core system utilities across embedded Linux distributions running on millions of industrial devices worldwide.

### 3. The GNU C Library (glibc) & Protocol Infrastructure

* **Enhancing glibc functionality:** The author authored patches for the **glibc** standard library implementing `getethertype` resolution via the Name Service Switch (NSS) mechanism, mapping Ethernet protocol names and numbers through `/etc/ethertypes`. This facility was widely adopted by traffic monitoring utilities, packet filters, and protocol decoders.
* Upstream contributions to wireless and peripheral subsystems (notably the `irda-usb` driver in the Linux kernel).

### 4. From Kernel Code to the Architecture of Evidence-Governed Expert Systems

This extensive experience in the lowest layers of operating system infrastructure directly defined the architectural ethos of this monograph:

* **Determinism over Stochasticity:** Just as the `ebtables` firewall cannot "probably" forward or drop a packet based on likelihood, an evidence-governed expert system must never base mission-critical conclusions on probabilistic conjectures.
* **The Admission Gateway (Gatekeeper Pattern):** The symbolic verification engine in a neuro-symbolic architecture acts exactly like a kernel firewall: it enforces unyielding mathematical invariants, deterministically rejecting invalid LLM proposals before they ever reach the execution runtime.
* **Direct Grounding in Silicon:** The trajectory from low-level Linux kernel device drivers to modern automotive silicon platforms (**Infineon AURIX TC4xx**) and hard real-time operating systems (**Zephyr RTOS**) provides a coherent, end-to-end perspective — spanning silicon transistors, lockstep execution pipelines, formal logic, ontologies, and artificial intelligence.

---

## Research Initiatives & Evangelism of Evidence-Based AI

Mykola is the author of an influential series of engineering research publications on the Ukrainian developer community DOU, addressing:

1. **The Documented Knowledge Crisis:** Overcoming the "PDF document graveyard" anti-pattern in large-scale enterprise R&D programs;
2. **Engineering Knowledge Graphs (EKG):** End-to-end, deterministic traceability linking regulatory standards, source code, verification suites, and physical silicon;
3. **Neuro-Symbolic Architecture (NeSy AI):** Designing hybrid expert systems where stochastic language models (SLMs/LLMs) are strictly segregated from the safety perimeter and subordinate to a deterministic Go and Datalog symbolic core;
4. **Automated Safety Certification (GSN):** Algorithmic synthesis of legally and technically robust safety proof trees (*Safety Cases* via Goal Structuring Notation).

This monograph represents the systematic synthesis of decades of production engineering experience bridging the gap between pure mathematical logic, generative AI, and the uncompromising realities of safety-critical industrial hardware.

---

## Engineering Creed

Two decades of production engineering — spanning bare-metal Linux drivers, severely constrained microcontroller memories, dual-core lockstep silicon architectures (Infineon AURIX), and ISO 26262 ASIL-D functional safety standards — forged an unyielding engineering conviction:

> **"True reliability begins where assumptions end and determinism begins."**

* **In the OS kernel**, interrupt handlers forgive zero unpredictable latency spikes, and the networking stack never speaks in terms of "transmission probabilities": a frame either strictly complies with bit-level protocol invariants or is deterministically dropped.
* **In automotive electronics and cyber-physical actuation systems**, safety is never quantified as a percentage of model confidence. When dual processor cores execute in cycle-by-cycle lockstep, a single bit mismatch is a hardware fault, not an "acceptable generative deviation." The physical world never forgives hallucinations.
* **In the era of Artificial Intelligence**, engineering rigor is more vital than ever. Neural networks and large language models are exceptional heuristic generators of plausible conjectures, but they must never wield unmediated authority over physical actuators. Ultimate decision authority must rest exclusively with formal verifiers, mathematical logic, and hardware-governed evidence:

> *"The neural network proposes — The symbolic core verifies — Deterministic hardware guarantees safety."*

---

## Contact & Connect

* **LinkedIn:** [https://www.linkedin.com/in/nickfedchik/](https://www.linkedin.com/in/nickfedchik/)
* **GitHub:** [https://github.com/nick-fedchik](https://github.com/nick-fedchik)
* **Monograph Repository:** [https://github.com/nick-fedchik/articles](https://github.com/nick-fedchik/articles)
