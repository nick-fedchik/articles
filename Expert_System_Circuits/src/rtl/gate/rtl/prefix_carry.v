// 5-LUT4-level parallel-prefix carry tree for 64-bit order checks in
// gate_core_prefix. Pairs of bits are packed into 4-input LUTs at level 1
// (32 groups), reduced to 4 groups at level 4, and split into upper/lower
// halves (c_hi, c_lo) at level 5 using a radix-4 32-bit upper group propagate
// (p_hi32 = &p0[31:16]) so the final OR c_hi | c_lo merges directly into the
// output 4-LUTs of result[0] and result[8] (total core depth = 7 LUT4 levels).
`timescale 1ns/1ps
`default_nettype none
/* verilator lint_off DECLFILENAME */
(* keep_hierarchy *)
module prefix_carry #(
    parameter integer W = 64,
    parameter [0:0]   CIN = 1'b0,
    parameter [0:0]   INV_B = 1'b0,
    parameter [0:0]   INV_MSB = 1'b0
) (
    input  wire [W-1:0] a_in,
    input  wire [W-1:0] b_in,
    input  wire         c_top,
    output wire         c_hi,
    output wire         c_lo,
    output wire         c_hi_dup,
    output wire         c_lo_dup
);
    wire [W-1:0] a = INV_MSB ? {~a_in[W-1], a_in[W-2:0]} : a_in;
    wire [W-1:0] b_msb = INV_MSB ? {~b_in[W-1], b_in[W-2:0]} : b_in;
    wire [W-1:0] b = INV_B ? ~b_msb : b_msb;

    wire [31:0] g0, p0;
    genvar m;
    generate
        for (m = 0; m < 32; m = m + 1) begin : l0
            if (m == 0 && CIN == 1'b1) begin : g_cin1
                assign g0[m] = (a[1] & b[1]) | ((a[1] | b[1]) & (a[0] | b[0]));
                assign p0[m] = 1'b0;
            end else begin : g_cin0
                assign g0[m] = (a[2*m+1] & b[2*m+1]) | ((a[2*m+1] | b[2*m+1]) & a[2*m] & b[2*m]);
                assign p0[m] = (a[2*m+1] | b[2*m+1]) & (a[2*m] | b[2*m]);
            end
        end
        wire [15:0] g1, p1;
        for (m = 0; m < 16; m = m + 1) begin : l1
            assign g1[m] = g0[2*m+1] | (p0[2*m+1] & g0[2*m]);
            assign p1[m] = p0[2*m+1] & p0[2*m];
        end
        wire [7:0] g2, p2;
        for (m = 0; m < 8; m = m + 1) begin : l2
            assign g2[m] = g1[2*m+1] | (p1[2*m+1] & g1[2*m]);
            assign p2[m] = p1[2*m+1] & p1[2*m];
        end
        /* verilator lint_off UNUSEDSIGNAL */ // p3[2] and p3[0] are unused at top split
        wire [3:0] g3, p3;
        /* verilator lint_on UNUSEDSIGNAL */
        for (m = 0; m < 4; m = m + 1) begin : l3
            assign g3[m] = g2[2*m+1] | (p2[2*m+1] & g2[2*m]);
            assign p3[m] = p2[2*m+1] & p2[2*m];
        end
    endgenerate
    wire p_hi32 = &p0[31:16];
    prefix_c_hi u_chi0 (.c_top(c_top), .g3_3(g3[3]), .p3_3(p3[3]), .g3_2(g3[2]), .c_hi(c_hi));
    prefix_c_lo u_clo0 (.p_hi32(p_hi32), .g3_1(g3[1]), .p3_1(p3[1]), .g3_0(g3[0]), .c_lo(c_lo));
    prefix_c_hi u_chi1 (.c_top(c_top), .g3_3(g3[3]), .p3_3(p3[3]), .g3_2(g3[2]), .c_hi(c_hi_dup));
    prefix_c_lo u_clo1 (.p_hi32(p_hi32), .g3_1(g3[1]), .p3_1(p3[1]), .g3_0(g3[0]), .c_lo(c_lo_dup));
endmodule

(* keep_hierarchy *)
module prefix_c_hi (
    input  wire c_top, g3_3, p3_3, g3_2,
    output wire c_hi
);
    assign c_hi = c_top | g3_3 | (p3_3 & g3_2);
endmodule

(* keep_hierarchy *)
module prefix_c_lo (
    input  wire p_hi32, g3_1, p3_1, g3_0,
    output wire c_lo
);
    assign c_lo = p_hi32 & (g3_1 | (p3_1 & g3_0));
endmodule

(* keep_hierarchy *)
module csa64 (
    input  wire [63:0] x,
    input  wire [63:0] y,
    input  wire [63:0] z_inv,
    output wire [63:0] s,
    output wire [63:0] k_shift,
    output wire        k_top
);
    wire [63:0] z = ~z_inv;
    wire [63:0] k = (x & y) | (x & z) | (y & z);
    assign s = x ^ y ^ z;
    assign k_shift = {k[62:0], 1'b1};
    assign k_top = k[63];
endmodule

(* keep_hierarchy *)
module final_ts_ul (
    input  wire ts_hi,
    input  wire ts_lo,
    input  wire ul_hi,
    input  wire ul_lo,
    output wire ts_and_ul
);
    assign ts_and_ul = (ts_hi | ts_lo) & (ul_hi | ul_lo);
endmodule

(* keep_hierarchy *)
module final_allow (
    input  wire age_lo,
    input  wire age_hi,
    input  wire ts_and_ul,
    input  wire ok_fast,
    output wire allow
);
    assign allow = ok_fast & ts_and_ul & (age_hi | age_lo);
endmodule

(* keep_hierarchy *)
module final_stale (
    input  wire age_lo,
    input  wire age_hi,
    input  wire ts_lo,
    input  wire ts_hi,
    output wire stale_fail
);
    assign stale_fail = ~((ts_hi | ts_lo) & (age_hi | age_lo));
endmodule

(* keep_hierarchy *)
module final_ul_fail (
    input  wire ul_hi,
    input  wire ul_lo,
    output wire ul_fail
);
    assign ul_fail = ~(ul_hi | ul_lo);
endmodule
/* verilator lint_on DECLFILENAME */
`default_nettype wire
