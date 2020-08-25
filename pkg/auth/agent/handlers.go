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
	default:
		// FIXME(jfs): log something, there is a corruption
		return proto.CharacterState_CharacterUnknown
	}
}

func userStateView(s uint) proto.UserState {
	switch s {
	case hegemonie_auth_backend.UserStateUnknown:
		return proto.UserState_UserUnknown
	case hegemonie_auth_backend.UserStateInvited:
		return proto.UserState_UserInvited
	case hegemonie_auth_backend.UserStateActive:
		return proto.UserState_UserActive
	case hegemonie_auth_backend.UserStateAdmin:
		return proto.UserState_UserActive
	case hegemonie_auth_backend.UserStateSuspended:
		return proto.UserState_UserSuspended
	default:
		// FIXME(jfs): log something, there is a corruption
		return proto.UserState_UserUnknown
	}
}

func userView(u hegemonie_auth_backend.User) *proto.UserView {
	return &proto.UserView{
		UserId: u.ID,
		Name:   u.Name,
		Mail:   u.Email,
		State:  userStateView(u.State),
	}
}

func characterView(u hegemonie_auth_backend.Character) *proto.CharacterView {
	return &proto.CharacterView{
		UserId: "NOT-SET",
		CharId: u.ID,
		Name:   u.Name,
		Region: u.Region,
		State:  characterStateView(u.State),
	}
}

func (srv *authService) UserList(req proto.UserListReq, stream proto.Auth_UserListServer) error {
	if req.Limit <= 0 {
		req.Limit = 1024
	}

	tab, err := srv.db.UsersList(req.Marker, req.Limit)
	if err != nil {
		return err
	}

	for _, u := range tab {
		err = stream.Send(userView(u))
		if err != nil {
			return err
		}
	}
	return nil
}

func (srv *authService) UserCreate(ctx context.Context, req *proto.UserCreateReq) (*proto.UserView, error) {
	var err error
	u, err := srv.db.UserCreate(req.Mail, req.Pass)
	if err != nil {
		return nil, err
	}
	return userView(u), err
}

func (srv *authService) UserShow(ctx context.Context, req *proto.UserId) (*proto.UserView, error) {
	u, err := srv.db.UserShow(req.Id)
	if err != nil {
		return nil, status.Error(codes.NotFound, "No such User")
	}
	return userView(u), nil
}

func (srv *authService) UserPromote(ctx context.Context, req *proto.UserId) (*proto.None, error) {
	err := srv.db.UserPromote(req.Id)
	if err != nil {
		return nil, status.Error(codes.NotFound, "No such User")
	}
	return none, nil
}

func (srv *authService) UserDemote(ctx context.Context, req *proto.UserId) (*proto.None, error) {
	err := srv.db.UserDemote(req.Id)
	if err != nil {
		return nil, status.Error(codes.NotFound, "No such User")
	}
	return none, nil
}

func (srv *authService) UserSuspend(ctx context.Context, req *proto.UserId) (*proto.None, error) {
	err := srv.db.UserSuspend(req.Id)
	if err != nil {
		return nil, status.Error(codes.NotFound, "No such User")
	}
	return none, nil
}

func (srv *authService) UserResume(ctx context.Context, req *proto.UserId) (*proto.None, error) {
	err := srv.db.UserResume(req.Id)
	if err != nil {
		return nil, status.Error(codes.NotFound, "No such User")
	}
	return none, nil
}

func (srv *authService) UserAuth(ctx context.Context, req *proto.UserAuthReq) (*proto.UserView, error) {
	u, err := srv.db.UserAuthenticate(req.Mail, req.Pass)
	if err != nil {
		return nil, status.Error(codes.NotFound, "No such User")
	}
	return userView(u), nil
}

func (srv *authService) CharacterList(ctx context.Context, req *proto.CharacterListReq) (*proto.UserView, error) {
	if req.User <= 0 || req.Character <= 0 {
		return nil, status.Error(codes.InvalidArgument, "Invalid user/role ID")
	}

	if u := srv.db.UserGet(req.User); u == nil {
		return nil, status.Error(codes.NotFound, "Not such User")
	} else {
		if u.Suspended || u.Deleted {
			return nil, status.Error(codes.PermissionDenied, "User suspended")
		}
		for _, c := range u.Characters {
			if c.ID == req.Character {
				uView := userView(u)
				uView.Characters = append(uView.Characters, &proto.CharacterView{
					Id: c.ID, Region: c.Region, Name: c.Name, Off: c.Off,
				})
				return uView, nil
			}
		}
		return nil, status.Error(codes.PermissionDenied, "Character mismatch")
	}
}

