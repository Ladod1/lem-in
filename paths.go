package main

// Each room is split in two nodes (in = 2*i, out = 2*i+1) joined by an edge
// of capacity 1, so a room can only be used by one path. Then we run
// Edmonds-Karp and after every new path we check how many turns the current
// set of paths needs, keeping the best set.

type edge struct {
	to, cap, flow, cost, rev int
}

type graph [][]edge

func (g graph) add(a, b, c, cost int) {
	g[a] = append(g[a], edge{to: b, cap: c, cost: cost, rev: len(g[b])})
	g[b] = append(g[b], edge{to: a, cap: 0, cost: -cost, rev: len(g[a]) - 1})
}

func findPaths(f *Farm) [][]int {
	n := len(f.names)
	g := make(graph, 2*n)
	for i := 0; i < n; i++ {
		c := 1
		if i == f.start || i == f.end {
			c = f.ants
		}
		g.add(2*i, 2*i+1, c, 0)
		for _, j := range f.links[i] {
			g.add(2*i+1, 2*j, 1, 1)
		}
	}
	src, sink := 2*f.start+1, 2*f.end

	var best [][]int
	bestTurns := -1
	for k := 0; k < f.ants; k++ {
		if !augment(g, src, sink) {
			break
		}
		paths := extractPaths(g, f, src)
		t := countTurns(paths, f.ants)
		if bestTurns == -1 || t < bestTurns {
			best, bestTurns = paths, t
		}
	}
	return best
}

// augment finds the shortest new path (counting tunnels) in the residual
// graph and pushes flow on it. Going back through a used tunnel counts -1,
// so it can also reroute old paths when that makes the total shorter.
func augment(g graph, src, sink int) bool {
	type step struct{ node, edge int }
	prev := make([]step, len(g))
	dist := make([]int, len(g))
	inQueue := make([]bool, len(g))
	for i := range prev {
		prev[i].node = -1
		dist[i] = 1 << 30
	}
	prev[src].node = src
	dist[src] = 0
	queue := []int{src}
	for len(queue) > 0 {
		u := queue[0]
		queue = queue[1:]
		inQueue[u] = false
		for i, e := range g[u] {
			if e.cap-e.flow > 0 && dist[u]+e.cost < dist[e.to] {
				dist[e.to] = dist[u] + e.cost
				prev[e.to] = step{u, i}
				if !inQueue[e.to] {
					inQueue[e.to] = true
					queue = append(queue, e.to)
				}
			}
		}
	}
	if prev[sink].node == -1 {
		return false
	}
	for v := sink; v != src; v = prev[v].node {
		e := &g[prev[v].node][prev[v].edge]
		e.flow++
		g[v][e.rev].flow--
	}
	return true
}

// extractPaths follows the flow from start to end and returns the room
// indexes of every path (without the start room), shortest first.
func extractPaths(g graph, f *Farm, src int) [][]int {
	var paths [][]int
	for _, e := range g[src] {
		if e.flow <= 0 || e.cap == 0 {
			continue
		}
		room := e.to / 2
		path := []int{room}
		for room != f.end {
			for _, e2 := range g[2*room+1] {
				if e2.flow > 0 && e2.cap > 0 {
					room = e2.to / 2
					break
				}
			}
			path = append(path, room)
		}
		paths = append(paths, path)
	}
	// sort by length (insertion sort, few paths)
	for i := 1; i < len(paths); i++ {
		for j := i; j > 0 && len(paths[j]) < len(paths[j-1]); j-- {
			paths[j], paths[j-1] = paths[j-1], paths[j]
		}
	}
	return paths
}

// assign gives each ant a path: the one where it will arrive first.
// It returns the path index of every ant.
func assign(paths [][]int, ants int) []int {
	count := make([]int, len(paths))
	result := make([]int, ants)
	for a := 0; a < ants; a++ {
		bestP := 0
		for p := 1; p < len(paths); p++ {
			if len(paths[p])+count[p] <= len(paths[bestP])+count[bestP] {
				bestP = p
			}
		}
		result[a] = bestP
		count[bestP]++
	}
	return result
}

func countTurns(paths [][]int, ants int) int {
	count := make([]int, len(paths))
	for _, p := range assign(paths, ants) {
		count[p]++
	}
	turns := 0
	for p := range paths {
		if count[p] > 0 && len(paths[p])+count[p]-1 > turns {
			turns = len(paths[p]) + count[p] - 1
		}
	}
	return turns
}
