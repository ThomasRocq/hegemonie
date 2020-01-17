// Copyright (C) 2018-2020 Hegemonie's AUTHORS
// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at http://mozilla.org/MPL/2.0/.

package world

import (
	"sync"
)

const (
	ResourceMax = 6
)

const (
	// Do nothing. Useful for waypoints
	CmdPause = 0
	// Start a fight or join a running fight on the side of the attackers
	CmdCityAttack = 2
	// Join a running fight on the side of the defenders, or Watch the City if
	CmdCityDefend = 3
	// Attack the City and become its overlord in case of victory
	CmdCityOverlord = 4
	// Attack the City and break a building in case of victory
	CmdCityBreak = 5
	// Attack the City and reduce its production for the next turn
	CmdCityMassacre = 6
	// Deposit all the resources of the Army to the local City
	CmdCityDeposit = 7
	// Disband the Army and transfer its units and resources to the local City
	CmdCityDisband = 8
)

type Resources [ResourceMax]uint64

type ResourcesIncrement [ResourceMax]int64

type ResourcesMultiplier [ResourceMax]float64

type ResourceModifiers struct {
	Mult ResourcesMultiplier
	Plus ResourcesIncrement
}

type Command struct {
	// The unique ID of the Cell to target
	Cell uint64

	// What to do once arrived at the given Cell.
	Action uint
}

type Army struct {
	// The unique ID of the current Army
	Id uint64

	// The ID of the City that controls the current Army
	City uint64 `json:",omitempty"`

	// The ID of the Fight this Army is involved in.
	Fight uint64 `json:",omitempty"`

	// The ID of the Cell the Army is on
	Cell uint64 `json:",omitempty"`

	// Is the current Army still alive?
	// This flag is used to save array handlings.
	Deleted bool `json:",omitempty"`

	// A display name for the current City
	Name string

	// How many resources are carried by that Army
	Stock Resources

	Units SetOfUnits

	// The IS of a Cell of the Map that is a goal of the current movement of the Army
	Targets []Command `json:",omitempty"`

	// An array of Postures against armies of other cities.
	// A positive value means "defend"
	// A negative value means "assault"
	Postures []int64 `json:",omitempty"`
}

type KnowledgeType struct {
	Id    uint64
	Name  string
	Ticks uint `json:",omitempty"`

	// Transient bonus of Popularity, when the Knowledge is present
	PopBonus int64

	// Permanent bonus of Popularity when the Knowledge is achieved
	PopBonusLearn int64

	// Permanent bonus of Popularity (to the owner) when the Knowledge is stolen
	PopBonusStealVictim int64

	// Permanent bonus of Popularity (to the robber) when the Knowledge is stolen
	PopBonusStealActor int64

	Cost      Resources
	Requires  []uint64
	Conflicts []uint64
}

type Knowledge struct {
	Id    uint64
	Type  uint64
	Ticks uint `json:",omitempty"`
}

type BuildingType struct {
	// Unique ID of the BuildingType
	Id uint64

	// Display name of the current BuildingType
	Name string

	// How many ticks for the construction
	Ticks uint `json:",omitempty"`

	// How much does the production cost
	Cost Resources

	// Has the building to be unique a the City
	Unique bool `json:",omitempty"`

	// Amount of total popularity required to start the construction of the building
	PopRequired int64

	// Transient bonus of Popularity, when the Building is alive
	PopBonus int64

	// Permanent bonus of Popularity given when the Building is achieved
	PopBonusBuild int64

	// Permanent bonus of Popularity given to the owner of the Building when it is destroyed.
	PopBonusFall int64

	// Permanent bonus of Popularity given to the destroyer of the Building
	PopBonusDestroy int64

	// Permanent bonus of Popularity given to the owner of the Building when it is dismantled.
	PopBonusDismantle int64

	// Impat of the current Building on the total storage capacity of the City.
	Stock ResourceModifiers

	// Increment of resources produced by this building.
	Prod ResourceModifiers

	// A set of KnowledgeType ID that must all be present in a City to let that City start
	// this kind of building.
	Requires []uint64

	// A set of KnowledgeType ID that must all be absent in a City to let that City start
	// this kind of building.
	Conflicts []uint64
}

