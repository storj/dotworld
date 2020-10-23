package main

import (
	"bytes"
	"context"
	"fmt"
	"image"
	"image/png"
	"io/ioutil"
	"math/rand"
	"os"
	"time"
)

func main() {
	ctx := context.Background()

	rand.Seed(time.Now().UnixNano())
	if err := run(ctx, "reference/world.png"); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

func run(_ context.Context, imagepath string) error {
	worlddata, err := ioutil.ReadFile(imagepath)
	if err != nil {
		return fmt.Errorf("failed to load %q: %w", imagepath, err)
	}

	rm, err := png.Decode(bytes.NewReader(worlddata))
	if err != nil {
		return fmt.Errorf("unable to decode: %w", err)
	}
	m := rm.(*image.Gray)

	grid := Grid{
		CountX:    64,
		Threshold: 0.2,

		Coord: PlateCarre{
			Width:  float32(m.Bounds().Dx()),
			Height: float32(m.Bounds().Dy()),
		},
	}

	noAntartica := m.SubImage(image.Rect(0, 100, 8192, 3385)).(*image.Gray)
	dotmap := grid.MapFromImage(noAntartica)

	const N = 80
	for k := 0; k < N; k++ {
		i := rand.Intn(len(dotmap.Locations))
		dotmap.Locations[i].Load += 1.0 / N
	}

	var b bytes.Buffer
	_ = dotmap.EncodeSVG(&b, 800, 400)
	ioutil.WriteFile("map.svg", b.Bytes(), 0666)

	return nil
}
