// Copyright (C) 2020 Hegemonie's AUTHORS
// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at http://mozilla.org/MPL/2.0/.

package main

import (
	"math"
	"math/rand"
	"testing"
)

type floatExpectation struct {
	d      float64
	v0, v1 *Vertex
}

func TestDistance(t *testing.T) {
	expectations := []floatExpectation{
		{1.0, V(0, 1), V(0, 0)},
		{1.0, V(0, 0), V(0, 1)},
		{1.0, V(1, 0), V(0, 0)},
		{1.0, V(0, 0), V(1, 0)},
		{0.0, V(0, 0), V(0, 0)},
	}
	p := makePhysics(rand.New(rand.NewSource(0)))
	for idx, e := range expectations {
		d := p.Distance(e.v0, e.v1)
		if e.d != d {
			t.Fatalf("idx=%v expected=%v got=%v v0=%v v1=%v", idx, e.d, d, e.v0, e.v1)
		}
	}
}

func TestVertex(t *testing.T) {
	expectations := []floatExpectation{
		// Ideal position
		{0.0, V(0, 20), V(0, 0)},
		{0.0, V(20, 0), V(0, 0)},
		{0.0, V(0, 0), V(0, 20)},
		{0.0, V(0, 0), V(20, 0)},
		// Maximal repulsion
		{12.0, V(0, 0), V(0, 0)},
	}
	p := makePhysics(rand.New(rand.NewSource(0)))
	for idx, e := range expectations {
		_, v := p.VertexForce(e.v0, e.v1)
		if e.d != v {
			t.Fatalf("idx=%v expected=%v got=%v v0=%v v1=%v", idx, e.d, v, e.v0, e.v1)
		}
	}
}

func TestEdge(t *testing.T) {
	expectations := []floatExpectation{
		// Ideal position
		{0.0, V(0, 50), V(0, 0)},
		{0.0, V(50, 0), V(0, 0)},
		{0.0, V(0, 0), V(0, 50)},
		{0.0, V(0, 0), V(50, 0)},
		// Maximal repulsion
		{20.0, V(0, 0), V(0, 0)},
	}
	p := makePhysics(rand.New(rand.NewSource(0)))
	for idx, e := range expectations {
		_, v := p.EdgeForce(e.v0, e.v1)
		if e.d != v {
			t.Fatalf("idx=%v expected=%v got=%v v0=%v v1=%v", idx, e.d, v, e.v0, e.v1)
		}
	}
}

type deltaExpectation struct {
	// Expected coordinates variation for the point that is not at <0,0>
	dx, dy float64

	// Relative coordinates of the point that is not at <0,0>
	x, y float64

	// Distance between points.
	dist float64

	// Distance delta (e.g. induced by a force) to be applied
	delta float64
}

func eq(v0, v1 float64) bool {
	return math.Abs(v0-v1) <= epsilon
}

func TestDelta(t *testing.T) {
	p := makePhysics(rand.New(rand.NewSource(0)))
	E := func(ex, ey, x, y, delta float64) deltaExpectation {
		return deltaExpectation{ex, ey, x, y, p.DistanceXY(x, y), delta}
	}
	expectations := []deltaExpectation{
		// Zero delta
		E(0, 0, 10, 0, 0),
		E(0, 0, 0, 10, 0),
		// Double distance along the axes
		E(10, 0, 10, 0, 10),
		E(0, 10, 0, 10, 10),
		// Double distance with non-special angle
		E(3.0, 4.0, 3.0, 4.0, 5.0),
		// Non-zero delta with zero distance
		E(math.NaN(), math.NaN(), 0.0, 0.0, 5.0),
	}

	for i, e := range expectations {
		dx, dy := p.GetDeltaXY(e.x, e.y, e.dist, e.delta)
		if math.IsNaN(e.dx) || math.IsNaN(e.dy) {
			if math.Abs(p.DistanceXY(dx, dy)-e.delta) > epsilon {
				t.Fatalf(
					"i=%d dx=%v dy=%v Expect{dx:%v, dy:%v, x:%v, y:%v, dist:%v, delta:%v}",
					i, dx, dy, e.dx, e.dy, e.x, e.y, e.dist, e.delta)
			}
		} else {
			if !eq(e.dx, dx) || !eq(e.dy, dy) {
				t.Fatalf(
					"i=%d dx=%v dy=%v Expect{dx:%v, dy:%v, x:%v, y:%v, dist:%v, delta:%v}",
					i, dx, dy, e.dx, e.dy, e.x, e.y, e.dist, e.delta)
			}
		}
	}
}
