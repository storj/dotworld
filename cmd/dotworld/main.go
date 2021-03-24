// Copyright (C) 2021 Storj Labs, Inc.
// See LICENSE for copying information.

package main

import (
	"bytes"
	"context"
	"flag"
	"fmt"
	"image"
	"image/png"
	"io/ioutil"
	"math/rand"
	"os"
	"time"

	"storj.io/dotworld"
)

// Config defines arguments for map generation.
type Config struct {
	Reference     string
	ReferenceArea image.Rectangle

	GridCountX    int
	GridThreshold float64

	Output string
}

func main() {
	ctx := context.Background()
	rand.Seed(time.Now().UnixNano())

	config := Config{}

	flag.StringVar(&config.Reference, "reference", "reference/world.png", "reference world image")
	flag.IntVar(&config.ReferenceArea.Min.X, "reference.min.x", 0, "rect bounds for reference image")
	flag.IntVar(&config.ReferenceArea.Min.Y, "reference.min.y", 100, "rect bounds for reference image")
	flag.IntVar(&config.ReferenceArea.Max.X, "reference.max.x", 8192, "rect bounds for reference image")
	flag.IntVar(&config.ReferenceArea.Max.Y, "reference.max.y", 3385, "rect bounds for reference image")

	flag.IntVar(&config.GridCountX, "grid.count.x", 64, "dots in x dimension")
	flag.Float64Var(&config.GridThreshold, "grid.threshold", 0.2, "luminosity threshold for grid")

	flag.Parse()

	if err := run(ctx, flag.Arg(0), config); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

func run(ctx context.Context, command string, config Config) error {
	switch command {
	case "random-map":
		return generateRandomMap(ctx, config)
	default:
		return fmt.Errorf("unknown command %q", command)
	}
}

func generateMap(_ context.Context, config Config) (*dotworld.Map, error) {
	worlddata, err := ioutil.ReadFile(config.Reference)
	if err != nil {
		return nil, fmt.Errorf("failed to load %q: %w", config.Reference, err)
	}

	rm, err := png.Decode(bytes.NewReader(worlddata))
	if err != nil {
		return nil, fmt.Errorf("unable to decode: %w", err)
	}
	m := rm.(*image.Gray)

	grid := dotworld.Grid{
		CountX:    config.GridCountX,
		Threshold: float32(config.GridThreshold),

		Coord: dotworld.PlateCarre{
			Width:  float32(m.Bounds().Dx()),
			Height: float32(m.Bounds().Dy()),
		},
	}

	subimage := m
	if !config.ReferenceArea.Empty() {
		subimage = m.SubImage(config.ReferenceArea).(*image.Gray)
	}

	return grid.MapFromImage(subimage), nil
}

func generateRandomMap(ctx context.Context, config Config) error {
	dotmap, err := generateMap(ctx, config)
	if err != nil {
		return err
	}

	const N = 80
	for k := 0; k < N; k++ {
		gridPos := dotmap.Lookup(dotworld.S2{
			Lat:  float32(rand.Intn(180) - 90),
			Long: float32(rand.Intn(360) - 180),
		})
		if nearest := dotmap.Nearest(gridPos); nearest != nil {
			nearest.Load += 1.0 / N
		}
	}

	_ = dotmap.EncodeSVG(os.Stdout, 800, 400)
	return nil
}
