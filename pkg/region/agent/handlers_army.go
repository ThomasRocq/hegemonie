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
)

type srvArmy struct {
	cfg *regionConfig
	w   *region.World
}

func (s *srvArmy) getAndCheckArmy(req *proto.ArmyId) (*region.City, *region.Army, error) {
	city, err := s.w.CityGetAndCheck(req.Character, req.City)
	if err != nil {
		return nil, nil, err
	}
	army := city.Armies.Get(req.Army)
	if army == nil {
		return nil, nil, status.Errorf(codes.NotFound, "Army Not found (id=%d)", req.Army)
	}
	return city, army, err
}

func (s *srvArmy) wlock(action func() error) error {
	s.w.WLock()
	defer s.w.WUnlock()
	return action()
}

func (s *srvArmy) rlock(action func() error) error {
	s.w.RLock()
	defer s.w.RUnlock()
	return action()
}

func (s *srvArmy) getAndDo(id *proto.ArmyId, action func(*region.City, *region.Army) error) error {
	city, army, err := s.getAndCheckArmy(id)
	if err == nil {
		err = action(city, army)
	}
	return err
}

func (s *srvArmy) wlockDo(id *proto.ArmyId, action func(*region.City, *region.Army) error) error {
	return s.wlock(func() error { return s.getAndDo(id, action) })
}

func (s *srvArmy) rlockDo(id *proto.ArmyId, action func(*region.City, *region.Army) error) error {
	return s.rlock(func() error { return s.getAndDo(id, action) })
}

func (s *srvArmy) Show(ctx context.Context, req *proto.ArmyId) (*proto.ArmyView, error) {
	var rc *proto.ArmyView
	err := s.rlockDo(req, func(_ *region.City, army *region.Army) error {
		rc = ShowArmy(s.w, army)
		return nil
	})
	return rc, err
}

func (s *srvArmy) Flea(ctx context.Context, req *proto.ArmyId) (*proto.None, error) {
	return &proto.None{}, s.wlockDo(req, func(_ *region.City, a *region.Army) error {
		return a.Flea(s.w)
	})
}

func (s *srvArmy) Flip(ctx context.Context, req *proto.ArmyId) (*proto.None, error) {
	return &proto.None{}, s.wlockDo(req, func(_ *region.City, a *region.Army) error {
		return a.Flip(s.w)
	})
}

func (s *srvArmy) Move(ctx context.Context, req *proto.ArmyMoveReq) (*proto.None, error) {
	return &proto.None{}, s.wlockDo(req.Id, func(_ *region.City, army *region.Army) error {
		return army.DeferMove(s.w, req.Target, region.ActionArgMove{})
	})
}

func (s *srvArmy) Attack(ctx context.Context, req *proto.ArmyAssaultReq) (*proto.None, error) {
	return &proto.None{}, s.wlockDo(req.Id, func(_ *region.City, army *region.Army) error {
		return army.DeferAttack(s.w, req.Target, region.ActionArgAssault{})
	})
}

func (s *srvArmy) Wait(ctx context.Context, req *proto.ArmyTarget) (*proto.None, error) {
	return &proto.None{}, s.wlockDo(req.Id, func(_ *region.City, army *region.Army) error {
		return army.DeferWait(s.w, req.Target)
	})
}

func (s *srvArmy) Defend(ctx context.Context, req *proto.ArmyTarget) (*proto.None, error) {
	return &proto.None{}, s.wlockDo(req.Id, func(_ *region.City, army *region.Army) error {
		return army.DeferDefend(s.w, req.Target)
	})
}

func (s *srvArmy) Disband(ctx context.Context, req *proto.ArmyTarget) (*proto.None, error) {
	return &proto.None{}, s.wlockDo(req.Id, func(_ *region.City, army *region.Army) error {
		return army.DeferDisband(s.w, req.Target)
	})
}

func (s *srvArmy) Cancel(ctx context.Context, req *proto.ArmyId) (*proto.None, error) {
	return &proto.None{}, s.wlockDo(req, func(_ *region.City, army *region.Army) error {
		return army.Cancel(s.w)
	})
}
