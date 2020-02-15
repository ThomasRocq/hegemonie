// Copyright (C) 2020 Hegemonie's AUTHORS
// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at http://mozilla.org/MPL/2.0/.

package main

import (
	"math/rand"
)

func Cluster(g ...Graph) Graph {
	rc := &clusteredGraph{}
	rc.clusters = append(rc.clusters, g...)
	return rc
}

func MakeCircle(n uint) Graph {
	if n < 3 {
		panic("A line with length < 2 makes no sense")
	}

	vertices := make([]*Vertex, 0)
	for i := uint(0); i < n; i++ {
		vertices = append(vertices, V0())
	}
	for i, v := range vertices {
		if (i % 2) == 0 {
			v.SetAnchor(true)
		}
	}

	edges := make([]*Edge, 0)
	for i := uint(0); i < n; i++ {
		e := E(vertices[i], vertices[int(i+1)%len(vertices)])
		edges = append(edges, e)
	}

	return &memGraph{vertices: vertices, edges: edges}
}

func MakeLine(n uint) Graph {
	if n < 2 {
		panic("A line with length < 2 makes no sense")
	}

	vertices := make([]*Vertex, 0)
	for i := uint(0); i < n; i++ {
		vertices = append(vertices, V0())
	}
	for i, v := range vertices {
		if (i % 2) == 0 {
			v.SetAnchor(true)
		}
	}
	vertices[n-1].SetAnchor(true)

	edges := make([]*Edge, 0)
	for i := uint(1); i < n; i++ {
		e := E(vertices[i-1], vertices[i])
		edges = append(edges, e)
	}

	return &memGraph{vertices: vertices, edges: edges}
}

func MakeStar(r *rand.Rand, nb, radius uint) Graph {
	if nb < 3 {
		panic("Star with less than 3 branches makes no sense")
	}
	if radius < 1 {
		panic("Star with radius zero makes no sense")
	}

	clusters := make([]Graph, 0)
	for i := uint(0); i < nb; i++ {
		clusters = append(clusters, MakeLine(radius))
	}

	return GlueStar(r, clusters...)
}

func GlueStar(r *rand.Rand, gv ...Graph) Graph {
	if len(gv) < 2 {
		panic("Glueing less than 2 Graphs makes no sense")
	}

	cluster := Cluster(gv...)
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

func GlueChain(r *rand.Rand, gv ...Graph) Graph {
	if len(gv) < 2 {
		panic("Glueing less than 2 Graphs makes no sense")
	}

	rc := Cluster(gv[0])
	for _, g := range gv[1:] {
		a := PeekAnchor(r, rc)
		a.SetAnchor(false)
		b := PeekAnchor(r, g)
		b.SetAnchor(false)
		rc.AddEdge(E(a, b))
		rc.(*clusteredGraph).Add(g)
	}
	return rc
}

func Loop(r *rand.Rand, g Graph) {
	a := PeekAnchor(r, g)
	a.SetAnchor(false)
	b := PeekAnchor(r, g)
	b.SetAnchor(false)
	g.AddEdge(E(a, b))
}

func PeekAnchor(r *rand.Rand, g Graph) *Vertex {
	anchors := g.GetAnchors()
	return anchors[r.Intn(len(anchors))]
}