type Building struct {
	// The unique ID of the current Building
	Id uint64

	// The unique ID of the BuildingType associated to the current Building
	Type uint64

	// How many construction rounds remain before the building's achievement
	Ticks uint `json:",omitempty"`

	// Has the building been ruined?
	Deleted bool `json:",omitempty"`
}

type Character struct {
	// The unique identifier of the current Character
	Id uint64

	// The unique identifier of the only User that controls the Character.
	User uint64

	// The display name of the current Character
	Name string
}

type City struct {
	// The unique ID of the current City
	Id uint64

	// The unique ID of the main Character in charge of the City.
	// The Manager may name a Deputy manager in the City.
	Owner uint64

	// The unique ID of a second Character in charge of the City.
	Deputy uint64 `json:",omitempty"`

	// The unique ID of a City who is the boss of the current City.
	// Used for resources production computations.
	Overlord uint64

	// Ratio of the produced resources automatically sent to the Overlord City.
	TaxRate ResourcesMultiplier

	// The unique ID of the Cell the current City is built on.
	// This is redundant with the City field in the Cell structure.
	// Both information must match.
	Cell uint64

	Assault *Fight `json:",omitempty"`

	// The display name of the current City
	Name string

	// Permanent Popularity of the current City
	// The total value is the permanent value plus several "transient" bonus
	Pop int64

	Chaotic    uint
	Alignement uint
	Race       uint
	Religion   uint

	// Resources stock owned by the current City
	Stock Resources

	// Maximum amounts of each resources that might be stored in the town hall
	// of the city. That limit doesn't consider the modifiers.
	StockCapacity Resources

	// Resources produced each round by the City, before the enforcing of
	// Production Boosts ans Production Multipliers
	Production Resources

	// Number of massacres the current City undergo.
	// It takes one production turn to recover one Massacre.
	LastMassacres uint `json:",omitempty"`

	// Is the city still usable
	Deleted bool `json:",omitempty"`

	// Tells if the City is in automatic mode.
	// The "auto" mode is intented for inactive or absent players.
	// The armies come home to defend the City, no new building or unit is spawned.
	// In the plans: a conservative behavior should be automated
	Auto bool `json:",omitempty"`

	Knowledges SetOfKnowledges

	Buildings SetOfBuildings

	// Units directly defending the current City
	Units SetOfUnits

	// PRIVATE
	// Armies under the responsibility of the current City
	armies SetOfArmies

	// PRIVATE
	// Pointer to the current Overlord of the current City
	pOverlord *City

	// PRIVATE
	// Pointer to cities we currently are the overlord of
	lieges SetOfCities
}

type UnitType struct {
	// Unique Id of the Unit Type
	Id uint64

	// The number of health point for that type of unit.
	// A health equal to 0 means the death of the unit.
	Health uint

	// How affected is that type of unit by a loss of Health.
	// Must be between 0 and 1.
	// 0 means that the capacity of the Unit isn't affected by a health reduction.
	// 1 means that the capacity of the Unit loses an equal percentage of its capacity
	// for a loss of health (in other words, a HealthFactor of 1 means that the Unit
	// will hit at 90% of its maximal power if it has 90% of its health points).
	HealthFactor float64

	// The display name of the Unit Type
	Name string

	// How many ticks
	Ticks uint

	// Instantiation cost of the current UnitType
	Cost Resources

	// Transient bonus of Popularity, when the Unit is alive
	PopBonus int64

	// Permanent bonus of Popularity given when the Unit's training is done
	PopBonusTrain int64

	// Permanent bonus of Popularity given to the owner of the Unit when it dies
	PopBonusDeath int64

	// Permanent bonus of Popularity given to the killer of the Unit
	PopBonusKill int64

	// Permanent bonus of Popularity given to the ownerof the Unit when it is disbanded.
	PopBonusDisband int64

	// Might positive (resource boost) or more commonly negative (maintenance cost)
	Prod ResourceModifiers

	// Required Popularity to start trzaining this type of troop
	ReqPop int64

	// A UnitType is only dependant on the presence of a Building of that BuildingType.
	RequiredBuilding uint64
}

