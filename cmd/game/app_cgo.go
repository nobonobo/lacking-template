//go:build !js

package main

import (
	"fmt"

	nativeapp "github.com/mokiat/lacking-native/app"
	nativegame "github.com/mokiat/lacking-native/game"
	nativeui "github.com/mokiat/lacking-native/ui"
	"github.com/mokiat/lacking/app"
	"github.com/mokiat/lacking/game"
	"github.com/mokiat/lacking/game/asset"
	"github.com/mokiat/lacking/ui"
	"github.com/mokiat/lacking/util/resource"
	gameui "github.com/nobonobo/lacking-template/internal/ui"
	"github.com/nobonobo/lacking-template/resources"
)

func runApplication() error {
	storage, err := asset.NewFSStorage(".")
	if err != nil {
		return err
	}
	formatter := asset.NewBlobFormatter()
	registry, err := asset.NewRegistry(storage, formatter)
	if err != nil {
		return fmt.Errorf("failed to initialize registry: %w", err)
	}
	gameController := game.NewController(registry, nativegame.NewShaderCollection(), nativegame.NewShaderBuilder())

	locator := ui.WrappedLocator(resource.NewFSLocator(resources.UI))

	uiController := ui.NewController(locator, nativeui.NewShaderCollection(), func(w *ui.Window) {
		gameui.BootstrapApplication(w)
	})

	cfg := nativeapp.NewConfig("Game", 1280, 800)
	cfg.SetFullscreen(false)
	cfg.SetMaximized(false)
	cfg.SetMinSize(1024, 576)
	cfg.SetVSync(true)
	cfg.SetIcon("ui/images/icon.png")
	cfg.SetLocator(locator)
	cfg.SetAudioEnabled(false)
	return nativeapp.Run(cfg, app.NewLayeredController(gameController, uiController))
}
