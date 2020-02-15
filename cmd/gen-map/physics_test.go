// Copyright (C) 2020 Hegemonie's AUTHORS
// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at http://mozilla.org/MPL/2.0/.

package main

import "testing"

func p() Physics {
	return &physics{edgeBalance: 10.0, vertexBalance: 10.0}
}

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
	p := p()
	for _, e := range expectations {
		if e.d != p.Distance(e.v0, e.v1) {
			t.Fatal()
		}
	}
}

func TestVertex(t *testing.T) {
	expectations := []floatExpectation{
		// Ideal position
		{0.0, V(0, 10), V(0, 0)},
		{0.0, V(0, 0), V(0, 10)},
		{0.0, V(10, 0), V(0, 0)},
		{0.0, V(0, 0), V(10, 0)},
		// Maximal repulsion
		{9.0, V(0, 0), V(0, 0)},
	}
	p := p()
	for _, e := range expectations {
		_, v := p.VertexForce(e.v0, e.v1)
		if e.d != v {
			t.Fatalf("VertexForce expected=%v got=%v v0=%v v1=%v", e.d, v, e.v0, e.v1)
		}
	}
}

func TestEdge(t *testing.T) {
	expectations := []floatExpectation{
		// Ideal position
		{0.0, V(0, 10), V(0, 0)},
		{0.0, V(0, 0), V(0, 10)},
		{0.0, V(10, 0), V(0, 0)},
		{0.0, V(0, 0), V(10, 0)},
		// Maximal repulsion
		{5.0, V(0, 0), V(0, 0)},
	}
	p := p()
	for _, e := range expectations {
		_, v := p.EdgeForce(E(e.v0, e.v1))
		if e.d != v {
			t.Fatalf("EdgeForce expected=%v got=%v v0=%v v1=%v", e.d, v, e.v0, e.v1)
		}
	}
}

func TestDelta(t *testing.T) {

}
