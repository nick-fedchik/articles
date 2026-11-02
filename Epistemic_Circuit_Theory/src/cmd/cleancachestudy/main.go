package main

import (
	"encoding/csv"
	"encoding/json"
	"errors"
	"flag"
	"fmt"
	"os"
	"path/filepath"
	"runtime"
	"strconv"
	"time"

	e "github.com/nick-fedchik/articles/Epistemic_Circuit_Theory/src"
)

var modes = []e.CleanCacheMode{e.CleanRecompute, e.CleanRevision, e.CleanCone}

type report struct {
	Events, Queries, CheckedAnswers, NaiveWrongClean int
	Reused                                           map[e.CleanCacheMode]int
	Measurements                                     int
	Go, Scope                                        string
}

type study struct {
	report          report
	events, timings [][]string
}

func n(value int) string  { return strconv.Itoa(value) }
func b(value bool) string { return strconv.FormatBool(value) }
func code(query e.SignedLiteral) int {
	value := 2 * query.Atom
	if query.Negative {
		value++
	}
	return value
}
func rule(id string, conclusion int, premises ...int) e.SignedRule {
	body := make([]e.SignedLiteral, len(premises))
	for i, value := range premises {
		body[i] = e.SignedLiteral{Atom: value / 2, Negative: value%2 != 0}
	}
	return e.SignedRule{ID: id, Version: "1", Premises: body, Conclusion: e.SignedLiteral{Atom: conclusion / 2, Negative: conclusion%2 != 0}}
}
func positiveStillValid(world e.SignedWorld, witness e.CleanWitnessCertificate) bool {
	for _, literal := range witness.Proof.Facts {
		state := world.Facts[literal/2]
		if literal%2 == 0 && state != e.True && state != e.Both || literal%2 != 0 && state != e.False && state != e.Both {
			return false
		}
	}
	for _, index := range witness.Proof.Rules {
		if !world.Enabled[index] {
			return false
		}
	}
	return true
}

func (s *study) correctness() error {
	rules := []e.SignedRule{rule("derive", 4, 0), rule("opposition", 1, 2), rule("alternative", 4, 6), rule("trigger", 2, 8)}
	p, err := e.CompileCleanSupport(5, rules)
	if err != nil {
		return err
	}
	updates := []e.HornUpdate{
		{Facts: []e.HornChange{{Index: 8, Present: true}}}, {Facts: []e.HornChange{{Index: 8, Present: false}}},
		{Rules: []e.HornChange{{Index: 1, Present: false}}}, {Rules: []e.HornChange{{Index: 1, Present: true}}},
		{Facts: []e.HornChange{{Index: 0, Present: false}}}, {Facts: []e.HornChange{{Index: 0, Present: true}}},
		{Facts: []e.HornChange{{Index: 1, Present: true}}}, {Facts: []e.HornChange{{Index: 1, Present: false}}},
		{Facts: []e.HornChange{{Index: 2, Present: true}}}, {Facts: []e.HornChange{{Index: 2, Present: false}}},
		{Facts: []e.HornChange{{Index: 2, Present: true}}}, {Facts: []e.HornChange{{Index: 2, Present: false}}},
	}
	queries := []e.SignedLiteral{{Atom: 2}, {Atom: 2, Negative: true}, {Atom: 0}}
	for initial := 0; initial < 256; initial++ {
		world := e.SignedWorld{Revision: 1, Facts: []e.State{e.State(initial & 3), e.State(initial >> 2 & 3), e.Unknown, e.State(initial >> 4 & 3), e.Unknown}, Enabled: []bool{true, initial&64 != 0, initial&128 != 0, true}}
		caches := []*e.CleanWitnessCache{}
		for _, mode := range modes {
			cache, err := e.NewCleanWitnessCache(p, world, mode)
			if err != nil {
				return err
			}
			caches = append(caches, cache)
			for _, query := range queries {
				if _, err := cache.Resolve(query); err != nil {
					return err
				}
			}
		}
		naive := map[int]e.CleanWitnessCertificate{}
		for _, query := range queries {
			answer, err := p.FreshCleanAnswer(world, query)
			if err != nil {
				return err
			}
			if answer.Witness != nil {
				naive[code(query)] = *answer.Witness
			}
		}
		for step, update := range updates {
			for _, cache := range caches {
				if err := cache.Update(update); err != nil {
					return err
				}
			}
			current, err := caches[0].World()
			if err != nil {
				return err
			}
			s.report.Events++
			for _, query := range queries {
				fresh, _, err := p.EvaluateClean(current, query)
				if err != nil {
					return err
				}
				if err := p.VerifyClean(current, query, fresh); err != nil {
					return err
				}
				ref, err := p.ReferenceClean(current, query, 100000)
				if err != nil {
					return err
				}
				if ref.ExistsClean != (fresh.Status == e.CleanSupported) {
					return errors.New("independent oracle mismatch")
				}
				s.report.Queries++
				claim := fresh.Status == e.CleanSupported
				if old, exists := naive[code(query)]; exists && positiveStillValid(current, old) {
					claim = true
				} else if fresh.Status == e.CleanSupported {
					answer, err := p.FreshCleanAnswer(current, query)
					if err != nil {
						return err
					}
					naive[code(query)] = *answer.Witness
				} else {
					delete(naive, code(query))
				}
				if claim && fresh.Status != e.CleanSupported {
					s.report.NaiveWrongClean++
				}
				for i, cache := range caches {
					result, err := cache.Resolve(query)
					if err != nil {
						return err
					}
					if result.Answer.Status != fresh.Status {
						return errors.New("cache/fresh status mismatch")
					}
					if err := p.VerifyCleanAnswer(current, query, result.Answer); err != nil {
						return err
					}
					s.report.CheckedAnswers++
					if result.Reused {
						s.report.Reused[modes[i]]++
					}
					s.events = append(s.events, []string{n(initial), n(step), n(code(query)), string(modes[i]), string(fresh.Status), b(claim), b(result.Reused), result.Invalidation, n(result.FactDependencies), n(result.RuleDependencies)})
				}
			}
		}
	}
	return nil
}

