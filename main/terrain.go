package main

import (
	"fmt"

	"github.com/KEINOS/go-noise"
)

const width = 1280
const height = 720
const scale = 1

func generate_terrain() [][]float64 {
	// Generate terrain at reduced resolution based on scale
	scaledHeight := height / scale
	scaledWidth := width / scale

	terrain_map := make([][]float64, scaledHeight)

	for i := range scaledHeight {
		terrain_map[i] = make([]float64, scaledWidth)
	}
	for y := 0; y < scaledHeight; y++ {
		for x := 0; x < scaledWidth; x++ {
			noise_value := generate_perlin_noise(float64(x), float64(y))
			terrain_map[y][x] = (noise_value + 1) / 2
		}
	}

	return terrain_map
}

func generate_perlin_noise(x float64, y float64) float64 {
	const seed = 492       // noise pattern ID
	const smoothness = 140 // noise smoothness

	n, err := noise.New(noise.Perlin, seed)

	if err != nil {
		fmt.Println("Error creating noise generator:", err)
	}

	v := n.Eval64(x/smoothness, y/smoothness) // yy is between -1.0 and 1.0 of float64

	return v

}
