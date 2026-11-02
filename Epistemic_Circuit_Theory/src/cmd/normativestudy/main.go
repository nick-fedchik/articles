package main

import (
	"context"
	"encoding/csv"
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"reflect"
	"slices"
	"strconv"
	"strings"
	"time"

	evidence "github.com/nick-fedchik/articles/Epistemic_Circuit_Theory/src"
)

type semanticRow struct {
	Stage, Case, Action, Facts, Closure string
	Decision                            evidence.NormDecision
	Raw, Applicable, Surviving          evidence.State
	SAT                                 bool
}

type cacheRow struct {
	Sequence, Step, Repeat  int
	Expected, Naive, Actual evidence.NormDecision
	Reused, Invalidated     bool
}

type report struct {
	Stages                   map[string]int
	SATChecks                int
	Decisions                map[evidence.NormDecision]int
	CacheEvents, CacheChecks int
	CacheHits, Invalidations int
	NaiveUnsafeEvents        int
	OppositePotentialAllows  int
	PriorityResolvedAllows   int
	Scope                    string
}

type study struct {
	ctx       context.Context
	solver    string
	rows      []semanticRow
	cacheRows []cacheRow
	report    report
}

func statesString(states []evidence.State) string {
	var out strings.Builder
	for _, state := range states {
		out.WriteByte(byte('0' + state))
	}
	return out.String()
}

func (s *study) compare(p *evidence.NormProgram, w evidence.NormWorld, action, stage, id string, useSAT bool) (evidence.NormCertificate, error) {
	certificate, err := p.Evaluate(w, action)
	if err != nil {
		return evidence.NormCertificate{}, err
	}
	closure, reference, err := p.Reference(w, action)
	if err != nil || !slices.Equal(closure, certificate.Closure) || !reflect.DeepEqual(reference, certificate.Result) {
		return evidence.NormCertificate{}, fmt.Errorf("reference mismatch %s/%s: %v", stage, id, err)
	}
	if useSAT {
		sat, err := p.SATReference(s.ctx, s.solver, w, action)
		if err != nil || !slices.Equal(sat.Closure, closure) || sat.Decision != reference.Decision ||
			sat.ConditionEvidence != reference.ConditionEvidence || sat.Applicable != reference.Applicable || sat.Surviving != reference.Surviving {
			return evidence.NormCertificate{}, fmt.Errorf("SAT mismatch %s/%s: %v", stage, id, err)
		}
		s.report.SATChecks++
	}
	if certificate.Result.Decision == evidence.NormAllow {
		if certificate.Result.ConditionEvidence == evidence.Both {
			s.report.OppositePotentialAllows++
			if certificate.Result.Applicable == evidence.Both && certificate.Result.Surviving == evidence.True {
				s.report.PriorityResolvedAllows++
			}
		}
		for _, rule := range certificate.Result.Rules {
			if strings.HasPrefix(rule.Status, "UNKNOWN_") || strings.HasPrefix(rule.Status, "CONFLICTING_") {
				return evidence.NormCertificate{}, fmt.Errorf("unsafe allowance %s/%s", stage, id)
			}
		}
	}
	s.rows = append(s.rows, semanticRow{stage, id, action, statesString(w.Facts), statesString(closure),
		reference.Decision, reference.ConditionEvidence, reference.Applicable, reference.Surviving, useSAT})
	s.report.Stages[stage]++
	s.report.Decisions[reference.Decision]++
	return certificate, nil
}

