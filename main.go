package main

import (
	"github.com/200sc/gel/gel"
	"github.com/oakmound/oak"
)

func main() {
	oak.SetupConfig.Screen = oak.Screen{
		Width:  800,
		Height: 600,
	}
	oak.SetupConfig.Assets = oak.Assets{
		AssetPath: "\\",
		ImagePath: "model",
	}
	oak.AddScene("gel", gel.Scene)
	oak.Init("gel")
}
