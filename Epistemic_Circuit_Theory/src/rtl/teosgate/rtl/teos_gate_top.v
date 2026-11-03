// FPGA pin wrapper (iCE40). The 462-bit word cannot be bonded to package pins,
// so it is loaded LSB first through a serial shift register (1 bit per clock);
// `load` then copies it into the gate's word register. This wrapper adds
// 462 flip-flops and the shift path; it is NOT part of the gate predicate and
// its loading time (462 clocks) is excluded from the one-tick claim.
// VARIANT=0: teos_gate, direct core (one tick, primary H5 design);
// VARIANT=1: teos_gate_pipe (2-clock latency, II=1);
// VARIANT=2: teos_gate with prefix-comparator core (one tick).
`timescale 1ns/1ps
`default_nettype none
module teos_gate_top #(
    parameter VARIANT = 0
) (
    input  wire       clk,
    input  wire       sin,
    input  wire       shift,
    input  wire       load,
    output wire [9:0] result
);
    reg [461:0] sr;
    initial sr = 462'd0;
    always @(posedge clk)
        if (shift)
            sr <= {sin, sr[461:1]};

    generate
        if (VARIANT == 0) begin : g_one_tick
            teos_gate #(.PREFIX(0)) gate (.clk(clk), .load(load), .word_in(sr), .result_q(result));
        end else if (VARIANT == 2) begin : g_one_tick_prefix
            teos_gate #(.PREFIX(1)) gate (.clk(clk), .load(load), .word_in(sr), .result_q(result));
        end else begin : g_pipe
            teos_gate_pipe gate (.clk(clk), .load(load), .word_in(sr), .result_q(result));
        end
    endgenerate
endmodule
`default_nettype wire
