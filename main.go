package main

import (
	"github.com/200sc/gel/gel"
	"github.com/oakmound/oak"
)

// In Go, main is the entry point for main packages. Code starts executing here.
func main() {
	// We set the width and height of the rendering screen to 800x600 pixels
	oak.SetupConfig.Screen = oak.Screen{
		Width:  800,
		Height: 600,
	}
	// Add our only scene, named "gel", to oak
	// gel/scene.go is the next place our application's code will execute.
	oak.AddScene("gel", gel.SceneStart, gel.SceneLoop, gel.SceneEnd)
	// Start oak at the scene "gel"
	oak.Init("gel")
}
