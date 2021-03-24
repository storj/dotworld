// Copyright (C) 2021 Storj Labs, Inc.
// See LICENSE for copying information.

package reference

//go:generate go run gen.go

import (
	"storj.io/dotworld"
)

// WorldMap is a good default world map to use.
func WorldMap() *dotworld.Map {
	return worldMap.Copy()
}
