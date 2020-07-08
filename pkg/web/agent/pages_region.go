// Copyright (C) 2018-2020 Hegemonie's AUTHORS
// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at http://mozilla.org/MPL/2.0/.

package hegemonie_web_agent

import (
	"github.com/go-macaron/session"
	mproto "github.com/jfsmig/hegemonie/pkg/map/proto"
	rproto "github.com/jfsmig/hegemonie/pkg/region/proto"
	"gopkg.in/macaron.v1"
)

type rawVertex struct {
	ID   uint64 `json:"id"`
	X    uint64 `json:"x"`
	Y    uint64 `json:"y"`
}

type rawEdge struct {
	Src uint64 `json:"src"`
	Dst uint64 `json:"dst"`
}

type rawCity struct {
	ID   uint64 `json:"id"`
	Name string `json:"name"`
}

type rawMap struct {
	Cells map[uint64]rawVertex `json:"cells"`
	Roads []rawEdge            `json:"roads"`
}

func serveRegionMap(f *frontService) NoFlashPage {
	return func(ctx *macaron.Context, sess session.Store) {
		id := ctx.Query("id")
		if id != "calaquyr" {
			ctx.Error(400, "Invalid region")
			return
		}

		m := rawMap{
			Cells: make(map[uint64]rawVertex),
			Roads: make([]rawEdge, 0),
		}
		cli := mproto.NewMapClient(f.cnxRegion)
		ctx0 := contextMacaronToGrpc(ctx, sess)

		// FIXME(jfs): iterate in case of a truncated result
		vertices, err := f.loadAllLocations(ctx0, cli)
		if err != nil {
			ctx.Error(502, err.Error())
			return
		}
		for _, v := range vertices {
			m.Cells[v.Id] = rawVertex{ID: v.Id, X: v.X, Y: v.Y}
		}

		// FIXME(jfs): iterate in case of a truncated result
		edges, err := f.loadAllRoads(ctx0, cli)
		if err != nil {
			ctx.Error(502, err.Error())
			return
		}
		for _, e := range edges {
			m.Roads = append(m.Roads, rawEdge{Src: e.Src, Dst: e.Dst})
		}

		ctx.JSON(200, m)
	}
}

func serveRegionCities(f *frontService) NoFlashPage {
	return func(ctx *macaron.Context, sess session.Store) {
		id := ctx.Query("id")
		if id != "calaquyr" {
			ctx.Error(400, "Invalid region")
			return
		}

		tab := make([]rawCity, 0)
		cli := rproto.NewCityClient(f.cnxRegion)

		// FIXME(jfs): iterate in case of a truncated result
		cities, err := f.loadAllCities(contextMacaronToGrpc(ctx, sess), cli)
		if err != nil {
			ctx.Error(502, err.Error())
			return
		}
		for _, v := range cities {
			tab = append(tab, rawCity{ID: v.Id, Name: v.Name})
		}

		ctx.JSON(200, tab)
	}
}
