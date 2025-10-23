package model

import (
	"github.com/mokiat/lacking/game"

	"github.com/nobonobo/lacking-template/internal/game/data"
)

func NewHomeModel() *HomeModel {
	return &HomeModel{}
}

type HomeModel struct {
	sceneData *data.HomeData
	scene     *HomeScene
}

func (h *HomeModel) Data() *data.HomeData {
	return h.sceneData
}

func (h *HomeModel) SetData(data *data.HomeData) {
	h.sceneData = data
}

func (h *HomeModel) Scene() *HomeScene {
	return h.scene
}

func (h *HomeModel) SetScene(scene *HomeScene) {
	h.scene = scene
}

type HomeScene struct {
	Scene *game.Scene
}
