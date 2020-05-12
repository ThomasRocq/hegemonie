// Copyright (C) 2020 Hegemonie's AUTHORS
// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at http://mozilla.org/MPL/2.0/.

package main

import (
	"math"
	"strconv"
)

type Graph struct {
	vertices []*Vertex
	edges    []*Edge
}

type Vertex struct {
	id             int
	uid            uint64
	x, y           float64
	center, anchor bool
}

type Edge struct {
	s, d *Vertex
}

func G() *Graph { return &Graph{} }

func E(s, d *Vertex) *Edge { return &Edge{s: s, d: d} }
func (e *Edge) S() *Vertex { return e.s }
func (e *Edge) D() *Vertex { return e.d }

func V0() *Vertex                  { return &Vertex{uid: getNextId()} }
func V(x, y float64) *Vertex       { return &Vertex{uid: getNextId(), x: x, y: y} }
func (v *Vertex) Id() string       { return strconv.FormatUint(v.uid, 16) }
func (v *Vertex) IsAnchor() bool   { return v.anchor }
func (v *Vertex) IsCenter() bool   { return v.center }
func (v *Vertex) SetAnchor(b bool) { v.anchor = b }
func (v *Vertex) SetCenter(b bool) { v.center = b }
func (v *Vertex) Distance(o *Vertex) float64 {
	return math.Sqrt(math.Pow(v.x-o.x, 2) + math.Pow(v.y-o.y, 2))
}

func (g *Graph) AddEdges(e ... *Edge) *Graph {
	g.edges = append(g.edges, e...)
	return g
}

func (g *Graph) AddVertices(v ... *Vertex) *Graph {
	g.vertices = append(g.vertices, v...)
	return g
}

func (g *Graph) AddEdge(e *Edge)        { g.edges = append(g.edges, e) }
func (g *Graph) AddVertex(v *Vertex)    { g.vertices = append(g.vertices, v) }
func (g *Graph) GetEdges() []*Edge      { return g.edges[:] }
func (g *Graph) GetVertices() []*Vertex { return g.vertices[:] }
func (g *Graph) GetAnchors() []*Vertex {
	rc := make([]*Vertex, 0)
	for _, v := range g.vertices {
		if v.IsAnchor() {
			rc = append(rc, v)
		}
	}
	return rc
}
