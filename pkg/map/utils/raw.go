// Copyright (C) 2020 Hegemonie's AUTHORS
// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at http://mozilla.org/MPL/2.0/.

package maputils

import (
	"fmt"
)

type SiteRaw struct {
	ID   string  `json:"id"`
	X    float64 `json:"x"`
	Y    float64 `json:"y"`
	City bool
}

type RoadRaw struct {
	Src string `json:"src"`
	Dst string `json:"dst"`
}

type MapRaw struct {
	ID    string    `json:"id"`
	Sites []SiteRaw `json:"sites"`
	Roads []RoadRaw `json:"roads"`
}

func MakeRawMap() MapRaw {
	return MapRaw{
		Sites: make([]SiteRaw, 0),
		Roads: make([]RoadRaw, 0),
	}
}

func (mr *MapRaw) Finalize() (Map, error) {
	var err error
	m := makeMap()

	for _, s := range mr.Sites {
		m.Sites[s.ID] = &Site{
			Raw:   s,
			Peers: make(map[*Site]bool),
		}
	}
	for _, r := range mr.Roads {
		if src, ok := m.Sites[r.Src]; !ok {
			err = fmt.Errorf("No such site [%s]", r.Src)
			break
		} else if dst, ok := m.Sites[r.Dst]; !ok {
			err = fmt.Errorf("No such site [%s]", r.Dst)
			break
		} else {
			src.Peers[dst] = true
			dst.Peers[src] = true
		}
	}
	return m, err
}
