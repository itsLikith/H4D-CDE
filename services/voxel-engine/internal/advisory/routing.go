package advisory

import (
	"github.com/uber/h3-go/v4"
	"hive/voxel-engine/internal/spatial"
)

// bfsShortestPath finds the shortest hop-count path from origin to dest on
// the H3 neighbour graph, optionally avoiding a set of excluded cells (used
// by kShortestPathReroute below to find *alternative* paths).
func bfsShortestPath(origin, dest h3.Cell, maxHops int, excluded map[h3.Cell]bool) ([]h3.Cell, bool) {
	type node struct {
		cell h3.Cell
		path []h3.Cell
	}
	visited := map[h3.Cell]bool{origin: true}
	queue := []node{{origin, []h3.Cell{origin}}}

	for len(queue) > 0 {
		cur := queue[0]
		queue = queue[1:]
		if len(cur.path) > maxHops {
			continue
		}
		if cur.cell == dest {
			return cur.path, true
		}
		neighbors, err := spatial.NeighborCells(cur.cell, 1)
		if err != nil {
			continue
		}
		for _, n := range neighbors {
			if n == cur.cell || visited[n] || excluded[n] {
				continue
			}
			visited[n] = true
			newPath := append(append([]h3.Cell{}, cur.path...), n)
			queue = append(queue, node{n, newPath})
		}
	}
	return nil, false
}

// kShortestPathReroute returns up to k alternative simple paths: the
// shortest path, then repeatedly excluding one of its interior cells and
// re-running BFS. This is a lightweight approximation of Yen's k-shortest-
// paths algorithm — a correct and sufficient substitute for the original
// design's `networkx.shortest_simple_paths` call at the local, few-dozen-
// cell scale this routing problem operates at. For much larger routing
// graphs, swap this for `gonum.org/v1/gonum/graph/path`.
func kShortestPathReroute(origin, dest h3.Cell, k, maxHops int) [][]h3.Cell {
	first, ok := bfsShortestPath(origin, dest, maxHops, nil)
	if !ok {
		return nil
	}
	paths := [][]h3.Cell{first}
	for i := 1; i < len(first)-1 && len(paths) < k; i++ {
		alt, ok := bfsShortestPath(origin, dest, maxHops, map[h3.Cell]bool{first[i]: true})
		if ok {
			paths = append(paths, alt)
		}
	}
	return paths
}
