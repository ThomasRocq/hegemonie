// Copyright (C) 2018-2020 Hegemonie's AUTHORS
// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at http://mozilla.org/MPL/2.0/.

package region

import (
	"errors"
	"fmt"
)

var (
	ErrNoSuchUnit         = errors.New("No such Unit")
	ErrNotEnoughResources = errors.New("Not enough resources")
)

func MakeCity() *City {
	return &City{
		ID:         0,
		Units:      make(SetOfUnits, 0),
		Buildings:  make(SetOfBuildings, 0),
		Knowledges: make(SetOfKnowledges, 0),
		Armies:     make(SetOfArmies, 0),
		lieges:     make(SetOfCities, 0),
	}
}

func CopyCity(original *City) *City {
	c := MakeCity()
	if original != nil {
		c.Stock.Set(original.Stock)
		c.Production.Set(original.Production)
		c.StockCapacity.Set(original.StockCapacity)
	}
	return c
}

// Return a Unit owned by the current City, given the Unit ID
func (c *City) Unit(id uint64) *Unit {
	return c.Units.Get(id)
}

// Return a Building owned by the current City, given the Building ID
func (c *City) Building(id uint64) *Building {
	return c.Buildings.Get(id)
}

// Return a Knowledge owned by the current City, given the Knowledge ID
func (c *City) Knowledge(id uint64) *Knowledge {
	return c.Knowledges.Get(id)
}

// Return total Popularity of the current City (permanent + transient)
func (c *City) GetActualPopularity(w *World) int64 {
	var pop int64 = c.PermanentPopularity

	// Add Transient values for Units in the Armies
	for _, a := range c.Armies {
		for _, u := range a.Units {
			ut := w.UnitTypeGet(u.Type)
			pop += ut.PopBonus
		}
		pop += w.Config.PopBonusArmyAlive
	}

	// Add Transient values for Units in the City
	for _, u := range c.Units {
		ut := w.UnitTypeGet(u.Type)
		pop += ut.PopBonus
	}

	// Add Transient values for Buildings
	for _, b := range c.Buildings {
		bt := w.BuildingTypeGet(b.Type)
		pop += bt.PopBonus
	}

	// Add Transient values for Knowledges
	for _, k := range c.Knowledges {
		kt := w.KnowledgeTypeGet(k.Type)
		pop += kt.PopBonus
	}

	return pop
}

func (c *City) GetProduction(w *World) *CityProduction {
	p := &CityProduction{
		Buildings: ResourceModifierNoop(),
		Knowledge: ResourceModifierNoop(),
	}

	for _, b := range c.Buildings {
		t := w.BuildingTypeGet(b.Type)
		p.Buildings.ComposeWith(t.Prod)
	}
	for _, u := range c.Knowledges {
		t := w.KnowledgeTypeGet(u.Type)
		p.Knowledge.ComposeWith(t.Prod)
	}

	p.Base = c.Production
	p.Actual = c.Production
	p.Actual.Apply(p.Buildings)
	p.Actual.Apply(p.Knowledge)
	return p
}

func (c *City) GetStock(w *World) *CityStock {
	p := &CityStock{
		Buildings: ResourceModifierNoop(),
		Knowledge: ResourceModifierNoop(),
	}

	for _, b := range c.Buildings {
		t := w.BuildingTypeGet(b.Type)
		p.Buildings.ComposeWith(t.Stock)
	}
	for _, b := range c.Knowledges {
		t := w.BuildingTypeGet(b.Type)
		p.Buildings.ComposeWith(t.Stock)
	}

	p.Base = c.StockCapacity
	p.Actual = c.StockCapacity
	p.Actual.Apply(p.Buildings)
	p.Actual.Apply(p.Knowledge)
	p.Usage = c.Stock
	return p
}

func (c *City) CreateEmptyArmy(w *World) *Army {
	aid := w.getNextID()
	a := &Army{
		ID:       aid,
		City:     c,
		Cell:     c.ID,
		Name:     fmt.Sprintf("A-%d", aid),
		Units:    make(SetOfUnits, 0),
		Postures: []int64{int64(c.ID)},
		Targets:  make([]Command, 0),
	}
	c.Armies.Add(a)
	return a
}

func unitsToIDs(uv []*Unit) (out []uint64) {
	for _, u := range uv {
		out = append(out, u.ID)
	}
	return out
}

func unitsFilterIdle(uv []*Unit) (out []*Unit) {
	for _, u := range uv {
		if u.Health > 0 && u.Ticks <= 0 {
			out = append(out, u)
		}
	}
	return out
}

// Create an Army made of some Unit of the City
func (c *City) CreateArmyFromUnit(w *World, units ...*Unit) (*Army, error) {
	return c.CreateArmyFromIds(w, unitsToIDs(unitsFilterIdle(units))...)
}

// Create an Army made of some Unit of the City
func (c *City) CreateArmyFromIds(w *World, ids ...uint64) (*Army, error) {
	a := c.CreateEmptyArmy(w)
	err := c.TransferOwnUnit(a, ids...)
	if err != nil { // Rollback
		a.Disband(w, c, false)
		return nil, err
	}
	return a, nil
}

