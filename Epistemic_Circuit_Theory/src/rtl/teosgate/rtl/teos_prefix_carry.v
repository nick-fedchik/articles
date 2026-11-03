// Kogge-Stone parallel-prefix carry: c[i] is the carry into bit i of
// a + b + cin (c[0] = cin, c[W] = carry out). Logarithmic depth instead of
// the iCE40 ripple SB_CARRY chain; used only by teos_gate_core_prefix.
`timescale 1ns/1ps
`default_nettype none
module teos_prefix_carry #(
    parameter integer W = 64
) (
    input  wire [W-1:0] a,
    input  wire [W-1:0] b,
    input  wire         cin,
    output wire [W:0]   c
);
    localparam integer N = W + 1;          // position 0 carries cin
    localparam integer L = $clog2(N);

    // Behavioural loops (unrolled by synthesis) keep event-driven simulators
    // from re-evaluating every bit-select assignment of a shared vector.
    reg [N-1:0] g, p, g_next, p_next;
    integer i, l;
    always @(*) begin
        g[0] = cin;
        p[0] = 1'b0;
        for (i = 1; i < N; i = i + 1) begin
            g[i] = a[i-1] & b[i-1];
            p[i] = a[i-1] ^ b[i-1];
        end
        for (l = 0; l < L; l = l + 1) begin
            for (i = 0; i < N; i = i + 1) begin
                if (i >= (1 << l)) begin
                    g_next[i] = g[i] | (p[i] & g[i - (1 << l)]);
                    p_next[i] = p[i] & p[i - (1 << l)];
                end else begin
                    g_next[i] = g[i];
                    p_next[i] = p[i];
                end
            end
            g = g_next;
            p = p_next;
        end
    end
    // Group generate over positions [i:0] is the carry into original bit i.
    assign c = g;
endmodule
`default_nettype wire
