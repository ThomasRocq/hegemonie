// Copyright (C) 2018-2020 Hegemonie's AUTHORS
// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at http://mozilla.org/MPL/2.0/.

package mapclient

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/google/uuid"
	proto "github.com/jfsmig/hegemonie/pkg/map/proto"
	"github.com/jfsmig/hegemonie/pkg/utils"
	"github.com/spf13/cobra"
	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
	"io"
	"os"
	"strconv"
	"time"
)

type mapClientConfig struct {
	endpoint string
}

func (cfg *mapClientConfig) Connect() (context.Context, *grpc.ClientConn, error) {
	ctx, _ := context.WithTimeout(context.Background(), 5*time.Second)

	sessionId := os.Getenv("HEGE_CLI_SESSIONID")
	if sessionId == "" {
		sessionId = "cli/" + uuid.New().String()
	}
	ctx = metadata.AppendToOutgoingContext(ctx, "session-id", sessionId)

	cnx, err := grpc.DialContext(ctx, cfg.endpoint, grpc.WithInsecure(), grpc.WithBlock())
	return ctx, cnx, err
}

func Command() *cobra.Command {
	cfg := mapClientConfig{}

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

	cities := &cobra.Command{
		Use:     "cities",
		Aliases: []string{"city"},
		Short:   "List the Cities when the map is instantiated",
		RunE: func(cmd *cobra.Command, args []string) error {
			return doCities(args, &cfg)
		},
	}

	edges := &cobra.Command{
		Use:     "edges",
		Aliases: []string{"edge", "e"},
		Short:   "List of the edges of the map",
		RunE: func(cmd *cobra.Command, args []string) error {
			return doEdges(args, &cfg)
		},
	}

	vertices := &cobra.Command{
		Use:     "vertices",
		Aliases: []string{"vertex", "v", "sites"},
		Short:   "List the vertices of the map",
		RunE: func(cmd *cobra.Command, args []string) error {
			return doVertices(args, &cfg)
		},
	}

	local := &cobra.Command{
		Use:     "tools",
		Aliases: []string{"init", "local"},
		Short:   "Local tools to help handle maps in various forms",
		RunE: func(cmd *cobra.Command, args []string) error {
			return errors.New("Argl")
		},
	}
	cmd.Flags().StringVar(&cfg.endpoint,
		"endpoint", utils.DefaultEndpointMap, "IP:PORT endpoint for the TCP/IP server")
	cmd.AddCommand(path, step, path, cities, edges, vertices, local)
	local.AddCommand(CommandInit(), CommandDot(), CommandNormalize(), CommandSplit(), CommandSvg())
	return cmd
}

func CommandNormalize() *cobra.Command {
	return &cobra.Command{
		Use:     "normalize",
		Aliases: []string{"check", "prepare", "sanitize"},
		Short:   "Normalize the positions in a map (stdin/stdout)",
		Long:    `Read the map description on the standard input, remap the positions of the vertices in the map graph so that they fit in the given boundaries and dump it to the standard output.`,
		RunE: func(cmd *cobra.Command, args []string) (err error) {
			decoder := json.NewDecoder(os.Stdin)

			var raw MapRaw
			err = decoder.Decode(&raw)
			if err != nil {
				return err
			}

			var m MapMem
			m, err = raw.Transform()
			if err != nil {
				return err
			}

			var xbound, ybound, xPad, yPad uint64 = 1920, 1080, 50, 50
			m.ResizeAdjust(xbound-2*xPad, ybound-2*yPad)
			m.Center(xbound, ybound)

			raw = m.Raw()
			encoder := json.NewEncoder(os.Stdout)
			encoder.SetIndent("", " ")
			return encoder.Encode(&raw)
		},
	}
}

func CommandSplit() *cobra.Command {
	var maxDist float64
	var noise float64

	cmd := &cobra.Command{
		Use:     "split",
		Aliases: []string{},
		Short:   "Split the long edges of a map (stdin/stdout)",
		Long:    `Read the map on the standard input, split all the edges that are longer to the given value and dump the new graph on the standard output.`,
		RunE: func(cmd *cobra.Command, args []string) (err error) {
			decoder := json.NewDecoder(os.Stdin)
			encoder := json.NewEncoder(os.Stdout)
			encoder.SetIndent("", " ")

			var raw MapRaw
			err = decoder.Decode(&raw)
			if err != nil {
				return err
			}

			var m MapMem
			m, err = raw.Transform()
			if err != nil {
				return err
			}

			if maxDist > 0 {
				m = m.SplitLongRoads(maxDist)
			}
			xmin, xmax, ymin, ymax := m.ComputeBox()

			if noise > 0 {
				m.Noise(float64(xmax-xmin)*(noise/100), float64(ymax-ymin)*(noise/100))
			}

			raw = m.Raw()
			return encoder.Encode(&raw)
		},
	}
	cmd.Flags().Float64VarP(&maxDist, "dist", "d", 60, "Max road length")
	cmd.Flags().Float64VarP(&noise, "noise", "n", 15, "Percent of the image dimension used as max noise variation on non-city nodes positions")
	return cmd
}

