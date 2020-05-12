// Copyright (C) 2020 Hegemonie's AUTHORS
// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at http://mozilla.org/MPL/2.0/.

package main

import (
	"math/rand"
)

func MakeLine(n uint) *Graph {
	if n < 2 {
		panic("A line with length < 2 makes no sense")
	}

	g := G()
	for i := uint(0); i < n; i++ {
		v := V0()
		v.SetAnchor(0 == (i%2) || i == n-1)
		v.x = 0
		v.y = float64(i) * 50.0
		g.AddVertex(v)
	}

	vertices := g.GetVertices()
	for i := 1; i < len(vertices); i++ {
		g.AddEdge(E(vertices[i-1], vertices[i]))
	}

	return g
}

func MakeStar(r *rand.Rand, nb, radius uint) *Graph {
	if nb < 3 {
		panic("Star with less than 3 branches makes no sense")
	}
	if radius < 1 {
		panic("Star with radius zero makes no sense")
	}

	clusters := make([]*Graph, 0)
	for i := uint(0); i < nb; i++ {
		clusters = append(clusters, MakeLine(radius))
	}

	return GlueStar(r, clusters...)
}

func MakeCircle(n uint) *Graph {
	if n < 3 {
		panic("A line with length < 2 makes no sense")
	}

	g := G()
	for i := uint(0); i < n; i++ {
		v := V0()
		v.SetAnchor(0 == (i%2))
		g.AddVertex(v)
	}

	vertices := g.GetVertices()
	for i := uint(0); i < n; i++ {
		g.AddEdge(E(vertices[i], vertices[int(i+1)%len(vertices)]))
	}

	return g
}

func GlueStar(r *rand.Rand, gv ...*Graph) *Graph {
	if len(gv) < 2 {
		panic("Glueing less than 2 Graphs makes no sense")
	}

	cluster := Merge(gv...)
	a := V0()
	cluster.AddVertex(a)
	for _, g := range gv {
		b := PeekAnchor(r, g)
		cluster.AddEdge(E(a, b))
		b.SetAnchor(false)
	}
	a.SetAnchor(false)
	return cluster
}

func GlueChain(r *rand.Rand, gv ...*Graph) *Graph {
	if len(gv) < 2 {
		panic("Glueing less than 2 Graphs makes no sense")
	}

	rc := gv[0]
	for _, g := range gv[1:] {
		rc = Juxtapose(rc, g)
	}
	return rc
}

func Juxtapose(g0, g1 *Graph) (g *Graph) {
	Origin(g0)
	Origin(g1)

	// Vertical alignment on the middle of each graph
	_, xmax0, _, ymax0 := Box(g0)
	_, _, _, ymax1 := Box(g1)
	if ymax0 > ymax1 {
		Translate(g1, 0, (ymax1-ymax0)/2)
	} else if ymax0 < ymax1 {
		Translate(g0, 0, (ymax0-ymax1)/2)
	}

	Translate(g1, xmax0+50, 0)

	// Add  a link between both graphs
	rc := Merge(g0, g1)
	rc.AddEdge(E(RightMost(g0), LeftMost(g1)))
	return rc
}

func Loop(r *rand.Rand, g *Graph) {
	a := PeekAnchor(r, g)
	a.SetAnchor(false)
	b := PeekAnchor(r, g)
	b.SetAnchor(false)
	g.AddEdge(E(a, b))
}

func PeekAnchor(r *rand.Rand, g *Graph) *Vertex {
	anchors := g.GetAnchors()
	return anchors[r.Intn(len(anchors))]
}
