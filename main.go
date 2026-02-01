package main

import (
	"embed"

	"github.com/wailsapp/wails/v2"
	"github.com/wailsapp/wails/v2/pkg/options"
	"github.com/wailsapp/wails/v2/pkg/options/assetserver"
	"github.com/wailsapp/wails/v2/pkg/options/mac"
)

//go:embed all:frontend/dist
var assets embed.FS

func main() {
	// Initialize logging FIRST
	if err := InitLogger(); err != nil {
		println("Warning: Could not initialize file logging:", err.Error())
	}
	defer CloseLogger()

	LogInfo("=== Folje Application Starting ===")

	// Create an instance of the app structure
	app := NewApp()

	// Create application with options
	err := wails.Run(&options.App{
		Title:  "Följe",
		Width:  1024,
		Height: 768,
		AssetServer: &assetserver.Options{
			Assets: assets,
		},
		BackgroundColour:         &options.RGBA{R: 8, G: 18, B: 34, A: 1},
		OnStartup:                app.startup,
		OnShutdown:               app.shutdown,
		WindowStartState:         options.Maximised,
		EnableDefaultContextMenu: false,
		Bind: []interface{}{
			app,
		},
		Mac: &mac.Options{
			TitleBar: &mac.TitleBar{
				FullSizeContent: true,
			},
		},
	})

	if err != nil {
		LogError("Wails error: %s", err.Error())
	}

	LogInfo("=== Folje Application Exiting ===")
}
