package main

import (
	"embed"
	"log"

	"github.com/wailsapp/wails/v2"
	"github.com/wailsapp/wails/v2/pkg/options"
	"github.com/wailsapp/wails/v2/pkg/options/assetserver"
)

//go:embed all:frontend/dist
var assets embed.FS

func main() {
	app := NewApp()

	err := wails.Run(&options.App{
		Title:  "BTC Price Tracker",
		Width:  420,
		Height: 320,
		AssetServer: &assetserver.Options{
			Assets: assets,
		},
		OnStartup:  app.startup,
		Bind:       []interface{}{app},
		Frameless:  false,
		MinWidth:   360,
		MinHeight:  280,
	})
	if err != nil {
		log.Fatal(err)
	}
}
