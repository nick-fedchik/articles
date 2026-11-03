# H5 flow: Go reference vectors -> iverilog/verilator simulation -> yosys
# synth_ice40 -> nextpnr-ice40 static timing. Run from src/rtl/teosgate.
SRC_ROOT := ../..
BUILD    := build
RESULTS  := $(SRC_ROOT)/testdata/teos-rtl
RTL      := rtl/teos_gate_core.v rtl/teos_prefix_carry.v rtl/teos_gate_core_prefix.v rtl/teos_gate.v rtl/teos_gate_pipe.v rtl/teos_gate_top.v
FREQ     ?= 100
SEEDS    ?= 1 2 3 4 5
RANDOM   ?= 10000
TOP_M    ?= 512
# device tag : nextpnr flags
DEVICES  := up5k-sg48 hx8k-ct256
up5k-sg48_FLAGS  := --up5k --package sg48
hx8k-ct256_FLAGS := --hx8k --package ct256
VARIANTS := p0 p1 p2

N = $(shell cat $(BUILD)/vectors.count)

.PHONY: all versions vectors sim sim-iverilog sim-verilator lint mutants synth synth-gate pnr summary clean
all: summary

$(BUILD):
	mkdir -p $(BUILD)

versions: | $(BUILD)
	{ go version; yosys -V; iverilog -V 2>&1 | head -1; verilator --version; \
	  nextpnr-ice40 --version 2>&1 | head -1; dpkg-query -W -f='$${Package} $${Version}\n' fpga-icestorm 2>/dev/null || true; \
	  python3 --version; } > $(BUILD)/versions.txt
	cat $(BUILD)/versions.txt

vectors: $(BUILD)/vectors.count
$(BUILD)/vectors.count: | $(BUILD)
	cd $(SRC_ROOT) && go run ./cmd/teosgatevectors -out rtl/teosgate/$(BUILD) -random $(RANDOM) > rtl/teosgate/$(BUILD)/vectors_summary.json

sim: sim-iverilog sim-verilator

sim-iverilog: vectors
	iverilog -g2005 -Wall -P tb_teos_gate.N=$(N) -s tb_teos_gate -o $(BUILD)/tb_gate.vvp $(RTL) tb/tb_teos_gate.v
	vvp -n $(BUILD)/tb_gate.vvp | tee $(BUILD)/sim_iverilog_gate.log
	for p in 0 1 2; do \
	  iverilog -g2005 -Wall -P tb_teos_gate_top.N=$(N) -P tb_teos_gate_top.M=$(TOP_M) -P tb_teos_gate_top.VARIANT=$$p \
	    -s tb_teos_gate_top -o $(BUILD)/tb_top_p$$p.vvp $(RTL) tb/tb_teos_gate_top.v && \
	  vvp -n $(BUILD)/tb_top_p$$p.vvp | tee $(BUILD)/sim_iverilog_top_p$$p.log; done

sim-verilator: vectors
	verilator --binary --timing -Wno-fatal -GN=$(N) --top-module tb_teos_gate -Mdir $(BUILD)/vl_gate \
	  $(RTL) tb/tb_teos_gate.v > $(BUILD)/verilator_build_gate.log 2>&1
	$(BUILD)/vl_gate/Vtb_teos_gate | tee $(BUILD)/sim_verilator_gate.log
	for p in 0 1 2; do \
	  verilator --binary --timing -Wno-fatal -GN=$(N) -GM=$(TOP_M) -GVARIANT=$$p --top-module tb_teos_gate_top \
	    -Mdir $(BUILD)/vl_top_p$$p $(RTL) tb/tb_teos_gate_top.v > $(BUILD)/verilator_build_top_p$$p.log 2>&1 && \
	  $(BUILD)/vl_top_p$$p/Vtb_teos_gate_top | tee $(BUILD)/sim_verilator_top_p$$p.log; done

lint: | $(BUILD)
	: > $(BUILD)/lint.log
	for p in 0 1 2; do verilator --lint-only -Wall -GVARIANT=$$p --top-module teos_gate_top $(RTL) >> $(BUILD)/lint.log 2>&1; \
	  echo "verilator lint VARIANT=$$p exit $$?" >> $(BUILD)/lint.log; done
	cat $(BUILD)/lint.log

mutants: vectors
	./scripts/mutants.sh $(BUILD)

# Gate alone (word register + predicate + result register), no pin wrapper.
# Cell counts only: 462+10 ports exceed package pins, so no place and route.
synth-gate: | $(BUILD)
	yosys -q -l $(BUILD)/yosys_gate_direct.log -p "read_verilog $(RTL); synth_ice40 -top teos_gate; tee -o $(BUILD)/yosys_gate_stat_direct.txt stat"
	yosys -q -l $(BUILD)/yosys_gate_prefix.log -p "read_verilog $(RTL); hierarchy -top teos_gate -chparam PREFIX 1; synth_ice40 -top teos_gate; tee -o $(BUILD)/yosys_gate_stat_prefix.txt stat"
	yosys -q -l $(BUILD)/yosys_gate_pipe.log -p "read_verilog $(RTL); synth_ice40 -top teos_gate_pipe; tee -o $(BUILD)/yosys_gate_stat_pipe.txt stat"

synth: $(foreach v,$(VARIANTS),$(BUILD)/synth_$(v).json)
$(BUILD)/synth_p%.json: $(RTL) | $(BUILD)
	yosys -q -l $(BUILD)/yosys_p$*.log -p "read_verilog $(RTL); hierarchy -top teos_gate_top -chparam VARIANT $*; \
	  synth_ice40 -top teos_gate_top -json $@; tee -o $(BUILD)/yosys_stat_p$*.txt stat"

PNR_TARGETS := $(foreach d,$(DEVICES),$(foreach v,$(VARIANTS),$(foreach s,$(SEEDS),$(BUILD)/pnr_$(d)_$(v)_s$(s).json)))
pnr: $(PNR_TARGETS)
define PNR_RULE
$(BUILD)/pnr_$(1)_$(2)_s$(3).json: $(BUILD)/synth_$(2).json
	nextpnr-ice40 $($(1)_FLAGS) --freq $(FREQ) --seed $(3) --pcf-allow-unconstrained --timing-allow-fail \
	  --json $$< --asc $(BUILD)/pnr_$(1)_$(2)_s$(3).asc --report $$@ -q -l $(BUILD)/pnr_$(1)_$(2)_s$(3).log
endef
$(foreach d,$(DEVICES),$(foreach v,$(VARIANTS),$(foreach s,$(SEEDS),$(eval $(call PNR_RULE,$(d),$(v),$(s))))))

summary: versions sim lint mutants synth-gate pnr
	mkdir -p $(RESULTS)
	python3 scripts/summarize.py --build $(BUILD) --freq $(FREQ) --out $(RESULTS)

clean:
	rm -rf $(BUILD)
