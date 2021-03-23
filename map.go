// Copyright (C) 2020 Storj Labs, Inc.
// See LICENSE for copying information.

package dotworld

import (
	"fmt"
	"io"
)

// Map defines a location map with dots.
type Map struct {
	CountX, CountY int
	Locations      []Location
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

	plot := func(loc Location) {
		p := pc.Forward(loc.S2)

		szf := loc.Land
		if loc.Load > 0 {
			szf = 1.5
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
