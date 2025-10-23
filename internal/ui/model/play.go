package model

import (
	"github.com/mokiat/lacking/game"

	"github.com/nobonobo/lacking-template/internal/game/data"
)

func NewPlayModel() *PlayModel {
	return &PlayModel{}
}

type PlayModel struct {
	sceneData *data.PlayData
	scene     *PlayScene
}

func (h *PlayModel) Data() *data.PlayData {
	return h.sceneData
}

func (h *PlayModel) SetData(data *data.PlayData) {
	h.sceneData = data
}

func (h *PlayModel) Scene() *PlayScene {
	return h.scene
}

func (h *PlayModel) SetScene(scene *PlayScene) {
	h.scene = scene
}

type PlayScene struct {
	Scene *game.Scene
}
