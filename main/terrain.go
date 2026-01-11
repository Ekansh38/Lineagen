package main

import (
	"fmt"

	"github.com/KEINOS/go-noise"
)

func generate_terrain() [][]bool {
	const width = 1280
	const height = 720

	terrain_map := make([][]float64, height)
	for i := range terrain_map {
		terrain_map[i] = make([]float64, width)
	}
	for y := 0; y < height; y++ {
		for x := 0; x < width; x++ {
			noise_value := generate_perlin_noise(float64(x), float64(y))
			terrain_map[y][x] = (noise_value + 1) / 2
		}
	}

	// convert to boolean map for land/water
	land_water_map := make([][]bool, height)
	for i := range land_water_map {
		land_water_map[i] = make([]bool, width)
	}
	for y := 0; y < height; y++ {
		for x := 0; x < width; x++ {
			if terrain_map[y][x] < 0.4 {
				land_water_map[y][x] = false // water
			} else {
				land_water_map[y][x] = true // land
			}
		}
	}

	// find number of true and false and print that.

	true_count := 0
	false_count := 0

	for y := 0; y < height; y++ {
		for x := 0; x < width; x++ {
			if land_water_map[y][x] {
				true_count++
			} else {
				false_count++
			}
		}
	}

	fmt.Println("Land count:", true_count)
	fmt.Println("Water count:", false_count)

	return land_water_map
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
