// Copyright (C) 2018-2020 Hegemonie's AUTHORS
// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at http://mozilla.org/MPL/2.0/.

package region

import (
	"errors"
	"sort"

	"github.com/jfsmig/hegemonie/pkg/utils"
)

func (w *World) SetNotifier(n Notifier) {
	w.notifier = LogEvent(n)
}

func (w *World) Init() {
	w.WLock()
	defer w.WUnlock()

	w.SetNotifier(&noEvt{})
	w.Regions = make(SetOfRegions, 0)
	w.Definitions.Units = make(SetOfUnitTypes, 0)
	w.Definitions.Buildings = make(SetOfBuildingTypes, 0)
	w.Definitions.Knowledges = make(SetOfKnowledgeTypes, 0)
}

func (w *World) Check() error {
	w.RLock()
	defer w.RUnlock()
	if w.notifier == nil || w.mapView == nil {
		return errInvalidState
	}
	if err := w.Definitions.Check(); err != nil {
		return err
	}
	for _, r := range w.Regions {
		if err := r.Check(); err != nil {
			return err
		}
	}
	return nil
}

func (d *DefinitionsBase) Check() error {
	if !sort.IsSorted(&d.Knowledges) {
		return errors.New("knowledge types unsorted")
	}
	if !sort.IsSorted(&d.Buildings) {
		return errors.New("building types unsorted")
	}
	if !sort.IsSorted(&d.Units) {
		return errors.New("unit types unsorted")
	}

	return nil
}

func (w *Region) Check() error {
	if !sort.IsSorted(&w.Cities) {
		return errors.New("cities unsorted")
	}
	if !sort.IsSorted(&w.Fights) {
		return errors.New("fights unsorted")
	}

	for _, a := range w.Fights {
		if !sort.IsSorted(&a.Attack) {
			return errors.New("fight attack unsorted")
		}
		if !sort.IsSorted(&a.Defense) {
			return errors.New("fight defense unsorted")
		}
	}
	for _, a := range w.Cities {
		if !sort.IsSorted(&a.Knowledges) {
			return errors.New("knowledge unsorted")
		}
		if !sort.IsSorted(&a.Buildings) {
			return errors.New("building unsorted")
		}
		if !sort.IsSorted(&a.Units) {
			return errors.New("unit sequence: unsorted")
		}
		if !sort.IsSorted(&a.lieges) {
			return errors.New("city lieges unsorted")
		}
		if !sort.IsSorted(&a.Armies) {
			return errors.New("city armies unsorted")
		}
		for _, a := range a.Armies {
			if !sort.IsSorted(&a.Units) {
				return errors.New("units unsorted")
			}
		}
	}

	return nil
}

func (d *DefinitionsBase) PostLoad() error {
	sort.Sort(&d.Knowledges)
	sort.Sort(&d.Buildings)
	sort.Sort(&d.Units)
	return nil
}

func (r *Region) PostLoad() error {
	// Sort all the lookup arrays
	sort.Sort(&r.Cities)
	sort.Sort(&r.Fights)

	for _, c := range r.Cities {
		sort.Sort(&c.Knowledges)
		sort.Sort(&c.Buildings)
		sort.Sort(&c.Units)
		if c.Armies == nil {
			c.Armies = make(SetOfArmies, 0)
		} else {
			sort.Sort(&c.Armies)
		}
		if c.lieges == nil {
			c.lieges = make(SetOfCities, 0)
		} else {
			sort.Sort(&c.lieges)
		}

		for _, a := range c.Armies {
			// Link Armies to their City
			a.City = c
			// FIXME: Link each Army to its Fight
		}
	}

	// Compute the highest unique ID
	maxID := r.nextID
	for _, c := range r.Cities {
		if c.ID > maxID {
			maxID = c.ID
		}
		for _, a := range c.Armies {
			if a.ID > maxID {
				maxID = a.ID
			}
			for _, u := range a.Units {
				if u.ID > maxID {
					maxID = u.ID
				}
			}
		}
		for _, u := range c.Units {
			if u.ID > maxID {
				maxID = u.ID
			}
		}
		for _, u := range c.Knowledges {
			if u.ID > maxID {
				maxID = u.ID
			}
		}
		for _, u := range c.Buildings {
			if u.ID > maxID {
				maxID = u.ID
			}
		}
	}

	r.nextID = maxID + 1
	return nil
}

func (w *World) PostLoad() error {
	w.Definitions.PostLoad()
	sort.Sort(&w.Regions)
	for _, r := range w.Regions {
		r.PostLoad()
	}
	return nil
}

func (d *DefinitionsBase) Sections(p string) utils.PersistencyMapping {
	if p == "" {
		panic("Invalid path")
	}
	return []utils.CfgSection{
		{p + "/units.json", &d.Units},
		{p + "/buildings.json", &d.Buildings},
		{p + "/knowledge.json", &d.Knowledges},
	}
}

func (r *Region) Sections(p string) utils.PersistencyMapping {
	if p == "" {
		panic("Invalid path")
	}
	return []utils.CfgSection{
		{p + "/cities.json", &r.Cities},
		{p + "/fights.json", &r.Fights},
	}
}

func (w *World) Sections(p string) utils.PersistencyMapping {
	if p == "" {
		panic("Invalid path")
	}
	sections := []utils.CfgSection{
		{p + "/config.json", &w.Config},
	}
	sections = append(sections, w.Definitions.Sections(p+"/_defs")...)
	for _, r := range w.Regions {
		sections = append(sections, r.Sections(p+"/"+r.Name)...)
	}
	return sections
}
