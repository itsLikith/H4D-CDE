// Copyright 2026 Likith Saragadam
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package advisory

import (
	"github.com/itsLikith/h4d-cde/services/voxel-engine/internal/spatial"
	"github.com/uber/h3-go/v4"
)

// bfsShortestPath performs Breadth-First Search on the unweighted H3 adjacency graph.
// Because all adjacent hexagonal edges have uniform distance, BFS guarantees the minimum hop-count shortest path.
func bfsShortestPath(origin, dest h3.Cell, maxHops int, excluded map[h3.Cell]bool) ([]h3.Cell, bool) {
	if origin == dest {
		return []h3.Cell{origin}, true
	}

	type node struct {
		cell h3.Cell
		path []h3.Cell
	}

	visited := make(map[h3.Cell]bool)
	visited[origin] = true

	queue := []node{{cell: origin, path: []h3.Cell{origin}}}

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

		for _, nbr := range neighbors {
			if nbr == cur.cell || visited[nbr] || (excluded != nil && excluded[nbr]) {
				continue
			}

			visited[nbr] = true
			nextPath := make([]h3.Cell, len(cur.path)+1)
			copy(nextPath, cur.path)
			nextPath[len(cur.path)] = nbr

			queue = append(queue, node{cell: nbr, path: nextPath})
		}
	}

	return nil, false
}

// kShortestPathReroute computes up to k alternative paths by iteratively masking interior path nodes.
// Lightweight Yen's algorithm approximation over the localized H3 airspace graph.
func kShortestPathReroute(origin, dest h3.Cell, k, maxHops int) [][]h3.Cell {
	firstPath, ok := bfsShortestPath(origin, dest, maxHops, nil)
	if !ok {
		return nil
	}

	paths := [][]h3.Cell{firstPath}

	for i := 1; i < len(firstPath)-1 && len(paths) < k; i++ {
		excluded := map[h3.Cell]bool{firstPath[i]: true}
		altPath, found := bfsShortestPath(origin, dest, maxHops, excluded)
		if found {
			paths = append(paths, altPath)
		}
	}

	return paths
}
