// Copyright (C) 2018-2020 Hegemonie's AUTHORS
// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at http://mozilla.org/MPL/2.0/.

package mapclient

import (
	"errors"
	"github.com/jfsmig/hegemonie/pkg/utils"
	"github.com/spf13/cobra"
)

type eventConfig struct {
	endpoint string
}

func Command() *cobra.Command {
	cfg := eventConfig{}

	cmd := &cobra.Command{
		Use:     "client",
		Aliases: []string{"cli"},
		Short:   "Event service client",
		RunE: func(cmd *cobra.Command, args []string) error {
			return errors.New("Missing subcommand")
		},
	}

	path := &cobra.Command{
		Use:   "path",
		Short: "Compute the path between two nodes",
		RunE: func(cmd *cobra.Command, args []string) error {
			return doPath(args, &cfg)
		},
	}

	step := &cobra.Command{
		Use:     "step",
		Aliases: []string{"next", "hop"},
		Short:   "Get the next step of the path between two nodes",
		RunE: func(cmd *cobra.Command, args []string) error {
			return doStep(args, &cfg)
		},
	}

	cmd.Flags().StringVar(&cfg.endpoint,
		"endpoint", utils.DefaultEndpointMap, "IP:PORT endpoint for the TCP/IP server")
	cmd.AddCommand(path, step)
	return cmd
}
