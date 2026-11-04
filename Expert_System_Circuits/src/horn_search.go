package evidence

import (
	"container/heap"
	"errors"
	"math"
)

var ErrHornBudget = errors.New("Horn: graph operation budget exhausted")

var HornMethods = []string{"fifo", "backward", "greedy", "heat", "pagerank"}

type HornSearchResult struct {
	Method                          string
	Proved, Exhausted               bool
	Operations, Applied, RankVisits int
	Proof                           HornProof
}

type hornMeter struct{ used, limit int }

func (m *hornMeter) take(n int) bool {
	if n > m.limit-m.used {
		return false
	}
	m.used += n
	return true
}

type hornAgenda struct {
	items  []int
	score  []float64
	rules  []HornRule
	greedy bool
}

func (q hornAgenda) Len() int { return len(q.items) }
func (q hornAgenda) Less(i, j int) bool {
	a, b := q.items[i], q.items[j]
	if q.score[a] != q.score[b] {
		return q.score[a] > q.score[b]
	}
	if q.greedy && len(q.rules[a].Premises) != len(q.rules[b].Premises) {
		return len(q.rules[a].Premises) < len(q.rules[b].Premises)
	}
	return a < b
}
func (q hornAgenda) Swap(i, j int) { q.items[i], q.items[j] = q.items[j], q.items[i] }
func (q *hornAgenda) Push(x any)   { q.items = append(q.items, x.(int)) }
func (q *hornAgenda) Pop() any {
	last := q.items[len(q.items)-1]
	q.items = q.items[:len(q.items)-1]
	return last
}

func (p *HornProgram) Search(snapshot HornSnapshot, query int, method string, budget int) (HornSearchResult, error) {
	if err := p.validate(snapshot, query); err != nil {
		return HornSearchResult{}, err
	}
	validMethod := false
	for _, name := range HornMethods {
		validMethod = validMethod || name == method
	}
	if !validMethod || budget < 1 {
		return HornSearchResult{}, errors.New("Horn search: unknown method or nonpositive budget")
	}
	m := &hornMeter{limit: budget}
	out := HornSearchResult{Method: method}
	stopped := func() (HornSearchResult, error) {
		out.Operations, out.Exhausted = m.used, true
		return out, nil
	}
	if !m.take(2 * p.atoms) {
		return stopped()
	}
	known := append([]bool(nil), snapshot.Facts...)
	parents := make([]int, p.atoms)
	for i := range parents {
		parents[i] = -1
	}
	finish := func() (HornSearchResult, error) {
		proof, err := p.certificate(snapshot, query, parents, m.take)
		if errors.Is(err, ErrHornBudget) {
			return stopped()
		}
		if err != nil {
			return HornSearchResult{}, err
		}
		cost := p.atoms + len(p.rules) + len(proof.Facts)
		for _, index := range proof.Rules {
			cost += 1 + len(p.rules[index].Premises)
		}
		if !m.take(cost) {
			return stopped()
		}
		if err := p.VerifyProof(snapshot, query, proof); err != nil {
			return HornSearchResult{}, err
		}
		out.Proved, out.Proof, out.Operations = true, proof, m.used
		return out, nil
	}
	if known[query] {
		return finish()
	}
	eligible := make([]bool, len(p.rules))
	distance := make([]int, p.atoms)
	if method == "fifo" {
		for i, enabled := range snapshot.Enabled {
			if !m.take(1) {
				return stopped()
			}
			eligible[i] = enabled
		}
	} else {
		for i := range distance {
			distance[i] = p.atoms + 1
		}
		distance[query] = 0
		queue := []int{query}
		for at := 0; at < len(queue); at++ {
			atom := queue[at]
			for _, index := range p.producers[atom] {
				if !m.take(1) {
					return stopped()
				}
				if !snapshot.Enabled[index] {
					continue
				}
				eligible[index] = true
				for _, parent := range p.rules[index].Premises {
					if !m.take(1) {
						return stopped()
					}
					if distance[parent] == p.atoms+1 {
						distance[parent] = distance[atom] + 1
						queue = append(queue, parent)
					}
				}
			}
		}
	}
	score := make([]float64, len(p.rules))
	if method == "greedy" {
		for i, rule := range p.rules {
			if !m.take(1) {
				return stopped()
			}
			score[i] = -float64(distance[rule.Conclusion])
		}
	} else if method == "heat" || method == "pagerank" {
		var err error
		score, out.RankVisits, err = p.diffusionRanks(query, eligible, method, m)
		if errors.Is(err, ErrHornBudget) {
			return stopped()
		}
		if err != nil {
			return HornSearchResult{}, err
		}
	}
	fifo := []int{}
	priorityQueue := &hornAgenda{score: score, rules: p.rules, greedy: method == "greedy"}
	useHeap := method != "fifo" && method != "backward"
	push := func(index int) bool {
		if !m.take(1) {
			return false
		}
		if useHeap {
			heap.Push(priorityQueue, index)
		} else {
			fifo = append(fifo, index)
		}
		return true
	}
	missing := make([]int, len(p.rules))
	for i, rule := range p.rules {
		if !m.take(1) {
			return stopped()
		}
		if !eligible[i] {
			continue
		}
		for _, atom := range rule.Premises {
			if !m.take(1) {
				return stopped()
			}
			if !known[atom] {
				missing[i]++
			}
		}
		if missing[i] == 0 && !push(i) {
			return stopped()
		}
	}
	at := 0
	for at < len(fifo) || priorityQueue.Len() > 0 {
		if !m.take(1) {
			return stopped()
		}
		index := 0
		if useHeap {
			index = heap.Pop(priorityQueue).(int)
		} else {
			index = fifo[at]
			at++
		}
		out.Applied++
		head := p.rules[index].Conclusion
		if known[head] {
			continue
		}
		known[head], parents[head] = true, index
		if head == query {
			return finish()
		}
		for _, consumer := range p.consumers[head] {
			if !m.take(1) {
				return stopped()
			}
			if !eligible[consumer] {
				continue
			}
			missing[consumer]--
			if missing[consumer] == 0 && !push(consumer) {
				return stopped()
			}
		}
	}
	out.Operations = m.used
	return out, nil
}

