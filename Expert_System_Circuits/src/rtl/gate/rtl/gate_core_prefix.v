// Separately labelled one-tick variant of gate_core with the same word
// layout and result encoding. Implements PrevalidatedGate (the 9-condition
// postcondition comparator over pre-computed flags), NOT the Horn/normative
// closure engine of kernel K. The 64-bit order checks are computed with
// 5-level parallel-prefix carry trees instead of cascaded ripple chains:
//   a <= b (unsigned)  <=> carry out of b + ~a + 1
//   a <= b (signed)    <=> same with both sign bits inverted
//   now - ts <= max_age, given ts <= now  <=>  max_age + ts + ~now + 1 >= 2^64,
//   evaluated via one 64-bit carry-save level (csa64) and one 5-level
//   prefix tree whose split outputs (age_hi, age_lo) merge in the final LUTs.
`timescale 1ns/1ps
`default_nettype none
(* keep_hierarchy *)
module gate_core_prefix (
    input  wire [461:0] word,
    output wire [9:0]   result
);
    wire [7:0]  evidence_code   = word[13:6];
    wire [31:0] data_revision   = word[45:14];
    wire [31:0] validated_rev   = word[77:46];
    wire [31:0] policy_revision = word[109:78];
    wire [31:0] expected_policy = word[141:110];
    wire [63:0] timestamp       = word[205:142];
    wire [63:0] now             = word[269:206];
    wire [63:0] max_age         = word[333:270];
    wire [63:0] upper_raw       = word[397:334];
    wire [63:0] limit_raw       = word[461:398];

    wire ts_hi, ts_lo, ts_hi_dup, ts_lo_dup;
    prefix_carry #(.W(64), .CIN(1'b1), .INV_B(1'b1), .INV_MSB(1'b0)) u_ts (
        .a_in(now), .b_in(timestamp), .c_top(1'b0),
        .c_hi(ts_hi), .c_lo(ts_lo), .c_hi_dup(ts_hi_dup), .c_lo_dup(ts_lo_dup)
    );

    wire ul_hi, ul_lo, ul_hi_dup, ul_lo_dup;
    prefix_carry #(.W(64), .CIN(1'b1), .INV_B(1'b1), .INV_MSB(1'b1)) u_ul (
        .a_in(limit_raw), .b_in(upper_raw), .c_top(1'b0),
        .c_hi(ul_hi), .c_lo(ul_lo), .c_hi_dup(ul_hi_dup), .c_lo_dup(ul_lo_dup)
    );

    wire [63:0] s, k_shift;
    wire k_top;
    csa64 u_csa (.x(max_age), .y(timestamp), .z_inv(now), .s(s), .k_shift(k_shift), .k_top(k_top));

    wire age_hi, age_lo, age_hi_dup, age_lo_dup;
    prefix_carry #(.W(64), .CIN(1'b0), .INV_B(1'b0), .INV_MSB(1'b0)) u_age (
        .a_in(s), .b_in(k_shift), .c_top(k_top),
        .c_hi(age_hi), .c_lo(age_lo), .c_hi_dup(age_hi_dup), .c_lo_dup(age_lo_dup)
    );

    wire [6:0] ok_low;
    assign ok_low[0] = word[0];
    assign ok_low[1] = word[1] && (data_revision == validated_rev);
    assign ok_low[2] = word[2];
    assign ok_low[3] = word[3];
    assign ok_low[4] = word[4] && !word[5];
    assign ok_low[5] = (evidence_code == 8'd1);
    assign ok_low[6] = (policy_revision == expected_policy);

    wire ok_fast = &ok_low;
    wire ts_and_ul, allow, stale_fail, ul_fail;

    final_ts_ul u_ts_ul (.ts_hi(ts_hi), .ts_lo(ts_lo), .ul_hi(ul_hi), .ul_lo(ul_lo), .ts_and_ul(ts_and_ul));
    final_stale u_stale (.ts_hi(ts_hi_dup), .ts_lo(ts_lo_dup), .age_hi(age_hi_dup), .age_lo(age_lo_dup), .stale_fail(stale_fail));
    final_allow u_allow (.ok_fast(ok_fast), .ts_and_ul(ts_and_ul), .age_hi(age_hi), .age_lo(age_lo), .allow(allow));
    final_ul_fail u_ul_fail (.ul_hi(ul_hi_dup), .ul_lo(ul_lo_dup), .ul_fail(ul_fail));

    assign result = {ul_fail, stale_fail, ~ok_low, allow};
endmodule
`default_nettype wire
