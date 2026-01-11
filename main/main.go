package main

import (
	"log"

	"image"
	"image/color"

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
	terrain_map := generate_terrain()
	img := createLandscapeImage(terrain_map)

	// Initialize camera centered on the world
	camera := NewCamera(float64(width)/2, float64(height)/2, 1.0)

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

func createLandscapeImage(data [][]float64) *ebiten.Image {
	terrainHeight := len(data)
	terrainWidth := len(data[0])

	// Create image at full screen resolution
	rgba := image.NewRGBA(image.Rect(0, 0, width, height))

	blue := color.RGBA{0, 0, 255, 255}
	middleBlue := color.RGBA{0, 128, 255, 255}
	green := color.RGBA{0, 255, 0, 255}

	deepWaterHeight := 0.4
	ShallowHeight := 0.45

	// Scale up each terrain pixel
	for ty := 0; ty < terrainHeight; ty++ {
		for tx := 0; tx < terrainWidth; tx++ {
			var pixelColor color.RGBA
			if data[ty][tx] < deepWaterHeight {
				pixelColor = blue
			} else if data[ty][tx] < ShallowHeight && data[ty][tx] >= deepWaterHeight {
				pixelColor = middleBlue
			} else {
				pixelColor = green
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

	// Convert to Ebiten image
	return ebiten.NewImageFromImage(rgba)
}
