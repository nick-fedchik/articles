// Separately labelled one-tick variant of teos_gate_core with the same word
// layout and result encoding. The four 64-bit order checks are computed with
// parallel-prefix carries instead of ripple chains; the predicate is unchanged:
//   a <= b (unsigned)  <=> carry out of b + ~a + 1
//   a <= b (signed)    <=> same with both sign bits inverted
//   now - ts <= max_age, given ts <= now  <=>  max_age + ts - now >= 0 in
//   66-bit two's complement, evaluated via a carry-save step and one prefix carry.
`timescale 1ns/1ps
`default_nettype none
module teos_gate_core_prefix (
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

    /* verilator lint_off UNUSEDSIGNAL */ // only the top carry is used
    wire [64:0] c_ts;
    /* verilator lint_on UNUSEDSIGNAL */
    teos_prefix_carry #(.W(64)) u_ts (.a(now), .b(~timestamp), .cin(1'b1), .c(c_ts));
    wire ts_le_now = c_ts[64];

    wire [63:0] upper_ob = {~upper_raw[63], upper_raw[62:0]};
    wire [63:0] limit_ob = {~limit_raw[63], limit_raw[62:0]};
    /* verilator lint_off UNUSEDSIGNAL */ // only the top carry is used
    wire [64:0] c_ul;
    /* verilator lint_on UNUSEDSIGNAL */
    teos_prefix_carry #(.W(64)) u_ul (.a(limit_ob), .b(~upper_ob), .cin(1'b1), .c(c_ul));
    wire upper_le_limit = c_ul[64];

    // Carry-save: x + y + z with x = max_age, y = timestamp, z = ~now (66-bit).
    wire [65:0] x = {2'b00, max_age};
    wire [65:0] y = {2'b00, timestamp};
    wire [65:0] z = ~{2'b00, now};
    wire [65:0] s = x ^ y ^ z;
    wire [65:0] k = {((x[64:0] & y[64:0]) | (x[64:0] & z[64:0]) | (y[64:0] & z[64:0])), 1'b0};
    /* verilator lint_off UNUSEDSIGNAL */ // only the top carry is used
    wire [65:0] c_age;
    /* verilator lint_on UNUSEDSIGNAL */
    teos_prefix_carry #(.W(65)) u_age (.a(s[64:0]), .b(k[64:0]), .cin(1'b1), .c(c_age));
    wire age_sign = s[65] ^ k[65] ^ c_age[65];
    wire age_le_max = ~age_sign;

    wire [8:0] ok;
    assign ok[0] = word[0];
    assign ok[1] = word[1] && (data_revision == validated_rev);
    assign ok[2] = word[2];
    assign ok[3] = word[3];
    assign ok[4] = word[4] && !word[5];
    assign ok[5] = (evidence_code == 8'd1);
    assign ok[6] = (policy_revision == expected_policy);
    assign ok[7] = ts_le_now && age_le_max;
    assign ok[8] = upper_le_limit;

    assign result = {~ok, &ok};
endmodule
`default_nettype wire
