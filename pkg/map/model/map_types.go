// Copyright (C) 2018-2020 Hegemonie's AUTHORS
// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at http://mozilla.org/MPL/2.0/.

package mapgraph

// A MapEdge is an edge if the transportation directed graph
type MapEdge struct {
	// Unique identifier of the source Cell
	S uint64 `json:"src"`

	// Unique identifier of the destination Cell
	D uint64 `json:"dst"`
}

// A MapVertex is a vertex in the transportation directed graph
type MapVertex struct {
	// The unique identifier of the current cell.
	ID uint64 `json:"id"`

	// // Biome in which the cell is
	// Biome uint64

	// Location of the Cell on the map. Used for rendering
	X uint64 `json:"x"`
	Y uint64 `json:"y"`

	// The unique ID of the city present at this location.
	City uint64 `json:"city,omitempty"`
}

// A Map is a directed graph destined to be used as a transport network,
// organised as an adjacency list.
type Map struct {
	Cells SetOfVertices `json:"cells"`
	Roads SetOfEdges    `json:"roads"`

	nextID uint64
	steps  map[vector]uint64
}

type SetOfEdges []*MapEdge

//go:generate go run github.com/jfsmig/hegemonie/cmd/gen-set -acc .ID mapgraph ./map_auto.go *MapVertex SetOfVertices
