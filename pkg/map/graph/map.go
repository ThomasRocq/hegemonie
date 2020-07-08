// Copyright (C) 2018-2020 Hegemonie's AUTHORS
// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at http://mozilla.org/MPL/2.0/.

package mapgraph

import (
	"errors"
	"sort"
	"strconv"
	"strings"
	"sync/atomic"
)

func (m *Map) Init() {
	m.Cells = make(SetOfVertices, 0)
	m.Roads = make(SetOfEdges, 0)
}

func (m *Map) Check() error {
	if !sort.IsSorted(&m.Cells) {
		return errors.New("locations unsorted")
	}
	if !sort.IsSorted(&m.Roads) {
		return errors.New("roads unsorted")
	}
	return nil
}

func (m *Map) PostLoad() error {
	sort.Sort(&m.Cells)
	sort.Sort(&m.Roads)
	return nil
}

func (m *Map) getNextID() uint64 {
	return atomic.AddUint64(&m.nextID, 1)
}

func (m *Map) CellGet(id uint64) *Vertex {
	return m.Cells.Get(id)
}

func (m *Map) CellHas(id uint64) bool {
	return m.Cells.Has(id)
}

func (m *Map) CellCreate() *Vertex {
	id := m.getNextID()
	c := &Vertex{ID: id}
	m.Cells.Add(c)
	return c
}

// Raw creation of an edge, with no check the Source and Destination exist
// The set of roads isn't sorted afterwards
func (m *Map) RoadCreateRaw(src, dst uint64) *Edge {
	if src == dst || src == 0 || dst == 0 {
		panic("Invalid Edge parameters")
	}

	e := &Edge{src, dst}
	m.Roads = append(m.Roads, e)
	return e
}

func (m *Map) RoadCreate(src, dst uint64, check bool) error {
	if src == dst || src == 0 || dst == 0 {
		return errors.New("EINVAL")
	}

	if check && !m.CellHas(src) {
		return errors.New("Source not found")
	}
	if check && !m.CellHas(dst) {
		return errors.New("Destination not found")
	}

	if r := m.Roads.Get(src, dst); r != nil {
		return errors.New("Edge exists")
	}
	m.Roads.Add(&Edge{src, dst})
	return nil
}

func (m *Map) PathNextStep(src, dst uint64) (uint64, error) {
	if src == dst || src == 0 || dst == 0 {
		return 0, errors.New("EINVAL")
	}

	next, ok := m.steps[vector{src, dst}]
	if ok {
		return next, nil
	}
	return 0, errors.New("No route")
}

func (m *Map) CellAdjacency(id uint64) []uint64 {
	adj := make([]uint64, 0)

	for i := m.Roads.First(id); i < len(m.Roads); i++ {
		r := m.Roads[i]
		if r.S != id {
			break
		}
		adj = append(adj, r.D)
	}

	return adj
}

func (m *Map) Dot() string {
	var sb strings.Builder
	sb.WriteString("digraph g {")
	for _, c := range m.Cells {
		sb.WriteString("n" + strconv.FormatUint(c.ID, 10))
		sb.WriteRune(';')
		sb.WriteRune('\n')
	}
	for _, r := range m.Roads {
		sb.WriteRune(' ')
		sb.WriteString("n" + strconv.FormatUint(r.S, 10))
		sb.WriteString(" -> ")
		sb.WriteString("n" + strconv.FormatUint(r.D, 10))
		sb.WriteRune(';')
		sb.WriteRune('\n')
	}
	sb.WriteString("}")
	return sb.String()
}

func (m *Map) Rehash() {
	next := make(map[vector]uint64)

	// Ensure the locations are sorted
	sort.Sort(&m.Roads)

	// Fill with the immediate neighbors
	for _, r := range m.Roads {
		next[vector{r.S, r.D}] = r.D
	}

	add := func(src, dst, step uint64) {
		_, found := next[vector{src, dst}]
		if !found {
			next[vector{src, dst}] = step
		}
	}

	// Call one DFS per node and shortcut when possible
	for _, cell := range m.Cells {
		already := make(map[uint64]bool)
		q := newQueue()

		// Bootstrap the DFS with adjacent nodes
		for _, next := range m.CellAdjacency(cell.ID) {
			q.push(next, next)
			already[next] = true
			// No need to add this in the known routes, we already did it
			// with an iteration on the roads (much faster)
		}

		for !q.empty() {
			current, first := q.pop()
			neighbors := m.CellAdjacency(current)
			// TODO(jfs): shuffle the neighbors
			for _, next := range neighbors {
				if !already[next] {
					// Avoid passing again in the neighbor
					already[next] = true
					// Tell to contine at that neighbor
					q.push(next, first)
					// We already learned the shortest path to that neighbor
					add(cell.ID, next, first)
				}
			}
		}
	}

	m.steps = next
}

type vector struct {
	src uint64
	dst uint64
}

type dfsTrack struct {
	current uint64
	first   uint64
}

type queue struct {
	tab   []dfsTrack
	start int
}

func newQueue() queue {
	var q queue
	q.tab = make([]dfsTrack, 0)
	return q
}

func (q *queue) push(node, first uint64) {
	q.tab = append(q.tab, dfsTrack{node, first})
}

func (q *queue) pop() (uint64, uint64) {
	v := q.tab[q.start]
	q.start++
	return v.current, v.first
}

func (q *queue) empty() bool {
	return q.start == len(q.tab)
}
