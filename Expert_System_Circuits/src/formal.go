package evidence

import (
	"bufio"
	"context"
	"crypto/sha256"
	"errors"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strconv"
	"strings"
)

type cnf struct {
	variables int
	clauses   [][]int
}

func (c *cnf) variable() int {
	c.variables++
	return c.variables
}

func (c *cnf) add(literals ...int) { c.clauses = append(c.clauses, literals) }

func (c *cnf) constant(value bool) int {
	v := c.variable()
	if value {
		c.add(v)
	} else {
		c.add(-v)
	}
	return v
}

func (c *cnf) and(a, b int) int {
	out := c.variable()
	c.add(-out, a)
	c.add(-out, b)
	c.add(out, -a, -b)
	return out
}

func (c *cnf) or(a, b int) int {
	out := c.variable()
	c.add(out, -a)
	c.add(out, -b)
	c.add(-out, a, b)
	return out
}

func (c *cnf) mux(condition, whenTrue, whenFalse int) int {
	return c.or(c.and(condition, whenTrue), c.and(-condition, whenFalse))
}

func (c *cnf) xor(a, b int) int {
	return c.or(c.and(a, -b), c.and(-a, b))
}

func (c *cnf) solve(ctx context.Context, solver string) (bool, []bool, error) {
	if ctx == nil || solver == "" || c.variables > 200000 || len(c.clauses) > 1000000 {
		return false, nil, errors.New("SAT: context, solver and bounded formula required")
	}
	if err := ctx.Err(); err != nil {
		return false, nil, err
	}
	dir, err := os.MkdirTemp("", "teos-sat-")
	if err != nil {
		return false, nil, err
	}
	defer os.RemoveAll(dir)
	var input strings.Builder
	fmt.Fprintf(&input, "p cnf %d %d\n", c.variables, len(c.clauses))
	for _, clause := range c.clauses {
		for _, literal := range clause {
			fmt.Fprintf(&input, "%d ", literal)
		}
		input.WriteString("0\n")
	}
	in, out := filepath.Join(dir, "query.cnf"), filepath.Join(dir, "answer.txt")
	if err := os.WriteFile(in, []byte(input.String()), 0600); err != nil {
		return false, nil, err
	}
	cmd := exec.CommandContext(ctx, solver, "-verb=0", in, out)
	log, runErr := cmd.CombinedOutput()
	if err := ctx.Err(); err != nil {
		return false, nil, fmt.Errorf("SAT: %w", err)
	}
	var exit *exec.ExitError
	if !errors.As(runErr, &exit) || (exit.ExitCode() != 10 && exit.ExitCode() != 20) {
		return false, nil, fmt.Errorf("SAT solver failed: %v: %.500s", runErr, log)
	}
	answer, err := os.ReadFile(out)
	if err != nil {
		return false, nil, fmt.Errorf("SAT output: %w", err)
	}
	fields := strings.Fields(string(answer))
	if exit.ExitCode() == 20 {
		if len(fields) != 1 || fields[0] != "UNSAT" {
			return false, nil, errors.New("SAT: malformed UNSAT response")
		}
		return false, nil, nil
	}
	if len(fields) < 2 || fields[0] != "SAT" || fields[len(fields)-1] != "0" {
		return false, nil, errors.New("SAT: malformed model")
	}
	values, seen := make([]bool, c.variables+1), make([]bool, c.variables+1)
	for _, field := range fields[1 : len(fields)-1] {
		v, err := strconv.Atoi(field)
		if err != nil || v == 0 || v < -c.variables || v > c.variables {
			return false, nil, errors.New("SAT: invalid model literal")
		}
		index := v
		if index < 0 {
			index = -index
		}
		if seen[index] {
			return false, nil, errors.New("SAT: duplicate assignment")
		}
		values[index], seen[index] = v > 0, true
	}
	for i := 1; i <= c.variables; i++ {
		if !seen[i] {
			return false, nil, errors.New("SAT: incomplete model")
		}
	}
	for _, clause := range c.clauses {
		holds := false
		for _, lit := range clause {
			if lit > 0 && values[lit] || lit < 0 && !values[-lit] {
				holds = true
			}
		}
		if !holds {
			return false, nil, errors.New("SAT: solver model fails independent clause check")
		}
	}
	return true, values, nil
}

type symbolicState struct{ high, low int }

func (c *cnf) is(state symbolicState, value State) int {
	hi, lo := state.high, state.low
	if value&2 == 0 {
		hi = -hi
	}
	if value&1 == 0 {
		lo = -lo
	}
	return c.and(hi, lo)
}

