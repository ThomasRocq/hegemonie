// Copyright (C) 2018-2020 Hegemonie's AUTHORS
// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at http://mozilla.org/MPL/2.0/.

package mapgraph

import (
	"errors"
	"os"
	"path/filepath"
	"strings"
	"sync"
)

var (
	errNoSuchMap = errors.New("No such Map")
)

type memRepository struct {
	maps setOfMaps
	rw   sync.RWMutex
}

func NewRepository() Repository {
	return &memRepository{}
}

func LoadDirectory(repo Repository, path string) error {
	return filepath.Walk(path, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		// Only accept non-hidden JSON files
		_, fn := filepath.Split(path)
		if info.IsDir() || info.Size() <= 0 {
			return nil
		}
		if len(fn) < 2 || fn[0] == '.' {
			return nil
		}
		if !strings.HasSuffix(fn, ".json") {
			return nil
		}

		var m Map
		m.Init()
		// The filename without its extension is the name of the map
		m.ID = fn[:len(fn)-5]

		err = m.Sections(path).Load()
		if err != nil {
			return err
		}

		err = m.PostLoad()
		if err != nil {
			return err
		}

		err = m.Check()
		if err != nil {
			return err
		}

		repo.Register(&m)
		return nil
	})
}

func SaveInDirectory(repo Repository, dataDir string) error {
	for m := range ListAllMaps(repo) {
		path := dataDir + "/" + m.ID + ".json"
		err := m.Sections(path).Dump()
		if err != nil {
			return err
		}
	}
	return nil
}

func ListAllMaps(repo Repository) <-chan *Map {
	out := make(chan *Map)
	go func(r Repository) {
		defer close(out)
		var last string
		for {
			tab, err := r.ListMaps(last, 10)
			if err != nil {
				return
			}
			if len(tab) <= 0 {
				return
			}
			for _, m := range tab {
				out <- m
			}
		}
	}(repo)
	return out
}

func (r *memRepository) GetMap(ID string) (*Map, error) {
	m := r.maps.Get(ID)
	if m == nil {
		return nil, errNoSuchMap
	}
	return m, nil
}

func (r *memRepository) ListMaps(marker string, max uint32) ([]*Map, error) {
	return r.maps.Slice(marker, max), nil
}

func (r *memRepository) Register(m *Map) {
	r.maps.Add(m)
}

func (r *memRepository) RLock()   { r.rw.RLock() }
func (r *memRepository) RUnlock() { r.rw.RUnlock() }

func (r *memRepository) WLock()   { r.rw.Lock() }
func (r *memRepository) WUnlock() { r.rw.Unlock() }

//go:generate go run github.com/jfsmig/hegemonie/cmd/gen-set ./map_auto.go mapgraph:setOfMaps:*Map ID:string
