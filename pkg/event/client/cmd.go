// Copyright (C) 2018-2020 Hegemonie's AUTHORS
// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at http://mozilla.org/MPL/2.0/.

package hegemonie_event_client

import (
	"context"
	"errors"
	"fmt"
	"github.com/google/uuid"
	proto "github.com/jfsmig/hegemonie/pkg/event/proto"
	"github.com/jfsmig/hegemonie/pkg/utils"
	"github.com/spf13/cobra"
	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
	"os"
	"strconv"
	"time"
)

type eventClientConfig struct {
	endpoint string
}

func Command() *cobra.Command {
	cfg := eventClientConfig{}

	cmd := &cobra.Command{
		Use:     "client",
		Aliases: []string{"cli"},
		Short:   "Event service client",
		RunE: func(cmd *cobra.Command, args []string) error {
			return errors.New("Missing subcommand")
		},
	}

	push := &cobra.Command{
		Use:     "push",
		Short:   "Push an event",
		Args:    cobra.MinimumNArgs(2),
		Example: `hege event push "$CHARACTER" 'Higher deaths rates are reported on your city'`,
		RunE: func(cmd *cobra.Command, args []string) error {
			return doPush(args, &cfg)
		},
	}

	list := &cobra.Command{
		Use:     "list",
		Aliases: []string{"ls"},
		Short:   "List the events",
		Args:    cobra.RangeArgs(1, 2),
		Example: `hege event list "$CHARACTER" "$AFTER""`,
		RunE: func(cmd *cobra.Command, args []string) error {
			return doList(args, &cfg)
		},
	}

	ack := &cobra.Command{
		Use:     "ack",
		Short:   "Acknowledge an event",
		Args:    cobra.ExactArgs(3),
		Example: `hege event ack "$CHARACTER" "$EVENT_TIMESTAMP" "$EVENT_UUID"`,
		RunE: func(cmd *cobra.Command, args []string) error {
			return doAck(args, &cfg)
		},
	}

	cmd.Flags().StringVar(&cfg.endpoint,
		"endpoint", utils.DefaultEndpointEvent, "IP:PORT endpoint for the TCP/IP server")
	cmd.AddCommand(push, ack, list)
	return cmd
}

func (cfg *eventClientConfig) Connect() (context.Context, *grpc.ClientConn, error) {
	ctx, _ := context.WithTimeout(context.Background(), 5*time.Second)

	sessionId := os.Getenv("HEGE_CLI_SESSIONID")
	if sessionId == "" {
		sessionId = "cli/" + uuid.New().String()
	}
	ctx = metadata.AppendToOutgoingContext(ctx, "session-id", sessionId)

	cnx, err := grpc.DialContext(ctx, cfg.endpoint, grpc.WithInsecure(), grpc.WithBlock())
	return ctx, cnx, err
}

func doPush(args []string, cfg *eventClientConfig) error {
	// TODO(jfs): validate the format of charID
	charID := args[0]

	ctx, cnx, err := cfg.Connect()
	if err != nil {
		return err
	}
	defer cnx.Close()

	client := proto.NewProducerClient(cnx)
	anyError := false
	for _, a := range args[1:] {
		id := uuid.New().String()
		req := proto.Push1Req{
			CharId:  charID,
			EvtId:   id,
			Payload: []byte(a),
		}
		_, err := client.Push1(ctx, &req)
		if err != nil {
			anyError = true
			utils.Logger.Error().Str("char", charID).Str("msg", a).Str("uuid", id).Err(err).Msg("PUSH")
		} else {
			utils.Logger.Info().Str("char", charID).Str("msg", a).Str("uuid", id).Msg("PUSH")
		}
	}
	if !anyError {
		return nil
	}
	return errors.New("Errors occured")
}

func doAck(args []string, cfg *eventClientConfig) error {
	var when uint64
	var charID, evtID string
	var err error

	// TODO(jfs): validate the format of charID
	charID = args[0]
	if when, err = strconv.ParseUint(args[1], 10, 64); err != nil {
		return err
	}
	if s, err := uuid.Parse(args[2]); err != nil {
		return err
	} else {
		evtID = s.String()
	}

	// Send the request
	ctx, cnx, err := cfg.Connect()
	if err != nil {
		return err
	}
	defer cnx.Close()

	client := proto.NewConsumerClient(cnx)
	req := proto.Ack1Req{CharId: charID, When: when, EvtId: evtID}
	_, err = client.Ack1(ctx, &req)

	if err != nil {
		return err
	}
	utils.Logger.Info().
		Str("char", charID).
		Uint64("when", when).
		Str("uuid", evtID).
		Msg("ACK")
	return nil
}

func doList(args []string, cfg *eventClientConfig) error {
	var when uint64
	var charId string
	var err error

	// TODO(jfs): validate the format of charID
	charId = args[0]

	if len(args) == 2 {
		when, err = strconv.ParseUint(args[1], 10, 64)
		if err != nil {
			return err
		}
	}

	// Send the request
	ctx, cnx, err := cfg.Connect()
	if err != nil {
		return err
	}
	defer cnx.Close()

	client := proto.NewConsumerClient(cnx)
	req := proto.ListReq{CharId: charId, Marker: when, Max: 100}
	rep, err := client.List(ctx, &req)

	if err != nil {
		return err
	}
	anyError := false
	for _, x := range rep.Items {
		fmt.Printf("%s %d %s %s\n", x.CharId, x.When, x.EvtId, x.Payload)
	}
	if anyError {
		return errors.New("Invalid events matched")
	}
	return nil
}
