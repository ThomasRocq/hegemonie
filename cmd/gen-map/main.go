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

func generate(ctx *fasthttp.RequestCtx, r *rand.Rand) Graph {
	blocks := []Graph{
		MakeStar(r, 5, 2),
		MakeCircle(9),
		MakeCircle(7),
		MakeCircle(7),
		MakeCircle(7),
		MakeCircle(7),
		MakeCircle(7),
		MakeStar(r, 4, 2),
		MakeStar(r, 3, 2),
	}

	g := GlueChain(r, blocks...)
	Loop(r, g)
	Loop(r, g)
	return g
}

func mainHandler(ctx *fasthttp.RequestCtx) {
	var err error
	var seed int64 = time.Now().UnixNano()

	strSeed := string(ctx.QueryArgs().Peek("seed"))
	if strSeed != "" {
		seed, err = strconv.ParseInt(strSeed, 10, 63)
		if err != nil {
			ctx.Response.Header.Add("X-Error", err.Error())
			ctx.Response.SetStatusCode(fasthttp.StatusBadRequest)
			return
		}
	}

	r := rand.New(rand.NewSource(seed))

	switch string(ctx.URI().Path()) {
	case "/dot":
		g := generate(ctx, r)
		ctx.SetContentType("text/plain")
		ctx.Response.SetStatusCode(fasthttp.StatusOK)
		DrawDot(g, ctx)
	case "/svg":
		g := Simplify(generate(ctx, r))
		Noise(r, g, 1024.0, 768.0)
		FDP(r, g, 612.0, 384.0)
		Normalize(g, 1024, 768)
		ctx.SetContentType("image/svg+xml")
		ctx.Response.SetStatusCode(fasthttp.StatusOK)
		DrawSvg(g, ctx, 1024, 768)
	default:
		ctx.Response.SetStatusCode(fasthttp.StatusBadRequest)
	}
}

func main() {
	var err error
	flag.Parse()

	err = fasthttp.ListenAndServe("127.0.0.1:8080", mainHandler)
	if err != nil {
		log.Fatalln("HTTP error:", err.Error())
	}
}
