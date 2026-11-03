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
	"slices"
	"strconv"
	"time"

	evidence "github.com/nick-fedchik/articles/Expert_System_Circuits/src"
)

type model struct {
	name  string
	atoms int
	rules []evidence.SignedRule
}

type summary struct {
	Cases, Clean, ConflictDependent, Unsupported int
	LeavesOnlyDisagreements, AllProofsOverblocks int
	ByFamily                                     map[string]int
	ScalingRows, EnumerationBudgetFailures       int
	Go, Platform, Scope                          string
}

type study struct {
	report   summary
	semantic [][]string
	scaling  [][]string
}

func newStudy() *study {
	return &study{report: summary{ByFamily: map[string]int{}, Go: runtime.Version(),
		Platform: fmt.Sprintf("%s/%s GOMAXPROCS=%d", runtime.GOOS, runtime.GOARCH, runtime.GOMAXPROCS(0)),
		Scope:    "synthetic signed acyclic models; clean support is an experimental query property, not a normative verdict or proof of real-world truth"}}
}

func rule(id string, conclusion evidence.SignedLiteral, premises ...evidence.SignedLiteral) evidence.SignedRule {
	return evidence.SignedRule{ID: id, Version: "1", Premises: premises, Conclusion: conclusion}
}

func families() []model {
	pos := func(atom int) evidence.SignedLiteral { return evidence.SignedLiteral{Atom: atom} }
	neg := func(atom int) evidence.SignedLiteral { return evidence.SignedLiteral{Atom: atom, Negative: true} }
	return []model{
		{"signed-chain", 3, []evidence.SignedRule{
			rule("a", pos(1), pos(0)), rule("b", neg(1), neg(0)), rule("c", pos(2), pos(1)),
		}},
		{"alternative", 3, []evidence.SignedRule{
			rule("a", pos(2), pos(0)), rule("b", pos(2), pos(1)), rule("c", neg(0), neg(1)),
		}},
		{"intermediate-conflict", 4, []evidence.SignedRule{
			rule("a", pos(2), pos(0)), rule("b", neg(2), pos(1)),
			rule("c", pos(3), pos(2), pos(1)), rule("d", neg(3), neg(0)),
		}},
	}
}

func b(value bool) string { return strconv.FormatBool(value) }
func n(value int) string  { return strconv.Itoa(value) }

func (s *study) semantics() error {
	for _, family := range families() {
		p, err := evidence.CompileCleanSupport(family.atoms, family.rules)
		if err != nil {
			return err
		}
		for factsCode := 0; factsCode < 1<<(2*family.atoms); factsCode++ {
			facts := make([]evidence.State, family.atoms)
			for atom := range facts {
				facts[atom] = evidence.State(factsCode >> (2 * atom) & 3)
			}
			for mask := 0; mask < 1<<len(family.rules); mask++ {
				enabled := make([]bool, len(family.rules))
				for i := range enabled {
					enabled[i] = mask&(1<<i) != 0
				}
				world := evidence.SignedWorld{Revision: 1, Facts: facts, Enabled: enabled}
				for code := 0; code < 2*family.atoms; code++ {
					query := evidence.SignedLiteral{Atom: code / 2, Negative: code%2 != 0}
					cert, work, err := p.EvaluateClean(world, query)
					if err != nil {
						return err
					}
					if err := p.VerifyClean(world, query, cert); err != nil {
						return err
					}
					ref, err := p.ReferenceClean(world, query, 100000)
					if err != nil {
						return err
					}
					clean := cert.Status == evidence.CleanSupported
					if !slices.Equal(ref.Closure, cert.Closure) || ref.ExistsClean != clean ||
						ref.Supported != cert.Closure[code] {
						return fmt.Errorf("clean/reference disagreement: %s/%d/%d/%d", family.name, factsCode, mask, code)
					}
					s.report.Cases++
					s.report.ByFamily[family.name]++
					switch cert.Status {
					case evidence.CleanSupported:
						s.report.Clean++
					case evidence.ConflictDependent:
						s.report.ConflictDependent++
					case evidence.NotSupported:
						s.report.Unsupported++
					default:
						return errors.New("unknown clean-support status")
					}
					if ref.LeavesOnly != clean {
						if clean || !ref.LeavesOnly {
							return errors.New("unexpected direction of leaves-only disagreement")
						}
						s.report.LeavesOnlyDisagreements++
					}
					if ref.AllProofsClean != clean {
						if !clean || ref.AllProofsClean {
							return errors.New("unexpected direction of all-proofs disagreement")
						}
						s.report.AllProofsOverblocks++
					}
					s.semantic = append(s.semantic, []string{family.name, n(factsCode), n(mask), n(code), string(cert.Status),
						b(ref.Supported), b(ref.LeavesOnly), b(ref.AllProofsClean), n(ref.Proofs), n(work.Nodes), n(work.Rules), n(work.Premises)})
				}
			}
		}
	}
	return nil
}

