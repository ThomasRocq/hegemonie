// Copyright (C) 2018-2020 Hegemonie's AUTHORS
// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at http://mozilla.org/MPL/2.0/.

package hegemonie_web_agent

import (
	"fmt"
	"github.com/go-macaron/session"
	region "github.com/jfsmig/hegemonie/pkg/region/proto"
	"gopkg.in/macaron.v1"
)

func expandCityView(f *frontService, lView *region.CityView) {
	f.rw.RLock()
	defer f.rw.RUnlock()

	for _, item := range lView.Assets.Units {
		item.Type = f.units[item.IdType]
	}
	for _, item := range lView.Assets.Buildings {
		item.Type = f.buildings[item.IdType]
	}
	for _, item := range lView.Assets.Knowledges {
		item.Type = f.knowledge[item.IdType]
	}
	for _, item := range lView.Assets.Armies {
		for _, u := range item.Units {
			u.Type = f.units[u.IdType]
		}
	}
}

func serveGameCityPage(f *frontService, template string) ActionPage {
	return func(ctx *macaron.Context, sess session.Store, flash *session.Flash) {
		rid := ctx.Query("rid")
		cid := ctx.Query("cid")
		uView, cView, err := f.authenticateCharacterFromSession(ctx, sess, rid, cid)
		if err != nil {
			flash.Warning(err.Error())
			ctx.Redirect("/game/user")
			return
		}

		// Load the chosen City
		cliReg := region.NewCityClient(f.cnxRegion)
		lView, err := cliReg.Show(contextMacaronToGrpc(ctx, sess),
			&region.CityId{Character: cView.Id, City: atou(ctx.Query("lid"))})
		if err != nil {
			flash.Warning("Region error: " + err.Error())
			ctx.Redirect("/game/character?cid=" + fmt.Sprint(cView.Id))
			return
		}

		expandCityView(f, lView)

		ctx.Data["Title"] = cView.Name + "|" + lView.Name
		ctx.Data["userid"] = utoa(uView.Id)
		ctx.Data["User"] = uView
		ctx.Data["cid"] = utoa(cView.Id)
		ctx.Data["Character"] = cView
		ctx.Data["lid"] = utoa(lView.Id)
		ctx.Data["Land"] = lView
		ctx.HTML(200, template)
	}
}

func serveGameCityOverview(f *frontService) ActionPage {
	return serveGameCityPage(f, "land_overview")
}

func serveGameCityBuildings(f *frontService) ActionPage {
	return serveGameCityPage(f, "land_buildings")
}

func serveGameCityKnowledges(f *frontService) ActionPage {
	return serveGameCityPage(f, "land_knowledges")
}

func serveGameCityUnits(f *frontService) ActionPage {
	return serveGameCityPage(f, "land_units")
}

func serveGameCityArmies(f *frontService) ActionPage {
	return serveGameCityPage(f, "land_armies")
}

func serveGameCityBudget(f *frontService) ActionPage {
	return serveGameCityPage(f, "land_budget")
}
