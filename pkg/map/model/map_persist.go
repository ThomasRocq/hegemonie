// Copyright (C) 2018-2020 Hegemonie's AUTHORS
// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at http://mozilla.org/MPL/2.0/.

package mapgraph

import (
	"github.com/jfsmig/hegemonie/pkg/utils"
	"os"
)

func mapSections(p string, m *Map) utils.PersistencyMapping {
	if p == "" {
		panic("Invalid path")
	}
	return []utils.CfgSection{
		{p + "/map.json", &m},
	}
}

func (m *Map) SaveToFiles(path string) error {
	err := os.MkdirAll(path, 0755)
	if err == nil {
		err = mapSections(path, m).Dump()
	}
	return err
}

func (m *Map) LoadFromFiles(path string) error {
	return mapSections(path, m).Load()
}