func alternatives(pairs int) (model, evidence.SignedWorld, evidence.SignedLiteral) {
	m := model{name: "and-alternatives", atoms: 3*pairs + 1}
	world := evidence.SignedWorld{Revision: 1, Facts: make([]evidence.State, m.atoms)}
	premises := []evidence.SignedLiteral{}
	for pair := 0; pair < pairs; pair++ {
		left, right, head := 3*pair, 3*pair+1, 3*pair+2
		world.Facts[left], world.Facts[right] = evidence.True, evidence.True
		target := evidence.SignedLiteral{Atom: head}
		m.rules = append(m.rules, rule(fmt.Sprintf("left-%d", pair), target, evidence.SignedLiteral{Atom: left}),
			rule(fmt.Sprintf("right-%d", pair), target, evidence.SignedLiteral{Atom: right}))
		premises = append(premises, target)
	}
	query := evidence.SignedLiteral{Atom: m.atoms - 1}
	m.rules = append(m.rules, rule("query", query, premises...))
	world.Enabled = make([]bool, len(m.rules))
	for i := range world.Enabled {
		world.Enabled[i] = true
	}
	return m, world, query
}

func (s *study) measure(m model, world evidence.SignedWorld, query evidence.SignedLiteral, regime string, want bool, enumerate bool) error {
	p, err := evidence.CompileCleanSupport(m.atoms, m.rules)
	if err != nil {
		return err
	}
	edges := 0
	for _, rule := range m.rules {
		edges += len(rule.Premises)
	}
	refStatus, refProofs, refWork, refNS := "not_run_size", "", "", ""
	if enumerate {
		start := time.Now()
		ref, err := p.ReferenceClean(world, query, 100000)
		refNS = strconv.FormatInt(time.Since(start).Nanoseconds(), 10)
		if errors.Is(err, evidence.ErrCleanProofBudget) {
			refStatus = "budget_exhausted"
			s.report.EnumerationBudgetFailures++
		} else if err != nil {
			return err
		} else {
			if ref.ExistsClean != want {
				return errors.New("scaling reference disagrees with constructed case")
			}
			refStatus, refProofs, refWork = "complete", n(ref.Proofs), n(ref.Work)
		}
	}
	for repeat := 0; repeat < 9; repeat++ {
		start := time.Now()
		cert, work, err := p.EvaluateClean(world, query)
		evaluateNS := time.Since(start).Nanoseconds()
		if err != nil {
			return err
		}
		start = time.Now()
		if err := p.VerifyClean(world, query, cert); err != nil {
			return err
		}
		verifyNS := time.Since(start).Nanoseconds()
		if (cert.Status == evidence.CleanSupported) != want || work.Nodes != 2*m.atoms ||
			work.Rules > len(m.rules) || work.Premises > edges {
			return errors.New("scaling clean result or abstract linear work bound violated")
		}
		witness := 0
		if cert.Proof != nil {
			witness = len(cert.Proof.Facts) + len(cert.Proof.Rules)
		}
		s.scaling = append(s.scaling, []string{m.name, n(m.atoms), n(len(m.rules)), n(edges), regime, n(repeat),
			string(cert.Status), n(work.Nodes), n(work.Rules), n(work.Premises), n(witness),
			strconv.FormatInt(evaluateNS, 10), strconv.FormatInt(verifyNS, 10), refStatus, refProofs, refWork, refNS})
		s.report.ScalingRows++
	}
	return nil
}