func (c *cnf) referencePriority(a, b symbolicState) symbolicState {
	f := c.or(c.is(a, False), c.is(b, False))
	n := c.and(-f, c.or(c.is(a, Unknown), c.is(b, Unknown)))
	conflict := c.and(c.and(-f, -n), c.or(c.is(a, Both), c.is(b, Both)))
	return symbolicState{c.and(-f, -n), c.or(f, conflict)}
}

func (c *cnf) tablePriority(a, b symbolicState, lut uint32) symbolicState {
	out := symbolicState{c.variable(), c.variable()}
	for x := Unknown; x <= Both; x++ {
		for y := Unknown; y <= Both; y++ {
			guard := c.and(c.is(a, x), c.is(b, y))
			value := State((lut >> (2 * (4*uint(x) + uint(y)))) & 3)
			for _, bit := range []struct {
				literal int
				mask    State
			}{{out.high, 2}, {out.low, 1}} {
				literal := bit.literal
				if value&bit.mask == 0 {
					literal = -literal
				}
				c.add(-guard, literal)
			}
		}
	}
	return out
}

type symbolicValue struct {
	state   symbolicState
	support []int
}

func (c *cnf) combine(op Op, a, b, cond symbolicValue, table bool, lut uint32) symbolicValue {
	join := c.referencePriority
	if table {
		join = func(a, b symbolicState) symbolicState { return c.tablePriority(a, b, lut) }
	}
	out := symbolicValue{state: join(a.state, b.state), support: make([]int, len(a.support))}
	for i := range out.support {
		out.support[i] = c.or(a.support[i], b.support[i])
	}
	if op == Select {
		t, f := c.is(cond.state, True), c.is(cond.state, False)
		both := symbolicState{c.constant(true), c.constant(true)}
		uncertain := join(out.state, both)
		out.state = symbolicState{
			c.mux(t, b.state.high, c.mux(f, a.state.high, uncertain.high)),
			c.mux(t, b.state.low, c.mux(f, a.state.low, uncertain.low)),
		}
		for i := range out.support {
			out.support[i] = c.or(cond.support[i], c.mux(t, b.support[i], c.mux(f, a.support[i], out.support[i])))
		}
	}
	return out
}

type FormalCheck struct {
	Equivalent bool `json:"equivalent"`
	Variables  int  `json:"variables"`
	Clauses    int  `json:"clauses"`
}

// VerifyGraphSAT checks every state assignment of one bounded graph, not all possible graphs.
func VerifyGraphSAT(ctx context.Context, solver string, root *Expr, sources []Source) (FormalCheck, error) {
	return verifyGraphSAT(ctx, solver, root, sources, priorityLUT)
}

func verifyGraphSAT(ctx context.Context, solver string, root *Expr, sources []Source, lut uint32) (FormalCheck, error) {
	program, err := Compile(root, sources)
	if err != nil {
		return FormalCheck{}, err
	}
	if len(program.code) > 100 || len(sources) > 32 {
		return FormalCheck{}, errors.New("SAT miter: graph limited to 100 instructions and 32 inputs")
	}
	c := &cnf{}
	inputs := make([]symbolicValue, len(sources))
	for i := range inputs {
		inputs[i] = symbolicValue{state: symbolicState{c.variable(), c.variable()}, support: make([]int, len(program.sources))}
		for k := range inputs[i].support {
			inputs[i].support[k] = c.constant(k == program.sourceIndex[i])
		}
	}
	memo := map[*Expr]symbolicValue{}
	var reference func(*Expr) symbolicValue
	reference = func(expr *Expr) symbolicValue {
		if v, ok := memo[expr]; ok {
			return v
		}
		if expr.Op == Input {
			return inputs[expr.Input]
		}
		a, b := reference(expr.Left), reference(expr.Right)
		var condition symbolicValue
		if expr.Op == Select {
			condition = reference(expr.Condition)
		}
		out := c.combine(expr.Op, a, b, condition, false, 0)
		memo[expr] = out
		return out
	}
	want := reference(root)
	compiled := make([]symbolicValue, len(program.code))
	for i, inst := range program.code {
		if inst.op == Input {
			compiled[i] = inputs[inst.input]
			continue
		}
		var condition symbolicValue
		if inst.op == Select {
			condition = compiled[inst.cond]
		}
		compiled[i] = c.combine(inst.op, compiled[inst.left], compiled[inst.right], condition, true, lut)
	}
	got := compiled[program.root]
	differences := []int{c.xor(got.state.high, want.state.high), c.xor(got.state.low, want.state.low)}
	for i := range got.support {
		differences = append(differences, c.xor(got.support[i], want.support[i]))
	}
	c.add(differences...)
	satisfiable, _, err := c.solve(ctx, solver)
	if err != nil {
		return FormalCheck{}, err
	}
	return FormalCheck{!satisfiable, c.variables, len(c.clauses)}, nil
}

