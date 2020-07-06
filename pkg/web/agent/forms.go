// Copyright (C) 2018-2020 Hegemonie's AUTHORS
// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at http://mozilla.org/MPL/2.0/.

package hegemonie_web_agent

import (
	"errors"
	"github.com/go-macaron/session"
	"github.com/google/uuid"
	auth "github.com/jfsmig/hegemonie/pkg/auth/proto"
	"gopkg.in/macaron.v1"
)

var (
	errAuthFailed = errors.New("No authenticated")
)

func (f *frontService) authenticateUserFromSession(ctx *macaron.Context, sess session.Store) (*auth.UserView, error) {
	// Validate the session data
	p := sess.Get("userid")
	if p == nil {
		return nil, errAuthFailed
	}
	userID, ok := p.(string)
	if !ok || userID == "" {
		return nil, errAuthFailed
	}

	// Authorize the character with the user
	cliAuth := auth.NewUserClient(f.cnxAuth)
	return cliAuth.Show(contextMacaronToGrpc(ctx, sess), &auth.UserId{Id: userID})
}

func (f *frontService) authenticateAdminFromSession(ctx *macaron.Context, sess session.Store) (*auth.UserView, error) {
	uView, err := f.authenticateUserFromSession(ctx, sess)
	if err != nil {
		return nil, err
	}
	if uView.State != auth.UserState_UserAdmin {
		return nil, errAuthFailed
	}
	return uView, nil
}

func (f *frontService) authenticateCharacterFromSession(ctx *macaron.Context, sess session.Store, idRegion, idChar string) (*auth.UserView, *auth.CharacterView, error) {
	// Validate the session data
	p := sess.Get("userid")
	if p == nil {
		return nil, nil, errAuthFailed
	}
	userID, ok := p.(string)
	if !ok || userID == "" {
		return nil, nil, errAuthFailed
	}

	// Authorize the character with the user
	ctx1 := contextMacaronToGrpc(ctx, sess)
	cliUser := auth.NewUserClient(f.cnxAuth)
	uView, err := cliUser.Show(ctx1, &auth.UserId{Id: userID})
	if err != nil {
		return nil, nil, err
	}
	cliCharacter := auth.NewCharacterClient(f.cnxAuth)
	cView, err := cliCharacter.Show(ctx1, &auth.CharacterId{Region: idRegion, Name: idChar})
	if err != nil {
		return nil, nil, err
	}
	return uView, cView, nil
}

type FormLogin struct {
	UserMail string `form:"email" binding:"Required"`
	UserPass string `form:"password" binding:"Required"`
}

func doLogin(f *frontService) macaron.Handler {
	return func(ctx *macaron.Context, flash *session.Flash, sess session.Store, info FormLogin) {
		// Cleanup a previous session
		sess.Flush()

		ctx1 := contextMacaronToGrpc(ctx, sess)

		// Authorize the character with the user
		cliAuth := auth.NewAuthClient(f.cnxAuth)
		jwt, err := cliAuth.Login(ctx1,
			&auth.LoginReq{Mail: info.UserMail, Pass: []byte(info.UserPass)})
		if err != nil {
			flash.Warning(err.Error())
			ctx.Redirect("/")
		}

		// FIXME(jfs): Need for a RPC or just get the info from the token?
		cliUser := auth.NewUserClient(f.cnxAuth)
		uView, err := cliUser.GetByMail(ctx1, &auth.UserMail{Mail: info.UserMail})

		if err != nil {
			flash.Warning(err.Error())
			ctx.Redirect("/")
		} else {
			ctx.SetSecureCookie("session", uView.UserId)
			sessionID := uuid.New().String()
			sess.Set("session-id", sessionID)
			sess.Set("userid", uView.UserId)
			sess.Set("authorization", jwt.Token)

			ctx.Redirect("/game/user")
		}
	}
}

func doLogout(f *frontService) macaron.Handler {
	return func(ctx *macaron.Context, s session.Store) {
		ctx.SetSecureCookie("session", "")
		s.Flush()
		ctx.Redirect("/")
	}
}
