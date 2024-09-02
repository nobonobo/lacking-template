//go:build js

package main

import (
	"fmt"

	jsapp "github.com/mokiat/lacking-js/app"
	jsgame "github.com/mokiat/lacking-js/game"
	jsui "github.com/mokiat/lacking-js/ui"
	"github.com/mokiat/lacking/app"
	"github.com/mokiat/lacking/game"
	"github.com/mokiat/lacking/game/asset"
	"github.com/mokiat/lacking/ui"
	"github.com/mokiat/lacking/util/resource"
	gameui "github.com/nobonobo/lacking-template/internal/ui"
	"github.com/nobonobo/lacking-template/resources"
)

func runApplication() error {
	storage, err := asset.NewWebStorage(".")
	if err != nil {
		return err
	}
	formatter := asset.NewBlobFormatter()
	registry, err := asset.NewRegistry(storage, formatter)
	if err != nil {
		return fmt.Errorf("failed to initialize registry: %w", err)
	}
	gameController := game.NewController(registry, jsgame.NewShaderCollection(), jsgame.NewShaderBuilder())

	locator := ui.WrappedLocator(resource.NewFSLocator(resources.UI))
	uiController := ui.NewController(locator, jsui.NewShaderCollection(), func(w *ui.Window) {
		gameui.BootstrapApplication(w)
	})

	cfg := jsapp.NewConfig("screen")
	cfg.AddGLExtension("EXT_color_buffer_float")
	cfg.SetFullscreen(false)
	return jsapp.Run(cfg, app.NewLayeredController(gameController, uiController))
}
