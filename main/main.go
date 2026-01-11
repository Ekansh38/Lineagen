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
	cfg         *Config
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
	return g.cfg.Window.Width, g.cfg.Window.Height
}

func main() {
	// Load configuration
	cfg := MustLoadConfig("../config.toml")

	var img *ebiten.Image
	var terrain_map [][]float64

	// Try to load from cache first
	if terrainCacheExists() {
		log.Println("Loading terrain from cache...")
		var err error
		img, terrain_map, err = loadTerrainCache()
		if err != nil {
			log.Printf("Failed to load cache: %v. Generating new terrain...", err)
			terrain_map, img = generateAndCacheTerrain(cfg)
		}
	} else {
		log.Println("No cache found. Generating terrain...")
		terrain_map, img = generateAndCacheTerrain(cfg)
	}

	// Initialize camera with config settings
	camera := NewCamera(
		cfg.Camera.InitialX,
		cfg.Camera.InitialY,
		cfg.Camera.InitialZoom,
		cfg,
	)

	ebiten.SetWindowSize(cfg.Window.Width, cfg.Window.Height)
	ebiten.SetWindowTitle(cfg.Window.Title)
	if err := ebiten.RunGame(&Game{
		img:         img,
		terrain_map: terrain_map,
		camera:      camera,
		cfg:         cfg,
	}); err != nil {
		log.Fatal(err)
	}
}

// generateAndCacheTerrain generates new terrain and saves it to cache
func generateAndCacheTerrain(cfg *Config) ([][]float64, *ebiten.Image) {
	terrain_map := generate_terrain(cfg)
	rgba := createLandscapeImageRGBA(terrain_map, cfg)

	// Save to cache
	if err := saveTerrainCache(rgba, terrain_map); err != nil {
		log.Printf("Warning: Failed to save terrain cache: %v", err)
	}

	// Convert to Ebiten image
	return terrain_map, ebiten.NewImageFromImage(rgba)
}

func createLandscapeImageRGBA(data [][]float64, cfg *Config) *image.RGBA {
	terrainHeight := len(data)
	terrainWidth := len(data[0])

	// Create image at full screen resolution
	rgba := image.NewRGBA(image.Rect(0, 0, cfg.Window.Width, cfg.Window.Height))

	// Terrain color palette from config
	deepWater := cfg.Terrain.Colors.DeepWater.ToColor()
	shallowWater := cfg.Terrain.Colors.ShallowWater.ToColor()
	beach := cfg.Terrain.Colors.Beach.ToColor()
	grass := cfg.Terrain.Colors.Grass.ToColor()
	forest := cfg.Terrain.Colors.Forest.ToColor()
	dryLand := cfg.Terrain.Colors.DryLand.ToColor()

	// Terrain thresholds from config
	deepWaterThreshold := cfg.Terrain.Biomes.DeepWater
	shallowWaterThreshold := cfg.Terrain.Biomes.ShallowWater
	beachThreshold := cfg.Terrain.Biomes.Beach
	grassThreshold := cfg.Terrain.Biomes.Grass
	forestThreshold := cfg.Terrain.Biomes.Forest

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
			for dy := 0; dy < cfg.Terrain.Scale; dy++ {
				for dx := 0; dx < cfg.Terrain.Scale; dx++ {
					screenX := tx*cfg.Terrain.Scale + dx
					screenY := ty*cfg.Terrain.Scale + dy
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
