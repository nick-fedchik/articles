# PYNQ-Z1 sensor feasibility for drift monitoring

## What interfaces and electrical limits does the PYNQ-Z1 provide?

### Takeaway
PYNQ-Z1 is plausible for a low-voltage bench prototype, not as a protected automotive front end. Confirm pin routes and ranges from the board manual before connecting sensors.

### Cited Findings
- **Board specifications (verification pending):** Digilent’s manual describes Zynq-7000, Arduino/chipKIT and Pmod headers, and XADC access. PYNQ-Z1 is specified with dual-core Cortex-A9 PS, XC7Z020 PL, and 512 MB DDR3. Digilent’s page returned HTTP 403 here, so exact board specs and pin mapping could not be checked directly. [Digilent manual](https://digilent.com/reference/programmable-logic/pynq-z1/reference-manual)
- Zynq-7000 combines ARM processing system and programmable logic; the 7-series XADC offers 12-bit conversion up to 1 MSPS. That is converter capability, not system sample throughput. Input range depends on the ADC pins and board circuitry. [AMD TRM](https://docs.amd.com/v/u/en-US/ug585-Zynq-7000-TRM); [AMD XADC guide](https://docs.amd.com/v/u/en-US/ug480_7Series_XADC)
- The board’s Arduino analog routes/JXADC expose XADC; Arduino/Pmod digital signals are low-voltage (nominal 3.3 V), not automotive-tolerant. Verify analog scaling, pin sharing, and limits in the schematic. [Digilent manual](https://digilent.com/reference/programmable-logic/pynq-z1/reference-manual); [PYNQ docs](https://pynq.readthedocs.io/en/latest/)

### Inferences
- Do not connect unknown vehicle/controller pins directly. Identify voltage, ground, protocol and protection needs first; an automotive AFE, level shifter, isolation, or CAN/LIN transceiver may be required.

### Gaps
- Manual retrieval returned 403, PYNQ board route returned 404, and AMD docs rendered an application shell. Exact ranges, resistor networks, connector map, current limits, OS image, and board revision remain unverified.

## Can an automotive capacitive HMI provide raw sensor samples?

### Takeaway
Raw electrode data from the TSG controller is unknown; treat access as a dependency, not a PYNQ feature.

### Cited Findings
- TI’s FDC1004 is a dedicated capacitance-to-digital converter that reports measurements over I²C: a possible bench AFE, not evidence about the TSG controller or automotive qualification. [TI datasheet](https://www.ti.com/lit/ds/symlink/fdc1004.pdf)
- No TSG model, register map, datasheet, or debug-interface documentation was available.

### Inferences
- Identify the controller and determine whether host registers expose raw counts/baseline/noise or only touch events. If only events are available, use an approved debug mode or a separate CDC/MCU; touch events alone cannot substantiate electrode drift. For the bench, acquire timestamped values through a documented AFE/MCU I²C/SPI interface, rather than treating XADC as a capacitance meter.

### Gaps
- Raw-data availability, bus, cadence, and voltage levels require the TSG datasheet or vendor confirmation.

## What should run on PS, PL, Go, and Ollama?

### Takeaway
Start with Go on PS as acquisition and software reference; add a small fixed-window PL feature/threshold block only after software works. Keep Ollama remote and advisory.

### Cited Findings
- Go documents Linux ARMv7 as `GOARCH=arm`, `GOARM=7`; deployment still depends on the board’s Linux ABI. [Go ARM support](https://go.dev/wiki/GoArm)
- Ollama documents a Linux ARM64 package, not ARMv7. Thus official docs do not establish Ollama support on the PYNQ’s 32-bit Cortex-A9 PS. [Ollama Linux](https://docs.ollama.com/linux)
- “Go ARMv7 host agent,” “BRAM/UIO check virtual machine,” and “Ollama on Jetson” are design claims of the research prototype, not verified operation; node problems prevented live inventory.

### Inferences
- PS/Go should timestamp, sequence, buffer/log AFE samples, batch-transfer to PL, and report to Jetson. PL’s first function can calculate fixed-window mean/baseline deviation, slope or noise, plus a threshold flag; Go remains the golden reference.
- BRAM/UIO is plausible but needs a working overlay, map, device-tree/kernel access, correct cache/order handling, and measured transfer overhead. SPI/I²C, Linux scheduling, copies, and reporting may cost more than simple arithmetic. FPGA use alone does not guarantee end-to-end hard real-time behavior.
- Send summaries asynchronously to Ollama; never put advisory inference in a safety decision path. Validate Jetson package/driver/model support separately.

### Gaps
- Board OS, Go runtime, overlay, UIO/BRAM path, network, and Jetson setup remain unverified.

## How to measure feasibility without inventing performance targets?

### Takeaway
Replay identical timestamped samples through Go-only and Go-plus-PL paths; compare correctness and measured system cost against owner-defined requirements.

### Cited Findings
- XADC conversion rate excludes acquisition, transport, and decision latency. Go ARMv7 is documented, whereas Ollama ARM64 is documented and ARMv7 support is not. [AMD XADC guide](https://docs.amd.com/v/u/en-US/ug480_7Series_XADC); [Go ARM support](https://go.dev/wiki/GoArm); [Ollama Linux](https://docs.ollama.com/linux)

### Inferences
- First capture a bench/approved-controller stream. Run the same sample IDs through Go and PL; compare each feature and decision (exact integer match or predeclared quantization tolerance), overflow/reset behavior, and dropped/reordered samples.
- Timestamp acquisition, transfer, window completion, and decision. Measure latency distribution, loss, CPU/PL load, and board input energy; measure Jetson separately. **Go** only if equivalence and sample integrity hold and measured benefit meets approved budgets beyond measurement uncertainty. Otherwise keep Go and resolve the interface or overhead issue.

### Gaps
- Accuracy, required cadence, delay/energy budgets, and environmental envelope are unspecified. Set them from HMI requirements before sizing or numerical pass/fail limits.
