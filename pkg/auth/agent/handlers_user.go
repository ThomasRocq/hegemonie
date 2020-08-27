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
	case hegemonie_auth_backend.UserStateDeleted:
		return proto.UserState_UserDeleted
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

func (srv *userService) List(req *proto.UserListReq, stream proto.User_ListServer) error {
	if req.Limit <= 0 {
		req.Limit = 1024
	}

	tab, err := srv.users.List(req.Marker, req.Limit)
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

func (srv *userService) Create(ctx context.Context, req *proto.UserCreateReq) (*proto.UserView, error) {
	var err error
	u, err := srv.users.Create(req.Mail, req.Pass)
	if err != nil {
		return nil, err
	}
	return userView(u), err
}

func (srv *userService) Delete(ctx context.Context, req *proto.UserId) (*proto.None, error) {
	err := srv.users.Delete(req.Id)
	if err != nil {
		return nil, err
	}
	return none, nil
}

func (srv *userService) Show(ctx context.Context, req *proto.UserId) (*proto.UserView, error) {
	u, err := srv.users.Show(req.Id)
	if err != nil {
		return nil, status.Error(codes.NotFound, "No such User")
	}
	return userView(u), nil
}

func (srv *userService) GetByMail(ctx context.Context, req *proto.UserMail) (*proto.UserView, error) {
	u, err := srv.users.GetByMail(req.Mail)
	if err != nil {
		return nil, status.Error(codes.NotFound, "No such User")
	}
	return userView(u), nil
}

func (srv *userService) Promote(ctx context.Context, req *proto.UserId) (*proto.None, error) {
	err := srv.users.Promote(req.Id)
	if err != nil {
		return nil, status.Error(codes.NotFound, "No such User")
	}
	return none, nil
}

func (srv *userService) Demote(ctx context.Context, req *proto.UserId) (*proto.None, error) {
	err := srv.users.Demote(req.Id)
	if err != nil {
		return nil, status.Error(codes.NotFound, "No such User")
	}
	return none, nil
}

func (srv *userService) Suspend(ctx context.Context, req *proto.UserId) (*proto.None, error) {
	err := srv.users.Suspend(req.Id)
	if err != nil {
		return nil, status.Error(codes.NotFound, "No such User")
	}
	return none, nil
}

func (srv *userService) Resume(ctx context.Context, req *proto.UserId) (*proto.None, error) {
	err := srv.users.Resume(req.Id)
	if err != nil {
		return nil, status.Error(codes.NotFound, "No such User")
	}
	return none, nil
}

func (srv *userService) Characters(req *proto.UserId, stream proto.User_CharactersServer) error {
	tab, err := srv.users.Characters(req.Id)
	if err != nil {
		return err
	}
	for _, c := range tab {
		err = stream.Send(characterView(c))
		if err == io.EOF {
			return nil
		}
		return err
	}
	return nil
}
