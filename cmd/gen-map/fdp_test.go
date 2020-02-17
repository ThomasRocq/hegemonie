// Copyright (C) 2020 Hegemonie's AUTHORS
// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at http://mozilla.org/MPL/2.0/.

package main

import (
	"fmt"
	"math/rand"
	"strings"
	"testing"
)

func TestForceBalance(t *testing.T) {
	r := rand.New(rand.NewSource(0))
	g := MakeCircle(3)
	Noise(r, g, 200, 200)

	p := makePhysics(r)
	dump := func() {
		var sb strings.Builder
		for _, e := range g.GetEdges() {
			sb.WriteString(fmt.Sprintf(" D=%.3f", p.Distance(e.S(), e.D())))
		}
		for _, v := range g.GetVertices() {
			sb.WriteString(fmt.Sprintf(" {%.3f+%.3f,%.3f+%.3f}", v.x, v.vx, v.y, v.vy))
		}
		t.Log(sb.String())
	}
	for i := 0; i < 100; i++ {
		dump()
		FDP(p, g, 500, 500, 1)
	}
	dump()
}

func TestForceProgress(t *testing.T) {
	r := rand.New(rand.NewSource(0))
	p := makePhysics(r)
	for i := float64(0); i < 200.0; i += 1.0 {
		v0 := V(0, 0)
		v1 := V(3*(i/5.0), 4*(i/5.0))

		var dist, force float64
		dist, force = p.VertexForce(v0, v1)
		p.ApplyForce(v0, v1, dist, force)
		dist, force = p.EdgeForce(v0, v1)
		p.ApplyForce(v0, v1, dist, force)
		v0.Move()
		v1.Move()
		t.Logf("force=%.3f dist=%.3f   v0{%.3f+%.3f, %.3f+%.3f}   v1{%.3f+%.3f, %.3f+%.3f}", force, dist, v0.x, v0.vx, v0.y, v0.vy, v1.x, v1.vx, v1.y, v1.vy)
	}
}
