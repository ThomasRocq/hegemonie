// Copyright (C) 2018-2020 Hegemonie's AUTHORS
// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at http://mozilla.org/MPL/2.0/.

package mapagent

import (
	"context"
	"github.com/jfsmig/hegemonie/pkg/map/graph"
	proto "github.com/jfsmig/hegemonie/pkg/map/proto"
)

type srvMap struct {
	config *regionConfig
	repo   mapgraph.Repository
}

func (s *srvMap) Vertices(ctx context.Context, req *proto.ListVerticesReq) (*proto.ListOfVertices, error) {
	s.repo.RLock()
	defer s.repo.RUnlock()

	m, err := s.repo.GetMap(req.MapName)
	if err != nil {
		return nil, err
	}

	rep := &proto.ListOfVertices{}
	for _, x := range m.Cells.Slice(req.Marker, req.Max) {
		rep.Items = append(rep.Items, &proto.Vertex{Id: x.ID, X: x.X, Y: x.Y})
	}
	return rep, nil
}

func (s *srvMap) Edges(ctx context.Context, req *proto.ListEdgesReq) (*proto.ListOfEdges, error) {
	s.repo.RLock()
	defer s.repo.RUnlock()

	m, err := s.repo.GetMap(req.MapName)
	if err != nil {
		return nil, err
	}

	rep := &proto.ListOfEdges{}
	for _, x := range m.Roads.Slice(req.MarkerSrc, req.MarkerDst, req.Max) {
		rep.Items = append(rep.Items, &proto.Edge{Src: x.S, Dst: x.D})
	}
	return rep, nil
}

func (s *srvMap) GetPath(ctx context.Context, req *proto.PathRequest) (*proto.PathReply, error) {
	s.repo.RLock()
	defer s.repo.RUnlock()

	m, err := s.repo.GetMap(req.MapName)
	if err != nil {
		return nil, err
	}

	rep := &proto.PathReply{}
	p := req.Src
	for {
		p, err = m.PathNextStep(p, req.Dst)
		if err != nil {
			return nil, err
		}
		rep.Steps = append(rep.Steps, p)
		if p == req.Dst || (req.Max > 0 && uint32(len(rep.Steps)) > req.Max) {
			break
		}
	}
	return rep, nil
}

func (s *srvMap) Cities(ctx context.Context, req *proto.ListCitiesReq) (*proto.ListOfCities, error) {
	s.repo.RLock()
	defer s.repo.RUnlock()

	m, err := s.repo.GetMap(req.MapName)
	if err != nil {
		return nil, err
	}

	if req.Marker <= 0 || req.Marker > 1024 {
		req.Marker = 1024
	}

	rep := &proto.ListOfCities{}
	next := req.Marker
	oneMore := func() bool { return uint64(len(rep.Items)) < req.Marker }

	for oneMore() {
		for _, v := range m.Cells.Slice(next, req.Max) {
			if v.CityHere {
				rep.Items = append(rep.Items, &proto.CityLocation{
					Id: v.ID, Name: v.City,
				})
				if !oneMore() {
					return rep, nil
				}
			}
			next = v.ID
		}
	}
	return rep, nil
}
