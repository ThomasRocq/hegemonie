// Copyright (C) 2020 Hegemonie's AUTHORS
// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at http://mozilla.org/MPL/2.0/.

package main

import (
	"encoding/json"
	"errors"
	auth "github.com/jfsmig/hegemonie/pkg/auth/model"
	mapgraph "github.com/jfsmig/hegemonie/pkg/map/graph"
	maputils "github.com/jfsmig/hegemonie/pkg/map/utils"
	region "github.com/jfsmig/hegemonie/pkg/region/model"
	"github.com/spf13/cobra"
	"log"
	"os"
)

func initDbAuthentication(path string) error {
	var aaa auth.Db
	aaa.Init()
	aaa.ReHash()

	u, err := aaa.CreateUser("admin@hegemonie.be")
	if err != nil {
		return err
	}
	u.Rename("Super Admin").SetRawPassword(":plop").Promote()

	_, err = aaa.CreateCharacter(u.ID, "Waku", "Calaquyr")
	if err != nil {
		return err
	}

	var f *os.File
	f, err = os.Create(path)
	if err != nil {
		return err
	}
	defer f.Close()

	encoder := json.NewEncoder(f)
	encoder.SetIndent("", " ")
	return encoder.Encode(aaa.UsersByID)
}

func loadMap(pathIn string) (mapgraph.Map, region.World, error) {
	var rawMap maputils.MapRaw
	var finalMap maputils.Map
	var graphMap mapgraph.Map
	var world region.World

	decoder := json.NewDecoder(os.Stdin)
	err := decoder.Decode(&rawMap)
	if err != nil {
		return graphMap, world, err
	}

	finalMap, err = rawMap.Finalize()
	if err != nil {
		return graphMap, world, err
	}

	world.Init()
	graphMap.Init()

	// Load the configuration, because we need models
	err = world.LoadDefinitionsFromFiles(pathIn)
	if err != nil {
		return graphMap, world, err
	}

	// Fill the world with cities and map cells
	site2cell := make(map[*maputils.Site]*mapgraph.Vertex)
	for site := range finalMap.SortedSites() {
		cell := graphMap.CellCreate()
		cell.X = uint64(site.Raw.X)
		cell.Y = uint64(site.Raw.Y)
		if site.Raw.City {
			city, err := world.CityCreateRandom(cell.ID)
			if err != nil {
				return graphMap, world, err
			}
			city.Name = site.Raw.ID
		}
		site2cell[site] = cell
	}
	for road := range finalMap.UniqueRoads() {
		src := site2cell[road.Src]
		dst := site2cell[road.Dst]
		if err = graphMap.RoadCreate(src.ID, dst.ID, true); err != nil {
			return graphMap, world, err
		}
		if err = graphMap.RoadCreate(dst.ID, src.ID, true); err != nil {
			return graphMap, world, err
		}
	}

	if err = world.PostLoad(); err != nil {
		return graphMap, world, err
	}
	if err = world.Check(); err != nil {
		return graphMap, world, err
	}

	return graphMap, world, nil
}

func CommandExport() *cobra.Command {
	var config string

	cmd := &cobra.Command{
		Use:     "init",
		Aliases: []string{"finish"},
		Short:   "Export the map as JSON files as expected by a Region agent",
		RunE: func(cmd *cobra.Command, args []string) (err error) {
			var dirOut string

			switch len(args) {
			case 0:
				return errors.New("Expected argument: path to the output directory")
			case 1:
				dirOut = args[0]
			default:
				return errors.New("")
			}

			finalMap, world, err := loadMap(config)

			err = finalMap.SaveToFile(dirOut + "/map")
			if err != nil {
				return err
			}

			err = world.SaveLiveToFiles(dirOut + "/live")
			if err != nil {
				return err
			}

			err = world.SaveDefinitionsToFiles(dirOut + "/definitions")
			if err != nil {
				return err
			}

			// Dump the authentication base
			if err != nil {
				err = initDbAuthentication(dirOut + "/auth.json")
			}

			return err
		},
	}

	cmd.Flags().StringVarP(&config, "config", "c", "", "Configuration Directory used to load the City patterns")
	return cmd
}

func main() {
	rootCmd := CommandExport()
	if err := rootCmd.Execute(); err != nil {
		log.Fatalln("Command error:", err)
	}
}
