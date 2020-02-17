// Copyright (C) 2020 Hegemonie's AUTHORS
// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at http://mozilla.org/MPL/2.0/.

package main

func FDP(p Physics, g Graph, centerX, centerY float64, rounds uint) {
	for i := uint(0); i < rounds; i++ {
		fdpRound(p, g, centerX, centerY)
	}
}

type pair struct{ x, y uint64 }

func fdpVertices(p Physics, g Graph) {
	done := make(map[pair]bool)
	// O(N^2) Nodes repulsion
	for _, v0 := range g.GetVertices() {
		for _, v1 := range g.GetVertices() {
			if v0 == v1 || done[pair{v0.id, v1.id}] || done[pair{v1.id, v0.id}] {
				continue
			}
			dist, force := p.VertexForce(v0, v1)
			//log.Printf("Vertices %v--%v dist=%v force=%v  {%.3f,%.3f}--{%.3f,%.3f}", v0.id, v1.id, dist, force, v0.x, v0.y, v1.x, v1.y)
			p.ApplyForce(v0, v1, dist, force)
			done[pair{v0.id, v1.id}] = true
		}
	}
}

func fdpEdges(p Physics, g Graph) {
	// O(M) Edges forces
	for _, e := range g.GetEdges() {
		dist, force := p.EdgeForce(e.S(), e.D())
		//log.Printf("Edges %v--%v dist=%v force=%v  {%.3f,%.3f}--{%.3f,%.3f}", e.S().id, e.D().id, dist, force, e.S().x, e.S().y, e.D().x, e.D().y)
		p.ApplyForce(e.S(), e.D(), dist, force)
	}
}

func fdpApply(g Graph) {
	// O(N) Apply the force
	for _, v := range g.GetVertices() {
		v.Move()
	}
}

func fdpCenter(p Physics, g Graph, centerX, centerY float64) {
	// O(N) Slightly pushing toward the center
	/*
		for _, v := range g.GetVertices() {
			center := Vertex{x: centerX, y: centerY}
			dist, delta := p.VertexForce(&center, v)
			p.ApplyForce(&center, v, dist, delta)
		}
	*/
}

func fdpReset(g Graph) {
	// O(N) Apply the force
	for _, v := range g.GetVertices() {
		v.Reset()
	}
}

func fdpRound(p Physics, g Graph, centerX, centerY float64) {
	fdpReset(g)
	fdpCenter(p, g, centerY, centerY)
	fdpVertices(p, g)
	fdpEdges(p, g)
	fdpApply(g)
}
