package main

// Hardcoded calculations for:
//
//   Projection: Plate Carree aka Geographic or "LatLong"
//   Earth ellipsoid: Sphere, radius 6370997 m
//   Extent: 180 West to 180 East, 90 North to 90 South
//   Size: 8,192 height samples wide x 4,096 high

type S2 struct {
	Lat  float32
	Long float32
}

type P2 struct {
	X float32
	Y float32
}

type PlateCarre struct {
	Width  float32
	Height float32
}

var Reference = PlateCarre{
	Width:  8192,
	Height: 4096,
}

func (pc *PlateCarre) Forward(s S2) P2 {
	return P2{
		X: pc.Width * (s.Long + 180) / 360,
		Y: pc.Height * (90 - s.Lat) / 180,
	}
}

func (pc *PlateCarre) Reverse(p P2) S2 {
	return S2{
		Long: (p.X/pc.Width)*360 - 180,
		Lat:  90 - (p.Y/pc.Height)*180,
	}
}
