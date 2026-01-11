package main

import (
	"encoding/gob"
	"log"
	"os"

	"image"
	"image/color"
	"image/png"

	"github.com/hajimehoshi/ebiten/v2"
)

type Game struct {
	img         *ebiten.Image
	terrain_map [][]float64
	camera      *Camera
}

func (g *Game) Update() error {
	g.camera.Update()
	return nil
}

func (g *Game) Draw(screen *ebiten.Image) {
	// Apply camera transformation to terrain
	opts := &ebiten.DrawImageOptions{}
	opts.GeoM = g.camera.GetTransform()
	screen.DrawImage(g.img, opts)
}

func (g *Game) Layout(outsideWidth, outsideHeight int) (screenWidth, screenHeight int) {
	return width, height
}

func main() {
	var img *ebiten.Image
	var terrain_map [][]float64

	// Try to load from cache first
	if terrainCacheExists() {
		log.Println("Loading terrain from cache...")
		var err error
		img, terrain_map, err = loadTerrainCache()
		if err != nil {
			log.Printf("Failed to load cache: %v. Generating new terrain...", err)
			terrain_map, img = generateAndCacheTerrain()
		}
	} else {
		log.Println("No cache found. Generating terrain...")
		terrain_map, img = generateAndCacheTerrain()
	}

	// Initialize camera centered on the world
	// World is terrainResolution times larger than screen
	//camera := NewCamera(float64(width*terrainResolution), float64(height*terrainResolution), 1.0)
	camera := NewCamera(0, 0, 1.0)

	ebiten.SetWindowSize(width, height)
	ebiten.SetWindowTitle("Lineagen!")
	if err := ebiten.RunGame(&Game{
		img:         img,
		terrain_map: terrain_map,
		camera:      camera,
	}); err != nil {
		log.Fatal(err)
	}
}

// generateAndCacheTerrain generates new terrain and saves it to cache
func generateAndCacheTerrain() ([][]float64, *ebiten.Image) {
	terrain_map := generate_terrain()
	rgba := createLandscapeImageRGBA(terrain_map)

	// Save to cache
	if err := saveTerrainCache(rgba, terrain_map); err != nil {
		log.Printf("Warning: Failed to save terrain cache: %v", err)
	}

	// Convert to Ebiten image
	return terrain_map, ebiten.NewImageFromImage(rgba)
}

func createLandscapeImageRGBA(data [][]float64) *image.RGBA {
	terrainHeight := len(data)
	terrainWidth := len(data[0])

	// Create image at full screen resolution
	rgba := image.NewRGBA(image.Rect(0, 0, width, height))

	// Terrain color palette - techy and simple
	deepWater := color.RGBA{20, 60, 140, 255}     // Dark blue - deep water
	shallowWater := color.RGBA{40, 120, 200, 255} // Medium blue - shallow water
	beach := color.RGBA{210, 200, 140, 255}       // Sandy tan - beach/shore
	grass := color.RGBA{80, 180, 80, 255}         // Green - grassland (most common)
	forest := color.RGBA{40, 120, 40, 255}        // Dark green - dense vegetation
	dryLand := color.RGBA{160, 140, 100, 255}     // Brown/tan - dry rocky areas

	// Terrain thresholds
	deepWaterThreshold := 0.35
	shallowWaterThreshold := 0.42
	beachThreshold := 0.45
	grassThreshold := 0.65
	forestThreshold := 0.75

	// Scale up each terrain pixel
	for ty := 0; ty < terrainHeight; ty++ {
		for tx := 0; tx < terrainWidth; tx++ {
			var pixelColor color.RGBA
			value := data[ty][tx]

			// Determine terrain type based on noise value
			if value < deepWaterThreshold {
				pixelColor = deepWater
			} else if value < shallowWaterThreshold {
				pixelColor = shallowWater
			} else if value < beachThreshold {
				pixelColor = beach
			} else if value < grassThreshold {
				pixelColor = grass
			} else if value < forestThreshold {
				pixelColor = forest
			} else {
				pixelColor = dryLand
			}

			// Draw this terrain pixel as a scale x scale block
			for dy := 0; dy < scale; dy++ {
				for dx := 0; dx < scale; dx++ {
					screenX := tx*scale + dx
					screenY := ty*scale + dy
					rgba.Set(screenX, screenY, pixelColor)
				}
			}
		}
	}

	return rgba
}

// saveTerrainCache saves both the terrain image and data array to disk
func saveTerrainCache(img *image.RGBA, terrainData [][]float64) error {
	// Save PNG image
	imgFile, err := os.Create("terrain_cache.png")
	if err != nil {
		return err
	}
	defer imgFile.Close()

	if err := png.Encode(imgFile, img); err != nil {
		return err
	}

	// Save terrain data array using gob encoding
	dataFile, err := os.Create("terrain_cache.dat")
	if err != nil {
		return err
	}
	defer dataFile.Close()

	encoder := gob.NewEncoder(dataFile)
	if err := encoder.Encode(terrainData); err != nil {
		return err
	}

	log.Println("Terrain cached successfully!")
	return nil
}

// loadTerrainCache loads both the terrain image and data array from disk
func loadTerrainCache() (*ebiten.Image, [][]float64, error) {
	// Load PNG image
	imgFile, err := os.Open("terrain_cache.png")
	if err != nil {
		return nil, nil, err
	}
	defer imgFile.Close()

	imgDecoded, err := png.Decode(imgFile)
	if err != nil {
		return nil, nil, err
	}

	// Load terrain data array
	dataFile, err := os.Open("terrain_cache.dat")
	if err != nil {
		return nil, nil, err
	}
	defer dataFile.Close()

	var terrainData [][]float64
	decoder := gob.NewDecoder(dataFile)
	if err := decoder.Decode(&terrainData); err != nil {
		return nil, nil, err
	}

	log.Println("Terrain loaded from cache!")
	return ebiten.NewImageFromImage(imgDecoded), terrainData, nil
}

// terrainCacheExists checks if cached terrain files exist
func terrainCacheExists() bool {
	_, err1 := os.Stat("terrain_cache.png")
	_, err2 := os.Stat("terrain_cache.dat")
	return err1 == nil && err2 == nil
}