type SATRepair struct {
	Repair
	Calls int `json:"solver_calls"`
}

// FindRepairSAT is an independent bounded-cost SAT oracle, not the production branch-and-bound algorithm.
func FindRepairSAT(ctx context.Context, solver string, problem RepairProblem) (SATRepair, error) {
	if ctx == nil || solver == "" {
		return SATRepair{}, errors.New("SAT repair: context and explicit solver required")
	}
	if err := problem.validate(); err != nil {
		return SATRepair{}, err
	}
	total := int64(0)
	for _, action := range problem.Actions {
		total += action.Cost
	}
	if total > 10000 || len(problem.Actions) > 40 {
		return SATRepair{}, errors.New("SAT repair: total cost <=10000 and actions <=40 required")
	}
	calls := 0
	query := func(bound int64) (bool, uint64, error) {
		calls++
		c := &cnf{}
		vars := make([]int, len(problem.Actions))
		for i, action := range problem.Actions {
			vars[i] = c.variable()
			if action.Weakens || !action.Verified {
				c.add(-vars[i])
			}
		}
		for _, cut := range problem.Cuts {
			clause := []int{}
			for _, index := range cut {
				clause = append(clause, vars[index])
			}
			c.add(clause...)
		}
		for _, pair := range problem.Conflicts {
			c.add(-vars[pair[0]], -vars[pair[1]])
		}
		yes, no := c.constant(true), c.constant(false)
		type key struct{ index, capacity int }
		cache := map[key]int{}
		var buildErr error
		var capacity func(int, int) int
		capacity = func(index, remaining int) int {
			if buildErr != nil {
				return no
			}
			if err := ctx.Err(); err != nil {
				buildErr = err
				return no
			}
			if c.variables > 190000 || len(c.clauses) > 900000 {
				buildErr = errors.New("SAT repair: encoding memory limit exceeded")
				return no
			}
			if remaining < 0 {
				return no
			}
			if index == len(vars) {
				return yes
			}
			k := key{index, remaining}
			if v, ok := cache[k]; ok {
				return v
			}
			v := c.mux(vars[index], capacity(index+1, remaining-int(problem.Actions[index].Cost)), capacity(index+1, remaining))
			cache[k] = v
			return v
		}
		c.add(capacity(0, int(bound)))
		if buildErr != nil {
			return false, 0, buildErr
		}
		sat, values, err := c.solve(ctx, solver)
		var mask uint64
		if sat && err == nil {
			for i, v := range vars {
				if values[v] {
					mask |= uint64(1) << uint(i)
				}
			}
		}
		return sat, mask, err
	}
	possible, mask, err := query(total)
	if err != nil {
		return SATRepair{}, err
	}
	if !possible {
		return SATRepair{}, ErrNoRepair
	}
	lo, hi := int64(0), total
	for lo < hi {
		mid := lo + (hi-lo)/2
		sat, candidate, err := query(mid)
		if err != nil {
			return SATRepair{}, err
		}
		if sat {
			hi, mask = mid, candidate
		} else {
			lo = mid + 1
		}
	}
	cost := int64(0)
	for i, action := range problem.Actions {
		if mask&(uint64(1)<<uint(i)) != 0 {
			cost += action.Cost
		}
	}
	if cost != lo {
		return SATRepair{}, errors.New("SAT repair: model cost does not equal proved bound")
	}
	return SATRepair{Repair: Repair{Mask: mask, Cost: cost}, Calls: calls}, nil
}

type SolverInfo struct {
	Name, Banner, SHA256 string
}

// IdentifySolver captures tool identity, not query contents or a private executable path.
func IdentifySolver(ctx context.Context, solver string) (SolverInfo, error) {
	if ctx == nil {
		return SolverInfo{}, errors.New("SAT solver: context required")
	}
	path, err := exec.LookPath(solver)
	if err != nil {
		return SolverInfo{}, err
	}
	binary, err := os.ReadFile(path)
	if err != nil {
		return SolverInfo{}, err
	}
	cmd := exec.CommandContext(ctx, solver, "--help")
	output, err := cmd.CombinedOutput()
	if err != nil {
		return SolverInfo{}, err
	}
	scanner := bufio.NewScanner(strings.NewReader(string(output)))
	if scanner.Scan() {
		return SolverInfo{filepath.Base(path), scanner.Text(), fmt.Sprintf("%x", sha256.Sum256(binary))}, nil
	}
	return SolverInfo{}, errors.New("SAT solver returned no identification")
}
