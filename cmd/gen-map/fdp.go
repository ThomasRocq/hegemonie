package main

import (
	"math/rand"
)

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

	// O(N^2) Nodes repulsion
	for _, v0 := range g.GetVertices() {
		//log.Printf("%p %v %v", v0, v0.vx, v0.vy)
		for _, v1 := range g.GetVertices() {
			if v0 == v1 {
				continue
			}
		}
	}

	// Apply the force
	for _, v := range g.GetVertices() {
		v.x += v.vx
		v.y += v.vy
	}
}