// Both Cell and City must not be 0, and have a non-0 value
type Unit struct {
	// Unique Id of the Unit
	Id uint64

	// A copy of the definition for the current Unit.
	Type uint64

	// How many ticks remain before the Troop training is finished
	Ticks uint

	// The number of health points of the unit, Health should be less or equal to HealthMax
	Health uint `json:"H,omitempty"`
}

type User struct {
	// The unique identifier of the current User
	Id uint64

	// The display name of the current User
	Name string

	// The unique email that authenticates the User.
	Email string

	// The hashed password that authenticates the User
	Password string

	// Has the current User the permission to manage the service.
	Admin bool `json:",omitempty"`

	// Can the user still login.
	Inactive bool `json:",omitempty"`
}

type Fight struct {
	// The unique ID of the
	Id uint64

	// The unique ID of the MapVertex the current Fight is happening on.
	Cell uint64

	// The set of Id of armies involved in the current Fight on the "attack" side
	// (the side that initiated the fight)
	Attack SetOfArmies

	/// The set of Id of armies involved in the current Fight on the "defence" side
	// the (side that has been pforce-pulled).
	Defense SetOfArmies
}

// A MapEdge is an edge if the transportation directed graph
type MapEdge struct {
	// Unique identifier of the source Cell
	S uint64

	// Unique identifier of the destination Cell
	D uint64

	// May the road be used by Units
	Deleted bool `json:",omitempty"`
}

// A MapVertex is a vertex in the transportation directed graph
type MapVertex struct {
	// The unique identifier of the current cell.
	Id uint64

	// Biome in which the cell is
	Biome uint64

	// The unique ID of the city present at this location.
	City uint64 `json:",omitempty"`
}

// A Map is a directed graph destined to be used as a transport network,
// organised as an adjacency list.
type Map struct {
	Cells  SetOfVertices
	Roads  SetOfEdges
	NextId uint64

	steps map[vector]uint64
}

type SetOfId []uint64

type SetOfArmies []*Army

type SetOfUnits []*Unit

type SetOfUnitTypes []*UnitType

type SetOfUsers []*User

type SetOfBuildings []*Building

type SetOfBuildingTypes []*BuildingType

type SetOfKnowledges []*Knowledge

type SetOfKnowledgeTypes []*KnowledgeType

type SetOfCharacters []*Character

type SetOfCities []*City

type SetOfFights []*Fight

type SetOfVertices []*MapVertex

type SetOfEdges []*MapEdge

type AuthBase struct {
	Users      SetOfUsers
	Characters SetOfCharacters
}

type DefinitionsBase struct {
	Units      SetOfUnitTypes
	Buildings  SetOfBuildingTypes
	Knowledges SetOfKnowledgeTypes

	// Ratio applied to the production of resources that is applied for each
	// Massacre underwent by any city. It only impacts the production of the City itself.
	MassacreImpact float64

	// Should resource transfers happen instantly or should an actual transport
	// be emitted by the sender? Set to `true` for an instant transfer or to
	// `false` for a transport.
	InstantTransfers bool

	// Permanent bonus to the Popularity when a City creates an Army
	PopBonusArmyCreate int64

	// Permanent bonus to the Popularity when a City disband an Army
	PopBonusArmyDisband int64

	// Transient bonus to the Popularity of a City for each of its live Army
	PopBonusArmyAlive int64

	// Default Overlord rate: percentage of the production of a City that is
	// taxed by its Overlord
	RateOverlord float64
}

type LiveBase struct {
	// Free armies on the map, not involved in any Fight
	Armies SetOfArmies

	// All the cities present on the Region
	Cities SetOfCities

	// Fights currently happening. The armies involved in the Fight are owned
	// By the Fight and do not appear in the "Armies" field.
	Fights SetOfFights
}

type World struct {
	Auth        AuthBase
	Definitions DefinitionsBase
	Live        LiveBase
	Places      Map

	NextId uint64
	Salt   string
	rw     sync.RWMutex
}
