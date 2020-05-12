// Copyright (C) 2020 Hegemonie's AUTHORS
// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at http://mozilla.org/MPL/2.0/.

package main

// #cgo pkg-config: igraph gsl
// #include <igraph.h>
// #include <stdio.h>
// #include "binding.h"
import "C"

import (
	"errors"
	"log"
)

func FDP(g *Graph, x, y float64, rounds uint) error {
	v0 := g.GetVertices()
	for i, v := range v0 {
		v.id = i
	}
	e0 := g.GetEdges()

	// Print some debug
	log.Println("V", len(v0), "E", len(e0))
	for _, v := range v0 { log.Println("V>", "id", v.id, "uid", v.uid) }
	for _, e := range e0 { log.Println("E>", "from", e.S().id, "to", e.D().id) }

	var vertices *C.vertex_array_t
	vertices = C.vertex_array_create(C.uint32_t(len(v0)))
	defer C.vertex_array_destroy(vertices)
	for idx, v0 := range v0 {
		v0.id = idx
		var v C.vertex_t
		v.x = C.double_t(v0.x)
		v.y = C.double_t(v0.y)
		C.vertex_array_set(vertices, C.uint32_t(idx), v)
	}

	var edges *C.edge_array_t
	edges = C.edge_array_create(C.uint32_t(0))
	defer C.edge_array_destroy(edges)
	for _, e := range e0 {
		C.edge_array_add(edges, C.uint32_t(e.S().id), C.uint32_t(e.D().id))
	}

	rc := C.igraph_fdp(vertices, edges, C.double_t(x), C.double_t(y), C.uint32_t(rounds))
	if rc == 0 {
		for idx, v := range v0 {
			v.x = float64(C.vertex_array_x(vertices, C.uint32_t(idx)))
			v.y = float64(C.vertex_array_y(vertices, C.uint32_t(idx)))
		}
	}

	if rc == 0 {
		return nil
	}

	return errors.New("FDP computation error")
}