func (srv *authService) CharacterCreate(ctx context.Context, req *proto.CharacterCreateReq) (*proto.None, error) {
	if req.User <= 0 || req.Character <= 0 {
		return nil, status.Error(codes.InvalidArgument, "Invalid user ID")
	}

	u := srv.db.UserGet(req.User)
	if u == nil {
		return nil, status.Error(codes.NotFound, "Not such User")
	}

	c := u.GetCharacter(req.Character)
	if c == nil {
		return nil, status.Error(codes.NotFound, "No such Role")
	}
	c.Name = req.Name
	c.Off = req.Off
	c.Deleted = req.Deleted

	return &proto.None{}, nil
}

func (srv *authService) CharacterDelete(ctx context.Context, req *proto.CharacterId) (*proto.None, error) {
	if req.User <= 0 || req.Character <= 0 {
		return nil, status.Error(codes.InvalidArgument, "Invalid user ID")
	}

	u := srv.db.UserGet(req.User)
	if u == nil {
		return nil, status.Error(codes.NotFound, "Not such User")
	}

	c := u.GetCharacter(req.Character)
	if c == nil {
		return nil, status.Error(codes.NotFound, "No such Role")
	}
	c.Region = req.Region

	return &proto.None{}, nil
}

func (srv *authService) CharacterShow(ctx context.Context, req *proto.CharacterId) (*proto.UserView, error) {
	if req.User <= 0 || req.Character <= 0 {
		return nil, status.Error(codes.InvalidArgument, "Invalid user/role ID")
	}

	if u := srv.db.UserGet(req.User); u == nil {
		return nil, status.Error(codes.NotFound, "Not such User")
	} else {
		if u.Suspended || u.Deleted {
			return nil, status.Error(codes.PermissionDenied, "User suspended")
		}
		for _, c := range u.Characters {
			if c.ID == req.Character {
				uView := userView(u)
				uView.Characters = append(uView.Characters, &proto.CharacterView{
					Id: c.ID, Region: c.Region, Name: c.Name, Off: c.Off,
				})
				return uView, nil
			}
		}
		return nil, status.Error(codes.PermissionDenied, "Character mismatch")
	}
}

func (srv *authService) CharacterSuper(ctx context.Context, req *proto.CharacterId) (*proto.None, error) {
	if req.User <= 0 || req.Character <= 0 {
		return nil, status.Error(codes.InvalidArgument, "Invalid user ID")
	}

	u := srv.db.UserGet(req.User)
	if u == nil {
		return nil, status.Error(codes.NotFound, "Not such User")
	}

	c := u.GetCharacter(req.Character)
	if c == nil {
		return nil, status.Error(codes.NotFound, "No such Role")
	}
	c.Region = req.Region

	return &proto.None{}, nil
}

func (srv *authService) CharacterNormal(ctx context.Context, req *proto.CharacterId) (*proto.None, error) {
	if req.User <= 0 || req.Character <= 0 {
		return nil, status.Error(codes.InvalidArgument, "Invalid user ID")
	}

	u := srv.db.UserGet(req.User)
	if u == nil {
		return nil, status.Error(codes.NotFound, "Not such User")
	}

	c := u.GetCharacter(req.Character)
	if c == nil {
		return nil, status.Error(codes.NotFound, "No such Role")
	}
	c.Region = req.Region

	return &proto.None{}, nil
}

func (srv *authService) CharacterSuspend(ctx context.Context, req *proto.CharacterId) (*proto.None, error) {
	if req.User <= 0 || req.Character <= 0 {
		return nil, status.Error(codes.InvalidArgument, "Invalid user ID")
	}

	u := srv.db.UserGet(req.User)
	if u == nil {
		return nil, status.Error(codes.NotFound, "Not such User")
	}

	c := u.GetCharacter(req.Character)
	if c == nil {
		return nil, status.Error(codes.NotFound, "No such Role")
	}
	c.Region = req.Region

	return &proto.None{}, nil
}

func (srv *authService) CharacterResume(ctx context.Context, req *proto.CharacterId) (*proto.None, error) {
	if req.User <= 0 || req.Character <= 0 {
		return nil, status.Error(codes.InvalidArgument, "Invalid user ID")
	}

	u := srv.db.UserGet(req.User)
	if u == nil {
		return nil, status.Error(codes.NotFound, "Not such User")
	}

	c := u.GetCharacter(req.Character)
	if c == nil {
		return nil, status.Error(codes.NotFound, "No such Role")
	}
	c.Region = req.Region

	return &proto.None{}, nil
}
