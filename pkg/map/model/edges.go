// Copyright (C) 2018-2020 Hegemonie's AUTHORS
// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at http://mozilla.org/MPL/2.0/.

package mapgraph

import (
	"sort"
)

func (ev SetOfEdges) Len() int      { return len(ev) }
func (ev SetOfEdges) Swap(i, j int) { ev[i], ev[j] = ev[j], ev[i] }
func (ev SetOfEdges) Less(i, j int) bool {
	return ev.edgeLess(i, *ev[j])
}

func (ev SetOfEdges) edgeLess(i int, d MapEdge) bool {
	s := ev[i]
	return s.S < d.S || (s.S == d.S && s.D < d.D)
}

func (ev SetOfEdges) First(at uint64) int {
	return sort.Search(len(ev), func(i int) bool { return ev[i].S >= at })
}

func (ev *SetOfEdges) Add(e *MapEdge) {
	*ev = append(*ev, e)
	if nb := len(*ev); nb > 2 && !sort.IsSorted((*ev)[nb-2:]) {
		sort.Sort(*ev)
	}
}

func (ev SetOfEdges) Get(src, dst uint64) *MapEdge {
	i := sort.Search(len(ev), func(i int) bool {
		return ev[i].S >= src || (ev[i].S == src && ev[i].D >= dst)
	})
	if i < len(ev) && ev[i].S == src && ev[i].D == dst {
		return ev[i]
	}
	return nil

}

func (ev SetOfEdges) Slice(markerSrc, markerDst uint64, max uint32) []MapEdge {
	tab := make([]MapEdge, 0)

	iMax := ev.Len()
	i := ev.First(markerSrc)
	if i < iMax && ev[i].S == markerSrc && ev[i].D == markerDst {
		i++
	}

	needle := MapEdge{S: markerSrc, D: markerDst}
	for ; i < iMax; i++ {
		if ev.edgeLess(i, needle) {
			continue
		}
		tab = append(tab, *ev[i])
		if uint32(len(tab)) >= max {
			break
		}
	}
	return tab
}