func (p *HornProgram) diffusionRanks(query int, eligible []bool, method string, m *hornMeter) ([]float64, int, error) {
	size := p.atoms + len(p.rules)
	if !m.take(size) {
		return nil, 0, ErrHornBudget
	}
	adj := make([][]int, size)
	active := make([]bool, size)
	nodes := []int{}
	addNode := func(v int) {
		if !active[v] {
			active[v] = true
			nodes = append(nodes, v)
		}
	}
	addNode(query)
	for i, allowed := range eligible {
		if !m.take(1) {
			return nil, 0, ErrHornBudget
		}
		if !allowed {
			continue
		}
		ruleNode := p.atoms + i
		for _, atom := range append(append([]int(nil), p.rules[i].Premises...), p.rules[i].Conclusion) {
			if !m.take(2) {
				return nil, 0, ErrHornBudget
			}
			adj[ruleNode] = append(adj[ruleNode], atom)
			adj[atom] = append(adj[atom], ruleNode)
			addNode(atom)
			addNode(ruleNode)
		}
	}
	maxDegree := 1
	for _, node := range nodes {
		if len(adj[node]) > maxDegree {
			maxDegree = len(adj[node])
		}
	}
	values, next := make([]float64, size), make([]float64, size)
	values[query] = 1
	visits := 0
	eta := .5 / float64(maxDegree)
	for iteration := 0; iteration < 32; iteration++ {
		for _, node := range nodes {
			if !m.take(1) {
				return nil, visits, ErrHornBudget
			}
			next[node] = 0
		}
		if method == "pagerank" {
			next[query] = .15
		}
		for _, node := range nodes {
			if !m.take(1) {
				return nil, visits, ErrHornBudget
			}
			degree := len(adj[node])
			if method == "heat" {
				next[node] += values[node] * (1 - eta*float64(degree))
			} else if degree == 0 {
				next[query] += .85 * values[node]
			}
			for _, neighbor := range adj[node] {
				if !m.take(1) {
					return nil, visits, ErrHornBudget
				}
				visits++
				if method == "heat" {
					next[neighbor] += eta * values[node]
				} else {
					next[neighbor] += .85 * values[node] / float64(degree)
				}
			}
		}
		values, next = next, values
	}
	score := make([]float64, len(p.rules))
	for i := range score {
		if !m.take(1) {
			return nil, visits, ErrHornBudget
		}
		score[i] = values[p.atoms+i]
		if math.IsNaN(score[i]) || math.IsInf(score[i], 0) {
			return nil, visits, errors.New("Horn diffusion: nonfinite rank")
		}
	}
	return score, visits, nil
}