func (s *study) signedStage() error {
	theory := evidence.NormTheory{
		Atoms: 4,
		Evidence: []evidence.SignedRule{
			{ID: "positive", Version: "1", Premises: []evidence.SignedLiteral{{Atom: 0}}, Conclusion: evidence.SignedLiteral{Atom: 1}},
			{ID: "negative", Version: "1", Premises: []evidence.SignedLiteral{{Atom: 0, Negative: true}}, Conclusion: evidence.SignedLiteral{Atom: 1, Negative: true}},
		},
		Norms: []evidence.NormRule{
			{ID: "permit", Version: "1", Action: "run", Effect: evidence.Permit, Conditions: []evidence.SignedLiteral{{Atom: 1}}},
			{ID: "ban", Version: "1", Action: "run", Effect: evidence.Prohibit, Conditions: []evidence.SignedLiteral{{Atom: 1, Negative: true}}},
		},
	}
	p, err := evidence.CompileNormative(theory)
	if err != nil {
		return err
	}
	for raw := 0; raw < 256; raw++ {
		facts := []evidence.State{evidence.State(raw & 3), evidence.State(raw >> 2 & 3),
			evidence.State(raw >> 4 & 3), evidence.State(raw >> 6 & 3)}
		world, err := p.NewWorld(facts)
		if err != nil {
			return err
		}
		certificate, err := s.compare(p, world, "run", "signed", strconv.Itoa(raw), raw%4 == raw/16%4)
		if err != nil {
			return err
		}
		if certificate.Closure[2] != facts[2] || certificate.Closure[3] != facts[3] {
			return fmt.Errorf("unrelated atom was derived by explosion")
		}
	}
	return nil
}

func (s *study) exceptionStage() error {
	for _, withException := range []bool{false, true} {
		for premise := evidence.Unknown; premise <= evidence.Both; premise++ {
			for guard := evidence.Unknown; guard <= evidence.Both; guard++ {
				for _, banEnabled := range []bool{false, true} {
					var raw, applicable evidence.State
					for priority := 0; priority < 3; priority++ {
						theory := simpleTheory()
						if !withException {
							theory.Norms[0].Exception = nil
						}
						setPriority(&theory, priority)
						p, err := evidence.CompileNormative(theory)
						if err != nil {
							return err
						}
						world, err := p.NewWorld([]evidence.State{premise, guard, evidence.Unknown, evidence.Unknown})
						if err != nil {
							return err
						}
						world.NormEnabled[1] = banEnabled
						id := fmt.Sprintf("exception=%t/p=%d/g=%d/ban=%t/order=%d", withException, premise, guard, banEnabled, priority)
						certificate, err := s.compare(p, world, "run", "exceptions", id, true)
						if err != nil {
							return err
						}
						if priority == 0 {
							raw, applicable = certificate.Result.ConditionEvidence, certificate.Result.Applicable
						} else if certificate.Result.ConditionEvidence != raw || certificate.Result.Applicable != applicable {
							return fmt.Errorf("priority erased contrary normative evidence")
						}
					}
				}
			}
		}
	}
	return nil
}

func (s *study) cacheStage() error {
	for sequence := 0; sequence < 16; sequence++ {
		cache := evidence.NewNormCache()
		naive := evidence.NormAllow
		for step := 0; step < 64; step++ {
			theory := simpleTheory()
			facts := []evidence.State{evidence.True, evidence.False, evidence.Unknown, evidence.State(sequence % 4)}
			ban := false
			switch step % 8 {
			case 0:
				// The declared exception is explicitly false in the initial facts.
			case 1:
				facts[1] = evidence.Unknown
			case 2:
				facts[1] = evidence.Both
			case 3:
				facts[1] = evidence.True
			case 4:
				ban = true
			case 5:
				ban = true
				setPriority(&theory, 1)
			case 6:
				ban = true
				setPriority(&theory, 2)
			case 7:
				ban = true
				facts[0] = evidence.Both
				setPriority(&theory, 1)
			}
			p, err := evidence.CompileNormative(theory)
			if err != nil {
				return err
			}
			world, err := p.NewWorld(facts)
			if err != nil {
				return err
			}
			world.Revision = uint64(1 + step)
			world.NormEnabled[1] = ban
			_, expected, err := p.Reference(world, "run")
			if err != nil {
				return err
			}
			if step == 0 && expected.Decision != evidence.NormAllow {
				return fmt.Errorf("naive baseline was not initialized by a valid allowance")
			}
			// Intentionally unsafe comparison: the positive prerequisite remains
			// supported, while guards, opposing norms and priorities are ignored.
			if facts[0]&2 == 0 {
				return fmt.Errorf("positive-only cache witness lost its positive prerequisite")
			}
			if naive == evidence.NormAllow && expected.Decision != evidence.NormAllow {
				s.report.NaiveUnsafeEvents++
			}
			s.report.CacheEvents++
			for repeat := 0; repeat < 2; repeat++ {
				result, err := cache.Resolve(p, world, "run")
				if err != nil || !reflect.DeepEqual(result.Certificate.Result, expected) {
					return fmt.Errorf("context cache mismatch: %v", err)
				}
				if repeat == 1 && !result.Reused {
					return fmt.Errorf("identical normative context did not reuse the certificate")
				}
				if result.Reused {
					s.report.CacheHits++
				}
				if result.Invalidated {
					s.report.Invalidations++
				}
				s.report.CacheChecks++
				s.cacheRows = append(s.cacheRows, cacheRow{sequence, step, repeat, expected.Decision, naive,
					result.Certificate.Result.Decision, result.Reused, result.Invalidated})
			}
		}
	}
	return nil
}

