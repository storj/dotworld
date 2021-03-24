// Copyright (C) 2020 Storj Labs, Inc.
// See LICENSE for copying information.

package dotworld

import (
	"fmt"
	"image"
	"io"
)

// Map defines a location map with dots.
type Map struct {
	CountX, CountY int
	Grid           *Grid
	Bounds         image.Rectangle
	Locations      map[GridPosition]*Location
}

// Copy returns a new Map with the same data
func (m *Map) Copy() *Map {
	copy := &Map{
		CountX:    m.CountX,
		CountY:    m.CountY,
		Grid:      m.Grid,
		Bounds:    m.Bounds,
		Locations: make(map[GridPosition]*Location, len(m.Locations)),
	}
	for k, v := range m.Locations {
		loc := *v
		copy.Locations[k] = &loc
	}
	return copy
}

// GridPosition represents the integer row and column of a location on a
// grid.
type GridPosition struct {
	Row int
	Col int
}

// Location defines a location in LatLong coordinates.
type Location struct {
	S2
	Land float32 // 0 to 1
	Load float32
}

// EncodeSVG encodes map as a svg.
func (m *Map) EncodeSVG(w io.Writer, width, height int) (err error) {
	writef := func(s string, args ...interface{}) {
		if err != nil {
			return
		}
		_, err = fmt.Fprintf(w, s, args...)
	}

	writef(`<?xml version="1.0" encoding="UTF-8" standalone="no"?>`)
	writef(`<svg width="%[1]v" height="%[2]v" viewBox="0 0 %[1]v %[2]v" xmlns="http://www.w3.org/2000/svg">`, width, height)
	defer writef(`</svg>`)

	writef(`<style>circle{fill:#D8DDE1} .sn{fill:#2582FF}</style>`)

	pc := PlateCarre{
		Width:  float32(width),
		Height: float32(height),
	}

	locr := 0.5 * minF32(
		float32(width)/float32(m.CountX),
		float32(height)/float32(m.CountY),
	)
	locr -= 1.0

	plot := func(loc *Location) {
		p := pc.Forward(loc.S2)

		szf := loc.Land
		if loc.Load > 0 {
			szf = 1.1 + 2*loc.Load
		}

		sz := locr * szf
		if sz < 2 {
			sz = 2
		}

		if loc.Load > 0 {
			writef(`<circle class="sn" cx="%.0f" cy="%.0f" r="%.1f"/>`, p.X, p.Y, sz)
		} else {
			writef(`<circle cx="%.0f" cy="%.0f" r="%.1f"/>`, p.X, p.Y, sz)
		}
	}

	for _, loc := range m.Locations {
		if loc.Load == 0 {
			plot(loc)
		}
	}
	for _, loc := range m.Locations {
		if loc.Load > 0 {
			plot(loc)
		}
	}

	return
}

// minF32 calculates min of floats.
func minF32(a, b float32) float32 {
	if a < b {
		return a
	}
	return b
}

// Lookup finds the appropriate row and column in the grid for the given S2
// location.
func (m *Map) Lookup(pos S2) GridPosition {
	p2 := m.Grid.Coord.Forward(pos)
	tileSize := m.Bounds.Dx() / m.Grid.CountX
	if p2.X > float32(m.Bounds.Max.X) {
		p2.X = float32(m.Bounds.Max.X)
	}
	if p2.Y > float32(m.Bounds.Max.Y) {
		p2.Y = float32(m.Bounds.Max.Y)
	}
	p2.X -= float32(m.Bounds.Min.X)
	p2.Y -= float32(m.Bounds.Min.Y)
	return GridPosition{
		Row: int(p2.Y / float32(tileSize)),
		Col: int(p2.X / float32(tileSize)),
	}
}

// Nearest returns the requested grid position, or a close by neighbor if
// the requested one isn't found. Returns nil if nothing is found.
func (m *Map) Nearest(pos GridPosition) *Location {
	for _, attempt := range []struct {
		RowDelta int
		ColDelta int
	}{
		{0, 0},
		{0, 1}, // do columns first, because there is more column fidelity
		{0, -1},
		{1, 0},
		{-1, 0},
		{1, 1},
		{1, -1},
		{-1, 1},
		{-1, -1},
	} {
		if loc, exists := m.Locations[GridPosition{
			Row: pos.Row + attempt.RowDelta,
			Col: pos.Col + attempt.RowDelta}]; exists {
			return loc
		}
	}
	return nil
}
