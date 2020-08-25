// Copyright (C) 2018-2020 Hegemonie's AUTHORS
// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at http://mozilla.org/MPL/2.0/.

package mapgraph

import (
	"github.com/jfsmig/hegemonie/pkg/utils"
)

func (m *Map) Sections(p string) utils.PersistencyMapping {
	if p == "" {
		panic("Invalid path")
	}
	return []utils.CfgSection{
		{p, &m},
	}
}
