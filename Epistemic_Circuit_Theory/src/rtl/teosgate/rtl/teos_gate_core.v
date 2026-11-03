// Combinational PREVALIDATED gate predicate. Bit-exact port of
// teoshardware.PrevalidatedGate (src/teoshardware/calculate.go); no policy added.
// Word layout (LSB first) matches cmd/teosgatevectors Layout:
//   [0] custody_validated  [1] binding_valid  [2] si_match  [3] numeric_valid
//   [4] norm_permitted     [5] undercut       [13:6] evidence_code (uint8)
//   [45:14] data_revision  [77:46] validated_revision
//   [109:78] policy_revision  [141:110] expected_policy_revision
//   [205:142] timestamp_ticks [269:206] now_ticks [333:270] max_age_ticks (uint64)
//   [397:334] prevalidated_upper_bound_raw [461:398] limit_raw (int64, two's complement)
// result[0] = allow; result[i+1] = reason i failed, in Go check order.
// Not a custody/SHA-256 validator: inputs must already be bound by trusted validation.
`timescale 1ns/1ps
`default_nettype none
module teos_gate_core (
    input  wire [461:0] word,
    output wire [9:0]   result
);
    wire        custody_validated = word[0];
    wire        binding_valid     = word[1];
    wire        si_match          = word[2];
    wire        numeric_valid     = word[3];
    wire        norm_permitted    = word[4];
    wire        undercut          = word[5];
    wire [7:0]  evidence_code     = word[13:6];
    wire [31:0] data_revision     = word[45:14];
    wire [31:0] validated_rev     = word[77:46];
    wire [31:0] policy_revision   = word[109:78];
    wire [31:0] expected_policy   = word[141:110];
    wire [63:0] timestamp         = word[205:142];
    wire [63:0] now               = word[269:206];
    wire [63:0] max_age           = word[333:270];
    wire signed [63:0] upper_raw  = word[397:334];
    wire signed [63:0] limit_raw  = word[461:398];

    // Now - Timestamp is only meaningful when Timestamp <= Now (Go guards with &&).
    wire [63:0] age = now - timestamp;

    wire [8:0] ok;
    assign ok[0] = custody_validated;
    assign ok[1] = binding_valid && (data_revision == validated_rev);
    assign ok[2] = si_match;
    assign ok[3] = numeric_valid;
    assign ok[4] = norm_permitted && !undercut;
    assign ok[5] = (evidence_code == 8'd1);
    assign ok[6] = (policy_revision == expected_policy);
    assign ok[7] = (timestamp <= now) && (age <= max_age);
    assign ok[8] = (upper_raw <= limit_raw);

    assign result = {~ok, &ok};
endmodule
`default_nettype wire
