// Copyright (C) 2018-2020 Hegemonie's AUTHORS
// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at http://mozilla.org/MPL/2.0/.

package region

import (
	"errors"
	"os"
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

	w.nextID = 1
	w.Live.Cities = make(SetOfCities, 0)
	w.Live.Fights = make(SetOfFights, 0)
	w.Definitions.Units = make(SetOfUnitTypes, 0)
	w.Definitions.Buildings = make(SetOfBuildingTypes, 0)
	w.Definitions.Knowledges = make(SetOfKnowledgeTypes, 0)
}

func (w *World) Check() error {
	if !sort.IsSorted(&w.Definitions.Knowledges) {
		return errors.New("knowledge types unsorted")
	}
	if !sort.IsSorted(&w.Definitions.Buildings) {
		return errors.New("building types unsorted")
	}
	if !sort.IsSorted(&w.Definitions.Units) {
		return errors.New("unit types unsorted")
	}

	if !sort.IsSorted(&w.Live.Cities) {
		return errors.New("cities unsorted")
	}
	if !sort.IsSorted(&w.Live.Fights) {
		return errors.New("fights unsorted")
	}

	for _, a := range w.Live.Fights {
		if !sort.IsSorted(&a.Attack) {
			return errors.New("fight attack unsorted")
		}
		if !sort.IsSorted(&a.Defense) {
			return errors.New("fight defense unsorted")
		}
	}
	for _, a := range w.Live.Cities {
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

func (w *World) PostLoad() error {
	// Sort all the lookup arrays
	sort.Sort(&w.Definitions.Knowledges)
	sort.Sort(&w.Definitions.Buildings)
	sort.Sort(&w.Definitions.Units)
	sort.Sort(&w.Live.Cities)
	sort.Sort(&w.Live.Fights)
	for _, c := range w.Live.Cities {
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
	maxID := w.nextID
	for _, u := range w.Definitions.Units {
		if u.ID > maxID {
			maxID = u.ID
		}
	}
	for _, u := range w.Definitions.Buildings {
		if u.ID > maxID {
			maxID = u.ID
		}
	}
	for _, u := range w.Definitions.Knowledges {
		if u.ID > maxID {
			maxID = u.ID
		}
	}
	for _, c := range w.Live.Cities {
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

	w.nextID = maxID + 1
	return nil
}

func liveSections(p string, w *World) utils.PersistencyMapping {
	if p == "" {
		panic("Invalid path")
	}
	return []utils.CfgSection{
		{p + "/cities.json", &w.Live.Cities},
		{p + "/fights.json", &w.Live.Fights},
	}
}

func defsSections(p string, w *World) utils.PersistencyMapping {
	if p == "" {
		panic("Invalid path")
	}
	return []utils.CfgSection{
		{p + "/config.json", &w.Config},
		{p + "/units.json", &w.Definitions.Units},
		{p + "/buildings.json", &w.Definitions.Buildings},
		{p + "/knowledge.json", &w.Definitions.Knowledges},
	}
}

// Save the current state of the World as a set of JSON objects, each for a dataset
// that participates to represent the World. Only the "LIVE" entities are concerned:
// - the map in use in the region
// - the cities spwaned on the map
// - the fights currently running on the map.
// Counter-part of LoadLiveFromFiles()
func (w *World) SaveLiveToFiles(basePath string) error {
	err := os.MkdirAll(basePath, 0755)
	if err == nil {
		err = liveSections(basePath, w).Dump()
	}
	return err
}

// Save the current state of the World as a set of JSON objects, each for a dataset
// that participates to represent the World. Only the "DEFINITIONS" entities are concerned:
// - the knowledge tree
// - the building definitions
// - the troops definitions
// - the general configuration of the World.
// Counter-part of LoadDefinitionsFromFiles()
func (w *World) SaveDefinitionsToFiles(basePath string) error {
	err := os.MkdirAll(basePath, 0755)
	if err == nil {
		err = defsSections(basePath, w).Dump()
	}
	return err
}

// Restore a state for the World, from a set of JSON objects, where each file/object
// participates to represent of dataset of the World. Only the "LIVE" entities
// are concerned:
// - the map in use in the region
// - the cities spwaned on the map
// - the fights currently running on the map.
// Counter-part of SaveLiveToFiles()
func (w *World) LoadLiveFromFiles(basePath string) error {
	return liveSections(basePath, w).Load()
}

// Restore a state for the World, from a set of JSON objects, where each file/object
// participates to represent of dataset of the World. Only the "DEFINITIONS" entities
// are concerned:
// - the knowledge tree
// - the building definitions
// - the troops definitions
// - the general configuration of the World.
// Counter-part of SaveDefinitionsToFiles()
func (w *World) LoadDefinitionsFromFiles(basePath string) error {
	return defsSections(basePath, w).Load()
}
