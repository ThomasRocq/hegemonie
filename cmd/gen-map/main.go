// Copyright (C) 2020 Hegemonie's AUTHORS
// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at http://mozilla.org/MPL/2.0/.

package main

import (
	"flag"
	"github.com/valyala/fasthttp"
	"log"
	"math/rand"
	"strconv"
	"time"
)

func generate(ctx *fasthttp.RequestCtx, r *rand.Rand) *Graph {
	blocks := []*Graph{
		MakeStar(r, 5, 2),
		MakeCircle(9),
		MakeCircle(9),
		MakeCircle(9),
		MakeCircle(9),
	}

	return GlueChain(r, blocks...)
}

func mainHandler(ctx *fasthttp.RequestCtx) {
	var err error
	var seed int64 = time.Now().UnixNano()
	var rounds uint = 1000
	var x, y int64 = 1024, 768

	s := string(ctx.QueryArgs().Peek("seed"))
	if s != "" {
		seed, err = strconv.ParseInt(s, 10, 63)
		if err != nil {
			ctx.Response.Header.Add("X-Error", err.Error())
			ctx.Response.SetStatusCode(fasthttp.StatusBadRequest)
			return
		}
	}
	s = string(ctx.QueryArgs().Peek("rounds"))
	if s != "" {
		u64, err := strconv.ParseUint(s, 10, 31)
		if err != nil {
			ctx.Response.Header.Add("X-Error", err.Error())
			ctx.Response.SetStatusCode(fasthttp.StatusBadRequest)
			return
		}
		rounds = uint(u64)
	}
	s = string(ctx.QueryArgs().Peek("x"))
	if s != "" {
		x, err = strconv.ParseInt(s, 10, 31)
		if err != nil {
			ctx.Response.Header.Add("X-Error", err.Error())
			ctx.Response.SetStatusCode(fasthttp.StatusBadRequest)
			return
		}
	}
	s = string(ctx.QueryArgs().Peek("y"))
	if s != "" {
		y, err = strconv.ParseInt(s, 10, 31)
		if err != nil {
			ctx.Response.Header.Add("X-Error", err.Error())
			ctx.Response.SetStatusCode(fasthttp.StatusBadRequest)
			return
		}
	}

	r := rand.New(rand.NewSource(seed))
	g := generate(ctx, r)

	switch string(ctx.URI().Path()) {
	case "/dot":
		ctx.SetContentType("text/plain")
		ctx.Response.SetStatusCode(fasthttp.StatusOK)
		DrawDot(g, ctx)
	case "/svg":
		Noise(r, g, float64(x), float64(y))
		FDP(g, float64(x), float64(y), rounds)
		Normalize(g, float64(x), float64(y))
		//LogGraph(p, g)
		ctx.SetContentType("image/svg+xml")
		ctx.Response.SetStatusCode(fasthttp.StatusOK)
		DrawSvg(g, ctx, uint(x), uint(y))
	default:
		ctx.Response.SetStatusCode(fasthttp.StatusBadRequest)
	}
}

func main() {
	var err error
	flag.Parse()

	err = fasthttp.ListenAndServe(":8080", mainHandler)
	if err != nil {
		log.Fatalln("HTTP error:", err.Error())
	}
}
