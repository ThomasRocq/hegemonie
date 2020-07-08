// Copyright (C) 2018-2020 Hegemonie's AUTHORS
// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at http://mozilla.org/MPL/2.0/.

package mapgraph

import (
	"github.com/jfsmig/hegemonie/pkg/region/model"
	"os"
	"testing"
	"time"
)

func TestMapInit(t *testing.T) {
	var m region.Map
	m.Init()

	if m.CellHas(1) {
		t.Fatal()
	}
	if m.getNextID() != 1 {
		t.Fatal()
	}

	cell := m.CellCreate()
	if cell.ID != 2 {
		t.Fatal()
	}
	if m.getNextID() != 3 {
		t.Fatal()
	}
	if !m.CellHas(cell.ID) {
		t.Fatal()
	}
}

func TestMapEinval(t *testing.T) {
	var m region.Map
	m.Init()

	// Test that identical, zero or inexistant locations return an error
	for _, src := range []uint64{0, 1, 2} {
		for _, dst := range []uint64{0, 1, 2} {
			if err := m.RoadCreate(src, dst, true); err == nil {
				t.Fatal()
			}
		}
	}
}

func TestMapMultiConnect(t *testing.T) {
	var err error
	var m region.Map
	m.Init()
	l0 := m.CellCreate()
	l1 := m.CellCreate()
	if err = m.RoadCreate(l0.ID, l1.ID, true); err != nil {
		t.Fatal()
	}
	if err = m.RoadCreate(l1.ID, l0.ID, true); err != nil {
		t.Fatal()
	}
	for i := 0; i < 5; i++ {
		if err = m.RoadCreate(l0.ID, l1.ID, true); err == nil {
			t.Logf("Cells %v", m.Cells)
			t.Logf("Roads %v", m.Roads)
			t.Fatal()
		}
		if err = m.RoadCreate(l1.ID, l0.ID, true); err == nil {
			t.Logf("Cells %v", m.Cells)
			t.Logf("Roads %v", m.Roads)
			t.Fatal()
		}
	}
}

func TestMapPathOneWay(t *testing.T) {
	var m region.Map
	m.Init()

	l0 := m.CellCreate()
	l1 := m.CellCreate()
	l2 := m.CellCreate()
	l3 := m.CellCreate()
	m.RoadCreate(l0.ID, l1.ID, true)
	m.RoadCreate(l1.ID, l2.ID, true)
	m.RoadCreate(l2.ID, l3.ID, true)

	m.Rehash()

	if step, err := m.PathNextStep(l0.ID, l3.ID); err != nil {
		t.Fatal()
	} else if step != l1.ID {
		t.Fatal()
	}

	if step, err := m.PathNextStep(l1.ID, l3.ID); err != nil {
		t.Fatal()
	} else if step != l2.ID {
		t.Fatal()
	}

	if step, err := m.PathNextStep(l2.ID, l3.ID); err != nil {
		t.Fatal()
	} else if step != l3.ID {
		t.Fatal()
	}

	if step, err := m.PathNextStep(l1.ID, l0.ID); err == nil {
		t.Fatal()
	} else if step != 0 {
		t.Fatal()
	}
}

func TestMapPathTwoWay(t *testing.T) {
	var m region.Map
	m.Init()

	biconnect := func(l0, l1 uint64) {
		m.RoadCreate(l0, l1, false)
		m.RoadCreate(l1, l0, false)
	}

	l0 := m.CellCreate()
	l1 := m.CellCreate()
	l2 := m.CellCreate()
	l3 := m.CellCreate()

	biconnect(l0.ID, l1.ID)
	biconnect(l1.ID, l2.ID)
	biconnect(l2.ID, l3.ID)

	m.Rehash()

	if step, err := m.PathNextStep(l3.ID, l0.ID); err != nil {
		t.Fatal()
	} else if step != l2.ID {
		t.Fatal()
	}

	if step, err := m.PathNextStep(l1.ID, l3.ID); err != nil {
		t.Fatal()
	} else if step != l2.ID {
		t.Fatal()
	}

	if step, err := m.PathNextStep(l2.ID, l3.ID); err != nil {
		t.Fatal()
	} else if step != l3.ID {
		t.Fatal()
	}
}

type grid struct {
	tab []uint64
	x   int
	y   int
}

func newGrid(x, y int) *grid {
	g := grid{x: x, y: y}
	g.tab = make([]uint64, x*y, x*y)
	return &g
}

func (g *grid) loc(i, j int) int {
	return i*g.y + j
}

func (g *grid) get(i, j int) uint64 {
	return g.tab[g.loc(i, j)]
}

func (g *grid) set(i, j int, v uint64) {
	g.tab[g.loc(i, j)] = v
}

func TestMapGrid(t *testing.T) {
	var m region.Map
	m.Init()

	x := 10
	y := 20
	t.Logf("Building the grid at %v", time.Now())
	grid := newGrid(x, y)
	for i := 0; i < x; i++ {
		for j := 0; j < y; j++ {
			v := m.CellCreate()
			grid.set(i, j, v.ID)
		}
	}

	offsets := []int{-1, 0, 1}
	neighbourhood := func(i, j int) {
		src := grid.get(i, j)
		// Iterate on the neighbors
		for _, id := range offsets {
			for _, jd := range offsets {
				if id == 0 && jd == 0 {
					// No self route
					continue
				}
				if id != 0 && jd != 0 {
					// No diagonals
					continue
				}
				if i+id < 0 || j+jd < 0 {
					// No underflow
					continue
				}
				if i+id >= x || j+jd >= y {
					// No overflow
					continue
				}
				dst := grid.get(i+id, j+jd)
				m.RoadCreateRaw(src, dst)
				m.RoadCreateRaw(dst, src)
			}
		}
	}

	// Even rows
	for i := 0; i < x; i += 2 {
		for j := 0; j < y; j += 2 {
			neighbourhood(i, j)
		}
	}
	// Odd rows
	for i := 1; i < x; i += 2 {
		for j := 1; j < y; j += 2 {
			neighbourhood(i, j)
		}
	}

	t.Logf("Rehashing at %v", time.Now())
	m.Rehash()

	t.Logf("Testing at %v", time.Now())

	dot := m.Dot()
	f, _ := os.Create("/tmp/dot")
	f.WriteString(dot)
	f.Close()

	t.Logf("Done at %v", time.Now())
}