func CommandDot() *cobra.Command {
	return &cobra.Command{
		Use:     "dot",
		Aliases: []string{},
		Short:   "Convert the JSON map to DOT (stdin/stdout)",
		RunE: func(cmd *cobra.Command, args []string) (err error) {
			decoder := json.NewDecoder(os.Stdin)

			var raw MapRaw
			err = decoder.Decode(&raw)
			if err != nil {
				return err
			}

			var m MapMem
			m, err = raw.Transform()
			if err != nil {
				return err
			}

			fmt.Println("graph g {")
			for r := range m.UniqueRoads() {
				fmt.Printf("%s -- %s;\n", r.Src.DotName(), r.Dst.DotName())
			}
			fmt.Println("}")
			return nil
		},
	}
}

func CommandSvg() *cobra.Command {
	var flagStandalone bool

	cmd := &cobra.Command{
		Use:     "svg",
		Aliases: []string{},
		Short:   "Convert the JSON map to SVG  (stdin/stdout)",
		RunE: func(cmd *cobra.Command, args []string) (err error) {
			decoder := json.NewDecoder(os.Stdin)

			var raw MapRaw
			err = decoder.Decode(&raw)
			if err != nil {
				return err
			}

			var m MapMem
			m, err = raw.Transform()
			if err != nil {
				return err
			}

			var xbound, ybound, xPad, yPad uint64 = 1920, 1080, 50, 50
			m.ResizeAdjust(xbound-2*xPad, ybound-2*yPad)
			m.Center(xbound, ybound)

			fmt.Printf(`<?xml version="1.0" encoding="UTF-8"?>
<!DOCTYPE svg PUBLIC "-//W3C//DTD SVG 1.1//EN" "http://www.w3.org/Graphics/SVG/1.1/DTD/svg11.dtd">
<svg xmlns="http://www.w3.org/2000/svg"
	style="background-color: rgb(255, 255, 255);"
	xmlns:xlink="http://www.w3.org/1999/xlink"
	version="1.1"
	width="%dpx" height="%dpx"
	viewBox="-0.5 -0.5 %d %d">
`, int64(xbound), int64(ybound), int64(xbound), int64(ybound))
			fmt.Println(`<g>`)
			for r := range m.UniqueRoads() {
				fmt.Printf(`<line x1="%d" y1="%d" x2="%d" y2="%d" stroke="black" stroke-width="1"/>
`, int64(r.Src.Raw.X), int64(r.Src.Raw.Y), int64(r.Dst.Raw.X), int64(r.Dst.Raw.Y))
			}
			fmt.Println(`</g>`)
			fmt.Println(`<g>`)
			for s := range m.SortedSites() {
				color := `white`
				radius := 5
				stroke := 1
				if s.Raw.City != "" {
					color = `gray`
					radius = 10
					stroke = 1
				}
				fmt.Printf(`<circle id="%s" class="clickable" cx="%d" cy="%d" r="%d" stroke="black" stroke-width="%d" fill="%s"/>
`, s.Raw.ID, int64(s.Raw.X), int64(s.Raw.Y), radius, stroke, color)
			}
			fmt.Println(`</g>`)
			fmt.Println(`</svg>`)
			return nil
		},
	}
	cmd.Flags().BoolVarP(&flagStandalone, "standalone", "1", false, "Also generate the xml header")
	return cmd
}

func CommandInit() *cobra.Command {
	var flagStandalone bool

	cmd := &cobra.Command{
		Use:     "init",
		Aliases: []string{"seed"},
		Short:   "Convert the JSON map seed to a JSON raw map (stdin/stdout)",
		RunE: func(cmd *cobra.Command, args []string) (err error) {
			decoder := json.NewDecoder(os.Stdin)

			var in MapSeed
			err = decoder.Decode(&in)
			if err != nil {
				return err
			}

			var out MapRaw
			out, err = in.Transform()
			if err != nil {
				return err
			}

			encoder := json.NewEncoder(os.Stdout)
			encoder.SetIndent("", " ")
			return encoder.Encode(&out)
		},
	}
	cmd.Flags().BoolVarP(&flagStandalone, "standalone", "1", false, "Also generate the xml header")
	return cmd
}

