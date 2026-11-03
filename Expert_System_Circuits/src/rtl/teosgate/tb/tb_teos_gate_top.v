// Exercises the FPGA pin wrapper: shifts a strided subset of vectors in serially
// (LSB first), pulses load, and checks the registered result pins.
`timescale 1ns/1ps
`default_nettype none
module tb_teos_gate_top;
    parameter integer N = 1;
    parameter integer M = 256;   // vectors checked through the serial loader
    parameter integer VARIANT = 0;  // 1 = pipelined (one extra clock)
    parameter VECFILE = "build/vectors.hex";

    reg [471:0] vec [0:N-1];
    reg clk = 1'b0, sin = 1'b0, shift = 1'b0, load = 1'b0;
    wire [9:0] result;

    teos_gate_top #(.VARIANT(VARIANT)) dut (.clk(clk), .sin(sin), .shift(shift), .load(load), .result(result));

    always #5 clk = ~clk;

    integer k, b, idx, stride, mism, checked;
    reg [471:0] v;

    initial begin
        $readmemh(VECFILE, vec);
        mism = 0; checked = 0;
        stride = (N > M) ? N / M : 1;
        for (k = 0; k < M && k * stride < N; k = k + 1) begin
            idx = k * stride;
            v = vec[idx];
            for (b = 0; b < 462; b = b + 1) begin
                @(negedge clk);
                sin = v[b]; shift = 1'b1; load = 1'b0;
            end
            @(negedge clk);
            shift = 1'b0; load = 1'b1;
            @(negedge clk);
            load = 1'b0;
            repeat (1 + (VARIANT == 1)) @(negedge clk);
            checked = checked + 1;
            if (result !== v[471:462]) begin
                mism = mism + 1;
                if (mism <= 10)
                    $display("MISMATCH top vector %0d: got %b expected %b", idx, result, v[471:462]);
            end
        end
        $display("RESULT top_variant%0d checked=%0d mismatches=%0d", VARIANT, checked, mism);
        $finish;
    end
endmodule
`default_nettype wire
