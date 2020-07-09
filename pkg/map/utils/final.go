// Copyright (C) 2020 Hegemonie's AUTHORS
// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at http://mozilla.org/MPL/2.0/.

package maputils

import (
	"math"
	"math/rand"
	"sort"
	"strconv"
)

type Road struct {
	Src, Dst *Site
}

type Map struct {
	Id    string
	Sites map[string]*Site
}

type Site struct {
	Raw   SiteRaw
	Peers map[*Site]bool
}

func makeMap() Map {
	return Map{
		Sites: make(map[string]*Site),
	}
}

func makeSite(raw SiteRaw) *Site {
	return &Site{
		Raw:   raw,
		Peers: make(map[*Site]bool),
	}
}

func (s *Site) DotName() string {
	if s.Raw.City {
		return s.Raw.ID
	}
	return "r" + s.Raw.ID
}

func (r *Road) Raw() RoadRaw {
	return RoadRaw{Src: r.Src.Raw.ID, Dst: r.Dst.Raw.ID}
}

func (m *Map) UniqueRoads() <-chan Road {
	out := make(chan Road)
	go func() {
		seen := make(map[RoadRaw]bool)
		for _, s := range m.Sites {
			for peer := range s.Peers {
				r0 := RoadRaw{Src: s.Raw.ID, Dst: peer.Raw.ID}
				r1 := RoadRaw{Src: peer.Raw.ID, Dst: s.Raw.ID}
				if !seen[r0] && !seen[r1] {
					seen[r0] = true
					seen[r1] = true
					out <- Road{s, peer}
				}
			}
		}
		close(out)
	}()
	return out
}

func (m *Map) SortedSites() <-chan *Site {
	out := make(chan *Site)
	go func() {
		keys := make([]string, 0, len(m.Sites))
		for k := range m.Sites {
			keys = append(keys, k)
		}
		sort.Strings(keys)
		for _, k := range keys {
			out <- m.Sites[k]
		}
		close(out)
	}()
	return out
}

func (m *Map) Raw() MapRaw {
	rm := MakeRawMap()
	for s := range m.SortedSites() {
		rm.Sites = append(rm.Sites, s.Raw)
	}
	for r := range m.UniqueRoads() {
		rm.Roads = append(rm.Roads, r.Raw())
	}
	return rm
}

func (m *Map) DeepCopy() Map {
	mFinal := makeMap()
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

func (m *Map) ComputeBox() (xmin, xmax, ymin, ymax float64) {
	const Max = math.MaxFloat64
	const Min = -Max
	xmin, ymin = Max, Max
	xmax, ymax = Min, Min
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

func (m *Map) ShiftAt(xabs, yabs float64) {
	xmin, _, ymin, _ := m.ComputeBox()
	m.Shift(xabs-xmin, yabs-ymin)
}

func (m *Map) Shift(xrel, yrel float64) {
	for _, s := range m.Sites {
		s.Raw.X += xrel
		s.Raw.Y += yrel
	}
}

func (m *Map) ResizeRatio(xratio, yratio float64) {
	for _, s := range m.Sites {
		s.Raw.X *= xratio
		s.Raw.Y *= yratio
	}
}

func (m *Map) ResizeStretch(x, y float64) {
	m.ShiftAt(0, 0)
	_, xmax, _, ymax := m.ComputeBox()
	m.ResizeRatio(x/xmax, y/ymax)
}

func (m *Map) ResizeAdjust(x, y float64) {
	m.ShiftAt(0, 0)
	_, xmax, _, ymax := m.ComputeBox()
	xRatio := x / xmax
	yRatio := y / ymax
	ratio := math.Min(xRatio, yRatio)
	m.ResizeRatio(ratio, ratio)
}

func (m *Map) Center(xbound, ybound float64) {
	xmin, xmax, ymin, ymax := m.ComputeBox()
	xdelta, ydelta := xbound-(xmax-xmin), ybound-(ymax-ymin)
	xpad, ypad := xdelta/2.0, ydelta/2.0
	m.Shift(xpad-xmin, ypad-ymin)
}

func (m *Map) splitOneRoad(src, dst *Site, nbSegments uint) {
	if nbSegments < 2 {
		panic("bug")
	}

	xinc := (dst.Raw.X - src.Raw.X) / float64(nbSegments)
	yinc := (dst.Raw.Y - src.Raw.Y) / float64(nbSegments)
	segments := make([]*Site, 0, nbSegments+1)

	delete(src.Peers, dst)
	delete(dst.Peers, src)

	// Create segment boundaries
	segments = append(segments, src)
	for i := uint(0); i < nbSegments-1; i++ {
		last := segments[len(segments)-1]
		x := math.Round(last.Raw.X + xinc)
		y := math.Round(last.Raw.Y + yinc)
		id := "x-" + strconv.FormatInt(int64(x), 10) + "-" + strconv.FormatInt(int64(y), 10)
		raw := SiteRaw{
			ID:   id,
			City: false,
			X:    x,
			Y:    y,
		}
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

func (m *Map) SplitLongRoads(max float64) Map {
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

func (m *Map) Noise(xjitter, yjitter float64) {
	for _, s := range m.Sites {
		if s.Raw.City {
			continue
		}
		s.Raw.X += (0.5 - rand.Float64()) * xjitter
		s.Raw.Y += (0.5 - rand.Float64()) * yjitter
	}
}

func distance(src, dst *Site) float64 {
	dx := (dst.Raw.X - src.Raw.X)
	dy := (dst.Raw.Y - src.Raw.Y)
	return math.Sqrt(dx*dx + dy*dy)
}
