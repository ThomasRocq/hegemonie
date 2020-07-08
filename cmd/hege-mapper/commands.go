// Copyright (C) 2020 Hegemonie's AUTHORS
// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at http://mozilla.org/MPL/2.0/.

package main

import (
	"encoding/json"
	"fmt"
	"github.com/jfsmig/hegemonie/pkg/map/utils"
	"os"

	"github.com/spf13/cobra"
)

func CommandNormalize() *cobra.Command {
	return &cobra.Command{
		Use:     "normalize",
		Aliases: []string{"check", "prepare", "sanitize"},
		Short:   "Normalize a map",
		RunE: func(cmd *cobra.Command, args []string) (err error) {
			decoder := json.NewDecoder(os.Stdin)
			encoder := json.NewEncoder(os.Stdout)

			var raw maputils.MapRaw
			err = decoder.Decode(&raw)
			if err != nil {
				return err
			}

			var m maputils.Map
			m, err = raw.Finalize()
			if err != nil {
				return err
			}

			xbound, ybound := 1024.0, 768.0
			xPad, yPad := 50.0, 50.0
			m.ResizeAdjust(xbound-2*xPad, ybound-2*yPad)
			m.Center(xbound, ybound)

			raw = m.Raw()
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
		Short:   "Split a map",
		RunE: func(cmd *cobra.Command, args []string) (err error) {
			decoder := json.NewDecoder(os.Stdin)
			encoder := json.NewEncoder(os.Stdout)
			encoder.SetIndent("", " ")

			var raw maputils.MapRaw
			err = decoder.Decode(&raw)
			if err != nil {
				return err
			}

			var m maputils.Map
			m, err = raw.Finalize()
			if err != nil {
				return err
			}

			if maxDist > 0 {
				m = m.SplitLongRoads(maxDist)
			}
			xmin, xmax, ymin, ymax := m.ComputeBox()

			if noise > 0 {
				m.Noise((xmax-xmin)*(noise/100), (ymax-ymin)*(noise/100))
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
		Short:   "Convert the JSON map to DOT (graphviz)",
		RunE: func(cmd *cobra.Command, args []string) (err error) {
			decoder := json.NewDecoder(os.Stdin)

			var raw maputils.MapRaw
			err = decoder.Decode(&raw)
			if err != nil {
				return err
			}

			var m maputils.Map
			m, err = raw.Finalize()
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
		Short:   "Convert the JSON map to SVG",
		RunE: func(cmd *cobra.Command, args []string) (err error) {
			decoder := json.NewDecoder(os.Stdin)

			var raw maputils.MapRaw
			err = decoder.Decode(&raw)
			if err != nil {
				return err
			}

			var m maputils.Map
			m, err = raw.Finalize()
			if err != nil {
				return err
			}

			xbound, ybound := 1024.0, 768.0
			xPad, yPad := 50.0, 50.0
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
				if s.Raw.City {
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
