package main

import (
	"bytes"
	"math"
	"strings"
	"testing"

	"github.com/nick-fedchik/articles/Expert_System_Circuits/src/hardware"
)

func TestLayoutIsContiguous(t *testing.T) {
	next := 0
	for _, f := range Layout {
		if f.lsb != next || f.width <= 0 {
			t.Fatalf("field %s at %d width %d, want lsb %d", f.name, f.lsb, f.width, next)
		}
		next += f.width
	}
	if next != WordBits {
		t.Fatalf("layout covers %d bits, want %d", next, WordBits)
	}
}

func TestPackUnpackRoundTrip(t *testing.T) {
	inputs := []hardware.GateInput{Base(), {}, {
		CustodyValidated: true, Undercut: true, Evidence: 0xff,
		DataRevision: math.MaxUint32, ValidatedRevision: 1, PolicyRevision: 1 << 31, ExpectedPolicy: 5,
		Timestamp: math.MaxUint64, Now: 1, MaxAge: 1 << 63,
		UpperRaw: math.MinInt64, LimitRaw: -1,
	}}
	for _, in := range inputs {
		if got := Unpack(Pack(in)); got != in {
			t.Fatalf("round trip %+v -> %+v", in, got)
		}
	}
}

func TestExpectedMatchesGate(t *testing.T) {
	if got, err := Expected(Base()); err != nil || got != 1 {
		t.Fatalf("base expected 1, got %d %v", got, err)
	}
	all := hardware.GateInput{Undercut: true, DataRevision: 1, PolicyRevision: 1, Timestamp: 2, Now: 1, UpperRaw: 1}
	got, err := Expected(all)
	if err != nil {
		t.Fatal(err)
	}
	if want := uint16(1<<ResultBits - 2); got != want {
		t.Fatalf("all checks failing: got %#x want %#x", got, want)
	}
}

func TestCeilRawMatchesReportExample(t *testing.T) {
	cases := map[string]int64{"10.4": 681_575, "9.5": 622_592, "10": 655_360, "-0.5": -32_768, "-0.0000001": 0, "0.0000001": 1}
	for decimal, want := range cases {
		if got, err := CeilRaw(decimal); err != nil || got != want {
			t.Fatalf("CeilRaw(%s) = %d, %v; want %d", decimal, got, err, want)
		}
	}
	if _, err := CeilRaw("1e30"); err == nil {
		t.Fatal("overflow accepted")
	}
}

func TestRenderDeterministic(t *testing.T) {
	var hex, summary [2]bytes.Buffer
	for i := range hex {
		vectors, err := Vectors(20261103, 50)
		if err != nil {
			t.Fatal(err)
		}
		if _, err := render(vectors, &hex[i], &summary[i]); err != nil {
			t.Fatal(err)
		}
	}
	if !bytes.Equal(hex[0].Bytes(), hex[1].Bytes()) || summary[0].String() != summary[1].String() {
		t.Fatal("generator is not deterministic")
	}
	for _, line := range strings.Split(strings.TrimSpace(hex[0].String()), "\n") {
		if len(line) != hexDigits {
			t.Fatalf("line length %d, want %d", len(line), hexDigits)
		}
	}
	if !strings.Contains(summary[0].String(), `"allow"`) {
		t.Fatalf("summary missing allow count: %s", summary[0].String())
	}
}

func TestRunRejectsArguments(t *testing.T) {
	var out, diag bytes.Buffer
	if err := run(&out, &diag, []string{"extra"}); err == nil {
		t.Fatal("positional argument accepted")
	}
	if err := run(&out, &diag, []string{"-random", "-1"}); err == nil {
		t.Fatal("negative count accepted")
	}
}
