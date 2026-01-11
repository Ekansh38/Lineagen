package main

import (
	"fmt"

	"github.com/KEINOS/go-noise"
)

const width = 2560
const height = 1440
const scale = 1
const terrainResolution = 3

func generate_terrain() [][]float64 {
	terrainHeight := (height * terrainResolution) / scale
	terrainWidth := (width * terrainResolution) / scale

	terrain_map := make([][]float64, terrainHeight)

	for i := range terrainHeight {
		terrain_map[i] = make([]float64, terrainWidth)
	}
	for y := 0; y < terrainHeight; y++ {
		for x := 0; x < terrainWidth; x++ {
			noise_value := generate_perlin_noise(float64(x), float64(y))
			terrain_map[y][x] = (noise_value + 1) / 2
		}
	}

	return terrain_map
}

func generate_perlin_noise(x float64, y float64) float64 {
	const seed = 6767                                // noise pattern ID
	const baseSmoothness = 140                       // base noise smoothness
	smoothness := baseSmoothness * terrainResolution // scale smoothness to maintain feature size

	n, err := noise.New(noise.Perlin, seed)

	if err != nil {
		fmt.Println("Error creating noise generator:", err)
	}

	v := n.Eval64(x/float64(smoothness), y/float64(smoothness)) // yy is between -1.0 and 1.0 of float64

	return v

}
