// Copyright (C) 2020 Hegemonie's AUTHORS
// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at http://mozilla.org/MPL/2.0/.

package main

import (
	"math"
	"math/rand"
)

const (
	epsilon = float64(0.000001)
)

type Physics interface {
	// Compute the distance between both vertices
	DistanceXY(x, y float64) float64

	// Compute the distance between both vertices
	Distance(v0, v1 *Vertex) float64

	// Compute the variation of <x,y> based on the distance variation (given by `delta`)
	GetDeltaXY(x, y float64, dist, delta float64) (float64, float64)

	// Update in place the position of v0 and v1 according to the target distance delta.
	ApplyForce(v0, v1 *Vertex, dist, delta float64)

	// Returns the actual distance and the next distance goal based on the Vertices interaction
	VertexForce(v0, v1 *Vertex) (float64, float64)

	// Returns the actual distance and the next distance goal based on the Edge interaction
	EdgeForce(v0, v1 *Vertex) (float64, float64)
}

type physics struct {
	r *rand.Rand

	// Distance threshold below wich the repulsion is triggered
	vertexBalance float64

	// Distance threshold below wich the repulsion is triggered
	edgeBalance float64
}

func makePhysics(r *rand.Rand) Physics {
	return &physics{r: r, edgeBalance: 50.0, vertexBalance: 20.0}
}

func (p *physics) DistanceXY(x, y float64) float64 {
	return math.Sqrt(math.Pow(x, 2) + math.Pow(y, 2))
}

func (p *physics) Distance(v0, v1 *Vertex) float64 {
	return p.DistanceXY(v1.x-v0.x, v1.y-v0.y)
}

func (p *physics) VertexForce(v0, v1 *Vertex) (dist, f float64) {
	dist = p.Distance(v0, v1)

	if dist < p.vertexBalance {
		// ultra strong Repulsion (Strong Interaction)
		f = 2.0 * (p.vertexBalance - dist)
	} else if dist > p.vertexBalance {
		// weaker attraction (Gravity)
		f = -100000.0 / math.Pow(p.vertexBalance-dist, 2.0)
	}

	return dist, f
}

func (p *physics) EdgeForce(v0, v1 *Vertex) (dist, f float64) {
	dist = p.Distance(v0, v1)

	if dist < p.edgeBalance {
		// very strong Repulsion
		f = 1 * (p.edgeBalance - dist)
	} else if dist > p.edgeBalance {
		// strong attraction
		f = 1 * (p.edgeBalance - dist)
	}

	return dist, f
}

func (p *physics) GetDeltaXY(x, y float64, dist, delta float64) (float64, float64) {
	x, y = math.Abs(x), math.Abs(y)
	dist = math.Abs(dist)

	if math.Abs(delta) <= epsilon {
		return 0, 0
	}

	if dist <= epsilon || (x <= epsilon && y <= epsilon) {
		theta := math.Pi * p.r.Float64() / 2.0
		sin := math.Sin(theta)
		cos := math.Cos(theta)
		return delta * cos, delta * sin
	} else if x <= epsilon {
		return 0.0, delta
	} else if y <= epsilon {
		return delta, 0.0
	} else {
		ratio := delta / dist
		return x * ratio, y * ratio
	}
}

func (p *physics) ApplyForce(v0, v1 *Vertex, dist, force float64) {
	dx, dy := p.GetDeltaXY(v0.x-v1.x, v0.y-v1.y, dist, force)
	dx2, dy2 := dx/2.0, dy/2.0
	if v0.x > v1.x {
		v0.vx += dx2
		v1.vx -= dx2
	} else {
		v1.vx += dx2
		v0.vx -= dx2
	}
	if v0.y > v1.y {
		v0.vy += dy2
		v1.vy -= dy2
	} else {
		v1.vy += dy2
		v0.vy -= dy2
	}
	//log.Println(*v0, *v1)
}
