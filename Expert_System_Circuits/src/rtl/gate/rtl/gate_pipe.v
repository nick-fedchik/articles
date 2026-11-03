// Separately labelled pipelined variant (NOT the one-tick design under H5):
// word register -> stage 1 (checks 0..6, 8; age = now - timestamp; timestamp <= now)
// -> stage 2 (age <= max_age, allow reduction) -> result register.
// Latency: 2 clocks after the word register is loaded; initiation interval 1.
// Same function as gate_core; verified against the same Go vectors.
`timescale 1ns/1ps
`default_nettype none
module gate_pipe (
    input  wire         clk,
    input  wire         load,
    input  wire [461:0] word_in,
    output reg  [9:0]   result_q
);
    reg  [461:0] word_q;

    wire [7:0]  evidence_code    = word_q[13:6];
    wire [31:0] data_revision    = word_q[45:14];
    wire [31:0] validated_rev    = word_q[77:46];
    wire [31:0] policy_revision  = word_q[109:78];
    wire [31:0] expected_policy  = word_q[141:110];
    wire [63:0] timestamp        = word_q[205:142];
    wire [63:0] now              = word_q[269:206];
    wire [63:0] max_age          = word_q[333:270];
    wire signed [63:0] upper_raw = word_q[397:334];
    wire signed [63:0] limit_raw = word_q[461:398];

    reg  [8:0]  ok_s1;      // bit 7 holds (timestamp <= now) until stage 2
    reg  [63:0] age_s1;
    reg  [63:0] max_age_s1;
    wire [8:0]  ok_s2 = {ok_s1[8], ok_s1[7] && (age_s1 <= max_age_s1), ok_s1[6:0]};

    initial begin
        word_q = 462'd0; ok_s1 = 9'd0; age_s1 = 64'd0; max_age_s1 = 64'd0; result_q = 10'd0;
    end

    always @(posedge clk) begin
        if (load)
            word_q <= word_in;
        ok_s1[0] <= word_q[0];
        ok_s1[1] <= word_q[1] && (data_revision == validated_rev);
        ok_s1[2] <= word_q[2];
        ok_s1[3] <= word_q[3];
        ok_s1[4] <= word_q[4] && !word_q[5];
        ok_s1[5] <= (evidence_code == 8'd1);
        ok_s1[6] <= (policy_revision == expected_policy);
        ok_s1[7] <= (timestamp <= now);
        ok_s1[8] <= (upper_raw <= limit_raw);
        age_s1     <= now - timestamp;
        max_age_s1 <= max_age;
        result_q   <= {~ok_s2, &ok_s2};
    end
endmodule
`default_nettype wire
