package camera

import (
	"main/utils"

	"github.com/hajimehoshi/ebiten/v2"
)

type CameraStruct struct {
	Pos            utils.Vec2
	Start_Move_Pos utils.Vec2
	Move_Pos       utils.Vec2
	Clicking       bool
}

func NewCamera(pos utils.Vec2) (camera CameraStruct) {
	camera.Pos = pos

	return camera
}

func (camera *CameraStruct) Update() {
	if ebiten.IsMouseButtonPressed(ebiten.MouseButtonRight) && !camera.Clicking {
		camera.Clicking = true
		camera.Start_Move_Pos = utils.Vec2{X: utils.Mouse_X, Y: utils.Mouse_Y}
	}

	if camera.Clicking {
		camera.Move_Pos = utils.Vec2{X: utils.Mouse_X, Y: utils.Mouse_Y}
	}

	if !ebiten.IsMouseButtonPressed(ebiten.MouseButtonRight) {
		camera.Clicking = false
		camera.Pos = utils.Vec2{X: camera.Pos.X + (camera.Start_Move_Pos.X - camera.Move_Pos.X), Y: camera.Pos.Y + (camera.Start_Move_Pos.Y - camera.Move_Pos.Y)}
		camera.Move_Pos = utils.Vec2{X: 0, Y: 0}
		camera.Start_Move_Pos = utils.Vec2{X: 0, Y: 0}
	}
}

var Camera = NewCamera(utils.Vec2{X: 0, Y: 0})