// Create an Army made of all the Units defending the City
func (c *City) CreateArmyDefence(w *World) (*Army, error) {
	ids := unitsToIDs(unitsFilterIdle(c.Units))
	if len(ids) <= 0 {
		return nil, ErrNoSuchUnit
	}
	return c.CreateArmyFromIds(w, ids...)
}

// Create an Army carrying resources you own
func (c *City) CreateTransport(w *World, r Resources) (*Army, error) {
	if !c.Stock.GreaterOrEqualTo(r) {
		return nil, ErrNotEnoughResources
	}

	a := c.CreateEmptyArmy(w)
	c.Stock.Remove(r)
	a.Stock.Add(r)
	return a, nil
}

// Play one round of local production and return the
func (c *City) ProduceLocally(w *World, p *CityProduction) Resources {
	var prod Resources = p.Actual
	if c.TicksMassacres > 0 {
		mult := MultiplierUniform(w.Config.MassacreImpact)
		for i := uint32(0); i < c.TicksMassacres; i++ {
			prod.Multiply(mult)
		}
		c.TicksMassacres--
	}
	return prod
}

func (c *City) Produce(w *World) {
	// Pre-compute the modified values of Stock and Production.
	// We just reuse a functon that already does it (despite it does more)
	prod0 := c.GetProduction(w)
	stock := c.GetStock(w)

	// Make the local City generate resources (and recover the massacres)
	prod := c.ProduceLocally(w, prod0)
	c.Stock.Add(prod)

	if c.Overlord != 0 {
		if c.pOverlord != nil {
			// Compute the expected Tax based on the local production
			var tax Resources = prod
			tax.Multiply(c.TaxRate)
			// Ensure the tax isn't superior to the actual production (to cope with
			// invalid tax rates)
			tax.TrimTo(c.Stock)
			// Then preempt the tax from the stock
			c.Stock.Remove(tax)

			// TODO(jfs): check for potential shortage
			//  shortage := c.Tax.GreaterThan(tax)

			if w.Config.InstantTransfers {
				c.pOverlord.Stock.Add(tax)
			} else {
				c.SendResourcesTo(w, c.pOverlord, tax)
			}

			// FIXME(jfs): notify overlord
			// FIXME(jfs): notify c
		}
	}

	// ATM the stock maybe still stores resources. We use them to make the assets evolve.
	// We arbitrarily give the preference to Units, then Buildings and eventually the
	// Knowledge.

	for _, u := range c.Units {
		if u.Ticks > 0 {
			ut := w.UnitTypeGet(u.Type)
			if c.Stock.GreaterOrEqualTo(ut.Cost) {
				c.Stock.Remove(ut.Cost)
				u.Ticks--
				if u.Ticks <= 0 {
					// FIXME(jfs): Notify the City
				}
			}
		}
	}

	for _, b := range c.Buildings {
		if b.Ticks > 0 {
			bt := w.BuildingTypeGet(b.ID)
			if c.Stock.GreaterOrEqualTo(bt.Cost) {
				c.Stock.Remove(bt.Cost)
				b.Ticks--
				if b.Ticks <= 0 {
					// FIXME(jfs): Notify the City
				}
			}
		}
	}

	for _, k := range c.Knowledges {
		if k.Ticks > 0 {
			bt := w.KnowledgeTypeGet(k.ID)
			if c.Stock.GreaterOrEqualTo(bt.Cost) {
				c.Stock.Remove(bt.Cost)
				k.Ticks--
			}
			if k.Ticks <= 0 {
				// FIXME(jfs): Notify the City
			}
		}
	}

	// At the end of the turn, ensure we do not hold more resources than the actual
	// stock capacity (with the effect of all the multipliers)
	c.Stock.TrimTo(stock.Actual)
}

// Set a tax rate on the current City, with the same ratio on every Resource.
func (c *City) SetUniformTaxRate(nb float64) {
	c.TaxRate = MultiplierUniform(nb)
}

// Set the given tax rate to the current City.
func (c *City) SetTaxRate(m ResourcesMultiplier) {
	c.TaxRate = m
}

func (c *City) LiberateCity(w *World, other *City) {
	pre := other.pOverlord
	if pre == nil {
		return
	}

	other.Overlord = 0
	other.pOverlord = nil

	// FIXME(jfs): Notify 'pre'
	// FIXME(jfs): Notify 'c'
	// FIXME(jfs): Notify 'other'
}

func (c *City) GainFreedom(w *World) {
	pre := c.pOverlord
	if pre == nil {
		return
	}

	c.Overlord = 0
	c.pOverlord = nil

	// FIXME(jfs): Notify 'pre'
	// FIXME(jfs): Notify 'c'
}

