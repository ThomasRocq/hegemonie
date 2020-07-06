// Copyright (C) 2018-2020 Hegemonie's AUTHORS
// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at http://mozilla.org/MPL/2.0/.

package main

import (
	"errors"
	hegemonie_event_agent "github.com/jfsmig/hegemonie/pkg/event/agent"
	hegemonie_map_agent "github.com/jfsmig/hegemonie/pkg/map/agent"
	"github.com/jfsmig/hegemonie/pkg/utils"
	"github.com/spf13/cobra"
	"log"
)

func main() {
	rootCmd := &cobra.Command{
		Use:   "heged",
		Short: "Hegemonie main CLI",
		Long:  "Hegemonie: main binary tool to start service agents, query clients and operation jobs.",
		RunE: func(cmd *cobra.Command, args []string) error {
			return errors.New("Missing subcommand")
		},
	}
	utils.PatchCommandLogs(rootCmd)

	mapCmd := hegemonie_map_agent.Command()
	mapCmd.Use = "map"
	mapCmd.Aliases = []string{}
	rootCmd.AddCommand(mapCmd)

	evtCmd := hegemonie_event_agent.Command()
	evtCmd.Use = "event"
	evtCmd.Aliases = []string{}
	rootCmd.AddCommand(evtCmd)

	/*
		regCmd := hegemonie_region_agent.Command()
		regCmd.Use = "region"
		regCmd.Aliases = []string{}
		rootCmd.AddCommand(regCmd)

		aaaCmd := hegemonie_auth_agent.Command()
		aaaCmd.Use = "auth"
		aaaCmd.Aliases = []string{}
		rootCmd.AddCommand(aaaCmd)

		webCmd := hegemonie_web_agent.Command()
		webCmd.Use = "web"
		webCmd.Aliases = []string{}
		rootCmd.AddCommand(webCmd)
	*/

	if err := rootCmd.Execute(); err != nil {
		log.Fatalln("Command error:", err)
	}
}
