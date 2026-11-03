// Streams one Go-generated vector per clock into gate_reg (one tick, direct
// and prefix cores) and gate_pipe (2-clock latency) and compares every registered result with
// the expected value computed by hardware.PrevalidatedGate.
// Vector line: {expected[9:0], word[461:0]} in hex (cmd/gatevectors).
`timescale 1ns/1ps
`default_nettype none
module tb_gate;
    parameter integer N = 1;
    parameter VECFILE = "build/vectors.hex";

    reg [471:0] vec [0:N-1];
    reg         clk = 1'b0;
    reg         load = 1'b0;
    reg [461:0] word_in = 462'd0;
    wire [9:0]  res_one, res_pipe, res_prefix;

    gate_reg #(.PREFIX(1)) dut_prefix (.clk(clk), .load(load), .word_in(word_in), .result_q(res_prefix));
    gate_reg      dut_one  (.clk(clk), .load(load), .word_in(word_in), .result_q(res_one));
    gate_pipe dut_pipe (.clk(clk), .load(load), .word_in(word_in), .result_q(res_pipe));

    always #5 clk = ~clk;

    integer i, mism_prefix, mism_one, mism_pipe, checked_one, checked_pipe, allow_seen;
    reg [9:0] exp1, exp2;

    initial begin
        $readmemh(VECFILE, vec);
        mism_prefix = 0; mism_one = 0; mism_pipe = 0; checked_one = 0; checked_pipe = 0; allow_seen = 0;
        // Vector i is driven at negedge i and captured into the word register at
        // the next posedge; the one-tick result register updates one posedge
        // later (observed at negedge i+2), the pipelined one at negedge i+3.
        for (i = 0; i < N + 3; i = i + 1) begin
            @(negedge clk);
            if (i >= 2 && i - 2 < N) begin
                exp1 = vec[i - 2][471:462];
                checked_one = checked_one + 1;
                if (exp1[0]) allow_seen = allow_seen + 1;
                if (res_one !== exp1) begin
                    mism_one = mism_one + 1;
                    if (mism_one <= 10)
                        $display("MISMATCH one-tick vector %0d: got %b expected %b", i - 2, res_one, exp1);
                end
                if (res_prefix !== exp1) begin
                    mism_prefix = mism_prefix + 1;
                    if (mism_prefix <= 10)
                        $display("MISMATCH one-tick-prefix vector %0d: got %b expected %b", i - 2, res_prefix, exp1);
                end
            end
            if (i >= 3) begin
                exp2 = vec[i - 3][471:462];
                checked_pipe = checked_pipe + 1;
                if (res_pipe !== exp2) begin
                    mism_pipe = mism_pipe + 1;
                    if (mism_pipe <= 10)
                        $display("MISMATCH pipelined vector %0d: got %b expected %b", i - 3, res_pipe, exp2);
                end
            end
            if (i < N) begin
                word_in = vec[i][461:0];
                load = 1'b1;
            end else begin
                load = 1'b0;
            end
        end
        $display("RESULT one_tick checked=%0d mismatches=%0d allow=%0d", checked_one, mism_one, allow_seen);
        $display("RESULT one_tick_prefix checked=%0d mismatches=%0d", checked_one, mism_prefix);
        $display("RESULT pipelined checked=%0d mismatches=%0d", checked_pipe, mism_pipe);
        $finish;
    end
endmodule
`default_nettype wire