func writeCSV(dir, name string, header []string, rows [][]string) error {
	if dir == "" {
		return nil
	}
	if err := os.MkdirAll(dir, 0755); err != nil {
		return err
	}
	file, err := os.Create(filepath.Join(dir, name))
	if err != nil {
		return err
	}
	w := csv.NewWriter(file)
	if err := w.Write(header); err != nil {
		file.Close()
		return err
	}
	for _, row := range rows {
		if err := w.Write(row); err != nil {
			file.Close()
			return err
		}
	}
	w.Flush()
	if err := w.Error(); err != nil {
		file.Close()
		return err
	}
	return file.Close()
}

func (s *study) write(dir string) error {
	rows := [][]string{}
	for _, row := range s.rows {
		rows = append(rows, []string{row.Stage, row.Case, row.Action, row.Facts, row.Closure, string(row.Decision),
			strconv.Itoa(int(row.Raw)), strconv.Itoa(int(row.Applicable)), strconv.Itoa(int(row.Surviving)), strconv.FormatBool(row.SAT)})
	}
	if err := writeCSV(dir, "semantics.csv", []string{"stage", "case", "action", "fact_states", "closure_states", "decision",
		"condition_evidence", "applicable_evidence", "surviving_evidence", "sat_checked"}, rows); err != nil {
		return err
	}
	rows = nil
	for _, row := range s.cacheRows {
		rows = append(rows, []string{strconv.Itoa(row.Sequence), strconv.Itoa(row.Step), strconv.Itoa(row.Repeat),
			string(row.Expected), string(row.Naive), string(row.Actual), strconv.FormatBool(row.Reused), strconv.FormatBool(row.Invalidated)})
	}
	return writeCSV(dir, "context-cache.csv", []string{"sequence", "step", "repeat", "expected", "unsafe_positive_only",
		"actual", "reused", "invalidated"}, rows)
}

func run(ctx context.Context, solver, dir string, writer io.Writer) error {
	if _, err := evidence.IdentifySolver(ctx, solver); err != nil {
		return fmt.Errorf("required normative SAT oracle: %w", err)
	}
	s := &study{ctx: ctx, solver: solver, report: report{
		Stages: map[string]int{}, Decisions: map[evidence.NormDecision]int{},
		Scope: "bounded signed Horn evidence and conservative normative decisions; synthetic models, not ASPIC+ or independent text evaluation",
	}}
	for _, stage := range []func() error{s.signedStage, s.exceptionStage, s.randomStage, s.cacheStage} {
		if err := stage(); err != nil {
			return err
		}
	}
	if err := s.write(dir); err != nil {
		return err
	}
	encoder := json.NewEncoder(writer)
	encoder.SetIndent("", "  ")
	return encoder.Encode(s.report)
}

func main() {
	solver := flag.String("solver", "minisat", "required external Minisat-compatible solver")
	dir := flag.String("csv-dir", "", "explicit directory for two new primary CSV files; existing files are replaced")
	flag.Parse()
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Minute)
	defer cancel()
	if err := run(ctx, *solver, *dir, os.Stdout); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}
