// Copyright (C) 2018-2020 Hegemonie's AUTHORS
// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at http://mozilla.org/MPL/2.0/.

package mapagent

import (
	"errors"
	"github.com/jfsmig/hegemonie/pkg/map/graph"
	proto "github.com/jfsmig/hegemonie/pkg/map/proto"
	"os"
	"path/filepath"
	"strings"
	"sync"
)

type srvMap struct {
	config *mapServiceConfig

	maps mapgraph.SetOfMaps
	rw   sync.RWMutex
}

func (s *srvMap) Vertices(req *proto.ListVerticesReq, stream proto.Map_VerticesServer) error {
	s.rw.RLock()
	defer s.rw.RUnlock()

	m := s.maps.Get(req.MapName)
	if m == nil {
		return errors.New("No Such Map")
	}

	next := req.Marker
	for {
		vertices := m.Cells.Slice(next, 100)
		if len(vertices) <= 0 {
			return nil
		}
		for _, x := range vertices {
			err := stream.Send(&proto.Vertex{Id: x.ID, X: x.X, Y: x.Y})
			if err != nil {
				return err
			}
			next = x.ID
		}
	}
}

func (s *srvMap) Edges(req *proto.ListEdgesReq, stream proto.Map_EdgesServer) error {
	s.rw.RLock()
	defer s.rw.RUnlock()

	m := s.maps.Get(req.MapName)
	if m == nil {
		return errors.New("No Such Map")
	}

	src, dst := req.MarkerSrc, req.MarkerDst
	for {
		edges := m.Roads.Slice(src, dst, 100)
		if len(edges) <= 0 {
			return nil
		}
		for _, x := range edges {
			err := stream.Send(&proto.Edge{Src: x.S, Dst: x.D})
			if err != nil {
				return err
			}
			src, dst = x.S, x.D
		}
	}
}

func (s *srvMap) GetPath(req *proto.PathRequest, stream proto.Map_GetPathServer) error {
	s.rw.RLock()
	defer s.rw.RUnlock()

	m := s.maps.Get(req.MapName)
	if m == nil {
		return errors.New("No Such Map")
	}

	src := req.Src
	for {
		next, err := m.PathNextStep(src, req.Dst)
		if err != nil {
			return err
		}
		err = stream.Send(&proto.PathElement{Id: src})
		if err != nil {
			return err
		}
		if next == req.Dst {
			return nil
		}
		src = next
	}
}

func (s *srvMap) Cities(req *proto.ListCitiesReq, stream proto.Map_CitiesServer) error {
	s.rw.RLock()
	defer s.rw.RUnlock()

	m := s.maps.Get(req.MapName)
	if m == nil {
		return errors.New("No Such Map")
	}

	next := req.Marker
	for {
		cities := m.Cells.Slice(next, 100)
		if len(cities) <= 0 {
			return nil
		}
		for _, v := range cities {
			if v.City != "" {
				err := stream.Send(&proto.CityLocation{
					Id: v.ID, Name: v.City,
				})
				if err != nil {
					return err
				}
			}
			next = v.ID
		}
	}
}

func (s *srvMap) Maps(req *proto.ListMapsReq, stream proto.Map_MapsServer) error {
	slice := func(marker string) []string {
		s.rw.RLock()
		defer s.rw.RUnlock()
		out := make([]string, 0)
		for _, m := range s.maps.Slice(marker, 100) {
			out = append(out, m.ID)
		}
		return out
	}

	next := req.Marker
	for {
		names := slice(next)
		if len(names) <= 0 {
			return nil
		}
		for _, v := range names {
			if err := stream.Send(&proto.MapName{Name: v}); err != nil {
				return err
			}
			if v > next {
				next = v
			}
		}
	}
}

func (s *srvMap) LoadDirectory(path string) error {
	return filepath.Walk(path, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		// Only accept non-hidden JSON files
		_, fn := filepath.Split(path)
		if info.IsDir() || info.Size() <= 0 {
			return nil
		}
		if len(fn) < 2 || fn[0] == '.' {
			return nil
		}
		if !strings.HasSuffix(fn, ".final.json") {
			return nil
		}

		m := mapgraph.NewMap()
		if f, err := os.Open(path); err != nil {
			return err
		} else {
			defer f.Close()
			if err = m.Load(f); err != nil {
				return err
			}
		}

		s.maps.Add(m)
		return nil
	})
}
