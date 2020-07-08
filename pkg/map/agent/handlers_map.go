// Copyright (C) 2018-2020 Hegemonie's AUTHORS
// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at http://mozilla.org/MPL/2.0/.

package hegemonie_map_agent

import (
	"context"
	"github.com/jfsmig/hegemonie/pkg/map/model"
	proto "github.com/jfsmig/hegemonie/pkg/map/proto"
)

type srvMap struct {
	cfg *regionConfig
	w   *mapgraph.Map
}

func (s *srvMap) Vertices(ctx context.Context, req *proto.ListVerticesReq) (*proto.ListOfVertices, error) {
	// FIXME(jfs): Need for a lock
	//s.w.RLock()
	//defer s.w.RUnlock()

	rep := &proto.ListOfVertices{}
	for _, x := range s.w.Cells.Slice(req.Marker, req.Max) {
		rep.Items = append(rep.Items, &proto.Vertex{
			Id: x.ID, X: x.X, Y: x.Y, CityId: x.City})
	}
	return rep, nil
}

func (s *srvMap) Edges(ctx context.Context, req *proto.ListEdgesReq) (*proto.ListOfEdges, error) {
	// FIXME(jfs): Need for a lock
	//s.w.RLock()
	//defer s.w.RUnlock()

	rep := &proto.ListOfEdges{}
	for _, x := range s.w.Roads.Slice(req.MarkerSrc, req.MarkerDst, req.Max) {
		rep.Items = append(rep.Items, &proto.Edge{Src: x.S, Dst: x.D})
	}
	return rep, nil
}

func (s *srvMap) GetPath(ctx context.Context, req *proto.PathRequest) (*proto.PathReply, error) {
	// FIXME(jfs): Need for a lock
	//s.w.RLock()
	//defer s.w.RUnlock()

	var err error
	rep := &proto.PathReply{}
	p := req.Src
	for {
		p, err = s.w.PathNextStep(p, req.Dst)
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
