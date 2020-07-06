// Copyright (C) 2020 Hegemonie's AUTHORS
// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at http://mozilla.org/MPL/2.0/.

package mapclient

import (
	"fmt"
	"math"
	"math/rand"
	"sort"
	"sync/atomic"
)

type SiteMem struct {
	Raw   SiteRaw
	Peers map[*SiteMem]bool
}

type RoadMem struct {
	Src, Dst *SiteMem
}

// Human-unfriendly representation of a Map
// - The Sites are indexed by a unique number
// - The roads are bidirectional.
// - Road can be duplicated
type MapMem struct {
	ID     string
	Sites  map[uint64]*SiteMem
	nextID uint64
}

func MakeMemMap() MapMem {
	return MapMem{
		Sites: make(map[uint64]*SiteMem),
	}
}

func makeSite(raw SiteRaw) *SiteMem {
	return &SiteMem{
		Raw:   raw,
		Peers: make(map[*SiteMem]bool),
	}
}

func (s *SiteMem) DotName() string {
	if s.Raw.City != "" {
		return s.Raw.City
	}
	return fmt.Sprintf("x%v", s.Raw.ID)
}

func (m *MapMem) UniqueRoads() <-chan RoadMem {
	out := make(chan RoadMem)
	go func() {
		seen := make(map[RoadRaw]bool)
		for _, s := range m.Sites {
			for peer := range s.Peers {
				r0 := RoadRaw{Src: s.Raw.ID, Dst: peer.Raw.ID}
				if !seen[r0] {
					seen[r0] = true
					out <- RoadMem{s, peer}
				}
			}
		}
		close(out)
	}()
	return out
}

func (m *MapMem) SortedSites() <-chan *SiteMem {
	out := make(chan *SiteMem)
	go func() {
		keys := make([]uint64, 0, len(m.Sites))
		for k := range m.Sites {
			keys = append(keys, k)
		}
		sort.Slice(keys, func(i, j int) bool { return keys[i] < keys[j] })
		for _, k := range keys {
			out <- m.Sites[k]
		}
		close(out)
	}()
	return out
}

// Produces a MapRaw with sorted sites and unique roads.
func (m *MapMem) Raw() MapRaw {
	rawMap := MakeRawMap()
	rawMap.ID = m.ID
	for s := range m.SortedSites() {
		rawMap.Sites = append(rawMap.Sites, s.Raw)
	}
	for r := range m.UniqueRoads() {
		rawRoad := RoadRaw{Src: r.Src.Raw.ID, Dst: r.Dst.Raw.ID}
		rawMap.Roads = append(rawMap.Roads, rawRoad)
	}
	return rawMap
}

func (m *MapMem) DeepCopy() MapMem {
	mFinal := MakeMemMap()
	for id, site := range m.Sites {
		mFinal.Sites[id] = makeSite(site.Raw)
	}
	for _, s := range m.Sites {
		src := mFinal.Sites[s.Raw.ID]
		for d := range s.Peers {
			dst := mFinal.Sites[d.Raw.ID]
			src.Peers[dst] = true
			dst.Peers[src] = true
		}
	}
	return mFinal
}

func (m *MapMem) ComputeBox() (xmin, xmax, ymin, ymax uint64) {
	const Max = math.MaxUint64
	const Min = 0
	xmin, ymin, xmax, ymax = Max, Max, Min, Min
	for _, s := range m.Sites {
		x, y := s.Raw.X, s.Raw.Y
		if x < xmin {
			xmin = x
		}
		if x > xmax {
			xmax = x
		}
		if y < ymin {
			ymin = y
		}
		if y > ymax {
			ymax = y
		}
	}
	if xmin == Max {
		xmin, xmax, ymin, ymax = 0, 0, 0, 0
	}
	return
}

func (m *MapMem) ShiftAt(xabs, yabs uint64) {
	xmin, _, ymin, _ := m.ComputeBox()
	m.Shift(xabs-xmin, yabs-ymin)
}

