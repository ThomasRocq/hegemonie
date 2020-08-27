// Copyright (C) 2018-2020 Hegemonie's AUTHORS
// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at http://mozilla.org/MPL/2.0/.

package hegemonie_auth_agent

import (
	"errors"
	"fmt"
	"github.com/jfsmig/hegemonie/pkg/auth/backend"
	proto "github.com/jfsmig/hegemonie/pkg/auth/proto"
	grpc_health_v1 "github.com/jfsmig/hegemonie/pkg/healthcheck"
	"github.com/jfsmig/hegemonie/pkg/utils"
	"github.com/spf13/cobra"
	"google.golang.org/grpc"
	"net"
)

type authConfig struct {
	endpoint string
	pathLive string
	pathSave string
}

type authService struct {
	users hegemonie_auth_backend.UserBackend
	cfg   *authConfig
}

type userService struct {
	users      hegemonie_auth_backend.UserBackend
	characters hegemonie_auth_backend.CharacterBackend
	cfg        *authConfig
}

type charactersService struct {
	users      hegemonie_auth_backend.UserBackend
	characters hegemonie_auth_backend.CharacterBackend
	cfg        *authConfig
}

func CommandAgent() *cobra.Command {
	cfg := authConfig{}

	agent := &cobra.Command{
		Use:     "agent",
		Aliases: []string{"srv", "server", "service", "worker"},
		Short:   "Authentication service",
		RunE: func(cmd *cobra.Command, args []string) error {
			return cfg.execute()
		},
	}

	agent.Flags().StringVar(
		&cfg.endpoint, "endpoint", "127.0.0.1:8080",
		"IP:PORT endpoint for the TCP/IP server")
	agent.Flags().StringVar(
		&cfg.pathLive, "backend", ":mem:",
		"Path of the DB backup to load at startup")

	return agent
}

func (cfg *authConfig) execute() error {
	if cfg.pathLive == "" {
		return errors.New("Missing: path to the live data directory")
	}

	var err error
	var u hegemonie_auth_backend.UserBackend
	var c hegemonie_auth_backend.CharacterBackend

	u, err = hegemonie_auth_backend.ConnectUserBackend(cfg.pathLive)
	if err != nil {
		return err
	}
	c, err = hegemonie_auth_backend.ConnectCharacterBackend(cfg.pathLive)
	if err != nil {
		return err
	}

	aSrv := authService{users: u, cfg: cfg}
	cSrv := charactersService{users: u, characters: c, cfg: cfg}
	uSrv := userService{users: u, characters: c, cfg: cfg}

	var lis net.Listener
	if lis, err = net.Listen("tcp", cfg.endpoint); err != nil {
		return fmt.Errorf("failed to listen: %v", err)
	}

	server := grpc.NewServer(utils.ServerUnaryInterceptorZerolog())
	proto.RegisterAuthServer(server, &aSrv)
	proto.RegisterUserServer(server, &uSrv)
	proto.RegisterCharacterServer(server, &cSrv)
	grpc_health_v1.RegisterHealthServer(server, &aSrv)
	return server.Serve(lis)
}
