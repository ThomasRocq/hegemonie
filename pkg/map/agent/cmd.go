// Copyright (C) 2018-2020 Hegemonie's AUTHORS
// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at http://mozilla.org/MPL/2.0/.

package mapagent

import (
	"errors"
	"fmt"
	grpc_health_v1 "github.com/jfsmig/hegemonie/pkg/healthcheck"
	mapgraph "github.com/jfsmig/hegemonie/pkg/map/graph"
	mproto "github.com/jfsmig/hegemonie/pkg/map/proto"
	"github.com/jfsmig/hegemonie/pkg/utils"
	"github.com/spf13/cobra"
	"google.golang.org/grpc"
	"net"
)

type regionConfig struct {
	endpoint       string
	pathRepository string
}

func Command() *cobra.Command {
	cfg := regionConfig{}

	agent := &cobra.Command{
		Use:     "agent",
		Aliases: []string{"srv", "service"},
		Short:   "Region service",
		RunE: func(cmd *cobra.Command, args []string) error {
			if cfg.pathRepository == "" {
				return errors.New("Missing path for live data")
			}

			repo := mapgraph.NewRepository()
			if err := mapgraph.LoadDirectory(repo, cfg.pathRepository); err != nil {
				return err
			}

			lis, err := net.Listen("tcp", cfg.endpoint)
			if err != nil {
				return fmt.Errorf("failed to listen: %v", err)
			}

			grpcServer := grpc.NewServer(utils.ServerUnaryInterceptorZerolog())
			mproto.RegisterMapServer(grpcServer, &srvMap{config: &cfg, repo: repo})
			grpc_health_v1.RegisterHealthServer(grpcServer, &srvHealth{repo: repo})

			if err := grpcServer.Serve(lis); err != nil {
				return fmt.Errorf("failed to serve: %v", err)
			}

			return nil

		},
	}
	agent.Flags().StringVar(&cfg.endpoint,
		"endpoint", utils.DefaultEndpointRegion, "IP:PORT endpoint for the TCP/IP server")
	agent.Flags().StringVar(&cfg.pathRepository,
		"data", "", "Path to the file with the state of the region.")

	return agent
}
