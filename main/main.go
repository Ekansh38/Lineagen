package main

import (
	"log"

	"image"
	"image/color"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
)

type Game struct {
	img *ebiten.Image
}

func (g *Game) Update() error {
	return nil
}

func (g *Game) Draw(screen *ebiten.Image) {
	screen.DrawImage(g.img, nil)

	ebitenutil.DebugPrint(screen, "Hello, World!")
}

func (g *Game) Layout(outsideWidth, outsideHeight int) (screenWidth, screenHeight int) {
	return 320, 240
}

func main() {
	landWaterMap := generate_terrain()
	img := createImageFromArray(landWaterMap)

	ebiten.SetWindowSize(1280, 720)
	ebiten.SetWindowTitle("Hello, World!")
	if err := ebiten.RunGame(&Game{img: img}); err != nil {
		log.Fatal(err)
	}
}

func createImageFromArray(data [][]bool) *ebiten.Image {
	height := len(data)
	width := len(data[0])

	// Create standard Go image
	rgba := image.NewRGBA(image.Rect(0, 0, width, height))

	blue := color.RGBA{0, 0, 255, 255}
	green := color.RGBA{0, 255, 0, 255}

	for y := 0; y < height; y++ {
		for x := 0; x < width; x++ {
			if data[y][x] {
				rgba.Set(x, y, green)
			} else {
				rgba.Set(x, y, blue)
			}
		}
	}

	// Convert to Ebiten image
	return ebiten.NewImageFromImage(rgba)
}
