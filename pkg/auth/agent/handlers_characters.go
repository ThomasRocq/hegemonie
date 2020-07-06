// Copyright (C) 2018-2020 Hegemonie's AUTHORS
// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at http://mozilla.org/MPL/2.0/.

package hegemonie_auth_agent

import (
	"context"
	"github.com/jfsmig/hegemonie/pkg/auth/backend"
	proto "github.com/jfsmig/hegemonie/pkg/auth/proto"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"io"
)

var none = &proto.None{}

func characterStateView(s uint) proto.CharacterState {
	switch s {
	case hegemonie_auth_backend.CharacterStateUnknown:
		return proto.CharacterState_CharacterUnknown
	case hegemonie_auth_backend.CharacterStateActive:
		return proto.CharacterState_CharacterActive
	case hegemonie_auth_backend.CharacterStateSuper:
		return proto.CharacterState_CharacterSuper
	case hegemonie_auth_backend.CharacterStateSuspended:
		return proto.CharacterState_CharacterSuspended
	case hegemonie_auth_backend.CharacterStateDeleted:
		return proto.CharacterState_CharacterDeleted
	default:
		// FIXME(jfs): log something, there is a corruption
		return proto.CharacterState_CharacterUnknown
	}
}

func characterView(u hegemonie_auth_backend.Character) *proto.CharacterView {
	return &proto.CharacterView{
		UserId: "NOT-SET",
		Name:   u.Name,
		Region: u.Region,
		State:  characterStateView(u.State),
	}
}

func (srv *charactersService) List(req *proto.CharacterListReq, stream proto.Character_ListServer) error {
	for {
		tab, err := srv.characters.List(req.Region, req.Marker, 100)
		if err != nil {
			return err
		}
		if req.User != "" {
			for _, c := range tab {
				if c.User != req.User {
					return status.Error(codes.DataLoss, "ghost character")
				}
			}
		}
		for _, c := range tab {
			err = stream.Send(characterView(c))
			if err == io.EOF {
				return nil
			}
			if err != nil {
				return err
			}
		}
	}
}

func (srv *charactersService) Create(ctx context.Context, req *proto.CharacterCreateReq) (*proto.CharacterView, error) {
	c, err := srv.characters.Create(req.User, req.Region, req.Name)
	if err != nil {
		return nil, err
	}
	return characterView(c), err
}

func (srv *charactersService) Delete(ctx context.Context, req *proto.CharacterId) (*proto.None, error) {
	if req.User != "" {
		u, err := srv.characters.Show(req.Region, req.Name)
		if err != nil {
			return nil, status.Error(codes.NotFound, "Not such User")
		}
		if u.User != req.User {
			return nil, status.Error(codes.Unauthenticated, "Not your character")
		}
	}

	err := srv.characters.Delete(req.Region, req.Name)
	if err != nil {
		return nil, err
	}
	return none, err
}

func (srv *charactersService) Show(ctx context.Context, req *proto.CharacterId) (*proto.CharacterView, error) {
	c, err := srv.characters.Show(req.Region, req.Name)
	if err != nil {
		return nil, err
	}
	if req.User != "" && req.User != c.User {
		return nil, status.Error(codes.FailedPrecondition, "User/Character mismatch")
	}
	return characterView(c), err
}

func (srv *charactersService) Promote(ctx context.Context, req *proto.CharacterId) (*proto.None, error) {
	if req.User != "" {
		u, err := srv.characters.Show(req.Region, req.Name)
		if err != nil {
			return nil, status.Error(codes.NotFound, "Not such User")
		}
		if u.User != req.User {
			return nil, status.Error(codes.Unauthenticated, "Not your character")
		}
	}

	err := srv.characters.Promote(req.Region, req.Name)
	if err != nil {
		return nil, err
	}
	return none, err
}

func (srv *charactersService) Demote(ctx context.Context, req *proto.CharacterId) (*proto.None, error) {
	if req.User != "" {
		u, err := srv.characters.Show(req.Region, req.Name)
		if err != nil {
			return nil, status.Error(codes.NotFound, "Not such User")
		}
		if u.User != req.User {
			return nil, status.Error(codes.Unauthenticated, "Not your character")
		}
	}

	err := srv.characters.Demote(req.Region, req.Name)
	if err != nil {
		return nil, err
	}
	return none, err
}

func (srv *charactersService) Suspend(ctx context.Context, req *proto.CharacterId) (*proto.None, error) {
	if req.User != "" {
		u, err := srv.characters.Show(req.Region, req.Name)
		if err != nil {
			return nil, status.Error(codes.NotFound, "Not such User")
		}
		if u.User != req.User {
			return nil, status.Error(codes.Unauthenticated, "Not your character")
		}
	}

	err := srv.characters.Suspend(req.Region, req.Name)
	return none, err
}

func (srv *charactersService) Resume(ctx context.Context, req *proto.CharacterId) (*proto.None, error) {
	if req.User != "" {
		u, err := srv.characters.Show(req.Region, req.Name)
		if err != nil {
			return nil, status.Error(codes.NotFound, "Not such User")
		}
		if u.User != req.User {
			return nil, status.Error(codes.Unauthenticated, "Not your character")
		}
	}

	err := srv.characters.Resume(req.Region, req.Name)
	return none, err
}
