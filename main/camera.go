package main

import (
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
)

type Camera struct {
	X, Y     float64 // Camera center position in world coordinates
	Zoom     float64 // Zoom level (1.0 = normal, 2.0 = zoomed in 2x)
	dragging bool
	lastX    int
	lastY    int
	cfg      *Config // Configuration reference
}

// NewCamera creates a new camera centered at the given world position
func NewCamera(x, y, zoom float64, cfg *Config) *Camera {
	return &Camera{
		X:    x,
		Y:    y,
		Zoom: zoom,
		cfg:  cfg,
	}
}

// GetTransform returns the geometry matrix for rendering
// This transforms world coordinates to screen coordinates
func (c *Camera) GetTransform() ebiten.GeoM {
	var geo ebiten.GeoM

	// 1. Translate so camera center is at origin
	geo.Translate(-c.X, -c.Y)

	// 2. Apply zoom
	geo.Scale(c.Zoom, c.Zoom)

	// 3. Translate to screen center
	geo.Translate(float64(c.cfg.Window.Width)/2, float64(c.cfg.Window.Height)/2)

	return geo
}

// ScreenToWorld converts screen coordinates to world coordinates
func (c *Camera) ScreenToWorld(screenX, screenY int) (float64, float64) {
	// Inverse of GetTransform
	worldX := (float64(screenX)-float64(c.cfg.Window.Width)/2)/c.Zoom + c.X
	worldY := (float64(screenY)-float64(c.cfg.Window.Height)/2)/c.Zoom + c.Y
	return worldX, worldY
}

// WorldToScreen converts world coordinates to screen coordinates
func (c *Camera) WorldToScreen(worldX, worldY float64) (int, int) {
	screenX := (worldX-c.X)*c.Zoom + float64(c.cfg.Window.Width)/2
	screenY := (worldY-c.Y)*c.Zoom + float64(c.cfg.Window.Height)/2
	return int(screenX), int(screenY)
}

// Update handles camera input (zoom and pan)
func (c *Camera) Update() {
	// Mouse wheel zoom
	_, scrollY := ebiten.Wheel()
	if scrollY != 0 {
		// Zoom in/out based on config
		zoomFactor := 1.0 + scrollY*c.cfg.Camera.ZoomFactor
		c.Zoom *= zoomFactor

		// Clamp zoom between configured limits
		if c.Zoom < c.cfg.Camera.ZoomMin {
			c.Zoom = c.cfg.Camera.ZoomMin
		} else if c.Zoom > c.cfg.Camera.ZoomMax {
			c.Zoom = c.cfg.Camera.ZoomMax
		}
	}

	// Click and drag to pan
	if inpututil.IsMouseButtonJustPressed(ebiten.MouseButtonLeft) {
		c.dragging = true
		c.lastX, c.lastY = ebiten.CursorPosition()
	}

	if inpututil.IsMouseButtonJustReleased(ebiten.MouseButtonLeft) {
		c.dragging = false
	}

	if c.dragging {
		currentX, currentY := ebiten.CursorPosition()
		deltaX := float64(currentX - c.lastX)
		deltaY := float64(currentY - c.lastY)

		// Move camera in opposite direction of drag (in world space)
		c.X -= deltaX / c.Zoom
		c.Y -= deltaY / c.Zoom

		c.lastX = currentX
		c.lastY = currentY
	}

	// Keyboard panning (arrow keys) - alternative/additional control
	panSpeed := c.cfg.Camera.PanSpeed / c.Zoom // Pan speed inversely proportional to zoom
	if ebiten.IsKeyPressed(ebiten.KeyArrowLeft) {
		c.X -= panSpeed
	}
	if ebiten.IsKeyPressed(ebiten.KeyArrowRight) {
		c.X += panSpeed
	}
	if ebiten.IsKeyPressed(ebiten.KeyArrowUp) {
		c.Y -= panSpeed
	}
	if ebiten.IsKeyPressed(ebiten.KeyArrowDown) {
		c.Y += panSpeed
	}

	// Reset camera with 'R' key
	if inpututil.IsKeyJustPressed(ebiten.KeyR) {
		c.X = c.cfg.Camera.InitialX
		c.Y = c.cfg.Camera.InitialY
		c.Zoom = c.cfg.Camera.InitialZoom
	}
}
