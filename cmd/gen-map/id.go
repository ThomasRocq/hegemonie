// Copyright (C) 2020 Hegemonie's AUTHORS
// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at http://mozilla.org/MPL/2.0/.

package main

import (
	"sync/atomic"
)

var nextId uint64

func getNextId() uint64 {
	return atomic.AddUint64(&nextId, 1)
}
