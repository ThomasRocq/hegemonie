// Copyright (C) 2020 Hegemonie's AUTHORS
// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at http://mozilla.org/MPL/2.0/.

package main

import (
	"math"
	"math/rand"
)

type Physics interface {
	// Compute the distance between both vertices
	Distance(v0, v1 *Vertex) float64

	// Returns the actual distance and the next distance goal based on the Vertices interaction
	VertexForce(v0, v1 *Vertex) (float64, float64)

	// Returns the actual distance and the next distance goal based on the Edge interaction
	EdgeForce(e *Edge) (float64, float64)
}

type physics struct {
	r *rand.Rand

	// Distance threshold below wich the repulsion is triggered
	vertexBalance float64

	// Distance threshold below wich the repulsion is triggered
	edgeBalance float64
}

func Noise(r *rand.Rand, g Graph, x, y float64) {
	for _, v := range g.GetVertices() {
		v.x += r.Float64() * x
		v.y += r.Float64() * y
	}
}

func Simplify(g Graph) Graph {
	return &memGraph{vertices: g.GetVertices(), edges: g.GetEdges()}
}

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

	// Stretch everything so that
	var xratio, yratio float64 = x / xmax, y / ymax
	for _, v := range g.GetVertices() {
		v.x, v.y = v.x*xratio, v.y*yratio
	}
}

func FDP(r *rand.Rand, g Graph, centerX, centerY float64) {
	for i := 0; i < 1000; i++ {
		fdpRound(r, g, centerX, centerY)
	}
}

func fdpRound(r *rand.Rand, g Graph, centerX, centerY float64) {
	// Reset the force resultants
	for _, v := range g.GetVertices() {
		v.vx, v.vy = 0, 0
	}

	// O(N) Slightly pushing toward the center
	for _, v := range g.GetVertices() {
		v.vx -= (centerX - v.x) / 100.0
		v.vy -= (centerY - v.y) / 100.0
	}

	// O(N) Nodes repulsion
	for _, v0 := range g.GetVertices() {
		//log.Printf("%p %v %v", v0, v0.vx, v0.vy)
		for _, v1 := range g.GetVertices() {
			if v0 == v1 {
				continue
			}

		}
	}

	// O(M) Edges forces
	for _, e := range g.GetEdges() {
		if e == nil {
			continue
		}
	}

	// Apply the force
	for _, v := range g.GetVertices() {
		v.x += v.vx
		v.y += v.vy
	}
}

func (p physics) Distance(v0, v1 *Vertex) float64 {
	x2 := math.Pow(v1.x-v0.x, 2)
	y2 := math.Pow(v1.y-v0.y, 2)
	return math.Sqrt(x2 + y2)
}

func (p physics) VertexForce(v0, v1 *Vertex) (float64, float64) {
	dist := p.Distance(v0, v1)

	delta := float64(0.0)
	if dist < p.vertexBalance {
		// very strong Repulsion
		delta = (p.vertexBalance - dist) * 0.9
	} else if dist > p.vertexBalance {
		// weaker attraction
		delta = -(dist - p.vertexBalance) / 2.0
	}

	return dist, delta
}

func (p physics) EdgeForce(e *Edge) (float64, float64) {
	dist := p.Distance(e.S(), e.D())

	delta := float64(0.0)
	if dist < p.vertexBalance {
		// strong Repulsion
		delta = (p.vertexBalance - dist) * 0.5
	} else if dist > p.vertexBalance {
		// weaker attraction
		delta = -(dist - p.edgeBalance) / 3.0
	}

	return dist, delta
}

func applyDeltaNorm(dx, dy float64, dist, delta float64) (float64, float64) {
	theta := math.Acos(dist / dx)
	dist += delta
	return (dist / math.Cos(theta)) - dx, (dist / math.Sin(theta)) - dy
}

func applyDelta(v0, v1 *Vertex, dist, delta float64) {
	dx, dy := applyDeltaNorm(math.Abs(v0.x-v1.x), math.Abs(v0.y-v1.y), dist, delta)
	dx2, dy2 := dx/2.0, dy/2.0
	v0.vx += dx2
	v1.vx -= dx2
	v0.vy += dy2
	v1.vy -= dy2
}
