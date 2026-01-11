package main

import (
	"fmt"

	"github.com/KEINOS/go-noise"
)

func generate_terrain(cfg *Config) [][]float64 {
	terrainHeight := (cfg.Window.Height * cfg.Terrain.Resolution) / cfg.Terrain.Scale
	terrainWidth := (cfg.Window.Width * cfg.Terrain.Resolution) / cfg.Terrain.Scale

	terrain_map := make([][]float64, terrainHeight)

	for i := range terrainHeight {
		terrain_map[i] = make([]float64, terrainWidth)
	}
	for y := 0; y < terrainHeight; y++ {
		for x := 0; x < terrainWidth; x++ {
			noise_value := generate_perlin_noise(float64(x), float64(y), cfg)
			terrain_map[y][x] = (noise_value + 1) / 2
		}
	}

	return terrain_map
}

func generate_perlin_noise(x float64, y float64, cfg *Config) float64 {
	seed := int64(cfg.Terrain.Noise.Seed)
	smoothness := cfg.Terrain.Noise.BaseSmoothness * cfg.Terrain.Resolution

	n, err := noise.New(noise.Perlin, seed)

	if err != nil {
		fmt.Println("Error creating noise generator:", err)
	}

	v := n.Eval64(x/float64(smoothness), y/float64(smoothness))

	return v
}
