package view

import (
	"time"

	"github.com/mokiat/gog/opt"
	"github.com/mokiat/gomath/dprec"
	"github.com/mokiat/gomath/sprec"
	"github.com/mokiat/lacking/debug/metric/metricui"
	"github.com/mokiat/lacking/game"
	"github.com/mokiat/lacking/game/graphics"
	"github.com/mokiat/lacking/game/physics"
	"github.com/mokiat/lacking/game/physics/acceleration"
	"github.com/mokiat/lacking/game/physics/collision"
	"github.com/mokiat/lacking/ui"
	co "github.com/mokiat/lacking/ui/component"
	"github.com/mokiat/lacking/ui/layout"
	"github.com/mokiat/lacking/ui/std"

	"github.com/nobonobo/lacking-template/internal/ui/global"
	"github.com/nobonobo/lacking-template/internal/ui/model"
)

var PlayScreen = co.Define(&playScreenComponent{})

type PlayScreenData struct {
	AppModel  *model.ApplicationModel
	PlayModel *model.PlayModel
}

type playScreenComponent struct {
	co.BaseComponent

	engine      *game.Engine
	resourceSet *game.ResourceSet

	appModel  *model.ApplicationModel
	playModel *model.PlayModel
	scene     *model.PlayScene

	debugVisible bool
}

var _ ui.ElementKeyboardHandler = (*playScreenComponent)(nil)

func (c *playScreenComponent) OnCreate() {
	globalContext := co.TypedValue[global.Context](c.Scope())
	c.engine = globalContext.Engine
	c.resourceSet = globalContext.ResourceSet

	data := co.GetData[PlayScreenData](c.Properties())
	c.appModel = data.AppModel
	c.playModel = data.PlayModel

	c.debugVisible = false

	c.scene = c.playModel.Scene()
	if c.scene == nil {
		c.scene = c.createScene()
		c.playModel.SetScene(c.scene)
	}
	c.engine.SetActiveScene(c.scene.Scene)
}

func (c *playScreenComponent) OnDelete() {
	c.engine.SetActiveScene(nil)
}

func (c *playScreenComponent) OnKeyboardEvent(element *ui.Element, event ui.KeyboardEvent) bool {
	switch event.Code {

	case ui.KeyCodeEscape:
		co.Window(c.Scope()).Close()
		return true

	case ui.KeyCodeTab:
		if event.Action == ui.KeyboardActionDown {
			c.debugVisible = !c.debugVisible
			c.Invalidate()
		}
		return true

	default:
		return false
	}
}

func (c *playScreenComponent) Render() co.Instance {
	return co.New(std.Element, func() {
		co.WithData(std.ElementData{
			Essence:   c,
			Focusable: opt.V(true),
			Focused:   opt.V(true),
			Layout:    layout.Anchor(),
		})

		if c.debugVisible {
			co.WithChild("flamegraph", co.New(metricui.FlameGraph, func() {
				co.WithData(metricui.FlameGraphData{
					UpdateInterval: time.Second,
				})
				co.WithLayoutData(layout.Data{
					Top:   opt.V(0),
					Left:  opt.V(0),
					Right: opt.V(0),
				})
			}))
		}
	})
}

func (c *playScreenComponent) createScene() *model.PlayScene {
	sceneData := c.playModel.Data()

	scene := c.engine.CreateScene()

	scene.CreateModel(game.ModelInfo{
		Name:       "Scene",
		Definition: sceneData.Scene,
		IsDynamic:  false,
	})

	boardModel := scene.CreateModel(game.ModelInfo{
		Name:       "Board",
		Definition: sceneData.Board,
		IsDynamic:  false,
	})
	scene.Root().AppendChild(boardModel.Root())

	camera := c.createCamera(scene.Graphics())
	scene.Graphics().SetActiveCamera(camera)

	if cameraNode := boardModel.FindNode("Camera"); cameraNode != nil {
		cameraNode.SetTarget(game.CameraNodeTarget{
			Camera: camera,
		})
	}

	ballModel := scene.CreateModel(game.ModelInfo{
		Name:       "Ball",
		Position:   opt.V(dprec.NewVec3(-2.0, 3.0, 2.0)),
		Definition: sceneData.Ball,
		IsDynamic:  true,
	})

	physicsScene := scene.Physics()
	ballBodyDef := physicsScene.Engine().CreateBodyDefinition(physics.BodyDefinitionInfo{
		Mass:                   1.0,
		MomentOfInertia:        physics.SolidSphereMomentOfInertia(1.0, 1.0),
		FrictionCoefficient:    0.5,
		RestitutionCoefficient: 0.5,
		DragFactor:             0.1,
		AngularDragFactor:      0.1,
		CollisionGroup:         1,
		CollisionSpheres: []collision.Sphere{
			collision.NewSphere(dprec.ZeroVec3(), 1.0),
		},
	})
	ballBody := physicsScene.CreateBody(physics.BodyInfo{
		Name:       "Ball",
		Definition: ballBodyDef,
		Position:   ballModel.Root().Position(),
		Rotation:   ballModel.Root().Rotation(),
	})
	ballBody.SetVelocity(dprec.NewVec3(0.0, 0.0, 3.0))
	ballModel.Root().SetSource(game.BodyNodeSource{
		Body: ballBody,
	})

	physicsScene.CreateGlobalAccelerator(acceleration.NewGravityDirection())

	return &model.PlayScene{
		Scene: scene,
	}
}

func (c *playScreenComponent) createCamera(scene *graphics.Scene) *graphics.Camera {
	result := scene.CreateCamera()
	result.SetFoVMode(graphics.FoVModeHorizontalPlus)
	result.SetFoV(sprec.Degrees(30))
	result.SetAutoExposure(false)
	result.SetExposure(1.0)
	result.SetAutoFocus(false)
	result.SetAutoExposureSpeed(0.1)
	result.SetCascadeDistances([]float32{32.0})
	return result
}
