// Copyright (C) 2018-2020 Hegemonie's AUTHORS
// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at http://mozilla.org/MPL/2.0/.

package mapagent

import (
	"fmt"
	grpc_health_v1 "github.com/jfsmig/hegemonie/pkg/healthcheck"
	mapgraph "github.com/jfsmig/hegemonie/pkg/map/graph"
	mproto "github.com/jfsmig/hegemonie/pkg/map/proto"
	"github.com/jfsmig/hegemonie/pkg/utils"
	"github.com/spf13/cobra"
	"google.golang.org/grpc"
	"net"
)

type mapServiceConfig struct {
	endpoint       string
	pathRepository string
}

func Command() *cobra.Command {
	cfg := mapServiceConfig{}

	agent := &cobra.Command{
		Use:     "agent",
		Aliases: []string{"server"},
		Short:   "Map service",
		Example: "heged map --endpoint=10.0.0.1:1234 /path/to/maps/directory",
		Args:    cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			cfg.pathRepository = args[0]

			srv := &srvMap{config: &cfg, maps: make(mapgraph.SetOfMaps, 0)}
			if err := srv.LoadDirectory(cfg.pathRepository); err != nil {
				return err
			}

			lis, err := net.Listen("tcp", cfg.endpoint)
			if err != nil {
				return fmt.Errorf("failed to listen: %v", err)
			}

			grpcServer := grpc.NewServer(
				utils.ServerUnaryInterceptorZerolog(),
				utils.ServerStreamInterceptorZerolog())
			grpc_health_v1.RegisterHealthServer(grpcServer, &srvHealth{srv: srv})
			mproto.RegisterMapServer(grpcServer, srv)

			utils.Logger.Info().
				Int("maps", srv.maps.Len()).
				Str("endpoint", cfg.endpoint).
				Msg("Starting")
			for _, m := range srv.maps {
				utils.Logger.Debug().
					Str("name", m.ID).
					Int("sites", m.Cells.Len()).
					Int("roads", m.Roads.Len()).
					Msg("map>")
			}
			if err := grpcServer.Serve(lis); err != nil {
				return fmt.Errorf("failed to serve: %v", err)
			}

			return nil
		},
	}
	agent.Flags().StringVar(&cfg.endpoint,
		"endpoint", utils.DefaultEndpointMap, "IP:PORT endpoint for the gRPC server")

	return agent
}
