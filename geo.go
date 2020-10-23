package main

// S2 defines a LatLong coordinate.
type S2 struct {
	Lat  float32
	Long float32
}

// P2 defines a pixel coordinate.
type P2 struct {
	X float32
	Y float32
}

// PlateCarree defines parameters for Plate Carree projection.
type PlateCarre struct {
	Width  float32
	Height float32
}

// Reference defines the coordinates for default reference image.
var Reference = PlateCarre{
	Width:  8192,
	Height: 4096,
}

// Forward converts from LatLong to pixel coordinates.
func (pc *PlateCarre) Forward(s S2) P2 {
	return P2{
		X: pc.Width * (s.Long + 180) / 360,
		Y: pc.Height * (90 - s.Lat) / 180,
	}
}

// Reverse converts from pixel coordinates to LatLong.
func (pc *PlateCarre) Reverse(p P2) S2 {
	return S2{
		Long: (p.X/pc.Width)*360 - 180,
		Lat:  90 - (p.Y/pc.Height)*180,
	}
}