func getPath(args []string, cfg *mapClientConfig, max uint32) ([]uint64, error) {
	var err error
	var out []uint64
	var cnx *grpc.ClientConn
	var req proto.PathRequest

	if len(args) != 3 {
		return out, errors.New("2 arguments expected: MAP SRC DST")
	}

	req.MapName = args[0]
	req.Src, err = strconv.ParseUint(args[1], 10, 64)
	if err != nil {
		return out, err
	}
	req.Dst, err = strconv.ParseUint(args[2], 10, 64)
	if err != nil {
		return out, err
	}

	ctx, cnx, err := cfg.Connect()
	if err != nil {
		return out, err
	}
	defer cnx.Close()
	client := proto.NewMapClient(cnx)

	rep, err := client.GetPath(ctx, &req)
	if err != nil {
		return out, err
	}

	for i := uint32(0); max <= 0 || i < max; i++ {
		x, err := rep.Recv()
		if err != nil {
			if err == io.EOF {
				break
			}
			return out, err
		}
		out = append(out, x.GetId())
	}
	return out, nil
}

func getAndPrintPath(args []string, cfg *mapClientConfig, max uint32) error {
	path, err := getPath(args, cfg, max)
	if err != nil {
		return err
	}
	encoder := json.NewEncoder(os.Stdout)
	encoder.SetIndent("", "  ")
	encoder.Encode(path)
	return nil
}

func doPath(args []string, cfg *mapClientConfig) error {
	return getAndPrintPath(args, cfg, 0)
}

func doStep(args []string, cfg *mapClientConfig) error {
	return getAndPrintPath(args, cfg, 1)
}

func doCities(args []string, cfg *mapClientConfig) error {
	var err error
	var cnx *grpc.ClientConn
	var req proto.ListCitiesReq

	if len(args) <= 0 || len(args) > 2 {
		return errors.New("2 arguments expected: MAP [MARKER]")
	}

	req.MapName = args[0]
	if len(args) > 1 {
		req.Marker, err = strconv.ParseUint(args[1], 10, 64)
		if err != nil {
			return err
		}
	}

	ctx, cnx, err := cfg.Connect()
	if err != nil {
		return err
	}
	defer cnx.Close()
	client := proto.NewMapClient(cnx)

	rep, err := client.Cities(ctx, &req)
	if err != nil {
		return err
	}

	out := make([]uint64, 0)
	for {
		x, err := rep.Recv()
		if err != nil {
			if err == io.EOF {
				break
			}
			return err
		}
		out = append(out, x.GetId())
	}

	encoder := json.NewEncoder(os.Stdout)
	encoder.SetIndent("", "  ")
	encoder.Encode(out)
	return nil
}

func doEdges(args []string, cfg *mapClientConfig) error {
	var err error
	var cnx *grpc.ClientConn
	var req proto.ListEdgesReq

	if len(args) <= 0 || len(args) > 3 {
		return errors.New("2 arguments expected: MAP [SRC [DST]]")
	}

	req.MapName = args[0]
	if len(args) > 1 {
		req.MarkerSrc, err = strconv.ParseUint(args[1], 10, 64)
		if err != nil {
			return err
		}
		if len(args) > 2 {
			req.MarkerDst, err = strconv.ParseUint(args[2], 10, 64)
			if err != nil {
				return err
			}
		}
	}

	ctx, cnx, err := cfg.Connect()
	if err != nil {
		return err
	}
	defer cnx.Close()
	client := proto.NewMapClient(cnx)

	rep, err := client.Edges(ctx, &req)
	if err != nil {
		return err
	}

	type Pair struct{ Src, Dst uint64 }
	out := make([]Pair, 0)
	for {
		x, err := rep.Recv()
		if err != nil {
			if err == io.EOF {
				break
			}
			return err
		}
		out = append(out, Pair{x.GetSrc(), x.GetDst()})
	}

	encoder := json.NewEncoder(os.Stdout)
	encoder.SetIndent("", "  ")
	encoder.Encode(out)
	return nil
}

func doVertices(args []string, cfg *mapClientConfig) error {
	var err error
	var cnx *grpc.ClientConn
	var req proto.ListVerticesReq

	if len(args) <= 0 || len(args) > 2 {
		return errors.New("2 arguments expected: MAP [SRC_MARKER [DST_MARKER]]")
	}

	req.MapName = args[0]
	if len(args) > 1 {
		req.Marker, err = strconv.ParseUint(args[1], 10, 64)
		if err != nil {
			return err
		}
	}

	ctx, cnx, err := cfg.Connect()
	if err != nil {
		return err
	}
	defer cnx.Close()
	client := proto.NewMapClient(cnx)

	rep, err := client.Vertices(ctx, &req)
	if err != nil {
		return err
	}

	type V struct{ Id, X, D uint64 }
	out := make([]V, 0)
	for {
		x, err := rep.Recv()
		if err != nil {
			if err == io.EOF {
				break
			}
			return err
		}
		out = append(out, V{x.GetId(), x.GetX(), x.GetY()})
	}

	encoder := json.NewEncoder(os.Stdout)
	encoder.SetIndent("", "  ")
	encoder.Encode(out)
	return nil
}
