// Copyright (C) 2018-2020 Hegemonie's AUTHORS
// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at http://mozilla.org/MPL/2.0/.

package hegemonie_map_agent

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
	endpoint string
	pathLive string
}

func Command() *cobra.Command {
	cfg := regionConfig{}

	agent := &cobra.Command{
		Use:     "agent",
		Aliases: []string{"srv", "service"},
		Short:   "Region service",
		RunE: func(cmd *cobra.Command, args []string) error {
			return cfg.execute()
		},
	}
	agent.Flags().StringVar(&cfg.endpoint,
		"endpoint", utils.DefaultEndpointRegion, "IP:PORT endpoint for the TCP/IP server")
	agent.Flags().StringVar(&cfg.pathLive,
		"live", "", "Path to the file with the state of the region.")

	return agent
}

func (cfg *regionConfig) execute() error {
	var err error

	w := mapgraph.Map{}
	w.Init()

	if cfg.pathLive == "" {
		return errors.New("Missing path for live data")
	}

	err = w.LoadFromFiles(cfg.pathLive)
	if err != nil {
		return err
	}

	err = w.PostLoad()
	if err != nil {
		return fmt.Errorf("Inconsistent Map from [%s] and [%s]: %v", cfg.pathLive, cfg.pathLive, err)
	}

	err = w.Check()
	if err != nil {
		return fmt.Errorf("Inconsistent World: %v", err)
	}

	lis, err := net.Listen("tcp", cfg.endpoint)
	if err != nil {
		return fmt.Errorf("failed to listen: %v", err)
	}

	srv := grpc.NewServer(utils.ServerUnaryInterceptorZerolog())
	mproto.RegisterMapServer(srv, &srvMap{cfg: cfg, w: &w})
	grpc_health_v1.RegisterHealthServer(srv, &srvHealth{w: &w})

	if err := srv.Serve(lis); err != nil {
		return fmt.Errorf("failed to serve: %v", err)
	}

	return nil
}
