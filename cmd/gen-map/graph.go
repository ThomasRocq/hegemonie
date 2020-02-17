// Copyright (C) 2020 Hegemonie's AUTHORS
// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at http://mozilla.org/MPL/2.0/.

package main

import (
	"fmt"
	"io"
	"math/rand"
	"strconv"
)

type Vertex struct {
	id     uint64
	center bool
	anchor bool

	x, y   float64
	vx, vy float64
}

type Edge struct {
	s, d *Vertex
}

type Graph interface {
	AddVertex(e *Vertex)
	AddEdge(v *Edge)

	GetVertices() []*Vertex
	GetEdges() []*Edge

	GetAnchors() []*Vertex

	DrawDot(o io.Writer)
}

type memGraph struct {
	vertices []*Vertex
	edges    []*Edge
}

type clusteredGraph struct {
	memGraph
	clusters []Graph
}

func E(s, d *Vertex) *Edge { return &Edge{s: s, d: d} }
func (e *Edge) S() *Vertex { return e.s }
func (e *Edge) D() *Vertex { return e.d }

const deltaMax = float64(2.0)

func absMax(v float64) float64 {
	if v < -deltaMax {
		return -deltaMax
	}
	if v > deltaMax {
		return deltaMax
	}
	return v
}

func V0() *Vertex                  { return &Vertex{id: getNextId()} }
func V(x, y float64) *Vertex       { return &Vertex{id: getNextId(), x: x, y: y} }
func (v *Vertex) Id() string       { return strconv.FormatUint(v.id, 16) }
func (v *Vertex) IsAnchor() bool   { return v.anchor }
func (v *Vertex) IsCenter() bool   { return v.center }
func (v *Vertex) SetAnchor(b bool) { v.anchor = b }
func (v *Vertex) SetCenter(b bool) { v.center = b }
func (v *Vertex) Move()            { v.x, v.y = absMax(v.x+v.vx), absMax(v.y+v.vy) }
func (v *Vertex) Reset()           { v.vx, v.vy = 0, 0 }

func (g *memGraph) AddEdge(e *Edge)        { g.edges = append(g.edges, e) }
func (g *memGraph) AddVertex(v *Vertex)    { g.vertices = append(g.vertices, v) }
func (g *memGraph) GetEdges() []*Edge      { return g.edges[:] }
func (g *memGraph) GetVertices() []*Vertex { return g.vertices[:] }
func (g *memGraph) GetAnchors() []*Vertex {
	rc := make([]*Vertex, 0)
	for _, v := range g.vertices {
		if v.IsAnchor() {
			rc = append(rc, v)
		}
	}
	return rc
}

func (cg *clusteredGraph) Add(g Graph)         { cg.clusters = append(cg.clusters, g) }
func (cg *clusteredGraph) AddEdge(e *Edge)     { cg.edges = append(cg.edges, e) }
func (cg *clusteredGraph) AddVertex(v *Vertex) { cg.vertices = append(cg.vertices, v) }
func (cg *clusteredGraph) GetEdges() []*Edge {
	rc := make([]*Edge, 0)
	for _, g := range cg.clusters {
		rc = append(rc, g.GetEdges()...)
	}
	return rc
}

func (cg *clusteredGraph) GetVertices() []*Vertex {
	rc := make([]*Vertex, 0)
	for _, g := range cg.clusters {
		rc = append(rc, g.GetVertices()...)
	}
	return rc
}

func (cg *clusteredGraph) GetAnchors() []*Vertex {
	rc := make([]*Vertex, 0)
	for _, g := range cg.clusters {
		rc = append(rc, g.GetAnchors()...)
	}
	return rc
}

func (cg *clusteredGraph) DrawDot(o io.Writer) {
	for _, g := range cg.clusters {
		g.DrawDot(o)
	}
	for _, e := range cg.edges {
		fmt.Fprintf(o, "\"%s\" -- \"%s\";\n", e.S().Id(), e.D().Id())
	}
}

func (g *memGraph) DrawDot(o io.Writer) {
	for _, e := range g.edges {
		fmt.Fprintf(o, "\"%s\" -- \"%s\";\n", e.S().Id(), e.D().Id())
	}
}

func DrawDot(g Graph, o io.Writer) {
	fmt.Fprint(o, "graph g {\n")
	for _, v := range g.GetAnchors() {
		if v.IsAnchor() {
			fmt.Fprintf(o, "\"%s\" [color=red];", v.Id())
		}
	}
	g.DrawDot(o)
	fmt.Fprint(o, "}\n")
}

func DrawSvg(g Graph, o io.Writer, x, y uint) {
	CR := []byte{'\n'}

	fmt.Fprintf(o, `<?xml version="1.0" encoding="UTF-8" standalone="no"?>
<svg width="%[1]d" height="%[2]d" viewBox="0 0 %[1]d %[2]d"
 xmlns="http://www.w3.org/2000/svg" xmlns:xlink="http://www.w3.org/1999/xlink">
<g>
<rect x="0" y="0" width="%[1]d" height="%[2]d" fill="#E0E0E0" fill-opacity="0.3" stroke="none"/>`, x, y)
	o.Write(CR)

	f := func(f float64) uint64 { return uint64(f) }
	for _, e := range g.GetEdges() {
		a := e.S()
		b := e.D()
		fmt.Fprintf(o, `<line x1="%d" y1="%d" x2="%d" y2="%d" stroke="#000000"/>`, f(a.x), f(a.y), f(b.x), f(b.y))
		o.Write(CR)
	}
	for _, v := range g.GetVertices() {
		fmt.Fprintf(o, `<circle cx="%d" cy="%d" r="5" fill="#FFFFFF" stroke="#000000"/>`, f(v.x), f(v.y))
		o.Write(CR)
	}

	fmt.Fprint(o, `</g></svg>`)
}

func Noise(r *rand.Rand, g Graph, x, y float64) {
	for _, v := range g.GetVertices() {
		v.x += r.Float64() * x
		v.y += r.Float64() * y
	}
}

// Reimplement the graph interface with a compact AdjacencyList implementation
func Simplify(g Graph) Graph {
	return &memGraph{vertices: g.GetVertices(), edges: g.GetEdges()}
}

// Refactor each node position to fit in a given rectangle.
func Normalize(g Graph, x, y float64) {
	var xmin, xmax, ymin, ymax float64 = 0, 0, 0, 0
	for _, v := range g.GetVertices() {
		if v.x < xmin {
			xmin = v.x
		}
		if v.x > xmax {
			xmax = v.x
		}
		if v.y < ymin {
			ymin = v.y
		}
		if v.y > ymax {
			ymax = v.y
		}
	}

	// Translate everything so that <xmin,ymin> is at <0,0>
	xmax, ymax = xmax-xmin, ymax-ymin
	for _, v := range g.GetVertices() {
		v.x, v.y = v.x-xmin, v.y-ymin
	}

	// Stretch everything so that <xmax,ymax> is at <x,y>
	var xratio, yratio float64 = x / xmax, y / ymax
	for _, v := range g.GetVertices() {
		v.x, v.y = v.x*xratio, v.y*yratio
	}
}
