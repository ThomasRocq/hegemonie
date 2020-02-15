// Copyright (C) 2020 Hegemonie's AUTHORS
// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at http://mozilla.org/MPL/2.0/.

package main

import (
	"fmt"
	"io"
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

func V0() *Vertex                  { return &Vertex{id: getNextId()} }
func V(x, y float64) *Vertex       { return &Vertex{id: getNextId(), x: x, y: y} }
func (v *Vertex) Id() string       { return strconv.FormatUint(v.id, 16) }
func (v *Vertex) IsAnchor() bool   { return v.anchor }
func (v *Vertex) IsCenter() bool   { return v.center }
func (v *Vertex) SetAnchor(b bool) { v.anchor = b }
func (v *Vertex) SetCenter(b bool) { v.center = b }

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
<!DOCTYPE svg PUBLIC "-//W3C//DTD SVG 1.1//EN"
 "http://www.w3.org/Graphics/SVG/1.1/DTD/svg11.dtd">
<svg
 width="%[1]d" height="%[2]d"
 viewBox="0.00 0.00 %[1]d %[2]d'"
 xmlns="http://www.w3.org/2000/svg"
 xmlns:xlink="http://www.w3.org/1999/xlink">
<g>
<rect x="0" y="0" width="%[1]d" height="%[2]d" fill="#E0E0E0" fill-opacity="0.3" stroke="none"/>`, x, y)
	o.Write(CR)

	for _, e := range g.GetEdges() {
		a := e.S()
		b := e.D()
		fmt.Fprintf(o, `<line x1="%f" y1="%f" x2="%f" y2="%f" stroke="#000000"/>`, a.x, a.y, b.x, b.y)
		o.Write(CR)
	}
	for _, v := range g.GetVertices() {
		fmt.Fprintf(o, `<circle cx="%f" cy="%f" r="5" fill="#FFFFFF" stroke="#000000"/>`, v.x, v.y)
		o.Write(CR)
	}

	fmt.Fprint(o, `</g></svg>`)
}
