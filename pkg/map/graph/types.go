// Copyright (C) 2018-2020 Hegemonie's AUTHORS
// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at http://mozilla.org/MPL/2.0/.

package mapgraph

// A Edge is an edge if the transportation directed graph
type Edge struct {
	// Unique identifier of the source Cell
	S uint64 `json:"src"`

	// Unique identifier of the destination Cell
	D uint64 `json:"dst"`
}

// A Vertex is a vertex in the transportation directed graph
type Vertex struct {
	// The unique identifier of the current cell.
	ID uint64 `json:"id"`

	// // Biome in which the cell is
	// Biome uint64

	// Location of the Cell on the map. Used for rendering
	X uint64 `json:"x"`
	Y uint64 `json:"y"`

	// Should the current location carry a city when the region starts,
	// and if yes, what should be the name of that city.
	CityHere bool   `json:"city,omitempty"`
	City     string `json:"cityName,omitempty"`
}

// A Map is a directed graph destined to be used as a transport network,
// organised as an adjacency list.
type Map struct {
	// The unique name of the map
	ID string `json:"id"`

	Cells SetOfVertices `json:"cells"`

	Roads SetOfEdges `json:"roads"`

	nextID uint64
	steps  map[vector]uint64
}

type Repository interface {
	RLock()
	RUnlock()

	WLock()
	WUnlock()

	GetMap(ID string) (*Map, error)
	ListMaps(marker string, max uint32) ([]*Map, error)
	Register(m *Map)
}

//go:generate go run github.com/jfsmig/hegemonie/cmd/gen-set ./map_auto.go mapgraph:SetOfVertices:*Vertex ID:uint64
//go:generate go run github.com/jfsmig/hegemonie/cmd/gen-set ./map_auto.go mapgraph:SetOfEdges:*Edge S:uint64 D:uint64
