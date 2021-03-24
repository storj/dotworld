// Copyright (C) 2020 Storj Labs, Inc.
// See LICENSE for copying information.

package dotworld

import (
	"image"
)

// Grid defines arguments for deriving a Map from an image.
type Grid struct {
	// CountX defines number of dots in X dimension.
	CountX int
	// Thshold defines minimum average Y value, where dot should be placed.
	Threshold float32

	// Coord define image coordinate system.
	Coord PlateCarre
}

// MapFromImage derives a Map from an image.
func (grid *Grid) MapFromImage(m *image.Gray) *Map {
	b := m.Bounds()
	tileSize := b.Dx() / grid.CountX

	dotmap := Map{
		CountX:    grid.CountX,
		Grid:      grid,
		Bounds:    b,
		Locations: map[GridPosition]*Location{},
	}

	for tileY0 := b.Min.Y; tileY0 < b.Max.Y; tileY0 += tileSize {
		dotmap.CountY++
		tileY1 := minInt(tileY0+tileSize, b.Max.Y)
		tileCountX := 0
		for tileX0 := b.Min.X; tileX0 < b.Max.X; tileX0 += tileSize {
			tileCountX++
			tileX1 := minInt(tileX0+tileSize, b.Max.X)

			tileCenter := P2{
				X: float32(tileX0+tileX1) / 2,
				Y: float32(tileY0+tileY1) / 2,
			}
			center := grid.Coord.Reverse(tileCenter)

			tileRect := image.Rect(tileX0, tileY0, tileX1, tileY1)
			land := 1 - AvgY(m, tileRect)
			if land > grid.Threshold {

				gridPos := GridPosition{
					Row: dotmap.CountY - 1,
					Col: tileCountX - 1,
				}

				dotmap.Locations[gridPos] = &Location{
					S2:   center,
					Land: land,
				}

			}
		}
	}

	return &dotmap
}

// AvgY calculates average Y (luminosity) for a grayscale image.
func AvgY(m *image.Gray, b image.Rectangle) float32 {
	mb := m.Bounds()
	totalY := int64(0)
	for y := b.Min.Y; y < b.Max.Y; y++ {
		iy := (y-mb.Min.Y)*m.Stride - mb.Min.X
		for x := b.Min.X; x < b.Max.X; x++ {
			totalY += int64(m.Pix[iy+x])
		}
	}

	return float32(totalY) / float32(255*b.Dx()*b.Dy())
}

func minInt(x, y int) int {
	if x < y {
		return x
	}
	return y
}
