// Copyright (C) 2018-2020 Hegemonie's AUTHORS
// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at http://mozilla.org/MPL/2.0/.

package hegemonie_auth_agent

import (
	"context"
	proto "github.com/jfsmig/hegemonie/pkg/auth/proto"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (srv *authService) Login(ctx context.Context, req *proto.LoginReq) (*proto.LoginRep, error) {
	u, err := srv.users.GetByMail(req.Mail)
	if err != nil {
		return nil, status.Error(codes.Unauthenticated, "Invalid password")
	}
	if !u.Authenticate(req.Pass) {
		return nil, status.Error(codes.Unauthenticated, "Invalid password")
	}
	return &proto.LoginRep{Token: "plop"}, nil
}
