// Copyright (C) 2018-2020 Hegemonie's AUTHORS
// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at http://mozilla.org/MPL/2.0/.

package hegemonie_map_client

import (
	"context"
	"encoding/json"
	"errors"
	"os"
	"strconv"
	"time"

	"google.golang.org/grpc"

	proto "github.com/jfsmig/hegemonie/pkg/map/proto"
)

func getPath(args []string, cfg *eventConfig, max uint32) ([]uint64, error) {
	var src, dst uint64
	var err error
	var out []uint64

	if len(args) != 2 {
		return out, errors.New("2 arguments expected: SRC DST")
	}
	src, err = strconv.ParseUint(args[0], 10, 64)
	if err != nil {
		return out, err
	}
	dst, err = strconv.ParseUint(args[1], 10, 64)
	if err != nil {
		return out, err
	}

	ctx, _ := context.WithTimeout(context.Background(), time.Second)

	cnx, err := grpc.DialContext(ctx, cfg.endpoint, grpc.WithInsecure(), grpc.WithBlock())
	if err != nil {
		return out, err
	}
	defer cnx.Close()
	client := proto.NewMapClient(cnx)

	req := proto.PathRequest{Src: src, Dst: dst, Max: max}
	rep, err := client.GetPath(ctx, &req)
	if err != nil {
		return out, err
	}

	for _, x := range rep.Steps {
		out = append(out, x)
	}
	return out, nil
}

func getAndPrintPath(args []string, cfg *eventConfig, max uint32) error {
	path, err := getPath(args, cfg, max)
	if err != nil {
		return err
	}
	encoder := json.NewEncoder(os.Stdout)
	encoder.SetIndent("", "  ")
	encoder.Encode(path)
	return nil
}

func doPath(args []string, cfg *eventConfig) error {
	return getAndPrintPath(args, cfg, 0)
}

func doStep(args []string, cfg *eventConfig) error {
	return getAndPrintPath(args, cfg, 1)
}
