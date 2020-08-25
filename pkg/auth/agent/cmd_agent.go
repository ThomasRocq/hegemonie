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
	db  hegemonie_auth_backend.Backend
	cfg *authConfig
}

func CommandAgent() *cobra.Command {
	cfg := authConfig{}

	agent := &cobra.Command{
		Use:     "agent",
		Aliases: []string{"srv", "server", "service", "worker"},
		Short:   "Authentication service",
		RunE: func(cmd *cobra.Command, args []string) error {
			srv := authService{cfg: &cfg}
			return srv.execute()
		},
	}

	agent.Flags().StringVar(
		&cfg.endpoint, "endpoint", "127.0.0.1:8080",
		"IP:PORT endpoint for the TCP/IP server")
	agent.Flags().StringVar(
		&cfg.pathLive, "live", "",
		"Path of the DB backup to load at startup")
	agent.Flags().StringVar(
		&cfg.pathSave, "save", "",
		"Path where to save the DB backup at exit")

	return agent
}

func (srv *authService) execute() error {
	if srv.cfg.pathLive == "" {
		return errors.New("Missing: path to the live data directory")
	}

	var err error
	srv.db, err = hegemonie_auth_backend.Connect(srv.cfg.pathLive)
	if err != nil {
		return err
	}

	var lis net.Listener
	if lis, err = net.Listen("tcp", srv.cfg.endpoint); err != nil {
		return fmt.Errorf("failed to listen: %v", err)
	}

	server := grpc.NewServer(utils.ServerUnaryInterceptorZerolog())
	proto.RegisterAuthServer(server, srv)
	grpc_health_v1.RegisterHealthServer(server, srv)
	return server.Serve(lis)
}