func workload(size int, broad bool) (*e.CleanProgram, e.SignedWorld, error) {
	rules := []e.SignedRule{}
	for atom := 0; atom < size-1; atom++ {
		rules = append(rules, rule(fmt.Sprintf("r%d", atom), 2*(atom+1), 2*atom))
	}
	if broad {
		rules = append(rules, rule("wide-opposition", 1, 2*(size-1), 2*size))
	}
	p, err := e.CompileCleanSupport(size+1, rules)
	world := e.SignedWorld{Revision: 1, Facts: make([]e.State, size+1), Enabled: make([]bool, len(rules))}
	world.Facts[0] = e.True
	for i := range world.Enabled {
		world.Enabled[i] = true
	}
	return p, world, err
}

func (s *study) measure(p *e.CleanProgram, world e.SignedWorld, family string, size, per, repeat int, mode e.CleanCacheMode) error {
	start := time.Now()
	setupStart := time.Now()
	cache, err := e.NewCleanWitnessCache(p, world, mode)
	if err != nil {
		return err
	}
	setupNS := time.Since(setupStart).Nanoseconds()
	var updateNS, queryNS, verifyNS int64
	reused, invalidated, queries, maxFacts, maxRules, maxAbsence := 0, 0, 0, 0, 0, 0
	resolve := func(atom int) error {
		query := e.SignedLiteral{Atom: atom}
		begin := time.Now()
		result, err := cache.Resolve(query)
		queryNS += time.Since(begin).Nanoseconds()
		if err != nil {
			return err
		}
		begin = time.Now()
		current, err := cache.World()
		if err != nil {
			return err
		}
		if err := p.VerifyCleanAnswer(current, query, result.Answer); err != nil {
			return err
		}
		verifyNS += time.Since(begin).Nanoseconds()
		if result.Answer.Status != e.CleanSupported {
			return errors.New("constructed performance query lost clean support")
		}
		queries++
		if result.Reused {
			reused++
		}
		if result.Invalidation != "" {
			invalidated++
		}
		maxFacts = max(maxFacts, result.FactDependencies)
		maxRules = max(maxRules, result.RuleDependencies)
		maxAbsence = max(maxAbsence, len(result.Answer.Witness.Absence))
		return nil
	}
	for _, atom := range []int{7, 15, 23, 31} {
		if err := resolve(atom); err != nil {
			return err
		}
	}
	for step := 0; step < 16; step++ {
		index := 2*(size-1) + 1
		if family == "broad" {
			index = 2 * (size / 2)
		}
		begin := time.Now()
		err := cache.Update(e.HornUpdate{Facts: []e.HornChange{{Index: index, Present: step%2 == 0}}})
		updateNS += time.Since(begin).Nanoseconds()
		if err != nil {
			return err
		}
		for i := 0; i < per; i++ {
			if err := resolve(7 + 8*((step*per+i)%4)); err != nil {
				return err
			}
		}
	}
	totalNS := time.Since(start).Nanoseconds()
	s.timings = append(s.timings, []string{family, n(size), n(per), n(repeat), string(mode), n(queries), n(reused), n(invalidated), n(maxFacts), n(maxRules), n(maxAbsence),
		strconv.FormatInt(setupNS, 10), strconv.FormatInt(updateNS, 10), strconv.FormatInt(queryNS, 10), strconv.FormatInt(verifyNS, 10), strconv.FormatInt(totalNS, 10)})
	s.report.Measurements++
	return nil
}