func (s *study) scale() error {
	for _, pairs := range []int{4, 8, 12, 16, 20, 24} {
		m, world, query := alternatives(pairs)
		for _, regime := range []string{"all-clean", "one-tainted-alternative", "blocked-pair"} {
			world.Facts[0], world.Facts[1] = evidence.True, evidence.True
			if regime != "all-clean" {
				world.Facts[0] = evidence.Both
			}
			if regime == "blocked-pair" {
				world.Facts[1] = evidence.Both
			}
			if err := s.measure(m, world, query, regime, regime != "blocked-pair", true); err != nil {
				return err
			}
		}
	}
	for _, size := range []int{32, 256, 2048, 8192} {
		m := model{name: "chain", atoms: size}
		world := evidence.SignedWorld{Revision: 1, Facts: make([]evidence.State, size)}
		for atom := 1; atom < size; atom++ {
			m.rules = append(m.rules, rule(fmt.Sprintf("r-%d", atom), evidence.SignedLiteral{Atom: atom}, evidence.SignedLiteral{Atom: atom - 1}))
			world.Enabled = append(world.Enabled, true)
		}
		for _, conflict := range []bool{false, true} {
			regime := "all-clean"
			world.Facts[0] = evidence.True
			if conflict {
				regime, world.Facts[0] = "conflicted-leaf", evidence.Both
			}
			if err := s.measure(m, world, evidence.SignedLiteral{Atom: size - 1}, regime, !conflict, size <= 128); err != nil {
				return err
			}
		}
	}
	return nil
}

var semanticHeader = []string{"family", "facts_code", "enabled_mask", "query_code", "status", "supported",
	"leaves_only", "all_proofs_clean", "proof_trees", "nodes", "rules", "premises"}
var scalingHeader = []string{"family", "atoms", "rules_total", "edges", "regime", "repeat", "status", "nodes", "rules", "premises",
	"witness_items", "evaluate_ns", "verify_ns", "reference_status", "reference_proofs", "reference_work", "reference_ns"}

func writeCSV(path string, header []string, rows [][]string) error {
	file, err := os.OpenFile(path, os.O_WRONLY|os.O_CREATE|os.O_EXCL, 0644)
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
	writer.Flush()
	return errors.Join(writer.Error(), file.Close())
}

func run(dir string) (summary, error) {
	s := newStudy()
	if err := s.semantics(); err != nil {
		return summary{}, err
	}
	if err := s.scale(); err != nil {
		return summary{}, err
	}
	if dir != "" {
		for _, name := range []string{"semantics.csv", "scaling.csv"} {
			if _, err := os.Lstat(filepath.Join(dir, name)); err == nil || !os.IsNotExist(err) {
				return summary{}, errors.New("clean study refuses to overwrite published results")
			}
		}
		if err := os.MkdirAll(dir, 0755); err != nil {
			return summary{}, err
		}
		if err := writeCSV(filepath.Join(dir, "semantics.csv"), semanticHeader, s.semantic); err != nil {
			return summary{}, err
		}
		if err := writeCSV(filepath.Join(dir, "scaling.csv"), scalingHeader, s.scaling); err != nil {
			return summary{}, err
		}
	}
	return s.report, nil
}

func main() {
	dir := flag.String("csv-dir", "", "explicit new results directory; omitted means no CSV files")
	flag.Parse()
	if flag.NArg() != 0 {
		fmt.Fprintln(os.Stderr, "unexpected positional clean-study arguments")
		os.Exit(2)
	}
	report, err := run(*dir)
	if err == nil {
		err = json.NewEncoder(os.Stdout).Encode(report)
	}
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}
