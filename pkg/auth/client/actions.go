// Copyright (C) 2018-2020 Hegemonie's AUTHORS
// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at http://mozilla.org/MPL/2.0/.

package hegemonie_auth_client

import (
	"context"
	"encoding/json"
	"errors"
	proto "github.com/jfsmig/hegemonie/pkg/auth/proto"
	"github.com/spf13/cobra"
	"google.golang.org/grpc"
	"log"
	"net/mail"
	"os"
)

func connectAndDo(cfg *authConfig, action func(proto.AuthClient) error) error {
	cnx, err := grpc.Dial(cfg.endpoint, grpc.WithInsecure(), grpc.WithBlock())
	if err != nil {
		return err
	}
	defer cnx.Close()
	client := proto.NewAuthClient(cnx)
	return action(client)
}

func doShow(cmd *cobra.Command, args []string, cfg *authConfig) error {
	return connectAndDo(cfg, func(client proto.AuthClient) error {
		for _, a := range args {
			view, err := client.UserShow(context.Background(),
				&proto.UserShowReq{Mail: a})
			if err != nil {
				log.Printf("%v : %v", a, err)
			} else {
				enc := json.NewEncoder(os.Stdout)
				enc.SetIndent("", " ")
				enc.Encode(view)
			}
		}
		return nil
	})
}

func doCreate(cmd *cobra.Command, args []string, cfg *authConfig) error {
	return connectAndDo(cfg, func(client proto.AuthClient) error {
		if len(args) <= 0 {
			return errors.New("Missing argument, at least 1 email address is expected")
		}
		for _, a := range args {
			addr, err := mail.ParseAddress(a)
			if err != nil {
				log.Printf("Invalid e-mail address (%s): %s", a, err.Error())
				continue
			}
			u, err := client.UserCreate(context.Background(),
				&proto.UserCreateReq{Mail: addr.Address})
			if err != nil {
				log.Println("ERR", a, err)
			} else {
				enc := json.NewEncoder(os.Stdout)
				enc.SetIndent("", " ")
				enc.Encode(u)
			}
		}
		return nil
	})
}

func doList(cmd *cobra.Command, args []string, cfg *authConfig) error {
	return connectAndDo(cfg, func(client proto.AuthClient) error {
		var last uint64
		for {
			rep, err := client.UserList(context.Background(),
				&proto.UserListReq{Marker: last, Limit: 100})
			if err != nil {
				return err
			}
			if len(rep.Items) <= 0 {
				break
			}
			for _, u := range rep.Items {
				if u.UserId > last {
					last = u.UserId
				}
				enc := json.NewEncoder(os.Stdout)
				enc.Encode(u)
			}
		}
		return nil
	})
}
