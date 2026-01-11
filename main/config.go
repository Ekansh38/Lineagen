package main

import (
	"fmt"
	"image/color"
	"log"
	"os"

	"github.com/BurntSushi/toml"
)

// Config represents the entire configuration file structure
type Config struct {
	Window  WindowConfig  `toml:"window"`
	Terrain TerrainConfig `toml:"terrain"`
	Camera  CameraConfig  `toml:"camera"`
}

type WindowConfig struct {
	Width  int    `toml:"width"`
	Height int    `toml:"height"`
	Title  string `toml:"title"`
}

type TerrainConfig struct {
	Scale      int                `toml:"scale"`
	Resolution int                `toml:"resolution"`
	Noise      NoiseConfig        `toml:"noise"`
	Biomes     BiomesConfig       `toml:"biomes"`
	Colors     TerrainColorsConfig `toml:"colors"`
}

type NoiseConfig struct {
	Seed           int `toml:"seed"`
	BaseSmoothness int `toml:"base_smoothness"`
}

type BiomesConfig struct {
	DeepWater    float64 `toml:"deep_water"`
	ShallowWater float64 `toml:"shallow_water"`
	Beach        float64 `toml:"beach"`
	Grass        float64 `toml:"grass"`
	Forest       float64 `toml:"forest"`
}

type TerrainColorsConfig struct {
	DeepWater    ColorRGBA `toml:"deep_water"`
	ShallowWater ColorRGBA `toml:"shallow_water"`
	Beach        ColorRGBA `toml:"beach"`
	Grass        ColorRGBA `toml:"grass"`
	Forest       ColorRGBA `toml:"forest"`
	DryLand      ColorRGBA `toml:"dry_land"`
}

type ColorRGBA struct {
	R uint8 `toml:"r"`
	G uint8 `toml:"g"`
	B uint8 `toml:"b"`
	A uint8 `toml:"a"`
}

// ToColor converts ColorRGBA to color.RGBA
func (c ColorRGBA) ToColor() color.RGBA {
	return color.RGBA{R: c.R, G: c.G, B: c.B, A: c.A}
}

type CameraConfig struct {
	InitialX   float64 `toml:"initial_x"`
	InitialY   float64 `toml:"initial_y"`
	InitialZoom float64 `toml:"initial_zoom"`
	ZoomMin    float64 `toml:"zoom_min"`
	ZoomMax    float64 `toml:"zoom_max"`
	ZoomFactor float64 `toml:"zoom_factor"`
	PanSpeed   float64 `toml:"pan_speed"`
}

// LoadConfig loads the configuration from the config.toml file
func LoadConfig(configPath string) (*Config, error) {
	// Check if file exists
	if _, err := os.Stat(configPath); os.IsNotExist(err) {
		return nil, fmt.Errorf("config file not found: %s", configPath)
	}

	var config Config
	if _, err := toml.DecodeFile(configPath, &config); err != nil {
		return nil, fmt.Errorf("failed to parse config file: %w", err)
	}

	log.Printf("Configuration loaded from %s", configPath)
	return &config, nil
}

// MustLoadConfig loads the config or panics if it fails
func MustLoadConfig(configPath string) *Config {
	config, err := LoadConfig(configPath)
	if err != nil {
		log.Fatalf("Failed to load configuration: %v", err)
	}
	return config
}