func (c *City) ConquerCity(w *World, other *City) {
	if other.pOverlord == c {
		c.pOverlord = nil
		c.Overlord = 0
		c.TaxRate = MultiplierUniform(0)
		return
	}

	//pre := other.pOverlord
	other.pOverlord = c
	other.Overlord = c.ID
	other.TaxRate = MultiplierUniform(w.Config.RateOverlord)

	// FIXME(jfs): Notify 'pre'
	// FIXME(jfs): Notify 'c'
	// FIXME(jfs): Notify 'other'
}

func (c *City) SendResourcesTo(w *World, overlord *City, amount Resources) error {
	// FIXME(jfs): NYI
	return errors.New("SendResourcesTo() not implemented")
}

func (c *City) TransferOwnResources(a *Army, r Resources) error {
	if a.City != c {
		return errors.New("Army not controlled by the City")
	}
	if !c.Stock.GreaterOrEqualTo(r) {
		return errors.New("Insufficient resources")
	}

	c.Stock.Remove(r)
	a.Stock.Add(r)
	return nil
}

func (c *City) TransferOwnUnit(a *Army, units ...uint64) error {
	if len(units) <= 0 || a == nil {
		panic("EINVAL")
	}

	if a.City != c {
		return errors.New("Army not controlled by the City")
	}

	allUnits := make(map[uint64]*Unit)
	for _, uid := range units {
		if _, ok := allUnits[uid]; ok {
			continue
		}
		if u := c.Units.Get(uid); u == nil {
			return errors.New("Unit not found")
		} else if u.Ticks > 0 || u.Health <= 0 {
			continue
		} else {
			allUnits[uid] = u
		}
	}

	for _, u := range allUnits {
		c.Units.Remove(u)
		a.Units.Add(u)
	}
	return nil
}

func (c *City) KnowledgeFrontier(w *World) []*KnowledgeType {
	return w.KnowledgeGetFrontier(c.Knowledges)
}

func (c *City) BuildingFrontier(w *World) []*BuildingType {
	return w.BuildingGetFrontier(c.GetActualPopularity(w), c.Buildings, c.Knowledges)
}

// Return a collection of UnitType that may be trained by the current City
// because all the requirements are met.
// Each UnitType 'p' returned validates 'c.UnitAllowed(p)'.
func (c *City) UnitFrontier(w *World) []*UnitType {
	return w.UnitGetFrontier(c.Buildings)
}

// Check the current City has all the requirements to train a Unti of the
// given UnitType.
func (c *City) UnitAllowed(pType *UnitType) bool {
	if pType.RequiredBuilding == 0 {
		return true
	}
	for _, b := range c.Buildings {
		if b.Type == pType.RequiredBuilding {
			return true
		}
	}
	return false
}

// Create a Unit of the given UnitType.
// No check is performed to verify the City has all the requirements.
func (c *City) UnitCreate(w *World, pType *UnitType) *Unit {
	id := w.getNextID()
	u := &Unit{ID: id, Type: pType.ID, Ticks: pType.Ticks, Health: pType.Health}
	c.Units.Add(u)
	return u
}

// Start the training of a Unit of the given UnitType (id).
// The whole chain of requirements will be checked.
func (c *City) Train(w *World, typeID uint64) (uint64, error) {
	pType := w.UnitTypeGet(typeID)
	if pType == nil {
		return 0, errors.New("Unit Type not found")
	}
	if !c.UnitAllowed(pType) {
		return 0, errors.New("Precondition Failed: no suitable building")
	}

	u := c.UnitCreate(w, pType)
	return u.ID, nil
}

func (c *City) Study(w *World, typeID uint64) (uint64, error) {
	kType := w.KnowledgeTypeGet(typeID)
	if kType == nil {
		return 0, errors.New("Knowledge Type not found")
	}
	for _, k := range c.Knowledges {
		if typeID == k.Type {
			return 0, errors.New("Already started")
		}
	}
	if !CheckKnowledgeDependencies(c.Knowledges, kType.Requires, kType.Conflicts) {
		return 0, errors.New("Conflict")
	}

	id := w.getNextID()
	c.Knowledges.Add(&Knowledge{ID: id, Type: typeID, Ticks: kType.Ticks})
	return id, nil
}

func (c *City) Build(w *World, bID uint64) (uint64, error) {
	bType := w.BuildingTypeGet(bID)
	if bType == nil {
		return 0, errors.New("Building Type not found")
	}
	if !bType.MultipleAllowed {
		for _, b := range c.Buildings {
			if b.Type == bID {
				return 0, errors.New("Building already present")
			}
		}
	}
	if !CheckKnowledgeDependencies(c.Knowledges, bType.Requires, bType.Conflicts) {
		return 0, errors.New("Conflict")
	}
	if !c.Stock.GreaterOrEqualTo(bType.Cost0) {
		return 0, errors.New("Not enough ressources")
	}

	id := w.getNextID()
	c.Buildings.Add(&Building{ID: id, Type: bID, Ticks: bType.Ticks})
	return id, nil
}

func (c *City) Lieges() []*City {
	return c.lieges[:]
}
