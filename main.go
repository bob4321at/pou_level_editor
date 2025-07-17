package main

import (
	"main/camera"
	"main/grid"
	"main/ui"
	"main/utils"
	"strconv"

	"github.com/ebitengine/debugui"
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
)

type Game struct {
	debugui debugui.DebugUI
}

func (g *Game) Update() error {
	var ics debugui.InputCapturingState
	var err error

	if ics, err = g.debugui.Update(func(ctx *debugui.Context) error {
		ui.Menu(ctx)
		if grid.SelectedEnemySpawner != nil {
			ui.EditEnemySpawnerUi(ctx)
		}
		if grid.SelectedBreakableTile != nil {
			ui.EditBreakableTileUi(ctx)
		}
		if grid.SelectedTriggerTile != nil {
			ui.EditTriggerTileUi(ctx)
		}
		if grid.SelectedGunTile != nil {
			ui.EditGunTileUi(ctx)
		}
		return nil
	}); err != nil {
		panic(err)
	}

	rmx, rmy := ebiten.CursorPosition()

	utils.Mouse_X = float64(rmx)
	utils.Mouse_Y = float64(rmy)

	camera.Camera.Update()

	if ics == 0 {
		grid.Current_Level.Update()
	}

	return nil
}

func (g *Game) Draw(screen *ebiten.Image) {
	grid.Current_Level.Draw(screen)

	g.debugui.Draw(screen)

	ebitenutil.DebugPrint(screen, strconv.Itoa(int(ebiten.ActualFPS())))
}

func (g *Game) Layout(ow, oh int) (sw, sh int) {
	return 640, 360
}

func main() {
	ebiten.SetWindowSize(1920, 1080)

	if err := ebiten.RunGame(&Game{}); err != nil {
		panic(err)
	}
}
