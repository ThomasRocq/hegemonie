// Copyright (C) 2018-2020 Hegemonie's AUTHORS
// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at http://mozilla.org/MPL/2.0/.

package hegemonie_region_agent

import (
	"errors"
	"fmt"
	grpc_health_v1 "github.com/jfsmig/hegemonie/pkg/healthcheck"
	"github.com/jfsmig/hegemonie/pkg/region/model"
	rproto "github.com/jfsmig/hegemonie/pkg/region/proto"
	"github.com/jfsmig/hegemonie/pkg/utils"
	"github.com/spf13/cobra"
	"google.golang.org/grpc"
	"net"
)

type regionConfig struct {
	endpoint      string
	endpointEvent string
	backend       string
}

func Command() *cobra.Command {
	cfg := regionConfig{}

	agent := &cobra.Command{
		Use:     "agent",
		Aliases: []string{"srvCity", "srv", "service"},
		Short:   "Region service",
		RunE: func(cmd *cobra.Command, args []string) error {
			return cfg.execute()
		},
	}
	agent.Flags().StringVar(&cfg.endpoint,
		"endpoint", utils.DefaultEndpointRegion, "IP:PORT endpoint for the TCP/IP server")
	agent.Flags().StringVar(&cfg.endpointEvent,
		"event", utils.DefaultEndpointEvent, "Address of the Event server to connect to.")
	agent.Flags().StringVar(&cfg.backend,
		"defs", "", "Path to the file with the definition of the world.")

	return agent
}

func (cfg *regionConfig) execute() error {
	var err error

	w := region.World{}
	w.Init()

	if cfg.backend == "" {
		return errors.New("Missing path for live data")
	}

	err = w.Sections(cfg.backend).Load()
	if err != nil {
		return err
	}

	err = w.PostLoad()
	if err != nil {
		return fmt.Errorf("Inconsistent World from [%s]: %v", cfg.backend, err)
	}

	err = w.Check()
	if err != nil {
		return fmt.Errorf("Inconsistent World from [%s]: %v", cfg.backend, err)
	}

	lis, err := net.Listen("tcp", cfg.endpoint)
	if err != nil {
		return fmt.Errorf("failed to listen: %v", err)
	}

	var cnxEvent *grpc.ClientConn
	cnxEvent, err = grpc.Dial(cfg.endpointEvent, grpc.WithInsecure())
	if err != nil {
		return err
	}
	defer cnxEvent.Close()
	w.SetNotifier(&EventStore{cnx: cnxEvent})

	srv := grpc.NewServer(utils.ServerUnaryInterceptorZerolog())
	rproto.RegisterCityServer(srv, &srvCity{cfg: cfg, w: &w})
	rproto.RegisterDefinitionsServer(srv, &srvDefinitions{cfg: cfg, w: &w})
	rproto.RegisterAdminServer(srv, &srvAdmin{cfg: cfg, w: &w})
	rproto.RegisterArmyServer(srv, &srvArmy{cfg: cfg, w: &w})
	grpc_health_v1.RegisterHealthServer(srv, &srvHealth{w: &w})

	if err := srv.Serve(lis); err != nil {
		return fmt.Errorf("failed to serve: %v", err)
	}

	return nil
}
