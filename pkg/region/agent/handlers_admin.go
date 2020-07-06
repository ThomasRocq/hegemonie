// Copyright (C) 2018-2020 Hegemonie's AUTHORS
// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at http://mozilla.org/MPL/2.0/.

package hegemonie_region_agent

import (
	"context"
	"github.com/jfsmig/hegemonie/pkg/region/model"
	proto "github.com/jfsmig/hegemonie/pkg/region/proto"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"io"
)

type srvAdmin struct {
	cfg *regionConfig
	w   *region.World
}

var none = &proto.None{}

func (srv *srvAdmin) rlockDo(action func() error) error {
	srv.w.RLock()
	defer srv.w.RUnlock()
	return action()
}

func (srv *srvAdmin) wlockDo(action func() error) error {
	srv.w.WLock()
	defer srv.w.WUnlock()
	return action()
}

func (s *srvAdmin) Produce(ctx context.Context, req *proto.RegionId) (*proto.None, error) {
	return none, s.rlockDo(func() error {
		r := s.w.Regions.Get(req.Region)
		if r == nil {
			return status.Error(codes.NotFound, "No such region")
		}
		r.Produce()
		return nil
	})
}

func (s *srvAdmin) Move(ctx context.Context, req *proto.RegionId) (*proto.None, error) {
	return none, s.rlockDo(func() error {
		r := s.w.Regions.Get(req.Region)
		if r == nil {
			return status.Error(codes.NotFound, "No such region")
		}
		r.Move()
		return nil
	})
}

func (s *srvAdmin) CreateRegion(ctx context.Context, req *proto.RegionCreateReq) (*proto.None, error) {
	return none, s.wlockDo(func() error {
		_, err := s.w.CreateRegion(req.Name, req.MapName)
		return err
	})
}

func (s *srvAdmin) GetScores(req *proto.RegionId, stream proto.Admin_GetScoresServer) error {
	return s.rlockDo(func() error {
		r := s.w.Regions.Get(req.Region)
		if r == nil {
			return status.Error(codes.NotFound, "No such region")
		}
		for _, c := range r.Cities {
			err := stream.Send(ShowCityPublic(s.w, c, true))
			if err == io.EOF {
				return nil
			}
			if err != nil {
				return err
			}
		}
		return nil
	})
}