func (m *MapMem) Shift(xrel, yrel uint64) {
	for _, s := range m.Sites {
		s.Raw.X += xrel
		s.Raw.Y += yrel
	}
}

func (m *MapMem) ResizeRatio(xratio, yratio float64) {
	for _, s := range m.Sites {
		s.Raw.X = uint64(math.Round(float64(s.Raw.X) * xratio))
		s.Raw.Y = uint64(math.Round(float64(s.Raw.Y) * yratio))
	}
}

func (m *MapMem) ResizeStretch(x, y float64) {
	m.ShiftAt(0, 0)
	_, xmax, _, ymax := m.ComputeBox()
	m.ResizeRatio(x/float64(xmax), y/float64(ymax))
}

func (m *MapMem) ResizeAdjust(x, y uint64) {
	m.ShiftAt(0, 0)
	_, xmax, _, ymax := m.ComputeBox()
	xRatio := float64(x) / float64(xmax)
	yRatio := float64(y) / float64(ymax)
	ratio := math.Min(xRatio, yRatio)
	m.ResizeRatio(ratio, ratio)
}

func (m *MapMem) Center(xbound, ybound uint64) {
	xmin, xmax, ymin, ymax := m.ComputeBox()
	xdelta, ydelta := xbound-(xmax-xmin), ybound-(ymax-ymin)
	xpad, ypad := xdelta/2.0, ydelta/2.0
	m.Shift(xpad-xmin, ypad-ymin)
}

func (m *MapMem) splitOneRoad(src, dst *SiteMem, nbSegments uint) {
	if nbSegments < 2 {
		panic("bug")
	}

	xinc := uint64(math.Round(float64(dst.Raw.X-src.Raw.X) / float64(nbSegments)))
	yinc := uint64(math.Round(float64(dst.Raw.Y-src.Raw.Y) / float64(nbSegments)))
	segments := make([]*SiteMem, 0, nbSegments+1)

	delete(src.Peers, dst)
	delete(dst.Peers, src)

	// Create segment boundaries
	segments = append(segments, src)
	for i := uint(0); i < nbSegments-1; i++ {
		last := segments[len(segments)-1]
		x := last.Raw.X + xinc
		y := last.Raw.Y + yinc
		id := atomic.AddUint64(&m.nextID, 1)
		raw := SiteRaw{ID: id, City: "", X: x, Y: y}
		middle := makeSite(raw)
		m.Sites[middle.Raw.ID] = middle
		segments = append(segments, middle)
	}
	segments = append(segments, dst)

	// Link the segment boundaries
	for i, end := range segments[1:] {
		start := segments[i]
		start.Peers[end] = true
		end.Peers[start] = true
	}
}

func (m *MapMem) SplitLongRoads(max float64) MapMem {
	// Work on a deep copy to iterate on the original map while we alter the copy
	mCopy := m.DeepCopy()
	for r := range m.UniqueRoads() {
		src := mCopy.Sites[r.Src.Raw.ID]
		dst := mCopy.Sites[r.Dst.Raw.ID]
		dist := distance(src, dst)
		if max < dist {
			mCopy.splitOneRoad(src, dst, uint(math.Ceil(dist/max)))
		}
	}
	return mCopy
}

func (m *MapMem) Noise(xjitter, yjitter float64) {
	for _, s := range m.Sites {
		if s.Raw.City != "" {
			continue
		}
		s.Raw.X += uint64(math.Round((0.5 - rand.Float64()) * xjitter))
		s.Raw.Y += uint64(math.Round((0.5 - rand.Float64()) * yjitter))
	}
}

func distance(src, dst *SiteMem) float64 {
	dx := (dst.Raw.X - src.Raw.X)
	dy := (dst.Raw.Y - src.Raw.Y)
	return math.Sqrt(float64(dx*dx) + float64(dy*dy))
}
