package main

import (
	"math"
	"testing"
)

func TestForwardReverse(t *testing.T) {
	for long := -180; long <= 180; long += 10 {
		for lat := -90; lat <= 90; lat += 10 {
			s := S2{
				Long: float32(long),
				Lat:  float32(lat),
			}
			p := Reference.Forward(s)
			if p.X < 0 || p.X > 8192 || p.Y < 0 || p.Y > 4096 {
				t.Errorf("p.X coordinate out of bounds %v", p)
				continue
			}

			r := Reference.Reverse(p)

			const tol = 0.0001
			dlong := math.Abs(float64(r.Long - s.Long))
			dlat := math.Abs(float64(r.Lat - s.Lat))
			if dlong > tol || dlat > tol {
				t.Errorf("%v above tolerance dlong=%v dlat=%v", s, dlong, dlat)
			}
		}
	}
}
