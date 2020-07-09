// Copyright (C) 2018-2020 Hegemonie's AUTHORS
// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at http://mozilla.org/MPL/2.0/.

package mapgraph

import (
	"errors"
	"sort"
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
