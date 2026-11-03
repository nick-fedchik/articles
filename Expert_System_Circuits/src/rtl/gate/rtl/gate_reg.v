// Register-to-register wrapper: input word register -> gate_core ->
// output register. One clock from captured word to registered result.
// PREFIX=0: gate_core (direct transcription, primary H5 design);
// PREFIX=1: gate_core_prefix (same function, parallel-prefix comparators).
`timescale 1ns/1ps
`default_nettype none
module gate_reg #(
    parameter PREFIX = 0
) (
    input  wire         clk,
    input  wire         load,      // capture word_in at this edge
    input  wire [461:0] word_in,
    output reg  [9:0]   result_q   // valid one clock after the word register is loaded
);
    reg  [461:0] word_q;
    wire [9:0]   result_d;

    generate
        if (PREFIX == 0) begin : g_direct
            gate_core core (.word(word_q), .result(result_d));
        end else begin : g_prefix
            gate_core_prefix core (.word(word_q), .result(result_d));
        end
    endgenerate

    initial begin
        word_q   = 462'd0;
        result_q = 10'd0;
    end

    always @(posedge clk) begin
        if (load)
            word_q <= word_in;
        result_q <= result_d;
    end
endmodule
`default_nettype wire
