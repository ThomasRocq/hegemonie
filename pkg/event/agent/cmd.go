// Copyright (C) 2018-2020 Hegemonie's AUTHORS
// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at http://mozilla.org/MPL/2.0/.

package hegemonie_event_agent

import (
	"errors"
	"fmt"
	grpc_health_v1 "github.com/jfsmig/hegemonie/pkg/healthcheck"
	"net"

	"github.com/spf13/cobra"
	"google.golang.org/grpc"

	back "github.com/jfsmig/hegemonie/pkg/event/backend-local"
	proto "github.com/jfsmig/hegemonie/pkg/event/proto"
	"github.com/jfsmig/hegemonie/pkg/utils"
)

type eventConfig struct {
	endpoint string
	pathBase string
}

type eventService struct {
	cfg     *eventConfig
	backend *back.Backend
}

func Command() *cobra.Command {
	cfg := eventConfig{}

	agent := &cobra.Command{
		Use:     "agent",
		Aliases: []string{"server"},
		Short:   "Authentication service",
		Example: "heged event --endpoint=10.0.0.1:2345 /path/to/event/rocksdb",
		Args:    cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			cfg.pathBase = args[0]
			srv := eventService{cfg: &cfg}
			return srv.execute()
		},
	}

	agent.Flags().StringVar(&cfg.endpoint,
		"endpoint", utils.DefaultEndpointEvent, "IP:PORT endpoint for the gRPC server")
	return agent
}

func (srv *eventService) execute() error {
	if srv.cfg.pathBase == "" {
		return errors.New("Missing: path to the live data directory")
	}

	var err error
	srv.backend, err = back.Open(srv.cfg.pathBase)
	if err != nil {
		return err
	}

	var lis net.Listener
	if lis, err = net.Listen("tcp", srv.cfg.endpoint); err != nil {
		return fmt.Errorf("failed to listen: %v", err)
	}

	server := grpc.NewServer(
		utils.ServerUnaryInterceptorZerolog(),
		utils.ServerStreamInterceptorZerolog())
	grpc_health_v1.RegisterHealthServer(server, srv)
	proto.RegisterProducerServer(server, srv)
	proto.RegisterConsumerServer(server, srv)

	utils.Logger.Info().
		Str("base", srv.cfg.pathBase).
		Str("url", srv.cfg.endpoint).
		Msg("starting")
	if err := server.Serve(lis); err != nil {
		return fmt.Errorf("failed to serve: %v", err)
	}

	return nil
}