func (s *study) performance() error {
	for _, broad := range []bool{false, true} {
		family := "local"
		if broad {
			family = "broad"
		}
		for _, size := range []int{512, 4096} {
			p, world, err := workload(size, broad)
			if err != nil {
				return err
			}
			for _, per := range []int{1, 16} {
				for repeat := 0; repeat < 9; repeat++ {
					for offset := range modes {
						mode := modes[(repeat+offset)%len(modes)]
						if err := s.measure(p, world, family, size, per, repeat, mode); err != nil {
							return err
						}
					}
				}
			}
		}
	}
	return nil
}

var eventHeader = []string{"initial", "step", "query_code", "mode", "expected", "positive_only_claim", "reused", "invalidation", "fact_dependencies", "rule_dependencies"}
var timingHeader = []string{"family", "atoms_in_chain", "queries_per_update", "repeat", "mode", "queries", "reused", "invalidated", "max_fact_dependencies", "max_rule_dependencies", "max_absence_literals", "setup_ns", "update_ns", "query_ns", "external_verify_ns", "total_ns"}

func write(path string, header []string, rows [][]string) error {
	file, err := os.OpenFile(path, os.O_CREATE|os.O_EXCL|os.O_WRONLY, 0644)
	if err != nil {
		return err
	}
	writer := csv.NewWriter(file)
	if err := writer.Write(header); err != nil {
		file.Close()
		return err
	}
	if err := writer.WriteAll(rows); err != nil {
		file.Close()
		return err
	}
	return errors.Join(writer.Error(), file.Close())
}

func run(dir string) (report, error) {
	s := &study{report: report{Reused: map[e.CleanCacheMode]int{}, Go: runtime.Version(), Scope: "synthetic fixed signed DAGs; query and verification include snapshot work, total includes setup and updates; normative policy unchanged"}}
	if err := s.correctness(); err != nil {
		return report{}, err
	}
	if err := s.performance(); err != nil {
		return report{}, err
	}
	if dir != "" {
		for _, name := range []string{"events.csv", "timings.csv"} {
			if _, err := os.Lstat(filepath.Join(dir, name)); err == nil || !os.IsNotExist(err) {
				return report{}, errors.New("refusing to overwrite research data")
			}
		}
		if err := os.MkdirAll(dir, 0755); err != nil {
			return report{}, err
		}
		if err := write(filepath.Join(dir, "events.csv"), eventHeader, s.events); err != nil {
			return report{}, err
		}
		if err := write(filepath.Join(dir, "timings.csv"), timingHeader, s.timings); err != nil {
			return report{}, err
		}
	}
	return s.report, nil
}
func main() {
	dir := flag.String("csv-dir", "", "explicit new output directory; omitted means no files")
	flag.Parse()
	if flag.NArg() != 0 {
		fmt.Fprintln(os.Stderr, "unexpected positional arguments")
		os.Exit(2)
	}
	result, err := run(*dir)
	if err == nil {
		err = json.NewEncoder(os.Stdout).Encode(result)
	}
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}
