// Copyright (C) 2018-2020 Hegemonie's AUTHORS
// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at http://mozilla.org/MPL/2.0/.

package hegemonie_web_agent

import (
	"github.com/go-macaron/session"
	region "github.com/jfsmig/hegemonie/pkg/region/proto"
	"gopkg.in/macaron.v1"
)

type FormCityStudy struct {
	RegionID    string `form:"reg" binding:"Required"`
	CharacterID string `form:"cid" binding:"Required"`
	CityID      uint64 `form:"lid" binding:"Required"`
	KnowledgeID uint64 `form:"kid" binding:"Required"`
}

type FormCityBuild struct {
	RegionID    string `form:"reg" binding:"Required"`
	CharacterID string `form:"cid" binding:"Required"`
	CityID      uint64 `form:"lid" binding:"Required"`
	BuildingID  uint64 `form:"bid" binding:"Required"`
}

type FormCityTrain struct {
	RegionID    string `form:"reg" binding:"Required"`
	CharacterID string `form:"cid" binding:"Required"`
	CityID      uint64 `form:"lid" binding:"Required"`
	UnitID      uint64 `form:"uid" binding:"Required"`
}

type FormCityUnitTransfer struct {
	RegionID    string `form:"reg" binding:"Required"`
	CharacterID string `form:"cid" binding:"Required"`
	CityID      uint64 `form:"lid" binding:"Required"`
	UnitID      string `form:"uid" binding:"Required"`
	ArmyID      string `form:"aid" binding:"Required"`
}

type FormCityStockTransfer struct {
	// Identifier of the city
	RegionID    string `form:"reg" binding:"Required"`
	CharacterID string `form:"cid" binding:"Required"`
	CityID      uint64 `form:"lid" binding:"Required"`
	ArmyID      string `form:"aid" binding:"Required"`

	// Resources to be transferred
	R0 int64 `form:"r0" binding:"Required"`
	R1 int64 `form:"r1" binding:"Required"`
	R2 int64 `form:"r2" binding:"Required"`
	R3 int64 `form:"r3" binding:"Required"`
	R4 int64 `form:"r4" binding:"Required"`
	R5 int64 `form:"r5" binding:"Required"`
}

type FormCityArmyCreate struct {
	CharacterID string `form:"cid" binding:"Required"`
	CityID      uint64 `form:"lid" binding:"Required"`
	Name        string `form:"name" binding:"Required"`
}

func doCityStudy(f *frontService) macaron.Handler {
	return func(ctx *macaron.Context, flash *session.Flash, sess session.Store, info FormCityStudy) {
		_, _, err := f.authenticateCharacterFromSession(ctx, sess, info.RegionID, info.CharacterID)
		if err != nil {
			flash.Warning(err.Error())
			ctx.Redirect("/game/user")
			return
		}

		cliReg := region.NewCityClient(f.cnxRegion)
		_, err = cliReg.Study(contextMacaronToGrpc(ctx, sess),
			&region.StudyReq{
				City: &region.CityId{
					Region:    info.RegionID,
					City:      info.CityID,
					Character: info.CharacterID,
				},
				KnowledgeType: info.KnowledgeID})
		if err != nil {
			flash.Warning(err.Error())
		}

		ctx.Redirect("/game/land/knowledges?cid=" + info.CharacterID + "&lid=" + utoa(info.CityID))
	}
}

func doCityBuild(f *frontService) macaron.Handler {
	return func(ctx *macaron.Context, flash *session.Flash, sess session.Store, info FormCityBuild) {
		_, _, err := f.authenticateCharacterFromSession(ctx, sess, info.RegionID, info.CharacterID)
		if err != nil {
			flash.Warning(err.Error())
			ctx.Redirect("/game/user")
			return
		}

		cliReg := region.NewCityClient(f.cnxRegion)
		_, err = cliReg.Build(contextMacaronToGrpc(ctx, sess),
			&region.BuildReq{
				City: &region.CityId{
					Region:    info.RegionID,
					City:      info.CityID,
					Character: info.CharacterID},
				BuildingType: info.BuildingID})
		if err != nil {
			flash.Warning(err.Error())
		}

		ctx.Redirect("/game/land/buildings?cid=" + info.CharacterID + "&lid=" + utoa(info.CityID))
	}
}

func doCityTrain(f *frontService) macaron.Handler {
	return func(ctx *macaron.Context, flash *session.Flash, sess session.Store, info FormCityTrain) {
		_, _, err := f.authenticateCharacterFromSession(ctx, sess, info.RegionID, info.CharacterID)
		if err != nil {
			flash.Warning(err.Error())
			ctx.Redirect("/game/user")
			return
		}

		cliReg := region.NewCityClient(f.cnxRegion)
		_, err = cliReg.Train(contextMacaronToGrpc(ctx, sess),
			&region.TrainReq{
				City: &region.CityId{
					Region:    info.RegionID,
					City:      info.CityID,
					Character: info.CharacterID},
				UnitType: info.UnitID})
		if err != nil {
			flash.Warning(err.Error())
		}

		ctx.Redirect("/game/land/units?cid=" + info.CharacterID + "&lid=" + utoa(info.CityID))
	}
}

func doCityArmyCreate(f *frontService) macaron.Handler {
	return func(ctx *macaron.Context, flash *session.Flash, sess session.Store, info FormCityArmyCreate) {
		ctx.Redirect("/game/land/overview?cid=" + info.CharacterID + "&lid=" + utoa(info.CityID))
	}
}

func doCityTransferUnit(f *frontService) macaron.Handler {
	return func(ctx *macaron.Context, flash *session.Flash, sess session.Store, info FormCityUnitTransfer) {
		ctx.Redirect("/game/land/overview?cid=" + info.CharacterID + "&lid=" + utoa(info.CityID))
	}
}

func doCityTransferResources(f *frontService) macaron.Handler {
	return func(ctx *macaron.Context, flash *session.Flash, sess session.Store, info FormCityUnitTransfer) {
		ctx.Redirect("/game/land/overview?cid=" + info.CharacterID + "&lid=" + utoa(info.CityID))
	}
}
