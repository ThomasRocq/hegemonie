// Copyright (C) 2018-2020 Hegemonie's AUTHORS
// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at http://mozilla.org/MPL/2.0/.

package region

import (
	"sync/atomic"
)

func (r *Region) Produce() {
	for _, c := range r.Cities {
		c.Produce(r)
	}
}

func (r *Region) Move() {
	for _, c := range r.Cities {
		for _, a := range c.Armies {
			a.Move(r)
		}
	}
}

func (r *Region) CityGet(id uint64) *City {
	return r.Cities.Get(id)
}

func (r *Region) CityGetAt(loc uint64) *City {
	return r.CityGet(loc)
}

func (r *Region) CityCheck(id uint64) bool {
	return r.CityGet(id) != nil
}

func (r *Region) CityCreateModel(loc uint64, model *City) (*City, error) {
	if r.Cities.Has(loc) {
		return nil, errCityExists
	}
	city := CopyCity(model)
	city.ID = loc
	r.Cities.Add(city)
	return city, nil
}

func (r *Region) CityCreate(loc uint64) (*City, error) {
	return r.CityCreateModel(loc, nil)
}

func (r *Region) CityGetAndCheck(charID, cityID uint64) (*City, error) {
	// Fetch + sanity checks about the city
	pCity := r.CityGet(cityID)
	if pCity == nil {
		return nil, errCityNotFound
	}
	if pCity.Deputy != charID && pCity.Owner != charID {
		return nil, errForbidden
	}

	return pCity, nil
}

func (r *Region) CitiesList(idChar uint64) []*City {
	rep := make([]*City, 0)
	for _, c := range r.Cities {
		if c.Owner == idChar || c.Deputy == idChar {
			rep = append(rep, c)
		}
	}
	return rep
}

func (w *Region) getNextID() uint64 {
	return atomic.AddUint64(&w.nextID, 1)
}
