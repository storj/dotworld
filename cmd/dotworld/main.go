// Copyright (C) 2021 Storj Labs, Inc.
// See LICENSE for copying information.

package main

import (
	"context"
	"flag"
	"fmt"
	"math/rand"
	"os"
	"time"

	"storj.io/dotworld"
	"storj.io/dotworld/reference"
)

func main() {
	ctx := context.Background()
	rand.Seed(time.Now().UnixNano())
	flag.Parse()

	if err := run(ctx, flag.Arg(0)); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

func run(ctx context.Context, command string) error {
	switch command {
	case "random-map":
		return generateRandomMap(ctx)
	default:
		return fmt.Errorf("unknown command %q", command)
	}
}

func generateRandomMap(ctx context.Context) error {
	dotmap := reference.WorldMap()

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
